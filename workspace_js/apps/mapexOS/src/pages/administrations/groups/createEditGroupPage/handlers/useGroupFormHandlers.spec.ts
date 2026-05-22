import { describe, it, expect, beforeEach, vi } from 'vitest';
import { ref } from 'vue';
import type { Ref } from 'vue';

import type { GroupFormData, RoleSelectionItem } from '../interfaces';
import type { UserSelectorItem } from '@components/drawers';
import { INITIAL_GROUP_FORM_DATA } from '../constants';

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
  }),
}));

import { apis } from '@services/mapex';
import { notifySuccess, notifyFail } from '@utils/alert/notify';

import { useGroupFormHandlers } from './useGroupFormHandlers';

/**
 * Creates a mock translations object
 */
function createMockTranslations() {
  return {
    createEditNotifications: {
      created: { value: 'Group created' },
      updated: { value: 'Group updated' },
      createFailed: { value: 'Create failed' },
      updateFailed: { value: 'Update failed' },
      alreadyExists: { value: 'Already exists' },
      forbidden: { value: 'Forbidden' },
    },
  };
}

describe('useGroupFormHandlers', () => {
  let groupData: ReturnType<typeof ref<GroupFormData>>;
  let selectedRoles: ReturnType<typeof ref<RoleSelectionItem[]>>;
  let selectedMembers: ReturnType<typeof ref<string[]>>;
  let pendingAdditions: ReturnType<typeof ref<any[]>>;
  let pendingRemovals: ReturnType<typeof ref<string[]>>;
  let currentStep: ReturnType<typeof ref<number>>;
  let isEditMode: ReturnType<typeof ref<boolean>>;
  let groupId: ReturnType<typeof ref<string | undefined>>;
  let isSaving: ReturnType<typeof ref<boolean>>;
  let step1FormRef: ReturnType<typeof ref<any>>;
  let step2RolesFormRef: ReturnType<typeof ref<any>>;
  let step3MembersFormRef: ReturnType<typeof ref<any>>;
  let t: ReturnType<typeof createMockTranslations>;

  beforeEach(() => {
    vi.clearAllMocks();
    groupData = ref({ ...INITIAL_GROUP_FORM_DATA });
    selectedRoles = ref<RoleSelectionItem[]>([]);
    selectedMembers = ref<string[]>([]);
    pendingAdditions = ref([]);
    pendingRemovals = ref([]);
    currentStep = ref(1);
    isEditMode = ref(false);
    groupId = ref(undefined);
    isSaving = ref(false);
    step1FormRef = ref(null);
    step2RolesFormRef = ref(null);
    step3MembersFormRef = ref(null);
    t = createMockTranslations();

    // Setup API mocks
    (apis.mapexOS as any).groups = {
      ...apis.mapexOS.groups,
      create: vi.fn().mockResolvedValue({ id: 'new-group-id' }),
      update: vi.fn().mockResolvedValue({}),
      addMember: vi.fn().mockResolvedValue({}),
      removeMember: vi.fn().mockResolvedValue({}),
    };
  });

  function setup() {
    return useGroupFormHandlers({
      groupData: groupData as Ref<GroupFormData>,
      selectedRoles: selectedRoles as Ref<RoleSelectionItem[]>,
      selectedMembers: selectedMembers as Ref<string[]>,
      pendingAdditions: pendingAdditions as Ref<UserSelectorItem[]>,
      pendingRemovals: pendingRemovals as Ref<string[]>,
      currentStep: currentStep as Ref<number>,
      isEditMode: isEditMode as Ref<boolean>,
      groupId,
      isSaving: isSaving as Ref<boolean>,
      step1FormRef,
      step2RolesFormRef,
      step3MembersFormRef,
      t,
    });
  }

  describe('selectedMembersCount', () => {
    it('returns 0 when no members selected', () => {
      const { selectedMembersCount } = setup();

      expect(selectedMembersCount.value).toBe(0);
    });

    it('returns correct count', () => {
      selectedMembers.value = ['u-1', 'u-2', 'u-3'];
      const { selectedMembersCount } = setup();

      expect(selectedMembersCount.value).toBe(3);
    });
  });

  describe('isNextButtonDisabled', () => {
    it('returns true on step 1 when name is too short', () => {
      currentStep.value = 1;
      groupData.value!.name = 'AB';
      const { isNextButtonDisabled } = setup();

      expect(isNextButtonDisabled.value).toBe(true);
    });

    it('returns false on step 1 when name is valid', () => {
      currentStep.value = 1;
      groupData.value!.name = 'Engineering Team';
      const { isNextButtonDisabled } = setup();

      expect(isNextButtonDisabled.value).toBe(false);
    });

    it('returns true on step 2 when no roles selected', () => {
      currentStep.value = 2;
      selectedRoles.value = [];
      const { isNextButtonDisabled } = setup();

      expect(isNextButtonDisabled.value).toBe(true);
    });

    it('returns false on step 2 when roles are selected', () => {
      currentStep.value = 2;
      selectedRoles.value = [{ id: 'r-1', name: 'Admin' }];
      const { isNextButtonDisabled } = setup();

      expect(isNextButtonDisabled.value).toBe(false);
    });

    it('returns false on step 3 (members are optional)', () => {
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

    it('blocks forward from step 2 when no roles selected', async () => {
      currentStep.value = 2;
      selectedRoles.value = [];
      const { changeStep } = setup();

      await changeStep(3);

      expect(currentStep.value).toBe(2);
    });

    it('allows forward from step 2 when roles are selected', async () => {
      currentStep.value = 2;
      selectedRoles.value = [{ id: 'r-1', name: 'Admin' }];
      const { changeStep } = setup();

      await changeStep(3);

      expect(currentStep.value).toBe(3);
    });

    it('allows backward navigation without validation', async () => {
      currentStep.value = 3;
      const { changeStep } = setup();

      await changeStep(1);

      expect(currentStep.value).toBe(1);
    });
  });

  describe('submitForm — CREATE mode', () => {
    it('calls create API with role IDs', async () => {
      groupData.value!.name = 'New Group';
      groupData.value!.description = 'Desc';
      selectedRoles.value = [{ id: 'r-1', name: 'Admin' }, { id: 'r-2', name: 'Editor' }];
      const { submitForm } = setup();

      await submitForm();

      expect(apis.mapexOS.groups.create).toHaveBeenCalledWith(
        expect.objectContaining({
          name: 'New Group',
          roleIds: ['r-1', 'r-2'],
          orgId: 'org-123',
        }),
      );
      expect(notifySuccess).toHaveBeenCalled();
      expect(isSaving.value).toBe(false);
    });

    it('adds pending members after creation', async () => {
      groupData.value!.name = 'New Group';
      selectedRoles.value = [{ id: 'r-1', name: 'Admin' }];
      pendingAdditions.value = [{ id: 'u-1', name: 'User 1' }];
      const { submitForm } = setup();

      await submitForm();

      expect(apis.mapexOS.groups.addMember).toHaveBeenCalledWith(
        { groupId: 'new-group-id' },
        { userId: 'u-1' },
      );
    });

    it('handles 409 conflict error', async () => {
      (apis.mapexOS as any).groups.create.mockRejectedValue({ response: { status: 409 } });
      groupData.value!.name = 'Duplicate';
      selectedRoles.value = [{ id: 'r-1', name: 'Admin' }];
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
      groupId.value = 'group-123';
    });

    it('calls update API with group ID', async () => {
      groupData.value!.name = 'Updated Group';
      const { submitForm } = setup();

      await submitForm();

      expect(apis.mapexOS.groups.update).toHaveBeenCalledWith(
        { groupId: 'group-123' },
        expect.objectContaining({ name: 'Updated Group' }),
      );
      expect(notifySuccess).toHaveBeenCalled();
    });

    it('persists member changes (removals then additions)', async () => {
      groupData.value!.name = 'Group';
      pendingRemovals.value = ['u-old'];
      pendingAdditions.value = [{ id: 'u-new', name: 'New User' }];
      const { submitForm } = setup();

      await submitForm();

      expect(apis.mapexOS.groups.removeMember).toHaveBeenCalledWith({
        groupId: 'group-123',
        userId: 'u-old',
      });
      expect(apis.mapexOS.groups.addMember).toHaveBeenCalledWith(
        { groupId: 'group-123' },
        { userId: 'u-new' },
      );
    });
  });
});
