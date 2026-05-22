import { computed } from 'vue';
import { useTS } from '@utils/translation';

/**
 * i18n translations composable for EventFieldInput component
 *
 * Provides translated strings for:
 * - Type options (Event, State, Variable, Literal)
 * - Labels and placeholders
 * - Empty states
 * - Tooltips
 * - Context menu items
 */
export function useEventFieldInputTranslations() {
  const tsRaw = useTS();
  const tsTitle = useTS({ titleCase: true });

  return {
    typeOptions: {
      event: computed(() => tsTitle('components.forms.eventFieldInput.typeOptions.event')),
      state: computed(() => tsTitle('components.forms.eventFieldInput.typeOptions.state')),
      variable: computed(() => tsTitle('components.forms.eventFieldInput.typeOptions.variable')),
      literal: computed(() => tsTitle('components.forms.eventFieldInput.typeOptions.literal')),
    },
    labels: {
      eventField: computed(() => tsTitle('components.forms.eventFieldInput.labels.eventField')),
      stateName: computed(() => tsTitle('components.forms.eventFieldInput.labels.stateName')),
    },
    placeholders: {
      selectOrTypeState: computed(() => tsRaw('components.forms.eventFieldInput.placeholders.selectOrTypeState')),
    },
    empty: {
      noStates: computed(() => tsRaw('components.forms.eventFieldInput.empty.noStates')),
      goToLocalState: computed(() => tsRaw('components.forms.eventFieldInput.empty.goToLocalState')),
    },
    tooltips: {
      clear: computed(() => tsTitle('components.forms.eventFieldInput.tooltips.clear')),
      templatesSelected: computed(() => tsRaw('components.forms.eventFieldInput.tooltips.templatesSelected')),
      options: computed(() => tsTitle('components.forms.eventFieldInput.tooltips.options')),
    },
    menu: {
      selectField: {
        title: computed(() => tsTitle('components.forms.eventFieldInput.menu.selectField.title')),
        description: computed(() => tsRaw('components.forms.eventFieldInput.menu.selectField.description')),
      },
      searchTemplates: {
        title: computed(() => tsTitle('components.forms.eventFieldInput.menu.searchTemplates.title')),
        description: computed(() => tsRaw('components.forms.eventFieldInput.menu.searchTemplates.description')),
      },
      manualMode: {
        title: computed(() => tsTitle('components.forms.eventFieldInput.menu.manualMode.title')),
        description: computed(() => tsRaw('components.forms.eventFieldInput.menu.manualMode.description')),
      },
      dynamicMode: {
        title: computed(() => tsTitle('components.forms.eventFieldInput.menu.dynamicMode.title')),
        description: computed(() => tsRaw('components.forms.eventFieldInput.menu.dynamicMode.description')),
      },
    },
  };
}
