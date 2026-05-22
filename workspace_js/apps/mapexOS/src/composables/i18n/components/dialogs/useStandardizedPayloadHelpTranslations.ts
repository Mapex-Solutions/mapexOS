import { computed } from 'vue';
import { useTS } from '@utils/translation';

/**
 * Custom composable for StandardizedPayloadHelpModal translations.
 * Provides type-safe reactive access to all translated strings with formatting utilities.
 *
 * Structure mirrors:
 * - File: src/components/dialogs/standardizedPayloadHelp/StandardizedPayloadHelpModal.vue
 * - JSON: src/i18n/{locale}/components/dialogs/standardizedPayloadHelp.json
 * - Composable: src/composables/i18n/components/dialogs/useStandardizedPayloadHelpTranslations.ts
 */
export function useStandardizedPayloadHelpTranslations() {
  const ts = useTS({ capitalize: true });
  const tsTitle = useTS({ titleCase: true });
  const tsRaw = useTS({ capitalize: false });

  return {
    /**
     * Header translations
     */
    header: {
      title: computed(() => tsTitle('components.dialogs.standardizedPayloadHelp.header.title')),
    },

    /**
     * Overview section translations
     */
    overview: {
      title: computed(() => ts('components.dialogs.standardizedPayloadHelp.overview.title')),
      description: computed(() => tsRaw('components.dialogs.standardizedPayloadHelp.overview.description')),
      whyTitle: computed(() => ts('components.dialogs.standardizedPayloadHelp.overview.whyTitle')),
      benefits: computed(() => [
        tsRaw('components.dialogs.standardizedPayloadHelp.overview.benefits.0'),
        tsRaw('components.dialogs.standardizedPayloadHelp.overview.benefits.1'),
        tsRaw('components.dialogs.standardizedPayloadHelp.overview.benefits.2'),
        tsRaw('components.dialogs.standardizedPayloadHelp.overview.benefits.3'),
      ]),
    },

    /**
     * Structure section translations
     */
    structure: {
      title: computed(() => ts('components.dialogs.standardizedPayloadHelp.structure.title')),
      copyTooltip: computed(() => tsRaw('components.dialogs.standardizedPayloadHelp.structure.copyTooltip')),
      comment1: computed(() => tsRaw('components.dialogs.standardizedPayloadHelp.structure.comment1')),
      comment2: computed(() => tsRaw('components.dialogs.standardizedPayloadHelp.structure.comment2')),
      comment3: computed(() => tsRaw('components.dialogs.standardizedPayloadHelp.structure.comment3')),
      comment4: computed(() => tsRaw('components.dialogs.standardizedPayloadHelp.structure.comment4')),
      comment5: computed(() => tsRaw('components.dialogs.standardizedPayloadHelp.structure.comment5')),
    },

    /**
     * Fields section translations
     */
    fields: {
      title: computed(() => ts('components.dialogs.standardizedPayloadHelp.fields.title')),
      required: computed(() => ts('components.dialogs.standardizedPayloadHelp.fields.required')),
      optional: computed(() => ts('components.dialogs.standardizedPayloadHelp.fields.optional')),
      eventType: {
        name: computed(() => ts('components.dialogs.standardizedPayloadHelp.fields.eventType.name')),
        type: computed(() => tsRaw('components.dialogs.standardizedPayloadHelp.fields.eventType.type')),
        description: computed(() => tsRaw('components.dialogs.standardizedPayloadHelp.fields.eventType.description')),
        examplesTitle: computed(() => ts('components.dialogs.standardizedPayloadHelp.fields.eventType.examplesTitle')),
        examples: computed(() => [
          tsRaw('components.dialogs.standardizedPayloadHelp.fields.eventType.examples.0'),
          tsRaw('components.dialogs.standardizedPayloadHelp.fields.eventType.examples.1'),
          tsRaw('components.dialogs.standardizedPayloadHelp.fields.eventType.examples.2'),
          tsRaw('components.dialogs.standardizedPayloadHelp.fields.eventType.examples.3'),
        ]),
      },
      eventId: {
        name: computed(() => ts('components.dialogs.standardizedPayloadHelp.fields.eventId.name')),
        type: computed(() => tsRaw('components.dialogs.standardizedPayloadHelp.fields.eventId.type')),
        description: computed(() => tsRaw('components.dialogs.standardizedPayloadHelp.fields.eventId.description')),
        examplesTitle: computed(() => ts('components.dialogs.standardizedPayloadHelp.fields.eventId.examplesTitle')),
        examples: computed(() => [
          tsRaw('components.dialogs.standardizedPayloadHelp.fields.eventId.examples.0'),
          tsRaw('components.dialogs.standardizedPayloadHelp.fields.eventId.examples.1'),
          tsRaw('components.dialogs.standardizedPayloadHelp.fields.eventId.examples.2'),
        ]),
      },
      data: {
        name: computed(() => ts('components.dialogs.standardizedPayloadHelp.fields.data.name')),
        type: computed(() => tsRaw('components.dialogs.standardizedPayloadHelp.fields.data.type')),
        description: computed(() => tsRaw('components.dialogs.standardizedPayloadHelp.fields.data.description')),
        exampleTitle: computed(() => ts('components.dialogs.standardizedPayloadHelp.fields.data.exampleTitle')),
      },
      metadata: {
        name: computed(() => ts('components.dialogs.standardizedPayloadHelp.fields.metadata.name')),
        type: computed(() => tsRaw('components.dialogs.standardizedPayloadHelp.fields.metadata.type')),
        description: computed(() => tsRaw('components.dialogs.standardizedPayloadHelp.fields.metadata.description')),
        exampleTitle: computed(() => ts('components.dialogs.standardizedPayloadHelp.fields.metadata.exampleTitle')),
      },
      created: {
        name: computed(() => ts('components.dialogs.standardizedPayloadHelp.fields.created.name')),
        type: computed(() => tsRaw('components.dialogs.standardizedPayloadHelp.fields.created.type')),
        description: computed(() => tsRaw('components.dialogs.standardizedPayloadHelp.fields.created.description')),
        validFormatsTitle: computed(() => ts('components.dialogs.standardizedPayloadHelp.fields.created.validFormatsTitle')),
        formats: computed(() => [
          tsRaw('components.dialogs.standardizedPayloadHelp.fields.created.formats.0'),
          tsRaw('components.dialogs.standardizedPayloadHelp.fields.created.formats.1'),
          tsRaw('components.dialogs.standardizedPayloadHelp.fields.created.formats.2'),
        ]),
      },
    },

    /**
     * Example section translations
     */
    example: {
      title: computed(() => ts('components.dialogs.standardizedPayloadHelp.example.title')),
      description: computed(() => tsRaw('components.dialogs.standardizedPayloadHelp.example.description')),
      copyTooltip: computed(() => tsRaw('components.dialogs.standardizedPayloadHelp.example.copyTooltip')),
    },

    /**
     * Common mistakes section translations
     */
    mistakes: {
      title: computed(() => ts('components.dialogs.standardizedPayloadHelp.mistakes.title')),
      items: computed(() => [
        {
          label: ts('components.dialogs.standardizedPayloadHelp.mistakes.items.0.label'),
          caption: tsRaw('components.dialogs.standardizedPayloadHelp.mistakes.items.0.caption'),
        },
        {
          label: ts('components.dialogs.standardizedPayloadHelp.mistakes.items.1.label'),
          caption: tsRaw('components.dialogs.standardizedPayloadHelp.mistakes.items.1.caption'),
        },
        {
          label: ts('components.dialogs.standardizedPayloadHelp.mistakes.items.2.label'),
          caption: tsRaw('components.dialogs.standardizedPayloadHelp.mistakes.items.2.caption'),
        },
        {
          label: ts('components.dialogs.standardizedPayloadHelp.mistakes.items.3.label'),
          caption: tsRaw('components.dialogs.standardizedPayloadHelp.mistakes.items.3.caption'),
        },
        {
          label: ts('components.dialogs.standardizedPayloadHelp.mistakes.items.4.label'),
          caption: tsRaw('components.dialogs.standardizedPayloadHelp.mistakes.items.4.caption'),
        },
      ]),
    },

    /**
     * Footer section translations
     */
    footer: {
      info: computed(() => tsRaw('components.dialogs.standardizedPayloadHelp.footer.info')),
      close: computed(() => ts('components.dialogs.standardizedPayloadHelp.footer.close')),
    },

    /**
     * Notification messages translations
     */
    notifications: {
      copySuccess: computed(() => tsRaw('components.dialogs.standardizedPayloadHelp.notifications.copySuccess')),
      copyFail: computed(() => tsRaw('components.dialogs.standardizedPayloadHelp.notifications.copyFail')),
    },
  };
}
