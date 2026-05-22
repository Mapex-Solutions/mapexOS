/**
 * Organization Types
 *
 * Defines types for organizations in the MapexOS hierarchical structure.
 * Organizations have 6 levels: Vendor → Customer → Site → Building → Floor → Zone
 *
 * Each organization has:
 * - Type: level in hierarchy
 * - PathKey: hierarchical path (e.g., "000001/000001/0001")
 * - Depth: level number (0-5)
 * - ChildCount: number of children
 * - Enabled: active/inactive status
 */

export type OrganizationType = 'vendor' | 'customer' | 'site' | 'building' | 'floor' | 'zone';

/**
 * Base organization interface matching backend OrganizationResponse DTO
 */
export interface Organization {
  id: string;
  name: string;
  type: OrganizationType;
  parentOrgId?: string;
  pathKey: string;
  depth: number;
  childCount: number;
  enabled: boolean;
}

/**
 * Organization tree node for q-tree component
 * Extends Organization with tree-specific properties
 */
export interface OrganizationTreeNode extends Organization {
  children?: OrganizationTreeNode[];
  lazy?: boolean;
  header?: string;
}

/**
 * Organization filters for API queries
 * Maps to backend OrganizationQuery DTO
 */
export interface OrganizationFilters {
  name?: string;
  types?: OrganizationType[];
  enabled?: 'all' | 'active' | 'inactive';
  depth?: number;
}
