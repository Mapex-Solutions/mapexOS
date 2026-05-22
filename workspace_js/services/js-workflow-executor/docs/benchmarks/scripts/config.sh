#!/bin/bash
# =============================================================
# Service Configuration — JS-Workflow-Executor Benchmark
#
# Constants + common bootstrap. Source this first in any script:
#   source "$SCRIPT_DIR/config.sh"
#
# Then source the common services you need:
#   source "$COMMON_DIR/services/nats.sh"
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

# ─── JS-Workflow-Executor Constants ──────────────────────────

# Environment
GO_ENV_VALUE="${GO_ENV:-dev}"

# Primary stream and consumer
STREAM="WORKFLOW-JS-CODE"
CONSUMER="js-workflow-executor-code"
SUBJECT="workflow.js.code"

# All NATS streams to purge during cleanup
ALL_STREAMS=(
    "WORKFLOW-JS-CODE"
    "WORKFLOW-RESUME"
    "FANOUT"
    "MAPEXOS-DLQ"
)

# ─── Workflow Constants ──────────────────────────────────────
# Used by seed scripts for MinIO and NATS payloads

BENCH_ORG_ID="000000000000000000000001"
BENCH_WORKFLOW_ID="bench-workflow-0001"
BENCH_NODE_ID="bench-code-node-0001"
BENCH_INSTANCE_ID_PREFIX="bench-instance"

# MinIO
MINIO_ALIAS="${MINIO_ALIAS:-local}"
MINIO_BUCKET="${MINIO_WORKFLOWS_BUCKET:-mapex-workflows}"

# Script source path in MinIO: {orgId}/{workflowId}/scripts/{nodeId}.json
MINIO_SCRIPT_OBJECT="${BENCH_ORG_ID}/${BENCH_WORKFLOW_ID}/scripts/${BENCH_NODE_ID}.json"

# ─── MinIO Helpers (workflow-specific, not using common/services/minio.sh) ───

# Upload workflow script source to MinIO
bench_seed_workflow_minio() {
    local seed_file="${1:-$SEED_DIR/minio/workflow-script-source.js}"
    local object_path="${2:-$MINIO_SCRIPT_OBJECT}"

    if [ ! -f "$seed_file" ]; then
        log_warn "Seed file not found: $seed_file"
        return 1
    fi

    log_info "Uploading script source to MinIO: ${MINIO_BUCKET}/${object_path}..."
    mc cp "$seed_file" "${MINIO_ALIAS}/${MINIO_BUCKET}/${object_path}" 2>/dev/null \
        && log_success "Script uploaded: ${object_path}" \
        || log_warn "MinIO upload failed — check mc alias '${MINIO_ALIAS}'"
}

# Remove workflow script source from MinIO
bench_cleanup_workflow_minio() {
    local object_path="${1:-$MINIO_SCRIPT_OBJECT}"

    log_info "Removing script source from MinIO: ${object_path}..."
    mc rm "${MINIO_ALIAS}/${MINIO_BUCKET}/${object_path}" 2>/dev/null \
        && log_success "Script removed: ${object_path}" \
        || log_warn "MinIO removal failed (object may not exist)"
}

# Ensure NATS streams exist (create if missing)
bench_ensure_streams() {
    for stream_name in "${ALL_STREAMS[@]}"; do
        local exists
        exists=$(nats_cmd stream info "$stream_name" --json 2>/dev/null | grep -c '"name"' || true)
        if [ "${exists:-0}" -eq 0 ]; then
            local subject
            case "$stream_name" in
                WORKFLOW-JS-CODE)  subject="workflow.js.code" ;;
                WORKFLOW-RESUME)   subject="workflow.resume.>" ;;
                FANOUT)            subject="fanout.>" ;;
                MAPEXOS-DLQ)       subject="dlq.>" ;;
                *)                 subject="${stream_name,,}.>" ;;
            esac
            nats_cmd stream add "$stream_name" \
                --subjects "$subject" \
                --storage file \
                --retention limits \
                --max-msgs=-1 \
                --max-bytes=-1 \
                --max-age=24h \
                --max-msg-size=-1 \
                --discard old \
                --dupe-window 2m \
                --replicas 1 \
                --no-deny-delete \
                --no-deny-purge \
                --allow-rollup \
                2>/dev/null && log_info "Created stream $stream_name" || true
        fi
    done
}

# ─── Service ───────────────────────────────────────────────

SERVICE_PORT=8001
METRICS_URL="http://localhost:${SERVICE_PORT}/metrics"

# ─── Benchmark Isolation ─────────────────────────────────────
# No dependencies on other services — kill everything
BENCH_KEEP_PORTS=()

# ─── Preflight ──────────────────────────────────────────────

BENCH_CLI_TOOLS="nats,mc,curl"
BENCH_SERVICE_CHECKS="NATS:nats_cmd stream list --json,MinIO:mc ls ${MINIO_ALIAS:-local}"
