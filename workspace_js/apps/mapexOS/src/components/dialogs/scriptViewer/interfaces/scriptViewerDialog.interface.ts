/**
 * ScriptViewerDialog Interfaces
 */

// Props interface
export interface ScriptViewerDialogProps {
  modelValue: boolean;
  title: string;
  scriptContent: string;
  language?: 'javascript' | 'json';
  copyTooltip?: string;
  closeTooltip?: string;
  copySuccessMessage?: string;
  copyFailMessage?: string;
}

// Emits interface
export interface ScriptViewerDialogEmits {
  (e: 'update:modelValue', value: boolean): void;
}
