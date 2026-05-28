import { computed } from 'vue';
import { useI18n } from 'vue-i18n';

/**
 * Composable for accessing global, app-wide error/notify message translations.
 * Holds the shared HTTP error copy plus generic notification strings reused
 * across drawers, forms and async actions.
 */
export function useCommonErrors() {
  const { t } = useI18n();

  return {
    http: {
      400: computed(() => t('common.errors.http.400')),
      401: computed(() => t('common.errors.http.401')),
      403: computed(() => t('common.errors.http.403')),
      404: computed(() => t('common.errors.http.404')),
      409: computed(() => t('common.errors.http.409')),
      422: computed(() => t('common.errors.http.422')),
      500: computed(() => t('common.errors.http.500')),
      503: computed(() => t('common.errors.http.503')),
      network: computed(() => t('common.errors.http.network')),
      unknown: computed(() => t('common.errors.http.unknown')),
    },
    apiNotInitialized: computed(() => t('common.errors.apiNotInitialized')),
    copyFailed: computed(() => t('common.errors.copyFailed')),
    loadFailed: computed(() => t('common.errors.loadFailed')),
    saveFailed: computed(() => t('common.errors.saveFailed')),
    deleteFailed: computed(() => t('common.errors.deleteFailed')),
  };
}
