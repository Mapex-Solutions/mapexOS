import { computed } from 'vue';
import { useTS } from '@utils/translation';

/**
 * Custom composable for AddAssetTemplate page translations.
 * Provides type-safe reactive access to all translated strings with formatting utilities.
 *
 * Structure mirrors:
 * - File: src/pages/assets/assetTemplates/addAssetTemplate/AddAssetTemplate.vue
 * - JSON: src/i18n/{locale}/pages/assets/addAssetTemplate.json
 * - Composable: src/composables/i18n/pages/assets/addAssetTemplate/useAddAssetTemplateTranslations.ts
 */
export function useAddAssetTemplateTranslations() {
  const ts = useTS({ capitalize: true });
  const tsTitle = useTS({ titleCase: true });
  const tsRaw = useTS({ capitalize: false });

  return {
    /**
     * Convention banner shown above every Monaco editor reminding users that
     * scripts must assign to `result` (not use `return`).
     */
    scriptConvention: {
      title: computed(() => ts('pages.assets.addAssetTemplate.scriptConvention.title')),
      textPrefix: computed(() => tsRaw('pages.assets.addAssetTemplate.scriptConvention.textPrefix')),
      textSuffix: computed(() => tsRaw('pages.assets.addAssetTemplate.scriptConvention.textSuffix')),
    },

    /**
     * Page header translations
     */
    page: {
      title: computed(() => ts('pages.assets.addAssetTemplate.page.title')),
      titleEdit: computed(() => ts('pages.assets.addAssetTemplate.page.titleEdit')),
      description: computed(() => tsRaw('pages.assets.addAssetTemplate.page.description')),
      button: computed(() => ts('pages.assets.addAssetTemplate.page.button')),
    },

    /**
     * Stepper header translations
     */
    stepper: {
      title: computed(() => tsTitle('pages.assets.addAssetTemplate.stepper.title')),
      subtitle: computed(() => tsRaw('pages.assets.addAssetTemplate.stepper.subtitle')),
      requiredInfo: computed(() => tsRaw('pages.assets.addAssetTemplate.stepper.requiredInfo')),
      currentStep: computed(() => ts('pages.assets.addAssetTemplate.stepper.currentStep')),
    },

    /**
     * Step-specific translations
     */
    steps: {
      /**
       * Step 1: Basic Information
       */
      step1: {
        label: computed(() => ts('pages.assets.addAssetTemplate.steps.step1.label')),
        description: computed(() => tsRaw('pages.assets.addAssetTemplate.steps.step1.description')),
        title: computed(() => ts('pages.assets.addAssetTemplate.steps.step1.title')),
        subtitle: computed(() => tsRaw('pages.assets.addAssetTemplate.steps.step1.subtitle')),
        fields: {
          name: {
            label: computed(() => ts('pages.assets.addAssetTemplate.steps.step1.fields.name.label')),
            placeholder: computed(() => tsRaw('pages.assets.addAssetTemplate.steps.step1.fields.name.placeholder')),
            hint: computed(() => tsRaw('pages.assets.addAssetTemplate.steps.step1.fields.name.hint')),
            required: computed(() => tsRaw('pages.assets.addAssetTemplate.steps.step1.fields.name.required')),
          },
          status: {
            label: computed(() => ts('pages.assets.addAssetTemplate.steps.step1.fields.status.label')),
            placeholder: computed(() => tsRaw('pages.assets.addAssetTemplate.steps.step1.fields.status.placeholder')),
            hint: computed(() => tsRaw('pages.assets.addAssetTemplate.steps.step1.fields.status.hint')),
            required: computed(() => tsRaw('pages.assets.addAssetTemplate.steps.step1.fields.status.required')),
          },
          description: {
            label: computed(() => ts('pages.assets.addAssetTemplate.steps.step1.fields.description.label')),
            placeholder: computed(() => tsRaw('pages.assets.addAssetTemplate.steps.step1.fields.description.placeholder')),
            hint: computed(() => tsRaw('pages.assets.addAssetTemplate.steps.step1.fields.description.hint')),
          },
          manufacture: {
            label: computed(() => ts('pages.assets.addAssetTemplate.steps.step1.fields.manufacture.label')),
            placeholder: computed(() => tsRaw('pages.assets.addAssetTemplate.steps.step1.fields.manufacture.placeholder')),
            hint: computed(() => tsRaw('pages.assets.addAssetTemplate.steps.step1.fields.manufacture.hint')),
            required: computed(() => tsRaw('pages.assets.addAssetTemplate.steps.step1.fields.manufacture.required')),
          },
          model: {
            label: computed(() => ts('pages.assets.addAssetTemplate.steps.step1.fields.model.label')),
            placeholder: computed(() => tsRaw('pages.assets.addAssetTemplate.steps.step1.fields.model.placeholder')),
            hint: computed(() => tsRaw('pages.assets.addAssetTemplate.steps.step1.fields.model.hint')),
            required: computed(() => tsRaw('pages.assets.addAssetTemplate.steps.step1.fields.model.required')),
          },
          version: {
            label: computed(() => ts('pages.assets.addAssetTemplate.steps.step1.fields.version.label')),
            placeholder: computed(() => tsRaw('pages.assets.addAssetTemplate.steps.step1.fields.version.placeholder')),
            hint: computed(() => tsRaw('pages.assets.addAssetTemplate.steps.step1.fields.version.hint')),
          },
          isTemplate: {
            label: computed(() => ts('pages.assets.addAssetTemplate.steps.step1.fields.isTemplate.label')),
            hint: computed(() => tsRaw('pages.assets.addAssetTemplate.steps.step1.fields.isTemplate.hint')),
          },
        },
      },

      /**
       * Step 2: Asset ID Path
       */
      step2: {
        label: computed(() => ts('pages.assets.addAssetTemplate.steps.step2.label')),
        description: computed(() => tsRaw('pages.assets.addAssetTemplate.steps.step2.description')),
        title: computed(() => ts('pages.assets.addAssetTemplate.steps.step2.title')),
        subtitle: computed(() => tsRaw('pages.assets.addAssetTemplate.steps.step2.subtitle')),
        banner: {
          title: computed(() => ts('pages.assets.addAssetTemplate.steps.step2.banner.title')),
          description: computed(() => tsRaw('pages.assets.addAssetTemplate.steps.step2.banner.description')),
          examplesTitle: computed(() => ts('pages.assets.addAssetTemplate.steps.step2.banner.examplesTitle')),
          examples: computed(() => [
            tsRaw('pages.assets.addAssetTemplate.steps.step2.banner.examples.0'),
            tsRaw('pages.assets.addAssetTemplate.steps.step2.banner.examples.1'),
            tsRaw('pages.assets.addAssetTemplate.steps.step2.banner.examples.2'),
          ]),
        },
        fields: {
          assetIdPath: {
            label: computed(() => ts('pages.assets.addAssetTemplate.steps.step2.fields.assetIdPath.label')),
            placeholder: computed(() => tsRaw('pages.assets.addAssetTemplate.steps.step2.fields.assetIdPath.placeholder')),
            hint: computed(() => tsRaw('pages.assets.addAssetTemplate.steps.step2.fields.assetIdPath.hint')),
            required: computed(() => tsRaw('pages.assets.addAssetTemplate.steps.step2.fields.assetIdPath.required')),
          },
        },
      },

      /**
       * Step 3: Preprocessor Script
       */
      step3: {
        label: computed(() => ts('pages.assets.addAssetTemplate.steps.step3.label')),
        description: computed(() => tsRaw('pages.assets.addAssetTemplate.steps.step3.description')),
        title: computed(() => ts('pages.assets.addAssetTemplate.steps.step3.title')),
        subtitle: computed(() => tsRaw('pages.assets.addAssetTemplate.steps.step3.subtitle')),
        guidelines: {
          title: computed(() => ts('pages.assets.addAssetTemplate.steps.step3.guidelines.title')),
          availableVariables: computed(() => ts('pages.assets.addAssetTemplate.steps.step3.guidelines.availableVariables')),
          payloadVar: computed(() => tsRaw('pages.assets.addAssetTemplate.steps.step3.guidelines.payloadVar')),
          consoleLog: computed(() => tsRaw('pages.assets.addAssetTemplate.steps.step3.guidelines.consoleLog')),
          expectedReturn: computed(() => ts('pages.assets.addAssetTemplate.steps.step3.guidelines.expectedReturn')),
          expectedReturnValue: computed(() => tsRaw('pages.assets.addAssetTemplate.steps.step3.guidelines.expectedReturnValue')),
          example: computed(() => ts('pages.assets.addAssetTemplate.steps.step3.guidelines.example')),
        },
        exampleCode: {
          comment1: computed(() => tsRaw('pages.assets.addAssetTemplate.steps.step3.exampleCode.comment1')),
          comment2: computed(() => tsRaw('pages.assets.addAssetTemplate.steps.step3.exampleCode.comment2')),
          comment3: computed(() => tsRaw('pages.assets.addAssetTemplate.steps.step3.exampleCode.comment3')),
          comment4: computed(() => tsRaw('pages.assets.addAssetTemplate.steps.step3.exampleCode.comment4')),
          comment5: computed(() => tsRaw('pages.assets.addAssetTemplate.steps.step3.exampleCode.comment5')),
          comment6: computed(() => tsRaw('pages.assets.addAssetTemplate.steps.step3.exampleCode.comment6')),
        },
      },

      /**
       * Step 4: Validation Script
       */
      step4: {
        label: computed(() => ts('pages.assets.addAssetTemplate.steps.step4.label')),
        description: computed(() => tsRaw('pages.assets.addAssetTemplate.steps.step4.description')),
        title: computed(() => ts('pages.assets.addAssetTemplate.steps.step4.title')),
        subtitle: computed(() => tsRaw('pages.assets.addAssetTemplate.steps.step4.subtitle')),
        guidelines: {
          title: computed(() => ts('pages.assets.addAssetTemplate.steps.step4.guidelines.title')),
          availableVariables: computed(() => ts('pages.assets.addAssetTemplate.steps.step4.guidelines.availableVariables')),
          payloadVar: computed(() => tsRaw('pages.assets.addAssetTemplate.steps.step4.guidelines.payloadVar')),
          expectedReturn: computed(() => ts('pages.assets.addAssetTemplate.steps.step4.guidelines.expectedReturn')),
          expectedReturnValue: computed(() => tsRaw('pages.assets.addAssetTemplate.steps.step4.guidelines.expectedReturnValue')),
          example: computed(() => ts('pages.assets.addAssetTemplate.steps.step4.guidelines.example')),
        },
        exampleCode: {
          comment1: computed(() => tsRaw('pages.assets.addAssetTemplate.steps.step4.exampleCode.comment1')),
          comment2: computed(() => tsRaw('pages.assets.addAssetTemplate.steps.step4.exampleCode.comment2')),
          comment3: computed(() => tsRaw('pages.assets.addAssetTemplate.steps.step4.exampleCode.comment3')),
          comment4: computed(() => tsRaw('pages.assets.addAssetTemplate.steps.step4.exampleCode.comment4')),
          comment5: computed(() => tsRaw('pages.assets.addAssetTemplate.steps.step4.exampleCode.comment5')),
          comment6: computed(() => tsRaw('pages.assets.addAssetTemplate.steps.step4.exampleCode.comment6')),
          comment7: computed(() => tsRaw('pages.assets.addAssetTemplate.steps.step4.exampleCode.comment7')),
          rangeError: computed(() => tsRaw('pages.assets.addAssetTemplate.steps.step4.exampleCode.rangeError')),
        },
      },

      /**
       * Step 5: Conversion Script
       */
      step5: {
        label: computed(() => ts('pages.assets.addAssetTemplate.steps.step5.label')),
        description: computed(() => tsRaw('pages.assets.addAssetTemplate.steps.step5.description')),
        title: computed(() => ts('pages.assets.addAssetTemplate.steps.step5.title')),
        subtitle: computed(() => tsRaw('pages.assets.addAssetTemplate.steps.step5.subtitle')),
        banner: {
          important: computed(() => tsRaw('pages.assets.addAssetTemplate.steps.step5.banner.important')),
        },
        guidelines: {
          title: computed(() => ts('pages.assets.addAssetTemplate.steps.step5.guidelines.title')),
          availableVariables: computed(() => ts('pages.assets.addAssetTemplate.steps.step5.guidelines.availableVariables')),
          payloadVar: computed(() => tsRaw('pages.assets.addAssetTemplate.steps.step5.guidelines.payloadVar')),
          expectedReturn: computed(() => ts('pages.assets.addAssetTemplate.steps.step5.guidelines.expectedReturn')),
          expectedReturnValue: computed(() => tsRaw('pages.assets.addAssetTemplate.steps.step5.guidelines.expectedReturnValue')),
          standardizedPayloadFormat: computed(() => ts('pages.assets.addAssetTemplate.steps.step5.guidelines.standardizedPayloadFormat')),
          standardizedPayloadDescription: computed(() => tsRaw('pages.assets.addAssetTemplate.steps.step5.guidelines.standardizedPayloadDescription')),
          requiredFields: computed(() => ts('pages.assets.addAssetTemplate.steps.step5.guidelines.requiredFields')),
          optionalFields: computed(() => ts('pages.assets.addAssetTemplate.steps.step5.guidelines.optionalFields')),
          eventType: computed(() => tsRaw('pages.assets.addAssetTemplate.steps.step5.guidelines.eventType')),
          eventId: computed(() => tsRaw('pages.assets.addAssetTemplate.steps.step5.guidelines.eventId')),
          data: computed(() => tsRaw('pages.assets.addAssetTemplate.steps.step5.guidelines.data')),
          created: computed(() => tsRaw('pages.assets.addAssetTemplate.steps.step5.guidelines.created')),
          metadata: computed(() => tsRaw('pages.assets.addAssetTemplate.steps.step5.guidelines.metadata')),
          example: computed(() => ts('pages.assets.addAssetTemplate.steps.step5.guidelines.example')),
        },
        exampleCode: {
          comment1: computed(() => tsRaw('pages.assets.addAssetTemplate.steps.step5.exampleCode.comment1')),
          comment2: computed(() => tsRaw('pages.assets.addAssetTemplate.steps.step5.exampleCode.comment2')),
          comment3: computed(() => tsRaw('pages.assets.addAssetTemplate.steps.step5.exampleCode.comment3')),
          comment4: computed(() => tsRaw('pages.assets.addAssetTemplate.steps.step5.exampleCode.comment4')),
          comment5: computed(() => tsRaw('pages.assets.addAssetTemplate.steps.step5.exampleCode.comment5')),
        },
      },

      /**
       * Step 6: Test Payload
       */
      step6: {
        label: computed(() => ts('pages.assets.addAssetTemplate.steps.step6.label')),
        description: computed(() => tsRaw('pages.assets.addAssetTemplate.steps.step6.description')),
        title: computed(() => ts('pages.assets.addAssetTemplate.steps.step6.title')),
        subtitle: computed(() => tsRaw('pages.assets.addAssetTemplate.steps.step6.subtitle')),
        fields: {
          testInput: {
            label: computed(() => ts('pages.assets.addAssetTemplate.steps.step6.fields.testInput.label')),
          },
        },
        buttons: {
          loadExample: computed(() => ts('pages.assets.addAssetTemplate.steps.step6.buttons.loadExample')),
        },
        banner: {
          info: computed(() => tsRaw('pages.assets.addAssetTemplate.steps.step6.banner.info')),
        },
        guidelines: {
          title: computed(() => tsTitle('pages.assets.addAssetTemplate.steps.step6.guidelines.title')),
          availableVariables: computed(() => tsRaw('pages.assets.addAssetTemplate.steps.step6.guidelines.availableVariables')),
          payloadVar: computed(() => tsRaw('pages.assets.addAssetTemplate.steps.step6.guidelines.payloadVar')),
          expectedReturn: computed(() => tsRaw('pages.assets.addAssetTemplate.steps.step6.guidelines.expectedReturn')),
          example: computed(() => tsRaw('pages.assets.addAssetTemplate.steps.step6.guidelines.example')),
        },
        exampleCode: {
          comment1: computed(() => tsRaw('pages.assets.addAssetTemplate.steps.step6.exampleCode.comment1')),
          comment2: computed(() => tsRaw('pages.assets.addAssetTemplate.steps.step6.exampleCode.comment2')),
          comment3: computed(() => tsRaw('pages.assets.addAssetTemplate.steps.step6.exampleCode.comment3')),
          comment4: computed(() => tsRaw('pages.assets.addAssetTemplate.steps.step6.exampleCode.comment4')),
          comment5: computed(() => tsRaw('pages.assets.addAssetTemplate.steps.step6.exampleCode.comment5')),
          comment6: computed(() => tsRaw('pages.assets.addAssetTemplate.steps.step6.exampleCode.comment6')),
        },
      },

      /**
       * Step 7: Testing & Review
       */
      step7: {
        label: computed(() => ts('pages.assets.addAssetTemplate.steps.step7.label')),
        description: computed(() => tsRaw('pages.assets.addAssetTemplate.steps.step7.description')),
        title: computed(() => ts('pages.assets.addAssetTemplate.steps.step7.title')),
        subtitle: computed(() => tsRaw('pages.assets.addAssetTemplate.steps.step7.subtitle')),
        banner: {
          testRequired: computed(() => tsRaw('pages.assets.addAssetTemplate.steps.step7.banner.testRequired')),
          standardizedPayloadInfo: computed(() => tsRaw('pages.assets.addAssetTemplate.steps.step7.banner.standardizedPayloadInfo')),
          viewFormat: computed(() => ts('pages.assets.addAssetTemplate.steps.step7.banner.viewFormat')),
        },
        testSection: {
          title: computed(() => ts('pages.assets.addAssetTemplate.steps.step7.testSection.title')),
          usingPayloadFrom: computed(() => tsRaw('pages.assets.addAssetTemplate.steps.step7.testSection.usingPayloadFrom')),
          usingPayloadFromInline: computed(() => tsRaw('pages.assets.addAssetTemplate.steps.step7.testSection.usingPayloadFromInline')),
          editLink: computed(() => ts('pages.assets.addAssetTemplate.steps.step7.testSection.editLink')),
          detailsLabel: computed(() => ts('pages.assets.addAssetTemplate.steps.step7.testSection.detailsLabel')),
          modifyPayloadHint: computed(() => tsRaw('pages.assets.addAssetTemplate.steps.step7.testSection.modifyPayloadHint')),
          statusLabel: computed(() => ts('pages.assets.addAssetTemplate.steps.step7.testSection.statusLabel')),
          payloadLabel: computed(() => ts('pages.assets.addAssetTemplate.steps.step7.testSection.payloadLabel')),
          newPayloadLabel: computed(() => ts('pages.assets.addAssetTemplate.steps.step7.testSection.newPayloadLabel')),
          errorDetailsLabel: computed(() => ts('pages.assets.addAssetTemplate.steps.step7.testSection.errorDetailsLabel')),
          testInput: {
            label: computed(() => ts('pages.assets.addAssetTemplate.steps.step7.testSection.testInput.label')),
            placeholder: computed(() => tsRaw('pages.assets.addAssetTemplate.steps.step7.testSection.testInput.placeholder')),
          },
          testResults: {
            label: computed(() => ts('pages.assets.addAssetTemplate.steps.step7.testSection.testResults.label')),
            emptyState: computed(() => tsRaw('pages.assets.addAssetTemplate.steps.step7.testSection.testResults.emptyState')),
            statusLabel: computed(() => ts('pages.assets.addAssetTemplate.steps.step7.testSection.testResults.statusLabel')),
            errorDetailsLabel: computed(() => ts('pages.assets.addAssetTemplate.steps.step7.testSection.testResults.errorDetailsLabel')),
          },
          buttons: {
            runTest: computed(() => ts('pages.assets.addAssetTemplate.steps.step7.testSection.buttons.runTest')),
            clear: computed(() => ts('pages.assets.addAssetTemplate.steps.step7.testSection.buttons.clear')),
            clearResults: computed(() => ts('pages.assets.addAssetTemplate.steps.step7.testSection.buttons.clearResults')),
          },
          status: {
            notExecuted: computed(() => tsRaw('pages.assets.addAssetTemplate.steps.step7.testSection.status.notExecuted')),
            allPassed: computed(() => ts('pages.assets.addAssetTemplate.steps.step7.testSection.status.allPassed')),
            testFailed: computed(() => ts('pages.assets.addAssetTemplate.steps.step7.testSection.status.testFailed')),
            running: computed(() => tsRaw('pages.assets.addAssetTemplate.steps.step7.testSection.status.running')),
            invalidJson: computed(() => tsRaw('pages.assets.addAssetTemplate.steps.step7.testSection.status.invalidJson')),
          },
          steps: {
            preprocessor: computed(() => ts('pages.assets.addAssetTemplate.steps.step7.testSection.steps.preprocessor')),
            validation: computed(() => ts('pages.assets.addAssetTemplate.steps.step7.testSection.steps.validation')),
            conversion: computed(() => ts('pages.assets.addAssetTemplate.steps.step7.testSection.steps.conversion')),
            testScript: computed(() => ts('pages.assets.addAssetTemplate.steps.step7.testSection.steps.testScript')),
            error: computed(() => ts('pages.assets.addAssetTemplate.steps.step7.testSection.steps.error')),
            skipped: computed(() => tsRaw('pages.assets.addAssetTemplate.steps.step7.testSection.steps.skipped')),
          },
          results: {
            finalOutput: computed(() => tsRaw('pages.assets.addAssetTemplate.steps.step7.testSection.results.finalOutput')),
            standardizedPayload: computed(() => ts('pages.assets.addAssetTemplate.steps.step7.testSection.results.standardizedPayload')),
            consoleOutput: computed(() => tsRaw('pages.assets.addAssetTemplate.steps.step7.testSection.results.consoleOutput')),
            validationFailed: computed(() => tsRaw('pages.assets.addAssetTemplate.steps.step7.testSection.results.validationFailed')),
          },
        },
        reviewSection: {
          title: computed(() => ts('pages.assets.addAssetTemplate.steps.step7.reviewSection.title')),
          description: computed(() => tsRaw('pages.assets.addAssetTemplate.steps.step7.reviewSection.description')),
          basicInfo: computed(() => ts('pages.assets.addAssetTemplate.steps.step7.reviewSection.basicInfo')),
          scriptSummary: computed(() => ts('pages.assets.addAssetTemplate.steps.step7.reviewSection.scriptSummary')),
          fields: {
            name: computed(() => ts('pages.assets.addAssetTemplate.steps.step7.reviewSection.fields.name')),
            status: computed(() => ts('pages.assets.addAssetTemplate.steps.step7.reviewSection.fields.status')),
            manufacturer: computed(() => ts('pages.assets.addAssetTemplate.steps.step7.reviewSection.fields.manufacturer')),
            model: computed(() => ts('pages.assets.addAssetTemplate.steps.step7.reviewSection.fields.model')),
            category: computed(() => ts('pages.assets.addAssetTemplate.steps.step7.reviewSection.fields.category')),
            version: computed(() => ts('pages.assets.addAssetTemplate.steps.step7.reviewSection.fields.version')),
            description: computed(() => ts('pages.assets.addAssetTemplate.steps.step7.reviewSection.fields.description')),
            assetIdPath: computed(() => ts('pages.assets.addAssetTemplate.steps.step7.reviewSection.fields.assetIdPath')),
            preprocessorScript: computed(() => ts('pages.assets.addAssetTemplate.steps.step7.reviewSection.fields.preprocessorScript')),
            validationScript: computed(() => ts('pages.assets.addAssetTemplate.steps.step7.reviewSection.fields.validationScript')),
            conversionScript: computed(() => ts('pages.assets.addAssetTemplate.steps.step7.reviewSection.fields.conversionScript')),
            testScript: computed(() => ts('pages.assets.addAssetTemplate.steps.step7.reviewSection.fields.testScript')),
            configured: computed(() => ts('pages.assets.addAssetTemplate.steps.step7.reviewSection.fields.configured')),
            notConfigured: computed(() => ts('pages.assets.addAssetTemplate.steps.step7.reviewSection.fields.notConfigured')),
          },
        },
      },

      /**
       * Step 8: Dynamic Fields
       */
      step8: {
        label: computed(() => ts('pages.assets.addAssetTemplate.steps.step8.label')),
        description: computed(() => tsRaw('pages.assets.addAssetTemplate.steps.step8.description')),
        title: computed(() => ts('pages.assets.addAssetTemplate.steps.step8.title')),
        subtitle: computed(() => tsRaw('pages.assets.addAssetTemplate.steps.step8.subtitle')),
        banner: {
          title: computed(() => ts('pages.assets.addAssetTemplate.steps.step8.banner.title')),
          description: computed(() => tsRaw('pages.assets.addAssetTemplate.steps.step8.banner.description')),
        },
        noFieldsWarning: {
          title: computed(() => ts('pages.assets.addAssetTemplate.steps.step8.noFieldsWarning.title')),
          description: computed(() => tsRaw('pages.assets.addAssetTemplate.steps.step8.noFieldsWarning.description')),
        },
        addField: {
          title: computed(() => ts('pages.assets.addAssetTemplate.steps.step8.addField.title')),
          fieldName: {
            label: computed(() => ts('pages.assets.addAssetTemplate.steps.step8.addField.fieldName.label')),
            placeholder: computed(() => tsRaw('pages.assets.addAssetTemplate.steps.step8.addField.fieldName.placeholder')),
          },
          fieldType: {
            label: computed(() => ts('pages.assets.addAssetTemplate.steps.step8.addField.fieldType.label')),
          },
          valuePath: {
            label: computed(() => ts('pages.assets.addAssetTemplate.steps.step8.addField.valuePath.label')),
            placeholder: computed(() => tsRaw('pages.assets.addAssetTemplate.steps.step8.addField.valuePath.placeholder')),
            noOptions: computed(() => tsRaw('pages.assets.addAssetTemplate.steps.step8.addField.valuePath.noOptions')),
          },
          latitudePath: {
            label: computed(() => ts('pages.assets.addAssetTemplate.steps.step8.addField.latitudePath.label')),
            placeholder: computed(() => tsRaw('pages.assets.addAssetTemplate.steps.step8.addField.latitudePath.placeholder')),
          },
          longitudePath: {
            label: computed(() => ts('pages.assets.addAssetTemplate.steps.step8.addField.longitudePath.label')),
            placeholder: computed(() => tsRaw('pages.assets.addAssetTemplate.steps.step8.addField.longitudePath.placeholder')),
          },
          addButton: computed(() => ts('pages.assets.addAssetTemplate.steps.step8.addField.addButton')),
        },
        mappedFields: {
          title: computed(() => ts('pages.assets.addAssetTemplate.steps.step8.mappedFields.title')),
          removeTooltip: computed(() => ts('pages.assets.addAssetTemplate.steps.step8.mappedFields.removeTooltip')),
        },
        emptyState: {
          title: computed(() => ts('pages.assets.addAssetTemplate.steps.step8.emptyState.title')),
          description: computed(() => tsRaw('pages.assets.addAssetTemplate.steps.step8.emptyState.description')),
        },
      },

      /**
       * Step 9: Review
       */
      step9: {
        label: computed(() => ts('pages.assets.addAssetTemplate.steps.step9.label')),
        description: computed(() => tsRaw('pages.assets.addAssetTemplate.steps.step9.description')),
        title: computed(() => ts('pages.assets.addAssetTemplate.steps.step9.title')),
        subtitle: computed(() => tsRaw('pages.assets.addAssetTemplate.steps.step9.subtitle')),
        completeBanner: computed(() => tsRaw('pages.assets.addAssetTemplate.steps.step9.completeBanner')),
        availableFields: {
          title: computed(() => ts('pages.assets.addAssetTemplate.steps.step9.availableFields.title')),
          rerunTooltip: computed(() => ts('pages.assets.addAssetTemplate.steps.step9.availableFields.rerunTooltip')),
        },
        noFieldsWarning: computed(() => tsRaw('pages.assets.addAssetTemplate.steps.step9.noFieldsWarning')),
        goToTesting: computed(() => ts('pages.assets.addAssetTemplate.steps.step9.goToTesting')),
        dynamicFields: {
          title: computed(() => ts('pages.assets.addAssetTemplate.steps.step9.dynamicFields.title')),
          editTooltip: computed(() => ts('pages.assets.addAssetTemplate.steps.step9.dynamicFields.editTooltip')),
          empty: computed(() => tsRaw('pages.assets.addAssetTemplate.steps.step9.dynamicFields.empty')),
        },
        reviewSection: {
          basicInfo: computed(() => ts('pages.assets.addAssetTemplate.steps.step9.reviewSection.basicInfo')),
          scriptSummary: computed(() => ts('pages.assets.addAssetTemplate.steps.step9.reviewSection.scriptSummary')),
          fields: {
            name: computed(() => ts('pages.assets.addAssetTemplate.steps.step9.reviewSection.fields.name')),
            status: computed(() => ts('pages.assets.addAssetTemplate.steps.step9.reviewSection.fields.status')),
            manufacturer: computed(() => ts('pages.assets.addAssetTemplate.steps.step9.reviewSection.fields.manufacturer')),
            model: computed(() => ts('pages.assets.addAssetTemplate.steps.step9.reviewSection.fields.model')),
            category: computed(() => ts('pages.assets.addAssetTemplate.steps.step9.reviewSection.fields.category')),
            version: computed(() => ts('pages.assets.addAssetTemplate.steps.step9.reviewSection.fields.version')),
            description: computed(() => ts('pages.assets.addAssetTemplate.steps.step9.reviewSection.fields.description')),
            assetIdPath: computed(() => ts('pages.assets.addAssetTemplate.steps.step9.reviewSection.fields.assetIdPath')),
            preprocessorScript: computed(() => ts('pages.assets.addAssetTemplate.steps.step9.reviewSection.fields.preprocessorScript')),
            validationScript: computed(() => ts('pages.assets.addAssetTemplate.steps.step9.reviewSection.fields.validationScript')),
            conversionScript: computed(() => ts('pages.assets.addAssetTemplate.steps.step9.reviewSection.fields.conversionScript')),
            testScript: computed(() => ts('pages.assets.addAssetTemplate.steps.step9.reviewSection.fields.testScript')),
            configured: computed(() => ts('pages.assets.addAssetTemplate.steps.step9.reviewSection.fields.configured')),
            notConfigured: computed(() => ts('pages.assets.addAssetTemplate.steps.step9.reviewSection.fields.notConfigured')),
          },
        },
      },
    },

    /**
     * Tour translations for Asset Template wizard
     */
    tour: {
      pageOverview: {
        title: computed(() => ts('pages.assets.addAssetTemplate.tour.pageOverview.title')),
        description: computed(() => tsRaw('pages.assets.addAssetTemplate.tour.pageOverview.description')),
      },
      stepper: {
        title: computed(() => ts('pages.assets.addAssetTemplate.tour.stepper.title')),
        description: computed(() => tsRaw('pages.assets.addAssetTemplate.tour.stepper.description')),
      },
      assetIdPath: {
        title: computed(() => ts('pages.assets.addAssetTemplate.tour.assetIdPath.title')),
        description: computed(() => tsRaw('pages.assets.addAssetTemplate.tour.assetIdPath.description')),
      },
      preprocessorScript: {
        title: computed(() => ts('pages.assets.addAssetTemplate.tour.preprocessorScript.title')),
        description: computed(() => tsRaw('pages.assets.addAssetTemplate.tour.preprocessorScript.description')),
      },
      validationScript: {
        title: computed(() => ts('pages.assets.addAssetTemplate.tour.validationScript.title')),
        description: computed(() => tsRaw('pages.assets.addAssetTemplate.tour.validationScript.description')),
      },
      conversionScript: {
        title: computed(() => ts('pages.assets.addAssetTemplate.tour.conversionScript.title')),
        description: computed(() => tsRaw('pages.assets.addAssetTemplate.tour.conversionScript.description')),
      },
      testPayload: {
        title: computed(() => ts('pages.assets.addAssetTemplate.tour.testPayload.title')),
        description: computed(() => tsRaw('pages.assets.addAssetTemplate.tour.testPayload.description')),
      },
      testing: {
        title: computed(() => ts('pages.assets.addAssetTemplate.tour.testing.title')),
        description: computed(() => tsRaw('pages.assets.addAssetTemplate.tour.testing.description')),
      },
      dynamicFields: {
        title: computed(() => ts('pages.assets.addAssetTemplate.tour.dynamicFields.title')),
        description: computed(() => tsRaw('pages.assets.addAssetTemplate.tour.dynamicFields.description')),
      },
      review: {
        title: computed(() => ts('pages.assets.addAssetTemplate.tour.review.title')),
        description: computed(() => tsRaw('pages.assets.addAssetTemplate.tour.review.description')),
      },
    },

    /**
     * Status options translations
     */
    statusOptions: {
      active: {
        label: computed(() => ts('pages.assets.addAssetTemplate.statusOptions.active')),
        value: true,
      },
      inactive: {
        label: computed(() => ts('pages.assets.addAssetTemplate.statusOptions.inactive')),
        value: false,
      },
    },

    /**
     * Navigation buttons translations
     */
    navigation: {
      back: computed(() => ts('pages.assets.addAssetTemplate.navigation.back')),
      previous: computed(() => ts('pages.assets.addAssetTemplate.navigation.previous')),
      next: computed(() => ts('pages.assets.addAssetTemplate.navigation.next')),
      review: computed(() => ts('pages.assets.addAssetTemplate.navigation.review')),
      save: computed(() => ts('pages.assets.addAssetTemplate.navigation.save')),
      update: computed(() => ts('pages.assets.addAssetTemplate.navigation.update')),
    },

    /**
     * Notification messages translations
     */
    notifications: {
      created: computed(() => tsRaw('pages.assets.addAssetTemplate.notifications.created')),
      updated: computed(() => tsRaw('pages.assets.addAssetTemplate.notifications.updated')),
      loading: computed(() => tsRaw('pages.assets.addAssetTemplate.notifications.loading')),
      loadFailed: computed(() => tsRaw('pages.assets.addAssetTemplate.notifications.loadFailed')),
      updateFailed: computed(() => tsRaw('pages.assets.addAssetTemplate.notifications.updateFailed')),
      validationError: computed(() => tsRaw('pages.assets.addAssetTemplate.notifications.validationError')),
      testSuccess: computed(() => tsRaw('pages.assets.addAssetTemplate.notifications.testSuccess')),
      testFailed: computed(() => tsRaw('pages.assets.addAssetTemplate.notifications.testFailed')),
      invalidScript: computed(() => tsRaw('pages.assets.addAssetTemplate.notifications.invalidScript')),
      conversionScriptRequired: computed(() => tsRaw('pages.assets.addAssetTemplate.notifications.conversionScriptRequired')),
      testsMustPass: computed(() => tsRaw('pages.assets.addAssetTemplate.notifications.testsMustPass')),
      creating: computed(() => tsRaw('pages.assets.addAssetTemplate.notifications.creating')),
      creationFailed: computed(() => tsRaw('pages.assets.addAssetTemplate.notifications.creationFailed')),
      validationFailed: computed(() => tsRaw('pages.assets.addAssetTemplate.notifications.validationFailed')),
      alreadyExists: computed(() => tsRaw('pages.assets.addAssetTemplate.notifications.alreadyExists')),
      networkError: computed(() => tsRaw('pages.assets.addAssetTemplate.notifications.networkError')),
    },
  };
}
