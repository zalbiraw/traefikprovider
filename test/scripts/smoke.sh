#!/usr/bin/env bash
set -euo pipefail

MAIN_HOST=localhost
WEB_PORT=80
TCP_PORT=8000
UDP_PORT=9000

section() { echo; echo "==> $1"; }
req() {
  local host=$1
  local path=${2:-/}
  echo "curl -s -H 'Host: ${host}' http://${MAIN_HOST}:${WEB_PORT}${path} | head -n 5"
  curl -s -H "Host: ${host}" "http://${MAIN_HOST}:${WEB_PORT}${path}" | head -n 20 || true
  echo
}

section "HTTP routes via traefik-main (provider1)"
req provider1.example.com /
req provider1.example.com /web
req admin.provider1.example.com /
req test.provider1.example.com /health

section "HTTP routes via traefik-main (provider2)"
req provider2.example.com /
req secure.provider2.example.com /
req metrics.provider2.example.com /metrics || true

section "TCP route via traefik-main"
# Send a line to tcp-echo service through Traefik TCP entrypoint
if command -v nc >/dev/null; then
  echo "echo 'hello tcp' | nc ${MAIN_HOST} ${TCP_PORT}"
  echo 'hello tcp' | nc ${MAIN_HOST} ${TCP_PORT} || true
else
  echo "nc not found; skipping TCP smoke"
fi

section "UDP route via traefik-main"
if command -v nc >/dev/null; then
  echo "echo 'hello udp' | nc -u -w1 ${MAIN_HOST} ${UDP_PORT}"
  echo 'hello udp' | nc -u -w1 ${MAIN_HOST} ${UDP_PORT} || true
else
  echo "nc not found; skipping UDP smoke"
fi

section "Done"
