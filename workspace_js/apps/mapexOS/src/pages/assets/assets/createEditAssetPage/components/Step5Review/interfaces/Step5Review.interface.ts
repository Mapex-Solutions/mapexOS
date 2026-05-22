import type { AssetFormData, AssetFormState } from '../../../interfaces';

export interface Step5ReviewProps {
  modelValue: AssetFormData;
  formState: AssetFormState;
}

export interface Step5ReviewEmits {
  (e: 'editSection', step: number): void;
}
