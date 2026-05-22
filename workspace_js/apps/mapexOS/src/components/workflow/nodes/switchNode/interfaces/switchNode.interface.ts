/** TYPE IMPORTS */
import type { WorkflowConditionGroup } from '../../conditionNode/interfaces';

/**
 * Single case in a Switch node.
 * Each case has a root condition group and maps to one output handle.
 */
export interface SwitchCase {
  /** Unique ID (also used as output handle ID) */
  id: string;

  /** Display name (user-editable, shown on canvas output handle) */
  name: string;

  /** Root condition group to evaluate for this case */
  condition: WorkflowConditionGroup;
}
