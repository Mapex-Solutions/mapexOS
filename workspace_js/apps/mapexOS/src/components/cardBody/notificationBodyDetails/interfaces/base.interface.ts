// interfaces.ts
export type ChannelType = 'slack' | 'teams' | 'email' | 'push';

export interface BaseChannel {
  type: ChannelType;
  name: string;           // e.g. "slack"
  label: string;          // e.g. "Slack"
  icon: string;           // mdi-slack, etc.
  enabled: boolean;
}