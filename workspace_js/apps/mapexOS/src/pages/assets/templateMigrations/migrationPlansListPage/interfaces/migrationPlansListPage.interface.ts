// MigrationPlansListPage Interfaces

import type { MigrationPlanStatus } from '@mapexos/schemas';

/**
 * Filter state for the migration plans list page
 */
export interface MigrationPlansListPageFilters {
  name: string | undefined;
  status: MigrationPlanStatus | undefined;
}

/**
 * Option entry for the status quick-filter select
 */
export interface MigrationStatusOption {
  label: string;
  value: MigrationPlanStatus | null;
}
