# Router Authentication & Authorization System

## Overview

This router implements a **Guard-First Architecture** where ALL authentication, authorization, and organization context logic is handled by the Router Guard (`guards/beforeEach.guard.ts`), not by components.

## Architecture Principles

### 1. Guard-First Pattern

All security logic is centralized in the Router Guard:
- Authentication check (JWT token validation)
- Organization coverage loading
- Organization context selection
- Permission/Role checks (future)

**App.vue is intentionally minimal** - it only renders `<router-view />` with NO auth/org logic.

### 2. Metadata-Driven Routes

Routes declare their access requirements via `meta` object:

```typescript
// Public route (no authentication required)
{
  path: '/login',
  meta: { isPublic: true },
  component: LoginPage,
}

// Protected route (authentication required)
{
  path: '/assets',
  meta: { isProtected: true },
  component: AssetsPage,
}
```

### 3. TypeScript Route Meta

Type definitions in `types/route-meta.d.ts` extend Vue Router's `RouteMeta`:

```typescript
interface RouteMeta {
  isPublic?: boolean;           // No auth required
  isProtected?: boolean;        // Auth required (default)
  permissions?: string[];       // FUTURE: ANY of these permissions
  requireAllPermissions?: string[];  // FUTURE: ALL of these
  roles?: string[];             // FUTURE: ANY of these roles
  requireAllRoles?: string[];   // FUTURE: ALL of these
}
```

## Route Categories

### Public Routes (isPublic: true)

Accessible without authentication:
- `/` - Login page
- `/errors/no-organization` - Error page
- `/:catchAll(.*)*` - 404 page

### Protected Routes (isProtected: true)

Require authentication (DEFAULT for all routes without explicit metadata):
- `/assets` - Assets management
- `/users` - User management
- `/business_rules` - Business rules
- `/home` - Dashboard
- All other application routes

## Guard Workflow

### For Public Routes

```
1. Check meta.isPublic === true
2. Allow access immediately (no checks)
```

### For Protected Routes

```
1. Check meta.isPublic === true → Allow if true
2. Hydrate auth from storage (if not already done)
3. Check authentication → Redirect to /login if not authenticated
4. Ensure organization coverage loaded (fetch if needed)
5. Ensure organization selected (restore from localStorage or intelligent selection)
6. FUTURE: Check permissions from meta.permissions
7. FUTURE: Check roles from meta.roles
8. Allow navigation
```

## Organization Context Management

### Intelligent Organization Selection

The guard uses a priority-based strategy to select organization:

1. **Single org** → Auto-select (simplest case)
2. **Restore from localStorage** → If user still has access (HIGHEST PRIORITY)
3. **Scope "recursive"** → Prioritize vendor org
4. **Type "vendor"** → Fallback to first vendor
5. **First in list** → Last fallback

### Organization Restoration Flow

On page reload or organization change:

```
1. Guard hydrates auth from storage
2. Guard fetches coverage (if not loaded)
3. Guard checks if org already selected in store
   - If yes and still valid → Keep it
   - If no → Try to restore from localStorage
4. If restored from localStorage → Validate user still has access
5. If not in localStorage or invalid → Apply intelligent selection
6. Save to localStorage for future restores
```

## Use Cases Handled

### 1. Initial Login
```
User logs in → Guard allows /login (public)
→ After login, redirect to protected route
→ Guard hydrates auth, loads coverage, selects org
→ User navigates normally
```

### 2. Page Reload (F5)
```
User reloads page
→ Guard hydrates auth from storage
→ Guard fetches coverage
→ Guard restores org from localStorage
→ User stays authenticated and in same org context
```

### 3. Organization Change
```
User selects different org in UI
→ Store updates selectedOrganizationId
→ Store saves to localStorage
→ App calls router.go(0) to reload
→ Guard hydrates auth from storage
→ Guard restores new org from localStorage
→ Page reloads with new context
```

### 4. Unauthenticated Access
```
User navigates to /assets without auth
→ Guard checks meta.isProtected === true
→ Guard hydrates auth from storage
→ No valid token found
→ Guard redirects to /login with redirect query
→ After login, user returns to /assets
```

### 5. No Organizations Available
```
User authenticates but has no org access
→ Guard fetches coverage
→ No organizations returned
→ Guard redirects to /errors/no-organization
```

## Future: Permissions & Roles

### Adding Permissions to Routes

```typescript
{
  path: '/assets',
  meta: {
    isProtected: true,
    permissions: ['assets.list', 'assets.read'],  // ANY of these
  }
}
```

Guard will check:
```typescript
if (to.meta.permissions) {
  if (!auth.hasAnyPermission(to.meta.permissions)) {
    return next('/errors/forbidden');
  }
}
```

### Adding Roles to Routes

```typescript
{
  path: '/admin/settings',
  meta: {
    isProtected: true,
    roles: ['admin', 'manager'],  // ANY of these
  }
}
```

Guard will check:
```typescript
if (to.meta.roles) {
  if (!auth.hasAnyRole(to.meta.roles)) {
    return next('/errors/forbidden');
  }
}
```

## Benefits

1. **Centralized Security** - All auth logic in one place (Guard)
2. **Clear Separation** - Components don't handle auth/org logic
3. **Type-Safe** - TypeScript validates route metadata
4. **Scalable** - Easy to add permissions/roles in future
5. **Testable** - Guard logic isolated and testable
6. **No Race Conditions** - Guard executes before component lifecycle
7. **Consistent Behavior** - All navigation scenarios handled uniformly
8. **DRY** - No duplication of auth logic across components

## Implementation Files

- `guards/beforeEach.guard.ts` - Main guard implementation
- `types/route-meta.d.ts` - TypeScript route metadata types
- `routes/**/*.ts` - All route definitions with metadata
- `App.vue` - Minimal component (only `<router-view />`)

## Migration Notes

### From App.vue to Guard

**BEFORE:**
- App.vue had `onMounted` hook
- App.vue checked auth and org
- App.vue called `initializeAfterLogin()`
- Race conditions possible

**AFTER:**
- App.vue is minimal (no logic)
- Guard handles everything
- Guard executes BEFORE component lifecycle
- No race conditions

### From Array to Metadata

**BEFORE:**
```typescript
const PUBLIC_ROUTES = ['/', '/login', '/errors/...'];
function isPublicRoute(path: string): boolean {
  return PUBLIC_ROUTES.includes(path);
}
```

**AFTER:**
```typescript
if (to.meta.isPublic === true) {
  return next();
}
```

More maintainable, type-safe, and scalable.

## Best Practices

### Adding New Routes

1. **Public routes** (login, errors):
   ```typescript
   {
     path: '/login',
     meta: { isPublic: true },
     component: LoginPage,
   }
   ```

2. **Protected routes** (application pages):
   ```typescript
   {
     path: '/assets',
     meta: { isProtected: true },
     component: AssetsPage,
   }
   ```

3. **Future permissions**:
   ```typescript
   {
     path: '/assets',
     meta: {
       isProtected: true,
       permissions: ['assets.list'],
     },
     component: AssetsPage,
   }
   ```

### Testing Routes

Always test these scenarios:
1. Unauthenticated access to protected route
2. Authenticated access to protected route
3. Public route access (with and without auth)
4. Page reload (F5) while authenticated
5. Organization change
6. Token expiration
7. Loss of organization access

## Troubleshooting

### "Redirecting to login despite being logged in"
- Check if token is in storage
- Check if `hydrateFromStorage()` is working
- Check if token is expired
- Check store's `isAuthenticated` getter

### "Organization not selected after login"
- Check if coverage fetch is successful
- Check if `ensureOrganizationContext()` is called
- Check localStorage for `selectedOrgId`
- Check if user has access to any organizations

### "Permissions not working"
- Verify permissions are in route metadata
- Verify auth store has `hasAnyPermission()` method
- Verify backend returns user permissions
- Uncomment permission checks in guard

## Summary

This Guard-First Architecture with metadata-driven routes provides a robust, scalable, and maintainable authentication/authorization system that handles all edge cases and prepares for future RBAC/ABAC requirements.
