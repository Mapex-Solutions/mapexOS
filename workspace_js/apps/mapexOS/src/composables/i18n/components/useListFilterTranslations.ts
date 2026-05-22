import { computed } from 'vue';
import { useTS } from '@utils/translation';

/**
 * ListFilter component translations
 *
 * Structure mirrors:
 * - File: src/components/filters/listFilter/ListFilter.vue
 * - JSON: src/i18n/{locale}/components/filters.json
 * - Composable: src/composables/i18n/components/useListFilterTranslations.ts
 *
 * Provides all translations for the ListFilter component including:
 * - Filter actions (apply, clear)
 * - Filter hints and labels
 * - Date range labels
 *
 * @example
 * ```ts
 * // In ListFilter.vue
 * const { listFilter } = useListFilterTranslations();
 *
 * <div>{{ listFilter.hint.value }}</div>
 * <q-btn :label="listFilter.apply.value" />
 * ```
 */
export function useListFilterTranslations() {
  const ts = useTS({ capitalize: true });

  return {
    /**
     * ListFilter translations
     * Mirrors: components.filters.listFilter
     */
    listFilter: {
      title: computed(() => ts('components.filters.listFilter.title')),
      apply: computed(() => ts('components.filters.listFilter.apply')),
      clear: computed(() => ts('components.filters.listFilter.clear')),
      hint: computed(() => ts('components.filters.listFilter.hint')),
      includeChildren: computed(() => ts('components.filters.listFilter.includeChildren')),
      includeChildrenTooltip: computed(() => ts('components.filters.listFilter.includeChildrenTooltip')),
      active: computed(() => ts('components.filters.listFilter.active')),
      close: computed(() => ts('components.filters.listFilter.close')),
      dateRange: {
        from: computed(() => ts('components.filters.listFilter.dateRange.from')),
        to: computed(() => ts('components.filters.listFilter.dateRange.to')),
        selectRange: computed(() => ts('components.filters.listFilter.dateRange.selectRange')),
      },
    },
  };
}
