import type { OrganizationFormData } from '../../../interfaces';

/** PROPS & EMITS */
export interface Step1BasicProps {
  /** Organization form data */
  modelValue: OrganizationFormData;

  /** Whether this org type supports phone field */
  hasPhone?: boolean;
}
