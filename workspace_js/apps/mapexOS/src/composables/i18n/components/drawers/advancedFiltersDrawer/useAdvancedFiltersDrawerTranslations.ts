import { computed } from 'vue';
import { useTS } from '@utils/translation';

/**
 * Translations composable for AdvancedFiltersDrawer component
 * Provides all translated strings for the advanced filters drawer
 */
export function useAdvancedFiltersDrawerTranslations() {
  const ts = useTS({ capitalize: true });

  return {
    /** Drawer title */
    title: computed(() => ts('components.drawers.advancedFiltersDrawer.title')),

    /** Close button tooltip */
    closeTooltip: computed(() => ts('components.drawers.advancedFiltersDrawer.closeTooltip')),

    /** Button labels and tooltips */
    buttons: {
      reset: computed(() => ts('components.drawers.advancedFiltersDrawer.buttons.reset')),
      resetTooltip: computed(() => ts('components.drawers.advancedFiltersDrawer.buttons.resetTooltip')),
      apply: computed(() => ts('components.drawers.advancedFiltersDrawer.buttons.apply')),
      applyTooltip: computed(() => ts('components.drawers.advancedFiltersDrawer.buttons.applyTooltip')),
    },

    /** Autocomplete field translations */
    autocomplete: {
      noOption: computed(() => ts('components.drawers.advancedFiltersDrawer.autocomplete.noOption')),
    },
  };
}
