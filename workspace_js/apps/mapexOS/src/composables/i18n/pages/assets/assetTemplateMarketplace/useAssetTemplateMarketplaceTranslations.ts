import { computed } from 'vue';
import { useTS } from '@utils/translation';

/**
 * Reactive translations for the asset template marketplace page, cards, detail modal, and install flow.
 * Every key maps to pages.assets.assetTemplateMarketplace in the locale JSON, keeping en-US and pt-BR in sync.
 */
export function useAssetTemplateMarketplaceTranslations() {
  const ts = useTS({ capitalize: true });
  const tsTitle = useTS({ titleCase: true });

  return {
    page: {
      title: computed(() => tsTitle('pages.assets.assetTemplateMarketplace.page.title')),
      subtitle: computed(() => ts('pages.assets.assetTemplateMarketplace.page.subtitle')),
    },

    empty: {
      title: computed(() => ts('pages.assets.assetTemplateMarketplace.empty.title')),
      description: computed(() => ts('pages.assets.assetTemplateMarketplace.empty.description')),
      button: computed(() => ts('pages.assets.assetTemplateMarketplace.empty.button')),
    },

    error: {
      title: computed(() => ts('pages.assets.assetTemplateMarketplace.error.title')),
      description: computed(() => ts('pages.assets.assetTemplateMarketplace.error.description')),
      button: computed(() => ts('pages.assets.assetTemplateMarketplace.error.button')),
    },

    search: {
      placeholder: computed(() => ts('pages.assets.assetTemplateMarketplace.search.placeholder')),
    },

    facets: {
      category: computed(() => ts('pages.assets.assetTemplateMarketplace.facets.category')),
      vendor: computed(() => ts('pages.assets.assetTemplateMarketplace.facets.vendor')),
      model: computed(() => ts('pages.assets.assetTemplateMarketplace.facets.model')),
      version: computed(() => ts('pages.assets.assetTemplateMarketplace.facets.version')),
      all: computed(() => ts('pages.assets.assetTemplateMarketplace.facets.all')),
    },

    card: {
      /**
       * Field-count label for a catalog card, e.g. "12 fields".
       * @param {number} count - Number of dynamic fields in the template.
       * @returns {string} Localized field-count label.
       */
      fields: (count: number): string =>
        ts('pages.assets.assetTemplateMarketplace.card.fields', { count }),
      hasScripts: computed(() => ts('pages.assets.assetTemplateMarketplace.card.hasScripts')),
      view: computed(() => ts('pages.assets.assetTemplateMarketplace.card.view')),
    },

    modal: {
      title: computed(() => ts('pages.assets.assetTemplateMarketplace.modal.title')),
      loading: computed(() => ts('pages.assets.assetTemplateMarketplace.modal.loading')),
      loadError: computed(() => ts('pages.assets.assetTemplateMarketplace.modal.loadError')),
      sections: {
        overview: computed(() => ts('pages.assets.assetTemplateMarketplace.modal.sections.overview')),
        dynamicFields: computed(() => ts('pages.assets.assetTemplateMarketplace.modal.sections.dynamicFields')),
        scripts: computed(() => ts('pages.assets.assetTemplateMarketplace.modal.sections.scripts')),
        availableFields: computed(() => ts('pages.assets.assetTemplateMarketplace.modal.sections.availableFields')),
      },
      overview: {
        manufacturer: computed(() => ts('pages.assets.assetTemplateMarketplace.modal.overview.manufacturer')),
        model: computed(() => ts('pages.assets.assetTemplateMarketplace.modal.overview.model')),
        version: computed(() => ts('pages.assets.assetTemplateMarketplace.modal.overview.version')),
        category: computed(() => ts('pages.assets.assetTemplateMarketplace.modal.overview.category')),
        assetIdPath: computed(() => ts('pages.assets.assetTemplateMarketplace.modal.overview.assetIdPath')),
        description: computed(() => ts('pages.assets.assetTemplateMarketplace.modal.overview.description')),
      },
      scripts: {
        test: computed(() => ts('pages.assets.assetTemplateMarketplace.modal.scripts.test')),
        processor: computed(() => ts('pages.assets.assetTemplateMarketplace.modal.scripts.processor')),
        validator: computed(() => ts('pages.assets.assetTemplateMarketplace.modal.scripts.validator')),
        conversion: computed(() => ts('pages.assets.assetTemplateMarketplace.modal.scripts.conversion')),
        configured: computed(() => ts('pages.assets.assetTemplateMarketplace.modal.scripts.configured')),
        notConfigured: computed(() => ts('pages.assets.assetTemplateMarketplace.modal.scripts.notConfigured')),
        view: computed(() => ts('pages.assets.assetTemplateMarketplace.modal.scripts.view')),
      },
      actions: {
        cancel: computed(() => ts('pages.assets.assetTemplateMarketplace.modal.actions.cancel')),
        close: computed(() => ts('pages.assets.assetTemplateMarketplace.modal.actions.close')),
      },
    },

    install: {
      button: computed(() => ts('pages.assets.assetTemplateMarketplace.install.button')),
      installing: computed(() => ts('pages.assets.assetTemplateMarketplace.install.installing')),
      installed: computed(() => ts('pages.assets.assetTemplateMarketplace.install.installed')),
      success: computed(() => ts('pages.assets.assetTemplateMarketplace.install.success')),
      checksumError: computed(() => ts('pages.assets.assetTemplateMarketplace.install.checksumError')),
      genericError: computed(() => ts('pages.assets.assetTemplateMarketplace.install.genericError')),
      shareWithChildren: computed(() => ts('pages.assets.assetTemplateMarketplace.install.shareWithChildren')),
      shareWithChildrenHint: computed(() => ts('pages.assets.assetTemplateMarketplace.install.shareWithChildrenHint')),
    },
  };
}
