import type { AssetTemplateData } from '../../../interfaces';

export interface Step4ValidationScriptProps {
  modelValue: AssetTemplateData;
}

export interface Step4ValidationScriptEmits {
  (e: 'update:modelValue', value: AssetTemplateData): void;
}
