/** TYPE IMPORTS */
import type { WorkflowNode, WorkflowEdge } from '../interfaces';

import { describe, it, expect, beforeEach, vi } from 'vitest';
import { setActivePinia, createPinia } from 'pinia';
import { defineComponent, createApp } from 'vue';

/**
 * Mock the plugin registry store (required by useWorkflowEditorState)
 */
vi.mock('@stores/pluginRegistry', () => ({
  usePluginRegistryStore: () => ({
    getNodeType: vi.fn(),
  }),
}));

/**
 * Mock utils used by useWorkflowEditorState
 */
vi.mock('../utils', () => ({
  resolveNodeHandles: vi.fn().mockReturnValue({ inputs: [], outputs: [] }),
  findOrphanedEdges: vi.fn().mockReturnValue([]),
  detectTightCycles: vi.fn().mockReturnValue(new Set<string>()),
}));

/**
 * Helper to invoke a composable inside a real Vue setup context.
 *
 * @param {() => T} composable - Composable factory
 * @returns {{ result: T; app: ReturnType<typeof createApp> }}
 */
function withSetup<T>(composable: () => T): { result: T; app: ReturnType<typeof createApp> } {
  let result!: T;
  const app = createApp(
    defineComponent({
      setup() {
        result = composable();
        return () => null;
      },
    }),
  );
  app.mount(document.createElement('div'));
  return { result, app };
}

/** Lazy imports — resolved after vi.mock */
// eslint-disable-next-line @typescript-eslint/consistent-type-imports
let useWorkflowHistory: typeof import('./useWorkflowHistory').useWorkflowHistory;
// eslint-disable-next-line @typescript-eslint/consistent-type-imports
let useWorkflowEditorState: typeof import('./useWorkflowEditorState').useWorkflowEditorState;

beforeEach(async () => {
  setActivePinia(createPinia());

  /**
   * Re-import to reset module-level singleton refs (undoStack, redoStack, isBatchRestoring).
   * Both composables share singletons so both must come from the same fresh module graph.
   */
  vi.resetModules();
  const historyMod = await import('./useWorkflowHistory');
  const editorMod = await import('./useWorkflowEditorState');
  useWorkflowHistory = historyMod.useWorkflowHistory;
  useWorkflowEditorState = editorMod.useWorkflowEditorState;
});

// ────────────────────────────────────────────────────────────────────────────
// Helpers
// ────────────────────────────────────────────────────────────────────────────

/**
 * Create a minimal WorkflowNode
 */
function makeNode(overrides: Partial<WorkflowNode> = {}): WorkflowNode {
  return {
    id: `n_${Math.random().toString(36).slice(2, 8)}`,
    type: 'core/delay',
    position: { x: 100, y: 100 },
    config: {},
    label: 'Delay',
    ...overrides,
  };
}

/**
 * Create a minimal WorkflowEdge
 */
function makeEdge(overrides: Partial<WorkflowEdge> = {}): WorkflowEdge {
  return {
    id: `e_${Math.random().toString(36).slice(2, 8)}`,
    source: 'a',
    target: 'b',
    ...overrides,
  };
}

// ════════════════════════════════════════════════════════════════════════════
// TESTS
// ════════════════════════════════════════════════════════════════════════════

describe('useWorkflowHistory', () => {
  // ────────────────────────────────────────────────────────────────────────
  // Initial state
  // ────────────────────────────────────────────────────────────────────────
  describe('initial state', () => {
    it('canUndo and canRedo are false', () => {
      const { result } = withSetup(() => useWorkflowHistory());

      expect(result.canUndo.value).toBe(false);
      expect(result.canRedo.value).toBe(false);
    });

    it('isBatchRestoring is false', () => {
      const { result } = withSetup(() => useWorkflowHistory());

      expect(result.isBatchRestoring.value).toBe(false);
    });
  });

  // ────────────────────────────────────────────────────────────────────────
  // pushSnapshot
  // ────────────────────────────────────────────────────────────────────────
  describe('pushSnapshot', () => {
    it('makes canUndo true after pushing', () => {
      const { result: editor } = withSetup(() => useWorkflowEditorState());
      const { result: history } = withSetup(() => useWorkflowHistory());

      editor.nodes.value = [makeNode({ id: 'n1' })];
      history.pushSnapshot('Add node');

      expect(history.canUndo.value).toBe(true);
    });

    it('clears redo stack on new push', () => {
      const { result: editor } = withSetup(() => useWorkflowEditorState());
      const { result: history } = withSetup(() => useWorkflowHistory());

      // Push, undo, then push again
      editor.nodes.value = [makeNode({ id: 'n1' })];
      history.pushSnapshot('Step 1');

      editor.nodes.value = [makeNode({ id: 'n2' })];
      history.undo();
      history.finishRestore();

      expect(history.canRedo.value).toBe(true);

      // New push invalidates redo
      history.pushSnapshot('Step 2');

      expect(history.canRedo.value).toBe(false);
    });

    it('deep-clones nodes and edges (mutation-safe)', () => {
      const { result: editor } = withSetup(() => useWorkflowEditorState());
      const { result: history } = withSetup(() => useWorkflowHistory());

      const node = makeNode({ id: 'n1', config: { delay: 1000 } });
      editor.nodes.value = [node];

      history.pushSnapshot('Before mutation');

      // Mutate the original node
      editor.nodes.value[0]!.config.delay = 9999;

      // Undo — should restore original config
      const snapshot = history.undo();

      expect(snapshot).not.toBeNull();
      expect(snapshot!.nodes[0]!.config.delay).toBe(1000);
    });

    it('is a no-op during batch restoring', () => {
      const { result: editor } = withSetup(() => useWorkflowEditorState());
      const { result: history } = withSetup(() => useWorkflowHistory());

      // Simulate a batch restore in progress
      editor.nodes.value = [makeNode({ id: 'n1' })];
      history.pushSnapshot('Step 1');

      history.undo(); // sets isBatchRestoring = true

      // This push should be ignored
      history.pushSnapshot('Should be ignored');

      // The undo stack should only contain what we pushed before undo
      // (undo pops one from undo stack, so after pushSnapshot('Step 1') then undo(),
      //  the undo stack is empty)
      expect(history.canUndo.value).toBe(false);

      history.finishRestore();
    });
  });

  // ────────────────────────────────────────────────────────────────────────
  // undo
  // ────────────────────────────────────────────────────────────────────────
  describe('undo', () => {
    it('returns null when undo stack is empty', () => {
      const { result: history } = withSetup(() => useWorkflowHistory());

      const snapshot = history.undo();

      expect(snapshot).toBeNull();
    });

    it('restores the previous state', () => {
      const { result: editor } = withSetup(() => useWorkflowEditorState());
      const { result: history } = withSetup(() => useWorkflowHistory());

      // Initial state: 1 node
      editor.nodes.value = [makeNode({ id: 'n1' })];
      editor.edges.value = [];

      history.pushSnapshot('Before adding n2');

      // Add a second node
      editor.nodes.value = [makeNode({ id: 'n1' }), makeNode({ id: 'n2' })];

      // Undo
      const snapshot = history.undo();
      history.finishRestore();

      expect(snapshot).not.toBeNull();
      expect(editor.nodes.value).toHaveLength(1);
      expect(editor.nodes.value[0]!.id).toBe('n1');
    });

    it('saves current state to redo stack', () => {
      const { result: editor } = withSetup(() => useWorkflowEditorState());
      const { result: history } = withSetup(() => useWorkflowHistory());

      editor.nodes.value = [makeNode({ id: 'n1' })];
      history.pushSnapshot('Step 1');

      editor.nodes.value = [makeNode({ id: 'n1' }), makeNode({ id: 'n2' })];

      history.undo();
      history.finishRestore();

      expect(history.canRedo.value).toBe(true);
    });

    it('sets isBatchRestoring to true during restore', () => {
      const { result: editor } = withSetup(() => useWorkflowEditorState());
      const { result: history } = withSetup(() => useWorkflowHistory());

      editor.nodes.value = [makeNode({ id: 'n1' })];
      history.pushSnapshot('Step 1');

      editor.nodes.value = [];

      history.undo();

      expect(history.isBatchRestoring.value).toBe(true);

      history.finishRestore();
      expect(history.isBatchRestoring.value).toBe(false);
    });
  });

  // ────────────────────────────────────────────────────────────────────────
  // redo
  // ────────────────────────────────────────────────────────────────────────
  describe('redo', () => {
    it('returns null when redo stack is empty', () => {
      const { result: history } = withSetup(() => useWorkflowHistory());

      const snapshot = history.redo();

      expect(snapshot).toBeNull();
    });

    it('re-applies the undone state', () => {
      const { result: editor } = withSetup(() => useWorkflowEditorState());
      const { result: history } = withSetup(() => useWorkflowHistory());

      // Step 1: 1 node
      editor.nodes.value = [makeNode({ id: 'n1' })];
      history.pushSnapshot('Before adding n2');

      // Step 2: 2 nodes
      const n2 = makeNode({ id: 'n2' });
      editor.nodes.value = [makeNode({ id: 'n1' }), n2];

      // Undo back to 1 node
      history.undo();
      history.finishRestore();

      expect(editor.nodes.value).toHaveLength(1);

      // Redo back to 2 nodes
      history.redo();
      history.finishRestore();

      expect(editor.nodes.value).toHaveLength(2);
    });

    it('saves current state to undo stack before redo', () => {
      const { result: editor } = withSetup(() => useWorkflowEditorState());
      const { result: history } = withSetup(() => useWorkflowHistory());

      editor.nodes.value = [makeNode({ id: 'n1' })];
      history.pushSnapshot('Step 1');

      editor.nodes.value = [];
      history.undo();
      history.finishRestore();

      // Redo
      history.redo();
      history.finishRestore();

      // Should be able to undo again
      expect(history.canUndo.value).toBe(true);
    });

    it('sets isBatchRestoring during redo', () => {
      const { result: editor } = withSetup(() => useWorkflowEditorState());
      const { result: history } = withSetup(() => useWorkflowHistory());

      editor.nodes.value = [makeNode({ id: 'n1' })];
      history.pushSnapshot('Step 1');
      editor.nodes.value = [];

      history.undo();
      history.finishRestore();

      history.redo();

      expect(history.isBatchRestoring.value).toBe(true);

      history.finishRestore();
      expect(history.isBatchRestoring.value).toBe(false);
    });
  });

  // ────────────────────────────────────────────────────────────────────────
  // undo/redo chain
  // ────────────────────────────────────────────────────────────────────────
  describe('undo/redo chain', () => {
    it('supports multiple undo/redo cycles correctly', () => {
      const { result: editor } = withSetup(() => useWorkflowEditorState());
      const { result: history } = withSetup(() => useWorkflowHistory());

      // Build up 3 states
      editor.nodes.value = [];
      history.pushSnapshot('Empty');

      editor.nodes.value = [makeNode({ id: 'n1' })];
      history.pushSnapshot('One node');

      editor.nodes.value = [makeNode({ id: 'n1' }), makeNode({ id: 'n2' })];
      history.pushSnapshot('Two nodes');

      editor.nodes.value = [makeNode({ id: 'n1' }), makeNode({ id: 'n2' }), makeNode({ id: 'n3' })];

      // Undo 3 times
      history.undo();
      history.finishRestore();
      expect(editor.nodes.value).toHaveLength(2);

      history.undo();
      history.finishRestore();
      expect(editor.nodes.value).toHaveLength(1);

      history.undo();
      history.finishRestore();
      expect(editor.nodes.value).toHaveLength(0);

      // Redo 2 times
      history.redo();
      history.finishRestore();
      expect(editor.nodes.value).toHaveLength(1);

      history.redo();
      history.finishRestore();
      expect(editor.nodes.value).toHaveLength(2);
    });
  });

  // ────────────────────────────────────────────────────────────────────────
  // History stack limit
  // ────────────────────────────────────────────────────────────────────────
  describe('history stack limit', () => {
    it('trims undo stack to MAX_HISTORY_SIZE (50)', () => {
      const { result: editor } = withSetup(() => useWorkflowEditorState());
      const { result: history } = withSetup(() => useWorkflowHistory());

      // Push 60 snapshots
      for (let i = 0; i < 60; i++) {
        editor.nodes.value = [makeNode({ id: `n_${i}` })];
        history.pushSnapshot(`Step ${i}`);
      }

      // Should be capped at 50
      // We can verify canUndo is true, and by doing 50 undos we exhaust the stack
      let undoCount = 0;
      while (history.canUndo.value) {
        history.undo();
        history.finishRestore();
        undoCount++;
      }

      expect(undoCount).toBe(50);
    });
  });

  // ────────────────────────────────────────────────────────────────────────
  // finishRestore
  // ────────────────────────────────────────────────────────────────────────
  describe('finishRestore', () => {
    it('resets isBatchRestoring to false', () => {
      const { result: editor } = withSetup(() => useWorkflowEditorState());
      const { result: history } = withSetup(() => useWorkflowHistory());

      editor.nodes.value = [makeNode()];
      history.pushSnapshot('Step 1');
      editor.nodes.value = [];

      history.undo();
      expect(history.isBatchRestoring.value).toBe(true);

      history.finishRestore();
      expect(history.isBatchRestoring.value).toBe(false);
    });
  });

  // ────────────────────────────────────────────────────────────────────────
  // clearHistory
  // ────────────────────────────────────────────────────────────────────────
  describe('clearHistory', () => {
    it('empties both undo and redo stacks', () => {
      const { result: editor } = withSetup(() => useWorkflowEditorState());
      const { result: history } = withSetup(() => useWorkflowHistory());

      editor.nodes.value = [makeNode()];
      history.pushSnapshot('Step 1');

      editor.nodes.value = [];
      history.undo();
      history.finishRestore();

      expect(history.canUndo.value).toBe(false);
      expect(history.canRedo.value).toBe(true);

      history.clearHistory();

      expect(history.canUndo.value).toBe(false);
      expect(history.canRedo.value).toBe(false);
    });
  });

  // ────────────────────────────────────────────────────────────────────────
  // Snapshot content integrity
  // ────────────────────────────────────────────────────────────────────────
  describe('snapshot content', () => {
    it('stores both nodes and edges in the snapshot', () => {
      const { result: editor } = withSetup(() => useWorkflowEditorState());
      const { result: history } = withSetup(() => useWorkflowHistory());

      const node = makeNode({ id: 'n1' });
      const edge = makeEdge({ id: 'e1', source: 'n1', target: 'n2' });

      editor.nodes.value = [node];
      editor.edges.value = [edge];

      history.pushSnapshot('Initial state');

      // Modify state
      editor.nodes.value = [];
      editor.edges.value = [];

      // Undo restores both
      const snapshot = history.undo();
      history.finishRestore();

      expect(snapshot).not.toBeNull();
      expect(snapshot!.nodes).toHaveLength(1);
      expect(snapshot!.nodes[0]!.id).toBe('n1');
      expect(snapshot!.edges).toHaveLength(1);
      expect(snapshot!.edges[0]!.id).toBe('e1');
    });

    it('snapshot label is preserved', () => {
      const { result: editor } = withSetup(() => useWorkflowEditorState());
      const { result: history } = withSetup(() => useWorkflowHistory());

      editor.nodes.value = [makeNode()];
      history.pushSnapshot('Delete node');

      editor.nodes.value = [];

      const snapshot = history.undo();

      expect(snapshot!.label).toBe('Delete node');
    });
  });
});
