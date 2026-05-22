/**
 * Comparison operators for condition evaluation
 */
export enum ComparisonOperator {
  Equals = 'equals',
  NotEquals = 'notEquals',
  Contains = 'contains',
  NotContains = 'notContains',
  StartsWith = 'startsWith',
  EndsWith = 'endsWith',
  GreaterThan = 'greaterThan',
  GreaterThanEquals = 'greaterThanEquals',
  LessThan = 'lessThan',
  LessThanEquals = 'lessThanEquals',
  In = 'in',
  NotIn = 'notIn',
  IsNull = 'isNull',
  IsNotNull = 'isNotNull',
  Regex = 'regex',
}

/**
 * Group logic options (AND/OR/NAND/NOR) for condition groups
 */
export const GROUP_LOGIC_OPTIONS = [
  { value: 'AND', label: 'All conditions', icon: 'check_circle', color: 'positive', description: 'All must be true' },
  { value: 'OR', label: 'Any condition', icon: 'check', color: 'primary', description: 'At least one must be true' },
  { value: 'NAND', label: 'Not all conditions', icon: 'highlight_off', color: 'orange', description: 'At least one must be false' },
  { value: 'NOR', label: 'No conditions', icon: 'cancel', color: 'negative', description: 'All must be false' },
] as const;

/**
 * Condition comparison operators with symbols for display
 */
export const CONDITION_OPERATOR_OPTIONS = [
  { value: ComparisonOperator.Equals, label: 'Equals', symbol: '=' },
  { value: ComparisonOperator.NotEquals, label: 'Not Equals', symbol: '≠' },
  { value: ComparisonOperator.Contains, label: 'Contains', symbol: '∋' },
  { value: ComparisonOperator.NotContains, label: 'Not Contains', symbol: '∌' },
  { value: ComparisonOperator.StartsWith, label: 'Starts With', symbol: '^' },
  { value: ComparisonOperator.EndsWith, label: 'Ends With', symbol: '$' },
  { value: ComparisonOperator.GreaterThan, label: 'Greater Than', symbol: '>' },
  { value: ComparisonOperator.GreaterThanEquals, label: 'Greater or Equal', symbol: '≥' },
  { value: ComparisonOperator.LessThan, label: 'Less Than', symbol: '<' },
  { value: ComparisonOperator.LessThanEquals, label: 'Less or Equal', symbol: '≤' },
  { value: ComparisonOperator.In, label: 'In (Array)', symbol: '∈' },
  { value: ComparisonOperator.NotIn, label: 'Not In (Array)', symbol: '∉' },
  { value: ComparisonOperator.IsNull, label: 'Is Null', symbol: '∅' },
  { value: ComparisonOperator.IsNotNull, label: 'Is Not Null', symbol: '!∅' },
  { value: ComparisonOperator.Regex, label: 'Regex', symbol: '~' },
] as const;
