#!/bin/bash
# =============================================================
# Seed Orchestrator — MapexOS Benchmark
#
# Manages the complete test data lifecycle:
#   setup    — Seed 7 entity types into MongoDB (orgs, roles, users, groups, memberships, group_members, group-memberships)
#   teardown — Purge NATS streams + delete seed data + flush Redis
#
# Uses modular functions from config + common/services/:
#   config.sh                       — Constants, common bootstrap
#   common/services/mapexos.sh      — MongoDB seed/cleanup
#   common/services/nats.sh         — Stream purge
#
# Seed data:
#   Users:       client{i}@test.com with bcrypt hash of "test@123"
#   Memberships: each user linked to an org with a role
#
# Usage:
#   ./seed.sh setup      — Insert benchmark data (idempotent)
#   ./seed.sh teardown   — Clean all benchmark data
#
# Prerequisites:
#   - MongoDB running at localhost:27017 (replica set: rs0)
#   - NATS running at localhost:4222 (user: service / service_secret)
#   - Redis running at localhost:6379
#   - CLIs: mongosh, nats, redis-cli
# =============================================================

set -euo pipefail

COMMAND="${1:-}"

# ─── Source Modules ──────────────────────────────────────────

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"

source "$SCRIPT_DIR/config.sh"
source "$COMMON_DIR/services/mapexos.sh"
source "$COMMON_DIR/services/nats.sh"

# ─── Commands ────────────────────────────────────────────────

cmd_setup() {
    echo ""
    echo "================================================================"
    echo "  MapexOS Benchmark Seed — Setup"
    echo "================================================================"
    echo "  MongoDB:       ${MONGO_DB}"
    echo "  Organizations: ${BENCH_ORG_COUNT}"
    echo "  Users:         ${BENCH_USER_COUNT}  (client{i}@test.com)"
    echo "  Groups:        ${BENCH_GROUP_COUNT}"
    echo "  Roles:         ${BENCH_ROLE_COUNT}"
    echo "  Memberships:   ${BENCH_MEMBERSHIP_COUNT}  (user) + ${BENCH_GROUP_COUNT} (group)"
    echo "  Group members: ${BENCH_USER_COUNT}  (1 per user → group)"
    echo "================================================================"
    echo ""

    # 1. Seed all 7 entity types into MongoDB
    log_step "1" "2" "Seeding MongoDB entities..."
    bench_seed_mapexos

    # 2. Flush Redis caches so benchmark starts clean
    echo ""
    log_step "2" "2" "Flushing Redis caches..."
    redis_flush 0
    redis_flush 5

    echo ""
    log_success "Setup complete."
    echo ""
    echo "  Login:     client1@test.com / ${BENCH_PASSWORD}"
    echo "  User ID:   ${BENCH_FIRST_USER_ID}"
    echo "  Org:       ${BENCH_FIRST_CUSTOMER_ORG_ID}"
    echo ""
}

cmd_teardown() {
    echo ""
    echo "================================================================"
    echo "  MapexOS Benchmark Seed — Teardown"
    echo "================================================================"
    echo ""

    # 1. Purge NATS streams (stop in-flight processing first)
    log_step "1" "2" "Purging NATS streams..."
    bench_purge_all_streams

    # 2. Cleanup MongoDB + Redis
    echo ""
    log_step "2" "2" "Cleaning MongoDB + Redis..."
    bench_cleanup_mapexos

    echo ""
    log_success "Teardown complete. All benchmark data removed."
    echo ""
}

# ─── Main ────────────────────────────────────────────────────

case "$COMMAND" in
    setup)
        cmd_setup
        ;;
    teardown)
        cmd_teardown
        ;;
    *)
        echo "Usage: $0 <setup|teardown>"
        echo ""
        echo "  setup      Seed orgs, roles, users, groups, memberships, group_members (idempotent)"
        echo "  teardown   Purge NATS + delete seed data + flush Redis"
        echo ""
        echo "  Seed data:"
        echo "    Users:   client{i}@test.com (password: test@123)"
        echo "    IDs:     zero-padded deterministic ObjectIds"
        echo ""
        echo "  Environment variables (all optional, have defaults):"
        echo "    GO_ENV           Default: dev"
        echo "    MONGO_DB         Default: \${GO_ENV}-mapexos"
        echo "    ORG_COUNT        Default: 10"
        echo "    USER_COUNT       Default: 1000"
        echo "    GROUP_COUNT      Default: 100"
        echo "    ROLE_COUNT       Default: 50"
        echo "    MEMBERSHIP_COUNT Default: 500"
        exit 1
        ;;
esac
