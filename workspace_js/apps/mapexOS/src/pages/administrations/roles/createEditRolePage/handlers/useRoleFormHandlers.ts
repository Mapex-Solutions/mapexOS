import type { Ref } from 'vue';
import type { QForm } from 'quasar';
import type { RoleFormData, ResourcePermission } from '../interfaces';

import { computed } from 'vue';
import { useRouter } from 'vue-router';

import { notifySuccess, notifyFail } from '@utils/alert/notify';
import { useLogger } from '@composables/useLogger';
import { apis } from '@services/mapex';
import { useOrganizationStore } from '@stores/organization';
import { NAME_MIN_LENGTH } from '../constants';

const logger = useLogger('useRoleFormHandlers');

/**
 * Step component ref with validate method
 */
interface StepComponentRef {
  formRef?: QForm | null;
  validate: () => boolean | Promise<boolean>;
}

/**
 * Parameters for useRoleFormHandlers composable
 */
interface UseRoleFormHandlersParams {
  /** Role form data ref */
  roleData: Ref<RoleFormData>;

  /** Resource permissions ref */
  resourcePermissions: Ref<ResourcePermission[]>;

  /** Current step ref */
  currentStep: Ref<number>;

  /** Whether in edit mode */
  isEditMode: Ref<boolean>;

  /** Role ID for edit mode */
  roleId: Ref<string | undefined>;

  /** Whether form is saving */
  isSaving: Ref<boolean>;

  /** Step 1 component ref */
  step1Ref: Ref<StepComponentRef | null>;

  /** Step 2 component ref */
  step2Ref: Ref<StepComponentRef | null>;

  /** Translations object */
  t: any;
}

/**
 * Composable for role form handlers
 * Handles validation, step navigation, and form submission
 *
 * @param {UseRoleFormHandlersParams} params - Handler parameters
 * @returns Form handler functions and computed properties
 */
export function useRoleFormHandlers(params: UseRoleFormHandlersParams) {
  const {
    roleData,
    resourcePermissions,
    currentStep,
    isEditMode,
    roleId,
    isSaving,
    step1Ref,
    step2Ref,
    t,
  } = params;

  const router = useRouter();
  const orgStore = useOrganizationStore();

  /**
   * Toggle resource enabled state
   *
   * @param {number} index - Resource index
   */
  function onResourceToggle(index: number): void {
    const resource = resourcePermissions.value[index];
    if (resource) {
      resource.enabled = !resource.enabled;
      // Sync all actions with resource enabled state
      resource.actions.forEach(action => {
        action.granted = resource.enabled;
      });
    }
  }

  /**
   * Toggle permission action within a resource
   * Auto-selects "list" when selecting create, read, update, or delete
   *
   * @param {number} resourceIndex - Resource index
   * @param {number} actionIndex - Action index
   */
  function onActionToggle(resourceIndex: number, actionIndex: number): void {
    const resource = resourcePermissions.value[resourceIndex];
    if (resource && resource.actions[actionIndex]) {
      const action = resource.actions[actionIndex];
      action.granted = !action.granted;

      // Auto-select "list" when selecting create, read, update, or delete
      if (action.granted && ['create', 'read', 'update', 'delete'].includes(action.name)) {
        const listAction = resource.actions.find(a => a.name === 'list');
        if (listAction && !listAction.granted) {
          listAction.granted = true;
        }
      }

      // Auto-enable resource if any action is granted
      const hasAnyGranted = resource.actions.some(a => a.granted);
      resource.enabled = hasAnyGranted;
    }
  }

  /**
   * Toggle all actions for a resource
   *
   * @param {number} resourceIndex - Resource index
   * @param {boolean} granted - Whether to grant or revoke all
   */
  function onToggleAllActions(resourceIndex: number, granted: boolean): void {
    const resource = resourcePermissions.value[resourceIndex];
    if (resource) {
      resource.enabled = granted;
      resource.actions.forEach(action => {
        action.granted = granted;
      });
    }
  }

  /**
   * Get selected permissions count
   */
  const selectedPermissionsCount = computed(() => {
    let count = 0;
    resourcePermissions.value.forEach(resource => {
      resource.actions.forEach(action => {
        if (action.granted) count++;
      });
    });
    return count;
  });

  /**
   * Check if Next button should be disabled
   */
  const isNextButtonDisabled = computed(() => {
    if (currentStep.value === 1) {
      // Step 1: Basic info - require name and scope
      const { name, scope } = roleData.value;
      if (!name || name.length < NAME_MIN_LENGTH) return true;
      if (!scope) return true;
      return false;
    }

    if (currentStep.value === 2) {
      // Step 2: Permissions - require at least 1 permission
      return selectedPermissionsCount.value === 0;
    }

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
      if (currentStep.value === 1 && step1Ref.value) {
        const valid = await step1Ref.value.validate();
        if (!valid) return;
      }

      if (currentStep.value === 2 && step2Ref.value) {
        const valid = await step2Ref.value.validate();
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
   * Build permissions array from resource permissions state
   *
   * @returns {string[]} Array of permission strings (e.g., ['users.list', 'users.create'])
   */
  function buildPermissionsArray(): string[] {
    const permissions: string[] = [];

    resourcePermissions.value.forEach(resource => {
      resource.actions.forEach(action => {
        if (action.granted) {
          permissions.push(action.permissionKey || `${resource.resource}.${action.name}`);
        }
      });
    });

    return permissions;
  }

  /**
   * Submit form - create or update based on mode
   * Handles both CREATE and EDIT operations in single function
   *
   * @returns {Promise<void>}
   */
  async function submitForm(): Promise<void> {
    if (!apis.mapexOS?.roles) {
      notifyFail({ message: t.errors.apiNotInitialized.value });
      return;
    }

    isSaving.value = true;

    try {
      // Build permissions array
      const permissions = buildPermissionsArray();

      if (permissions.length === 0) {
        notifyFail({ message: t.notifications.noPermissions.value });
        isSaving.value = false;
        return;
      }

      if (isEditMode.value && roleId.value) {
        // UPDATE existing role
        const updatePayload = {
          name: roleData.value.name,
          description: roleData.value.description || undefined,
          permissions,
        };

        logger.debug('Updating Role:', { roleId: roleId.value, payload: updatePayload });

        await apis.mapexOS.roles.update(
          { roleId: roleId.value },
          updatePayload
        );

        notifySuccess({
          message: t.notifications.updated.value,
          timeout: 3000,
        });
      } else {
        // CREATE new role
        // pathKey is derived from current organization's pathKey
        const currentOrg = orgStore.flatList.find(org => org.id === orgStore.selectedOrganizationId);
        const pathKey = currentOrg?.pathKey || '';

        if (!pathKey) {
          notifyFail({ message: t.errors.orgPathKeyMissing.value });
          isSaving.value = false;
          return;
        }

        const createPayload = {
          name: roleData.value.name,
          description: roleData.value.description || undefined,
          permissions,
          scope: roleData.value.scope as 'global' | 'local',
          isSystem: false, // Users cannot create system roles
          isTemplate: roleData.value.isTemplate,
          pathKey,
        };

        logger.debug('Creating Role:', createPayload);

        await apis.mapexOS.roles.create(createPayload);

        notifySuccess({
          message: t.notifications.created.value,
          timeout: 3000,
        });
      }

      // Navigate to roles list after success
      await router.push('/roles');
    } catch (error: any) {
      logger.error('Form submission error:', error);

      // Handle specific error codes
      const errorCode = error?.response?.status || error?.code;
      let errorMessage = isEditMode.value
        ? t.notifications.updateFailed.value
        : t.notifications.createFailed.value;

      if (errorCode === 409) {
        errorMessage = t.notifications.alreadyExists.value;
      } else if (errorCode === 403) {
        errorMessage = t.notifications.forbidden.value;
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
    onResourceToggle,
    onActionToggle,
    onToggleAllActions,
    selectedPermissionsCount,
    isNextButtonDisabled,
    changeStep,
    handleStepChange,
    buildPermissionsArray,
    submitForm,
  };
}
