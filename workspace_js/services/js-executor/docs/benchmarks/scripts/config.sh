#!/bin/bash
# =============================================================
# Service Configuration — JS-Executor Benchmark
#
# Constants + common bootstrap. Source this first in any script:
#   source "$SCRIPT_DIR/config.sh"
#
# Then source the common services you need:
#   source "$COMMON_DIR/services/nats.sh"
#   source "$COMMON_DIR/services/assets.sh"
# =============================================================

# Guard: only source once
[ -n "${_BENCH_CONFIG_LOADED:-}" ] && return 0
_BENCH_CONFIG_LOADED=1

# ─── Paths ──────────────────────────────────────────────────

_SCRIPTS_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
COMMON_DIR="$(cd "$_SCRIPTS_DIR/../../../../../.." && pwd)/scripts/benchmarks/common"
SEED_DIR="$_SCRIPTS_DIR/../seed"

# ─── Source Common Modules ──────────────────────────────────

source "$COMMON_DIR/init.sh"

# ─── JS-Executor Constants ──────────────────────────────────

# Environment
GO_ENV_VALUE="${GO_ENV:-dev}"

# Primary stream and consumer
STREAM="PROCESSOR-JS-EXECUTE"
CONSUMER="processor-js-execute"
SUBJECT="processor.js.execute"

# All NATS streams to purge during cleanup
ALL_STREAMS=(
    "PROCESSOR-JS-EXECUTE"
    "EVENTS-JSEXEC"
    "EVENTS-RAW"
    "EVENTS-ROUTER"
    "EVENTS-TRIGGER"
    "EVENTS-BUSINESSRULE"
    "EVENTS"
    "ROUTE-GROUPS"
    "RULE-ENGINE"
    "TRIGGERS"
    "FANOUT"
    "MAPEXOS-DLQ"
    "MAPEXOS"
    "MAPEXOS-LISTS"
    "MAPEXOS_CACHE_INVALIDATION"
    "MQTT-DATA"
)

# ─── Assets Constants ────────────────────────────────────────
# Used by common/services/assets.sh

BENCH_ORG_ID="000000000000000000000001"
BENCH_TMPL_ID="000000000000000000000001"
BENCH_ASSET_ID="000000000000000000000002"
BENCH_ASSET_UUID="bench-asset-uuid-0001"

# MongoDB (dev- prefix for development environment)
ASSETS_DB="${ASSETS_DB:-${GO_ENV_VALUE}-assets}"

# MinIO
MINIO_ALIAS="${MINIO_ALIAS:-local}"
MINIO_ASSETS_BUCKET="${MINIO_ASSETS_BUCKET:-mapex-assets}"
MINIO_TEMPLATES_BUCKET="${MINIO_TEMPLATES_BUCKET:-mapex-templates}"

# ─── Service ───────────────────────────────────────────────

SERVICE_PORT=8000
METRICS_URL="http://localhost:${SERVICE_PORT}/metrics"

# ─── Benchmark Isolation ─────────────────────────────────────
# No dependencies on other services — kill everything
BENCH_KEEP_PORTS=()

# ─── Preflight ──────────────────────────────────────────────

BENCH_CLI_TOOLS="nats,mongosh,mc,curl"
BENCH_SERVICE_CHECKS="NATS:nats_cmd stream list --json,Redis:redis_check,MongoDB:mongo_check,MinIO:mc ls ${MINIO_ALIAS:-local}"
