export interface ChannelPushProps {
  appName: string;
  deviceCount: number;
  apiKey: string;
  serviceProvider: string;
  priority: string;
  ttl: number;
  badge: boolean;
  sound: string;
  clickAction: string;
}