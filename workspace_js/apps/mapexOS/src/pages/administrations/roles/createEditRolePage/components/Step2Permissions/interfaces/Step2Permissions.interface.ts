/** TYPE IMPORTS */
import type { ResourcePermission } from '../../../interfaces';

/**
 * Grouped resource item with its global index in the flat array
 */
export interface GroupedResourceItem {
  /** Resource permission object */
  resource: ResourcePermission;
  /** Index in the flat resourcePermissions array (for emitting events) */
  globalIndex: number;
}

/**
 * Grouped permissions section for rendering
 */
export interface GroupedPermissionSection {
  /** Group display label */
  label: string;
  /** Group icon */
  icon: string;
  /** Resources belonging to this group with their global indices */
  items: GroupedResourceItem[];
}
