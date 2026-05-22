/**
 * AssetSelectorDrawer component interfaces
 * Provides a dialog for selecting a single asset with filtering capabilities
 */
import type { AssetResponse } from '@mapexos/schemas';

/**
 * Props for AssetSelectorDrawer component
 * Controls drawer visibility and manages asset selection
 */
export interface AssetSelectorDrawerProps {
  /** Dialog/drawer visibility state */
  modelValue: boolean;

  /** ID of the currently selected asset (optional) */
  selectedAssetId?: string | null;

  /** Enable multi-select mode (future feature, currently defaults to false) */
  multiSelect?: boolean;

  /** Pre-filter by asset template ID — only shows assets of this template type */
  assetTemplateId?: string;
}

/**
 * Emits for AssetSelectorDrawer component
 * Handles drawer state changes and asset selection events
 */
export interface AssetSelectorDrawerEmits {
  /** Update drawer visibility */
  (e: 'update:modelValue', value: boolean): void;

  /** Emit selected asset */
  (e: 'select', asset: AssetResponse): void;

  /** Emit cancel action */
  (e: 'cancel'): void;
}
