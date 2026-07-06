import type { AssetTemplateCatalogItem } from '@mapexos/schemas';

/**
 * Props for the MarketplaceCard component.
 */
export interface MarketplaceCardProps {
  /** The catalog item rendered by this card. */
  item: AssetTemplateCatalogItem;
}

/**
 * Events emitted by the MarketplaceCard component.
 */
export interface MarketplaceCardEmits {
  /**
   * Emitted when the card is clicked to preview the template.
   * @param e - Event name.
   * @param item - The catalog item the card represents.
   */
  (e: 'view', item: AssetTemplateCatalogItem): void;
}
