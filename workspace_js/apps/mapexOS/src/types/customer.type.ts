export interface CustomerData {
  firstName: string;
  lastName: string;
  email: string;
  phone: string;
  company: string;
  industry: SelectOption | null;
  type: CustomerTypeOption | null;
  assignedTo: SelectOption | null;
  status: SelectOption | null;
}

export interface SelectOption {
  label: string;
  value: string | number;
  color?: string;
}

export interface CustomerTypeOption extends SelectOption {
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
