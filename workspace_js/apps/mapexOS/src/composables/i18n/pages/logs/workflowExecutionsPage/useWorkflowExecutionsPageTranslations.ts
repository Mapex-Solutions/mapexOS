import { computed } from 'vue';
import { useTS } from '@utils/translation';

/**
 * Translation composable for Workflow Executions Page
 * Provides all translated strings for the workflow executions log interface
 * @returns Translation object with page header, filters, columns, modal, etc.
 */
export function useWorkflowExecutionsPageTranslations() {
  const ts = useTS({ capitalize: true });
  const tsTitle = useTS({ titleCase: true });
  const tsRaw = useTS({ capitalize: false });

  return {
    pageHeader: {
      title: computed(() => tsTitle('pages.logs.workflowExecutionsPage.pageHeader.title')),
      description: computed(() => tsRaw('pages.logs.workflowExecutionsPage.pageHeader.description')),
    },
    listTitle: computed(() => tsTitle('pages.logs.workflowExecutionsPage.listTitle')),
    itemLabel: computed(() => ts('pages.logs.workflowExecutionsPage.itemLabel')),
    itemLabelPlural: computed(() => ts('pages.logs.workflowExecutionsPage.itemLabelPlural')),

    filters: {
      label: computed(() => ts('pages.logs.workflowExecutionsPage.filters.label')),
      searchPlaceholder: computed(() => ts('pages.logs.workflowExecutionsPage.filters.searchPlaceholder')),
      advancedFilters: computed(() => ts('pages.logs.workflowExecutionsPage.filters.advancedFilters')),
      pendingFilters: computed(() => ts('pages.logs.workflowExecutionsPage.filters.pendingFilters')),
      clearAll: computed(() => ts('pages.logs.workflowExecutionsPage.filters.clearAll')),
      allStatus: computed(() => ts('pages.logs.workflowExecutionsPage.filters.allStatus')),
      dateRange: computed(() => ts('pages.logs.workflowExecutionsPage.filters.dateRange')),
      startDate: computed(() => ts('pages.logs.workflowExecutionsPage.filters.startDate')),
      endDate: computed(() => ts('pages.logs.workflowExecutionsPage.filters.endDate')),
      status: computed(() => ts('pages.logs.workflowExecutionsPage.filters.status')),
      instanceId: computed(() => ts('pages.logs.workflowExecutionsPage.filters.instanceId')),
      instanceIdPlaceholder: computed(() => ts('pages.logs.workflowExecutionsPage.filters.instanceIdPlaceholder')),
      definitionId: computed(() => ts('pages.logs.workflowExecutionsPage.filters.definitionId')),
      definitionIdPlaceholder: computed(() => ts('pages.logs.workflowExecutionsPage.filters.definitionIdPlaceholder')),
      workflow: computed(() => ts('pages.logs.workflowExecutionsPage.filters.workflow')),
      workflowPlaceholder: computed(() => ts('pages.logs.workflowExecutionsPage.filters.workflowPlaceholder')),
      includeChildren: computed(() => ts('pages.logs.workflowExecutionsPage.filters.includeChildren')),
      includeChildrenOrgs: computed(() => ts('pages.logs.workflowExecutionsPage.filters.includeChildrenOrgs')),
      options: {
        yes: computed(() => ts('pages.logs.workflowExecutionsPage.filters.options.yes')),
        no: computed(() => ts('pages.logs.workflowExecutionsPage.filters.options.no')),
      },
    },

    statusOptions: {
      all: computed(() => ts('pages.logs.workflowExecutionsPage.statusOptions.all')),
      running: computed(() => ts('pages.logs.workflowExecutionsPage.statusOptions.running')),
      waiting: computed(() => ts('pages.logs.workflowExecutionsPage.statusOptions.waiting')),
      created: computed(() => ts('pages.logs.workflowExecutionsPage.statusOptions.created')),
      completed: computed(() => ts('pages.logs.workflowExecutionsPage.statusOptions.completed')),
      failed: computed(() => ts('pages.logs.workflowExecutionsPage.statusOptions.failed')),
      cancelled: computed(() => ts('pages.logs.workflowExecutionsPage.statusOptions.cancelled')),
    },

    columns: {
      workflowName: computed(() => ts('pages.logs.workflowExecutionsPage.columns.workflowName')),
      instanceName: computed(() => ts('pages.logs.workflowExecutionsPage.columns.instanceName')),
      definitionName: computed(() => ts('pages.logs.workflowExecutionsPage.columns.definitionName')),
      triggerSource: computed(() => ts('pages.logs.workflowExecutionsPage.columns.triggerSource')),
      status: computed(() => ts('pages.logs.workflowExecutionsPage.columns.status')),
      duration: computed(() => ts('pages.logs.workflowExecutionsPage.columns.duration')),
      instanceId: computed(() => ts('pages.logs.workflowExecutionsPage.columns.instanceId')),
      timestamp: computed(() => ts('pages.logs.workflowExecutionsPage.columns.timestamp')),
    },

    statusBadge: {
      running: computed(() => tsRaw('pages.logs.workflowExecutionsPage.statusBadge.running')),
      waiting: computed(() => tsRaw('pages.logs.workflowExecutionsPage.statusBadge.waiting')),
      created: computed(() => tsRaw('pages.logs.workflowExecutionsPage.statusBadge.created')),
      completed: computed(() => tsRaw('pages.logs.workflowExecutionsPage.statusBadge.completed')),
      failed: computed(() => tsRaw('pages.logs.workflowExecutionsPage.statusBadge.failed')),
      cancelled: computed(() => tsRaw('pages.logs.workflowExecutionsPage.statusBadge.cancelled')),
    },

    modal: {
      title: computed(() => tsTitle('pages.logs.workflowExecutionsPage.modal.title')),
      tabs: {
        dag: computed(() => tsRaw('pages.logs.workflowExecutionsPage.modal.tabs.dag')),
        event: computed(() => ts('pages.logs.workflowExecutionsPage.modal.tabs.event')),
        state: computed(() => ts('pages.logs.workflowExecutionsPage.modal.tabs.state')),
        logs: computed(() => ts('pages.logs.workflowExecutionsPage.modal.tabs.logs')),
      },
      legend: {
        completed: computed(() => ts('pages.logs.workflowExecutionsPage.modal.legend.completed')),
        waiting: computed(() => ts('pages.logs.workflowExecutionsPage.modal.legend.waiting')),
        retrying: computed(() => ts('pages.logs.workflowExecutionsPage.modal.legend.retrying')),
        timeout: computed(() => ts('pages.logs.workflowExecutionsPage.modal.legend.timeout')),
        error: computed(() => ts('pages.logs.workflowExecutionsPage.modal.legend.error')),
        cancelled: computed(() => ts('pages.logs.workflowExecutionsPage.modal.legend.cancelled')),
        notReached: computed(() => ts('pages.logs.workflowExecutionsPage.modal.legend.notReached')),
      },
      states: {
        workflowState: computed(() => ts('pages.logs.workflowExecutionsPage.modal.states.workflowState')),
        eventPayload: computed(() => ts('pages.logs.workflowExecutionsPage.modal.states.eventPayload')),
        nodeOutputs: computed(() => ts('pages.logs.workflowExecutionsPage.modal.states.nodeOutputs')),
        colName: computed(() => ts('pages.logs.workflowExecutionsPage.modal.states.colName')),
        colType: computed(() => ts('pages.logs.workflowExecutionsPage.modal.states.colType')),
        colDefault: computed(() => ts('pages.logs.workflowExecutionsPage.modal.states.colDefault')),
        colCurrent: computed(() => ts('pages.logs.workflowExecutionsPage.modal.states.colCurrent')),
        durable: computed(() => ts('pages.logs.workflowExecutionsPage.modal.states.durable')),
      },
      drawer: {
        status: computed(() => ts('pages.logs.workflowExecutionsPage.modal.drawer.status')),
        nodeType: computed(() => ts('pages.logs.workflowExecutionsPage.modal.drawer.nodeType')),
        nodeId: computed(() => ts('pages.logs.workflowExecutionsPage.modal.drawer.nodeId')),
        duration: computed(() => ts('pages.logs.workflowExecutionsPage.modal.drawer.duration')),
        outputHandle: computed(() => ts('pages.logs.workflowExecutionsPage.modal.drawer.outputHandle')),
        configuration: computed(() => ts('pages.logs.workflowExecutionsPage.modal.drawer.configuration')),
        configurationEmpty: computed(() => ts('pages.logs.workflowExecutionsPage.modal.drawer.configurationEmpty')),
        error: computed(() => ts('pages.logs.workflowExecutionsPage.modal.drawer.error')),
        outputs: computed(() => ts('pages.logs.workflowExecutionsPage.modal.drawer.outputs')),
      },
      noData: computed(() => tsRaw('pages.logs.workflowExecutionsPage.modal.noData')),
    },

    empty: {
      title: computed(() => tsTitle('pages.logs.workflowExecutionsPage.empty.title')),
      description: computed(() => tsRaw('pages.logs.workflowExecutionsPage.empty.description')),
    },

    pagination: {
      newer: computed(() => ts('pages.logs.workflowExecutionsPage.pagination.newer')),
      older: computed(() => ts('pages.logs.workflowExecutionsPage.pagination.older')),
    },

    messages: {
      loadFailed: computed(() => tsRaw('pages.logs.workflowExecutionsPage.messages.loadFailed')),
    },

    defaults: {
      unknown: computed(() => tsRaw('pages.logs.workflowExecutionsPage.defaults.unknown')),
      notAvailable: computed(() => tsRaw('pages.logs.workflowExecutionsPage.defaults.notAvailable')),
    },
  };
}
