/** TYPE IMPORTS */
import type { AssetTemplateResponse } from '@mapexos/schemas';

/**
 * Props for AssetTemplateSelector component
 */
export interface AssetTemplateSelectorProps {
  /** Selected template ID (v-model) */
  modelValue: string | null;
}

/**
 * Emits for AssetTemplateSelector component
 */
export interface AssetTemplateSelectorEmits {
  (e: 'update:modelValue', value: string | null): void;
  (e: 'update:selectedTemplate', value: AssetTemplateResponse | null): void;
}
