/** TYPE IMPORTS */
import type { RouteGroupResponse } from '@mapexos/schemas';

/**
 * Props for RouteGroupSelectorDrawer component
 */
export interface RouteGroupSelectorDrawerProps {
  /** Dialog visibility (v-model) */
  modelValue: boolean;

  /** Pre-selected route group IDs */
  selectedRouteGroupIds?: string[];

  /** Whether multiple selection is allowed */
  multiSelect?: boolean;

  /**
   * Optional whitelist of router kinds. When set, only RouteGroups whose every
   * router.kind is contained in this array are listed (strict semantic).
   * When undefined, all RouteGroups are listed (default).
   * Used by the asset wizard's Health step with ['trigger', 'workflow'].
   */
  allowedRouterKinds?: string[];
}

/**
 * Emits for RouteGroupSelectorDrawer component
 */
export interface RouteGroupSelectorDrawerEmits {
  (e: 'update:modelValue', value: boolean): void;
  (e: 'select', routeGroups: RouteGroupResponse[]): void;
  (e: 'cancel'): void;
}
