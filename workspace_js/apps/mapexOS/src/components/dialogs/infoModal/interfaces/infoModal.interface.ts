export interface InfoItem {
	icon?: string;
	color?: string;
	title?: string;
	text: string;
}

export interface InfoModalProps {
	modelValue: boolean;
	icon?: string | undefined;
	title: string;
	description: string;
	items?: InfoItem[] | undefined;
	docsUrl?: string | undefined;
	docsLabel?: string | undefined;
	closeLabel?: string | undefined;
	showActions?: boolean | undefined;
}

export interface InfoModalEmits {
	(e: 'update:modelValue', value: boolean): void;
}
