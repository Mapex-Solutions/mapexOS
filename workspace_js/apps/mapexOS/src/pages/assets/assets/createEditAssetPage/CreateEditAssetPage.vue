<script setup lang="ts">
defineOptions({
  name: 'CreateEditAssetPage'
});

/** TYPE IMPORTS */
import type { FormCardHeader } from '@components/cards';
import type { AssetFormState } from './interfaces';

/** VUE IMPORTS */
import { ref, computed, watch, onMounted } from 'vue';
import { useRouter, useRoute } from 'vue-router';

/** COMPOSABLES */
import { useStepperNavigation } from '@composables/shared/form';
import { useAddAssetTranslations } from '@src/composables/i18n/pages/assets/addAsset/useAddAssetTranslations';
import { useLogger } from '@composables/useLogger';

/** COMPONENTS */
import { StepperVertical } from '@components/steppers';
import { PageHeader } from '@components/headers';
import { FormCard } from '@components/cards';
import {
  Step1Identification,
  Step2AssetTemplate,
  Step3RouteGroups,
  Step4Connectivity,
  Step5Review,
  HealthMonitoringSection,
} from './components';
import { GenerateCertificateDialog } from './components/GenerateCertificateDialog';

/** LOCAL IMPORTS */
import { INITIAL_ASSET_FORM_DATA, TOTAL_STEPS, STEP } from './constants';
import { useAssetFormHandlers } from './handlers';

/** STORES */

/** SERVICES */
import { apis } from '@services/mapex';

/** UTILS */
import { handleApiError } from '@utils/error';
import { downloadCertZip, decodeBase64ToBytes } from '@utils/zipDownload';
import { notifyWarning } from '@utils/alert/notify';

/** COMPOSABLES & STORES */
const t = useAddAssetTranslations();
const router = useRouter();
const route = useRoute();
const logger = useLogger('CreateEditAssetPage');

/** EDIT MODE DETECTION (MANDATORY) */
const isEditMode = ref(!!route.params.id);
const assetId = ref(route.params.id as string | undefined);

/** LOADING STATES */
const isLoading = ref(false);  // Loading data in EDIT mode
const isSaving = ref(false);   // Saving/updating data (replaces isCreating)

/** STATE */
const step1Ref = ref<InstanceType<typeof Step1Identification> | null>(null);
const step2Ref = ref<InstanceType<typeof Step2AssetTemplate> | null>(null);
const step3Ref = ref<InstanceType<typeof Step3RouteGroups> | null>(null);
const step4Ref = ref<InstanceType<typeof Step4Connectivity> | null>(null);

const currentStep = ref(1);
const assetData = ref({ ...INITIAL_ASSET_FORM_DATA });
const hasCurrentCert = ref(false);
const formState = ref<AssetFormState>({
  selectedTemplate: null,
  selectedRouteGroups: [],
  isCreating: false,
  currentStep: currentStep.value,
});

/**
 * Post-save cert dialog state. Opened automatically after creating an
 * MQTT asset so the operator can issue the device certificate and
 * download the private key (which the platform discards immediately
 * after returning it). Skipping is allowed — the list page will then
 * surface a warning chip next to the protocol cell.
 */
const showCertDialog = ref(false);
const issuingCert = ref(false);
const justCreatedUUID = ref('');
const justCreatedName = ref('');

/** FUNCTIONS */

/**
 * Load asset data for edit mode
 * Fetches data from API and populates form state
 * Only executes in edit mode with valid ID
 *
 * @returns {Promise<void>}
 */
async function loadAssetData(): Promise<void> {
  if (!isEditMode.value || !assetId.value) return;

  isLoading.value = true;
  try {
    const data = await apis.assets.asset.getById({
      assetId: assetId.value
    });

    // Populate assetData from API response
    assetData.value.name = data.name || '';
    assetData.value.assetId = data.assetUUID || '';
    assetData.value.enabled = data.enabled ?? true;
    assetData.value.debugEnabled = data.debugEnabled ?? false;
    assetData.value.description = data.description || '';
    assetData.value.assetTemplateId = data.assetTemplateId || null;
    assetData.value.routeGroupIds = data.routeGroupIds || [];
    assetData.value.protocol = data.protocol?.type?.toUpperCase() || 'HTTP';
    assetData.value.latitude = data.latitude ?? null;
    assetData.value.longitude = data.longitude ?? null;

    // Populate MQTT config if protocol is MQTT. Password is intentionally
    // left blank on edit — the existing bcrypt hash stays on the asset
    // unless the operator types a new password (see Step4Connectivity
    // password rules + useAssetFormHandlers.buildAssetPayload).
    if (data.protocol?.mqtt) {
      // certTTL is hydrated only when the asset carries it so the
      // optional-field contract is preserved (legacy assets without
      // certTTL fall back to INITIAL_MQTT_CONFIG's default at form
      // mount). exactOptionalPropertyTypes forbids assigning
      // `undefined` to an optional field, hence the conditional spread.
      assetData.value.mqttConfig = {
        clientId: data.protocol.mqtt.clientId || '',
        username: data.protocol.mqtt.username || '',
        authType: data.protocol.mqtt.authType || 'cert',
        password: '',
        ...(data.protocol.mqtt.certTTL
          ? { certTTL: { value: data.protocol.mqtt.certTTL.value, unit: data.protocol.mqtt.certTTL.unit } }
          : {}),
      };
    }

    // Track whether the asset already has an active cert so Step4
    // opens on the right auth-mode radio (cert vs password). A
    // non-empty serial is the signal — empty/absent means no cert.
    hasCurrentCert.value = !!data.currentCert?.serial;

    // Populate health monitoring config
    if (data.healthMonitor) {
      assetData.value.healthMonitor = {
        enabled: data.healthMonitor.enabled ?? false,
        thresholdMinutes: data.healthMonitor.thresholdMinutes ?? 10,
        requiredMisses: data.healthMonitor.requiredMisses ?? 3,
        heartbeatMode: data.healthMonitor.heartbeatMode ?? 'implicit',
        offlineRouteGroupIds: data.healthMonitor.offlineRouteGroupIds || [],
        onlineRouteGroupIds: data.healthMonitor.onlineRouteGroupIds || [],
        selectedOfflineRouteGroups: [],
        selectedOnlineRouteGroups: [],
      };

      // Fetch offline route groups for display
      if (data.healthMonitor.offlineRouteGroupIds?.length) {
        try {
          const offlineRgs = await Promise.all(
            data.healthMonitor.offlineRouteGroupIds.map((id: string) =>
              apis.router.routegroup.getById({ routeGroupId: id })
            )
          );
          assetData.value.healthMonitor.selectedOfflineRouteGroups = offlineRgs;
        } catch (error: any) {
          logger.error('Failed to load offline route groups:', error);
        }
      }

      // Fetch online route groups for display
      if (data.healthMonitor.onlineRouteGroupIds?.length) {
        try {
          const onlineRgs = await Promise.all(
            data.healthMonitor.onlineRouteGroupIds.map((id: string) =>
              apis.router.routegroup.getById({ routeGroupId: id })
            )
          );
          assetData.value.healthMonitor.selectedOnlineRouteGroups = onlineRgs;
        } catch (error: any) {
          logger.error('Failed to load online route groups:', error);
        }
      }
    }

    // Fetch and populate selected template if assetTemplateId exists
    if (data.assetTemplateId) {
      try {
        const templateData = await apis.assets.assetTemplate.getById({
          assetTemplateId: data.assetTemplateId
        });
        assetData.value.selectedTemplate = templateData;
        // Also update formState for Step5Review to display correctly
        formState.value.selectedTemplate = templateData;
      } catch (error: any) {
        logger.error('Failed to load asset template:', error);
        // Continue even if template fetch fails
      }
    }

    // Fetch and populate selected route groups if routeGroupIds exist
    if (data.routeGroupIds && data.routeGroupIds.length > 0) {
      try {
        const routeGroupsPromises = data.routeGroupIds.map((routeGroupId: string) =>
          apis.router.routegroup.getById({ routeGroupId })
        );
        const routeGroupsData = await Promise.all(routeGroupsPromises);
        assetData.value.selectedRouteGroups = routeGroupsData;
        // Also update formState for Step5Review to display correctly
        formState.value.selectedRouteGroups = routeGroupsData;
      } catch (error: any) {
        logger.error('Failed to load route groups:', error);
        // Continue even if route groups fetch fails
      }
    }

    logger.debug('Loaded Asset for editing:', {
      assetId: assetId.value,
      assetData: assetData.value,
      fullData: data
    });

    // In EDIT mode, skip to Review step by default
    currentStep.value = STEP.REVIEW;

  } catch (error: any) {
    handleApiError(error, {
      defaultMessage: t.notifications.loadFailed.value,
      timeout: 5000,
    });

    // Navigate back to list on load error
    await router.push('/assets');
  } finally {
    isLoading.value = false;
  }
}

/**
 * Update asset data with partial updates from step components
 * Merges partial data into existing asset data
 * @param {Partial<typeof assetData.value>} partialData - Partial asset data to merge
 * @returns {void}
 */
function updateAssetData(partialData: Partial<typeof assetData.value>): void {
  assetData.value = {
    ...assetData.value,
    ...partialData,
  };
}

/**
 * Save flow orchestrator — wraps the handler's submitForm result.
 * For MQTT creates, defers navigation and opens the certificate
 * dialog so the operator can either issue (and download) or skip.
 * For every other case (HTTP creates, any edit), navigates straight
 * back to the list.
 */
async function onSubmit(): Promise<void> {
  const result = await handlers.submitForm();
  if (!result.ok) return;

  // The cert dialog is only meaningful when the asset is being
  // created in cert mode — password-mode assets have nothing to
  // generate at this point, and edits already render the cert
  // section inside the details drawer.
  if (result.isCreate && result.isMqtt && result.isCert && result.assetUUID) {
    justCreatedUUID.value = result.assetUUID;
    justCreatedName.value = result.assetName ?? '';
    showCertDialog.value = true;
    return;
  }

  await router.push('/assets');
}

/**
 * Cert dialog `@issued` handler — issue + download in one shot.
 * The PEMs are base64 JSON bytes returned ONLY at issue time; the
 * server discards them immediately, so this download is the only
 * chance to capture them. On failure the operator stays on the
 * dialog (close button still available) and can retry.
 */
async function onCertIssued(): Promise<void> {
  if (!justCreatedUUID.value || issuingCert.value) return;
  issuingCert.value = true;
  try {
    const res = await apis.assets.mqttcerts.issueCert(justCreatedUUID.value);
    await downloadCertZip({
      filename: `${justCreatedName.value || justCreatedUUID.value}-mqtt-cert.zip`,
      certPEM: decodeBase64ToBytes(res.certPEM),
      keyPEM: decodeBase64ToBytes(res.keyPEM),
      caChainPEM: decodeBase64ToBytes(res.caChainPEM),
    });
    showCertDialog.value = false;
    await router.push('/assets');
  } catch (error: any) {
    handleApiError(error, {
      defaultMessage: t.notifications.certIssueFailed.value,
      timeout: 5000,
    });
  } finally {
    issuingCert.value = false;
  }
}

/**
 * Cert dialog close handler — the operator chose to skip cert
 * issuance for now. Warn that the asset will not be able to
 * connect via mTLS until a cert is generated (the assets list
 * surfaces this through a warning chip next to the protocol).
 */
async function onCertSkipped(): Promise<void> {
  showCertDialog.value = false;
  notifyWarning({
    message: t.notifications.assetCreatedWithoutCert.value,
    timeout: 5000,
  });
  await router.push('/assets');
}

/** COMPUTED */

/**
 * Page title (dynamic based on mode)
 */
const pageTitle = computed(() =>
  isEditMode.value ? t.page.titleEdit.value : t.page.title.value
);

/**
 * Translated steps array with reactive translations
 */
const translatedSteps = computed(() => [
  {
    title: t.steps.step1.label.value,
    icon: 'mdi-fingerprint',
    description: t.steps.step1.description.value,
  },
  {
    title: t.steps.step2.label.value,
    icon: 'mdi-file-document',
    description: t.steps.step2.description.value,
  },
  {
    title: t.steps.step3.label.value,
    icon: 'mdi-routes',
    description: t.steps.step3.description.value,
  },
  {
    title: t.steps.step4.label.value,
    icon: 'mdi-wifi',
    description: t.steps.step4.description.value,
  },
  {
    title: t.steps.step5.label.value,
    icon: 'mdi-heart-pulse',
    description: t.steps.step5.description.value,
  },
  {
    title: t.steps.step6.label.value,
    icon: 'mdi-clipboard-check',
    description: t.steps.step6.description.value,
  },
]);

/**
 * Asset form handlers composable
 */
const handlers = useAssetFormHandlers({
  assetData,
  formState,
  currentStep,
  isEditMode,
  assetId,
  isSaving,
  step1FormRef: computed(() => step1Ref.value?.formRef ?? null),
  step2FormRef: computed(() => step2Ref.value?.formRef ?? null),
  step3FormRef: computed(() => step3Ref.value?.formRef ?? null),
  step4FormRef: computed(() => step4Ref.value?.formRef ?? null),
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
  disableNextButton: handlers.isNextButtonDisabled.value || isSaving.value,
  disableSaveButton: handlers.isNextButtonDisabled.value || isSaving.value,
  loadingSaveButton: isSaving.value,
}));

/**
 * Button labels with reactive translations
 */
const buttonLabels = computed(() => ({
  previous: t.navigation.previous.value,
  next: t.navigation.next.value,
  save: isEditMode.value ? t.navigation.update.value : t.navigation.save.value,
}));

/**
 * Labels for the post-save certificate dialog — passed as a single
 * resolved-strings object so the dialog component itself stays free
 * of any specific i18n composable wiring (it can be reused from the
 * drawer and the list page later).
 */
const certDialogLabels = computed(() => ({
  title: t.steps.step4.certDialog.title.value,
  warning: t.steps.step4.certDialog.warning.value,
  replaceWarning: t.steps.step4.certDialog.replaceWarning.value,
  generateButton: t.steps.step4.certDialog.generateButton.value,
  skipButton: t.steps.step4.certDialog.skipButton.value,
}));

/** COMPOSABLES USAGE */
useStepperNavigation({
  currentStep,
  totalSteps: TOTAL_STEPS,
  changeStep: handlers.handleStepChange,
});

/** WATCHERS */
watch(() => assetData.value.assetTemplateId, (newValue, oldValue) => {
  logger.debug('assetTemplateId changed:', {
    old: oldValue,
    new: newValue,
  });
});

/** LIFECYCLE HOOKS */
onMounted(() => {
  void loadAssetData(); // Loads data only if edit mode
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
      <PageHeader
        icon="devices"
        iconColor="primary"
        :title="pageTitle"
        :description="t.page.description.value"
        :button="{ label: t.page.button.value, icon: 'arrow_back', flat: true, to: '/assets' }"
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
            @save="onSubmit"
          >
            <!-- FORM BODY -->
            <template #form>
              <!-- STEP 1: IDENTIFICATION -->
              <Step1Identification
                v-if="currentStep === STEP.IDENTIFICATION"
                ref="step1Ref"
                :model-value="assetData"
                @update:model-value="updateAssetData"
              />

              <!-- STEP 2: ASSET TEMPLATE -->
              <Step2AssetTemplate
                v-else-if="currentStep === STEP.ASSET_TEMPLATE"
                ref="step2Ref"
                :model-value="assetData"
                @update:model-value="updateAssetData"
                @template-selected="handlers.onTemplateSelected"
              />

              <!-- STEP 3: ROUTE GROUPS -->
              <Step3RouteGroups
                v-else-if="currentStep === STEP.ROUTE_GROUPS"
                ref="step3Ref"
                :model-value="assetData"
                @update:model-value="updateAssetData"
                @route-groups-selected="handlers.onRouteGroupsSelected"
              />

              <!-- STEP 4: CONNECTIVITY -->
              <Step4Connectivity
                v-else-if="currentStep === STEP.CONNECTIVITY"
                ref="step4Ref"
                :model-value="assetData"
                :is-edit-mode="isEditMode"
                @update:model-value="updateAssetData"
              />

              <!-- STEP 5: HEALTH MONITORING -->
              <HealthMonitoringSection
                v-else-if="currentStep === STEP.HEALTH_MONITORING"
                :model-value="assetData.healthMonitor"
                :asset-u-u-i-d="assetData.assetId"
                :protocol="assetData.protocol"
                @update:model-value="(val) => updateAssetData({ healthMonitor: { ...assetData.healthMonitor, ...val } })"
              />

              <!-- STEP 6: REVIEW -->
              <Step5Review
                v-else-if="currentStep === STEP.REVIEW"
                :model-value="assetData"
                :form-state="formState"
                @edit-section="handlers.changeStep"
              />
            </template>
          </FormCard>
        </div>
      </div>
    </div>

    <!-- Post-save certificate dialog. Opens automatically after an
         MQTT asset is created so the operator can issue and download
         the device certificate in one click; skipping is allowed and
         flagged on the assets list. -->
    <GenerateCertificateDialog
      v-model:show="showCertDialog"
      :asset-uuid="justCreatedUUID"
      :asset-name="justCreatedName"
      :has-existing-cert="false"
      :labels="certDialogLabels"
      @issued="onCertIssued"
      @update:show="(v) => { if (!v && !issuingCert) void onCertSkipped(); }"
    />

  </q-page>
</template>
