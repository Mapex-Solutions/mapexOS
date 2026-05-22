import { watch } from 'vue';
import { useOrganizationStore } from '@stores/organization';
import { useLogger } from '@composables/useLogger';

const logger = useLogger('useOrgChangeRefresh');

/**
 * Composable to handle automatic data refresh when organization context changes
 *
 * This composable watches for changes in the selected organization and automatically
 * calls the provided callback function to refresh page data.
 *
 * Purpose:
 * - Ensures data stays consistent with selected organization context
 * - Prevents stale data when user switches organizations
 * - Provides a standardized pattern across all pages that need org-aware data
 *
 * Implementation Details:
 * - Only triggers on actual organization changes (not on initial mount)
 * - Ignores transitions from/to null (page initialization)
 * - Automatically cleans up watch when component unmounts
 *
 * Usage Example:
 * ```typescript
 * import { useOrgChangeRefresh } from '@composables/organizations';
 *
 * // In your component setup:
 * const { loading } = useOrgChangeRefresh(async () => {
 *   await fetchAssets();
 * });
 *
 * // Optional: use returned loading ref for UI feedback
 * <q-spinner v-if="loading" />
 * ```
 *
 * Use Cases:
 * - Asset list pages (refetch assets for new org)
 * - Template list pages (refetch templates for new org)
 * - Dashboard pages (refetch metrics for new org)
 * - Any page that displays org-specific data
 *
 * Performance:
 * - Callback only executes when organization actually changes
 * - No unnecessary calls during initial page load
 * - Supports both sync and async callbacks
 *
 * @param callback - Function to call when organization changes (can be async)
 * @param options - Optional configuration
 * @param options.immediate - If true, triggers callback immediately on mount (default: false)
 * @returns Object with loading ref for optional UI feedback
 */
export function useOrgChangeRefresh(
  callback: () => void | Promise<void>,
  options?: {
    immediate?: boolean;
  }
) {
  const orgStore = useOrganizationStore();

  // Watch for organization changes
  watch(
    () => orgStore.selectedOrganizationId,
    async (newOrgId, oldOrgId) => {
      // Only refetch if org actually changed (not initial load)
      // This prevents double-fetching: once in onMounted, once in watch
      if (newOrgId && oldOrgId && newOrgId !== oldOrgId) {
        logger.debug(
          `Organization changed from ${oldOrgId} to ${newOrgId}, triggering refresh...`
        );

        await callback();
      }
    },
    {
      immediate: options?.immediate ?? false,
    }
  );

  // Return empty object to maintain consistent return type
  // Future enhancement could include loading state management
  return {};
}
