import { computed } from 'vue';
import { useTS } from '@utils/translation';

/**
 * Translations composable for ListHeaderMenu component
 * Provides all translated strings for the list header menu
 */
export function useListHeaderMenuTranslations() {
  const ts = useTS({ capitalize: true });

  return {
    /** Items per page label */
    itemsPerPage: computed(() => ts('components.headers.listHeaderMenu.itemsPerPage')),

    /** Visible columns section header */
    visibleColumns: computed(() => ts('components.headers.listHeaderMenu.visibleColumns')),

    /** Filtered indicator suffix */
    filtered: computed(() => ts('components.headers.listHeaderMenu.filtered')),

    /** Columns section header */
    columns: computed(() => ts('components.headers.listHeaderMenu.columns')),

    /** Refresh inline link label */
    refresh: computed(() => ts('components.headers.listHeaderMenu.refresh')),

    /** "Updated now" label (diff < 5s) */
    lastUpdatedNow: computed(() => ts('components.headers.listHeaderMenu.lastUpdated.now')),

    /** "Updated {n}s ago" — call with the number of seconds */
    lastUpdatedSeconds: (n: number) => ts('components.headers.listHeaderMenu.lastUpdated.seconds', { n }),

    /** "Updated {n}m ago" — call with the number of minutes */
    lastUpdatedMinutes: (n: number) => ts('components.headers.listHeaderMenu.lastUpdated.minutes', { n }),

    /** "Updated {n}h ago" — call with the number of hours */
    lastUpdatedHours: (n: number) => ts('components.headers.listHeaderMenu.lastUpdated.hours', { n }),
  };
}
