import type { ExternalVariable } from '@src/components/workflow/interfaces';

/**
 * Local form type with relaxed defaultValue for v-model compatibility
 */
export type ExternalVariableForm = Omit<ExternalVariable, 'defaultValue'> & {
  /** Relaxed default value for v-model binding */
  defaultValue: any;

  /** Asset template ID — only when type = 'assetFromTemplate' */
  assetTemplateId?: string;

  /** Asset field path — only when type = 'assetFromTemplate' */
  fieldPath?: string;
};
