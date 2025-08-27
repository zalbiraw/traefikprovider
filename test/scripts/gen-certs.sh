#!/usr/bin/env bash
set -euo pipefail

# Generate example CA, server, and client certificates for mTLS testing.
# Outputs to test/certs/ with filenames expected by test configs.

ROOT_DIR="$(cd "$(dirname "$0")"/.. && pwd)"
CERT_DIR="$ROOT_DIR/certs"
mkdir -p "$CERT_DIR"

# OpenSSL config for SAN support
make_openssl_cnf() {
  local cn=$1
  local san_csv=$2
  cat <<EOF
[ req ]
prompt = no
distinguished_name = dn
req_extensions = req_ext

[ dn ]
CN = ${cn}
O = Example
C = US

[ req_ext ]
subjectAltName = @alt_names

[ alt_names ]
$(awk -F, '{ for (i=1; i<=NF; i++) printf("DNS.%d = %s\n", i, $i) }' <<<"$san_csv")
EOF
}

# Generate CA
if [[ ! -f "$CERT_DIR/ca.key" ]]; then
  openssl genrsa -out "$CERT_DIR/ca.key" 4096
fi
if [[ ! -f "$CERT_DIR/ca.crt" ]]; then
  openssl req -x509 -new -nodes -key "$CERT_DIR/ca.key" -sha256 -days 3650 \
    -subj "/CN=Test Root CA/O=Example/C=US" \
    -out "$CERT_DIR/ca.crt"
fi

# Function to create server cert signed by CA
make_server_cert() {
  local name=$1
  local cn=$2
  local sans=$3

  local cnf
  cnf=$(mktemp)
  make_openssl_cnf "$cn" "$sans" > "$cnf"

  openssl genrsa -out "$CERT_DIR/${name}.key" 2048
  openssl req -new -key "$CERT_DIR/${name}.key" -out "$CERT_DIR/${name}.csr" -config "$cnf"
  openssl x509 -req -in "$CERT_DIR/${name}.csr" -CA "$CERT_DIR/ca.crt" -CAkey "$CERT_DIR/ca.key" -CAcreateserial \
    -out "$CERT_DIR/${name}.crt" -days 825 -sha256 -extensions req_ext -extfile "$cnf"
  rm -f "$cnf" "$CERT_DIR/${name}.csr"
}

# Function to create client cert signed by CA
make_client_cert() {
  local name=$1
  openssl genrsa -out "$CERT_DIR/${name}.key" 2048
  openssl req -new -key "$CERT_DIR/${name}.key" -subj "/CN=${name}/O=Example/C=US" -out "$CERT_DIR/${name}.csr"
  openssl x509 -req -in "$CERT_DIR/${name}.csr" -CA "$CERT_DIR/ca.crt" -CAkey "$CERT_DIR/ca.key" -CAcreateserial \
    -out "$CERT_DIR/${name}.crt" -days 825 -sha256
  rm -f "$CERT_DIR/${name}.csr"
}

# Servers for provider1 and provider2 (match docker container hostnames)
make_server_cert server1 "traefik-provider1" "traefik-provider1,localhost"
make_server_cert server2 "traefik-provider2" "traefik-provider2,localhost"

# Client for mTLS tunnels
make_client_cert client

# Print summary
ls -l "$CERT_DIR"
echo "\nCertificates generated under $CERT_DIR"
echo "Server1: server1.crt/server1.key"
echo "Server2: server2.crt/server2.key"
echo "Client:  client.crt/client.key"
echo "CA:      ca.crt"
