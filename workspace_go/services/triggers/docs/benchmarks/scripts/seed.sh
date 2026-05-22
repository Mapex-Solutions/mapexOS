#!/bin/bash
# =============================================================
# Seed Orchestrator — Triggers Benchmark
#
# Manages the complete test data lifecycle:
#   setup    — Seed 5 trigger documents into MongoDB (one per executor)
#   teardown — Purge NATS streams + delete seed data + flush Redis
#
# Uses modular functions from config + common/services/:
#   config.sh                       — Constants, common bootstrap
#   common/services/triggers.sh     — MongoDB seed/cleanup
#   common/services/nats.sh         — Stream purge
#
# Usage:
#   ./seed.sh setup      — Insert benchmark data (idempotent)
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
source "$COMMON_DIR/services/triggers.sh"
source "$COMMON_DIR/services/nats.sh"

# ─── Commands ────────────────────────────────────────────────

cmd_setup() {
    echo ""
    echo "================================================================"
    echo "  Triggers Benchmark Seed — Setup"
    echo "================================================================"
    echo "  MongoDB:     ${MONGO_DB}.${MONGO_COLLECTION}"
    echo "  Scenarios:   http, mqtt, nats, rabbitmq, email"
    echo "  Org ID:      ${BENCH_ORG_ID}"
    echo "================================================================"
    echo ""

    # 1. Seed MongoDB (5 Trigger documents)
    log_step "1" "1" "Seeding MongoDB Trigger documents..."
    bench_seed_triggers

    echo ""
    log_success "Setup complete."
    echo ""
}

cmd_teardown() {
    echo ""
    echo "================================================================"
    echo "  Triggers Benchmark Seed — Teardown"
    echo "================================================================"
    echo ""

    # 1. Purge NATS streams (stop in-flight processing first)
    log_step "1" "2" "Purging NATS streams..."
    bench_purge_all_streams

    # 2. Cleanup Triggers (MongoDB + Redis)
    echo ""
    log_step "2" "2" "Cleaning Trigger data..."
    bench_cleanup_triggers

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
        echo "  setup      Seed 5 trigger documents into MongoDB (idempotent)"
        echo "  teardown   Purge NATS + delete seed data + flush Redis"
        echo ""
        echo "  Examples:"
        echo "    $0 setup       # create all benchmark data"
        echo "    $0 teardown    # clean everything"
        echo ""
        echo "  Environment variables (all optional, have defaults):"
        echo "    GO_ENV         Default: dev"
        echo "    MONGO_DB       Default: \${GO_ENV}-triggers"
        exit 1
        ;;
esac
