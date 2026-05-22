import type { Ref } from 'vue';
import type { Router as VueRouter } from 'vue-router';
import type {
  RouteGroupCreate,
  MatchRule,
} from '@interfaces/routing/routeGroups.interface';
import type { RouterFormState } from '../interfaces';
import { apis } from '@services/mapex';
import { notifySuccess } from '@utils/alert/notify';
import { handleApiError } from '@utils/error';
import { useLogger } from '@composables/useLogger';
import {
  DEFAULT_ROUTER_KIND,
  DEFAULT_MATCH_POLICY,
  DEFAULT_MATCH_OPERATOR,
  STEP,
} from '../constants';

const logger = useLogger('routeGroupHandler');

/**
 * Add a new router to the router forms array
 *
 * @param {Ref<RouterFormState[]>} routerForms - Router forms array ref
 * @returns {void}
 */
export function handleAddRouter(routerForms: Ref<RouterFormState[]>): void {
  const newRouter: RouterFormState = {
    id: `router-${Date.now()}-${Math.random()}`,
    kind: DEFAULT_ROUTER_KIND,
    hasConditionalRouting: false,
    saveEvent: {},
  };
  routerForms.value.push(newRouter);
}

/**
 * Remove a router from the router forms array
 *
 * @param {Ref<RouterFormState[]>} routerForms - Router forms array ref
 * @param {string} routerId - ID of the router to remove
 * @returns {void}
 */
export function handleRemoveRouter(
  routerForms: Ref<RouterFormState[]>,
  routerId: string,
): void {
  routerForms.value = routerForms.value.filter((r) => r.id !== routerId);
}

/**
 * Handle router kind change and initialize appropriate destination config
 *
 * @param {Ref<RouterFormState[]>} routerForms - Router forms array ref
 * @param {string} routerId - ID of the router
 * @param {string} newKind - New router kind value
 * @returns {void}
 */
export function handleRouterKindChange(
  routerForms: Ref<RouterFormState[]>,
  routerId: string,
  newKind: string,
): void {
  const routerForm = routerForms.value.find((r) => r.id === routerId);
  if (!routerForm) return;

  // Clear all destination configs
  delete routerForm.lakeHouse;
  delete routerForm.notification;
  delete routerForm.saveEvent;
  delete routerForm.workflow;

  // Initialize the appropriate destination config
  if (newKind === 'lake_house') {
    routerForm.lakeHouse = { lakeHouseId: '', metadata: {} };
  } else if (newKind === 'notification') {
    routerForm.notification = { notificationId: '', metadata: {} };
  } else if (newKind === 'save_event') {
    routerForm.saveEvent = { metadata: {} };
  } else if (newKind === 'workflow') {
    routerForm.workflow = { mode: 'newInstance', data: {} };
  }
}

/**
 * Add a match rule to a router's conditional routing configuration
 *
 * @param {Ref<RouterFormState[]>} routerForms - Router forms array ref
 * @param {string} routerId - ID of the router
 * @returns {void}
 */
export function handleAddMatchRule(
  routerForms: Ref<RouterFormState[]>,
  routerId: string,
): void {
  const routerForm = routerForms.value.find((r) => r.id === routerId);
  if (!routerForm) return;

  if (!routerForm.match) {
    routerForm.match = {
      policy: DEFAULT_MATCH_POLICY,
      rules: [],
    };
  }

  const newRule: MatchRule = {
    field: '',
    operator: DEFAULT_MATCH_OPERATOR,
    value: '',
  };

  routerForm.match.rules.push(newRule);
}

/**
 * Remove a match rule from a router's conditional routing configuration
 *
 * @param {Ref<RouterFormState[]>} routerForms - Router forms array ref
 * @param {string} routerId - ID of the router
 * @param {number} ruleIndex - Index of the rule to remove
 * @returns {void}
 */
export function handleRemoveMatchRule(
  routerForms: Ref<RouterFormState[]>,
  routerId: string,
  ruleIndex: number,
): void {
  const routerForm = routerForms.value.find((r) => r.id === routerId);
  if (!routerForm || !routerForm.match) return;

  routerForm.match.rules.splice(ruleIndex, 1);

  // Clean up match config if no rules left
  if (routerForm.match.rules.length === 0) {
    routerForm.hasConditionalRouting = false;
    delete routerForm.match;
  }
}

/**
 * Toggle conditional routing for a router
 *
 * @param {Ref<RouterFormState[]>} routerForms - Router forms array ref
 * @param {string} routerId - ID of the router
 * @param {boolean} enabled - Whether to enable conditional routing
 * @returns {void}
 */
export function handleToggleConditionalRouting(
  routerForms: Ref<RouterFormState[]>,
  routerId: string,
  enabled: boolean,
): void {
  const routerForm = routerForms.value.find((r) => r.id === routerId);
  if (!routerForm) return;

  routerForm.hasConditionalRouting = enabled;

  if (enabled && !routerForm.match) {
    // Add ONE default rule when enabling conditional routing
    const defaultRule: MatchRule = {
      field: '',
      operator: DEFAULT_MATCH_OPERATOR,
      value: '',
    };

    routerForm.match = {
      policy: DEFAULT_MATCH_POLICY,
      rules: [defaultRule],
    };
  } else if (!enabled) {
    delete routerForm.match;
  }
}

/**
 * Change the current step in the form
 * Validates current step before allowing navigation in create mode
 *
 * @param {Ref<number>} currentStep - Current step ref
 * @param {number} step - Target step number
 * @param {Ref<any>} step1FormRef - Step 1 form ref for validation
 * @param {Ref<RouterFormState[]>} routerForms - Router forms array ref
 * @returns {Promise<void>}
 */
export async function handleChangeStep(
  currentStep: Ref<number>,
  step: number,
  step1FormRef: Ref<any>,
  routerForms: Ref<RouterFormState[]>,
): Promise<void> {
  // Validate current step before proceeding
  if (currentStep.value === 1 && step > 1) {
    if (step1FormRef.value) {
      const valid = await step1FormRef.value.validate();
      if (!valid) return;
    }
  }

  // Validate Step 2 (at least one router)
  if (currentStep.value === 2 && step > 2) {
    if (routerForms.value.length === 0) {
      // Show warning
      return;
    }
    // Validate each router
    for (const routerForm of routerForms.value) {
      if (routerForm.hasConditionalRouting && routerForm.match) {
        if (routerForm.match.rules.length === 0) {
          // Show error: must have at least one rule
          return;
        }
      }
    }
  }

  currentStep.value = step;
}

/**
 * Save or update the route group
 *
 * @param {Ref<boolean>} isSaving - Saving state ref
 * @param {Ref<boolean>} isEditMode - Edit mode ref
 * @param {Ref<string | undefined>} routeGroupId - Route group ID ref (for edit mode)
 * @param {Ref<RouteGroupCreate>} formData - Form data ref
 * @param {Ref<RouterFormState[]>} routerForms - Router forms array ref
 * @param {VueRouter} router - Vue router instance
 * @param {any} t - Translation object
 * @returns {Promise<void>}
 */
export async function handleSave(
  isSaving: Ref<boolean>,
  isEditMode: Ref<boolean>,
  routeGroupId: Ref<string | undefined>,
  formData: Ref<RouteGroupCreate>,
  routerForms: Ref<RouterFormState[]>,
  router: VueRouter,
  t: any,
): Promise<void> {
  isSaving.value = true;

  try {
    // Prepare routers data
    const routers: any[] = routerForms.value.map((rf) => {
      const routerObj: any = {
        kind: rf.kind,
      };

      // Add match config if conditional routing is enabled
      if (rf.hasConditionalRouting && rf.match) {
        routerObj.match = rf.match;
      }

      // Add destination config based on kind
      if (rf.kind === 'lake_house' && rf.lakeHouse) {
        routerObj.lakeHouse = rf.lakeHouse;
      } else if (rf.kind === 'notification' && rf.notification) {
        routerObj.notification = rf.notification;
      } else if (rf.kind === 'save_event') {
        routerObj.saveEvent = rf.saveEvent || {};
      } else if (rf.kind === 'workflow' && rf.workflow) {
        routerObj.workflow = rf.workflow;
      }

      return routerObj;
    });

    // Build payload
    const payload: any = {
      name: formData.value.name,
      version: '1.0.0',
      enabled: formData.value.enabled,
      isTemplate: formData.value.isTemplate,
    };

    // Only add description if not empty
    if (formData.value.description && formData.value.description.trim() !== '') {
      payload.description = formData.value.description;
    }

    // Only add routers if there are any
    if (routers.length > 0) {
      payload.routers = routers;
    }

    if (isEditMode.value && routeGroupId.value) {
      // Update existing route group
      await apis.router.routegroup.update(
        { routeGroupId: routeGroupId.value },
        payload,
      );

      notifySuccess({
        message: t.notifications.updated.value,
        timeout: 3000,
      });
    } else {
      // Create new route group
      await apis.router.routegroup.create(payload);

      notifySuccess({
        message: t.notifications.created.value,
        timeout: 3000,
      });
    }

    // Navigate to list
    await router.push('/routing/route_groups');
  } catch (error: any) {
    handleApiError(error, {
      customMessages: {
        409: t.notifications.alreadyExists,
        422: t.notifications.validationFailed,
        network: t.notifications.networkError,
      },
      defaultMessage: isEditMode.value
        ? t.notifications.updateFailed
        : t.notifications.creationFailed,
      timeout: 5000,
    });
  } finally {
    isSaving.value = false;
  }
}

/**
 * Fetch trigger/notification name by ID
 *
 * @param {string} triggerId - Trigger ID
 * @returns {Promise<string | undefined>} Trigger name or undefined
 */
async function fetchTriggerName(triggerId: string): Promise<string | undefined> {
  try {
    const response = await apis.triggers.trigger.getById({ triggerId });
    return response.name;
  } catch {
    logger.warn('Failed to fetch trigger name:', triggerId);
    return undefined;
  }
}

/**
 * Load route group data for edit mode
 *
 * @param {Ref<boolean>} isLoading - Loading state ref
 * @param {Ref<boolean>} isEditMode - Edit mode ref
 * @param {Ref<string | undefined>} routeGroupId - Route group ID ref
 * @param {Ref<RouteGroupCreate>} formData - Form data ref
 * @param {Ref<RouterFormState[]>} routerForms - Router forms array ref
 * @param {Ref<number>} currentStep - Current step ref
 * @param {VueRouter} router - Vue router instance
 * @param {any} t - Translation object
 * @returns {Promise<void>}
 */
export async function handleLoadRouteGroup(
  isLoading: Ref<boolean>,
  isEditMode: Ref<boolean>,
  routeGroupId: Ref<string | undefined>,
  formData: Ref<RouteGroupCreate>,
  routerForms: Ref<RouterFormState[]>,
  currentStep: Ref<number>,
  router: VueRouter,
  t: any,
): Promise<void> {
  if (!isEditMode.value || !routeGroupId.value) return;

  isLoading.value = true;
  try {
    const data = await apis.router.routegroup.getById({
      routeGroupId: routeGroupId.value,
    });

    // Populate form data
    formData.value = {
      name: data.name || '',
      description: data.description || '',
      version: '1.0.0',
      enabled: data.enabled ?? true,
      isTemplate: data.isTemplate ?? false,
    };

    // Convert routers to form state
    if (data.routers && data.routers.length > 0) {
      const routerFormsData: RouterFormState[] = [];

      for (let index = 0; index < data.routers.length; index++) {
        const routerData = data.routers[index] as any;

        const routerForm: RouterFormState = {
          id: `router-${Date.now()}-${index}`,
          kind: routerData.kind,
          hasConditionalRouting: !!routerData.match,
          match: routerData.match,
          lakeHouse: routerData.lakeHouse,
          notification: routerData.notification,
          saveEvent: routerData.saveEvent,
          workflow: routerData.workflow,
        };

        // Fetch names for selected destinations
        if (routerData.kind === 'notification' && routerData.notification?.notificationId) {
          const name = await fetchTriggerName(routerData.notification.notificationId);
          if (name) routerForm.notificationName = name;
        }

        // TODO: Fetch data lake name when API is available
        // if (routerData.kind === 'lake_house' && routerData.lakeHouse?.lakeHouseId) {
        //   routerForm.lakeHouseName = await fetchLakeHouseName(routerData.lakeHouse.lakeHouseId);
        // }

        routerFormsData.push(routerForm);
      }

      routerForms.value = routerFormsData;
    }

    // In EDIT mode, skip to Review step by default
    currentStep.value = STEP.REVIEW;

  } catch (error: any) {
    handleApiError(error, {
      defaultMessage: t.notifications.loadFailed.value,
      timeout: 5000,
    });
    void router.push('/routing/route_groups');
  } finally {
    isLoading.value = false;
  }
}
