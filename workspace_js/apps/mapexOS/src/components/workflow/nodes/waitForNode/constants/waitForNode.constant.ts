/**
 * Comparison operators for Wait For condition evaluation
 */
export const WAIT_FOR_OPERATORS = [
  { value: 'equals', label: 'Equals', symbol: '=' },
  { value: 'notEquals', label: 'Not Equals', symbol: '≠' },
  { value: 'greaterThan', label: 'Greater Than', symbol: '>' },
  { value: 'greaterThanEquals', label: 'Greater or Equal', symbol: '≥' },
  { value: 'lessThan', label: 'Less Than', symbol: '<' },
  { value: 'lessThanEquals', label: 'Less or Equal', symbol: '≤' },
  { value: 'contains', label: 'Contains', symbol: '∋' },
  { value: 'notContains', label: 'Not Contains', symbol: '∌' },
  { value: 'isEmpty', label: 'Is Empty', symbol: '∅' },
  { value: 'isNotEmpty', label: 'Is Not Empty', symbol: '!∅' },
  { value: 'isTrue', label: 'Is True', symbol: '✓' },
  { value: 'isFalse', label: 'Is False', symbol: '✗' },
] as const;
