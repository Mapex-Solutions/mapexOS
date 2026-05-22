/** TYPE IMPORTS */
import type { DataRowColumn } from '@components/cards';

/** VUE IMPORTS */
import { computed } from 'vue';

/** UTILS */
import { useTS } from '@utils/translation';

const bp = 'pages.automations.workflowInstanceList';

/**
 * Translations composable for the Workflow Instance List page
 *
 * @returns {object} Reactive translation objects for page header, filters, columns, drawer, etc.
 */
export function useWorkflowInstanceListPageTranslations() {
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
          { title: ts(`${bp}.page.info.items.binding.title`), text: ts(`${bp}.page.info.items.binding.text`) },
          { title: ts(`${bp}.page.info.items.inputs.title`), text: ts(`${bp}.page.info.items.inputs.text`) },
          { title: ts(`${bp}.page.info.items.trigger.title`), text: ts(`${bp}.page.info.items.trigger.text`) },
          { title: ts(`${bp}.page.info.items.reuse.title`), text: ts(`${bp}.page.info.items.reuse.text`) },
        ],
        docsLabel: ts(`${bp}.page.info.docsLabel`),
      })),
    },
    buttons: {
      newInstance: computed(() => tsTitle(`${bp}.buttons.newInstance`)),
      createInstance: computed(() => tsTitle(`${bp}.buttons.createInstance`)),
    },
    filters: {
      label: computed(() => ts(`${bp}.filters.label`)),
      searchPlaceholder: computed(() => tsRaw(`${bp}.filters.searchPlaceholder`)),
      advancedFilters: computed(() => ts(`${bp}.filters.advancedFilters`)),
      pendingFilters: computed(() => ts(`${bp}.filters.pendingFilters`)),
      clearAll: computed(() => tsRaw(`${bp}.filters.clearAll`)),
      allStatus: computed(() => tsTitle(`${bp}.filters.allStatus`)),
      instanceName: computed(() => ts(`${bp}.filters.instanceName`)),
      status: computed(() => tsTitle(`${bp}.filters.status`)),
      workflow: computed(() => tsTitle(`${bp}.filters.workflow`)),
      allWorkflows: computed(() => tsTitle(`${bp}.filters.allWorkflows`)),
      uniqueExecution: computed(() => ts(`${bp}.filters.uniqueExecution`)),
      options: {
        enabled: computed(() => ts(`${bp}.filters.options.enabled`)),
        disabled: computed(() => ts(`${bp}.filters.options.disabled`)),
        yes: computed(() => ts(`${bp}.filters.options.yes`)),
        no: computed(() => ts(`${bp}.filters.options.no`)),
      },
    },
    menuColumns: {
      definitionName: computed(() => tsTitle(`${bp}.menuColumns.definitionName`)),
      inputsCount: computed(() => tsTitle(`${bp}.menuColumns.inputsCount`)),
      uniqueExecution: computed(() => tsTitle(`${bp}.menuColumns.uniqueExecution`)),
      createdAt: computed(() => tsTitle(`${bp}.menuColumns.createdAt`)),
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
    executeResult: {
      title: computed(() => tsTitle(`${bp}.executeResult.title`)),
      executing: computed(() => ts(`${bp}.executeResult.executing`)),
      statusLabel: computed(() => tsTitle(`${bp}.executeResult.statusLabel`)),
      uuidLabel: computed(() => tsTitle(`${bp}.executeResult.uuidLabel`)),
      viewExecutionsButton: computed(() => tsTitle(`${bp}.executeResult.viewExecutionsButton`)),
    },
    drawer: {
      title: computed(() => tsTitle(`${bp}.drawer.title`)),
      close: computed(() => ts(`${bp}.drawer.close`)),
      edit: computed(() => tsTitle(`${bp}.drawer.edit`)),
      loading: computed(() => ts(`${bp}.drawer.loading`)),
      error: computed(() => ts(`${bp}.drawer.error`)),
      sections: {
        basicInfo: computed(() => tsTitle(`${bp}.drawer.sections.basicInfo`)),
        workflow: computed(() => tsTitle(`${bp}.drawer.sections.workflow`)),
        inputs: computed(() => tsTitle(`${bp}.drawer.sections.inputs`)),
        timestamps: computed(() => tsTitle(`${bp}.drawer.sections.timestamps`)),
      },
      fields: {
        name: computed(() => tsTitle(`${bp}.drawer.fields.name`)),
        description: computed(() => tsTitle(`${bp}.drawer.fields.description`)),
        status: computed(() => tsTitle(`${bp}.drawer.fields.status`)),
        unique: computed(() => tsTitle(`${bp}.drawer.fields.unique`)),
        instanceUUID: computed(() => tsTitle(`${bp}.drawer.fields.instanceUUID`)),
        definition: computed(() => tsTitle(`${bp}.drawer.fields.definition`)),
        version: computed(() => tsTitle(`${bp}.drawer.fields.version`)),
        inputsCount: computed(() => tsTitle(`${bp}.drawer.fields.inputsCount`)),
        noInputs: computed(() => tsRaw(`${bp}.drawer.fields.noInputs`)),
        created: computed(() => tsTitle(`${bp}.drawer.fields.created`)),
        updated: computed(() => tsTitle(`${bp}.drawer.fields.updated`)),
      },
      values: {
        enabled: computed(() => ts(`${bp}.drawer.values.enabled`)),
        disabled: computed(() => ts(`${bp}.drawer.values.disabled`)),
        unique: computed(() => ts(`${bp}.drawer.values.unique`)),
        notApplicable: computed(() => ts(`${bp}.drawer.values.notApplicable`)),
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
        icon: () => 'play_circle',
        color: (_val: unknown, row: Record<string, unknown>) =>
          row.enabled ? 'teal-7' : 'grey-5',
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
        key: 'definitionName',
        label: ts(`${bp}.columns.definition`),
        type: 'chip',
        visible: 'always',
        width: 180,
        format: (val: unknown) => (typeof val === 'string' && val) ? val : '—',
        color: () => 'blue-6',
      },
      {
        key: 'inputsCount',
        label: ts(`${bp}.columns.inputsCount`),
        type: 'text',
        visible: 'laptop',
        width: 80,
        format: (val: unknown) => String(Number(val) || 0),
        align: 'center',
      },
      {
        key: 'uniqueExecution',
        label: ts(`${bp}.columns.unique`),
        type: 'chip',
        visible: 'laptop',
        width: 100,
        format: (val: unknown) => val ? ts(`${bp}.columns.uniqueYes`) : ts(`${bp}.columns.uniqueNo`),
        color: (val: unknown) => val ? 'orange-7' : 'grey-5',
        tooltip: (_val: unknown, row: Record<string, unknown>) =>
          row.uniqueExecution
            ? (row.workflowUUID as string) || ts(`${bp}.columns.uniqueYes`)
            : ts(`${bp}.columns.uniqueNoTooltip`),
      },
    ]),
  };
}
