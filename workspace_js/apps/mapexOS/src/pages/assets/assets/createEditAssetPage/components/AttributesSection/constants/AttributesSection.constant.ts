import type { AssetAttributeForm, AssetAttributeKind } from '../../../interfaces';

/** Selectable attribute kinds, in display order. Labels resolve from i18n. */
export const ATTRIBUTE_KIND_VALUES: AssetAttributeKind[] = [
  'string',
  'integer',
  'boolean',
  'date',
  'geo',
];

/**
 * The empty value for a kind, applied when a row is added or its kind changes.
 * @param {AssetAttributeKind} kind - The attribute kind.
 * @returns {AssetAttributeForm['value']} The zero value for that kind.
 */
export function defaultValueForKind(kind: AssetAttributeKind): AssetAttributeForm['value'] {
  switch (kind) {
    case 'boolean':
      return false;
    case 'geo':
      return { lat: null, lon: null };
    case 'integer':
      return null;
    case 'date':
    case 'string':
    default:
      return '';
  }
}
