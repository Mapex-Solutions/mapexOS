import type { PluginCredentialDefinition } from '@components/workflow/interfaces';

/**
 * Props for the CredentialForm component
 */
export interface CredentialFormProps {
  /** Credential definition from the plugin manifest */
  credentialDef: PluginCredentialDefinition;

  /** Initial name for the credential (empty for new) */
  initialName: string;

  /** Whether this is an edit (secrets show placeholder instead of empty) */
  isEdit: boolean;
}

/**
 * Emits for the CredentialForm component
 */
export interface CredentialFormEmits {
  (e: 'save', payload: { name: string; data: Record<string, unknown> }): void;
  (e: 'cancel'): void;
}
