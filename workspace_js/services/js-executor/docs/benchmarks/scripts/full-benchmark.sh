#!/bin/bash
# =============================================================
# Full Automated CPU Benchmark — Piscina Worker Threads
#
# Uses cgroup v2 cpuset shield to truly isolate CPU cores.
# System processes are confined to cores 0-15, benchmark runs
# exclusively on cores 16-31 (up to 16 isolated cores).
#
# Architecture:
#   CPU_LIMIT=N → batchWorkers = N-1, batchSize = N × 500
#   Each worker: own V8 Isolate + NATS conn + scriptCache
#
# Flow:
#   1. Preflight checks (infra health + user confirmation)
#   2. Seed assets (MongoDB + MinIO) — once before all tests
#   3. For each CPU config:
#      a. Purge NATS streams
#      b. Populate stream with N messages
#      c. Start js-executor in isolated cgroup
#      d. Monitor drain until complete
#      e. Capture Prometheus metrics
#      f. Stop js-executor
#   4. Teardown assets
#
# Prerequisites:
#   sudo bash scripts/benchmarks/common/cgroup/setup-cgroup-shield.sh   (run once)
#
# Usage:
#   ./full-benchmark.sh [event_count] [cpu_list] [tag] [extra_env]
#
# Examples:
#   ./full-benchmark.sh                          # 1M events, CPU 16,8,4,2,1
#   ./full-benchmark.sh 500000 "8 4 2 1"         # 500K, specific CPUs
#   ./full-benchmark.sh 1000 "1" test             # Quick test
# =============================================================

set -euo pipefail

# ─── NVM Setup ────────────────────────────────────────────
export NVM_DIR="$HOME/.nvm"
[ -s "$NVM_DIR/nvm.sh" ] && . "$NVM_DIR/nvm.sh"
nvm use 24 > /dev/null 2>&1
echo "Using Node.js $(node --version)"

# ─── Source Modules ────────────────────────────────────────

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"

source "$SCRIPT_DIR/config.sh"
source "$COMMON_DIR/services/assets.sh"
source "$COMMON_DIR/services/nats.sh"
source "$COMMON_DIR/cgroup/cgroup.sh"

# ─── Configuration ───────────────────────────────────────────
EVENT_COUNT="${1:-1000000}"
CPU_LIST="${2:-16 8 4 2 1}"
BENCH_TAG="${3:-baseline}"
EXTRA_ENV="${4:-}"   # Extra env vars for the service (e.g. "CONSUMER_PREFETCH=true")
POLL_INTERVAL=2
SUDO_PASS="${SUDO_PASS:-4670}"

# Paths
SERVICE_DIR="$(cd "$SCRIPT_DIR/../../.." && pwd)"
RESULTS_DIR="$SCRIPT_DIR/../results"

# ─── Local Helper Functions ────────────────────────────────

wait_for_stream_populate() {
    local target="$1"
    local cpu="$2"
    local actual
    actual=$(get_msg_count "$STREAM")
    actual="${actual:-0}"

    echo ""
    echo "  ╔═══════════════════════════════════════════════════════════════╗"
    echo "  ║  AGUARDANDO POPULACAO DO STREAM                             ║"
    echo "  ║                                                             ║"
    echo "  ║  Teste:    CPU_LIMIT=$cpu ($(( cpu > 1 ? cpu - 1 : 0 )) batch workers)"
    echo "  ║  Stream:   $STREAM"
    echo "  ║  Atual:    $actual mensagens"
    echo "  ║  Esperado: $target mensagens"
    echo "  ║                                                             ║"
    echo "  ║  Polling a cada 5s ate detectar mensagens no stream...      ║"
    echo "  ╚═══════════════════════════════════════════════════════════════╝"
    echo ""

    # Auto-detect: wait for >= 90% of target AND stable (not growing for 15s)
    local threshold=$(( target * 90 / 100 ))
    local prev=0
    local stable_count=0
    while true; do
        actual=$(get_msg_count "$STREAM")
        actual="${actual:-0}"
        if [ "$actual" -ge "$threshold" ]; then
            if [ "$actual" -eq "$prev" ]; then
                stable_count=$((stable_count + 1))
                if [ "$stable_count" -ge 3 ]; then
                    break
                fi
            else
                stable_count=0
            fi
        fi
        prev="$actual"
        printf "\r      Aguardando... %s / %s msgs (%d%%)   " "$actual" "$target" "$(( actual * 100 / target ))"
        sleep 5
    done

    echo ""
    echo "      Stream pronto e estavel: $actual mensagens. Iniciando teste..."
}

wait_for_drain() {
    local total="$1"

    # Wait for consumer to start (count starts decreasing)
    echo "      Waiting for consumer to start draining..."
    local prev="$total"
    local wait_count=0
    while true; do
        sleep "$POLL_INTERVAL"
        local count
        count=$(get_msg_count "$STREAM")
        [ -z "$count" ] && continue
        wait_count=$((wait_count + 1))
        printf "\r      Waiting for drain to start... %s msgs remaining (%ds)   " "$count" "$((wait_count * POLL_INTERVAL))"
        if [ "$count" -lt "$prev" ]; then
            echo ""
            echo "      Consumer started draining!"
            break
        fi
        prev="$count"
    done

    local start_time
    start_time=$(date +%s%N)

    while true; do
        sleep "$POLL_INTERVAL"
        local count
        count=$(get_msg_count "$STREAM")
        [ -z "$count" ] && continue

        local processed=$((total - count))
        local pct=$((processed * 100 / total))
        local now_s=$(date +%s%N)
        local elapsed_s=$(( (now_s - start_time) / 1000000000 ))
        local rate=0
        [ "$elapsed_s" -gt 0 ] && rate=$((processed / elapsed_s))
        printf "\r      [%dm%02ds] %d / %d (%d%%) | %d ev/s | remaining: %d   " \
            $((elapsed_s/60)) $((elapsed_s%60)) "$processed" "$total" "$pct" "$rate" "$count"

        if [ "$count" -eq 0 ]; then
            echo ""
            break
        fi
    done

    local end_time
    end_time=$(date +%s%N)
    local elapsed_ns=$((end_time - start_time))
    DRAIN_SECONDS=$((elapsed_ns / 1000000000))
    DRAIN_MS=$((elapsed_ns / 1000000))
    if [ "$DRAIN_SECONDS" -gt 0 ]; then
        DRAIN_RATE=$((total / DRAIN_SECONDS))
    else
        DRAIN_RATE="$total"
    fi
}

# ─── Kill Competing Services ───────────────────────────────
kill_competing_services

# ─── Preflight Checks ───────────────────────────────────────

# 1. Infrastructure health checks
run_preflight_checks "" "$BENCH_CLI_TOOLS" "$BENCH_SERVICE_CHECKS"

# 2. Verify cgroup shield exists
verify_shield || exit 1

# 3. Validate CPU list doesn't exceed shield capacity
validate_cpu_list "$CPU_LIST"

# 4. Validate SERVICE_DIR exists
if [ ! -f "$SERVICE_DIR/src/main.ts" ]; then
    log_error "Service not found at $SERVICE_DIR/src/main.ts"
    echo "Expected: workspace_js/services/js-executor/src/main.ts"
    exit 1
fi
echo "Service dir:     $SERVICE_DIR"

# ─── Seed Assets (once) ─────────────────────────────────────

echo ""
log_info "Seeding benchmark assets (MongoDB + MinIO)..."
bench_seed_assets_mongodb
bench_seed_assets_minio
echo ""

# ─── Main ────────────────────────────────────────────────────

RESULTS_DIR="$RESULTS_DIR/$BENCH_TAG"
mkdir -p "$RESULTS_DIR"
declare -a SUMMARY_ROWS=()

echo ""
echo "================================================================"
echo "  JS-Executor Piscina Benchmark (cgroup v2 isolated)"
echo "================================================================"
echo "  Architecture:    Piscina Worker Threads"
echo "  Events per test: $(format_number "$EVENT_COUNT")"
echo "  CPU list:        $CPU_LIST"
echo "  Tag:             $BENCH_TAG"
[ -n "$EXTRA_ENV" ] && echo "  Extra env:       $EXTRA_ENV"
echo "  Isolation:       cgroup v2 cpuset (cores $SHIELD_CORE_START-$SHIELD_CORE_END)"
echo "  Results dir:     $RESULTS_DIR"
echo "================================================================"
echo ""

TEST_NUM=0
TOTAL_TESTS=$(echo "$CPU_LIST" | wc -w)

for CPU in $CPU_LIST; do
    TEST_NUM=$((TEST_NUM + 1))

    # Piscina auto-tuning (matches consumer.constant.ts)
    if [ "$CPU" -eq 1 ]; then
        WORKERS=0
        BATCH_SIZE=500
    else
        WORKERS=$((CPU - 1))
        BATCH_SIZE=$((CPU * 500))
    fi
    EVENTS_PER_WORKER=500

    # Set shield to use exactly N cores
    CORE_RANGE=$(set_shield_cpus "$CPU")

    echo ""
    echo "╔══════════════════════════════════════════════════════════════════╗"
    echo "║  TEST $TEST_NUM/$TOTAL_TESTS — CPU_LIMIT=$CPU  workers=$WORKERS  batch=$BATCH_SIZE"
    echo "║  Isolated cores: $CORE_RANGE (cgroup v2 shield)"
    echo "╚══════════════════════════════════════════════════════════════════╝"
    echo ""

    # 0. Kill any leftover js-executor
    kill_service_on_port "$SERVICE_PORT"

    # 1. Purge ALL NATS streams
    log_step "1" "6" "Cleaning ALL NATS streams..."
    bench_purge_all_streams
    echo "      Done."

    # 2. Populate the stream
    log_step "2" "6" "Populating stream..."
    bench_remove_consumer
    bench_seed_messages "$EVENT_COUNT" "$SEED_DIR/nats/js-execute-http.json"
    bench_verify_stream "$EVENT_COUNT"
    wait_for_stream_populate "$EVENT_COUNT" "$CPU"

    ACTUAL_COUNT=$(get_msg_count "$STREAM")
    echo "      Stream has $ACTUAL_COUNT messages."

    # 3. Start js-executor in isolated cgroup
    log_step "3" "6" "Starting js-executor (CPU_LIMIT=$CPU, workers=$WORKERS, cores $CORE_RANGE)..."
    SERVICE_LOG="$RESULTS_DIR/test-cpu${CPU}-output.log"
    cd "$SERVICE_DIR"
    env CPU_LIMIT="$CPU" LOG_LEVEL=silent $EXTRA_ENV node -r ts-node/register -r tsconfig-paths/register src/main.ts \
        > "$SERVICE_LOG" 2>&1 &
    JS_PID=$!

    # Verify process is alive
    sleep 1
    if ! kill -0 "$JS_PID" 2>/dev/null; then
        log_error "Service died immediately. Last 20 lines:"
        tail -20 "$SERVICE_LOG" 2>/dev/null || true
        continue
    fi

    # Move process into the isolated cgroup
    move_to_shield "$JS_PID" || log_warn "Failed to move PID $JS_PID to cgroup (check SUDO_PASS)"
    # Also move child processes (Piscina workers)
    sleep 2
    for child_pid in $(pgrep -P "$JS_PID" 2>/dev/null); do
        move_to_shield "$child_pid" 2>/dev/null || true
    done
    echo "      PID: $JS_PID → cgroup benchmark (cores $CORE_RANGE)"

    # Verify isolation
    ACTUAL_CPUS=$(grep Cpus_allowed_list /proc/$JS_PID/status 2>/dev/null | awk '{print $2}') || true
    echo "      Verified: /proc/$JS_PID/status → Cpus_allowed_list: ${ACTUAL_CPUS:-N/A}"

    # Wait for service to be ready (metrics endpoint)
    echo "      Waiting for service to start..."
    if ! wait_for_service_ready "$METRICS_URL" 120 "$JS_PID"; then
        echo "      SKIPPING test CPU=$CPU (service failed to start)."
        kill "$JS_PID" 2>/dev/null || true
        wait "$JS_PID" 2>/dev/null || true
        continue
    fi

    # 4. Wait for drain
    log_step "4" "6" "Monitoring drain..."
    DRAIN_START_COUNT="${ACTUAL_COUNT:-$EVENT_COUNT}"
    wait_for_drain "$DRAIN_START_COUNT"
    echo "      Drain complete: ${DRAIN_SECONDS}s ($((DRAIN_SECONDS/60))m$((DRAIN_SECONDS%60))s) → ${DRAIN_RATE} ev/s"

    # 5. Capture metrics (wait a few seconds for final metrics flush)
    log_step "5" "6" "Capturing Prometheus metrics..."
    sleep 5
    METRICS=$(curl -s "$METRICS_URL" 2>/dev/null)
    echo "$METRICS" > "$RESULTS_DIR/test-cpu${CPU}-metrics.txt"
    echo "      Saved: test-cpu${CPU}-metrics.txt ($(echo "$METRICS" | wc -l) lines)"

    # ── Parse metrics ──

    # Latency
    LATENCY_AVG=$(extract_histogram_avg "$METRICS" "jsexec_event_duration_seconds" 'consumer="js_execute"')

    # Events
    EVENTS_SUCCESS=$(extract_metric "$METRICS" "jsexec_events_processed_total" 'status="success"')
    EVENTS_FAILED=$(extract_metric "$METRICS" "jsexec_events_processed_total" 'status="error"')
    SCRIPT_ERRORS=$(extract_metric "$METRICS" "jsexec_script_errors_total")
    OOM_COUNT=$(extract_metric "$METRICS" "jsexec_pool_recycled_total" 'reason="oom"')

    # Memory
    RSS_BYTES=$(extract_metric "$METRICS" "process_resident_memory_bytes")
    HEAP_USED_BYTES=$(extract_metric "$METRICS" "nodejs_heap_size_used_bytes")
    HEAP_TOTAL_BYTES=$(extract_metric "$METRICS" "nodejs_heap_size_total_bytes")
    EXTERNAL_BYTES=$(extract_metric "$METRICS" "nodejs_external_memory_bytes")

    # Event loop
    EL_P99=$(extract_metric "$METRICS" "nodejs_eventloop_lag_p99_seconds")

    # GC
    GC_MINOR=$(extract_metric "$METRICS" "nodejs_gc_duration_seconds_count" 'kind="minor"')
    GC_MAJOR=$(extract_metric "$METRICS" "nodejs_gc_duration_seconds_count" 'kind="major"')

    # Derived values
    RSS_MB=$(echo "${RSS_BYTES:-0}" | awk '{printf "%.0f", $1/1048576}')
    HEAP_USED_MB=$(echo "${HEAP_USED_BYTES:-0}" | awk '{printf "%.0f", $1/1048576}')
    HEAP_TOTAL_MB=$(echo "${HEAP_TOTAL_BYTES:-0}" | awk '{printf "%.0f", $1/1048576}')
    EXTERNAL_MB=$(echo "${EXTERNAL_BYTES:-0}" | awk '{printf "%.1f", $1/1048576}')
    EL_P99_MS=$(echo "${EL_P99:-0}" | awk '{printf "%.1f", $1*1000}')

    # Latency percentiles from histogram buckets
    EVT_COUNT=$(echo "$METRICS" | grep '^jsexec_event_duration_seconds_count{.*js_execute' 2>/dev/null | head -1 | awk '{print $NF}') || true
    EVT_COUNT="${EVT_COUNT:-0}"
    HALF_COUNT=$(echo "$EVT_COUNT" | awk '{printf "%.0f", $1*0.5}')
    P50_BUCKET=$(echo "$METRICS" | grep '^jsexec_event_duration_seconds_bucket{.*js_execute' 2>/dev/null | \
        awk -v target="$HALF_COUNT" '{val=$NF; if(val+0 >= target+0 && !found){gsub(/.*le="/,""); gsub(/".*/,""); printf "%.1f", $0*1000; found=1}}') || true
    P99_TARGET=$(echo "$EVT_COUNT" | awk '{printf "%.0f", $1*0.99}')
    P99_BUCKET=$(echo "$METRICS" | grep '^jsexec_event_duration_seconds_bucket{.*js_execute' 2>/dev/null | \
        awk -v target="$P99_TARGET" '{val=$NF; if(val+0 >= target+0 && !found){gsub(/.*le="/,""); gsub(/".*/,""); printf "%.1f", $0*1000; found=1}}') || true

    # 6. Stop js-executor
    log_step "6" "6" "Stopping js-executor..."
    kill "$JS_PID" 2>/dev/null || true
    wait "$JS_PID" 2>/dev/null || true
    sleep 2
    # Kill any orphan workers
    kill_service_on_port "$SERVICE_PORT"
    echo "      Stopped."

    # Collect summary row
    SUMMARY_ROWS+=("$(printf "  │ %3s │ %7s │ %5s │ %7s/s │ %5ss │ %6sMB │ %6sMB │ %5sms │" \
        "$CPU" "$WORKERS" "$BATCH_SIZE" "$DRAIN_RATE" "$DRAIN_SECONDS" "${RSS_MB:-0}" "${HEAP_USED_MB:-0}" "${EL_P99_MS:-N/A}")")

    # Print summary
    echo ""
    echo "  ┌──────────────────────────────────────────────────────────────┐"
    echo "  │ CPU=$CPU (${WORKERS}w) → ${DRAIN_RATE} ev/s (${DRAIN_SECONDS}s) [cores $CORE_RANGE]"
    echo "  │ Latency: avg=${LATENCY_AVG:-N/A}ms p50=${P50_BUCKET:-N/A}ms p99=${P99_BUCKET:-N/A}ms"
    echo "  │ Memory:  RSS=${RSS_MB}MB | Heap=${HEAP_USED_MB}/${HEAP_TOTAL_MB}MB | Ext=${EXTERNAL_MB}MB"
    echo "  │ GC: minor=${GC_MINOR:-0} major=${GC_MAJOR:-0} | OOMs=${OOM_COUNT:-0} | EL_p99=${EL_P99_MS}ms"
    echo "  └──────────────────────────────────────────────────────────────┘"

    sleep 5  # Cool down between tests
done

# ─── Teardown Assets ────────────────────────────────────────

echo ""
log_info "Cleaning up benchmark assets..."
bench_cleanup_assets

# ─── Final Summary ───────────────────────────────────────────

echo ""
echo ""
echo "════════════════════════════════════════════════════════════════════"
echo "  BENCHMARK COMPLETE — $TOTAL_TESTS tests (Piscina, cgroup v2)"
echo "════════════════════════════════════════════════════════════════════"
echo ""
echo "  Results: $RESULTS_DIR/"
echo ""
echo "  Summary:"
echo "  ┌─────┬─────────┬───────┬────────────┬──────────┬──────────┬──────────┬──────────┐"
printf "  │ %3s │ %7s │ %5s │ %10s │ %8s │ %8s │ %8s │ %8s │\n" \
    "CPU" "Workers" "Batch" "Throughput" "Drain" "RSS(MB)" "Heap(MB)" "EL_p99"
echo "  ├─────┼─────────┼───────┼────────────┼──────────┼──────────┼──────────┼──────────┤"

for ROW in "${SUMMARY_ROWS[@]}"; do
    echo "$ROW"
done

echo "  └─────┴─────────┴───────┴────────────┴──────────┴──────────┴──────────┴──────────┘"
echo ""
echo "  Per-test files:"
echo "    Metrics: $RESULTS_DIR/test-cpu<N>-metrics.txt  (full Prometheus dump)"
echo "    Logs:    $RESULTS_DIR/test-cpu<N>-output.log   (service output)"
echo ""
