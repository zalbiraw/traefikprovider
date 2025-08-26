#!/usr/bin/env bash
# Integration shell test for preview environment
# Mirrors checks in integration_test.go using curl + jq
set -euo pipefail

COMPOSE_FILE="$(dirname "$0")/docker-compose.yml"
TIMEOUT_SECONDS=${TIMEOUT_SECONDS:-60}
SLEEP_SECONDS=2

need() {
  command -v "$1" >/dev/null 2>&1 || { echo "Error: required command '$1' not found in PATH" >&2; exit 97; }
}

wait_for() {
  local url="$1"
  local deadline=$(( $(date +%s) + TIMEOUT_SECONDS ))
  while :; do
    if code=$(curl -s -o /dev/null -w "%{http_code}" "$url" || true); then
      if [[ "$code" == "200" ]]; then
        echo "Ready: $url"
        return 0
      fi
    fi
    if (( $(date +%s) > deadline )); then
      echo "Timeout waiting for $url" >&2
      return 1
    fi
    sleep "$SLEEP_SECONDS"
  done
}

assert_key_exists() {
  local json="$1"; shift
  local jq_path="$1"; shift
  local key="$1"; shift
  if ! echo "$json" | jq -e "$jq_path | has(\"$key\")" >/dev/null; then
    echo "Assertion failed: key '$key' not found at $jq_path" >&2
    return 1
  fi
}

main() {
  need docker-compose
  need curl
  need jq

  echo "Starting services with docker-compose..."
  docker-compose -f "$COMPOSE_FILE" up -d

  echo "Waiting for provider1 and provider2 APIs to be ready..."
  wait_for "http://localhost:8081/api/rawdata"
  wait_for "http://localhost:8082/api/rawdata"

  echo "Fetching rawdata from provider1..."
  p1_json=$(curl -s "http://localhost:8081/api/rawdata")

  echo "Validating expected routers for provider1..."
  for r in provider1-api provider1-web provider1-admin provider1-test; do
    assert_key_exists "$p1_json" '.http.routers' "$r"
  done

  echo "Validating expected services for provider1..."
  for s in provider1-service provider1-web-service provider1-admin-service provider1-test-service; do
    assert_key_exists "$p1_json" '.http.services' "$s"
  done

  echo "Fetching rawdata from provider2..."
  p2_json=$(curl -s "http://localhost:8082/api/rawdata")

  echo "Validating expected routers for provider2..."
  for r in provider2-dashboard provider2-api provider2-secure provider2-metrics; do
    assert_key_exists "$p2_json" '.http.routers' "$r"
  done

  echo "Validating expected services for provider2..."
  for s in provider2-service provider2-api-service provider2-secure-service provider2-metrics-service; do
    assert_key_exists "$p2_json" '.http.services' "$s"
  done

  echo "All shell integration checks passed."
}

main "$@"
