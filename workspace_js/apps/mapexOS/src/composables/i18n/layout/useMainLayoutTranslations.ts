import { computed } from 'vue';
import { useTS } from '@utils/translation';

/**
 * MainLayout translations
 *
 * Structure mirrors:
 * - File: src/layouts/main/MainLayout.vue
 * - JSON: src/i18n/{locale}/layout/mainLayout.json
 * - Composable: src/composables/i18n/layout/useMainLayoutTranslations.ts
 *
 * Provides all translations for the main layout including:
 * - Language selector messages
 * - Customer selector messages
 * - User menu items
 * - Notification titles
 *
 * @example
 * ```ts
 * // In MainLayout.vue
 * const { languageSelector, customerSelector, userMenu } = useMainLayoutTranslations();
 *
 * <div>{{ languageSelector.changed.value }} {{ languageName }}</div>
 * <q-btn :label="userMenu.logout.value" />
 * ```
 */
export function useMainLayoutTranslations() {
  const ts = useTS({ capitalize: true });

  return {
    /**
     * Language selector translations
     * Mirrors: layout.mainLayout.languageSelector
     */
    languageSelector: {
      changed: computed(() => ts('layout.mainLayout.languageSelector.changed')),
      languages: {
        english: computed(() => ts('layout.mainLayout.languageSelector.languages.english')),
        portuguese: computed(() => ts('layout.mainLayout.languageSelector.languages.portuguese')),
      },
    },

    /**
     * Customer selector translations
     * Mirrors: layout.mainLayout.customerSelector
     */
    customerSelector: {
      select: computed(() => ts('layout.mainLayout.customerSelector.select')),
      loggedOut: computed(() => ts('layout.mainLayout.customerSelector.loggedOut')),
      loggedIn: computed(() => ts('layout.mainLayout.customerSelector.loggedIn')),
      exit: computed(() => ts('layout.mainLayout.customerSelector.exit')),
    },

    /**
     * Organization tree translations
     * Mirrors: layout.mainLayout.orgTree
     */
    orgTree: {
      tooltip: computed(() => ts('layout.mainLayout.orgTree.tooltip')),
    },

    /**
     * Organization indicator translations
     * Mirrors: layout.mainLayout.orgIndicator
     */
    orgIndicator: {
      label: computed(() => ts('layout.mainLayout.orgIndicator.label')),
      tooltip: computed(() => ts('layout.mainLayout.orgIndicator.tooltip')),
    },

    /**
     * Notifications translations
     * Mirrors: layout.mainLayout.notifications
     */
    notifications: {
      title: computed(() => ts('layout.mainLayout.notifications.title')),
    },

    /**
     * User menu translations
     * Mirrors: layout.mainLayout.userMenu
     */
    userMenu: {
      profile: computed(() => ts('layout.mainLayout.userMenu.profile')),
      settings: computed(() => ts('layout.mainLayout.userMenu.settings')),
      docs: computed(() => ts('layout.mainLayout.userMenu.docs')),
      startTour: computed(() => ts('layout.mainLayout.userMenu.startTour')),
      logout: computed(() => ts('layout.mainLayout.userMenu.logout')),
    },

    /**
     * Sidebar menu translations
     * Mirrors: layout.mainLayout.menu
     */
    menu: {
      dashboard: computed(() => ts('layout.mainLayout.menu.dashboard')),
      assets: computed(() => ts('layout.mainLayout.menu.assets')),
      assetsItems: computed(() => ts('layout.mainLayout.menu.assetsItems')),
      assetsTemplate: computed(() => ts('layout.mainLayout.menu.assetsTemplate')),
      data: computed(() => ts('layout.mainLayout.menu.data')),
      http: computed(() => ts('layout.mainLayout.menu.http')),
      automation: computed(() => ts('layout.mainLayout.menu.automation')),
      workflows: computed(() => ts('layout.mainLayout.menu.workflows')),
      workflowInstances: computed(() => ts('layout.mainLayout.menu.workflowInstances')),
      triggers: computed(() => ts('layout.mainLayout.menu.triggers')),
      routing: computed(() => ts('layout.mainLayout.menu.routing')),
      routeGroups: computed(() => ts('layout.mainLayout.menu.routeGroups')),
      logsAndExecutions: computed(() => ts('layout.mainLayout.menu.logsAndExecutions')),
      eventTracer: computed(() => ts('layout.mainLayout.menu.eventTracer')),
      rawEvents: computed(() => ts('layout.mainLayout.menu.rawEvents')),
      assetConnectivity: computed(() => ts('layout.mainLayout.menu.assetConnectivity')),
      jsExecutor: computed(() => ts('layout.mainLayout.menu.jsExecutor')),
      router: computed(() => ts('layout.mainLayout.menu.router')),
      triggerLogs: computed(() => ts('layout.mainLayout.menu.triggerLogs')),
      workflowExecutions: computed(() => ts('layout.mainLayout.menu.workflowExecutions')),
      dlq: computed(() => ts('layout.mainLayout.menu.dlq')),
      events: computed(() => ts('layout.mainLayout.menu.events')),
      eventStore: computed(() => ts('layout.mainLayout.menu.eventStore')),
      administration: computed(() => ts('layout.mainLayout.menu.administration')),
      customers: computed(() => ts('layout.mainLayout.menu.customers')),
      users: computed(() => ts('layout.mainLayout.menu.users')),
      roles: computed(() => ts('layout.mainLayout.menu.roles')),
      groups: computed(() => ts('layout.mainLayout.menu.groups')),
      accessAudit: computed(() => ts('layout.mainLayout.menu.accessAudit')),
      lists: computed(() => ts('layout.mainLayout.menu.lists')),
      settings: computed(() => ts('layout.mainLayout.menu.settings')),
    },

    /**
     * Version info
     * Mirrors: layout.mainLayout.version
     */
    version: computed(() => ts('layout.mainLayout.version')),

    /**
     * Breadcrumbs translations
     * Mirrors: layout.mainLayout.breadcrumbs
     */
    breadcrumbs: {
      home: computed(() => ts('layout.mainLayout.breadcrumbs.home')),
    },
  };
}
