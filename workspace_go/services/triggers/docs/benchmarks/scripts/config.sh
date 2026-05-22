#!/bin/bash
# =============================================================
# Service Configuration — Triggers Benchmark
#
# Constants + common bootstrap. Source this first in any script:
#   source "$SCRIPT_DIR/config.sh"
#
# Then source the common services you need:
#   source "$COMMON_DIR/services/nats.sh"
#   source "$COMMON_DIR/services/triggers.sh"
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

# ─── Triggers Constants ─────────────────────────────────────

# Environment
GO_ENV_VALUE="${GO_ENV:-dev}"

# NATS — primary stream consumed by triggers
STREAM="TRIGGERS"
SUBJECT="trigger.bench.execute"
CONSUMER="triggers-trigger-executor"

# All NATS streams to purge during cleanup (source + downstream)
ALL_STREAMS=(
    "TRIGGERS"
    "MAPEXOS-DLQ"
)

# MongoDB — matches GO_ENV prefix convention: {go_env}-triggers
MONGO_DB="${MONGO_DB:-${GO_ENV_VALUE}-triggers}"
MONGO_COLLECTION="triggers"

# Service
SERVICE_PORT=5006
METRICS_URL="http://localhost:${SERVICE_PORT}/metrics"

# ─── Deterministic Seed IDs (24-char hex, zero-padded) ───────

# Org (shared across services)
BENCH_ORG_ID="000000000000000000000001"

# Per-scenario trigger IDs — one per executor type
BENCH_TRIGGER_HTTP_ID="000000000000000000000030"
BENCH_TRIGGER_MQTT_ID="000000000000000000000031"
BENCH_TRIGGER_NATS_ID="000000000000000000000032"
BENCH_TRIGGER_RABBITMQ_ID="000000000000000000000033"
BENCH_TRIGGER_EMAIL_ID="000000000000000000000034"

# ─── Mock Server Ports ──────────────────────────────────────

MOCK_HTTP_PORT=9999
MOCK_MQTT_PORT=1884
MOCK_NATS_PORT=4333
MOCK_RABBITMQ_PORT=5673
MOCK_SMTP_PORT=2525

# ─── Benchmark Isolation ─────────────────────────────────────
# No dependencies need to stay alive for triggers benchmarks
BENCH_KEEP_PORTS=()

# ─── Preflight ──────────────────────────────────────────────

BENCH_CLI_TOOLS="mongosh,nats,redis-cli,curl"
BENCH_SERVICE_CHECKS="MongoDB:mongo_check,NATS:nats_cmd stream list --json,Redis App:redis_cmd PING"
