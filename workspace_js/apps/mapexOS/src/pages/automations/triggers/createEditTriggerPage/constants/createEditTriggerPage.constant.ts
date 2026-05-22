import type {
  CategoryOption,
  TriggerTypeOption,
  Trigger,
  TriggerFormState,
} from '../interfaces';
import { TriggerTypeEnum, TriggerCategoryEnum } from '@mapexos/schemas';

/**
 * Total number of steps in the trigger creation wizard
 */
export const TOTAL_STEPS = 5;

/**
 * Step enum for better readability
 */
export const STEP = {
  CATEGORY: 1,
  TYPE: 2,
  BASIC_INFO: 3,
  CONFIGURATION: 4,
  REVIEW: 5,
} as const;

/**
 * Category options for Step 1
 */
export const CATEGORY_OPTIONS: CategoryOption[] = [
  {
    value: TriggerCategoryEnum.TECHNICAL,
    label: 'Technical Triggers',
    description: 'Infrastructure and system-level triggers (HTTP, MQTT, TCP, etc.)',
    icon: 'mdi-server-network',
    emoji: '📡',
  },
  {
    value: TriggerCategoryEnum.COMMUNICATION,
    label: 'Communication Triggers',
    description: 'Human-facing notification triggers (Email, Slack, Teams, etc.)',
    icon: 'mdi-message-badge',
    emoji: '💬',
  },
];

/**
 * Trigger type options for Step 2
 */
export const TRIGGER_TYPE_OPTIONS: TriggerTypeOption[] = [
  // Technical Triggers
  {
    value: TriggerTypeEnum.HTTP,
    label: 'HTTP',
    description: 'Send HTTP/HTTPS requests to external endpoints',
    icon: 'mdi-web',
    category: TriggerCategoryEnum.TECHNICAL,
  },
  {
    value: TriggerTypeEnum.MQTT,
    label: 'MQTT',
    description: 'Publish messages to MQTT broker topics',
    icon: 'mdi-transit-connection-variant',
    category: TriggerCategoryEnum.TECHNICAL,
  },
  {
    value: TriggerTypeEnum.RABBITMQ,
    label: 'RabbitMQ',
    description: 'Send messages to RabbitMQ queues or exchanges',
    icon: 'mdi-rabbit',
    category: TriggerCategoryEnum.TECHNICAL,
  },
  {
    value: TriggerTypeEnum.NATS,
    label: 'NATS',
    description: 'Publish to NATS messaging system',
    icon: 'mdi-email-fast',
    category: TriggerCategoryEnum.TECHNICAL,
  },
  {
    value: TriggerTypeEnum.WEBSOCKET,
    label: 'WebSocket',
    description: 'Send WebSocket messages to connected clients',
    icon: 'mdi-swap-horizontal',
    category: TriggerCategoryEnum.TECHNICAL,
  },

  // Communication Triggers
  {
    value: TriggerTypeEnum.EMAIL,
    label: 'Email',
    description: 'Send email notifications to recipients',
    icon: 'mdi-email',
    category: TriggerCategoryEnum.COMMUNICATION,
  },
  {
    value: TriggerTypeEnum.TEAMS,
    label: 'Microsoft Teams',
    description: 'Send messages to Teams channels via webhook',
    icon: 'mdi-microsoft-teams',
    category: TriggerCategoryEnum.COMMUNICATION,
  },
  {
    value: TriggerTypeEnum.SLACK,
    label: 'Slack',
    description: 'Send messages to Slack channels or users',
    icon: 'mdi-slack',
    category: TriggerCategoryEnum.COMMUNICATION,
  },
];

/**
 * Initial trigger form data
 */
export const INITIAL_TRIGGER_FORM_DATA: Trigger = {
  name: '',
  description: '',
  triggerType: TriggerTypeEnum.HTTP,
  category: TriggerCategoryEnum.TECHNICAL,
  enabled: true,
  config: {},
};

/**
 * Initial form state
 */
export const INITIAL_FORM_STATE: TriggerFormState = {
  selectedCategory: null,
  selectedType: null,
  isCreating: false,
  currentStep: 1,
};

