#!/bin/bash
# Seed data for workflow benchmarks
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
source "$SCRIPT_DIR/config.sh"

echo "=== Seeding workflow benchmark data ==="

# Seed MongoDB with test workflow definitions
echo "[1/2] Seeding MongoDB..."
for f in "$SCRIPT_DIR/../seed/mongodb/"*.json; do
  [ -f "$f" ] || continue
  echo "  → $(basename "$f")"
  mongoimport --uri="$MONGO_URI" --db="$MONGO_DATABASE" --collection="workflows" --file="$f" --mode=upsert
done

# Seed NATS with test trigger messages
echo "[2/2] Seeding NATS..."
for f in "$SCRIPT_DIR/../seed/nats/"*.json; do
  [ -f "$f" ] || continue
  echo "  → $(basename "$f")"
  nats pub --server="$NATS_URL" --user="$NATS_USERNAME" --password="$NATS_PASSWORD" \
    "workflow.trigger.benchmark" "$(cat "$f")"
done

echo "=== Seeding complete ==="
