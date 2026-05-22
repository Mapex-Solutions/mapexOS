/** TYPE IMPORTS */
import type { FieldSourceValue } from '@src/components/workflow/interfaces';
import type { ComparisonOperator } from '../constants/conditionNode.constant';

/** Logical operator for combining items within a group */
export type GroupLogicOperator = 'AND' | 'OR' | 'NAND' | 'NOR';

/**
 * Single condition (field + operator + value)
 */
export interface WorkflowConditionItem {
  /** Unique ID */
  id: string;

  /** Display name (user-editable) */
  name: string;

  /** Left side of comparison */
  field: FieldSourceValue;

  /** Comparison operator */
  operator: ComparisonOperator;

  /** Right side of comparison */
  value: FieldSourceValue;
}

/**
 * Discriminated union: an item inside a group is either a condition or a sub-group
 */
export type ConditionGroupItem =
  | { type: 'condition'; data: WorkflowConditionItem }
  | { type: 'group'; data: WorkflowConditionGroup };

/**
 * Condition group containing mixed items (conditions and/or sub-groups).
 * Sub-groups should only contain conditions (max 2 levels deep).
 */
export interface WorkflowConditionGroup {
  /** Unique ID (required for sub-groups, omitted at root level) */
  id?: string;

  /** Display name (required for sub-groups, omitted at root level) */
  name?: string;

  /** Logical operator for combining items */
  logic: GroupLogicOperator;

  /** Items in this group (conditions or sub-groups) */
  items: ConditionGroupItem[];
}
