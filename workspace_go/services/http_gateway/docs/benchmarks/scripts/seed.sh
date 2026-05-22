#!/bin/bash
# =============================================================
# Seed Orchestrator — HTTP Gateway Benchmark
#
# Manages the complete test data lifecycle:
#   setup    — Seed 4 DataSources into MongoDB (all auth types)
#   teardown — Purge NATS streams + delete seed data + flush Redis
#
# Uses modular functions from config + common/services/:
#   config.sh                        — Constants, common bootstrap
#   common/services/datasources.sh   — MongoDB seed/cleanup
#   common/services/nats.sh          — Stream purge
#
# Usage:
#   ./seed.sh setup      — Insert benchmark DataSources (idempotent)
#   ./seed.sh teardown   — Clean all benchmark data
#
# Prerequisites:
#   - MongoDB running at localhost:27017 (replica set: rs0)
#   - NATS running at localhost:4222 (user: service / service_secret)
#   - Redis running at localhost:6379
#   - CLIs: mongosh, nats, redis-cli
# =============================================================

set -euo pipefail

COMMAND="${1:-}"

# ─── Source Modules ──────────────────────────────────────────

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"

source "$SCRIPT_DIR/config.sh"
source "$COMMON_DIR/services/datasources.sh"
source "$COMMON_DIR/services/nats.sh"

# ─── Commands ────────────────────────────────────────────────

cmd_setup() {
    echo ""
    echo "================================================================"
    echo "  HTTP Gateway Seed — Setup"
    echo "================================================================"
    echo "  MongoDB:     ${MONGO_DB}.${MONGO_COLLECTION}"
    echo "  Auth types:  jwt, apiKey, ip_whitelist, none"
    echo "  Org ID:      ${BENCH_ORG_ID}"
    echo "================================================================"
    echo ""

    # Seed MongoDB (4 DataSources)
    log_step "1" "1" "Seeding MongoDB DataSources..."
    bench_seed_datasources

    echo ""
    log_success "Setup complete."
    echo ""
}

cmd_teardown() {
    echo ""
    echo "================================================================"
    echo "  HTTP Gateway Seed — Teardown"
    echo "================================================================"
    echo ""

    # 1. Purge NATS streams (stop in-flight processing first)
    log_step "1" "2" "Purging NATS streams..."
    bench_purge_all_streams

    # 2. Cleanup DataSources (MongoDB + Redis)
    echo ""
    log_step "2" "2" "Cleaning DataSources..."
    bench_cleanup_datasources

    echo ""
    log_success "Teardown complete. All benchmark data removed."
    echo ""
}

# ─── Main ────────────────────────────────────────────────────

case "$COMMAND" in
    setup)
        cmd_setup
        ;;
    teardown)
        cmd_teardown
        ;;
    *)
        echo "Usage: $0 <setup|teardown>"
        echo ""
        echo "  setup      Seed 4 DataSources into MongoDB (idempotent)"
        echo "  teardown   Purge NATS + delete seed data + flush Redis"
        echo ""
        echo "  Examples:"
        echo "    $0 setup       # create benchmark DataSources"
        echo "    $0 teardown    # clean everything"
        echo ""
        echo "  Environment variables (all optional, have defaults):"
        echo "    GO_ENV         Default: dev"
        echo "    MONGO_URI      Default: mongodb://localhost:27017/?replicaSet=rs0"
        echo "    MONGO_DB       Default: \${GO_ENV}-http_gateway"
        echo "    REDIS_HOST     Default: 127.0.0.1"
        echo "    REDIS_PORT     Default: 6379"
        exit 1
        ;;
esac
