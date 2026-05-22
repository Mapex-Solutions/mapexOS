import { computed } from 'vue';
import { useTS } from '@utils/translation';
import type { DataRowColumn } from '@components/cards';
import type { FilterListItem } from '@components/filters';
import type { PageHeaderInfo } from '@components/headers';
import { useOrganizationStore } from '@stores/organization';

/**
 * Roles list page translations
 *
 * Structure mirrors:
 * - File: src/pages/administrations/roles/rolesListPage/RolesListPage.vue
 * - JSON: src/i18n/{locale}/pages/administrations/roles.json
 * - Composable: src/composables/i18n/pages/administrations/roles/useRolesTranslations.ts
 */
export function useRolesTranslations() {
  const ts = useTS({ capitalize: true });
  const tsTitle = useTS({ titleCase: true });
  const organizationStore = useOrganizationStore();

  return {
    /**
     * Page header translations
     */
    page: {
      title: computed(() => tsTitle('pages.administrations.roles.title')),
      titleEdit: computed(() => tsTitle('pages.administrations.roles.titleEdit')),
      titleCreate: computed(() => tsTitle('pages.administrations.roles.titleCreate')),
      description: computed(() => ts('pages.administrations.roles.description')),
      descriptionCreate: computed(() => ts('pages.administrations.roles.descriptionCreate')),
      backButton: computed(() => ts('pages.administrations.roles.backButton')),
      addButton: computed(() => ts('pages.administrations.roles.addButton')),
      listTitle: computed(() => tsTitle('pages.administrations.roles.listTitle')),
      itemLabel: computed(() => ts('pages.administrations.roles.itemLabel')),
      itemLabelPlural: computed(() => ts('pages.administrations.roles.itemLabelPlural')),
      info: computed((): PageHeaderInfo => ({
        title: ts('pages.administrations.roles.info.title'),
        description: ts('pages.administrations.roles.info.description'),
        items: [
          {
            icon: 'admin_panel_settings',
            color: 'blue-6',
            title: ts('pages.administrations.roles.info.items.permissions.title'),
            text: ts('pages.administrations.roles.info.items.permissions.text'),
          },
          {
            icon: 'lock',
            color: 'green-6',
            title: ts('pages.administrations.roles.info.items.systemCustom.title'),
            text: ts('pages.administrations.roles.info.items.systemCustom.text'),
          },
          {
            icon: 'verified_user',
            color: 'orange-6',
            title: ts('pages.administrations.roles.info.items.accessControl.title'),
            text: ts('pages.administrations.roles.info.items.accessControl.text'),
          },
          {
            icon: 'group',
            color: 'purple-6',
            title: ts('pages.administrations.roles.info.items.userAssignment.title'),
            text: ts('pages.administrations.roles.info.items.userAssignment.text'),
          },
          {
            icon: 'tune',
            color: 'indigo-6',
            title: ts('pages.administrations.roles.info.items.granular.title'),
            text: ts('pages.administrations.roles.info.items.granular.text'),
          },
        ],
        docsUrl: 'https://docs.mapexos.com/administration/roles',
        docsLabel: ts('pages.administrations.roles.info.docsLabel'),
      })),
    },

    /**
     * Filter translations
     */
    filters: {
      label: computed(() => ts('pages.administrations.roles.filters.label')),
      searchPlaceholder: computed(() => ts('pages.administrations.roles.filters.searchPlaceholder')),
      allStatus: computed(() => ts('pages.administrations.roles.filters.allStatus')),
      advancedFilters: computed(() => ts('pages.administrations.roles.filters.advancedFilters')),
      pendingFilters: computed(() => ts('pages.administrations.roles.filters.pendingFilters')),
      clearAll: computed(() => ts('pages.administrations.roles.filters.clearAll')),
      name: computed(() => ts('pages.administrations.roles.filters.name')),
      isSystem: computed(() => ts('pages.administrations.roles.filters.isSystem')),
      includeChildren: computed(() => ts('pages.administrations.roles.filters.includeChildren')),
      includeChildrenOrgs: computed(() => ts('pages.administrations.roles.filters.includeChildrenOrgs')),
      scope: computed(() => ts('pages.administrations.roles.filters.scope')),
      permission: computed(() => ts('pages.administrations.roles.filters.permission')),
      filterByPermission: computed(() => ts('pages.administrations.roles.filters.filterByPermission')),
      isTemplate: computed(() => ts('pages.administrations.roles.filters.isTemplate')),

      options: {
        yes: computed(() => ts('pages.administrations.roles.filters.options.yes')),
        no: computed(() => ts('pages.administrations.roles.filters.options.no')),
        system: computed(() => ts('pages.administrations.roles.filters.options.system')),
        custom: computed(() => ts('pages.administrations.roles.filters.options.custom')),
        global: computed(() => ts('pages.administrations.roles.filters.options.global')),
        local: computed(() => ts('pages.administrations.roles.filters.options.local')),
        templates: computed(() => ts('pages.administrations.roles.filters.options.templates')),
      },
    },

    /**
     * Menu column labels
     */
    menuColumns: {
      organization: computed(() => ts('pages.administrations.roles.menuColumns.organization')),
      description: computed(() => ts('pages.administrations.roles.menuColumns.description')),
      permissions: computed(() => ts('pages.administrations.roles.menuColumns.permissions')),
      scope: computed(() => ts('pages.administrations.roles.menuColumns.scope')),
      templateSource: computed(() => ts('pages.administrations.roles.columns.templateSource')),
      created: computed(() => ts('pages.administrations.roles.menuColumns.created')),
    },

    /**
     * Empty state translations
     */
    empty: {
      title: computed(() => ts('pages.administrations.roles.empty.title')),
      description: computed(() => ts('pages.administrations.roles.empty.description')),
    },

    /**
     * Scope labels
     */
    scope: {
      global: computed(() => ts('pages.administrations.roles.scope.global')),
      local: computed(() => ts('pages.administrations.roles.scope.local')),
    },

    /**
     * Type labels
     */
    type: {
      system: computed(() => ts('pages.administrations.roles.type.system')),
      custom: computed(() => ts('pages.administrations.roles.type.custom')),
    },

    /**
     * Dialog translations
     */
    dialog: {
      deleteTitle: computed(() => ts('pages.administrations.roles.dialog.deleteTitle')),
    },

    /**
     * Notification translations
     */
    notifications: {
      systemEdit: computed(() => ts('pages.administrations.roles.notifications.systemEdit')),
      systemDelete: computed(() => ts('pages.administrations.roles.notifications.systemDelete')),
      sharedEdit: computed(() => ts('pages.administrations.roles.notifications.sharedEdit')),
      sharedDelete: computed(() => ts('pages.administrations.roles.notifications.sharedDelete')),
      created: computed(() => ts('pages.administrations.roles.createEdit.notifications.created')),
      updated: computed(() => ts('pages.administrations.roles.createEdit.notifications.updated')),
      createFailed: computed(() => ts('pages.administrations.roles.createEdit.notifications.createFailed')),
      updateFailed: computed(() => ts('pages.administrations.roles.createEdit.notifications.updateFailed')),
      loadFailed: computed(() => ts('pages.administrations.roles.createEdit.notifications.loadFailed')),
      noPermissions: computed(() => ts('pages.administrations.roles.createEdit.notifications.noPermissions')),
      alreadyExists: computed(() => ts('pages.administrations.roles.createEdit.notifications.alreadyExists')),
      forbidden: computed(() => ts('pages.administrations.roles.createEdit.notifications.forbidden')),
    },

    /**
     * Error translations
     */
    errors: {
      apiNotInitialized: computed(() => ts('pages.administrations.roles.errors.apiNotInitialized')),
      idMissing: computed(() => ts('pages.administrations.roles.errors.idMissing')),
      orgPathKeyMissing: computed(() => ts('pages.administrations.roles.errors.orgPathKeyMissing')),
    },

    /**
     * CreateEdit page translations - Steps
     */
    steps: {
      basicInfo: computed(() => tsTitle('pages.administrations.roles.createEdit.steps.basicInfo')),
      permissions: computed(() => tsTitle('pages.administrations.roles.createEdit.steps.permissions')),
      review: computed(() => tsTitle('pages.administrations.roles.createEdit.steps.review')),
    },

    /**
     * CreateEdit page translations - Step descriptions
     */
    stepDescriptions: {
      basicInfo: computed(() => ts('pages.administrations.roles.createEdit.stepDescriptions.basicInfo')),
      permissions: computed(() => ts('pages.administrations.roles.createEdit.stepDescriptions.permissions')),
      review: computed(() => ts('pages.administrations.roles.createEdit.stepDescriptions.review')),
    },

    /**
     * CreateEdit page translations - Sections
     */
    sections: {
      basicInfo: computed(() => tsTitle('pages.administrations.roles.createEdit.sections.basicInfo')),
      permissions: computed(() => tsTitle('pages.administrations.roles.createEdit.sections.permissions')),
      progressSteps: computed(() => tsTitle('pages.administrations.roles.createEdit.sections.progressSteps')),
    },

    /**
     * CreateEdit page translations - Form descriptions
     */
    formDescriptions: {
      basicInfo: computed(() => ts('pages.administrations.roles.createEdit.formDescriptions.basicInfo')),
      permissions: computed(() => ts('pages.administrations.roles.createEdit.formDescriptions.permissions')),
      isTemplate: computed(() => ts('pages.administrations.roles.createEdit.formDescriptions.isTemplate')),
    },

    /**
     * CreateEdit page translations - Fields
     */
    fields: {
      name: computed(() => ts('pages.administrations.roles.createEdit.fields.name')),
      description: computed(() => ts('pages.administrations.roles.createEdit.fields.description')),
      scope: computed(() => ts('pages.administrations.roles.createEdit.fields.scope')),
      isTemplate: computed(() => ts('pages.administrations.roles.createEdit.fields.isTemplate')),
    },

    /**
     * CreateEdit page translations - Scope options
     */
    scopeOptions: {
      global: computed(() => tsTitle('pages.administrations.roles.createEdit.scopeOptions.global')),
      local: computed(() => tsTitle('pages.administrations.roles.createEdit.scopeOptions.local')),
    },

    /**
     * CreateEdit page translations - Validation
     */
    validation: {
      nameRequired: computed(() => ts('pages.administrations.roles.createEdit.validation.nameRequired')),
      nameMinLength: computed(() => ts('pages.administrations.roles.createEdit.validation.nameMinLength')),
      nameMaxLength: computed(() => ts('pages.administrations.roles.createEdit.validation.nameMaxLength')),
      descriptionMaxLength: computed(() => ts('pages.administrations.roles.createEdit.validation.descriptionMaxLength')),
      scopeRequired: computed(() => ts('pages.administrations.roles.createEdit.validation.scopeRequired')),
      permissionsRequired: computed(() => ts('pages.administrations.roles.createEdit.validation.permissionsRequired')),
    },

    /**
     * CreateEdit page translations - Labels
     */
    labels: {
      selectAll: computed(() => ts('pages.administrations.roles.createEdit.labels.selectAll')),
      deselectAll: computed(() => ts('pages.administrations.roles.createEdit.labels.deselectAll')),
      actionsSelected: computed(() => ts('pages.administrations.roles.createEdit.labels.actionsSelected')),
      permissionSelected: computed(() => ts('pages.administrations.roles.createEdit.labels.permissionSelected')),
      permissionsSelected: computed(() => ts('pages.administrations.roles.createEdit.labels.permissionsSelected')),
      permission: computed(() => ts('pages.administrations.roles.createEdit.labels.permission')),
      permissions: computed(() => ts('pages.administrations.roles.createEdit.labels.permissions')),
    },

    /**
     * CreateEdit page translations - Buttons
     */
    buttons: {
      back: computed(() => ts('pages.administrations.roles.createEdit.buttons.back')),
      next: computed(() => ts('pages.administrations.roles.createEdit.buttons.next')),
      createRole: computed(() => ts('pages.administrations.roles.createEdit.buttons.createRole')),
      updateRole: computed(() => ts('pages.administrations.roles.createEdit.buttons.updateRole')),
    },

    /**
     * CreateEdit page translations - Messages
     */
    messages: {
      deletedSuccessfully: computed(() => ts('pages.administrations.roles.messages.deletedSuccessfully')),
      confirmDelete: (name: string) => ts('pages.administrations.roles.messages.confirmDelete', { name }),
      loadingRole: computed(() => ts('pages.administrations.roles.createEdit.messages.loadingRole')),
      completeAllSteps: computed(() => ts('pages.administrations.roles.createEdit.messages.completeAllSteps')),
      allFieldsRequired: computed(() => ts('pages.administrations.roles.createEdit.messages.allFieldsRequired')),
      currentStep: computed(() => ts('pages.administrations.roles.createEdit.messages.currentStep')),
      systemRoleWarning: computed(() => ts('pages.administrations.roles.createEdit.messages.systemRoleWarning')),
    },

    /**
     * CreateEdit page translations - Review step
     */
    reviewStep: {
      subtitle: computed(() => ts('pages.administrations.roles.createEdit.reviewStep.subtitle')),
      successMessage: computed(() => ts('pages.administrations.roles.createEdit.reviewStep.successMessage')),
      successMessageEdit: computed(() => ts('pages.administrations.roles.createEdit.reviewStep.successMessageEdit')),
      sections: {
        basicInfo: computed(() => tsTitle('pages.administrations.roles.createEdit.reviewStep.sections.basicInfo')),
        permissions: computed(() => tsTitle('pages.administrations.roles.createEdit.reviewStep.sections.permissions')),
      },
      fields: {
        name: computed(() => ts('pages.administrations.roles.createEdit.reviewStep.fields.name')),
        description: computed(() => ts('pages.administrations.roles.createEdit.reviewStep.fields.description')),
        scope: computed(() => ts('pages.administrations.roles.createEdit.reviewStep.fields.scope')),
        isTemplate: computed(() => ts('pages.administrations.roles.createEdit.reviewStep.fields.isTemplate')),
        totalPermissions: computed(() => ts('pages.administrations.roles.createEdit.reviewStep.fields.totalPermissions')),
        enabledResources: computed(() => ts('pages.administrations.roles.createEdit.reviewStep.fields.enabledResources')),
      },
      values: {
        notProvided: computed(() => ts('pages.administrations.roles.createEdit.reviewStep.values.notProvided')),
        notSelected: computed(() => ts('pages.administrations.roles.createEdit.reviewStep.values.notSelected')),
        noPermissions: computed(() => ts('pages.administrations.roles.createEdit.reviewStep.values.noPermissions')),
      },
    },

    /**
     * Drawer translations
     */
    drawer: {
      title: computed(() => ts('pages.administrations.roles.drawer.title')),
      close: computed(() => ts('pages.administrations.roles.drawer.close')),
      edit: computed(() => ts('pages.administrations.roles.drawer.edit')),
      loading: computed(() => ts('pages.administrations.roles.drawer.loading')),
      error: computed(() => ts('pages.administrations.roles.drawer.error')),
      systemRoleWarning: computed(() => ts('pages.administrations.roles.drawer.systemRoleWarning')),
      systemRoleTooltip: computed(() => ts('pages.administrations.roles.drawer.systemRoleTooltip')),

      sections: {
        basicInfo: computed(() => ts('pages.administrations.roles.drawer.sections.basicInfo')),
        permissions: computed(() => ts('pages.administrations.roles.drawer.sections.permissions')),
        scope: computed(() => ts('pages.administrations.roles.drawer.sections.scope')),
        timestamps: computed(() => ts('pages.administrations.roles.drawer.sections.timestamps')),
      },

      fields: {
        name: computed(() => ts('pages.administrations.roles.drawer.fields.name')),
        description: computed(() => ts('pages.administrations.roles.drawer.fields.description')),
        type: computed(() => ts('pages.administrations.roles.drawer.fields.type')),
        scope: computed(() => ts('pages.administrations.roles.drawer.fields.scope')),
        organization: computed(() => ts('pages.administrations.roles.drawer.fields.organization')),
        pathKey: computed(() => ts('pages.administrations.roles.drawer.fields.pathKey')),
        permissionsCount: computed(() => ts('pages.administrations.roles.drawer.fields.permissionsCount')),
        created: computed(() => ts('pages.administrations.roles.drawer.fields.created')),
        updated: computed(() => ts('pages.administrations.roles.drawer.fields.updated')),
      },

      empty: {
        description: computed(() => ts('pages.administrations.roles.drawer.empty.description')),
        permissions: computed(() => ts('pages.administrations.roles.drawer.empty.permissions')),
      },

      system: {
        yes: computed(() => ts('pages.administrations.roles.drawer.system.yes')),
        no: computed(() => ts('pages.administrations.roles.drawer.system.no')),
      },
    },

    /**
     * DataRow column definitions with reactive translations
     */
    columns: computed(() => {
      return [
        {
          key: 'avatar',
          label: '',
          type: 'avatar',
          visible: 'always',
          width: 56,
          icon: () => 'admin_panel_settings',
          color: (value: any, row: any) => row.isSystem ? 'purple-6' : 'primary',
        },
        {
          key: 'name',
          label: ts('pages.administrations.roles.columns.name'),
          type: 'text',
          visible: 'always',
          width: 250,
          ellipsis: true,
          secondaryKey: 'description',
        },
        {
          key: 'organizationName',
          label: ts('pages.administrations.roles.columns.organization'),
          type: 'chip',
          visible: 'laptop',
          width: 180,
          ellipsis: true,
          color: 'indigo-6',
          icon: 'domain',
        },
        {
          key: 'permissions',
          label: ts('pages.administrations.roles.columns.permissions'),
          type: 'chip',
          visible: 'laptop',
          width: 150,
          format: (value: any, row: any) => {
            const count = row.permissions?.length || 0;
            return `${count}`;
          },
          color: 'blue-6',
        },
        {
          key: 'scope',
          label: ts('pages.administrations.roles.columns.scope'),
          type: 'badge',
          visible: 'laptop',
          width: 100,
          format: (value: any) => value ? value.toUpperCase() : 'N/A',
          color: (value: any) => value === 'global' ? 'purple-6' : 'orange-6',
        },
        {
          key: 'type',
          label: ts('pages.administrations.roles.columns.type'),
          type: 'badge',
          visible: 'laptop',
          width: 100,
          format: (value: any, row: any) => row.isSystem ? 'SYSTEM' : 'CUSTOM',
          color: (value: any, row: any) => row.isSystem ? 'purple-6' : 'blue-6',
        },
        {
          key: 'isTemplate',
          label: ts('pages.administrations.roles.columns.templateSource'),
          type: 'chip',
          visible: 'laptop',
          width: 120,
          format: (value: any) => value
            ? ts('pages.administrations.roles.filters.options.templates').toUpperCase()
            : ts('pages.administrations.roles.filters.options.local').toUpperCase(),
          color: (value: any) => value ? 'orange-6' : 'green-6',
          icon: (value: any) => value ? 'content_copy' : 'folder',
        },
        {
          key: 'created',
          label: ts('pages.administrations.roles.columns.created'),
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

    /**
     * Filter items with reactive translations
     */
    filterItems: computed(() => {
      const baseFilters: FilterListItem[] = [
        // Row 1: Standard pattern (6 + 3 + 3 = 12 cols)
        {
          key: 'name',
          type: 'input',
          label: ts('pages.administrations.roles.filters.name'),
          icon: 'search',
          grid: 'col-12 col-md-6'
        },
        {
          key: 'isSystem',
          type: 'select',
          label: ts('pages.administrations.roles.filters.isSystem'),
          icon: 'lock',
          options: [
            { label: ts('pages.administrations.roles.filters.options.all'), value: null },
            { label: ts('pages.administrations.roles.filters.options.system'), value: true },
            { label: ts('pages.administrations.roles.filters.options.custom'), value: false }
          ],
          grid: 'col-6 col-md-3'
        },
        {
          key: 'includeChildren',
          type: 'select',
          label: ts('pages.administrations.roles.filters.includeChildren'),
          icon: 'account_tree',
          options: [
            { label: ts('pages.administrations.roles.filters.options.all'), value: null },
            { label: ts('pages.administrations.roles.filters.options.yes'), value: true },
            { label: ts('pages.administrations.roles.filters.options.no'), value: false }
          ],
          grid: 'col-6 col-md-3'
        },
        // Row 2: Domain-specific filters (6 + 6 = 12 cols)
        {
          key: 'scope',
          type: 'select',
          label: ts('pages.administrations.roles.filters.scope'),
          icon: 'public',
          options: [
            { label: ts('pages.administrations.roles.filters.options.all'), value: null },
            { label: ts('pages.administrations.roles.filters.options.global'), value: 'global' },
            { label: ts('pages.administrations.roles.filters.options.local'), value: 'local' }
          ],
          grid: 'col-12 col-md-6'
        },
        {
          key: 'permission',
          type: 'input',
          label: ts('pages.administrations.roles.filters.permission'),
          icon: 'vpn_key',
          grid: 'col-12 col-md-6'
        },
      ];

      // Add isTemplate filter only for Customer and Site organizations
      if (organizationStore.isCustomer || organizationStore.isSite) {
        baseFilters.push({
          key: 'isTemplate',
          type: 'select',
          label: ts('pages.administrations.roles.filters.isTemplate'),
          icon: 'content_copy',
          options: [
            { label: ts('pages.administrations.roles.filters.options.all'), value: null },
            { label: ts('pages.administrations.roles.filters.options.templates'), value: true },
            { label: ts('pages.administrations.roles.filters.options.local'), value: false }
          ],
          grid: 'col-12 col-md-6'
        });
      }

      return baseFilters;
    }),
  };
}
