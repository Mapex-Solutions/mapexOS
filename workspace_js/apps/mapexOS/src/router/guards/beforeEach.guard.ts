import type { NavigationGuardNext, RouteLocationNormalized } from 'vue-router';
import { useAuthStore } from '@stores/auth';
import { useOrganizationStore } from '@stores/organization';
import { usePermissionStore } from '@stores/permission';
import { useLogger } from '@composables/useLogger';

const logger = useLogger('RouterGuard');

/**
 * Select initial organization with intelligent strategy
 * Priority order:
 * 1. If user has only 1 org → auto-select it
 * 2. Restore last selected org from localStorage (if user still has access)
 * 3. Prioritize org with scope "recursive" (vendor)
 * 4. Fallback: first vendor (type === 'vendor')
 * 5. Last fallback: first org in list
 */
function selectInitialOrganization(organizations: { id: string; name: string; scope?: string; type?: string }[]): string {
  if (organizations.length === 0) {
    throw new Error('No organizations available');
  }

  // 1. If only 1 org → auto-select (simplest case)
  if (organizations.length === 1) {
    return organizations[0]!.id;
  }

  // 2. Restore last selected org from localStorage (HIGHEST PRIORITY after single org)
  const lastOrgId = localStorage.getItem('selectedOrgId');
  if (lastOrgId) {
    const hasAccess = organizations.some(org => org.id === lastOrgId);
    if (hasAccess) {
      return lastOrgId;
    }
    // User lost access to previously selected org, clear localStorage
    localStorage.removeItem('selectedOrgId');
    localStorage.removeItem('selectedOrgName');
  }

  // 3. Prioritize org with scope "recursive" (usually vendor)
  const recursiveOrg = organizations.find(org => org.scope === 'recursive');
  if (recursiveOrg) {
    return recursiveOrg.id;
  }

  // 4. Fallback: first vendor
  const vendor = organizations.find(org => org.type === 'vendor');
  if (vendor) {
    return vendor.id;
  }

  // 5. Last fallback: first org in list
  return organizations[0]!.id;
}

/**
 * Ensure organization is selected and coverage is loaded
 * Handles restoration from localStorage and intelligent selection
 */
async function ensureOrganizationContext(org: ReturnType<typeof useOrganizationStore>): Promise<void> {
  // 1. Guarantee coverage is loaded
  if (org.flatList.length === 0) {
    await org.fetchCoverage();
  }

  // 2. Check if organization already selected
  if (org.selectedOrganizationId) {
    // Validate it still exists in coverage
    const orgExists = org.flatList.some(o => o.id === org.selectedOrganizationId);
    if (orgExists) {
      return; // Already selected and valid
    }
    // Org no longer exists, clear and re-select
    org.selectedOrganizationId = null;
    org.selectedOrganizationName = null;
  }

  // 3. Try to restore from localStorage
  const lastOrgId = localStorage.getItem('selectedOrgId');
  const lastOrgName = localStorage.getItem('selectedOrgName');

  if (lastOrgId) {
    const orgExists = org.flatList.some(o => o.id === lastOrgId);
    if (orgExists) {
      org.selectedOrganizationId = lastOrgId;
      org.selectedOrganizationName = lastOrgName;
      return;
    }
    // Org no longer accessible, clear localStorage
    localStorage.removeItem('selectedOrgId');
    localStorage.removeItem('selectedOrgName');
  }

  // 4. No valid selection, apply intelligent selection strategy
  const selectedId = selectInitialOrganization(org.flatList);
  org.selectedOrganizationId = selectedId;

  const selectedOrg = org.flatList.find(o => o.id === selectedId);
  org.selectedOrganizationName = selectedOrg?.name || null;

  // 5. Save to localStorage
  localStorage.setItem('selectedOrgId', selectedId);
  if (selectedOrg?.name) {
    localStorage.setItem('selectedOrgName', selectedOrg.name);
  }
}

/**
 * Global beforeEach navigation guard.
 *
 * CRITICAL: Guard is SELF-SUFFICIENT and METADATA-DRIVEN
 * - Hydrates auth store from storage BEFORE checking any routes
 * - Routes declare their access requirements via meta.isPublic or meta.isProtected
 * - Redirects authenticated users away from login page to dashboard
 * - Ensures organization context is loaded and selected for protected routes
 * - Does NOT depend on App.vue or any other component
 * - Ensures auth state is available for ALL route transitions
 *
 * Route Metadata:
 * - meta.isPublic = true → No authentication required (login, errors)
 * - meta.isProtected = true → Authentication required (default for routes without metadata)
 *
 * Workflow:
 * 1. Hydrate auth from storage FIRST (if not already done)
 * 2. PUBLIC ROUTES:
 *    - If accessing "/" (login) with valid token → redirect to /home (dashboard)
 *    - Otherwise allow access
 * 3. PROTECTED ROUTES:
 *    - Check authentication → redirect to / if not authenticated
 *    - Fetch coverage (if not already loaded)
 *    - Ensure organization selected (restore from localStorage or intelligent selection)
 *    - Check route-level permissions from Pinia cache → redirect to /errors/forbidden if denied
 *    - Permissions are fetched ONCE on org selection, guard only reads cache
 *    - Allow navigation
 *
 * Handles:
 * - Initial login
 * - Page reload (F5)
 * - Organization change (router.go(0))
 * - Navigation between routes
 * - Authentication expiration (401 interceptor handles token refresh)
 * - Redirect authenticated users from login page to dashboard
 */
export default async function beforeEach(
  to: RouteLocationNormalized,
  from: RouteLocationNormalized,
  next: NavigationGuardNext
): Promise<void> {
  const auth = useAuthStore();
  const org = useOrganizationStore();

  // 1. HYDRATE AUTH FROM STORAGE FIRST (if not already done)
  // This must happen BEFORE any route checks to ensure we have the latest auth state
  // Handles page reload (F5) and direct navigation
  if (!auth.accessToken) {
    auth.hydrateFromStorage();
  }

  // 2. PUBLIC ROUTES (like login page)
  if (to.meta.isPublic === true) {
    // Allow change-password page for authenticated users with the flag
    if (to.path === '/change-password' && auth.isAuthenticated) {
      if (!auth.user?.changePasswordNextLogin) {
        return next('/home'); // Flag is false, no need to be here
      }
      return next(); // Allow access
    }

    // Special case: If user is already authenticated and trying to access login page ("/")
    // Redirect to dashboard instead
    if (to.path === '/' && auth.isAuthenticated) {
      return next('/home');
    }
    // Allow access to other public routes
    return next();
  }

  // 3. PROTECTED ROUTES (default behavior for routes without explicit metadata)
  // Routes are protected by default unless explicitly marked as public

  // 3.1 Check authentication
  if (!auth.isAuthenticated) {
    // Redirect to login with redirect query param for post-login navigation
    return next({
      path: '/',
      query: { redirect: to.fullPath }
    });
  }

  // 3.1.5 Force password change check
  if (auth.user?.changePasswordNextLogin && to.path !== '/change-password') {
    return next({ path: '/change-password', query: { redirect: to.fullPath } });
  }

  // 3.2 Ensure organization context is loaded and selected
  try {
    await ensureOrganizationContext(org);
  } catch (error) {
    logger.error('Failed to ensure organization context', error);

    // If no organizations available, redirect to error page
    if (!org.flatList || org.flatList.length === 0) {
      return next('/errors/no-organization');
    }

    // For other errors, log but allow navigation (user might already have context)
    logger.error('Organization context error, allowing navigation', error);
  }

  // 3.3 Ensure permissions are cached (ONE request after F5/page reload only)
  // Normal flow: permissions are fetched when user selects org (org store actions).
  // F5 flow: ensureOrganizationContext restores org from localStorage but skips selectOrganization(),
  //          so we need a single fetch here if permissions haven't been loaded yet.
  const permStore = usePermissionStore();

  if (!permStore.isLoaded && org.selectedOrganizationId) {
    await permStore.fetchPermissions();
  }

  if (to.meta.permissions?.length && permStore.isLoaded) {
    if (!permStore.hasAnyPermission(to.meta.permissions)) {
      logger.warn(`Access denied to ${to.path} — missing permissions: ${to.meta.permissions.join(', ')}`);
      return next('/errors/forbidden');
    }
  }

  if (to.meta.requireAllPermissions?.length && permStore.isLoaded) {
    if (!permStore.hasAllPermissions(to.meta.requireAllPermissions)) {
      logger.warn(`Access denied to ${to.path} — missing all permissions: ${to.meta.requireAllPermissions.join(', ')}`);
      return next('/errors/forbidden');
    }
  }

  // 3.4 FUTURE: Check roles from route metadata
  // if (to.meta.roles) {
  //   if (!auth.hasAnyRole(to.meta.roles)) {
  //     return next('/errors/forbidden');
  //   }
  // }
  //
  // if (to.meta.requireAllRoles) {
  //   if (!auth.hasAllRoles(to.meta.requireAllRoles)) {
  //     return next('/errors/forbidden');
  //   }
  // }

  // 3.5 User is authenticated and has organization context, allow navigation
  return next();
}
