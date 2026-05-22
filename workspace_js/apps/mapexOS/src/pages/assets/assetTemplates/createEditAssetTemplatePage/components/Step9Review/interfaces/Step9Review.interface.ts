import type { AssetTemplateData } from '../../../interfaces';

/** PROPS & EMITS */
export interface Step9ReviewProps {
	modelValue: AssetTemplateData;
}

export interface Step9ReviewEmits {
	(e: 'editSection', stepNumber: number): void;
}
