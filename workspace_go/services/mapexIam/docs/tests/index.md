# Tests

## Prerequisites
- Go 1.25+
- No running infrastructure required for unit tests (use mocks)
- MongoDB and Redis for integration tests

## Run
```bash
go test ./... -count=1
```

## Test Files
| File | Module | Description |
|---|---|---|
| `src/modules/auth/application/services/auth_service_test.go` | auth | Unit tests for AuthService (login, token refresh, session management) |
| `src/modules/groups/application/services/group_service_test.go` | groups | Unit tests for GroupService (CRUD, member management, cache invalidation) |
| `src/modules/groups/application/services/group_query_service_test.go` | groups | Unit tests for GroupQueryService (list, search, filtering) |
| `src/modules/memberships/application/services/membership_service_test.go` | memberships | Unit tests for MembershipService (CRUD, cache invalidation) |
| `src/modules/roles/application/services/role_service_test.go` | roles | Unit tests for RoleService (CRUD, permission changes, cache invalidation) |
| `src/modules/users/application/services/user_service_test.go` | users | Unit tests for UserService (CRUD, profile, multi-tenant filtering) |

## Mocks
- `src/shared/mocks/` — shared infrastructure mocks
