import type { PluginCredentialDefinition } from '@components/workflow/interfaces';

/**
 * Props for the CredentialSelector component
 */
export interface CredentialSelectorProps {
  /** Plugin ID (e.g., 'telegram') */
  pluginId: string;

  /** Plugin display name (for dialog title) */
  pluginName: string;

  /** All credential definitions from the plugin manifest */
  credentialDefs: PluginCredentialDefinition[];

  /** Currently selected credential ID (null = none) */
  modelValue: string | null;
}

/**
 * Emits for the CredentialSelector component
 */
export interface CredentialSelectorEmits {
  (e: 'update:modelValue', value: string | null): void;
}
