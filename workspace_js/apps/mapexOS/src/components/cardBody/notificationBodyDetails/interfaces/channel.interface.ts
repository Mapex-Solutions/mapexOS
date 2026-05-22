import type { EmailConfig } from './email.interface';
import type { SlackConfig } from './slack.interface';
import type { PushConfig } from './push.interface';
import type { TeamsConfig } from './teams.interface';

export type ChannelConfig =
	| SlackConfig
	| TeamsConfig
	| EmailConfig
	| PushConfig;
