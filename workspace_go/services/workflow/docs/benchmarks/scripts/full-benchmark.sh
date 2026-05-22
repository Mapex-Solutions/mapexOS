#!/bin/bash
# Full benchmark suite for workflow service
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
source "$SCRIPT_DIR/config.sh"

RESULTS_TAG_DIR="$RESULTS_DIR/$TAG"
mkdir -p "$RESULTS_TAG_DIR"

echo "=== Workflow Service Benchmark ==="
echo "Tag: $TAG"
echo "Results: $RESULTS_TAG_DIR"
echo ""

# Health check
echo "[0/3] Health check..."
curl -sf "$SERVICE_URL/health" | jq .status
echo ""

# Test 1: Definition CRUD throughput
echo "[1/3] Definition CRUD throughput..."
hey -n "$REQUESTS" -c "$CONCURRENCY" -m GET \
  -H "Authorization: Bearer $AUTH_TOKEN" \
  "$SERVICE_URL/api/v1/workflows" \
  > "$RESULTS_TAG_DIR/test-definitions-list.txt" 2>&1

# Test 2: Instance listing throughput
echo "[2/3] Instance listing throughput..."
hey -n "$REQUESTS" -c "$CONCURRENCY" -m GET \
  -H "Authorization: Bearer $AUTH_TOKEN" \
  "$SERVICE_URL/api/v1/workflow-instances" \
  > "$RESULTS_TAG_DIR/test-instances-list.txt" 2>&1

# Test 3: Metrics endpoint
echo "[3/3] Collecting metrics..."
curl -sf "$SERVICE_URL/metrics" > "$RESULTS_TAG_DIR/metrics.txt"

echo ""
echo "=== Benchmark complete ==="
echo "Results saved to: $RESULTS_TAG_DIR"
