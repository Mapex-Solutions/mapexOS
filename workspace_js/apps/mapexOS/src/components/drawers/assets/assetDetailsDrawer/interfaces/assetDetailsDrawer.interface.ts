/**
 * AssetDetailsDrawer component interfaces
 * Displays comprehensive information about an asset in a side drawer
 */

/**
 * Props for AssetDetailsDrawer component
 * Controls drawer visibility and specifies which asset to display
 */
export interface AssetDetailsDrawerProps {
  /** Drawer visibility state */
  modelValue: boolean;

  /** ID of the asset to display details for */
  assetId: string | null;
}

/**
 * Emits for AssetDetailsDrawer component
 * Handles drawer state changes and user actions
 */
export interface AssetDetailsDrawerEmits {
  /** Update drawer visibility */
  (e: 'update:modelValue', value: boolean): void;

  /** Emit edit action with asset ID */
  (e: 'edit', assetId: string): void;

  /** Emit duplicate action with asset ID */
  (e: 'duplicate', assetId: string): void;
}
