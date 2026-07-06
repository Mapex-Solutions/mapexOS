import type { PageHeaderInfo } from '@components/headers';

import { computed } from 'vue';
import { useTS } from '@utils/translation';

/**
 * Migration plans list page translations.
 *
 * Structure mirrors:
 * - File: src/pages/assets/templateMigrations/migrationPlansListPage/MigrationPlansListPage.vue
 * - JSON: src/i18n/{locale}/pages/templateMigrations.json
 */
export function useMigrationPlansTranslations() {
  const ts = useTS({ capitalize: true });
  const tsTitle = useTS({ titleCase: true });
  const tsRaw = useTS({ capitalize: false });

  const base = 'pages.templateMigrations';

  return {
    pageHeader: {
      title: computed(() => tsTitle(`${base}.pageHeader.title`)),
      description: computed(() => ts(`${base}.pageHeader.description`)),
      button: computed(() => ts(`${base}.pageHeader.button`)),
      info: computed((): PageHeaderInfo => ({
        title: ts(`${base}.pageHeader.info.title`),
        description: ts(`${base}.pageHeader.info.description`),
        items: [
          {
            icon: 'view_module',
            color: 'blue-6',
            title: ts(`${base}.pageHeader.info.items.batches.title`),
            text: ts(`${base}.pageHeader.info.items.batches.text`),
          },
          {
            icon: 'event',
            color: 'green-6',
            title: ts(`${base}.pageHeader.info.items.scheduling.title`),
            text: ts(`${base}.pageHeader.info.items.scheduling.text`),
          },
          {
            icon: 'insights',
            color: 'orange-6',
            title: ts(`${base}.pageHeader.info.items.tracking.title`),
            text: ts(`${base}.pageHeader.info.items.tracking.text`),
          },
        ],
        docsUrl: 'https://docs.mapexos.com/asset-templates',
        docsLabel: ts(`${base}.pageHeader.info.docsLabel`),
      })),
    },

    filters: {
      label: computed(() => ts(`${base}.filters.label`)),
      searchPlaceholder: computed(() => ts(`${base}.filters.searchPlaceholder`)),
      status: computed(() => ts(`${base}.filters.status`)),
      allStatus: computed(() => ts(`${base}.filters.allStatus`)),
      clearAll: computed(() => ts(`${base}.filters.clearAll`)),
    },

    listHeader: {
      title: computed(() => tsTitle(`${base}.listHeader.title`)),
      itemLabel: computed(() => tsRaw(`${base}.listHeader.itemLabel`)),
      itemLabelPlural: computed(() => tsRaw(`${base}.listHeader.itemLabelPlural`)),
    },

    columns: {
      name: computed(() => ts(`${base}.columns.name`)),
      route: computed(() => ts(`${base}.columns.route`)),
      schedule: computed(() => ts(`${base}.columns.schedule`)),
      status: computed(() => ts(`${base}.columns.status`)),
      progress: computed(() => ts(`${base}.columns.progress`)),
    },

    schedule: {
      immediate: computed(() => ts(`${base}.schedule.immediate`)),
    },

    /** Status display label, keyed by the wire status value */
    statusLabel: (status: string) => ts(`${base}.statuses.${status}`),

    actions: {
      view: computed(() => ts(`${base}.actions.view`)),
      viewHint: computed(() => ts(`${base}.actions.viewHint`)),
      edit: computed(() => ts(`${base}.actions.edit`)),
      editHint: computed(() => ts(`${base}.actions.editHint`)),
      remove: computed(() => ts(`${base}.actions.remove`)),
      removeHint: computed(() => ts(`${base}.actions.removeHint`)),
    },

    cancelDialog: {
      title: computed(() => ts(`${base}.cancelDialog.title`)),
      message: computed(() => ts(`${base}.cancelDialog.message`)),
      confirm: computed(() => ts(`${base}.cancelDialog.confirm`)),
      success: computed(() => ts(`${base}.cancelDialog.success`)),
      error: computed(() => ts(`${base}.cancelDialog.error`)),
    },

    empty: {
      title: computed(() => ts(`${base}.empty.title`)),
      description: computed(() => ts(`${base}.empty.description`)),
    },

    errors: {
      apiNotInitialized: computed(() => ts(`${base}.errors.apiNotInitialized`)),
    },

    notifications: {
      loadError: computed(() => ts(`${base}.notifications.loadError`)),
    },
  };
}
