#!/bin/bash
# =============================================================
# Triggers Service Full CPU Benchmark (cgroup v2 isolated)
#
# Runs ALL trigger executor scenarios across multiple CPU
# configurations using NATS JetStream stream drain.
#
# Scenarios (executor-focused):
#   http_trigger     — HTTP POST executor
#   mqtt_trigger     — MQTT publish executor
#   nats_trigger     — NATS publish executor
#   rabbitmq_trigger — RabbitMQ publish executor
#   email_trigger    — Email SMTP executor
#
# Flow:
#   1. Preflight checks (infra health + cgroup)
#   2. Build Go binary + mock servers
#   3. Seed MongoDB (trigger documents)
#   4. For each scenario x CPU config:
#        a. Purge NATS streams + seed N messages
#        b. Start service in isolated cgroup
#        c. Poll drain until processed >= N
#        d. Capture metrics (Prometheus + Go runtime)
#        e. Stop + cool down
#   5. Teardown via seed.sh
#
# Prerequisites:
#   1. cgroup shield: sudo bash scripts/benchmarks/common/cgroup/setup-cgroup-shield.sh
#   2. MongoDB, NATS, Redis running
#   3. nats, mongosh, redis-cli, curl CLIs
# =============================================================

set -euo pipefail

# ─── Source Modules ──────────────────────────────────────────

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"

source "$SCRIPT_DIR/config.sh"
source "$COMMON_DIR/services/triggers.sh"
source "$COMMON_DIR/services/nats.sh"
source "$COMMON_DIR/cgroup/cgroup.sh"

# ─── Help ───────────────────────────────────────────────────

show_help() {
    cat <<'HELP'
Triggers Service Full CPU Benchmark

Usage:
  ./full-benchmark.sh [message_count] [cpu_list] [nats_batch_size]

Arguments:
  message_count     Messages per test                  (default: 1000000)
  cpu_list          CPU configs to test, quoted         (default: "1 2 4 8 16")
  nats_batch_size   NATS consumer batch size            (default: 500)

Environment variables:
  SCENARIOS         Override which scenarios to run     (default: all)
                    Values: http_trigger mqtt_trigger nats_trigger rabbitmq_trigger email_trigger
  GO_ENV            Service environment                 (default: dev)
  MONGO_DB          MongoDB database name               (default: ${GO_ENV}-triggers)
  SUDO_PASS         Sudo password for cgroup ops        (prompted if not set)

Examples:
  ./full-benchmark.sh                                         # all, 1M, CPU 1-16
  ./full-benchmark.sh 1000000 "16 12 8 4 2 1"                # all, 1M, CPUs desc
  ./full-benchmark.sh 100000 "4 2 1" 1000                    # all, 100K, batch 1000
  ./full-benchmark.sh 1000 "1"                                # all, 1K, 1 CPU (smoke)
  SCENARIOS="http_trigger" ./full-benchmark.sh 50000 "4 8"    # only HTTP executor
HELP
    exit 0
}

case "${1:-}" in -h|--help|help) show_help ;; esac

# ─── Configuration ───────────────────────────────────────────

MESSAGE_COUNT="${1:-1000000}"
CPU_LIST="${2:-1 2 4 8 16}"
NATS_BATCH="${3:-500}"
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

# Tag — auto-generated from timestamp
BENCH_TAG="$(date +%Y%m%d-%H%M%S)"

# Scenarios — all 5 executor types by default
SCENARIOS="${SCENARIOS:-http_trigger mqtt_trigger nats_trigger rabbitmq_trigger email_trigger}"

# Validate scenarios
for _sc in $SCENARIOS; do
    case "$_sc" in
        http_trigger|mqtt_trigger|nats_trigger|rabbitmq_trigger|email_trigger) ;;
        *) log_error "Invalid scenario '$_sc'. Use: http_trigger | mqtt_trigger | nats_trigger | rabbitmq_trigger | email_trigger"; exit 1 ;;
    esac
done

# Paths
SERVICE_DIR="$(cd "$SCRIPT_DIR/../../.." && pwd)"
RESULTS_DIR="$SCRIPT_DIR/../results/$BENCH_TAG"
BINARY_PATH="$SCRIPT_DIR/triggers_bench"
MOCK_SERVERS_BIN="/tmp/bench_mock_servers"

# Drain polling
POLL_INTERVAL=2

# ─── Local Helper Functions ──────────────────────────────────

# poll_drain waits until the triggers service has processed all messages.
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
        ack=$(echo "$metrics" | grep 'triggers_message_total{result="ack"}' 2>/dev/null | head -1 | awk '{print $NF}') || true
        nack=$(echo "$metrics" | grep 'triggers_message_total{result="nack"}' 2>/dev/null | head -1 | awk '{print $NF}') || true
        reject=$(echo "$metrics" | grep 'triggers_message_total{result="reject"}' 2>/dev/null | head -1 | awk '{print $NF}') || true

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
# Ensure no other MapexOS services interfere
kill_competing_services

# ─── Preflight Checks ───────────────────────────────────────

# Infrastructure health checks
run_preflight_checks "" "$BENCH_CLI_TOOLS" "$BENCH_SERVICE_CHECKS"

# ─── Build Binary ────────────────────────────────────────────

echo ""
log_info "Compiling triggers service..."
cd "$SERVICE_DIR"
GOWORK=off CGO_ENABLED=0 go build -o "$BINARY_PATH" ./src/ 2>&1
log_success "Binary: $BINARY_PATH"

# ─── Start Mock Servers ──────────────────────────────────────

echo ""
log_info "Building mock servers (HTTP, MQTT, NATS, RabbitMQ, SMTP)..."

# Kill any leftovers on mock ports
for _port in $MOCK_HTTP_PORT $MOCK_MQTT_PORT $MOCK_NATS_PORT $MOCK_RABBITMQ_PORT $MOCK_SMTP_PORT; do
    kill_service_on_port "$_port"
done

GOWORK=off go build -o "$MOCK_SERVERS_BIN" -ldflags="-s -w" "$SCRIPT_DIR/mock-http-server.go" 2>&1
"$MOCK_SERVERS_BIN" &
MOCK_PID=$!
sleep 1
if kill -0 "$MOCK_PID" 2>/dev/null; then
    log_success "Mock servers PID=$MOCK_PID (HTTP:$MOCK_HTTP_PORT MQTT:$MOCK_MQTT_PORT NATS:$MOCK_NATS_PORT RabbitMQ:$MOCK_RABBITMQ_PORT SMTP:$MOCK_SMTP_PORT)"
else
    log_error "Mock servers failed to start"
    exit 1
fi

# ─── Seed Data ──────────────────────────────────────────────

echo ""
log_info "Seeding benchmark data (MongoDB) via seed.sh..."
bash "$SCRIPT_DIR/seed.sh" setup

# Validate payload files
for _sc in $SCENARIOS; do
    local_payload="$SEED_DIR/nats/${_sc}.json"
    if [ ! -f "$local_payload" ]; then
        log_error "Payload file not found: $local_payload"
        exit 1
    fi
done
PAYLOAD_SIZE=$(wc -c < "$SEED_DIR/nats/http_trigger.json")
log_success "Payloads OK (${PAYLOAD_SIZE} bytes each)"

# Kill stale processes
kill_service_on_port "$SERVICE_PORT"

# ─── Main Benchmark Loop ─────────────────────────────────────

mkdir -p "$RESULTS_DIR"

# Count total tests: SCENARIOS x CPU_LIST
NUM_SCENARIOS=$(echo "$SCENARIOS" | wc -w)
NUM_CPU=$(echo "$CPU_LIST" | wc -w)
TOTAL_TESTS=$((NUM_SCENARIOS * NUM_CPU))

echo ""
echo "================================================================"
echo "  Triggers Service Stream Drain Benchmark (cgroup v2)"
echo "================================================================"
echo "  Scenarios:      $SCENARIOS"
echo "  Messages/test:  $(format_number $MESSAGE_COUNT)"
echo "  CPU list:       $CPU_LIST"
echo "  Total tests:    $TOTAL_TESTS (${NUM_SCENARIOS} scenarios x ${NUM_CPU} CPU)"
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
        log_step "3" "7" "Starting triggers (GOMAXPROCS=$CPU, NATS_BATCH_SIZE=$NATS_BATCH, LOG_LEVEL=silent)..."
        SERVICE_LOG="$RESULTS_DIR/test-${SCENARIO}-cpu${CPU}-output.log"
        env GO_ENV="$GO_ENV_VALUE" LOG_LEVEL=silent GOMAXPROCS="$CPU" \
            NATS_BATCH_SIZE="$NATS_BATCH" \
            METRICS_GO_COLLECTOR=true METRICS_PROCESS_COLLECTOR=true \
            "$BINARY_PATH" > "$SERVICE_LOG" 2>&1 &
        TRIGGERS_PID=$!
        sleep 1
        move_to_shield "$TRIGGERS_PID" || log_warn "Failed to move PID $TRIGGERS_PID to cgroup"
        echo "      PID=$TRIGGERS_PID -> cores $CORE_RANGE"

        ACTUAL_CPUS=$(grep Cpus_allowed_list /proc/$TRIGGERS_PID/status 2>/dev/null | awk '{print $2}') || true
        echo "      Verified: $ACTUAL_CPUS"

        # 4. Wait for ready
        log_step "4" "7" "Waiting for service..."
        if ! wait_for_service_ready "$METRICS_URL" 60 "$TRIGGERS_PID"; then
            echo "      SKIP scenario=$SCENARIO CPU=$CPU"
            kill "$TRIGGERS_PID" 2>/dev/null || true; wait "$TRIGGERS_PID" 2>/dev/null || true
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

        # Parse service metrics — trigger processing
        TRIGGERS_OK=$(extract_metric "$METRICS" "triggers_trigger_processed_total" 'status="success"')
        TRIGGERS_ERR=$(extract_metric "$METRICS" "triggers_trigger_processed_total" 'status="error"')

        # Message lifecycle
        MSGS_ACK=$(extract_metric "$METRICS" "triggers_message_total" 'result="ack"')
        MSGS_NACK=$(extract_metric "$METRICS" "triggers_message_total" 'result="nack"')
        MSGS_REJECT=$(extract_metric "$METRICS" "triggers_message_total" 'result="reject"')

        # NATS publish
        PUB_OK=$(extract_metric "$METRICS" "triggers_event_published_total" 'status="ok"')
        PUB_ERR=$(extract_metric "$METRICS" "triggers_event_published_total" 'status="error"')

        # Parse Go runtime (via common module)
        parse_go_runtime_metrics "$METRICS"

        # Compute triggers/s from drain time
        TRIGGERS_PER_SEC="0"
        if [ "$DRAIN_SECS" != "0" ] && [ "$DRAIN_SECS" != "0.0" ]; then
            TRIGGERS_PER_SEC=$(awk "BEGIN{printf \"%.0f\", ${TRIGGERS_OK}/${DRAIN_SECS}}")
        fi

        # 7. Stop service
        log_step "7" "7" "Stopping service..."
        kill "$TRIGGERS_PID" 2>/dev/null || true
        wait "$TRIGGERS_PID" 2>/dev/null || true
        echo "      Done."

        # Summary row for final table (with Scenario column)
        SUMMARY_ROWS+=("$(printf "  │ %-16s │ %3s │ %8s │ %7s/s │ %5ss │ %8s │ %8s │ %5sMB │" \
            "$SCENARIO" "$CPU" "$CORE_RANGE" "$(format_number "$TRIGGERS_PER_SEC")" "$DRAIN_SECS" "$(format_number "$MSGS_ACK")" "$(format_number "$PUB_OK")" "$GO_RSS_MB")")

        # Inline summary
        echo ""
        echo "  ┌──────────────────────────────────────────────────────────┐"
        echo "  │ scenario=$SCENARIO CPU=$CPU -> $(format_number "$TRIGGERS_PER_SEC")/s  (${DRAIN_SECS}s drain)"
        echo "  │ Ack=$(format_number "$MSGS_ACK")  Pub=$(format_number "$PUB_OK")"
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
TABLE_HEADER=$(printf "  │ %-16s │ %3s │ %8s │ %10s │ %8s │ %8s │ %8s │ %7s │" \
    "Scenario" "CPU" "Cores" "Trigs/s" "Drain" "Ack" "Pub OK" "RSS MB")

echo "  Summary:"
echo "  ┌──────────────────┬─────┬──────────┬────────────┬──────────┬──────────┬──────────┬─────────┐"
echo "$TABLE_HEADER"
echo "  ├──────────────────┼─────┼──────────┼────────────┼──────────┼──────────┼──────────┼─────────┤"

for ROW in "${SUMMARY_ROWS[@]}"; do
    echo "$ROW"
done

echo "  └──────────────────┴─────┴──────────┴────────────┴──────────┴──────────┴──────────┴─────────┘"
echo ""

# ─── Cleanup Benchmark Data ─────────────────────────────────
# Delegate to seed.sh teardown (separation of concerns per benchmark standard)

# Stop mock servers first
if [ -n "${MOCK_PID:-}" ]; then
    kill "$MOCK_PID" 2>/dev/null || true
    wait "$MOCK_PID" 2>/dev/null || true
    log_info "Mock servers: stopped"
fi

log_info "Removing benchmark data via seed.sh teardown..."
bash "$SCRIPT_DIR/seed.sh" teardown

# Remove temp binaries
rm -f "$BINARY_PATH"
rm -f "$MOCK_SERVERS_BIN"

log_success "Done!"
echo "  Results stored at: $RESULTS_DIR"
