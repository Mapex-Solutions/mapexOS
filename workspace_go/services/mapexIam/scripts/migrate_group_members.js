/**
 * Migration Script: Group.members[] -> group_members collection
 *
 * This script migrates the embedded members array from the groups collection
 * to a separate junction collection for scalability (100K+ tenants).
 *
 * Usage:
 *   mongosh "mongodb://localhost:27017/mapexos" migrate_group_members.js
 *
 * Or with authentication:
 *   mongosh "mongodb://user:pass@host:port/mapexos?authSource=admin" migrate_group_members.js
 */

// Configuration
const DRY_RUN = false; // Set to true to preview changes without modifying data
const BATCH_SIZE = 100;

print("=".repeat(60));
print("[MIGRATION] Group.members[] -> group_members collection");
print("=".repeat(60));
print(`[CONFIG] Dry run: ${DRY_RUN}`);
print(`[CONFIG] Batch size: ${BATCH_SIZE}`);
print("");

// Step 1: Backup verification
print("[STEP 1] Verifying collections...");

const groupsCount = db.groups.countDocuments();
const existingMembersCount = db.group_members.countDocuments();

print(`  - groups collection: ${groupsCount} documents`);
print(`  - group_members collection: ${existingMembersCount} documents`);

if (existingMembersCount > 0) {
  print("[WARNING] group_members collection already has data!");
  print("  This migration may create duplicates if run again.");
  print("  Consider dropping the collection first if re-running.");
  print("");
}

// Step 2: Count groups with members to migrate
print("[STEP 2] Analyzing groups with members...");

const groupsWithMembers = db.groups.find({
  members: { $exists: true, $ne: [], $type: "array" }
}).toArray();

print(`  - Groups with members to migrate: ${groupsWithMembers.length}`);

let totalMembersToMigrate = 0;
groupsWithMembers.forEach(group => {
  if (group.members && Array.isArray(group.members)) {
    totalMembersToMigrate += group.members.length;
  }
});

print(`  - Total member records to create: ${totalMembersToMigrate}`);
print("");

if (groupsWithMembers.length === 0) {
  print("[INFO] No groups with members found. Nothing to migrate.");
  print("[DONE] Migration complete.");
  quit(0);
}

// Step 3: Create indexes on group_members collection
print("[STEP 3] Creating indexes on group_members collection...");

if (!DRY_RUN) {
  // Unique compound index - prevents duplicate memberships
  db.group_members.createIndex(
    { groupId: 1, userId: 1 },
    { unique: true, name: "idx_group_user_unique" }
  );

  // Index for "user's groups in org" queries
  db.group_members.createIndex(
    { userId: 1, orgId: 1 },
    { name: "idx_user_org" }
  );

  // Index for "group members" queries
  db.group_members.createIndex(
    { groupId: 1 },
    { name: "idx_group" }
  );

  print("  - Created idx_group_user_unique (unique)");
  print("  - Created idx_user_org");
  print("  - Created idx_group");
} else {
  print("  - [DRY RUN] Would create 3 indexes");
}
print("");

// Step 4: Migrate members data
print("[STEP 4] Migrating members data...");

let migratedCount = 0;
let skippedCount = 0;
let errorCount = 0;

groupsWithMembers.forEach((group, index) => {
  if (!group.members || !Array.isArray(group.members)) {
    return;
  }

  const groupId = group._id;
  const orgId = group.orgId;
  const groupCreated = group.created || new Date();

  group.members.forEach(userId => {
    const memberDoc = {
      groupId: groupId,
      userId: userId,
      orgId: orgId,
      addedAt: groupCreated, // Use group creation as fallback
      addedBy: null, // Historical data - no audit trail
      created: new Date(),
      updated: new Date()
    };

    if (!DRY_RUN) {
      try {
        db.group_members.insertOne(memberDoc);
        migratedCount++;
      } catch (e) {
        if (e.code === 11000) {
          // Duplicate key - member already exists
          skippedCount++;
        } else {
          print(`  [ERROR] Failed to migrate member ${userId} in group ${groupId}: ${e.message}`);
          errorCount++;
        }
      }
    } else {
      migratedCount++;
    }
  });

  // Progress indicator
  if ((index + 1) % BATCH_SIZE === 0 || index === groupsWithMembers.length - 1) {
    print(`  - Processed ${index + 1}/${groupsWithMembers.length} groups...`);
  }
});

print("");
print(`  - Migrated: ${migratedCount} member records`);
print(`  - Skipped (duplicates): ${skippedCount}`);
print(`  - Errors: ${errorCount}`);
print("");

// Step 5: Remove members field from groups collection
print("[STEP 5] Removing members field from groups collection...");

if (!DRY_RUN) {
  const updateResult = db.groups.updateMany(
    { members: { $exists: true } },
    { $unset: { members: "" } }
  );
  print(`  - Updated ${updateResult.modifiedCount} documents`);
} else {
  const docsWithMembers = db.groups.countDocuments({ members: { $exists: true } });
  print(`  - [DRY RUN] Would update ${docsWithMembers} documents`);
}
print("");

// Step 6: Verification
print("[STEP 6] Verification...");

const finalGroupMembersCount = db.group_members.countDocuments();
const groupsStillWithMembers = db.groups.countDocuments({ members: { $exists: true } });

print(`  - group_members collection: ${finalGroupMembersCount} documents`);
print(`  - groups with members field: ${groupsStillWithMembers}`);

if (!DRY_RUN && groupsStillWithMembers === 0 && errorCount === 0) {
  print("");
  print("[SUCCESS] Migration completed successfully!");
} else if (DRY_RUN) {
  print("");
  print("[DRY RUN] Preview complete. Set DRY_RUN = false to execute migration.");
} else {
  print("");
  print("[WARNING] Migration completed with issues. Please review the output.");
}

print("=".repeat(60));
