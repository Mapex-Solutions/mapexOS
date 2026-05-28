import { computed } from 'vue';
import { useTS } from '@utils/translation';

/**
 * Add/Edit User page translations
 *
 * Structure mirrors:
 * - File: src/pages/administrations/users/createEditUserPage/CreateEditUserPage.vue
 * - JSON: src/i18n/{locale}/pages/administrations/addUser.json
 * - Composable: src/composables/i18n/pages/administrations/users/useAddUserTranslations.ts
 */
export function useAddUserTranslations() {
  const ts = useTS({ capitalize: true });
  const tsTitle = useTS({ titleCase: true });

  return {
    page: {
      title: computed(() => tsTitle('pages.administrations.addUser.title')),
      titleEdit: computed(() => tsTitle('pages.administrations.addUser.titleEdit')),
      description: computed(() => ts('pages.administrations.addUser.description')),
      descriptionEdit: computed(() => ts('pages.administrations.addUser.descriptionEdit')),
      backButton: computed(() => ts('pages.administrations.addUser.backButton')),
    },

    steps: {
      personal: computed(() => tsTitle('pages.administrations.addUser.steps.personal')),
      security: computed(() => tsTitle('pages.administrations.addUser.steps.security')),
      access: computed(() => tsTitle('pages.administrations.addUser.steps.access')),
      review: computed(() => tsTitle('pages.administrations.addUser.steps.review')),
    },

    stepDescriptions: {
      personal: computed(() => ts('pages.administrations.addUser.stepDescriptions.personal')),
      security: computed(() => ts('pages.administrations.addUser.stepDescriptions.security')),
      access: computed(() => ts('pages.administrations.addUser.stepDescriptions.access')),
      review: computed(() => ts('pages.administrations.addUser.stepDescriptions.review')),
    },

    formDescriptions: {
      personal: computed(() => ts('pages.administrations.addUser.formDescriptions.personal')),
      security: computed(() => ts('pages.administrations.addUser.formDescriptions.security')),
      access: computed(() => ts('pages.administrations.addUser.formDescriptions.access')),
      review: computed(() => ts('pages.administrations.addUser.formDescriptions.review')),
    },

    fields: {
      firstName: computed(() => ts('pages.administrations.addUser.fields.firstName')),
      lastName: computed(() => ts('pages.administrations.addUser.fields.lastName')),
      email: computed(() => ts('pages.administrations.addUser.fields.email')),
      phone: computed(() => ts('pages.administrations.addUser.fields.phone')),
      jobTitle: computed(() => ts('pages.administrations.addUser.fields.jobTitle')),
      enabled: computed(() => ts('pages.administrations.addUser.fields.enabled')),
      authProvider: computed(() => ts('pages.administrations.addUser.fields.authProvider')),
      password: computed(() => ts('pages.administrations.addUser.fields.password')),
      confirmPassword: computed(() => ts('pages.administrations.addUser.fields.confirmPassword')),
      changePasswordNextLogin: computed(() => ts('pages.administrations.addUser.fields.changePasswordNextLogin')),
      externalId: computed(() => ts('pages.administrations.addUser.fields.externalId')),
      organization: computed(() => ts('pages.administrations.addUser.fields.organization')),
      roles: computed(() => ts('pages.administrations.addUser.fields.roles')),
      group: computed(() => ts('pages.administrations.addUser.fields.group')),
      groups: computed(() => ts('pages.administrations.addUser.fields.groups')),
      scope: computed(() => ts('pages.administrations.addUser.fields.scope')),
      accessType: computed(() => ts('pages.administrations.addUser.fields.accessType')),
      selectGroup: computed(() => ts('pages.administrations.addUser.fields.selectGroup')),
      directMemberships: computed(() => ts('pages.administrations.addUser.fields.directMemberships')),
      newGroupName: computed(() => ts('pages.administrations.addUser.fields.newGroupName')),
      newGroupDescription: computed(() => ts('pages.administrations.addUser.fields.newGroupDescription')),
      newGroupRoles: computed(() => ts('pages.administrations.addUser.fields.newGroupRoles')),
    },

    labels: {
      optional: computed(() => ts('pages.administrations.addUser.labels.optional')),
      recommended: computed(() => ts('pages.administrations.addUser.labels.recommended')),
      newGroup: computed(() => ts('pages.administrations.addUser.labels.newGroup')),
      local: computed(() => ts('pages.administrations.addUser.labels.local')),
      recursive: computed(() => ts('pages.administrations.addUser.labels.recursive')),
    },

    placeholders: {
      selectOrganization: computed(() => ts('pages.administrations.addUser.placeholders.selectOrganization')),
      selectOrganizationFirst: computed(() => ts('pages.administrations.addUser.placeholders.selectOrganizationFirst')),
      selectRoles: computed(() => ts('pages.administrations.addUser.placeholders.selectRoles')),
      selectGroup: computed(() => ts('pages.administrations.addUser.placeholders.selectGroup')),
      noGroupsSelected: computed(() => ts('pages.administrations.addUser.placeholders.noGroupsSelected')),
      noMembershipsSelected: computed(() => ts('pages.administrations.addUser.placeholders.noMembershipsSelected')),
      newGroupName: computed(() => ts('pages.administrations.addUser.placeholders.newGroupName')),
      newGroupDescription: computed(() => ts('pages.administrations.addUser.placeholders.newGroupDescription')),
    },

    hints: {
      phoneFormat: computed(() => ts('pages.administrations.addUser.hints.phoneFormat')),
      enabled: computed(() => ts('pages.administrations.addUser.hints.enabled')),
      passwordRequirements: computed(() => ts('pages.administrations.addUser.hints.passwordRequirements')),
      changePasswordNextLogin: computed(() => ts('pages.administrations.addUser.hints.changePasswordNextLogin')),
      externalId: computed(() => ts('pages.administrations.addUser.hints.externalId')),
      organization: computed(() => ts('pages.administrations.addUser.hints.organization')),
      roles: computed(() => ts('pages.administrations.addUser.hints.roles')),
      group: computed(() => ts('pages.administrations.addUser.hints.group')),
      scopeLocal: computed(() => ts('pages.administrations.addUser.hints.scopeLocal')),
      scopeRecursive: computed(() => ts('pages.administrations.addUser.hints.scopeRecursive')),
      groupSelection: computed(() => ts('pages.administrations.addUser.hints.groupSelection')),
      groupInheritance: computed(() => ts('pages.administrations.addUser.hints.groupInheritance')),
      directWarning: computed(() => ts('pages.administrations.addUser.hints.directWarning')),
      orgFromContext: computed(() => ts('pages.administrations.addUser.hints.orgFromContext')),
      multipleGroups: computed(() => ts('pages.administrations.addUser.hints.multipleGroups')),
      multipleMemberships: computed(() => ts('pages.administrations.addUser.hints.multipleMemberships')),
    },

    sections: {
      personalInfo: computed(() => tsTitle('pages.administrations.addUser.sections.personalInfo')),
      security: computed(() => tsTitle('pages.administrations.addUser.sections.security')),
      access: computed(() => tsTitle('pages.administrations.addUser.sections.access')),
      review: computed(() => tsTitle('pages.administrations.addUser.sections.review')),
      progressSteps: computed(() => tsTitle('pages.administrations.addUser.sections.progressSteps')),
    },

    actions: {
      remove: computed(() => ts('pages.administrations.addUser.actions.remove')),
      addGroup: computed(() => ts('pages.administrations.addUser.actions.addGroup')),
      addMembership: computed(() => ts('pages.administrations.addUser.actions.addMembership')),
      cancel: computed(() => ts('pages.administrations.addUser.actions.cancel')),
      add: computed(() => ts('pages.administrations.addUser.actions.add')),
      clearAll: computed(() => ts('pages.administrations.addUser.actions.clearAll')),
    },

    dialogs: {
      addMembership: {
        title: computed(() => tsTitle('pages.administrations.addUser.dialogs.addMembership.title')),
      },
    },

    buttons: {
      back: computed(() => ts('pages.administrations.addUser.buttons.back')),
      next: computed(() => ts('pages.administrations.addUser.buttons.next')),
      edit: computed(() => ts('pages.administrations.addUser.buttons.edit')),
      createUser: computed(() => ts('pages.administrations.addUser.buttons.createUser')),
      updateUser: computed(() => ts('pages.administrations.addUser.buttons.updateUser')),
    },

    messages: {
      allFieldsRequired: computed(() => ts('pages.administrations.addUser.messages.allFieldsRequired')),
      currentStep: computed(() => ts('pages.administrations.addUser.messages.currentStep')),
      created: computed(() => ts('pages.administrations.addUser.messages.created')),
      updated: computed(() => ts('pages.administrations.addUser.messages.updated')),
      createFailed: computed(() => ts('pages.administrations.addUser.messages.createFailed')),
      updateFailed: computed(() => ts('pages.administrations.addUser.messages.updateFailed')),
      loadFailed: computed(() => ts('pages.administrations.addUser.messages.loadFailed')),
      emailExists: computed(() => ts('pages.administrations.addUser.messages.emailExists')),
      forbidden: computed(() => ts('pages.administrations.addUser.messages.forbidden')),
      completeAllSteps: computed(() => ts('pages.administrations.addUser.messages.completeAllSteps')),
      loading: computed(() => ts('pages.administrations.addUser.messages.loading')),
      ssoInfo: computed(() => ts('pages.administrations.addUser.messages.ssoInfo')),
      notApplicable: computed(() => ts('pages.administrations.addUser.messages.notApplicable')),
      reviewCreateSummary: computed(() => ts('pages.administrations.addUser.messages.reviewCreateSummary')),
      reviewEditSummary: computed(() => ts('pages.administrations.addUser.messages.reviewEditSummary')),
      noAccessConfigured: computed(() => ts('pages.administrations.addUser.messages.noAccessConfigured')),
      noOrgContext: computed(() => ts('pages.administrations.addUser.messages.noOrgContext')),
    },

    status: {
      enabled: computed(() => ts('pages.administrations.addUser.status.enabled')),
      disabled: computed(() => ts('pages.administrations.addUser.status.disabled')),
      yes: computed(() => ts('pages.administrations.addUser.status.yes')),
      no: computed(() => ts('pages.administrations.addUser.status.no')),
    },

    /**
     * Tour translations
     * Mirrors: pages.administrations.addUser.tour
     */
    tour: {
      overview: {
        title: computed(() => ts('pages.administrations.addUser.tour.overview.title')),
        description: computed(() => ts('pages.administrations.addUser.tour.overview.description')),
      },
      step1: {
        title: computed(() => ts('pages.administrations.addUser.tour.step1.title')),
        description: computed(() => ts('pages.administrations.addUser.tour.step1.description')),
      },
      step2: {
        title: computed(() => ts('pages.administrations.addUser.tour.step2.title')),
        description: computed(() => ts('pages.administrations.addUser.tour.step2.description')),
      },
      step3: {
        title: computed(() => ts('pages.administrations.addUser.tour.step3.title')),
        description: computed(() => ts('pages.administrations.addUser.tour.step3.description')),
      },
      step4: {
        title: computed(() => ts('pages.administrations.addUser.tour.step4.title')),
        description: computed(() => ts('pages.administrations.addUser.tour.step4.description')),
      },
    },

    validation: {
      firstNameRequired: computed(() => ts('pages.administrations.addUser.validation.firstNameRequired')),
      firstNameMinLength: computed(() => ts('pages.administrations.addUser.validation.firstNameMinLength')),
      firstNameMaxLength: computed(() => ts('pages.administrations.addUser.validation.firstNameMaxLength')),
      lastNameRequired: computed(() => ts('pages.administrations.addUser.validation.lastNameRequired')),
      lastNameMinLength: computed(() => ts('pages.administrations.addUser.validation.lastNameMinLength')),
      lastNameMaxLength: computed(() => ts('pages.administrations.addUser.validation.lastNameMaxLength')),
      emailRequired: computed(() => ts('pages.administrations.addUser.validation.emailRequired')),
      emailInvalid: computed(() => ts('pages.administrations.addUser.validation.emailInvalid')),
      emailMaxLength: computed(() => ts('pages.administrations.addUser.validation.emailMaxLength')),
      jobTitleMaxLength: computed(() => ts('pages.administrations.addUser.validation.jobTitleMaxLength')),
      passwordRequired: computed(() => ts('pages.administrations.addUser.validation.passwordRequired')),
      passwordMinLength: computed(() => ts('pages.administrations.addUser.validation.passwordMinLength')),
      passwordMaxLength: computed(() => ts('pages.administrations.addUser.validation.passwordMaxLength')),
      passwordMismatch: computed(() => ts('pages.administrations.addUser.validation.passwordMismatch')),
      organizationRequired: computed(() => ts('pages.administrations.addUser.validation.organizationRequired')),
      rolesRequired: computed(() => ts('pages.administrations.addUser.validation.rolesRequired')),
      newGroupNameRequired: computed(() => ts('pages.administrations.addUser.validation.newGroupNameRequired')),
      newGroupNameMinLength: computed(() => ts('pages.administrations.addUser.validation.newGroupNameMinLength')),
    },

    errors: {
      apiNotInitialized: computed(() => ts('pages.administrations.addUser.errors.apiNotInitialized')),
      onboardingApiNotInitialized: computed(() => ts('pages.administrations.addUser.errors.onboardingApiNotInitialized')),
      idMissing: computed(() => ts('pages.administrations.addUser.errors.idMissing')),
    },

    reviewLabels: {
      groupMode: computed(() => ts('pages.administrations.addUser.reviewLabels.groupMode')),
      newGroupName: computed(() => ts('pages.administrations.addUser.reviewLabels.newGroupName')),
      description: computed(() => ts('pages.administrations.addUser.reviewLabels.description')),
    },
  };
}
