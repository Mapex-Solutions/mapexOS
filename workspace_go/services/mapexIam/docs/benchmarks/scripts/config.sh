#!/bin/bash
# =============================================================
# Service Configuration — MapexOS Benchmark
#
# Constants + common bootstrap. Source this first in any script:
#   source "$SCRIPT_DIR/config.sh"
#
# Then source the common services you need:
#   source "$COMMON_DIR/services/nats.sh"
#   source "$COMMON_DIR/services/mapexos.sh"
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

# ─── MapexOS Constants ─────────────────────────────────────

# Environment
GO_ENV_VALUE="${GO_ENV:-dev}"

# NATS — mapexos uses cache invalidation stream + DLQ
STREAM="MAPEXOS_CACHE_INVALIDATION"
CONSUMER="${GO_ENV_VALUE}-mapexos-cache-invalidation-consumer"

# All NATS streams to purge during cleanup
ALL_STREAMS=(
    "MAPEXOS_CACHE_INVALIDATION"
    "MAPEXOS-DLQ"
)

# MongoDB — matches GO_ENV prefix convention: {go_env}-mapexos
MONGO_DB="${MONGO_DB:-${GO_ENV_VALUE}-mapexos}"

# Service
SERVICE_PORT=5000
METRICS_URL="http://localhost:${SERVICE_PORT}/metrics"

# JWT — used by benchmark to authenticate requests
JWT_SECRET="${JWT_SECRET:-bench-secret-key-at-least-256-bits-long-for-hs256}"

# Auth — real bcrypt hash of "test@123" (cost 10, same as production)
# Generated via: go run scripts/general/generate-password-hash.go 'test@123'
BENCH_PASSWORD="test@123"
BENCH_PASSWORD_HASH='$2a$10$Ilj/twFiKOiSY4ZOJyC86.3TT2NmDItAx2tV2RXC48JutHVba0o1a'

# ─── Deterministic Seed IDs (24-char hex, zero-padded) ───────

# ID ranges per entity (non-overlapping for clean teardown):
#   Organizations:  000000000000000000010001 - 000000000000000000010XXX
#   Users:          000000000000000000020001 - 000000000000000000020XXX
#   Groups:         000000000000000000030001 - 000000000000000000030XXX
#   Roles:          000000000000000000040001 - 000000000000000000040XXX
#   Memberships:    000000000000000000050001 - 000000000000000000050XXX
#   GroupMembers:   000000000000000000060001 - 000000000000000000060XXX
#   GroupMemberships: 000000000000000000070001 - 000000000000000000070XXX

# First user ID — used for JWT generation and auth_login payload
BENCH_FIRST_USER_ID="000000000000000000020001"
# First customer org — used for X-Org-Context header
BENCH_FIRST_CUSTOMER_ORG_ID="000000000000000000010002"

# Seed counts (override via env)
BENCH_ORG_COUNT="${ORG_COUNT:-10}"
BENCH_USER_COUNT="${USER_COUNT:-1000}"
BENCH_GROUP_COUNT="${GROUP_COUNT:-100}"
BENCH_ROLE_COUNT="${ROLE_COUNT:-50}"
BENCH_MEMBERSHIP_COUNT="${MEMBERSHIP_COUNT:-500}"

# ─── Benchmark Isolation ─────────────────────────────────────
# MapexOS has no required dependency services for benchmark
BENCH_KEEP_PORTS=()

# ─── Preflight ──────────────────────────────────────────────

BENCH_CLI_TOOLS="mongosh,nats,redis-cli,curl"
BENCH_SERVICE_CHECKS="MongoDB:mongo_check,NATS:nats_cmd stream list --json,Redis App:redis_cmd PING,Redis Shared:redis_cmd -n 5 PING"
