#!/bin/bash

# E2E Test Setup Script
# Clean environment + seed initial test data for E2E tests

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

echo "🚀 Starting E2E environment setup..."
echo ""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Step 1: Clean environment
echo -e "${BLUE}══════════════════════════════════════════════════${NC}"
echo -e "${BLUE}  STEP 1: Cleaning Environment${NC}"
echo -e "${BLUE}══════════════════════════════════════════════════${NC}"
echo ""

# Call e2e-cleanup.sh
if [ -f "$SCRIPT_DIR/e2e-cleanup.sh" ]; then
    bash "$SCRIPT_DIR/e2e-cleanup.sh"
else
    echo -e "${RED}  ✗ e2e-cleanup.sh not found${NC}"
    exit 1
fi

echo ""

# Step 2: Seed initial data
echo -e "${BLUE}══════════════════════════════════════════════════${NC}"
echo -e "${BLUE}  STEP 2: Seeding Test Data${NC}"
echo -e "${BLUE}══════════════════════════════════════════════════${NC}"
echo ""

# Seed data using mongosh
echo -e "${YELLOW}📦 Seeding MongoDB with test data...${NC}"

mongosh --quiet --eval "
const mapexosDb = db.getSiblingDB('mapexos');

// Fixed IDs for test data (matching E2E test expectations)
const mapexosOrgId = ObjectId('68f5bbce1aef22967c3ebb30');
const rootUserId = ObjectId('68f689755f27dc70c50a6322');
const adminRoleId = ObjectId('68f5bbce1aef22967c3ebb40');

// 1. Create Mapexos Organization (root organization for tests)
print('  ✓ Creating Mapexos organization...');
mapexosDb.organizations.insertOne({
    _id: mapexosOrgId,
    name: 'Mapexos',
    customerId: mapexosOrgId,  // Root customer
    orgId: null,               // Root org
    pathKey: '/',              // Root path
    scope: 'global',
    status: true,
    created: new Date(),
    updated: new Date(),
});

// 2. Create Admin Role (with all permissions)
print('  ✓ Creating Admin role...');
mapexosDb.roles.insertOne({
    _id: adminRoleId,
    name: 'Admin',
    description: 'System administrator with all permissions',
    permissions: [
        'users:create', 'users:read', 'users:update', 'users:delete', 'users:list',
        'roles:create', 'roles:read', 'roles:update', 'roles:delete', 'roles:list',
        'organizations:create', 'organizations:read', 'organizations:update', 'organizations:delete', 'organizations:list',
        'groups:create', 'groups:read', 'groups:update', 'groups:delete', 'groups:list',
        'memberships:create', 'memberships:read', 'memberships:update', 'memberships:delete', 'memberships:list',
        'assets:create', 'assets:read', 'assets:update', 'assets:delete', 'assets:list',
        'assettemplates:create', 'assettemplates:read', 'assettemplates:update', 'assettemplates:delete', 'assettemplates:list',
        'routegroups:create', 'routegroups:read', 'routegroups:update', 'routegroups:delete', 'routegroups:list',
        'datasources:create', 'datasources:read', 'datasources:update', 'datasources:delete', 'datasources:list',
    ],
    customerId: mapexosOrgId,
    orgId: null,               // Global role
    pathKey: '/',
    scope: 'global',
    isSystem: true,
    status: true,
    created: new Date(),
    updated: new Date(),
});

// 3. Create Root User (root@mapex.global)
print('  ✓ Creating root user...');
mapexosDb.users.insertOne({
    _id: rootUserId,
    name: 'Root User',
    email: 'root@mapex.global',
    password: '\$2a\$10\$rN.hOKUqH5FvZWJLZWKZCO5xVqF8xQxYxYxYxYxYxYxYxYxYxYx', // hashed: 'root123'
    status: true,
    isEmailVerified: true,
    created: new Date(),
    updated: new Date(),
});

// 4. Create Membership (root user in Mapexos org with Admin role)
print('  ✓ Creating root user membership...');
mapexosDb.memberships.insertOne({
    userId: rootUserId,
    orgId: mapexosOrgId,
    customerId: mapexosOrgId,
    pathKey: '/',
    roleIds: [adminRoleId],
    groupIds: [],
    status: true,
    created: new Date(),
    updated: new Date(),
});

print('');
print('✅ Test data seeded successfully!');
print('');
print('Seeded data:');
print('  • Organization: Mapexos (ID: 68f5bbce1aef22967c3ebb30)');
print('  • Role: Admin (ID: 68f5bbce1aef22967c3ebb40)');
print('  • User: root@mapex.global (ID: 68f689755f27dc70c50a6322)');
print('  • Membership: root user in Mapexos org with Admin role');
print('');
" && echo -e "${GREEN}  ✓ MongoDB seeding completed${NC}" || echo -e "${RED}  ✗ MongoDB seeding failed${NC}"

echo ""
echo -e "${GREEN}══════════════════════════════════════════════════${NC}"
echo -e "${GREEN}  ✅ E2E Setup Completed Successfully!${NC}"
echo -e "${GREEN}══════════════════════════════════════════════════${NC}"
echo ""
echo "Environment is ready for E2E tests:"
echo "  • All databases cleaned (MongoDB, Redis, NATS)"
echo "  • Test organization created: Mapexos (68f5bbce1aef22967c3ebb30)"
echo "  • Root user created: root@mapex.global (68f689755f27dc70c50a6322)"
echo "  • Admin role created with all permissions"
echo "  • Root user has Admin role in Mapexos org"
echo ""
echo "You can now run E2E tests!"
echo ""
