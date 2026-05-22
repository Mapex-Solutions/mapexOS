/** VUE IMPORTS */
import { computed } from 'vue';

/** UTILS */
import { useTS } from '@utils/translation';

const bp = 'pages.dashboards.dashboardAdm';

/**
 * Translations composable for the admin Dashboard page
 *
 * Structure mirrors:
 * - File: src/pages/dashboards/dashboardAdm/DashboardAdm.vue
 * - JSON: src/i18n/{locale}/pages/dashboards/dashboardAdm.json
 * - Composable: src/composables/i18n/pages/dashboards/dashboardAdm/useDashboardAdmTranslations.ts
 *
 * @returns {object} Reactive translation objects for header, sections, KPIs and quick actions
 */
export function useDashboardAdmTranslations() {
  const ts = useTS({ capitalize: true });
  const tsTitle = useTS({ titleCase: true });
  const tsRaw = useTS({ capitalize: false });

  return {
    /**
     * Page header translations
     */
    header: {
      title: computed(() => tsTitle(`${bp}.header.title`)),
      subtitlePrefix: computed(() => ts(`${bp}.header.subtitlePrefix`)),
      notAvailable: computed(() => tsRaw(`${bp}.header.notAvailable`)),
    },

    /**
     * Section titles
     */
    sections: {
      quickActions: computed(() => tsTitle(`${bp}.sections.quickActions`)),
      platformInfo: computed(() => tsTitle(`${bp}.sections.platformInfo`)),
    },

    /**
     * Platform info card labels
     */
    platformInfo: {
      organization: computed(() => tsTitle(`${bp}.platformInfo.organization`)),
      type: computed(() => tsTitle(`${bp}.platformInfo.type`)),
      totalOrganizations: computed(() => tsTitle(`${bp}.platformInfo.totalOrganizations`)),
      notAvailable: computed(() => tsRaw(`${bp}.platformInfo.notAvailable`)),
    },

    /**
     * Primary KPI labels
     */
    primaryKpis: {
      assets: computed(() => tsTitle(`${bp}.primaryKpis.assets`)),
      users: computed(() => tsTitle(`${bp}.primaryKpis.users`)),
      triggers: computed(() => tsTitle(`${bp}.primaryKpis.triggers`)),
      wfDefinitions: computed(() => tsTitle(`${bp}.primaryKpis.wfDefinitions`)),
    },

    /**
     * Secondary KPI labels
     */
    secondaryKpis: {
      assetTemplates: computed(() => tsTitle(`${bp}.secondaryKpis.assetTemplates`)),
      groups: computed(() => tsTitle(`${bp}.secondaryKpis.groups`)),
      routeGroups: computed(() => tsTitle(`${bp}.secondaryKpis.routeGroups`)),
      wfInstances: computed(() => tsTitle(`${bp}.secondaryKpis.wfInstances`)),
    },

    /**
     * Quick actions labels
     */
    quickActions: {
      newAsset: computed(() => tsTitle(`${bp}.quickActions.newAsset`)),
      newUser: computed(() => tsTitle(`${bp}.quickActions.newUser`)),
      newTemplate: computed(() => tsTitle(`${bp}.quickActions.newTemplate`)),
      newRouteGroup: computed(() => tsTitle(`${bp}.quickActions.newRouteGroup`)),
      viewLogs: computed(() => tsTitle(`${bp}.quickActions.viewLogs`)),
    },
  };
}
