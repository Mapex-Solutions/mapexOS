<script setup lang="ts">
defineOptions({ name: 'CreateEditWorkflowInstancePage' });

/** TYPE IMPORTS */
import type { FormCardHeader } from '@components/cards';
import type { WorkflowInstanceFormState } from './interfaces';

/** VUE IMPORTS */
import { ref, computed, onMounted } from 'vue';
import { useRouter, useRoute } from 'vue-router';

/** COMPONENTS */
import { StepperVertical } from '@components/steppers';
import { PageHeader } from '@components/headers';
import { FormCard } from '@components/cards';
import {
  Step1Identification,
  Step2Definition,
  Step3ExternalInputs,
  Step4Review,
} from './components';

/** COMPOSABLES */
import { useStepperNavigation } from '@composables/shared/form';
import { useCreateEditWorkflowInstanceTranslations } from '@src/composables/i18n/pages/automations/workflowInstances/createEditWorkflowInstancePage/useCreateEditWorkflowInstanceTranslations';

/** LOCAL IMPORTS */
import { INITIAL_FORM_DATA, TOTAL_STEPS, STEP } from './constants';
import { useWorkflowInstanceFormHandlers } from './handlers';

/** SERVICES */
import { apis } from '@services/mapex';

/** UTILS */
import { handleApiError } from '@utils/error';

/** COMPOSABLES & STORES */
const t = useCreateEditWorkflowInstanceTranslations();
const router = useRouter();
const route = useRoute();

/** EDIT MODE DETECTION */
const isEditMode = ref(!!route.params.id);
const instanceId = ref(route.params.id as string | undefined);

/** LOADING STATES */
const isLoading = ref(false);
const isSaving = ref(false);

/** STATE */
const step1Ref = ref<InstanceType<typeof Step1Identification> | null>(null);
const step2Ref = ref<InstanceType<typeof Step2Definition> | null>(null);
const step3Ref = ref<InstanceType<typeof Step3ExternalInputs> | null>(null);

const currentStep = ref(1);
const formData = ref({ ...INITIAL_FORM_DATA });
const formState = ref<WorkflowInstanceFormState>({
  selectedDefinition: null,
  isSaving: false,
  currentStep: 1,
});

/** FUNCTIONS */

/**
 * Load instance data for edit mode.
 * Fetches instance from API and populates form state.
 * @returns {Promise<void>}
 */
async function loadInstanceData(): Promise<void> {
  if (!isEditMode.value || !instanceId.value) return;

  isLoading.value = true;
  try {
    const data = await apis.workflows.instance.getById({
      instanceId: instanceId.value,
    });

    formData.value.name = data.name || '';
    formData.value.description = data.description || '';
    formData.value.enabled = data.enabled ?? true;
    formData.value.definitionId = data.definitionId || null;
    formData.value.definitionVersion = data.definitionVersion || 1;
    formData.value.externalInputs = data.externalInputs || {};

    // Load the linked definition for display
    if (data.definitionId) {
      try {
        const defData = await apis.workflows.definition.getById({
          workflowId: data.definitionId,
        });
        formData.value.selectedDefinition = defData;
        formState.value.selectedDefinition = defData;
      } catch {
        // Continue even if definition fetch fails
      }
    }

    // In edit mode, skip to Review step
    currentStep.value = STEP.REVIEW;
  } catch (error: any) {
    handleApiError(error, {
      defaultMessage: t.notifications.loadFailed.value,
      timeout: 5000,
    });
    await router.push('/workflow_instances');
  } finally {
    isLoading.value = false;
  }
}

/**
 * Update form data with partial updates from step components
 * @param {Partial<typeof formData.value>} partialData - Partial data to merge
 * @returns {void}
 */
function updateFormData(partialData: Partial<typeof formData.value>): void {
  formData.value = { ...formData.value, ...partialData };
}

/** COMPUTED */

const pageTitle = computed(() =>
  isEditMode.value ? t.page.titleEdit.value : t.page.title.value
);

const translatedSteps = computed(() => [
  {
    title: t.steps.step1.label.value,
    icon: 'mdi-fingerprint',
    description: t.steps.step1.description.value,
  },
  {
    title: t.steps.step2.label.value,
    icon: 'mdi-file-tree',
    description: t.steps.step2.description.value,
  },
  {
    title: t.steps.step3.label.value,
    icon: 'mdi-form-textbox',
    description: t.steps.step3.description.value,
  },
  {
    title: t.steps.step4.label.value,
    icon: 'mdi-clipboard-check',
    description: t.steps.step4.description.value,
  },
]);

const handlers = useWorkflowInstanceFormHandlers({
  formData,
  formState,
  currentStep,
  isEditMode,
  instanceId,
  isSaving,
  step1FormRef: computed(() => step1Ref.value?.formRef ?? null),
  step2FormRef: computed(() => step2Ref.value?.formRef ?? null),
  step3FormRef: computed(() => step3Ref.value?.formRef ?? null),
});

const formNavigation = computed(() => ({
  currentStep: currentStep.value,
  totalSteps: TOTAL_STEPS,
  showPreviousButton: true,
  showNextButton: true,
  showSaveButton: true,
  showCancelButton: true,
  disableNextButton: handlers.isNextButtonDisabled.value || isSaving.value,
  disableSaveButton: handlers.isNextButtonDisabled.value || isSaving.value,
  loadingSaveButton: isSaving.value,
}));

const buttonLabels = computed(() => ({
  previous: t.navigation.previous.value,
  next: t.navigation.next.value,
  save: isEditMode.value ? t.navigation.update.value : t.navigation.save.value,
}));

/** COMPOSABLES USAGE */
useStepperNavigation({
  currentStep,
  totalSteps: TOTAL_STEPS,
  changeStep: (step: number) => { void handlers.handleStepChange(step); },
});

/** LIFECYCLE HOOKS */
onMounted(() => {
  void loadInstanceData();
});
</script>

<template>
  <q-page class="q-pa-lg">
    <!-- Loading State (Edit Mode) -->
    <div v-if="isLoading" class="flex flex-center q-pa-xl">
      <q-spinner color="primary" size="50px" />
      <span class="q-ml-md text-grey-7">{{ t.notifications.loading.value }}</span>
    </div>

    <!-- Form Content -->
    <div v-else>
      <!-- Header -->
      <PageHeader
        icon="memory"
        iconColor="primary"
        :title="pageTitle"
        :description="t.page.description.value"
        :button="{ label: t.page.button.value, icon: 'arrow_back', flat: true, to: '/workflow_instances' }"
      />

      <!-- Content -->
      <div class="row q-col-gutter-lg">
        <!-- Stepper -->
        <div class="col-12 col-md-4">
          <StepperVertical
            :title="t.stepper.title.value"
            :subtitle="t.stepper.subtitle.value"
            :info-text="t.stepper.requiredInfo.value"
            :current-step-label="t.stepper.currentStep.value"
            :current-step="currentStep"
            :steps="translatedSteps"
            :allow-step-navigation="isEditMode"
            @step-click="handlers.changeStep"
          />
        </div>

        <!-- Form Card -->
        <div class="col-12 col-md-8">
          <FormCard
            :header="translatedSteps[currentStep - 1] as unknown as FormCardHeader"
            :navigation="formNavigation"
            :button-labels="buttonLabels"
            @previous="handlers.changeStep"
            @next="handlers.changeStep"
            @save="handlers.submitForm"
          >
            <template #form>
              <!-- STEP 1: IDENTIFICATION -->
              <Step1Identification
                v-if="currentStep === STEP.IDENTIFICATION"
                ref="step1Ref"
                :model-value="formData"
                @update:model-value="updateFormData"
              />

              <!-- STEP 2: DEFINITION -->
              <Step2Definition
                v-else-if="currentStep === STEP.DEFINITION"
                ref="step2Ref"
                :model-value="formData"
                @update:model-value="updateFormData"
                @definition-selected="handlers.onDefinitionSelected"
              />

              <!-- STEP 3: EXTERNAL INPUTS -->
              <Step3ExternalInputs
                v-else-if="currentStep === STEP.EXTERNAL_INPUTS"
                ref="step3Ref"
                :model-value="formData"
                @update:model-value="updateFormData"
              />

              <!-- STEP 4: REVIEW -->
              <Step4Review
                v-else-if="currentStep === STEP.REVIEW"
                :model-value="formData"
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
