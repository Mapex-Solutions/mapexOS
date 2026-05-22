import { computed } from 'vue';
import { useI18n } from 'vue-i18n';

export function useAssetTemplateSelectorTranslations() {
  const { t } = useI18n();

  return {
    search: {
      placeholder: computed(() => t('components.selectors.assetTemplateSelector.search.placeholder')),
    },
    filters: {
      manufacturer: computed(() => t('components.selectors.assetTemplateSelector.filters.manufacturer')),
      model: computed(() => t('components.selectors.assetTemplateSelector.filters.model')),
      status: computed(() => t('components.selectors.assetTemplateSelector.filters.status')),
      templateType: computed(() => t('components.selectors.assetTemplateSelector.filters.templateType')),
      templateSource: computed(() => t('components.selectors.assetTemplateSelector.filters.templateSource')),
      statusOptions: {
        all: computed(() => t('components.selectors.assetTemplateSelector.filters.statusOptions.all')),
        active: computed(() => t('components.selectors.assetTemplateSelector.filters.statusOptions.active')),
        inactive: computed(() => t('components.selectors.assetTemplateSelector.filters.statusOptions.inactive')),
      },
      templateTypeOptions: {
        all: computed(() => t('components.selectors.assetTemplateSelector.filters.templateTypeOptions.all')),
        system: computed(() => t('components.selectors.assetTemplateSelector.filters.templateTypeOptions.system')),
        custom: computed(() => t('components.selectors.assetTemplateSelector.filters.templateTypeOptions.custom')),
      },
      templateSourceOptions: {
        all: computed(() => t('components.selectors.assetTemplateSelector.filters.templateSourceOptions.all')),
        shared: computed(() => t('components.selectors.assetTemplateSelector.filters.templateSourceOptions.shared')),
        local: computed(() => t('components.selectors.assetTemplateSelector.filters.templateSourceOptions.local')),
      },
    },
    results: {
      found: computed(() => t('components.selectors.assetTemplateSelector.results.found')),
    },
    badge: {
      active: computed(() => t('components.selectors.assetTemplateSelector.badge.active')),
      inactive: computed(() => t('components.selectors.assetTemplateSelector.badge.inactive')),
    },
    noDescription: computed(() => t('components.selectors.assetTemplateSelector.noDescription')),
    loading: {
      more: computed(() => t('components.selectors.assetTemplateSelector.loading.more')),
    },
    emptyState: {
      title: computed(() => t('components.selectors.assetTemplateSelector.emptyState.title')),
      subtitle: computed(() => t('components.selectors.assetTemplateSelector.emptyState.subtitle')),
    },
    selected: {
      label: computed(() => t('components.selectors.assetTemplateSelector.selected.label')),
    },
    errors: {
      loadFailed: computed(() => t('components.selectors.assetTemplateSelector.errors.loadFailed')),
    },
  };
}
