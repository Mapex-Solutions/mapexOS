# E2E Tests - MapexOS

End-to-end tests for MapexOS microservices.

## Prerequisites

1. **Running services:**
   - **mapexos: `http://localhost:5000`** (required - includes auth endpoints)
   - router: `http://localhost:5003` (optional for some tests)
   - assets: `http://localhost:5002` (optional for some tests)
   - http_gateway: `http://localhost:5001` (optional for some tests)

2. **Databases:**
   - MongoDB: `mongodb://localhost:27017`
   - Redis: `localhost:6379`

3. **Admin user created:**
   - Email: `admin@mapex.global`
   - Password: `mapex123`
   - Tests perform real login with this user via `/auth/login`

4. **Start services:**
   ```bash
   cd workspace_go
   make run
   ```

## Structure

```
e2eTests/
├── common/                    # Shared code
│   ├── constants/            # URLs, timeouts, configs
│   ├── handlers/             # HTTPClient
│   ├── types/                # Common types (StandardResponse)
│   └── utils/                # JWT generation, assertions
│
├── services/                  # Tests per service/module
│   ├── mapexos/
│   │   ├── organizations/    # ✅ CRUD + Hierarchy
│   │   ├── roles/            # ✅ CRUD + Permissions
│   │   ├── groups/           # ✅ CRUD
│   │   ├── users/            # ✅ CRUD
│   │   └── memberships/      # 🚧 TODO
│   ├── assets/
│   │   ├── assets/           # 🚧 TODO
│   │   └── assettemplates/   # 🚧 TODO
│   ├── router/
│   │   └── routegroups/      # 🚧 TODO
│   └── http_gateway/
│       └── datasources/      # 🚧 TODO
│
└── journeys/                  # Integration tests (multiple services)
    ├── user_onboarding/      # 🚧 TODO
    ├── authorization_flow/   # 🚧 TODO
    └── permission_validation/# 🚧 TODO
```

## How to Run

### 🚀 Recommended Method: Using the Helper Script

The `run-tests.sh` script provides a friendly interface for running tests:

```bash
cd workspace_go/packages/e2eTests

# ============================================
# NEW FORMAT: SERVICE MODULE
# ============================================

# Run tests for a specific module
./run-tests.sh mapexos organizations    # Organizations from mapexos
./run-tests.sh mapexos roles            # Roles from mapexos
./run-tests.sh mapexos users            # Users from mapexos
./run-tests.sh assets assets            # Assets from assets service
./run-tests.sh router routegroups       # Routegroups from router

# Run all tests from a service
./run-tests.sh mapexos                  # All tests from mapexos
./run-tests.sh assets                   # All tests from assets
./run-tests.sh router                   # All tests from router
./run-tests.sh http_gateway             # All tests from http_gateway

# Run ALL E2E tests
./run-tests.sh all

# ============================================
# USEFUL COMMANDS
# ============================================

# View complete help
./run-tests.sh help

# List available services and modules
./run-tests.sh list

# Check if services are running
./run-tests.sh check

# ============================================
# OPTIONS
# ============================================

# Run without verbose (quiet)
./run-tests.sh mapexos users -q

# Run in parallel (4 workers)
./run-tests.sh mapexos users -p 4

# With custom timeout
./run-tests.sh mapexos users -t 10m

# Combining options
./run-tests.sh mapexos organizations -q -p 4 -t 5m
```

### 🎯 Autocomplete (Bash Completion)

To enable autocomplete with TAB:

```bash
# Add to your ~/.bashrc or ~/.zshrc:
source /path/to/e2eTests/.run-tests-completion.bash

# Or directly in the current session:
cd workspace_go/packages/e2eTests
source .run-tests-completion.bash

# Now use TAB to autocomplete:
./run-tests.sh [TAB]           # Shows: all check list mapexos assets router http_gateway
./run-tests.sh mapexos [TAB]   # Shows: organizations roles groups users memberships
```

### 📋 View Available Options

```bash
# View detailed help with examples
./run-tests.sh help

# List all available services and modules
./run-tests.sh list
```

### ⚙️ Alternative Method: Direct Go Commands

If you prefer to run tests directly with `go test`:

```bash
# Tests authenticate themselves at startup; no login script step is needed.
go test ./services/mapexos/organizations -v
go test ./services/mapexos/roles -v
go test ./services/mapexos/... -v
```

### 🔍 Run a Specific Test

```bash
# Only the customer creation test
go test ./services/mapexos/organizations -v -run TestCreateOrganization_Customer

# Only hierarchy tests
go test ./services/mapexos/organizations -v -run TestOrganizationHierarchy
```

### 4. Run ALL E2E Tests

```bash
go test ./... -v
```

### 5. Run with Longer Timeout (for slow tests)

```bash
go test ./services/mapexos/organizations -v -timeout 5m
```

### 6. Run in Parallel (faster)

```bash
go test ./... -v -parallel 4
```

## Environment Variables

You can customize service URLs:

```bash
export MAPEXOS_URL=http://localhost:5000
export ROUTER_URL=http://localhost:5003
export ASSETS_URL=http://localhost:5002
export GATEWAY_URL=http://localhost:5001
export MONGO_URI=mongodb://localhost:27017
export MONGO_DATABASE=mapexos_test
export REDIS_HOST=localhost
export REDIS_PORT=6379
```

## Service Ports

| Service | Port | Required | Note |
|---------|------|----------|------|
| mapexos | 5000 | ✅ Yes | Includes /auth/* endpoints |
| router | 5003 | ⚠️ Optional | For routegroups tests |
| assets | 5002 | ⚠️ Optional | For assets tests |
| http_gateway | 5001 | ⚠️ Optional | For datasources tests |

## Test Coverage

### ✅ Implemented

- **Organizations:** 18 tests
  - CREATE: Customer, Site, Building, Minimal, Validations
  - GET: ById, List, NotFound
  - UPDATE: Name, Disable, Full
  - DELETE: Delete, NotFound
  - HIERARCHY: PathKey propagation, CustomerID inheritance

- **Roles:** 20 tests
  - CREATE: System, Org, Minimal, Validations
  - GET: ById, List, NotFound
  - UPDATE: Name, Permissions, Disable, Full
  - DELETE: Delete, NotFound
  - PERMISSIONS: Wildcards, Admin

- **Groups:** ~15 tests (existing)
- **Users:** ~12 tests (existing)

### 🚧 TODO

- Memberships (mapexos)
- Assets (assets service)
- AssetTemplates (assets service)
- DataSources (http_gateway)
- RouteGroups (router)
- Journey Tests (3 integration tests)

## Tips

1. **Clean test data:**
   ```bash
   # Tests use automatic cleanup, but if something gets stuck:
   mongo mapexos_test --eval "db.organizations.deleteMany({slug: /test-org-e2e/})"
   ```

2. **Debug a test:**
   ```bash
   # Add -v for verbose
   go test ./services/mapexos/organizations -v -run TestCreateOrganization_Customer
   ```

3. **View only failures:**
   ```bash
   go test ./... -v | grep -E "FAIL|PASS"
   ```

4. **Run with race detector:**
   ```bash
   go test ./... -race -v
   ```

## Fixtures

Each module has its `fixtures/` folder with example JSON payloads:

```
organizations/fixtures/
├── create_customer.json
├── create_site.json
├── create_building.json
├── update_name.json
└── ...

roles/fixtures/
├── create_system_role.json
├── create_org_role.json
├── update_permissions.json
└── ...
```

Fixtures support placeholders:
- `{{ORG_ID}}` - Replaced by test organization ID
- `{{PARENT_ID}}` - Replaced by parent ID in hierarchy
- `{{USER_ID}}` - Replaced by test user ID

## Conventions

1. **Test naming:**
   - `TestCreate{Resource}_{Scenario}`
   - `TestGet{Resource}ById`
   - `TestUpdate{Resource}_{Field}`
   - `TestDelete{Resource}`

2. **Cleanup:**
   - Always use `t.Cleanup(func() { cleanup...})` or `defer cleanup()`
   - Tests must be idempotent

3. **Fixtures:**
   - Use fixtures for complex payloads
   - Minimal fixtures for required field tests

4. **Assertions:**
   - Use `utils.AssertCreated()`, `utils.AssertOK()`, etc.
   - Use `require.*` for fatal errors
   - Use `assert.*` for non-critical validations
