export const NOTIFICATIONS_LOG_LIST_STUB = [
  {
    id: '1',
    notificationType: 'slack',
    notificationName: 'CreateUser',
    status: 'success',
    tenantId: 'Mapex',
    created: '2025-06-24T08:12:34Z',
    actor: 'thiago.anselmo',
    resource: 'UserService',
    details: 'Created user record with id 1024'
  },
  {
    id: '2',
    notificationType: 'email',
    notificationName: 'EditDataSource',
    status: 'success',
    tenantId: 'Mapex',
    created: '2025-06-24T09:05:12Z',
    actor: 'api.service',
    resource: 'DataSourceModule',
    details: 'Edited connection string for tenant “Mapex”'
  },
  {
    id: '3',
    notificationType: 'webhook',
    notificationName: 'DeleteAsset',
    status: 'failed',
    tenantId: 'TraffordCentre',
    created: '2025-06-24T10:22:47Z',
    actor: 'maria.silva',
    resource: 'AssetManager',
    details: 'Failed to delete asset “Printer-42” due to permission'
  },
  {
    id: '4',
    notificationType: 'push',
    notificationName: 'CreateJSONParser',
    status: 'success',
    tenantId: 'Default',
    created: '2025-06-24T11:17:08Z',
    actor: 'julia.souza',
    resource: 'PayloadProcessor',
    details: 'Created new payload handler “JSONParser”'
  },
  {
    id: '5',
    notificationType: 'telegram',
    notificationName: 'UpdateRuleThreshold',
    status: 'success',
    tenantId: 'Mapex',
    created: '2025-06-24T12:03:55Z',
    actor: 'system.scheduler',
    details: 'Updated rule threshold from 10 to 15'
  },
  {
    id: '6',
    notificationType: 'teams',
    notificationName: 'DeleteHighTemperatureAlert',
    status: 'success',
    tenantId: 'ControlCenter',
    created: '2025-06-24T12:45:29Z',
    actor: 'john.doe',
    resource: 'TriggerService',
    details: 'Deleted trigger “HighTemperatureAlert”'
  },
  {
    id: '7',
    notificationType: 'slack',
    notificationName: 'CreateWeeklySummary',
    status: 'success',
    tenantId: 'Analytics',
    created: '2025-06-24T13:28:10Z',
    actor: 'alice.wong',
    resource: 'TemplateRepo',
    details: 'Created template “WeeklySummary”'
  },
  {
    id: '8',
    notificationType: 'slack',
    notificationName: 'EditUserRole',
    status: 'success',
    tenantId: 'Auth',
    created: '2025-06-24T14:09:03Z',
    actor: 'charlie.brown',
    resource: 'UserAdmin',
    details: 'Changed role of user “eve” to “admin”'
  },
  {
    id: '9',
    notificationType: 'push',
    notificationName: 'DeleteCustomerRecord',
    status: 'failed',
    tenantId: 'CRM',
    created: '2025-06-24T14:55:47Z',
    actor: 'eva.green',
    resource: 'CustomerDB',
    details: 'Failed to delete customer “ACME Corp.” due to foreign key'
  }
];
