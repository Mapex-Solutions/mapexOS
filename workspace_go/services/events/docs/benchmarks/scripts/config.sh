#!/bin/bash
# =============================================================
# Service Configuration — Events Benchmark
#
# Constants + common bootstrap. Source this first in any script:
#   source "$SCRIPT_DIR/config.sh"
#
# Then source the common services you need:
#   source "$COMMON_DIR/services/nats.sh"
# =============================================================

# Guard: only source once
[ -n "${_BENCH_CONFIG_LOADED:-}" ] && return 0
_BENCH_CONFIG_LOADED=1

# ─── Paths ──────────────────────────────────────────────────

_SCRIPTS_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
COMMON_DIR="$(cd "$_SCRIPTS_DIR/../../../../../.." && pwd)/scripts/benchmarks/common"
SEED_DIR="$_SCRIPTS_DIR/../seed"
SEEDS_SQL_DIR="$(cd "$_SCRIPTS_DIR/../../../../../packages/seeds/clickhouse/events/v1" && pwd)"

# ─── Source Common Modules ──────────────────────────────────

source "$COMMON_DIR/init.sh"

# ─── Events Constants ───────────────────────────────────────

# Environment
GO_ENV_VALUE="${GO_ENV:-dev}"

# Service
SERVICE_PORT=5004
METRICS_URL="http://localhost:${SERVICE_PORT}/metrics"

# MongoDB — matches GO_ENV prefix convention: events (no prefix)
MONGO_DB="${MONGO_DB:-events}"

# ClickHouse
CH_HOST="${CLICKHOUSE_HOST:-localhost}"
CH_PORT="${CLICKHOUSE_PORT:-9440}"
CH_DB="${CLICKHOUSE_DATABASE:-mapexos}"
CH_USER="${CLICKHOUSE_USERNAME:-mapexos_user}"
CH_PASS="${CLICKHOUSE_PASSWORD:-mapexos_password}"

# ─── Deterministic Seed IDs (24-char hex, zero-padded) ───────

BENCH_ORG_ID="000000000000000000000001"
BENCH_TEMPLATE_ID="bench-template-001"

# MinIO (TieredCache L2 for templates)
MINIO_ALIAS="${MINIO_ALIAS:-local}"
MINIO_TEMPLATES_BUCKET="${MINIO_TEMPLATES_BUCKET:-mapex-templates}"
MINIO_TEMPLATE_PATH="templates/${BENCH_ORG_ID}/${BENCH_TEMPLATE_ID}.json"

# ─── Scenario Configuration ─────────────────────────────────
#
# Maps each benchmark scenario to its NATS stream/subject,
# consumer name, ClickHouse table, and Prometheus metric labels.

# All available scenarios
ALL_SCENARIOS="save_raw_event save_jsexec_event save_router_event save_businessrule_event save_trigger_event save_event save_dlq_event"

# Scenario → NATS Config: STREAM|SUBJECT|CONSUMER_SUFFIX
declare -A SCENARIO_NATS=(
    [save_raw_event]="EVENTS-RAW|events.raw|events-raw"
    [save_jsexec_event]="EVENTS-JSEXEC|events.logs.jsexecutor|events-jsexec"
    [save_router_event]="EVENTS-ROUTER|events.router|events-router"
    [save_businessrule_event]="EVENTS-BUSINESSRULE|events.businessrule|events-businessrule"
    [save_trigger_event]="EVENTS-TRIGGER|events.trigger|events-trigger"
    [save_event]="EVENTS|events.save|events-save"
    [save_dlq_event]="MAPEXOS-DLQ|dlq.mapexos|events-dlq"
)

# Scenario → Prometheus Metric Labels: CONSUMER_LABEL|TABLE_LABEL
declare -A SCENARIO_METRICS=(
    [save_raw_event]="raw|eventsRaw"
    [save_jsexec_event]="jsexec|eventsJsExecutor"
    [save_router_event]="router|eventsRouter"
    [save_businessrule_event]="businessrule|eventsBusinessRule"
    [save_trigger_event]="trigger|eventsTrigger"
    [save_event]="store|events"
    [save_dlq_event]="dlq|eventsDLQ"
)

# Scenario → ClickHouse Table Name (for seed.sh)
declare -A SCENARIO_CH_TABLE=(
    [save_raw_event]="events_raw"
    [save_jsexec_event]="events_jsexecutor"
    [save_router_event]="events_router"
    [save_businessrule_event]="events_businessrule"
    [save_trigger_event]="events_trigger"
    [save_event]="events"
    [save_dlq_event]="events_dlq"
)

# ALL_STREAMS — used by bench_purge_all_streams (services/nats.sh)
ALL_STREAMS=(
    "EVENTS-RAW"
    "EVENTS-JSEXEC"
    "EVENTS-ROUTER"
    "EVENTS-BUSINESSRULE"
    "EVENTS-TRIGGER"
    "EVENTS"
    "MAPEXOS-DLQ"
)

# ─── Scenario Helpers ────────────────────────────────────────

get_nats_stream()     { echo "${SCENARIO_NATS[$1]}" | cut -d'|' -f1; }
get_nats_subject()    { echo "${SCENARIO_NATS[$1]}" | cut -d'|' -f2; }
get_consumer_suffix() { echo "${SCENARIO_NATS[$1]}" | cut -d'|' -f3; }
get_metric_consumer() { echo "${SCENARIO_METRICS[$1]}" | cut -d'|' -f1; }
get_metric_table()    { echo "${SCENARIO_METRICS[$1]}" | cut -d'|' -f2; }
get_ch_table()        { echo "${SCENARIO_CH_TABLE[$1]}"; }

validate_scenario() {
    local scenario="$1"
    if [ -z "${SCENARIO_NATS[$scenario]+x}" ]; then
        log_error "Invalid scenario '$scenario'"
        echo "  Available: $ALL_SCENARIOS"
        return 1
    fi
    return 0
}

# ─── ClickHouse Helper ───────────────────────────────────────

ch_query() {
    if command -v clickhouse-client &> /dev/null; then
        clickhouse-client \
            --host="$CH_HOST" \
            --port="$CH_PORT" \
            --database="$CH_DB" \
            --user="$CH_USER" \
            --password="$CH_PASS" \
            --query="$1" 2>/dev/null
    else
        echo "${SUDO_PASS:-}" | sudo -S docker exec mapexos-clickhouse clickhouse-client \
            --database="$CH_DB" \
            --user="$CH_USER" \
            --password="$CH_PASS" \
            --query="$1" 2>/dev/null
    fi
}

# ─── Benchmark Isolation ─────────────────────────────────────
# Events is self-contained — no other services needed
BENCH_KEEP_PORTS=()

# ─── Preflight ──────────────────────────────────────────────

BENCH_CLI_TOOLS="mongosh,nats,redis-cli,curl,mc"

ch_check() { ch_query "SELECT 1" > /dev/null 2>&1; }

BENCH_SERVICE_CHECKS="MongoDB:mongo_check,NATS:nats_cmd stream list --json,Redis App:redis_cmd PING,ClickHouse:ch_check"
