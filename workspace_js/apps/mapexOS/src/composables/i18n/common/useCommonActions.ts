import { computed } from 'vue';
import { useTS } from '@utils/translation';

/**
 * Common action translations (buttons, links)
 * Used across all pages for consistency
 *
 * @example
 * ```ts
 * const { actions } = useCommonActions();
 * <q-btn :label="actions.save" />
 * ```
 */
export function useCommonActions() {
  const ts = useTS({ capitalize: true });

  return {
    actions: {
      // Basic actions
      add: computed(() => ts('common.actions.add')),
      addNew: computed(() => ts('common.actions.addNew')),
      edit: computed(() => ts('common.actions.edit')),
      delete: computed(() => ts('common.actions.delete')),
      save: computed(() => ts('common.actions.save')),
      saveChanges: computed(() => ts('common.actions.saveChanges')),
      cancel: computed(() => ts('common.actions.cancel')),
      confirm: computed(() => ts('common.actions.confirm')),
      close: computed(() => ts('common.actions.close')),

      // Search and filter
      search: computed(() => ts('common.actions.search')),
      filter: computed(() => ts('common.actions.filter')),
      clear: computed(() => ts('common.actions.clear')),
      apply: computed(() => ts('common.actions.apply')),

      // View actions
      view: computed(() => ts('common.actions.view')),
      viewDetails: computed(() => ts('common.actions.viewDetails')),

      // Data actions
      export: computed(() => ts('common.actions.export')),
      import: computed(() => ts('common.actions.import')),
      refresh: computed(() => ts('common.actions.refresh')),

      // Navigation
      back: computed(() => ts('common.actions.back')),
      next: computed(() => ts('common.actions.next')),
      previous: computed(() => ts('common.actions.previous')),
      finish: computed(() => ts('common.actions.finish')),
    },
  };
}
