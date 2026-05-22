import { describe, it, expect, beforeEach, vi } from 'vitest';
import { ref } from 'vue';
import type { Ref } from 'vue';

import type { UserFormData } from '../interfaces';
import { VALIDATION, INITIAL_USER_FORM_DATA } from '../constants';

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

import { apis } from '@services/mapex';
import { notifySuccess, notifyFail } from '@utils/alert/notify';

import { useUserFormHandlers } from './useUserFormHandlers';

/**
 * Creates a mock translations object matching what the handler expects
 */
function createMockTranslations() {
  return {
    messages: {
      created: { value: 'User created' },
      updated: { value: 'User updated' },
      createFailed: { value: 'Create failed' },
      updateFailed: { value: 'Update failed' },
      emailExists: { value: 'Email already exists' },
      forbidden: { value: 'Forbidden' },
    },
  };
}

describe('useUserFormHandlers', () => {
  let userData: ReturnType<typeof ref<UserFormData>>;
  let currentStep: ReturnType<typeof ref<number>>;
  let isEditMode: ReturnType<typeof ref<boolean>>;
  let userId: ReturnType<typeof ref<string | undefined>>;
  let isSaving: ReturnType<typeof ref<boolean>>;
  let isTourMode: ReturnType<typeof ref<boolean>>;
  let step1FormRef: ReturnType<typeof ref<any>>;
  let step2FormRef: ReturnType<typeof ref<any>>;
  let step3FormRef: ReturnType<typeof ref<any>>;
  let t: ReturnType<typeof createMockTranslations>;

  beforeEach(() => {
    vi.clearAllMocks();
    userData = ref({ ...INITIAL_USER_FORM_DATA });
    currentStep = ref(1);
    isEditMode = ref(false);
    userId = ref(undefined);
    isSaving = ref(false);
    isTourMode = ref(false);
    step1FormRef = ref(null);
    step2FormRef = ref(null);
    step3FormRef = ref(null);
    t = createMockTranslations();
  });

  function setup() {
    return useUserFormHandlers({
      userData: userData as Ref<UserFormData>,
      currentStep: currentStep as Ref<number>,
      isEditMode: isEditMode as Ref<boolean>,
      userId,
      isSaving: isSaving as Ref<boolean>,
      isTourMode: isTourMode as Ref<boolean>,
      step1FormRef,
      step2FormRef,
      step3FormRef,
      t,
    });
  }

  describe('isNextButtonDisabled', () => {
    describe('Step 1 — Personal', () => {
      it('returns true when firstName is too short', () => {
        currentStep.value = 1;
        userData.value!.firstName = 'A';
        userData.value!.lastName = 'Smith';
        userData.value!.email = 'a@b.com';
        const { isNextButtonDisabled } = setup();

        expect(isNextButtonDisabled.value).toBe(true);
      });

      it('returns true when lastName is empty', () => {
        currentStep.value = 1;
        userData.value!.firstName = 'John';
        userData.value!.lastName = '';
        userData.value!.email = 'a@b.com';
        const { isNextButtonDisabled } = setup();

        expect(isNextButtonDisabled.value).toBe(true);
      });

      it('returns true when email is empty', () => {
        currentStep.value = 1;
        userData.value!.firstName = 'John';
        userData.value!.lastName = 'Smith';
        userData.value!.email = '';
        const { isNextButtonDisabled } = setup();

        expect(isNextButtonDisabled.value).toBe(true);
      });

      it('returns false when all required fields are valid', () => {
        currentStep.value = 1;
        userData.value!.firstName = 'John';
        userData.value!.lastName = 'Smith';
        userData.value!.email = 'john@example.com';
        const { isNextButtonDisabled } = setup();

        expect(isNextButtonDisabled.value).toBe(false);
      });
    });

    describe('Step 2 — Security', () => {
      it('returns true in create mode when password is too short', () => {
        currentStep.value = 2;
        isEditMode.value = false;
        userData.value!.password = 'short';
        const { isNextButtonDisabled } = setup();

        expect(isNextButtonDisabled.value).toBe(true);
      });

      it('returns false in create mode when password meets minimum length', () => {
        currentStep.value = 2;
        isEditMode.value = false;
        userData.value!.password = 'a'.repeat(VALIDATION.PASSWORD_MIN_LENGTH);
        const { isNextButtonDisabled } = setup();

        expect(isNextButtonDisabled.value).toBe(false);
      });

      it('returns false in edit mode even without password', () => {
        currentStep.value = 2;
        isEditMode.value = true;
        userData.value!.password = '';
        const { isNextButtonDisabled } = setup();

        expect(isNextButtonDisabled.value).toBe(false);
      });
    });

    describe('Step 3 — Access', () => {
      it('returns true when no access type has valid data', () => {
        currentStep.value = 3;
        userData.value!.accessType = 'group';
        userData.value!.selectedGroup = undefined;
        userData.value!.selectedGroups = [];
        const { isNextButtonDisabled } = setup();

        expect(isNextButtonDisabled.value).toBe(true);
      });

      it('returns false for group access with valid existing group', () => {
        currentStep.value = 3;
        userData.value!.accessType = 'group';
        userData.value!.selectedGroups = [
          { mode: 'existing', existingGroup: { groupId: 'g-1', groupName: 'G1' } },
        ];
        const { isNextButtonDisabled } = setup();

        expect(isNextButtonDisabled.value).toBe(false);
      });

      it('returns false for direct access with valid membership', () => {
        currentStep.value = 3;
        userData.value!.accessType = 'direct';
        userData.value!.directMemberships = [
          { orgId: 'o-1', orgName: 'Org', roleIds: ['r-1'], roleNames: ['Admin'], scope: 'local' },
        ];
        const { isNextButtonDisabled } = setup();

        expect(isNextButtonDisabled.value).toBe(false);
      });

      it('returns true for group access with new group missing roles', () => {
        currentStep.value = 3;
        userData.value!.accessType = 'group';
        userData.value!.selectedGroups = [
          { mode: 'new', newGroup: { name: 'New Group', roleIds: [], roleNames: [] } },
        ];
        const { isNextButtonDisabled } = setup();

        expect(isNextButtonDisabled.value).toBe(true);
      });
    });
  });

  describe('changeStep', () => {
    it('changes step when no form ref exists', async () => {
      const { changeStep } = setup();

      await changeStep(2);

      expect(currentStep.value).toBe(2);
    });

    it('blocks forward navigation when form validation fails', async () => {
      step1FormRef.value = { validate: vi.fn().mockResolvedValue(false) };
      currentStep.value = 1;
      const { changeStep } = setup();

      await changeStep(2);

      expect(currentStep.value).toBe(1);
    });

    it('allows forward navigation when form validation passes', async () => {
      step1FormRef.value = { validate: vi.fn().mockResolvedValue(true) };
      currentStep.value = 1;
      const { changeStep } = setup();

      await changeStep(2);

      expect(currentStep.value).toBe(2);
    });

    it('allows backward navigation without validation', async () => {
      currentStep.value = 3;
      const { changeStep } = setup();

      await changeStep(1);

      expect(currentStep.value).toBe(1);
    });

    it('skips validation in tour mode', async () => {
      isTourMode.value = true;
      step1FormRef.value = { validate: vi.fn().mockResolvedValue(false) };
      currentStep.value = 1;
      const { changeStep } = setup();

      await changeStep(2);

      expect(currentStep.value).toBe(2);
      expect(step1FormRef.value.validate).not.toHaveBeenCalled();
    });
  });

  describe('submitForm — CREATE mode', () => {
    beforeEach(() => {
      (apis.mapexOS as any).onboarding = {
        createUserWithMemberships: vi.fn().mockResolvedValue({}),
        updateUserWithAccess: vi.fn().mockResolvedValue({}),
      };
    });

    it('calls createUserWithMemberships with correct payload', async () => {
      userData.value = {
        ...INITIAL_USER_FORM_DATA,
        email: 'john@test.com',
        firstName: 'John',
        lastName: 'Doe',
        password: 'SecureP@ss123',
        enabled: true,
        changePasswordNextLogin: false,
        accessType: 'group',
        selectedGroup: {
          mode: 'existing',
          existingGroup: { groupId: 'g-1', groupName: 'Team' },
        },
      };
      const { submitForm } = setup();

      await submitForm();

      expect(apis.mapexOS.onboarding.createUserWithMemberships).toHaveBeenCalledWith(
        expect.objectContaining({
          email: 'john@test.com',
          firstName: 'John',
          lastName: 'Doe',
          password: 'SecureP@ss123',
          groups: [{ existingGroup: { groupId: 'g-1' } }],
        }),
      );
      expect(notifySuccess).toHaveBeenCalled();
      expect(isSaving.value).toBe(false);
    });

    it('handles 409 error with email exists message', async () => {
      (apis.mapexOS as any).onboarding.createUserWithMemberships.mockRejectedValue({
        response: { status: 409 },
      });
      userData.value!.accessType = 'group';
      userData.value!.selectedGroup = {
        mode: 'existing',
        existingGroup: { groupId: 'g-1', groupName: 'T' },
      };
      const { submitForm } = setup();

      await submitForm();

      expect(notifyFail).toHaveBeenCalledWith(
        expect.objectContaining({ message: 'Email already exists' }),
      );
    });
  });

  describe('submitForm — EDIT mode', () => {
    beforeEach(() => {
      isEditMode.value = true;
      userId.value = 'user-123';
      (apis.mapexOS as any).onboarding = {
        createUserWithMemberships: vi.fn(),
        updateUserWithAccess: vi.fn().mockResolvedValue({}),
      };
    });

    it('calls updateUserWithAccess with user ID', async () => {
      userData.value = {
        ...INITIAL_USER_FORM_DATA,
        firstName: 'Jane',
        lastName: 'Doe',
        email: 'jane@test.com',
        accessType: 'direct',
        directMembership: {
          orgId: 'o-1',
          orgName: 'Org',
          roleIds: ['r-1'],
          roleNames: ['Admin'],
          scope: 'local',
        },
      };
      const { submitForm } = setup();

      await submitForm();

      expect(apis.mapexOS.onboarding.updateUserWithAccess).toHaveBeenCalledWith(
        { userId: 'user-123' },
        expect.objectContaining({
          firstName: 'Jane',
          lastName: 'Doe',
          memberships: [expect.objectContaining({ roles: ['r-1'] })],
        }),
      );
      expect(notifySuccess).toHaveBeenCalled();
    });
  });

  describe('submitForm — Tour mode', () => {
    it('skips API call and navigates away', async () => {
      isTourMode.value = true;
      (apis.mapexOS as any).onboarding = {
        createUserWithMemberships: vi.fn(),
      };
      const { submitForm } = setup();

      await submitForm();

      expect(apis.mapexOS.onboarding.createUserWithMemberships).not.toHaveBeenCalled();
    });
  });
});
