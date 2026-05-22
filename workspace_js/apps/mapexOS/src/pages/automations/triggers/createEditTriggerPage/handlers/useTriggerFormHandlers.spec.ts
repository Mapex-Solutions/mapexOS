import { describe, it, expect, beforeEach, vi } from 'vitest';
import { ref, computed } from 'vue';
import type { Ref } from 'vue';

import type { Trigger, TriggerFormState } from '../interfaces';

/** Mock dependencies */
vi.mock('@utils/alert', () => ({
  notifySuccess: vi.fn(),
  notifyFail: vi.fn(),
}));

vi.mock('@utils/error', () => ({
  handleApiError: vi.fn(),
}));

/**
 * Mock @mapexos/schemas enums used by trigger constants
 */
vi.mock('@mapexos/schemas', () => ({
  TriggerTypeEnum: {
    HTTP: 'HTTP',
    MQTT: 'MQTT',
    RABBITMQ: 'RABBITMQ',
    NATS: 'NATS',
    WEBSOCKET: 'WEBSOCKET',
    EMAIL: 'EMAIL',
    TEAMS: 'TEAMS',
    SLACK: 'SLACK',
  },
  TriggerCategoryEnum: {
    TECHNICAL: 'TECHNICAL',
    COMMUNICATION: 'COMMUNICATION',
  },
}));

import { apis } from '@services/mapex';
import { notifySuccess, notifyFail } from '@utils/alert';
import { handleApiError } from '@utils/error';
import { STEP } from '../constants';

import { useTriggerFormHandlers } from './useTriggerFormHandlers';

/**
 * Factory for default Trigger data
 */
function createTriggerData(overrides: Partial<Trigger> = {}): Trigger {
  return {
    name: '',
    description: '',
    triggerType: 'HTTP' as any,
    category: 'TECHNICAL' as any,
    enabled: true,
    config: {},
    ...overrides,
  };
}

/**
 * Factory for default TriggerFormState
 */
function createFormState(overrides: Partial<TriggerFormState> = {}): TriggerFormState {
  return {
    selectedCategory: null,
    selectedType: null,
    isCreating: false,
    currentStep: 1,
    ...overrides,
  };
}

describe('useTriggerFormHandlers', () => {
  let triggerData: ReturnType<typeof ref<Trigger>>;
  let formState: ReturnType<typeof ref<TriggerFormState>>;
  let currentStep: ReturnType<typeof ref<number>>;
  let step3FormRef: ReturnType<typeof ref<any>>;
  let step4FormRef: ReturnType<typeof ref<any>>;
  let isEditMode: ReturnType<typeof ref<boolean>>;
  let triggerId: ReturnType<typeof ref<string | undefined>>;

  beforeEach(() => {
    vi.clearAllMocks();
    triggerData = ref(createTriggerData());
    formState = ref(createFormState());
    currentStep = ref(1);
    step3FormRef = ref(null);
    step4FormRef = ref(null);
    isEditMode = ref(false);
    triggerId = ref(undefined);

    // Setup API mocks
    (apis as any).triggers = {
      trigger: {
        create: vi.fn().mockResolvedValue({}),
        update: vi.fn().mockResolvedValue({}),
        getById: vi.fn().mockResolvedValue({
          name: 'Test Trigger',
          description: 'Desc',
          triggerType: 'HTTP',
          category: 'TECHNICAL',
          enabled: true,
          config: { HTTP: { url: 'https://example.com' } },
          id: 'trig-1',
        }),
      },
    };
  });

  /**
   * Mock translations object — mirrors the shape returned by
   * useCreateEditTriggerTranslations for the keys the handler uses.
   */
  const tMock = {
    validation: {
      selectCategory: { value: 'Please select a trigger category' },
      selectType: { value: 'Please select a trigger type' },
      fillRequired: { value: 'Please fill in all required fields' },
      completeConfig: { value: 'Please complete the configuration' },
    },
    notifications: {
      created: { value: 'Trigger created successfully!' },
      updated: { value: 'Trigger updated successfully!' },
    },
  };

  function setup() {
    return useTriggerFormHandlers({
      triggerData: triggerData as Ref<Trigger>,
      formState: formState as Ref<TriggerFormState>,
      currentStep: currentStep as Ref<number>,
      step3FormRef: computed(() => step3FormRef.value),
      step4FormRef: computed(() => step4FormRef.value),
      isEditMode: isEditMode as Ref<boolean>,
      triggerId,
      t: tMock,
    });
  }

  describe('isNextButtonDisabled', () => {
    it('returns true on category step when no category selected', () => {
      currentStep.value = STEP.CATEGORY;
      formState.value!.selectedCategory = null;
      const { isNextButtonDisabled } = setup();

      expect(isNextButtonDisabled.value).toBe(true);
    });

    it('returns false on category step when category is selected', () => {
      currentStep.value = STEP.CATEGORY;
      formState.value!.selectedCategory = 'TECHNICAL' as any;
      const { isNextButtonDisabled } = setup();

      expect(isNextButtonDisabled.value).toBe(false);
    });

    it('returns true on type step when no type selected', () => {
      currentStep.value = STEP.TYPE;
      formState.value!.selectedType = null;
      const { isNextButtonDisabled } = setup();

      expect(isNextButtonDisabled.value).toBe(true);
    });

    it('returns false on type step when type is selected', () => {
      currentStep.value = STEP.TYPE;
      formState.value!.selectedType = 'HTTP' as any;
      const { isNextButtonDisabled } = setup();

      expect(isNextButtonDisabled.value).toBe(false);
    });

    it('returns true on basic info step when name is empty', () => {
      currentStep.value = STEP.BASIC_INFO;
      triggerData.value!.name = '';
      const { isNextButtonDisabled } = setup();

      expect(isNextButtonDisabled.value).toBe(true);
    });

    it('returns false on basic info step when name is set', () => {
      currentStep.value = STEP.BASIC_INFO;
      triggerData.value!.name = 'My Trigger';
      const { isNextButtonDisabled } = setup();

      expect(isNextButtonDisabled.value).toBe(false);
    });

    it('returns false on configuration step', () => {
      currentStep.value = STEP.CONFIGURATION;
      const { isNextButtonDisabled } = setup();

      expect(isNextButtonDisabled.value).toBe(false);
    });
  });

  describe('onCategorySelected', () => {
    it('sets the selected category in form state', () => {
      const { onCategorySelected } = setup();

      onCategorySelected('TECHNICAL' as any);

      expect(formState.value!.selectedCategory).toBe('TECHNICAL');
    });
  });

  describe('onTypeSelected', () => {
    it('sets the selected type in form state', () => {
      const { onTypeSelected } = setup();

      onTypeSelected('MQTT' as any);

      expect(formState.value!.selectedType).toBe('MQTT');
    });
  });

  describe('changeStep', () => {
    it('validates category step before moving forward in create mode', async () => {
      currentStep.value = STEP.CATEGORY;
      formState.value!.selectedCategory = null;
      const { changeStep } = setup();

      await changeStep(STEP.TYPE);

      expect(currentStep.value).toBe(STEP.CATEGORY);
      expect(notifyFail).toHaveBeenCalledWith(
        expect.objectContaining({ message: 'Please select a trigger category' }),
      );
    });

    it('allows forward when category is selected', async () => {
      currentStep.value = STEP.CATEGORY;
      formState.value!.selectedCategory = 'TECHNICAL' as any;
      const { changeStep } = setup();

      await changeStep(STEP.TYPE);

      expect(currentStep.value).toBe(STEP.TYPE);
    });

    it('allows free navigation in edit mode', async () => {
      isEditMode.value = true;
      currentStep.value = STEP.CATEGORY;
      formState.value!.selectedCategory = null;
      const { changeStep } = setup();

      await changeStep(STEP.CONFIGURATION);

      expect(currentStep.value).toBe(STEP.CONFIGURATION);
    });

    it('clamps step to valid range', async () => {
      isEditMode.value = true;
      const { changeStep } = setup();

      await changeStep(99);

      // Should remain unchanged since 99 > TOTAL_STEPS
      expect(currentStep.value).toBe(1);
    });
  });

  describe('validateCurrentStep', () => {
    it('validates step 3 form ref when present', async () => {
      currentStep.value = STEP.BASIC_INFO;
      step3FormRef.value = { validate: vi.fn().mockResolvedValue(true) };
      const { validateCurrentStep } = setup();

      const result = await validateCurrentStep();

      expect(result).toBe(true);
      expect(step3FormRef.value.validate).toHaveBeenCalled();
    });

    it('returns false and notifies when step 3 form ref fails', async () => {
      currentStep.value = STEP.BASIC_INFO;
      step3FormRef.value = { validate: vi.fn().mockResolvedValue(false) };
      const { validateCurrentStep } = setup();

      const result = await validateCurrentStep();

      expect(result).toBe(false);
      expect(notifyFail).toHaveBeenCalled();
    });
  });

  describe('loadTriggerData', () => {
    it('returns null when not in edit mode', async () => {
      isEditMode.value = false;
      const { loadTriggerData } = setup();

      const result = await loadTriggerData();

      expect(result).toBeNull();
    });

    it('returns null when triggerId is undefined', async () => {
      isEditMode.value = true;
      triggerId.value = undefined;
      const { loadTriggerData } = setup();

      const result = await loadTriggerData();

      expect(result).toBeNull();
    });

    it('loads and transforms trigger data from API', async () => {
      isEditMode.value = true;
      triggerId.value = 'trig-1';
      const { loadTriggerData } = setup();

      const result = await loadTriggerData();

      expect(result).toBeDefined();
      expect(result!.name).toBe('Test Trigger');
      expect(result!.triggerType).toBe('HTTP');
      // Config should be unwrapped from { HTTP: { url: ... } } to { url: ... }
      expect(result!.config).toEqual({ url: 'https://example.com' });
    });

    it('calls handleApiError on failure', async () => {
      isEditMode.value = true;
      triggerId.value = 'trig-1';
      (apis as any).triggers.trigger.getById.mockRejectedValue(new Error('Network'));
      const { loadTriggerData } = setup();

      const result = await loadTriggerData();

      expect(result).toBeNull();
      expect(handleApiError).toHaveBeenCalled();
    });
  });

  describe('submitForm — CREATE mode', () => {
    it('calls create API with wrapped config', async () => {
      triggerData.value = createTriggerData({
        name: 'HTTP Trigger',
        triggerType: 'HTTP' as any,
        category: 'TECHNICAL' as any,
        config: { url: 'https://example.com', method: 'POST' },
      });
      formState.value!.selectedCategory = 'TECHNICAL' as any;
      formState.value!.selectedType = 'HTTP' as any;
      currentStep.value = STEP.CONFIGURATION;
      const { submitForm } = setup();

      await submitForm();

      expect(apis.triggers.trigger.create).toHaveBeenCalledWith(
        expect.objectContaining({
          name: 'HTTP Trigger',
          triggerType: 'HTTP',
          config: { HTTP: { url: 'https://example.com', method: 'POST' } },
        }),
      );
      expect(notifySuccess).toHaveBeenCalled();
      expect(formState.value!.isCreating).toBe(false);
    });

    it('handles API error on create failure', async () => {
      (apis as any).triggers.trigger.create.mockRejectedValue({ response: { status: 409 } });
      triggerData.value = createTriggerData({ name: 'Dup' });
      formState.value!.selectedCategory = 'TECHNICAL' as any;
      currentStep.value = STEP.CONFIGURATION;
      const { submitForm } = setup();

      await submitForm();

      expect(handleApiError).toHaveBeenCalled();
      expect(formState.value!.isCreating).toBe(false);
    });
  });

  describe('submitForm — EDIT mode', () => {
    it('calls update API with trigger ID', async () => {
      isEditMode.value = true;
      triggerId.value = 'trig-123';
      triggerData.value = createTriggerData({
        name: 'Updated Trigger',
        triggerType: 'MQTT' as any,
        category: 'TECHNICAL' as any,
        config: { broker: 'tcp://localhost:1883' },
      });
      formState.value!.selectedCategory = 'TECHNICAL' as any;
      formState.value!.selectedType = 'MQTT' as any;
      currentStep.value = STEP.CONFIGURATION;
      const { submitForm } = setup();

      await submitForm();

      expect(apis.triggers.trigger.update).toHaveBeenCalledWith(
        { triggerId: 'trig-123' },
        expect.objectContaining({
          name: 'Updated Trigger',
          config: { MQTT: { broker: 'tcp://localhost:1883' } },
        }),
      );
      expect(notifySuccess).toHaveBeenCalled();
    });
  });
});
