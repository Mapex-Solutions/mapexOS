/** TYPE IMPORTS */
import type { Ref, ComputedRef } from 'vue';
import type { QForm } from 'quasar';
import type { Trigger, TriggerFormState, TriggerCategory, TriggerType } from '../interfaces';

/** VUE IMPORTS */
import { computed } from 'vue';
import { useRouter } from 'vue-router';

/** SERVICES */
import { apis } from '@services/mapex';

/** UTILS */
import { notifySuccess, notifyFail } from '@utils/alert';
import { handleApiError } from '@utils/error';

/** LOCAL IMPORTS */
import { TOTAL_STEPS, STEP } from '../constants';

/**
 * Handler composable parameters
 */
interface UseTriggerFormHandlersParams {
  triggerData: Ref<Trigger>;
  formState: Ref<TriggerFormState>;
  currentStep: Ref<number>;
  step3FormRef: ComputedRef<QForm | null>;
  step4FormRef: ComputedRef<QForm | null>;
  isEditMode: Ref<boolean>;
  triggerId: Ref<string | undefined>;
  /**
   * Translation composable instance (from useCreateEditTriggerTranslations)
   * Used for localized validation and notification messages
   */
  t: any;
}

/**
 * Trigger form handlers composable
 * Manages form navigation, validation, and submission for CREATE and EDIT modes
 * @param {UseTriggerFormHandlersParams} params - Handler parameters
 * @returns Handlers and computed properties
 */
export function useTriggerFormHandlers(params: UseTriggerFormHandlersParams) {
  const {
    triggerData,
    formState,
    currentStep,
    step3FormRef,
    step4FormRef,
    isEditMode,
    triggerId,
    t,
  } = params;
  const router = useRouter();

  /**
   * Check if next button should be disabled
   */
  const isNextButtonDisabled = computed(() => {
    switch (currentStep.value) {
      case STEP.CATEGORY:
        return !formState.value.selectedCategory;
      case STEP.TYPE:
        return !formState.value.selectedType;
      case STEP.BASIC_INFO:
        return !triggerData.value.name;
      case STEP.CONFIGURATION:
        return false;
      default:
        return false;
    }
  });

  /**
   * Handle category selection from Step 1
   * @param {TriggerCategory} category - Selected category
   * @returns {void}
   */
  function onCategorySelected(category: TriggerCategory): void {
    formState.value.selectedCategory = category;
  }

  /**
   * Handle trigger type selection from Step 2
   * @param {TriggerType} type - Selected trigger type
   * @returns {void}
   */
  function onTypeSelected(type: TriggerType): void {
    formState.value.selectedType = type;
  }

  /**
   * Change step with validation
   * In EDIT mode, allows free navigation without validation
   * In CREATE mode, validates current step before moving forward
   * @param {number} targetStep - Target step number
   * @returns {Promise<void>}
   */
  async function changeStep(targetStep: number): Promise<void> {
    // In EDIT mode, allow free navigation without validation
    // In CREATE mode, validate current step before moving forward
    if (!isEditMode.value && targetStep > currentStep.value) {
      const isValid = await validateCurrentStep();
      if (!isValid) {
        return;
      }
    }

    // Update step
    if (targetStep >= 1 && targetStep <= TOTAL_STEPS) {
      currentStep.value = targetStep;
      formState.value.currentStep = targetStep;
    }
  }

  /**
   * Handle step change from stepper click
   * @param {number} step - Clicked step number
   * @returns {void}
   */
  function handleStepChange(step: number): void {
    void changeStep(step);
  }

  /**
   * Validate current step
   * @returns {Promise<boolean>} Validation result
   */
  async function validateCurrentStep(): Promise<boolean> {
    switch (currentStep.value) {
      case STEP.CATEGORY:
        if (!formState.value.selectedCategory) {
          notifyFail({ message: t.validation.selectCategory.value });
          return false;
        }
        return true;

      case STEP.TYPE:
        if (!formState.value.selectedType) {
          notifyFail({ message: t.validation.selectType.value });
          return false;
        }
        return true;

      case STEP.BASIC_INFO:
        if (step3FormRef.value) {
          const isValid = await step3FormRef.value.validate();
          if (!isValid) {
            notifyFail({ message: t.validation.fillRequired.value });
            return false;
          }
        }
        return true;

      case STEP.CONFIGURATION:
        if (step4FormRef.value) {
          const isValid = await step4FormRef.value.validate();
          if (!isValid) {
            notifyFail({ message: t.validation.completeConfig.value });
            return false;
          }
        }
        return true;

      default:
        return true;
    }
  }

  /**
   * Load trigger data from API in EDIT mode
   * @returns {Promise<Trigger | null>} Loaded trigger data or null on error
   */
  async function loadTriggerData(): Promise<Trigger | null> {
    if (!isEditMode.value || !triggerId.value) return null;

    try {
      if (!apis.triggers) {
        throw new Error('Triggers API is not initialized');
      }

      const response = await apis.triggers.trigger.getById({ triggerId: triggerId.value });

      // Extract config from wrapped format if needed
      const rawConfig = response.config as Record<string, unknown>;
      let extractedConfig = rawConfig;

      // If config is wrapped like { HTTP: { ... } }, extract the inner config
      if (response.triggerType && rawConfig[response.triggerType]) {
        extractedConfig = rawConfig[response.triggerType] as Record<string, unknown>;
      }

      const trigger: Trigger = {
        name: response.name ?? '',
        description: response.description || '',
        triggerType: response.triggerType as TriggerType,
        category: response.category as TriggerCategory,
        enabled: response.enabled ?? true,
        config: extractedConfig as Record<string, any>,
      };

      // Add optional properties if they exist
      if (response.id) trigger.id = response.id;
      if (response.isTemplate !== undefined) trigger.isTemplate = response.isTemplate;
      if (response.created) trigger.createdAt = response.created;
      if (response.updated) trigger.updatedAt = response.updated;

      return trigger;
    } catch (error) {
      handleApiError(error, {
        defaultMessage: 'Failed to load trigger data',
        timeout: 5000,
      });
      return null;
    }
  }

  /**
   * Submit form - CREATE or UPDATE based on mode
   * Transforms frontend data to match backend schema and saves to database
   * @returns {Promise<void>}
   */
  async function submitForm(): Promise<void> {
    formState.value.isCreating = true;

    try {
      // Final validation
      const isValid = await validateCurrentStep();
      if (!isValid) {
        formState.value.isCreating = false;
        return;
      }

      // Wrap config in the correct union type structure based on triggerType
      const wrappedConfig = {
        [triggerData.value.triggerType]: triggerData.value.config,
      };

      // Build payload matching TriggerCreate/TriggerUpdate schema
      const payload = {
        name: triggerData.value.name,
        description: triggerData.value.description || undefined,
        triggerType: triggerData.value.triggerType,
        category: triggerData.value.category,
        enabled: triggerData.value.enabled,
        isSystem: false,
        isTemplate: triggerData.value.isTemplate ?? false,
        config: wrappedConfig,
      };

      if (!apis.triggers) {
        throw new Error('Triggers API is not initialized');
      }

      if (isEditMode.value && triggerId.value) {
        // UPDATE existing trigger
        await apis.triggers.trigger.update({ triggerId: triggerId.value }, payload);
        notifySuccess({ message: t.notifications.updated.value });
      } else {
        // CREATE new trigger
        await apis.triggers.trigger.create(payload);
        notifySuccess({ message: t.notifications.created.value });
      }

      // Navigate back to triggers list
      void router.push('/triggers');
    } catch (error) {
      const action = isEditMode.value ? 'update' : 'create';
      handleApiError(error, {
        customMessages: {
          409: 'Trigger with this name already exists',
          422: 'Invalid trigger configuration. Please check all fields.',
          network: 'Network error. Please check your connection.',
        },
        defaultMessage: `Failed to ${action} trigger. Please try again.`,
        timeout: 5000,
      });
    } finally {
      formState.value.isCreating = false;
    }
  }

  return {
    isNextButtonDisabled,
    onCategorySelected,
    onTypeSelected,
    changeStep,
    handleStepChange,
    validateCurrentStep,
    loadTriggerData,
    submitForm,
  };
}
