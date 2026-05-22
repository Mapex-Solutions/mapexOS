#!/bin/bash
# =============================================================
# Seed Orchestrator — Assets Benchmark
#
# Manages the complete test data lifecycle:
#   setup       — Seed template + 1000 HTTP assets + N MQTT assets
#                 into MongoDB, upload MinIO read models, seed SharedCache
#   teardown    — Purge NATS + delete MongoDB + flush Redis + clean SharedCache + remove MinIO
#
# Redis auth cache is NOT seeded manually — the service populates it
# via AppCache.SetEx() during warmup requests. This guarantees the
# cache format (key prefix + JSON serialization) matches exactly.
#
# Uses modular functions from config + common modules:
#   config.sh              — Constants, common bootstrap
#   common/services/nats.sh — Stream purge
#
# MQTT asset count is controlled by BENCH_MQTT_ASSET_COUNT env var
# (default 10,000 — set in config.sh).
#
# Usage:
#   ./seed.sh setup        — Insert benchmark data (idempotent)
#   ./seed.sh teardown     — Clean all benchmark data
#
# Prerequisites:
#   - MongoDB running at localhost:27017 (replica set: rs0)
#   - Redis running at localhost:6379
#   - NATS running at localhost:4222
#   - MinIO accessible via mc alias "local"
#   - CLIs: mongosh, redis-cli, nats, mc
# =============================================================

set -euo pipefail

COMMAND="${1:-}"

# ─── Source Modules ──────────────────────────────────────────

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"

source "$SCRIPT_DIR/config.sh"
source "$COMMON_DIR/services/nats.sh"

# ─── Commands ────────────────────────────────────────────────

cmd_setup() {
    echo ""
    echo "================================================================"
    echo "  Assets Benchmark — Seed Setup"
    echo "================================================================"
    echo "  MongoDB:      ${MONGO_DB}"
    echo "  HTTP assets:  ${SEED_ASSET_COUNT} documents (org=${BENCH_ORG_ID})"
    echo "  MQTT assets:  ${BENCH_MQTT_ASSET_COUNT} documents (org=${BENCH_MQTT_ORG_ID})"
    echo "  Template:     ID ${BENCH_TEMPLATE_ID}"
    echo "  HTTP Org:     ID ${BENCH_ORG_ID} pathKey=${BENCH_ORG_PATHKEY}"
    echo "  MQTT Org:     ID ${BENCH_MQTT_ORG_ID} pathKey=${BENCH_MQTT_ORG_PATHKEY}"
    echo "================================================================"
    echo ""

    # 1. MongoDB: insert benchmark asset template
    log_step "1" "5" "Inserting benchmark asset template..."
    mongo_eval "$MONGO_DB" "
        var templateId = ObjectId('${BENCH_TEMPLATE_ID}');
        if (!db.assets_templates.findOne({ _id: templateId })) {
            db.assets_templates.insertOne({
                _id: templateId,
                name: 'bench-template',
                enabled: true,
                description: 'Benchmark asset template — do not use in production',
                categoryId: null,
                categoryName: 'Benchmark',
                manufacturerId: null,
                manufacturerName: 'BenchCorp',
                modelId: null,
                modelName: 'BenchModel-v1',
                version: '1.0.0',
                assetIdPath: 'metadata.deviceId',
                scriptConversion: 'return payload;',
                scriptValidator: 'return true;',
                scriptProcessor: null,
                scriptTest: null,
                availableFields: ['data.temperature', 'data.humidity', 'metadata.deviceId'],
                dynamicFields: [
                    { fieldId: 1, field: 'temperature', value: 'data.temperature', type: 'number', status: 1 },
                    { fieldId: 2, field: 'humidity',    value: 'data.humidity',    type: 'number', status: 1 }
                ],
                nextFieldId: 3,
                isSystem: false,
                isTemplate: false,
                orgId: ObjectId('${BENCH_ORG_ID}'),
                pathKey: '${BENCH_ORG_PATHKEY}',
                created: new Date(),
                updated: new Date()
            });
            print('Inserted bench-template');
        } else {
            print('bench-template already exists, skipping');
        }
    "
    log_success "Template ready."

    # 2. MongoDB: insert 1000 HTTP benchmark assets
    echo ""
    log_step "2" "5" "Inserting ${SEED_ASSET_COUNT} HTTP benchmark assets..."
    mongo_eval "$MONGO_DB" "
        var orgId       = ObjectId('${BENCH_ORG_ID}');
        var templateId  = ObjectId('${BENCH_TEMPLATE_ID}');
        var customerId  = ObjectId('${BENCH_CUSTOMER_ID}');
        var now         = new Date();
        var count       = ${SEED_ASSET_COUNT};

        var docs = Array.from({ length: count }, function(_, i) {
            var n = i + 1;
            var id = ObjectId((n + 256).toString(16).padStart(24, '0'));
            return {
                _id: id,
                name: 'bench-asset-' + n,
                enabled: true,
                debugEnabled: false,
                description: 'Benchmark asset ' + n,
                assetUUID: 'bench-device-' + n.toString().padStart(6, '0'),
                assetTemplateId: templateId,
                orgId: orgId,
                pathKey: '${BENCH_ORG_PATHKEY}',
                customerId: customerId,
                routeGroupIds: [],
                protocol: { type: 'http', http: {} },
                latitude: null,
                longitude: null,
                created: now,
                updated: now
            };
        });

        try {
            var result = db.assets.insertMany(docs, { ordered: false });
            print('Inserted ' + Object.keys(result.insertedIds).length + ' assets');
        } catch (e) {
            if (e.code === 11000 || (e.writeErrors && e.writeErrors.length > 0)) {
                print('Some assets already exist (idempotent re-run), continuing');
            } else {
                throw e;
            }
        }
    "
    log_success "HTTP assets ready."

    # 3. MongoDB: insert MQTT benchmark assets (for auth callout)
    echo ""
    log_step "3" "5" "Inserting ${BENCH_MQTT_ASSET_COUNT} MQTT benchmark assets..."
    mongo_eval "$MONGO_DB" "
        var __COUNT__ = ${BENCH_MQTT_ASSET_COUNT};
        var __START_HEX__ = ${BENCH_MQTT_ASSET_START};
        var __ORG_ID__ = '${BENCH_MQTT_ORG_ID}';
        var __TEMPLATE_ID__ = '${BENCH_TEMPLATE_ID}';
        var __CUSTOMER_ID__ = '${BENCH_MQTT_CUSTOMER_ID}';
        var __PATHKEY__ = '${BENCH_MQTT_ORG_PATHKEY}';
        $(cat "$SEED_DIR/mongodb/mqtt-assets.js")
    "
    log_success "MQTT assets ready."

    # 4. MinIO: upload asset read models (first 100 HTTP assets)
    echo ""
    log_step "4" "5" "Uploading AssetReadModel objects to MinIO (first 100 assets)..."
    for i in $(seq 1 100); do
        UUID="bench-device-$(printf '%06d' "$i")"
        ASSET_ID=$(printf '%024x' $((i + 256)))
        READ_MODEL="{\"id\":\"${ASSET_ID}\",\"uuid\":\"${UUID}\",\"orgId\":\"${BENCH_ORG_ID}\",\"pathKey\":\"${BENCH_ORG_PATHKEY}\",\"enabled\":true,\"debugEnabled\":false,\"name\":\"bench-asset-${i}\",\"assetTemplateId\":\"${BENCH_TEMPLATE_ID}\",\"routeGroupIds\":[]}"
        TMPFILE=$(mktemp /tmp/bench-asset-XXXXXX.json)
        echo "$READ_MODEL" > "$TMPFILE"
        if command -v mc &> /dev/null; then
            mc cp "$TMPFILE" "${MINIO_ALIAS}/${MINIO_ASSETS_BUCKET}/${BENCH_ORG_ID}/${UUID}.json" --quiet > /dev/null 2>&1 || true
        fi
        rm -f "$TMPFILE"
    done
    log_success "MinIO upload complete."

    # 5. SharedCache (Redis DB 5): seed permission + coverage for benchmark user
    echo ""
    log_step "5" "5" "Seeding SharedCache (Redis DB ${SHARED_CACHE_DB}) — permissions + coverage..."

    # Permission: HTTP org-scoped (for X-Org-Context requests)
    redis_cmd -n "$SHARED_CACHE_DB" SET \
        "${SHARED_CACHE_PREFIX}:auth:org:${BENCH_ORG_ID}:user:${BENCH_USER_ID}:ver" "1" > /dev/null
    redis_cmd -n "$SHARED_CACHE_DB" SET \
        "${SHARED_CACHE_PREFIX}:auth:org:${BENCH_ORG_ID}:user:${BENCH_USER_ID}:v1" '["mapex.*"]' > /dev/null

    # Permission: MQTT org-scoped (for auth callout requests)
    redis_cmd -n "$SHARED_CACHE_DB" SET \
        "${SHARED_CACHE_PREFIX}:auth:org:${BENCH_MQTT_ORG_ID}:user:${BENCH_USER_ID}:ver" "1" > /dev/null
    redis_cmd -n "$SHARED_CACHE_DB" SET \
        "${SHARED_CACHE_PREFIX}:auth:org:${BENCH_MQTT_ORG_ID}:user:${BENCH_USER_ID}:v1" '["mapex.*"]' > /dev/null

    # Permission: global (for requests without org context)
    redis_cmd -n "$SHARED_CACHE_DB" SET \
        "${SHARED_CACHE_PREFIX}:auth:org:global:user:${BENCH_USER_ID}:ver" "1" > /dev/null
    redis_cmd -n "$SHARED_CACHE_DB" SET \
        "${SHARED_CACHE_PREFIX}:auth:org:global:user:${BENCH_USER_ID}:v1" '["mapex.*"]' > /dev/null

    # Coverage: user org access — includes both HTTP and MQTT orgs
    local coverage_json
    coverage_json="{\"userId\":\"${BENCH_USER_ID}\",\"accessibleOrgIds\":[\"${BENCH_ORG_ID}\",\"${BENCH_MQTT_ORG_ID}\"],\"organizations\":[{\"id\":\"${BENCH_ORG_ID}\",\"name\":\"Benchmark Org\",\"type\":\"vendor\",\"pathKey\":\"${BENCH_ORG_PATHKEY}\",\"scope\":\"local\",\"membershipId\":\"${BENCH_USER_ID}\",\"roleIds\":[\"admin\"]},{\"id\":\"${BENCH_MQTT_ORG_ID}\",\"name\":\"Benchmark MQTT Org\",\"type\":\"vendor\",\"pathKey\":\"${BENCH_MQTT_ORG_PATHKEY}\",\"scope\":\"local\",\"membershipId\":\"${BENCH_USER_ID}\",\"roleIds\":[\"admin\"]}],\"lastUpdated\":\"2025-01-01T00:00:00Z\",\"version\":1}"
    redis_cmd -n "$SHARED_CACHE_DB" SET \
        "${SHARED_CACHE_PREFIX}:coverage:user:${BENCH_USER_ID}" "$coverage_json" > /dev/null

    log_success "SharedCache seeded (user=${BENCH_USER_ID}, httpOrg=${BENCH_ORG_ID}, mqttOrg=${BENCH_MQTT_ORG_ID})."

    echo ""
    log_success "Seed complete. Run './seed.sh teardown' when done."
    echo ""
}

cmd_teardown() {
    echo ""
    echo "================================================================"
    echo "  Assets Benchmark — Seed Teardown"
    echo "================================================================"
    echo ""

    # 1. NATS: purge all streams
    log_step "1" "5" "Purging NATS streams..."
    bench_purge_all_streams

    # 2. MongoDB: delete seed documents
    # Seed IDs are zero-padded (000000...) which are always less than real ObjectIds
    # (real ObjectIds start with a Unix timestamp, e.g., 0x67... for 2025+).
    # $regex does NOT work on ObjectId fields, so we use $lt range instead.
    echo ""
    log_step "2" "5" "Deleting seed documents from MongoDB..."
    mongo_eval "$MONGO_DB" "
        var cutoff = ObjectId('000001000000000000000000');
        var assetResult = db.assets.deleteMany({ _id: { \$lt: cutoff } });
        print('Deleted assets: ' + assetResult.deletedCount);

        var templateResult = db.assets_templates.deleteMany({ _id: { \$lt: cutoff } });
        print('Deleted templates: ' + templateResult.deletedCount);
    "
    log_success "MongoDB seed documents deleted."

    # 3. Redis: flush AppCache (DB 0)
    echo ""
    log_step "3" "5" "Flushing Redis AppCache..."
    redis_flush 0

    # 4. SharedCache (Redis DB 5): delete benchmark keys (NOT flushdb — shared with other services)
    echo ""
    log_step "4" "5" "Removing benchmark keys from SharedCache (Redis DB ${SHARED_CACHE_DB})..."
    redis_cmd -n "$SHARED_CACHE_DB" DEL \
        "${SHARED_CACHE_PREFIX}:auth:org:${BENCH_ORG_ID}:user:${BENCH_USER_ID}:ver" \
        "${SHARED_CACHE_PREFIX}:auth:org:${BENCH_ORG_ID}:user:${BENCH_USER_ID}:v1" \
        "${SHARED_CACHE_PREFIX}:auth:org:${BENCH_MQTT_ORG_ID}:user:${BENCH_USER_ID}:ver" \
        "${SHARED_CACHE_PREFIX}:auth:org:${BENCH_MQTT_ORG_ID}:user:${BENCH_USER_ID}:v1" \
        "${SHARED_CACHE_PREFIX}:auth:org:global:user:${BENCH_USER_ID}:ver" \
        "${SHARED_CACHE_PREFIX}:auth:org:global:user:${BENCH_USER_ID}:v1" \
        "${SHARED_CACHE_PREFIX}:coverage:user:${BENCH_USER_ID}" > /dev/null 2>&1 || true
    log_success "SharedCache benchmark keys removed."

    # 5. MinIO: remove seed objects
    echo ""
    log_step "5" "5" "Removing seed objects from MinIO..."
    if command -v mc &> /dev/null; then
        mc rm --recursive --force --quiet \
            "${MINIO_ALIAS}/${MINIO_ASSETS_BUCKET}/${BENCH_ORG_ID}/" > /dev/null 2>&1 || {
            log_warn "Could not remove MinIO objects (continuing)"
        }
        log_success "MinIO objects removed."
    else
        log_warn "mc not installed, skipping MinIO cleanup."
    fi

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
        echo "  setup        Seed MongoDB (template + HTTP + MQTT assets),"
        echo "               upload MinIO read models, seed SharedCache."
        echo ""
        echo "  teardown     Purge NATS + delete MongoDB + flush Redis + clean SharedCache + remove MinIO."
        echo ""
        echo "  NOTE: Redis auth cache is populated by the service itself during"
        echo "        benchmark warmup — NOT seeded manually."
        echo ""
        echo "  Environment variables (all optional, have defaults):"
        echo "    GO_ENV                  Default: dev"
        echo "    MONGO_DB                Default: \${GO_ENV}-assets"
        echo "    MINIO_ALIAS             Default: local"
        echo "    BENCH_MQTT_ASSET_COUNT  Default: 10000 (MQTT devices for auth callout)"
        exit 1
        ;;
esac
