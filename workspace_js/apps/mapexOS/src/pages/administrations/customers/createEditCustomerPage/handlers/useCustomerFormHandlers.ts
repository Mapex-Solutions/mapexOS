import type { Ref, ComputedRef } from 'vue';
import type { QForm } from 'quasar';
import type { OrganizationFormData, OrganizationType, OrgTypeConfig } from '../interfaces';

import { computed } from 'vue';
import { useRouter } from 'vue-router';

import { notifySuccess, notifyFail } from '@utils/alert/notify';
import { useLogger } from '@composables/useLogger';
import { apis } from '@services/mapex';
import { useOrganizationStore } from '@stores/organization';
import { NAME_MIN_LENGTH } from '../constants';

const logger = useLogger('useCustomerFormHandlers');

/**
 * Parameters for useCustomerFormHandlers composable
 */
interface UseCustomerFormHandlersParams {
  /** Organization form data ref */
  formData: Ref<OrganizationFormData>;

  /** Current step ref */
  currentStep: Ref<number>;

  /** Whether in edit mode */
  isEditMode: Ref<boolean>;

  /** Organization ID for edit mode */
  organizationId: Ref<string | undefined>;

  /** Whether form is saving */
  isSaving: Ref<boolean>;

  /** Organization type being created/edited */
  orgType: ComputedRef<OrganizationType>;

  /** Parent organization ID (for create mode) */
  parentOrgId: ComputedRef<string | undefined>;

  /** Type config for the current org type */
  typeConfig: ComputedRef<OrgTypeConfig>;

  /** Step 1 form ref */
  step1FormRef: Ref<QForm | null>;

  /** Step 2 form ref (address - only for types with address) */
  step2FormRef: Ref<QForm | null>;

  /** Step 3 form ref (access policy - step number varies) */
  accessPolicyFormRef: Ref<QForm | null>;

  /** Translations object */
  t: any;
}

/**
 * Composable for organization form handlers
 * Handles validation, step navigation, and form submission
 *
 * @param {UseCustomerFormHandlersParams} params - Handler parameters
 * @returns Form handler functions and computed properties
 */
export function useCustomerFormHandlers(params: UseCustomerFormHandlersParams) {
  const {
    formData,
    currentStep,
    isEditMode,
    organizationId,
    isSaving,
    orgType,
    parentOrgId,
    typeConfig,
    step1FormRef,
    step2FormRef,
    accessPolicyFormRef,
    t,
  } = params;

  const router = useRouter();
  const orgStore = useOrganizationStore();

  /**
   * Get the step number for access policy based on type config
   *
   * @returns {number} Step number for access policy
   */
  const accessPolicyStep = computed(() =>
    typeConfig.value.hasAddress ? 3 : 2,
  );

  /**
   * Get the step number for review based on type config
   *
   * @returns {number} Step number for review
   */
  const reviewStep = computed(() =>
    typeConfig.value.hasAddress ? 4 : 3,
  );

  /**
   * Check if Next button should be disabled
   */
  const isNextButtonDisabled = computed(() => {
    if (currentStep.value === 1) {
      // Step 1: Basic info - require name
      const { name } = formData.value;
      if (!name || name.length < NAME_MIN_LENGTH) return true;
      return false;
    }

    // Step 2 for types with address: all fields optional
    if (typeConfig.value.hasAddress && currentStep.value === 2) {
      return false;
    }

    // Access Policy step: V1 always internal, no validation needed
    if (currentStep.value === accessPolicyStep.value) {
      return false;
    }

    return false;
  });

  /**
   * Validate and change step
   *
   * @param {number} step - Target step number
   */
  async function changeStep(step: number): Promise<void> {
    logger.debug('changeStep called:', { currentStep: currentStep.value, targetStep: step });

    // Validate current step form before proceeding forward
    if (step > currentStep.value) {
      if (currentStep.value === 1 && step1FormRef.value) {
        const valid = await step1FormRef.value.validate();
        if (!valid) return;
      }

      if (typeConfig.value.hasAddress && currentStep.value === 2 && step2FormRef.value) {
        const valid = await step2FormRef.value.validate();
        if (!valid) return;
      }

      if (currentStep.value === accessPolicyStep.value && accessPolicyFormRef.value) {
        const valid = await accessPolicyFormRef.value.validate();
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
   *
   * @returns {Promise<void>}
   */
  async function submitForm(): Promise<void> {
    if (!apis.mapexOS?.organizations) {
      notifyFail({ message: t.errors.apiNotInitialized.value });
      return;
    }

    isSaving.value = true;

    try {
      // Build address only if type supports it and all required fields are set
      let address: Record<string, string> | undefined;
      if (typeConfig.value.hasAddress) {
        const addressData = formData.value.address;
        const hasValidAddress = addressData.city && addressData.state && addressData.country && addressData.zipCode;
        if (hasValidAddress) {
          address = {
            city: addressData.city,
            state: addressData.state,
            country: addressData.country,
            zipCode: addressData.zipCode,
          };
        }
      }

      // V1: AuthConfig is always internal - override regardless of form state
      const authConfig = {
        providerType: 'internal' as const,
      };

      // Build access policy
      const accessPolicy = {
        rolePolicy: formData.value.accessPolicy.rolePolicy,
        defaultScope: formData.value.accessPolicy.defaultScope,
      };

      if (isEditMode.value && organizationId.value) {
        // UPDATE existing organization
        const updatePayload: Record<string, any> = {
          name: formData.value.name,
          enabled: formData.value.enabled,
          authConfig,
          accessPolicy,
        };
        if (typeConfig.value.hasPhone && formData.value.phone) {
          updatePayload.phone = formData.value.phone;
        }
        if (address) {
          updatePayload.address = address;
        }

        logger.debug('Updating Organization:', { organizationId: organizationId.value, payload: updatePayload });

        await apis.mapexOS.organizations.update(
          { organizationId: organizationId.value },
          updatePayload as any,
        );

        notifySuccess({
          message: t.notifications.updated.value,
          timeout: 3000,
        });
      } else {
        // CREATE new organization
        const createPayload: Record<string, any> = {
          name: formData.value.name,
          type: orgType.value,
          enabled: formData.value.enabled,
          authConfig,
          accessPolicy,
        };
        if (parentOrgId.value) {
          createPayload.parentOrgId = parentOrgId.value;
        }
        if (typeConfig.value.hasPhone && formData.value.phone) {
          createPayload.phone = formData.value.phone;
        }
        if (address) {
          createPayload.address = address;
        }

        logger.debug('Creating Organization:', createPayload);

        const createdOrg = await apis.mapexOS.organizations.create(createPayload as any);

        // Add new organization to the store tree (no extra backend request needed)
        if (createdOrg?.id && createdOrg?.pathKey) {
          orgStore.addOrganizationToTree({
            id: createdOrg.id,
            name: createdOrg.name || formData.value.name,
            type: createdOrg.type || orgType.value,
            pathKey: createdOrg.pathKey,
          });
        }

        notifySuccess({
          message: t.notifications.created.value,
          timeout: 3000,
        });
      }

      // Navigate to customers list after success
      await router.push('/customers');
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
    isNextButtonDisabled,
    accessPolicyStep,
    reviewStep,
    changeStep,
    handleStepChange,
    submitForm,
  };
}
