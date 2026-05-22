import { computed } from 'vue';
import { useTS } from '@utils/translation';

/**
 * User Profile page translations
 *
 * Structure mirrors:
 * - File: src/pages/administrations/users/userProfilePage/UserProfilePage.vue
 * - JSON: src/i18n/{locale}/pages/administrations/userProfile.json
 * - Composable: src/composables/i18n/pages/administrations/users/useUserProfileTranslations.ts
 */
export function useUserProfileTranslations() {
  const ts = useTS({ capitalize: true });
  const tsTitle = useTS({ titleCase: true });

  return {
    page: {
      title: computed(() => tsTitle('pages.administrations.userProfile.title')),
      description: computed(() => ts('pages.administrations.userProfile.description')),
      backButton: computed(() => ts('pages.administrations.userProfile.backButton')),
    },

    navigation: {
      personal: {
        label: computed(() => tsTitle('pages.administrations.userProfile.navigation.personal.label')),
        description: computed(() => ts('pages.administrations.userProfile.navigation.personal.description')),
      },
      password: {
        label: computed(() => tsTitle('pages.administrations.userProfile.navigation.password.label')),
        description: computed(() => ts('pages.administrations.userProfile.navigation.password.description')),
      },
      groupsAccess: {
        label: computed(() => tsTitle('pages.administrations.userProfile.navigation.groupsAccess.label')),
        description: computed(() => ts('pages.administrations.userProfile.navigation.groupsAccess.description')),
      },
      review: {
        label: computed(() => tsTitle('pages.administrations.userProfile.navigation.review.label')),
        description: computed(() => ts('pages.administrations.userProfile.navigation.review.description')),
      },
    },

    sections: {
      personalInfo: computed(() => tsTitle('pages.administrations.userProfile.sections.personalInfo')),
      password: computed(() => tsTitle('pages.administrations.userProfile.sections.password')),
      groupsAccess: computed(() => tsTitle('pages.administrations.userProfile.sections.groupsAccess')),
      profileSettings: computed(() => tsTitle('pages.administrations.userProfile.sections.profileSettings')),
      review: computed(() => tsTitle('pages.administrations.userProfile.sections.review')),
    },

    sectionDescriptions: {
      personalInfo: computed(() => ts('pages.administrations.userProfile.sectionDescriptions.personalInfo')),
      password: computed(() => ts('pages.administrations.userProfile.sectionDescriptions.password')),
      groupsAccess: computed(() => ts('pages.administrations.userProfile.sectionDescriptions.groupsAccess')),
      profileSettings: computed(() => ts('pages.administrations.userProfile.sectionDescriptions.profileSettings')),
      review: computed(() => ts('pages.administrations.userProfile.sectionDescriptions.review')),
    },

    stepper: {
      infoText: computed(() => ts('pages.administrations.userProfile.stepper.infoText')),
      currentSection: computed(() => ts('pages.administrations.userProfile.stepper.currentSection')),
    },

    fields: {
      firstName: computed(() => ts('pages.administrations.userProfile.fields.firstName')),
      lastName: computed(() => ts('pages.administrations.userProfile.fields.lastName')),
      email: computed(() => ts('pages.administrations.userProfile.fields.email')),
      phone: computed(() => ts('pages.administrations.userProfile.fields.phone')),
      jobTitle: computed(() => ts('pages.administrations.userProfile.fields.jobTitle')),
      currentPassword: computed(() => ts('pages.administrations.userProfile.fields.currentPassword')),
      newPassword: computed(() => ts('pages.administrations.userProfile.fields.newPassword')),
      confirmPassword: computed(() => ts('pages.administrations.userProfile.fields.confirmPassword')),
      required: computed(() => ts('pages.administrations.userProfile.fields.required')),
    },

    buttons: {
      saveChanges: computed(() => ts('pages.administrations.userProfile.buttons.saveChanges')),
      updatePassword: computed(() => ts('pages.administrations.userProfile.buttons.updatePassword')),
      cancel: computed(() => ts('pages.administrations.userProfile.buttons.cancel')),
      save: computed(() => ts('pages.administrations.userProfile.buttons.save')),
    },

    groupsAccess: {
      groupsTitle: computed(() => ts('pages.administrations.userProfile.groupsAccess.groupsTitle')),
      groupsCount: computed(() => ts('pages.administrations.userProfile.groupsAccess.groupsCount')),
      membershipsTitle: computed(() => ts('pages.administrations.userProfile.groupsAccess.membershipsTitle')),
      membershipsCount: computed(() => ts('pages.administrations.userProfile.groupsAccess.membershipsCount')),
      noGroups: computed(() => ts('pages.administrations.userProfile.groupsAccess.noGroups')),
      noMemberships: computed(() => ts('pages.administrations.userProfile.groupsAccess.noMemberships')),
      noDescription: computed(() => ts('pages.administrations.userProfile.groupsAccess.noDescription')),
      columns: {
        organization: computed(() => ts('pages.administrations.userProfile.groupsAccess.columns.organization')),
        roles: computed(() => ts('pages.administrations.userProfile.groupsAccess.columns.roles')),
        scope: computed(() => ts('pages.administrations.userProfile.groupsAccess.columns.scope')),
        via: computed(() => ts('pages.administrations.userProfile.groupsAccess.columns.via')),
      },
    },

    validation: {
      firstNameRequired: computed(() => ts('pages.administrations.userProfile.validation.firstNameRequired')),
      lastNameRequired: computed(() => ts('pages.administrations.userProfile.validation.lastNameRequired')),
      emailRequired: computed(() => ts('pages.administrations.userProfile.validation.emailRequired')),
      emailInvalid: computed(() => ts('pages.administrations.userProfile.validation.emailInvalid')),
      currentPasswordRequired: computed(() => ts('pages.administrations.userProfile.validation.currentPasswordRequired')),
      newPasswordRequired: computed(() => ts('pages.administrations.userProfile.validation.newPasswordRequired')),
      passwordMinLength: computed(() => ts('pages.administrations.userProfile.validation.passwordMinLength')),
      confirmPasswordRequired: computed(() => ts('pages.administrations.userProfile.validation.confirmPasswordRequired')),
      passwordsDoNotMatch: computed(() => ts('pages.administrations.userProfile.validation.passwordsDoNotMatch')),
    },

    messages: {
      personalInfoUpdated: computed(() => ts('pages.administrations.userProfile.messages.personalInfoUpdated')),
      passwordUpdated: computed(() => ts('pages.administrations.userProfile.messages.passwordUpdated')),
      profileSaved: computed(() => ts('pages.administrations.userProfile.messages.profileSaved')),
      errorLoading: computed(() => ts('pages.administrations.userProfile.messages.errorLoading')),
    },

    review: {
      title: computed(() => ts('pages.administrations.userProfile.review.title')),
      subtitle: computed(() => ts('pages.administrations.userProfile.review.subtitle')),
      sections: {
        personal: computed(() => ts('pages.administrations.userProfile.review.sections.personal')),
        password: computed(() => ts('pages.administrations.userProfile.review.sections.password')),
        groupsAccess: computed(() => ts('pages.administrations.userProfile.review.sections.groupsAccess')),
      },
      fields: {
        firstName: computed(() => ts('pages.administrations.userProfile.review.fields.firstName')),
        lastName: computed(() => ts('pages.administrations.userProfile.review.fields.lastName')),
        email: computed(() => ts('pages.administrations.userProfile.review.fields.email')),
        phone: computed(() => ts('pages.administrations.userProfile.review.fields.phone')),
        jobTitle: computed(() => ts('pages.administrations.userProfile.review.fields.jobTitle')),
        passwordChanged: computed(() => ts('pages.administrations.userProfile.review.fields.passwordChanged')),
        yes: computed(() => ts('pages.administrations.userProfile.review.fields.yes')),
        no: computed(() => ts('pages.administrations.userProfile.review.fields.no')),
        groupsCount: computed(() => ts('pages.administrations.userProfile.review.fields.groupsCount')),
        membershipsCount: computed(() => ts('pages.administrations.userProfile.review.fields.membershipsCount')),
      },
    },
  };
}
