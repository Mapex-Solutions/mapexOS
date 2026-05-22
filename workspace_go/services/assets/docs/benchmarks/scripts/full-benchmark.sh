#!/bin/bash
# =============================================================
# Assets Service Full CPU Benchmark (cgroup v2 isolated)
#
# Runs ALL scenarios (HTTP + Auth Callout) across multiple CPU
# configurations (1, 2, 4, 8, 16 cores).
#
# HTTP scenarios use hey (HTTP flood pattern).
# Auth scenarios use auth-bench Go tool (MQTT connect).
#
# Flow:
#   1. Preflight checks (infra health + cgroup)
#   2. Build Go binaries (service + auth-bench)
#   3. Seed MongoDB (template + 1000 HTTP + 10K MQTT assets)
#   4. For each scenario x CPU config:
#        a. Start service in isolated cgroup
#        b. Warmup
#        c. Benchmark (N requests @ C concurrency)
#        d. Capture metrics (Prometheus + Go runtime)
#        e. Stop + cool down
#   5. Teardown via seed.sh
#
# Prerequisites:
#   1. cgroup shield: sudo bash scripts/benchmarks/common/cgroup/setup-cgroup-shield.sh
#   2. MongoDB, Redis, NATS, MinIO running
#   3. mongosh, redis-cli, nats, mc, curl, hey, openssl CLIs
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
Assets Service Full CPU Benchmark

Usage:
  ./full-benchmark.sh [request_count] [cpu_list] [concurrency]

Arguments:
  request_count   Total requests per test              (default: 100000)
  cpu_list        CPU configs to test, quoted           (default: "1 2 4 8 16")
  concurrency     Concurrent workers                    (default: 200)

Environment variables:
  SCENARIOS       Override which scenarios to run        (default: all)
                  HTTP:  list-assets get-asset create-asset get-template
                  Auth:  auth-cache-hit auth-cache-miss
  GO_ENV          Service environment                    (default: dev)
  MONGO_DB        MongoDB database name                  (default: ${GO_ENV}-assets)
  JWT_SECRET      JWT signing secret for HTTP auth       (default: a-string-secret-...)
  MQTT_BROKER     MQTT broker address for auth callout    (default: localhost:1883)
  MINIO_ALIAS     MinIO client alias                     (default: local)
  BENCH_MQTT_ASSET_COUNT  MQTT devices to seed (auth callout)  (default: 10000)
  SUDO_PASS       sudo password for cgroup operations    (prompted if not set)

Examples:
  ./full-benchmark.sh                                    # all, 100K, CPU 1-16
  ./full-benchmark.sh 1000000 "16 12 8 4 2 1"           # all, 1M, CPUs desc
  ./full-benchmark.sh 1000000 "16 12 8 4 2 1" 100       # all, 1M, concurrency 100
  ./full-benchmark.sh 10000 "1"                          # all, 10K, 1 CPU (smoke test)
  SCENARIOS="get-asset" ./full-benchmark.sh 50000 "4 8"  # only get-asset
  SCENARIOS="auth-cache-hit auth-cache-miss" ./full-benchmark.sh 10000 "1 2 4"  # only auth
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
SCENARIOS="${SCENARIOS:-list-assets get-asset create-asset get-template auth-cache-hit auth-cache-miss}"

# Validate scenarios
for _sc in $SCENARIOS; do
    case "$_sc" in
        list-assets|get-asset|create-asset|get-template) ;;
        auth-cache-hit|auth-cache-miss) ;;
        *) log_error "Invalid scenario '$_sc'. Use: list-assets | get-asset | create-asset | get-template | auth-cache-hit | auth-cache-miss"; exit 1 ;;
    esac
done

# Detect if auth scenarios are requested
NEEDS_AUTH=false
for _sc in $SCENARIOS; do
    case "$_sc" in auth-cache-hit|auth-cache-miss) NEEDS_AUTH=true ;; esac
done

# Paths
SERVICE_DIR="$(cd "$SCRIPT_DIR/../../.." && pwd)"
RESULTS_DIR="$SCRIPT_DIR/../results/$BENCH_TAG"
BINARY_PATH="$SCRIPT_DIR/assets_bench"
PAYLOADS_DIR="$SEED_DIR/payloads"
AUTH_BENCH_DIR="$SCRIPT_DIR/../tools/auth-bench"
AUTH_BENCH_BIN="$AUTH_BENCH_DIR/auth-bench"

# Service
BASE_URL="http://localhost:${SERVICE_PORT}"
JWT_TOKEN=""

# ─── Local Helper Functions ──────────────────────────────────

generate_jwt() {
    local secret="$1"
    local exp_time
    exp_time=$(($(date +%s) + 86400))

    local header
    header=$(echo -n '{"alg":"HS256","typ":"JWT"}' | openssl base64 -A | tr '+/' '-_' | tr -d '=')

    local payload
    payload=$(echo -n "{\"userId\":\"${BENCH_USER_ID}\",\"sub\":\"benchmark\",\"iss\":\"assets-bench\",\"iat\":$(date +%s),\"exp\":${exp_time},\"roles\":[\"admin\"]}" \
        | openssl base64 -A | tr '+/' '-_' | tr -d '=')

    local signature
    signature=$(echo -n "${header}.${payload}" \
        | openssl dgst -sha256 -hmac "$secret" -binary \
        | openssl base64 -A | tr '+/' '-_' | tr -d '=')

    echo "${header}.${payload}.${signature}"
}

parse_hey_output() {
    local file="$1"
    local output
    output=$(<"$file")

    HEY_RPS=$(echo "$output" | grep "Requests/sec:" | awk '{print $2}') || true
    HEY_RPS="${HEY_RPS:-0}"

    HEY_DURATION=$(echo "$output" | grep "Total:" | head -1 | awk '{print $2}') || true
    HEY_DURATION="${HEY_DURATION:-0}"

    HEY_LAT_AVG=$(echo "$output" | grep "Average:" | head -1 | awk '{print $2}') || true
    HEY_LAT_AVG="${HEY_LAT_AVG:-0}"

    HEY_LAT_P50=$(echo "$output" | grep "50%% in" | head -1 | awk '{print $3}') || true
    HEY_LAT_P50="${HEY_LAT_P50:-0}"

    HEY_LAT_P95=$(echo "$output" | grep "95%% in" | head -1 | awk '{print $3}') || true
    HEY_LAT_P95="${HEY_LAT_P95:-0}"

    HEY_LAT_P99=$(echo "$output" | grep "99%% in" | head -1 | awk '{print $3}') || true
    HEY_LAT_P99="${HEY_LAT_P99:-0}"

    HEY_LAT_AVG_MS=$(awk "BEGIN{printf \"%.2f\", ${HEY_LAT_AVG}*1000}")
    HEY_LAT_P50_MS=$(awk "BEGIN{printf \"%.2f\", ${HEY_LAT_P50}*1000}")
    HEY_LAT_P95_MS=$(awk "BEGIN{printf \"%.2f\", ${HEY_LAT_P95}*1000}")
    HEY_LAT_P99_MS=$(awk "BEGIN{printf \"%.2f\", ${HEY_LAT_P99}*1000}")

    HEY_2XX=$({ echo "$output" | grep -E "^\s*\[2[0-9]{2}\]" || true; } | awk '{sum+=$2} END {print sum+0}')
    HEY_4XX=$({ echo "$output" | grep -E "^\s*\[4[0-9]{2}\]" || true; } | awk '{sum+=$2} END {print sum+0}')
    HEY_5XX=$({ echo "$output" | grep -E "^\s*\[5[0-9]{2}\]" || true; } | awk '{sum+=$2} END {print sum+0}')
}

parse_auth_bench_output() {
    local file="$1"
    AUTH_RPS=$(grep "^AUTH_BENCH_RPS=" "$file" | cut -d= -f2) || true
    AUTH_AVG_MS=$(grep "^AUTH_BENCH_AVG_MS=" "$file" | cut -d= -f2) || true
    AUTH_P50_MS=$(grep "^AUTH_BENCH_P50_MS=" "$file" | cut -d= -f2) || true
    AUTH_P95_MS=$(grep "^AUTH_BENCH_P95_MS=" "$file" | cut -d= -f2) || true
    AUTH_P99_MS=$(grep "^AUTH_BENCH_P99_MS=" "$file" | cut -d= -f2) || true
    AUTH_OK=$(grep "^AUTH_BENCH_OK=" "$file" | cut -d= -f2) || true
    AUTH_FAIL=$(grep "^AUTH_BENCH_FAIL=" "$file" | cut -d= -f2) || true
    AUTH_TIMEOUT=$(grep "^AUTH_BENCH_TIMEOUT=" "$file" | cut -d= -f2) || true
    AUTH_DURATION_S=$(grep "^AUTH_BENCH_DURATION_S=" "$file" | cut -d= -f2) || true
}

# get_http_scenario_config sets target_url, http_method, payload_file for a scenario.
get_http_scenario_config() {
    local scenario="$1"

    case "$scenario" in
        list-assets)
            target_url="${BASE_URL}/api/v1/assets?page=1&perPage=20"
            http_method="GET"
            payload_file=""
            ;;
        get-asset)
            target_url="${BASE_URL}/api/v1/assets/${BENCH_ASSET_ID}"
            http_method="GET"
            payload_file=""
            ;;
        create-asset)
            target_url="${BASE_URL}/api/v1/assets"
            http_method="POST"
            payload_file="${PAYLOADS_DIR}/create-asset.json"
            ;;
        get-template)
            target_url="${BASE_URL}/api/v1/asset_templates/${BENCH_TEMPLATE_ID}"
            http_method="GET"
            payload_file=""
            ;;
    esac
}

# is_http_scenario returns 0 if the scenario is an HTTP scenario.
is_http_scenario() {
    case "$1" in
        list-assets|get-asset|create-asset|get-template) return 0 ;;
        *) return 1 ;;
    esac
}

# is_auth_scenario returns 0 if the scenario is an auth scenario.
is_auth_scenario() {
    case "$1" in
        auth-cache-hit|auth-cache-miss) return 0 ;;
        *) return 1 ;;
    esac
}

# get_request_count returns the effective request count for a scenario.
# Write scenarios (create-asset) use 10x fewer requests.
get_request_count() {
    local scenario="$1"
    if [ "$scenario" = "create-asset" ]; then
        local count=$(( REQUEST_COUNT / 10 ))
        [ "$count" -lt 1000 ] && count=1000
        echo "$count"
    else
        echo "$REQUEST_COUNT"
    fi
}

# ─── Kill Competing Services ────────────────────────────────
kill_competing_services

# ─── Preflight Checks ────────────────────────────────────────

run_preflight_checks "" "$BENCH_CLI_TOOLS" "$BENCH_SERVICE_CHECKS"

# ─── Build Binaries ──────────────────────────────────────────

echo ""
log_info "Compiling assets service..."
cd "$SERVICE_DIR"
GOWORK=off CGO_ENABLED=0 go build -o "$BINARY_PATH" ./src/ 2>&1
log_success "Binary: $BINARY_PATH"

if [ "$NEEDS_AUTH" = "true" ]; then
    log_info "Compiling auth-bench tool..."
    cd "$AUTH_BENCH_DIR"
    GOWORK=off CGO_ENABLED=0 go build -o "$AUTH_BENCH_BIN" . 2>&1
    log_success "Binary: $AUTH_BENCH_BIN"
fi

# ─── Seed Data ───────────────────────────────────────────────

echo ""
log_info "Seeding benchmark data via seed.sh..."
bash "$SCRIPT_DIR/seed.sh" setup

# Generate JWT for HTTP scenarios
JWT_TOKEN=$(generate_jwt "$JWT_SECRET")
log_success "JWT token generated (24h)."

# Kill stale processes
kill_service_on_port "$SERVICE_PORT"

mkdir -p "$RESULTS_DIR"

# ─── Count total tests ───────────────────────────────────────

NUM_SCENARIOS=$(echo "$SCENARIOS" | wc -w)
NUM_CPU=$(echo "$CPU_LIST" | wc -w)
TOTAL_TESTS=$((NUM_SCENARIOS * NUM_CPU))

echo ""
echo "================================================================"
echo "  Assets Service Full CPU Benchmark (cgroup v2)"
echo "================================================================"
echo "  Scenarios:      $SCENARIOS"
echo "  Requests/test:  $(format_number "$REQUEST_COUNT")"
echo "  CPU list:       $CPU_LIST"
echo "  Total tests:    $TOTAL_TESTS (${NUM_SCENARIOS} scenarios x ${NUM_CPU} CPU)"
echo "  MongoDB:        $MONGO_DB"
echo "  Results:        $RESULTS_DIR"
echo "================================================================"
echo ""

# ─── Main Benchmark Loop ─────────────────────────────────────

declare -a SUMMARY_ROWS=()
TEST_NUM=0

for SCENARIO in $SCENARIOS; do
    log_info "Scenario block: $SCENARIO"

    SCENARIO_REQ_COUNT=$(get_request_count "$SCENARIO")

    for CPU in $CPU_LIST; do
        TEST_NUM=$((TEST_NUM + 1))
        CORE_RANGE=$(set_shield_cpus "$CPU")

        echo ""
        echo "════════════════════════════════════════════════════════════════"
        echo "  TEST $TEST_NUM/$TOTAL_TESTS — scenario=$SCENARIO  GOMAXPROCS=$CPU  cores=$CORE_RANGE"
        echo "  $(format_number "$SCENARIO_REQ_COUNT") requests"
        echo "════════════════════════════════════════════════════════════════"
        echo ""

        # 0. Kill any leftover process
        kill_service_on_port "$SERVICE_PORT"

        # ── Auth scenario: flush Redis (clean state) ───────────
        # Redis auth cache is populated by the SERVICE during warmup,
        # not by us. We only flush to ensure a clean starting point.
        if is_auth_scenario "$SCENARIO"; then
            log_step "1" "7" "Flushing Redis (clean state for auth scenario)..."
            redis_flush 0
            log_success "Redis flushed."
        else
            log_step "1" "7" "Purging NATS streams..."
            bench_purge_all_streams
        fi

        # 2. Start service
        log_step "2" "7" "Starting assets service (GOMAXPROCS=$CPU, LOG_LEVEL=silent)..."
        env GO_ENV="$GO_ENV_VALUE" LOG_LEVEL=silent GOMAXPROCS="$CPU" \
            CTX_TIMEOUT=30 \
            METRICS_GO_COLLECTOR=true METRICS_PROCESS_COLLECTOR=true \
            "$BINARY_PATH" > /dev/null 2>&1 &
        SVC_PID=$!
        sleep 1
        move_to_shield "$SVC_PID" || log_warn "Failed to move PID $SVC_PID to cgroup"
        echo "      PID=$SVC_PID -> cores $CORE_RANGE"

        ACTUAL_CPUS=$(grep Cpus_allowed_list /proc/$SVC_PID/status 2>/dev/null | awk '{print $2}') || true
        echo "      Verified: $ACTUAL_CPUS"

        # 3. Wait for ready
        log_step "3" "7" "Waiting for service..."
        if ! wait_for_service_ready "$METRICS_URL" 60 "$SVC_PID"; then
            echo "      SKIP scenario=$SCENARIO CPU=$CPU"
            kill "$SVC_PID" 2>/dev/null || true; wait "$SVC_PID" 2>/dev/null || true
            continue
        fi

        # ── HTTP scenario: warmup + hey ──────────────────────
        if is_http_scenario "$SCENARIO"; then
            target_url="" http_method="" payload_file=""
            get_http_scenario_config "$SCENARIO"

            # 4. Warmup
            log_step "4" "7" "Warmup (1000 requests)..."
            if [ "$http_method" = "POST" ] && [ -n "$payload_file" ]; then
                hey -n 1000 -c 50 -m POST \
                    -H "Authorization: Bearer ${JWT_TOKEN}" \
                    -H "X-Org-Context: ${BENCH_ORG_ID}" \
                    -H "Content-Type: application/json" \
                    -D "$payload_file" "$target_url" > /dev/null 2>&1 || true
            else
                hey -n 1000 -c 50 -m GET \
                    -H "Authorization: Bearer ${JWT_TOKEN}" \
                    -H "X-Org-Context: ${BENCH_ORG_ID}" \
                    "$target_url" > /dev/null 2>&1 || true
            fi
            sleep 2
            echo "      Done."

            # 5. Benchmark
            # Clamp concurrency to request count (hey exits non-zero if -c > -n)
            effective_concurrency=$CONCURRENCY
            if [ "$SCENARIO_REQ_COUNT" -lt "$CONCURRENCY" ]; then
                effective_concurrency=$SCENARIO_REQ_COUNT
            fi

            log_step "5" "7" "Running: $(format_number "$SCENARIO_REQ_COUNT") requests..."
            HEY_FILE="$RESULTS_DIR/test-${SCENARIO}-cpu${CPU}-hey.txt"

            if [ "$http_method" = "POST" ] && [ -n "$payload_file" ]; then
                hey -n "$SCENARIO_REQ_COUNT" -c "$effective_concurrency" -m POST \
                    -H "Authorization: Bearer ${JWT_TOKEN}" \
                    -H "X-Org-Context: ${BENCH_ORG_ID}" \
                    -H "Content-Type: application/json" \
                    -D "$payload_file" "$target_url" > "$HEY_FILE" 2>&1
            else
                hey -n "$SCENARIO_REQ_COUNT" -c "$effective_concurrency" -m GET \
                    -H "Authorization: Bearer ${JWT_TOKEN}" \
                    -H "X-Org-Context: ${BENCH_ORG_ID}" \
                    "$target_url" > "$HEY_FILE" 2>&1
            fi
            echo "      Done."
        fi

        # ── Auth scenario: auth-bench ────────────────────────
        if is_auth_scenario "$SCENARIO"; then
            # cache-hit: Warmup = ALL users (service populates Redis via AppCache.SetEx).
            #            Then benchmark measures pure Redis fast path.
            # cache-miss: Warmup = 0 (Redis already flushed above).
            #            user-count = request-count → every request hits MongoDB.
            auth_bench_flags=""
            auth_warmup=0
            if [ "$SCENARIO" = "auth-cache-hit" ]; then
                # Warmup must cover ALL unique users so service populates Redis.
                auth_warmup="$BENCH_MQTT_ASSET_COUNT"
                auth_bench_flags="--scenario cache-hit --user-count $BENCH_MQTT_ASSET_COUNT"
            else
                # user-count = request count → each request uses a unique user → 100% DB hits.
                # Cap at BENCH_MQTT_ASSET_COUNT (seeded in MongoDB).
                auth_warmup=0
                auth_user_count="$SCENARIO_REQ_COUNT"
                if [ "$auth_user_count" -gt "$BENCH_MQTT_ASSET_COUNT" ]; then
                    auth_user_count="$BENCH_MQTT_ASSET_COUNT"
                fi
                auth_bench_flags="--scenario cache-miss --user-count $auth_user_count"
            fi

            log_step "4" "7" "Auth warmup ($auth_warmup requests — service populates Redis)..."
            echo "      Done."

            log_step "5" "7" "Running auth-bench ($(format_number "$SCENARIO_REQ_COUNT") requests)..."
            AUTH_FILE="$RESULTS_DIR/test-${SCENARIO}-cpu${CPU}-auth.txt"

            # shellcheck disable=SC2086
            "$AUTH_BENCH_BIN" \
                --count "$SCENARIO_REQ_COUNT" \
                --concurrency "$CONCURRENCY" \
                --warmup "$auth_warmup" \
                --mqtt-broker "$MQTT_BROKER" \
                $auth_bench_flags \
                > "$AUTH_FILE" 2>&1
            echo "      Done."
        fi

        # 6. Collect metrics
        log_step "6" "7" "Collecting metrics..."
        sleep 2
        METRICS_FILE="$RESULTS_DIR/test-${SCENARIO}-cpu${CPU}-metrics.txt"
        METRICS=$(curl -s "$METRICS_URL" 2>/dev/null)
        echo "$METRICS" > "$METRICS_FILE"

        # Parse Go runtime (via common module)
        parse_go_runtime_metrics "$METRICS"

        # 7. Stop service
        log_step "7" "7" "Stopping service..."
        kill "$SVC_PID" 2>/dev/null || true
        wait "$SVC_PID" 2>/dev/null || true
        echo "      Done."

        # ── Build summary row ────────────────────────────────
        if is_http_scenario "$SCENARIO"; then
            parse_hey_output "$HEY_FILE"

            ASSET_OPS_OK=$(extract_metric "$METRICS" "assets_asset_operations_total" 'status="success"')
            ASSET_OPS_ERR=$(extract_metric "$METRICS" "assets_asset_operations_total" 'status="error"')
            ASSET_CACHE_HIT=$(extract_metric "$METRICS" "assets_asset_cache_total" 'result="hit"')
            ASSET_CACHE_MISS=$(extract_metric "$METRICS" "assets_asset_cache_total" 'result="miss"')

            SUMMARY_ROWS+=("$(printf "  │ %-16s │ %3s │ %8s │ %7s/s │ %5sms │ %5sms │ %5sMB │" \
                "$SCENARIO" "$CPU" "$CORE_RANGE" "$HEY_RPS" "$HEY_LAT_P50_MS" "$HEY_LAT_P99_MS" "$GO_RSS_MB")")

            echo ""
            echo "  ┌──────────────────────────────────────────────────────────┐"
            echo "  │ scenario=$SCENARIO CPU=$CPU -> ${HEY_RPS} req/s  (${HEY_DURATION}s)"
            echo "  │ Lat: avg=${HEY_LAT_AVG_MS}ms p50=${HEY_LAT_P50_MS}ms p99=${HEY_LAT_P99_MS}ms"
            echo "  │ RSS=${GO_RSS_MB}MB  Heap=${GO_HEAP_MB}MB  Goroutines=${GO_GOROUTINES}"
            echo "  │ 2xx=${HEY_2XX} 4xx=${HEY_4XX} 5xx=${HEY_5XX}"
            echo "  │ AssetOps: ok=${ASSET_OPS_OK} err=${ASSET_OPS_ERR} cache=${ASSET_CACHE_HIT}hit/${ASSET_CACHE_MISS}miss"
            echo "  └──────────────────────────────────────────────────────────┘"
        fi

        if is_auth_scenario "$SCENARIO"; then
            parse_auth_bench_output "$AUTH_FILE"

            AUTH_CALLOUT_OK=$(extract_metric "$METRICS" "assets_auth_callout_total" 'status="success"')
            AUTH_CACHE_HIT=$(extract_metric "$METRICS" "assets_auth_cache_total" 'result="hit"')
            AUTH_CACHE_MISS=$(extract_metric "$METRICS" "assets_auth_cache_total" 'result="miss"')
            AUTH_CALLOUT_AVG_MS=$(extract_histogram_avg "$METRICS" "assets_auth_callout_duration_seconds")

            SUMMARY_ROWS+=("$(printf "  │ %-16s │ %3s │ %8s │ %7s/s │ %5sms │ %5sms │ %5sMB │" \
                "$SCENARIO" "$CPU" "$CORE_RANGE" "${AUTH_RPS:-0}" "${AUTH_P50_MS:-0}" "${AUTH_P99_MS:-0}" "$GO_RSS_MB")")

            echo ""
            echo "  ┌──────────────────────────────────────────────────────────┐"
            echo "  │ scenario=$SCENARIO CPU=$CPU -> ${AUTH_RPS:-0} req/s  (${AUTH_DURATION_S:-0}s)"
            echo "  │ Lat: avg=${AUTH_AVG_MS:-0}ms p50=${AUTH_P50_MS:-0}ms p99=${AUTH_P99_MS:-0}ms"
            echo "  │ RSS=${GO_RSS_MB}MB  Heap=${GO_HEAP_MB}MB  Goroutines=${GO_GOROUTINES}"
            echo "  │ OK=${AUTH_OK:-0} Fail=${AUTH_FAIL:-0} Timeout=${AUTH_TIMEOUT:-0}"
            echo "  │ Service: callout=${AUTH_CALLOUT_OK} cache=${AUTH_CACHE_HIT}hit/${AUTH_CACHE_MISS}miss avg=${AUTH_CALLOUT_AVG_MS}ms"
            echo "  └──────────────────────────────────────────────────────────┘"
        fi

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

TABLE_HEADER=$(printf "  │ %-16s │ %3s │ %8s │ %10s │ %8s │ %8s │ %7s │" \
    "Scenario" "CPU" "Cores" "Req/s" "Lat p50" "Lat p99" "RSS MB")

echo "  Summary:"
echo "  ┌──────────────────┬─────┬──────────┬────────────┬──────────┬──────────┬─────────┐"
echo "$TABLE_HEADER"
echo "  ├──────────────────┼─────┼──────────┼────────────┼──────────┼──────────┼─────────┤"

for ROW in "${SUMMARY_ROWS[@]}"; do
    echo "$ROW"
done

echo "  └──────────────────┴─────┴──────────┴────────────┴──────────┴──────────┴─────────┘"
echo ""

# ─── Teardown ─────────────────────────────────────────────────

log_info "Removing benchmark data via seed.sh teardown..."
bash "$SCRIPT_DIR/seed.sh" teardown

# Remove temp binaries
rm -f "$BINARY_PATH"
[ "$NEEDS_AUTH" = "true" ] && rm -f "$AUTH_BENCH_BIN"

log_success "Done!"
echo "  Results stored at: $RESULTS_DIR"
