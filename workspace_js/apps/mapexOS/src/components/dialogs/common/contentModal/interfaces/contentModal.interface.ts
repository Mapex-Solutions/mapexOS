/**
 * Generic content modal props interface
 * A container-only dialog whose body is provided by the consumer via slot.
 */
export interface ContentModalProps {
  /** Controls the open/close state (v-model) */
  modelValue: boolean;

  /** Header title text */
  title: string;

  /** Optional header icon name (Material Icons) */
  icon?: string | undefined;
}

/**
 * Generic content modal emits interface
 */
export interface ContentModalEmits {
  (e: 'update:modelValue', value: boolean): void;
}
