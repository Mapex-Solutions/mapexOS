export interface AssetClassification {
  categoryId: string;
  manufacturerId: string;
  modelId: string;
  version: string;
  categoryName?: string | undefined;
  manufacturerName?: string | undefined;
  modelName?: string | undefined;
}

export interface AssetClassificationSelectorProps {
  modelValue?: AssetClassification | undefined;
  disabled?: boolean | undefined;
  required?: boolean | undefined;
}

export interface AssetClassificationSelectorEmits {
  (e: 'update:modelValue', value: AssetClassification | undefined): void;
}
