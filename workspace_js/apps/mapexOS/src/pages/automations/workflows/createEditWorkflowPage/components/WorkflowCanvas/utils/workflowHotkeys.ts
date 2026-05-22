/**
 * Callbacks for workflow canvas hotkeys
 */
export interface WorkflowHotkeyCallbacks {
  /** Called when Delete key is pressed */
  onDelete: () => void;

  /** Called when Shift+D is pressed */
  onDuplicate: () => void;

  /** Called when Ctrl+Z is pressed */
  onUndo?: () => void;

  /** Called when Ctrl+Shift+Z is pressed */
  onRedo?: () => void;

  /** Called when Ctrl+C is pressed */
  onCopy?: () => void;

  /** Called when Ctrl+V is pressed */
  onPaste?: () => void;

  /** Called when Escape is pressed */
  onEscape?: () => void;
}

/**
 * Check if a keyboard event targets a form input (should be ignored by hotkeys)
 *
 * @param {KeyboardEvent} e - Keyboard event
 * @returns {boolean} True if event targets a form input
 */
function isFormInput(e: KeyboardEvent): boolean {
  const tag = (e.target as HTMLElement)?.tagName;
  return tag === 'INPUT' || tag === 'TEXTAREA' || tag === 'SELECT';
}

/**
 * Create workflow canvas hotkey handler.
 * Returns attach/detach functions for lifecycle management.
 *
 * Hotkeys:
 * - Escape → onEscape callback
 * - Ctrl+Z → onUndo callback
 * - Ctrl+Shift+Z → onRedo callback
 * - Delete → onDelete callback
 * - Shift+D → onDuplicate callback
 *
 * @param {WorkflowHotkeyCallbacks} callbacks - Callbacks for each hotkey
 * @returns {{ attach: () => void; detach: () => void }} Lifecycle functions
 */
export function createWorkflowHotkeyHandler(callbacks: WorkflowHotkeyCallbacks) {
  function handler(e: KeyboardEvent): void {
    // Escape → always fires (even from form inputs)
    if (e.key === 'Escape') {
      callbacks.onEscape?.();
      return;
    }

    if (isFormInput(e)) return;

    // Ctrl+Z → Undo
    if ((e.ctrlKey || e.metaKey) && !e.shiftKey && e.key === 'z') {
      e.preventDefault();
      callbacks.onUndo?.();
      return;
    }

    // Ctrl+Shift+Z → Redo
    if ((e.ctrlKey || e.metaKey) && e.shiftKey && (e.key === 'Z' || e.key === 'z')) {
      e.preventDefault();
      callbacks.onRedo?.();
      return;
    }

    // Ctrl+C → Copy
    if ((e.ctrlKey || e.metaKey) && !e.shiftKey && (e.key === 'c' || e.key === 'C')) {
      e.preventDefault();
      callbacks.onCopy?.();
      return;
    }

    // Ctrl+V → Paste
    if ((e.ctrlKey || e.metaKey) && !e.shiftKey && (e.key === 'v' || e.key === 'V')) {
      e.preventDefault();
      callbacks.onPaste?.();
      return;
    }

    if (e.key === 'Delete') {
      callbacks.onDelete();
      return;
    }

    if (e.shiftKey && (e.key === 'D' || e.key === 'd')) {
      e.preventDefault();
      callbacks.onDuplicate();
    }
  }

  return {
    attach: () => document.addEventListener('keydown', handler),
    detach: () => document.removeEventListener('keydown', handler),
  };
}
