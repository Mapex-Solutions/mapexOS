export interface AuditLogProps {
  id: string;
  type:
    | 'userLog'
    | 'dataSource'
    | 'assets'
    | 'payloadHandler'
    | 'businessRule'
    | 'triggers'
    | 'ruleTemplate'
    | 'users'
    | 'customers';
  actor: string;
  action: 'Create' | 'Update' | 'Edit' | 'Delete';
  resource: string;
  status: 'success' | 'failure' | 'warning' | 'info';
  created: string;
  details?: string;
  [key: string]: any;
}