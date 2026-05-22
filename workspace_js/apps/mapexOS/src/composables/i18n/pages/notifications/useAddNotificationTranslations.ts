import { computed } from 'vue';
import { useTS } from '@utils/translation';

export function useAddNotificationTranslations() {
  const ts = useTS({ capitalize: true });
  const tsTitle = useTS({ titleCase: true });
  const tsRaw = useTS({ capitalize: false });

  return {
    page: {
      title: computed(() => tsTitle('pages.notifications.addNotification.title')),
      description: computed(() => tsRaw('pages.notifications.addNotification.description')),
      backButton: computed(() => tsRaw('pages.notifications.addNotification.backButton')),
    },
    steps: {
      general: {
        label: computed(() => ts('pages.notifications.addNotification.steps.general')),
        description: computed(() => tsRaw('pages.notifications.addNotification.stepDescriptions.general')),
        formDescription: computed(() => tsRaw('pages.notifications.addNotification.formDescriptions.general')),
      },
      configuration: {
        label: computed(() => ts('pages.notifications.addNotification.steps.configuration')),
        description: computed(() => tsRaw('pages.notifications.addNotification.stepDescriptions.configuration')),
        formDescription: computed(() => tsRaw('pages.notifications.addNotification.formDescriptions.configuration')),
      },
      review: {
        label: computed(() => ts('pages.notifications.addNotification.steps.review')),
        description: computed(() => tsRaw('pages.notifications.addNotification.stepDescriptions.review')),
        formDescription: computed(() => tsRaw('pages.notifications.addNotification.formDescriptions.review')),
      },
    },
    stepper: {
      title: computed(() => ts('pages.notifications.addNotification.stepper.title')),
      subtitle: computed(() => tsRaw('pages.notifications.addNotification.stepper.subtitle')),
    },
    fields: {
      type: computed(() => tsTitle('pages.notifications.addNotification.fields.type')),
      name: computed(() => tsTitle('pages.notifications.addNotification.fields.name')),
      description: computed(() => tsTitle('pages.notifications.addNotification.fields.description')),
      status: computed(() => tsTitle('pages.notifications.addNotification.fields.status')),
      created: computed(() => tsTitle('pages.notifications.addNotification.fields.created')),
    },
    errors: {
      unsupportedType: computed(() => ts('pages.notifications.addNotification.errors.unsupportedType')),
    },
    review: {
      generalInformation: computed(() => ts('pages.notifications.addNotification.review.generalInformation')),
      type: computed(() => tsTitle('pages.notifications.addNotification.review.type')),
      name: computed(() => tsTitle('pages.notifications.addNotification.review.name')),
      status: computed(() => tsTitle('pages.notifications.addNotification.review.status')),
      description: computed(() => tsTitle('pages.notifications.addNotification.review.description')),
      created: computed(() => tsTitle('pages.notifications.addNotification.review.created')),
    },
  };
}
