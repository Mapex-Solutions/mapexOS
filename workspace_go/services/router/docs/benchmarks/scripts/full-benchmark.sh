#!/bin/bash
# =============================================================
# Router Service Full CPU Benchmark (cgroup v2 isolated)
#
# Strategy: "Stream Drain" — seed N messages into NATS JetStream,
# start the router service, measure time to drain all messages.
#
# Self-contained: collects ALL metrics from /metrics endpoint.
# No Grafana or Prometheus reset needed.
#
# Runs ALL scenarios (save_event, rule_engine, trigger) by default.
# Override via SCENARIOS env var.
#
# Flow:
#   1. Preflight checks (infra health + user confirmation)
#   2. Build Go binary
#   3. Seed MongoDB (3 RouteGroups) + MinIO (asset) — once before all tests
#   4. For each scenario:
#      For each CPU config:
#        a. Kill old process → purge ALL NATS streams
#        b. Seed N messages to NATS (for current scenario)
#        c. Start service in isolated cgroup
#        d. Poll drain until ack+nack+reject >= N
#        e. Capture metrics (Prometheus + Go runtime)
#        f. Stop → cool down
#   5. Teardown via seed.sh
#
# Prerequisites:
#   1. cgroup shield: sudo bash scripts/benchmarks/common/cgroup/setup-cgroup-shield.sh
#   2. MongoDB, NATS, Redis, MinIO running
#   3. nats, mongosh, mc, redis-cli, curl CLIs installed
#
# Usage:
#   ./full-benchmark.sh [message_count] [cpu_list] [tag] [nats_batch_size]
#
# Examples:
#   ./full-benchmark.sh                                     # all scenarios, 1M, CPU 1 2 4 8 16
#   ./full-benchmark.sh 100000 "4 2 1"                      # all scenarios, 100K, CPU 4 2 1
#   ./full-benchmark.sh 1000 "1" smoke-test                 # all scenarios, 1K, 1 CPU
#   SCENARIOS="save_event" ./full-benchmark.sh 1000 "1" quick  # single scenario
# =============================================================

set -euo pipefail

# ─── Source Modules ──────────────────────────────────────────

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"

source "$SCRIPT_DIR/config.sh"
source "$COMMON_DIR/services/routegroups.sh"
source "$COMMON_DIR/services/minio.sh"
source "$COMMON_DIR/services/nats.sh"
source "$COMMON_DIR/cgroup/cgroup.sh"

# ─── Configuration ───────────────────────────────────────────
MESSAGE_COUNT="${1:-1000000}"
CPU_LIST="${2:-1 2 4 8 16}"
BENCH_TAG="${3:-baseline}"
NATS_BATCH="${4:-8000}"
# Prompt for sudo password (required for cgroup CPU isolation)
if [ -z "${SUDO_PASS:-}" ]; then
    read -s -p "[sudo] password for cgroup operations: " SUDO_PASS
    echo ""
fi
export SUDO_PASS

# Verify sudo + cgroup shield (fail fast before compile/seed)
if ! echo "$SUDO_PASS" | sudo -S true 2>/dev/null; then
    log_error "Invalid sudo password. Cannot manage cgroup."
    exit 1
fi
verify_shield || exit 1
validate_cpu_list "$CPU_LIST"

# Scenarios — run all by default, override via env var
SCENARIOS="${SCENARIOS:-save_event rule_engine trigger}"

# Validate scenarios
for _sc in $SCENARIOS; do
    case "$_sc" in
        save_event|rule_engine|trigger) ;;
        *) log_error "Invalid scenario '$_sc'. Use: save_event | rule_engine | trigger"; exit 1 ;;
    esac
done

# Paths
SERVICE_DIR="$(cd "$SCRIPT_DIR/../../.." && pwd)"
RESULTS_DIR="$SCRIPT_DIR/../results/$BENCH_TAG"
BINARY_PATH="$SCRIPT_DIR/router_bench"

# Drain polling
POLL_INTERVAL=2

# ─── Local Helper Functions ──────────────────────────────────

# poll_drain waits until the router has processed all messages.
# Polls the /metrics endpoint and sums ack+nack+reject counters.
poll_drain() {
    local target="$1"
    local start_time="$2"
    local last_total=0
    local stall_count=0

    while true; do
        local metrics
        metrics=$(curl -s "$METRICS_URL" 2>/dev/null) || { sleep "$POLL_INTERVAL"; continue; }

        local ack nack reject total
        ack=$(echo "$metrics" | grep 'router_message_total{result="ack"}' 2>/dev/null | head -1 | awk '{print $NF}') || true
        nack=$(echo "$metrics" | grep 'router_message_total{result="nack"}' 2>/dev/null | head -1 | awk '{print $NF}') || true
        reject=$(echo "$metrics" | grep 'router_message_total{result="reject"}' 2>/dev/null | head -1 | awk '{print $NF}') || true

        # awk handles scientific notation (1e+06)
        total=$(awk "BEGIN{printf \"%.0f\", ${ack:-0}+${nack:-0}+${reject:-0}}")

        local elapsed
        elapsed=$(( ($(date +%s%N) - start_time) / 1000000000 ))
        local rate=0
        [ "$elapsed" -gt 0 ] && rate=$(awk "BEGIN{printf \"%.0f\", ${total}/${elapsed}}")

        printf "\r      Progress: %s / %s  (%s/s)  %ds" \
            "$(format_number "$total")" \
            "$(format_number "$target")" \
            "$(format_number "$rate")" \
            "$elapsed"

        if [ "$total" -ge "$target" ]; then
            echo ""
            return 0
        fi

        # Stall detection: if total hasn't changed in 10 polls, abort
        if [ "$total" -eq "$last_total" ]; then
            stall_count=$((stall_count + 1))
            if [ "$stall_count" -ge 10 ]; then
                echo ""
                echo "      WARNING: Drain stalled at $total / $target (no progress for $((stall_count * POLL_INTERVAL))s)"
                return 1
            fi
        else
            stall_count=0
            last_total=$total
        fi

        sleep "$POLL_INTERVAL"
    done
}

# ─── Kill Competing Services ───────────────────────────────
# Ensure no other MapexOS services interfere (keeps Assets alive for tiered cache)
kill_competing_services

# ─── Preflight Checks ───────────────────────────────────────

# Infrastructure health checks
run_preflight_checks "" "$BENCH_CLI_TOOLS" "$BENCH_SERVICE_CHECKS"

# ─── Build Binary ────────────────────────────────────────────

echo ""
log_info "Compiling router service..."
cd "$SERVICE_DIR"
GOWORK=off CGO_ENABLED=0 go build -o "$BINARY_PATH" ./src/ 2>&1
log_success "Binary: $BINARY_PATH"

# ─── Seed Data ──────────────────────────────────────────────

echo ""
log_info "Seeding benchmark data (MongoDB + MinIO) via seed.sh..."
bash "$SCRIPT_DIR/seed.sh" setup

# Validate payload files
for _sc in $SCENARIOS; do
    local_payload="$SEED_DIR/nats/${_sc}.json"
    if [ ! -f "$local_payload" ]; then
        log_error "Payload file not found: $local_payload"
        exit 1
    fi
done
PAYLOAD_SIZE=$(wc -c < "$SEED_DIR/nats/save_event.json")
log_success "Payloads OK (${PAYLOAD_SIZE} bytes each)"

# Kill stale processes
kill_service_on_port "$SERVICE_PORT"

# ─── Main Benchmark Loop ─────────────────────────────────────

mkdir -p "$RESULTS_DIR"

# Count total tests: SCENARIOS × CPU_LIST
NUM_SCENARIOS=$(echo "$SCENARIOS" | wc -w)
NUM_CPU=$(echo "$CPU_LIST" | wc -w)
TOTAL_TESTS=$((NUM_SCENARIOS * NUM_CPU))

echo ""
echo "================================================================"
echo "  Router Service Stream Drain Benchmark (cgroup v2)"
echo "================================================================"
echo "  Scenarios:      $SCENARIOS"
echo "  Messages/test:  $(format_number $MESSAGE_COUNT)"
echo "  CPU list:       $CPU_LIST"
echo "  Total tests:    $TOTAL_TESTS (${NUM_SCENARIOS} scenarios x ${NUM_CPU} CPU)"
echo "  Tag:            $BENCH_TAG"
echo "  NATS batch:     $NATS_BATCH"
echo "  MongoDB:        $MONGO_DB"
echo "  NATS stream:    $STREAM -> $SUBJECT"
echo "  Results:        $RESULTS_DIR"
echo "================================================================"
echo ""

declare -a SUMMARY_ROWS=()
TEST_NUM=0

for SCENARIO in $SCENARIOS; do
    log_info "Scenario block: $SCENARIO"

    for CPU in $CPU_LIST; do
        TEST_NUM=$((TEST_NUM + 1))
        CORE_RANGE=$(set_shield_cpus "$CPU")

        echo ""
        echo "════════════════════════════════════════════════════════════════"
        echo "  TEST $TEST_NUM/$TOTAL_TESTS — scenario=$SCENARIO  GOMAXPROCS=$CPU  cores=$CORE_RANGE"
        echo "  $(format_number $MESSAGE_COUNT) messages (stream drain)"
        echo "════════════════════════════════════════════════════════════════"
        echo ""

        # 0. Kill any leftover process
        kill_service_on_port "$SERVICE_PORT"

        # 1. Purge ALL NATS streams
        log_step "1" "7" "Purging ALL NATS streams..."
        bench_purge_all_streams

        # 2. Seed messages for this scenario
        log_step "2" "7" "Seeding NATS stream (${SCENARIO})..."
        SEED_START=$(date +%s%N)
        bench_seed_messages "$MESSAGE_COUNT" "$SEED_DIR/nats/${SCENARIO}.json"
        SEED_END=$(date +%s%N)
        SEED_SECS=$(( (SEED_END - SEED_START) / 1000000000 ))
        echo "      Seeded in ${SEED_SECS}s"

        # 3. Start service (LOG_LEVEL=silent for zero noise)
        log_step "3" "7" "Starting router (GOMAXPROCS=$CPU, NATS_BATCH_SIZE=$NATS_BATCH, LOG_LEVEL=silent)..."
        SERVICE_LOG="$RESULTS_DIR/test-${SCENARIO}-cpu${CPU}-output.log"
        env GO_ENV="$GO_ENV_VALUE" LOG_LEVEL=silent GOMAXPROCS="$CPU" \
            NATS_BATCH_SIZE="$NATS_BATCH" \
            METRICS_GO_COLLECTOR=true METRICS_PROCESS_COLLECTOR=true \
            "$BINARY_PATH" > "$SERVICE_LOG" 2>&1 &
        ROUTER_PID=$!
        sleep 1
        move_to_shield "$ROUTER_PID" || log_warn "Failed to move PID $ROUTER_PID to cgroup"
        echo "      PID=$ROUTER_PID -> cores $CORE_RANGE"

        ACTUAL_CPUS=$(grep Cpus_allowed_list /proc/$ROUTER_PID/status 2>/dev/null | awk '{print $2}') || true
        echo "      Verified: $ACTUAL_CPUS"

        # 4. Wait for ready
        log_step "4" "7" "Waiting for service..."
        if ! wait_for_service_ready "$METRICS_URL" 60 "$ROUTER_PID"; then
            echo "      SKIP scenario=$SCENARIO CPU=$CPU"
            kill "$ROUTER_PID" 2>/dev/null || true; wait "$ROUTER_PID" 2>/dev/null || true
            continue
        fi

        # 5. Poll drain
        log_step "5" "7" "Draining $(format_number $MESSAGE_COUNT) messages..."
        DRAIN_START=$(date +%s%N)
        DRAIN_OK=true
        poll_drain "$MESSAGE_COUNT" "$DRAIN_START" || DRAIN_OK=false
        DRAIN_END=$(date +%s%N)
        DRAIN_SECS=$(awk "BEGIN{printf \"%.1f\", (${DRAIN_END}-${DRAIN_START})/1000000000}")

        # 6. Collect metrics
        log_step "6" "7" "Collecting metrics..."
        sleep 2
        METRICS_FILE="$RESULTS_DIR/test-${SCENARIO}-cpu${CPU}-metrics.txt"
        METRICS=$(curl -s "$METRICS_URL" 2>/dev/null)
        echo "$METRICS" > "$METRICS_FILE"

        # Parse service metrics — event processing
        EVENTS_OK=$(extract_metric "$METRICS" "router_event_processed_total" 'status="success"')
        EVENTS_ERR=$(extract_metric "$METRICS" "router_event_processed_total" 'status="error"')
        PROC_AVG_MS=$(extract_histogram_avg "$METRICS" "router_event_processing_duration_seconds")

        # Message lifecycle
        MSGS_ACK=$(extract_metric "$METRICS" "router_message_total" 'result="ack"')
        MSGS_NACK=$(extract_metric "$METRICS" "router_message_total" 'result="nack"')
        MSGS_REJECT=$(extract_metric "$METRICS" "router_message_total" 'result="reject"')

        # NATS publish
        PUB_OK=$(extract_metric "$METRICS" "router_event_published_total" 'status="success"')
        PUB_ERR=$(extract_metric "$METRICS" "router_event_published_total" 'status="error"')

        # Parse Go runtime (via common module)
        parse_go_runtime_metrics "$METRICS"

        # Compute events/s from drain time
        EVENTS_PER_SEC="0"
        if [ "$DRAIN_SECS" != "0" ] && [ "$DRAIN_SECS" != "0.0" ]; then
            EVENTS_PER_SEC=$(awk "BEGIN{printf \"%.0f\", ${EVENTS_OK}/${DRAIN_SECS}}")
        fi

        # 7. Stop service
        log_step "7" "7" "Stopping service..."
        kill "$ROUTER_PID" 2>/dev/null || true
        wait "$ROUTER_PID" 2>/dev/null || true
        echo "      Done."

        # Summary row for final table (with Scenario column)
        SUMMARY_ROWS+=("$(printf "  │ %-12s │ %3s │ %8s │ %7s/s │ %5ss │ %8s │ %8s │ %5sMB │" \
            "$SCENARIO" "$CPU" "$CORE_RANGE" "$(format_number "$EVENTS_PER_SEC")" "$DRAIN_SECS" "$(format_number "$MSGS_ACK")" "$(format_number "$PUB_OK")" "$GO_RSS_MB")")

        # Inline summary
        echo ""
        echo "  ┌──────────────────────────────────────────────────────────┐"
        echo "  │ scenario=$SCENARIO CPU=$CPU -> $(format_number "$EVENTS_PER_SEC")/s  (${DRAIN_SECS}s drain)"
        echo "  │ Ack=$(format_number "$MSGS_ACK")  Pub=$(format_number "$PUB_OK")  Proc=${PROC_AVG_MS}ms"
        echo "  │ RSS=${GO_RSS_MB}MB  Heap=${GO_HEAP_MB}MB  Goroutines=${GO_GOROUTINES}"
        echo "  └──────────────────────────────────────────────────────────┘"

        sleep 5
    done
done

# ─── Final Summary ───────────────────────────────────────────

echo ""
echo ""
echo "================================================================"
echo "  COMPLETE — ${TOTAL_TESTS} tests | scenarios: ${SCENARIOS}"
echo "================================================================"
echo ""

# Summary table (with Scenario column)
TABLE_HEADER=$(printf "  │ %-12s │ %3s │ %8s │ %10s │ %8s │ %8s │ %8s │ %7s │" \
    "Scenario" "CPU" "Cores" "Events/s" "Drain" "Ack" "Pub OK" "RSS MB")

echo "  Summary:"
echo "  ┌──────────────┬─────┬──────────┬────────────┬──────────┬──────────┬──────────┬─────────┐"
echo "$TABLE_HEADER"
echo "  ├──────────────┼─────┼──────────┼────────────┼──────────┼──────────┼──────────┼─────────┤"

for ROW in "${SUMMARY_ROWS[@]}"; do
    echo "$ROW"
done

echo "  └──────────────┴─────┴──────────┴────────────┴──────────┴──────────┴──────────┴─────────┘"
echo ""

# ─── Cleanup Benchmark Data ─────────────────────────────────
# Delegate to seed.sh teardown (separation of concerns per benchmark standard)

log_info "Removing benchmark data via seed.sh teardown..."
bash "$SCRIPT_DIR/seed.sh" teardown

# Remove temp binary
rm -f "$BINARY_PATH"

log_success "Done!"
echo "  Results stored at: $RESULTS_DIR"
