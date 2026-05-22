# MapexOS Permissions Hierarchy

This package defines all permissions for the MapexOS platform following a hierarchical wildcard structure.

## Permission Hierarchy

```
mapex.*                              # Root platform access
├── mapexos.*                        # MapexOS service
│   ├── auth.*
│   │   ├── auth.login
│   │   ├── auth.logout
│   │   └── auth.refresh
│   ├── organizations.*
│   │   ├── organization.list
│   │   ├── organization.create
│   │   ├── organization.read
│   │   ├── organization.update
│   │   └── organization.delete
│   ├── users.*
│   │   ├── user.list
│   │   ├── user.create
│   │   ├── user.read
│   │   ├── user.update
│   │   └── user.delete
│   ├── roles.*
│   ├── groups.*
│   ├── memberships.*
│   └── lists.*
│
├── router.*                         # Router service
│   └── routegroups.*
│       ├── routegroup.list
│       ├── routegroup.create
│       ├── routegroup.read
│       ├── routegroup.update
│       └── routegroup.delete
│
└── assets.*                         # HTTP Gateway service
    ├── datasources.*
    │   ├── datasource.list
    │   ├── datasource.create
    │   ├── datasource.read
    │   ├── datasource.update
    │   └── datasource.delete
    ├── assets.*
    │   ├── asset.list
    │   ├── asset.create
    │   ├── asset.read
    │   ├── asset.update
    │   └── asset.delete
    └── assettemplates.*
        ├── assettemplate.list
        ├── assettemplate.create
        ├── assettemplate.read
        ├── assettemplate.update
        └── assettemplate.delete
```

## Usage

### Global Permissions
- `mapex.*` - Root access to entire platform (for ROOT users)

### Service-Level Permissions
- `mapexos.*` - Access to all MapexOS modules
- `router.*` - Access to all Router modules
- `assets.*` - Access to all HTTP Gateway modules

### Module-Level Permissions
- `routegroup.*` - Access to all routegroup operations
- `datasource.*` - Access to all datasource operations
- `organization.*` - Access to all organization operations

### Operation-Level Permissions
- `routegroup.list` - List route groups
- `datasource.create` - Create data sources
- `user.update` - Update users

## Wildcard Matching

The permission system supports wildcard matching:
- `mapex.*` grants access to ALL operations on ALL services
- `router.*` grants access to ALL operations on Router service
- `routegroup.*` grants access to ALL routegroup operations

## Package Structure

```
permissions/
├── global.go                # Platform-level: mapex.*
├── mapexos/
│   ├── mapexos.go          # Service-level: mapexos.*
│   ├── auth.go             # Module-level: auth.*
│   ├── organizations.go
│   ├── users.go
│   ├── roles.go
│   ├── groups.go
│   ├── memberships.go
│   └── lists.go
├── router/
│   ├── router.go           # Service-level: router.*
│   └── routegroups.go      # Module-level: routegroup.*
└── assets/
    ├── assets_service.go   # Service-level: assets.*
    ├── datasources.go      # Module-level: datasource.*
    ├── assets.go           # Module-level: asset.*
    └── assettemplates.go   # Module-level: assettemplate.*
```

## Examples

### ROOT User
```go
// ROOT user has global access
permissions := []string{"mapex.*"}
// Can access: datasource.create, routegroup.list, user.update, etc.
```

### Service Admin
```go
// Admin for Router service
permissions := []string{"router.*"}
// Can access: routegroup.create, routegroup.delete, etc.
// Cannot access: datasource.create, user.update, etc.
```

### Module-Specific User
```go
// Can only manage datasources
permissions := []string{"datasource.*"}
// Can access: datasource.list, datasource.create, etc.
// Cannot access: routegroup.create, asset.list, etc.
```

### Operation-Specific User
```go
// Can only view routegroups
permissions := []string{"routegroup.list", "routegroup.read"}
// Can access: routegroup.list, routegroup.read
// Cannot access: routegroup.create, routegroup.delete, etc.
```
