/**
 * ScriptEditorDialog Interfaces
 */

/**
 * Props for ScriptEditorDialog component
 */
export interface ScriptEditorDialogProps {
  /** Dialog visibility (v-model) */
  modelValue: boolean;

  /** Dialog header title */
  title: string;

  /** Script content to edit */
  scriptContent: string;

  /** Monaco editor language */
  language?: 'javascript' | 'json' | 'typescript';

  /** Close button tooltip text */
  closeTooltip?: string;

  /** Optional info banner items displayed below header */
  guidelines?: ScriptEditorGuideline[];
}

/**
 * A single guideline item displayed in the info bar
 */
export interface ScriptEditorGuideline {
  /** Code keyword (displayed in <code> tag) */
  code: string;

  /** Description text */
  description: string;
}

/**
 * Emits for ScriptEditorDialog component
 */
export interface ScriptEditorDialogEmits {
  (e: 'update:modelValue', value: boolean): void;
  (e: 'update:scriptContent', value: string): void;
}
