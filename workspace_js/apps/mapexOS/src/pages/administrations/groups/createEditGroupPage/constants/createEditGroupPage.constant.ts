/**
 * CreateEditGroupPage Constants
 */

import type { GroupFormData } from '../interfaces';

/**
 * Total number of steps in the form
 */
export const TOTAL_STEPS = 4;

/**
 * Step numbers enum for better readability
 */
export const STEP = {
  BASIC_INFO: 1,
  ROLES: 2,
  MEMBERS: 3,
  REVIEW: 4,
} as const;

/**
 * Initial form data values
 */
export const INITIAL_GROUP_FORM_DATA: GroupFormData = {
  name: '',
  description: '',
  enabled: true,
};

/**
 * Name validation rules
 */
export const NAME_MIN_LENGTH = 3;
export const NAME_MAX_LENGTH = 150;

/**
 * Description validation rules
 */
export const DESCRIPTION_MAX_LENGTH = 500;
