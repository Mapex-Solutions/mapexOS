#!/bin/bash
# =============================================================
# Seed Orchestrator — Events Benchmark
#
# Manages the complete test data lifecycle:
#   setup       — Create all 7 ClickHouse tables, seed MongoDB
#                 retention policy, flush Redis, purge NATS
#   teardown    — Purge NATS + delete MongoDB + flush Redis
#                 (ClickHouse TTL handles automatic cleanup)
#
# Uses modular functions from config + common modules:
#   config.sh              — Constants, common bootstrap, CH helpers
#   common/services/nats.sh — Stream purge
#
# Usage:
#   ./seed.sh setup        — Insert benchmark data (idempotent)
#   ./seed.sh teardown     — Clean all benchmark data
#
# Prerequisites:
#   - ClickHouse running at localhost:9000
#   - MongoDB running at localhost:27017 (replica set: rs0)
#   - Redis running at localhost:6379
#   - NATS running at localhost:4222
#   - CLIs: clickhouse-client, mongosh, redis-cli, nats
# =============================================================

set -euo pipefail

COMMAND="${1:-}"

# ─── Source Modules ──────────────────────────────────────────

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"

source "$SCRIPT_DIR/config.sh"
source "$COMMON_DIR/services/nats.sh"

# ─── ClickHouse Table Setup ──────────────────────────────────

# ensure_ch_table creates a table IF NOT EXISTS using the SQL seed file.
ensure_ch_table() {
    local table_name="$1"
    local sql_file="$SEEDS_SQL_DIR/${table_name}.sql"

    if [ ! -f "$sql_file" ]; then
        log_warn "SQL seed not found: $sql_file — skipping $table_name"
        return
    fi

    # Extract CREATE TABLE block and make idempotent (IF NOT EXISTS)
    local create_sql
    create_sql=$(sed -n '/^CREATE TABLE /,/;$/p' "$sql_file" | sed 's/^CREATE TABLE /CREATE TABLE IF NOT EXISTS /')

    if [ -n "$create_sql" ]; then
        ch_query "$create_sql" && log_success "$table_name table verified" \
            || log_warn "Could not create $table_name (may already exist with different schema)"
    else
        log_warn "No CREATE TABLE found in $sql_file"
    fi
}

# ─── Commands ────────────────────────────────────────────────

cmd_setup() {
    echo ""
    echo "================================================================"
    echo "  Events Benchmark — Seed Setup"
    echo "================================================================"
    echo "  ClickHouse:    ${CH_DB} (all 7 event tables)"
    echo "  MongoDB:       ${MONGO_DB} (retention policy)"
    echo "  MinIO:         ${MINIO_TEMPLATES_BUCKET} (benchmark template)"
    echo "  Benchmark org: ${BENCH_ORG_ID}"
    echo "================================================================"
    echo ""

    # 1. ClickHouse — ensure all 7 tables exist
    log_step "1" "5" "ClickHouse: verifying all event tables..."
    if ch_query "SELECT 1" > /dev/null 2>&1; then
        log_success "ClickHouse connection OK"

        local tables=("events_raw" "events" "events_jsexecutor" "events_router" "events_businessrule" "events_trigger" "events_dlq")
        for table in "${tables[@]}"; do
            ensure_ch_table "$table"
        done
    else
        log_warn "ClickHouse not reachable. Events will not persist — metrics still work."
    fi

    # 2. MinIO — upload benchmark template for TieredCache (save_event EVA resolution)
    echo ""
    log_step "2" "5" "MinIO: uploading benchmark template..."
    if command -v mc &> /dev/null; then
        mc cp "$SEED_DIR/payloads/bench-template.json" \
            "${MINIO_ALIAS}/${MINIO_TEMPLATES_BUCKET}/${MINIO_TEMPLATE_PATH}" \
            --quiet 2>/dev/null \
            && log_success "Template uploaded: ${MINIO_TEMPLATE_PATH}" \
            || log_warn "MinIO upload failed — save_event will use fallback HTTP."
    else
        log_warn "mc CLI not found — skipping template upload."
    fi

    # 3. MongoDB — insert retention policy for the benchmark org
    echo ""
    log_step "3" "5" "MongoDB: upserting retention policy for org ${BENCH_ORG_ID}..."
    mongo_eval "$MONGO_DB" "
        db.retentionpolicies.updateOne(
            { _id: ObjectId('${BENCH_ORG_ID}') },
            {
                \$setOnInsert: {
                    _id: ObjectId('${BENCH_ORG_ID}'),
                    orgId: '${BENCH_ORG_ID}',
                    createdAt: new Date(),
                },
                \$set: {
                    updatedAt: new Date(),
                    policies: [
                        { table: 'eventsRaw',          retentionDays: 1  },
                        { table: 'events',             retentionDays: 30 },
                        { table: 'eventsJsExecutor',   retentionDays: 1  },
                        { table: 'eventsDLQ',          retentionDays: 30 },
                        { table: 'eventsRouter',       retentionDays: 7  },
                        { table: 'eventsBusinessRule', retentionDays: 7  },
                        { table: 'eventsTrigger',      retentionDays: 7  },
                    ],
                }
            },
            { upsert: true }
        );
        print('retention policy upserted');
    "
    log_success "MongoDB retention policy ready."

    # 4. Redis — flush app cache (rebuilds on miss)
    echo ""
    log_step "4" "5" "Redis: flushing app cache..."
    redis_flush 0

    # 5. NATS — purge all event streams
    echo ""
    log_step "5" "5" "NATS: purging all event streams..."
    bench_purge_all_streams

    echo ""
    log_success "Seed complete."
    echo ""
    echo "  Next step: run the benchmark"
    echo "    bash full-benchmark.sh [scenario] [message_count] [cpu_list] [nats_batch_size]"
    echo ""
    echo "  Available scenarios: $ALL_SCENARIOS"
    echo "  Or use 'all' to run all scenarios sequentially."
    echo ""
}

cmd_teardown() {
    echo ""
    echo "================================================================"
    echo "  Events Benchmark — Seed Teardown"
    echo "================================================================"
    echo ""

    # 1. NATS — purge all streams
    log_step "1" "4" "Purging NATS streams..."
    bench_purge_all_streams

    # 2. MinIO — remove benchmark template
    echo ""
    log_step "2" "4" "Removing benchmark template from MinIO..."
    if command -v mc &> /dev/null; then
        mc rm "${MINIO_ALIAS}/${MINIO_TEMPLATES_BUCKET}/${MINIO_TEMPLATE_PATH}" \
            --quiet 2>/dev/null \
            && log_success "Template removed." \
            || log_warn "Template not found (may already be removed)."
    else
        log_warn "mc CLI not found — skipping."
    fi

    # 3. MongoDB — remove seed retention policy
    echo ""
    log_step "3" "4" "Deleting seed documents from MongoDB..."
    mongo_eval "$MONGO_DB" "
        var result = db.retentionpolicies.deleteOne({ _id: ObjectId('${BENCH_ORG_ID}') });
        print('deleted: ' + result.deletedCount);
    "
    log_success "MongoDB seed documents deleted."

    # 4. Redis — flush app cache
    echo ""
    log_step "4" "4" "Flushing Redis AppCache..."
    redis_flush 0

    echo ""
    log_info "ClickHouse data will expire automatically via TTL (retention_days)."
    log_info "To force immediate cleanup: OPTIMIZE TABLE ${CH_DB}.<table> FINAL"
    echo ""
    log_success "Teardown complete."
    echo ""
}

# ─── Main ────────────────────────────────────────────────────

case "$COMMAND" in
    setup)
        cmd_setup
        ;;
    teardown|down|clean)
        cmd_teardown
        ;;
    *)
        echo "Usage: $0 <setup|teardown>"
        echo ""
        echo "  setup        Create all 7 ClickHouse tables, seed MongoDB retention"
        echo "               policy, flush Redis, purge NATS streams."
        echo ""
        echo "  teardown     Purge NATS streams, delete MongoDB seed docs, flush Redis."
        echo "               ClickHouse data expires via TTL automatically."
        echo ""
        echo "  Environment variables (all optional, have defaults):"
        echo "    GO_ENV                  Default: dev"
        echo "    MONGO_DB                Default: events"
        echo "    CLICKHOUSE_HOST         Default: localhost"
        echo "    CLICKHOUSE_PORT         Default: 9000"
        echo "    CLICKHOUSE_DATABASE     Default: mapexos"
        exit 1
        ;;
esac
