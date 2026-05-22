import type { FilterListItem } from '@components/filters';
import type { DataRowColumn } from '@components/cards';
import type { PageHeaderInfo } from '@components/headers';

import { computed } from 'vue';
import { useTS } from '@utils/translation';
import { useOrganizationStore } from '@stores/organization';

export function useAssetTemplatesTranslations() {
  const ts = useTS({ capitalize: true });
  const tsTitle = useTS({ titleCase: true });
  const organizationStore = useOrganizationStore();

  return {
    pageHeader: {
      title: computed(() => tsTitle('pages.assets.assetTemplates.pageHeader.title')),
      description: computed(() => ts('pages.assets.assetTemplates.pageHeader.description')),
      button: computed(() => ts('pages.assets.assetTemplates.pageHeader.button')),
      info: computed((): PageHeaderInfo => ({
        title: ts('pages.assets.assetTemplates.pageHeader.info.title'),
        description: ts('pages.assets.assetTemplates.pageHeader.info.description'),
        items: [
          {
            icon: 'layers',
            color: 'blue-6',
            title: ts('pages.assets.assetTemplates.pageHeader.info.items.reusability.title'),
            text: ts('pages.assets.assetTemplates.pageHeader.info.items.reusability.text'),
          },
          {
            icon: 'settings',
            color: 'green-6',
            title: ts('pages.assets.assetTemplates.pageHeader.info.items.processing.title'),
            text: ts('pages.assets.assetTemplates.pageHeader.info.items.processing.text'),
          },
          {
            icon: 'factory',
            color: 'orange-6',
            title: ts('pages.assets.assetTemplates.pageHeader.info.items.manufacturer.title'),
            text: ts('pages.assets.assetTemplates.pageHeader.info.items.manufacturer.text'),
          },
          {
            icon: 'lock',
            color: 'purple-6',
            title: ts('pages.assets.assetTemplates.pageHeader.info.items.systemTemplates.title'),
            text: ts('pages.assets.assetTemplates.pageHeader.info.items.systemTemplates.text'),
          },
          {
            icon: 'tune',
            color: 'indigo-6',
            title: ts('pages.assets.assetTemplates.pageHeader.info.items.customization.title'),
            text: ts('pages.assets.assetTemplates.pageHeader.info.items.customization.text'),
          },
        ],
        docsUrl: 'https://docs.mapexos.com/asset-templates',
        docsLabel: ts('pages.assets.assetTemplates.pageHeader.info.docsLabel'),
      })),
    },

    /**
     * Quick filter translations for enterprise pattern
     */
    quickFilters: {
      label: computed(() => ts('pages.assets.assetTemplates.filters.label')),
      searchPlaceholder: computed(() => ts('pages.assets.assetTemplates.filters.searchPlaceholder')),
      advancedFilters: computed(() => ts('pages.assets.assetTemplates.filters.advancedFilters')),
      pendingFilters: computed(() => ts('pages.assets.assetTemplates.filters.pendingFilters')),
      clearAll: computed(() => ts('pages.assets.assetTemplates.filters.clearAll')),
      allStatus: computed(() => ts('pages.assets.assetTemplates.filters.allStatus')),
      status: computed(() => ts('pages.assets.assetTemplates.filters.status')),
      name: computed(() => ts('pages.assets.assetTemplates.filters.name')),
      category: computed(() => ts('pages.assets.assetTemplates.filters.category')),
      manufacturer: computed(() => ts('pages.assets.assetTemplates.filters.manufacturer')),
      model: computed(() => ts('pages.assets.assetTemplates.filters.model')),
      isSystem: computed(() => ts('pages.assets.assetTemplates.filters.isSystem')),
      isTemplate: computed(() => ts('pages.assets.assetTemplates.filters.isTemplate')),
      includeChildren: computed(() => ts('pages.assets.assetTemplates.filters.includeChildren')),
      filterByCategory: computed(() => ts('pages.assets.assetTemplates.filters.filterByCategory')),
      filterByManufacturer: computed(() => ts('pages.assets.assetTemplates.filters.filterByManufacturer')),
      filterByModel: computed(() => ts('pages.assets.assetTemplates.filters.filterByModel')),
      options: {
        all: computed(() => ts('pages.assets.assetTemplates.filters.statusOptions.all')),
        enabled: computed(() => ts('pages.assets.assetTemplates.filters.statusOptions.enabled')),
        disabled: computed(() => ts('pages.assets.assetTemplates.filters.statusOptions.disabled')),
        active: computed(() => ts('pages.assets.assetTemplates.filters.statusOptions.active')),
        inactive: computed(() => ts('pages.assets.assetTemplates.filters.statusOptions.inactive')),
        yes: computed(() => ts('pages.assets.assetTemplates.filters.includeChildrenOptions.yes')),
        no: computed(() => ts('pages.assets.assetTemplates.filters.includeChildrenOptions.no')),
        system: computed(() => ts('pages.assets.assetTemplates.filters.isSystemOptions.system')),
        custom: computed(() => ts('pages.assets.assetTemplates.filters.isSystemOptions.custom')),
        shared: computed(() => ts('pages.assets.assetTemplates.filters.isTemplateOptions.templates')),
        local: computed(() => ts('pages.assets.assetTemplates.filters.isTemplateOptions.local')),
      },
    },

    filters: computed((): FilterListItem[] => {
      const baseFilters: FilterListItem[] = [
        // Row 1: Common Filters (Standard Pattern)
        {
          key: 'name',
          type: 'input',
          label: ts('pages.assets.assetTemplates.filters.name'),
          icon: 'search',
          grid: 'col-sm-12 col-md-6'
        },
        {
          key: 'status',
          type: 'select',
          label: ts('pages.assets.assetTemplates.filters.status'),
          icon: 'toggle_on',
          options: [
            { label: ts('pages.assets.assetTemplates.filters.statusOptions.all'), value: undefined },
            { label: ts('pages.assets.assetTemplates.filters.statusOptions.active'), value: true },
            { label: ts('pages.assets.assetTemplates.filters.statusOptions.inactive'), value: false },
          ],
          grid: 'col-sm-12 col-6 col-md-3'
        },
        {
          key: 'includeChildren',
          type: 'select',
          label: ts('pages.assets.assetTemplates.filters.includeChildren'),
          icon: 'account_tree',
          options: [
            { label: ts('pages.assets.assetTemplates.filters.statusOptions.all'), value: null },
            { label: ts('pages.assets.assetTemplates.filters.includeChildrenOptions.yes'), value: true },
            { label: ts('pages.assets.assetTemplates.filters.includeChildrenOptions.no'), value: false },
          ],
          grid: 'col-sm-12 col-6 col-md-3'
        },

        // Row 2: Domain-Specific Filters
        // Category, Manufacturer, Model, Template Type
        {
          key: 'category',
          type: 'input',
          label: ts('pages.assets.assetTemplates.filters.category'),
          icon: 'category',
          grid: 'col-sm-12 col-6 col-md-3'
        },
        {
          key: 'manufacture',
          type: 'input',
          label: ts('pages.assets.assetTemplates.filters.manufacturer'),
          icon: 'factory',
          grid: 'col-sm-12 col-6 col-md-3'
        },
        {
          key: 'model',
          type: 'input',
          label: ts('pages.assets.assetTemplates.filters.model'),
          icon: 'memory',
          grid: 'col-sm-12 col-6 col-md-3'
        },
        {
          key: 'isSystem',
          type: 'select',
          label: ts('pages.assets.assetTemplates.filters.isSystem'),
          icon: 'lock',
          options: [
            { label: ts('pages.assets.assetTemplates.filters.statusOptions.all'), value: null },
            { label: ts('pages.assets.assetTemplates.filters.isSystemOptions.system'), value: true },
            { label: ts('pages.assets.assetTemplates.filters.isSystemOptions.custom'), value: false },
          ],
          grid: 'col-sm-12 col-6 col-md-3',
        },
      ];

      // Add isTemplate filter only for Customer and Site organizations
      if (organizationStore.isCustomer || organizationStore.isSite) {
        baseFilters.splice(baseFilters.length - 1, 0, {
          key: 'isTemplate',
          type: 'select',
          label: ts('pages.assets.assetTemplates.filters.isTemplate'),
          icon: 'content_copy',
          options: [
            { label: ts('pages.assets.assetTemplates.filters.statusOptions.all'), value: null },
            { label: ts('pages.assets.assetTemplates.filters.isTemplateOptions.templates'), value: true },
            { label: ts('pages.assets.assetTemplates.filters.isTemplateOptions.local'), value: false },
          ],
          grid: 'col-sm-12 col-6 col-md-3',
        });
      }

      return baseFilters;
    }),

    listHeader: {
      title: computed(() => tsTitle('pages.assets.assetTemplates.listHeader.title')),
      itemLabel: computed(() => ts('pages.assets.assetTemplates.listHeader.itemLabel')),
      itemLabelPlural: computed(() => ts('pages.assets.assetTemplates.listHeader.itemLabelPlural')),
    },

    menuColumns: {
      organization: computed(() => ts('pages.assets.assetTemplates.menuColumns.organization')),
      manufacturerModel: computed(() => ts('pages.assets.assetTemplates.columns.manufacturer')),
      version: computed(() => ts('pages.assets.assetTemplates.columns.version')),
      templateType: computed(() => ts('pages.assets.assetTemplates.columns.templateType')),
      templateSource: computed(() => ts('pages.assets.assetTemplates.columns.templateSource')),
    },

    columns: computed((): DataRowColumn[] => [
      {
        key: 'icon',
        label: '',
        type: 'avatar',
        visible: 'always',
        width: 56,
        icon: (value: any, row: any) => row.icon || 'memory',
        color: (value: any, row: any) => row.enabled ? 'primary' : 'grey-5',
        tooltip: (value: any, row: any) =>
          row.enabled
            ? ts('pages.assets.assetTemplates.status.active')
            : ts('pages.assets.assetTemplates.status.inactive'),
      },
      {
        key: 'name',
        label: ts('pages.assets.assetTemplates.columns.name'),
        type: 'text',
        visible: 'always',
        width: 240,
        ellipsis: true,
        secondaryKey: 'description',
      },
      {
        key: 'organizationName',
        label: ts('pages.assets.assetTemplates.columns.organization'),
        type: 'chip',
        visible: 'laptop',
        width: 180,
        ellipsis: true,
        color: 'indigo-6',
        icon: 'domain',
      },
      {
        key: 'manufacturer',
        label: ts('pages.assets.assetTemplates.columns.manufacturer'),
        type: 'text',
        visible: 'laptop', // 💻🖥️ Laptop & Desktop (>= 1024px)
        width: 280,
        ellipsis: true,
        secondaryKey: 'deviceModel',
      },
      {
        key: 'version',
        label: ts('pages.assets.assetTemplates.columns.version'),
        type: 'chip',
        visible: 'laptop', // 💻🖥️ Laptop & Desktop (>= 1024px) - VERSION IS IMPORTANT!
        width: 80,
        color: 'purple-6',
        align: 'center',
      },
      {
        key: 'isSystem',
        label: ts('pages.assets.assetTemplates.columns.templateType'),
        type: 'chip',
        visible: 'laptop',
        width: 80,
        format: (value: any) => value
          ? ts('pages.assets.assetTemplates.filters.isSystemOptions.system').toUpperCase()
          : ts('pages.assets.assetTemplates.filters.isSystemOptions.custom').toUpperCase(),
        color: (value: any) => value ? 'purple-6' : 'blue-6',
        icon: (value: any) => value ? 'lock' : 'edit',
        align: 'center',
      },
      {
        key: 'isTemplate',
        label: ts('pages.assets.assetTemplates.columns.templateSource'),
        type: 'chip',
        visible: 'laptop',
        width: 80,
        format: (value: any) => value
          ? ts('pages.assets.assetTemplates.filters.isTemplateOptions.templates').toUpperCase()
          : ts('pages.assets.assetTemplates.filters.isTemplateOptions.local').toUpperCase(),
        color: (value: any) => value ? 'orange-6' : 'green-6',
        icon: (value: any) => value ? 'content_copy' : 'folder',
        align: 'center',
      },
    ]),

    empty: {
      title: computed(() => ts('pages.assets.assetTemplates.empty.title')),
      description: computed(() => ts('pages.assets.assetTemplates.empty.description')),
    },

    dialog: {
      confirmDelete: {
        title: computed(() => ts('pages.assets.assetTemplates.dialog.confirmDelete.title')),
        message: (name: string) => ts('pages.assets.assetTemplates.dialog.confirmDelete.message', { name }),
      },
    },

    notifications: {
      deleted: computed(() => ts('pages.assets.assetTemplates.notifications.deleted')),
      deleteError: computed(() => ts('pages.assets.assetTemplates.notifications.deleteError')),
      systemTemplateEdit: computed(() => ts('pages.assets.assetTemplates.notifications.systemTemplateEdit')),
      systemTemplateDelete: computed(() => ts('pages.assets.assetTemplates.notifications.systemTemplateDelete')),
      sharedTemplateEdit: computed(() => ts('pages.assets.assetTemplates.notifications.sharedTemplateEdit')),
      sharedTemplateDelete: computed(() => ts('pages.assets.assetTemplates.notifications.sharedTemplateDelete')),
    },

    errors: {
      apiNotInitialized: computed(() => ts('pages.assets.assetTemplates.errors.apiNotInitialized')),
      idMissing: computed(() => ts('pages.assets.assetTemplates.errors.idMissing')),
    },

    actions: {
      edit: computed(() => ts('pages.assets.assetTemplates.actions.edit')),
      view: computed(() => ts('pages.assets.assetTemplates.actions.view')),
      delete: computed(() => ts('pages.assets.assetTemplates.actions.delete')),
    },

    status: {
      active: computed(() => ts('pages.assets.assetTemplates.status.active')),
      inactive: computed(() => ts('pages.assets.assetTemplates.status.inactive')),
    },

    drawer: {
      title: computed(() => ts('pages.assets.assetTemplates.drawer.title')),
      close: computed(() => ts('pages.assets.assetTemplates.drawer.close')),
      edit: computed(() => ts('pages.assets.assetTemplates.drawer.edit')),
      duplicate: computed(() => ts('pages.assets.assetTemplates.drawer.duplicate')),
      loading: computed(() => ts('pages.assets.assetTemplates.drawer.loading')),
      error: computed(() => ts('pages.assets.assetTemplates.drawer.error')),
      systemTemplateWarning: computed(() => ts('pages.assets.assetTemplates.drawer.systemTemplateWarning')),
      systemTemplateTooltip: computed(() => ts('pages.assets.assetTemplates.drawer.systemTemplateTooltip')),

      sections: {
        basicInfo: computed(() => ts('pages.assets.assetTemplates.drawer.sections.basicInfo')),
        configuration: computed(() => ts('pages.assets.assetTemplates.drawer.sections.configuration')),
        scripts: computed(() => ts('pages.assets.assetTemplates.drawer.sections.scripts')),
        timestamps: computed(() => ts('pages.assets.assetTemplates.drawer.sections.timestamps')),
      },

      fields: {
        name: computed(() => ts('pages.assets.assetTemplates.drawer.fields.name')),
        status: computed(() => ts('pages.assets.assetTemplates.drawer.fields.status')),
        description: computed(() => ts('pages.assets.assetTemplates.drawer.fields.description')),
        isSystem: computed(() => ts('pages.assets.assetTemplates.drawer.fields.isSystem')),
        manufacturer: computed(() => ts('pages.assets.assetTemplates.drawer.fields.manufacturer')),
        model: computed(() => ts('pages.assets.assetTemplates.drawer.fields.model')),
        version: computed(() => ts('pages.assets.assetTemplates.drawer.fields.version')),
        assetIdPath: computed(() => ts('pages.assets.assetTemplates.drawer.fields.assetIdPath')),
        scriptTest: computed(() => ts('pages.assets.assetTemplates.drawer.fields.scriptTest')),
        scriptProcessor: computed(() => ts('pages.assets.assetTemplates.drawer.fields.scriptProcessor')),
        scriptValidator: computed(() => ts('pages.assets.assetTemplates.drawer.fields.scriptValidator')),
        scriptConversion: computed(() => ts('pages.assets.assetTemplates.drawer.fields.scriptConversion')),
        scriptsSummary: computed(() => ts('pages.assets.assetTemplates.drawer.fields.scriptsSummary')),
        created: computed(() => ts('pages.assets.assetTemplates.drawer.fields.created')),
        updated: computed(() => ts('pages.assets.assetTemplates.drawer.fields.updated')),
      },

      empty: {
        description: computed(() => ts('pages.assets.assetTemplates.drawer.empty.description')),
      },

      system: {
        yes: computed(() => ts('pages.assets.assetTemplates.drawer.system.yes')),
        no: computed(() => ts('pages.assets.assetTemplates.drawer.system.no')),
      },

      scripts: {
        configured: computed(() => ts('pages.assets.assetTemplates.drawer.scripts.configured')),
        notConfigured: computed(() => ts('pages.assets.assetTemplates.drawer.scripts.notConfigured')),
      },

      scriptViewer: {
        viewScript: computed(() => ts('pages.assets.assetTemplates.drawer.scriptViewer.viewScript')),
        copyScript: computed(() => ts('pages.assets.assetTemplates.drawer.scriptViewer.copyScript')),
        close: computed(() => ts('pages.assets.assetTemplates.drawer.scriptViewer.close')),
        copySuccess: computed(() => ts('pages.assets.assetTemplates.drawer.scriptViewer.copySuccess')),
        copyFail: computed(() => ts('pages.assets.assetTemplates.drawer.scriptViewer.copyFail')),
      },
    },
  };
}
