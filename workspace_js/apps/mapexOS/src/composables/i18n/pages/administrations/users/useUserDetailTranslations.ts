import { computed } from 'vue';
import { useTS } from '@utils/translation';

/**
 * User detail page translations
 *
 * Provides all translations for the User Detail page including:
 * - Page header (title, description, buttons)
 * - Tab labels
 * - Profile tab (fields, sections, status)
 * - Access tab (roles, permissions)
 * - Groups tab (group list)
 */
export function useUserDetailTranslations() {
  const ts = useTS({ capitalize: true });
  const tsTitle = useTS({ titleCase: true });

  return {
    /**
     * Page header translations
     */
    page: {
      title: computed(() => tsTitle('pages.administrations.users.detail.title')),
      description: computed(() => ts('pages.administrations.users.detail.description')),
      backButton: computed(() => ts('pages.administrations.users.detail.backButton')),
      editButton: computed(() => ts('pages.administrations.users.detail.editButton')),
      loading: computed(() => ts('pages.administrations.users.detail.loading')),
      loadError: computed(() => ts('pages.administrations.users.detail.loadError')),
    },

    /**
     * Tab translations
     */
    tabs: {
      profile: computed(() => ts('pages.administrations.users.detail.tabs.profile')),
      access: computed(() => ts('pages.administrations.users.detail.tabs.access')),
      groups: computed(() => ts('pages.administrations.users.detail.tabs.groups')),
    },

    /**
     * Profile tab translations
     */
    profile: {
      loading: computed(() => ts('pages.administrations.users.detail.profile.loading')),
      error: computed(() => ts('pages.administrations.users.detail.profile.error')),
      sections: {
        basicInfo: computed(() => ts('pages.administrations.users.detail.profile.sections.basicInfo')),
        contact: computed(() => ts('pages.administrations.users.detail.profile.sections.contact')),
        authentication: computed(() => ts('pages.administrations.users.detail.profile.sections.authentication')),
        timestamps: computed(() => ts('pages.administrations.users.detail.profile.sections.timestamps')),
      },
      fields: {
        name: computed(() => ts('pages.administrations.users.detail.profile.fields.name')),
        email: computed(() => ts('pages.administrations.users.detail.profile.fields.email')),
        phone: computed(() => ts('pages.administrations.users.detail.profile.fields.phone')),
        jobTitle: computed(() => ts('pages.administrations.users.detail.profile.fields.jobTitle')),
        organization: computed(() => ts('pages.administrations.users.detail.profile.fields.organization')),
        status: computed(() => ts('pages.administrations.users.detail.profile.fields.status')),
        authProvider: computed(() => ts('pages.administrations.users.detail.profile.fields.authProvider')),
        changePasswordNextLogin: computed(() => ts('pages.administrations.users.detail.profile.fields.changePasswordNextLogin')),
        externalId: computed(() => ts('pages.administrations.users.detail.profile.fields.externalId')),
        created: computed(() => ts('pages.administrations.users.detail.profile.fields.created')),
        updated: computed(() => ts('pages.administrations.users.detail.profile.fields.updated')),
      },
      status: {
        active: computed(() => ts('pages.administrations.users.detail.profile.status.active')),
        inactive: computed(() => ts('pages.administrations.users.detail.profile.status.inactive')),
      },
      values: {
        yes: computed(() => ts('pages.administrations.users.detail.profile.values.yes')),
        no: computed(() => ts('pages.administrations.users.detail.profile.values.no')),
      },
      authProviders: {
        internal: computed(() => ts('pages.administrations.users.detail.profile.authProviders.internal')),
        google: computed(() => ts('pages.administrations.users.detail.profile.authProviders.google')),
        github: computed(() => ts('pages.administrations.users.detail.profile.authProviders.github')),
        microsoft: computed(() => ts('pages.administrations.users.detail.profile.authProviders.microsoft')),
        keycloak: computed(() => ts('pages.administrations.users.detail.profile.authProviders.keycloak')),
      },
    },

    /**
     * Access tab translations
     */
    access: {
      title: computed(() => ts('pages.administrations.users.detail.access.title')),
      loading: computed(() => ts('pages.administrations.users.detail.access.loading')),
      error: computed(() => ts('pages.administrations.users.detail.access.error')),
      count: computed(() => ts('pages.administrations.users.detail.access.count')),
      empty: {
        title: computed(() => ts('pages.administrations.users.detail.access.empty.title')),
        description: computed(() => ts('pages.administrations.users.detail.access.empty.description')),
      },
      apiNotAvailable: {
        title: computed(() => ts('pages.administrations.users.detail.access.apiNotAvailable.title')),
        description: computed(() => ts('pages.administrations.users.detail.access.apiNotAvailable.description')),
      },
      columns: {
        organization: computed(() => ts('pages.administrations.users.detail.access.columns.organization')),
        roles: computed(() => ts('pages.administrations.users.detail.access.columns.roles')),
        scope: computed(() => ts('pages.administrations.users.detail.access.columns.scope')),
        via: computed(() => ts('pages.administrations.users.detail.access.columns.via')),
      },
      fields: {
        scope: computed(() => ts('pages.administrations.users.detail.access.fields.scope')),
        permissions: computed(() => ts('pages.administrations.users.detail.access.fields.permissions')),
        organization: computed(() => ts('pages.administrations.users.detail.access.fields.organization')),
        roles: computed(() => ts('pages.administrations.users.detail.access.fields.roles')),
        via: computed(() => ts('pages.administrations.users.detail.access.fields.via')),
      },
      type: {
        system: computed(() => ts('pages.administrations.users.detail.access.type.system')),
        custom: computed(() => ts('pages.administrations.users.detail.access.type.custom')),
      },
    },

    /**
     * Groups tab translations
     */
    groups: {
      loading: computed(() => ts('pages.administrations.users.detail.groups.loading')),
      error: computed(() => ts('pages.administrations.users.detail.groups.error')),
      count: computed(() => ts('pages.administrations.users.detail.groups.count')),
      noDescription: computed(() => ts('pages.administrations.users.detail.groups.noDescription')),
      empty: {
        title: computed(() => ts('pages.administrations.users.detail.groups.empty.title')),
        description: computed(() => ts('pages.administrations.users.detail.groups.empty.description')),
      },
      type: {
        system: computed(() => ts('pages.administrations.users.detail.groups.type.system')),
        custom: computed(() => ts('pages.administrations.users.detail.groups.type.custom')),
      },
      status: {
        enabled: computed(() => ts('pages.administrations.users.detail.groups.status.enabled')),
        disabled: computed(() => ts('pages.administrations.users.detail.groups.status.disabled')),
      },
      template: {
        yes: computed(() => ts('pages.administrations.users.detail.groups.template.yes')),
        no: computed(() => ts('pages.administrations.users.detail.groups.template.no')),
      },
    },

    /**
     * Error translations
     * Mirrors: pages.administrations.users.errors
     */
    errors: {
      apiNotInitialized: computed(() => ts('pages.administrations.users.errors.apiNotInitialized')),
      idMissing: computed(() => ts('pages.administrations.users.errors.idMissing')),
    },
  };
}
