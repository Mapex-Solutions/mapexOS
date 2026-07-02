import type { AssetTemplateData, DynamicFieldType } from '../../../interfaces';

/** PROPS & EMITS */
export interface Step8DynamicFieldsProps {
	modelValue: AssetTemplateData;
}

export interface Step8DynamicFieldsEmits {
	(e: 'update:modelValue', value: AssetTemplateData): void;
}

/**
 * Flattened entry for the field-name suggestion select.
 *
 * The list interleaves non-selectable category headers (`isHeader`) with
 * selectable canonical fields, preserving the backend group order so Quasar
 * renders a visual separator per category.
 */
export interface FieldNameOption {
	/**
	 * Canonical, English-only field name used as the option value and label.
	 * For headers this carries the localized category label and is never picked.
	 */
	value: string;

	/**
	 * Localized hint shown as the option caption (empty for headers).
	 */
	hint: string;

	/**
	 * Unit shown alongside the hint (empty when dimensionless or for headers).
	 */
	unit: string;

	/**
	 * Dynamic-field type pre-selected when this suggestion is picked.
	 * Undefined for headers.
	 */
	type?: DynamicFieldType;

	/**
	 * Marks a non-selectable category separator row.
	 */
	isHeader: boolean;
}
