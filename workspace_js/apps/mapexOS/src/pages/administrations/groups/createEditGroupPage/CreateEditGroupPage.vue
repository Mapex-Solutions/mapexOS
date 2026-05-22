<script setup lang="ts">
defineOptions({
  name: 'CreateEditGroupPage',
});

/** TYPE IMPORTS */
import type { FormCardHeader } from '@components/cards';
import type { GroupResponse } from '@mapexos/schemas';
import type { UserSelectorItem } from '@components/drawers';
import type { RoleSelectionItem } from './interfaces';

/** VUE IMPORTS */
import { ref, computed, onMounted } from 'vue';
import { useRouter, useRoute } from 'vue-router';

/** COMPOSABLES */
import { useStepperNavigation } from '@composables/shared/form';
import { useGroupsTranslations } from '@composables/i18n';
import { useLogger } from '@composables/useLogger';

/** COMPONENTS */
import { StepperVertical } from '@components/steppers';
import { PageHeader } from '@components/headers';
import { FormCard } from '@components/cards';
import {
  Step1BasicInfo,
  Step2Roles,
  Step3Members,
  Step4Review,
} from './components';

/** SERVICES */
import { apis } from '@services/mapex';

/** UTILS */
import { notifyFail } from '@utils/alert';

/** LOCAL IMPORTS */
import {
  INITIAL_GROUP_FORM_DATA,
  TOTAL_STEPS,
  STEP,
} from './constants';
import { useGroupFormHandlers } from './handlers';

/** COMPOSABLES & STORES */
const t = useGroupsTranslations();
const router = useRouter();
const route = useRoute();
const logger = useLogger('CreateEditGroupPage');

/** EDIT MODE DETECTION (MANDATORY) */
const isEditMode = ref(!!route.params.id);
const groupId = ref(route.params.id as string | undefined);

/** LOADING STATES */
const isLoading = ref(false);
const isSaving = ref(false);

/** STATE */
const step1Ref = ref<InstanceType<typeof Step1BasicInfo> | null>(null);
const step2RolesRef = ref<InstanceType<typeof Step2Roles> | null>(null);
const step3MembersRef = ref<InstanceType<typeof Step3Members> | null>(null);

const currentStep = ref(1);
const groupData = ref({ ...INITIAL_GROUP_FORM_DATA });
const selectedRoles = ref<RoleSelectionItem[]>([]);
const selectedMembers = ref<string[]>([]);

/** PENDING MEMBER CHANGES */
const pendingAdditions = ref<UserSelectorItem[]>([]);
const pendingRemovals = ref<string[]>([]);

/** INITIAL MEMBERS COUNT (from API in edit mode) */
const initialMembersCount = ref(0);

/** FUNCTIONS */

/**
 * Load group data for edit mode
 * Fetches data from API and populates form state
 * Members are loaded separately by Step3Members component
 * Only executes in edit mode with valid ID
 *
 * @returns {Promise<void>}
 */
async function loadGroupData(): Promise<void> {
  if (!isEditMode.value || !groupId.value) return;

  if (!apis.mapexOS?.groups) {
    notifyFail({ message: t.errors.apiNotInitialized.value });
    return;
  }

  isLoading.value = true;
  try {
    const data: GroupResponse = await apis.mapexOS.groups.getById({ groupId: groupId.value });

    // Populate groupData from API response
    groupData.value = {
      name: data.name || '',
      description: data.description || '',
      enabled: data.enabled ?? true,
    };

    // Store initial members count from API
    initialMembersCount.value = data.membersCount || 0;

    // Load existing roles if available
    const roleIds = (data as any).roleIds as string[] | undefined;
    if (roleIds?.length && apis.mapexOS?.roles) {
      try {
        // Fetch all roles and filter by roleIds
        const rolesResponse = await apis.mapexOS.roles.list({ perPage: 100 });
        const roleIdsSet = new Set(roleIds);
        selectedRoles.value = (rolesResponse.items || [])
          .filter((role: any) => roleIdsSet.has(role.id))
          .map((role: any) => ({
            id: role.id,
            name: role.name || '',
          }));
        logger.debug('Loaded existing roles:', { count: selectedRoles.value.length });
      } catch (rolesError) {
        logger.warn('Failed to load existing roles:', rolesError);
      }
    }

    logger.debug('Loaded Group for editing:', {
      groupId: groupId.value,
      groupData: groupData.value,
      membersCount: initialMembersCount.value,
      rolesCount: selectedRoles.value.length,
    });

    // In EDIT mode, skip to Review step by default
    currentStep.value = STEP.REVIEW;

  } catch (error: any) {
    logger.error('Failed to load group:', error);
    notifyFail({ message: t.createEditNotifications.loadFailed.value });

    // Navigate back to list on load error
    await router.push('/groups');
  } finally {
    isLoading.value = false;
  }
}

/**
 * Update group data with partial updates from step components
 *
 * @param {Partial<typeof groupData.value>} partialData - Partial group data to merge
 */
function updateGroupData(partialData: Partial<typeof groupData.value>): void {
  groupData.value = {
    ...groupData.value,
    ...partialData,
  };
}

/**
 * Handle pending member changes from Step3Members
 *
 * @param {Object} changes - Pending changes object
 * @param {UserSelectorItem[]} changes.additions - Users to add
 * @param {string[]} changes.removals - User IDs to remove
 */
function onPendingChanges(changes: { additions: UserSelectorItem[]; removals: string[] }): void {
  pendingAdditions.value = changes.additions;
  pendingRemovals.value = changes.removals;
  logger.debug('Pending member changes updated:', {
    additions: changes.additions.length,
    removals: changes.removals.length,
  });
}

/**
 * Update selected members list
 *
 * @param {string[]} members - Updated member IDs
 */
function onSelectedMembersUpdate(members: string[]): void {
  selectedMembers.value = members;
}

/**
 * Update selected roles list
 *
 * @param {RoleSelectionItem[]} roles - Updated roles
 */
function onSelectedRolesUpdate(roles: RoleSelectionItem[]): void {
  selectedRoles.value = roles;
  logger.debug('Selected roles updated:', { count: roles.length });
}

/** COMPUTED */

/**
 * Page title (dynamic based on mode)
 */
const pageTitle = computed(() =>
  isEditMode.value ? t.page.titleEdit.value : t.page.titleCreate.value
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
    icon: 'groups',
    description: t.stepDescriptions.basicInfo.value,
  },
  {
    title: t.steps.roles.value,
    icon: 'admin_panel_settings',
    description: t.stepDescriptions.roles.value,
  },
  {
    title: t.steps.members.value,
    icon: 'person_add',
    description: t.stepDescriptions.members.value,
  },
  {
    title: t.steps.review.value,
    icon: 'check_circle',
    description: t.stepDescriptions.review.value,
  },
]);

/**
 * Group form handlers composable
 */
const handlers = useGroupFormHandlers({
  groupData,
  selectedRoles,
  selectedMembers,
  pendingAdditions,
  pendingRemovals,
  currentStep,
  isEditMode,
  groupId,
  isSaving,
  step1FormRef: computed(() => step1Ref.value?.formRef ?? null),
  step2RolesFormRef: computed(() => step2RolesRef.value?.formRef ?? null),
  step3MembersFormRef: computed(() => step3MembersRef.value?.formRef ?? null),
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
  save: isEditMode.value ? t.buttons.updateGroup.value : t.buttons.createGroup.value,
}));

/** COMPOSABLES USAGE */
useStepperNavigation({
  currentStep,
  totalSteps: TOTAL_STEPS,
  changeStep: handlers.handleStepChange,
});

/** LIFECYCLE HOOKS */
onMounted(() => {
  void loadGroupData();
});
</script>

<template>
  <q-page class="q-pa-lg">
    <!-- Loading State (Edit Mode) -->
    <div v-if="isLoading" class="flex flex-center q-pa-xl">
      <q-spinner color="primary" size="50px" />
      <span class="q-ml-md text-grey-7">{{ t.messages.loadingGroup.value }}</span>
    </div>

    <!-- Form Content -->
    <div v-else>
      <!-- Header Section -->
      <PageHeader
        :icon="pageIcon"
        icon-color="primary"
        :title="pageTitle"
        :description="t.page.descriptionCreate.value"
        :button="{
          label: t.page.backButton.value,
          icon: 'arrow_back',
          flat: true,
          to: '/groups',
        }"
      />

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
                :model-value="groupData"
                :is-edit-mode="isEditMode"
                @update:model-value="updateGroupData"
              />

              <!-- STEP 2: ROLES -->
              <Step2Roles
                v-else-if="currentStep === STEP.ROLES"
                ref="step2RolesRef"
                :selected-roles="selectedRoles"
                @update:selected-roles="onSelectedRolesUpdate"
              />

              <!-- STEP 3: MEMBERS -->
              <Step3Members
                v-else-if="currentStep === STEP.MEMBERS"
                ref="step3MembersRef"
                :is-edit-mode="isEditMode"
                :group-id="groupId"
                :selected-members="selectedMembers"
                @pending-changes="onPendingChanges"
                @update:selected-members="onSelectedMembersUpdate"
              />

              <!-- STEP 4: REVIEW -->
              <Step4Review
                v-else-if="currentStep === STEP.REVIEW"
                :model-value="groupData"
                :selected-roles="selectedRoles"
                :selected-members="selectedMembers"
                :initial-members-count="initialMembersCount"
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
