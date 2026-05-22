/** TYPE IMPORTS */
import type { AssetResponse } from '@mapexos/schemas';

/**
 * Props for AssetSelector component
 */
export interface AssetSelectorProps {
  /** Selected asset ID (v-model) */
  modelValue: string | null;

  /** Label for the selector */
  label?: string;

  /** Whether selection is required */
  required?: boolean;
}

/**
 * Emits for AssetSelector component
 */
export interface AssetSelectorEmits {
  (e: 'update:modelValue', value: string | null): void;
  (e: 'update:selectedAsset', value: AssetResponse | null): void;
}
