type JsonData = object | string;

export interface JsonDrawerProps {
	show: boolean;
	jsonData: JsonData;
	editable?: boolean;
	title: string;
	subtitle?: string;
}

export interface JsonDrawerEmit {
	(e: 'update:show', value: boolean): void;

	(e: 'save', updated: JsonData): void;
}