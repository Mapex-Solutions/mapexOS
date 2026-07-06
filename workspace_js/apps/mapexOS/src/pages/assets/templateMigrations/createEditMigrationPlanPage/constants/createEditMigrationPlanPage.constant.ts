/**
 * CreateEditMigrationPlanPage Constants
 */

import type { MigrationPlanForm } from '../interfaces';

/**
 * Base route for migration plan navigation
 */
export const MIGRATION_ROUTE_BASE = '/assets_template/migrations';

/**
 * Empty starting state for the create migration plan wizard
 */
export const DEFAULT_MIGRATION_PLAN_FORM: MigrationPlanForm = {
  name: '',
  description: '',
  scheduleAt: null,
  fromTemplateId: null,
  toTemplateId: null,
  assetIds: [],
};

/**
 * Total number of wizard steps: Details, Schedule, From template, Assets, To template
 */
export const MIGRATION_WIZARD_TOTAL_STEPS = 5;

/**
 * Total steps in edit mode: only Details and Schedule are mutable after creation.
 */
export const MIGRATION_EDIT_TOTAL_STEPS = 2;

/**
 * Plan statuses that still allow editing (the backend rejects edits otherwise).
 */
export const MIGRATION_EDITABLE_STATUSES = ['pending', 'scheduled'];
