import type { BaseChannel } from './base.interface';

export interface TeamsConfig extends BaseChannel {
  type: 'teams';
  teamName: string;
  channelName: string;
}
