/**
 * Props for the marketplace detail modal. `vendor`/`slug` identify the template
 * whose bundle is fetched and previewed; they are `null` while no card is
 * selected.
 */
export interface AssetTemplateMarketplaceDetailModalProps {
  /** Controls dialog visibility (v-model). */
  modelValue: boolean;
  /** Vendor key of the template to preview, or null when none is selected. */
  vendor: string | null;
  /** Slug of the template to preview, or null when none is selected. */
  slug: string | null;
}

/**
 * The vendor/slug pair carried by the install events.
 */
export interface AssetTemplateMarketplaceInstallPayload {
  /** Template vendor key. */
  vendor: string;
  /** Template slug. */
  slug: string;
}

/**
 * Emits for the marketplace detail modal.
 */
export interface AssetTemplateMarketplaceDetailModalEmits {
  /** Sync dialog visibility. */
  (e: 'update:modelValue', value: boolean): void;
  /** Fired when the user triggers an install for the previewed template. */
  (e: 'install', payload: AssetTemplateMarketplaceInstallPayload): void;
  /** Fired after a successful install of the previewed template. */
  (e: 'installed', payload: AssetTemplateMarketplaceInstallPayload): void;
}
