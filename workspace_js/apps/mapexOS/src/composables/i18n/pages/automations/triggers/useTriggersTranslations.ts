import type { FilterField } from '@components/drawers';
import type { DataRowColumn } from '@components/cards';
import type { ListHeaderMenuColumn } from '@components/headers';

import { computed } from 'vue';
import { useTS } from '@utils/translation';

export function useTriggersTranslations() {
  const ts = useTS({ capitalize: true });
  const tsTitle = useTS({ titleCase: true });
  const tsRaw = useTS({ capitalize: false });

  return {
    page: {
      title: computed(() => tsTitle('pages.automations.triggers.title')),
      description: computed(() => tsRaw('pages.automations.triggers.description')),
      listTitle: computed(() => tsTitle('pages.automations.triggers.listTitle')),
      button: {
        add: computed(() => ts('pages.automations.triggers.button.add')),
      },
    },
    triggers: {
      title: computed(() => tsTitle('pages.automations.triggers.addRulePage.title')),
      description: computed(() => tsRaw('pages.automations.triggers.addRulePage.description')),
      templateWarning: computed(() => tsRaw('pages.automations.triggers.addRulePage.templateWarning')),
      tabs: {
        individual: computed(() => ts('pages.automations.triggers.addRulePage.tabs.individual')),
        groups: computed(() => ts('pages.automations.triggers.addRulePage.tabs.groups')),
        external: computed(() => ts('pages.automations.triggers.addRulePage.tabs.external')),
      },
      individual: {
        title: computed(() => tsTitle('pages.automations.triggers.addRulePage.individual.title')),
        addTrigger: computed(() => ts('pages.automations.triggers.addRulePage.individual.addTrigger')),
        addTriggerTooltip: computed(() => ts('pages.automations.triggers.addRulePage.individual.addTriggerTooltip')),
        empty: {
          title: computed(() => tsRaw('pages.automations.triggers.addRulePage.individual.empty.title')),
          description: computed(() => tsRaw('pages.automations.triggers.addRulePage.individual.empty.description')),
        },
      },
      groups: {
        title: computed(() => tsTitle('pages.automations.triggers.addRulePage.groups.title')),
        addGroup: computed(() => ts('pages.automations.triggers.addRulePage.groups.addGroup')),
        addGroupTooltip: computed(() => ts('pages.automations.triggers.addRulePage.groups.addGroupTooltip')),
        addTriggerTooltip: computed(() => ts('pages.automations.triggers.addRulePage.groups.addTriggerTooltip')),
        deleteGroupTooltip: computed(() => ts('pages.automations.triggers.addRulePage.groups.deleteGroupTooltip')),
        groupNamePlaceholder: computed(() => ts('pages.automations.triggers.addRulePage.groups.groupNamePlaceholder')),
        empty: {
          title: computed(() => tsRaw('pages.automations.triggers.addRulePage.groups.empty.title')),
          description: computed(() => tsRaw('pages.automations.triggers.addRulePage.groups.empty.description')),
          noTriggersInGroup: computed(() => tsRaw('pages.automations.triggers.addRulePage.groups.empty.noTriggersInGroup')),
        },
      },
      external: {
        title: computed(() => tsTitle('pages.automations.triggers.addRulePage.external.title')),
        addTrigger: computed(() => ts('pages.automations.triggers.addRulePage.external.addTrigger')),
        addTriggerTooltip: computed(() => ts('pages.automations.triggers.addRulePage.external.addTriggerTooltip')),
        empty: {
          title: computed(() => tsRaw('pages.automations.triggers.addRulePage.external.empty.title')),
          description: computed(() => tsRaw('pages.automations.triggers.addRulePage.external.empty.description')),
        },
      },
      triggerVariableField: {
        labels: {
          variable: computed(() => tsRaw('pages.automations.triggers.addRulePage.triggerVariableField.labels.variable')),
          from: computed(() => tsRaw('pages.automations.triggers.addRulePage.triggerVariableField.labels.from')),
        },
      },
    },
    columns: computed(() => {
      return [
        {
          key: 'name',
          label: ts('pages.automations.triggers.columns.name'),
          type: 'text',
          visible: 'always',
          width: 250,
        },
        {
          key: 'triggerType',
          label: ts('pages.automations.triggers.columns.triggerType'),
          type: 'chip',
          visible: 'laptop',
          width: 120,
        },
        {
          key: 'category',
          label: ts('pages.automations.triggers.columns.category'),
          type: 'chip',
          visible: 'laptop',
          width: 150,
        },
        {
          key: 'status',
          label: ts('pages.automations.triggers.columns.status'),
          type: 'badge',
          visible: 'always',
          width: 100,
        },
      ] as DataRowColumn[];
    }),
    menuColumns: computed(() => {
      return [
        {
          key: 'triggerType',
          label: ts('pages.automations.triggers.columns.triggerType'),
          visible: true,
        },
        {
          key: 'category',
          label: ts('pages.automations.triggers.columns.category'),
          visible: true,
        },
      ] as ListHeaderMenuColumn[];
    }),
    menuLabels: {
      singular: computed(() => ts('pages.automations.triggers.menuLabels.singular')),
      plural: computed(() => ts('pages.automations.triggers.menuLabels.plural')),
    },
    filters: {
      label: computed(() => ts('pages.automations.triggers.filters.label')),
      searchPlaceholder: computed(() => ts('pages.automations.triggers.filters.searchPlaceholder')),
      allStatus: computed(() => ts('pages.automations.triggers.filters.allStatus')),
      advancedFilters: computed(() => ts('pages.automations.triggers.filters.advancedFilters')),
      pendingFilters: computed(() => ts('pages.automations.triggers.filters.pendingFilters')),
      clearAll: computed(() => ts('pages.automations.triggers.filters.clearAll')),
      includeChildren: computed(() => ts('pages.automations.triggers.filters.includeChildren')),
      includeChildrenOrgs: computed(() => ts('pages.automations.triggers.filters.includeChildrenOrgs')),
      category: computed(() => ts('pages.automations.triggers.filters.category')),
      triggerType: computed(() => ts('pages.automations.triggers.filters.triggerType')),

      options: {
        yes: computed(() => ts('pages.automations.triggers.filters.options.yes')),
        no: computed(() => ts('pages.automations.triggers.filters.options.no')),
        active: computed(() => ts('pages.automations.triggers.filters.options.active')),
        inactive: computed(() => ts('pages.automations.triggers.filters.options.inactive')),
        allCategories: computed(() => ts('pages.automations.triggers.filters.options.allCategories')),
        technical: computed(() => ts('pages.automations.triggers.filters.options.technical')),
        communication: computed(() => ts('pages.automations.triggers.filters.options.communication')),
        allTypes: computed(() => ts('pages.automations.triggers.filters.options.allTypes')),
      },
    },

    filterItems: computed((): FilterField[] => {
      return [
        {
          key: 'includeChildren',
          type: 'toggle',
          label: ts('pages.automations.triggers.filters.includeChildrenOrgs'),
          icon: 'account_tree',
          options: [
            { label: ts('pages.automations.triggers.filters.allStatus'), value: null },
            { label: ts('pages.automations.triggers.filters.options.yes'), value: true },
            { label: ts('pages.automations.triggers.filters.options.no'), value: false },
          ],
        },
        {
          key: 'category',
          type: 'select',
          label: ts('pages.automations.triggers.filters.category'),
          icon: 'list_alt',
          options: [
            { label: ts('pages.automations.triggers.filters.options.all'), value: null },
            { label: ts('pages.automations.triggers.filters.options.technical'), value: 'technical' },
            { label: ts('pages.automations.triggers.filters.options.communication'), value: 'communication' },
          ],
        },
        {
          key: 'triggerType',
          type: 'select',
          label: ts('pages.automations.triggers.filters.triggerType'),
          icon: 'category',
          options: [
            { label: ts('pages.automations.triggers.filters.options.all'), value: null },
            { label: 'HTTP', value: 'http' },
            { label: 'MQTT', value: 'mqtt' },
            { label: 'RabbitMQ', value: 'rabbitmq' },
            { label: 'NATS', value: 'nats' },
            { label: 'WebSocket', value: 'websocket' },
            { label: 'Email', value: 'email' },
            { label: 'Teams', value: 'teams' },
            { label: 'Slack', value: 'slack' },
          ],
        },
      ];
    }),

    empty: {
      title: computed(() => tsRaw('pages.automations.triggers.empty.title')),
      description: computed(() => tsRaw('pages.automations.triggers.empty.description')),
    },
    dialog: {
      confirmDelete: {
        title: computed(() => ts('pages.automations.triggers.dialog.confirmDelete.title')),
        message: (name: string) => tsRaw('pages.automations.triggers.dialog.confirmDelete.message', { name }),
      },
    },
    notifications: {
      deleteSuccess: computed(() => ts('pages.automations.triggers.notifications.deleteSuccess')),
    },
  };
}
