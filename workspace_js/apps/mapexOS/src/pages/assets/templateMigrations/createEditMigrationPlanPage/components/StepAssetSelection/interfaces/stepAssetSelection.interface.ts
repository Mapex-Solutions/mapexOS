/**
 * Props for StepAssetSelection.
 * - fromTemplateId scopes the asset list to a single source template.
 * - assetIds is the selected asset id slice of the form model (v-model:assetIds).
 */
export interface StepAssetSelectionProps {
  fromTemplateId: string | null;
  assetIds: string[];
}

/**
 * Emits for StepAssetSelection.
 * - update:assetIds emits the selected asset ids.
 * - update:valid signals that at least one asset is selected.
 */
export interface StepAssetSelectionEmits {
  (e: 'update:assetIds', value: string[]): void;
  (e: 'update:valid', value: boolean): void;
}

/** A single selectable row in the asset list. */
export interface SelectableAsset {
  id: string;
  name: string;
  assetUUID: string;
}
