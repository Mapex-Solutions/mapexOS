/** TYPE IMPORTS */
import type { HistorySnapshot } from '../interfaces';

/** VUE IMPORTS */
import { ref, computed, toRaw } from 'vue';

/** COMPOSABLES */
import { useWorkflowEditorState } from './useWorkflowEditorState';

/** LOCAL IMPORTS (constants and handlers ONLY) */
import { MAX_HISTORY_SIZE } from '../constants';

/** STATE (module-level singleton — shared across all callers) */

/**
 * Undo stack — most recent snapshot at the end
 */
const undoStack = ref<HistorySnapshot[]>([]);

/**
 * Redo stack — most recent snapshot at the end
 */
const redoStack = ref<HistorySnapshot[]>([]);

/**
 * Flag set during undo/redo restoration.
 * Guards WorkflowCanvas watchers from firing and overriding restored positions.
 * Caller must call finishRestore() after Vue Flow sync (via nextTick).
 */
const isBatchRestoring = ref(false);

/**
 * Composable for workflow canvas undo/redo history.
 * Uses a snapshot-based approach: before each mutation, the current
 * state of nodes + edges is deep-cloned and pushed to the undo stack.
 *
 * Module-level singleton pattern — all callers share the same stacks.
 *
 * @returns {object} History state, computed values, and methods
 */
export function useWorkflowHistory() {
  /** COMPOSABLES & STORES */
  const { nodes, edges } = useWorkflowEditorState();

  /** COMPUTED */

  /**
   * Whether undo is available
   * @returns {boolean} True if undo stack has entries
   */
  const canUndo = computed(() => undoStack.value.length > 0);

  /**
   * Whether redo is available
   * @returns {boolean} True if redo stack has entries
   */
  const canRedo = computed(() => redoStack.value.length > 0);

  /** FUNCTIONS */

  /**
   * Capture a snapshot of the current canvas state and push to undo stack.
   * Clears the redo stack (new mutation invalidates redo history).
   * Must be called BEFORE the mutation happens.
   * No-op if a batch restoration is in progress.
   *
   * @param {string} label - Human-readable action label (e.g., 'Add node')
   * @returns {void}
   */
  function pushSnapshot(label: string): void {
    if (isBatchRestoring.value) return;

    const snapshot: HistorySnapshot = {
      nodes: JSON.parse(JSON.stringify(toRaw(nodes.value))),
      edges: JSON.parse(JSON.stringify(toRaw(edges.value))),
      label,
    };

    undoStack.value.push(snapshot);

    // Trim oldest entries if exceeding max size
    if (undoStack.value.length > MAX_HISTORY_SIZE) {
      undoStack.value = undoStack.value.slice(-MAX_HISTORY_SIZE);
    }

    // New mutation invalidates redo history
    redoStack.value = [];
  }

  /**
   * Undo the last mutation.
   * Saves current state to redo stack, pops from undo stack,
   * restores composable state, and sets isBatchRestoring flag.
   * Caller MUST call finishRestore() after Vue Flow sync (via nextTick).
   *
   * @returns {HistorySnapshot | null} Restored snapshot, or null if nothing to undo
   */
  function undo(): HistorySnapshot | null {
    if (!undoStack.value.length) return null;

    // Save current state to redo stack before overwriting
    redoStack.value.push({
      nodes: JSON.parse(JSON.stringify(toRaw(nodes.value))),
      edges: JSON.parse(JSON.stringify(toRaw(edges.value))),
      label: 'redo',
    });

    const snapshot = undoStack.value.pop()!;

    // Set flag to guard watchers during restore
    isBatchRestoring.value = true;

    // Restore composable state
    nodes.value = snapshot.nodes;
    edges.value = snapshot.edges;

    return snapshot;
  }

  /**
   * Redo the last undone mutation.
   * Saves current state to undo stack, pops from redo stack,
   * restores composable state, and sets isBatchRestoring flag.
   * Caller MUST call finishRestore() after Vue Flow sync (via nextTick).
   *
   * @returns {HistorySnapshot | null} Restored snapshot, or null if nothing to redo
   */
  function redo(): HistorySnapshot | null {
    if (!redoStack.value.length) return null;

    // Save current state to undo stack before overwriting
    undoStack.value.push({
      nodes: JSON.parse(JSON.stringify(toRaw(nodes.value))),
      edges: JSON.parse(JSON.stringify(toRaw(edges.value))),
      label: 'undo',
    });

    const snapshot = redoStack.value.pop()!;

    // Set flag to guard watchers during restore
    isBatchRestoring.value = true;

    // Restore composable state
    nodes.value = snapshot.nodes;
    edges.value = snapshot.edges;

    return snapshot;
  }

  /**
   * Reset the isBatchRestoring flag after Vue Flow sync completes.
   * Must be called in nextTick after restoreSnapshot/syncSnapshotToVueFlow.
   *
   * @returns {void}
   */
  function finishRestore(): void {
    isBatchRestoring.value = false;
  }

  /**
   * Clear all undo/redo history.
   * Called when loading a new workflow or resetting editor state.
   *
   * @returns {void}
   */
  function clearHistory(): void {
    undoStack.value = [];
    redoStack.value = [];
  }

  return {
    // State
    canUndo,
    canRedo,
    isBatchRestoring,

    // Methods
    pushSnapshot,
    undo,
    redo,
    finishRestore,
    clearHistory,
  };
}
