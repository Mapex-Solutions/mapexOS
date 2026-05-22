/**
 * Organization Icon Utilities
 *
 * Provides consistent icon and color mapping for organization types.
 * Used by OrganizationTreeDrawer and other organization-related components.
 */

import type { OrganizationType } from '@stores/organization/types';

/**
 * Maps organization type to Material Design icon name
 *
 * Icon choices:
 * - vendor: store (commercial vendor)
 * - customer: business (business entity)
 * - site: place (physical location)
 * - building: apartment (multi-floor building)
 * - floor: layers (stacked levels)
 * - zone: room (area within floor)
 */
export function getOrganizationIcon(type: OrganizationType): string {
  const icons: Record<OrganizationType, string> = {
    vendor: 'store',
    customer: 'business',
    site: 'place',
    building: 'apartment',
    floor: 'layers',
    zone: 'room',
  };

  return icons[type] || 'help_outline';
}

/**
 * Maps organization type to Quasar color name
 *
 * Color scheme:
 * - vendor: purple (highest level)
 * - customer: primary (main business level)
 * - site: orange (physical locations)
 * - building: blue (structural level)
 * - floor: teal (intermediate level)
 * - zone: green (lowest level)
 */
export function getOrganizationColor(type: OrganizationType): string {
  const colors: Record<OrganizationType, string> = {
    vendor: 'purple',
    customer: 'primary',
    site: 'orange',
    building: 'blue',
    floor: 'teal',
    zone: 'green',
  };

  return colors[type] || 'grey';
}
