import type { AssetTemplateData } from '../../../interfaces';

export interface Step2AssetIdPathProps {
  modelValue: AssetTemplateData;
}

export interface Step2AssetIdPathEmits {
  (e: 'update:modelValue', value: AssetTemplateData): void;
}
