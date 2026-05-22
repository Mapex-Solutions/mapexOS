/**
 * AssetTemplateSelectorDrawer component interfaces
 * Provides a dialog for selecting one or multiple asset templates with filtering capabilities
 */
import type { AssetTemplateResponse } from '@mapexos/schemas';

/**
 * Props for AssetTemplateSelectorDrawer component
 * Allows selecting one or multiple asset templates with pre-selected items
 */
export interface AssetTemplateSelectorDrawerProps {
  /** Dialog visibility state */
  modelValue: boolean;

  /** Array of pre-selected template IDs */
  selectedTemplateIds?: string[];

  /** Enable multi-select mode (default: true) */
  multiSelect?: boolean;
}

/**
 * Emits for AssetTemplateSelectorDrawer component
 * Handles dialog state and template selection events
 */
export interface AssetTemplateSelectorDrawerEmits {
  /** Update dialog visibility */
  (e: 'update:modelValue', value: boolean): void;

  /** Emit selected templates (array for multi-select, single item for single-select) */
  (e: 'select', templates: AssetTemplateResponse[]): void;

  /** Emit cancel action */
  (e: 'cancel'): void;
}
