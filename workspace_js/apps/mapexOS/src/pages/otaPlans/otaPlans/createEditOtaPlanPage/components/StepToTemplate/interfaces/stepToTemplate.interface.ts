/**
 * Props for StepToTemplate: the target-template slice of the form model.
 */
export interface StepToTemplateProps {
	modelValue: string | null;
}

/**
 * Emits for StepToTemplate: two-way binding for the selected template id.
 */
export interface StepToTemplateEmits {
	(e: 'update:modelValue', value: string | null): void;
}
