import type { RouteGroupCreate } from '@interfaces/routing/routeGroups.interface';
import type { RouterFormState } from '../../../interfaces';

/**
 * Props for Step3Review component
 */
export interface Step3ReviewProps {
  /** Form data for route group */
  formData: RouteGroupCreate;

  /** Array of router form states */
  routerForms: RouterFormState[];
}
