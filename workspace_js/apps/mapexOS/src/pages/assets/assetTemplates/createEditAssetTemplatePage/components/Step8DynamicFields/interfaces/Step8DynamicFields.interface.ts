import type { AssetTemplateData } from '../../../interfaces';

/** PROPS & EMITS */
export interface Step8DynamicFieldsProps {
	modelValue: AssetTemplateData;
}

export interface Step8DynamicFieldsEmits {
	(e: 'update:modelValue', value: AssetTemplateData): void;
}
