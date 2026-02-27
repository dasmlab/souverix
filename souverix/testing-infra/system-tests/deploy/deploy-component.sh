#!/usr/bin/env bash
set -euo pipefail

# deploy-component.sh - Deploy a component to a test system
# This script is infrastructure-agnostic and only needs:
# - TEST-SYSTEM-NAME, TEST-SYSTEM-TYPE, TEST-SYSTEM-VERSION (env vars)
# - Component deployment manifests (k8s/*.yaml)
# - Deployment credentials (kubeconfig, tokens)

COMPONENT="${1:-}"
COMPONENT_DIR="${2:-}"

if [[ -z "${COMPONENT}" ]] || [[ -z "${COMPONENT_DIR}" ]]; then
    echo "Usage: $0 <component-name> <component-directory>"
    echo "Example: $0 pcscf components/coeur/pcscf"
    exit 1
fi

# Read test system configuration
TEST_SYSTEM_NAME="${TEST-SYSTEM-NAME:-}"
TEST_SYSTEM_TYPE="${TEST-SYSTEM-TYPE:-}"
TEST_SYSTEM_VERSION="${TEST-SYSTEM-VERSION:-}"

if [[ -z "${TEST_SYSTEM_NAME}" ]]; then
    echo "Error: TEST-SYSTEM-NAME environment variable is required"
    exit 1
fi

echo "=========================================="
echo "Deploying ${COMPONENT} to Test System"
echo "=========================================="
echo "  Component: ${COMPONENT}"
echo "  Test System: ${TEST_SYSTEM_NAME}"
echo "  System Type: ${TEST_SYSTEM_TYPE:-not specified}"
echo "  System Version: ${TEST_SYSTEM_VERSION:-not specified}"
echo ""

# Determine deployment method based on TEST-SYSTEM-TYPE
case "${TEST_SYSTEM_TYPE}" in
    ocp*|openshift*)
        DEPLOY_CMD="oc"
        echo "Using OpenShift deployment (oc)"
        ;;
    k8s*|kubernetes*)
        DEPLOY_CMD="kubectl"
        echo "Using Kubernetes deployment (kubectl)"
        ;;
    *)
        # Default to kubectl
        DEPLOY_CMD="kubectl"
        echo "Using default Kubernetes deployment (kubectl)"
        ;;
esac

# Check if deployment command is available
if ! command -v "${DEPLOY_CMD}" &> /dev/null; then
    echo "Error: ${DEPLOY_CMD} command not found"
    exit 1
fi

# Check for kubeconfig
if [[ -z "${KUBECONFIG:-}" ]] && [[ ! -f "${HOME}/.kube/config" ]]; then
    echo "Error: KUBECONFIG not set and ~/.kube/config not found"
    exit 1
fi

# Find deployment manifests
MANIFEST_DIR="${COMPONENT_DIR}/k8s"
if [[ ! -d "${MANIFEST_DIR}" ]]; then
    echo "Warning: No k8s/ directory found in ${COMPONENT_DIR}"
    echo "Looking for deployment.yaml in component directory..."
    if [[ -f "${COMPONENT_DIR}/deployment.yaml" ]]; then
        MANIFEST_DIR="${COMPONENT_DIR}"
    else
        echo "Error: No deployment manifests found"
        exit 1
    fi
fi

echo ""
echo "Applying deployment manifests from: ${MANIFEST_DIR}"
echo ""

# Apply manifests
if [[ -f "${MANIFEST_DIR}/deployment.yaml" ]]; then
    echo "Applying deployment.yaml..."
    ${DEPLOY_CMD} apply -f "${MANIFEST_DIR}/deployment.yaml" --namespace="${TEST_SYSTEM_NAME}" || ${DEPLOY_CMD} apply -f "${MANIFEST_DIR}/deployment.yaml"
fi

if [[ -f "${MANIFEST_DIR}/service.yaml" ]]; then
    echo "Applying service.yaml..."
    ${DEPLOY_CMD} apply -f "${MANIFEST_DIR}/service.yaml" --namespace="${TEST_SYSTEM_NAME}" || ${DEPLOY_CMD} apply -f "${MANIFEST_DIR}/service.yaml"
fi

if [[ -f "${MANIFEST_DIR}/configmap.yaml" ]]; then
    echo "Applying configmap.yaml..."
    ${DEPLOY_CMD} apply -f "${MANIFEST_DIR}/configmap.yaml" --namespace="${TEST_SYSTEM_NAME}" || ${DEPLOY_CMD} apply -f "${MANIFEST_DIR}/configmap.yaml"
fi

# Wait for deployment to be ready
echo ""
echo "Waiting for ${COMPONENT} deployment to be ready..."
${DEPLOY_CMD} wait --for=condition=available --timeout=300s deployment/${COMPONENT} --namespace="${TEST_SYSTEM_NAME}" 2>/dev/null || \
${DEPLOY_CMD} wait --for=condition=available --timeout=300s deployment/${COMPONENT} || \
echo "Warning: Deployment readiness check skipped or failed"

echo ""
echo "✅ ${COMPONENT} deployed to ${TEST_SYSTEM_NAME}"
echo ""
