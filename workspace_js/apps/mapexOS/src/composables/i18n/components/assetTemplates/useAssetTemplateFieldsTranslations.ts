import { computed } from 'vue';
import { useI18n } from 'vue-i18n';

/**
 * Translations for the read-only asset-template field components
 * (dynamic fields table and available fields list).
 */
export function useAssetTemplateFieldsTranslations() {
  const { t } = useI18n();

  return {
    dynamicFields: {
      headers: {
        field: computed(() => t('components.assetTemplates.assetTemplateFields.dynamicFields.headers.field')),
        type: computed(() => t('components.assetTemplates.assetTemplateFields.dynamicFields.headers.type')),
        value: computed(() => t('components.assetTemplates.assetTemplateFields.dynamicFields.headers.value')),
        status: computed(() => t('components.assetTemplates.assetTemplateFields.dynamicFields.headers.status')),
      },
      status: {
        active: computed(() => t('components.assetTemplates.assetTemplateFields.dynamicFields.status.active')),
        deprecated: computed(() => t('components.assetTemplates.assetTemplateFields.dynamicFields.status.deprecated')),
      },
      empty: computed(() => t('components.assetTemplates.assetTemplateFields.dynamicFields.empty')),
    },
    availableFields: {
      empty: computed(() => t('components.assetTemplates.assetTemplateFields.availableFields.empty')),
    },
    count: (count: number) => t('components.assetTemplates.assetTemplateFields.count', { count }),
  };
}
