import type { AssetTemplateResponse } from '@mapexos/schemas';

/**
 * A single dynamic field entry of an asset template.
 * Derived from the response schema, which does not export a standalone type.
 */
export type DynamicField = NonNullable<AssetTemplateResponse['dynamicFields']>[number];

/**
 * Read-only dynamic fields table props interface.
 */
export interface DynamicFieldsTableProps {
  /** Dynamic fields to render */
  fields: DynamicField[];
}
