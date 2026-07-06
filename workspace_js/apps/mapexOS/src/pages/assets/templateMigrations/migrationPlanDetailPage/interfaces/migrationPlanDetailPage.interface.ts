// MigrationPlanDetailPage Interfaces

import type { MigrationExecutionStatus } from '@mapexos/schemas';


/**
 * Option entry for the executions status filter select
 */
export interface MigrationExecutionStatusOption {
  label: string;
  value: MigrationExecutionStatus | null;
}
