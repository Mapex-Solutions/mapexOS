#!/bin/bash
# =============================================================
# Seed Orchestrator — JS-Workflow-Executor Benchmark
#
# Manages the complete test data lifecycle:
#   setup    — Seed MinIO + NATS stream
#   teardown — Purge streams + cleanup MinIO
#
# Uses modular functions from config + common/services/:
#   config.sh                    — Constants, common bootstrap, MinIO helpers
#   common/services/nats.sh      — Stream purge/seed/verify
#
# Usage:
#   ./seed.sh setup [message_count]   # Default: 1,000,000
#   ./seed.sh teardown                # Clean all benchmark data
#
# Prerequisites:
#   - MinIO accessible via mc alias "local"
#   - NATS running at localhost:4222 (user: service / service_secret)
#   - CLIs: nats, mc
# =============================================================

set -euo pipefail

COMMAND="${1:-}"
MESSAGE_COUNT="${2:-1000000}"

# ─── Source Modules ────────────────────────────────────────

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"

source "$SCRIPT_DIR/config.sh"
source "$COMMON_DIR/services/nats.sh"

# ─── Commands ──────────────────────────────────────────────

cmd_setup() {
    echo ""
    echo "================================================================"
    echo "  JS-Workflow-Executor Seed — Setup"
    echo "================================================================"
    echo "  Messages:      $(format_number "$MESSAGE_COUNT")"
    echo "  Stream:        $STREAM"
    echo "  Subject:       $SUBJECT"
    echo "  Workflow ID:   $BENCH_WORKFLOW_ID"
    echo "  Node ID:       $BENCH_NODE_ID"
    echo "================================================================"
    echo ""

    # Seed MinIO (TieredCache L2 — workflow script source)
    log_step "1" "2" "Seeding MinIO (workflow script source)..."
    bench_seed_workflow_minio

    # Ensure NATS streams exist + populate
    echo ""
    log_step "2" "2" "Populating NATS stream..."
    bench_ensure_streams
    bench_remove_consumer
    bench_seed_messages "$MESSAGE_COUNT" "$SEED_DIR/nats/workflow-code-execute.json"
    bench_verify_stream "$MESSAGE_COUNT"

    echo ""
    log_success "Setup complete. Run the benchmark:"
    echo "    bash docs/benchmarks/scripts/full-benchmark.sh $MESSAGE_COUNT"
    echo ""
}

cmd_teardown() {
    echo ""
    echo "================================================================"
    echo "  JS-Workflow-Executor Seed — Teardown"
    echo "================================================================"
    echo ""

    # Purge NATS streams
    log_step "1" "2" "Purging NATS streams..."
    bench_purge_all_streams

    # Cleanup MinIO
    echo ""
    log_step "2" "2" "Cleaning MinIO..."
    bench_cleanup_workflow_minio

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
        echo "  setup [N]   Seed MinIO + publish N messages (default: 1,000,000)"
        echo "  teardown    Purge streams + cleanup MinIO"
        echo ""
        echo "  Examples:"
        echo "    $0 setup              # 1M messages"
        echo "    $0 setup 10000        # 10K messages (quick test)"
        echo "    $0 teardown           # clean everything"
        exit 1
        ;;
esac
