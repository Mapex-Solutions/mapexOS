/**
 * SystemSettingsPage Interfaces
 *
 * This file defines TypeScript interfaces for the System Settings Page component.
 * Follows MapexOS architecture pattern: interfaces in separate folder.
 *
 * @see src/pages/administrations/settings/settingsPage/SystemSettingsPage.vue
 */

/**
 * Payload sent from GeneralSettings to parent when saving
 * Contains only editable fields that should be sent to API
 */
export interface GeneralSettingsSavePayload {
  name: string;
  phone: string;
  enabled: boolean;
  address: {
    city?: string;
    state?: string;
    country?: string;
    zipCode?: string;
  };
  accessPolicy: {
    rolePolicy: 'merge' | 'strict';
    defaultScope: 'recursive' | 'local';
  };
}
