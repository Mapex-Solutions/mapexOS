/** TYPE IMPORTS */
import type { WorkflowListPageFilters, WorkflowListPageColumnVisibility } from '../interfaces';

/** Default number of items per page */
export const DEFAULT_ITEMS_PER_PAGE = 15;

/** Default column visibility state */
export const WORKFLOW_COLUMN_VISIBILITY_DEFAULTS: WorkflowListPageColumnVisibility = {
  version: true,
  nodesCount: true,
  pluginsCount: true,
};

/** Default filter values */
export const WORKFLOW_FILTER_DEFAULTS: WorkflowListPageFilters = {
  name: undefined,
  status: undefined,
};

