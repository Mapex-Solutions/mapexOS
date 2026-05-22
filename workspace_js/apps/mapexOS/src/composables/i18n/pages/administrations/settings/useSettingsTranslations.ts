import { computed } from 'vue';
import { useTS } from '@utils/translation';
import type { PageHeaderInfo } from '@components/headers';

/**
 * Settings page translations
 *
 * Structure mirrors:
 * - File: src/pages/administrations/settings/settingsPage/SystemSettingsPage.vue
 * - JSON: src/i18n/{locale}/pages/administrations/settings.json
 * - Composable: src/composables/i18n/pages/administrations/settings/useSettingsTranslations.ts
 *
 * Provides all translations for the System Settings page including:
 * - Page header (title, description)
 * - Tab labels (General, Persistence)
 * - Section titles
 * - Form labels and validation
 * - Success/error messages
 *
 * @example
 * ```ts
 * // In SystemSettingsPage.vue
 * const {
 *   page,
 *   tabs,
 *   general,
 *   persistence
 * } = useSettingsTranslations();
 *
 * <PageHeader :title="page.title.value" :description="page.description.value" />
 * ```
 */
export function useSettingsTranslations() {
  const ts = useTS({ capitalize: true });
  const tsTitle = useTS({ titleCase: true });
  const tsRaw = useTS({ capitalize: false });

  return {
    /**
     * Page header translations
     * Mirrors: pages.administrations.settings
     */
    page: {
      title: computed(() => tsTitle('pages.administrations.settings.title')),
      description: computed(() => ts('pages.administrations.settings.description')),
      info: computed((): PageHeaderInfo => ({
        title: ts('pages.administrations.settings.info.title'),
        description: ts('pages.administrations.settings.info.description'),
        items: [
          {
            icon: 'settings',
            color: 'blue-6',
            title: ts('pages.administrations.settings.info.items.policies.title'),
            text: ts('pages.administrations.settings.info.items.policies.text'),
          },
          {
            icon: 'storage',
            color: 'green-6',
            title: ts('pages.administrations.settings.info.items.retention.title'),
            text: ts('pages.administrations.settings.info.items.retention.text'),
          },
          {
            icon: 'schedule',
            color: 'orange-6',
            title: ts('pages.administrations.settings.info.items.systemWide.title'),
            text: ts('pages.administrations.settings.info.items.systemWide.text'),
          },
          {
            icon: 'speed',
            color: 'purple-6',
            title: ts('pages.administrations.settings.info.items.performance.title'),
            text: ts('pages.administrations.settings.info.items.performance.text'),
          },
          {
            icon: 'backup',
            color: 'indigo-6',
            title: ts('pages.administrations.settings.info.items.backup.title'),
            text: ts('pages.administrations.settings.info.items.backup.text'),
          },
        ],
        docsUrl: 'https://docs.mapexos.com/administration/settings',
        docsLabel: ts('pages.administrations.settings.info.docsLabel'),
      })),
    },

    /**
     * Tab labels
     * Mirrors: pages.administrations.settings.tabs
     */
    tabs: {
      general: computed(() => ts('pages.administrations.settings.tabs.general')),
      persistence: computed(() => ts('pages.administrations.settings.tabs.persistence')),
    },

    /**
     * General tab translations
     * Mirrors: pages.administrations.settings.general
     */
    general: {
      title: computed(() => ts('pages.administrations.settings.general.title')),

      sections: {
        hierarchy: computed(() => ts('pages.administrations.settings.general.sections.hierarchy')),
        basicInfo: computed(() => ts('pages.administrations.settings.general.sections.basicInfo')),
        address: computed(() => ts('pages.administrations.settings.general.sections.address')),
        accessPolicy: computed(() => ts('pages.administrations.settings.general.sections.accessPolicy')),
        system: computed(() => ts('pages.administrations.settings.general.sections.system')),
      },

      fields: {
        parentOrganization: computed(() => ts('pages.administrations.settings.general.fields.parentOrganization')),
        rootOrganization: computed(() => ts('pages.administrations.settings.general.fields.rootOrganization')),
        type: computed(() => ts('pages.administrations.settings.general.fields.type')),
        childOrganizations: computed(() => ts('pages.administrations.settings.general.fields.childOrganizations')),
        code: computed(() => ts('pages.administrations.settings.general.fields.code')),
        path: computed(() => ts('pages.administrations.settings.general.fields.path')),
        depth: computed(() => ts('pages.administrations.settings.general.fields.depth')),
        status: computed(() => ts('pages.administrations.settings.general.fields.status')),
        enabled: computed(() => ts('pages.administrations.settings.general.fields.enabled')),
        disabled: computed(() => ts('pages.administrations.settings.general.fields.disabled')),
        organizationName: computed(() => ts('pages.administrations.settings.general.fields.organizationName')),
        email: computed(() => ts('pages.administrations.settings.general.fields.email')),
        description: computed(() => ts('pages.administrations.settings.general.fields.description')),
        phone: computed(() => ts('pages.administrations.settings.general.fields.phone')),
        street: computed(() => ts('pages.administrations.settings.general.fields.street')),
        address: computed(() => ts('pages.administrations.settings.general.fields.address')),
        city: computed(() => ts('pages.administrations.settings.general.fields.city')),
        state: computed(() => ts('pages.administrations.settings.general.fields.state')),
        country: computed(() => ts('pages.administrations.settings.general.fields.country')),
        zipCode: computed(() => ts('pages.administrations.settings.general.fields.zipCode')),
        authProvider: computed(() => ts('pages.administrations.settings.general.fields.authProvider')),
        rolePolicy: computed(() => ts('pages.administrations.settings.general.fields.rolePolicy')),
        defaultScope: computed(() => ts('pages.administrations.settings.general.fields.defaultScope')),
        created: computed(() => ts('pages.administrations.settings.general.fields.created')),
        lastUpdated: computed(() => ts('pages.administrations.settings.general.fields.lastUpdated')),
      },

      accessPolicyModal: {
        title: computed(() => ts('pages.administrations.settings.general.accessPolicyModal.title')),
        rolePolicyTitle: computed(() => ts('pages.administrations.settings.general.accessPolicyModal.rolePolicyTitle')),
        rolePolicyDescription: computed(() => ts('pages.administrations.settings.general.accessPolicyModal.rolePolicyDescription')),
        rolePolicyMerge: computed(() => ts('pages.administrations.settings.general.accessPolicyModal.rolePolicyMerge')),
        rolePolicyMergeDescription: computed(() => ts('pages.administrations.settings.general.accessPolicyModal.rolePolicyMergeDescription')),
        rolePolicyStrict: computed(() => ts('pages.administrations.settings.general.accessPolicyModal.rolePolicyStrict')),
        rolePolicyStrictDescription: computed(() => ts('pages.administrations.settings.general.accessPolicyModal.rolePolicyStrictDescription')),
        defaultScopeTitle: computed(() => ts('pages.administrations.settings.general.accessPolicyModal.defaultScopeTitle')),
        defaultScopeDescription: computed(() => ts('pages.administrations.settings.general.accessPolicyModal.defaultScopeDescription')),
        defaultScopeRecursive: computed(() => ts('pages.administrations.settings.general.accessPolicyModal.defaultScopeRecursive')),
        defaultScopeRecursiveDescription: computed(() => ts('pages.administrations.settings.general.accessPolicyModal.defaultScopeRecursiveDescription')),
        defaultScopeLocal: computed(() => ts('pages.administrations.settings.general.accessPolicyModal.defaultScopeLocal')),
        defaultScopeLocalDescription: computed(() => ts('pages.administrations.settings.general.accessPolicyModal.defaultScopeLocalDescription')),
        closeButton: computed(() => ts('pages.administrations.settings.general.accessPolicyModal.closeButton')),
      },

      disableModal: {
        title: computed(() => ts('pages.administrations.settings.general.disableModal.title')),
        message: computed(() => ts('pages.administrations.settings.general.disableModal.message')),
        confirmButton: computed(() => ts('pages.administrations.settings.general.disableModal.confirmButton')),
        cancelButton: computed(() => ts('pages.administrations.settings.general.disableModal.cancelButton')),
      },

      validation: {
        nameRequired: computed(() => ts('pages.administrations.settings.general.validation.nameRequired')),
        emailRequired: computed(() => ts('pages.administrations.settings.general.validation.emailRequired')),
      },

      messages: {
        savedSuccessfully: computed(() => ts('pages.administrations.settings.general.messages.savedSuccessfully')),
        loading: computed(() => ts('pages.administrations.settings.general.messages.loading')),
        infoTooltip: computed(() => ts('pages.administrations.settings.general.messages.infoTooltip')),
      },

      buttons: {
        saveChanges: computed(() => ts('pages.administrations.settings.general.buttons.saveChanges')),
      },
    },

    /**
     * Persistence tab translations
     * Mirrors: pages.administrations.settings.persistence
     */
    persistence: {
      title: computed(() => ts('pages.administrations.settings.persistence.title')),
      description: computed(() => ts('pages.administrations.settings.persistence.description')),
      resetButton: computed(() => ts('pages.administrations.settings.persistence.resetButton')),

      eventTypes: {
        eventsRaw: computed(() => ts('pages.administrations.settings.persistence.eventTypes.eventsRaw')),
        eventsJsExecutor: computed(() => ts('pages.administrations.settings.persistence.eventTypes.eventsJsExecutor')),
        eventsRouter: computed(() => ts('pages.administrations.settings.persistence.eventTypes.eventsRouter')),
        events: computed(() => ts('pages.administrations.settings.persistence.eventTypes.events')),
        eventsWorkflow: computed(() => ts('pages.administrations.settings.persistence.eventTypes.eventsWorkflow')),
        eventsAudit: computed(() => ts('pages.administrations.settings.persistence.eventTypes.eventsAudit')),
        eventsNotifications: computed(() => ts('pages.administrations.settings.persistence.eventTypes.eventsNotifications')),
      },

      eventDescriptions: {
        eventsRaw: {
          title: computed(() => ts('pages.administrations.settings.persistence.eventDescriptions.eventsRaw.title')),
          description: computed(() => ts('pages.administrations.settings.persistence.eventDescriptions.eventsRaw.description')),
        },
        eventsJsExecutor: {
          title: computed(() => ts('pages.administrations.settings.persistence.eventDescriptions.eventsJsExecutor.title')),
          description: computed(() => ts('pages.administrations.settings.persistence.eventDescriptions.eventsJsExecutor.description')),
        },
        eventsRouter: {
          title: computed(() => ts('pages.administrations.settings.persistence.eventDescriptions.eventsRouter.title')),
          description: computed(() => ts('pages.administrations.settings.persistence.eventDescriptions.eventsRouter.description')),
        },
        events: {
          title: computed(() => ts('pages.administrations.settings.persistence.eventDescriptions.events.title')),
          description: computed(() => ts('pages.administrations.settings.persistence.eventDescriptions.events.description')),
        },
        eventsWorkflow: {
          title: computed(() => ts('pages.administrations.settings.persistence.eventDescriptions.eventsWorkflow.title')),
          description: computed(() => ts('pages.administrations.settings.persistence.eventDescriptions.eventsWorkflow.description')),
        },
        eventsAudit: {
          title: computed(() => ts('pages.administrations.settings.persistence.eventDescriptions.eventsAudit.title')),
          description: computed(() => ts('pages.administrations.settings.persistence.eventDescriptions.eventsAudit.description')),
        },
        eventsNotifications: {
          title: computed(() => ts('pages.administrations.settings.persistence.eventDescriptions.eventsNotifications.title')),
          description: computed(() => ts('pages.administrations.settings.persistence.eventDescriptions.eventsNotifications.description')),
        },
      },

      groups: {
        pipelineEvents: computed(() => ts('pages.administrations.settings.persistence.groups.pipelineEvents')),
        pipelineDescription: computed(() => ts('pages.administrations.settings.persistence.groups.pipelineDescription')),
        storageCompliance: computed(() => ts('pages.administrations.settings.persistence.groups.storageCompliance')),
        storageDescription: computed(() => ts('pages.administrations.settings.persistence.groups.storageDescription')),
      },

      columns: {
        eventType: computed(() => ts('pages.administrations.settings.persistence.columns.eventType')),
        retentionDays: computed(() => ts('pages.administrations.settings.persistence.columns.retentionDays')),
        limits: computed(() => ts('pages.administrations.settings.persistence.columns.limits')),
      },

      form: {
        daysLabel: computed(() => ts('pages.administrations.settings.persistence.form.daysLabel')),
        daysHint: computed(() => tsRaw('pages.administrations.settings.persistence.form.daysHint')),
        defaultLabel: computed(() => tsRaw('pages.administrations.settings.persistence.form.defaultLabel')),
      },

      validation: {
        required: computed(() => ts('pages.administrations.settings.persistence.validation.required')),
        minValue: computed(() => ts('pages.administrations.settings.persistence.validation.minValue')),
        maxValue: computed(() => ts('pages.administrations.settings.persistence.validation.maxValue')),
        integer: computed(() => ts('pages.administrations.settings.persistence.validation.integer')),
      },

      messages: {
        savedSuccessfully: computed(() => ts('pages.administrations.settings.persistence.messages.savedSuccessfully')),
        resetSuccessfully: computed(() => ts('pages.administrations.settings.persistence.messages.resetSuccessfully')),
        confirmReset: computed(() => ts('pages.administrations.settings.persistence.messages.confirmReset')),
        complianceWarning: computed(() => ts('pages.administrations.settings.persistence.messages.complianceWarning')),
        lakeHouseWarning: computed(() => ts('pages.administrations.settings.persistence.messages.lakeHouseWarning')),
      },

      buttons: {
        close: computed(() => ts('pages.administrations.settings.persistence.buttons.close')),
      },

      tooltips: {
        moreInfo: computed(() => ts('pages.administrations.settings.persistence.tooltips.moreInfo')),
      },
    },

    /**
     * Error translations
     * Mirrors: pages.administrations.settings.errors
     */
    errors: {
      apiNotInitialized: computed(() => ts('pages.administrations.settings.errors.apiNotInitialized')),
      noOrganization: computed(() => ts('pages.administrations.settings.errors.noOrganization')),
    },

    /**
     * Page-level notification translations
     * Mirrors: pages.administrations.settings.notifications
     */
    notifications: {
      savedSuccessfully: computed(() => ts('pages.administrations.settings.notifications.savedSuccessfully')),
    },

    /**
     * Overview/section translations
     * Mirrors: pages.administrations.settings.overview
     */
    overview: {
      organizationOverview: computed(() => ts('pages.administrations.settings.overview.organizationOverview')),
    },

  };
}
