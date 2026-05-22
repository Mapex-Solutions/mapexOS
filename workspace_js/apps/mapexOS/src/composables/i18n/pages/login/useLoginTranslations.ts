import { computed } from 'vue';
import { useTS } from '@utils/translation';

/**
 * Login page translations
 *
 * Structure mirrors:
 * - File: src/pages/login/LoginPage.vue
 * - JSON: src/i18n/{locale}/pages/login.json
 * - Composable: src/composables/i18n/pages/login/useLoginTranslations.ts
 *
 * Provides all translations for the Login page including:
 * - Welcome section (title, subtitle, features)
 * - Form labels and placeholders
 * - Validation messages
 * - Error messages
 * - Language options
 * - ARIA labels for accessibility
 *
 * @example
 * ```ts
 * // In LoginPage.vue
 * const {
 *   welcome,
 *   features,
 *   form,
 *   validation,
 *   errors,
 *   languages
 * } = useLoginTranslations();
 *
 * <h2>{{ welcome.title.value }}</h2>
 * <q-input :label="form.email.label.value" />
 * ```
 */
export function useLoginTranslations() {
  const ts = useTS({ capitalize: true });

  return {
    /**
     * Welcome section translations
     * Mirrors: pages.login.welcome
     */
    welcome: {
      title: computed(() => ts('pages.login.welcome.title')),
      tagline: computed(() => ts('pages.login.welcome.tagline')),
      subtitle: computed(() => ts('pages.login.welcome.subtitle')),
    },

    /**
     * Feature list translations
     * Mirrors: pages.login.features
     */
    features: {
      monitoring: computed(() => ts('pages.login.features.monitoring')),
      analytics: computed(() => ts('pages.login.features.analytics')),
      security: computed(() => ts('pages.login.features.security')),
    },

    /**
     * Form field translations
     * Mirrors: pages.login.form
     */
    form: {
      title: computed(() => ts('pages.login.form.title')),

      email: {
        label: computed(() => ts('pages.login.form.email.label')),
        placeholder: computed(() => ts('pages.login.form.email.placeholder')),
      },

      password: {
        label: computed(() => ts('pages.login.form.password.label')),
        placeholder: computed(() => ts('pages.login.form.password.placeholder')),
      },

      rememberMe: computed(() => ts('pages.login.form.rememberMe')),
      forgotPassword: computed(() => ts('pages.login.form.forgotPassword')),
      signIn: computed(() => ts('pages.login.form.signIn')),
    },

    /**
     * Validation message translations
     * Mirrors: pages.login.validation
     */
    validation: {
      emailRequired: computed(() => ts('pages.login.validation.emailRequired')),
      emailInvalid: computed(() => ts('pages.login.validation.emailInvalid')),
      passwordRequired: computed(() => ts('pages.login.validation.passwordRequired')),
    },

    /**
     * Error message translations
     * Mirrors: pages.login.errors
     */
    errors: {
      title: computed(() => ts('pages.login.errors.title')),
      default: computed(() => ts('pages.login.errors.default')),
    },

    /**
     * Language option translations
     * Mirrors: pages.login.languages
     */
    languages: {
      english: computed(() => ts('pages.login.languages.english')),
      portuguese: computed(() => ts('pages.login.languages.portuguese')),
    },

    /**
     * ARIA labels for accessibility
     * Mirrors: pages.login.aria
     */
    aria: {
      closeError: computed(() => ts('pages.login.aria.closeError')),
      languageSelector: computed(() => ts('pages.login.aria.languageSelector')),
    },
  };
}
