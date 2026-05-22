/**
 * Constants for CreateEditRouteGroupPage
 */

/** Default enabled status for new route groups */
export const DEFAULT_ENABLED = true;

/** Default isTemplate status for new route groups */
export const DEFAULT_IS_TEMPLATE = false;

/** Default router kind */
export const DEFAULT_ROUTER_KIND = 'save_event';

/** Total number of steps in the form */
export const TOTAL_STEPS = 3;

/** Initial step number */
export const INITIAL_STEP = 1;

/** Default match policy */
export const DEFAULT_MATCH_POLICY = 'all';

/** Default match operator */
export const DEFAULT_MATCH_OPERATOR = 'eq';

/** Step enum for clearer step navigation */
export const STEP = {
  BASIC_INFO: 1,
  ROUTERS_CONFIG: 2,
  REVIEW: 3,
} as const;
