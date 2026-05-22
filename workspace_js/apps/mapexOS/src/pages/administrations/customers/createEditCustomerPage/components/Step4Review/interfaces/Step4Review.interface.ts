import type { OrganizationFormData, OrgTypeConfig, OrganizationType } from '../../../interfaces';

/** PROPS & EMITS */
export interface Step4ReviewProps {
  /** Organization form data to display in review */
  modelValue: OrganizationFormData;

  /** Whether in edit mode */
  isEditMode?: boolean;

  /** Type configuration for current org type */
  typeConfig: OrgTypeConfig;

  /** Organization type being created/edited */
  orgType: OrganizationType;
}

export interface Step4ReviewEmits {
  /** Emitted when user clicks edit button on a section */
  (e: 'editSection', step: number): void;
}
