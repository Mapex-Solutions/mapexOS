/**
 * CreateEditUserPage Constants
 *
 * Based on contract: workspace_go/packages/contracts/services/mapexos/onboarding/dtos.go
 *
 * V1: AuthProvider removed - always internal auth.
 * Next version: auth provider will be determined by the customer's Organization.AuthConfig
 */

import type { UserFormData, ScopeType, AccessType, GroupAccessMode } from '../interfaces';
import type { TourStepDefinition } from '@composables/tour';

/**
 * Total number of steps in the form
 * Personal, Security, Access, Review = 4 steps
 */
export const TOTAL_STEPS = 4;

/**
 * Step numbers enum for better readability
 */
export const STEP = {
  PERSONAL: 1,
  SECURITY: 2,
  ACCESS: 3,
  REVIEW: 4,
} as const;

/**
 * Initial form data values
 */
export const INITIAL_USER_FORM_DATA: UserFormData = {
  email: '',
  password: '',
  changePasswordNextLogin: false,
  firstName: '',
  lastName: '',
  phone: '',
  jobTitle: '',
  enabled: true,
  avatar: '',
  accessType: 'group', // Default to group (recommended)
  // Legacy single values (create mode)
  selectedGroup: undefined,
  directMembership: undefined,
  // Array values (edit mode with multiple)
  selectedGroups: [],
  directMemberships: [],
};

/**
 * Access type options for Step 3
 */
export const ACCESS_TYPE_OPTIONS: {
  label: string;
  value: AccessType;
  icon: string;
  description: string;
  recommended?: boolean;
  warning?: string;
}[] = [
  {
    label: 'Group access',
    value: 'group',
    icon: 'group',
    description: 'User joins a group and inherits its permissions',
    recommended: true,
  },
  {
    label: 'Direct assignment',
    value: 'direct',
    icon: 'person',
    description: 'Assign roles directly to this user',
    warning: 'Harder to manage at scale. Use only for exceptions.',
  },
];

/**
 * Group access mode options (when accessType = 'group')
 */
export const GROUP_ACCESS_MODE_OPTIONS: {
  label: string;
  value: GroupAccessMode;
  icon: string;
  description: string;
  recommended?: boolean;
}[] = [
  {
    label: 'Use existing group',
    value: 'existing',
    icon: 'group',
    description: 'Select an existing group. User inherits its roles.',
    recommended: true,
  },
  {
    label: 'Create new group',
    value: 'new',
    icon: 'group_add',
    description: 'Create a new group with specific roles.',
  },
];

/**
 * Scope options for direct membership
 */
export const SCOPE_OPTIONS: {
  label: string;
  value: ScopeType;
  icon: string;
  description: string;
}[] = [
  {
    label: 'Local',
    value: 'local',
    icon: 'location_on',
    description: 'Access only to the selected organization',
  },
  {
    label: 'Recursive',
    value: 'recursive',
    icon: 'account_tree',
    description: 'Access to the organization and all its children',
  },
];

/**
 * Demo user form data for tour mode
 * Pre-fills the form with realistic example data
 */
export const DEMO_USER_FORM_DATA: UserFormData = {
  email: 'john.doe@example.com',
  password: 'SecureP@ss123',
  changePasswordNextLogin: false,
  firstName: 'John',
  lastName: 'Doe',
  phone: '+1234567890',
  jobTitle: 'IoT Engineer',
  enabled: true,
  avatar: '',
  accessType: 'group',
  selectedGroup: {
    mode: 'existing',
    existingGroup: {
      groupId: 'demo-group-id',
      groupName: 'Engineering Team',
    },
  },
  selectedGroups: [
    {
      mode: 'existing',
      existingGroup: {
        groupId: 'demo-group-id',
        groupName: 'Engineering Team',
      },
    },
  ],
  directMemberships: [],
};

/**
 * Tour step definitions for the Create User page
 * Text comes from translations, these define targeting and positioning
 */
export const CREATE_USER_TOUR_STEPS: TourStepDefinition[] = [
  {
    element: '#stepper-sidebar',
    translationKey: 'overview',
    side: 'right',
    align: 'start',
  },
  {
    element: '#form-card',
    translationKey: 'step1',
    side: 'left',
    align: 'start',
  },
  {
    element: '#form-card',
    translationKey: 'step2',
    side: 'left',
    align: 'start',
  },
  {
    element: '#form-card',
    translationKey: 'step3',
    side: 'left',
    align: 'start',
  },
  {
    element: '#form-card',
    translationKey: 'step4',
    side: 'left',
    align: 'start',
  },
];

/**
 * Validation constraints from contract
 */
export const VALIDATION = {
  EMAIL_MAX_LENGTH: 254,
  PASSWORD_MIN_LENGTH: 8,
  PASSWORD_MAX_LENGTH: 72,
  FIRST_NAME_MIN_LENGTH: 2,
  FIRST_NAME_MAX_LENGTH: 100,
  LAST_NAME_MIN_LENGTH: 2,
  LAST_NAME_MAX_LENGTH: 100,
  JOB_TITLE_MAX_LENGTH: 120,
} as const;

/**
 * Email validation regex pattern
 */
export const EMAIL_VALIDATION_REGEX = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;

/**
 * Phone validation regex (E.164 format)
 */
export const PHONE_VALIDATION_REGEX = /^\+[1-9]\d{1,14}$/;
