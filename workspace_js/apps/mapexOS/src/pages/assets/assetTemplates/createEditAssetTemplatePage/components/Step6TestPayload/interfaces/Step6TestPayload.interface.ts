import type { AssetTemplateData } from '../../../interfaces';

export interface Step6TestPayloadProps {
  modelValue: AssetTemplateData;
}

export interface Step6TestPayloadEmits {
  (e: 'update:modelValue', value: AssetTemplateData): void;
}
