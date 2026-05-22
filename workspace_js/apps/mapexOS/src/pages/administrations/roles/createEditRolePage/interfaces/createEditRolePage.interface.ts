/**
 * CreateEditRolePage Interfaces
 */

/**
 * Permission action item
 */
export interface PermissionAction {
  /** Action name (e.g., 'list', 'create', 'read', 'update', 'delete') */
  name: string;

  /** Display label */
  label: string;

  /** Whether action is granted */
  granted: boolean;

  /** Override permission string (default: `${resource}.${name}`) */
  permissionKey?: string;
}

/**
 * Resource permission group
 */
export interface ResourcePermission {
  /** Resource key (e.g., 'users', 'assets') */
  resource: string;

  /** Display name */
  label: string;

  /** Material icon name */
  icon: string;

  /** Whether resource is enabled (has any permission) */
  enabled: boolean;

  /** Available actions for this resource */
  actions: PermissionAction[];
}

/**
 * Permission group that mirrors sidebar menu sections
 */
export interface PermissionGroup {
  /** Group label matching sidebar section name */
  label: string;

  /** Material icon name matching sidebar section icon */
  icon: string;

  /** Resource keys that belong to this group */
  resources: string[];
}

/**
 * Role form data structure
 */
export interface RoleFormData {
  /** Role name (3-100 chars) */
  name: string;

  /** Role description (max 500 chars) */
  description: string;

  /** Scope: 'global' (inherits to children) or 'local' (this org only) */
  scope: 'global' | 'local' | null;

  /** Whether role is a template (can be shared with child orgs) */
  isTemplate: boolean;
}

/**
 * Role form state
 */
export interface RoleFormState {
  /** Resource permissions configuration */
  resourcePermissions: ResourcePermission[];

  /** Whether form is being saved */
  isSaving: boolean;

  /** Current step number */
  currentStep: number;
}

/**
 * Step validation errors state
 */
export interface StepValidationErrors {
  /** Step 1 has errors */
  step1: boolean;

  /** Step 2 has errors */
  step2: boolean;

  /** Step 3 has errors */
  step3: boolean;
}

/**
 * Props for Step1BasicInfo component
 */
export interface Step1BasicInfoProps {
  /** Current form data */
  modelValue: RoleFormData;
}

/**
 * Emits for Step1BasicInfo component
 */
export interface Step1BasicInfoEmits {
  (e: 'update:modelValue', value: Partial<RoleFormData>): void;
}

/**
 * Props for Step2Permissions component
 */
export interface Step2PermissionsProps {
  /** Current permissions state */
  modelValue: ResourcePermission[];
}

/**
 * Emits for Step2Permissions component
 */
export interface Step2PermissionsEmits {
  (e: 'update:modelValue', value: ResourcePermission[]): void;
}

/**
 * Props for Step3Review component
 */
export interface Step3ReviewProps {
  /** Role form data */
  roleData: RoleFormData;

  /** Resource permissions */
  resourcePermissions: ResourcePermission[];
}

/**
 * Emits for Step3Review component
 */
export interface Step3ReviewEmits {
  (e: 'editSection', step: number): void;
}
