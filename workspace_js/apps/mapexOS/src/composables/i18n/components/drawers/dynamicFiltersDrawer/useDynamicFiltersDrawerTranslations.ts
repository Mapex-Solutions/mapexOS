import { computed } from 'vue';
import { useTS } from '@utils/translation';

/**
 * Translations composable for DynamicFiltersDrawer component
 * Provides all translated strings for the dynamic filters drawer
 */
export function useDynamicFiltersDrawerTranslations() {
  const ts = useTS({ capitalize: true });
  const tsRaw = useTS({ capitalize: false });

  return {
    /** Drawer title */
    title: computed(() => ts('components.drawers.dynamicFiltersDrawer.title')),

    /** Close button tooltip */
    closeTooltip: computed(() => ts('components.drawers.dynamicFiltersDrawer.closeTooltip')),

    /** Source selector label */
    sourceLabel: computed(() => ts('components.drawers.dynamicFiltersDrawer.sourceLabel')),

    /** Source type labels */
    sourceAsset: computed(() => ts('components.drawers.dynamicFiltersDrawer.sourceAsset')),
    sourceAssetTemplate: computed(() => ts('components.drawers.dynamicFiltersDrawer.sourceAssetTemplate')),

    /** Search placeholders */
    searchAsset: computed(() => ts('components.drawers.dynamicFiltersDrawer.searchAsset')),
    searchAssetTemplate: computed(() => ts('components.drawers.dynamicFiltersDrawer.searchAssetTemplate')),

    /** No results message */
    noResults: computed(() => ts('components.drawers.dynamicFiltersDrawer.noResults')),

    /** Template resolved label */
    templateResolved: computed(() => ts('components.drawers.dynamicFiltersDrawer.templateResolved')),

    /** Loading fields message */
    loadingFields: computed(() => ts('components.drawers.dynamicFiltersDrawer.loadingFields')),

    /** No fields message */
    noFields: computed(() => tsRaw('components.drawers.dynamicFiltersDrawer.noFields')),

    /** No template assigned to the selected asset */
    noTemplateForAsset: computed(() => tsRaw('components.drawers.dynamicFiltersDrawer.noTemplateForAsset')),

    /** Resolving template message */
    resolvingTemplate: computed(() => ts('components.drawers.dynamicFiltersDrawer.resolvingTemplate')),

    /** Dynamic fields section title */
    dynamicFieldsTitle: computed(() => ts('components.drawers.dynamicFiltersDrawer.dynamicFieldsTitle')),

    /** Filters section label */
    filtersSection: computed(() => ts('components.drawers.dynamicFiltersDrawer.filtersSection')),

    /** Filters count template */
    filtersCount: computed(() => tsRaw('components.drawers.dynamicFiltersDrawer.filtersCount')),

    /** Add filter button label */
    addFilter: computed(() => ts('components.drawers.dynamicFiltersDrawer.addFilter')),

    /** Search fields placeholder */
    searchFields: computed(() => tsRaw('components.drawers.dynamicFiltersDrawer.searchFields')),

    /** Empty state title */
    emptyTitle: computed(() => ts('components.drawers.dynamicFiltersDrawer.emptyTitle')),

    /** Empty state description */
    emptyDescription: computed(() => tsRaw('components.drawers.dynamicFiltersDrawer.emptyDescription')),

    /** All fields added message */
    allFieldsAdded: computed(() => ts('components.drawers.dynamicFiltersDrawer.allFieldsAdded')),

    /** Field type group headers */
    fieldTypeHeaders: {
      number: computed(() => ts('components.drawers.dynamicFiltersDrawer.fieldTypeHeaders.number')),
      string: computed(() => ts('components.drawers.dynamicFiltersDrawer.fieldTypeHeaders.string')),
      boolean: computed(() => ts('components.drawers.dynamicFiltersDrawer.fieldTypeHeaders.boolean')),
      date: computed(() => ts('components.drawers.dynamicFiltersDrawer.fieldTypeHeaders.date')),
    },

    /** Operator labels */
    operators: {
      equals: computed(() => ts('components.drawers.dynamicFiltersDrawer.operators.equals')),
      notEquals: computed(() => ts('components.drawers.dynamicFiltersDrawer.operators.notEquals')),
      greaterThan: computed(() => ts('components.drawers.dynamicFiltersDrawer.operators.greaterThan')),
      greaterThanEquals: computed(() => ts('components.drawers.dynamicFiltersDrawer.operators.greaterThanEquals')),
      lessThan: computed(() => ts('components.drawers.dynamicFiltersDrawer.operators.lessThan')),
      lessThanEquals: computed(() => ts('components.drawers.dynamicFiltersDrawer.operators.lessThanEquals')),
      range: computed(() => ts('components.drawers.dynamicFiltersDrawer.operators.range')),
      startsWith: computed(() => ts('components.drawers.dynamicFiltersDrawer.operators.startsWith')),
      rangeFrom: computed(() => ts('components.drawers.dynamicFiltersDrawer.operators.rangeFrom')),
      rangeTo: computed(() => ts('components.drawers.dynamicFiltersDrawer.operators.rangeTo')),
    },

    /** Boolean labels */
    booleanTrue: computed(() => ts('components.drawers.dynamicFiltersDrawer.booleanTrue')),
    booleanFalse: computed(() => ts('components.drawers.dynamicFiltersDrawer.booleanFalse')),

    /** Button labels and tooltips */
    buttons: {
      reset: computed(() => ts('components.drawers.dynamicFiltersDrawer.buttons.reset')),
      resetTooltip: computed(() => ts('components.drawers.dynamicFiltersDrawer.buttons.resetTooltip')),
      apply: computed(() => ts('components.drawers.dynamicFiltersDrawer.buttons.apply')),
      applyTooltip: computed(() => ts('components.drawers.dynamicFiltersDrawer.buttons.applyTooltip')),
    },
  };
}
