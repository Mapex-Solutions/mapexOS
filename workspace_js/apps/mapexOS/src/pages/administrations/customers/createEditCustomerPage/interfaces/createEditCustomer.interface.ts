/**
 * CreateEditOrganizationPage Interfaces
 * Maps to workspace_go/packages/contracts/services/mapexos/organizations/dtos.go
 */

/**
 * Organization type in the hierarchy
 */
export type OrganizationType = 'vendor' | 'customer' | 'site' | 'building' | 'floor' | 'zone';

/**
 * Auth provider type for organization authentication
 */
export type AuthProviderType = 'keycloak' | 'internal';

/**
 * Role policy type
 */
export type RolePolicyType = 'strict' | 'merge';

/**
 * Default scope type
 */
export type DefaultScopeType = 'local' | 'recursive';

/**
 * Address structure for organization
 */
export interface Address {
  /** City name */
  city: string;

  /** State/Province name */
  state: string;

  /** Country name */
  country: string;

  /** Postal/ZIP code */
  zipCode: string;
}

/**
 * Auth configuration for organization
 */
export interface AuthConfig {
  /** Authentication provider type */
  providerType: AuthProviderType;

  /** SSO issuer URL (for keycloak) */
  issuerUrl: string;

  /** SSO client ID (for keycloak) */
  clientId: string;

  /** JWT claim mappings (optional, advanced) */
  jwtClaimMappings?: Record<string, string>;

  /** Additional metadata (optional, advanced) */
  metadata?: Record<string, unknown>;
}

/**
 * Access policy configuration for organization
 */
export interface AccessPolicy {
  /** How roles are applied (strict or merge) */
  rolePolicy: RolePolicyType;

  /** Default scope for role assignments */
  defaultScope: DefaultScopeType;
}

/**
 * Organization form data structure
 * Maps to OrganizationCreate DTO
 */
export interface OrganizationFormData {
  /** Organization name (required, min 3, max 150) */
  name: string;

  /** Phone number (optional, e164 format) - only for customer and site */
  phone: string;

  /** Whether the organization is enabled */
  enabled: boolean;

  /** Address information - only for customer and site */
  address: Address;

  /** Authentication configuration - V1: always internal, locked in UI */
  authConfig: AuthConfig;

  /** Access policy configuration */
  accessPolicy: AccessPolicy;
}

/**
 * Organization type configuration for UI rendering
 */
export interface OrgTypeConfig {
  /** Display label for the type */
  label: string;

  /** Icon for the type */
  icon: string;

  /** Icon color */
  iconColor: string;

  /** Whether this type has address fields */
  hasAddress: boolean;

  /** Whether this type has phone field */
  hasPhone: boolean;
}

/** @deprecated Use OrganizationFormData instead */
export type CustomerFormData = OrganizationFormData;

/**
 * Organization form state
 */
export interface OrganizationFormState {
  /** Whether form is being saved */
  isSaving: boolean;

  /** Current step number */
  currentStep: number;
}

/** @deprecated Use OrganizationFormState instead */
export type CustomerFormState = OrganizationFormState;
