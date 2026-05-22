import { computed } from 'vue';
import { useTS } from '@utils/translation';

/**
 * Add/Edit Customer page translations
 *
 * Structure mirrors:
 * - File: src/pages/administrations/customers/createEditCustomerPage/CreateEditCustomerPage.vue
 * - JSON: src/i18n/{locale}/pages/administrations/addCustomer.json
 * - Composable: src/composables/i18n/pages/administrations/customers/useAddCustomerTranslations.ts
 */
export function useAddCustomerTranslations() {
  const ts = useTS({ capitalize: true });
  const tsTitle = useTS({ titleCase: true });

  return {
    page: {
      title: computed(() => tsTitle('pages.administrations.addCustomer.title')),
      titleEdit: computed(() => tsTitle('pages.administrations.addCustomer.titleEdit')),
      description: computed(() => ts('pages.administrations.addCustomer.description')),
      descriptionEdit: computed(() => ts('pages.administrations.addCustomer.descriptionEdit')),
      backButton: computed(() => ts('pages.administrations.addCustomer.backButton')),
    },

    steps: {
      basic: computed(() => tsTitle('pages.administrations.addCustomer.steps.basic')),
      address: computed(() => tsTitle('pages.administrations.addCustomer.steps.address')),
      accessPolicy: computed(() => tsTitle('pages.administrations.addCustomer.steps.accessPolicy')),
      review: computed(() => tsTitle('pages.administrations.addCustomer.steps.review')),
    },

    stepDescriptions: {
      basic: computed(() => ts('pages.administrations.addCustomer.stepDescriptions.basic')),
      address: computed(() => ts('pages.administrations.addCustomer.stepDescriptions.address')),
      accessPolicy: computed(() => ts('pages.administrations.addCustomer.stepDescriptions.accessPolicy')),
      review: computed(() => ts('pages.administrations.addCustomer.stepDescriptions.review')),
    },

    formDescriptions: {
      basic: computed(() => ts('pages.administrations.addCustomer.formDescriptions.basic')),
      address: computed(() => ts('pages.administrations.addCustomer.formDescriptions.address')),
      authConfig: computed(() => ts('pages.administrations.addCustomer.formDescriptions.authConfig')),
      accessPolicy: computed(() => ts('pages.administrations.addCustomer.formDescriptions.accessPolicy')),
      review: computed(() => ts('pages.administrations.addCustomer.formDescriptions.review')),
    },

    fields: {
      name: computed(() => ts('pages.administrations.addCustomer.fields.name')),
      phone: computed(() => ts('pages.administrations.addCustomer.fields.phone')),
      enabled: computed(() => ts('pages.administrations.addCustomer.fields.enabled')),
      city: computed(() => ts('pages.administrations.addCustomer.fields.city')),
      state: computed(() => ts('pages.administrations.addCustomer.fields.state')),
      country: computed(() => ts('pages.administrations.addCustomer.fields.country')),
      zipCode: computed(() => ts('pages.administrations.addCustomer.fields.zipCode')),
      authProvider: computed(() => ts('pages.administrations.addCustomer.fields.authProvider')),
      issuerUrl: computed(() => ts('pages.administrations.addCustomer.fields.issuerUrl')),
      clientId: computed(() => ts('pages.administrations.addCustomer.fields.clientId')),
      rolePolicy: computed(() => ts('pages.administrations.addCustomer.fields.rolePolicy')),
      defaultScope: computed(() => ts('pages.administrations.addCustomer.fields.defaultScope')),
    },

    hints: {
      phoneFormat: computed(() => ts('pages.administrations.addCustomer.hints.phoneFormat')),
      enabled: computed(() => ts('pages.administrations.addCustomer.hints.enabled')),
      issuerUrl: computed(() => ts('pages.administrations.addCustomer.hints.issuerUrl')),
      clientId: computed(() => ts('pages.administrations.addCustomer.hints.clientId')),
    },

    sections: {
      basicInfo: computed(() => tsTitle('pages.administrations.addCustomer.sections.basicInfo')),
      address: computed(() => tsTitle('pages.administrations.addCustomer.sections.address')),
      authConfig: computed(() => tsTitle('pages.administrations.addCustomer.sections.authConfig')),
      accessPolicy: computed(() => tsTitle('pages.administrations.addCustomer.sections.accessPolicy')),
      progressSteps: computed(() => tsTitle('pages.administrations.addCustomer.sections.progressSteps')),
    },

    buttons: {
      back: computed(() => ts('pages.administrations.addCustomer.buttons.back')),
      next: computed(() => ts('pages.administrations.addCustomer.buttons.next')),
      review: computed(() => ts('pages.administrations.addCustomer.buttons.review')),
      createCustomer: computed(() => ts('pages.administrations.addCustomer.buttons.createCustomer')),
      updateCustomer: computed(() => ts('pages.administrations.addCustomer.buttons.updateCustomer')),
    },

    messages: {
      allFieldsRequired: computed(() => ts('pages.administrations.addCustomer.messages.allFieldsRequired')),
      currentStep: computed(() => ts('pages.administrations.addCustomer.messages.currentStep')),
      completeAllSteps: computed(() => ts('pages.administrations.addCustomer.messages.completeAllSteps')),
      loading: computed(() => ts('pages.administrations.addCustomer.messages.loading')),
      loadFailed: computed(() => ts('pages.administrations.addCustomer.messages.loadFailed')),
      reviewCreateSummary: computed(() => ts('pages.administrations.addCustomer.messages.reviewCreateSummary')),
      reviewEditSummary: computed(() => ts('pages.administrations.addCustomer.messages.reviewEditSummary')),
    },

    notifications: {
      created: computed(() => ts('pages.administrations.addCustomer.notifications.created')),
      updated: computed(() => ts('pages.administrations.addCustomer.notifications.updated')),
      createFailed: computed(() => ts('pages.administrations.addCustomer.notifications.createFailed')),
      updateFailed: computed(() => ts('pages.administrations.addCustomer.notifications.updateFailed')),
      alreadyExists: computed(() => ts('pages.administrations.addCustomer.notifications.alreadyExists')),
      forbidden: computed(() => ts('pages.administrations.addCustomer.notifications.forbidden')),
    },

    errors: {
      apiNotInitialized: computed(() => ts('pages.administrations.addCustomer.errors.apiNotInitialized')),
      idMissing: computed(() => ts('pages.administrations.addCustomer.errors.idMissing')),
    },

    validation: {
      nameRequired: computed(() => ts('pages.administrations.addCustomer.validation.nameRequired')),
      nameMinLength: computed(() => ts('pages.administrations.addCustomer.validation.nameMinLength')),
      nameMaxLength: computed(() => ts('pages.administrations.addCustomer.validation.nameMaxLength')),
      cityMaxLength: computed(() => ts('pages.administrations.addCustomer.validation.cityMaxLength')),
      stateMaxLength: computed(() => ts('pages.administrations.addCustomer.validation.stateMaxLength')),
      countryMaxLength: computed(() => ts('pages.administrations.addCustomer.validation.countryMaxLength')),
      zipCodeMaxLength: computed(() => ts('pages.administrations.addCustomer.validation.zipCodeMaxLength')),
      issuerUrlRequired: computed(() => ts('pages.administrations.addCustomer.validation.issuerUrlRequired')),
      clientIdRequired: computed(() => ts('pages.administrations.addCustomer.validation.clientIdRequired')),
    },

    status: {
      enabled: computed(() => ts('pages.administrations.addCustomer.status.enabled')),
      disabled: computed(() => ts('pages.administrations.addCustomer.status.disabled')),
    },
  };
}
