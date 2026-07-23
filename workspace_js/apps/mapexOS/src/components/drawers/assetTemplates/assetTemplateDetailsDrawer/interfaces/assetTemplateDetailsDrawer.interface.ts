/**
 * AssetTemplateDetailsDrawer component interfaces
 * Displays detailed information about an asset template in a side drawer
 */

/**
 * Props for AssetTemplateDetailsDrawer component
 * Controls drawer visibility and specifies which template to display
 */
export interface AssetTemplateDetailsDrawerProps {
  /** Dialog/drawer visibility state */
  modelValue: boolean;

  /** ID of the asset template to display details for */
  templateId: string | null;
}

/**
 * Emits for AssetTemplateDetailsDrawer component
 * Handles drawer state changes and user actions
 */
export interface AssetTemplateDetailsDrawerEmits {
  /** Update drawer visibility */
  (e: 'update:modelValue', value: boolean): void;

  /** Emit edit action with template ID */
  (e: 'edit', templateId: string): void;

  /** Emit after a marketplace template is cloned into a new local template. */
  (e: 'cloned', templateId: string): void;
}
