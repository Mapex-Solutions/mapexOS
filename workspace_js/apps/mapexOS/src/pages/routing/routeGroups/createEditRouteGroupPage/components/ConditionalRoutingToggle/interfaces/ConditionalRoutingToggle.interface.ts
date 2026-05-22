/**
 * Props for ConditionalRoutingToggle component
 */
export interface ConditionalRoutingToggleProps {
  /** Current toggle state */
  modelValue: boolean;

  /** Translation composable */
  t: any;
}

/**
 * Emits for ConditionalRoutingToggle component
 */
export interface ConditionalRoutingToggleEmits {
  (e: 'update:modelValue', value: boolean): void;
}
