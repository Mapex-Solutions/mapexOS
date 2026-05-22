/**
 * RouteGroupDetailsDrawer Interfaces
 */

/**
 * Props for RouteGroupDetailsDrawer component
 */
export interface RouteGroupDetailsDrawerProps {
  /** Whether the drawer is open */
  modelValue: boolean;

  /** Route Group ID to display */
  routeGroupId?: string | undefined;
}

/**
 * Emits for RouteGroupDetailsDrawer component
 */
export interface RouteGroupDetailsDrawerEmits {
  (e: 'update:modelValue', value: boolean): void;
}

/**
 * Extended router type with match field from API response
 */
export interface RouterWithMatch {
  kind: string;
  match?: {
    policy: 'all' | 'any';
    rules?: Array<{ field: string; operator: string; value: unknown }>;
  };
  lakeHouse?: { lakeHouseId: string };
  notification?: { notificationId: string };
  saveEvent?: Record<string, unknown>;
}
