#!/usr/bin/env bash
set -euo pipefail

echo "=== Complai Dev Bootstrap ==="

check_cmd() {
  if ! command -v "$1" &>/dev/null; then
    echo "ERROR: $1 is not installed. Please install it first."
    return 1
  fi
  echo "  ✓ $1 found: $(command -v "$1")"
}

echo ""
echo "Checking prerequisites..."
check_cmd go
check_cmd node
check_cmd pnpm
check_cmd docker

GO_VERSION=$(go version | grep -oP '\d+\.\d+' | head -1)
NODE_VERSION=$(node --version | grep -oP '\d+' | head -1)

echo ""
echo "Versions:"
echo "  Go:   $(go version)"
echo "  Node: $(node --version)"
echo "  pnpm: $(pnpm --version)"

echo ""
echo "Installing Node dependencies..."
pnpm install

echo ""
echo "Syncing Go workspace..."
go work sync

echo ""
echo "Starting dev services..."
make dev

echo ""
echo "Waiting for services to be healthy (up to 120s)..."
TIMEOUT=120
ELAPSED=0
while [ $ELAPSED -lt $TIMEOUT ]; do
  HEALTHY=$(docker compose -f docker-compose.dev.yml ps --format json 2>/dev/null | grep -c '"healthy"' || echo "0")
  TOTAL=$(docker compose -f docker-compose.dev.yml ps --format json 2>/dev/null | wc -l | tr -d ' ')
  echo "  Healthy: $HEALTHY / $TOTAL (${ELAPSED}s elapsed)"
  if [ "$HEALTHY" -ge 8 ]; then
    echo "  All critical services healthy!"
    break
  fi
  sleep 5
  ELAPSED=$((ELAPSED + 5))
done

echo ""
echo "=== Bootstrap complete ==="
echo ""
echo "Services available at:"
echo "  Postgres:     localhost:5432  (user: complai, pass: complai_dev)"
echo "  Redis:        localhost:6379"
echo "  LocalStack:   localhost:4566"
echo "  OpenSearch:   localhost:9200"
echo "  Keycloak:     http://localhost:8080  (admin/admin)"
echo "  Temporal:     localhost:7233"
echo "  Temporal UI:  http://localhost:8088"
echo "  Mailpit:      http://localhost:8025"
echo "  Jaeger:       http://localhost:16686"
echo "  OTel:         localhost:4317 (gRPC), localhost:4318 (HTTP)"
