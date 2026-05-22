import type { FilterField } from '@components/drawers';

import { computed } from 'vue';
import { useTS } from '@utils/translation';

/**
 * Translation composable for Asset Connectivity Logs Page.
 * Provides all translated strings for the asset connectivity transitions
 * (offline/online) logs interface.
 *
 * @returns Translation object with page header, filters, columns, etc.
 */
export function useAssetConnectivityLogsPageTranslations() {
  const ts = useTS({ capitalize: true });
  const tsTitle = useTS({ titleCase: true });
  const tsRaw = useTS({ capitalize: false });

  return {
    pageHeader: {
      title: computed(() => tsTitle('pages.logs.assetConnectivityLogsPage.pageHeader.title')),
      description: computed(() => tsRaw('pages.logs.assetConnectivityLogsPage.pageHeader.description')),
    },
    listTitle: computed(() => tsTitle('pages.logs.assetConnectivityLogsPage.listTitle')),
    itemLabel: computed(() => ts('pages.logs.assetConnectivityLogsPage.itemLabel')),
    itemLabelPlural: computed(() => ts('pages.logs.assetConnectivityLogsPage.itemLabelPlural')),
    filters: {
      label: computed(() => ts('pages.logs.assetConnectivityLogsPage.filters.label')),
      searchPlaceholder: computed(() => ts('pages.logs.assetConnectivityLogsPage.filters.searchPlaceholder')),
      advancedFilters: computed(() => ts('pages.logs.assetConnectivityLogsPage.filters.advancedFilters')),
      pendingFilters: computed(() => ts('pages.logs.assetConnectivityLogsPage.filters.pendingFilters')),
      clearAll: computed(() => ts('pages.logs.assetConnectivityLogsPage.filters.clearAll')),
      allStatus: computed(() => ts('pages.logs.assetConnectivityLogsPage.filters.allStatus')),
      dateRange: computed(() => ts('pages.logs.assetConnectivityLogsPage.filters.dateRange')),
      eventType: computed(() => ts('pages.logs.assetConnectivityLogsPage.filters.eventType')),
      assetUUID: computed(() => ts('pages.logs.assetConnectivityLogsPage.filters.assetUUID')),
      includeChildren: computed(() => ts('pages.logs.assetConnectivityLogsPage.filters.includeChildren')),
      includeChildrenOrgs: computed(() => ts('pages.logs.assetConnectivityLogsPage.filters.includeChildrenOrgs')),
      options: {
        yes: computed(() => ts('pages.logs.assetConnectivityLogsPage.filters.options.yes')),
        no: computed(() => ts('pages.logs.assetConnectivityLogsPage.filters.options.no')),
        offline: computed(() => ts('pages.logs.assetConnectivityLogsPage.filters.options.offline')),
        online: computed(() => ts('pages.logs.assetConnectivityLogsPage.filters.options.online')),
      },
    },

    advancedFilters: computed((): FilterField[] => [
      {
        key: 'includeChildren',
        type: 'toggle',
        label: ts('pages.logs.assetConnectivityLogsPage.filters.includeChildrenOrgs'),
        icon: 'account_tree',
        options: [
          { label: ts('pages.logs.assetConnectivityLogsPage.filters.allStatus'), value: null },
          { label: ts('pages.logs.assetConnectivityLogsPage.filters.options.yes'), value: true },
          { label: ts('pages.logs.assetConnectivityLogsPage.filters.options.no'), value: false },
        ],
      },
      {
        key: 'startDate',
        type: 'input',
        label: ts('pages.logs.assetConnectivityLogsPage.filters.startDate'),
        icon: 'event',
        inputType: 'date',
      },
      {
        key: 'endDate',
        type: 'input',
        label: ts('pages.logs.assetConnectivityLogsPage.filters.endDate'),
        icon: 'event',
        inputType: 'date',
      },
      {
        key: 'eventType',
        type: 'select',
        label: ts('pages.logs.assetConnectivityLogsPage.filters.eventType'),
        icon: 'wifi',
        options: [
          { label: ts('pages.logs.assetConnectivityLogsPage.filters.allStatus'), value: null },
          { label: ts('pages.logs.assetConnectivityLogsPage.filters.options.offline'), value: 'offline' },
          { label: ts('pages.logs.assetConnectivityLogsPage.filters.options.online'), value: 'online' },
        ],
      },
      {
        key: 'assetUUID',
        type: 'input',
        label: ts('pages.logs.assetConnectivityLogsPage.filters.assetUUID'),
        icon: 'link',
        placeholder: ts('pages.logs.assetConnectivityLogsPage.filters.assetUUIDPlaceholder'),
      },
    ]),

    columns: {
      asset: computed(() => ts('pages.logs.assetConnectivityLogsPage.columns.asset')),
      assetUUID: computed(() => ts('pages.logs.assetConnectivityLogsPage.columns.assetUUID')),
      eventType: computed(() => ts('pages.logs.assetConnectivityLogsPage.columns.eventType')),
      lastSeenAt: computed(() => ts('pages.logs.assetConnectivityLogsPage.columns.lastSeenAt')),
      missCount: computed(() => ts('pages.logs.assetConnectivityLogsPage.columns.missCount')),
      thresholdMinutes: computed(() => ts('pages.logs.assetConnectivityLogsPage.columns.thresholdMinutes')),
      timestamp: computed(() => ts('pages.logs.assetConnectivityLogsPage.columns.timestamp')),
    },
    eventTypeOptions: {
      offline: computed(() => ts('pages.logs.assetConnectivityLogsPage.eventTypeOptions.offline')),
      online: computed(() => ts('pages.logs.assetConnectivityLogsPage.eventTypeOptions.online')),
    },
    statusBadge: {
      offline: computed(() => tsRaw('pages.logs.assetConnectivityLogsPage.statusBadge.offline')),
      online: computed(() => tsRaw('pages.logs.assetConnectivityLogsPage.statusBadge.online')),
    },
    drawer: {
      title: computed(() => tsTitle('pages.logs.assetConnectivityLogsPage.drawer.title')),
    },
    empty: {
      title: computed(() => tsTitle('pages.logs.assetConnectivityLogsPage.empty.title')),
      description: computed(() => tsRaw('pages.logs.assetConnectivityLogsPage.empty.description')),
    },
    pagination: {
      newer: computed(() => ts('pages.logs.assetConnectivityLogsPage.pagination.newer')),
      older: computed(() => ts('pages.logs.assetConnectivityLogsPage.pagination.older')),
    },
    messages: {
      loadFailed: computed(() => tsRaw('pages.logs.assetConnectivityLogsPage.messages.loadFailed')),
    },
    defaults: {
      unknown: computed(() => tsRaw('pages.logs.assetConnectivityLogsPage.defaults.unknown')),
      notAvailable: computed(() => tsRaw('pages.logs.assetConnectivityLogsPage.defaults.notAvailable')),
    },
    actions: {
      trackEvent: computed(() => ts('pages.logs.assetConnectivityLogsPage.actions.trackEvent')),
      viewJson: computed(() => ts('pages.logs.assetConnectivityLogsPage.actions.viewJson')),
      viewAsset: computed(() => ts('pages.logs.assetConnectivityLogsPage.actions.viewAsset')),
    },
  };
}
