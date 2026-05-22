<script setup lang="ts">
defineOptions({
  name: 'CreateEditUserPage',
});

/** TYPE IMPORTS */
import type { FormCardHeader } from '@components/cards';
import type { UserResponse } from '@mapexos/schemas';
import type { PageTourStep } from '@composables/tour';
import type { UserFormData } from './interfaces';

/** VUE IMPORTS */
import { ref, computed, onMounted } from 'vue';
import { useRouter, useRoute } from 'vue-router';

/** COMPOSABLES */
import { useStepperNavigation } from '@composables/shared/form';
import { useAddUserTranslations } from '@composables/i18n';
import { usePageTour } from '@composables/tour';
import { useLogger } from '@composables/useLogger';

/** COMPONENTS */
import { StepperVertical } from '@components/steppers';
import { PageHeader } from '@components/headers';
import { FormCard } from '@components/cards';
import {
  Step1Personal,
  Step2Security,
  Step3Access,
  Step4Review,
} from './components';

/** SERVICES */
import { apis } from '@services/mapex';

/** UTILS */
import { notifyFail } from '@utils/alert';

/** LOCAL IMPORTS */
import {
  INITIAL_USER_FORM_DATA,
  DEMO_USER_FORM_DATA,
  CREATE_USER_TOUR_STEPS,
  TOTAL_STEPS,
  STEP,
} from './constants';
import { useUserFormHandlers } from './handlers';

/** COMPOSABLES & STORES */
const t = useAddUserTranslations();
const router = useRouter();
const route = useRoute();
const logger = useLogger('CreateEditUserPage');

/** EDIT MODE DETECTION (MANDATORY) */
const isEditMode = ref(!!route.params.id);
const userId = ref(route.params.id as string | undefined);

/** TOUR MODE DETECTION */
const isTourDemoMode = computed(() => route.query.tour === 'true');

/** LOADING STATES */
const isLoading = ref(false);
const isSaving = ref(false);

/** STATE */
const step1Ref = ref<InstanceType<typeof Step1Personal> | null>(null);
const step2Ref = ref<InstanceType<typeof Step2Security> | null>(null);
const step3Ref = ref<InstanceType<typeof Step3Access> | null>(null);

const currentStep = ref(1);
const userData = ref(isTourDemoMode.value ? { ...DEMO_USER_FORM_DATA } : { ...INITIAL_USER_FORM_DATA });

/** FUNCTIONS */

/**
 * Load user data for edit mode
 * Fetches user data and current access configuration from API
 * Only executes in edit mode with valid ID
 *
 * V1: AuthProvider removed - always internal auth.
 * Backend defaults to internal when not provided.
 *
 * @returns {Promise<void>}
 */
async function loadUserData(): Promise<void> {
  if (isTourDemoMode.value) return; // Tour mode: use pre-filled demo data
  if (!isEditMode.value || !userId.value) return;

  if (!apis.mapexOS?.users) {
    notifyFail({ message: t.errors.apiNotInitialized.value });
    return;
  }

  isLoading.value = true;
  try {
    // Load user data
    const data: UserResponse = await apis.mapexOS.users.getById({ userId: userId.value });

    // Initialize form data with user info
    // V1: AuthProvider not loaded - always internal auth
    const formData: UserFormData = {
      email: data.email || '',
      password: '', // Password is not returned from API
      changePasswordNextLogin: data.changePasswordNextLogin || false,
      firstName: data.firstName || '',
      lastName: data.lastName || '',
      phone: data.phone || '',
      jobTitle: data.jobTitle || '',
      enabled: data.enabled ?? true,
      avatar: data.avatar || '',
      accessType: 'group', // Default, will be updated below
    };

    // Load current access configuration from the response data
    // The API already returns groups[] and memberships[] in the user response
    loadUserAccessFromResponse(formData, data);

    userData.value = formData;

    logger.debug('Loaded User for editing:', {
      userId: userId.value,
      userData: userData.value,
    });

    // In EDIT mode, skip to Review step by default
    currentStep.value = STEP.REVIEW;

  } catch (error: any) {
    logger.error('Failed to load user:', error);
    notifyFail({ message: t.messages.loadFailed.value });

    // Navigate back to list on load error
    await router.push('/users');
  } finally {
    isLoading.value = false;
  }
}

/**
 * Load user's current access configuration from the API response
 * The getById endpoint already returns groups[] and memberships[] in the response
 *
 * @param {UserFormData} formData - Form data to populate with access info
 * @param {UserResponse} data - API response data containing groups and memberships
 */
function loadUserAccessFromResponse(formData: UserFormData, data: UserResponse): void {
  // Type assertion for extended response fields
  const responseData = data as UserResponse & {
    groups?: Array<{ id: string; name: string; description?: string }>;
    memberships?: Array<{
      orgId: string;
      orgName: string;
      orgType?: string;
      scope: 'local' | 'recursive';
      roleNames: string[];
      via?: string;
    }>;
  };

  const groups = responseData.groups || [];
  const memberships = responseData.memberships || [];

  // Separate direct memberships from group-based memberships
  // Direct memberships don't have "via" field (or via doesn't start with "Group:")
  const directMemberships = memberships.filter(m => !m.via || !m.via.startsWith('Group:'));
  const hasGroupMemberships = groups.length > 0;
  const hasDirectMemberships = directMemberships.length > 0;

  logger.debug('Loading access from response:', {
    groups: groups.length,
    totalMemberships: memberships.length,
    directMemberships: directMemberships.length,
  });

  // Populate groups array for edit mode
  if (hasGroupMemberships) {
    formData.selectedGroups = groups.map((group) => ({
      mode: 'existing' as const,
      existingGroup: {
        groupId: group.id,
        groupName: group.name,
      },
    }));

    // Also populate legacy single value (first group)
    formData.selectedGroup = formData.selectedGroups[0];

    logger.debug('Loaded group memberships for user:', { count: groups.length, groups });
  }

  // Populate direct memberships array for edit mode
  if (hasDirectMemberships) {
    formData.directMemberships = directMemberships.map((membership) => ({
      orgId: membership.orgId,
      orgName: membership.orgName,
      roleIds: [], // IDs not available in this response, only names
      roleNames: membership.roleNames || [],
      scope: membership.scope || 'local',
    }));

    // Also populate legacy single value (first membership)
    formData.directMembership = formData.directMemberships[0];

    logger.debug('Loaded direct memberships for user:', directMemberships.length);
  }

  // Determine access type based on what was found
  if (hasDirectMemberships && hasGroupMemberships) {
    formData.accessType = 'both';
  } else if (hasDirectMemberships) {
    formData.accessType = 'direct';
  } else if (hasGroupMemberships) {
    formData.accessType = 'group';
  } else {
    // No access configuration found - user needs to configure
    logger.debug('No existing access configuration found for user');
    formData.accessType = 'group';
  }
}

/**
 * Update user data with partial updates from step components
 *
 * @param {Partial<typeof userData.value>} partialData - Partial user data to merge
 */
function updateUserData(partialData: Partial<typeof userData.value>): void {
  userData.value = {
    ...userData.value,
    ...partialData,
  };
}

/**
 * Build tour steps with resolved translations and stepper sync hooks
 *
 * @returns {PageTourStep[]} Tour steps with resolved text and callbacks
 */
function buildCreateTourSteps(): PageTourStep[] {
  return CREATE_USER_TOUR_STEPS.map((step) => {
    const key = step.translationKey as keyof typeof t.tour;
    const translation = t.tour[key];

    // Map step keys to stepper steps for auto-sync
    const stepMapping: Record<string, number> = {
      step1: STEP.PERSONAL,
      step2: STEP.SECURITY,
      step3: STEP.ACCESS,
      step4: STEP.REVIEW,
    };

    const targetStep = stepMapping[step.translationKey];

    const result: PageTourStep = {
      element: step.element,
      title: translation.title.value,
      description: translation.description.value,
    };
    if (step.side) result.side = step.side;
    if (step.align) result.align = step.align;
    if (targetStep) {
      result.onHighlightStarted = () => { currentStep.value = targetStep; };
    }
    return result;
  });
}

/** PAGE TOUR (only active in tour mode) */
const { isTourMode } = usePageTour({
  tourId: 'create-user',
  steps: buildCreateTourSteps,
  onTourEnd: () => {
    // Navigate back to users list when tour finishes or user cancels
    if (isTourDemoMode.value) {
      void router.push('/users');
    }
  },
});

/** COMPUTED */

/**
 * Page title (dynamic based on mode)
 */
const pageTitle = computed(() =>
  isEditMode.value ? t.page.titleEdit.value : t.page.title.value
);

/**
 * Page icon (dynamic based on mode)
 */
const pageIcon = computed(() =>
  isEditMode.value ? 'edit' : 'person_add'
);

/**
 * Translated steps array with reactive translations
 */
const translatedSteps = computed(() => [
  {
    title: t.steps.personal.value,
    icon: 'person',
    description: t.stepDescriptions.personal.value,
  },
  {
    title: t.steps.security.value,
    icon: 'security',
    description: t.stepDescriptions.security.value,
  },
  {
    title: t.steps.access.value,
    icon: 'admin_panel_settings',
    description: t.stepDescriptions.access.value,
  },
  {
    title: t.steps.review.value,
    icon: 'check_circle',
    description: t.stepDescriptions.review.value,
  },
]);

/**
 * User form handlers composable
 */
const handlers = useUserFormHandlers({
  userData,
  currentStep,
  isEditMode,
  userId,
  isSaving,
  isTourMode: isTourMode,
  step1FormRef: computed(() => step1Ref.value?.formRef ?? null),
  step2FormRef: computed(() => step2Ref.value?.formRef ?? null),
  step3FormRef: computed(() => step3Ref.value?.formRef ?? null),
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
  save: isEditMode.value ? t.buttons.updateUser.value : t.buttons.createUser.value,
}));

/** COMPOSABLES USAGE */
useStepperNavigation({
  currentStep,
  totalSteps: TOTAL_STEPS,
  changeStep: handlers.handleStepChange,
});

/** LIFECYCLE HOOKS */
onMounted(() => {
  void loadUserData();
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
        icon-color="primary"
        :title="pageTitle"
        :description="t.page.description.value"
        :button="{
          label: t.page.backButton.value,
          icon: 'arrow_back',
          flat: true,
          to: '/users',
        }"
      />

      <!-- Content -->
      <div class="row q-col-gutter-lg">
        <!-- Progress Stepper Vertical -->
        <div id="stepper-sidebar" class="col-12 col-md-4">
          <StepperVertical
            :title="t.sections.progressSteps.value"
            :subtitle="t.messages.completeAllSteps.value"
            :info-text="t.messages.allFieldsRequired.value"
            :current-step-label="t.messages.currentStep.value"
            :current-step="currentStep"
            :steps="translatedSteps"
            :allow-step-navigation="isEditMode || isTourDemoMode"
            @step-click="handlers.changeStep"
          />
        </div>

        <!-- Form Card -->
        <div id="form-card" class="col-12 col-md-8">
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
              <!-- STEP 1: PERSONAL -->
              <Step1Personal
                v-if="currentStep === STEP.PERSONAL"
                ref="step1Ref"
                :model-value="userData"
                :is-edit-mode="isEditMode"
                @update:model-value="updateUserData"
              />

              <!-- STEP 2: SECURITY -->
              <Step2Security
                v-else-if="currentStep === STEP.SECURITY"
                ref="step2Ref"
                :model-value="userData"
                :is-edit-mode="isEditMode"
                @update:model-value="updateUserData"
              />

              <!-- STEP 3: ACCESS -->
              <Step3Access
                v-else-if="currentStep === STEP.ACCESS"
                ref="step3Ref"
                :model-value="userData"
                :is-edit-mode="isEditMode"
                @update:model-value="updateUserData"
              />

              <!-- STEP 4: REVIEW -->
              <Step4Review
                v-else-if="currentStep === STEP.REVIEW"
                :model-value="userData"
                :is-edit-mode="isEditMode"
                @edit-section="handlers.handleStepChange"
              />
            </template>
          </FormCard>
        </div>
      </div>
    </div>
  </q-page>
</template>
