import { computed } from 'vue';
import { useTS } from '@utils/translation';

/**
 * Custom composable for AddAsset page translations.
 * Provides type-safe reactive access to all translated strings with formatting utilities.
 *
 * Structure mirrors:
 * - File: src/pages/assets/assets/assetAddPage/AddAsset.vue
 * - JSON: src/i18n/{locale}/pages/assets/addAsset.json
 * - Composable: src/composables/i18n/pages/assets/addAsset/useAddAssetTranslations.ts
 */
export function useAddAssetTranslations() {
  const ts = useTS({ capitalize: true });
  const tsTitle = useTS({ titleCase: true });
  const tsRaw = useTS({ capitalize: false });

  return {
    /**
     * Page header translations
     */
    page: {
      title: computed(() => ts('pages.assets.addAsset.page.title')),
      titleEdit: computed(() => ts('pages.assets.addAsset.page.titleEdit')),
      description: computed(() => tsRaw('pages.assets.addAsset.page.description')),
      button: computed(() => ts('pages.assets.addAsset.page.button')),
    },

    /**
     * Stepper header translations
     */
    stepper: {
      title: computed(() => tsTitle('pages.assets.addAsset.stepper.title')),
      subtitle: computed(() => tsRaw('pages.assets.addAsset.stepper.subtitle')),
      requiredInfo: computed(() => tsRaw('pages.assets.addAsset.stepper.requiredInfo')),
      currentStep: computed(() => ts('pages.assets.addAsset.stepper.currentStep')),
    },

    /**
     * Step-specific translations
     */
    steps: {
      /**
       * Step 1: Identification
       */
      step1: {
        label: computed(() => ts('pages.assets.addAsset.steps.step1.label')),
        description: computed(() => tsRaw('pages.assets.addAsset.steps.step1.description')),
        title: computed(() => ts('pages.assets.addAsset.steps.step1.title')),
        subtitle: computed(() => tsRaw('pages.assets.addAsset.steps.step1.subtitle')),
        fields: {
          name: {
            label: computed(() => ts('pages.assets.addAsset.steps.step1.fields.name.label')),
            placeholder: computed(() => tsRaw('pages.assets.addAsset.steps.step1.fields.name.placeholder')),
            hint: computed(() => tsRaw('pages.assets.addAsset.steps.step1.fields.name.hint')),
            required: computed(() => tsRaw('pages.assets.addAsset.steps.step1.fields.name.required')),
          },
          assetId: {
            label: computed(() => ts('pages.assets.addAsset.steps.step1.fields.assetId.label')),
            placeholder: computed(() => tsRaw('pages.assets.addAsset.steps.step1.fields.assetId.placeholder')),
            hint: computed(() => tsRaw('pages.assets.addAsset.steps.step1.fields.assetId.hint')),
            required: computed(() => tsRaw('pages.assets.addAsset.steps.step1.fields.assetId.required')),
          },
          status: {
            label: computed(() => ts('pages.assets.addAsset.steps.step1.fields.status.label')),
            placeholder: computed(() => tsRaw('pages.assets.addAsset.steps.step1.fields.status.placeholder')),
            hint: computed(() => tsRaw('pages.assets.addAsset.steps.step1.fields.status.hint')),
            required: computed(() => tsRaw('pages.assets.addAsset.steps.step1.fields.status.required')),
          },
          description: {
            label: computed(() => ts('pages.assets.addAsset.steps.step1.fields.description.label')),
            placeholder: computed(() => tsRaw('pages.assets.addAsset.steps.step1.fields.description.placeholder')),
            hint: computed(() => tsRaw('pages.assets.addAsset.steps.step1.fields.description.hint')),
          },
        },
      },

      /**
       * Step 2: Asset Template
       */
      step2: {
        label: computed(() => ts('pages.assets.addAsset.steps.step2.label')),
        description: computed(() => tsRaw('pages.assets.addAsset.steps.step2.description')),
        title: computed(() => ts('pages.assets.addAsset.steps.step2.title')),
        subtitle: computed(() => tsRaw('pages.assets.addAsset.steps.step2.subtitle')),
        fields: {
          assetTemplateId: {
            label: computed(() => ts('pages.assets.addAsset.steps.step2.fields.assetTemplateId.label')),
            placeholder: computed(() => tsRaw('pages.assets.addAsset.steps.step2.fields.assetTemplateId.placeholder')),
            hint: computed(() => tsRaw('pages.assets.addAsset.steps.step2.fields.assetTemplateId.hint')),
            required: computed(() => tsRaw('pages.assets.addAsset.steps.step2.fields.assetTemplateId.required')),
            clearTooltip: computed(() => ts('pages.assets.addAsset.steps.step2.fields.assetTemplateId.clearTooltip')),
            searchTooltip: computed(() => ts('pages.assets.addAsset.steps.step2.fields.assetTemplateId.searchTooltip')),
          },
        },
        banner: {
          title: computed(() => ts('pages.assets.addAsset.steps.step2.banner.title')),
          description: computed(() => tsRaw('pages.assets.addAsset.steps.step2.banner.description')),
        },
        preview: {
          title: computed(() => ts('pages.assets.addAsset.steps.step2.preview.title')),
          detailsTitle: computed(() => ts('pages.assets.addAsset.steps.step2.preview.detailsTitle')),
          nameLabel: computed(() => ts('pages.assets.addAsset.steps.step2.preview.nameLabel')),
          descriptionLabel: computed(() => ts('pages.assets.addAsset.steps.step2.preview.descriptionLabel')),
          uuidPathLabel: computed(() => tsRaw('pages.assets.addAsset.steps.step2.preview.uuidPathLabel')),
          manufacturer: computed(() => ts('pages.assets.addAsset.steps.step2.preview.manufacturer')),
          model: computed(() => ts('pages.assets.addAsset.steps.step2.preview.model')),
          version: computed(() => ts('pages.assets.addAsset.steps.step2.preview.version')),
          status: computed(() => ts('pages.assets.addAsset.steps.step2.preview.status')),
          statusActive: computed(() => tsRaw('pages.assets.addAsset.steps.step2.preview.statusActive')),
          statusInactive: computed(() => tsRaw('pages.assets.addAsset.steps.step2.preview.statusInactive')),
          noTemplateSelected: computed(() => tsRaw('pages.assets.addAsset.steps.step2.preview.noTemplateSelected')),
        },
      },

      /**
       * Step 3: Route Groups
       */
      step3: {
        label: computed(() => ts('pages.assets.addAsset.steps.step3.label')),
        description: computed(() => ts('pages.assets.addAsset.steps.step3.description')),
        title: computed(() => ts('pages.assets.addAsset.steps.step3.title')),
        subtitle: computed(() => ts('pages.assets.addAsset.steps.step3.subtitle')),
        banner: {
          title: computed(() => ts('pages.assets.addAsset.steps.step3.banner.title')),
          description: computed(() => ts('pages.assets.addAsset.steps.step3.banner.description')),
        },
        labels: {
          systemTemplate: computed(() => ts('pages.assets.addAsset.steps.step3.labels.systemTemplate')),
          sharedTemplate: computed(() => ts('pages.assets.addAsset.steps.step3.labels.sharedTemplate')),
          selectedHeader: computed(() => ts('pages.assets.addAsset.steps.step3.labels.selectedHeader')),
          noneSelected: computed(() => ts('pages.assets.addAsset.steps.step3.labels.noneSelected')),
          routersConfigured: (count: number) => tsRaw('pages.assets.addAsset.steps.step3.labels.routersConfigured', { count }),
          selectedCount: (count: number) =>
            tsRaw(
              count === 1
                ? 'pages.assets.addAsset.steps.step3.labels.selectedCountSingular'
                : 'pages.assets.addAsset.steps.step3.labels.selectedCountPlural',
              { count },
            ),
          routingExplanation: computed(() => tsRaw('pages.assets.addAsset.steps.step3.labels.routingExplanation')),
        },
        buttons: {
          select: computed(() => ts('pages.assets.addAsset.steps.step3.buttons.select')),
          clearAll: computed(() => ts('pages.assets.addAsset.steps.step3.buttons.clearAll')),
        },
      },

      /**
       * Step 4: Connectivity
       */
      step4: {
        label: computed(() => ts('pages.assets.addAsset.steps.step4.label')),
        description: computed(() => tsRaw('pages.assets.addAsset.steps.step4.description')),
        title: computed(() => ts('pages.assets.addAsset.steps.step4.title')),
        subtitle: computed(() => tsRaw('pages.assets.addAsset.steps.step4.subtitle')),
        fields: {
          protocol: {
            label: computed(() => ts('pages.assets.addAsset.steps.step4.fields.protocol.label')),
            placeholder: computed(() => tsRaw('pages.assets.addAsset.steps.step4.fields.protocol.placeholder')),
            hint: computed(() => tsRaw('pages.assets.addAsset.steps.step4.fields.protocol.hint')),
            required: computed(() => tsRaw('pages.assets.addAsset.steps.step4.fields.protocol.required')),
          },
          mqttUsername: {
            label: computed(() => ts('pages.assets.addAsset.steps.step4.fields.mqttUsername.label')),
            placeholder: computed(() => tsRaw('pages.assets.addAsset.steps.step4.fields.mqttUsername.placeholder')),
            hint: computed(() => tsRaw('pages.assets.addAsset.steps.step4.fields.mqttUsername.hint')),
            hintDerived: computed(() => tsRaw('pages.assets.addAsset.steps.step4.fields.mqttUsername.hintDerived')),
            required: computed(() => tsRaw('pages.assets.addAsset.steps.step4.fields.mqttUsername.required')),
            minLength: computed(() => tsRaw('pages.assets.addAsset.steps.step4.fields.mqttUsername.minLength')),
          },
          mqttPassword: {
            label: computed(() => ts('pages.assets.addAsset.steps.step4.fields.mqttPassword.label')),
            placeholder: computed(() => tsRaw('pages.assets.addAsset.steps.step4.fields.mqttPassword.placeholder')),
            placeholderEdit: computed(() => tsRaw('pages.assets.addAsset.steps.step4.fields.mqttPassword.placeholderEdit')),
            hint: computed(() => tsRaw('pages.assets.addAsset.steps.step4.fields.mqttPassword.hint')),
            hintEdit: computed(() => tsRaw('pages.assets.addAsset.steps.step4.fields.mqttPassword.hintEdit')),
            required: computed(() => tsRaw('pages.assets.addAsset.steps.step4.fields.mqttPassword.required')),
            minLength: computed(() => tsRaw('pages.assets.addAsset.steps.step4.fields.mqttPassword.minLength')),
          },
          authType: {
            label: computed(() => ts('pages.assets.addAsset.steps.step4.fields.authType.label')),
            hint: computed(() => tsRaw('pages.assets.addAsset.steps.step4.fields.authType.hint')),
            required: computed(() => tsRaw('pages.assets.addAsset.steps.step4.fields.authType.required')),
            optionPassword: computed(() => ts('pages.assets.addAsset.steps.step4.fields.authType.optionPassword')),
            optionCert: computed(() => ts('pages.assets.addAsset.steps.step4.fields.authType.optionCert')),
            certBannerTitle: computed(() => ts('pages.assets.addAsset.steps.step4.fields.authType.certBannerTitle')),
            certBannerBody: computed(() => tsRaw('pages.assets.addAsset.steps.step4.fields.authType.certBannerBody')),
          },
          certTTL: {
            value: {
              label: computed(() => ts('pages.assets.addAsset.steps.step4.fields.certTTL.value.label')),
              hint: computed(() => tsRaw('pages.assets.addAsset.steps.step4.fields.certTTL.value.hint')),
              required: computed(() => tsRaw('pages.assets.addAsset.steps.step4.fields.certTTL.value.required')),
              maxDays: computed(() => tsRaw('pages.assets.addAsset.steps.step4.fields.certTTL.value.maxDays')),
            },
            unit: {
              label: computed(() => ts('pages.assets.addAsset.steps.step4.fields.certTTL.unit.label')),
              hint: computed(() => tsRaw('pages.assets.addAsset.steps.step4.fields.certTTL.unit.hint')),
              required: computed(() => tsRaw('pages.assets.addAsset.steps.step4.fields.certTTL.unit.required')),
            },
            units: {
              day: computed(() => ts('pages.assets.addAsset.steps.step4.fields.certTTL.units.day')),
              week: computed(() => ts('pages.assets.addAsset.steps.step4.fields.certTTL.units.week')),
              month: computed(() => ts('pages.assets.addAsset.steps.step4.fields.certTTL.units.month')),
              year: computed(() => ts('pages.assets.addAsset.steps.step4.fields.certTTL.units.year')),
            },
          },
          mqttClientId: {
            label: computed(() => ts('pages.assets.addAsset.steps.step4.fields.mqttClientId.label')),
            placeholder: computed(() => tsRaw('pages.assets.addAsset.steps.step4.fields.mqttClientId.placeholder')),
            hint: computed(() => tsRaw('pages.assets.addAsset.steps.step4.fields.mqttClientId.hint')),
            hintDerived: computed(() => tsRaw('pages.assets.addAsset.steps.step4.fields.mqttClientId.hintDerived')),
            required: computed(() => tsRaw('pages.assets.addAsset.steps.step4.fields.mqttClientId.required')),
            minLength: computed(() => tsRaw('pages.assets.addAsset.steps.step4.fields.mqttClientId.minLength')),
          },
          latitude: {
            label: computed(() => ts('pages.assets.addAsset.steps.step4.fields.latitude.label')),
            placeholder: computed(() => tsRaw('pages.assets.addAsset.steps.step4.fields.latitude.placeholder')),
            hint: computed(() => tsRaw('pages.assets.addAsset.steps.step4.fields.latitude.hint')),
          },
          longitude: {
            label: computed(() => ts('pages.assets.addAsset.steps.step4.fields.longitude.label')),
            placeholder: computed(() => tsRaw('pages.assets.addAsset.steps.step4.fields.longitude.placeholder')),
            hint: computed(() => tsRaw('pages.assets.addAsset.steps.step4.fields.longitude.hint')),
          },
          debugEnabled: {
            label: computed(() => ts('pages.assets.addAsset.steps.step4.fields.debugEnabled.label')),
            hint: computed(() => tsRaw('pages.assets.addAsset.steps.step4.fields.debugEnabled.hint')),
          },
        },
        protocolOptions: {
          http: computed(() => ts('pages.assets.addAsset.steps.step4.protocolOptions.http')),
          mqtt: computed(() => ts('pages.assets.addAsset.steps.step4.protocolOptions.mqtt')),
          lorawan: computed(() => ts('pages.assets.addAsset.steps.step4.protocolOptions.lorawan')),
        },
        banner: {
          info: computed(() => tsRaw('pages.assets.addAsset.steps.step4.banner.info')),
          httpInfo: computed(() => tsRaw('pages.assets.addAsset.steps.step4.banner.httpInfo')),
          mqttInfo: computed(() => tsRaw('pages.assets.addAsset.steps.step4.banner.mqttInfo')),
        },
        certDialog: {
          title: computed(() => ts('pages.assets.addAsset.steps.step4.certDialog.title')),
          warning: computed(() => tsRaw('pages.assets.addAsset.steps.step4.certDialog.warning')),
          replaceWarning: computed(() => tsRaw('pages.assets.addAsset.steps.step4.certDialog.replaceWarning')),
          generateButton: computed(() => ts('pages.assets.addAsset.steps.step4.certDialog.generateButton')),
          skipButton: computed(() => ts('pages.assets.addAsset.steps.step4.certDialog.skipButton')),
        },
        certificate: {
          activeTitle: computed(() => ts('pages.assets.addAsset.steps.step4.certificate.activeTitle')),
          revokedTitle: computed(() => ts('pages.assets.addAsset.steps.step4.certificate.revokedTitle')),
          revokedEmpty: computed(() => tsRaw('pages.assets.addAsset.steps.step4.certificate.revokedEmpty')),
          noActiveCert: computed(() => tsRaw('pages.assets.addAsset.steps.step4.certificate.noActiveCert')),
          fieldSerial: computed(() => ts('pages.assets.addAsset.steps.step4.certificate.fieldSerial')),
          fieldFingerprint: computed(() => ts('pages.assets.addAsset.steps.step4.certificate.fieldFingerprint')),
          fieldSubjectCN: computed(() => ts('pages.assets.addAsset.steps.step4.certificate.fieldSubjectCN')),
          fieldIssued: computed(() => ts('pages.assets.addAsset.steps.step4.certificate.fieldIssued')),
          fieldExpires: computed(() => ts('pages.assets.addAsset.steps.step4.certificate.fieldExpires')),
          revokeButton: computed(() => ts('pages.assets.addAsset.steps.step4.certificate.revokeButton')),
          revokeConfirm: computed(() => tsRaw('pages.assets.addAsset.steps.step4.certificate.revokeConfirm')),
        },
      },

      /**
       * Step 5: Health Monitoring
       */
      step5: {
        label: computed(() => ts('pages.assets.addAsset.steps.step5.label')),
        description: computed(() => tsRaw('pages.assets.addAsset.steps.step5.description')),
        title: computed(() => ts('pages.assets.addAsset.steps.step5.title')),
        subtitle: computed(() => tsRaw('pages.assets.addAsset.steps.step5.subtitle')),
      },

      /**
       * Step 6: Review
       */
      step6: {
        label: computed(() => ts('pages.assets.addAsset.steps.step6.label')),
        description: computed(() => tsRaw('pages.assets.addAsset.steps.step6.description')),
        title: computed(() => ts('pages.assets.addAsset.steps.step6.title')),
        subtitle: computed(() => tsRaw('pages.assets.addAsset.steps.step6.subtitle')),
        successMessage: computed(() => tsRaw('pages.assets.addAsset.steps.step6.successMessage')),
        sections: {
          identification: computed(() => ts('pages.assets.addAsset.steps.step6.sections.identification')),
          assetTemplate: computed(() => ts('pages.assets.addAsset.steps.step6.sections.assetTemplate')),
          routeGroups: computed(() => ts('pages.assets.addAsset.steps.step6.sections.routeGroups')),
          connectivity: computed(() => ts('pages.assets.addAsset.steps.step6.sections.connectivity')),
          healthMonitoring: computed(() => ts('pages.assets.addAsset.steps.step6.sections.healthMonitoring')),
        },
        fields: {
          name: computed(() => ts('pages.assets.addAsset.steps.step6.fields.name')),
          assetId: computed(() => ts('pages.assets.addAsset.steps.step6.fields.assetId')),
          status: computed(() => ts('pages.assets.addAsset.steps.step6.fields.status')),
          description: computed(() => ts('pages.assets.addAsset.steps.step6.fields.description')),
          assetTemplate: computed(() => ts('pages.assets.addAsset.steps.step6.fields.assetTemplate')),
          manufacturer: computed(() => ts('pages.assets.addAsset.steps.step6.fields.manufacturer')),
          model: computed(() => ts('pages.assets.addAsset.steps.step6.fields.model')),
          version: computed(() => ts('pages.assets.addAsset.steps.step6.fields.version')),
          routeGroups: computed(() => ts('pages.assets.addAsset.steps.step6.fields.routeGroups')),
          protocol: computed(() => ts('pages.assets.addAsset.steps.step6.fields.protocol')),
          mqttUsername: computed(() => ts('pages.assets.addAsset.steps.step6.fields.mqttUsername')),
          mqttClientId: computed(() => ts('pages.assets.addAsset.steps.step6.fields.mqttClientId')),
          latitude: computed(() => ts('pages.assets.addAsset.steps.step6.fields.latitude')),
          longitude: computed(() => ts('pages.assets.addAsset.steps.step6.fields.longitude')),
          healthMonitoringEnabled: computed(() => ts('pages.assets.addAsset.steps.step6.fields.healthMonitoringEnabled')),
          threshold: computed(() => ts('pages.assets.addAsset.steps.step6.fields.threshold')),
          requiredMisses: computed(() => ts('pages.assets.addAsset.steps.step6.fields.requiredMisses')),
          offlineRouteGroups: computed(() => ts('pages.assets.addAsset.steps.step6.fields.offlineRouteGroups')),
          onlineRouteGroups: computed(() => ts('pages.assets.addAsset.steps.step6.fields.onlineRouteGroups')),
        },
      },
    },

    /**
     * Status options translations
     */
    statusOptions: {
      active: {
        label: computed(() => ts('pages.assets.addAsset.statusOptions.active')),
        value: true,
      },
      inactive: {
        label: computed(() => ts('pages.assets.addAsset.statusOptions.inactive')),
        value: false,
      },
    },

    /**
     * Navigation buttons translations
     */
    navigation: {
      previous: computed(() => ts('pages.assets.addAsset.navigation.previous')),
      next: computed(() => ts('pages.assets.addAsset.navigation.next')),
      save: computed(() => ts('pages.assets.addAsset.navigation.save')),
      update: computed(() => ts('pages.assets.addAsset.navigation.update')),
    },

    /**
     * Notification messages translations
     */
    notifications: {
      created: computed(() => tsRaw('pages.assets.addAsset.notifications.created')),
      creationFailed: computed(() => tsRaw('pages.assets.addAsset.notifications.creationFailed')),
      validationFailed: computed(() => tsRaw('pages.assets.addAsset.notifications.validationFailed')),
      networkError: computed(() => tsRaw('pages.assets.addAsset.notifications.networkError')),
      alreadyExists: computed(() => tsRaw('pages.assets.addAsset.notifications.alreadyExists')),
      fillRequiredFields: computed(() => tsRaw('pages.assets.addAsset.notifications.fillRequiredFields')),
      loadingTemplates: computed(() => tsRaw('pages.assets.addAsset.notifications.loadingTemplates')),
      loading: computed(() => tsRaw('pages.assets.addAsset.notifications.loading')),
      loadFailed: computed(() => tsRaw('pages.assets.addAsset.notifications.loadFailed')),
      updated: computed(() => tsRaw('pages.assets.addAsset.notifications.updated')),
      updateFailed: computed(() => tsRaw('pages.assets.addAsset.notifications.updateFailed')),
      certIssueFailed: computed(() => tsRaw('pages.assets.addAsset.notifications.certIssueFailed')),
      assetCreatedWithoutCert: computed(() => tsRaw('pages.assets.addAsset.notifications.assetCreatedWithoutCert')),
    },
  };
}
