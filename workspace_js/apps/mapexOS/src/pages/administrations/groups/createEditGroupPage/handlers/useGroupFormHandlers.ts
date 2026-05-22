import type { Ref } from 'vue';
import type { QForm } from 'quasar';
import type { GroupFormData, RoleSelectionItem } from '../interfaces';
import type { UserSelectorItem } from '@components/drawers';

import { computed } from 'vue';
import { useRouter } from 'vue-router';

import { notifySuccess, notifyFail } from '@utils/alert/notify';
import { useLogger } from '@composables/useLogger';
import { apis } from '@services/mapex';
import { useOrganizationStore } from '@stores/organization';
import { NAME_MIN_LENGTH } from '../constants';

const logger = useLogger('useGroupFormHandlers');

/**
 * Parameters for useGroupFormHandlers composable
 */
interface UseGroupFormHandlersParams {
  /** Group form data ref */
  groupData: Ref<GroupFormData>;

  /** Selected roles ref */
  selectedRoles: Ref<RoleSelectionItem[]>;

  /** Selected member IDs ref */
  selectedMembers: Ref<string[]>;

  /** Pending member additions */
  pendingAdditions: Ref<UserSelectorItem[]>;

  /** Pending member removals (user IDs) */
  pendingRemovals: Ref<string[]>;

  /** Current step ref */
  currentStep: Ref<number>;

  /** Whether in edit mode */
  isEditMode: Ref<boolean>;

  /** Group ID for edit mode */
  groupId: Ref<string | undefined>;

  /** Whether form is saving */
  isSaving: Ref<boolean>;

  /** Step 1 form ref */
  step1FormRef: Ref<QForm | null>;

  /** Step 2 Roles form ref */
  step2RolesFormRef: Ref<QForm | null>;

  /** Step 3 Members form ref */
  step3MembersFormRef: Ref<QForm | null>;

  /** Translations object */
  t: any;
}

/**
 * Composable for group form handlers
 * Handles validation, step navigation, and form submission
 *
 * @param {UseGroupFormHandlersParams} params - Handler parameters
 * @returns Form handler functions and computed properties
 */
export function useGroupFormHandlers(params: UseGroupFormHandlersParams) {
  const {
    groupData,
    selectedRoles,
    selectedMembers,
    pendingAdditions,
    pendingRemovals,
    currentStep,
    isEditMode,
    groupId,
    isSaving,
    step1FormRef,
    step2RolesFormRef,
    step3MembersFormRef,
    t,
  } = params;

  const router = useRouter();
  const orgStore = useOrganizationStore();

  /**
   * Get selected members count
   */
  const selectedMembersCount = computed(() => selectedMembers.value.length);

  /**
   * Check if Next button should be disabled
   */
  const isNextButtonDisabled = computed(() => {
    if (currentStep.value === 1) {
      // Step 1: Basic info - require name
      const { name } = groupData.value;
      if (!name || name.length < NAME_MIN_LENGTH) return true;
      return false;
    }

    if (currentStep.value === 2) {
      // Step 2: Roles - at least one role required
      if (selectedRoles.value.length === 0) return true;
      return false;
    }

    // Step 3: Members - no minimum requirement (group can be empty)
    return false;
  });

  /**
   * Validate and change step
   *
   * @param {number} step - Target step number
   */
  async function changeStep(step: number): Promise<void> {
    // Validate current step form before proceeding forward
    if (step > currentStep.value) {
      if (currentStep.value === 1 && step1FormRef.value) {
        const valid = await step1FormRef.value.validate();
        if (!valid) return;
      }

      if (currentStep.value === 2) {
        // Step 2: Roles - validate that at least one role is selected
        if (selectedRoles.value.length === 0) {
          logger.warn('At least one role is required');
          return;
        }
        if (step2RolesFormRef.value) {
          const valid = await step2RolesFormRef.value.validate();
          if (!valid) return;
        }
      }

      if (currentStep.value === 3 && step3MembersFormRef.value) {
        const valid = await step3MembersFormRef.value.validate();
        if (!valid) return;
      }
    }

    currentStep.value = step;
  }

  /**
   * Wrapper for step navigation (non-async)
   *
   * @param {number} step - Target step number
   */
  function handleStepChange(step: number): void {
    void changeStep(step);
  }

  /**
   * Persist member changes (add/remove) for a group
   *
   * @param {string} targetGroupId - Group ID to update
   * @returns {Promise<void>}
   */
  async function persistMemberChanges(targetGroupId: string): Promise<void> {
    if (!apis.mapexOS?.groups) return;

    // Process removals first
    for (const userId of pendingRemovals.value) {
      try {
        await apis.mapexOS.groups.removeMember({
          groupId: targetGroupId,
          userId,
        });
        logger.debug('Removed member:', { groupId: targetGroupId, userId });
      } catch (error: any) {
        logger.warn('Failed to remove member:', { userId, error: error?.message });
      }
    }

    // Then process additions
    for (const user of pendingAdditions.value) {
      try {
        await apis.mapexOS.groups.addMember(
          { groupId: targetGroupId },
          { userId: user.id },
        );
        logger.debug('Added member:', { groupId: targetGroupId, userId: user.id });
      } catch (error: any) {
        logger.warn('Failed to add member:', { userId: user.id, error: error?.message });
      }
    }
  }

  /**
   * Submit form - create or update based on mode
   * Handles both CREATE and EDIT operations in single function
   *
   * @returns {Promise<void>}
   */
  async function submitForm(): Promise<void> {
    if (!apis.mapexOS?.groups) {
      notifyFail({ message: t.errors.apiNotInitialized.value });
      return;
    }

    isSaving.value = true;

    try {
      if (isEditMode.value && groupId.value) {
        // UPDATE existing group
        const updatePayload = {
          name: groupData.value.name,
          description: groupData.value.description || undefined,
          enabled: groupData.value.enabled,
        };

        logger.debug('Updating Group:', { groupId: groupId.value, payload: updatePayload });

        await apis.mapexOS.groups.update(
          { groupId: groupId.value },
          updatePayload
        );

        // Persist member changes
        if (pendingAdditions.value.length > 0 || pendingRemovals.value.length > 0) {
          await persistMemberChanges(groupId.value);
        }

        notifySuccess({
          message: t.createEditNotifications.updated.value,
          timeout: 3000,
        });
      } else {
        // CREATE new group
        const createPayload = {
          name: groupData.value.name,
          description: groupData.value.description || undefined,
          enabled: groupData.value.enabled,
          isSystem: false,
          orgId: orgStore.selectedOrganizationId || undefined,
          roleIds: selectedRoles.value.map(r => r.id),
        };

        logger.debug('Creating Group:', createPayload);

        const createdGroup = await apis.mapexOS.groups.create(createPayload);

        // Add members after creation if any selected
        if (pendingAdditions.value.length > 0 && createdGroup?.id) {
          await persistMemberChanges(createdGroup.id);
        }

        notifySuccess({
          message: t.createEditNotifications.created.value,
          timeout: 3000,
        });
      }

      // Navigate to groups list after success
      await router.push('/groups');
    } catch (error: any) {
      logger.error('Form submission error:', error);

      // Handle specific error codes
      const errorCode = error?.response?.status || error?.code;
      let errorMessage = isEditMode.value
        ? t.createEditNotifications.updateFailed.value
        : t.createEditNotifications.createFailed.value;

      if (errorCode === 409) {
        errorMessage = t.createEditNotifications.alreadyExists.value;
      } else if (errorCode === 403) {
        errorMessage = t.createEditNotifications.forbidden.value;
      }

      notifyFail({
        message: errorMessage,
        timeout: 5000,
      });
    } finally {
      isSaving.value = false;
    }
  }

  return {
    selectedMembersCount,
    isNextButtonDisabled,
    changeStep,
    handleStepChange,
    submitForm,
  };
}
