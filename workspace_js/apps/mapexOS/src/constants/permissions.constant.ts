/**
 * Frontend mirror of backend Go permission constants.
 * Source: workspace_go/packages/permissions/
 *
 * IMPORTANT: Keep in sync with backend constants.
 * Each key maps exactly to the Go constant string value.
 */
export const PERMISSIONS = {
  /** Root wildcard — access to EVERYTHING */
  MAPEX_ALL: 'mapex.*',

  /** Admin wildcards — full access within their scope */
  ADMIN_VENDOR_ALL: 'admin_vendor.*',
  ADMIN_CUSTOMER_ALL: 'admin_customer.*',
  ADMIN_ALL: 'admin.*',

  /** Asset permissions */
  ASSETS: {
    LIST: 'assets.list',
    CREATE: 'assets.create',
    READ: 'assets.read',
    UPDATE: 'assets.update',
    DELETE: 'assets.delete',
    ALL: 'assets.*',
  },

  /** Asset template permissions */
  ASSET_TEMPLATES: {
    LIST: 'assettemplates.list',
    CREATE: 'assettemplates.create',
    READ: 'assettemplates.read',
    UPDATE: 'assettemplates.update',
    DELETE: 'assettemplates.delete',
    ALL: 'assettemplates.*',
  },

  /** DataSource permissions */
  DATASOURCES: {
    LIST: 'datasources.list',
    CREATE: 'datasources.create',
    READ: 'datasources.read',
    UPDATE: 'datasources.update',
    DELETE: 'datasources.delete',
    ALL: 'datasources.*',
  },

  /** Route group permissions */
  ROUTE_GROUPS: {
    LIST: 'routegroups.list',
    CREATE: 'routegroups.create',
    READ: 'routegroups.read',
    UPDATE: 'routegroups.update',
    DELETE: 'routegroups.delete',
    ALL: 'routegroups.*',
  },

  /** Trigger permissions */
  TRIGGERS: {
    LIST: 'triggers.list',
    CREATE: 'triggers.create',
    READ: 'triggers.read',
    UPDATE: 'triggers.update',
    DELETE: 'triggers.delete',
    ALL: 'triggers.*',
  },

  /** Workflow permissions */
  WORKFLOWS: {
    LIST: 'workflows.list',
    CREATE: 'workflows.create',
    READ: 'workflows.read',
    UPDATE: 'workflows.update',
    DELETE: 'workflows.delete',
    ALL: 'workflows.*',
    INSTANCES: {
      LIST: 'workflows.instances.list',
      READ: 'workflows.instances.read',
      CANCEL: 'workflows.instances.cancel',
      SIGNAL: 'workflows.instances.signal',
    },
  },

  /** Event permissions (general + per type) */
  EVENTS: {
    LIST: 'events.list',
    READ: 'events.read',
    CREATE: 'events.create',
    DELETE: 'events.delete',
    RAW: { LIST: 'events.raw.list', READ: 'events.raw.read', CREATE: 'events.raw.create', DELETE: 'events.raw.delete' },
    PROCESSED: { LIST: 'events.processed.list', READ: 'events.processed.read', CREATE: 'events.processed.create', DELETE: 'events.processed.delete' },
    JS_EXECUTOR: { LIST: 'events.js_executor.list', READ: 'events.js_executor.read', CREATE: 'events.js_executor.create', DELETE: 'events.js_executor.delete' },
    ROUTER: { LIST: 'events.router.list', READ: 'events.router.read', CREATE: 'events.router.create', DELETE: 'events.router.delete' },
    TRIGGER: { LIST: 'events.trigger.list', READ: 'events.trigger.read', CREATE: 'events.trigger.create', DELETE: 'events.trigger.delete' },
    AUDIT: { LIST: 'events.audit.list', READ: 'events.audit.read', CREATE: 'events.audit.create', DELETE: 'events.audit.delete' },
    NOTIFICATIONS: { LIST: 'events.notifications.list', READ: 'events.notifications.read', CREATE: 'events.notifications.create', DELETE: 'events.notifications.delete' },
    DLQ: { LIST: 'events.dlq.list', READ: 'events.dlq.read' },
    ASSET_STATUS: { LIST: 'events.asset_status.list', READ: 'events.asset_status.read' },
  },

  /** User permissions */
  USERS: {
    LIST: 'users.list',
    CREATE: 'users.create',
    READ: 'users.read',
    UPDATE: 'users.update',
    DELETE: 'users.delete',
    ALL: 'users.*',
  },

  /** Role permissions */
  ROLES: {
    LIST: 'roles.list',
    CREATE: 'roles.create',
    READ: 'roles.read',
    UPDATE: 'roles.update',
    DELETE: 'roles.delete',
    ALL: 'roles.*',
  },

  /** Group permissions */
  GROUPS: {
    LIST: 'groups.list',
    CREATE: 'groups.create',
    READ: 'groups.read',
    UPDATE: 'groups.update',
    DELETE: 'groups.delete',
    ALL: 'groups.*',
  },

  /** Organization permissions */
  ORGANIZATIONS: {
    LIST: 'organizations.list',
    CREATE: 'organizations.create',
    READ: 'organizations.read',
    UPDATE: 'organizations.update',
    DELETE: 'organizations.delete',
    ALL: 'organizations.*',
  },

  /** Membership permissions */
  MEMBERSHIPS: {
    LIST: 'memberships.list',
    CREATE: 'memberships.create',
    READ: 'memberships.read',
    UPDATE: 'memberships.update',
    DELETE: 'memberships.delete',
    ALL: 'memberships.*',
  },

  /** List permissions */
  LISTS: {
    LIST: 'lists.lists',
    CREATE: 'lists.create',
    READ: 'lists.read',
    UPDATE: 'lists.update',
    DELETE: 'lists.delete',
    ALL: 'lists.*',
  },

  /** Scheduler job permissions */
  JOBS: {
    LIST: 'jobs.list',
    CREATE: 'jobs.create',
    READ: 'jobs.read',
    UPDATE: 'jobs.update',
    DELETE: 'jobs.delete',
    ALL: 'jobs.*',
  },

  /** Retention permissions */
  RETENTION: {
    LIST: 'retention.list',
    READ: 'retention.read',
    UPDATE: 'retention.update',
    ALL: 'retention.*',
  },

  /** Auth permissions */
  AUTH: {
    LOGIN: 'auth.login',
    LOGOUT: 'auth.logout',
    REFRESH: 'auth.refresh',
    CHANGE_PASSWORD: 'auth.changepassword',
    RESET_PASSWORD: 'auth.resetpassword',
    ALL: 'auth.*',
  },
} as const;
