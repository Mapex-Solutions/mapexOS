import type { AssetTemplateResponse } from '@mapexos/schemas';

export interface AssetTemplateDetailsModalProps {
	/** Two-way open state. */
	modelValue: boolean;
	/** The template to inspect. Supplied by the caller — the modal never fetches. */
	template: AssetTemplateResponse | null;
}

export interface AssetTemplateDetailsModalEmits {
	(e: 'update:modelValue', value: boolean): void;
}

/**
 * One tab in the details viewer. `always` tabs render even with no data;
 * data tabs render only when `showWhen` returns true (hidden when empty).
 */
export interface AssetTemplateDetailTab {
	key: string;
	icon: string;
	label: string;
	always?: boolean;
	showWhen?: () => boolean;
}
