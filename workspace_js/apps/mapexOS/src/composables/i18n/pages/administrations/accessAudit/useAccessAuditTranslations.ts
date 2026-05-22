import { computed } from 'vue';
import { useTS } from '@utils/translation';
import type { DataRowColumn } from '@components/cards';
import type { PageHeaderInfo } from '@components/headers';

/**
 * Access Audit page translations
 *
 * Structure mirrors:
 * - File: src/pages/administrations/accessAudit/accessAuditPage/AccessAuditPage.vue
 * - JSON: src/i18n/{locale}/pages/administrations/accessAudit.json
 * - Composable: src/composables/i18n/pages/administrations/accessAudit/useAccessAuditTranslations.ts
 *
 * Provides all translations for the Access Audit page including:
 * - Page header (title, description)
 * - Filter items (labels, options)
 * - DataRow column definitions (reactive)
 * - Menu column labels
 * - Empty state
 */
export function useAccessAuditTranslations() {
  const ts = useTS({ capitalize: true });
  const tsTitle = useTS({ titleCase: true });

  return {
    /**
     * Page header translations
     * Mirrors: pages.administrations.accessAudit
     */
    page: {
      title: computed(() => tsTitle('pages.administrations.accessAudit.title')),
      description: computed(() => ts('pages.administrations.accessAudit.description')),
      listTitle: computed(() => tsTitle('pages.administrations.accessAudit.listTitle')),
      itemLabel: computed(() => ts('pages.administrations.accessAudit.itemLabel')),
      itemLabelPlural: computed(() => ts('pages.administrations.accessAudit.itemLabelPlural')),
      info: computed((): PageHeaderInfo => ({
        title: ts('pages.administrations.accessAudit.info.title'),
        description: ts('pages.administrations.accessAudit.info.description'),
        items: [
          {
            icon: 'security',
            color: 'blue-6',
            title: ts('pages.administrations.accessAudit.info.items.memberships.title'),
            text: ts('pages.administrations.accessAudit.info.items.memberships.text'),
          },
          {
            icon: 'person',
            color: 'green-6',
            title: ts('pages.administrations.accessAudit.info.items.users.title'),
            text: ts('pages.administrations.accessAudit.info.items.users.text'),
          },
          {
            icon: 'group',
            color: 'orange-6',
            title: ts('pages.administrations.accessAudit.info.items.groups.title'),
            text: ts('pages.administrations.accessAudit.info.items.groups.text'),
          },
          {
            icon: 'domain',
            color: 'purple-6',
            title: ts('pages.administrations.accessAudit.info.items.scope.title'),
            text: ts('pages.administrations.accessAudit.info.items.scope.text'),
          },
        ],
        docsUrl: 'https://docs.mapexos.com/administration/access-audit',
        docsLabel: ts('pages.administrations.accessAudit.info.docsLabel'),
      })),
    },

    /**
     * Filter translations with reactive options
     * Mirrors: pages.administrations.accessAudit.filters
     */
    filters: {
      label: computed(() => ts('pages.administrations.accessAudit.filters.label')),
      searchPlaceholder: computed(() => ts('pages.administrations.accessAudit.filters.searchPlaceholder')),
      advancedFilters: computed(() => ts('pages.administrations.accessAudit.filters.advancedFilters')),
      pendingFilters: computed(() => ts('pages.administrations.accessAudit.filters.pendingFilters')),
      clearAll: computed(() => ts('pages.administrations.accessAudit.filters.clearAll')),
      allStatus: computed(() => ts('pages.administrations.accessAudit.filters.allStatus')),
      filterByAssignee: computed(() => ts('pages.administrations.accessAudit.filters.filterByAssignee')),
      assigneeType: computed(() => ts('pages.administrations.accessAudit.filters.assigneeType')),
      assigneeId: computed(() => ts('pages.administrations.accessAudit.filters.assigneeId')),
      roleId: computed(() => ts('pages.administrations.accessAudit.filters.roleId')),
      scope: computed(() => ts('pages.administrations.accessAudit.filters.scope')),
      status: computed(() => ts('pages.administrations.accessAudit.filters.status')),
      includeChildren: computed(() => ts('pages.administrations.accessAudit.filters.includeChildren')),

      options: {
        all: computed(() => ts('pages.administrations.accessAudit.filters.options.all')),
        user: computed(() => ts('pages.administrations.accessAudit.filters.options.user')),
        group: computed(() => ts('pages.administrations.accessAudit.filters.options.group')),
        local: computed(() => ts('pages.administrations.accessAudit.filters.options.local')),
        recursive: computed(() => ts('pages.administrations.accessAudit.filters.options.recursive')),
        enabled: computed(() => ts('pages.administrations.accessAudit.filters.options.enabled')),
        disabled: computed(() => ts('pages.administrations.accessAudit.filters.options.disabled')),
        yes: computed(() => ts('pages.administrations.accessAudit.filters.options.yes')),
        no: computed(() => ts('pages.administrations.accessAudit.filters.options.no')),
      },
    },

    /**
     * Menu column labels (for ListHeaderMenu)
     * Mirrors: pages.administrations.accessAudit.menuColumns
     */
    menuColumns: {
      organization: computed(() => ts('pages.administrations.accessAudit.menuColumns.organization')),
      roles: computed(() => ts('pages.administrations.accessAudit.menuColumns.roles')),
      scope: computed(() => ts('pages.administrations.accessAudit.menuColumns.scope')),
    },

    /**
     * Empty state translations
     * Mirrors: pages.administrations.accessAudit.empty
     */
    empty: {
      title: computed(() => ts('pages.administrations.accessAudit.empty.title')),
      description: computed(() => ts('pages.administrations.accessAudit.empty.description')),
    },

    /**
     * Scope labels
     * Mirrors: pages.administrations.accessAudit.scope
     */
    scope: {
      local: computed(() => ts('pages.administrations.accessAudit.scope.local')),
      recursive: computed(() => ts('pages.administrations.accessAudit.scope.recursive')),
    },

    /**
     * Assignee type labels
     * Mirrors: pages.administrations.accessAudit.assigneeType
     */
    assigneeType: {
      user: computed(() => ts('pages.administrations.accessAudit.assigneeType.user')),
      group: computed(() => ts('pages.administrations.accessAudit.assigneeType.group')),
    },

    /**
     * DataRow column definitions with reactive translations
     * IMPORTANT: Returns a computed ref so columns update when language changes
     *
     * Usage: <DataRow :columns="columns.value" />
     */
    columns: computed(() => {
      return [
        {
          key: 'avatar',
          label: '',
          type: 'avatar',
          visible: 'always',
          width: 56,
          icon: (value: any, row: any) => row.assigneeType === 'user' ? 'person' : 'group',
          color: (value: any, row: any) => row.enabled ? 'primary' : 'grey-5',
          tooltip: (value: any, row: any) =>
            row.assigneeType === 'user'
              ? ts('pages.administrations.accessAudit.assigneeType.user')
              : ts('pages.administrations.accessAudit.assigneeType.group'),
        },
        {
          key: 'assigneeName',
          label: ts('pages.administrations.accessAudit.columns.assignee'),
          type: 'text',
          visible: 'always',
          width: 200,
          ellipsis: true,
          secondaryKey: 'assigneeType',
          secondaryFormat: (value: any) =>
            value === 'user'
              ? ts('pages.administrations.accessAudit.assigneeType.user')
              : ts('pages.administrations.accessAudit.assigneeType.group'),
        },
        {
          key: 'orgName',
          label: ts('pages.administrations.accessAudit.columns.organization'),
          type: 'chip',
          visible: 'laptop',
          width: 180,
          ellipsis: true,
          color: 'indigo-6',
          icon: 'domain',
        },
        {
          key: 'roleNames',
          label: ts('pages.administrations.accessAudit.columns.roles'),
          type: 'chips',
          visible: 'laptop',
          width: 200,
          color: 'teal-6',
          icon: 'badge',
        },
        {
          key: 'scope',
          label: ts('pages.administrations.accessAudit.columns.scope'),
          type: 'chip',
          visible: 'always',
          width: 120,
          format: (value: any) =>
            value === 'local'
              ? ts('pages.administrations.accessAudit.scope.local')
              : ts('pages.administrations.accessAudit.scope.recursive'),
          color: (value: any) => value === 'recursive' ? 'orange-6' : 'blue-6',
          icon: (value: any) => value === 'recursive' ? 'account_tree' : 'stop',
        },
      ] as DataRowColumn[];
    }),

    /**
     * Error translations
     */
    errors: {
      assigneeIdMissing: computed(() => ts('pages.administrations.accessAudit.errors.assigneeIdMissing')),
    },

  };
}
