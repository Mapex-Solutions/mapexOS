/**
 * Role info for display in the TabRoles component
 */
export interface RoleInfo {
  id: string;
  name: string;
  description?: string;
  isSystem?: boolean;
  isTemplate?: boolean;
  scope?: string;
}
