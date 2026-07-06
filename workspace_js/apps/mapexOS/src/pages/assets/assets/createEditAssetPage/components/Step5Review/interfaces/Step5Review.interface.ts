import type { AssetFormData, AssetFormState } from '../../../interfaces';

export interface Step5ReviewProps {
  modelValue: AssetFormData;
  formState: AssetFormState;
  /** Visible step ids in display order; maps each review section to its live position. */
  visibleStepIds: string[];
}

export interface Step5ReviewEmits {
  (e: 'editSection', step: number): void;
}
