import type { MenuItem } from '../interfaces';
import type { ComputedRef } from 'vue';
import { PERMISSIONS } from '@src/constants/permissions.constant';

/** Menu translation keys provided by useMainLayoutTranslations().menu */
interface MenuTranslations {
  dashboard: ComputedRef<string>;
  assets: ComputedRef<string>;
  assetsItems: ComputedRef<string>;
  assetsTemplate: ComputedRef<string>;
  data: ComputedRef<string>;
  http: ComputedRef<string>;
  automation: ComputedRef<string>;
  workflows: ComputedRef<string>;
  workflowInstances: ComputedRef<string>;
  triggers: ComputedRef<string>;
  routing: ComputedRef<string>;
  routeGroups: ComputedRef<string>;
  logsAndExecutions: ComputedRef<string>;
  eventTracer: ComputedRef<string>;
  rawEvents: ComputedRef<string>;
  assetConnectivity: ComputedRef<string>;
  jsExecutor: ComputedRef<string>;
  router: ComputedRef<string>;
  triggerLogs: ComputedRef<string>;
  workflowExecutions: ComputedRef<string>;
  dlq: ComputedRef<string>;
  events: ComputedRef<string>;
  eventStore: ComputedRef<string>;
  administration: ComputedRef<string>;
  customers: ComputedRef<string>;
  users: ComputedRef<string>;
  roles: ComputedRef<string>;
  groups: ComputedRef<string>;
  accessAudit: ComputedRef<string>;
  lists: ComputedRef<string>;
  settings: ComputedRef<string>;
}

/**
 * Build the sidebar menu list with translated labels
 *
 * @param {MenuTranslations} m - Menu translations from useMainLayoutTranslations().menu
 * @returns {MenuItem[]} Translated menu items
 */
export function buildMenuList(m: MenuTranslations): MenuItem[] {
  return [
    { icon: 'dashboard', label: m.dashboard.value, to: '/home' },

    {
      icon: 'device_hub',
      label: m.assets.value,
      permissions: [PERMISSIONS.ASSETS.LIST, PERMISSIONS.ASSET_TEMPLATES.LIST],
      children: [
        { label: m.assetsItems.value, to: '/assets', icon: 'devices', permissions: [PERMISSIONS.ASSETS.LIST] },
        { label: m.assetsTemplate.value, to: '/assets_template', icon: 'memory', permissions: [PERMISSIONS.ASSET_TEMPLATES.LIST] },
      ],
    },

    {
      icon: 'settings_input_antenna',
      label: m.data.value,
      permissions: [PERMISSIONS.DATASOURCES.LIST, PERMISSIONS.RETENTION.LIST],
      children: [
        { label: m.http.value, to: '/data_sources/http', icon: 'http', permissions: [PERMISSIONS.DATASOURCES.LIST] },
      ],
    },

    {
      icon: 'psychology',
      label: m.automation.value,
      permissions: [PERMISSIONS.WORKFLOWS.LIST, PERMISSIONS.WORKFLOWS.INSTANCES.LIST, PERMISSIONS.TRIGGERS.LIST],
      children: [
        { label: m.workflows.value, to: '/workflows', icon: 'account_tree', permissions: [PERMISSIONS.WORKFLOWS.LIST] },
        { label: m.workflowInstances.value, to: '/workflow_instances', icon: 'play_circle', permissions: [PERMISSIONS.WORKFLOWS.INSTANCES.LIST] },
        { separator: true },
        { label: m.triggers.value, to: '/triggers', icon: 'flash_on', permissions: [PERMISSIONS.TRIGGERS.LIST] },
      ],
    },

    {
      icon: 'route',
      label: m.routing.value,
      permissions: [PERMISSIONS.ROUTE_GROUPS.LIST],
      children: [
        { label: m.routeGroups.value, to: '/routing/route_groups', icon: 'alt_route', permissions: [PERMISSIONS.ROUTE_GROUPS.LIST] },
      ],
    },

    {
      icon: 'list_alt',
      label: m.logsAndExecutions.value,
      permissions: [
        PERMISSIONS.EVENTS.LIST,
        PERMISSIONS.EVENTS.RAW.LIST,
        PERMISSIONS.EVENTS.ASSET_STATUS.LIST,
        PERMISSIONS.EVENTS.JS_EXECUTOR.LIST,
        PERMISSIONS.EVENTS.ROUTER.LIST,
        PERMISSIONS.EVENTS.TRIGGER.LIST,
        PERMISSIONS.EVENTS.DLQ.LIST,
      ],
      children: [
        { label: m.eventTracer.value, to: '/logs/event_tracer', icon: 'account_tree', permissions: [PERMISSIONS.EVENTS.LIST] },
        { separator: true },
        { label: m.rawEvents.value, to: '/logs/raw_logs', icon: 'terminal', permissions: [PERMISSIONS.EVENTS.RAW.LIST] },
        { label: m.assetConnectivity.value, to: '/logs/connectivity', icon: 'wifi', permissions: [PERMISSIONS.EVENTS.ASSET_STATUS.LIST] },
        { label: m.jsExecutor.value, to: '/logs/js_exec_logs', icon: 'code', permissions: [PERMISSIONS.EVENTS.JS_EXECUTOR.LIST] },
        { label: m.router.value, to: '/logs/router_logs', icon: 'route', permissions: [PERMISSIONS.EVENTS.ROUTER.LIST] },
        { separator: true },
        { label: m.triggerLogs.value, to: '/logs/triggers_log', icon: 'flash_on', permissions: [PERMISSIONS.EVENTS.TRIGGER.LIST] },
        { label: m.workflowExecutions.value, to: '/logs/workflow_executions', icon: 'play_circle', permissions: [PERMISSIONS.WORKFLOWS.INSTANCES.LIST] },
        { separator: true },
        { label: m.dlq.value, to: '/logs/dlq', icon: 'report_problem', permissions: [PERMISSIONS.EVENTS.DLQ.LIST] },
      ],
    },

    {
      icon: 'event_note',
      label: m.events.value,
      permissions: [PERMISSIONS.EVENTS.LIST],
      children: [
        { label: m.eventStore.value, to: '/events/store', icon: 'storage', permissions: [PERMISSIONS.EVENTS.LIST] },
      ],
    },

    {
      icon: 'admin_panel_settings',
      label: m.administration.value,
      permissions: [
        PERMISSIONS.ORGANIZATIONS.LIST,
        PERMISSIONS.USERS.LIST,
        PERMISSIONS.ROLES.LIST,
        PERMISSIONS.GROUPS.LIST,
        PERMISSIONS.EVENTS.AUDIT.LIST,
        PERMISSIONS.LISTS.LIST,
      ],
      children: [
        { label: m.customers.value, to: '/customers', icon: 'domain', permissions: [PERMISSIONS.ORGANIZATIONS.LIST] },
        { separator: true },
        { label: m.users.value, to: '/users', icon: 'group', permissions: [PERMISSIONS.USERS.LIST] },
        { label: m.roles.value, to: '/roles', icon: 'admin_panel_settings', permissions: [PERMISSIONS.ROLES.LIST] },
        { label: m.groups.value, to: '/groups', icon: 'groups', permissions: [PERMISSIONS.GROUPS.LIST] },
        { separator: true },
        { label: m.accessAudit.value, to: '/access-audit', icon: 'security', permissions: [PERMISSIONS.EVENTS.AUDIT.LIST] },
        { separator: true },
        { label: m.lists.value, to: '/admin/lists', icon: 'list', permissions: [PERMISSIONS.LISTS.LIST] },
        { label: m.settings.value, to: '/admin/settings', icon: 'settings' },
      ],
    },
  ];
}
