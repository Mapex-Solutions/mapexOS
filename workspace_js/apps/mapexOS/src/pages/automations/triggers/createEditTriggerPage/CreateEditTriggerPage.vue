<script setup lang="ts">
defineOptions({
  name: 'CreateEditTriggerPage'
});

/** TYPE IMPORTS */
import type { FormCardHeader } from '@components/cards';

/** VUE IMPORTS */
import { ref, computed, onMounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';

/** COMPONENTS */
import { StepperVertical } from '@components/steppers';
import { PageHeader } from '@components/headers';
import { FormCard } from '@components/cards';
import {
  Step1Category,
  Step2Type,
  Step3BasicInfo,
  Step4Configuration,
  Step5Review,
} from './components';

/** COMPOSABLES */
import { useStepperNavigation } from '@composables/shared/form';
import { useCreateEditTriggerTranslations } from '@composables/i18n/pages/automations/triggers/createEditTrigger';

/** UTILS */
import { notifyFail } from '@utils/alert';

/** LOCAL IMPORTS */
import {
  INITIAL_TRIGGER_FORM_DATA,
  INITIAL_FORM_STATE,
  TOTAL_STEPS,
  STEP,
} from './constants';

import { useTriggerFormHandlers } from './handlers';

/** ROUTER */
const route = useRoute();
const router = useRouter();

/** COMPOSABLES & STORES */
const t = useCreateEditTriggerTranslations();

/** EDIT MODE DETECTION */
const isEditMode = ref(!!route.params.id);
const triggerId = ref(route.params.id as string | undefined);

/** STATE */

/**
 * Loading state for fetching trigger data in edit mode
 */
const isLoading = ref(false);

/**
 * Reference to Step 3 basic info form
 */
const step3Ref = ref<InstanceType<typeof Step3BasicInfo> | null>(null);

/**
 * Reference to Step 4 configuration form
 */
const step4Ref = ref<InstanceType<typeof Step4Configuration> | null>(null);

/**
 * Current step in the wizard
 */
const currentStep = ref(1);

/**
 * Trigger form data
 */
const triggerData = ref({ ...INITIAL_TRIGGER_FORM_DATA });

/**
 * Form state for managing navigation and validation
 */
const formState = ref({ ...INITIAL_FORM_STATE });

/** COMPUTED */

/**
 * Dynamic page title based on mode
 */
const pageTitle = computed(() =>
  isEditMode.value ? t.page.titleEdit.value : t.page.title.value
);

/**
 * Dynamic page description based on mode
 */
const pageDescription = computed(() => t.page.description.value);

/**
 * Dynamic save button label based on mode
 */
const saveButtonLabel = computed(() =>
  isEditMode.value ? t.navigation.update.value : t.navigation.save.value
);

/**
 * Steps configuration for the stepper
 */
const steps = computed(() => [
  {
    title: t.steps.step1.label.value,
    icon: 'mdi-shape',
    description: t.steps.step1.description.value,
  },
  {
    title: t.steps.step2.label.value,
    icon: 'mdi-format-list-bulleted-type',
    description: t.steps.step2.description.value,
  },
  {
    title: t.steps.step3.label.value,
    icon: 'mdi-information',
    description: t.steps.step3.description.value,
  },
  {
    title: t.steps.step4.label.value,
    icon: 'mdi-cog',
    description: t.steps.step4.description.value,
  },
  {
    title: t.steps.step5.label.value,
    icon: 'mdi-clipboard-check',
    description: t.steps.step5.description.value,
  },
]);

/**
 * Trigger form handlers
 */
const handlers = useTriggerFormHandlers({
  triggerData,
  formState,
  currentStep,
  step3FormRef: computed(() => step3Ref.value?.formRef ?? null),
  step4FormRef: computed(() => step4Ref.value?.formRef ?? null),
  isEditMode,
  triggerId,
  t,
});

/**
 * Form navigation configuration
 */
const formNavigation = computed(() => ({
  currentStep: currentStep.value,
  totalSteps: TOTAL_STEPS,
  showPreviousButton: true,
  showNextButton: true,
  showSaveButton: true,
  showCancelButton: true,
  disableNextButton: handlers.isNextButtonDisabled.value || formState.value.isCreating,
  disableSaveButton: handlers.isNextButtonDisabled.value || formState.value.isCreating,
  loadingSaveButton: formState.value.isCreating,
}));

/**
 * Button labels
 */
const buttonLabels = computed(() => ({
  previous: t.navigation.previous.value,
  next: t.navigation.next.value,
  save: saveButtonLabel.value,
}));

/** COMPOSABLES USAGE */
useStepperNavigation({
  currentStep,
  totalSteps: TOTAL_STEPS,
  changeStep: handlers.handleStepChange,
});

/** FUNCTIONS */

/**
 * Load trigger data from API in EDIT mode
 * @returns {Promise<void>}
 */
async function loadTriggerData(): Promise<void> {
  if (!isEditMode.value || !triggerId.value) return;

  isLoading.value = true;
  try {
    const loadedTrigger = await handlers.loadTriggerData();

    if (loadedTrigger) {
      // Update trigger data
      triggerData.value = loadedTrigger;

      // Update form state with loaded values
      formState.value.selectedCategory = loadedTrigger.category;
      formState.value.selectedType = loadedTrigger.triggerType;

      // In EDIT mode, skip to Review step by default
      currentStep.value = STEP.REVIEW;
      formState.value.currentStep = STEP.REVIEW;
    } else {
      // Failed to load - redirect to list
      notifyFail({ message: t.notifications.loadFailed.value });
      void router.push('/triggers');
    }
  } finally {
    isLoading.value = false;
  }
}

/**
 * Update trigger data with partial updates from step components
 * @param {Partial<typeof triggerData.value>} partialData - Partial trigger data to merge
 * @returns {void}
 */
function updateTriggerData(partialData: Partial<typeof triggerData.value>): void {
  triggerData.value = {
    ...triggerData.value,
    ...partialData,
  };
}

/** LIFECYCLE HOOKS */
onMounted(() => {
  if (isEditMode.value) {
    void loadTriggerData();
  }
});
</script>

<template>
  <q-page class="q-pa-lg">
    <!-- Loading Spinner for EDIT mode data fetch -->
    <div v-if="isLoading" class="flex flex-center q-pa-xl">
      <div class="column items-center">
        <q-spinner color="primary" size="50px" />
        <span class="q-mt-md text-grey-7">{{ t.notifications.loading.value }}</span>
      </div>
    </div>

    <!-- Main Content -->
    <div v-else>
      <!-- Header Section -->
      <PageHeader
        icon="flash_on"
        iconColor="primary"
        :title="pageTitle"
        :description="pageDescription"
        :button="{ label: t.page.button.value, icon: 'arrow_back', flat: true, to: '/triggers' }"
      />

      <!-- Content -->
      <div class="row q-col-gutter-lg">
        <!-- Progress Stepper Vertical -->
        <div class="col-12 col-md-4">
          <StepperVertical
            :title="t.stepper.title.value"
            :subtitle="t.stepper.subtitle.value"
            :info-text="t.stepper.requiredInfo.value"
            :current-step-label="t.stepper.currentStep.value"
            :current-step="currentStep"
            :steps="steps"
            :allow-step-navigation="isEditMode"
            @step-click="handlers.changeStep"
          />
        </div>

        <!-- Form Card -->
        <div class="col-12 col-md-8">
          <FormCard
            :header="steps[currentStep - 1] as unknown as FormCardHeader"
            :navigation="formNavigation"
            :button-labels="buttonLabels"
            @previous="handlers.changeStep"
            @next="handlers.changeStep"
            @save="handlers.submitForm"
          >
            <!-- FORM BODY -->
            <template #form>
              <!-- STEP 1: CATEGORY -->
              <Step1Category
                v-if="currentStep === STEP.CATEGORY"
                :model-value="triggerData"
                @update:model-value="updateTriggerData"
                @category-selected="handlers.onCategorySelected"
              />

              <!-- STEP 2: TYPE -->
              <Step2Type
                v-else-if="currentStep === STEP.TYPE"
                :model-value="triggerData"
                @update:model-value="updateTriggerData"
                @type-selected="handlers.onTypeSelected"
              />

              <!-- STEP 3: BASIC INFO -->
              <Step3BasicInfo
                v-else-if="currentStep === STEP.BASIC_INFO"
                ref="step3Ref"
                :model-value="triggerData"
                @update:model-value="updateTriggerData"
              />

              <!-- STEP 4: CONFIGURATION -->
              <Step4Configuration
                v-else-if="currentStep === STEP.CONFIGURATION"
                ref="step4Ref"
                :model-value="triggerData"
                @update:model-value="updateTriggerData"
              />

              <!-- STEP 5: REVIEW -->
              <Step5Review
                v-else-if="currentStep === STEP.REVIEW"
                :model-value="triggerData"
                :form-state="formState"
                @edit-section="handlers.changeStep"
              />
            </template>
          </FormCard>
        </div>
      </div>
    </div>
  </q-page>
</template>

<style lang="scss" scoped>
</style>
