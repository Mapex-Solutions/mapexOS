/**
 * Re-export canonical constants — SOURCE_TYPE_OPTIONS and DEFAULT_FIELD_SOURCE_VALUE
 * are defined in @src/components/workflow/constants for cross-plugin reusability.
 */
export { SOURCE_TYPE_OPTIONS, DEFAULT_FIELD_SOURCE_VALUE } from '@src/components/workflow/constants';

/**
 * Re-export SourceTypeOption interface for consumers that import from this file
 */
export type { SourceTypeOption } from '@src/components/workflow/interfaces';
