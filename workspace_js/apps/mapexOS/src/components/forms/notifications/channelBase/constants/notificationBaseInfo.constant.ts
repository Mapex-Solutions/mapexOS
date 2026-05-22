export const CHANNEL_TYPES = [
  {
    value: 'slack',
    label: 'Slack',
    icon: 'mdi-slack',
    description: 'Send notifications to Slack channels',
  },
  {
    value: 'teams',
    label: 'Microsoft Teams',
    icon: 'mdi-microsoft-teams',
    description: 'Post messages to Microsoft Teams channels',
  },
  {
    value: 'email',
    label: 'Email',
    icon: 'mdi-email',
    description: 'Send emails to configured recipients',
  },
  {
    value: 'push',
    label: 'Push Notifications',
    icon: 'mdi-bell-ring',
    description: 'Push notifications to mobile devices',
  },
  {
    value: 'telegram',
    label: 'Telegram',
    icon: 'mdi-send',
    description: 'Send messages via Telegram bot',
  },
  {
    value: 'webhook',
    label: 'Webhook',
    icon: 'mdi-webhook',
    description: 'Make HTTP calls to external endpoints',
  },
];
