#!/bin/bash
set -e

# GitHub Actions Runner startup script for Kubernetes
# Container runs as root, but we switch to runner user for the GitHub Actions runner
# Setup for Kaniko

# Ensure PATH includes buildah/podman location
export PATH="/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:${PATH}"

# Install sudo if not present
if ! command -v sudo &> /dev/null; then
    apt-get update && apt-get install -y sudo && rm -rf /var/lib/apt/lists/* || true
fi

# Ensure runner user exists (should already exist from Dockerfile)
if ! id runner &>/dev/null; then
    useradd -m -u 1000 runner
    # Configure sudo if user was just created (dev/internal cluster - allow all)
    echo "runner ALL=(ALL) NOPASSWD: ALL" > /etc/sudoers.d/runner && \
    chmod 0440 /etc/sudoers.d/runner || true
fi

cd /home/runner

# Switch to runner user and start the runner
# GitHub Actions runner must run as non-root
# Use 'su runner' (not 'su - runner') to preserve environment and sudo access
if [ -f ".runner" ]; then
    echo "Runner already configured, starting as runner user..."
    # Ensure SSL certificates are available
    export SSL_CERT_DIR=${SSL_CERT_DIR:-/etc/ssl/certs}
    export REQUESTS_CA_BUNDLE=${REQUESTS_CA_BUNDLE:-/etc/ssl/certs/ca-certificates.crt}
    export CURL_CA_BUNDLE=${CURL_CA_BUNDLE:-/etc/ssl/certs/ca-certificates.crt}
    su runner -c "cd /home/runner && SSL_CERT_DIR=${SSL_CERT_DIR} REQUESTS_CA_BUNDLE=${REQUESTS_CA_BUNDLE} CURL_CA_BUNDLE=${CURL_CA_BUNDLE} ./run.sh"
else
    echo "Configuring runner as runner user..."

    # Validate required environment variables
    if [ -z "$GITHUB_TOKEN" ]; then
        echo "ERROR: GITHUB_TOKEN environment variable is required"
        exit 1
    fi

    # Support both REPO_URL and GITHUB_REPO_URL for compatibility
    REPO_URL=${REPO_URL:-$GITHUB_REPO_URL}
    if [ -z "$REPO_URL" ]; then
        echo "ERROR: REPO_URL or GITHUB_REPO_URL environment variable is required (e.g., https://github.com/org/repo)"
        exit 1
    fi

    # Set default runner name if not provided
    RUNNER_NAME=${RUNNER_NAME:-"k8s-runner-$(hostname)"}
    
    # Set default labels if not provided
    RUNNER_LABELS=${RUNNER_LABELS:-"self-hosted,linux,x64"}

    # Ensure SSL certificates are available for runner
    # The runner needs to connect to GitHub API over HTTPS
    export SSL_CERT_DIR=${SSL_CERT_DIR:-/etc/ssl/certs}
    export REQUESTS_CA_BUNDLE=${REQUESTS_CA_BUNDLE:-/etc/ssl/certs/ca-certificates.crt}
    export CURL_CA_BUNDLE=${CURL_CA_BUNDLE:-/etc/ssl/certs/ca-certificates.crt}

    # Configure the runner as runner user with labels
    su runner -c "cd /home/runner && SSL_CERT_DIR=${SSL_CERT_DIR} REQUESTS_CA_BUNDLE=${REQUESTS_CA_BUNDLE} CURL_CA_BUNDLE=${CURL_CA_BUNDLE} ./config.sh --url \"$REPO_URL\" --pat \"$GITHUB_TOKEN\" --name \"$RUNNER_NAME\" --labels \"$RUNNER_LABELS\" --work \"_work\" --replace --unattended"

    echo "Runner configured successfully. Starting as runner user..."

    # Start the runner as runner user
    su runner -c "cd /home/runner && ./run.sh"
fi


