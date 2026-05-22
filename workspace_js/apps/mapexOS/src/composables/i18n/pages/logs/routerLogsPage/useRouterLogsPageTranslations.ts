import { computed } from 'vue';
import { useTS } from '@utils/translation';

/**
 * Translation composable for Router Logs Page
 * Provides all translated strings for the router logs interface
 * @returns {Object} Translation object with page header, filters, columns, etc.
 */
export function useRouterLogsPageTranslations() {
  const ts = useTS({ capitalize: true });
  const tsTitle = useTS({ titleCase: true });
  const tsRaw = useTS({ capitalize: false });

  return {
    pageHeader: {
      title: computed(() => tsTitle('pages.logs.routerLogsPage.pageHeader.title')),
      description: computed(() => tsRaw('pages.logs.routerLogsPage.pageHeader.description')),
    },
    listTitle: computed(() => tsTitle('pages.logs.routerLogsPage.listTitle')),
    itemLabel: computed(() => ts('pages.logs.routerLogsPage.itemLabel')),
    itemLabelPlural: computed(() => ts('pages.logs.routerLogsPage.itemLabelPlural')),
    filters: {
      label: computed(() => ts('pages.logs.routerLogsPage.filters.label')),
      searchPlaceholder: computed(() => ts('pages.logs.routerLogsPage.filters.searchPlaceholder')),
      advancedFilters: computed(() => ts('pages.logs.routerLogsPage.filters.advancedFilters')),
      pendingFilters: computed(() => ts('pages.logs.routerLogsPage.filters.pendingFilters')),
      clearAll: computed(() => ts('pages.logs.routerLogsPage.filters.clearAll')),
      allStatus: computed(() => ts('pages.logs.routerLogsPage.filters.allStatus')),
      dateRange: computed(() => ts('pages.logs.routerLogsPage.filters.dateRange')),
      status: computed(() => ts('pages.logs.routerLogsPage.filters.status')),
      uuid: computed(() => ts('pages.logs.routerLogsPage.filters.uuid')),
      assetId: computed(() => ts('pages.logs.routerLogsPage.filters.assetId')),
      assetIdPlaceholder: computed(() => ts('pages.logs.routerLogsPage.filters.assetIdPlaceholder')),
      routerId: computed(() => ts('pages.logs.routerLogsPage.filters.routerId')),
      routerIdPlaceholder: computed(() => ts('pages.logs.routerLogsPage.filters.routerIdPlaceholder')),
      includeChildren: computed(() => ts('pages.logs.routerLogsPage.filters.includeChildren')),
      includeChildrenOrgs: computed(() => ts('pages.logs.routerLogsPage.filters.includeChildrenOrgs')),
      options: {
        yes: computed(() => ts('pages.logs.routerLogsPage.filters.options.yes')),
        no: computed(() => ts('pages.logs.routerLogsPage.filters.options.no')),
      },
    },
    columns: {
      routeGroup: computed(() => ts('pages.logs.routerLogsPage.columns.routeGroup')),
      uuid: computed(() => ts('pages.logs.routerLogsPage.columns.uuid')),
      status: computed(() => ts('pages.logs.routerLogsPage.columns.status')),
      totalRouters: computed(() => ts('pages.logs.routerLogsPage.columns.totalRouters')),
      matchedCount: computed(() => ts('pages.logs.routerLogsPage.columns.matchedCount')),
      publishedCount: computed(() => ts('pages.logs.routerLogsPage.columns.publishedCount')),
      timestamp: computed(() => ts('pages.logs.routerLogsPage.columns.timestamp')),
    },
    statusOptions: {
      success: computed(() => ts('pages.logs.routerLogsPage.statusOptions.success')),
      failed: computed(() => ts('pages.logs.routerLogsPage.statusOptions.failed')),
    },
    statusBadge: {
      success: computed(() => tsRaw('pages.logs.routerLogsPage.statusBadge.success')),
      failed: computed(() => tsRaw('pages.logs.routerLogsPage.statusBadge.failed')),
    },
    drawer: {
      title: computed(() => tsTitle('pages.logs.routerLogsPage.drawer.title')),
      subtitleSuccess: computed(() => ts('pages.logs.routerLogsPage.drawer.subtitleSuccess')),
      subtitleFailed: computed(() => ts('pages.logs.routerLogsPage.drawer.subtitleFailed')),
    },
    empty: {
      title: computed(() => tsTitle('pages.logs.routerLogsPage.empty.title')),
      description: computed(() => tsRaw('pages.logs.routerLogsPage.empty.description')),
    },
    pagination: {
      newer: computed(() => ts('pages.logs.routerLogsPage.pagination.newer')),
      older: computed(() => ts('pages.logs.routerLogsPage.pagination.older')),
    },
    messages: {
      loadFailed: computed(() => tsRaw('pages.logs.routerLogsPage.messages.loadFailed')),
    },
    defaults: {
      unknown: computed(() => tsRaw('pages.logs.routerLogsPage.defaults.unknown')),
      notAvailable: computed(() => tsRaw('pages.logs.routerLogsPage.defaults.notAvailable')),
    },
  };
}
