export interface RawNotificationProps {
  id: string;
  notificationType: 'slack' | 'teams' | 'email' | 'push' | 'telegram' | 'webhook';
  notificationName: string;
  status: 'success' | 'failed';
  tenantId?: string;
  created: string;              // ISO timestamp
  [key: string]: any;           // allow any additional properties at the root
}