import type { BaseChannel } from './base.interface';

export interface EmailConfig extends BaseChannel {
	type: 'email';
	from: string;           // sender
	to: string[];           // recipients
}
