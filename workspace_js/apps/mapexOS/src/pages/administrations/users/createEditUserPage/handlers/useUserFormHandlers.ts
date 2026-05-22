import type { Ref } from 'vue';
import type { QForm } from 'quasar';
import type { UserFormData } from '../interfaces';

import { computed } from 'vue';
import { useRouter } from 'vue-router';

import { notifySuccess, notifyFail } from '@utils/alert/notify';
import { useLogger } from '@composables/useLogger';
import { apis } from '@services/mapex';
import { VALIDATION } from '../constants';

const logger = useLogger('useUserFormHandlers');

/**
 * Parameters for useUserFormHandlers composable
 */
interface UseUserFormHandlersParams {
  /** User form data ref */
  userData: Ref<UserFormData>;

  /** Current step ref */
  currentStep: Ref<number>;

  /** Whether in edit mode */
  isEditMode: Ref<boolean>;

  /** User ID for edit mode */
  userId: Ref<string | undefined>;

  /** Whether form is saving */
  isSaving: Ref<boolean>;

  /** Whether in tour demo mode (skips validation and API calls) */
  isTourMode?: Ref<boolean>;

  /** Step 1 form ref (Personal) */
  step1FormRef: Ref<QForm | null>;

  /** Step 2 form ref (Security) */
  step2FormRef: Ref<QForm | null>;

  /** Step 3 form ref (Access) */
  step3FormRef: Ref<QForm | null>;

  /** Translations object */
  t: any;
}

/**
 * Composable for user form handlers
 * Handles validation, step navigation, and form submission
 *
 * V1: AuthProvider removed - always internal auth.
 * Steps: Personal (1), Security (2), Access (3), Review (4)
 *
 * @param {UseUserFormHandlersParams} params - Handler parameters
 * @returns Form handler functions and computed properties
 */
export function useUserFormHandlers(params: UseUserFormHandlersParams) {
  const {
    userData,
    currentStep,
    isEditMode,
    userId,
    isSaving,
    isTourMode,
    step1FormRef,
    step2FormRef,
    step3FormRef,
    t,
  } = params;

  const router = useRouter();

  /**
   * Check if Next button should be disabled
   */
  const isNextButtonDisabled = computed(() => {
    if (currentStep.value === 1) {
      // Step 1: Personal
      const { firstName, lastName, email } = userData.value;
      if (!firstName || firstName.length < VALIDATION.FIRST_NAME_MIN_LENGTH) return true;
      if (!lastName || lastName.length < VALIDATION.LAST_NAME_MIN_LENGTH) return true;
      if (!email) return true;
      return false;
    }

    if (currentStep.value === 2) {
      // Step 2: Security
      const { password } = userData.value;
      // V1: Always internal auth - require password for create mode
      if (!isEditMode.value) {
        if (!password || password.length < VALIDATION.PASSWORD_MIN_LENGTH) return true;
      }
      return false;
    }

    if (currentStep.value === 3) {
      // Step 3: Access - validate based on access type
      const { accessType, selectedGroup, selectedGroups, directMembership, directMemberships } = userData.value;

      // Helper to check if groups are valid
      const hasValidGroups = (): boolean => {
        // Check new array format first
        if (selectedGroups && selectedGroups.length > 0) {
          return selectedGroups.every(g => {
            if (g.mode === 'existing') return !!g.existingGroup?.groupId;
            if (g.mode === 'new') return g.newGroup?.name && g.newGroup.name.length >= 3 && g.newGroup.roleIds?.length > 0;
            return false;
          });
        }
        // Fallback to legacy single group
        if (selectedGroup?.mode === 'existing') return !!selectedGroup.existingGroup?.groupId;
        if (selectedGroup?.mode === 'new') {
          return !!(selectedGroup.newGroup?.name && selectedGroup.newGroup.name.length >= 3 && selectedGroup.newGroup.roleIds?.length > 0);
        }
        return false;
      };

      // Helper to check if direct memberships are valid
      const hasValidMemberships = (): boolean => {
        // Check new array format first
        if (directMemberships && directMemberships.length > 0) {
          return directMemberships.every(m => m.roleIds && m.roleIds.length > 0);
        }
        // Fallback to legacy single membership
        return !!(directMembership?.roleIds && directMembership.roleIds.length > 0);
      };

      if (accessType === 'group') {
        return !hasValidGroups();
      }

      if (accessType === 'direct') {
        return !hasValidMemberships();
      }

      if (accessType === 'both') {
        // Both requires at least one group OR one membership
        return !hasValidGroups() && !hasValidMemberships();
      }

      // No access type selected
      return true;
    }

    return false;
  });

  /**
   * Validate and change step
   * In tour mode, validation is skipped
   *
   * @param {number} step - Target step number
   */
  async function changeStep(step: number): Promise<void> {
    // Skip validation in tour mode
    if (isTourMode?.value) {
      currentStep.value = step;
      return;
    }

    // Validate current step form before proceeding forward
    if (step > currentStep.value) {
      if (currentStep.value === 1 && step1FormRef.value) {
        const valid = await step1FormRef.value.validate();
        if (!valid) return;
      }

      if (currentStep.value === 2 && step2FormRef.value) {
        const valid = await step2FormRef.value.validate();
        if (!valid) return;
      }

      if (currentStep.value === 3 && step3FormRef.value) {
        const valid = await step3FormRef.value.validate();
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
   * Submit form - create or update based on mode
   * Handles both CREATE and EDIT operations in single function
   * CREATE uses onboarding API to create user with memberships atomically
   * In tour mode, skips API call and navigates to users list
   *
   * @returns {Promise<void>}
   */
  async function submitForm(): Promise<void> {
    // In tour mode, skip API call and just navigate back
    if (isTourMode?.value) {
      logger.debug('Tour mode: skipping API call');
      await router.push('/users');
      return;
    }

    isSaving.value = true;

    try {
      if (isEditMode.value && userId.value) {
        // UPDATE existing user with access (uses onboarding API for atomic update)
        if (!apis.mapexOS?.onboarding) {
          notifyFail({ message: t.errors.onboardingApiNotInitialized.value });
          return;
        }

        // Build the update payload with user data
        const updatePayload: Record<string, unknown> = {
          firstName: userData.value.firstName,
          lastName: userData.value.lastName,
          phone: userData.value.phone || undefined,
          jobTitle: userData.value.jobTitle || undefined,
          enabled: userData.value.enabled,
          avatar: userData.value.avatar || undefined,
        };

        // Only include password if provided
        if (userData.value.password) {
          updatePayload.password = userData.value.password;
        }

        // V1: Always internal auth
        updatePayload.changePasswordNextLogin = userData.value.changePasswordNextLogin;

        // Handle access configuration (supports multiple groups + multiple memberships)
        const { accessType, selectedGroup, selectedGroups, directMembership, directMemberships } = userData.value;

        // Build groups array
        const groupsToSend: any[] = [];

        if (selectedGroups && selectedGroups.length > 0) {
          selectedGroups.forEach(g => {
            if (g.mode === 'existing' && g.existingGroup?.groupId) {
              groupsToSend.push({ existingGroup: { groupId: g.existingGroup.groupId } });
            } else if (g.mode === 'new' && g.newGroup?.name && g.newGroup?.roleIds?.length) {
              groupsToSend.push({
                newGroup: {
                  name: g.newGroup.name,
                  description: g.newGroup.description || undefined,
                  roleIds: g.newGroup.roleIds,
                },
              });
            }
          });
        } else if ((accessType === 'group' || accessType === 'both') && selectedGroup) {
          if (selectedGroup.mode === 'existing' && selectedGroup.existingGroup?.groupId) {
            groupsToSend.push({ existingGroup: { groupId: selectedGroup.existingGroup.groupId } });
          } else if (selectedGroup.mode === 'new' && selectedGroup.newGroup?.name && selectedGroup.newGroup?.roleIds?.length) {
            groupsToSend.push({
              newGroup: {
                name: selectedGroup.newGroup.name,
                description: selectedGroup.newGroup.description || undefined,
                roleIds: selectedGroup.newGroup.roleIds,
              },
            });
          }
        }

        // Build memberships array
        const membershipsToSend: any[] = [];

        if (directMemberships && directMemberships.length > 0) {
          directMemberships.forEach(m => {
            if (m.roleIds && m.roleIds.length > 0) {
              membershipsToSend.push({
                roles: m.roleIds,
                scope: m.scope || undefined,
              });
            }
          });
        } else if ((accessType === 'direct' || accessType === 'both') && directMembership?.roleIds?.length) {
          membershipsToSend.push({
            roles: directMembership.roleIds,
            scope: directMembership.scope || undefined,
          });
        }

        // Only include arrays if they have items
        if (groupsToSend.length > 0 || accessType === 'group' || accessType === 'both') {
          updatePayload.groups = groupsToSend;
        }
        if (membershipsToSend.length > 0 || accessType === 'direct' || accessType === 'both') {
          updatePayload.memberships = membershipsToSend;
        }

        logger.debug('Updating User with Access (Onboarding):', { userId: userId.value, payload: updatePayload });

        await apis.mapexOS.onboarding.updateUserWithAccess(
          { userId: userId.value },
          updatePayload,
        );

        notifySuccess({
          message: t.messages.updated.value,
          timeout: 3000,
        });
      } else {
        // CREATE new user with membership (uses onboarding API)
        if (!apis.mapexOS?.onboarding) {
          notifyFail({ message: t.errors.onboardingApiNotInitialized.value });
          return;
        }

        // Build the onboarding payload
        // V1: AuthProvider not sent - backend defaults to internal
        const createPayload: Record<string, unknown> = {
          email: userData.value.email,
          firstName: userData.value.firstName,
          lastName: userData.value.lastName,
          enabled: userData.value.enabled,
          changePasswordNextLogin: userData.value.changePasswordNextLogin,
          password: userData.value.password,
        };

        // Add optional fields
        if (userData.value.phone) {
          createPayload.phone = userData.value.phone;
        }
        if (userData.value.jobTitle) {
          createPayload.jobTitle = userData.value.jobTitle;
        }
        if (userData.value.avatar) {
          createPayload.avatar = userData.value.avatar;
        }

        // Handle access type: group OR direct
        const { accessType: createAccessType, selectedGroup: createSelectedGroup, directMembership: createDirectMembership } = userData.value;

        if (createAccessType === 'group' && createSelectedGroup) {
          if (createSelectedGroup.mode === 'existing' && createSelectedGroup.existingGroup?.groupId) {
            createPayload.groups = [
              {
                existingGroup: {
                  groupId: createSelectedGroup.existingGroup.groupId,
                },
              },
            ];
          } else if (createSelectedGroup.mode === 'new' && createSelectedGroup.newGroup?.name && createSelectedGroup.newGroup?.roleIds?.length) {
            createPayload.groups = [
              {
                newGroup: {
                  name: createSelectedGroup.newGroup.name,
                  description: createSelectedGroup.newGroup.description || undefined,
                  roleIds: createSelectedGroup.newGroup.roleIds,
                },
              },
            ];
          }
        } else if (createAccessType === 'direct' && createDirectMembership?.roleIds?.length) {
          createPayload.memberships = [
            {
              roles: createDirectMembership.roleIds,
              scope: createDirectMembership.scope || undefined,
            },
          ];
        }

        logger.debug('Creating User with Membership (Onboarding):', createPayload);

        await apis.mapexOS.onboarding.createUserWithMemberships(createPayload as any);

        notifySuccess({
          message: t.messages.created.value,
          timeout: 3000,
        });
      }

      // Navigate to users list after success
      await router.push('/users');
    } catch (error: any) {
      logger.error('Form submission error:', error);

      // Handle specific error codes
      const errorCode = error?.response?.status || error?.code;
      let errorMessage = isEditMode.value
        ? t.messages.updateFailed.value
        : t.messages.createFailed.value;

      if (errorCode === 409) {
        errorMessage = t.messages.emailExists.value;
      } else if (errorCode === 403) {
        errorMessage = t.messages.forbidden.value;
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
    isNextButtonDisabled,
    changeStep,
    handleStepChange,
    submitForm,
  };
}
