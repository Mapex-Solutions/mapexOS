/**
 * Props for StepToTemplate: the target template id slice of the form model.
 */
export interface StepToTemplateProps {
  modelValue: string | null;
}

/**
 * Emits for StepToTemplate.
 * - update:modelValue emits the chosen target template id.
 * - update:valid signals that a target template is selected.
 */
export interface StepToTemplateEmits {
  (e: 'update:modelValue', value: string | null): void;
  (e: 'update:valid', value: boolean): void;
}
