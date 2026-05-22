import { computed } from 'vue';

import { useTS } from '@utils/translation';

export function useAssetClassificationSelectorTranslations() {
  const ts = useTS({ capitalize: true });

  return {
    title: computed(() => ts('components.forms.assetClassificationSelector.title')),

    guided: {
      category: {
        label: computed(() => ts('components.forms.assetClassificationSelector.guided.category.label')),
        placeholder: computed(() => ts('components.forms.assetClassificationSelector.guided.category.placeholder')),
        hint: computed(() => ts('components.forms.assetClassificationSelector.guided.category.hint')),
        required: computed(() => ts('components.forms.assetClassificationSelector.guided.category.required')),
      },
      manufacturer: {
        label: computed(() => ts('components.forms.assetClassificationSelector.guided.manufacturer.label')),
        placeholder: computed(() => ts('components.forms.assetClassificationSelector.guided.manufacturer.placeholder')),
        hint: computed(() => ts('components.forms.assetClassificationSelector.guided.manufacturer.hint')),
        disabled: computed(() => ts('components.forms.assetClassificationSelector.guided.manufacturer.disabled')),
        required: computed(() => ts('components.forms.assetClassificationSelector.guided.manufacturer.required')),
      },
      model: {
        label: computed(() => ts('components.forms.assetClassificationSelector.guided.model.label')),
        placeholder: computed(() => ts('components.forms.assetClassificationSelector.guided.model.placeholder')),
        hint: computed(() => ts('components.forms.assetClassificationSelector.guided.model.hint')),
        disabled: computed(() => ts('components.forms.assetClassificationSelector.guided.model.disabled')),
        required: computed(() => ts('components.forms.assetClassificationSelector.guided.model.required')),
      },
      version: {
        label: computed(() => ts('components.forms.assetClassificationSelector.guided.version.label')),
        placeholder: computed(() => ts('components.forms.assetClassificationSelector.guided.version.placeholder')),
        hint: computed(() => ts('components.forms.assetClassificationSelector.guided.version.hint')),
        required: computed(() => ts('components.forms.assetClassificationSelector.guided.version.required')),
      },
      loading: computed(() => ts('components.forms.assetClassificationSelector.guided.loading')),
      noOptions: computed(() => ts('components.forms.assetClassificationSelector.guided.noOptions')),
    },

    errors: {
      loadFailed: computed(() => ts('components.forms.assetClassificationSelector.errors.loadFailed')),
    },
  };
}
