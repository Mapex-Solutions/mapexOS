import type { RouteGroupCreate } from '@interfaces/routing/routeGroups.interface';

/**
 * Props for Step1BasicInfo component
 */
export interface Step1BasicInfoProps {
  /** Form data for route group */
  formData: RouteGroupCreate;

  /** Status options for the enabled field */
  statusOptions: Array<{ label: string; value: boolean }>;

  /** Whether the user can create templates */
  canCreateTemplate: boolean;

  /** Translation composable */
  t: any;
}

/**
 * Emits for Step1BasicInfo component
 */
export interface Step1BasicInfoEmits {
  (e: 'update:formData', value: RouteGroupCreate): void;
}
