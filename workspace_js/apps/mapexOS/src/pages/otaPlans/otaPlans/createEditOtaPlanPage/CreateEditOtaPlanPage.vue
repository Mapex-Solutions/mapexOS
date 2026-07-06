<script setup lang="ts">
defineOptions({
  name: 'CreateEditOtaPlanPage'
});

/** TYPE IMPORTS */
import type { FormCardHeader, FormCardNavigation, FormCardButtonLabels } from '@components/cards';
import type { OtaPlanFormData } from './interfaces/createEditOtaPlan.interface';

/** VUE IMPORTS */
import { ref, computed } from 'vue';
import { useRouter } from 'vue-router';

/** COMPONENTS */
import { StepperVertical, buildStepperTree } from '@components/steppers';
import { PageHeader } from '@components/headers';
import { FormCard } from '@components/cards';
import { StepDetails, StepSchedule, StepFromTemplate, StepDevices, StepToTemplate, StepFirmware, StepReview } from './components';

/** COMPOSABLES */
import { useStepperNavigation } from '@composables/shared/form';
import { useCreateEditOtaPlanTranslations } from '@composables/i18n/pages/otaPlans/createEditOtaPlan/useCreateEditOtaPlanTranslations';

/** LOCAL IMPORTS */
import { STEP, TOTAL_STEPS, createInitialOtaPlanFormData } from './constants/createEditOtaPlan.constant';
import { useOtaPlanFormHandlers } from './handlers/useOtaPlanFormHandlers';

/** COMPOSABLES & STORES */
const t = useCreateEditOtaPlanTranslations();
const router = useRouter();

/** STATE */
const currentStep = ref(1);
const isSaving = ref(false);
const formData = ref<OtaPlanFormData>(createInitialOtaPlanFormData());
// The schedule step owns its own validity because "run now" and an empty
// scheduled datetime both leave startAt null but differ in validity.
const scheduleValid = ref(false);

const handlers = useOtaPlanFormHandlers({ formData, currentStep, isSaving, scheduleValid });

/** COMPUTED */

// Grouped tree for the vertical stepper; navigation still runs on the flat leaves.
const stepperTree = computed(() => buildStepperTree(handlers.translatedSteps.value));

/** FormCard header follows the active leaf step */
const formHeader = computed((): FormCardHeader => {
  const step = handlers.translatedSteps.value[currentStep.value - 1];
  return {
    icon: step?.icon || 'mdi-information',
    iconColor: 'primary',
    title: step?.title || '',
    description: step?.description || '',
  };
});

/** FormCard navigation state */
const formNavigation = computed((): FormCardNavigation => ({
  currentStep: currentStep.value,
  totalSteps: TOTAL_STEPS,
  showPreviousButton: currentStep.value > 1,
  showNextButton: currentStep.value < TOTAL_STEPS,
  showSaveButton: currentStep.value === TOTAL_STEPS,
  disableNextButton: handlers.isNextButtonDisabled.value,
  loadingSaveButton: isSaving.value,
}));

/** Wizard button labels */
const buttonLabels = computed((): FormCardButtonLabels => ({
  previous: t.buttons.back.value,
  next: t.buttons.next.value,
  save: t.buttons.submit.value,
}));

/** FUNCTIONS */

/**
 * Replace the form state (steps emit the full merged object)
 */
function updateFormData(value: OtaPlanFormData): void {
  formData.value = value;
}

/**
 * Changing the source template invalidates the device selection made against
 * the previous template's fleet.
 */
function onSourceTemplateChange(templateId: string | null): void {
  const cleared = templateId === formData.value.sourceTemplateId ? formData.value.selectedAssetIds : [];
  updateFormData({ ...formData.value, sourceTemplateId: templateId, selectedAssetIds: cleared });
}

/**
 * Create the plan and return to the list on success
 */
async function onSubmit(): Promise<void> {
  const ok = await handlers.submitForm();
  if (ok) {
    void router.push('/ota_plans');
  }
}

/** Keyboard arrows move between steps (inputs keep focus behavior) */
useStepperNavigation({
  currentStep,
  totalSteps: TOTAL_STEPS,
  changeStep: handlers.changeStep,
});
</script>

<template>
  <q-page class="q-pa-lg">
    <!-- Header Section -->
    <PageHeader
      icon="system_update"
      iconColor="primary"
      :title="t.page.title.value"
      :description="t.page.description.value"
      :button="{ label: t.page.button.value, icon: 'arrow_back', flat: true, to: '/ota_plans' }"
    />

    <!-- Content -->
    <div class="row q-col-gutter-lg">
      <!-- Progress Stepper Vertical (grouped: Setup / Rollout / Finalization) -->
      <div class="col-12 col-md-4">
        <StepperVertical
          :current-step="currentStep"
          :steps="stepperTree"
          @step-click="handlers.changeStep"
        />
      </div>

      <!-- Form Card -->
      <div class="col-12 col-md-8">
        <FormCard
          :header="formHeader"
          :navigation="formNavigation"
          :button-labels="buttonLabels"
          @previous="handlers.changeStep"
          @next="handlers.changeStep"
          @save="onSubmit"
        >
          <template #form>
            <!-- DETAILS -->
            <StepDetails
              v-if="currentStep === STEP.DETAILS"
              :name="formData.name"
              :description="formData.description"
              @update:name="(v) => updateFormData({ ...formData, name: v })"
              @update:description="(v) => updateFormData({ ...formData, description: v })"
            />

            <!-- SCHEDULE -->
            <StepSchedule
              v-else-if="currentStep === STEP.SCHEDULE"
              :start-at="formData.startAt"
              :max-time="formData.maxTime"
              :rollout-config="formData.rolloutConfig"
              @update:start-at="(v) => updateFormData({ ...formData, startAt: v })"
              @update:max-time="(v) => updateFormData({ ...formData, maxTime: v })"
              @update:rollout-config="(v) => updateFormData({ ...formData, rolloutConfig: v })"
              @update:valid="(v) => (scheduleValid = v)"
            />

            <!-- SOURCE TEMPLATE (from) -->
            <StepFromTemplate
              v-else-if="currentStep === STEP.FROM_TEMPLATE"
              :model-value="formData.sourceTemplateId"
              @update:model-value="onSourceTemplateChange"
            />

            <!-- DEVICES -->
            <StepDevices
              v-else-if="currentStep === STEP.DEVICES"
              :model-value="formData"
              @update:model-value="updateFormData"
            />

            <!-- TARGET TEMPLATE (to) -->
            <StepToTemplate
              v-else-if="currentStep === STEP.TO_TEMPLATE"
              :model-value="formData.targetTemplateId"
              @update:model-value="(v) => updateFormData({ ...formData, targetTemplateId: v })"
            />

            <!-- FIRMWARE -->
            <StepFirmware
              v-else-if="currentStep === STEP.FIRMWARE"
              :model-value="formData"
              @update:model-value="updateFormData"
            />

            <!-- REVIEW -->
            <StepReview
              v-else-if="currentStep === STEP.REVIEW"
              :model-value="formData"
            />
          </template>
        </FormCard>
      </div>
    </div>
  </q-page>
</template>
