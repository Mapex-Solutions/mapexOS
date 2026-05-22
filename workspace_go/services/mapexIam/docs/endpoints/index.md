# Endpoints

## HTTP API

### Auth (base: `/api/v1/auth`)
| Method | Path | Description |
|---|---|---|
| POST | `/login` | Login |
| POST | `/logout` | Logout |
| POST | `/refresh` | Refresh token |
| GET | `/users/me/coverage` | Get current user coverage |
| GET | `/me/permissions` | Get current user permissions |

### Users (base: `/api/v1/users`)
| Method | Path | Description |
|---|---|---|
| GET | `/me` | Get current user |
| PATCH | `/me` | Update current user |
| PATCH | `/me/tour` | Disable onboarding tour |
| GET | `/` | List users |
| GET | `/counter` | Count users |
| POST | `/` | Create user |
| GET | `/:userId` | Get user by ID |
| PATCH | `/:userId` | Update user |
| DELETE | `/:userId` | Delete user |

### Organizations (base: `/api/v1/organizations`)
| Method | Path | Description |
|---|---|---|
| GET | `/` | List organizations |
| GET | `/tree` | Get organizations tree |
| POST | `/` | Create organization |
| GET | `/:organizationId` | Get organization |
| PATCH | `/:organizationId` | Update organization |
| DELETE | `/:organizationId` | Delete organization |

### Roles (base: `/api/v1/roles`)
| Method | Path | Description |
|---|---|---|
| GET | `/` | List roles |
| POST | `/` | Create role |
| GET | `/:roleId` | Get role |
| PATCH | `/:roleId` | Update role |
| DELETE | `/:roleId` | Delete role |

### Groups (base: `/api/v1/groups`)
| Method | Path | Description |
|---|---|---|
| GET | `/` | List groups |
| GET | `/counter` | Count groups |
| POST | `/` | Create group |
| GET | `/:groupId` | Get group |
| PATCH | `/:groupId` | Update group |
| DELETE | `/:groupId` | Delete group |
| GET | `/:groupId/members` | List group members |
| POST | `/:groupId/members` | Add group member |
| DELETE | `/:groupId/members/:userId` | Remove group member |

### Memberships (base: `/api/v1/memberships`)
| Method | Path | Description |
|---|---|---|
| GET | `/` | List memberships |
| POST | `/` | Create membership |
| GET | `/:membershipId` | Get membership |
| PATCH | `/:membershipId` | Update membership |
| DELETE | `/:membershipId` | Delete membership |

### Memberships (me) (base: `/api/v1/memberships/me`)
| Method | Path | Description |
|---|---|---|
| GET | `/coverage` | Get my coverage |

### Lists (base: `/api/v1/lists`)
| Method | Path | Description |
|---|---|---|
| GET | `/` | List lists |
| POST | `/` | Create list |
| GET | `/:listId` | Get list |
| PATCH | `/:listId` | Update list |
| DELETE | `/:listId` | Delete list |

### Onboarding (base: `/api/v1/onboarding`)
| Method | Path | Description |
|---|---|---|
| POST | `/users` | Create user with memberships |
| PATCH | `/users/:userId` | Update user access configuration |

### Internal Auth (base: `/internal/auth`)
| Method | Path | Description |
|---|---|---|
| POST | `/build-authorization` | Build authorization cache |
| POST | `/build-coverage` | Build coverage cache |

## NATS

### Outbound (publish)
| Subject | Stream | Description |
|---|---|---|
| `mapexos.cache.invalidation.*` | `MAPEXOS_CACHE_INVALIDATION` | Cache invalidation events (role/org/group/membership changes) |
| `mapexos.lists.name_updated` | `MAPEXOS-LISTS` | List name sync for downstream services (asset templates) |

### Inbound (subscribe)
| Subject | Stream | Description |
|---|---|---|
| `mapexos.cache.invalidation.>` | `MAPEXOS_CACHE_INVALIDATION` | Consume invalidation events to refresh auth/coverage Redis caches |

### Cache Invalidation Event Types
| Event Type | Trigger |
|---|---|
| `role.permissions.changed` | Role permissions updated |
| `role.deleted` | Role deleted |
| `organization.access_policy.changed` | Org access policy updated |
| `organization.hierarchy.changed` | Org hierarchy modified |
| `membership.changed` | Membership created/updated |
| `membership.deleted` | Membership deleted |
| `group.changed` | Group created/updated or member added/removed |
| `group.deleted` | Group deleted |

## Observability
- Health: `GET /health`
- Metrics: `GET /metrics`
