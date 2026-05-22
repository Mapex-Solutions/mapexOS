import type { PluginCredentialDefinition } from '@components/workflow/interfaces';

/**
 * Props for the CredentialManagerDialog component
 */
export interface CredentialManagerDialogProps {
  /** Whether the dialog is visible */
  modelValue: boolean;

  /** Plugin ID (e.g., 'telegram') */
  pluginId: string;

  /** Plugin display name (for dialog title) */
  pluginName: string;

  /** All credential definitions from the plugin manifest */
  credentialDefs: PluginCredentialDefinition[];
}

/**
 * Emits for the CredentialManagerDialog component
 */
export interface CredentialManagerDialogEmits {
  (e: 'update:modelValue', value: boolean): void;
}

/**
 * Single credential item in the list view
 */
export interface CredentialListItem {
  /** MongoDB ObjectID */
  id: string;

  /** User-given label */
  name: string;

  /** ISO date string */
  created: string;

  /** Test status: null = untested, true = passed, false = failed */
  testStatus: boolean | null;
}
