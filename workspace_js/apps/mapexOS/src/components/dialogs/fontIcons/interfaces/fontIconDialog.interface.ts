export interface FontIconDialogProps {
	modelValue: string,
	show: boolean
}

export interface FontIconDialogEvents {
	(e: 'update:modelValue', value: string): void,
	(e: 'update:show', value: boolean): void
}

export interface FontIconData {
	name: string;
	categories: string[];
}