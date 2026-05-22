import type { AssetTemplateResponse } from '@mapexos/schemas';

/**
 * Props for AssetTemplateSelectorDialog component
 * Centered modal for selecting asset templates (used in workflow context)
 */
export interface AssetTemplateSelectorDialogProps {
  /** Dialog visibility state */
  modelValue: boolean;

  /** Array of pre-selected template IDs */
  selectedTemplateIds?: string[];

  /** Enable multi-select mode (default: true) */
  multiSelect?: boolean;
}

/**
 * Emits for AssetTemplateSelectorDialog component
 */
export interface AssetTemplateSelectorDialogEmits {
  /** Update dialog visibility */
  (e: 'update:modelValue', value: boolean): void;

  /** Emit selected templates */
  (e: 'select', templates: AssetTemplateResponse[]): void;

  /** Emit cancel action */
  (e: 'cancel'): void;
}
