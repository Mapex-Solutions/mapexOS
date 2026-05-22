#!/bin/bash
# =============================================================
# Seed Orchestrator — Router Benchmark
#
# Manages the complete test data lifecycle:
#   setup    — Seed ALL 3 RouteGroups into MongoDB + asset into MinIO
#   teardown — Purge NATS streams + delete seed data + flush Redis
#
# Uses modular functions from config + common/services/:
#   config.sh                        — Constants, common bootstrap
#   common/services/routegroups.sh   — MongoDB seed/cleanup
#   common/services/minio.sh         — MinIO seed/cleanup
#   common/services/nats.sh          — Stream purge
#
# Usage:
#   ./seed.sh setup      — Insert benchmark data (idempotent)
#   ./seed.sh teardown   — Clean all benchmark data
#
# Prerequisites:
#   - MongoDB running at localhost:27017 (replica set: rs0)
#   - MinIO accessible via mc alias "local"
#   - NATS running at localhost:4222 (user: service / service_secret)
#   - Redis running at localhost:6379
#   - CLIs: mongosh, mc, nats, redis-cli
# =============================================================

set -euo pipefail

COMMAND="${1:-}"

# ─── Source Modules ──────────────────────────────────────────

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"

source "$SCRIPT_DIR/config.sh"
source "$COMMON_DIR/services/routegroups.sh"
source "$COMMON_DIR/services/minio.sh"
source "$COMMON_DIR/services/nats.sh"

# ─── Commands ────────────────────────────────────────────────

cmd_setup() {
    echo ""
    echo "================================================================"
    echo "  Router Benchmark Seed — Setup"
    echo "================================================================"
    echo "  MongoDB:     ${MONGO_DB}.${MONGO_COLLECTION}"
    echo "  Scenarios:   save_event, rule_engine, trigger"
    echo "  Org ID:      ${BENCH_ORG_ID}"
    echo "  Assets:      1 per scenario (save_event, rule_engine, trigger)"
    echo "================================================================"
    echo ""

    # 1. Seed MongoDB (3 RouteGroups)
    log_step "1" "2" "Seeding MongoDB RouteGroups..."
    bench_seed_routegroups

    # 2. Seed MinIO (1 asset per scenario, each with 1 RouteGroup)
    echo ""
    log_step "2" "2" "Seeding MinIO assets (1 per scenario)..."
    bench_seed_minio "$SEED_DIR/minio/asset-save-event.json" "$BENCH_ASSET_SAVE_EVENT_UUID"
    bench_seed_minio "$SEED_DIR/minio/asset-rule-engine.json" "$BENCH_ASSET_RULE_ENGINE_UUID"
    bench_seed_minio "$SEED_DIR/minio/asset-trigger.json" "$BENCH_ASSET_TRIGGER_UUID"

    echo ""
    log_success "Setup complete."
    echo ""
}

cmd_teardown() {
    echo ""
    echo "================================================================"
    echo "  Router Benchmark Seed — Teardown"
    echo "================================================================"
    echo ""

    # 1. Purge NATS streams (stop in-flight processing first)
    log_step "1" "3" "Purging NATS streams..."
    bench_purge_all_streams

    # 2. Cleanup RouteGroups (MongoDB + Redis)
    echo ""
    log_step "2" "3" "Cleaning RouteGroups..."
    bench_cleanup_routegroups

    # 3. Cleanup MinIO assets (all 3 scenarios)
    echo ""
    log_step "3" "3" "Cleaning MinIO assets..."
    bench_cleanup_minio "$BENCH_ASSET_SAVE_EVENT_UUID"
    bench_cleanup_minio "$BENCH_ASSET_RULE_ENGINE_UUID"
    bench_cleanup_minio "$BENCH_ASSET_TRIGGER_UUID"

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
        echo "  setup      Seed 3 RouteGroups + MinIO asset (idempotent)"
        echo "  teardown   Purge NATS + delete seed data + flush Redis + remove asset"
        echo ""
        echo "  Examples:"
        echo "    $0 setup       # create all benchmark data"
        echo "    $0 teardown    # clean everything"
        echo ""
        echo "  Environment variables (all optional, have defaults):"
        echo "    GO_ENV         Default: dev"
        echo "    MONGO_DB       Default: \${GO_ENV}-router"
        echo "    MINIO_ALIAS    Default: local"
        echo "    MINIO_BUCKET   Default: mapex-assets"
        exit 1
        ;;
esac
