/** TYPE IMPORTS */
import type { AssetTemplateResponse } from '@mapexos/schemas';

/**
 * Props for AssetTemplateMultiSelector component
 */
export interface AssetTemplateMultiSelectorProps {
  /** Array of selected template IDs (v-model) */
  modelValue: string[];

  /** Label for the selector */
  label?: string;
}

/**
 * Emits for AssetTemplateMultiSelector component
 */
export interface AssetTemplateMultiSelectorEmits {
  (e: 'update:modelValue', value: string[]): void;
  (e: 'update:selectedTemplates', value: AssetTemplateResponse[]): void;
  (e: 'update:extractedPaths', value: Array<{ templateId: string; templateName: string; assetIdPath: string }>): void;
}
