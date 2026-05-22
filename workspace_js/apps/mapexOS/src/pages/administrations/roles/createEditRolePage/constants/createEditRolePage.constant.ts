/**
 * CreateEditRolePage Constants
 */

import type { RoleFormData, ResourcePermission, PermissionGroup } from '../interfaces';

/**
 * Total number of steps in the form
 */
export const TOTAL_STEPS = 3;

/**
 * Step numbers enum for better readability
 */
export const STEP = {
  BASIC_INFO: 1,
  PERMISSIONS: 2,
  REVIEW: 3,
} as const;

/**
 * Initial form data values
 */
export const INITIAL_ROLE_FORM_DATA: RoleFormData = {
  name: '',
  description: '',
  scope: null,
  isTemplate: false,
};

/**
 * Scope options for role configuration
 */
export const SCOPE_OPTIONS = [
  {
    label: 'Global',
    value: 'global',
    description: 'Role permissions inherit to child organizations',
    icon: 'public',
  },
  {
    label: 'Local',
    value: 'local',
    description: 'Role permissions apply only to this organization',
    icon: 'place',
  },
] as const;

/**
 * Standard actions available for most resources
 */
export const STANDARD_ACTIONS = ['list', 'create', 'read', 'update', 'delete'] as const;

/**
 * Default resource permissions configuration
 * Based on MapexOS permission system
 * Ordered by sidebar grouping: Device Management → Data → Automation → Routing → Logs → Security → Administration → System
 */
export const DEFAULT_RESOURCE_PERMISSIONS: ResourcePermission[] = [
  // ── Device Management ──────────────────────────────────────────────
  {
    resource: 'assets',
    label: 'Assets',
    icon: 'router',
    enabled: false,
    actions: [
      { name: 'list', label: 'List', granted: false },
      { name: 'create', label: 'Create', granted: false },
      { name: 'read', label: 'Read', granted: false },
      { name: 'update', label: 'Update', granted: false },
      { name: 'delete', label: 'Delete', granted: false },
    ],
  },
  {
    resource: 'assettemplates',
    label: 'Asset Templates',
    icon: 'memory',
    enabled: false,
    actions: [
      { name: 'list', label: 'List', granted: false },
      { name: 'create', label: 'Create', granted: false },
      { name: 'read', label: 'Read', granted: false },
      { name: 'update', label: 'Update', granted: false },
      { name: 'delete', label: 'Delete', granted: false },
    ],
  },

  // ── Data ───────────────────────────────────────────────────────────
  {
    resource: 'datasources',
    label: 'Data Sources',
    icon: 'storage',
    enabled: false,
    actions: [
      { name: 'list', label: 'List', granted: false },
      { name: 'create', label: 'Create', granted: false },
      { name: 'read', label: 'Read', granted: false },
      { name: 'update', label: 'Update', granted: false },
      { name: 'delete', label: 'Delete', granted: false },
    ],
  },
  {
    resource: 'retention',
    label: 'Retention Policies',
    icon: 'cloud_upload',
    enabled: false,
    actions: [
      { name: 'list', label: 'List', granted: false },
      { name: 'read', label: 'Read', granted: false },
      { name: 'update', label: 'Update', granted: false },
    ],
  },

  // ── Automation ─────────────────────────────────────────────────────
  {
    resource: 'triggers',
    label: 'Triggers',
    icon: 'flash_on',
    enabled: false,
    actions: [
      { name: 'list', label: 'List', granted: false },
      { name: 'create', label: 'Create', granted: false },
      { name: 'read', label: 'Read', granted: false },
      { name: 'update', label: 'Update', granted: false },
      { name: 'delete', label: 'Delete', granted: false },
    ],
  },

  // ── Routing ────────────────────────────────────────────────────────
  {
    resource: 'routegroups',
    label: 'Route Groups',
    icon: 'account_tree',
    enabled: false,
    actions: [
      { name: 'list', label: 'List', granted: false },
      { name: 'create', label: 'Create', granted: false },
      { name: 'read', label: 'Read', granted: false },
      { name: 'update', label: 'Update', granted: false },
      { name: 'delete', label: 'Delete', granted: false },
    ],
  },

  // ── Logs & Executions (read-only) ─────────────────────────────────
  {
    resource: 'events',
    label: 'Event Tracer',
    icon: 'account_tree',
    enabled: false,
    actions: [
      { name: 'list', label: 'List', granted: false },
      { name: 'read', label: 'Read', granted: false },
    ],
  },
  {
    resource: 'events.raw',
    label: 'Raw Events',
    icon: 'terminal',
    enabled: false,
    actions: [
      { name: 'list', label: 'List', granted: false, permissionKey: 'events.raw.list' },
      { name: 'read', label: 'Read', granted: false, permissionKey: 'events.raw.read' },
    ],
  },
  {
    resource: 'events.processed',
    label: 'Event Store',
    icon: 'done_all',
    enabled: false,
    actions: [
      { name: 'list', label: 'List', granted: false, permissionKey: 'events.processed.list' },
      { name: 'read', label: 'Read', granted: false, permissionKey: 'events.processed.read' },
    ],
  },
  {
    resource: 'events.js_executor',
    label: 'JS Executor',
    icon: 'code',
    enabled: false,
    actions: [
      { name: 'list', label: 'List', granted: false, permissionKey: 'events.js_executor.list' },
      { name: 'read', label: 'Read', granted: false, permissionKey: 'events.js_executor.read' },
    ],
  },
  {
    resource: 'events.router',
    label: 'Router Logs',
    icon: 'alt_route',
    enabled: false,
    actions: [
      { name: 'list', label: 'List', granted: false, permissionKey: 'events.router.list' },
      { name: 'read', label: 'Read', granted: false, permissionKey: 'events.router.read' },
    ],
  },
  {
    resource: 'events.trigger',
    label: 'Trigger Logs',
    icon: 'history_toggle_off',
    enabled: false,
    actions: [
      { name: 'list', label: 'List', granted: false, permissionKey: 'events.trigger.list' },
      { name: 'read', label: 'Read', granted: false, permissionKey: 'events.trigger.read' },
    ],
  },

  // ── Security ───────────────────────────────────────────────────────
  {
    resource: 'events.audit',
    label: 'Access Audit',
    icon: 'history',
    enabled: false,
    actions: [
      { name: 'list', label: 'List', granted: false, permissionKey: 'events.audit.list' },
      { name: 'read', label: 'Read', granted: false, permissionKey: 'events.audit.read' },
    ],
  },
  // ── Administration ─────────────────────────────────────────────────
  {
    resource: 'organizations',
    label: 'Customers',
    icon: 'domain',
    enabled: false,
    actions: [
      { name: 'list', label: 'List', granted: false },
      { name: 'create', label: 'Create', granted: false },
      { name: 'read', label: 'Read', granted: false },
      { name: 'update', label: 'Update', granted: false },
      { name: 'delete', label: 'Delete', granted: false },
    ],
  },
  {
    resource: 'users',
    label: 'Users',
    icon: 'person',
    enabled: false,
    actions: [
      { name: 'list', label: 'List', granted: false },
      { name: 'create', label: 'Create', granted: false },
      { name: 'read', label: 'Read', granted: false },
      { name: 'update', label: 'Update', granted: false },
      { name: 'delete', label: 'Delete', granted: false },
    ],
  },
  {
    resource: 'roles',
    label: 'Roles',
    icon: 'admin_panel_settings',
    enabled: false,
    actions: [
      { name: 'list', label: 'List', granted: false },
      { name: 'create', label: 'Create', granted: false },
      { name: 'read', label: 'Read', granted: false },
      { name: 'update', label: 'Update', granted: false },
      { name: 'delete', label: 'Delete', granted: false },
    ],
  },
  {
    resource: 'groups',
    label: 'Groups',
    icon: 'group',
    enabled: false,
    actions: [
      { name: 'list', label: 'List', granted: false },
      { name: 'create', label: 'Create', granted: false },
      { name: 'read', label: 'Read', granted: false },
      { name: 'update', label: 'Update', granted: false },
      { name: 'delete', label: 'Delete', granted: false },
    ],
  },
  {
    resource: 'memberships',
    label: 'Memberships',
    icon: 'card_membership',
    enabled: false,
    actions: [
      { name: 'list', label: 'List', granted: false },
      { name: 'create', label: 'Create', granted: false },
      { name: 'read', label: 'Read', granted: false },
      { name: 'update', label: 'Update', granted: false },
      { name: 'delete', label: 'Delete', granted: false },
    ],
  },
  {
    resource: 'lists',
    label: 'Lists',
    icon: 'list_alt',
    enabled: false,
    actions: [
      { name: 'list', label: 'List', granted: false, permissionKey: 'lists.lists' },
      { name: 'create', label: 'Create', granted: false },
      { name: 'read', label: 'Read', granted: false },
      { name: 'update', label: 'Update', granted: false },
      { name: 'delete', label: 'Delete', granted: false },
    ],
  },

];

/**
 * Permission groups matching sidebar menu sections
 * Used to visually group resources in Step2Permissions
 */
export const PERMISSION_GROUPS: PermissionGroup[] = [
  {
    label: 'Assets',
    icon: 'device_hub',
    resources: ['assets', 'assettemplates'],
  },
  {
    label: 'Data',
    icon: 'settings_input_antenna',
    resources: ['datasources', 'retention'],
  },
  {
    label: 'Automation',
    icon: 'psychology',
    resources: ['businessrules', 'rules', 'triggers'],
  },
  {
    label: 'Routing',
    icon: 'route',
    resources: ['routegroups'],
  },
  {
    label: 'Logs & Executions',
    icon: 'list_alt',
    resources: ['events', 'events.raw', 'events.processed', 'events.js_executor', 'events.router', 'events.trigger'],
  },
  {
    label: 'Administration',
    icon: 'admin_panel_settings',
    resources: ['organizations', 'users', 'roles', 'groups', 'memberships', 'events.audit', 'lists'],
  },
];

/**
 * Name validation rules
 */
export const NAME_MIN_LENGTH = 3;
export const NAME_MAX_LENGTH = 100;

/**
 * Description validation rules
 */
export const DESCRIPTION_MAX_LENGTH = 500;
