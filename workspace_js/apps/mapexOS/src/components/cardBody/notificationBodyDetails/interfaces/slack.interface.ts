import type { BaseChannel } from './base.interface';

export interface SlackConfig extends BaseChannel {
  type: 'slack';
  workspace: string;      // e.g. "Acme Corp"
  channelName: string;    // e.g. "#alerts"
}