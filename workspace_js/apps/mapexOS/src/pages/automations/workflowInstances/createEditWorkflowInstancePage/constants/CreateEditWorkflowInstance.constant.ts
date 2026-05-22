import type { WorkflowInstanceFormData } from '../interfaces';

/**
 * Initial form data values for workflow instance create/edit
 */
export const INITIAL_FORM_DATA: WorkflowInstanceFormData = {
  name: '',
  description: '',
  enabled: true,
  isTemplate: false,
  definitionId: null,
  definitionVersion: 1,
  externalInputs: {},
  uniqueExecution: false,
  workflowUUID: '',
  selectedDefinition: null,
};

/**
 * Total number of steps in the form
 */
export const TOTAL_STEPS = 4;

/**
 * Step numbers for better readability
 */
export const STEP = {
  IDENTIFICATION: 1,
  DEFINITION: 2,
  EXTERNAL_INPUTS: 3,
  REVIEW: 4,
} as const;
