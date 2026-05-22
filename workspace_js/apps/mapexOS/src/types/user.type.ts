export interface UserData {
  firstName: string;
  lastName: string;
  email: string;
  phone: string;
  jobTitle: string;
  company: SelectOption | null;
  department: SelectOption | null;
  role: RoleOption | null;
  reportingTo: SelectOption | null;
  accessLevel: SelectOption | null;
}

export interface SelectOption {
  label: string;
  value: string | number;
}

export interface RoleOption extends SelectOption {
  color: string;
}

export interface Permission {
  name: string;
  granted: boolean;
}

export interface ModulePermission {
  name: string;
  icon: string;
  enabled: boolean;
  permissions: Permission[];
}