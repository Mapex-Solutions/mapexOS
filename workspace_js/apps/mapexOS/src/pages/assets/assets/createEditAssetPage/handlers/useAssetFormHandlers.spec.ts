import { describe, it, expect, beforeEach, vi } from 'vitest';
import { ref } from 'vue';
import type { Ref } from 'vue';

import type { AssetFormData, AssetFormState } from '../interfaces';

/** Mock dependencies */
vi.mock('@utils/alert/notify', () => ({
  notifySuccess: vi.fn(),
  notifyFail: vi.fn(),
}));

vi.mock('@utils/error', () => ({
  handleApiError: vi.fn(),
}));

vi.mock('@composables/useLogger', () => ({
  useLogger: () => ({
    debug: vi.fn(),
    warn: vi.fn(),
    error: vi.fn(),
  }),
}));

vi.mock('@src/composables/i18n/pages/assets/addAsset/useAddAssetTranslations', () => ({
  useAddAssetTranslations: () => ({
    notifications: {
      created: { value: 'Asset created' },
      updated: { value: 'Asset updated' },
      alreadyExists: 'Already exists',
      validationFailed: 'Validation failed',
      networkError: 'Network error',
      creationFailed: 'Creation failed',
      updateFailed: { value: 'Update failed' },
    },
  }),
}));

vi.mock('@stores/organization', () => ({
  useOrganizationStore: () => ({
    selectedOrganizationId: 'org-123',
  }),
}));

import { apis } from '@services/mapex';
import { notifySuccess } from '@utils/alert/notify';
import { handleApiError } from '@utils/error';

import { useAssetFormHandlers } from './useAssetFormHandlers';

/**
 * Factory for default AssetFormData
 */
function createAssetFormData(overrides: Partial<AssetFormData> = {}): AssetFormData {
  return {
    name: '',
    assetId: '',
    enabled: true,
    description: '',
    assetTemplateId: null,
    routeGroupIds: [],
    protocol: 'HTTP',
    latitude: null,
    longitude: null,
    mqttConfig: { clientId: '', username: '', authType: 'cert' as const, password: '' },
    debugEnabled: false,
    healthMonitor: {
      enabled: false,
      thresholdMinutes: 10,
      requiredMisses: 3,
      offlineRouteGroupIds: [],
      onlineRouteGroupIds: [],
      selectedOfflineRouteGroups: [],
      selectedOnlineRouteGroups: [],
    },
    ...overrides,
  };
}

function createFormState(overrides: Partial<AssetFormState> = {}): AssetFormState {
  return {
    selectedTemplate: null,
    selectedRouteGroups: [],
    isCreating: false,
    currentStep: 1,
    ...overrides,
  };
}

describe('useAssetFormHandlers', () => {
  let assetData: ReturnType<typeof ref<AssetFormData>>;
  let formState: ReturnType<typeof ref<AssetFormState>>;
  let currentStep: ReturnType<typeof ref<number>>;
  let isEditMode: ReturnType<typeof ref<boolean>>;
  let assetId: ReturnType<typeof ref<string | undefined>>;
  let isSaving: ReturnType<typeof ref<boolean>>;
  let step1FormRef: ReturnType<typeof ref<any>>;
  let step2FormRef: ReturnType<typeof ref<any>>;
  let step3FormRef: ReturnType<typeof ref<any>>;
  let step4FormRef: ReturnType<typeof ref<any>>;

  beforeEach(() => {
    vi.clearAllMocks();
    assetData = ref(createAssetFormData());
    formState = ref(createFormState());
    currentStep = ref(1);
    isEditMode = ref(false);
    assetId = ref(undefined);
    isSaving = ref(false);
    step1FormRef = ref(null);
    step2FormRef = ref(null);
    step3FormRef = ref(null);
    step4FormRef = ref(null);
  });

  function setup() {
    return useAssetFormHandlers({
      assetData: assetData as Ref<AssetFormData>,
      formState: formState as Ref<AssetFormState>,
      currentStep: currentStep as Ref<number>,
      isEditMode: isEditMode as Ref<boolean>,
      assetId,
      isSaving: isSaving as Ref<boolean>,
      step1FormRef,
      step2FormRef,
      step3FormRef,
      step4FormRef,
    });
  }

  describe('onTemplateSelected', () => {
    it('sets the selected template in form state', () => {
      const { onTemplateSelected } = setup();
      const template = { id: 'tpl-1', name: 'Test' } as any;

      onTemplateSelected(template);

      expect(formState.value!.selectedTemplate).toEqual(template);
    });

    it('handles null template selection', () => {
      const { onTemplateSelected } = setup();

      onTemplateSelected(null);

      expect(formState.value!.selectedTemplate).toBeNull();
    });
  });

  describe('onRouteGroupsSelected', () => {
    it('sets the selected route groups in form state', () => {
      const { onRouteGroupsSelected } = setup();
      const routeGroups = [{ id: 'rg-1', name: 'RG1' }] as any;

      onRouteGroupsSelected(routeGroups);

      expect(formState.value!.selectedRouteGroups).toEqual(routeGroups);
    });
  });

  describe('isNextButtonDisabled', () => {
    it('returns false for step 1 (form validates itself)', () => {
      currentStep.value = 1;
      const { isNextButtonDisabled } = setup();

      expect(isNextButtonDisabled.value).toBe(false);
    });

    it('returns true for step 2 when no template selected', () => {
      currentStep.value = 2;
      assetData.value!.assetTemplateId = null;
      const { isNextButtonDisabled } = setup();

      expect(isNextButtonDisabled.value).toBe(true);
    });

    it('returns false for step 2 when template selected', () => {
      currentStep.value = 2;
      assetData.value!.assetTemplateId = 'tpl-1';
      const { isNextButtonDisabled } = setup();

      expect(isNextButtonDisabled.value).toBe(false);
    });

    it('returns true for step 3 when no route groups selected', () => {
      currentStep.value = 3;
      assetData.value!.routeGroupIds = [];
      const { isNextButtonDisabled } = setup();

      expect(isNextButtonDisabled.value).toBe(true);
    });

    it('returns false for step 3 when route groups selected', () => {
      currentStep.value = 3;
      assetData.value!.routeGroupIds = ['rg-1'];
      const { isNextButtonDisabled } = setup();

      expect(isNextButtonDisabled.value).toBe(false);
    });
  });

  describe('changeStep', () => {
    it('changes step when no form ref to validate', async () => {
      const { changeStep } = setup();

      await changeStep(2);

      expect(currentStep.value).toBe(2);
    });

    it('blocks step change when form validation fails', async () => {
      step1FormRef.value = { validate: vi.fn().mockResolvedValue(false) };
      currentStep.value = 1;
      const { changeStep } = setup();

      await changeStep(2);

      expect(currentStep.value).toBe(1);
    });

    it('allows step change when form validation passes', async () => {
      step1FormRef.value = { validate: vi.fn().mockResolvedValue(true) };
      currentStep.value = 1;
      const { changeStep } = setup();

      await changeStep(2);

      expect(currentStep.value).toBe(2);
    });
  });

  describe('submitForm — CREATE mode', () => {
    it('calls create API and navigates on success', async () => {
      assetData.value = createAssetFormData({
        name: 'My Asset',
        assetId: 'asset-uuid',
        assetTemplateId: 'tpl-1',
        routeGroupIds: ['rg-1'],
        protocol: 'HTTP',
      });
      (apis.assets as any) = {
        asset: {
          create: vi.fn().mockResolvedValue({ id: 'new-id' }),
          update: vi.fn(),
        },
      };
      const { submitForm } = setup();

      await submitForm();

      expect(apis.assets.asset.create).toHaveBeenCalledWith(
        expect.objectContaining({
          name: 'My Asset',
          assetTemplateId: 'tpl-1',
          routeGroupIds: ['rg-1'],
          orgId: 'org-123',
        }),
      );
      expect(notifySuccess).toHaveBeenCalled();
      expect(isSaving.value).toBe(false);
    });

    it('calls handleApiError on failure', async () => {
      (apis.assets as any) = {
        asset: {
          create: vi.fn().mockRejectedValue({ response: { status: 500 } }),
          update: vi.fn(),
        },
      };
      assetData.value = createAssetFormData({ assetTemplateId: 'tpl-1' });
      const { submitForm } = setup();

      await submitForm();

      expect(handleApiError).toHaveBeenCalled();
      expect(isSaving.value).toBe(false);
    });
  });

  describe('submitForm — EDIT mode', () => {
    it('calls update API with the asset ID', async () => {
      isEditMode.value = true;
      assetId.value = 'asset-123';
      assetData.value = createAssetFormData({
        name: 'Updated Asset',
        assetTemplateId: 'tpl-1',
        routeGroupIds: ['rg-1'],
        protocol: 'HTTP',
      });
      (apis.assets as any) = {
        asset: {
          create: vi.fn(),
          update: vi.fn().mockResolvedValue({}),
        },
      };
      const { submitForm } = setup();

      await submitForm();

      expect(apis.assets.asset.update).toHaveBeenCalledWith(
        { assetId: 'asset-123' },
        expect.objectContaining({ name: 'Updated Asset' }),
      );
      expect(notifySuccess).toHaveBeenCalled();
    });

    it('does not include orgId in UPDATE payload', async () => {
      isEditMode.value = true;
      assetId.value = 'asset-123';
      assetData.value = createAssetFormData({ assetTemplateId: 'tpl-1' });
      (apis.assets as any) = {
        asset: {
          create: vi.fn(),
          update: vi.fn().mockResolvedValue({}),
        },
      };
      const { submitForm } = setup();

      await submitForm();

      const updatePayload = (apis.assets.asset.update as any).mock.calls[0][1];
      expect(updatePayload.orgId).toBeUndefined();
    });
  });
});
