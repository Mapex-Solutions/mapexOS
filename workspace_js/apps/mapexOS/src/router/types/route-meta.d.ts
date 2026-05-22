import 'vue-router';

/**
 * Extend Vue Router's RouteMeta interface to include custom metadata
 * for authentication and authorization
 */
declare module 'vue-router' {
  interface RouteMeta {
    /**
     * Route is public and can be accessed without authentication
     * Examples: /login, /errors/:type
     */
    isPublic?: boolean;

    /**
     * Route is protected and requires authentication
     * This is the DEFAULT behavior for routes without explicit metadata
     */
    isProtected?: boolean;

    /**
     * User must have ANY of these permissions to access route.
     * Guard redirects to /errors/forbidden if check fails.
     * Example: ['assets.list', 'assets.read']
     */
    permissions?: string[];

    /**
     * User must have ALL of these permissions to access route.
     * Guard redirects to /errors/forbidden if check fails.
     * Example: ['assets.list', 'assets.delete']
     */
    requireAllPermissions?: string[];

    /**
     * FUTURE: User must have ANY of these roles to access route
     * Example: ['admin', 'manager']
     */
    roles?: string[];

    /**
     * FUTURE: User must have ALL of these roles to access route
     * Example: ['admin', 'superuser']
     */
    requireAllRoles?: string[];
  }
}

export {};
