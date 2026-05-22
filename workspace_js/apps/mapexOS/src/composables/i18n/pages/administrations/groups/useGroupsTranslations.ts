import { computed } from 'vue';
import { useTS } from '@utils/translation';
import type { DataRowColumn } from '@components/cards';
import type { FilterListItem } from '@components/filters';
import type { PageHeaderInfo } from '@components/headers';

/**
 * Groups list page translations
 *
 * Structure mirrors:
 * - File: src/pages/administrations/groups/groupsListPage/GroupsListPage.vue
 * - JSON: src/i18n/{locale}/pages/administrations/groups.json
 * - Composable: src/composables/i18n/pages/administrations/groups/useGroupsTranslations.ts
 */
export function useGroupsTranslations() {
  const ts = useTS({ capitalize: true });
  const tsTitle = useTS({ titleCase: true });

  return {
    /**
     * Page header translations
     */
    page: {
      title: computed(() => tsTitle('pages.administrations.groups.title')),
      titleEdit: computed(() => ts('pages.administrations.groups.createEdit.titleEdit')),
      titleCreate: computed(() => ts('pages.administrations.groups.createEdit.titleCreate')),
      description: computed(() => ts('pages.administrations.groups.description')),
      descriptionCreate: computed(() => ts('pages.administrations.groups.createEdit.descriptionCreate')),
      addButton: computed(() => ts('pages.administrations.groups.addButton')),
      backButton: computed(() => ts('pages.administrations.groups.createEdit.backButton')),
      listTitle: computed(() => tsTitle('pages.administrations.groups.listTitle')),
      itemLabel: computed(() => ts('pages.administrations.groups.itemLabel')),
      itemLabelPlural: computed(() => ts('pages.administrations.groups.itemLabelPlural')),
      info: computed((): PageHeaderInfo => ({
        title: ts('pages.administrations.groups.info.title'),
        description: ts('pages.administrations.groups.info.description'),
        items: [
          {
            icon: 'route',
            color: 'blue-6',
            title: ts('pages.administrations.groups.info.items.routing.title'),
            text: ts('pages.administrations.groups.info.items.routing.text'),
          },
          {
            icon: 'hub',
            color: 'green-6',
            title: ts('pages.administrations.groups.info.items.destinations.title'),
            text: ts('pages.administrations.groups.info.items.destinations.text'),
          },
          {
            icon: 'link',
            color: 'orange-6',
            title: ts('pages.administrations.groups.info.items.association.title'),
            text: ts('pages.administrations.groups.info.items.association.text'),
          },
          {
            icon: 'toggle_on',
            color: 'purple-6',
            title: ts('pages.administrations.groups.info.items.control.title'),
            text: ts('pages.administrations.groups.info.items.control.text'),
          },
          {
            icon: 'tune',
            color: 'indigo-6',
            title: ts('pages.administrations.groups.info.items.configuration.title'),
            text: ts('pages.administrations.groups.info.items.configuration.text'),
          },
        ],
        docsUrl: 'https://docs.mapexos.com/routing/groups',
        docsLabel: ts('pages.administrations.groups.info.docsLabel'),
      })),
    },

    /**
     * Filter translations
     */
    filters: {
      label: computed(() => ts('pages.administrations.groups.filters.label')),
      searchPlaceholder: computed(() => ts('pages.administrations.groups.filters.searchPlaceholder')),
      advancedFilters: computed(() => ts('pages.administrations.groups.filters.advancedFilters')),
      pendingFilters: computed(() => ts('pages.administrations.groups.filters.pendingFilters')),
      clearAll: computed(() => ts('pages.administrations.groups.filters.clearAll')),
      allStatus: computed(() => ts('pages.administrations.groups.filters.allStatus')),
      name: computed(() => ts('pages.administrations.groups.filters.name')),
      status: computed(() => ts('pages.administrations.groups.filters.status')),
      includeChildren: computed(() => ts('pages.administrations.groups.filters.includeChildren')),
      includeChildrenOrgs: computed(() => ts('pages.administrations.groups.filters.includeChildrenOrgs')),
      member: computed(() => ts('pages.administrations.groups.filters.member')),
      filterByMember: computed(() => ts('pages.administrations.groups.filters.filterByMember')),
      searchMemberPlaceholder: computed(() => ts('pages.administrations.groups.filters.searchMemberPlaceholder')),

      options: {
        enabled: computed(() => ts('pages.administrations.groups.filters.options.enabled')),
        disabled: computed(() => ts('pages.administrations.groups.filters.options.disabled')),
        yes: computed(() => ts('pages.administrations.groups.filters.options.yes')),
        no: computed(() => ts('pages.administrations.groups.filters.options.no')),
      },
    },

    /**
     * Menu column labels
     */
    menuColumns: {
      organization: computed(() => ts('pages.administrations.groups.menuColumns.organization')),
      description: computed(() => ts('pages.administrations.groups.menuColumns.description')),
      members: computed(() => ts('pages.administrations.groups.menuColumns.members')),
      created: computed(() => ts('pages.administrations.groups.menuColumns.created')),
    },

    /**
     * Empty state translations
     */
    empty: {
      title: computed(() => ts('pages.administrations.groups.empty.title')),
      description: computed(() => ts('pages.administrations.groups.empty.description')),
    },

    /**
     * Status labels
     */
    status: {
      active: computed(() => ts('pages.administrations.groups.status.active')),
      inactive: computed(() => ts('pages.administrations.groups.status.inactive')),
    },

    /**
     * Dialog translations
     */
    dialog: {
      deleteTitle: computed(() => ts('pages.administrations.groups.dialog.deleteTitle')),
    },

    /**
     * Notification translations
     */
    notifications: {
      sharedEdit: computed(() => ts('pages.administrations.groups.notifications.sharedEdit')),
      sharedDelete: computed(() => ts('pages.administrations.groups.notifications.sharedDelete')),
    },

    /**
     * Error translations
     */
    errors: {
      apiNotInitialized: computed(() => ts('pages.administrations.groups.errors.apiNotInitialized')),
      idMissing: computed(() => ts('pages.administrations.groups.errors.idMissing')),
    },

    /**
     * Drawer translations
     */
    drawer: {
      title: computed(() => ts('pages.administrations.groups.drawer.title')),
      close: computed(() => ts('pages.administrations.groups.drawer.close')),
      edit: computed(() => ts('pages.administrations.groups.drawer.edit')),
      loading: computed(() => ts('pages.administrations.groups.drawer.loading')),
      error: computed(() => ts('pages.administrations.groups.drawer.error')),
      sections: {
        basicInfo: computed(() => ts('pages.administrations.groups.drawer.sections.basicInfo')),
        roles: computed(() => ts('pages.administrations.groups.drawer.sections.roles')),
        members: computed(() => ts('pages.administrations.groups.drawer.sections.members')),
        organization: computed(() => ts('pages.administrations.groups.drawer.sections.organization')),
        timestamps: computed(() => ts('pages.administrations.groups.drawer.sections.timestamps')),
      },
      fields: {
        name: computed(() => ts('pages.administrations.groups.drawer.fields.name')),
        description: computed(() => ts('pages.administrations.groups.drawer.fields.description')),
        enabled: computed(() => ts('pages.administrations.groups.drawer.fields.enabled')),
        isTemplate: computed(() => ts('pages.administrations.groups.drawer.fields.isTemplate')),
        organization: computed(() => ts('pages.administrations.groups.drawer.fields.organization')),
        pathKey: computed(() => ts('pages.administrations.groups.drawer.fields.pathKey')),
        membersCount: computed(() => ts('pages.administrations.groups.drawer.fields.membersCount')),
        created: computed(() => ts('pages.administrations.groups.drawer.fields.created')),
        updated: computed(() => ts('pages.administrations.groups.drawer.fields.updated')),
      },
      empty: {
        description: computed(() => ts('pages.administrations.groups.drawer.empty.description')),
        roles: computed(() => ts('pages.administrations.groups.drawer.empty.roles')),
        members: computed(() => ts('pages.administrations.groups.drawer.empty.members')),
      },
      status: {
        enabled: computed(() => ts('pages.administrations.groups.drawer.status.enabled')),
        disabled: computed(() => ts('pages.administrations.groups.drawer.status.disabled')),
      },
      template: {
        yes: computed(() => ts('pages.administrations.groups.drawer.template.yes')),
        no: computed(() => ts('pages.administrations.groups.drawer.template.no')),
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
          icon: () => 'groups',
          color: (value: any, row: any) => row.enabled ? 'primary' : 'grey-5',
          tooltip: (value: any, row: any) =>
            row.enabled
              ? ts('pages.administrations.groups.status.active')
              : ts('pages.administrations.groups.status.inactive'),
        },
        {
          key: 'name',
          label: ts('pages.administrations.groups.columns.name'),
          type: 'text',
          visible: 'always',
          width: 250,
          ellipsis: true,
          secondaryKey: 'description',
        },
        {
          key: 'organizationName',
          label: ts('pages.administrations.groups.columns.organization'),
          type: 'chip',
          visible: 'laptop',
          width: 180,
          ellipsis: true,
          color: 'indigo-6',
          icon: 'domain',
        },
        {
          key: 'membersCount',
          label: ts('pages.administrations.groups.columns.members'),
          type: 'chip',
          visible: 'laptop',
          width: 120,
          format: (value: any) => {
            return `${value || 0}`;
          },
          color: 'blue-6',
        },
        {
          key: 'created',
          label: ts('pages.administrations.groups.columns.created'),
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
          label: ts('pages.administrations.groups.filters.name'),
          icon: 'search',
          grid: 'col-12 col-md-6'
        },
        {
          key: 'enabled',
          type: 'select',
          label: ts('pages.administrations.groups.filters.status'),
          icon: 'toggle_on',
          options: [
            { label: ts('pages.administrations.groups.filters.options.all'), value: null },
            { label: ts('pages.administrations.groups.filters.options.enabled'), value: true },
            { label: ts('pages.administrations.groups.filters.options.disabled'), value: false }
          ],
          grid: 'col-6 col-md-3'
        },
        {
          key: 'includeChildren',
          type: 'select',
          label: ts('pages.administrations.groups.filters.includeChildren'),
          icon: 'account_tree',
          options: [
            { label: ts('pages.administrations.groups.filters.options.all'), value: null },
            { label: ts('pages.administrations.groups.filters.options.yes'), value: true },
            { label: ts('pages.administrations.groups.filters.options.no'), value: false }
          ],
          grid: 'col-6 col-md-3'
        },
        // Row 2: Member filter (user selector drawer)
        {
          key: 'memberId',
          type: 'user-select',
          label: ts('pages.administrations.groups.filters.member'),
          icon: 'person_search',
          placeholder: ts('pages.administrations.groups.filters.memberPlaceholder'),
          grid: 'col-12 col-md-6'
        },
      ];

      return baseFilters;
    }),

    /**
     * Step translations
     */
    steps: {
      basicInfo: computed(() => ts('pages.administrations.groups.createEdit.steps.basicInfo')),
      roles: computed(() => ts('pages.administrations.groups.createEdit.steps.roles')),
      members: computed(() => ts('pages.administrations.groups.createEdit.steps.members')),
      review: computed(() => ts('pages.administrations.groups.createEdit.steps.review')),
    },

    /**
     * Step description translations
     */
    stepDescriptions: {
      basicInfo: computed(() => ts('pages.administrations.groups.createEdit.stepDescriptions.basicInfo')),
      roles: computed(() => ts('pages.administrations.groups.createEdit.stepDescriptions.roles')),
      members: computed(() => ts('pages.administrations.groups.createEdit.stepDescriptions.members')),
      review: computed(() => ts('pages.administrations.groups.createEdit.stepDescriptions.review')),
    },

    /**
     * Section translations
     */
    sections: {
      basicInfo: computed(() => ts('pages.administrations.groups.createEdit.sections.basicInfo')),
      roles: computed(() => ts('pages.administrations.groups.createEdit.sections.roles')),
      members: computed(() => ts('pages.administrations.groups.createEdit.sections.members')),
      progressSteps: computed(() => ts('pages.administrations.groups.createEdit.sections.progressSteps')),
    },

    /**
     * Form description translations
     */
    formDescriptions: {
      basicInfo: computed(() => ts('pages.administrations.groups.createEdit.formDescriptions.basicInfo')),
      roles: computed(() => ts('pages.administrations.groups.createEdit.formDescriptions.roles')),
      members: computed(() => ts('pages.administrations.groups.createEdit.formDescriptions.members')),
      name: computed(() => ts('pages.administrations.groups.createEdit.formDescriptions.name')),
      description: computed(() => ts('pages.administrations.groups.createEdit.formDescriptions.description')),
      enabled: computed(() => ts('pages.administrations.groups.createEdit.formDescriptions.enabled')),
    },

    /**
     * Field translations
     */
    fields: {
      name: computed(() => ts('pages.administrations.groups.createEdit.fields.name')),
      description: computed(() => ts('pages.administrations.groups.createEdit.fields.description')),
      status: computed(() => ts('pages.administrations.groups.createEdit.fields.status')),
      enabled: computed(() => ts('pages.administrations.groups.createEdit.fields.enabled')),
    },

    /**
     * Status options for select
     */
    statusOptions: {
      active: computed(() => ts('pages.administrations.groups.createEdit.statusOptions.active')),
      inactive: computed(() => ts('pages.administrations.groups.createEdit.statusOptions.inactive')),
    },

    /**
     * Validation translations
     */
    validation: {
      nameRequired: computed(() => ts('pages.administrations.groups.createEdit.validation.nameRequired')),
      nameMinLength: computed(() => ts('pages.administrations.groups.createEdit.validation.nameMinLength')),
      nameMaxLength: computed(() => ts('pages.administrations.groups.createEdit.validation.nameMaxLength')),
      descriptionMaxLength: computed(() => ts('pages.administrations.groups.createEdit.validation.descriptionMaxLength')),
    },

    /**
     * Label translations
     */
    labels: {
      selectAll: computed(() => ts('pages.administrations.groups.createEdit.labels.selectAll')),
      deselectAll: computed(() => ts('pages.administrations.groups.createEdit.labels.deselectAll')),
      membersSelected: computed(() => ts('pages.administrations.groups.createEdit.labels.membersSelected')),
      memberSelected: computed(() => ts('pages.administrations.groups.createEdit.labels.memberSelected')),
      member: computed(() => ts('pages.administrations.groups.createEdit.labels.member')),
      members: computed(() => ts('pages.administrations.groups.createEdit.labels.members')),
      searchMembers: computed(() => ts('pages.administrations.groups.createEdit.labels.searchMembers')),
      noMembersFound: computed(() => ts('pages.administrations.groups.createEdit.labels.noMembersFound')),
      loadingMembers: computed(() => ts('pages.administrations.groups.createEdit.labels.loadingMembers')),
      noMembersMatch: computed(() => ts('pages.administrations.groups.createEdit.labels.noMembersMatch')),
      noMembersAvailable: computed(() => ts('pages.administrations.groups.createEdit.labels.noMembersAvailable')),
      membersOptional: computed(() => ts('pages.administrations.groups.createEdit.labels.membersOptional')),
      addMembers: computed(() => ts('pages.administrations.groups.createEdit.labels.addMembers')),
      addMoreMembers: computed(() => ts('pages.administrations.groups.createEdit.labels.addMoreMembers')),
      noMembersInGroup: computed(() => ts('pages.administrations.groups.createEdit.labels.noMembersInGroup')),
      noMembersSelectedYet: computed(() => ts('pages.administrations.groups.createEdit.labels.noMembersSelectedYet')),
      clickAddMembers: computed(() => ts('pages.administrations.groups.createEdit.labels.clickAddMembers')),
      removeMember: computed(() => ts('pages.administrations.groups.createEdit.labels.removeMember')),
      undoRemoval: computed(() => ts('pages.administrations.groups.createEdit.labels.undoRemoval')),
      badgeNew: computed(() => ts('pages.administrations.groups.createEdit.labels.badgeNew')),
      badgeRemoving: computed(() => ts('pages.administrations.groups.createEdit.labels.badgeRemoving')),
      // Role-related labels
      rolesSelected: computed(() => ts('pages.administrations.groups.createEdit.labels.rolesSelected')),
      roleSelected: computed(() => ts('pages.administrations.groups.createEdit.labels.roleSelected')),
      role: computed(() => ts('pages.administrations.groups.createEdit.labels.role')),
      roles: computed(() => ts('pages.administrations.groups.createEdit.labels.roles')),
      addRoles: computed(() => ts('pages.administrations.groups.createEdit.labels.addRoles')),
      removeRole: computed(() => ts('pages.administrations.groups.createEdit.labels.removeRole')),
      noRolesSelected: computed(() => ts('pages.administrations.groups.createEdit.labels.noRolesSelected')),
      clickAddRoles: computed(() => ts('pages.administrations.groups.createEdit.labels.clickAddRoles')),
      rolesRequired: computed(() => ts('pages.administrations.groups.createEdit.labels.rolesRequired')),
    },

    /**
     * Button translations
     */
    buttons: {
      back: computed(() => ts('pages.administrations.groups.createEdit.buttons.back')),
      next: computed(() => ts('pages.administrations.groups.createEdit.buttons.next')),
      createGroup: computed(() => ts('pages.administrations.groups.createEdit.buttons.createGroup')),
      updateGroup: computed(() => ts('pages.administrations.groups.createEdit.buttons.updateGroup')),
    },

    /**
     * Message translations (extended)
     */
    messages: {
      deletedSuccessfully: computed(() => ts('pages.administrations.groups.messages.deletedSuccessfully')),
      confirmDelete: (name: string) => ts('pages.administrations.groups.messages.confirmDelete', { name }),
      loadingGroup: computed(() => ts('pages.administrations.groups.createEdit.messages.loadingGroup')),
      completeAllSteps: computed(() => ts('pages.administrations.groups.createEdit.messages.completeAllSteps')),
      allFieldsRequired: computed(() => ts('pages.administrations.groups.createEdit.messages.allFieldsRequired')),
      currentStep: computed(() => ts('pages.administrations.groups.createEdit.messages.currentStep')),
    },

    /**
     * Review step translations
     */
    reviewStep: {
      subtitle: computed(() => ts('pages.administrations.groups.createEdit.reviewStep.subtitle')),
      successMessage: computed(() => ts('pages.administrations.groups.createEdit.reviewStep.successMessage')),
      successMessageEdit: computed(() => ts('pages.administrations.groups.createEdit.reviewStep.successMessageEdit')),
      sections: {
        basicInfo: computed(() => ts('pages.administrations.groups.createEdit.reviewStep.sections.basicInfo')),
        roles: computed(() => ts('pages.administrations.groups.createEdit.reviewStep.sections.roles')),
        members: computed(() => ts('pages.administrations.groups.createEdit.reviewStep.sections.members')),
      },
      fields: {
        name: computed(() => ts('pages.administrations.groups.createEdit.reviewStep.fields.name')),
        description: computed(() => ts('pages.administrations.groups.createEdit.reviewStep.fields.description')),
        enabled: computed(() => ts('pages.administrations.groups.createEdit.reviewStep.fields.enabled')),
        status: computed(() => ts('pages.administrations.groups.createEdit.reviewStep.fields.status')),
        isTemplate: computed(() => ts('pages.administrations.groups.createEdit.reviewStep.fields.isTemplate')),
        totalRoles: computed(() => ts('pages.administrations.groups.createEdit.reviewStep.fields.totalRoles')),
        totalMembers: computed(() => ts('pages.administrations.groups.createEdit.reviewStep.fields.totalMembers')),
      },
      values: {
        notProvided: computed(() => ts('pages.administrations.groups.createEdit.reviewStep.values.notProvided')),
        yes: computed(() => ts('pages.administrations.groups.createEdit.reviewStep.values.yes')),
        no: computed(() => ts('pages.administrations.groups.createEdit.reviewStep.values.no')),
        noMembers: computed(() => ts('pages.administrations.groups.createEdit.reviewStep.values.noMembers')),
        enabled: computed(() => ts('pages.administrations.groups.createEdit.reviewStep.values.enabled')),
        disabled: computed(() => ts('pages.administrations.groups.createEdit.reviewStep.values.disabled')),
      },
    },

    /**
     * CreateEdit notifications
     */
    createEditNotifications: {
      created: computed(() => ts('pages.administrations.groups.createEdit.notifications.created')),
      updated: computed(() => ts('pages.administrations.groups.createEdit.notifications.updated')),
      createFailed: computed(() => ts('pages.administrations.groups.createEdit.notifications.createFailed')),
      updateFailed: computed(() => ts('pages.administrations.groups.createEdit.notifications.updateFailed')),
      loadFailed: computed(() => ts('pages.administrations.groups.createEdit.notifications.loadFailed')),
      alreadyExists: computed(() => ts('pages.administrations.groups.createEdit.notifications.alreadyExists')),
      forbidden: computed(() => ts('pages.administrations.groups.createEdit.notifications.forbidden')),
    },
  };
}
