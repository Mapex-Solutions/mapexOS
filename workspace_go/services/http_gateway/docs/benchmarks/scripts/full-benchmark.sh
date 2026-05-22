#!/bin/bash
# =============================================================
# HTTP Gateway Full CPU Benchmark (cgroup v2 isolated)
#
# Self-contained: collects ALL metrics from /metrics endpoint.
# No Grafana or Prometheus reset needed.
#
# Runs ALL auth types (jwt, apiKey, ipWhiteList, none) by default.
# Override via AUTH_TYPES env var.
#
# Flow:
#   1. Preflight checks (infra health + user confirmation)
#   2. Build Go binary
#   3. Seed MongoDB (4 DataSources) — once before all tests
#   4. For each auth type:
#      For each CPU config:
#        a. Kill old process → purge NATS
#        b. Start service in isolated cgroup
#        c. HTTP flood (hey -n N -c C)
#        d. Capture metrics (Prometheus + hey + Go runtime)
#        e. Stop → cool down
#   5. Teardown via seed.sh
#
# Prerequisites:
#   1. cgroup shield: sudo bash scripts/benchmarks/common/cgroup/setup-cgroup-shield.sh
#   2. MongoDB, NATS, Redis running
#   3. hey installed: go install github.com/rakyll/hey@latest
#
# Usage:
#   ./full-benchmark.sh [request_count] [cpu_list] [tag] [concurrency]
#
# Examples:
#   ./full-benchmark.sh                              # all auth, 1M, CPU 16 8 4 2 1
#   ./full-benchmark.sh 100000 "4 2 1"               # all auth, 100K, CPU 4 2 1
#   ./full-benchmark.sh 1000 "1" smoke-test           # all auth, 1K, 1 CPU
#   AUTH_TYPES="jwt none" ./full-benchmark.sh 500000   # only jwt + none
# =============================================================

set -euo pipefail

# ─── Source Modules ──────────────────────────────────────────

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"

source "$SCRIPT_DIR/config.sh"
source "$COMMON_DIR/services/datasources.sh"
source "$COMMON_DIR/services/nats.sh"
source "$COMMON_DIR/cgroup/cgroup.sh"

# ─── Configuration ───────────────────────────────────────────
REQUEST_COUNT="${1:-1000000}"
CPU_LIST="${2:-16 8 4 2 1}"
BENCH_TAG="${3:-baseline}"
CONCURRENCY="${4:-200}"
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

# Auth types — run all by default, override via env var
AUTH_TYPES="${AUTH_TYPES:-jwt apiKey ipWhiteList none}"

# Validate auth types
for _at in $AUTH_TYPES; do
    case "$_at" in
        jwt|apiKey|ipWhiteList|none) ;;
        *) log_error "Invalid auth type '$_at'. Use: jwt | apiKey | ipWhiteList | none"; exit 1 ;;
    esac
done

# Service
TARGET_URL="http://localhost:${SERVICE_PORT}/api/v1/events"

# Paths
SERVICE_DIR="$(cd "$SCRIPT_DIR/../../.." && pwd)"
RESULTS_DIR="$SCRIPT_DIR/../results/$BENCH_TAG"
BINARY_PATH="$SCRIPT_DIR/http_gateway_bench"
PAYLOAD_FILE="$SERVICE_DIR/docs/benchmarks/seed/http/event-iot-network-status.json"

# ─── Local Helper Functions ──────────────────────────────────

run_hey() {
    local auth="$1"
    local count="$2"
    local conc="$3"
    local url="$4"
    local output_file="${5:-/dev/null}"

    local -a hey_args=(-n "$count" -c "$conc" -m POST -H "Content-Type: application/json" -D "$PAYLOAD_FILE")

    case "$auth" in
        jwt)         hey_args+=(-H "x-dt-signature: ${JWT_TOKEN_STATIC}") ;;
        apiKey)      hey_args+=(-H "x-api-key: ${API_KEY_VALUE}") ;;
    esac

    hey "${hey_args[@]}" "$url" > "$output_file" 2>&1
}

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

    HEY_LAT_P50=$(echo "$output" | grep "50%% in" | head -1 | awk '{print $3}')
    HEY_LAT_P50="${HEY_LAT_P50:-0}"

    HEY_LAT_P95=$(echo "$output" | grep "95%% in" | head -1 | awk '{print $3}')
    HEY_LAT_P95="${HEY_LAT_P95:-0}"

    HEY_LAT_P99=$(echo "$output" | grep "99%% in" | head -1 | awk '{print $3}')
    HEY_LAT_P99="${HEY_LAT_P99:-0}"

    # Convert secs to ms
    HEY_LAT_AVG_MS=$(awk "BEGIN{printf \"%.2f\", ${HEY_LAT_AVG}*1000}")
    HEY_LAT_P50_MS=$(awk "BEGIN{printf \"%.2f\", ${HEY_LAT_P50}*1000}")
    HEY_LAT_P95_MS=$(awk "BEGIN{printf \"%.2f\", ${HEY_LAT_P95}*1000}")
    HEY_LAT_P99_MS=$(awk "BEGIN{printf \"%.2f\", ${HEY_LAT_P99}*1000}")

    # Status codes (grep may find 0 matches → exit 1 → use { ... || true; } to survive pipefail)
    HEY_2XX=$({ echo "$output" | grep -E "^\s*\[2[0-9]{2}\]" || true; } | awk '{sum+=$2} END {print sum+0}')
    HEY_4XX=$({ echo "$output" | grep -E "^\s*\[4[0-9]{2}\]" || true; } | awk '{sum+=$2} END {print sum+0}')
    HEY_5XX=$({ echo "$output" | grep -E "^\s*\[5[0-9]{2}\]" || true; } | awk '{sum+=$2} END {print sum+0}')
    HEY_ERRORS=$({ echo "$output" | grep -c "Error distribution:" || true; } | head -1)
}

# ─── Kill Competing Services ───────────────────────────────
kill_competing_services

# ─── Preflight Checks ───────────────────────────────────────

# Infrastructure health checks
run_preflight_checks "" "$BENCH_CLI_TOOLS" "$BENCH_SERVICE_CHECKS"

# ─── Build Binary ────────────────────────────────────────────

echo ""
log_info "Compiling http_gateway..."
cd "$SERVICE_DIR"
GOWORK=off CGO_ENABLED=0 go build -o "$BINARY_PATH" ./src/ 2>&1
log_success "Binary: $BINARY_PATH"

# ─── Seed MongoDB ────────────────────────────────────────────

echo ""
log_info "Seeding benchmark DataSources via seed.sh..."
bash "$SCRIPT_DIR/seed.sh" setup

# Validate payload file
if [ ! -f "$PAYLOAD_FILE" ]; then
    log_error "Payload file not found: $PAYLOAD_FILE"
    exit 1
fi
PAYLOAD_SIZE=$(wc -c < "$PAYLOAD_FILE")
log_success "Payload: ${PAYLOAD_SIZE} bytes ($(basename "$PAYLOAD_FILE"))"

# Kill stale processes
kill_service_on_port "$SERVICE_PORT"

# ─── Main Benchmark Loop ─────────────────────────────────────

mkdir -p "$RESULTS_DIR"

# Count total tests: AUTH_TYPES × CPU_LIST
NUM_AUTH=$(echo "$AUTH_TYPES" | wc -w)
NUM_CPU=$(echo "$CPU_LIST" | wc -w)
TOTAL_TESTS=$((NUM_AUTH * NUM_CPU))

echo ""
echo "================================================================"
echo "  HTTP Gateway CPU Benchmark (cgroup v2)"
echo "================================================================"
echo "  Auth types:     $AUTH_TYPES"
echo "  Requests/test:  $(format_number $REQUEST_COUNT)"
echo "  Concurrency:    $CONCURRENCY"
echo "  CPU list:       $CPU_LIST"
echo "  Total tests:    $TOTAL_TESTS (${NUM_AUTH} auth x ${NUM_CPU} CPU)"
echo "  Tag:            $BENCH_TAG"
echo "  MongoDB:        $MONGO_DB"
echo "  Results:        $RESULTS_DIR"
echo "================================================================"
echo ""

declare -a SUMMARY_ROWS=()
TEST_NUM=0

for AUTH_TYPE in $AUTH_TYPES; do
    # Resolve DataSource ID for this auth type
    DS_ID=$(bench_get_ds_id "$AUTH_TYPE")
    FULL_URL="${TARGET_URL}?ds=${DS_ID}"
    log_info "Auth block: $AUTH_TYPE → DS: $DS_ID"

    for CPU in $CPU_LIST; do
        TEST_NUM=$((TEST_NUM + 1))
        CORE_RANGE=$(set_shield_cpus "$CPU")

        echo ""
        echo "════════════════════════════════════════════════════════════════"
        echo "  TEST $TEST_NUM/$TOTAL_TESTS — auth=$AUTH_TYPE  GOMAXPROCS=$CPU  cores=$CORE_RANGE"
        echo "  $(format_number $REQUEST_COUNT) requests @ ${CONCURRENCY} concurrency"
        echo "════════════════════════════════════════════════════════════════"
        echo ""

        # 0. Kill any leftover process
        kill_service_on_port "$SERVICE_PORT"

        # 1. Purge NATS
        log_step "1" "5" "Purging NATS streams..."
        bench_purge_all_streams

        # 2. Start service (GO_ENV=dev for correct DB, LOG_LEVEL=silent for zero noise)
        log_step "2" "5" "Starting http_gateway (GOMAXPROCS=$CPU, LOG_LEVEL=silent)..."
        GO_ENV="${GO_ENV_VALUE}" LOG_LEVEL=silent GOMAXPROCS="$CPU" "$BINARY_PATH" \
            > /dev/null 2>&1 &
        GW_PID=$!
        sleep 1
        move_to_shield "$GW_PID" || log_warn "Failed to move PID $GW_PID to cgroup"
        echo "      PID=$GW_PID → cores $CORE_RANGE"

        ACTUAL_CPUS=$(grep Cpus_allowed_list /proc/$GW_PID/status 2>/dev/null | awk '{print $2}') || true
        echo "      Verified: $ACTUAL_CPUS"

        # 3. Wait for ready
        log_step "3" "5" "Waiting for service..."
        if ! wait_for_service_ready "$METRICS_URL" 60 "$GW_PID"; then
            echo "      SKIP auth=$AUTH_TYPE CPU=$CPU"
            kill "$GW_PID" 2>/dev/null || true; wait "$GW_PID" 2>/dev/null || true
            continue
        fi

        # 4. Benchmark
        log_step "4" "5" "Running: $(format_number $REQUEST_COUNT) requests..."
        HEY_FILE="$RESULTS_DIR/test-${AUTH_TYPE}-cpu${CPU}-hey.txt"

        BENCH_START=$(date +%s%N)
        run_hey "$AUTH_TYPE" "$REQUEST_COUNT" "$CONCURRENCY" "$FULL_URL" "$HEY_FILE"
        BENCH_END=$(date +%s%N)
        BENCH_SECS=$(( (BENCH_END - BENCH_START) / 1000000000 ))
        echo "      Done in ${BENCH_SECS}s"

        # 5. Collect metrics
        log_step "5" "5" "Collecting metrics..."
        sleep 2
        METRICS_FILE="$RESULTS_DIR/test-${AUTH_TYPE}-cpu${CPU}-metrics.txt"
        METRICS=$(curl -s "$METRICS_URL" 2>/dev/null)
        echo "$METRICS" > "$METRICS_FILE"

        # Parse hey output
        parse_hey_output "$HEY_FILE"

        # Parse service metrics (HTTP Gateway specific)
        AUTH_SUCCESS=$(extract_metric "$METRICS" "httpgw_event_auth_total" "result=\"success\"")
        AUTH_AVG_MS=$(extract_histogram_avg "$METRICS" "httpgw_event_auth_duration_seconds" "auth_type=\"${AUTH_TYPE}\"")
        EVENTS_OK=$(extract_metric "$METRICS" "httpgw_event_processed_total" 'status="success"')
        EVENTS_ERR=$(extract_metric "$METRICS" "httpgw_event_processed_total" 'status="error"')
        NATS_OK=$(extract_metric "$METRICS" "httpgw_event_published_total" 'status="success"')
        NATS_ERR=$(extract_metric "$METRICS" "httpgw_event_published_total" 'status="error"')
        PROC_AVG=$(extract_histogram_avg "$METRICS" "httpgw_event_processing_duration_seconds")

        # Parse Go runtime (via common module)
        parse_go_runtime_metrics "$METRICS"

        # Stop service
        kill "$GW_PID" 2>/dev/null || true
        wait "$GW_PID" 2>/dev/null || true

        # Summary row for final table (with Auth column)
        SUMMARY_ROWS+=("$(printf "  │ %-12s │ %3s │ %8s │ %7s/s │ %5sms │ %5sms │ %5sms │ %5sMB │ %5sMB │ %6s │ %6s │" \
            "$AUTH_TYPE" "$CPU" "$CORE_RANGE" "$HEY_RPS" "$HEY_LAT_AVG_MS" "$HEY_LAT_P50_MS" "$HEY_LAT_P99_MS" "$GO_RSS_MB" "$GO_HEAP_MB" "$HEY_2XX" "$HEY_5XX")")

        # Inline summary
        echo ""
        echo "  ┌──────────────────────────────────────────────────────────┐"
        echo "  │ auth=$AUTH_TYPE CPU=$CPU → ${HEY_RPS} req/s  (${HEY_DURATION}s)"
        echo "  │ Lat: avg=${HEY_LAT_AVG_MS}ms p50=${HEY_LAT_P50_MS}ms p99=${HEY_LAT_P99_MS}ms"
        echo "  │ RSS=${GO_RSS_MB}MB  Heap=${GO_HEAP_MB}MB  Goroutines=${GO_GOROUTINES}"
        echo "  │ 2xx=${HEY_2XX} 4xx=${HEY_4XX} 5xx=${HEY_5XX}"
        echo "  └──────────────────────────────────────────────────────────┘"

        sleep 5
    done
done

# ─── Final Summary ───────────────────────────────────────────

echo ""
echo ""
echo "================================================================"
echo "  COMPLETE — ${TOTAL_TESTS} tests | auth types: ${AUTH_TYPES}"
echo "================================================================"
echo ""

# Summary table (with Auth column)
TABLE_HEADER=$(printf "  │ %-12s │ %3s │ %8s │ %10s │ %8s │ %8s │ %8s │ %7s │ %7s │ %6s │ %6s │" \
    "Auth" "CPU" "Cores" "Req/s" "Lat avg" "Lat p50" "Lat p99" "RSS MB" "Heap" "2xx" "5xx")

echo "  Summary:"
echo "  ┌──────────────┬─────┬──────────┬────────────┬──────────┬──────────┬──────────┬─────────┬─────────┬────────┬────────┐"
echo "$TABLE_HEADER"
echo "  ├──────────────┼─────┼──────────┼────────────┼──────────┼──────────┼──────────┼─────────┼─────────┼────────┼────────┤"

for ROW in "${SUMMARY_ROWS[@]}"; do
    echo "$ROW"
done

echo "  └──────────────┴─────┴──────────┴────────────┴──────────┴──────────┴──────────┴─────────┴─────────┴────────┴────────┘"
echo ""

# ─── Cleanup Benchmark Data ─────────────────────────────────
# Delegate to seed.sh teardown (separation of concerns per benchmark standard)

log_info "Removing benchmark data via seed.sh teardown..."
bash "$SCRIPT_DIR/seed.sh" teardown

# Remove temp binary
rm -f "$BINARY_PATH"

log_success "Done!"
echo "  Results stored at: $RESULTS_DIR"
