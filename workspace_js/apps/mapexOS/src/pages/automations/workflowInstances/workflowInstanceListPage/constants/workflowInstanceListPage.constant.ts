/** TYPE IMPORTS */
import type {
  WorkflowInstanceListPageFilters,
  WorkflowInstanceListPageColumnVisibility,
} from '../interfaces';

/** Default number of items per page */
export const DEFAULT_ITEMS_PER_PAGE = 15;

/** Default column visibility state */
export const INSTANCE_COLUMN_VISIBILITY_DEFAULTS: WorkflowInstanceListPageColumnVisibility = {
  definitionName: true,
  inputsCount: true,
  uniqueExecution: true,
};

/** Default filter values */
export const INSTANCE_FILTER_DEFAULTS: WorkflowInstanceListPageFilters = {
  name: undefined,
  status: undefined,
};

/** Projection fields for list API call (only fields needed for the list view) */
export const LIST_PROJECTION = 'name,description,enabled,definitionName,externalInputs,uniqueExecution,workflowUUID';
