import { computed } from 'vue';
import { useTS } from '@utils/translation';

const BASE = 'pages.automations.createEditWorkflowInstance';

/**
 * Composable for CreateEditWorkflowInstance page translations.
 * Provides type-safe reactive access to all translated strings.
 *
 * Structure mirrors:
 * - JSON: src/i18n/{locale}/pages/automations/createEditWorkflowInstance.json
 * - Page: src/pages/automations/workflowInstances/createEditWorkflowInstancePage/
 */
export function useCreateEditWorkflowInstanceTranslations() {
  const ts = useTS({ capitalize: true });
  const tsTitle = useTS({ titleCase: true });
  const tsRaw = useTS({ capitalize: false });

  return {
    /** Page header translations */
    page: {
      title: computed(() => ts(`${BASE}.page.title`)),
      titleEdit: computed(() => ts(`${BASE}.page.titleEdit`)),
      description: computed(() => tsRaw(`${BASE}.page.description`)),
      button: computed(() => ts(`${BASE}.page.button`)),
    },

    /** Stepper header translations */
    stepper: {
      title: computed(() => tsTitle(`${BASE}.stepper.title`)),
      subtitle: computed(() => tsRaw(`${BASE}.stepper.subtitle`)),
      requiredInfo: computed(() => tsRaw(`${BASE}.stepper.requiredInfo`)),
      currentStep: computed(() => ts(`${BASE}.stepper.currentStep`)),
    },

    /** Step labels and descriptions */
    steps: {
      step1: {
        label: computed(() => ts(`${BASE}.steps.step1.label`)),
        description: computed(() => tsRaw(`${BASE}.steps.step1.description`)),
      },
      step2: {
        label: computed(() => ts(`${BASE}.steps.step2.label`)),
        description: computed(() => tsRaw(`${BASE}.steps.step2.description`)),
      },
      step3: {
        label: computed(() => ts(`${BASE}.steps.step3.label`)),
        description: computed(() => tsRaw(`${BASE}.steps.step3.description`)),
      },
      step4: {
        label: computed(() => ts(`${BASE}.steps.step4.label`)),
        description: computed(() => tsRaw(`${BASE}.steps.step4.description`)),
      },
    },

    /** Form field labels and placeholders */
    fields: {
      name: computed(() => ts(`${BASE}.fields.name`)),
      namePlaceholder: computed(() => tsRaw(`${BASE}.fields.namePlaceholder`)),
      description: computed(() => ts(`${BASE}.fields.description`)),
      descriptionPlaceholder: computed(() => tsRaw(`${BASE}.fields.descriptionPlaceholder`)),
      enabled: computed(() => ts(`${BASE}.fields.enabled`)),
      enabledHint: computed(() => tsRaw(`${BASE}.fields.enabledHint`)),
      isTemplate: computed(() => ts(`${BASE}.fields.isTemplate`)),
      isTemplateHint: computed(() => tsRaw(`${BASE}.fields.isTemplateHint`)),
      uniqueExecution: computed(() => ts(`${BASE}.fields.uniqueExecution`)),
      uniqueExecutionHint: computed(() => tsRaw(`${BASE}.fields.uniqueExecutionHint`)),
      workflowUUID: computed(() => ts(`${BASE}.fields.workflowUUID`)),
      workflowUUIDPlaceholder: computed(() => tsRaw(`${BASE}.fields.workflowUUIDPlaceholder`)),
      workflowUUIDHint: computed(() => tsRaw(`${BASE}.fields.workflowUUIDHint`)),
      generateUUID: computed(() => ts(`${BASE}.fields.generateUUID`)),
      allStatus: computed(() => ts(`${BASE}.fields.allStatus`)),
      active: computed(() => ts(`${BASE}.fields.active`)),
      inactive: computed(() => ts(`${BASE}.fields.inactive`)),
      definition: computed(() => ts(`${BASE}.fields.definition`)),
      definitionPlaceholder: computed(() => tsRaw(`${BASE}.fields.definitionPlaceholder`)),
      definitionDrawerTitle: computed(() => ts(`${BASE}.fields.definitionDrawerTitle`)),
      definitionDrawerSearch: computed(() => tsRaw(`${BASE}.fields.definitionDrawerSearch`)),
      noExternalInputs: computed(() => tsRaw(`${BASE}.fields.noExternalInputs`)),
      externalInputsTitle: computed(() => ts(`${BASE}.fields.externalInputsTitle`)),
      externalInputsDescription: computed(() => tsRaw(`${BASE}.fields.externalInputsDescription`)),
      requiredField: computed(() => tsRaw(`${BASE}.fields.requiredField`)),
      optionalField: computed(() => tsRaw(`${BASE}.fields.optionalField`)),
    },

    /** Review step labels */
    review: {
      identification: computed(() => ts(`${BASE}.review.identification`)),
      definition: computed(() => ts(`${BASE}.review.definition`)),
      externalInputs: computed(() => ts(`${BASE}.review.externalInputs`)),
      name: computed(() => ts(`${BASE}.review.name`)),
      description: computed(() => ts(`${BASE}.review.description`)),
      enabled: computed(() => ts(`${BASE}.review.enabled`)),
      isTemplate: computed(() => ts(`${BASE}.review.isTemplate`)),
      version: computed(() => ts(`${BASE}.review.version`)),
      status: computed(() => ts(`${BASE}.review.status`)),
      noInputs: computed(() => tsRaw(`${BASE}.review.noInputs`)),
      yes: computed(() => ts(`${BASE}.review.yes`)),
      no: computed(() => ts(`${BASE}.review.no`)),
      health: computed(() => ts(`${BASE}.review.health`)),
      nodes: computed(() => ts(`${BASE}.review.nodes`)),
      states: computed(() => ts(`${BASE}.review.states`)),
      inputs: computed(() => ts(`${BASE}.review.inputs`)),
      plugins: computed(() => ts(`${BASE}.review.plugins`)),
      coreOnly: computed(() => tsRaw(`${BASE}.review.coreOnly`)),
      deselect: computed(() => ts(`${BASE}.review.deselect`)),
      healthValid: computed(() => ts(`${BASE}.review.healthValid`)),
      healthPluginMissing: computed(() => ts(`${BASE}.review.healthPluginMissing`)),
      healthInvalid: computed(() => ts(`${BASE}.review.healthInvalid`)),
    },

    /** Navigation button labels */
    navigation: {
      previous: computed(() => ts(`${BASE}.navigation.previous`)),
      next: computed(() => ts(`${BASE}.navigation.next`)),
      save: computed(() => ts(`${BASE}.navigation.save`)),
      update: computed(() => ts(`${BASE}.navigation.update`)),
      cancel: computed(() => ts(`${BASE}.navigation.cancel`)),
    },

    /** Notification messages */
    notifications: {
      created: computed(() => ts(`${BASE}.notifications.created`)),
      updated: computed(() => ts(`${BASE}.notifications.updated`)),
      loading: computed(() => ts(`${BASE}.notifications.loading`)),
      loadFailed: computed(() => ts(`${BASE}.notifications.loadFailed`)),
      creationFailed: computed(() => ts(`${BASE}.notifications.creationFailed`)),
      updateFailed: computed(() => ts(`${BASE}.notifications.updateFailed`)),
    },
  };
}
