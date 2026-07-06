/**
 * Props for StepDetails: the plan name and description slice of the form model.
 */
export interface StepDetailsProps {
	name: string;
	description: string;
}

/**
 * Emits for StepDetails: two-way bindings for name and description.
 */
export interface StepDetailsEmits {
	(e: 'update:name', value: string): void;
	(e: 'update:description', value: string): void;
}
