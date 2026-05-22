import { computed } from 'vue';
import { useTS } from '@utils/translation';

/**
 * Error page translations
 *
 * Structure mirrors:
 * - File: src/pages/erros/errorNotFound/ErrorNotFound.vue
 * - JSON: src/i18n/{locale}/pages/error.json
 * - Composable: src/composables/i18n/pages/error/useErrorTranslations.ts
 *
 * Provides all translations for error pages including:
 * - 404 Not Found page content
 * - Error codes and descriptions
 * - Action buttons
 * - ARIA labels for accessibility
 *
 * @example
 * ```ts
 * // In ErrorNotFound.vue
 * const { notFound, aria } = useErrorTranslations();
 *
 * <div>{{ notFound.code.value }}</div>
 * <h5>{{ notFound.title.value }}</h5>
 * <q-btn :label="notFound.goBack.value" />
 * ```
 */
export function useErrorTranslations() {
  const ts = useTS({ capitalize: true });
  const tsTitle = useTS({ titleCase: true });

  return {
    /**
     * 404 Not Found page translations
     * Mirrors: pages.error.notFound
     */
    notFound: {
      code: computed(() => ts('pages.error.notFound.code')),
      title: computed(() => tsTitle('pages.error.notFound.title')),
      description: computed(() => ts('pages.error.notFound.description')),
      goBack: computed(() => ts('pages.error.notFound.goBack')),
    },

    /**
     * 403 Forbidden page translations
     * Mirrors: pages.error.forbidden
     */
    forbidden: {
      code: computed(() => ts('pages.error.forbidden.code')),
      title: computed(() => tsTitle('pages.error.forbidden.title')),
      description: computed(() => ts('pages.error.forbidden.description')),
      goBack: computed(() => ts('pages.error.forbidden.goBack')),
      dashboard: computed(() => tsTitle('pages.error.forbidden.dashboard')),
    },

    /**
     * ARIA labels for accessibility
     * Mirrors: pages.error.aria
     */
    aria: {
      logoAlt: computed(() => ts('pages.error.aria.logoAlt')),
    },
  };
}
