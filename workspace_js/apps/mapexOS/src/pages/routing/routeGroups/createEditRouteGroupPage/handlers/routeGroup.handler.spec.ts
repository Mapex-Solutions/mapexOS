import { describe, it, expect, vi, beforeEach } from 'vitest';
import { ref } from 'vue';
import type { RouterFormState } from '../interfaces';
import {
  handleAddRouter,
  handleRemoveRouter,
  handleRouterKindChange,
  handleAddMatchRule,
  handleRemoveMatchRule,
  handleToggleConditionalRouting,
  handleChangeStep,
  handleSave,
  handleLoadRouteGroup,
} from './routeGroup.handler';
import { apis } from '@services/mapex';
import { DEFAULT_ROUTER_KIND, DEFAULT_MATCH_POLICY, DEFAULT_MATCH_OPERATOR, STEP } from '../constants';

vi.mock('@utils/alert/notify', () => ({
  notifySuccess: vi.fn(),
}));
vi.mock('@utils/error', () => ({
  handleApiError: vi.fn(),
}));
vi.mock('@composables/useLogger', () => ({
  useLogger: () => ({
    warn: vi.fn(),
    error: vi.fn(),
    debug: vi.fn(),
  }),
}));

function makeRouter(overrides: Partial<RouterFormState> = {}): RouterFormState {
  return {
    id: `router-${Date.now()}`,
    kind: DEFAULT_ROUTER_KIND,
    hasConditionalRouting: false,
    saveEvent: {},
    ...overrides,
  };
}

describe('routeGroup.handler', () => {
  // ── handleAddRouter ──────────────────────────────────────────────────
  describe('handleAddRouter', () => {
    it('should add a new router with default kind to the array', () => {
      const routerForms = ref<RouterFormState[]>([]);

      handleAddRouter(routerForms);

      expect(routerForms.value).toHaveLength(1);
      expect(routerForms.value[0]!.kind).toBe(DEFAULT_ROUTER_KIND);
      expect(routerForms.value[0]!.hasConditionalRouting).toBe(false);
      expect(routerForms.value[0]!.id).toContain('router-');
    });

    it('should append to existing routers', () => {
      const existing = makeRouter({ id: 'existing-1' });
      const routerForms = ref<RouterFormState[]>([existing]);

      handleAddRouter(routerForms);

      expect(routerForms.value).toHaveLength(2);
      expect(routerForms.value[0]!.id).toBe('existing-1');
    });
  });

  // ── handleRemoveRouter ───────────────────────────────────────────────
  describe('handleRemoveRouter', () => {
    it('should remove the router with the given id', () => {
      const routerForms = ref<RouterFormState[]>([
        makeRouter({ id: 'r1' }),
        makeRouter({ id: 'r2' }),
        makeRouter({ id: 'r3' }),
      ]);

      handleRemoveRouter(routerForms, 'r2');

      expect(routerForms.value).toHaveLength(2);
      expect(routerForms.value.map((r) => r.id)).toEqual(['r1', 'r3']);
    });

    it('should do nothing when id does not exist', () => {
      const routerForms = ref<RouterFormState[]>([makeRouter({ id: 'r1' })]);

      handleRemoveRouter(routerForms, 'nonexistent');

      expect(routerForms.value).toHaveLength(1);
    });
  });

  // ── handleRouterKindChange ───────────────────────────────────────────
  describe('handleRouterKindChange', () => {
    it('should initialize lake_house config', () => {
      const routerForms = ref<RouterFormState[]>([makeRouter({ id: 'r1' })]);

      handleRouterKindChange(routerForms, 'r1', 'lake_house');

      expect(routerForms.value[0]!.lakeHouse).toEqual({ lakeHouseId: '', metadata: {} });
    });

    it('should initialize notification config', () => {
      const routerForms = ref<RouterFormState[]>([makeRouter({ id: 'r1' })]);

      handleRouterKindChange(routerForms, 'r1', 'notification');

      expect(routerForms.value[0]!.notification).toEqual({ notificationId: '', metadata: {} });
    });

    it('should initialize save_event config', () => {
      const routerForms = ref<RouterFormState[]>([makeRouter({ id: 'r1' })]);

      handleRouterKindChange(routerForms, 'r1', 'save_event');

      expect(routerForms.value[0]!.saveEvent).toEqual({ metadata: {} });
    });

    it('should initialize workflow config', () => {
      const routerForms = ref<RouterFormState[]>([makeRouter({ id: 'r1' })]);

      handleRouterKindChange(routerForms, 'r1', 'workflow');

      expect(routerForms.value[0]!.workflow).toEqual({
        mode: 'newInstance',
        workflowId: '',
        workflowUUID: '',
        signalName: '',
        metadata: {},
      });
    });

    it('should do nothing when router id is not found', () => {
      const routerForms = ref<RouterFormState[]>([makeRouter({ id: 'r1' })]);

      handleRouterKindChange(routerForms, 'nonexistent', 'lake_house');

      expect(routerForms.value[0]!.saveEvent).toBeDefined();
    });
  });

  // ── handleAddMatchRule ───────────────────────────────────────────────
  describe('handleAddMatchRule', () => {
    it('should create match config with default policy and add a rule', () => {
      const routerForms = ref<RouterFormState[]>([makeRouter({ id: 'r1' })]);

      handleAddMatchRule(routerForms, 'r1');

      const router = routerForms.value[0]!;
      expect(router.match).toBeDefined();
      expect(router.match!.policy).toBe(DEFAULT_MATCH_POLICY);
      expect(router.match!.rules).toHaveLength(1);
      expect(router.match!.rules[0]).toEqual({
        field: '',
        operator: DEFAULT_MATCH_OPERATOR,
        value: '',
      });
    });

    it('should append rules when match already exists', () => {
      const routerForms = ref<RouterFormState[]>([
        makeRouter({
          id: 'r1',
          match: {
            policy: 'all',
            rules: [{ field: 'a', operator: 'eq', value: '1' }],
          },
        }),
      ]);

      handleAddMatchRule(routerForms, 'r1');

      expect(routerForms.value[0]!.match!.rules).toHaveLength(2);
    });

    it('should do nothing when router id is not found', () => {
      const routerForms = ref<RouterFormState[]>([makeRouter({ id: 'r1' })]);

      handleAddMatchRule(routerForms, 'nonexistent');

      expect(routerForms.value[0]!.match).toBeUndefined();
    });
  });

  // ── handleRemoveMatchRule ────────────────────────────────────────────
  describe('handleRemoveMatchRule', () => {
    it('should remove the rule at the given index', () => {
      const routerForms = ref<RouterFormState[]>([
        makeRouter({
          id: 'r1',
          hasConditionalRouting: true,
          match: {
            policy: 'all',
            rules: [
              { field: 'a', operator: 'eq', value: '1' },
              { field: 'b', operator: 'eq', value: '2' },
            ],
          },
        }),
      ]);

      handleRemoveMatchRule(routerForms, 'r1', 0);

      expect(routerForms.value[0]!.match!.rules).toHaveLength(1);
      expect(routerForms.value[0]!.match!.rules[0]!.field).toBe('b');
    });

    it('should clean up match config and disable conditional routing when last rule removed', () => {
      const routerForms = ref<RouterFormState[]>([
        makeRouter({
          id: 'r1',
          hasConditionalRouting: true,
          match: {
            policy: 'all',
            rules: [{ field: 'a', operator: 'eq', value: '1' }],
          },
        }),
      ]);

      handleRemoveMatchRule(routerForms, 'r1', 0);

      expect(routerForms.value[0]!.hasConditionalRouting).toBe(false);
      expect(routerForms.value[0]!.match).toBeUndefined();
    });

    it('should do nothing when router has no match config', () => {
      const routerForms = ref<RouterFormState[]>([makeRouter({ id: 'r1' })]);

      handleRemoveMatchRule(routerForms, 'r1', 0);

      expect(routerForms.value[0]!.match).toBeUndefined();
    });
  });

  // ── handleToggleConditionalRouting ───────────────────────────────────
  describe('handleToggleConditionalRouting', () => {
    it('should enable conditional routing and create a default rule', () => {
      const routerForms = ref<RouterFormState[]>([makeRouter({ id: 'r1' })]);

      handleToggleConditionalRouting(routerForms, 'r1', true);

      const router = routerForms.value[0]!;
      expect(router.hasConditionalRouting).toBe(true);
      expect(router.match).toBeDefined();
      expect(router.match!.policy).toBe(DEFAULT_MATCH_POLICY);
      expect(router.match!.rules).toHaveLength(1);
    });

    it('should disable conditional routing and remove match config', () => {
      const routerForms = ref<RouterFormState[]>([
        makeRouter({
          id: 'r1',
          hasConditionalRouting: true,
          match: {
            policy: 'all',
            rules: [{ field: 'a', operator: 'eq', value: '1' }],
          },
        }),
      ]);

      handleToggleConditionalRouting(routerForms, 'r1', false);

      expect(routerForms.value[0]!.hasConditionalRouting).toBe(false);
      expect(routerForms.value[0]!.match).toBeUndefined();
    });

    it('should not create match if already exists when enabling', () => {
      const existingMatch = {
        policy: 'any' as const,
        rules: [{ field: 'x', operator: 'eq' as const, value: 'y' }],
      };
      const routerForms = ref<RouterFormState[]>([
        makeRouter({ id: 'r1', match: existingMatch }),
      ]);

      handleToggleConditionalRouting(routerForms, 'r1', true);

      // Should keep existing match, not overwrite
      expect(routerForms.value[0]!.match!.policy).toBe('any');
      expect(routerForms.value[0]!.match!.rules[0]!.field).toBe('x');
    });
  });

  // ── handleChangeStep ─────────────────────────────────────────────────
  describe('handleChangeStep', () => {
    it('should change step when going from step 1 with valid form', async () => {
      const currentStep = ref(1);
      const step1FormRef = ref({ validate: vi.fn().mockResolvedValue(true) });
      const routerForms = ref<RouterFormState[]>([]);

      await handleChangeStep(currentStep, 2, step1FormRef, routerForms);

      expect(currentStep.value).toBe(2);
    });

    it('should not change step when step 1 validation fails', async () => {
      const currentStep = ref(1);
      const step1FormRef = ref({ validate: vi.fn().mockResolvedValue(false) });
      const routerForms = ref<RouterFormState[]>([]);

      await handleChangeStep(currentStep, 2, step1FormRef, routerForms);

      expect(currentStep.value).toBe(1);
    });

    it('should not advance from step 2 if no routers exist', async () => {
      const currentStep = ref(2);
      const step1FormRef = ref(null);
      const routerForms = ref<RouterFormState[]>([]);

      await handleChangeStep(currentStep, 3, step1FormRef, routerForms);

      expect(currentStep.value).toBe(2);
    });

    it('should not advance from step 2 if a conditional router has no rules', async () => {
      const currentStep = ref(2);
      const step1FormRef = ref(null);
      const routerForms = ref<RouterFormState[]>([
        makeRouter({
          id: 'r1',
          hasConditionalRouting: true,
          match: { policy: 'all', rules: [] },
        }),
      ]);

      await handleChangeStep(currentStep, 3, step1FormRef, routerForms);

      expect(currentStep.value).toBe(2);
    });

    it('should allow going backward without validation', async () => {
      const currentStep = ref(2);
      const step1FormRef = ref(null);
      const routerForms = ref<RouterFormState[]>([]);

      await handleChangeStep(currentStep, 1, step1FormRef, routerForms);

      expect(currentStep.value).toBe(1);
    });
  });

  // ── handleSave ───────────────────────────────────────────────────────
  describe('handleSave', () => {
    const mockRouter = { push: vi.fn() };
    const mockT = {
      notifications: {
        created: { value: 'Created' },
        updated: { value: 'Updated' },
        alreadyExists: 'Already exists',
        validationFailed: 'Validation failed',
        networkError: 'Network error',
        creationFailed: 'Creation failed',
        updateFailed: 'Update failed',
      },
    };

    beforeEach(() => {
      vi.clearAllMocks();
      // Add missing API methods for this test
      (apis.router.routegroup as any).create = vi.fn().mockResolvedValue({});
      (apis.router.routegroup as any).update = vi.fn().mockResolvedValue({});
    });

    it('should create a new route group in create mode', async () => {
      const isSaving = ref(false);
      const isEditMode = ref(false);
      const routeGroupId = ref<string | undefined>(undefined);
      const formData = ref({
        name: 'Test Group',
        description: 'Desc',
        version: '1.0.0',
        enabled: true,
        isTemplate: false,
      });
      const routerForms = ref<RouterFormState[]>([
        makeRouter({ id: 'r1', kind: 'save_event', saveEvent: {} }),
      ]);

      await handleSave(
        isSaving,
        isEditMode,
        routeGroupId,
        formData,
        routerForms,
        mockRouter as any,
        mockT,
      );

      expect(apis.router.routegroup.create).toHaveBeenCalledWith(
        expect.objectContaining({ name: 'Test Group' }),
      );
      expect(mockRouter.push).toHaveBeenCalledWith('/routing/route_groups');
      expect(isSaving.value).toBe(false);
    });

    it('should update an existing route group in edit mode', async () => {
      const isSaving = ref(false);
      const isEditMode = ref(true);
      const routeGroupId = ref<string | undefined>('rg-123');
      const formData = ref({
        name: 'Updated Group',
        description: '',
        version: '1.0.0',
        enabled: true,
        isTemplate: false,
      });
      const routerForms = ref<RouterFormState[]>([]);

      await handleSave(
        isSaving,
        isEditMode,
        routeGroupId,
        formData,
        routerForms,
        mockRouter as any,
        mockT,
      );

      expect((apis.router.routegroup as any).update).toHaveBeenCalledWith(
        { routeGroupId: 'rg-123' },
        expect.objectContaining({ name: 'Updated Group' }),
      );
      expect(isSaving.value).toBe(false);
    });

    it('should set isSaving to false even on error', async () => {
      (apis.router.routegroup as any).create = vi.fn().mockRejectedValue(new Error('fail'));
      const isSaving = ref(false);
      const isEditMode = ref(false);
      const routeGroupId = ref<string | undefined>(undefined);
      const formData = ref({
        name: 'Test',
        description: '',
        version: '1.0.0',
        enabled: true,
        isTemplate: false,
      });
      const routerForms = ref<RouterFormState[]>([]);

      await handleSave(
        isSaving,
        isEditMode,
        routeGroupId,
        formData,
        routerForms,
        mockRouter as any,
        mockT,
      );

      expect(isSaving.value).toBe(false);
    });

    it('should not include empty description in payload', async () => {
      const isSaving = ref(false);
      const isEditMode = ref(false);
      const routeGroupId = ref<string | undefined>(undefined);
      const formData = ref({
        name: 'Test',
        description: '  ',
        version: '1.0.0',
        enabled: true,
        isTemplate: false,
      });
      const routerForms = ref<RouterFormState[]>([]);

      await handleSave(
        isSaving,
        isEditMode,
        routeGroupId,
        formData,
        routerForms,
        mockRouter as any,
        mockT,
      );

      const payload = (apis.router.routegroup.create as any).mock.calls[0][0];
      expect(payload.description).toBeUndefined();
    });

    it('should include match config for conditional routers', async () => {
      const isSaving = ref(false);
      const isEditMode = ref(false);
      const routeGroupId = ref<string | undefined>(undefined);
      const formData = ref({
        name: 'Test',
        description: '',
        version: '1.0.0',
        enabled: true,
        isTemplate: false,
      });
      const routerForms = ref<RouterFormState[]>([
        makeRouter({
          id: 'r1',
          hasConditionalRouting: true,
          match: {
            policy: 'all',
            rules: [{ field: 'x', operator: 'eq', value: '1' }],
          },
        }),
      ]);

      await handleSave(
        isSaving,
        isEditMode,
        routeGroupId,
        formData,
        routerForms,
        mockRouter as any,
        mockT,
      );

      const payload = (apis.router.routegroup.create as any).mock.calls[0][0];
      expect(payload.routers[0].match).toBeDefined();
    });
  });

  // ── handleLoadRouteGroup ─────────────────────────────────────────────
  describe('handleLoadRouteGroup', () => {
    const mockRouter = { push: vi.fn() };
    const mockT = {
      notifications: {
        loadFailed: { value: 'Load failed' },
      },
    };

    beforeEach(() => {
      vi.clearAllMocks();
    });

    it('should do nothing when not in edit mode', async () => {
      const isLoading = ref(false);
      const isEditMode = ref(false);
      const routeGroupId = ref<string | undefined>('rg-1');
      const formData = ref({} as any);
      const routerForms = ref<RouterFormState[]>([]);
      const currentStep = ref(1);

      await handleLoadRouteGroup(
        isLoading, isEditMode, routeGroupId, formData, routerForms, currentStep, mockRouter as any, mockT,
      );

      expect(apis.router.routegroup.getById).not.toHaveBeenCalled();
    });

    it('should populate form data and navigate to review step on success', async () => {
      vi.mocked(apis.router.routegroup.getById).mockResolvedValue({
        name: 'Loaded Group',
        description: 'Desc',
        enabled: true,
        isTemplate: false,
        routers: [],
      });

      const isLoading = ref(false);
      const isEditMode = ref(true);
      const routeGroupId = ref<string | undefined>('rg-1');
      const formData = ref({} as any);
      const routerForms = ref<RouterFormState[]>([]);
      const currentStep = ref(1);

      await handleLoadRouteGroup(
        isLoading, isEditMode, routeGroupId, formData, routerForms, currentStep, mockRouter as any, mockT,
      );

      expect(formData.value.name).toBe('Loaded Group');
      expect(currentStep.value).toBe(STEP.REVIEW);
      expect(isLoading.value).toBe(false);
    });

    it('should navigate back on error', async () => {
      vi.mocked(apis.router.routegroup.getById).mockRejectedValue(new Error('not found'));

      const isLoading = ref(false);
      const isEditMode = ref(true);
      const routeGroupId = ref<string | undefined>('rg-bad');
      const formData = ref({} as any);
      const routerForms = ref<RouterFormState[]>([]);
      const currentStep = ref(1);

      await handleLoadRouteGroup(
        isLoading, isEditMode, routeGroupId, formData, routerForms, currentStep, mockRouter as any, mockT,
      );

      expect(mockRouter.push).toHaveBeenCalledWith('/routing/route_groups');
      expect(isLoading.value).toBe(false);
    });
  });
});
