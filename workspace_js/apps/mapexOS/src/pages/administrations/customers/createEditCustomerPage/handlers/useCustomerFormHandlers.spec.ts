import { describe, it, expect, beforeEach, vi } from 'vitest';
import { ref, computed } from 'vue';
import type { Ref, ComputedRef } from 'vue';

import type { OrganizationFormData, OrganizationType, OrgTypeConfig } from '../interfaces';

import { INITIAL_ORGANIZATION_FORM_DATA, ORG_TYPE_CONFIG } from '../constants';

/** Mock dependencies */
vi.mock('@utils/alert/notify', () => ({
  notifySuccess: vi.fn(),
  notifyFail: vi.fn(),
}));

vi.mock('@composables/useLogger', () => ({
  useLogger: () => ({
    debug: vi.fn(),
    warn: vi.fn(),
    error: vi.fn(),
  }),
}));

vi.mock('@stores/organization', () => ({
  useOrganizationStore: () => ({
    selectedOrganizationId: 'org-123',
    addOrganizationToTree: vi.fn(),
  }),
}));

import { apis } from '@services/mapex';
import { notifySuccess, notifyFail } from '@utils/alert/notify';

import { useCustomerFormHandlers } from './useCustomerFormHandlers';

/**
 * Creates a mock translations object
 */
function createMockTranslations() {
  return {
    notifications: {
      created: { value: 'Organization created' },
      updated: { value: 'Organization updated' },
      createFailed: { value: 'Create failed' },
      updateFailed: { value: 'Update failed' },
      alreadyExists: { value: 'Already exists' },
      forbidden: { value: 'Forbidden' },
    },
  };
}

describe('useCustomerFormHandlers', () => {
  let formData: ReturnType<typeof ref<OrganizationFormData>>;
  let currentStep: ReturnType<typeof ref<number>>;
  let isEditMode: ReturnType<typeof ref<boolean>>;
  let organizationId: ReturnType<typeof ref<string | undefined>>;
  let isSaving: ReturnType<typeof ref<boolean>>;
  let orgTypeRef: ReturnType<typeof ref<OrganizationType>>;
  let parentOrgIdRef: ReturnType<typeof ref<string | undefined>>;
  let typeConfigRef: ReturnType<typeof ref<OrgTypeConfig>>;
  let step1FormRef: ReturnType<typeof ref<any>>;
  let step2FormRef: ReturnType<typeof ref<any>>;
  let accessPolicyFormRef: ReturnType<typeof ref<any>>;
  let t: ReturnType<typeof createMockTranslations>;

  beforeEach(() => {
    vi.clearAllMocks();
    formData = ref({ ...INITIAL_ORGANIZATION_FORM_DATA, address: { ...INITIAL_ORGANIZATION_FORM_DATA.address } });
    currentStep = ref(1);
    isEditMode = ref(false);
    organizationId = ref(undefined);
    isSaving = ref(false);
    orgTypeRef = ref<OrganizationType>('customer');
    parentOrgIdRef = ref<string | undefined>('parent-org-123');
    typeConfigRef = ref<OrgTypeConfig>(ORG_TYPE_CONFIG.customer);
    step1FormRef = ref(null);
    step2FormRef = ref(null);
    accessPolicyFormRef = ref(null);
    t = createMockTranslations();

    // Setup API mocks
    (apis.mapexOS as any).organizations = {
      ...apis.mapexOS.organizations,
      create: vi.fn().mockResolvedValue({ id: 'new-org', name: 'Test', type: 'customer', pathKey: '/v/c' }),
      update: vi.fn().mockResolvedValue({}),
    };
  });

  function setup() {
    return useCustomerFormHandlers({
      formData: formData as Ref<OrganizationFormData>,
      currentStep: currentStep as Ref<number>,
      isEditMode: isEditMode as Ref<boolean>,
      organizationId,
      isSaving: isSaving as Ref<boolean>,
      orgType: computed(() => orgTypeRef.value) as ComputedRef<OrganizationType>,
      parentOrgId: computed(() => parentOrgIdRef.value),
      typeConfig: computed(() => typeConfigRef.value) as ComputedRef<OrgTypeConfig>,
      step1FormRef,
      step2FormRef,
      accessPolicyFormRef,
      t,
    });
  }

  describe('accessPolicyStep and reviewStep', () => {
    it('returns step 3 for access policy when type has address', () => {
      typeConfigRef.value = ORG_TYPE_CONFIG.customer; // hasAddress: true
      const { accessPolicyStep, reviewStep } = setup();

      expect(accessPolicyStep.value).toBe(3);
      expect(reviewStep.value).toBe(4);
    });

    it('returns step 2 for access policy when type has no address', () => {
      typeConfigRef.value = ORG_TYPE_CONFIG.building; // hasAddress: false
      const { accessPolicyStep, reviewStep } = setup();

      expect(accessPolicyStep.value).toBe(2);
      expect(reviewStep.value).toBe(3);
    });
  });

  describe('isNextButtonDisabled', () => {
    it('returns true on step 1 when name is too short', () => {
      currentStep.value = 1;
      formData.value!.name = 'AB';
      const { isNextButtonDisabled } = setup();

      expect(isNextButtonDisabled.value).toBe(true);
    });

    it('returns false on step 1 when name is valid', () => {
      currentStep.value = 1;
      formData.value!.name = 'Customer Corp';
      const { isNextButtonDisabled } = setup();

      expect(isNextButtonDisabled.value).toBe(false);
    });

    it('returns false for address step (all fields optional)', () => {
      typeConfigRef.value = ORG_TYPE_CONFIG.customer;
      currentStep.value = 2;
      const { isNextButtonDisabled } = setup();

      expect(isNextButtonDisabled.value).toBe(false);
    });

    it('returns false for access policy step', () => {
      typeConfigRef.value = ORG_TYPE_CONFIG.customer;
      currentStep.value = 3;
      const { isNextButtonDisabled } = setup();

      expect(isNextButtonDisabled.value).toBe(false);
    });
  });

  describe('changeStep', () => {
    it('changes step when no form ref exists', async () => {
      const { changeStep } = setup();

      await changeStep(2);

      expect(currentStep.value).toBe(2);
    });

    it('blocks forward when step 1 form validation fails', async () => {
      step1FormRef.value = { validate: vi.fn().mockResolvedValue(false) };
      currentStep.value = 1;
      const { changeStep } = setup();

      await changeStep(2);

      expect(currentStep.value).toBe(1);
    });

    it('allows backward navigation without validation', async () => {
      currentStep.value = 3;
      const { changeStep } = setup();

      await changeStep(1);

      expect(currentStep.value).toBe(1);
    });
  });

  describe('submitForm — CREATE mode', () => {
    it('calls create API with correct payload', async () => {
      formData.value!.name = 'New Customer';
      formData.value!.enabled = true;
      formData.value!.address = { city: 'NYC', state: 'NY', country: 'US', zipCode: '10001' };
      const { submitForm } = setup();

      await submitForm();

      expect(apis.mapexOS.organizations.create).toHaveBeenCalledWith(
        expect.objectContaining({
          name: 'New Customer',
          type: 'customer',
          parentOrgId: 'parent-org-123',
          authConfig: { providerType: 'internal' },
          address: expect.objectContaining({ city: 'NYC' }),
        }),
      );
      expect(notifySuccess).toHaveBeenCalled();
      expect(isSaving.value).toBe(false);
    });

    it('omits address for types without address', async () => {
      typeConfigRef.value = ORG_TYPE_CONFIG.building;
      formData.value!.name = 'Building A';
      const { submitForm } = setup();

      await submitForm();

      const payload = (apis.mapexOS.organizations.create as any).mock.calls[0][0];
      expect(payload.address).toBeUndefined();
    });

    it('handles 409 conflict error', async () => {
      (apis.mapexOS as any).organizations.create.mockRejectedValue({ response: { status: 409 } });
      formData.value!.name = 'Duplicate';
      const { submitForm } = setup();

      await submitForm();

      expect(notifyFail).toHaveBeenCalledWith(
        expect.objectContaining({ message: 'Already exists' }),
      );
    });
  });

  describe('submitForm — EDIT mode', () => {
    beforeEach(() => {
      isEditMode.value = true;
      organizationId.value = 'org-456';
    });

    it('calls update API with organization ID', async () => {
      formData.value!.name = 'Updated Customer';
      const { submitForm } = setup();

      await submitForm();

      expect(apis.mapexOS.organizations.update).toHaveBeenCalledWith(
        { organizationId: 'org-456' },
        expect.objectContaining({ name: 'Updated Customer' }),
      );
      expect(notifySuccess).toHaveBeenCalled();
    });
  });
});
