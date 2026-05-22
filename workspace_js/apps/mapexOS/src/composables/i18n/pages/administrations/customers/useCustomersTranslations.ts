import { computed } from 'vue';
import { useTS } from '@utils/translation';
import type { DataRowColumn } from '@components/cards';
import type { FilterListItem } from '@components/filters';
import type { PageHeaderInfo } from '@components/headers';

/**
 * Customers list page translations
 *
 * Structure mirrors:
 * - File: src/pages/administrations/customers/customerListPage/CustomersListPage.vue
 * - JSON: src/i18n/{locale}/pages/administrations/customers.json
 * - Composable: src/composables/i18n/pages/administrations/customers/useCustomersTranslations.ts
 */
export function useCustomersTranslations() {
  const ts = useTS({ capitalize: true });
  const tsTitle = useTS({ titleCase: true });

  return {
    page: {
      title: computed(() => tsTitle('pages.administrations.customers.title')),
      description: computed(() => ts('pages.administrations.customers.description')),
      addButton: computed(() => ts('pages.administrations.customers.addButton')),
      listTitle: computed(() => tsTitle('pages.administrations.customers.listTitle')),
      itemLabel: computed(() => ts('pages.administrations.customers.itemLabel')),
      itemLabelPlural: computed(() => ts('pages.administrations.customers.itemLabelPlural')),
      info: computed((): PageHeaderInfo => ({
        title: ts('pages.administrations.customers.info.title'),
        description: ts('pages.administrations.customers.info.description'),
        items: [
          {
            icon: 'business',
            color: 'blue-6',
            title: ts('pages.administrations.customers.info.items.multiSite.title'),
            text: ts('pages.administrations.customers.info.items.multiSite.text'),
          },
          {
            icon: 'location_city',
            color: 'green-6',
            title: ts('pages.administrations.customers.info.items.hierarchy.title'),
            text: ts('pages.administrations.customers.info.items.hierarchy.text'),
          },
          {
            icon: 'settings',
            color: 'orange-6',
            title: ts('pages.administrations.customers.info.items.configuration.title'),
            text: ts('pages.administrations.customers.info.items.configuration.text'),
          },
          {
            icon: 'inventory',
            color: 'purple-6',
            title: ts('pages.administrations.customers.info.items.assetOwnership.title'),
            text: ts('pages.administrations.customers.info.items.assetOwnership.text'),
          },
          {
            icon: 'group',
            color: 'indigo-6',
            title: ts('pages.administrations.customers.info.items.userAccess.title'),
            text: ts('pages.administrations.customers.info.items.userAccess.text'),
          },
        ],
        docsUrl: 'https://docs.mapexos.com/organizations/customers',
        docsLabel: ts('pages.administrations.customers.info.docsLabel'),
      })),
    },

    filters: {
      label: computed(() => ts('pages.administrations.customers.filters.label')),
      searchPlaceholder: computed(() => ts('pages.administrations.customers.filters.searchPlaceholder')),
      allStatus: computed(() => ts('pages.administrations.customers.filters.allStatus')),
      advancedFilters: computed(() => ts('pages.administrations.customers.filters.advancedFilters')),
      pendingFilters: computed(() => ts('pages.administrations.customers.filters.pendingFilters')),
      clearAll: computed(() => ts('pages.administrations.customers.filters.clearAll')),
      organizationName: computed(() => ts('pages.administrations.customers.filters.organizationName')),
      status: computed(() => ts('pages.administrations.customers.filters.status')),
      includeChildren: computed(() => ts('pages.administrations.customers.filters.includeChildren')),
      includeChildrenOrgs: computed(() => ts('pages.administrations.customers.filters.includeChildrenOrgs')),
      type: computed(() => ts('pages.administrations.customers.filters.type')),

      options: {
        yes: computed(() => ts('pages.administrations.customers.filters.options.yes')),
        no: computed(() => ts('pages.administrations.customers.filters.options.no')),
        enabled: computed(() => ts('pages.administrations.customers.filters.options.enabled')),
        disabled: computed(() => ts('pages.administrations.customers.filters.options.disabled')),
        customer: computed(() => ts('pages.administrations.customers.filters.options.customer')),
        site: computed(() => ts('pages.administrations.customers.filters.options.site')),
        building: computed(() => ts('pages.administrations.customers.filters.options.building')),
        floor: computed(() => ts('pages.administrations.customers.filters.options.floor')),
        zone: computed(() => ts('pages.administrations.customers.filters.options.zone')),
      },
    },

    menuColumns: {
      organization: computed(() => ts('pages.administrations.customers.menuColumns.organization')),
      address: computed(() => ts('pages.administrations.customers.menuColumns.address')),
      created: computed(() => ts('pages.administrations.customers.menuColumns.created')),
    },

    empty: {
      title: computed(() => ts('pages.administrations.customers.empty.title')),
      description: computed(() => ts('pages.administrations.customers.empty.description')),
    },

    status: {
      active: computed(() => ts('pages.administrations.customers.status.active')),
      inactive: computed(() => ts('pages.administrations.customers.status.inactive')),
    },

    messages: {
      deletedSuccessfully: computed(() => ts('pages.administrations.customers.messages.deletedSuccessfully')),
      confirmDelete: (name: string) => ts('pages.administrations.customers.messages.confirmDelete', { name }),
    },

    dialog: {
      deleteTitle: computed(() => ts('pages.administrations.customers.dialog.deleteTitle')),
    },

    errors: {
      apiNotInitialized: computed(() => ts('pages.administrations.customers.errors.apiNotInitialized')),
      idMissing: computed(() => ts('pages.administrations.customers.errors.idMissing')),
    },

    drawer: {
      title: computed(() => ts('pages.administrations.customers.drawer.title')),
      close: computed(() => ts('pages.administrations.customers.drawer.close')),
      edit: computed(() => ts('pages.administrations.customers.drawer.edit')),
      loading: computed(() => ts('pages.administrations.customers.drawer.loading')),
      error: computed(() => ts('pages.administrations.customers.drawer.error')),

      sections: {
        basicInfo: computed(() => ts('pages.administrations.customers.drawer.sections.basicInfo')),
        hierarchy: computed(() => ts('pages.administrations.customers.drawer.sections.hierarchy')),
        address: computed(() => ts('pages.administrations.customers.drawer.sections.address')),
        contact: computed(() => ts('pages.administrations.customers.drawer.sections.contact')),
        timestamps: computed(() => ts('pages.administrations.customers.drawer.sections.timestamps')),
      },

      fields: {
        name: computed(() => ts('pages.administrations.customers.drawer.fields.name')),
        type: computed(() => ts('pages.administrations.customers.drawer.fields.type')),
        status: computed(() => ts('pages.administrations.customers.drawer.fields.status')),
        logo: computed(() => ts('pages.administrations.customers.drawer.fields.logo')),
        parentOrganization: computed(() => ts('pages.administrations.customers.drawer.fields.parentOrganization')),
        customerOf: computed(() => ts('pages.administrations.customers.drawer.fields.customerOf')),
        code: computed(() => ts('pages.administrations.customers.drawer.fields.code')),
        pathKey: computed(() => ts('pages.administrations.customers.drawer.fields.pathKey')),
        depth: computed(() => ts('pages.administrations.customers.drawer.fields.depth')),
        childCount: computed(() => ts('pages.administrations.customers.drawer.fields.childCount')),
        street: computed(() => ts('pages.administrations.customers.drawer.fields.street')),
        city: computed(() => ts('pages.administrations.customers.drawer.fields.city')),
        state: computed(() => ts('pages.administrations.customers.drawer.fields.state')),
        country: computed(() => ts('pages.administrations.customers.drawer.fields.country')),
        zipCode: computed(() => ts('pages.administrations.customers.drawer.fields.zipCode')),
        phone: computed(() => ts('pages.administrations.customers.drawer.fields.phone')),
        created: computed(() => ts('pages.administrations.customers.drawer.fields.created')),
        updated: computed(() => ts('pages.administrations.customers.drawer.fields.updated')),
      },

      empty: {
        address: computed(() => ts('pages.administrations.customers.drawer.empty.address')),
        phone: computed(() => ts('pages.administrations.customers.drawer.empty.phone')),
      },

      type: {
        customer: computed(() => ts('pages.administrations.customers.drawer.type.customer')),
        site: computed(() => ts('pages.administrations.customers.drawer.type.site')),
        building: computed(() => ts('pages.administrations.customers.drawer.type.building')),
        floor: computed(() => ts('pages.administrations.customers.drawer.type.floor')),
        zone: computed(() => ts('pages.administrations.customers.drawer.type.zone')),
      },
    },

    columns: computed(() => {
      // Helper function to get icon based on organization type
      const getTypeIcon = (type: string): string => {
        const iconMap: Record<string, string> = {
          customer: 'domain',
          site: 'location_on',
          building: 'apartment',
          floor: 'layers',
          zone: 'place',
        };
        return iconMap[type] || 'domain';
      };

      // Helper function to get color based on status
      const getStatusColor = (enabled: boolean): string => {
        return enabled ? 'primary' : 'grey-5';
      };

      // Helper function to get status tooltip
      const getStatusTooltip = (enabled: boolean): string => {
        return enabled
          ? ts('pages.administrations.customers.status.active')
          : ts('pages.administrations.customers.status.inactive');
      };

      return [
        {
          key: 'icon',
          label: '',
          type: 'avatar',
          visible: 'always',
          width: 56,
          icon: (value: any, row: any) => row.logo || getTypeIcon(row.type),
          color: (value: any, row: any) => getStatusColor(row.enabled),
          tooltip: (value: any, row: any) => getStatusTooltip(row.enabled),
        },
        {
          key: 'name',
          label: ts('pages.administrations.customers.columns.name'),
          type: 'text',
          visible: 'always',
          width: 250,
          ellipsis: true,
          secondaryKey: 'phone',
        },
        {
          key: 'organizationName',
          label: ts('pages.administrations.customers.columns.organization'),
          type: 'chip',
          visible: 'always',
          width: 180,
          ellipsis: true,
          color: 'indigo-6',
          icon: 'domain',
        },
        {
          key: 'address.city',
          label: ts('pages.administrations.customers.columns.address'),
          type: 'text',
          visible: 'laptop',
          width: 200,
          ellipsis: true,
          format: (value: any, row: any) => row.address?.city || 'N/A',
          secondary: (value: any, row: any) => {
            if (!row.address) return 'N/A, N/A';
            const state = row.address.state || 'N/A';
            const country = row.address.country || 'N/A';
            return `${state}, ${country}`;
          },
        },
        {
          key: 'created',
          label: ts('pages.administrations.customers.columns.created'),
          type: 'text',
          visible: 'laptop',
          width: 120,
          format: (value: any) => {
            if (!value) return 'N/A';
            return new Date(value).toLocaleDateString();
          },
        },
      ] as DataRowColumn[];
    }),

    filterItems: computed(() => {
      return [
        // Row 1: Standard pattern (6 + 3 + 3 = 12 cols)
        {
          key: 'name',
          type: 'input',
          label: ts('pages.administrations.customers.filters.organizationName'),
          icon: 'search',
          grid: 'col-12 col-md-6'
        },
        {
          key: 'enabled',
          type: 'select',
          label: ts('pages.administrations.customers.filters.status'),
          icon: 'toggle_on',
          options: [
            { label: ts('pages.administrations.customers.filters.options.all'), value: null },
            { label: ts('pages.administrations.customers.status.active'), value: true },
            { label: ts('pages.administrations.customers.status.inactive'), value: false }
          ],
          grid: 'col-6 col-md-3'
        },
        {
          key: 'includeChildren',
          type: 'select',
          label: ts('pages.administrations.customers.filters.includeChildren'),
          icon: 'account_tree',
          options: [
            { label: ts('pages.administrations.customers.filters.options.yes'), value: true },
            { label: ts('pages.administrations.customers.filters.options.no'), value: false }
          ],
          grid: 'col-6 col-md-3',
          disabled: true, // Always enabled, user cannot change
          defaultValue: true, // Default to "Yes"
        },
        // Row 2: Domain-specific filter (12 cols)
        {
          key: 'type',
          type: 'select',
          label: ts('pages.administrations.customers.filters.type'),
          icon: 'category',
          options: [
            { label: ts('pages.administrations.customers.filters.options.all'), value: null },
            { label: ts('pages.administrations.customers.filters.options.customer'), value: 'customer', icon: 'domain', color: 'green-6' },
            { label: ts('pages.administrations.customers.filters.options.site'), value: 'site', icon: 'location_on', color: 'orange-6' },
            { label: ts('pages.administrations.customers.filters.options.building'), value: 'building', icon: 'apartment', color: 'blue-6' },
            { label: ts('pages.administrations.customers.filters.options.floor'), value: 'floor', icon: 'layers', color: 'teal-6' },
            { label: ts('pages.administrations.customers.filters.options.zone'), value: 'zone', icon: 'place', color: 'green-7' },
          ],
          grid: 'col-12 col-md-6'
        },
      ] as FilterListItem[];
    }),
  };
}
