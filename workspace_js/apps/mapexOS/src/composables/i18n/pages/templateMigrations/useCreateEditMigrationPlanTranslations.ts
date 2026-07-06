import { computed } from 'vue';
import { useTS } from '@utils/translation';

/**
 * Create/edit migration plan wizard translations.
 *
 * Structure mirrors:
 * - Page: src/pages/assets/templateMigrations/createEditMigrationPlanPage/CreateEditMigrationPlanPage.vue
 * - JSON: src/i18n/{locale}/pages/templateMigrations.json (createPlan section)
 */
export function useCreateEditMigrationPlanTranslations() {
  const ts = useTS({ capitalize: true });
  const tsTitle = useTS({ titleCase: true });

  const base = 'pages.templateMigrations.createPlan';

  return {
    page: {
      title: computed(() => tsTitle(`${base}.page.title`)),
      description: computed(() => ts(`${base}.page.description`)),
      editTitle: computed(() => tsTitle(`${base}.page.editTitle`)),
      editDescription: computed(() => ts(`${base}.page.editDescription`)),
      back: computed(() => ts(`${base}.page.back`)),
    },

    stepper: {
      title: computed(() => tsTitle(`${base}.stepper.title`)),
      subtitle: computed(() => ts(`${base}.stepper.subtitle`)),
      requiredInfo: computed(() => ts(`${base}.stepper.requiredInfo`)),
      currentStep: computed(() => ts(`${base}.stepper.currentStep`)),
    },

    navigation: {
      previous: computed(() => ts(`${base}.navigation.previous`)),
      next: computed(() => ts(`${base}.navigation.next`)),
      save: computed(() => ts(`${base}.navigation.save`)),
    },

    groups: {
      setup: computed(() => tsTitle(`${base}.groups.setup`)),
      migration: computed(() => tsTitle(`${base}.groups.migration`)),
    },

    steps: {
      details: {
        label: computed(() => tsTitle(`${base}.steps.details.label`)),
        description: computed(() => ts(`${base}.steps.details.description`)),
        title: computed(() => tsTitle(`${base}.steps.details.title`)),
        subtitle: computed(() => ts(`${base}.steps.details.subtitle`)),
        fields: {
          name: {
            label: computed(() => ts(`${base}.steps.details.fields.name.label`)),
            placeholder: computed(() => ts(`${base}.steps.details.fields.name.placeholder`)),
            hint: computed(() => ts(`${base}.steps.details.fields.name.hint`)),
            required: computed(() => ts(`${base}.steps.details.fields.name.required`)),
          },
          description: {
            label: computed(() => ts(`${base}.steps.details.fields.description.label`)),
            placeholder: computed(() => ts(`${base}.steps.details.fields.description.placeholder`)),
            hint: computed(() => ts(`${base}.steps.details.fields.description.hint`)),
          },
        },
      },
      schedule: {
        label: computed(() => tsTitle(`${base}.steps.schedule.label`)),
        description: computed(() => ts(`${base}.steps.schedule.description`)),
        title: computed(() => tsTitle(`${base}.steps.schedule.title`)),
        subtitle: computed(() => ts(`${base}.steps.schedule.subtitle`)),
        modeLabel: computed(() => ts(`${base}.steps.schedule.modeLabel`)),
        options: {
          now: computed(() => ts(`${base}.steps.schedule.options.now`)),
          scheduled: computed(() => ts(`${base}.steps.schedule.options.scheduled`)),
        },
        field: {
          label: computed(() => ts(`${base}.steps.schedule.field.label`)),
          hint: computed(() => ts(`${base}.steps.schedule.field.hint`)),
          placeholder: computed(() => ts(`${base}.steps.schedule.field.placeholder`)),
          required: computed(() => ts(`${base}.steps.schedule.field.required`)),
          future: computed(() => ts(`${base}.steps.schedule.field.future`)),
          close: computed(() => ts(`${base}.steps.schedule.field.close`)),
        },
      },
      fromTemplate: {
        label: computed(() => tsTitle(`${base}.steps.fromTemplate.label`)),
        description: computed(() => ts(`${base}.steps.fromTemplate.description`)),
        title: computed(() => tsTitle(`${base}.steps.fromTemplate.title`)),
        subtitle: computed(() => ts(`${base}.steps.fromTemplate.subtitle`)),
      },
      assets: {
        label: computed(() => tsTitle(`${base}.steps.assets.label`)),
        description: computed(() => ts(`${base}.steps.assets.description`)),
        title: computed(() => tsTitle(`${base}.steps.assets.title`)),
        subtitle: computed(() => ts(`${base}.steps.assets.subtitle`)),
        selectAll: (total: number) => ts(`${base}.steps.assets.selectAll`, { total }),
        selectPage: computed(() => ts(`${base}.steps.assets.selectPage`)),
        clearSelection: computed(() => ts(`${base}.steps.assets.clearSelection`)),
        counter: (selected: number, total: number) => ts(`${base}.steps.assets.counter`, { selected, total }),
        columns: {
          name: computed(() => ts(`${base}.steps.assets.columns.name`)),
          identifier: computed(() => ts(`${base}.steps.assets.columns.identifier`)),
        },
        empty: computed(() => ts(`${base}.steps.assets.empty`)),
        loadError: computed(() => ts(`${base}.steps.assets.loadError`)),
      },
      toTemplate: {
        label: computed(() => tsTitle(`${base}.steps.toTemplate.label`)),
        description: computed(() => ts(`${base}.steps.toTemplate.description`)),
        title: computed(() => tsTitle(`${base}.steps.toTemplate.title`)),
        subtitle: computed(() => ts(`${base}.steps.toTemplate.subtitle`)),
      },
    },

    notifications: {
      createSuccess: computed(() => ts(`${base}.notifications.createSuccess`)),
      createError: computed(() => ts(`${base}.notifications.createError`)),
      updateSuccess: computed(() => ts(`${base}.notifications.updateSuccess`)),
      updateError: computed(() => ts(`${base}.notifications.updateError`)),
      loadError: computed(() => ts(`${base}.notifications.loadError`)),
      notEditable: computed(() => ts(`${base}.notifications.notEditable`)),
      apiNotInitialized: computed(() => ts(`${base}.notifications.apiNotInitialized`)),
    },
  };
}
