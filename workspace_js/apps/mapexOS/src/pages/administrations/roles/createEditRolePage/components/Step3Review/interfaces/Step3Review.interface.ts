import type { RoleFormData, ResourcePermission } from '../../../interfaces';

/** INTERFACES */
export interface Step3ReviewProps {
  modelValue: RoleFormData;
  resourcePermissions: ResourcePermission[];
  isEditMode?: boolean;
}

export interface Step3ReviewEmits {
  (e: 'editSection', step: number): void;
}
