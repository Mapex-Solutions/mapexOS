/**
 * DynamicFiltersDrawer Constants
 */
import type { DynamicSourceType, EvaOperatorOption } from '../interfaces';

/** Default source type when drawer opens */
export const DEFAULT_SOURCE_TYPE: DynamicSourceType = 'asset';

/**
 * Map EVA dynamic field types to DynamicFilterField types
 * Converts backend type names to frontend filter field types
 */
export const EVA_TYPE_MAP: Record<string, 'string' | 'number' | 'boolean' | 'date'> = {
  string: 'string',
  number: 'number',
  bool: 'boolean',
  date: 'date',
} as const;

/** Minimum characters required before autocomplete search triggers */
export const AUTOCOMPLETE_MIN_CHARS = 2;

/** Debounce delay in ms for autocomplete search */
export const AUTOCOMPLETE_DEBOUNCE = 300;

/** Number of results per page for autocomplete */
export const AUTOCOMPLETE_PER_PAGE = 20;

/** Default operator per field type */
export const DEFAULT_OPERATOR_BY_TYPE: Record<string, string> = {
  number: 'eq',
  string: 'eq',
  boolean: 'eq',
  date: 'gte',
} as const;

/** Operators available for each EVA field type */
export const EVA_OPERATORS_BY_TYPE: Record<string, EvaOperatorOption[]> = {
  number: [
    { value: 'eq', labelKey: 'operators.equals', icon: 'drag_handle' },
    { value: 'neq', labelKey: 'operators.notEquals', icon: 'not_equal' },
    { value: 'gt', labelKey: 'operators.greaterThan', icon: 'chevron_right' },
    { value: 'gte', labelKey: 'operators.greaterThanEquals', icon: 'keyboard_double_arrow_right' },
    { value: 'lt', labelKey: 'operators.lessThan', icon: 'chevron_left' },
    { value: 'lte', labelKey: 'operators.lessThanEquals', icon: 'keyboard_double_arrow_left' },
    { value: 'between', labelKey: 'operators.range', icon: 'swap_horiz' },
  ],
  string: [
    { value: 'eq', labelKey: 'operators.equals', icon: 'drag_handle' },
    { value: 'neq', labelKey: 'operators.notEquals', icon: 'not_equal' },
    { value: 'like', labelKey: 'operators.startsWith', icon: 'start' },
  ],
  boolean: [
    { value: 'eq', labelKey: 'operators.equals', icon: 'drag_handle' },
  ],
  date: [
    { value: 'eq', labelKey: 'operators.equals', icon: 'drag_handle' },
    { value: 'neq', labelKey: 'operators.notEquals', icon: 'not_equal' },
    { value: 'gte', labelKey: 'operators.greaterThanEquals', icon: 'keyboard_double_arrow_right' },
    { value: 'lte', labelKey: 'operators.lessThanEquals', icon: 'keyboard_double_arrow_left' },
    { value: 'between', labelKey: 'operators.range', icon: 'swap_horiz' },
  ],
} as const;
