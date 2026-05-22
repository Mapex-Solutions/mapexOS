#!/bin/bash
# =============================================================
# Seed Orchestrator — JS-Executor Benchmark
#
# Manages the complete test data lifecycle:
#   setup    — Seed MongoDB + MinIO + NATS stream
#   teardown — Purge streams + cleanup assets + flush Redis
#
# Uses modular functions from config + common/services/:
#   config.sh                    — Constants, common bootstrap
#   common/services/assets.sh    — MongoDB + MinIO seed/cleanup
#   common/services/nats.sh      — Stream purge/seed/verify
#
# Usage:
#   ./seed.sh setup [message_count]   # Default: 1,000,000
#   ./seed.sh teardown                # Clean all benchmark data
#
# Prerequisites:
#   - MongoDB running at localhost:27017 (replica set: rs0)
#   - MinIO accessible via mc alias "local"
#   - NATS running at localhost:4222 (user: service / service_secret)
#   - Redis running at localhost:6379
#   - CLIs: nats, mongosh, mc, redis-cli
# =============================================================

set -euo pipefail

COMMAND="${1:-}"
MESSAGE_COUNT="${2:-1000000}"

# ─── Source Modules ────────────────────────────────────────

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"

source "$SCRIPT_DIR/config.sh"
source "$COMMON_DIR/services/assets.sh"
source "$COMMON_DIR/services/nats.sh"

# ─── Commands ──────────────────────────────────────────────

cmd_setup() {
    echo ""
    echo "================================================================"
    echo "  JS-Executor Seed — Setup"
    echo "================================================================"
    echo "  Messages:      $(format_number "$MESSAGE_COUNT")"
    echo "  Stream:        $STREAM"
    echo "  Subject:       $SUBJECT"
    echo "  Asset UUID:    $BENCH_ASSET_UUID"
    echo "  Template ID:   $BENCH_TMPL_ID"
    echo "================================================================"
    echo ""

    # 1. Seed MongoDB (AssetTemplate + Asset)
    log_step "1" "3" "Seeding MongoDB..."
    bench_seed_assets_mongodb

    # 2. Seed MinIO (TieredCache L2)
    echo ""
    log_step "2" "3" "Seeding MinIO..."
    bench_seed_assets_minio

    # 3. Populate NATS stream
    echo ""
    log_step "3" "3" "Populating NATS stream..."
    bench_remove_consumer
    bench_seed_messages "$MESSAGE_COUNT" "$SEED_DIR/nats/js-execute-http.json"
    bench_verify_stream "$MESSAGE_COUNT"

    echo ""
    log_success "Setup complete. Run the benchmark:"
    echo "    bash docs/benchmarks/scripts/full-benchmark.sh $MESSAGE_COUNT"
    echo ""
}

cmd_teardown() {
    echo ""
    echo "================================================================"
    echo "  JS-Executor Seed — Teardown"
    echo "================================================================"
    echo ""

    # 1. Purge NATS streams
    log_step "1" "2" "Purging NATS streams..."
    bench_purge_all_streams

    # 2. Cleanup assets (MongoDB + MinIO + Redis)
    echo ""
    log_step "2" "2" "Cleaning assets..."
    bench_cleanup_assets

    echo ""
    log_success "Teardown complete. All benchmark data removed."
    echo ""
}

# ─── Main ──────────────────────────────────────────────────

case "$COMMAND" in
    setup)
        cmd_setup
        ;;
    teardown)
        cmd_teardown
        ;;
    *)
        echo "Usage: $0 <setup|teardown> [message_count]"
        echo ""
        echo "  setup [N]   Seed MongoDB + MinIO + publish N messages (default: 1,000,000)"
        echo "  teardown    Purge streams + cleanup assets + flush Redis"
        echo ""
        echo "  Examples:"
        echo "    $0 setup              # 1M messages"
        echo "    $0 setup 10000        # 10K messages (quick test)"
        echo "    $0 teardown           # clean everything"
        exit 1
        ;;
esac
