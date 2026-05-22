import { computed } from 'vue';
import { useI18n } from 'vue-i18n';

/**
 * Composable for accessing global error translations
 * Provides i18n messages for common HTTP error codes
 */
export function useErrorTranslations() {
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
    }
  };
}
