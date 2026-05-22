import type { Ref, ComputedRef } from 'vue';
import type { QForm } from 'quasar';
import type { WorkflowInstanceFormData, WorkflowInstanceFormState } from '../interfaces';
import type { DefinitionResponse } from '@mapexos/schemas';

import { computed } from 'vue';
import { useRouter } from 'vue-router';

import { apis } from '@services/mapex';
import { notifySuccess } from '@utils/alert/notify';
import { handleApiError } from '@utils/error';
import { useCreateEditWorkflowInstanceTranslations } from '@src/composables/i18n/pages/automations/workflowInstances/createEditWorkflowInstancePage/useCreateEditWorkflowInstanceTranslations';

import { TOTAL_STEPS, STEP } from '../constants';

interface UseWorkflowInstanceFormHandlersParams {
  formData: Ref<WorkflowInstanceFormData>;
  formState: Ref<WorkflowInstanceFormState>;
  currentStep: Ref<number>;
  isEditMode: Ref<boolean>;
  instanceId: Ref<string | undefined>;
  isSaving: Ref<boolean>;
  step1FormRef: ComputedRef<QForm | null>;
  step2FormRef: ComputedRef<QForm | null>;
  step3FormRef: ComputedRef<QForm | null>;
}

/**
 * Composable that encapsulates all form logic for the CreateEditWorkflowInstance page.
 * Handles step validation, navigation, definition selection, and form submission.
 *
 * @param {UseWorkflowInstanceFormHandlersParams} params - Refs to form state and UI elements
 * @returns Form handler functions and computed properties
 */
export function useWorkflowInstanceFormHandlers(params: UseWorkflowInstanceFormHandlersParams) {
  const {
    formData,
    formState,
    currentStep,
    isEditMode,
    instanceId,
    isSaving,
    step1FormRef,
    step2FormRef,
    step3FormRef,
  } = params;

  const router = useRouter();
  const t = useCreateEditWorkflowInstanceTranslations();

  /**
   * Handle definition selection from Step2
   * @param {DefinitionResponse} definition - The selected definition
   * @returns {void}
   */
  function onDefinitionSelected(definition: DefinitionResponse): void {
    formState.value.selectedDefinition = definition;
  }

  /**
   * Check if Next button should be disabled based on current step
   */
  const isNextButtonDisabled = computed(() => {
    if (currentStep.value === STEP.DEFINITION) {
      return !formData.value.definitionId;
    }
    return false;
  });

  /**
   * Get the form ref for the current step
   * @returns {QForm | null}
   */
  function getCurrentFormRef(): QForm | null {
    switch (currentStep.value) {
      case STEP.IDENTIFICATION: return step1FormRef.value;
      case STEP.DEFINITION: return step2FormRef.value;
      case STEP.EXTERNAL_INPUTS: return step3FormRef.value;
      default: return null;
    }
  }

  /**
   * Validate the current step before advancing
   * @returns {Promise<boolean>}
   */
  async function validateCurrentStep(): Promise<boolean> {
    const formRef = getCurrentFormRef();
    if (!formRef) return true;
    return await formRef.validate();
  }

  /**
   * Handle step change (next/previous/jump)
   * @param {number} newStep - The target step number
   * @returns {Promise<void>}
   */
  async function handleStepChange(newStep: number): Promise<void> {
    if (newStep < 1 || newStep > TOTAL_STEPS) return;

    // Validate before advancing (not when going back)
    if (newStep > currentStep.value) {
      const isValid = await validateCurrentStep();
      if (!isValid) return;
    }

    currentStep.value = newStep;
    formState.value.currentStep = newStep;
  }

  /**
   * Change step (alias for handleStepChange, used by StepperVertical)
   * @param {number} step - Target step
   * @returns {void}
   */
  function changeStep(step: number): void {
    void handleStepChange(step);
  }

  /**
   * Submit the form (create or update)
   * @returns {Promise<void>}
   */
  async function submitForm(): Promise<void> {
    isSaving.value = true;
    try {
      const payload = {
        name: formData.value.name,
        description: formData.value.description || '',
        definitionId: formData.value.definitionId!,
        definitionVersion: formData.value.definitionVersion,
        definitionName: formData.value.selectedDefinition?.name || '',
        pathKey: '',
        externalInputs: formData.value.externalInputs,
        isTemplate: formData.value.isTemplate,
        uniqueExecution: formData.value.uniqueExecution,
        workflowUUID: formData.value.uniqueExecution ? formData.value.workflowUUID : '',
      };

      if (isEditMode.value && instanceId.value) {
        await apis.workflows.instance.update(
          { instanceId: instanceId.value },
          {
            name: payload.name,
            description: payload.description,
            externalInputs: payload.externalInputs,
            isTemplate: payload.isTemplate,
            uniqueExecution: payload.uniqueExecution,
            workflowUUID: payload.workflowUUID,
          },
        );
        notifySuccess({ message: t.notifications.updated.value });
      } else {
        await apis.workflows.instance.create(payload);
        notifySuccess({ message: t.notifications.created.value });
      }

      await router.push('/workflow_instances');
    } catch (error: any) {
      handleApiError(error, {
        defaultMessage: isEditMode.value
          ? t.notifications.updateFailed.value
          : t.notifications.creationFailed.value,
      });
    } finally {
      isSaving.value = false;
    }
  }

  return {
    onDefinitionSelected,
    isNextButtonDisabled,
    handleStepChange,
    changeStep,
    submitForm,
  };
}
