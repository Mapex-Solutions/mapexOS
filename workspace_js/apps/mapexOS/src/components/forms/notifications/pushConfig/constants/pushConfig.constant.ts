export const SERVICE_PROVIDERS = [
  { value: 'firebase', label: 'Firebase Cloud Messaging', description: "Google's notification service" },
  { value: 'apns', label: 'Apple Push Notification Service', description: "Apple's notification service" },
  { value: 'azure', label: 'Azure Notification Hubs', description: "Microsoft's notification service" },
  { value: 'onesignal', label: 'OneSignal', description: 'Cross-platform notification platform' },
  { value: 'amazon-sns', label: 'Amazon SNS', description: "Amazon's notification service" },
];

export const PRIORITY_OPTIONS = [
  { value: 'high', label: 'High', description: 'Immediate delivery, may consume more battery' },
  { value: 'normal', label: 'Normal', description: 'Standard delivery balancing battery usage and speed' },
  { value: 'low', label: 'Low', description: 'Battery-saving delivery, lower priority' },
];

export const SOUND_OPTIONS = [
  'default',
  'none',
  'custom1',
  'custom2',
  'notification.mp3',
];
