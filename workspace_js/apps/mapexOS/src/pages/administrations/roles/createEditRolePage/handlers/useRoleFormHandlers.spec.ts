import { describe, it, expect, beforeEach, vi } from 'vitest';
import { ref } from 'vue';
import type { Ref } from 'vue';

import type { RoleFormData, ResourcePermission } from '../interfaces';
import { INITIAL_ROLE_FORM_DATA } from '../constants';

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
    flatList: [{ id: 'org-123', pathKey: '/vendor/customer-1' }],
  }),
}));

import { apis } from '@services/mapex';
import { notifySuccess, notifyFail } from '@utils/alert/notify';

import { useRoleFormHandlers } from './useRoleFormHandlers';

/**
 * Creates a mock translations object
 */
function createMockTranslations() {
  return {
    notifications: {
      created: { value: 'Role created' },
      updated: { value: 'Role updated' },
      createFailed: { value: 'Create failed' },
      updateFailed: { value: 'Update failed' },
      alreadyExists: { value: 'Already exists' },
      forbidden: { value: 'Forbidden' },
      noPermissions: { value: 'No permissions selected' },
    },
  };
}

/**
 * Creates a minimal resource permission for testing
 */
function createResourcePermission(overrides: Partial<ResourcePermission> = {}): ResourcePermission {
  return {
    resource: 'users',
    label: 'Users',
    icon: 'person',
    enabled: false,
    actions: [
      { name: 'list', label: 'List', granted: false },
      { name: 'create', label: 'Create', granted: false },
      { name: 'read', label: 'Read', granted: false },
      { name: 'update', label: 'Update', granted: false },
      { name: 'delete', label: 'Delete', granted: false },
    ],
    ...overrides,
  };
}

describe('useRoleFormHandlers', () => {
  let roleData: ReturnType<typeof ref<RoleFormData>>;
  let resourcePermissions: ReturnType<typeof ref<ResourcePermission[]>>;
  let currentStep: ReturnType<typeof ref<number>>;
  let isEditMode: ReturnType<typeof ref<boolean>>;
  let roleId: ReturnType<typeof ref<string | undefined>>;
  let isSaving: ReturnType<typeof ref<boolean>>;
  let step1Ref: ReturnType<typeof ref<any>>;
  let step2Ref: ReturnType<typeof ref<any>>;
  let t: ReturnType<typeof createMockTranslations>;

  beforeEach(() => {
    vi.clearAllMocks();
    roleData = ref({ ...INITIAL_ROLE_FORM_DATA });
    resourcePermissions = ref([createResourcePermission()]);
    currentStep = ref(1);
    isEditMode = ref(false);
    roleId = ref(undefined);
    isSaving = ref(false);
    step1Ref = ref(null);
    step2Ref = ref(null);
    t = createMockTranslations();
  });

  function setup() {
    return useRoleFormHandlers({
      roleData: roleData as Ref<RoleFormData>,
      resourcePermissions: resourcePermissions as Ref<ResourcePermission[]>,
      currentStep: currentStep as Ref<number>,
      isEditMode: isEditMode as Ref<boolean>,
      roleId,
      isSaving: isSaving as Ref<boolean>,
      step1Ref,
      step2Ref,
      t,
    });
  }

  describe('onResourceToggle', () => {
    it('toggles resource enabled and syncs all actions', () => {
      const { onResourceToggle } = setup();

      onResourceToggle(0);

      expect(resourcePermissions.value![0]!.enabled).toBe(true);
      expect(resourcePermissions.value![0]!.actions.every(a => a.granted)).toBe(true);
    });

    it('disables all actions when toggling off', () => {
      resourcePermissions.value![0]!.enabled = true;
      resourcePermissions.value![0]!.actions.forEach(a => (a.granted = true));
      const { onResourceToggle } = setup();

      onResourceToggle(0);

      expect(resourcePermissions.value![0]!.enabled).toBe(false);
      expect(resourcePermissions.value![0]!.actions.every(a => !a.granted)).toBe(true);
    });
  });

  describe('onActionToggle', () => {
    it('grants action and auto-selects list', () => {
      const { onActionToggle } = setup();
      const createIndex = resourcePermissions.value![0]!.actions.findIndex(a => a.name === 'create');

      onActionToggle(0, createIndex);

      expect(resourcePermissions.value![0]!.actions[createIndex]!.granted).toBe(true);
      const listAction = resourcePermissions.value![0]!.actions.find(a => a.name === 'list');
      expect(listAction?.granted).toBe(true);
      expect(resourcePermissions.value![0]!.enabled).toBe(true);
    });

    it('revokes action and disables resource if no actions remain', () => {
      // Grant only "list"
      resourcePermissions.value![0]!.actions[0]!.granted = true;
      resourcePermissions.value![0]!.enabled = true;
      const { onActionToggle } = setup();

      // Toggle "list" off
      onActionToggle(0, 0);

      expect(resourcePermissions.value![0]!.actions[0]!.granted).toBe(false);
      expect(resourcePermissions.value![0]!.enabled).toBe(false);
    });
  });

  describe('onToggleAllActions', () => {
    it('grants all actions for a resource', () => {
      const { onToggleAllActions } = setup();

      onToggleAllActions(0, true);

      expect(resourcePermissions.value![0]!.enabled).toBe(true);
      expect(resourcePermissions.value![0]!.actions.every(a => a.granted)).toBe(true);
    });

    it('revokes all actions for a resource', () => {
      resourcePermissions.value![0]!.enabled = true;
      resourcePermissions.value![0]!.actions.forEach(a => (a.granted = true));
      const { onToggleAllActions } = setup();

      onToggleAllActions(0, false);

      expect(resourcePermissions.value![0]!.enabled).toBe(false);
      expect(resourcePermissions.value![0]!.actions.every(a => !a.granted)).toBe(true);
    });
  });

  describe('selectedPermissionsCount', () => {
    it('returns 0 when no permissions are granted', () => {
      const { selectedPermissionsCount } = setup();

      expect(selectedPermissionsCount.value).toBe(0);
    });

    it('counts granted actions across all resources', () => {
      resourcePermissions.value![0]!.actions[0]!.granted = true;
      resourcePermissions.value![0]!.actions[1]!.granted = true;
      const { selectedPermissionsCount } = setup();

      expect(selectedPermissionsCount.value).toBe(2);
    });
  });

  describe('isNextButtonDisabled', () => {
    it('returns true on step 1 when name is too short', () => {
      currentStep.value = 1;
      roleData.value!.name = 'AB';
      roleData.value!.scope = 'global';
      const { isNextButtonDisabled } = setup();

      expect(isNextButtonDisabled.value).toBe(true);
    });

    it('returns true on step 1 when scope is null', () => {
      currentStep.value = 1;
      roleData.value!.name = 'Valid Name';
      roleData.value!.scope = null;
      const { isNextButtonDisabled } = setup();

      expect(isNextButtonDisabled.value).toBe(true);
    });

    it('returns false on step 1 when name and scope are valid', () => {
      currentStep.value = 1;
      roleData.value!.name = 'Admin Role';
      roleData.value!.scope = 'global';
      const { isNextButtonDisabled } = setup();

      expect(isNextButtonDisabled.value).toBe(false);
    });

    it('returns true on step 2 when no permissions selected', () => {
      currentStep.value = 2;
      const { isNextButtonDisabled } = setup();

      expect(isNextButtonDisabled.value).toBe(true);
    });

    it('returns false on step 2 when permissions are selected', () => {
      currentStep.value = 2;
      resourcePermissions.value![0]!.actions[0]!.granted = true;
      const { isNextButtonDisabled } = setup();

      expect(isNextButtonDisabled.value).toBe(false);
    });
  });

  describe('buildPermissionsArray', () => {
    it('returns empty array when no permissions granted', () => {
      const { buildPermissionsArray } = setup();

      expect(buildPermissionsArray()).toEqual([]);
    });

    it('builds correct permission strings', () => {
      resourcePermissions.value![0]!.actions[0]!.granted = true; // users.list
      resourcePermissions.value![0]!.actions[1]!.granted = true; // users.create
      const { buildPermissionsArray } = setup();

      expect(buildPermissionsArray()).toEqual(['users.list', 'users.create']);
    });

    it('uses permissionKey override when present', () => {
      resourcePermissions.value = [
        createResourcePermission({
          resource: 'events.raw',
          actions: [
            { name: 'list', label: 'List', granted: true, permissionKey: 'events.raw.list' },
          ],
        }),
      ];
      const { buildPermissionsArray } = setup();

      expect(buildPermissionsArray()).toEqual(['events.raw.list']);
    });
  });

  describe('changeStep', () => {
    it('allows backward navigation without validation', async () => {
      currentStep.value = 2;
      const { changeStep } = setup();

      await changeStep(1);

      expect(currentStep.value).toBe(1);
    });

    it('blocks forward navigation when step validation fails', async () => {
      step1Ref.value = { validate: vi.fn().mockResolvedValue(false) };
      currentStep.value = 1;
      const { changeStep } = setup();

      await changeStep(2);

      expect(currentStep.value).toBe(1);
    });
  });

  describe('submitForm — CREATE mode', () => {
    beforeEach(() => {
      (apis.mapexOS as any).roles = {
        create: vi.fn().mockResolvedValue({}),
        update: vi.fn().mockResolvedValue({}),
      };
    });

    it('calls create API with permissions array', async () => {
      roleData.value = { name: 'Admin', description: 'Desc', scope: 'global', isTemplate: false };
      resourcePermissions.value![0]!.actions[0]!.granted = true; // users.list
      const { submitForm } = setup();

      await submitForm();

      expect(apis.mapexOS.roles.create).toHaveBeenCalledWith(
        expect.objectContaining({
          name: 'Admin',
          permissions: ['users.list'],
          scope: 'global',
          pathKey: '/vendor/customer-1',
        }),
      );
      expect(notifySuccess).toHaveBeenCalled();
      expect(isSaving.value).toBe(false);
    });

    it('fails with notification when no permissions selected', async () => {
      roleData.value = { name: 'Admin', description: '', scope: 'global', isTemplate: false };
      const { submitForm } = setup();

      await submitForm();

      expect(notifyFail).toHaveBeenCalledWith(
        expect.objectContaining({ message: 'No permissions selected' }),
      );
      expect(apis.mapexOS.roles.create).not.toHaveBeenCalled();
    });
  });

  describe('submitForm — EDIT mode', () => {
    beforeEach(() => {
      isEditMode.value = true;
      roleId.value = 'role-123';
      (apis.mapexOS as any).roles = {
        create: vi.fn(),
        update: vi.fn().mockResolvedValue({}),
      };
    });

    it('calls update API with role ID', async () => {
      roleData.value = { name: 'Updated Role', description: '', scope: 'local', isTemplate: false };
      resourcePermissions.value![0]!.actions[0]!.granted = true;
      const { submitForm } = setup();

      await submitForm();

      expect(apis.mapexOS.roles.update).toHaveBeenCalledWith(
        { roleId: 'role-123' },
        expect.objectContaining({
          name: 'Updated Role',
          permissions: ['users.list'],
        }),
      );
    });

    it('handles 409 conflict error', async () => {
      (apis.mapexOS as any).roles.update.mockRejectedValue({ response: { status: 409 } });
      roleData.value = { name: 'Dup', description: '', scope: 'local', isTemplate: false };
      resourcePermissions.value![0]!.actions[0]!.granted = true;
      const { submitForm } = setup();

      await submitForm();

      expect(notifyFail).toHaveBeenCalledWith(
        expect.objectContaining({ message: 'Already exists' }),
      );
    });
  });
});
