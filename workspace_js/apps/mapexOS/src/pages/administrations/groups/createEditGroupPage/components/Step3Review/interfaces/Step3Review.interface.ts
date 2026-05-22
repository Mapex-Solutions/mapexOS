import type { GroupFormData, RoleSelectionItem } from '../../../interfaces';

/** INTERFACES */
export interface Step3ReviewProps {
  modelValue: GroupFormData;
  selectedRoles: RoleSelectionItem[];
  selectedMembers: string[];
  /** Initial members count from API (used in edit mode before visiting Step 2) */
  initialMembersCount?: number;
  isEditMode?: boolean;
}

export interface Step3ReviewEmits {
  (e: 'editSection', step: number): void;
}
