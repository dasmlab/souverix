#!/bin/bash
# Generate self-signed certificates for docker:dind

set -e

CERT_DIR=$(mktemp -d)
trap "rm -rf $CERT_DIR" EXIT

cd "$CERT_DIR"

# Create CA key and cert
openssl genrsa -out ca-key.pem 4096
openssl req -new -x509 -days 365 -key ca-key.pem -sha256 -out ca.pem -subj "/CN=docker-dind-ca"

# Create server key
openssl genrsa -out server-key.pem 4096
openssl req -subj "/CN=docker-dind-server" -sha256 -new -key server-key.pem -out server.csr

# Create server cert
echo subjectAltName = IP:127.0.0.1,IP:0.0.0.0,DNS:localhost > extfile.cnf
openssl x509 -req -days 365 -sha256 -in server.csr -CA ca.pem -CAkey ca-key.pem \
  -CAcreateserial -out server-cert.pem -extfile extfile.cnf

# Create client key
openssl genrsa -out key.pem 4096
openssl req -subj '/CN=client' -new -key key.pem -out client.csr

# Create client cert
echo extendedKeyUsage = clientAuth > extfile-client.cnf
openssl x509 -req -days 365 -sha256 -in client.csr -CA ca.pem -CAkey ca-key.pem \
  -CAcreateserial -out cert.pem -extfile extfile-client.cnf

# Create the secret
echo "Creating Kubernetes secret with dind certificates..."
oc create secret generic dind-certs \
  --from-file=ca.pem \
  --from-file=server-cert.pem \
  --from-file=server-key.pem \
  --from-file=cert.pem \
  --from-file=key.pem \
  --namespace=github-runner \
  --dry-run=client -o yaml | oc apply -f -

echo "✅ Dind certificates created and stored in secret 'dind-certs' in namespace 'github-runner'"
