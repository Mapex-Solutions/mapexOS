#!/bin/bash
# =============================================================
# MapexOS Full CPU Benchmark (cgroup v2 isolated)
#
# Runs ALL auth scenarios across multiple CPU configurations.
# Uses a custom load generator (loadgen) that rotates through
# N different seeded users per request for realistic benchmarks.
#
# Scenarios (auth-focused):
#   auth_login    — POST /api/v1/auth/login (bcrypt verify + JWT generation)
#   auth_coverage — GET  /api/v1/auth/users/me/coverage (org tree + cache)
#
# Flow:
#   1. Preflight checks (infra health + cgroup)
#   2. Build Go binaries (service + loadgen)
#   3. Seed MongoDB (orgs, roles, users, groups, memberships)
#   4. For each scenario x CPU config:
#        a. Start service in isolated cgroup
#        b. Warmup (1K requests)
#        c. Benchmark (N requests @ C concurrency)
#        d. Capture metrics (Prometheus + Go runtime)
#        e. Stop + cool down
#   5. Teardown via seed.sh
#
# Prerequisites:
#   1. cgroup shield: sudo bash scripts/benchmarks/common/cgroup/setup-cgroup-shield.sh
#   2. MongoDB, NATS, Redis running
#   3. mongosh, nats, redis-cli, curl CLIs
# =============================================================

set -euo pipefail

# ─── Source Modules ──────────────────────────────────────────

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"

source "$SCRIPT_DIR/config.sh"
source "$COMMON_DIR/services/mapexos.sh"
source "$COMMON_DIR/services/nats.sh"
source "$COMMON_DIR/cgroup/cgroup.sh"

# ─── Help ───────────────────────────────────────────────────

show_help() {
    cat <<'HELP'
MapexOS Full CPU Benchmark

Usage:
  ./full-benchmark.sh [request_count] [cpu_list] [concurrency]

Arguments:
  request_count   Total requests per test           (default: 100000)
  cpu_list        CPU configs to test, quoted        (default: "1 2 4 8 16")
  concurrency     Concurrent workers                 (default: 200)

Environment variables:
  SCENARIOS       Override which scenarios to run     (default: all)
                  Values: auth_login auth_coverage
  GO_ENV          Service environment                 (default: dev)
  USER_COUNT      Seeded users for rotation           (default: 1000)
  ORG_COUNT       Seeded organizations                (default: 10)
  GROUP_COUNT     Seeded groups                       (default: 100)
  ROLE_COUNT      Seeded roles                        (default: 50)
  MEMBERSHIP_COUNT  Seeded memberships                (default: 500)
  MONGO_DB        MongoDB database name               (default: ${GO_ENV}-mapexos)
  JWT_SECRET      JWT signing secret                  (default: bench-secret-...)

Examples:
  ./full-benchmark.sh                                 # all, 100K, CPU 1-16
  ./full-benchmark.sh 1000000 "16 12 8 4 2 1"        # all, 1M, CPUs desc
  ./full-benchmark.sh 1000000 "16 12 8 4 2 1" 100    # all, 1M, concurrency 100
  ./full-benchmark.sh 10000 "1"                       # all, 10K, 1 CPU (smoke test)
  SCENARIOS="auth_login" ./full-benchmark.sh 50000 "4 8"  # only login
HELP
    exit 0
}

case "${1:-}" in -h|--help|help) show_help ;; esac

# ─── Configuration ───────────────────────────────────────────

REQUEST_COUNT="${1:-100000}"
CPU_LIST="${2:-1 2 4 8 16}"
CONCURRENCY="${3:-200}"

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

# Scenarios — run all by default, override via env var
SCENARIOS="${SCENARIOS:-auth_login auth_coverage}"

# Validate scenarios
for _sc in $SCENARIOS; do
    case "$_sc" in
        auth_login|auth_coverage) ;;
        *) log_error "Invalid scenario '$_sc'. Use: auth_login | auth_coverage"; exit 1 ;;
    esac
done

# Paths
SERVICE_DIR="$(cd "$SCRIPT_DIR/../../.." && pwd)"
RESULTS_DIR="$SCRIPT_DIR/../results/$BENCH_TAG"
BINARY_PATH="$SCRIPT_DIR/mapexos_bench"
LOADGEN_PATH="$RESULTS_DIR/loadgen"

# ─── Local Helper Functions ──────────────────────────────────

# get_scenario_config sets TARGET_URL and LOADGEN_MODE for a scenario.
get_scenario_config() {
    local scenario="$1"
    local base_url="http://localhost:${SERVICE_PORT}"

    case "$scenario" in
        auth_login)
            TARGET_URL="${base_url}/auth/login"
            LOADGEN_MODE="login"
            ;;
        auth_coverage)
            TARGET_URL="${base_url}/auth/users/me/coverage"
            LOGIN_URL="${base_url}/auth/login"
            LOADGEN_MODE="coverage"
            ;;
    esac
}

# parse_hey_output extracts metrics from loadgen output (hey-compatible format).
parse_hey_output() {
    local file="$1"
    local output
    output=$(<"$file")

    HEY_RPS=$(echo "$output" | grep "Requests/sec:" | awk '{print $2}')
    HEY_RPS="${HEY_RPS:-0}"

    HEY_DURATION=$(echo "$output" | grep "Total:" | head -1 | awk '{print $2}')
    HEY_DURATION="${HEY_DURATION:-0}"

    HEY_LAT_AVG=$(echo "$output" | grep "Average:" | head -1 | awk '{print $2}')
    HEY_LAT_AVG="${HEY_LAT_AVG:-0}"

    HEY_LAT_P50=$(echo "$output" | grep "50% in" | head -1 | awk '{print $3}')
    HEY_LAT_P50="${HEY_LAT_P50:-0}"

    HEY_LAT_P95=$(echo "$output" | grep "95% in" | head -1 | awk '{print $3}')
    HEY_LAT_P95="${HEY_LAT_P95:-0}"

    HEY_LAT_P99=$(echo "$output" | grep "99% in" | head -1 | awk '{print $3}')
    HEY_LAT_P99="${HEY_LAT_P99:-0}"

    # Convert secs to ms
    HEY_LAT_AVG_MS=$(awk "BEGIN{printf \"%.2f\", ${HEY_LAT_AVG}*1000}")
    HEY_LAT_P50_MS=$(awk "BEGIN{printf \"%.2f\", ${HEY_LAT_P50}*1000}")
    HEY_LAT_P95_MS=$(awk "BEGIN{printf \"%.2f\", ${HEY_LAT_P95}*1000}")
    HEY_LAT_P99_MS=$(awk "BEGIN{printf \"%.2f\", ${HEY_LAT_P99}*1000}")

    # Status codes
    HEY_2XX=$({ echo "$output" | grep -E "^\s*\[2[0-9]{2}\]" || true; } | awk '{sum+=$2} END {print sum+0}')
    HEY_4XX=$({ echo "$output" | grep -E "^\s*\[4[0-9]{2}\]" || true; } | awk '{sum+=$2} END {print sum+0}')
    HEY_5XX=$({ echo "$output" | grep -E "^\s*\[5[0-9]{2}\]" || true; } | awk '{sum+=$2} END {print sum+0}')
}

# run_loadgen executes the custom load generator for the current scenario.
run_loadgen() {
    local count="$1"
    local concurrency="$2"
    local output_file="$3"

    local extra_args=()
    if [ "$LOADGEN_MODE" = "coverage" ]; then
        extra_args=(-login-url "$LOGIN_URL")
    fi

    "$LOADGEN_PATH" \
        -mode "$LOADGEN_MODE" \
        -url "$TARGET_URL" \
        -n "$count" \
        -c "$concurrency" \
        -users "$BENCH_USER_COUNT" \
        -password "$BENCH_PASSWORD" \
        "${extra_args[@]}" > "$output_file" 2>&1
}

# ─── Kill Competing Services ───────────────────────────────
kill_competing_services

# ─── Preflight Checks ───────────────────────────────────────

# 1. Infrastructure health checks
run_preflight_checks "" "$BENCH_CLI_TOOLS" "$BENCH_SERVICE_CHECKS"

# ─── Build Binaries ──────────────────────────────────────────

mkdir -p "$RESULTS_DIR"

echo ""
log_info "Compiling mapexos service..."
cd "$SERVICE_DIR"
GOWORK=off CGO_ENABLED=0 go build -o "$BINARY_PATH" ./src/ 2>&1
log_success "Binary: $BINARY_PATH"

log_info "Compiling loadgen..."
go build -o "$LOADGEN_PATH" "$SCRIPT_DIR/loadgen.go" 2>&1
log_success "Binary: $LOADGEN_PATH"

# ─── Seed Data ──────────────────────────────────────────────

echo ""
log_info "Seeding benchmark data via seed.sh..."
bash "$SCRIPT_DIR/seed.sh" setup

# Kill stale processes
kill_service_on_port "$SERVICE_PORT"

# ─── Main Benchmark Loop ─────────────────────────────────────

# Count total tests: SCENARIOS x CPU_LIST
NUM_SCENARIOS=$(echo "$SCENARIOS" | wc -w)
NUM_CPU=$(echo "$CPU_LIST" | wc -w)
TOTAL_TESTS=$((NUM_SCENARIOS * NUM_CPU))

echo ""
echo "================================================================"
echo "  MapexOS Auth Benchmark (cgroup v2)"
echo "================================================================"
echo "  Scenarios:      $SCENARIOS"
echo "  Requests/test:  $(format_number $REQUEST_COUNT)"
echo "  Concurrency:    $CONCURRENCY"
echo "  CPU list:       $CPU_LIST"
echo "  Total tests:    $TOTAL_TESTS (${NUM_SCENARIOS} scenarios x ${NUM_CPU} CPU)"
echo "  Seeded users:   $BENCH_USER_COUNT"
echo "  MongoDB:        $MONGO_DB"
echo "  Results:        $RESULTS_DIR"
echo "================================================================"
echo ""

declare -a SUMMARY_ROWS=()
TEST_NUM=0

for SCENARIO in $SCENARIOS; do
    get_scenario_config "$SCENARIO"
    log_info "Scenario block: $SCENARIO (mode=$LOADGEN_MODE $TARGET_URL)"

    for CPU in $CPU_LIST; do
        TEST_NUM=$((TEST_NUM + 1))
        CORE_RANGE=$(set_shield_cpus "$CPU")

        # Adaptive concurrency: for CPU-bound scenarios (bcrypt), scale
        # concurrency with cores to avoid excessive queuing.
        # For I/O-bound scenarios (coverage with cache), use full concurrency.
        if [ "$SCENARIO" = "auth_login" ]; then
            EFFECTIVE_CONCURRENCY=$(( CPU * 25 ))
            [ "$EFFECTIVE_CONCURRENCY" -gt "$CONCURRENCY" ] && EFFECTIVE_CONCURRENCY="$CONCURRENCY"
        else
            EFFECTIVE_CONCURRENCY="$CONCURRENCY"
        fi

        echo ""
        echo "════════════════════════════════════════════════════════════════"
        echo "  TEST $TEST_NUM/$TOTAL_TESTS — scenario=$SCENARIO  GOMAXPROCS=$CPU  cores=$CORE_RANGE"
        echo "  $(format_number $REQUEST_COUNT) requests @ ${EFFECTIVE_CONCURRENCY} concurrency"
        echo "════════════════════════════════════════════════════════════════"
        echo ""

        # 0. Kill any leftover process
        kill_service_on_port "$SERVICE_PORT"

        # 1. Start service
        log_step "1" "6" "Starting mapexos (GOMAXPROCS=$CPU, LOG_LEVEL=silent)..."
        SERVICE_LOG="$RESULTS_DIR/test-${SCENARIO}-cpu${CPU}-output.log"
        env GO_ENV="$GO_ENV_VALUE" LOG_LEVEL=silent GOMAXPROCS="$CPU" \
            METRICS_GO_COLLECTOR=true METRICS_PROCESS_COLLECTOR=true \
            CTX_TIMEOUT=30 \
            "$BINARY_PATH" > "$SERVICE_LOG" 2>&1 &
        SVC_PID=$!
        sleep 2
        move_to_shield "$SVC_PID" || log_warn "Failed to move PID $SVC_PID to cgroup"
        echo "      PID=$SVC_PID -> cores $CORE_RANGE"

        ACTUAL_CPUS=$(grep Cpus_allowed_list /proc/$SVC_PID/status 2>/dev/null | awk '{print $2}') || true
        echo "      Verified: $ACTUAL_CPUS"

        # 2. Wait for ready
        log_step "2" "6" "Waiting for service..."
        if ! wait_for_service_ready "$METRICS_URL" 60 "$SVC_PID"; then
            echo "      SKIP scenario=$SCENARIO CPU=$CPU"
            kill "$SVC_PID" 2>/dev/null || true; wait "$SVC_PID" 2>/dev/null || true
            continue
        fi

        # 3. Warmup
        log_step "3" "6" "Warmup (1000 requests)..."
        run_loadgen 1000 50 /dev/null || true
        sleep 2
        echo "      Done."

        # 4. Benchmark
        log_step "4" "6" "Running: $(format_number $REQUEST_COUNT) requests..."
        HEY_FILE="$RESULTS_DIR/test-${SCENARIO}-cpu${CPU}-hey.txt"

        BENCH_START=$(date +%s%N)
        run_loadgen "$REQUEST_COUNT" "$EFFECTIVE_CONCURRENCY" "$HEY_FILE"
        BENCH_END=$(date +%s%N)
        BENCH_SECS=$(( (BENCH_END - BENCH_START) / 1000000000 ))
        echo "      Done in ${BENCH_SECS}s"

        # 5. Collect metrics
        log_step "5" "6" "Collecting metrics..."
        sleep 2
        METRICS_FILE="$RESULTS_DIR/test-${SCENARIO}-cpu${CPU}-metrics.txt"
        METRICS=$(curl -s "$METRICS_URL" 2>/dev/null)
        echo "$METRICS" > "$METRICS_FILE"

        # Parse hey output
        parse_hey_output "$HEY_FILE"

        # Parse service metrics
        AUTH_OK=$(extract_metric "$METRICS" "mapexos_auth_attempts_total" 'status="success"')
        AUTH_ERR=$(extract_metric "$METRICS" "mapexos_auth_attempts_total" 'status="failure"')
        AUTH_AVG_MS=$(extract_histogram_avg "$METRICS" "mapexos_auth_duration_seconds")
        CACHE_HIT=$(extract_metric "$METRICS" "mapexos_cache_total" 'result="hit"')
        CACHE_MISS=$(extract_metric "$METRICS" "mapexos_cache_total" 'result="miss"')

        # Parse Go runtime (via common module)
        parse_go_runtime_metrics "$METRICS"

        # 6. Stop service
        log_step "6" "6" "Stopping service..."
        kill "$SVC_PID" 2>/dev/null || true
        wait "$SVC_PID" 2>/dev/null || true
        echo "      Done."

        # Summary row for final table
        SUMMARY_ROWS+=("$(printf "  │ %-14s │ %3s │ %8s │ %7s/s │ %5sms │ %5sms │ %5sms │ %5sMB │ %6s │ %6s │" \
            "$SCENARIO" "$CPU" "$CORE_RANGE" "$HEY_RPS" "$HEY_LAT_AVG_MS" "$HEY_LAT_P50_MS" "$HEY_LAT_P99_MS" "$GO_RSS_MB" "$HEY_2XX" "$HEY_5XX")")

        # Inline summary
        echo ""
        echo "  ┌──────────────────────────────────────────────────────────┐"
        echo "  │ scenario=$SCENARIO CPU=$CPU -> ${HEY_RPS} req/s (${HEY_DURATION}s)"
        echo "  │ Lat: avg=${HEY_LAT_AVG_MS}ms p50=${HEY_LAT_P50_MS}ms p99=${HEY_LAT_P99_MS}ms"
        echo "  │ RSS=${GO_RSS_MB}MB  Heap=${GO_HEAP_MB}MB  Goroutines=${GO_GOROUTINES}"
        echo "  │ Auth OK=${AUTH_OK} Err=${AUTH_ERR}  Cache hit=${CACHE_HIT} miss=${CACHE_MISS}"
        echo "  │ 2xx=${HEY_2XX} 4xx=${HEY_4XX} 5xx=${HEY_5XX}"
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

# Summary table
TABLE_HEADER=$(printf "  │ %-14s │ %3s │ %8s │ %10s │ %8s │ %8s │ %8s │ %7s │ %6s │ %6s │" \
    "Scenario" "CPU" "Cores" "Req/s" "Lat avg" "Lat p50" "Lat p99" "RSS MB" "2xx" "5xx")

echo "  Summary:"
echo "  ┌────────────────┬─────┬──────────┬────────────┬──────────┬──────────┬──────────┬─────────┬────────┬────────┐"
echo "$TABLE_HEADER"
echo "  ├────────────────┼─────┼──────────┼────────────┼──────────┼──────────┼──────────┼─────────┼────────┼────────┤"

for ROW in "${SUMMARY_ROWS[@]}"; do
    echo "$ROW"
done

echo "  └────────────────┴─────┴──────────┴────────────┴──────────┴──────────┴──────────┴─────────┴────────┴────────┘"
echo ""

# ─── Cleanup Benchmark Data ─────────────────────────────────
# Delegate to seed.sh teardown (separation of concerns per benchmark standard)

log_info "Removing benchmark data via seed.sh teardown..."
bash "$SCRIPT_DIR/seed.sh" teardown

# Remove temp binaries
rm -f "$BINARY_PATH"
rm -f "$LOADGEN_PATH"

log_success "Done!"
echo "  Results stored at: $RESULTS_DIR"
