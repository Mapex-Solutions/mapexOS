import type { CanvasToolbarState } from '../../../interfaces/CreateEditWorkflow.interface';

/**
 * Default canvas toolbar state
 */
export const DEFAULT_TOOLBAR_STATE: CanvasToolbarState = {
  showMinimap: true,
  showGrid: true,
  locked: false,
  maximized: false,
};

/**
 * Hotkey definitions — icon, color, key combo (fixed), and i18n key for description
 */
export const HOTKEY_DEFINITIONS = [
  { icon: 'close', color: 'grey-7', title: 'Escape', i18nKey: 'hotkeyEscape' },
  { icon: 'content_copy', color: 'primary', title: 'Ctrl + C', i18nKey: 'hotkeyCopy' },
  { icon: 'content_paste', color: 'primary', title: 'Ctrl + V', i18nKey: 'hotkeyPaste' },
  { icon: 'undo', color: 'primary', title: 'Ctrl + Z', i18nKey: 'hotkeyUndo' },
  { icon: 'redo', color: 'primary', title: 'Ctrl + Shift + Z', i18nKey: 'hotkeyRedo' },
  { icon: 'delete', color: 'negative', title: 'Delete', i18nKey: 'hotkeyDelete' },
  { icon: 'copy_all', color: 'primary', title: 'Shift + D', i18nKey: 'hotkeyDuplicate' },
  { icon: 'select_all', color: 'accent', title: 'Shift + Drag', i18nKey: 'hotkeyBoxSelect' },
  { icon: 'add_circle', color: 'accent', title: 'Ctrl + Click', i18nKey: 'hotkeyMultiSelect' },
] as const;
