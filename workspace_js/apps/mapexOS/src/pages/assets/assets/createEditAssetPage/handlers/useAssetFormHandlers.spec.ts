import { describe, it, expect, beforeEach, vi } from 'vitest';
import { ref } from 'vue';
import type { Ref } from 'vue';

import type { AssetFormData, AssetFormState } from '../interfaces';
import { INITIAL_LORAWAN_CONFIG } from '../constants';

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
    steps: {
      groups: {
        setup: { value: 'Setup' },
        identification: { value: 'Identification' },
        connectivity: { value: 'Connectivity' },
        routing: { value: 'Routing' },
        finalization: { value: 'Finalization' },
      },
      stepType: { label: { value: 'type' }, description: { value: '' } },
      step1: {
        label: { value: 'identification' },
        description: { value: '' },
        leafLabel: { value: 'Information' },
        leafDescription: { value: '' },
        attributesLabel: { value: 'Attributes' },
        attributesDescription: { value: '' },
      },
      step2: { label: { value: 'assetTemplate' }, description: { value: '' } },
      step3: { label: { value: 'routeGroups' }, description: { value: '' } },
      step4: { label: { value: 'connectivity' }, description: { value: '' } },
      step5: { label: { value: 'health' }, description: { value: '' } },
      step6: { label: { value: 'review' }, description: { value: '' } },
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
    attributes: [],
    protocol: 'HTTP',
    latitude: null,
    longitude: null,
    mqttConfig: { clientId: '', username: '', authType: 'cert' as const, password: '' },
    lorawanConfig: { ...INITIAL_LORAWAN_CONFIG, gateway: { ...INITIAL_LORAWAN_CONFIG.gateway } },
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

  function gatewayLorawanConfig() {
    return { ...INITIAL_LORAWAN_CONFIG, kind: 'gateway' as const, gateway: { ...INITIAL_LORAWAN_CONFIG.gateway } };
  }

  describe('visibleSteps', () => {
    it('includes all seven steps for a non-gateway asset', () => {
      assetData.value = createAssetFormData({ protocol: 'MQTT' });
      const { visibleSteps } = setup();

      expect(visibleSteps.value.map((s) => s.id)).toEqual([
        'type', 'assetTemplate', 'identification', 'attributes', 'connectivity', 'health', 'routeGroups', 'review',
      ]);
    });

    it('skips assetTemplate + routeGroups for a LoRaWAN gateway', () => {
      assetData.value = createAssetFormData({ protocol: 'LORAWAN', lorawanConfig: gatewayLorawanConfig() });
      const { visibleSteps } = setup();

      expect(visibleSteps.value.map((s) => s.id)).toEqual([
        'type', 'identification', 'attributes', 'connectivity', 'health', 'review',
      ]);
    });

    it('re-expands when switching from gateway to device', () => {
      assetData.value = createAssetFormData({ protocol: 'LORAWAN', lorawanConfig: gatewayLorawanConfig() });
      const { visibleSteps } = setup();

      expect(visibleSteps.value).toHaveLength(6);

      assetData.value.lorawanConfig.kind = 'device';

      expect(visibleSteps.value).toHaveLength(8);
      expect(visibleSteps.value.map((s) => s.id)).toContain('assetTemplate');
    });
  });

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

    it('returns true on the template step when no template selected', () => {
      // Order: TYPE(1), ASSET_TEMPLATE(2), IDENTIFICATION(3), ATTRIBUTES(4), CONNECTIVITY(5), HEALTH(6), ROUTE_GROUPS(7), ...
      currentStep.value = 2;
      assetData.value!.assetTemplateId = null;
      const { isNextButtonDisabled } = setup();

      expect(isNextButtonDisabled.value).toBe(true);
    });

    it('returns false on the template step when template selected', () => {
      currentStep.value = 2;
      assetData.value!.assetTemplateId = 'tpl-1';
      const { isNextButtonDisabled } = setup();

      expect(isNextButtonDisabled.value).toBe(false);
    });

    it('returns true on the route-groups step when none selected', () => {
      // ROUTE_GROUPS is position 7 in the non-gateway flow.
      currentStep.value = 7;
      assetData.value!.routeGroupIds = [];
      const { isNextButtonDisabled } = setup();

      expect(isNextButtonDisabled.value).toBe(true);
    });

    it('returns false on the route-groups step when selected', () => {
      currentStep.value = 7;
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
      // Position 3 is IDENTIFICATION, whose form is step1FormRef.
      step1FormRef.value = { validate: vi.fn().mockResolvedValue(false) };
      currentStep.value = 3;
      const { changeStep } = setup();

      await changeStep(4);

      expect(currentStep.value).toBe(3);
    });

    it('allows step change when form validation passes', async () => {
      step1FormRef.value = { validate: vi.fn().mockResolvedValue(true) };
      currentStep.value = 3;
      const { changeStep } = setup();

      await changeStep(4);

      expect(currentStep.value).toBe(4);
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

  describe('submitForm — LoRaWAN gateway key mode', () => {
    function setupCreateMock() {
      (apis.assets as any) = {
        asset: { create: vi.fn().mockResolvedValue({ id: 'new-id' }), update: vi.fn() },
      };
    }

    function gatewayConfig(gw: Partial<typeof INITIAL_LORAWAN_CONFIG.gateway>) {
      return {
        ...INITIAL_LORAWAN_CONFIG,
        kind: 'gateway' as const,
        gateway: { ...INITIAL_LORAWAN_CONFIG.gateway, frequencyPlanId: 'EU_863_870', ...gw },
      };
    }

    function sentGateway() {
      return (apis.assets.asset.create as any).mock.calls[0][0].protocol.lorawan.gateway;
    }

    it('sends apiKey when authMode=key and the token is filled', async () => {
      setupCreateMock();
      assetData.value = createAssetFormData({
        assetTemplateId: 'tpl-1',
        protocol: 'LORAWAN',
        lorawanConfig: gatewayConfig({ authMode: 'key', apiKey: 'ABCDEF0123' }),
      });
      const { submitForm } = setup();

      await submitForm();

      expect(sentGateway().apiKey).toBe('ABCDEF0123');
    });

    it('omits apiKey when authMode=key but the token is blank', async () => {
      setupCreateMock();
      assetData.value = createAssetFormData({
        assetTemplateId: 'tpl-1',
        protocol: 'LORAWAN',
        lorawanConfig: gatewayConfig({ authMode: 'key', apiKey: '' }),
      });
      const { submitForm } = setup();

      await submitForm();

      expect(sentGateway()).not.toHaveProperty('apiKey');
    });

    it('omits apiKey for cert and eui modes', async () => {
      setupCreateMock();
      assetData.value = createAssetFormData({
        assetTemplateId: 'tpl-1',
        protocol: 'LORAWAN',
        lorawanConfig: gatewayConfig({ authMode: 'cert', apiKey: 'should-not-send' }),
      });
      const { submitForm } = setup();

      await submitForm();

      expect(sentGateway()).not.toHaveProperty('apiKey');
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
