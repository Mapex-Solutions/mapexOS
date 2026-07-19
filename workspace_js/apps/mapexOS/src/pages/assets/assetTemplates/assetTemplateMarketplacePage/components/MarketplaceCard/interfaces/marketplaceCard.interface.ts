import type { AssetTemplateCatalogItem } from '@mapexos/schemas';

/**
 * Props for the MarketplaceCard component.
 */
export interface MarketplaceCardProps {
  /** The catalog item rendered by this card. */
  item: AssetTemplateCatalogItem;
  /**
   * Friendly category label resolved from the facets (e.g. "IoT Platforms").
   * Falls back to the raw category slug when the facet is not yet loaded.
   */
  categoryLabel?: string;
  /**
   * Whether this template is currently mid-request (install or uninstall). Drives
   * the footer button spinner; the same flag covers both actions since only one
   * runs at a time per card.
   */
  installing?: boolean;
  /**
   * Whether this template is already installed in the current organization. When
   * true the footer shows the uninstall action instead of install.
   */
  installed?: boolean;
}

/**
 * Events emitted by the MarketplaceCard component.
 */
export interface MarketplaceCardEmits {
  /**
   * Emitted when the card body is clicked to preview the template.
   * @param e - Event name.
   * @param item - The catalog item the card represents.
   */
  (e: 'view', item: AssetTemplateCatalogItem): void;
  /**
   * Emitted when the footer install button is clicked to install the template.
   * @param e - Event name.
   * @param item - The catalog item to install.
   */
  (e: 'install', item: AssetTemplateCatalogItem): void;
  /**
   * Emitted when the footer uninstall button is clicked to remove the installed
   * template. The parent runs the request so the card stays presentational.
   * @param e - Event name.
   * @param item - The catalog item to uninstall.
   */
  (e: 'uninstall', item: AssetTemplateCatalogItem): void;
}
