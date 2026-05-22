/**
 * PersistenceSettings Component Interfaces
 *
 * This file defines TypeScript interfaces for the PersistenceSettings component.
 * PersistenceSettings is now self-contained: it fetches and saves retention
 * policies directly via the events API (no props/emits needed).
 *
 * @see src/pages/administrations/settings/components/PersistenceSettings.vue
 */

/**
 * Retention policy limits for a specific event type
 */
export interface RetentionPolicyLimits {
  defaultDays: number;
  minDays: number;
  maxDays: number;
  name: string;
}

/**
 * Retention policy configuration for display and editing
 */
export interface RetentionPolicy {
  key: string;
  name: string;
  currentDays: number;
  defaultDays: number;
  minDays: number;
  maxDays: number;
}
