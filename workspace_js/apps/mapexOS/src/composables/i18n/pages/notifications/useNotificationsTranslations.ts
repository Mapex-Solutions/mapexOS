import { computed } from 'vue';
import type { DataRowColumn } from '@components/cards';
import type { FilterListItem } from '@components/filters';
import type { ListHeaderMenuColumn } from '@components/headers';
import { useTS } from '@utils/translation';

export function useNotificationsTranslations() {
  const ts = useTS({ capitalize: true });
  const tsTitle = useTS({ titleCase: true });
  const tsRaw = useTS({ capitalize: false });

  return {
    page: {
      title: computed(() => tsTitle('pages.notifications.notifications.title')),
      description: computed(() => tsRaw('pages.notifications.notifications.description')),
      listTitle: computed(() => tsTitle('pages.notifications.notifications.listTitle')),
      button: {
        add: computed(() => ts('pages.notifications.notifications.button.add')),
      },
    },
    columns: computed(() => {
      return [
        {
          key: 'name',
          label: ts('pages.notifications.notifications.columns.name'),
          type: 'text',
          visible: 'always',
          width: 250,
        },
        {
          key: 'type',
          label: ts('pages.notifications.notifications.columns.type'),
          type: 'chip',
          visible: 'laptop', // 💻🖥️ Laptop & Desktop (>= 1024px)
          width: 120,
        },
        {
          key: 'channels',
          label: ts('pages.notifications.notifications.columns.channels'),
          type: 'chip',
          visible: 'laptop', // 💻🖥️ Laptop & Desktop (>= 1024px)
          width: 150,
        },
        {
          key: 'created',
          label: ts('pages.notifications.notifications.columns.created'),
          type: 'text',
          visible: 'laptop', // 🖥️ Desktop only (>= 1440px)
          width: 140,
        },
        {
          key: 'status',
          label: ts('pages.notifications.notifications.columns.status'),
          type: 'badge',
          visible: 'always',
          width: 100,
        },
      ] as DataRowColumn[];
    }),
    menuColumns: computed(() => {
      return [
        {
          key: 'type',
          label: ts('pages.notifications.notifications.columns.type'),
          visible: true,
        },
        {
          key: 'channels',
          label: ts('pages.notifications.notifications.columns.channels'),
          visible: true,
        },
        {
          key: 'created',
          label: ts('pages.notifications.notifications.columns.created'),
          visible: true,
        },
      ] as ListHeaderMenuColumn[];
    }),
    menuLabels: {
      singular: computed(() => ts('pages.notifications.notifications.menuLabels.singular')),
      plural: computed(() => ts('pages.notifications.notifications.menuLabels.plural')),
    },
    filters: computed(() => {
      return [
        {
          key: 'search',
          type: 'input',
          label: ts('pages.notifications.notifications.filters.search'),
          icon: 'search',
          grid: 'col-xs-12 col-sm-6 col-md-6',
        },
        {
          key: 'tenant',
          type: 'select',
          label: ts('pages.notifications.notifications.filters.tenant'),
          icon: 'domain',
          options: [],
          grid: 'col-xs-12 col-sm-6 col-md-6',
        },
        {
          key: 'type',
          type: 'select',
          label: ts('pages.notifications.notifications.filters.type'),
          icon: 'category',
          options: [],
          grid: 'col-xs-12 col-sm-6 col-md-6',
        },
      ] as FilterListItem[];
    }),
    empty: {
      title: computed(() => tsRaw('pages.notifications.notifications.empty.title')),
      description: computed(() => tsRaw('pages.notifications.notifications.empty.description')),
    },
    dialog: {
      confirmDelete: {
        title: computed(() => ts('pages.notifications.notifications.dialog.confirmDelete.title')),
        message: (name: string) => tsRaw('pages.notifications.notifications.dialog.confirmDelete.message', { name }),
      },
    },
    notifications: {
      deleteSuccess: computed(() => ts('pages.notifications.notifications.notifications.deleteSuccess')),
    },
    channelCount: {
      singular: computed(() => tsRaw('pages.notifications.notifications.channelCount.singular')),
      plural: computed(() => tsRaw('pages.notifications.notifications.channelCount.plural')),
      none: computed(() => ts('pages.notifications.notifications.channelCount.none')),
    },
  };
}
