import type { AssetClassification } from '../../../interfaces';

export interface GuidedModeSelectorProps {
  modelValue?: AssetClassification | undefined;
  disabled?: boolean | undefined;
  required?: boolean | undefined;
}

export interface GuidedModeSelectorEmits {
  (e: 'update:modelValue', value: AssetClassification | undefined): void;
}

export interface ListOption {
  id: string;
  name: string;
  value: string;
}
