<script setup lang="ts">
import type { FormCardHeader } from '@components/cards';
import type { AssetTemplateData, TestResults } from './interfaces';
import type { PageTourStep } from '@composables/tour';

import { ref, computed, onMounted } from 'vue';
import { QForm } from 'quasar';
import { DEFAULT_ASSET_TEMPLATE_DATA, ASSET_TEMPLATE_TOUR_STEPS } from './constants';

import { useStepperNavigation } from '@composables/shared/form';
import { useAddAssetTemplateTranslations } from '@src/composables/i18n/pages/assets/addAssetTemplate/useAddAssetTemplateTranslations';
import { usePageTour } from '@composables/tour';
import { useLogger } from '@composables/useLogger';
import { apis } from '@services/mapex';
import { notifySuccess } from '@utils/alert/notify';
import { handleApiError } from '@utils/error';
import { useRouter, useRoute } from 'vue-router';
import { useOrganizationStore } from '@stores/organization';

defineOptions({
  name: 'CreateEditAssetTemplatePage'
});

/** COMPONENTS */
import { StepperVertical } from '@components/steppers';
import { PageHeader } from '@components/headers';
import { FormCard } from '@components/cards';
import { StandardizedPayloadHelpModal } from '@components/dialogs/standardizedPayloadHelp';

// Step Components
import {
  Step1BasicInfo,
  Step2AssetIdPath,
  Step3PreprocessorScript,
  Step4ValidationScript,
  Step5ConversionScript,
  Step6TestPayload,
  Step7Testing,
  Step8DynamicFields,
  Step9Review,
} from './components';

// i18n
const t = useAddAssetTemplateTranslations();

// Logger
const logger = useLogger('CreateEditAssetTemplatePage');

/**
 * Build tour steps with resolved translations
 * Maps step definitions to PageTourStep with translated text
 * Adds onHighlightStarted callbacks to navigate wizard steps during tour
 *
 * @returns {PageTourStep[]} Tour steps with resolved translations
 */
function buildTourSteps(): PageTourStep[] {
  return ASSET_TEMPLATE_TOUR_STEPS.map((step) => {
    const key = step.translationKey as keyof typeof t.tour;
    const translation = t.tour[key];
    const result: PageTourStep = {
      element: step.element,
      title: translation.title.value,
      description: translation.description.value,
    };
    if (step.side) result.side = step.side;
    if (step.align) result.align = step.align;

    // Add onHighlightStarted callback to navigate wizard steps during tour
    // Extract step number from element selector (e.g., '#step-2' -> 2)
    const stepMatch = step.element.match(/^#step-(\d+)$/);
    if (stepMatch && stepMatch[1]) {
      const stepNumber = parseInt(stepMatch[1], 10);
      result.onHighlightStarted = () => {
        currentStep.value = stepNumber;
      };
    } else if (step.element === '#stepper-section') {
      // When highlighting stepper overview, go to step 1
      result.onHighlightStarted = () => {
        currentStep.value = 1;
      };
    }

    return result;
  });
}

/** PAGE TOUR */
const { startTour } = usePageTour({
  tourId: 'asset-template-wizard',
  steps: buildTourSteps,
  onTourEnd: () => {
    // Reset to step 1 so user can start filling the form
    currentStep.value = 1;
  },
});

// Router & Route
const router = useRouter();
const route = useRoute();

// Organization store
const organizationStore = useOrganizationStore();

// Check if can create templates (only Vendor and Customer)
const canCreateTemplate = computed(() =>
  organizationStore.isVendor || organizationStore.isCustomer
);

/** EDIT MODE DETECTION (MANDATORY) */
const isEditMode = ref(!!route.params.id);
const assetTemplateId = ref(route.params.id as string | undefined);

/** LOADING STATES */
const isLoading = ref(false);  // Loading data (edit mode)
const isSaving = ref(false);   // Saving/updating (replaces isCreating)

// Translated STEPS array
const translatedSteps = computed(() => [
  {
    title: t.steps.step1.label.value,
    icon: 'mdi-information',
    description: t.steps.step1.description.value,
  },
  {
    title: t.steps.step2.label.value,
    icon: 'mdi-routes',
    description: t.steps.step2.description.value,
  },
  {
    title: t.steps.step3.label.value,
    icon: 'mdi-code-braces',
    description: t.steps.step3.description.value,
  },
  {
    title: t.steps.step4.label.value,
    icon: 'mdi-shield-check',
    description: t.steps.step4.description.value,
  },
  {
    title: t.steps.step5.label.value,
    icon: 'mdi-swap-horizontal',
    description: t.steps.step5.description.value,
  },
  {
    title: t.steps.step6.label.value,
    icon: 'mdi-flask',
    description: t.steps.step6.description.value,
  },
  {
    title: t.steps.step7.label.value,
    icon: 'mdi-test-tube',
    description: t.steps.step7.description.value,
  },
  {
    title: t.steps.step8.label.value,
    icon: 'mdi-database-cog',
    description: t.steps.step8.description.value,
  },
  {
    title: t.steps.step9.label.value,
    icon: 'mdi-clipboard-check',
    description: t.steps.step9.description.value,
  },
]);

// Form and stepper state
const step1FormRef = ref<QForm | null>(null);
const step2FormRef = ref<QForm | null>(null);
const currentStep = ref(1);

// Validation error states
const conversionScriptError = ref<string>('');

// Data model
const assetTemplateData = ref<AssetTemplateData>({ ...DEFAULT_ASSET_TEMPLATE_DATA });

// Testing state
const testResults = ref<TestResults>({
  executed: false,
  success: false,
  steps: [],
  output: null,
  logs: [],
});

// Help modal state
const showStandardizedPayloadHelp = ref(false);

// Check if Next button should be disabled
const isNextButtonDisabled = computed(() => {
  // No validation needed for these steps (no forms or editors)
  if (currentStep.value === 3 || currentStep.value === 4 || currentStep.value === 6 || currentStep.value === 8) {
    return false;
  }

  // Step 1: Check if form is valid (will trigger validation check)
  if (currentStep.value === 1) {
    return false; // Let validate() handle it in changeStep
  }

  // Step 2: Check if form is valid
  if (currentStep.value === 2) {
    return false; // Let validate() handle it in changeStep
  }

  // Step 5: Check if conversion script exists
  if (currentStep.value === 5) {
    return !assetTemplateData.value.scriptConversion || assetTemplateData.value.scriptConversion.trim() === '';
  }

  // Step 7: Check if tests passed
  if (currentStep.value === 7) {
    return !testResults.value.executed || !testResults.value.success;
  }

  return false;
});

// Form navigation computed properties
const formNavigation = computed(() => ({
  currentStep: currentStep.value,
  totalSteps: translatedSteps.value.length,
  showPreviousButton: true,
  showNextButton: true,
  showSaveButton: true,
  showCancelButton: true,
  disableNextButton: isNextButtonDisabled.value || isSaving.value,
  disableSaveButton: isNextButtonDisabled.value || isSaving.value,
  loadingSaveButton: isSaving.value,
}));

// Button labels for FormCard navigation
const buttonLabels = computed(() => ({
  previous: t.navigation.previous.value,
  next: t.navigation.next.value,
  save: isEditMode.value ? t.navigation.update.value : t.navigation.save.value,
}));

// Page title (dynamic based on mode)
const pageTitle = computed(() =>
  isEditMode.value ? t.page.titleEdit.value : t.page.title.value
);

// Methods
async function changeStep(step: number) {
  // Clear previous errors
  conversionScriptError.value = '';

  // Validate current step form before proceeding
  if (currentStep.value === 1 && step > 1 && step1FormRef.value) {
    const valid = await step1FormRef.value.validate();
    if (!valid) return;
  }

  if (currentStep.value === 2 && step > 2 && step2FormRef.value) {
    const valid = await step2FormRef.value.validate();
    if (!valid) return;
  }

  // Validate Conversion Script is required before moving past Step 5
  if (currentStep.value === 5 && step > 5) {
    if (!assetTemplateData.value.scriptConversion || assetTemplateData.value.scriptConversion.trim() === '') {
      conversionScriptError.value = t.notifications.conversionScriptRequired.value;
      return;
    }
  }

  // Validate tests passed before moving past Step 7
  if (currentStep.value === 7 && step > 7) {
    if (!testResults.value.executed || !testResults.value.success) {
      return;
    }
  }

  // Step 8 (Dynamic Fields) has no validation - optional step

  currentStep.value = step;
}

// Wrapper for step navigation (non-async)
function handleStepChange(step: number): void {
  void changeStep(step);
}

// Step navigation
useStepperNavigation({
  currentStep,
  totalSteps: translatedSteps.value.length,
  changeStep: handleStepChange,
});

/**
 * Load asset template data for edit mode
 * Fetches data from API and populates form state
 *
 * @returns {Promise<void>}
 */
async function loadAssetTemplateData(): Promise<void> {
  if (!isEditMode.value || !assetTemplateId.value) return;

  isLoading.value = true;
  try {
    const data = await apis.assets.assetTemplate.getById({
      assetTemplateId: assetTemplateId.value
    });

    // Populate form data from API response
    assetTemplateData.value.name = data.name || '';
    assetTemplateData.value.enabled = data.enabled ?? true;
    assetTemplateData.value.description = data.description;

    // Asset Classification
    assetTemplateData.value.categoryId = data.categoryId;
    assetTemplateData.value.categoryName = data.categoryName;
    assetTemplateData.value.manufacturerId = data.manufacturerId;
    assetTemplateData.value.manufacturerName = data.manufacturerName;
    assetTemplateData.value.modelId = data.modelId;
    assetTemplateData.value.modelName = data.modelName;
    assetTemplateData.value.version = data.version;

    assetTemplateData.value.isSystem = data.isSystem ?? false;
    assetTemplateData.value.isTemplate = data.isTemplate ?? false;

    // Scripts
    assetTemplateData.value.assetIdPath = data.assetIdPath || '';
    // Optional scripts - always set (empty if not provided)
    assetTemplateData.value.scriptProcessor = data.scriptProcessor || '';
    assetTemplateData.value.scriptValidator = data.scriptValidator || '';
    // Required script
    assetTemplateData.value.scriptConversion = data.scriptConversion || '';
    if (data.scriptTest) {
      assetTemplateData.value.scriptTest = data.scriptTest;
    }

    // Available Fields
    if (data.availableFields) {
      assetTemplateData.value.availableFields = data.availableFields;
    }

    // Dynamic Fields
    if ((data as any).dynamicFields) {
      assetTemplateData.value.dynamicFields = (data as any).dynamicFields;
    }

    logger.debug('Loaded Asset Template for editing:', {
      scriptTest: data.scriptTest,
      availableFields: data.availableFields,
      fullData: data
    });

    // In EDIT mode, skip to Review step (Step 9) by default
    currentStep.value = 9;

  } catch (error: any) {
    handleApiError(error, {
      defaultMessage: t.notifications.loadFailed,
      timeout: 5000,
    });

    // Navigate back to list on load error
    await router.push('/assets_template');
  } finally {
    isLoading.value = false;
  }
}

/**
 * Submit form - create or update based on mode
 *
 * @returns {Promise<void>}
 */
async function submitForm(): Promise<void> {
  isSaving.value = true;

  try {
    // Prepare data according to AssetTemplateCreate schema (FLAT structure)
    const templateData: any = {
      name: assetTemplateData.value.name,
      enabled: assetTemplateData.value.enabled,

      // Asset Classification (FLAT)
      categoryId: assetTemplateData.value.categoryId,
      categoryName: assetTemplateData.value.categoryName,
      manufacturerId: assetTemplateData.value.manufacturerId,
      manufacturerName: assetTemplateData.value.manufacturerName,
      modelId: assetTemplateData.value.modelId,
      modelName: assetTemplateData.value.modelName,
      version: assetTemplateData.value.version,

      assetIdPath: assetTemplateData.value.assetIdPath,
      scriptConversion: assetTemplateData.value.scriptConversion,
      scriptTest: assetTemplateData.value.scriptTest || '', // MANDATORY for CREATE
      isSystem: assetTemplateData.value.isSystem,
      isTemplate: assetTemplateData.value.isTemplate,

      // Available Fields - ALWAYS include (even if empty array or undefined)
      availableFields: assetTemplateData.value.availableFields || [],

      // Dynamic Fields - field mappings for ClickHouse storage
      dynamicFields: assetTemplateData.value.dynamicFields || [],
    };

    // Add optional fields only if they exist
    if (assetTemplateData.value.description) {
      templateData.description = assetTemplateData.value.description;
    }

    // Optional scripts - always send (empty string clears the field)
    templateData.scriptProcessor = assetTemplateData.value.scriptProcessor || '';
    templateData.scriptValidator = assetTemplateData.value.scriptValidator || '';

    if (isEditMode.value && assetTemplateId.value) {
      // UPDATE existing asset template
      await apis.assets.assetTemplate.update(
        { assetTemplateId: assetTemplateId.value },
        templateData
      );

      notifySuccess({
        message: t.notifications.updated.value,
        timeout: 3000,
      });
    } else {
      // CREATE new asset template
      logger.debug('Creating Asset Template:', templateData);

      const created = await apis.assets.assetTemplate.create(templateData);

      logger.debug('Asset Template created:', created);

      notifySuccess({
        message: t.notifications.created.value,
        timeout: 3000
      });
    }

    // Navigate to asset templates list
    await router.push('/assets_template');

  } catch (error: any) {
    // Use centralized error handler with custom messages for specific status codes
    handleApiError(error, {
      customMessages: {
        409: t.notifications.alreadyExists,
        422: t.notifications.validationFailed,
        network: t.notifications.networkError,
      },
      defaultMessage: isEditMode.value
        ? t.notifications.updateFailed
        : t.notifications.creationFailed,
      timeout: 5000
    });
  } finally {
    isSaving.value = false;
  }
}

/** LIFECYCLE HOOKS */
onMounted(() => {
  void loadAssetTemplateData(); // Load data if edit mode
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
      <!-- Header Section -->
      <div id="page-header-section">
        <PageHeader
          icon="memory"
          iconColor="primary"
          :title="pageTitle"
          :description="t.page.description.value"
          :button="{ label: t.page.button.value, icon: 'arrow_back', flat: true, to: '/assets_template' }"
          :tour="{ enabled: true }"
          @start-tour="startTour"
        />
      </div>

      <!-- Content -->
      <div class="row q-col-gutter-lg">
      <!-- Progress Stepper Vertical -->
      <div id="stepper-section" class="col-12 col-md-4">
        <StepperVertical
          :title="t.stepper.title.value"
          :subtitle="t.stepper.subtitle.value"
          :info-text="t.stepper.requiredInfo.value"
          :current-step-label="t.stepper.currentStep.value"
          :current-step="currentStep"
          :steps="translatedSteps"
          :allow-step-navigation="isEditMode"
          step-id-prefix="step"
          @step-click="changeStep"
        />
      </div>

      <!-- Form Card -->
      <div class="col-12 col-md-8">
        <FormCard
          :header="translatedSteps[currentStep - 1] as unknown as FormCardHeader"
          :navigation="formNavigation"
          :button-labels="buttonLabels"
          @previous="changeStep"
          @next="changeStep"
          @save="submitForm"
        >
          <!-- FORM BODY -->
          <template #form>
            <!-- STEP 1: BASIC INFORMATION -->
            <q-form v-if="currentStep === 1" ref="step1FormRef" greedy>
              <Step1BasicInfo
                v-model="assetTemplateData"
                :can-create-template="canCreateTemplate"
              />
            </q-form>

            <!-- STEP 2: ASSET ID PATH -->
            <q-form v-else-if="currentStep === 2" ref="step2FormRef" greedy>
              <Step2AssetIdPath v-model="assetTemplateData" />
            </q-form>

            <!-- STEP 3: PREPROCESSOR SCRIPT -->
            <Step3PreprocessorScript
              v-else-if="currentStep === 3"
              v-model="assetTemplateData"
            />

            <!-- STEP 4: VALIDATION SCRIPT -->
            <Step4ValidationScript
              v-else-if="currentStep === 4"
              v-model="assetTemplateData"
            />

            <!-- STEP 5: CONVERSION SCRIPT -->
            <Step5ConversionScript
              v-else-if="currentStep === 5"
              v-model="assetTemplateData"
              :error-message="conversionScriptError"
            />

            <!-- STEP 6: TEST PAYLOAD -->
            <Step6TestPayload
              v-else-if="currentStep === 6"
              v-model="assetTemplateData"
            />

            <!-- STEP 7: TESTING -->
            <Step7Testing
              v-else-if="currentStep === 7"
              :model-value="assetTemplateData"
              :test-results="testResults"
              @update:model-value="assetTemplateData = $event"
              @update:test-results="testResults = $event"
              @show-standardized-payload-help="showStandardizedPayloadHelp = true"
            />

            <!-- STEP 8: DYNAMIC FIELDS -->
            <Step8DynamicFields
              v-else-if="currentStep === 8"
              v-model="assetTemplateData"
            />

            <!-- STEP 9: REVIEW -->
            <Step9Review
              v-else-if="currentStep === 9"
              :model-value="assetTemplateData"
              @edit-section="changeStep"
            />
          </template>
        </FormCard>
      </div>
    </div>

      <!-- StandardizedPayload Helper Modal -->
      <StandardizedPayloadHelpModal v-model="showStandardizedPayloadHelp" />
    </div>
  </q-page>
</template>
