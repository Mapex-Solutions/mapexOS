// Types
export type ChannelType = 'slack' | 'teams' | 'email' | 'push' | 'telegram' | 'webhook';

export interface ChannelBaseProps {
  type: ChannelType;
  name: string;
  description: string;
  icon: string;
  status: 'Active' | 'Inactive';
  created: string;
  updated?: string;
}