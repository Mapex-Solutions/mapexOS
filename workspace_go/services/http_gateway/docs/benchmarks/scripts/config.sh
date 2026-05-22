#!/bin/bash
# =============================================================
# Service Configuration — HTTP Gateway Benchmark
#
# Constants + common bootstrap. Source this first in any script:
#   source "$SCRIPT_DIR/config.sh"
#
# Then source the common services you need:
#   source "$COMMON_DIR/services/nats.sh"
#   source "$COMMON_DIR/services/datasources.sh"
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

# ─── HTTP Gateway Constants ─────────────────────────────────

# Environment
GO_ENV_VALUE="${GO_ENV:-dev}"

# NATS — primary stream (downstream from http_gateway)
STREAM="PROCESSOR-JS-EXECUTE"

# All NATS streams to purge during cleanup (downstream pipeline)
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

# MongoDB — matches GO_ENV prefix convention: {go_env}-http_gateway
MONGO_DB="${MONGO_DB:-${GO_ENV_VALUE}-http_gateway}"
MONGO_COLLECTION="data_sources"

# Service
SERVICE_PORT=5001
METRICS_URL="http://localhost:${SERVICE_PORT}/metrics"

# ─── Deterministic Seed IDs (24-char hex, zero-padded) ───────

BENCH_ORG_ID="000000000000000000000001"
BENCH_DS_JWT_ID="000000000000000000000001"
BENCH_DS_APIKEY_ID="000000000000000000000002"
BENCH_DS_IPWHITELIST_ID="000000000000000000000003"
BENCH_DS_NONE_ID="000000000000000000000004"
BENCH_ASSET_ID="000000000000000000000002"

# ─── Auth Credentials ────────────────────────────────────────
# Must match what the seeded DataSources expect

JWT_SECRET="${JWT_SECRET:-2ca98cd7-4619-4f80-b80a-94df59e74ef7}"
API_KEY_VALUE="${API_KEY_VALUE:-bench-api-key-2026-secure-token}"

# Static JWT token (exp: 1802372218 ~2027-02-11, signed with JWT_SECRET)
JWT_TOKEN_STATIC="eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJPbmxpbmUgSldUIEJ1aWxkZXIiLCJpYXQiOjE3NzA4MzYyMTgsImV4cCI6MTgwMjM3MjIxOCwiYXVkIjoid3d3LmV4YW1wbGUuY29tIiwic3ViIjoianJvY2tldEBleGFtcGxlLmNvbSIsIkdpdmVuTmFtZSI6IkpvaG5ueSIsIlN1cm5hbWUiOiJSb2NrZXQiLCJFbWFpbCI6Impyb2NrZXRAZXhhbXBsZS5jb20iLCJSb2xlIjpbIk1hbmFnZXIiLCJQcm9qZWN0IEFkbWluaXN0cmF0b3IiXX0.4nRIcCH5IHF-hG1_f6HHgNPmKKtSQrMK1Lox99oZGN4"

# ─── Benchmark Isolation ─────────────────────────────────────
# No dependencies on other services — kill everything
BENCH_KEEP_PORTS=()

# ─── Preflight ──────────────────────────────────────────────

BENCH_CLI_TOOLS="hey,mongosh,nats,redis-cli,curl"
BENCH_SERVICE_CHECKS="MongoDB:mongo_check,NATS:nats_cmd stream list --json,Redis App:redis_cmd PING,Redis Shared:redis_cmd -n 5 PING"
