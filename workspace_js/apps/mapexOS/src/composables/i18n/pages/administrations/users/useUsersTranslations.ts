import { computed } from 'vue';
import { useTS } from '@utils/translation';
import type { DataRowColumn } from '@components/cards';
import type { FilterListItem } from '@components/filters';
import type { PageHeaderInfo } from '@components/headers';

/**
 * Users list page translations
 *
 * Structure mirrors:
 * - File: src/pages/administrations/users/userListPage/UserListPage.vue
 * - JSON: src/i18n/{locale}/pages/administrations/users.json
 * - Composable: src/composables/i18n/pages/administrations/users/useUsersTranslations.ts
 *
 * Provides all translations for the Users List page including:
 * - Page header (title, description, button)
 * - Filter items (labels, options)
 * - DataRow column definitions (reactive)
 * - Menu column labels
 * - Empty state
 * - Success/error messages
 * - Dialog content
 */
export function useUsersTranslations() {
  const ts = useTS({ capitalize: true });
  const tsTitle = useTS({ titleCase: true });

  return {
    /**
     * Page header translations
     * Mirrors: pages.administrations.users
     */
    page: {
      title: computed(() => tsTitle('pages.administrations.users.title')),
      description: computed(() => ts('pages.administrations.users.description')),
      addButton: computed(() => ts('pages.administrations.users.addButton')),
      listTitle: computed(() => tsTitle('pages.administrations.users.listTitle')),
      itemLabel: computed(() => ts('pages.administrations.users.itemLabel')),
      itemLabelPlural: computed(() => ts('pages.administrations.users.itemLabelPlural')),
      info: computed((): PageHeaderInfo => ({
        title: ts('pages.administrations.users.info.title'),
        description: ts('pages.administrations.users.info.description'),
        items: [
          {
            icon: 'person',
            color: 'blue-6',
            title: ts('pages.administrations.users.info.items.authentication.title'),
            text: ts('pages.administrations.users.info.items.authentication.text'),
          },
          {
            icon: 'security',
            color: 'green-6',
            title: ts('pages.administrations.users.info.items.rolePermissions.title'),
            text: ts('pages.administrations.users.info.items.rolePermissions.text'),
          },
          {
            icon: 'badge',
            color: 'orange-6',
            title: ts('pages.administrations.users.info.items.multiOrg.title'),
            text: ts('pages.administrations.users.info.items.multiOrg.text'),
          },
          {
            icon: 'settings',
            color: 'purple-6',
            title: ts('pages.administrations.users.info.items.profile.title'),
            text: ts('pages.administrations.users.info.items.profile.text'),
          },
          {
            icon: 'history',
            color: 'indigo-6',
            title: ts('pages.administrations.users.info.items.activity.title'),
            text: ts('pages.administrations.users.info.items.activity.text'),
          },
        ],
        docsUrl: 'https://docs.mapexos.com/administration/users',
        docsLabel: ts('pages.administrations.users.info.docsLabel'),
      })),
    },

    /**
     * Filter translations with reactive options
     * Mirrors: pages.administrations.users.filters
     */
    filters: {
      label: computed(() => ts('pages.administrations.users.filters.label')),
      searchPlaceholder: computed(() => ts('pages.administrations.users.filters.searchPlaceholder')),
      advancedFilters: computed(() => ts('pages.administrations.users.filters.advancedFilters')),
      pendingFilters: computed(() => ts('pages.administrations.users.filters.pendingFilters')),
      clearAll: computed(() => ts('pages.administrations.users.filters.clearAll')),
      allStatus: computed(() => ts('pages.administrations.users.filters.allStatus')),
      email: computed(() => ts('pages.administrations.users.filters.email')),
      firstName: computed(() => ts('pages.administrations.users.filters.firstName')),
      lastName: computed(() => ts('pages.administrations.users.filters.lastName')),
      status: computed(() => ts('pages.administrations.users.filters.status')),
      includeChildren: computed(() => ts('pages.administrations.users.filters.includeChildren')),
      includeChildrenOrgs: computed(() => ts('pages.administrations.users.filters.includeChildrenOrgs')),
      filterByEmail: computed(() => ts('pages.administrations.users.filters.filterByEmail')),
      filterByFirstName: computed(() => ts('pages.administrations.users.filters.filterByFirstName')),
      filterByLastName: computed(() => ts('pages.administrations.users.filters.filterByLastName')),

      options: {
        enabled: computed(() => ts('pages.administrations.users.filters.options.enabled')),
        disabled: computed(() => ts('pages.administrations.users.filters.options.disabled')),
        yes: computed(() => ts('pages.administrations.users.filters.options.yes')),
        no: computed(() => ts('pages.administrations.users.filters.options.no')),
      },
    },

    /**
     * Menu column labels (for ListHeaderMenu)
     * Mirrors: pages.administrations.users.menuColumns
     */
    menuColumns: {
      organization: computed(() => ts('pages.administrations.users.menuColumns.organization')),
      email: computed(() => ts('pages.administrations.users.menuColumns.email')),
      jobTitle: computed(() => ts('pages.administrations.users.menuColumns.jobTitle')),
      groups: computed(() => ts('pages.administrations.users.menuColumns.groups')),
    },

    /**
     * Empty state translations
     * Mirrors: pages.administrations.users.empty
     */
    empty: {
      title: computed(() => ts('pages.administrations.users.empty.title')),
      description: computed(() => ts('pages.administrations.users.empty.description')),
    },

    /**
     * Status labels
     * Mirrors: pages.administrations.users.status
     */
    status: {
      active: computed(() => ts('pages.administrations.users.status.active')),
      inactive: computed(() => ts('pages.administrations.users.status.inactive')),
    },

    /**
     * Message translations
     * Mirrors: pages.administrations.users.messages
     */
    messages: {
      deletedSuccessfully: computed(() => ts('pages.administrations.users.messages.deletedSuccessfully')),
      confirmDelete: (name: string) => ts('pages.administrations.users.messages.confirmDelete', { name }),
    },

    /**
     * Dialog translations
     * Mirrors: pages.administrations.users.dialog
     */
    dialog: {
      deleteTitle: computed(() => ts('pages.administrations.users.dialog.deleteTitle')),
    },

    /**
     * Error translations
     * Mirrors: pages.administrations.users.errors
     */
    errors: {
      apiNotInitialized: computed(() => ts('pages.administrations.users.errors.apiNotInitialized')),
      idMissing: computed(() => ts('pages.administrations.users.errors.idMissing')),
    },

    /**
     * Tour translations
     * Mirrors: pages.administrations.users.tour
     */
    tour: {
      tourButton: {
        title: computed(() => ts('pages.administrations.users.tour.tourButton.title')),
        description: computed(() => ts('pages.administrations.users.tour.tourButton.description')),
      },
      header: {
        title: computed(() => ts('pages.administrations.users.tour.header.title')),
        description: computed(() => ts('pages.administrations.users.tour.header.description')),
      },
      filters: {
        title: computed(() => ts('pages.administrations.users.tour.filters.title')),
        description: computed(() => ts('pages.administrations.users.tour.filters.description')),
      },
      advancedFiltersBtn: {
        title: computed(() => ts('pages.administrations.users.tour.advancedFiltersBtn.title')),
        description: computed(() => ts('pages.administrations.users.tour.advancedFiltersBtn.description')),
      },
      advancedFiltersOpen: {
        title: computed(() => ts('pages.administrations.users.tour.advancedFiltersOpen.title')),
        description: computed(() => ts('pages.administrations.users.tour.advancedFiltersOpen.description')),
      },
      results: {
        title: computed(() => ts('pages.administrations.users.tour.results.title')),
        description: computed(() => ts('pages.administrations.users.tour.results.description')),
      },
      rowActions: {
        title: computed(() => ts('pages.administrations.users.tour.rowActions.title')),
        description: computed(() => ts('pages.administrations.users.tour.rowActions.description')),
      },
      addNew: {
        title: computed(() => ts('pages.administrations.users.tour.addNew.title')),
        description: computed(() => ts('pages.administrations.users.tour.addNew.description')),
      },
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
          icon: (value: any, row: any) => row.avatar || 'person',
          color: (value: any, row: any) => row.enabled ? 'primary' : 'grey-5',
          tooltip: (value: any, row: any) =>
            row.enabled
              ? ts('pages.administrations.users.status.active')
              : ts('pages.administrations.users.status.inactive'),
        },
        {
          key: 'name',
          label: ts('pages.administrations.users.columns.name'),
          type: 'text',
          visible: 'always',
          width: 250,
          ellipsis: true,
          format: (value: any, row: any) => `${row.firstName || ''} ${row.lastName || ''}`.trim() || row.email || 'Unknown',
          secondaryKey: 'jobTitle',
        },
        {
          key: 'organizationName',
          label: ts('pages.administrations.users.columns.organization'),
          type: 'chip',
          visible: 'laptop',
          width: 180,
          ellipsis: true,
          color: 'indigo-6',
          icon: 'domain',
        },
        {
          key: 'email',
          label: ts('pages.administrations.users.columns.email'),
          type: 'text',
          visible: 'laptop',
          width: 200,
          ellipsis: true,
        },
        {
          key: 'jobTitle',
          label: ts('pages.administrations.users.columns.jobTitle'),
          type: 'chip',
          visible: 'always',
          width: 150,
          format: (value: any) => value || 'N/A',
          color: 'blue-6',
        },
        {
          key: 'groupsCount',
          label: ts('pages.administrations.users.columns.groups'),
          type: 'chip',
          visible: 'laptop',
          width: 100,
          format: (value: any) => value ?? 0,
          color: 'teal-6',
          icon: 'groups',
        },
      ] as DataRowColumn[];
    }),

    /**
     * Filter items with reactive translations
     * Returns computed FilterListItem array for ListFilter component
     *
     * Usage: <ListFilter :items="filterItems.value" />
     */
    filterItems: computed(() => {
      return [
        // Row 1: Standard pattern (6 + 3 + 3 = 12 cols)
        {
          key: 'email',
          type: 'input',
          label: ts('pages.administrations.users.filters.email'),
          icon: 'search',
          grid: 'col-12 col-md-6'
        },
        {
          key: 'enabled',
          type: 'select',
          label: ts('pages.administrations.users.filters.status'),
          icon: 'toggle_on',
          options: [
            { label: ts('pages.administrations.users.filters.options.all'), value: null },
            { label: ts('pages.administrations.users.filters.options.enabled'), value: true },
            { label: ts('pages.administrations.users.filters.options.disabled'), value: false }
          ],
          grid: 'col-6 col-md-3'
        },
        {
          key: 'includeChildren',
          type: 'select',
          label: ts('pages.administrations.users.filters.includeChildren'),
          icon: 'account_tree',
          options: [
            { label: ts('pages.administrations.users.filters.options.all'), value: null },
            { label: ts('pages.administrations.users.filters.options.yes'), value: true },
            { label: ts('pages.administrations.users.filters.options.no'), value: false }
          ],
          grid: 'col-6 col-md-3'
        },
        // Row 2: Domain-specific filters (6 + 6 = 12 cols)
        {
          key: 'firstName',
          type: 'input',
          label: ts('pages.administrations.users.filters.firstName'),
          icon: 'person',
          grid: 'col-12 col-md-6'
        },
        {
          key: 'lastName',
          type: 'input',
          label: ts('pages.administrations.users.filters.lastName'),
          icon: 'person_outline',
          grid: 'col-12 col-md-6'
        },
      ] as FilterListItem[];
    }),

    /**
     * User Details Drawer translations
     * Mirrors: pages.administrations.users.drawer
     */
    drawer: {
      title: computed(() => ts('pages.administrations.users.drawer.title')),
      close: computed(() => ts('pages.administrations.users.drawer.close')),
      edit: computed(() => ts('pages.administrations.users.drawer.edit')),
      loading: computed(() => ts('pages.administrations.users.drawer.loading')),
      error: computed(() => ts('pages.administrations.users.drawer.error')),

      sections: {
        basicInfo: computed(() => ts('pages.administrations.users.drawer.sections.basicInfo')),
        contact: computed(() => ts('pages.administrations.users.drawer.sections.contact')),
        authentication: computed(() => ts('pages.administrations.users.drawer.sections.authentication')),
        timestamps: computed(() => ts('pages.administrations.users.drawer.sections.timestamps')),
        groups: computed(() => ts('pages.administrations.users.drawer.sections.groups')),
        orgAccess: computed(() => ts('pages.administrations.users.drawer.sections.orgAccess')),
      },

      columns: {
        organization: computed(() => ts('pages.administrations.users.drawer.columns.organization')),
        scope: computed(() => ts('pages.administrations.users.drawer.columns.scope')),
        roles: computed(() => ts('pages.administrations.users.drawer.columns.roles')),
      },

      accessStats: {
        viaGroups: computed(() => ts('pages.administrations.users.drawer.accessStats.viaGroups')),
        direct: computed(() => ts('pages.administrations.users.drawer.accessStats.direct')),
      },

      viaGroups: {
        label: computed(() => ts('pages.administrations.users.drawer.viaGroups.label')),
        seeMore: computed(() => ts('pages.administrations.users.drawer.viaGroups.seeMore')),
        allGroups: computed(() => ts('pages.administrations.users.drawer.viaGroups.allGroups')),
      },

      directAccess: {
        allRoles: computed(() => ts('pages.administrations.users.drawer.directAccess.allRoles')),
      },

      fields: {
        name: computed(() => ts('pages.administrations.users.drawer.fields.name')),
        email: computed(() => ts('pages.administrations.users.drawer.fields.email')),
        status: computed(() => ts('pages.administrations.users.drawer.fields.status')),
        avatar: computed(() => ts('pages.administrations.users.drawer.fields.avatar')),
        phone: computed(() => ts('pages.administrations.users.drawer.fields.phone')),
        jobTitle: computed(() => ts('pages.administrations.users.drawer.fields.jobTitle')),
        authProvider: computed(() => ts('pages.administrations.users.drawer.fields.authProvider')),
        changePasswordNextLogin: computed(() => ts('pages.administrations.users.drawer.fields.changePasswordNextLogin')),
        externalId: computed(() => ts('pages.administrations.users.drawer.fields.externalId')),
        created: computed(() => ts('pages.administrations.users.drawer.fields.created')),
        updated: computed(() => ts('pages.administrations.users.drawer.fields.updated')),
      },

      authProviders: {
        internal: computed(() => ts('pages.administrations.users.drawer.authProviders.internal')),
        google: computed(() => ts('pages.administrations.users.drawer.authProviders.google')),
        github: computed(() => ts('pages.administrations.users.drawer.authProviders.github')),
        microsoft: computed(() => ts('pages.administrations.users.drawer.authProviders.microsoft')),
        keycloak: computed(() => ts('pages.administrations.users.drawer.authProviders.keycloak')),
      },

      empty: {
        phone: computed(() => ts('pages.administrations.users.drawer.empty.phone')),
        jobTitle: computed(() => ts('pages.administrations.users.drawer.empty.jobTitle')),
        externalId: computed(() => ts('pages.administrations.users.drawer.empty.externalId')),
        avatar: computed(() => ts('pages.administrations.users.drawer.empty.avatar')),
      },
    },
  };
}
