/** TYPE IMPORTS */
import type { RouteGroupResponse } from '@mapexos/schemas';

/**
 * Props for RouteGroupSelector component
 */
export interface RouteGroupSelectorProps {
  /** Array of selected route group IDs (v-model) */
  modelValue: string[];
}

/**
 * Emits for RouteGroupSelector component
 */
export interface RouteGroupSelectorEmits {
  (e: 'update:modelValue', value: string[]): void;
  (e: 'update:selectedRouteGroups', value: RouteGroupResponse[]): void;
}
