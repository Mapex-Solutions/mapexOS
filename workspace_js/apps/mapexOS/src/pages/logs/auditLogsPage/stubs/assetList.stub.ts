export const AUDIT_LOG_LIST_STUB = [
  {
    id: 1,
    type: 'userLog',
    actor: 'thiago.anselmo',
    action: 'Create',
    resource: 'UserService',
    status: 'success',
    created: '2025-06-24T08:12:34Z',
    details: 'Created user record with id 1024'
  },
  {
    id: 2,
    type: 'dataSource',
    actor: 'api.service',
    action: 'Edit',
    resource: 'DataSourceModule',
    status: 'success',
    created: '2025-06-24T09:05:12Z',
    details: 'Edited connection string for tenant “Mapex”'
  },
  {
    id: 3,
    type: 'assets',
    actor: 'maria.silva',
    action: 'Delete',
    resource: 'AssetManager',
    status: 'failure',
    created: '2025-06-24T10:22:47Z',
    details: 'Failed to delete asset “Printer-42” due to permission'
  },
  {
    id: 4,
    type: 'payloadHandler',
    actor: 'julia.souza',
    action: 'Create',
    resource: 'PayloadProcessor',
    status: 'success',
    created: '2025-06-24T11:17:08Z',
    details: 'Created new payload handler “JSONParser”'
  },
  {
    id: 6,
    type: 'triggers',
    actor: 'john.doe',
    action: 'Delete',
    resource: 'TriggerService',
    status: 'success',
    created: '2025-06-24T12:45:29Z',
    details: 'Deleted trigger “HighTemperatureAlert”'
  },
  {
    id: 7,
    type: 'ruleTemplate',
    actor: 'alice.wong',
    action: 'Create',
    resource: 'TemplateRepo',
    status: 'success',
    created: '2025-06-24T13:28:10Z',
    details: 'Created template “WeeklySummary”'
  },
  {
    id: 8,
    type: 'users',
    actor: 'charlie.brown',
    action: 'Edit',
    resource: 'UserAdmin',
    status: 'success',
    created: '2025-06-24T14:09:03Z',
    details: 'Changed role of user “eve” to “admin”'
  },
  {
    id: 9,
    type: 'customers',
    actor: 'eva.green',
    action: 'Delete',
    resource: 'CustomerDB',
    status: 'failure',
    created: '2025-06-24T14:55:47Z',
    details: 'Failed to delete customer “ACME Corp.” due to foreign key'
  }
];
