#!/bin/bash
# =============================================================
# Service Configuration — Assets Benchmark
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

# ─── Source Common Modules ──────────────────────────────────

source "$COMMON_DIR/init.sh"

# ─── Assets Constants ───────────────────────────────────────

# Environment
GO_ENV_VALUE="${GO_ENV:-dev}"

# Service
SERVICE_PORT=5002
METRICS_URL="http://localhost:${SERVICE_PORT}/metrics"

# Auth — JWT signed with the default auth_secret from config.go
JWT_SECRET="${JWT_SECRET:-a-string-secret-at-least-256-bits-long}"

# MongoDB — matches GO_ENV prefix convention: {go_env}-assets
MONGO_DB="${MONGO_DB:-${GO_ENV_VALUE}-assets}"

# NATS — fanout stream for cache invalidation
NATS_FANOUT_STREAM="${NATS_FANOUT_STREAM:-ASSETS-FANOUT}"

ALL_STREAMS=(
    "ASSETS-FANOUT"
)

# MQTT Broker — NKey-JWT decentralized auth at the leaf (port 1883).
MQTT_BROKER="${MQTT_BROKER:-localhost:1883}"

# ─── SharedCache (Redis DB 5) — used by permission/coverage middleware ──
SHARED_CACHE_DB="${REDIS_SHARED_DB:-5}"
SHARED_CACHE_PREFIX="${GO_ENV_VALUE}:shared"

# ─── Deterministic Seed IDs (24-char hex, zero-padded) ───────

# Benchmark user (for HTTP JWT auth)
BENCH_USER_ID="000000000000000000000003"

# Org — HTTP benchmark assets
BENCH_ORG_ID="000000000000000000000001"
BENCH_CUSTOMER_ID="000000000000000000000002"
BENCH_ORG_PATHKEY="bench"

# Org — MQTT auth callout assets (separate org to avoid polluting list-assets)
BENCH_MQTT_ORG_ID="000000000000000000000099"
BENCH_MQTT_CUSTOMER_ID="000000000000000000000098"
BENCH_MQTT_ORG_PATHKEY="bench-mqtt"

# Template
BENCH_TEMPLATE_ID="000000000000000000000010"

# First HTTP asset: 0x101 = 257 → 000000000000000000000101
BENCH_ASSET_ID=$(printf '%024x' 257)

# MQTT assets: 0x2001 (8193) .. (8193 + count - 1)
# Default 10,000 devices for auth callout benchmarks.
# Override with BENCH_MQTT_ASSET_COUNT env var.
BENCH_MQTT_ASSET_START=8193   # 0x2001
BENCH_MQTT_ASSET_COUNT="${BENCH_MQTT_ASSET_COUNT:-10000}"

# HTTP seed assets
SEED_ASSET_COUNT=1000

# ─── MinIO ──────────────────────────────────────────────────

MINIO_ALIAS="${MINIO_ALIAS:-local}"
MINIO_ASSETS_BUCKET="${MINIO_ASSETS_BUCKET:-mapex-assets}"
MINIO_TEMPLATES_BUCKET="${MINIO_TEMPLATES_BUCKET:-mapex-templates}"

# ─── Benchmark Isolation ─────────────────────────────────────
# Assets is self-contained — no other services needed
BENCH_KEEP_PORTS=()

# ─── Preflight ──────────────────────────────────────────────

BENCH_CLI_TOOLS="mongosh,nats,redis-cli,mc,curl,hey,openssl"
# mqtt_broker_check tests MQTT port reachability (mapex-mqtt-broker)
mqtt_broker_check() { timeout 2 bash -c "echo > /dev/tcp/${MQTT_BROKER%%:*}/${MQTT_BROKER##*:}" 2>/dev/null; }

BENCH_SERVICE_CHECKS="MongoDB:mongo_check,NATS:nats_cmd stream list --json,MQTT Broker:mqtt_broker_check,Redis App:redis_cmd PING,MinIO:mc ls ${MINIO_ALIAS}/${MINIO_ASSETS_BUCKET}"
