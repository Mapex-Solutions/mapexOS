<script setup lang="ts">
import type { FormCardHeader } from '@components/cards';

import { ref, computed } from 'vue';
import { QForm } from 'quasar';
import { STEPS } from './constants';
import { buildLakeHousePreview } from './handlers';

import { useStepperNavigation } from '@composables/shared/form';
import { useAddLakeHouseTranslations } from '@composables/i18n/pages/lakeHouse/useAddLakeHouseTranslations';
import { useLogger } from '@composables/useLogger';

import { DEFAULT_AWS_DATA } from '@components/forms/lakeHouse/lakeHouseProviderSelection/constants';

const logger = useLogger('AddLakeHouse');
import { DEFAULT_PATH_CONFIG } from '@components/forms/lakeHouse/lakeHousePathConfig/constants';
import { DEFAULT_FREQUENCY_CONFIG } from '@components/forms/lakeHouse/lakeHouseFrequency/constants';

defineOptions({
  name: 'AddLakeHouse'
});

/** COMPONENTS */
import { StepperVertical } from '@components/steppers';
import { PageHeader } from '@components/headers';
import { FormCard } from '@components/cards';

import {
  LakeHouseBaseInfo,
  LakeHouseProviderSelection,
  LakeHouseCredentials,
  LakeHousePathConfig,
  LakeHouseFrequency,
  FormReview,
} from '@components/forms';

const t = useAddLakeHouseTranslations();

const formRef = ref<QForm | null>(null);
const currentStep = ref(1);
const isEditMode = ref(false);

const lakeHouseData = ref<any>({
  type: 'aws-s3',
  name: '',
  description: '',
  status: true,
  credentials: DEFAULT_AWS_DATA,
  pathConfig: DEFAULT_PATH_CONFIG,
  frequency: DEFAULT_FREQUENCY_CONFIG,
});

const previewLakeHouse = computed(() => buildLakeHousePreview(lakeHouseData.value));

//computed properties
const formNavigation = computed(() => ({
  currentStep: currentStep.value,
  totalSteps: STEPS.length,
  showPreviousButton: true,
  showNextButton: true,
  showSaveButton: true,
  showCancelButton: true,
}));

// Step change handler
useStepperNavigation({
  currentStep,
  totalSteps: STEPS.length,
  changeStep,
});

// Methods
function changeStep(step: number) {
  currentStep.value = step;
}

function submitForm() {
  logger.debug('Form submitted:', lakeHouseData.value);
}

</script>

<template>
  <q-page class="q-pa-lg">
    <!-- Header Section -->
    <PageHeader
        icon="cloud_upload"
        iconColor="primary"
        :title="t.page.title.value"
        :description="t.page.description.value"
        :button="{ label: t.page.backButton.value, icon: 'arrow_back', flat: true, to: '/lakeHouses' }"
    />

    <!-- Content -->
    <div class="row q-col-gutter-lg">
      <!-- Progress Stepper Vertical -->
      <div class="col-12 col-md-4">
        <StepperVertical
            :title="t.stepper.title.value"
            :subtitle="t.stepper.subtitle.value"
            :current-step="currentStep"
            :steps="STEPS"
            :mode="isEditMode ? 'editing' : 'creating'"
            @step-click="changeStep"
        />
      </div>

      <!-- Form Card -->
      <div class="col-12 col-md-8">
        <!--        {{ lakeHouseData }}-->
        <FormCard
            :header="STEPS[currentStep - 1] as unknown as FormCardHeader"
            :navigation="formNavigation"
            @previous="changeStep"
            @next="changeStep"
            @save="submitForm"
        >
          <!-- FORM BODY -->
          <template #form>
            <q-form ref="formRef">
              <!-- STEP 1: GENERAL INFORMATION -->
              <div v-if="currentStep === 1">
                <LakeHouseBaseInfo v-model="lakeHouseData"/>
              </div>

              <!-- STEP 2: PROVIDER SELECITON -->
              <div v-else-if="currentStep === 2">
                <LakeHouseProviderSelection v-model="lakeHouseData"/>
              </div>

              <!-- STEP 3: CREDENTIALS CONFIG -->
              <div v-else-if="currentStep === 3">
                <LakeHouseCredentials v-model="lakeHouseData"/>
              </div>

              <!-- STEP 4: CREDENTIALS CONFIG -->
              <div v-else-if="currentStep === 4">
                <LakeHousePathConfig v-model="lakeHouseData"/>
              </div>

              <!-- STEP 5: FREQUENCY CONFIG -->
              <div v-else-if="currentStep === 5">
                <LakeHouseFrequency v-model="lakeHouseData"/>
              </div>

              <!-- STEP : REVIEW -->
              <div v-else-if="currentStep === 6">
                <FormReview :sections="previewLakeHouse" @edit-section="changeStep"/>
              </div>
            </q-form>
          </template>
        </FormCard>
      </div>
    </div>
  </q-page>
</template>