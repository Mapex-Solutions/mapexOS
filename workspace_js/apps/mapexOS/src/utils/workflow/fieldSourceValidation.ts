/** TYPE IMPORTS */
import type { FieldSourceValue } from '@src/components/workflow/interfaces';

/**
 * Check if a field source value is incomplete (type selected but value missing).
 * This is the canonical validation for any FieldSourceValue across all plugins.
 *
 * @param {FieldSourceValue | undefined} source - Field source object
 * @returns {boolean} True if the source is empty or incomplete
 */
export function isFieldSourceEmpty(source: FieldSourceValue | undefined): boolean {
  if (!source) return true;
  if (!source.type) return true;

  const val = source.value?.toString().trim() ?? '';
  if (!val) return true;

  // nodeOutput requires nodeId too
  if (source.type === 'nodeOutput' && !source.nodeId?.trim()) return true;

  // fetchOptions with empty value is valid (user may not have selected yet)
  if (source.type === 'fetchOptions') return !source.value?.toString().trim();

  return false;
}

/**
 * Recursively validate all condition items in a group/items structure.
 * Checks that each condition's field and value sources are complete.
 *
 * Works with the `ConditionGroupItem` discriminated union pattern:
 * `{ type: 'condition', data: { field, value, ... } }` or
 * `{ type: 'group', data: { items: [...] } }`
 *
 * @param {unknown[]} items - Array of ConditionGroupItem-like objects
 * @param {string[]} errors - Error array to push into
 * @param {string} prefix - Label prefix for error messages
 */
export function validateConditionItems(items: unknown[], errors: string[], prefix: string): void {
  for (const raw of items) {
    const item = raw as { type?: string; data?: Record<string, unknown> };
    if (!item?.data) continue;

    if (item.type === 'condition') {
      const field = item.data.field as FieldSourceValue | undefined;
      const value = item.data.value as FieldSourceValue | undefined;
      const name = (item.data.name as string) || 'Condition';

      if (isFieldSourceEmpty(field)) {
        errors.push(`${prefix}${name}: field source is incomplete`);
      }
      if (isFieldSourceEmpty(value)) {
        errors.push(`${prefix}${name}: comparison value is incomplete`);
      }
    } else if (item.type === 'group') {
      const subItems = item.data.items as unknown[] | undefined;
      const groupName = (item.data.name as string) || 'Group';
      if (subItems?.length) {
        validateConditionItems(subItems, errors, `${prefix}${groupName} > `);
      }
    }
  }
}
