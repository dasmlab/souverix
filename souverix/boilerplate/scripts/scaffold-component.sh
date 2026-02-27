#!/usr/bin/env bash
set -euo pipefail

# scaffold-component.sh - Scaffold a new Souverix component from boilerplate templates
# Usage: ./scaffold-component.sh <component> <component-name> <internal-path> [ports]

COMPONENT="${1:-}"
COMPONENT_NAME="${2:-}"
INTERNAL_PATH="${3:-}"
PORTS="${4:-8081}"

if [[ -z "$COMPONENT" || -z "$COMPONENT_NAME" || -z "$INTERNAL_PATH" ]]; then
    echo "Usage: $0 <component> <component-name> <internal-path> [ports]"
    echo "Example: $0 rempart 'Souverix Rempart' 'internal/rempart' '5060:5060 8081:8081'"
    exit 1
fi

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BOILERPLATE_DIR="$(cd "${SCRIPT_DIR}/.." && pwd)"
SOUVERIX_ROOT="$(cd "${BOILERPLATE_DIR}/.." && pwd)"

echo "🔨 Scaffolding ${COMPONENT_NAME} (${COMPONENT})..."

# Component package name (convert path to package)
COMPONENT_DIR="${INTERNAL_PATH#internal/}"
COMPONENT_PKG=$(echo "$COMPONENT_DIR" | tr '/' '_')

# Determine component initialization based on path
# For now, create stub components that will compile
if [[ "$INTERNAL_PATH" == "internal/coeur"* ]]; then
    COMPONENT_INIT="// Initialize Coeur component\n\tcomponentCfg := ${COMPONENT_PKG}.ConfigFromGouverne(cfg)\n\tcomponent := ${COMPONENT_PKG}.New(componentCfg, log)"
    COMPONENT_START="component.Start(ctx)"
    COMPONENT_STOP="component.Stop(shutdownCtx)"
elif [[ "$INTERNAL_PATH" == "internal/rempart"* ]]; then
    COMPONENT_INIT="// Initialize Rempart component\n\tcomponent, err := ${COMPONENT_PKG}.NewSBC(cfg, log)\n\tif err != nil {\n\t\tlog.WithError(err).Fatal(\"failed to create Rempart\")\n\t}"
    COMPONENT_START="component.Start(ctx)"
    COMPONENT_STOP="component.Stop(shutdownCtx)"
else
    # Generic initialization - create stub component
    COMPONENT_INIT="// Initialize ${COMPONENT_NAME} (stub)\n\tcomponent := ${COMPONENT_PKG}.New(cfg, log)"
    COMPONENT_START="component.Start(ctx)"
    COMPONENT_STOP="component.Stop(shutdownCtx)"
fi

# Parse ports for Dockerfile and k8s
EXPOSE_LINES=""
K8S_PORTS=""
K8S_SERVICE_PORTS=""
NETWORK_PORTS=""

IFS=' ' read -ra PORT_ARRAY <<< "$PORTS"
for port_pair in "${PORT_ARRAY[@]}"; do
    IFS=':' read -r host_port container_port <<< "$port_pair"
    if [[ -z "$container_port" ]]; then
        container_port="$host_port"
    fi
    
    # Dockerfile EXPOSE
    if [[ "$container_port" == *"/udp" ]]; then
        EXPOSE_LINES="${EXPOSE_LINES}EXPOSE ${container_port}\n"
    else
        EXPOSE_LINES="${EXPOSE_LINES}EXPOSE ${container_port}/tcp\n"
    fi
    
    # K8s container ports
    port_name=$(echo "$container_port" | sed 's|/.*||' | tr '-' '_')
    K8S_PORTS="${K8S_PORTS}        - name: ${port_name}\n          containerPort: ${container_port%%/*}\n          protocol: TCP\n"
    
    # K8s service ports
    K8S_SERVICE_PORTS="${K8S_SERVICE_PORTS}  - name: ${port_name}\n    port: ${host_port}\n    targetPort: ${container_port%%/*}\n    protocol: TCP\n"
    
    # Docker network ports
    if [[ "$container_port" == *"/udp" ]]; then
        NETWORK_PORTS="${NETWORK_PORTS}    -p ${host_port}:${container_port%%/*}/udp \\\\\n"
    else
        NETWORK_PORTS="${NETWORK_PORTS}    -p ${host_port}:${container_port%%/*} \\\\\n"
    fi
done

# Always add health port
EXPOSE_LINES="${EXPOSE_LINES}EXPOSE 8081/tcp\n"
K8S_PORTS="${K8S_PORTS}        - name: health\n          containerPort: 8081\n          protocol: TCP\n"
K8S_SERVICE_PORTS="${K8S_SERVICE_PORTS}  - name: health\n    port: 8081\n    targetPort: 8081\n    protocol: TCP\n"
NETWORK_PORTS="${NETWORK_PORTS}    -p 8081:8081"

# Create component directory structure
COMPONENT_DIR="${INTERNAL_PATH#internal/}"
mkdir -p "${SOUVERIX_ROOT}/internal/${COMPONENT_DIR}/cmd"
mkdir -p "${SOUVERIX_ROOT}/k8s/${COMPONENT}"

# Function to process template
process_template() {
    local template_file="$1"
    local output_file="$2"
    
    sed -e "s|{{COMPONENT}}|${COMPONENT}|g" \
        -e "s|{{COMPONENT_NAME}}|${COMPONENT_NAME}|g" \
        -e "s|{{COMPONENT_PATH}}|${INTERNAL_PATH}|g" \
        -e "s|{{COMPONENT_DIR}}|${COMPONENT_DIR}|g" \
        -e "s|{{COMPONENT_PKG}}|${COMPONENT_PKG}|g" \
        -e "s|{{COMPONENT_INIT}}|${COMPONENT_INIT}|g" \
        -e "s|{{COMPONENT_START}}|${COMPONENT_START}|g" \
        -e "s|{{COMPONENT_STOP}}|${COMPONENT_STOP}|g" \
        -e "s|{{EXPOSE_PORTS}}|${EXPOSE_LINES}|g" \
        -e "s|{{K8S_PORTS}}|${K8S_PORTS}|g" \
        -e "s|{{K8S_SERVICE_PORTS}}|${K8S_SERVICE_PORTS}|g" \
        -e "s|{{NETWORK_PORTS}}|${NETWORK_PORTS}|g" \
        "$template_file" > "$output_file"
}

# Copy and process templates
echo "  📝 Creating Dockerfile..."
process_template "${BOILERPLATE_DIR}/templates/Dockerfile.COMPONENT" "${SOUVERIX_ROOT}/Dockerfile.${COMPONENT}"

echo "  📝 Creating cmd/main.go..."
process_template "${BOILERPLATE_DIR}/templates/cmd/main.go" "${SOUVERIX_ROOT}/internal/${COMPONENT_DIR}/cmd/main.go"

echo "  📝 Creating buildme.sh..."
process_template "${BOILERPLATE_DIR}/templates/buildme.sh" "${SOUVERIX_ROOT}/buildme-${COMPONENT}.sh"
chmod +x "${SOUVERIX_ROOT}/buildme-${COMPONENT}.sh"

echo "  📝 Creating pushme.sh..."
process_template "${BOILERPLATE_DIR}/templates/pushme.sh" "${SOUVERIX_ROOT}/pushme-${COMPONENT}.sh"
chmod +x "${SOUVERIX_ROOT}/pushme-${COMPONENT}.sh"

echo "  📝 Creating pullme.sh..."
process_template "${BOILERPLATE_DIR}/templates/pullme.sh" "${SOUVERIX_ROOT}/pullme-${COMPONENT}.sh"
chmod +x "${SOUVERIX_ROOT}/pullme-${COMPONENT}.sh"

echo "  📝 Creating runme-local.sh..."
process_template "${BOILERPLATE_DIR}/templates/runme-local.sh" "${SOUVERIX_ROOT}/runme-local-${COMPONENT}.sh"
chmod +x "${SOUVERIX_ROOT}/runme-local-${COMPONENT}.sh"

echo "  📝 Creating k8s/deployment.yaml..."
process_template "${BOILERPLATE_DIR}/templates/k8s/deployment.yaml" "${SOUVERIX_ROOT}/k8s/${COMPONENT}/deployment.yaml"

echo "  📝 Creating k8s/configmap.yaml..."
process_template "${BOILERPLATE_DIR}/templates/k8s/configmap.yaml" "${SOUVERIX_ROOT}/k8s/${COMPONENT}/configmap.yaml"

echo "✅ ${COMPONENT_NAME} scaffolded successfully!"
echo ""
echo "Created files:"
echo "  - Dockerfile.${COMPONENT}"
echo "  - internal/${INTERNAL_PATH#internal/}/cmd/main.go"
echo "  - buildme-${COMPONENT}.sh"
echo "  - pushme-${COMPONENT}.sh"
echo "  - pullme-${COMPONENT}.sh"
echo "  - runme-local-${COMPONENT}.sh"
echo "  - k8s/${COMPONENT}/deployment.yaml"
echo "  - k8s/${COMPONENT}/configmap.yaml"
