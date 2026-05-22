import type { AssetTemplateData } from '../../../interfaces';

export interface Step1BasicInfoProps {
  modelValue: AssetTemplateData;
  canCreateTemplate: boolean;
}

export interface Step1BasicInfoEmits {
  (e: 'update:modelValue', value: AssetTemplateData): void;
}
