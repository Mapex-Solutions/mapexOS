import type { WorkflowVariable } from '@src/components/workflow/interfaces';

/**
 * Local form type with relaxed defaultValue for v-model compatibility
 */
export type VariableForm = Omit<WorkflowVariable, 'defaultValue'> & {
  defaultValue: any;
};
