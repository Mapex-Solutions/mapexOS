import type { AssetTemplateData } from '../../../interfaces';

export interface Step3PreprocessorScriptProps {
  modelValue: AssetTemplateData;
}

export interface Step3PreprocessorScriptEmits {
  (e: 'update:modelValue', value: AssetTemplateData): void;
}
