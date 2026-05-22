import type { AssetFormData } from '../../../interfaces';
import type { RouteGroupResponse } from '@mapexos/schemas';

export interface Step3RouteGroupsProps {
  modelValue: AssetFormData;
}

export interface Step3RouteGroupsEmits {
  (e: 'update:modelValue', value: Partial<AssetFormData>): void;
  (e: 'routeGroupsSelected', routeGroups: RouteGroupResponse[]): void;
}
