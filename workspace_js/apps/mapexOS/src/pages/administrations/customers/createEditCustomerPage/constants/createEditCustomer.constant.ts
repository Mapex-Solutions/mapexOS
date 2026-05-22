/**
 * CreateEditOrganizationPage Constants
 * Maps to workspace_go/packages/contracts/services/mapexos/organizations/dtos.go
 */

import type { OrganizationFormData, OrganizationType, OrgTypeConfig } from '../interfaces';

/**
 * Maps parent org type to the valid child type
 * Enforces: vendor → customer → site → building → floor → zone
 * Zone is leaf (not present as key - cannot have children)
 */
export const CHILD_TYPE_MAP: Record<string, OrganizationType> = {
  vendor: 'customer',
  customer: 'site',
  site: 'building',
  building: 'floor',
  floor: 'zone',
} as const;

/**
 * Configuration per organization type for UI rendering
 */
export const ORG_TYPE_CONFIG: Record<OrganizationType, OrgTypeConfig> = {
  vendor: {
    label: 'Vendor',
    icon: 'domain',
    iconColor: 'deep-purple',
    hasAddress: true,
    hasPhone: true,
  },
  customer: {
    label: 'Customer',
    icon: 'business',
    iconColor: 'primary',
    hasAddress: true,
    hasPhone: true,
  },
  site: {
    label: 'Site',
    icon: 'location_on',
    iconColor: 'teal',
    hasAddress: true,
    hasPhone: true,
  },
  building: {
    label: 'Building',
    icon: 'apartment',
    iconColor: 'orange',
    hasAddress: false,
    hasPhone: false,
  },
  floor: {
    label: 'Floor',
    icon: 'layers',
    iconColor: 'cyan',
    hasAddress: false,
    hasPhone: false,
  },
  zone: {
    label: 'Zone',
    icon: 'grid_view',
    iconColor: 'pink',
    hasAddress: false,
    hasPhone: false,
  },
} as const;

/**
 * Total number of steps in the form (dynamic based on type)
 * Types with address: Basic → Address → AccessPolicy → Review (4 steps)
 * Types without address: Basic → AccessPolicy → Review (3 steps)
 */
export const TOTAL_STEPS_WITH_ADDRESS = 4;
export const TOTAL_STEPS_WITHOUT_ADDRESS = 3;

/**
 * Step numbers for types WITH address (customer, site)
 */
export const STEP_WITH_ADDRESS = {
  BASIC: 1,
  ADDRESS: 2,
  ACCESS_POLICY: 3,
  REVIEW: 4,
} as const;

/**
 * Step numbers for types WITHOUT address (building, floor, zone)
 */
export const STEP_WITHOUT_ADDRESS = {
  BASIC: 1,
  ACCESS_POLICY: 2,
  REVIEW: 3,
} as const;

/** @deprecated Use STEP_WITH_ADDRESS instead */
export const STEP = STEP_WITH_ADDRESS;

/** @deprecated Use TOTAL_STEPS_WITH_ADDRESS instead */
export const TOTAL_STEPS = TOTAL_STEPS_WITH_ADDRESS;

/**
 * Validation constants from Go contract
 */
export const NAME_MIN_LENGTH = 3;
export const NAME_MAX_LENGTH = 150;
export const CITY_MAX_LENGTH = 100;
export const STATE_MAX_LENGTH = 100;
export const COUNTRY_MAX_LENGTH = 100;
export const ZIPCODE_MAX_LENGTH = 20;

/**
 * Initial form data values
 */
export const INITIAL_ORGANIZATION_FORM_DATA: OrganizationFormData = {
  name: '',
  phone: '',
  enabled: true,
  address: {
    city: '',
    state: '',
    country: '',
    zipCode: '',
  },
  authConfig: {
    providerType: 'internal',
    issuerUrl: '',
    clientId: '',
  },
  accessPolicy: {
    rolePolicy: 'strict',
    defaultScope: 'local',
  },
};

/** @deprecated Use INITIAL_ORGANIZATION_FORM_DATA instead */
export const INITIAL_CUSTOMER_FORM_DATA = INITIAL_ORGANIZATION_FORM_DATA;

/**
 * Auth provider options
 * V1: Only internal is enabled. Next version: unlock based on Organization.AuthConfig
 */
export const AUTH_PROVIDER_OPTIONS = [
  {
    value: 'internal',
    label: 'Internal',
    icon: 'lock',
    description: 'Use internal authentication managed by MapexOS',
    disabled: false,
  },
  {
    value: 'keycloak',
    label: 'Keycloak SSO',
    icon: 'vpn_key',
    description: 'Integrate with Keycloak for single sign-on',
    disabled: true,
  },
] as const;

/**
 * Role policy options
 */
export const ROLE_POLICY_OPTIONS = [
  {
    value: 'strict',
    label: 'Strict',
    icon: 'security',
    description: 'Only explicitly assigned roles are applied. More restrictive.',
  },
  {
    value: 'merge',
    label: 'Merge',
    icon: 'merge',
    description: 'Roles from parent organizations are merged with local roles.',
  },
] as const;

/**
 * Default scope options
 */
export const DEFAULT_SCOPE_OPTIONS = [
  {
    value: 'local',
    label: 'Local',
    icon: 'location_on',
    description: 'Roles only apply to this organization.',
  },
  {
    value: 'recursive',
    label: 'Recursive',
    icon: 'account_tree',
    description: 'Roles apply to this organization and all children.',
  },
] as const;
