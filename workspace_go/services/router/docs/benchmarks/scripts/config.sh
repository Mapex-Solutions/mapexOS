#!/bin/bash
# =============================================================
# Service Configuration — Router Benchmark
#
# Constants + common bootstrap. Source this first in any script:
#   source "$SCRIPT_DIR/config.sh"
#
# Then source the common services you need:
#   source "$COMMON_DIR/services/nats.sh"
#   source "$COMMON_DIR/services/routegroups.sh"
#   source "$COMMON_DIR/services/minio.sh"
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

# ─── Router Constants ───────────────────────────────────────

# Environment
GO_ENV_VALUE="${GO_ENV:-dev}"

# NATS — primary stream consumed by router
STREAM="ROUTE-GROUPS"
SUBJECT="route.bench.execute"
CONSUMER="${GO_ENV_VALUE}-router-execute"

# All NATS streams to purge during cleanup (source + downstream)
ALL_STREAMS=(
    "ROUTE-GROUPS"
    "EVENTS-ROUTER"
    "EVENTS-TRIGGER"
    "EVENTS-BUSINESSRULE"
    "EVENTS"
    "RULE-ENGINE"
    "TRIGGERS"
    "MAPEXOS-DLQ"
)

# MongoDB — matches GO_ENV prefix convention: {go_env}-router
MONGO_DB="${MONGO_DB:-${GO_ENV_VALUE}-router}"
MONGO_COLLECTION="routegroups"

# Service
SERVICE_PORT=5003
METRICS_URL="http://localhost:${SERVICE_PORT}/metrics"

# ─── Deterministic Seed IDs (24-char hex, zero-padded) ───────

# Org (shared with http_gateway)
BENCH_ORG_ID="000000000000000000000001"

# Per-scenario asset UUIDs (1 RouteGroup per asset for isolated benchmarks)
BENCH_ASSET_SAVE_EVENT_UUID="bench-asset-save-event"
BENCH_ASSET_RULE_ENGINE_UUID="bench-asset-rule-engine"
BENCH_ASSET_TRIGGER_UUID="bench-asset-trigger"

# RouteGroup IDs — one per scenario
BENCH_RG_SAVE_EVENT_ID="000000000000000000000010"
BENCH_RG_RULE_ENGINE_ID="000000000000000000000011"
BENCH_RG_TRIGGER_ID="000000000000000000000012"

# Referenced entity IDs (rule engine + trigger targets)
BENCH_BUSINESS_RULE_ID="000000000000000000000020"
BENCH_TRIGGER_REF_ID="000000000000000000000021"

# ─── MinIO ──────────────────────────────────────────────────

MINIO_ALIAS="${MINIO_ALIAS:-local}"
MINIO_BUCKET="${MINIO_BUCKET:-mapex-assets}"

# ─── Benchmark Isolation ─────────────────────────────────────
# Assets (5002) must stay alive — router uses it for tiered cache fallback
BENCH_KEEP_PORTS=(5002)

# ─── Preflight ──────────────────────────────────────────────

BENCH_CLI_TOOLS="mongosh,nats,redis-cli,mc,curl"
BENCH_SERVICE_CHECKS="MongoDB:mongo_check,NATS:nats_cmd stream list --json,Redis App:redis_cmd PING,Redis Shared:redis_cmd -n 5 PING,MinIO:mc ls ${MINIO_ALIAS}/${MINIO_BUCKET}"
