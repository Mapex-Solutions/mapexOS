import type { DefinitionResponse } from '@mapexos/schemas';

/**
 * External input definition from workflow definition.
 * Describes one input that the user must fill when configuring an instance.
 */
export interface ExternalInputDefinition {
  /** Key used in the externalInputs map */
  field: string;

  /** Display label for the input */
  label: string;

  /** Material icon name */
  icon: string;

  /** Input type: string | number | boolean | json */
  type: string;

  /** Helper text shown below the input */
  description: string;

  /** Pre-fill value */
  defaultValue: any;

  /** Whether this input is required */
  required: boolean;

  /** Asset template ID — only when type is "assetFromTemplate" */
  assetTemplateId: string;

  /** Field path to extract from the selected asset (e.g. "assetUUID") */
  fieldPath: string;
}

/**
 * Workflow instance form data structure.
 * Holds all fields needed for create/edit operations.
 */
export interface WorkflowInstanceFormData {
  /** Instance display name */
  name: string;

  /** Instance description */
  description: string;

  /** Whether the instance is enabled */
  enabled: boolean;

  /** Whether this instance is shared as template with child orgs */
  isTemplate: boolean;

  /** Selected workflow definition ID */
  definitionId: string | null;

  /** Definition version at time of selection */
  definitionVersion: number;

  /** External inputs filled by the user (key-value map) */
  externalInputs: Record<string, any>;

  /** Whether this instance uses a fixed UUID (single execution at a time) */
  uniqueExecution: boolean;

  /** Fixed UUID for unique execution mode (user-defined) */
  workflowUUID: string;

  /** Full definition response object (for display in steps 2-4) */
  selectedDefinition: DefinitionResponse | null;
}

/**
 * Workflow instance form state (external selections and UI state).
 */
export interface WorkflowInstanceFormState {
  /** Full definition response for review display */
  selectedDefinition: DefinitionResponse | null;

  /** Whether a save operation is in progress */
  isSaving: boolean;

  /** Current active step */
  currentStep: number;
}
