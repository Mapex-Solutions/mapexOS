import type { AssetTemplateData } from '../../../interfaces';

export interface Step5ConversionScriptProps {
  modelValue: AssetTemplateData;
  errorMessage?: string;
}

export interface Step5ConversionScriptEmits {
  (e: 'update:modelValue', value: AssetTemplateData): void;
}
