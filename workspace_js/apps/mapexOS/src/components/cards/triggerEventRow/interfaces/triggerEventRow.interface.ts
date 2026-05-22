export interface RawTriggerProps {
  id: string;
  triggerType: 'HTTP' | 'Notification' | 'Incident' | 'Task' | 'Workflow' | 'MQTT';
  triggerName: string;
  status: 'success' | 'failed';
  tenantName?: string;
  created: string;              // ISO timestamp
  [key: string]: any;           // allow any additional properties at the root
}