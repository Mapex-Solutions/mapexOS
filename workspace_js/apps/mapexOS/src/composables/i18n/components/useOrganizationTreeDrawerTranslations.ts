import { computed } from 'vue';
import { useTS } from '@utils/translation';

/**
 * OrganizationTreeDrawer component translations
 *
 * Structure mirrors:
 * - File: src/components/drawers/organizationTree/OrganizationTreeDrawer.vue
 * - JSON: src/i18n/{locale}/components/organizationTreeDrawer.json
 * - Composable: src/composables/i18n/components/useOrganizationTreeDrawerTranslations.ts
 *
 * Provides all translations for the organization tree drawer including:
 * - Drawer title
 * - Filter labels and options
 * - Organization type names
 * - Empty and loading states
 *
 * @example
 * ```ts
 * // In OrganizationTreeDrawer.vue
 * const t = useOrganizationTreeDrawerTranslations();
 *
 * <q-toolbar-title>{{ t.title.value }}</q-toolbar-title>
 * <q-input :label="t.filters.searchLabel.value" />
 * ```
 */
export function useOrganizationTreeDrawerTranslations() {
  const ts = useTS({ capitalize: true });

  return {
    /**
     * Drawer title
     * Mirrors: components.organizationTreeDrawer.title
     */
    title: computed(() => ts('components.organizationTreeDrawer.title')),

    /**
     * Close button tooltip
     * Mirrors: components.organizationTreeDrawer.close
     */
    close: computed(() => ts('components.organizationTreeDrawer.close')),

    /**
     * Legend tooltip
     * Mirrors: components.organizationTreeDrawer.legendTooltip
     */
    legendTooltip: computed(() => ts('components.organizationTreeDrawer.legendTooltip')),

    /**
     * Filter translations
     * Mirrors: components.organizationTreeDrawer.filters
     */
    filters: {
      searchLabel: computed(() => ts('components.organizationTreeDrawer.filters.searchLabel')),
      searchPlaceholder: computed(() => ts('components.organizationTreeDrawer.filters.searchPlaceholder')),
      typeLabel: computed(() => ts('components.organizationTreeDrawer.filters.typeLabel')),
      enabledAll: computed(() => ts('components.organizationTreeDrawer.filters.enabledAll')),
      enabledActive: computed(() => ts('components.organizationTreeDrawer.filters.enabledActive')),
      enabledInactive: computed(() => ts('components.organizationTreeDrawer.filters.enabledInactive')),
    },

    /**
     * Organization type names
     * Mirrors: components.organizationTreeDrawer.types
     */
    types: {
      vendor: computed(() => ts('components.organizationTreeDrawer.types.vendor')),
      customer: computed(() => ts('components.organizationTreeDrawer.types.customer')),
      site: computed(() => ts('components.organizationTreeDrawer.types.site')),
      building: computed(() => ts('components.organizationTreeDrawer.types.building')),
      floor: computed(() => ts('components.organizationTreeDrawer.types.floor')),
      zone: computed(() => ts('components.organizationTreeDrawer.types.zone')),
    },

    /**
     * Legend modal translations
     * Mirrors: components.organizationTreeDrawer.legend
     */
    legend: {
      title: computed(() => ts('components.organizationTreeDrawer.legend.title')),
      description: computed(() => ts('components.organizationTreeDrawer.legend.description')),
      typesTitle: computed(() => ts('components.organizationTreeDrawer.legend.typesTitle')),
      vendorDesc: computed(() => ts('components.organizationTreeDrawer.legend.vendorDesc')),
      customerDesc: computed(() => ts('components.organizationTreeDrawer.legend.customerDesc')),
      siteDesc: computed(() => ts('components.organizationTreeDrawer.legend.siteDesc')),
      buildingDesc: computed(() => ts('components.organizationTreeDrawer.legend.buildingDesc')),
      floorDesc: computed(() => ts('components.organizationTreeDrawer.legend.floorDesc')),
      zoneDesc: computed(() => ts('components.organizationTreeDrawer.legend.zoneDesc')),
      actionsTitle: computed(() => ts('components.organizationTreeDrawer.legend.actionsTitle')),
      singleClick: computed(() => ts('components.organizationTreeDrawer.legend.singleClick')),
      singleClickDesc: computed(() => ts('components.organizationTreeDrawer.legend.singleClickDesc')),
      doubleClick: computed(() => ts('components.organizationTreeDrawer.legend.doubleClick')),
      doubleClickDesc: computed(() => ts('components.organizationTreeDrawer.legend.doubleClickDesc')),
      statusTitle: computed(() => ts('components.organizationTreeDrawer.legend.statusTitle')),
      inactiveDesc: computed(() => ts('components.organizationTreeDrawer.legend.inactiveDesc')),
      close: computed(() => ts('components.organizationTreeDrawer.legend.close')),
    },

    /**
     * State messages
     * Mirrors: components.organizationTreeDrawer
     */
    empty: computed(() => ts('components.organizationTreeDrawer.empty')),
    loading: computed(() => ts('components.organizationTreeDrawer.loading')),
    inactive: computed(() => ts('components.organizationTreeDrawer.inactive')),
  };
}
