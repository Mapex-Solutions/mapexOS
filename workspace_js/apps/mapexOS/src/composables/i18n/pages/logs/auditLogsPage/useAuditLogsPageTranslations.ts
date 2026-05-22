import { computed } from 'vue';
import { useTS } from '@utils/translation';
import type { DataRowColumn } from '@components/cards';

/**
 * Translation composable for Audit Logs Page
 *
 * Structure mirrors:
 * - File: src/pages/logs/auditLogsPage/AuditLogsPage.vue
 * - JSON: src/i18n/{locale}/pages/logs/auditLogsPage.json
 * - Composable: src/composables/i18n/pages/logs/auditLogsPage/useAuditLogsPageTranslations.ts
 *
 * @returns {Object} Translation object with page header, filters, columns, etc.
 */
export function useAuditLogsPageTranslations() {
  const ts = useTS({ capitalize: true });
  const tsTitle = useTS({ titleCase: true });
  const tsRaw = useTS({ capitalize: false });

  return {
    /**
     * Page header translations
     */
    pageHeader: {
      title: computed(() => tsTitle('pages.logs.auditLogsPage.pageHeader.title')),
      description: computed(() => tsRaw('pages.logs.auditLogsPage.pageHeader.description')),
    },

    /**
     * List section translations
     */
    listTitle: computed(() => tsTitle('pages.logs.auditLogsPage.listTitle')),
    itemLabel: computed(() => ts('pages.logs.auditLogsPage.itemLabel')),
    itemLabelPlural: computed(() => ts('pages.logs.auditLogsPage.itemLabelPlural')),

    /**
     * Filter translations
     */
    filters: {
      label: computed(() => ts('pages.logs.auditLogsPage.filters.label')),
      searchPlaceholder: computed(() => ts('pages.logs.auditLogsPage.filters.searchPlaceholder')),
      advancedFilters: computed(() => ts('pages.logs.auditLogsPage.filters.advancedFilters')),
      pendingFilters: computed(() => ts('pages.logs.auditLogsPage.filters.pendingFilters')),
      clearAll: computed(() => ts('pages.logs.auditLogsPage.filters.clearAll')),
      allStatus: computed(() => ts('pages.logs.auditLogsPage.filters.allStatus')),
      status: computed(() => ts('pages.logs.auditLogsPage.filters.status')),
      actor: computed(() => ts('pages.logs.auditLogsPage.filters.actor')),
      action: computed(() => ts('pages.logs.auditLogsPage.filters.action')),
      resourceType: computed(() => ts('pages.logs.auditLogsPage.filters.resourceType')),
      includeChildren: computed(() => ts('pages.logs.auditLogsPage.filters.includeChildren')),
      includeChildrenOrgs: computed(() => ts('pages.logs.auditLogsPage.filters.includeChildrenOrgs')),
      dateRange: computed(() => ts('pages.logs.auditLogsPage.filters.dateRange')),
      options: {
        success: computed(() => ts('pages.logs.auditLogsPage.filters.options.success')),
        failure: computed(() => ts('pages.logs.auditLogsPage.filters.options.failure')),
        yes: computed(() => ts('pages.logs.auditLogsPage.filters.options.yes')),
        no: computed(() => ts('pages.logs.auditLogsPage.filters.options.no')),
      },
    },

    /**
     * Column translations
     */
    columns: {
      actor: computed(() => ts('pages.logs.auditLogsPage.columns.actor')),
      action: computed(() => ts('pages.logs.auditLogsPage.columns.action')),
      resource: computed(() => ts('pages.logs.auditLogsPage.columns.resource')),
      status: computed(() => ts('pages.logs.auditLogsPage.columns.status')),
      timestamp: computed(() => ts('pages.logs.auditLogsPage.columns.timestamp')),
    },

    /**
     * Menu column labels
     */
    menuColumns: {
      action: computed(() => ts('pages.logs.auditLogsPage.menuColumns.action')),
      resource: computed(() => ts('pages.logs.auditLogsPage.menuColumns.resource')),
      timestamp: computed(() => ts('pages.logs.auditLogsPage.menuColumns.timestamp')),
    },

    /**
     * Status options for filter
     */
    statusOptions: {
      success: computed(() => ts('pages.logs.auditLogsPage.statusOptions.success')),
      failure: computed(() => ts('pages.logs.auditLogsPage.statusOptions.failure')),
    },

    /**
     * Action options for filter
     */
    actionOptions: {
      create: computed(() => ts('pages.logs.auditLogsPage.actionOptions.create')),
      edit: computed(() => ts('pages.logs.auditLogsPage.actionOptions.edit')),
      delete: computed(() => ts('pages.logs.auditLogsPage.actionOptions.delete')),
    },

    /**
     * Resource type options for filter
     */
    resourceTypeOptions: {
      userLog: computed(() => ts('pages.logs.auditLogsPage.resourceTypeOptions.userLog')),
      dataSource: computed(() => ts('pages.logs.auditLogsPage.resourceTypeOptions.dataSource')),
      assets: computed(() => ts('pages.logs.auditLogsPage.resourceTypeOptions.assets')),
      payloadHandler: computed(() => ts('pages.logs.auditLogsPage.resourceTypeOptions.payloadHandler')),
      triggers: computed(() => ts('pages.logs.auditLogsPage.resourceTypeOptions.triggers')),
      users: computed(() => ts('pages.logs.auditLogsPage.resourceTypeOptions.users')),
      customers: computed(() => ts('pages.logs.auditLogsPage.resourceTypeOptions.customers')),
    },

    /**
     * Drawer translations
     */
    drawer: {
      title: computed(() => tsTitle('pages.logs.auditLogsPage.drawer.title')),
    },

    /**
     * Empty state translations
     */
    empty: {
      title: computed(() => tsTitle('pages.logs.auditLogsPage.empty.title')),
      description: computed(() => tsRaw('pages.logs.auditLogsPage.empty.description')),
    },

    /**
     * Pagination translations
     */
    pagination: {
      newer: computed(() => ts('pages.logs.auditLogsPage.pagination.newer')),
      older: computed(() => ts('pages.logs.auditLogsPage.pagination.older')),
    },

    /**
     * Message translations
     */
    messages: {
      loadFailed: computed(() => tsRaw('pages.logs.auditLogsPage.messages.loadFailed')),
    },

    /**
     * Default value translations
     */
    defaults: {
      unknown: computed(() => tsRaw('pages.logs.auditLogsPage.defaults.unknown')),
      notAvailable: computed(() => tsRaw('pages.logs.auditLogsPage.defaults.notAvailable')),
    },

    /**
     * DataRow column definitions with reactive translations
     */
    dataRowColumns: computed((): DataRowColumn[] => {
      const getTypeIcon = (type: string): string => {
        const icons: Record<string, string> = {
          userLog: 'person',
          dataSource: 'storage',
          assets: 'devices',
          payloadHandler: 'code',
          triggers: 'flash_on',
          users: 'group',
          customers: 'business',
        };
        return icons[type] || 'folder';
      };

      const getActionColor = (action: string): string => {
        const colors: Record<string, string> = {
          Create: 'green-6',
          Edit: 'blue-6',
          Delete: 'red-6',
        };
        return colors[action] || 'grey-6';
      };

      return [
        {
          key: 'icon',
          label: '',
          type: 'avatar',
          visible: 'always',
          width: 56,
          icon: (_value: any, row: any) => getTypeIcon(row.type),
          color: (_value: any, row: any) => row.status === 'success' ? 'primary' : 'red-5',
        },
        {
          key: 'actor',
          label: ts('pages.logs.auditLogsPage.columns.actor'),
          type: 'text',
          visible: 'always',
          width: 200,
          ellipsis: true,
          secondaryKey: 'details',
        },
        {
          key: 'action',
          label: ts('pages.logs.auditLogsPage.columns.action'),
          type: 'chip',
          visible: 'laptop',
          width: 100,
          color: (value: any) => getActionColor(value),
        },
        {
          key: 'resource',
          label: ts('pages.logs.auditLogsPage.columns.resource'),
          type: 'text',
          visible: 'laptop',
          width: 180,
          ellipsis: true,
        },
        {
          key: 'status',
          label: ts('pages.logs.auditLogsPage.columns.status'),
          type: 'badge',
          visible: 'laptop',
          width: 110,
          format: (value: any) => value ? value.toUpperCase() : 'UNKNOWN',
          color: (value: any) => value === 'success' ? 'green-6' : 'red-6',
        },
        {
          key: 'timestamp',
          label: ts('pages.logs.auditLogsPage.columns.timestamp'),
          type: 'text',
          visible: 'laptop',
          width: 180,
          format: (value: any) => {
            if (!value) return 'N/A';
            const date = new Date(value);
            return date.toLocaleDateString('en-US', {
              day: '2-digit',
              month: 'short',
              year: 'numeric',
              hour: '2-digit',
              minute: '2-digit',
            });
          },
        },
      ];
    }),
  };
}
