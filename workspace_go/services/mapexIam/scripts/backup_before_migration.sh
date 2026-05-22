#!/bin/bash

# Backup Script: Pre-migration database backup
#
# This script creates a backup of the mapexos database before running
# the group_members migration.
#
# Usage:
#   ./backup_before_migration.sh [mongodb_uri] [backup_dir]
#
# Example:
#   ./backup_before_migration.sh "mongodb://localhost:27017" "./backups"

set -e

# Configuration
MONGODB_URI="${1:-mongodb://localhost:27017}"
BACKUP_DIR="${2:-./backups}"
DATABASE="mapexos"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
BACKUP_PATH="${BACKUP_DIR}/mapexos_backup_${TIMESTAMP}"

echo "============================================================"
echo "[BACKUP] Pre-migration database backup"
echo "============================================================"
echo "MongoDB URI: ${MONGODB_URI}"
echo "Database: ${DATABASE}"
echo "Backup path: ${BACKUP_PATH}"
echo ""

# Create backup directory
mkdir -p "${BACKUP_DIR}"

# Run mongodump
echo "[STEP 1] Running mongodump..."
mongodump \
  --uri="${MONGODB_URI}" \
  --db="${DATABASE}" \
  --out="${BACKUP_PATH}" \
  --gzip

echo ""
echo "[STEP 2] Verifying backup..."
if [ -d "${BACKUP_PATH}/${DATABASE}" ]; then
  echo "  - Backup directory created successfully"
  echo "  - Contents:"
  ls -la "${BACKUP_PATH}/${DATABASE}/"
else
  echo "[ERROR] Backup directory not found!"
  exit 1
fi

echo ""
echo "[STEP 3] Creating collections summary..."
SUMMARY_FILE="${BACKUP_PATH}/backup_summary.txt"
{
  echo "Backup Summary"
  echo "=============="
  echo "Date: $(date)"
  echo "Database: ${DATABASE}"
  echo "MongoDB URI: ${MONGODB_URI}"
  echo ""
  echo "Collections backed up:"
  ls -1 "${BACKUP_PATH}/${DATABASE}/"
} > "${SUMMARY_FILE}"

echo "  - Summary saved to: ${SUMMARY_FILE}"

echo ""
echo "============================================================"
echo "[SUCCESS] Backup completed successfully!"
echo ""
echo "To restore this backup, run:"
echo "  mongorestore --uri=\"${MONGODB_URI}\" --gzip \"${BACKUP_PATH}\""
echo "============================================================"
