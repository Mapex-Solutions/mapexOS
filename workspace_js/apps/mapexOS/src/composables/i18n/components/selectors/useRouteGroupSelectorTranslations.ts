import { computed } from 'vue';
import { useI18n } from 'vue-i18n';

export function useRouteGroupSelectorTranslations() {
  const { t } = useI18n();

  return {
    title: {
      singular: computed(() => t('components.selectors.routeGroupSelector.title.singular')),
      plural: computed(() => t('components.selectors.routeGroupSelector.title.plural')),
    },

    selectedCount: computed(() => t('components.selectors.routeGroupSelector.selectedCount')),

    search: {
      label: computed(() => t('components.selectors.routeGroupSelector.search.label')),
      placeholder: computed(() => t('components.selectors.routeGroupSelector.search.placeholder')),
    },

    filters: {
      title: computed(() => t('components.selectors.routeGroupSelector.filters.title')),
      templateType: computed(() => t('components.selectors.routeGroupSelector.filters.templateType')),
      templateSource: computed(() => t('components.selectors.routeGroupSelector.filters.templateSource')),
      search: {
        label: computed(() => t('components.selectors.routeGroupSelector.filters.search.label')),
        placeholder: computed(() => t('components.selectors.routeGroupSelector.filters.search.placeholder')),
      },
      templateTypeOptions: {
        all: computed(() => t('components.selectors.routeGroupSelector.filters.templateTypeOptions.all')),
        system: computed(() => t('components.selectors.routeGroupSelector.filters.templateTypeOptions.system')),
        custom: computed(() => t('components.selectors.routeGroupSelector.filters.templateTypeOptions.custom')),
      },
      templateSourceOptions: {
        all: computed(() => t('components.selectors.routeGroupSelector.filters.templateSourceOptions.all')),
        shared: computed(() => t('components.selectors.routeGroupSelector.filters.templateSourceOptions.shared')),
        local: computed(() => t('components.selectors.routeGroupSelector.filters.templateSourceOptions.local')),
      },
    },

    resultsTitle: computed(() => t('components.selectors.routeGroupSelector.resultsTitle')),

    results: {
      found: computed(() => t('components.selectors.routeGroupSelector.results.found')),
    },

    chips: {
      system: computed(() => t('components.selectors.routeGroupSelector.chips.system')),
      shared: computed(() => t('components.selectors.routeGroupSelector.chips.shared')),
      routers: computed(() => t('components.selectors.routeGroupSelector.chips.routers')),
    },

    status: {
      active: computed(() => t('components.selectors.routeGroupSelector.status.active')),
      inactive: computed(() => t('components.selectors.routeGroupSelector.status.inactive')),
    },

    badge: {
      active: computed(() => t('components.selectors.routeGroupSelector.badge.active')),
      inactive: computed(() => t('components.selectors.routeGroupSelector.badge.inactive')),
    },

    totalLoaded: computed(() => t('components.selectors.routeGroupSelector.totalLoaded')),

    routers: computed(() => t('components.selectors.routeGroupSelector.routers')),
    noDescription: computed(() => t('components.selectors.routeGroupSelector.noDescription')),

    loading: {
      value: computed(() => t('components.selectors.routeGroupSelector.loading')),
      more: computed(() => t('components.selectors.routeGroupSelector.loading.more')),
    },
    empty: computed(() => t('components.selectors.routeGroupSelector.empty')),

    emptyState: {
      title: computed(() => t('components.selectors.routeGroupSelector.emptyState.title')),
      subtitle: computed(() => t('components.selectors.routeGroupSelector.emptyState.subtitle')),
    },

    selected: {
      label: computed(() => t('components.selectors.routeGroupSelector.selected.label')),
    },

    actions: {
      cancel: computed(() => t('components.selectors.routeGroupSelector.actions.cancel')),
      confirm: computed(() => t('components.selectors.routeGroupSelector.actions.confirm')),
    },

    errors: {
      loadFailed: computed(() => t('components.selectors.routeGroupSelector.errors.loadFailed')),
    },
  };
}
