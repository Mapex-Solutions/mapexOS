import type { ActiveItemStyle } from '@components/dialogs/common/genericSelectorDialog/interfaces/genericSelectorDialog.interface';

/**
 * Active item style for trigger selector
 */
export const TRIGGER_ACTIVE_ITEM_STYLE: ActiveItemStyle = {
  backgroundColor: 'rgba(255, 152, 0, 0.08)',
  borderColor: 'var(--q-warning)',
};

/**
 * Category icon map for trigger categories
 */
export const CATEGORY_ICON_MAP: Record<string, string> = {
  email: 'email',
  slack: 'chat',
  teams: 'groups',
  http: 'http',
  mqtt: 'router',
  rabbitmq: 'cloud_queue',
  nats: 'cloud',
  websocket: 'cable',
  custom: 'code',
};

/**
 * Category color map for trigger categories
 */
export const CATEGORY_COLOR_MAP: Record<string, string> = {
  email: 'blue',
  slack: 'purple',
  teams: 'indigo',
  http: 'green',
  mqtt: 'orange',
  rabbitmq: 'deep-orange',
  nats: 'cyan',
  websocket: 'teal',
  custom: 'grey',
};
