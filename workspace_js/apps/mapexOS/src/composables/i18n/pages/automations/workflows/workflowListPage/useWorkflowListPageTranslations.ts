/** TYPE IMPORTS */
import type { DataRowColumn } from '@components/cards';

/** VUE IMPORTS */
import { computed } from 'vue';

/** UTILS */
import { useTS } from '@utils/translation';

const bp = 'pages.automations.workflowList';

/**
 * Translations composable for the Workflow List page
 *
 * @returns {object} Reactive translation objects for page header, filters, columns, drawer, etc.
 */
export function useWorkflowListPageTranslations() {
  const ts = useTS({ capitalize: true });
  const tsTitle = useTS({ titleCase: true });
  const tsRaw = useTS({ capitalize: false });

  return {
    pageHeader: {
      title: computed(() => tsTitle(`${bp}.pageHeader.title`)),
      description: computed(() => ts(`${bp}.pageHeader.description`)),
    },
    page: {
      title: computed(() => tsTitle(`${bp}.pageHeader.title`)),
      description: computed(() => ts(`${bp}.pageHeader.description`)),
      listTitle: computed(() => tsTitle(`${bp}.page.listTitle`)),
      itemLabel: computed(() => tsRaw(`${bp}.page.itemLabel`)),
      itemLabelPlural: computed(() => tsRaw(`${bp}.page.itemLabelPlural`)),
      addButton: computed(() => tsTitle(`${bp}.page.addButton`)),
      info: computed(() => ({
        title: ts(`${bp}.page.info.title`),
        description: ts(`${bp}.page.info.description`),
        items: [
          { title: ts(`${bp}.page.info.items.visual.title`), text: ts(`${bp}.page.info.items.visual.text`) },
          { title: ts(`${bp}.page.info.items.plugins.title`), text: ts(`${bp}.page.info.items.plugins.text`) },
          { title: ts(`${bp}.page.info.items.triggers.title`), text: ts(`${bp}.page.info.items.triggers.text`) },
          { title: ts(`${bp}.page.info.items.variables.title`), text: ts(`${bp}.page.info.items.variables.text`) },
        ],
        docsLabel: ts(`${bp}.page.info.docsLabel`),
      })),
    },
    buttons: {
      newWorkflow: computed(() => tsTitle(`${bp}.buttons.newWorkflow`)),
      createWorkflow: computed(() => tsTitle(`${bp}.buttons.createWorkflow`)),
    },
    filters: {
      label: computed(() => ts(`${bp}.filters.label`)),
      searchPlaceholder: computed(() => tsRaw(`${bp}.filters.searchPlaceholder`)),
      advancedFilters: computed(() => ts(`${bp}.filters.advancedFilters`)),
      pendingFilters: computed(() => ts(`${bp}.filters.pendingFilters`)),
      clearAll: computed(() => tsRaw(`${bp}.filters.clearAll`)),
      allStatus: computed(() => tsTitle(`${bp}.filters.allStatus`)),
      workflowName: computed(() => ts(`${bp}.filters.workflowName`)),
      status: computed(() => tsTitle(`${bp}.filters.status`)),
      isTemplate: computed(() => tsTitle(`${bp}.filters.isTemplate`)),
      version: computed(() => tsTitle(`${bp}.filters.version`)),
      versionPlaceholder: computed(() => tsRaw(`${bp}.filters.versionPlaceholder`)),
      health: computed(() => tsTitle(`${bp}.filters.health`)),
      nodesCount: computed(() => tsTitle(`${bp}.filters.nodesCount`)),
      nodesCountPlaceholder: computed(() => tsRaw(`${bp}.filters.nodesCountPlaceholder`)),
      pluginType: computed(() => tsTitle(`${bp}.filters.pluginType`)),
      pluginTypePlaceholder: computed(() => tsRaw(`${bp}.filters.pluginTypePlaceholder`)),
      options: {
        enabled: computed(() => ts(`${bp}.filters.options.enabled`)),
        disabled: computed(() => ts(`${bp}.filters.options.disabled`)),
        yes: computed(() => ts(`${bp}.filters.options.yes`)),
        no: computed(() => ts(`${bp}.filters.options.no`)),
        valid: computed(() => ts(`${bp}.filters.options.valid`)),
        pluginMissing: computed(() => ts(`${bp}.filters.options.pluginMissing`)),
        invalid: computed(() => ts(`${bp}.filters.options.invalid`)),
      },
    },
    menuColumns: {
      version: computed(() => tsTitle(`${bp}.menuColumns.version`)),
      nodesCount: computed(() => tsTitle(`${bp}.menuColumns.nodesCount`)),
      pluginsCount: computed(() => tsTitle(`${bp}.menuColumns.pluginsCount`)),
    },
    status: {
      enabled: computed(() => ts(`${bp}.status.enabled`)),
      disabled: computed(() => ts(`${bp}.status.disabled`)),
    },
    empty: {
      title: computed(() => tsTitle(`${bp}.empty.title`)),
      description: computed(() => ts(`${bp}.empty.description`)),
    },
    emptyState: {
      title: computed(() => tsTitle(`${bp}.emptyState.title`)),
      description: computed(() => ts(`${bp}.emptyState.description`)),
    },
    badges: {
      active: computed(() => tsTitle(`${bp}.badges.active`)),
      disabled: computed(() => tsTitle(`${bp}.badges.disabled`)),
    },
    messages: {
      deletedSuccessfully: computed(() => ts(`${bp}.messages.deletedSuccessfully`)),
      confirmDelete: (name: string) => ts(`${bp}.messages.confirmDelete`).replace('{name}', name),
    },
    dialog: {
      deleteTitle: computed(() => tsTitle(`${bp}.dialog.deleteTitle`)),
    },
    drawer: {
      title: computed(() => tsTitle(`${bp}.drawer.title`)),
      close: computed(() => ts(`${bp}.drawer.close`)),
      edit: computed(() => tsTitle(`${bp}.drawer.edit`)),
      loading: computed(() => ts(`${bp}.drawer.loading`)),
      error: computed(() => ts(`${bp}.drawer.error`)),
      sections: {
        basicInfo: computed(() => tsTitle(`${bp}.drawer.sections.basicInfo`)),
        configuration: computed(() => tsTitle(`${bp}.drawer.sections.configuration`)),
        timestamps: computed(() => tsTitle(`${bp}.drawer.sections.timestamps`)),
      },
      fields: {
        name: computed(() => tsTitle(`${bp}.drawer.fields.name`)),
        description: computed(() => tsTitle(`${bp}.drawer.fields.description`)),
        status: computed(() => tsTitle(`${bp}.drawer.fields.status`)),
        version: computed(() => tsTitle(`${bp}.drawer.fields.version`)),
        nodesCount: computed(() => tsTitle(`${bp}.drawer.fields.nodesCount`)),
        edgesCount: computed(() => tsTitle(`${bp}.drawer.fields.edgesCount`)),
        timezone: computed(() => tsTitle(`${bp}.drawer.fields.timezone`)),
        isTemplate: computed(() => tsTitle(`${bp}.drawer.fields.isTemplate`)),
        created: computed(() => tsTitle(`${bp}.drawer.fields.created`)),
        updated: computed(() => tsTitle(`${bp}.drawer.fields.updated`)),
      },
      values: {
        enabled: computed(() => ts(`${bp}.drawer.values.enabled`)),
        disabled: computed(() => ts(`${bp}.drawer.values.disabled`)),
        yes: computed(() => ts(`${bp}.drawer.values.yes`)),
        no: computed(() => ts(`${bp}.drawer.values.no`)),
        noDescription: computed(() => ts(`${bp}.drawer.values.noDescription`)),
      },
    },
    columns: computed((): DataRowColumn[] => [
      {
        key: 'icon',
        label: '',
        type: 'avatar',
        visible: 'always',
        width: 56,
        icon: () => 'account_tree',
        color: (_val: unknown, row: Record<string, unknown>) =>
          row.enabled ? 'primary' : 'grey-5',
        tooltip: (_val: unknown, row: Record<string, unknown>) =>
          row.enabled
            ? ts(`${bp}.status.enabled`)
            : ts(`${bp}.status.disabled`),
      },
      {
        key: 'name',
        label: ts(`${bp}.columns.name`),
        type: 'text',
        visible: 'always',
        width: 250,
        ellipsis: true,
        secondaryKey: 'description',
      },
      {
        key: 'status',
        label: ts(`${bp}.columns.definitionStatus`),
        type: 'chip',
        visible: 'always',
        width: 120,
        format: (val: unknown) => {
          if (val === 'plugin_missing') return ts(`${bp}.definitionStatus.pluginMissing`).toUpperCase();
          if (val === 'invalid') return ts(`${bp}.definitionStatus.invalid`).toUpperCase();
          return ts(`${bp}.definitionStatus.valid`).toUpperCase();
        },
        color: (val: unknown) => {
          if (val === 'plugin_missing') return 'warning';
          if (val === 'invalid') return 'negative';
          return 'positive';
        },
      },
      {
        key: 'definitionVersion',
        label: ts(`${bp}.columns.version`),
        type: 'chip',
        visible: 'laptop',
        width: 80,
        format: (val: unknown) => `v${String(val)}`,
        color: () => 'blue-6',
      },
      {
        key: 'nodesCount',
        label: ts(`${bp}.columns.nodesCount`),
        type: 'text',
        visible: 'laptop',
        width: 80,
        format: (val: unknown) => String(val),
        align: 'center',
      },
      {
        key: 'pluginsCount',
        label: ts(`${bp}.columns.pluginsCount`),
        type: 'text',
        visible: 'laptop',
        width: 80,
        format: (val: unknown) => String(Number(val) || 0),
        align: 'center',
      },
    ]),
  };
}
