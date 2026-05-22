import type { UserFormData } from '../../../interfaces';

/** PROPS & EMITS */
export interface Step1PersonalProps {
  /** User form data */
  modelValue: UserFormData;

  /** Whether in edit mode */
  isEditMode?: boolean;
}
