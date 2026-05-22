import { computed } from 'vue';
import { useTS } from '@utils/translation';

/**
 * Change password page translations
 *
 * Structure mirrors:
 * - File: src/pages/changePassword/ChangePasswordPage.vue
 * - JSON: src/i18n/{locale}/pages/changePassword.json
 * - Composable: src/composables/i18n/pages/changePassword/useChangePasswordTranslations.ts
 *
 * Provides all translations for the Change Password page including:
 * - Welcome section (title, tagline, subtitle)
 * - Password tips (features)
 * - Form labels and placeholders
 * - Validation messages
 * - Error and success messages
 * - ARIA labels for accessibility
 *
 * @example
 * ```ts
 * // In ChangePasswordPage.vue
 * const {
 *   welcome,
 *   features,
 *   form,
 *   validation,
 *   errors,
 *   success
 * } = useChangePasswordTranslations();
 *
 * <h2>{{ welcome.title.value }}</h2>
 * <q-input :label="form.newPassword.label.value" />
 * ```
 */
export function useChangePasswordTranslations() {
  const ts = useTS({ capitalize: true });

  return {
    /**
     * Welcome section translations
     * Mirrors: pages.changePassword.welcome
     */
    welcome: {
      title: computed(() => ts('pages.changePassword.welcome.title')),
      tagline: computed(() => ts('pages.changePassword.welcome.tagline')),
      subtitle: computed(() => ts('pages.changePassword.welcome.subtitle')),
    },

    /**
     * Password tips translations
     * Mirrors: pages.changePassword.features
     */
    features: {
      tip1: computed(() => ts('pages.changePassword.features.tip1')),
      tip2: computed(() => ts('pages.changePassword.features.tip2')),
      tip3: computed(() => ts('pages.changePassword.features.tip3')),
    },

    /**
     * Form field translations
     * Mirrors: pages.changePassword.form
     */
    form: {
      title: computed(() => ts('pages.changePassword.form.title')),

      newPassword: {
        label: computed(() => ts('pages.changePassword.form.newPassword.label')),
        placeholder: computed(() => ts('pages.changePassword.form.newPassword.placeholder')),
      },

      confirmPassword: {
        label: computed(() => ts('pages.changePassword.form.confirmPassword.label')),
        placeholder: computed(() => ts('pages.changePassword.form.confirmPassword.placeholder')),
      },

      submit: computed(() => ts('pages.changePassword.form.submit')),
    },

    /**
     * Validation message translations
     * Mirrors: pages.changePassword.validation
     */
    validation: {
      passwordRequired: computed(() => ts('pages.changePassword.validation.passwordRequired')),
      passwordMinLength: computed(() => ts('pages.changePassword.validation.passwordMinLength')),
      passwordMaxLength: computed(() => ts('pages.changePassword.validation.passwordMaxLength')),
      passwordsMustMatch: computed(() => ts('pages.changePassword.validation.passwordsMustMatch')),
    },

    /**
     * Error message translations
     * Mirrors: pages.changePassword.errors
     */
    errors: {
      title: computed(() => ts('pages.changePassword.errors.title')),
      default: computed(() => ts('pages.changePassword.errors.default')),
    },

    /**
     * Success message translations
     * Mirrors: pages.changePassword.success
     */
    success: {
      title: computed(() => ts('pages.changePassword.success.title')),
      message: computed(() => ts('pages.changePassword.success.message')),
    },

    /**
     * ARIA labels for accessibility
     * Mirrors: pages.changePassword.aria
     */
    aria: {
      closeError: computed(() => ts('pages.changePassword.aria.closeError')),
    },
  };
}
