import type { FilterField } from '@components/drawers';
import type { DataRowColumn } from '@components/cards';

import { computed } from 'vue';
import { useTS } from '@utils/translation';

export function useEventStorePageTranslations() {
  const ts = useTS({ capitalize: true });
  const tsTitle = useTS({ titleCase: true });
  const tsRaw = useTS({ capitalize: false });

  return {
    pageHeader: {
      title: computed(() => tsTitle('pages.events.eventStore.pageHeader.title')),
      description: computed(() => tsRaw('pages.events.eventStore.pageHeader.description')),
    },

    listHeader: {
      title: computed(() => tsTitle('pages.events.eventStore.listHeader.title')),
      itemLabel: computed(() => ts('pages.events.eventStore.listHeader.itemLabel')),
      itemLabelPlural: computed(() => ts('pages.events.eventStore.listHeader.itemLabelPlural')),
    },

    filters: {
      label: computed(() => ts('pages.events.eventStore.filters.label')),
      searchPlaceholder: computed(() => ts('pages.events.eventStore.filters.searchPlaceholder')),
      allStatus: computed(() => ts('pages.events.eventStore.filters.allStatus')),
      advancedFilters: computed(() => ts('pages.events.eventStore.filters.advancedFilters')),
      dynamicFilters: computed(() => ts('pages.events.eventStore.filters.dynamicFilters')),
      pendingFilters: computed(() => ts('pages.events.eventStore.filters.pendingFilters')),
      clearAll: computed(() => ts('pages.events.eventStore.filters.clearAll')),
      includeChildrenOrgs: computed(() => ts('pages.events.eventStore.filters.includeChildrenOrgs')),
      startDate: computed(() => ts('pages.events.eventStore.filters.startDate')),
      endDate: computed(() => ts('pages.events.eventStore.filters.endDate')),
      threadId: computed(() => ts('pages.events.eventStore.filters.threadId')),
      assetId: computed(() => ts('pages.events.eventStore.filters.assetId')),
      assetTemplateId: computed(() => ts('pages.events.eventStore.filters.assetTemplateId')),
      asset: computed(() => ts('pages.events.eventStore.filters.asset')),
      assetTemplate: computed(() => ts('pages.events.eventStore.filters.assetTemplate')),
      eventType: computed(() => ts('pages.events.eventStore.filters.eventType')),
      source: computed(() => ts('pages.events.eventStore.filters.source')),
      options: {
        yes: computed(() => ts('pages.events.eventStore.filters.options.yes')),
        no: computed(() => ts('pages.events.eventStore.filters.options.no')),
      },
    },

    advancedFilters: computed((): FilterField[] => [
      {
        key: 'includeChildren',
        type: 'toggle',
        label: ts('pages.events.eventStore.filters.includeChildrenOrgs'),
        icon: 'account_tree',
        options: [
          { label: ts('pages.events.eventStore.filters.allStatus'), value: null },
          { label: ts('pages.events.eventStore.filters.options.yes'), value: true },
          { label: ts('pages.events.eventStore.filters.options.no'), value: false },
        ],
      },
      {
        key: 'startDate',
        type: 'input',
        label: ts('pages.events.eventStore.filters.startDate'),
        icon: 'event',
        inputType: 'date',
      },
      {
        key: 'endDate',
        type: 'input',
        label: ts('pages.events.eventStore.filters.endDate'),
        icon: 'event',
        inputType: 'date',
      },
      {
        key: 'threadId',
        type: 'input',
        label: ts('pages.events.eventStore.filters.threadId'),
        icon: 'fingerprint',
      },
      {
        key: 'source',
        type: 'input',
        label: ts('pages.events.eventStore.filters.source'),
        icon: 'dns',
      },
    ]),

    columns: computed((): DataRowColumn[] => [
      {
        key: 'icon',
        label: '',
        type: 'avatar',
        visible: 'always',
        width: 56,
        icon: () => 'storage',
        color: () => 'primary',
      },
      {
        key: 'assetName',
        label: ts('pages.events.eventStore.columns.assetName'),
        type: 'text',
        visible: 'laptop',
        width: 160,
        ellipsis: true,
      },
      {
        key: 'templateName',
        label: ts('pages.events.eventStore.columns.templateName'),
        type: 'text',
        visible: 'laptop',
        width: 160,
        ellipsis: true,
      },
      {
        key: 'threadId',
        label: ts('pages.events.eventStore.columns.threadId'),
        type: 'text',
        visible: 'always',
        width: 200,
        ellipsis: true,
      },
      {
        key: 'source',
        label: ts('pages.events.eventStore.columns.source'),
        type: 'chip',
        visible: 'laptop',
        width: 120,
        color: (value: string) => {
          if (value === 'asset') return 'blue';
          if (value === 'rule') return 'orange';
          return 'grey';
        },
      },
      {
        key: 'created',
        label: ts('pages.events.eventStore.columns.created'),
        type: 'text',
        visible: 'laptop',
        width: 180,
      },
    ]),

    menuColumns: {
      assetName: computed(() => ts('pages.events.eventStore.columns.assetName')),
      templateName: computed(() => ts('pages.events.eventStore.columns.templateName')),
      threadId: computed(() => ts('pages.events.eventStore.columns.threadId')),
      source: computed(() => ts('pages.events.eventStore.columns.source')),
      created: computed(() => ts('pages.events.eventStore.columns.created')),
    },

    empty: {
      title: computed(() => tsRaw('pages.events.eventStore.empty.title')),
      description: computed(() => tsRaw('pages.events.eventStore.empty.description')),
    },

    drawer: {
      title: computed(() => ts('pages.events.eventStore.drawer.title')),
    },

    pagination: {
      previous: computed(() => ts('pages.events.eventStore.pagination.previous')),
      next: computed(() => ts('pages.events.eventStore.pagination.next')),
    },

    notifications: {
      loadFailed: computed(() => ts('pages.events.eventStore.notifications.loadFailed')),
    },

    tour: {
      header: {
        title: computed(() => ts('pages.events.eventStore.tour.header.title')),
        description: computed(() => tsRaw('pages.events.eventStore.tour.header.description')),
      },
      searchInput: {
        title: computed(() => ts('pages.events.eventStore.tour.searchInput.title')),
        description: computed(() => tsRaw('pages.events.eventStore.tour.searchInput.description')),
      },
      advancedFiltersBtn: {
        title: computed(() => ts('pages.events.eventStore.tour.advancedFiltersBtn.title')),
        description: computed(() => tsRaw('pages.events.eventStore.tour.advancedFiltersBtn.description')),
      },
      advancedFiltersOpen: {
        title: computed(() => ts('pages.events.eventStore.tour.advancedFiltersOpen.title')),
        description: computed(() => tsRaw('pages.events.eventStore.tour.advancedFiltersOpen.description')),
      },
      dynamicFiltersBtn: {
        title: computed(() => ts('pages.events.eventStore.tour.dynamicFiltersBtn.title')),
        description: computed(() => tsRaw('pages.events.eventStore.tour.dynamicFiltersBtn.description')),
      },
      dynamicFiltersOpen: {
        title: computed(() => ts('pages.events.eventStore.tour.dynamicFiltersOpen.title')),
        description: computed(() => tsRaw('pages.events.eventStore.tour.dynamicFiltersOpen.description')),
      },
      results: {
        title: computed(() => ts('pages.events.eventStore.tour.results.title')),
        description: computed(() => tsRaw('pages.events.eventStore.tour.results.description')),
      },
    },
  };
}
