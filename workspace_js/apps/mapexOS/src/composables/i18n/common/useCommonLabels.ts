import { computed } from 'vue';
import { useTS } from '@utils/translation';

/**
 * Common label translations
 * Used for generic field labels across the application
 *
 * @example
 * ```ts
 * const { labels } = useCommonLabels();
 * <q-input :label="labels.name" />
 * ```
 */
export function useCommonLabels() {
  const ts = useTS({ capitalize: true });

  return {
    labels: {
      name: computed(() => ts('common.labels.name')),
      description: computed(() => ts('common.labels.description')),
      status: computed(() => ts('common.labels.status')),
      active: computed(() => ts('common.labels.active')),
      inactive: computed(() => ts('common.labels.inactive')),
      created: computed(() => ts('common.labels.created')),
      updated: computed(() => ts('common.labels.updated')),
      createdAt: computed(() => ts('common.labels.createdAt')),
      updatedAt: computed(() => ts('common.labels.updatedAt')),
      actions: computed(() => ts('common.labels.actions')),
      yes: computed(() => ts('common.labels.yes')),
      no: computed(() => ts('common.labels.no')),
      all: computed(() => ts('common.labels.all')),
    },
  };
}
