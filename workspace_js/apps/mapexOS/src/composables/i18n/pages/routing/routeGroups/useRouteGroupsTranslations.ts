import type { FilterField } from '@components/drawers';
import type { DataRowColumn } from '@components/cards';
import type { PageHeaderInfo } from '@components/headers';

import { computed } from 'vue';
import { useTS } from '@utils/translation';
import { useOrganizationStore } from '@stores/organization';

export function useRouteGroupsTranslations() {
  const ts = useTS({ capitalize: true });
  const tsTitle = useTS({ titleCase: true });
  const organizationStore = useOrganizationStore();

  return {
    pageHeader: {
      title: computed(() => tsTitle('pages.routing.routeGroups.pageHeader.title')),
      description: computed(() => ts('pages.routing.routeGroups.pageHeader.description')),
      button: computed(() => ts('pages.routing.routeGroups.pageHeader.button')),
      info: computed((): PageHeaderInfo => ({
        title: ts('pages.routing.routeGroups.pageHeader.info.title'),
        description: ts('pages.routing.routeGroups.pageHeader.info.description'),
        items: [
          {
            icon: 'alt_route',
            color: 'blue-6',
            title: ts('pages.routing.routeGroups.pageHeader.info.items.conditionalRouting.title'),
            text: ts('pages.routing.routeGroups.pageHeader.info.items.conditionalRouting.text'),
          },
          {
            icon: 'hub',
            color: 'green-6',
            title: ts('pages.routing.routeGroups.pageHeader.info.items.multipleDestinations.title'),
            text: ts('pages.routing.routeGroups.pageHeader.info.items.multipleDestinations.text'),
          },
          {
            icon: 'history',
            color: 'purple-6',
            title: ts('pages.routing.routeGroups.pageHeader.info.items.versionControl.title'),
            text: ts('pages.routing.routeGroups.pageHeader.info.items.versionControl.text'),
          },
          {
            icon: 'toggle_on',
            color: 'orange-6',
            title: ts('pages.routing.routeGroups.pageHeader.info.items.enableDisable.title'),
            text: ts('pages.routing.routeGroups.pageHeader.info.items.enableDisable.text'),
          },
          {
            icon: 'account_tree',
            color: 'indigo-6',
            title: ts('pages.routing.routeGroups.pageHeader.info.items.andOrLogic.title'),
            text: ts('pages.routing.routeGroups.pageHeader.info.items.andOrLogic.text'),
          },
        ],
        docsUrl: 'https://docs.mapexos.com/route-groups',
        docsLabel: ts('pages.routing.routeGroups.pageHeader.info.docsLabel'),
      })),
    },

    filters: {
      label: computed(() => ts('pages.routing.routeGroups.filters.label')),
      searchPlaceholder: computed(() => ts('pages.routing.routeGroups.filters.searchPlaceholder')),
      allStatus: computed(() => ts('pages.routing.routeGroups.filters.allStatus')),
      advancedFilters: computed(() => ts('pages.routing.routeGroups.filters.advancedFilters')),
      pendingFilters: computed(() => ts('pages.routing.routeGroups.filters.pendingFilters')),
      clearAll: computed(() => ts('pages.routing.routeGroups.filters.clearAll')),
      name: computed(() => ts('pages.routing.routeGroups.filters.name')),
      enabled: computed(() => ts('pages.routing.routeGroups.filters.enabled')),
      includeChildren: computed(() => ts('pages.routing.routeGroups.filters.includeChildren')),
      includeChildrenOrgs: computed(() => ts('pages.routing.routeGroups.filters.includeChildrenOrgs')),
      isTemplate: computed(() => ts('pages.routing.routeGroups.filters.isTemplate')),
      options: {
        yes: computed(() => ts('pages.routing.routeGroups.filters.options.yes')),
        no: computed(() => ts('pages.routing.routeGroups.filters.options.no')),
        active: computed(() => ts('pages.routing.routeGroups.filters.options.active')),
        inactive: computed(() => ts('pages.routing.routeGroups.filters.options.inactive')),
        templates: computed(() => ts('pages.routing.routeGroups.filters.options.templates')),
        local: computed(() => ts('pages.routing.routeGroups.filters.options.local')),
      },
    },

    advancedFilters: computed((): FilterField[] => {
      const baseFilters: FilterField[] = [
        {
          key: 'includeChildren',
          type: 'toggle',
          label: ts('pages.routing.routeGroups.filters.includeChildrenOrgs'),
          icon: 'account_tree',
          options: [
            { label: ts('pages.routing.routeGroups.filters.allStatus'), value: null },
            { label: ts('pages.routing.routeGroups.filters.options.yes'), value: true },
            { label: ts('pages.routing.routeGroups.filters.options.no'), value: false },
          ],
        },
      ];

      // Add isTemplate filter only for Customer and Site organizations
      if (organizationStore.isCustomer || organizationStore.isSite) {
        baseFilters.push({
          key: 'isTemplate',
          type: 'select',
          label: ts('pages.routing.routeGroups.filters.isTemplate'),
          icon: 'content_copy',
          options: [
            { label: ts('pages.routing.routeGroups.filters.allStatus'), value: null },
            { label: ts('pages.routing.routeGroups.filters.options.templates'), value: true },
            { label: ts('pages.routing.routeGroups.filters.options.local'), value: false }
          ],
        });
      }

      return baseFilters;
    }),

    listHeader: {
      title: computed(() => tsTitle('pages.routing.routeGroups.listHeader.title')),
      itemLabel: computed(() => ts('pages.routing.routeGroups.listHeader.itemLabel')),
      itemLabelPlural: computed(() => ts('pages.routing.routeGroups.listHeader.itemLabelPlural')),
    },

    menuColumns: {
      organization: computed(() => ts('pages.routing.routeGroups.menuColumns.organization')),
      routers: computed(() => ts('pages.routing.routeGroups.columns.routers')),
      templateSource: computed(() => ts('pages.routing.routeGroups.columns.templateSource')),
    },

    columns: computed((): DataRowColumn[] => [
      {
        key: 'icon',
        label: '',
        type: 'avatar',
        visible: 'always',
        width: 56,
        icon: 'route',
        color: (value: any, row: any) => row.enabled ? 'primary' : 'grey-5',
      },
      {
        key: 'name',
        label: ts('pages.routing.routeGroups.columns.name'),
        type: 'text',
        visible: 'always',
        width: 250,
        ellipsis: true,
        secondaryKey: 'description',
      },
      {
        key: 'organizationName',
        label: ts('pages.routing.routeGroups.columns.organization'),
        type: 'chip',
        visible: 'laptop',
        width: 180,
        ellipsis: true,
        color: 'indigo-6',
        icon: 'domain',
      },
      {
        key: 'routersCount',
        label: ts('pages.routing.routeGroups.columns.routers'),
        type: 'chip',
        visible: 'laptop',
        width: 120,
        color: 'purple-6',
        format: (value: any) => value ? `${value} routers` : '0 routers',
      },
      {
        key: 'isTemplate',
        label: ts('pages.routing.routeGroups.columns.templateSource'),
        type: 'chip',
        visible: 'laptop',
        width: 120,
        format: (value: any) => value
          ? ts('pages.routing.routeGroups.filters.options.templates').toUpperCase()
          : ts('pages.routing.routeGroups.filters.options.local').toUpperCase(),
        color: (value: any) => value ? 'orange-6' : 'green-6',
        icon: (value: any) => value ? 'content_copy' : 'folder',
      },
      {
        key: 'enabled',
        label: ts('pages.routing.routeGroups.columns.status'),
        type: 'badge',
        visible: 'always',
        width: 100,
        format: (value: any) => value ? ts('pages.routing.routeGroups.status.active').toUpperCase() : ts('pages.routing.routeGroups.status.inactive').toUpperCase(),
        color: (value: any) => value ? 'green-6' : 'red-6',
      },
    ]),

    empty: {
      title: computed(() => ts('pages.routing.routeGroups.empty.title')),
      description: computed(() => ts('pages.routing.routeGroups.empty.description')),
    },

    dialog: {
      confirmDelete: {
        title: computed(() => ts('pages.routing.routeGroups.dialog.confirmDelete.title')),
        message: (name: string) => ts('pages.routing.routeGroups.dialog.confirmDelete.message', { name }),
      },
    },

    notifications: {
      created: computed(() => ts('pages.routing.routeGroups.notifications.created')),
      updated: computed(() => ts('pages.routing.routeGroups.notifications.updated')),
      deleted: computed(() => ts('pages.routing.routeGroups.notifications.deleted')),
      deleteError: computed(() => ts('pages.routing.routeGroups.notifications.deleteError')),
      creationFailed: computed(() => ts('pages.routing.routeGroups.notifications.creationFailed')),
      updateFailed: computed(() => ts('pages.routing.routeGroups.notifications.updateFailed')),
      loadFailed: computed(() => ts('pages.routing.routeGroups.notifications.loadFailed')),
      alreadyExists: computed(() => ts('pages.routing.routeGroups.notifications.alreadyExists')),
      validationFailed: computed(() => ts('pages.routing.routeGroups.notifications.validationFailed')),
      networkError: computed(() => ts('pages.routing.routeGroups.notifications.networkError')),
      sharedEdit: computed(() => ts('pages.routing.routeGroups.notifications.sharedEdit')),
      sharedDelete: computed(() => ts('pages.routing.routeGroups.notifications.sharedDelete')),
    },

    errors: {
      apiNotInitialized: computed(() => ts('pages.routing.routeGroups.errors.apiNotInitialized')),
      idMissing: computed(() => ts('pages.routing.routeGroups.errors.idMissing')),
    },

    actions: {
      edit: computed(() => ts('pages.routing.routeGroups.actions.edit')),
      view: computed(() => ts('pages.routing.routeGroups.actions.view')),
      delete: computed(() => ts('pages.routing.routeGroups.actions.delete')),
    },

    status: {
      active: computed(() => ts('pages.routing.routeGroups.status.active')),
      inactive: computed(() => ts('pages.routing.routeGroups.status.inactive')),
    },

    drawer: {
      title: computed(() => ts('pages.routing.routeGroups.drawer.title')),
      close: computed(() => ts('pages.routing.routeGroups.drawer.close')),
      edit: computed(() => ts('pages.routing.routeGroups.drawer.edit')),
      duplicate: computed(() => ts('pages.routing.routeGroups.drawer.duplicate')),
      loading: computed(() => ts('pages.routing.routeGroups.drawer.loading')),
      error: computed(() => ts('pages.routing.routeGroups.drawer.error')),

      sections: {
        basicInfo: computed(() => ts('pages.routing.routeGroups.drawer.sections.basicInfo')),
        routers: computed(() => ts('pages.routing.routeGroups.drawer.sections.routers')),
        timestamps: computed(() => ts('pages.routing.routeGroups.drawer.sections.timestamps')),
      },

      fields: {
        name: computed(() => ts('pages.routing.routeGroups.drawer.fields.name')),
        enabled: computed(() => ts('pages.routing.routeGroups.drawer.fields.enabled')),
        description: computed(() => ts('pages.routing.routeGroups.drawer.fields.description')),
        created: computed(() => ts('pages.routing.routeGroups.drawer.fields.created')),
        updated: computed(() => ts('pages.routing.routeGroups.drawer.fields.updated')),
      },

      empty: {
        description: computed(() => ts('pages.routing.routeGroups.drawer.empty.description')),
        routers: computed(() => ts('pages.routing.routeGroups.drawer.empty.routers')),
      },
    },

    // Create/Edit Form translations (new pattern)
    createEdit: {
      page: {
        title: computed(() => ts('pages.routing.routeGroups.createEdit.page.title')),
        titleEdit: computed(() => ts('pages.routing.routeGroups.createEdit.page.titleEdit')),
        description: computed(() => ts('pages.routing.routeGroups.createEdit.page.description')),
        descriptionEdit: computed(() => ts('pages.routing.routeGroups.createEdit.page.descriptionEdit')),
        backButton: computed(() => ts('pages.routing.routeGroups.createEdit.page.backButton')),
        loading: computed(() => ts('pages.routing.routeGroups.createEdit.page.loading')),
      },

      stepper: {
        title: computed(() => ts('pages.routing.routeGroups.createEdit.stepper.title')),
        subtitle: computed(() => ts('pages.routing.routeGroups.createEdit.stepper.subtitle')),
        requiredInfo: computed(() => ts('pages.routing.routeGroups.createEdit.stepper.requiredInfo')),
        currentStep: computed(() => ts('pages.routing.routeGroups.createEdit.stepper.currentStep')),
      },

      navigation: {
        previous: computed(() => ts('pages.routing.routeGroups.createEdit.navigation.previous')),
        next: computed(() => ts('pages.routing.routeGroups.createEdit.navigation.next')),
        save: computed(() => ts('pages.routing.routeGroups.createEdit.navigation.save')),
        update: computed(() => ts('pages.routing.routeGroups.createEdit.navigation.update')),
        cancel: computed(() => ts('pages.routing.routeGroups.createEdit.navigation.cancel')),
      },

      steps: computed(() => {
        const stepsArray = ts('pages.routing.routeGroups.createEdit.steps');
        if (Array.isArray(stepsArray)) {
          return stepsArray;
        }
        // Fallback if translation not found
        return [
          { title: 'Basic Information', icon: 'mdi-information', description: 'Define route group name and status' },
          { title: 'Routers Configuration', icon: 'mdi-routes', description: 'Configure routing destinations and rules' },
          { title: 'Review', icon: 'mdi-check-circle', description: 'Review all configuration before saving' },
        ];
      }),

      basicInfoStep: {
        title: computed(() => ts('pages.routing.routeGroups.createEdit.basicInfoStep.title')),
        subtitle: computed(() => ts('pages.routing.routeGroups.createEdit.basicInfoStep.subtitle')),
        fields: {
          name: {
            label: computed(() => ts('pages.routing.routeGroups.createEdit.basicInfoStep.fields.name.label')),
            placeholder: computed(() => ts('pages.routing.routeGroups.createEdit.basicInfoStep.fields.name.placeholder')),
            hint: computed(() => ts('pages.routing.routeGroups.createEdit.basicInfoStep.fields.name.hint')),
            required: computed(() => ts('pages.routing.routeGroups.createEdit.basicInfoStep.fields.name.required')),
          },
          description: {
            label: computed(() => ts('pages.routing.routeGroups.createEdit.basicInfoStep.fields.description.label')),
            placeholder: computed(() => ts('pages.routing.routeGroups.createEdit.basicInfoStep.fields.description.placeholder')),
            hint: computed(() => ts('pages.routing.routeGroups.createEdit.basicInfoStep.fields.description.hint')),
          },
          enabled: {
            label: computed(() => ts('pages.routing.routeGroups.createEdit.basicInfoStep.fields.enabled.label')),
            hint: computed(() => ts('pages.routing.routeGroups.createEdit.basicInfoStep.fields.enabled.hint')),
            active: computed(() => ts('pages.routing.routeGroups.createEdit.basicInfoStep.fields.enabled.active')),
            inactive: computed(() => ts('pages.routing.routeGroups.createEdit.basicInfoStep.fields.enabled.inactive')),
          },
          isTemplate: {
            label: computed(() => ts('pages.routing.routeGroups.createEdit.basicInfoStep.fields.isTemplate.label')),
            hint: computed(() => ts('pages.routing.routeGroups.createEdit.basicInfoStep.fields.isTemplate.hint')),
          },
        },
      },

      routersStep: {
        title: computed(() => ts('pages.routing.routeGroups.createEdit.routersStep.title')),
        subtitle: computed(() => ts('pages.routing.routeGroups.createEdit.routersStep.subtitle')),
        addRouter: computed(() => ts('pages.routing.routeGroups.createEdit.routersStep.addRouter')),
        removeRouter: computed(() => ts('pages.routing.routeGroups.createEdit.routersStep.removeRouter')),
        noRouters: computed(() => ts('pages.routing.routeGroups.createEdit.routersStep.noRouters')),
        noRoutersWarning: computed(() => ts('pages.routing.routeGroups.createEdit.routersStep.noRoutersWarning')),
        noRoutersHint: computed(() => ts('pages.routing.routeGroups.createEdit.routersStep.noRoutersHint')),
        infoBanner: {
          title: computed(() => ts('pages.routing.routeGroups.createEdit.routersStep.infoBanner.title')),
          message: computed(() => ts('pages.routing.routeGroups.createEdit.routersStep.infoBanner.message')),
        },
        conditionalRouting: {
          label: computed(() => ts('pages.routing.routeGroups.createEdit.routersStep.conditionalRouting.label')),
          hint: computed(() => ts('pages.routing.routeGroups.createEdit.routersStep.conditionalRouting.hint')),
        },
        routerCard: {
          title: computed(() => ts('pages.routing.routeGroups.createEdit.routersStep.routerCard.title')),
          kind: {
            label: computed(() => ts('pages.routing.routeGroups.createEdit.routersStep.routerCard.kind.label')),
            placeholder: computed(() => ts('pages.routing.routeGroups.createEdit.routersStep.routerCard.kind.placeholder')),
            required: computed(() => ts('pages.routing.routeGroups.createEdit.routersStep.routerCard.kind.required')),
          },
          conditionalRouting: {
            label: computed(() => ts('pages.routing.routeGroups.createEdit.routersStep.routerCard.conditionalRouting.label')),
            hint: computed(() => ts('pages.routing.routeGroups.createEdit.routersStep.routerCard.conditionalRouting.hint')),
            enable: computed(() => ts('pages.routing.routeGroups.createEdit.routersStep.routerCard.conditionalRouting.enable')),
            disable: computed(() => ts('pages.routing.routeGroups.createEdit.routersStep.routerCard.conditionalRouting.disable')),
            policy: {
              label: computed(() => ts('pages.routing.routeGroups.createEdit.routersStep.routerCard.conditionalRouting.policy.label')),
            },
            addRule: computed(() => ts('pages.routing.routeGroups.createEdit.routersStep.routerCard.conditionalRouting.addRule')),
            removeRule: computed(() => ts('pages.routing.routeGroups.createEdit.routersStep.routerCard.conditionalRouting.removeRule')),
            noRules: computed(() => ts('pages.routing.routeGroups.createEdit.routersStep.routerCard.conditionalRouting.noRules')),
            rule: {
              field: {
                label: computed(() => ts('pages.routing.routeGroups.createEdit.routersStep.routerCard.conditionalRouting.rule.field.label')),
                placeholder: computed(() => ts('pages.routing.routeGroups.createEdit.routersStep.routerCard.conditionalRouting.rule.field.placeholder')),
                hint: computed(() => ts('pages.routing.routeGroups.createEdit.routersStep.routerCard.conditionalRouting.rule.field.hint')),
              },
              operator: {
                label: computed(() => ts('pages.routing.routeGroups.createEdit.routersStep.routerCard.conditionalRouting.rule.operator.label')),
                placeholder: computed(() => ts('pages.routing.routeGroups.createEdit.routersStep.routerCard.conditionalRouting.rule.operator.placeholder')),
              },
              value: {
                label: computed(() => ts('pages.routing.routeGroups.createEdit.routersStep.routerCard.conditionalRouting.rule.value.label')),
                placeholder: computed(() => ts('pages.routing.routeGroups.createEdit.routersStep.routerCard.conditionalRouting.rule.value.placeholder')),
                hint: computed(() => ts('pages.routing.routeGroups.createEdit.routersStep.routerCard.conditionalRouting.rule.value.hint')),
                arrayHint: computed(() => ts('pages.routing.routeGroups.createEdit.routersStep.routerCard.conditionalRouting.rule.value.arrayHint')),
              },
            },
          },
          businessRule: {
            label: computed(() => ts('pages.routing.routeGroups.createEdit.routersStep.routerCard.businessRule.label')),
            placeholder: computed(() => ts('pages.routing.routeGroups.createEdit.routersStep.routerCard.businessRule.placeholder')),
            required: computed(() => ts('pages.routing.routeGroups.createEdit.routersStep.routerCard.businessRule.required')),
          },
          lakeHouse: {
            label: computed(() => ts('pages.routing.routeGroups.createEdit.routersStep.routerCard.lakeHouse.label')),
            placeholder: computed(() => ts('pages.routing.routeGroups.createEdit.routersStep.routerCard.lakeHouse.placeholder')),
            required: computed(() => ts('pages.routing.routeGroups.createEdit.routersStep.routerCard.lakeHouse.required')),
          },
          notification: {
            label: computed(() => ts('pages.routing.routeGroups.createEdit.routersStep.routerCard.notification.label')),
            placeholder: computed(() => ts('pages.routing.routeGroups.createEdit.routersStep.routerCard.notification.placeholder')),
            required: computed(() => ts('pages.routing.routeGroups.createEdit.routersStep.routerCard.notification.required')),
          },
          workflow: {
            mode: computed(() => ts('pages.routing.routeGroups.createEdit.routersStep.routerCard.workflow.mode')),
            modeHint: computed(() => ts('pages.routing.routeGroups.createEdit.routersStep.routerCard.workflow.modeHint')),
            modeOptions: {
              newInstance: computed(() => ts('pages.routing.routeGroups.createEdit.routersStep.routerCard.workflow.modeOptions.newInstance')),
              signal: computed(() => ts('pages.routing.routeGroups.createEdit.routersStep.routerCard.workflow.modeOptions.signal')),
              signalOrStart: computed(() => ts('pages.routing.routeGroups.createEdit.routersStep.routerCard.workflow.modeOptions.signalOrStart')),
            },
            modeDescriptions: {
              newInstance: computed(() => ts('pages.routing.routeGroups.createEdit.routersStep.routerCard.workflow.modeDescriptions.newInstance')),
              signal: computed(() => ts('pages.routing.routeGroups.createEdit.routersStep.routerCard.workflow.modeDescriptions.signal')),
              signalOrStart: computed(() => ts('pages.routing.routeGroups.createEdit.routersStep.routerCard.workflow.modeDescriptions.signalOrStart')),
            },
            workflowId: computed(() => ts('pages.routing.routeGroups.createEdit.routersStep.routerCard.workflow.workflowId')),
            workflowIdHint: computed(() => ts('pages.routing.routeGroups.createEdit.routersStep.routerCard.workflow.workflowIdHint')),
            workflowUUID: computed(() => ts('pages.routing.routeGroups.createEdit.routersStep.routerCard.workflow.workflowUUID')),
            workflowUUIDHint: computed(() => ts('pages.routing.routeGroups.createEdit.routersStep.routerCard.workflow.workflowUUIDHint')),
            signalName: computed(() => ts('pages.routing.routeGroups.createEdit.routersStep.routerCard.workflow.signalName')),
            signalNameHint: computed(() => ts('pages.routing.routeGroups.createEdit.routersStep.routerCard.workflow.signalNameHint')),
            metadata: computed(() => ts('pages.routing.routeGroups.createEdit.routersStep.routerCard.workflow.metadata')),
            metadataHint: computed(() => ts('pages.routing.routeGroups.createEdit.routersStep.routerCard.workflow.metadataHint')),
          },
        },
      },

      reviewStep: {
        title: computed(() => ts('pages.routing.routeGroups.createEdit.reviewStep.title')),
        subtitle: computed(() => ts('pages.routing.routeGroups.createEdit.reviewStep.subtitle')),
        successMessage: computed(() => ts('pages.routing.routeGroups.createEdit.reviewStep.successMessage')),
        sections: {
          basicInfo: computed(() => ts('pages.routing.routeGroups.createEdit.reviewStep.sections.basicInfo')),
          routers: computed(() => ts('pages.routing.routeGroups.createEdit.reviewStep.sections.routers')),
        },
        fields: {
          name: computed(() => ts('pages.routing.routeGroups.createEdit.reviewStep.fields.name')),
          description: computed(() => ts('pages.routing.routeGroups.createEdit.reviewStep.fields.description')),
          status: computed(() => ts('pages.routing.routeGroups.createEdit.reviewStep.fields.status')),
          templateSource: computed(() => ts('pages.routing.routeGroups.createEdit.reviewStep.fields.templateSource')),
          routerType: computed(() => ts('pages.routing.routeGroups.createEdit.reviewStep.fields.routerType')),
          conditionalRouting: computed(() => ts('pages.routing.routeGroups.createEdit.reviewStep.fields.conditionalRouting')),
          matchPolicy: computed(() => ts('pages.routing.routeGroups.createEdit.reviewStep.fields.matchPolicy')),
          matchRules: computed(() => ts('pages.routing.routeGroups.createEdit.reviewStep.fields.matchRules')),
          lakeHouse: computed(() => ts('pages.routing.routeGroups.createEdit.reviewStep.fields.lakeHouse')),
          notification: computed(() => ts('pages.routing.routeGroups.createEdit.reviewStep.fields.notification')),
        },
        values: {
          enabled: computed(() => ts('pages.routing.routeGroups.createEdit.reviewStep.values.enabled')),
          disabled: computed(() => ts('pages.routing.routeGroups.createEdit.reviewStep.values.disabled')),
          shared: computed(() => ts('pages.routing.routeGroups.createEdit.reviewStep.values.shared')),
          local: computed(() => ts('pages.routing.routeGroups.createEdit.reviewStep.values.local')),
          none: computed(() => ts('pages.routing.routeGroups.createEdit.reviewStep.values.none')),
          notConfigured: computed(() => ts('pages.routing.routeGroups.createEdit.reviewStep.values.notConfigured')),
        },
      },
    },

    routerKinds: {
      trigger: {
        label: computed(() => ts('pages.routing.routeGroups.routerKinds.trigger.label')),
        description: computed(() => ts('pages.routing.routeGroups.routerKinds.trigger.description')),
      },
      lake_house: {
        label: computed(() => ts('pages.routing.routeGroups.routerKinds.lake_house.label')),
        description: computed(() => ts('pages.routing.routeGroups.routerKinds.lake_house.description')),
      },
      notification: {
        label: computed(() => ts('pages.routing.routeGroups.routerKinds.notification.label')),
        description: computed(() => ts('pages.routing.routeGroups.routerKinds.notification.description')),
      },
      save_event: {
        label: computed(() => ts('pages.routing.routeGroups.routerKinds.save_event.label')),
        description: computed(() => ts('pages.routing.routeGroups.routerKinds.save_event.description')),
      },
      workflow: {
        label: computed(() => ts('pages.routing.routeGroups.routerKinds.workflow.label')),
        description: computed(() => ts('pages.routing.routeGroups.routerKinds.workflow.description')),
      },
    },

    matchPolicies: {
      all: {
        label: computed(() => ts('pages.routing.routeGroups.matchPolicies.all.label')),
        description: computed(() => ts('pages.routing.routeGroups.matchPolicies.all.description')),
      },
      any: {
        label: computed(() => ts('pages.routing.routeGroups.matchPolicies.any.label')),
        description: computed(() => ts('pages.routing.routeGroups.matchPolicies.any.description')),
      },
    },

    matchOperators: {
      eq: computed(() => ts('pages.routing.routeGroups.matchOperators.eq')),
      neq: computed(() => ts('pages.routing.routeGroups.matchOperators.neq')),
      gt: computed(() => ts('pages.routing.routeGroups.matchOperators.gt')),
      gte: computed(() => ts('pages.routing.routeGroups.matchOperators.gte')),
      lt: computed(() => ts('pages.routing.routeGroups.matchOperators.lt')),
      lte: computed(() => ts('pages.routing.routeGroups.matchOperators.lte')),
      in: computed(() => ts('pages.routing.routeGroups.matchOperators.in')),
      nin: computed(() => ts('pages.routing.routeGroups.matchOperators.nin')),
    },

    // Options computed (deprecated - use createEdit.steps instead)
    translatedSteps: computed(() => {
      const stepsArray = ts('pages.routing.routeGroups.createEdit.steps');
      if (Array.isArray(stepsArray)) {
        return stepsArray;
      }
      // Fallback
      return [
        { title: 'Basic Information', icon: 'mdi-information', description: 'Define route group name and status' },
        { title: 'Routers Configuration', icon: 'mdi-routes', description: 'Configure routing destinations and rules' },
        { title: 'Review', icon: 'mdi-check-circle', description: 'Review all configuration before saving' },
      ];
    }),

    routerKindOptions: computed(() => [
      {
        label: ts('pages.routing.routeGroups.routerKinds.trigger.label'),
        value: 'trigger',
        icon: 'flash_on',
        color: 'amber-8',
        description: ts('pages.routing.routeGroups.routerKinds.trigger.description'),
      },
      {
        label: ts('pages.routing.routeGroups.routerKinds.lake_house.label'),
        value: 'lake_house',
        icon: 'storage',
        color: 'purple-6',
        description: ts('pages.routing.routeGroups.routerKinds.lake_house.description'),
      },
      {
        label: ts('pages.routing.routeGroups.routerKinds.notification.label'),
        value: 'notification',
        icon: 'notifications',
        color: 'orange-6',
        description: ts('pages.routing.routeGroups.routerKinds.notification.description'),
      },
      {
        label: ts('pages.routing.routeGroups.routerKinds.save_event.label'),
        value: 'save_event',
        icon: 'save',
        color: 'green-6',
        description: ts('pages.routing.routeGroups.routerKinds.save_event.description'),
      },
      {
        label: ts('pages.routing.routeGroups.routerKinds.workflow.label'),
        value: 'workflow',
        icon: 'account_tree',
        color: 'teal-6',
        description: ts('pages.routing.routeGroups.routerKinds.workflow.description'),
      },
    ]),

    matchPolicyOptions: computed(() => [
      {
        label: ts('pages.routing.routeGroups.matchPolicies.all.label'),
        value: 'all',
        description: ts('pages.routing.routeGroups.matchPolicies.all.description'),
      },
      {
        label: ts('pages.routing.routeGroups.matchPolicies.any.label'),
        value: 'any',
        description: ts('pages.routing.routeGroups.matchPolicies.any.description'),
      },
    ]),

    matchOperatorOptions: computed(() => [
      { label: ts('pages.routing.routeGroups.matchOperators.eq'), value: 'eq', description: 'Equal (=)' },
      { label: ts('pages.routing.routeGroups.matchOperators.neq'), value: 'neq', description: 'Not Equal (!=)' },
      { label: ts('pages.routing.routeGroups.matchOperators.gt'), value: 'gt', description: 'Greater Than (>)' },
      { label: ts('pages.routing.routeGroups.matchOperators.gte'), value: 'gte', description: 'Greater or Equal (>=)' },
      { label: ts('pages.routing.routeGroups.matchOperators.lt'), value: 'lt', description: 'Less Than (<)' },
      { label: ts('pages.routing.routeGroups.matchOperators.lte'), value: 'lte', description: 'Less or Equal (<=)' },
      { label: ts('pages.routing.routeGroups.matchOperators.in'), value: 'in', description: 'In Array' },
      { label: ts('pages.routing.routeGroups.matchOperators.nin'), value: 'nin', description: 'Not In Array' },
    ]),

    statusOptions: computed(() => [
      {
        label: ts('pages.routing.routeGroups.createEdit.basicInfoStep.fields.enabled.active'),
        value: true,
      },
      {
        label: ts('pages.routing.routeGroups.createEdit.basicInfoStep.fields.enabled.inactive'),
        value: false,
      },
    ]),

    // Tour translations
    tour: {
      welcome: {
        title: computed(() => tsTitle('pages.routing.routeGroups.tour.welcome.title')),
        description: computed(() => ts('pages.routing.routeGroups.tour.welcome.description')),
      },
      stepperOverview: {
        title: computed(() => tsTitle('pages.routing.routeGroups.tour.stepperOverview.title')),
        description: computed(() => ts('pages.routing.routeGroups.tour.stepperOverview.description')),
      },
      step1Overview: {
        title: computed(() => tsTitle('pages.routing.routeGroups.tour.step1Overview.title')),
        description: computed(() => ts('pages.routing.routeGroups.tour.step1Overview.description')),
      },
      fieldName: {
        title: computed(() => tsTitle('pages.routing.routeGroups.tour.fieldName.title')),
        description: computed(() => ts('pages.routing.routeGroups.tour.fieldName.description')),
      },
      fieldStatus: {
        title: computed(() => tsTitle('pages.routing.routeGroups.tour.fieldStatus.title')),
        description: computed(() => ts('pages.routing.routeGroups.tour.fieldStatus.description')),
      },
      step2Overview: {
        title: computed(() => tsTitle('pages.routing.routeGroups.tour.step2Overview.title')),
        description: computed(() => ts('pages.routing.routeGroups.tour.step2Overview.description')),
      },
      addRouter: {
        title: computed(() => tsTitle('pages.routing.routeGroups.tour.addRouter.title')),
        description: computed(() => ts('pages.routing.routeGroups.tour.addRouter.description')),
      },
      routerCard: {
        title: computed(() => tsTitle('pages.routing.routeGroups.tour.routerCard.title')),
        description: computed(() => ts('pages.routing.routeGroups.tour.routerCard.description')),
      },
      routerKind: {
        title: computed(() => tsTitle('pages.routing.routeGroups.tour.routerKind.title')),
        description: computed(() => ts('pages.routing.routeGroups.tour.routerKind.description')),
      },
      conditionalRouting: {
        title: computed(() => tsTitle('pages.routing.routeGroups.tour.conditionalRouting.title')),
        description: computed(() => ts('pages.routing.routeGroups.tour.conditionalRouting.description')),
      },
      step3Overview: {
        title: computed(() => tsTitle('pages.routing.routeGroups.tour.step3Overview.title')),
        description: computed(() => ts('pages.routing.routeGroups.tour.step3Overview.description')),
      },
      saveButton: {
        title: computed(() => tsTitle('pages.routing.routeGroups.tour.saveButton.title')),
        description: computed(() => ts('pages.routing.routeGroups.tour.saveButton.description')),
      },
    },
  };
}
