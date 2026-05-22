import type { BaseChannel } from './base.interface';

export interface PushConfig extends BaseChannel {
  type: 'push';
  appName: string;
  deviceCount?: number;
}
