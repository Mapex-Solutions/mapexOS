import { computed } from 'vue';
import { useTS } from '@utils/translation';

/**
 * Onboarding translations
 *
 * Structure mirrors:
 * - JSON: src/i18n/{locale}/composables/onboarding.json
 * - Composable: src/composables/i18n/composables/useOnboardingTranslations.ts
 *
 * Provides all translations for the onboarding tour including:
 * - Step titles and descriptions
 * - Button labels (next, previous, finish, skip)
 * - Menu item label
 *
 * @example
 * ```ts
 * const t = useOnboardingTranslations();
 * console.log(t.steps.sidebar.title.value); // "Navigation menu"
 * console.log(t.buttons.next.value); // "Next"
 * ```
 */
export function useOnboardingTranslations() {
  const ts = useTS({ capitalize: true });

  return {
    /**
     * Tour step translations
     * Mirrors: composables.onboarding.steps
     */
    steps: {
      sidebar: {
        title: computed(() => ts('composables.onboarding.steps.sidebar.title')),
        description: computed(() => ts('composables.onboarding.steps.sidebar.description')),
      },
      breadcrumbs: {
        title: computed(() => ts('composables.onboarding.steps.breadcrumbs.title')),
        description: computed(() => ts('composables.onboarding.steps.breadcrumbs.description')),
      },
      orgSelector: {
        title: computed(() => ts('composables.onboarding.steps.orgSelector.title')),
        description: computed(() => ts('composables.onboarding.steps.orgSelector.description')),
      },
      langSelector: {
        title: computed(() => ts('composables.onboarding.steps.langSelector.title')),
        description: computed(() => ts('composables.onboarding.steps.langSelector.description')),
      },
      userMenu: {
        title: computed(() => ts('composables.onboarding.steps.userMenu.title')),
        description: computed(() => ts('composables.onboarding.steps.userMenu.description')),
      },
    },

    /**
     * Button translations
     * Mirrors: composables.onboarding.buttons
     */
    buttons: {
      next: computed(() => ts('composables.onboarding.buttons.next')),
      previous: computed(() => ts('composables.onboarding.buttons.previous')),
      finish: computed(() => ts('composables.onboarding.buttons.finish')),
      continue: computed(() => ts('composables.onboarding.buttons.continue')),
      skip: computed(() => ts('composables.onboarding.buttons.skip')),
    },

    /**
     * Menu translations
     * Mirrors: composables.onboarding.menu
     */
    menu: {
      startTour: computed(() => ts('composables.onboarding.menu.startTour')),
    },
  };
}
