import { computed } from 'vue';
import type { DataRowColumn } from '@components/cards';
import type { FilterField } from '@components/drawers';
import type { ListHeaderMenuColumn } from '@components/headers';
import { useTS } from '@utils/translation';

export function useLakeHouseTranslations() {
  const ts = useTS({ capitalize: true });
  const tsTitle = useTS({ titleCase: true });
  const tsRaw = useTS({ capitalize: false });

  return {
    page: {
      title: computed(() => tsTitle('pages.lakeHouse.lakeHouse.title')),
      description: computed(() => tsRaw('pages.lakeHouse.lakeHouse.description')),
      listTitle: computed(() => tsTitle('pages.lakeHouse.lakeHouse.listTitle')),
      button: {
        add: computed(() => ts('pages.lakeHouse.lakeHouse.button.add')),
      },
    },
    columns: computed(() => {
      return [
        {
          key: 'name',
          label: ts('pages.lakeHouse.lakeHouse.columns.name'),
          type: 'text',
          visible: 'always',
          width: 250,
        },
        {
          key: 'bucket',
          label: ts('pages.lakeHouse.lakeHouse.columns.bucket'),
          type: 'text',
          visible: 'laptop',
          width: 180,
        },
        {
          key: 'region',
          label: ts('pages.lakeHouse.lakeHouse.columns.region'),
          type: 'chip',
          visible: 'laptop',
          width: 120,
        },
        {
          key: 'maxSize',
          label: ts('pages.lakeHouse.lakeHouse.columns.maxSize'),
          type: 'chip',
          visible: 'laptop',
          width: 120,
        },
        {
          key: 'frequency',
          label: ts('pages.lakeHouse.lakeHouse.columns.frequency'),
          type: 'text',
          visible: 'laptop',
          width: 140,
        },
        {
          key: 'status',
          label: ts('pages.lakeHouse.lakeHouse.columns.status'),
          type: 'badge',
          visible: 'always',
          width: 100,
        },
      ] as DataRowColumn[];
    }),
    menuColumns: computed(() => {
      return [
        {
          key: 'bucket',
          label: ts('pages.lakeHouse.lakeHouse.columns.bucket'),
          visible: true,
        },
        {
          key: 'region',
          label: ts('pages.lakeHouse.lakeHouse.columns.region'),
          visible: true,
        },
        {
          key: 'maxSize',
          label: ts('pages.lakeHouse.lakeHouse.columns.maxSize'),
          visible: true,
        },
        {
          key: 'frequency',
          label: ts('pages.lakeHouse.lakeHouse.columns.frequency'),
          visible: true,
        },
      ] as ListHeaderMenuColumn[];
    }),
    menuLabels: {
      singular: computed(() => ts('pages.lakeHouse.lakeHouse.menuLabels.singular')),
      plural: computed(() => ts('pages.lakeHouse.lakeHouse.menuLabels.plural')),
    },
    filters: {
      label: computed(() => ts('pages.lakeHouse.lakeHouse.filters.label')),
      searchPlaceholder: computed(() => ts('pages.lakeHouse.lakeHouse.filters.searchPlaceholder')),
      allStatus: computed(() => ts('pages.lakeHouse.lakeHouse.filters.allStatus')),
      advancedFilters: computed(() => ts('pages.lakeHouse.lakeHouse.filters.advancedFilters')),
      pendingFilters: computed(() => ts('pages.lakeHouse.lakeHouse.filters.pendingFilters')),
      clearAll: computed(() => ts('pages.lakeHouse.lakeHouse.filters.clearAll')),
      name: computed(() => ts('pages.lakeHouse.lakeHouse.filters.name')),
      type: computed(() => ts('pages.lakeHouse.lakeHouse.filters.type')),
      region: computed(() => ts('pages.lakeHouse.lakeHouse.filters.region')),
      status: computed(() => ts('pages.lakeHouse.lakeHouse.filters.status')),
      includeChildren: computed(() => ts('pages.lakeHouse.lakeHouse.filters.includeChildren')),
      includeChildrenOrgs: computed(() => ts('pages.lakeHouse.lakeHouse.filters.includeChildrenOrgs')),
      options: {
        all: computed(() => ts('pages.lakeHouse.lakeHouse.filters.options.all')),
        yes: computed(() => ts('pages.lakeHouse.lakeHouse.filters.options.yes')),
        no: computed(() => ts('pages.lakeHouse.lakeHouse.filters.options.no')),
        active: computed(() => ts('pages.lakeHouse.lakeHouse.filters.options.active')),
        inactive: computed(() => ts('pages.lakeHouse.lakeHouse.filters.options.inactive')),
        awsS3: computed(() => ts('pages.lakeHouse.lakeHouse.filters.options.awsS3')),
        azureBlob: computed(() => ts('pages.lakeHouse.lakeHouse.filters.options.azureBlob')),
        gcpStorage: computed(() => ts('pages.lakeHouse.lakeHouse.filters.options.gcpStorage')),
        minio: computed(() => ts('pages.lakeHouse.lakeHouse.filters.options.minio')),
      },
    },
    advancedFilters: computed((): FilterField[] => [
      {
        key: 'includeChildren',
        type: 'toggle',
        label: ts('pages.lakeHouse.lakeHouse.filters.includeChildrenOrgs'),
        icon: 'account_tree',
        options: [
          { label: ts('pages.lakeHouse.lakeHouse.filters.allStatus'), value: null },
          { label: ts('pages.lakeHouse.lakeHouse.filters.options.yes'), value: true },
          { label: ts('pages.lakeHouse.lakeHouse.filters.options.no'), value: false },
        ],
      },
      {
        key: 'type',
        type: 'select',
        label: ts('pages.lakeHouse.lakeHouse.filters.type'),
        icon: 'category',
        options: [
          { label: ts('pages.lakeHouse.lakeHouse.filters.allStatus'), value: null },
          { label: ts('pages.lakeHouse.lakeHouse.filters.options.awsS3'), value: 'aws-s3' },
          { label: ts('pages.lakeHouse.lakeHouse.filters.options.azureBlob'), value: 'azure-blob' },
          { label: ts('pages.lakeHouse.lakeHouse.filters.options.gcpStorage'), value: 'gcp-storage' },
          { label: ts('pages.lakeHouse.lakeHouse.filters.options.minio'), value: 'minio' },
        ],
      },
      {
        key: 'region',
        type: 'input',
        label: ts('pages.lakeHouse.lakeHouse.filters.region'),
        icon: 'public',
      },
    ]),
    empty: {
      title: computed(() => tsRaw('pages.lakeHouse.lakeHouse.empty.title')),
      description: computed(() => tsRaw('pages.lakeHouse.lakeHouse.empty.description')),
    },
    dialog: {
      confirmDelete: {
        title: computed(() => ts('pages.lakeHouse.lakeHouse.dialog.confirmDelete.title')),
        message: (name: string) => tsRaw('pages.lakeHouse.lakeHouse.dialog.confirmDelete.message', { name }),
      },
    },
    notifications: {
      deleteSuccess: computed(() => ts('pages.lakeHouse.lakeHouse.notifications.deleteSuccess')),
    },
    status: {
      active: computed(() => ts('pages.lakeHouse.lakeHouse.status.active')),
      inactive: computed(() => ts('pages.lakeHouse.lakeHouse.status.inactive')),
      notAvailable: computed(() => ts('pages.lakeHouse.lakeHouse.status.notAvailable')),
    },
  };
}
