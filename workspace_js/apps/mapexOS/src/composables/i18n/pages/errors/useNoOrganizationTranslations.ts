import { computed } from 'vue';
import { useTS } from '@utils/translation';

export function useNoOrganizationTranslations() {
  const ts = useTS({ capitalize: true });
  const tsTitle = useTS({ titleCase: true });

  return {
    page: {
      title: computed(() => tsTitle('pages.errors.noOrganization.title')),
      description: computed(() => ts('pages.errors.noOrganization.description')),
      message: computed(() => ts('pages.errors.noOrganization.message')),
    },

    actions: {
      logout: computed(() => ts('pages.errors.noOrganization.actions.logout')),
      support: computed(() => ts('pages.errors.noOrganization.actions.support')),
    },
  };
}
