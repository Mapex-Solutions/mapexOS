import type { FilterField } from '@components/drawers';

import { computed } from 'vue';
import { useTS } from '@utils/translation';

/**
 * Translation composable for Asset Raw Logs Page
 * Provides all translated strings for the raw events logs interface
 * @returns {Object} Translation object with page header, filters, columns, etc.
 */
export function useAssetRawLogsPageTranslations() {
  const ts = useTS({ capitalize: true });
  const tsTitle = useTS({ titleCase: true });
  const tsRaw = useTS({ capitalize: false });

  return {
    pageHeader: {
      title: computed(() => tsTitle('pages.logs.assetRawLogsPage.pageHeader.title')),
      description: computed(() => tsRaw('pages.logs.assetRawLogsPage.pageHeader.description')),
    },
    listTitle: computed(() => tsTitle('pages.logs.assetRawLogsPage.listTitle')),
    itemLabel: computed(() => ts('pages.logs.assetRawLogsPage.itemLabel')),
    itemLabelPlural: computed(() => ts('pages.logs.assetRawLogsPage.itemLabelPlural')),
    filters: {
      label: computed(() => ts('pages.logs.assetRawLogsPage.filters.label')),
      searchPlaceholder: computed(() => ts('pages.logs.assetRawLogsPage.filters.searchPlaceholder')),
      advancedFilters: computed(() => ts('pages.logs.assetRawLogsPage.filters.advancedFilters')),
      pendingFilters: computed(() => ts('pages.logs.assetRawLogsPage.filters.pendingFilters')),
      clearAll: computed(() => ts('pages.logs.assetRawLogsPage.filters.clearAll')),
      allStatus: computed(() => ts('pages.logs.assetRawLogsPage.filters.allStatus')),
      dateRange: computed(() => ts('pages.logs.assetRawLogsPage.filters.dateRange')),
      source: computed(() => ts('pages.logs.assetRawLogsPage.filters.source')),
      status: computed(() => ts('pages.logs.assetRawLogsPage.filters.status')),
      uuid: computed(() => ts('pages.logs.assetRawLogsPage.filters.uuid')),
      includeChildren: computed(() => ts('pages.logs.assetRawLogsPage.filters.includeChildren')),
      includeChildrenOrgs: computed(() => ts('pages.logs.assetRawLogsPage.filters.includeChildrenOrgs')),
      options: {
        yes: computed(() => ts('pages.logs.assetRawLogsPage.filters.options.yes')),
        no: computed(() => ts('pages.logs.assetRawLogsPage.filters.options.no')),
        success: computed(() => ts('pages.logs.assetRawLogsPage.filters.options.success')),
        failed: computed(() => ts('pages.logs.assetRawLogsPage.filters.options.failed')),
      },
    },

    advancedFilters: computed((): FilterField[] => [
      {
        key: 'includeChildren',
        type: 'toggle',
        label: ts('pages.logs.assetRawLogsPage.filters.includeChildrenOrgs'),
        icon: 'account_tree',
        options: [
          { label: ts('pages.logs.assetRawLogsPage.filters.allStatus'), value: null },
          { label: ts('pages.logs.assetRawLogsPage.filters.options.yes'), value: true },
          { label: ts('pages.logs.assetRawLogsPage.filters.options.no'), value: false },
        ],
      },
      {
        key: 'startDate',
        type: 'input',
        label: ts('pages.logs.assetRawLogsPage.filters.startDate'),
        icon: 'event',
        inputType: 'date',
      },
      {
        key: 'endDate',
        type: 'input',
        label: ts('pages.logs.assetRawLogsPage.filters.endDate'),
        icon: 'event',
        inputType: 'date',
      },
      {
        key: 'source',
        type: 'select',
        label: ts('pages.logs.assetRawLogsPage.filters.source'),
        icon: 'dns',
        options: [
          { label: ts('pages.logs.assetRawLogsPage.filters.allStatus'), value: null },
          { label: tsRaw('pages.logs.assetRawLogsPage.sourceOptions.httpGateway'), value: 'http_gateway' },
          { label: tsRaw('pages.logs.assetRawLogsPage.sourceOptions.mqttGateway'), value: 'mqtt_gateway' },
        ],
      },
      {
        key: 'threadId',
        type: 'input',
        label: ts('pages.logs.assetRawLogsPage.filters.uuid'),
        icon: 'link',
        placeholder: ts('pages.logs.assetRawLogsPage.filters.uuidPlaceholder'),
      },
    ]),

    columns: {
      dataSource: computed(() => ts('pages.logs.assetRawLogsPage.columns.dataSource')),
      uuid: computed(() => ts('pages.logs.assetRawLogsPage.columns.uuid')),
      status: computed(() => ts('pages.logs.assetRawLogsPage.columns.status')),
      source: computed(() => ts('pages.logs.assetRawLogsPage.columns.source')),
      timestamp: computed(() => ts('pages.logs.assetRawLogsPage.columns.timestamp')),
      retention: computed(() => ts('pages.logs.assetRawLogsPage.columns.retention')),
    },
    sourceOptions: {
      httpGateway: computed(() => tsRaw('pages.logs.assetRawLogsPage.sourceOptions.httpGateway')),
      mqttGateway: computed(() => tsRaw('pages.logs.assetRawLogsPage.sourceOptions.mqttGateway')),
      lorawanGateway: computed(() => tsRaw('pages.logs.assetRawLogsPage.sourceOptions.lorawanGateway')),
    },
    statusOptions: {
      success: computed(() => ts('pages.logs.assetRawLogsPage.statusOptions.success')),
      failed: computed(() => ts('pages.logs.assetRawLogsPage.statusOptions.failed')),
    },
    statusBadge: {
      success: computed(() => tsRaw('pages.logs.assetRawLogsPage.statusBadge.success')),
      failed: computed(() => tsRaw('pages.logs.assetRawLogsPage.statusBadge.failed')),
    },
    drawer: {
      title: computed(() => tsTitle('pages.logs.assetRawLogsPage.drawer.title')),
    },
    empty: {
      title: computed(() => tsTitle('pages.logs.assetRawLogsPage.empty.title')),
      description: computed(() => tsRaw('pages.logs.assetRawLogsPage.empty.description')),
    },
    pagination: {
      newer: computed(() => ts('pages.logs.assetRawLogsPage.pagination.newer')),
      older: computed(() => ts('pages.logs.assetRawLogsPage.pagination.older')),
    },
    messages: {
      loadFailed: computed(() => tsRaw('pages.logs.assetRawLogsPage.messages.loadFailed')),
    },
    defaults: {
      unknown: computed(() => tsRaw('pages.logs.assetRawLogsPage.defaults.unknown')),
      notAvailable: computed(() => tsRaw('pages.logs.assetRawLogsPage.defaults.notAvailable')),
    },
    actions: {
      trackEvent: computed(() => ts('pages.logs.assetRawLogsPage.actions.trackEvent')),
      viewJson: computed(() => ts('pages.logs.assetRawLogsPage.actions.viewJson')),
    },
  };
}
