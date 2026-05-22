/**
 * Portuguese (Brazil) - pt-BR Locale
 *
 * Organized structure:
 * - common: Shared translations (buttons, messages, labels, validation)
 * - components: Component-specific translations
 * - pages: Page-specific translations organized by module
 */

import common from './common.json';

// Layout
import mainLayout from './layout/mainLayout.json';

// Components
import headersComponent from './components/headers.json';
import cardsComponent from './components/cards.json';
import filtersComponent from './components/filters.json';
import organizationTreeDrawerComponent from './components/organizationTreeDrawer.json';
import advancedFiltersDrawerComponent from './components/drawers/advancedFiltersDrawer.json';
import dynamicFiltersDrawerComponent from './components/drawers/dynamicFiltersDrawer.json';
import assetTemplateSelectorComponent from './components/selectors/assetTemplateSelector.json';
import routeGroupSelectorComponent from './components/selectors/routeGroupSelector.json';
import assetClassificationSelectorComponent from './components/forms/assetClassificationSelector.json';
import eventFieldInputComponent from './components/forms/eventFieldInput.json';
import fieldSourceSelectorComponent from './components/forms/fieldSourceSelector.json';
import standardizedPayloadHelpDialog from './components/dialogs/standardizedPayloadHelp.json';
import genericSelectorDialog from './components/dialogs/genericSelector.json';
import triggerSelectorDialog from './components/dialogs/triggerSelector.json';
import workflowSelectorDialog from './components/dialogs/workflowSelector.json';
import scriptEditorDialog from './components/dialogs/scriptEditor.json';

// Composables
import onboardingComposable from './composables/onboarding.json';

// Pages - Administrations
import settingsPage from './pages/administrations/settings.json';
import groupsPage from './pages/administrations/groups.json';
import rolesPage from './pages/administrations/roles.json';
import listsPage from './pages/administrations/lists.json';
import usersPage from './pages/administrations/users.json';
import addUserPage from './pages/administrations/addUser.json';
import userProfilePage from './pages/administrations/userProfile.json';
import customersPage from './pages/administrations/customers.json';
import addCustomerPage from './pages/administrations/addCustomer.json';
import accessAuditPage from './pages/administrations/accessAudit.json';

// Pages - Assets
import assetsPage from './pages/assets/assets.json';
import assetTemplatesPage from './pages/assets/assetTemplates.json';
import addAssetPage from './pages/assets/addAsset.json';
import addAssetTemplatePage from './pages/assets/addAssetTemplate.json';

// Pages - Dashboards
import dashboardAdmPage from './pages/dashboards/dashboardAdm.json';

// Pages - DataSources
import httpDataSourcesPage from './pages/datasources/http.json';

// Pages - Automations
import triggersPage from './pages/automations/triggers.json';
import createEditTriggerPage from './pages/automations/createEditTrigger.json';
import createEditWorkflowPage from './pages/automations/createEditWorkflow.json';
import workflowListPage from './pages/automations/workflowList.json';
import workflowInstanceListPage from './pages/automations/workflowInstanceList.json';
import createEditWorkflowInstancePage from './pages/automations/createEditWorkflowInstance.json';

// Pages - Notifications
import notificationsPage from './pages/notifications/notifications.json';
import addNotificationPage from './pages/notifications/addNotification.json';

// Pages - LakeHouse
import lakeHousePage from './pages/lakeHouse/lakeHouse.json';
import addLakeHousePage from './pages/lakeHouse/addLakeHouse.json';

// Pages - Routing
import routeGroupsPage from './pages/routing/routeGroups.json';

// Pages - Events
import eventStorePage from './pages/events/eventStore.json';

// Pages - Logs
import assetRawLogsPage from './pages/logs/assetRawLogsPage.json';
import assetConnectivityLogsPage from './pages/logs/assetConnectivityLogsPage.json';
import jsExecLogsPage from './pages/logs/jsExecLogsPage.json';
import routerLogsPage from './pages/logs/routerLogsPage.json';
import eventTracerPage from './pages/logs/eventTracerPage.json';
import triggerLogsPage from './pages/logs/triggerLogsPage.json';
import workflowExecutionsPage from './pages/logs/workflowExecutionsPage.json';
import dlqLogsPage from './pages/logs/dlqLogsPage.json';

// Pages - Auth & Errors
import loginPage from './pages/login.json';
import changePasswordPage from './pages/changePassword.json';
import errorPage from './pages/error.json';
import noOrganizationPage from './pages/errors/noOrganization.json';

export default {
  // Common translations (actions, messages, labels, status, pagination, validation)
  common,

  // Layout translations
  layout: {
    mainLayout,
  },

  // Composable translations
  composables: {
    onboarding: onboardingComposable,
  },

  // Component translations
  components: {
    headers: headersComponent,
    cards: cardsComponent,
    filters: filtersComponent,
    organizationTreeDrawer: organizationTreeDrawerComponent,
    selectors: {
      assetTemplateSelector: assetTemplateSelectorComponent,
      routeGroupSelector: routeGroupSelectorComponent,
    },
    forms: {
      assetClassificationSelector: assetClassificationSelectorComponent,
      eventFieldInput: eventFieldInputComponent,
      fieldSourceSelector: fieldSourceSelectorComponent,
    },
    dialogs: {
      standardizedPayloadHelp: standardizedPayloadHelpDialog,
      genericSelector: genericSelectorDialog,
      triggerSelector: triggerSelectorDialog,
      workflowSelector: workflowSelectorDialog,
      scriptEditor: scriptEditorDialog,
    },
    drawers: {
      advancedFiltersDrawer: advancedFiltersDrawerComponent,
      dynamicFiltersDrawer: dynamicFiltersDrawerComponent,
    },
  },

  // Page translations
  pages: {
    administrations: {
      settings: settingsPage,
      groups: groupsPage,
      roles: rolesPage,
      lists: listsPage,
      users: usersPage,
      addUser: addUserPage,
      userProfile: userProfilePage,
      customers: customersPage,
      addCustomer: addCustomerPage,
      accessAudit: accessAuditPage,
    },
    assets: {
      assets: assetsPage,
      assetTemplates: assetTemplatesPage,
      addAsset: addAssetPage,
      addAssetTemplate: addAssetTemplatePage,
    },
    dashboards: {
      dashboardAdm: dashboardAdmPage,
    },
    datasources: {
      http: httpDataSourcesPage,
    },
    automations: {
      triggers: triggersPage,
      createEditTrigger: createEditTriggerPage,
      createEditWorkflow: createEditWorkflowPage,
      workflowList: workflowListPage,
      workflowInstanceList: workflowInstanceListPage,
      createEditWorkflowInstance: createEditWorkflowInstancePage,
    },
    notifications: {
      notifications: notificationsPage,
      addNotification: addNotificationPage,
    },
    lakeHouse: {
      lakeHouse: lakeHousePage,
      addLakeHouse: addLakeHousePage,
    },
    routing: {
      routeGroups: routeGroupsPage,
    },
    events: {
      eventStore: eventStorePage,
    },
    logs: {
      assetRawLogsPage: assetRawLogsPage,
      assetConnectivityLogsPage: assetConnectivityLogsPage,
      jsExecLogsPage: jsExecLogsPage,
      routerLogsPage: routerLogsPage,
      eventTracerPage: eventTracerPage,
      triggerLogsPage: triggerLogsPage,
      workflowExecutionsPage: workflowExecutionsPage,
      dlqLogsPage: dlqLogsPage,
    },
    login: loginPage,
    changePassword: changePasswordPage,
    error: errorPage,
    errors: {
      noOrganization: noOrganizationPage,
    },
  },
};
