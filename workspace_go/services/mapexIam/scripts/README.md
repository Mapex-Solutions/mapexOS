# Migration: Group.members[] -> group_members Collection

This migration replaces the embedded `members[]` array in the `groups` collection with a separate `group_members` junction collection for BIGTECH-scale scalability (100K+ tenants).

## Prerequisites

- MongoDB 5.0+
- `mongodump` and `mongorestore` installed (for backup)
- `mongosh` installed (for migration script)

## Migration Steps

### Step 1: Backup Database

**ALWAYS backup before migration.**

```bash
./backup_before_migration.sh "mongodb://localhost:27017" "./backups"
```

Or with authentication:

```bash
./backup_before_migration.sh "mongodb://user:pass@host:port?authSource=admin" "./backups"
```

### Step 2: Run Migration (Dry Run)

First, run in dry-run mode to preview changes:

1. Edit `migrate_group_members.js` and ensure `DRY_RUN = true`
2. Run the script:

```bash
mongosh "mongodb://localhost:27017/mapexos" migrate_group_members.js
```

Or with authentication:

```bash
mongosh "mongodb://user:pass@host:port/mapexos?authSource=admin" migrate_group_members.js
```

### Step 3: Execute Migration

After verifying the dry run output:

1. Edit `migrate_group_members.js` and set `DRY_RUN = false`
2. Run the script again:

```bash
mongosh "mongodb://localhost:27017/mapexos" migrate_group_members.js
```

### Step 4: Verify Migration

The script will output verification results. Ensure:
- `group_members` collection has the expected number of documents
- `groups` collection has no documents with `members` field
- All 3 indexes exist on `group_members`:
  - `idx_group_user_unique` (unique)
  - `idx_user_org`
  - `idx_group`

### Rollback Procedure

If something goes wrong:

```bash
# Restore from backup
mongorestore --uri="mongodb://localhost:27017" --gzip "./backups/mapexos_backup_TIMESTAMP"

# Drop the partially migrated collection (optional)
mongosh "mongodb://localhost:27017/mapexos" --eval "db.group_members.drop()"
```

## Indexes Created

| Index Name | Fields | Unique | Purpose |
|------------|--------|--------|---------|
| `idx_group_user_unique` | `{ groupId: 1, userId: 1 }` | Yes | Prevent duplicate memberships |
| `idx_user_org` | `{ userId: 1, orgId: 1 }` | No | Query: user's groups in org |
| `idx_group` | `{ groupId: 1 }` | No | Query: all members of a group |

## Document Schema

### Before (embedded in groups)

```json
{
  "_id": ObjectId("..."),
  "name": "Group Name",
  "orgId": ObjectId("..."),
  "members": [ObjectId("user1"), ObjectId("user2")],
  "created": ISODate("..."),
  "updated": ISODate("...")
}
```

### After (separate collection)

**groups collection:**
```json
{
  "_id": ObjectId("..."),
  "name": "Group Name",
  "orgId": ObjectId("..."),
  "created": ISODate("..."),
  "updated": ISODate("...")
}
```

**group_members collection:**
```json
{
  "_id": ObjectId("..."),
  "groupId": ObjectId("..."),
  "userId": ObjectId("..."),
  "orgId": ObjectId("..."),
  "addedAt": ISODate("..."),
  "addedBy": ObjectId("...") | null,
  "created": ISODate("..."),
  "updated": ISODate("...")
}
```
