/**
 * Trigger category for UI organization
 */
export type TriggerCategory = 'technical' | 'communication';

/**
 * All available trigger types
 */
export type TriggerType =
  | 'http'
  | 'mqtt'
  | 'rabbitmq'
  | 'nats'
  | 'tcp'
  | 'websocket'
  | 'email'
  | 'teams'
  | 'slack'
  | 'push'
  | 'sms';

/**
 * Trigger status
 */
export type TriggerStatus = 'Active' | 'Inactive';

/**
 * Trigger list item interface (deprecated - use TriggerResponse from @mapexos/schemas)
 * @deprecated This interface is no longer used. Use TriggerResponse from @mapexos/schemas instead.
 */
export interface TriggerListItem {
  id: string;
  name: string;
  description?: string;
  triggerType: TriggerType;
  category: TriggerCategory;
  status: TriggerStatus;
  icon?: string;
  createdAt?: string;
  updatedAt?: string;
}

/**
 * Trigger list page filters state
 * Maps to TriggerQuery fields from backend
 */
export interface TriggerListPageFilters {
  /** Search by trigger name (partial match) */
  name: string | undefined;

  /** Filter by trigger type */
  triggerType: string | undefined;

  /** Filter by category (technical/communication) */
  category: string | undefined;

  /** Filter by status (active/inactive) */
  status: boolean | undefined;
}

/**
 * Result returned from fetchTriggersHandler
 */
export interface FetchTriggersResult {
  /** List of triggers from API */
  triggers: any[];

  /** Total number of pages */
  totalPages: number;

  /** Total number of items across all pages */
  totalItems: number;
}
