import { computed } from 'vue';
import { useTS } from '@utils/translation';

export function useAddLakeHouseTranslations() {
  const ts = useTS({ capitalize: true });
  const tsTitle = useTS({ titleCase: true });
  const tsRaw = useTS({ capitalize: false });

  return {
    page: {
      title: computed(() => tsTitle('pages.lakeHouse.addLakeHouse.title')),
      description: computed(() => tsRaw('pages.lakeHouse.addLakeHouse.description')),
      backButton: computed(() => tsRaw('pages.lakeHouse.addLakeHouse.backButton')),
    },
    steps: {
      general: {
        label: computed(() => ts('pages.lakeHouse.addLakeHouse.steps.general')),
        description: computed(() => tsRaw('pages.lakeHouse.addLakeHouse.stepDescriptions.general')),
        formDescription: computed(() => tsRaw('pages.lakeHouse.addLakeHouse.formDescriptions.general')),
      },
      provider: {
        label: computed(() => ts('pages.lakeHouse.addLakeHouse.steps.provider')),
        description: computed(() => tsRaw('pages.lakeHouse.addLakeHouse.stepDescriptions.provider')),
        formDescription: computed(() => tsRaw('pages.lakeHouse.addLakeHouse.formDescriptions.provider')),
      },
      credentials: {
        label: computed(() => ts('pages.lakeHouse.addLakeHouse.steps.credentials')),
        description: computed(() => tsRaw('pages.lakeHouse.addLakeHouse.stepDescriptions.credentials')),
        formDescription: computed(() => tsRaw('pages.lakeHouse.addLakeHouse.formDescriptions.credentials')),
      },
      pathConfig: {
        label: computed(() => ts('pages.lakeHouse.addLakeHouse.steps.pathConfig')),
        description: computed(() => tsRaw('pages.lakeHouse.addLakeHouse.stepDescriptions.pathConfig')),
        formDescription: computed(() => tsRaw('pages.lakeHouse.addLakeHouse.formDescriptions.pathConfig')),
      },
      frequency: {
        label: computed(() => ts('pages.lakeHouse.addLakeHouse.steps.frequency')),
        description: computed(() => tsRaw('pages.lakeHouse.addLakeHouse.stepDescriptions.frequency')),
        formDescription: computed(() => tsRaw('pages.lakeHouse.addLakeHouse.formDescriptions.frequency')),
      },
      review: {
        label: computed(() => ts('pages.lakeHouse.addLakeHouse.steps.review')),
        description: computed(() => tsRaw('pages.lakeHouse.addLakeHouse.stepDescriptions.review')),
        formDescription: computed(() => tsRaw('pages.lakeHouse.addLakeHouse.formDescriptions.review')),
      },
    },
    stepper: {
      title: computed(() => ts('pages.lakeHouse.addLakeHouse.stepper.title')),
      subtitle: computed(() => tsRaw('pages.lakeHouse.addLakeHouse.stepper.subtitle')),
    },
    fields: {
      type: computed(() => tsTitle('pages.lakeHouse.addLakeHouse.fields.type')),
      name: computed(() => tsTitle('pages.lakeHouse.addLakeHouse.fields.name')),
      description: computed(() => tsTitle('pages.lakeHouse.addLakeHouse.fields.description')),
      status: computed(() => tsTitle('pages.lakeHouse.addLakeHouse.fields.status')),
    },
  };
}
