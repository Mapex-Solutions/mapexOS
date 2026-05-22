<script setup lang="ts">
defineOptions({
  name: 'CreateEditCustomerPage',
});

/** TYPE IMPORTS */
import type { FormCardHeader } from '@components/cards';
import type { OrganizationResponse } from '@mapexos/schemas';
import type { OrganizationType } from './interfaces';

/** VUE IMPORTS */
import { ref, computed, onMounted } from 'vue';
import { useRouter, useRoute } from 'vue-router';

/** COMPOSABLES */
import { useStepperNavigation } from '@composables/shared/form';
import { useAddCustomerTranslations } from '@composables/i18n';
import { useLogger } from '@composables/useLogger';

/** COMPONENTS */
import { StepperVertical } from '@components/steppers';
import { PageHeader } from '@components/headers';
import { FormCard } from '@components/cards';
import {
  Step1Basic,
  Step2Address,
  Step3AccessPolicy,
  Step4Review,
} from './components';

/** SERVICES */
import { apis } from '@services/mapex';

/** STORES */
import { useOrganizationStore } from '@stores/organization';

/** UTILS */
import { notifyFail } from '@utils/alert';

/** LOCAL IMPORTS */
import {
  INITIAL_ORGANIZATION_FORM_DATA,
  CHILD_TYPE_MAP,
  ORG_TYPE_CONFIG,
  TOTAL_STEPS_WITH_ADDRESS,
  TOTAL_STEPS_WITHOUT_ADDRESS,
} from './constants';
import { useCustomerFormHandlers } from './handlers';

/** COMPOSABLES & STORES */
const t = useAddCustomerTranslations();
const router = useRouter();
const route = useRoute();
const logger = useLogger('CreateEditCustomerPage');
const orgStore = useOrganizationStore();

/** EDIT MODE DETECTION (MANDATORY) */
const isEditMode = ref(!!route.params.id);
const organizationId = ref(route.params.id as string | undefined);

/**
 * Organization type resolved from edit data (only used in edit mode)
 */
const editOrgType = ref<OrganizationType | null>(null);

/** LOADING STATES */
const isLoading = ref(false);
const isSaving = ref(false);

/** STATE */
const step1Ref = ref<InstanceType<typeof Step1Basic> | null>(null);
const step2Ref = ref<InstanceType<typeof Step2Address> | null>(null);
const step3Ref = ref<InstanceType<typeof Step3AccessPolicy> | null>(null);

const currentStep = ref(1);
const formData = ref({ ...INITIAL_ORGANIZATION_FORM_DATA });

/** COMPUTED */

/**
 * Parent organization ID from the organization store (currently selected org)
 */
const parentOrgId = computed(() =>
  orgStore.selectedOrganizationId || undefined,
);

/**
 * Parent organization type resolved from the organization store getter
 */
const parentOrgType = computed((): OrganizationType | undefined => {
  const parent = orgStore.selectedOrganization;
  return parent?.type;
});

/**
 * The organization type being created or edited.
 * - Create mode: derived from parent type via CHILD_TYPE_MAP
 * - Edit mode: from loaded organization data
 * - Fallback: defaults to 'customer' for backward compatibility
 */
const orgType = computed((): OrganizationType => {
  if (isEditMode.value && editOrgType.value) {
    return editOrgType.value;
  }
  if (parentOrgType.value) {
    return (CHILD_TYPE_MAP[parentOrgType.value] as OrganizationType) || 'customer';
  }
  return 'customer';
});

/**
 * Type configuration for the current organization type
 */
const typeConfig = computed(() => ORG_TYPE_CONFIG[orgType.value]);

/**
 * Total number of steps based on whether type has address
 */
const totalSteps = computed(() =>
  typeConfig.value.hasAddress ? TOTAL_STEPS_WITH_ADDRESS : TOTAL_STEPS_WITHOUT_ADDRESS,
);

/**
 * Page title (dynamic based on mode and type)
 */
const pageTitle = computed(() =>
  isEditMode.value ? t.page.titleEdit.value : t.page.title.value,
);

/**
 * Page icon (dynamic based on mode and type)
 */
const pageIcon = computed(() =>
  isEditMode.value ? 'edit' : typeConfig.value.icon,
);

/**
 * Translated steps array - dynamic based on org type
 */
const translatedSteps = computed(() => {
  const steps = [
    {
      title: t.steps.basic.value,
      icon: typeConfig.value.icon,
      description: t.stepDescriptions.basic.value,
    },
  ];

  // Only add address step for types that support it
  if (typeConfig.value.hasAddress) {
    steps.push({
      title: t.steps.address.value,
      icon: 'location_on',
      description: t.stepDescriptions.address.value,
    });
  }

  steps.push(
    {
      title: t.steps.accessPolicy.value,
      icon: 'policy',
      description: t.stepDescriptions.accessPolicy.value,
    },
    {
      title: t.steps.review.value,
      icon: 'check_circle',
      description: t.stepDescriptions.review.value,
    },
  );

  return steps;
});

/**
 * Organization form handlers composable
 */
const handlers = useCustomerFormHandlers({
  formData,
  currentStep,
  isEditMode,
  organizationId,
  isSaving,
  orgType,
  parentOrgId,
  typeConfig,
  step1FormRef: computed(() => step1Ref.value?.formRef ?? null),
  step2FormRef: computed(() => step2Ref.value?.formRef ?? null),
  accessPolicyFormRef: computed(() => step3Ref.value?.formRef ?? null),
  t,
});

/**
 * Form navigation configuration
 */
const formNavigation = computed(() => ({
  currentStep: currentStep.value,
  totalSteps: totalSteps.value,
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
  previous: t.buttons.back.value,
  next: t.buttons.next.value,
  save: isEditMode.value ? t.buttons.updateCustomer.value : t.buttons.createCustomer.value,
}));

/** FUNCTIONS */

/**
 * Load organization data for edit mode
 * Fetches data from API and populates form state
 * Only executes in edit mode with valid ID
 *
 * @returns {Promise<void>}
 */
async function loadOrganizationData(): Promise<void> {
  if (!isEditMode.value || !organizationId.value) return;

  if (!apis.mapexOS?.organizations) {
    notifyFail({ message: t.errors.apiNotInitialized.value });
    return;
  }

  isLoading.value = true;
  try {
    const data: OrganizationResponse = await apis.mapexOS.organizations.getById({
      organizationId: organizationId.value,
    });

    // Set the org type from loaded data
    editOrgType.value = (data.type as OrganizationType) || 'customer';

    // Populate formData from API response
    formData.value = {
      name: data.name || '',
      phone: data.phone || '',
      enabled: data.enabled ?? true,
      address: {
        city: data.address?.city || '',
        state: data.address?.state || '',
        country: data.address?.country || '',
        zipCode: data.address?.zipCode || '',
      },
      authConfig: {
        providerType: 'internal',
        issuerUrl: '',
        clientId: '',
      },
      accessPolicy: {
        rolePolicy: (data.accessPolicy?.rolePolicy as 'strict' | 'merge') || 'strict',
        defaultScope: (data.accessPolicy?.defaultScope as 'local' | 'recursive') || 'local',
      },
    };

    logger.debug('Loaded Organization for editing:', {
      organizationId: organizationId.value,
      orgType: editOrgType.value,
      formData: formData.value,
    });

    // In EDIT mode, skip to Review step by default
    currentStep.value = handlers.reviewStep.value;

  } catch (error: any) {
    logger.error('Failed to load organization:', error);
    notifyFail({ message: t.messages.loadFailed.value });

    // Navigate back to list on load error
    await router.push('/customers');
  } finally {
    isLoading.value = false;
  }
}

/**
 * Update form data with partial updates from step components
 *
 * @param {Partial<typeof formData.value>} partialData - Partial form data to merge
 */
function updateFormData(partialData: Partial<typeof formData.value>): void {
  formData.value = {
    ...formData.value,
    ...partialData,
  };
}

/** COMPOSABLES USAGE */
useStepperNavigation({
  currentStep,
  totalSteps: totalSteps.value,
  changeStep: handlers.handleStepChange,
});

/** LIFECYCLE HOOKS */
onMounted(() => {
  void loadOrganizationData();
});
</script>

<template>
  <q-page class="q-pa-lg">
    <!-- Loading State (Edit Mode) -->
    <div v-if="isLoading" class="flex flex-center q-pa-xl">
      <q-spinner color="primary" size="50px" />
      <span class="q-ml-md text-grey-7">{{ t.messages.loading.value }}</span>
    </div>

    <!-- Form Content -->
    <div v-else>
      <!-- Header Section -->
      <PageHeader
        :icon="pageIcon"
        :icon-color="typeConfig.iconColor"
        :title="pageTitle"
        :description="t.page.description.value"
        :button="{
          label: t.page.backButton.value,
          icon: 'arrow_back',
          flat: true,
          to: '/customers',
        }"
      >
        <template #after-title>
          <q-badge
            :color="typeConfig.iconColor"
            :label="typeConfig.label"
            class="q-ml-sm"
          />
        </template>
      </PageHeader>

      <!-- Content -->
      <div class="row q-col-gutter-lg">
        <!-- Progress Stepper Vertical -->
        <div class="col-12 col-md-4">
          <StepperVertical
            :title="t.sections.progressSteps.value"
            :subtitle="t.messages.completeAllSteps.value"
            :info-text="t.messages.allFieldsRequired.value"
            :current-step-label="t.messages.currentStep.value"
            :current-step="currentStep"
            :steps="translatedSteps"
            :allow-step-navigation="isEditMode"
            @step-click="handlers.changeStep"
          />
        </div>

        <!-- Form Card -->
        <div class="col-12 col-md-8">
          <FormCard
            :header="(translatedSteps[currentStep - 1] as unknown as FormCardHeader)"
            :navigation="formNavigation"
            :button-labels="buttonLabels"
            @previous="handlers.changeStep"
            @next="handlers.changeStep"
            @save="handlers.submitForm"
          >
            <!-- FORM BODY -->
            <template #form>
              <!-- STEP 1: BASIC -->
              <Step1Basic
                v-if="currentStep === 1"
                ref="step1Ref"
                :model-value="formData"
                :has-phone="typeConfig.hasPhone"
                @update:model-value="updateFormData"
              />

              <!-- STEP 2: ADDRESS (only for types with address) -->
              <Step2Address
                v-else-if="typeConfig.hasAddress && currentStep === 2"
                ref="step2Ref"
                :model-value="formData"
                @update:model-value="updateFormData"
              />

              <!-- STEP: ACCESS POLICY (step number varies) -->
              <Step3AccessPolicy
                v-else-if="currentStep === handlers.accessPolicyStep.value"
                ref="step3Ref"
                :model-value="formData"
                @update:model-value="updateFormData"
              />

              <!-- STEP: REVIEW (last step, number varies) -->
              <Step4Review
                v-else-if="currentStep === handlers.reviewStep.value"
                :model-value="formData"
                :is-edit-mode="isEditMode"
                :type-config="typeConfig"
                :org-type="orgType"
                @edit-section="handlers.changeStep"
              />
            </template>
          </FormCard>
        </div>
      </div>
    </div>
  </q-page>
</template>
