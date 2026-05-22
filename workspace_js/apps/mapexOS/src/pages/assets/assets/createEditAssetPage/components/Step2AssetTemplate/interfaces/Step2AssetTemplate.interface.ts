import type { AssetFormData } from '../../../interfaces';
import type { AssetTemplateResponse } from '@mapexos/schemas';

/** INTERFACES */
export interface Step2AssetTemplateProps {
  modelValue: AssetFormData;
}

export interface Step2AssetTemplateEmits {
  (e: 'update:modelValue', value: Partial<AssetFormData>): void;
  (e: 'templateSelected', template: AssetTemplateResponse | null): void;
}
