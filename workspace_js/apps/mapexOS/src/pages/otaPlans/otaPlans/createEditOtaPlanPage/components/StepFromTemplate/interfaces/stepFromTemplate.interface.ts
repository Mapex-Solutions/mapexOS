/**
 * Props for StepFromTemplate: the source-template slice of the form model.
 */
export interface StepFromTemplateProps {
	modelValue: string | null;
}

/**
 * Emits for StepFromTemplate: two-way binding for the selected template id.
 */
export interface StepFromTemplateEmits {
	(e: 'update:modelValue', value: string | null): void;
}
