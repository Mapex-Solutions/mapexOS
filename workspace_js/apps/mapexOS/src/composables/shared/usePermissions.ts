import { computed } from 'vue';
import { usePermissionStore } from '@stores/permission';

/**
 * Composable for reactive permission checks in components.
 * Wraps the permission store getters in computed refs for template reactivity.
 *
 * @returns Object with permission check methods and loading state
 *
 * @example
 * ```vue
 * <script setup lang="ts">
 * const { canCreate, canDelete, permissionsLoaded } = usePermissions();
 * const canCreateAsset = canCreate('assets');
 * const canDeleteAsset = canDelete('assets');
 * // In script: use .value → canCreateAsset.value
 * // In template: auto-unwrapped → canCreateAsset (no .value)
 * </script>
 *
 * <template>
 *   <PageHeader :button="canCreateAsset ? addConfig : undefined" />
 *   <q-btn v-if="canDeleteAsset" icon="delete" />
 * </template>
 * ```
 */
export function usePermissions() {
  const store = usePermissionStore();

  return {
    /**
     * Check if user has a specific permission
     *
     * @param {string} perm - Permission string to check
     * @returns {import('vue').ComputedRef<boolean>} Reactive boolean
     */
    hasPermission: (perm: string) => computed(() => store.hasPermission(perm)),

    /**
     * Check if user has ANY of the required permissions
     *
     * @param {string[]} perms - Array of permission strings
     * @returns {import('vue').ComputedRef<boolean>} Reactive boolean
     */
    hasAnyPermission: (perms: string[]) => computed(() => store.hasAnyPermission(perms)),

    /**
     * Check if user has ALL of the required permissions
     *
     * @param {string[]} perms - Array of permission strings
     * @returns {import('vue').ComputedRef<boolean>} Reactive boolean
     */
    hasAllPermissions: (perms: string[]) => computed(() => store.hasAllPermissions(perms)),

    /**
     * Shorthand: check if user can list a resource
     *
     * @param {string} resource - Resource name (e.g., 'assets')
     * @returns {import('vue').ComputedRef<boolean>} Reactive boolean
     */
    canList: (resource: string) => computed(() => store.hasPermission(`${resource}.list`)),

    /**
     * Shorthand: check if user can create a resource
     *
     * @param {string} resource - Resource name (e.g., 'assets')
     * @returns {import('vue').ComputedRef<boolean>} Reactive boolean
     */
    canCreate: (resource: string) => computed(() => store.hasPermission(`${resource}.create`)),

    /**
     * Shorthand: check if user can read a resource
     *
     * @param {string} resource - Resource name (e.g., 'assets')
     * @returns {import('vue').ComputedRef<boolean>} Reactive boolean
     */
    canRead: (resource: string) => computed(() => store.hasPermission(`${resource}.read`)),

    /**
     * Shorthand: check if user can update a resource
     *
     * @param {string} resource - Resource name (e.g., 'assets')
     * @returns {import('vue').ComputedRef<boolean>} Reactive boolean
     */
    canUpdate: (resource: string) => computed(() => store.hasPermission(`${resource}.update`)),

    /**
     * Shorthand: check if user can delete a resource
     *
     * @param {string} resource - Resource name (e.g., 'assets')
     * @returns {import('vue').ComputedRef<boolean>} Reactive boolean
     */
    canDelete: (resource: string) => computed(() => store.hasPermission(`${resource}.delete`)),

    /** Whether permissions are loading */
    permissionsLoading: computed(() => store.loading),

    /** Whether permissions have been loaded */
    permissionsLoaded: computed(() => store.isLoaded),
  };
}
