<script setup lang="ts">
defineOptions({
  name: 'CreateEditRolePage',
});

/** TYPE IMPORTS */
import type { FormCardHeader } from '@components/cards';
import type { RoleResponse } from '@mapexos/schemas';

/** VUE IMPORTS */
import { ref, computed, onMounted } from 'vue';
import { useRouter, useRoute } from 'vue-router';

/** COMPOSABLES */
import { useStepperNavigation } from '@composables/shared/form';
import { useRolesTranslations } from '@composables/i18n';
import { useLogger } from '@composables/useLogger';

/** COMPONENTS */
import { StepperVertical } from '@components/steppers';
import { PageHeader } from '@components/headers';
import { FormCard } from '@components/cards';
import {
  Step1BasicInfo,
  Step2Permissions,
  Step3Review,
} from './components';

/** SERVICES */
import { apis } from '@services/mapex';

/** UTILS */
import { notifyFail } from '@utils/alert';

/** LOCAL IMPORTS */
import {
  INITIAL_ROLE_FORM_DATA,
  TOTAL_STEPS,
  STEP,
  DEFAULT_RESOURCE_PERMISSIONS,
} from './constants';
import { useRoleFormHandlers } from './handlers';

/** COMPOSABLES & STORES */
const t = useRolesTranslations();
const router = useRouter();
const route = useRoute();
const logger = useLogger('CreateEditRolePage');

/** EDIT MODE DETECTION (MANDATORY) */
const isEditMode = ref(!!route.params.id);
const roleId = ref(route.params.id as string | undefined);

/** LOADING STATES */
const isLoading = ref(false);
const isSaving = ref(false);

/** SYSTEM ROLE CHECK */
const isSystemRole = ref(false);

/** STATE */
const step1Ref = ref<InstanceType<typeof Step1BasicInfo> | null>(null);
const step2Ref = ref<InstanceType<typeof Step2Permissions> | null>(null);

const currentStep = ref(1);
const roleData = ref({ ...INITIAL_ROLE_FORM_DATA });
const resourcePermissions = ref(JSON.parse(JSON.stringify(DEFAULT_RESOURCE_PERMISSIONS)));

/** FUNCTIONS */

/**
 * Load role data for edit mode
 * Fetches data from API and populates form state
 * Only executes in edit mode with valid ID
 *
 * @returns {Promise<void>}
 */
async function loadRoleData(): Promise<void> {
  if (!isEditMode.value || !roleId.value) return;

  if (!apis.mapexOS?.roles) {
    notifyFail({ message: t.errors.apiNotInitialized.value });
    return;
  }

  isLoading.value = true;
  try {
    const data: RoleResponse = await apis.mapexOS.roles.getById({ roleId: roleId.value });

    // Check if system role
    isSystemRole.value = data.isSystem || false;

    // Populate roleData from API response
    roleData.value = {
      name: data.name || '',
      description: data.description || '',
      scope: data.scope as 'global' | 'local' || null,
      isTemplate: data.isTemplate || false,
    };

    // Populate permissions from API response
    if (data.permissions && data.permissions.length > 0) {
      populatePermissionsFromArray(data.permissions);
    }

    logger.debug('Loaded Role for editing:', {
      roleId: roleId.value,
      roleData: roleData.value,
      isSystemRole: isSystemRole.value,
    });

    // In EDIT mode, skip to Review step by default
    currentStep.value = STEP.REVIEW;

  } catch (error: any) {
    logger.error('Failed to load role:', error);
    notifyFail({ message: t.notifications.loadFailed.value });

    // Navigate back to list on load error
    await router.push('/roles');
  } finally {
    isLoading.value = false;
  }
}

/**
 * Populate resource permissions from permissions array
 * Converts ['users.list', 'events.raw.list'] to UI state
 * Supports both simple (resource.action) and compound (resource.subtype.action) permission keys
 *
 * @param {string[]} permissions - Array of permission strings from backend
 */
function populatePermissionsFromArray(permissions: string[]): void {
  // Reset all permissions first
  resourcePermissions.value.forEach((resource: typeof resourcePermissions.value[0]) => {
    resource.enabled = false;
    resource.actions.forEach((action: typeof resource.actions[0]) => {
      action.granted = false;
    });
  });

  // Set granted permissions by matching against permissionKey or default pattern
  permissions.forEach((permission: string) => {
    for (const resource of resourcePermissions.value) {
      for (const action of resource.actions) {
        const expectedKey = action.permissionKey || `${resource.resource}.${action.name}`;
        if (permission === expectedKey) {
          action.granted = true;
          resource.enabled = true;
        }
      }
    }
  });
}

/**
 * Update role data with partial updates from step components
 *
 * @param {Partial<typeof roleData.value>} partialData - Partial role data to merge
 */
function updateRoleData(partialData: Partial<typeof roleData.value>): void {
  roleData.value = {
    ...roleData.value,
    ...partialData,
  };
}

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
  isEditMode.value ? 'edit' : 'add_circle'
);

/**
 * Translated steps array with reactive translations
 */
const translatedSteps = computed(() => [
  {
    title: t.steps.basicInfo.value,
    icon: 'badge',
    description: t.stepDescriptions.basicInfo.value,
  },
  {
    title: t.steps.permissions.value,
    icon: 'vpn_key',
    description: t.stepDescriptions.permissions.value,
  },
  {
    title: t.steps.review.value,
    icon: 'check_circle',
    description: t.stepDescriptions.review.value,
  },
]);

/**
 * Role form handlers composable
 */
const handlers = useRoleFormHandlers({
  roleData,
  resourcePermissions,
  currentStep,
  isEditMode,
  roleId,
  isSaving,
  step1Ref: computed(() => step1Ref.value),
  step2Ref: computed(() => step2Ref.value),
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
  save: isEditMode.value ? t.buttons.updateRole.value : t.buttons.createRole.value,
}));

/** COMPOSABLES USAGE */
useStepperNavigation({
  currentStep,
  totalSteps: TOTAL_STEPS,
  changeStep: handlers.handleStepChange,
});

/** LIFECYCLE HOOKS */
onMounted(() => {
  void loadRoleData();
});
</script>

<template>
  <q-page class="q-pa-lg">
    <!-- Loading State (Edit Mode) -->
    <div v-if="isLoading" class="flex flex-center q-pa-xl">
      <q-spinner color="primary" size="50px" />
      <span class="q-ml-md text-grey-7">{{ t.messages.loadingRole.value }}</span>
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
          to: '/roles',
        }"
      />

      <!-- System Role Warning -->
      <q-banner
        v-if="isEditMode && isSystemRole"
        rounded
        class="bg-warning text-white q-mb-lg"
      >
        <template #avatar>
          <q-icon name="lock" color="white" />
        </template>
        {{ t.messages.systemRoleWarning.value }}
      </q-banner>

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
              <!-- STEP 1: BASIC INFO -->
              <Step1BasicInfo
                v-if="currentStep === STEP.BASIC_INFO"
                ref="step1Ref"
                :model-value="roleData"
                :is-edit-mode="isEditMode"
                :is-system-role="isSystemRole"
                @update:model-value="updateRoleData"
              />

              <!-- STEP 2: PERMISSIONS -->
              <Step2Permissions
                v-else-if="currentStep === STEP.PERMISSIONS"
                ref="step2Ref"
                :resource-permissions="resourcePermissions"
                @resource-toggle="handlers.onResourceToggle"
                @action-toggle="handlers.onActionToggle"
                @toggle-all-actions="handlers.onToggleAllActions"
              />

              <!-- STEP 3: REVIEW -->
              <Step3Review
                v-else-if="currentStep === STEP.REVIEW"
                :model-value="roleData"
                :resource-permissions="resourcePermissions"
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
