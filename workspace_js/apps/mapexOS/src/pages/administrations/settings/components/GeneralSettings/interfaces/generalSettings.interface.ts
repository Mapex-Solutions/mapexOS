import type { OrganizationResponse } from '@mapexos/schemas';
import type { GeneralSettingsSavePayload } from '../../../settingsPage/interfaces';

/**
 * GeneralSettings Component Interfaces
 *
 * This file defines TypeScript interfaces for the GeneralSettings component.
 * Follows MapexOS architecture pattern: Props Down, Events Up.
 *
 * @see src/pages/administrations/settings/components/GeneralSettings.vue
 */

/**
 * Props received from parent (SystemSettingsPage)
 */
export interface GeneralSettingsProps {
  /**
   * Complete organization data fetched from API
   * null when loading or no data
   */
  organizationData: OrganizationResponse | null;

  /**
   * Loading state managed by parent
   * Shows spinner while fetching or saving
   */
  loading?: boolean;
}

/**
 * Local editable form state.
 * Mirrors the editable subset of OrganizationResponse plus readonly display fields.
 */
export interface LocalGeneralData {
  name: string;
  phone: string;
  enabled: boolean;
  city: string;
  state: string;
  country: string;
  zipCode: string;
  rolePolicy: 'merge' | 'strict';
  defaultScope: 'recursive' | 'local';
  type: string;
  code: string;
  pathKey: string;
  depth: number;
  childCount: number;
  parentOrgId: string | null;
  authProviderType: string;
  created: string;
  updated: string;
}

/**
 * Events emitted to parent (SystemSettingsPage)
 */
export interface GeneralSettingsEmits {
  /**
   * Emitted when user clicks save button
   * Payload contains only editable fields
   *
   * @param payload - Sanitized data ready for API
   */
  (e: 'save', payload: GeneralSettingsSavePayload): void;

  /**
   * Emitted when dirty state changes
   * Used to track unsaved changes
   *
   * @param isDirty - Whether form has unsaved changes
   */
  (e: 'update:dirty', isDirty: boolean): void;
}
