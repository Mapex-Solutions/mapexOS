#!/bin/bash
# =============================================================
# Events Service Full CPU Benchmark (cgroup v2 isolated)
#
# Strategy: "Stream Drain" — seed N messages into NATS JetStream,
# start the events service, measure time to drain all messages.
#
# Supports ALL 7 event consumers:
#   save_raw_event, save_jsexec_event, save_router_event,
#   save_businessrule_event, save_trigger_event, save_event,
#   save_dlq_event
#
# Flow:
#   1. Sudo + cgroup verification (fail-fast)
#   2. Preflight checks (infra health)
#   3. Build Go binary
#   4. Seed data (ClickHouse tables, MongoDB retention, Redis flush)
#   5. For each scenario x CPU config:
#        a. Purge NATS streams
#        b. Seed N messages into scenario's stream
#        c. Start service in isolated cgroup
#        d. Poll drain until all messages consumed
#        e. Capture metrics (Prometheus + Go runtime)
#        f. Stop + cool down
#   6. Teardown via seed.sh
#
# Prerequisites:
#   1. cgroup shield: sudo bash scripts/benchmarks/common/cgroup/setup-cgroup-shield.sh
#   2. ClickHouse, NATS, MongoDB, Redis running
#   3. nats, mongosh, redis-cli, curl CLIs
# =============================================================

set -euo pipefail

# ─── Source Modules ──────────────────────────────────────────

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"

source "$SCRIPT_DIR/config.sh"
source "$COMMON_DIR/services/nats.sh"
source "$COMMON_DIR/cgroup/cgroup.sh"

# ─── Help ───────────────────────────────────────────────────

show_help() {
    cat <<'HELP'
Events Service Full CPU Benchmark (Stream Drain)

Usage:
  ./full-benchmark.sh [message_count] [cpu_list] [nats_batch_size]

Arguments:
  message_count   Messages per test                       (default: 1000000)
  cpu_list        CPU configs to test, quoted             (default: "1 2 4 8 16")
  nats_batch_size NATS batch size for the service         (default: 10000)

Environment variables:
  SCENARIOS       Override which scenarios to run          (default: all)
                  save_raw_event          Raw events from HTTP/MQTT gateways (lightest)
                  save_jsexec_event       JS executor debug logs
                  save_router_event       Router execution history
                  save_businessrule_event Business rule evaluation (heaviest — 5x JSON marshal)
                  save_trigger_event      Trigger execution logs
                  save_event              EventStore with EVA field resolution
                  save_dlq_event          Dead Letter Queue (special: never nacks)
  GO_ENV          Service environment                     (default: dev)
  SUDO_PASS       sudo password for cgroup operations     (prompted if not set)

Examples:
  ./full-benchmark.sh                                                         # all, 1M, CPU 1-16
  ./full-benchmark.sh 10000 "1 2"                                             # all, 10K, quick test
  ./full-benchmark.sh 1000000 "16 12 8 4 2 1"                                # all, 1M, CPUs desc
  ./full-benchmark.sh 1000000 "1 2 4 8 16" 10000                             # all, 1M, batch 10K
  SCENARIOS="save_raw_event" ./full-benchmark.sh 10000 "1 2"                  # only raw
  SCENARIOS="save_businessrule_event" ./full-benchmark.sh 1000000 "1 2 4 8"   # heaviest consumer
  SCENARIOS="save_raw_event save_event" ./full-benchmark.sh 100000 "4 8"      # two scenarios
HELP
    exit 0
}

case "${1:-}" in -h|--help|help) show_help ;; esac

# ─── Configuration ───────────────────────────────────────────

MESSAGE_COUNT="${1:-1000000}"
CPU_LIST="${2:-1 2 4 8 16}"
NATS_BATCH="${3:-10000}"

# Scenarios — run all by default, override via env var
SCENARIOS_TO_RUN="${SCENARIOS:-$ALL_SCENARIOS}"

# Validate scenarios
for _sc in $SCENARIOS_TO_RUN; do
    validate_scenario "$_sc" || exit 1
done

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

# Paths
SERVICE_DIR="$(cd "$SCRIPT_DIR/../../.." && pwd)"
RESULTS_DIR="$SCRIPT_DIR/../results/$BENCH_TAG"
BINARY_PATH="$SCRIPT_DIR/events_bench"
PAYLOAD_DIR="$SEED_DIR/payloads"

# Drain polling
POLL_INTERVAL=2

# ─── Kill Competing Services ────────────────────────────────
kill_competing_services

# ─── Preflight Checks ────────────────────────────────────────

run_preflight_checks "" "$BENCH_CLI_TOOLS" "$BENCH_SERVICE_CHECKS"

# Verify payload files exist
for scen in $SCENARIOS_TO_RUN; do
    if [ ! -f "$PAYLOAD_DIR/${scen}.json" ]; then
        log_error "Payload not found: $PAYLOAD_DIR/${scen}.json"
        exit 1
    fi
done
log_success "All payload files verified."

# ─── Build Binary ────────────────────────────────────────────

echo ""
log_info "Compiling events service..."
cd "$SERVICE_DIR"
GOWORK=off CGO_ENABLED=0 go build -o "$BINARY_PATH" ./src/ 2>&1
log_success "Binary: $BINARY_PATH"

# ─── Seed Data ───────────────────────────────────────────────

echo ""
log_info "Seeding benchmark data via seed.sh..."
bash "$SCRIPT_DIR/seed.sh" setup

# Kill stale processes
kill_service_on_port "$SERVICE_PORT"

mkdir -p "$RESULTS_DIR"

# ─── Poll Drain ──────────────────────────────────────────────

# poll_drain waits until the events service has processed all messages.
# Polls the /metrics endpoint and sums ack+nack+reject+dlq counters.
poll_drain() {
    local target="$1"
    local start_time="$2"
    local consumer_label="$3"
    local last_total=0
    local stall_count=0

    while true; do
        local metrics
        metrics=$(curl -s "$METRICS_URL" 2>/dev/null) || { sleep "$POLL_INTERVAL"; continue; }

        local ack nack reject dlq_count total
        ack=$(echo "$metrics" | grep "events_message_total{.*consumer=\"${consumer_label}\".*result=\"ack\"}" 2>/dev/null | head -1 | awk '{print $NF}') || true
        nack=$(echo "$metrics" | grep "events_message_total{.*consumer=\"${consumer_label}\".*result=\"nack\"}" 2>/dev/null | head -1 | awk '{print $NF}') || true
        reject=$(echo "$metrics" | grep "events_message_total{.*consumer=\"${consumer_label}\".*result=\"reject\"}" 2>/dev/null | head -1 | awk '{print $NF}') || true
        dlq_count=$(echo "$metrics" | grep "events_message_total{.*consumer=\"${consumer_label}\".*result=\"dlq\"}" 2>/dev/null | head -1 | awk '{print $NF}') || true

        # awk handles scientific notation (1e+06)
        total=$(awk "BEGIN{printf \"%.0f\", ${ack:-0}+${nack:-0}+${reject:-0}+${dlq_count:-0}}")

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
                log_warn "Drain stalled at $total / $target (no progress for $((stall_count * POLL_INTERVAL))s)"
                return 1
            fi
        else
            stall_count=0
            last_total=$total
        fi

        sleep "$POLL_INTERVAL"
    done
}

# ─── Count total tests ───────────────────────────────────────

NUM_SCENARIOS=$(echo "$SCENARIOS_TO_RUN" | wc -w)
NUM_CPU=$(echo "$CPU_LIST" | wc -w)
TOTAL_TESTS=$((NUM_SCENARIOS * NUM_CPU))

echo ""
echo "================================================================"
echo "  Events Service Stream Drain Benchmark (cgroup v2)"
echo "================================================================"
echo "  Scenarios:      $SCENARIOS_TO_RUN"
echo "  Messages/test:  $(format_number "$MESSAGE_COUNT")"
echo "  Batch size:     $(format_number "$NATS_BATCH")"
echo "  CPU list:       $CPU_LIST"
echo "  Total tests:    $TOTAL_TESTS (${NUM_SCENARIOS} scenarios x ${NUM_CPU} CPU)"
echo "  Results:        $RESULTS_DIR"
echo "================================================================"
echo ""

# ─── Main Benchmark Loop ─────────────────────────────────────

declare -a SUMMARY_ROWS=()
TEST_NUM=0

for CURRENT_SCENARIO in $SCENARIOS_TO_RUN; do
    log_info "Scenario block: $CURRENT_SCENARIO"

    # Resolve scenario config
    STREAM=$(get_nats_stream "$CURRENT_SCENARIO")
    SUBJECT=$(get_nats_subject "$CURRENT_SCENARIO")
    CONSUMER_SUFFIX=$(get_consumer_suffix "$CURRENT_SCENARIO")
    CONSUMER_LABEL=$(get_metric_consumer "$CURRENT_SCENARIO")
    TABLE_LABEL=$(get_metric_table "$CURRENT_SCENARIO")
    CH_TABLE=$(get_ch_table "$CURRENT_SCENARIO")
    CONSUMER="events-${CONSUMER_SUFFIX}"

    for CPU in $CPU_LIST; do
        TEST_NUM=$((TEST_NUM + 1))
        CORE_RANGE=$(set_shield_cpus "$CPU")

        echo ""
        echo "════════════════════════════════════════════════════════════════"
        echo "  TEST $TEST_NUM/$TOTAL_TESTS — scenario=$CURRENT_SCENARIO  GOMAXPROCS=$CPU  cores=$CORE_RANGE"
        echo "  $(format_number "$MESSAGE_COUNT") messages → $STREAM ($SUBJECT)"
        echo "════════════════════════════════════════════════════════════════"
        echo ""

        # 0. Kill any leftover process
        kill_service_on_port "$SERVICE_PORT"

        # 1. Purge NATS streams
        log_step "1" "7" "Purging NATS streams..."
        purge_streams "$STREAM" "MAPEXOS-DLQ"
        remove_consumer "$STREAM" "$CONSUMER"

        # 2. Seed messages
        log_step "2" "7" "Seeding $(format_number "$MESSAGE_COUNT") messages..."
        SEED_START=$(date +%s%N)
        seed_messages "$MESSAGE_COUNT" "$SUBJECT" "$PAYLOAD_DIR/${CURRENT_SCENARIO}.json"
        SEED_END=$(date +%s%N)
        SEED_SECS=$(( (SEED_END - SEED_START) / 1000000000 ))
        log_success "Seeded in ${SEED_SECS}s"

        # 3. Start service
        log_step "3" "7" "Starting events service (GOMAXPROCS=$CPU, NATS_BATCH_SIZE=$NATS_BATCH, LOG_LEVEL=silent)..."
        SERVICE_LOG="$RESULTS_DIR/test-${CURRENT_SCENARIO}-cpu${CPU}-output.log"
        env GO_ENV="$GO_ENV_VALUE" LOG_LEVEL=silent GOMAXPROCS="$CPU" \
            NATS_BATCH_SIZE="$NATS_BATCH" \
            METRICS_GO_COLLECTOR=true METRICS_PROCESS_COLLECTOR=true \
            "$BINARY_PATH" > "$SERVICE_LOG" 2>&1 &
        SVC_PID=$!
        sleep 1
        move_to_shield "$SVC_PID" || log_warn "Failed to move PID $SVC_PID to cgroup"
        echo "      PID=$SVC_PID -> cores $CORE_RANGE"

        ACTUAL_CPUS=$(grep Cpus_allowed_list /proc/$SVC_PID/status 2>/dev/null | awk '{print $2}') || true
        echo "      Verified: $ACTUAL_CPUS"

        # 4. Wait for ready
        log_step "4" "7" "Waiting for service..."
        if ! wait_for_service_ready "$METRICS_URL" 60 "$SVC_PID"; then
            echo "      SKIP scenario=$CURRENT_SCENARIO CPU=$CPU"
            kill "$SVC_PID" 2>/dev/null || true; wait "$SVC_PID" 2>/dev/null || true
            continue
        fi

        # 5. Poll drain
        log_step "5" "7" "Draining $(format_number "$MESSAGE_COUNT") messages..."
        DRAIN_START=$(date +%s%N)
        DRAIN_OK=true
        poll_drain "$MESSAGE_COUNT" "$DRAIN_START" "$CONSUMER_LABEL" || DRAIN_OK=false
        DRAIN_END=$(date +%s%N)
        DRAIN_SECS=$(awk "BEGIN{printf \"%.1f\", (${DRAIN_END}-${DRAIN_START})/1000000000}")

        # 6. Collect metrics
        log_step "6" "7" "Collecting metrics..."
        sleep 2
        METRICS_FILE="$RESULTS_DIR/test-${CURRENT_SCENARIO}-cpu${CPU}-metrics.txt"
        METRICS=$(curl -s "$METRICS_URL" 2>/dev/null)
        echo "$METRICS" > "$METRICS_FILE"

        # Parse Go runtime (via common module)
        parse_go_runtime_metrics "$METRICS"

        # Parse event processing metrics (using scenario labels)
        EVENTS_OK=$(extract_metric "$METRICS" "events_event_processed_total" "consumer=\"${CONSUMER_LABEL}\",status=\"success\"")
        EVENTS_ERR=$(extract_metric "$METRICS" "events_event_processed_total" "consumer=\"${CONSUMER_LABEL}\",status=\"error\"")
        PROC_AVG_MS=$(extract_histogram_avg "$METRICS" "events_event_processing_duration_seconds" "consumer=\"${CONSUMER_LABEL}\"")

        # Message lifecycle
        MSGS_ACK=$(extract_metric "$METRICS" "events_message_total" "consumer=\"${CONSUMER_LABEL}\",result=\"ack\"")
        MSGS_NACK=$(extract_metric "$METRICS" "events_message_total" "consumer=\"${CONSUMER_LABEL}\",result=\"nack\"")
        MSGS_REJECT=$(extract_metric "$METRICS" "events_message_total" "consumer=\"${CONSUMER_LABEL}\",result=\"reject\"")

        # ClickHouse insert performance (labels alphabetical: status before table)
        CH_INSERT_OK=$(extract_metric "$METRICS" "events_clickhouse_insert_total" "status=\"ok\",table=\"${TABLE_LABEL}\"")
        CH_INSERT_ERR=$(extract_metric "$METRICS" "events_clickhouse_insert_total" "status=\"error\",table=\"${TABLE_LABEL}\"")
        CH_INSERT_AVG_MS=$(extract_histogram_avg "$METRICS" "events_clickhouse_insert_duration_seconds" "table=\"${TABLE_LABEL}\"")

        # CH batch size: extract avg directly (sum/count, no *1000 since it's not time)
        ch_batch_sum=$(echo "$METRICS" | grep "^events_clickhouse_insert_batch_size_sum{.*table=\"${TABLE_LABEL}\"" 2>/dev/null | head -1 | awk '{print $NF}') || true
        ch_batch_count=$(echo "$METRICS" | grep "^events_clickhouse_insert_batch_size_count{.*table=\"${TABLE_LABEL}\"" 2>/dev/null | head -1 | awk '{print $NF}') || true
        if [ -n "$ch_batch_sum" ] && [ -n "$ch_batch_count" ] && [ "$ch_batch_count" != "0" ]; then
            CH_BATCH_AVG=$(awk "BEGIN{printf \"%.0f\", ${ch_batch_sum}/${ch_batch_count}}")
        else
            CH_BATCH_AVG="N/A"
        fi

        # Compute events/s: prefer drain time, fallback to processing duration
        EVENTS_PER_SEC="0"
        if [ "$DRAIN_SECS" != "0" ] && [ "$DRAIN_SECS" != "0.0" ]; then
            EVENTS_PER_SEC=$(awk "BEGIN{printf \"%.0f\", ${EVENTS_OK}/${DRAIN_SECS}}")
        elif [ "$PROC_AVG_MS" != "N/A" ] && [ "$PROC_AVG_MS" != "0" ]; then
            # Drain too fast for polling — estimate from Prometheus processing duration
            EVENTS_PER_SEC=$(awk "BEGIN{printf \"%.0f\", ${EVENTS_OK}/(${PROC_AVG_MS}/1000)}")
            DRAIN_SECS=$(awk "BEGIN{printf \"%.1f\", ${PROC_AVG_MS}/1000}")
        fi

        # 7. Stop service
        log_step "7" "7" "Stopping service..."
        kill "$SVC_PID" 2>/dev/null || true
        wait "$SVC_PID" 2>/dev/null || true
        echo "      Done."

        # Summary row
        SUMMARY_ROWS+=("$(printf "  │ %-24s │ %3s │ %8s │ %7s/s │ %5ss │ %8s │ %5sMB │" \
            "$CURRENT_SCENARIO" "$CPU" "$CORE_RANGE" "$(format_number "$EVENTS_PER_SEC")" "$DRAIN_SECS" "$(format_number "$CH_INSERT_OK")" "$GO_RSS_MB")")

        # Inline summary
        echo ""
        echo "  ┌──────────────────────────────────────────────────────────┐"
        echo "  │ scenario=$CURRENT_SCENARIO CPU=$CPU -> $(format_number "$EVENTS_PER_SEC")/s  (${DRAIN_SECS}s drain)"
        echo "  │ Ack=$(format_number "$MSGS_ACK")  Proc=${PROC_AVG_MS}ms"
        echo "  │ CH Insert: ok=$(format_number "$CH_INSERT_OK") err=$(format_number "$CH_INSERT_ERR") avg=${CH_INSERT_AVG_MS}ms"
        echo "  │ CH Batch avg=${CH_BATCH_AVG}"
        echo "  │ RSS=${GO_RSS_MB}MB  Heap=${GO_HEAP_MB}MB  Goroutines=${GO_GOROUTINES}"
        echo "  └──────────────────────────────────────────────────────────┘"

        sleep 5
    done
done

# ─── Final Summary ───────────────────────────────────────────

echo ""
echo ""
echo "================================================================"
echo "  COMPLETE — ${TOTAL_TESTS} tests | scenarios: ${SCENARIOS_TO_RUN}"
echo "================================================================"
echo ""

TABLE_HEADER=$(printf "  │ %-24s │ %3s │ %8s │ %10s │ %8s │ %8s │ %7s │" \
    "Scenario" "CPU" "Cores" "Events/s" "Drain" "CH OK" "RSS MB")

echo "  Summary:"
echo "  ┌──────────────────────────┬─────┬──────────┬────────────┬──────────┬──────────┬─────────┐"
echo "$TABLE_HEADER"
echo "  ├──────────────────────────┼─────┼──────────┼────────────┼──────────┼──────────┼─────────┤"

for ROW in "${SUMMARY_ROWS[@]}"; do
    echo "$ROW"
done

echo "  └──────────────────────────┴─────┴──────────┴────────────┴──────────┴──────────┴─────────┘"
echo ""

# ─── Teardown ─────────────────────────────────────────────────

log_info "Removing benchmark data via seed.sh teardown..."
bash "$SCRIPT_DIR/seed.sh" teardown

# Remove temp binary
rm -f "$BINARY_PATH"

log_success "Done!"
echo "  Results stored at: $RESULTS_DIR"
