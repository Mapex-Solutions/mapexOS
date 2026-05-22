<script setup lang="ts">
defineOptions({
  name: 'CreateEditHttpDataSourcePage'
});

/** TYPE IMPORTS */
import type { HttpDataSource } from './interfaces/httpDataSource.interface';
import type { FormCardHeader } from '@components/cards';
import type { QForm } from 'quasar';
import type { DataSourceCreate } from '@mapexos/schemas';

/** VUE IMPORTS */
import { ref, computed, onMounted } from 'vue';
import { useRouter, useRoute } from 'vue-router';

/** COMPOSABLES */
import { useStepperNavigation } from '@composables/shared/form';
import { useHttpDataSourceCreateEditTranslations } from '@composables/i18n/pages/datasources/http';
import { useLogger } from '@composables/useLogger';

/** COMPONENTS */
import { StepperVertical } from '@components/steppers';
import { PageHeader } from '@components/headers';
import { FormCard } from '@components/cards';
import {
  Step1BasicInfo,
  Step3WorkingHours,
  Step4Authentication,
  Step5AssetBinding,
  Step6Review,
} from './components';

/** LOCAL IMPORTS */
import { HTTP_DATASOURCE_DEFAULTS, INITIAL_STEP, TOTAL_STEPS, STEP } from './constants';
import { useHttpDataSourceHandlers } from './handlers/useHttpDataSourceHandlers';
import { notifySuccess } from '@utils/alert/notify';
import { handleApiError } from '@utils/error';

/** API */
import { apis } from '@services/mapex';

/** COMPOSABLES & STORES */
const route = useRoute();
const router = useRouter();
const t = useHttpDataSourceCreateEditTranslations();
const logger = useLogger('CreateEditHttpDataSourcePage');

/** EDIT MODE DETECTION */
const isEditMode = ref(!!route.params.id);
const dataSourceId = ref(route.params.id as string | undefined);

/** STATE */
const isLoading = ref(false);  // Loading data in EDIT mode
const isSaving = ref(false);   // Saving/updating data
const step1FormRef = ref<QForm | null>(null);
const step2FormRef = ref<QForm | null>(null);
const step3FormRef = ref<QForm | null>(null);
const step4FormRef = ref<QForm | null>(null);
const currentStep = ref(INITIAL_STEP);

const dataSource = ref<HttpDataSource>({ ...HTTP_DATASOURCE_DEFAULTS });

/** FUNCTIONS */

/**
 * Validate and change to a specific step
 * Performs validation on current step before allowing navigation
 * @param {number} step - Target step number to navigate to
 * @returns {Promise<void>}
 */
async function changeStep(step: number): Promise<void> {
  // Validate current step form before proceeding
  if (currentStep.value === 1 && step > 1 && step1FormRef.value) {
    const valid = await step1FormRef.value.validate();
    if (!valid) return;
  }

  if (currentStep.value === 2 && step > 2 && step2FormRef.value) {
    const valid = await step2FormRef.value.validate();
    if (!valid) return;
  }

  if (currentStep.value === 3 && step > 3 && step3FormRef.value) {
    const valid = await step3FormRef.value.validate();
    if (!valid) return;
  }

  if (currentStep.value === 4 && step > 4 && step4FormRef.value) {
    const valid = await step4FormRef.value.validate();
    if (!valid) return;
  }

  currentStep.value = step;
}

/**
 * Non-async wrapper for step navigation
 * Required for compatibility with components that don't support async handlers
 * @param {number} step - Target step number to navigate to
 * @returns {void}
 */
function handleStepChange(step: number): void {
  void changeStep(step);
}

/**
 * Load HTTP data source data for edit mode
 * Fetches data from API and populates form state
 * Only executes in edit mode with valid ID
 *
 * @returns {Promise<void>}
 */
async function loadHttpDataSourceData(): Promise<void> {
  if (!isEditMode.value || !dataSourceId.value) return;

  isLoading.value = true;
  try {
    const data = await apis.httpGateway.datasource.getById({
      dataSourceId: dataSourceId.value
    });

    // Populate basic information
    dataSource.value.name = data.name || '';
    dataSource.value.enabled = data.enabled ?? true;
    dataSource.value.description = data.description || '';
    dataSource.value.mode = data.mode?.toUpperCase() || 'PUSH';
    dataSource.value.protocol = data.protocol?.toUpperCase() || 'HTTP';

    // Populate authentication
    dataSource.value.authType = data.auth?.type || 'none';
    if (data.auth?.type === 'apiKey' && data.auth.apiKey) {
      dataSource.value.apiKey = {
        headerApiKey: data.auth.apiKey.fieldName || '',
        valueApiKey: data.auth.apiKey.key || '',
      };
    } else if (data.auth?.type === 'jwt' && data.auth.jwt) {
      dataSource.value.jwt = {
        secretKey: data.auth.jwt.secret || '',
        headerName: data.auth.jwt.headerName || 'Authorization',
      };
    } else if (data.auth?.type === 'ip_whitelist' && data.auth.ipWhitelist) {
      dataSource.value.ipWhitelist = {
        addresses: data.auth.ipWhitelist.cidrs || [],
      };
    } else if (data.auth?.type === 'oauth2' && data.auth.oauth2) {
      dataSource.value.oauth2 = {
        jwksUrl: data.auth.oauth2.jwksURL || '',
      };
    }

    // Populate working hours
    if (data.workingHours) {
      dataSource.value.enableWorkingHours = data.workingHours.enabled || false;
      dataSource.value.daysOfWeek = (data.workingHours.days || []).map(String);
      dataSource.value.timeIntervals = [{
        startTime: data.workingHours.startAt || '09:00',
        endTime: data.workingHours.endAt || '17:00',
      }];
      dataSource.value.timezone = data.workingHours.timeZone || 'UTC';
    }

    // Populate rate limit
    if (data.rateLimit) {
      dataSource.value.enableRateLimit = true;
      dataSource.value.rateLimitType = data.rateLimit.type || null;
      dataSource.value.rateLimitValue = data.rateLimit.value || 0;
      dataSource.value.burstCapacity = data.rateLimit.burstCapacity || 0;
      dataSource.value.actionOnExceed = data.rateLimit.actionOnExceed || null;
    }

    // Populate asset binding
    if (data.assetBind) {
      dataSource.value.bindingMode = data.assetBind.type || null;
      if (data.assetBind.type === 'fixedAssetId') {
        // assetId = MongoDB ID (for display/edit)
        dataSource.value.directAssetId = data.assetBind.data?.assetId || null;
        // uuidField[0] = path to extract UUID from payload
        dataSource.value.directAssetIdPath = data.assetBind.data?.uuidField?.[0] || null;
      } else if (data.assetBind.type === 'uuidField' && data.assetBind.data?.uuidField) {
        dataSource.value.finalUuidPaths = data.assetBind.data.uuidField;
        dataSource.value.customUuidPaths = (data.assetBind.data.uuidField || []).map(path => ({ path }));
      }
    }

    // In EDIT mode, navigate to Review step after loading data
    currentStep.value = STEP.REVIEW;

  } catch (error: any) {
    handleApiError(error, {
      defaultMessage: t.notifications.loadFailed.value,
      timeout: 5000,
    });

    // Navigate back to list on load error
    await router.push('/data_sources/http');
  } finally {
    isLoading.value = false;
  }
}

/**
 * Transform UI data structure to backend contract (DataSourceCreate)
 * Converts UI field names to match backend expectations
 * @param {HttpDataSource} uiData - UI form data
 * @returns {DataSourceCreate} Backend-compliant payload
 */
function transformToBackendPayload(uiData: HttpDataSource): DataSourceCreate {
  const payload: DataSourceCreate = {
    name: uiData.name,
    enabled: uiData.enabled,
    description: uiData.description || undefined,
    mode: uiData.mode?.toLowerCase() as 'push' | 'pull' | 'X',
    protocol: uiData.protocol?.toLowerCase() as 'http' | 'mqtt',

    // Auth transformation
    auth: {
      type: uiData.authType as 'apiKey' | 'jwt' | 'ip_whitelist' | 'oauth2' | 'none',
      ...(uiData.authType === 'apiKey' && {
        apiKey: {
          type: 'header',
          fieldName: uiData.apiKey.headerApiKey,
          key: uiData.apiKey.valueApiKey,
        }
      }),
      ...(uiData.authType === 'jwt' && {
        jwt: {
          secret: uiData.jwt.secretKey,
          algorithms: 'HS256' as const,
          headerName: uiData.jwt.headerName || 'Authorization',
        }
      }),
      ...(uiData.authType === 'ip_whitelist' && {
        ipWhitelist: {
          cidrs: uiData.ipWhitelist.addresses,
        }
      }),
      ...(uiData.authType === 'oauth2' && {
        oauth2: {
          jwksURL: uiData.oauth2.jwksUrl,
        }
      }),
      ...(uiData.authType === 'none' && {
        none: {}
      }),
    },

    // Asset Bind transformation
    assetBind: {
      type: uiData.bindingMode as 'fixedAssetId' | 'uuidField',
      data: {
        ...(uiData.bindingMode === 'fixedAssetId' && {
          assetId: uiData.directAssetId || '',
          // uuidField contains the path to extract UUID from payload (js-executor uses uuidField[0])
          // Only include if path exists (schema requires min 1 item when present)
          ...(uiData.directAssetIdPath && { uuidField: [uiData.directAssetIdPath] }),
        }),
        ...(uiData.bindingMode === 'uuidField' && {
          uuidField: uiData.finalUuidPaths || [],
        }),
      }
    },
  };

  // Add working hours if enabled
  if (uiData.enableWorkingHours && uiData.timeIntervals.length > 0) {
    const firstInterval = uiData.timeIntervals[0];
    payload.workingHours = {
      enabled: true,
      days: uiData.daysOfWeek.map(day => parseInt(day)),
      startAt: firstInterval?.startTime || '09:00',
      endAt: firstInterval?.endTime || '17:00',
      timeZone: uiData.timezone,
    };
  }

  // Add rate limit if enabled
  if (uiData.enableRateLimit && uiData.rateLimitType) {
    payload.rateLimit = {
      type: uiData.rateLimitType as 'second' | 'minute' | 'hour',
      value: uiData.rateLimitValue,
      burstCapacity: uiData.burstCapacity,
      actionOnExceed: uiData.actionOnExceed as 'drop' | 'queue',
    };
  }

  return payload;
}

/**
 * Submit the HTTP data source form - create or update based on mode
 * Handles both CREATE and EDIT operations in single function
 * @returns {Promise<void>}
 */
async function submitForm(): Promise<void> {
  isSaving.value = true;

  try {
    // Transform UI data to backend contract
    const payload = transformToBackendPayload(dataSource.value);

    logger.debug('Transformed payload:', payload);

    if (isEditMode.value && dataSourceId.value) {
      // UPDATE existing data source
      logger.debug('Updating HTTP Data Source:', payload);

      await apis.httpGateway.datasource.update(
        { dataSourceId: dataSourceId.value },
        payload
      );

      logger.debug('Data Source updated successfully');

      notifySuccess({
        message: t.notifications.updateSuccess.value,
        timeout: 3000
      });
    } else {
      // CREATE new data source
      logger.debug('Creating HTTP Data Source:', payload);

      const created = await apis.httpGateway.datasource.create(payload);

      logger.debug('Data Source created:', created);

      notifySuccess({
        message: t.notifications.createSuccess.value,
        timeout: 3000
      });
    }

    // Navigate to HTTP data sources list after success
    await router.push('/data_sources/http');

  } catch (error: any) {
    handleApiError(error, {
      defaultMessage: isEditMode.value
        ? t.notifications.updateFailed.value
        : t.notifications.createFailed.value,
      timeout: 5000
    });
  } finally {
    isSaving.value = false;
  }
}

/**
 * Handle asset binding updates from Step5AssetBinding component
 * Merges partial updates and triggers mapping test if requested
 * @param {Partial<HttpDataSource>} value - Partial data source update
 * @returns {void}
 */
function handleAssetBindingUpdate(value: Partial<HttpDataSource>): void {
  dataSource.value = { ...dataSource.value, ...value };
  if ((value as any).testMapping) {
    void testMapping();
  }
}

/** COMPUTED */

/**
 * Page title changes based on mode
 */
const pageTitle = computed(() =>
  isEditMode.value ? t.page.titleEdit.value : t.page.title.value
);

/**
 * Page description changes based on mode
 */
const pageDescription = computed(() =>
  isEditMode.value ? t.page.descriptionEdit.value : t.page.description.value
);

/**
 * Translated steps array with reactive configuration
 */
const translatedSteps = computed(() => t.steps.value);

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
  disableNextButton: isSaving.value,
  disableSaveButton: isSaving.value,
  loadingSaveButton: isSaving.value,
}));

/**
 * Button labels for FormCard navigation - changes based on mode
 */
const buttonLabels = computed(() => ({
  previous: t.navigation.back.value,
  next: t.navigation.next.value,
  save: isEditMode.value ? t.navigation.update.value : t.navigation.create.value,
}));

/** COMPOSABLES USAGE */

const { testMapping } = useHttpDataSourceHandlers(dataSource, t);

useStepperNavigation({
  currentStep,
  totalSteps: TOTAL_STEPS,
  changeStep: handleStepChange,
});

/** LIFECYCLE HOOKS */
onMounted(() => {
  void loadHttpDataSourceData(); // Loads data only if edit mode
});

/** WATCHERS */
// No watchers currently needed
</script>

<template>
  <q-page class="q-pa-lg">
    <!-- Loading State (Edit Mode Only) -->
    <div v-if="isLoading" class="flex flex-center q-pa-xl" style="min-height: 400px;">
      <div class="column items-center">
        <q-spinner color="primary" size="50px" />
        <span class="q-mt-md text-grey-7">{{ t.page.loading.value }}</span>
      </div>
    </div>

    <!-- Form Content (Both Modes) -->
    <div v-else>
      <!-- Header Section -->
      <PageHeader
        icon="source"
        iconColor="primary"
        :title="pageTitle"
        :description="pageDescription"
        :button="{ label: t.page.backButton.value, icon: 'arrow_back', flat: true, to: '/data_sources/http' }"
      />

      <!-- Content -->
      <div class="row q-col-gutter-lg">
      <!-- Progress Stepper Vertical -->
      <div class="col-12 col-md-4">
        <StepperVertical
          :title="t.stepper.title.value"
          :subtitle="t.stepper.subtitle.value"
          :info-text="t.stepper.infoText.value"
          :current-step-label="t.stepper.currentStepLabel.value"
          :current-step="currentStep"
          :steps="translatedSteps"
          :allow-step-navigation="isEditMode"
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
            <!-- STEP 1: Basic Information -->
            <div v-if="currentStep === STEP.BASIC_INFO">
              <q-form ref="step1FormRef" greedy>
                <Step1BasicInfo v-model="dataSource" />
              </q-form>
            </div>

            <!-- STEP 2: Working Hours & Rate Limit -->
            <div v-if="currentStep === STEP.WORKING_HOURS">
              <q-form ref="step2FormRef" greedy>
                <Step3WorkingHours v-model="dataSource" />
              </q-form>
            </div>

            <!-- STEP 3: Authentication -->
            <div v-if="currentStep === STEP.AUTHENTICATION">
              <q-form ref="step3FormRef" greedy>
                <Step4Authentication v-model="dataSource" />
              </q-form>
            </div>

            <!-- STEP 4: Asset Binding -->
            <div v-if="currentStep === STEP.ASSET_BINDING">
              <q-form ref="step4FormRef" greedy>
                <Step5AssetBinding v-model="dataSource" @update:model-value="handleAssetBindingUpdate" />
              </q-form>
            </div>

            <!-- STEP 5: Review -->
            <div v-if="currentStep === STEP.REVIEW">
              <Step6Review :data-source="dataSource" @edit-section="handleStepChange" />
            </div>
          </template>
        </FormCard>
      </div>
    </div>
    </div>
  </q-page>
</template>

<style scoped>
.rounded-borders {
  border-radius: var(--mapex-radius-md);
}

.form-card {
  height: 100%;
}
</style>
