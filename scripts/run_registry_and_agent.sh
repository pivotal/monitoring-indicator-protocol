#!/usr/bin/env bash
set -efu

export CERTS=./test_fixtures
export REGISTRY_PORT=8091
export REGISTRY_HOST="https://localhost:${REGISTRY_PORT}"
export INDICATOR_DOCUMENTS='./example_indicators*.yml'

echo "Starting registry on PORT $REGISTRY_PORT"
go run cmd/registry/main.go \
  -tls-key-path ${CERTS}/server.key \
  -tls-pem-path ${CERTS}/server.pem \
  -tls-root-ca-pem ${CERTS}/ca.pem \
  -port ${REGISTRY_PORT} \
  -indicator-expiration 1m \
  -config example_config.yml &
sleep 3
echo "Starting registry agent for $REGISTRY_HOST"
go run cmd/registry_agent/main.go \
  -tls-key-path ${CERTS}/client.key \
  -interval 5s \
  -tls-pem-path ${CERTS}/client.pem \
  -tls-root-ca-pem ${CERTS}/ca.pem \
  -tls-server-cn localhost \
  -registry ${REGISTRY_HOST} \
  -documents-glob ${INDICATOR_DOCUMENTS} &

wait
