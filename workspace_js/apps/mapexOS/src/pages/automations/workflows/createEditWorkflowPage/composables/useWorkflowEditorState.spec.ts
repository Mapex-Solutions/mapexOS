/** TYPE IMPORTS */
import type {
  WorkflowVariable,
  CaptureField,
  ExternalVariable,
  WorkflowNode,
  WorkflowEdge,
  WorkflowDefinition,
} from '../interfaces';

import { describe, it, expect, beforeEach, vi } from 'vitest';
import { setActivePinia, createPinia } from 'pinia';
import { defineComponent, createApp } from 'vue';

/**
 * Mock the plugin registry store.
 * Individual tests override `getNodeType` as needed.
 */
const mockGetNodeType = vi.fn();

vi.mock('@stores/pluginRegistry', () => ({
  usePluginRegistryStore: () => ({
    getNodeType: mockGetNodeType,
  }),
}));

/**
 * Mock utils used by updateNodeConfig and validateAllNodes
 */
const mockResolveNodeHandles = vi.fn().mockReturnValue({ inputs: [], outputs: [] });
const mockFindOrphanedEdges = vi.fn().mockReturnValue([]);
const mockDetectTightCycles = vi.fn().mockReturnValue(new Set<string>());

vi.mock('../utils', () => ({
  resolveNodeHandles: (...args: unknown[]) => mockResolveNodeHandles(...args),
  findOrphanedEdges: (...args: unknown[]) => mockFindOrphanedEdges(...args),
  detectTightCycles: (...args: unknown[]) => mockDetectTightCycles(...args),
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

/** Lazy import — resolved after vi.mock */
// eslint-disable-next-line @typescript-eslint/consistent-type-imports
let useWorkflowEditorState: typeof import('./useWorkflowEditorState').useWorkflowEditorState;

beforeEach(async () => {
  setActivePinia(createPinia());
  mockGetNodeType.mockReset();
  mockResolveNodeHandles.mockReset().mockReturnValue({ inputs: [], outputs: [] });
  mockFindOrphanedEdges.mockReset().mockReturnValue([]);
  mockDetectTightCycles.mockReset().mockReturnValue(new Set<string>());

  /**
   * Re-import to reset module-level singleton refs.
   * vi.resetModules clears the module cache so module-level state is fresh.
   */
  vi.resetModules();
  const mod = await import('./useWorkflowEditorState');
  useWorkflowEditorState = mod.useWorkflowEditorState;
});

// ────────────────────────────────────────────────────────────────────────────
// Helpers
// ────────────────────────────────────────────────────────────────────────────

/**
 * Create a minimal WorkflowNode
 */
function makeNode(overrides: Partial<WorkflowNode> = {}): WorkflowNode {
  return {
    id: `n_${Date.now()}`,
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
    id: `e_${Date.now()}`,
    source: 'a',
    target: 'b',
    ...overrides,
  };
}

/**
 * Create a minimal WorkflowVariable
 */
function makeVariable(overrides: Partial<WorkflowVariable> = {}): WorkflowVariable {
  return {
    field: 'counter',
    type: 'number',
    defaultValue: 0,
    durable: false,
    ...overrides,
  };
}

/**
 * Create a minimal CaptureField
 */
function makeCaptureField(overrides: Partial<CaptureField> = {}): CaptureField {
  return {
    field: 'temperature',
    type: 'number',
    description: 'Sensor temperature',
    ...overrides,
  };
}

/**
 * Create a minimal ExternalVariable
 */
function makeExternalInput(overrides: Partial<ExternalVariable> = {}): ExternalVariable {
  return {
    field: 'assetId',
    label: 'Asset ID',
    icon: 'fingerprint',
    type: 'string',
    defaultValue: '',
    required: false,
    ...overrides,
  };
}

// ════════════════════════════════════════════════════════════════════════════
// TESTS
// ════════════════════════════════════════════════════════════════════════════

describe('useWorkflowEditorState', () => {
  // ────────────────────────────────────────────────────────────────────────
  // Initial state
  // ────────────────────────────────────────────────────────────────────────
  describe('initial state', () => {
    it('provides empty arrays for nodes, edges, states, captureFields, externalInputs', () => {
      const { result } = withSetup(() => useWorkflowEditorState());

      expect(result.nodes.value).toEqual([]);
      expect(result.edges.value).toEqual([]);
      expect(result.states.value).toEqual([]);
      expect(result.captureFields.value).toEqual([]);
      expect(result.externalInputs.value).toEqual([]);
      expect(result.installedPlugins.value).toEqual([]);
    });

    it('has default general settings', () => {
      const { result } = withSetup(() => useWorkflowEditorState());
      const gs = result.generalSettings.value;

      expect(gs.name).toBe('');
      expect(gs.description).toBe('');
      expect(gs.enabled).toBe(true);
      expect(gs.isTemplate).toBe(false);
    });

    it('variablesCount is zero initially', () => {
      const { result } = withSetup(() => useWorkflowEditorState());
      expect(result.variablesCount.value).toBe(0);
    });

    it('nodesCount is zero initially', () => {
      const { result } = withSetup(() => useWorkflowEditorState());
      expect(result.nodesCount.value).toBe(0);
    });

    it('viewport defaults to x:0 y:0 zoom:1', () => {
      const { result } = withSetup(() => useWorkflowEditorState());
      expect(result.viewport.value).toEqual({ x: 0, y: 0, zoom: 1 });
    });

    it('nodeConfigVersion starts at 0', () => {
      const { result } = withSetup(() => useWorkflowEditorState());
      expect(result.nodeConfigVersion.value).toBe(0);
    });

    it('nodeValidationErrors is empty', () => {
      const { result } = withSetup(() => useWorkflowEditorState());
      expect(result.nodeValidationErrors.value).toEqual({});
    });
  });

  // ────────────────────────────────────────────────────────────────────────
  // resetAllStates
  // ────────────────────────────────────────────────────────────────────────
  describe('resetAllStates', () => {
    it('resets to defaults and adds a Start node', () => {
      const { result } = withSetup(() => useWorkflowEditorState());

      // Dirty the state first
      result.nodes.value = [makeNode()];
      result.edges.value = [makeEdge()];
      result.states.value = [makeVariable()];

      result.resetAllStates();

      expect(result.nodes.value).toHaveLength(1);
      expect(result.nodes.value[0]!.type).toBe('core/start');
      expect(result.edges.value).toEqual([]);
      expect(result.states.value).toEqual([]);
      expect(result.captureFields.value).toEqual([]);
      expect(result.externalInputs.value).toEqual([]);
      expect(result.generalSettings.value.name).toBe('');
    });
  });

  // ────────────────────────────────────────────────────────────────────────
  // setAllStates
  // ────────────────────────────────────────────────────────────────────────
  describe('setAllStates', () => {
    it('loads a complete workflow definition', () => {
      const { result } = withSetup(() => useWorkflowEditorState());

      const workflow: WorkflowDefinition = {
        name: 'My Workflow',
        description: 'Test',
        enabled: true,
        isTemplate: false,
        definitionVersion: 1,
        timezone: { type: 'literal', value: 'America/Sao_Paulo' },
        retryPolicy: {
          enabled: true,
          maxAttempts: 5,
          initialInterval: '2s',
          backoffMultiplier: 1.5,
          maxInterval: '10m',
          nonRetryableErrors: [],
        },
        states: [makeVariable()],
        captureFields: [makeCaptureField()],
        externalInputs: [makeExternalInput()],
        externalSignals: [],
        installedPlugins: ['telegram'],
        nodes: [
          { id: '__start__', type: 'core/start', position: { x: 0, y: 0 }, config: {} },
          makeNode({ id: 'n1' }),
        ],
        edges: [makeEdge({ id: 'e1', source: '__start__', target: 'n1' })],
        metadata: { canvasViewport: { x: 100, y: 200, zoom: 1.5 } },
      };

      result.setAllStates(workflow);

      expect(result.generalSettings.value.name).toBe('My Workflow');
      expect(result.generalSettings.value.timezone.value).toBe('America/Sao_Paulo');
      expect(result.states.value).toHaveLength(1);
      expect(result.captureFields.value).toHaveLength(1);
      expect(result.externalInputs.value).toHaveLength(1);
      expect(result.installedPlugins.value).toEqual(['telegram']);
      expect(result.nodes.value).toHaveLength(2);
      expect(result.edges.value).toHaveLength(1);
      expect(result.viewport.value).toEqual({ x: 100, y: 200, zoom: 1.5 });
    });

    it('injects a Start node if not present in loaded definition', () => {
      const { result } = withSetup(() => useWorkflowEditorState());

      const workflow: WorkflowDefinition = {
        name: 'No Start',
        description: '',
        enabled: true,
        isTemplate: false,
        definitionVersion: 1,
        timezone: { type: 'literal', value: 'UTC' },
        retryPolicy: {
          enabled: false,
          maxAttempts: 3,
          initialInterval: '1s',
          backoffMultiplier: 2.0,
          maxInterval: '5m',
          nonRetryableErrors: [],
        },
        states: [],
        captureFields: [],
        externalInputs: [],
        externalSignals: [],
        installedPlugins: [],
        nodes: [makeNode({ id: 'n1' })],
        edges: [],
        metadata: { canvasViewport: { x: 0, y: 0, zoom: 1 } },
      };

      result.setAllStates(workflow);

      expect(result.nodes.value[0]!.type).toBe('core/start');
      expect(result.nodes.value).toHaveLength(2);
    });
  });

  // ────────────────────────────────────────────────────────────────────────
  // updateGeneral
  // ────────────────────────────────────────────────────────────────────────
  describe('updateGeneral', () => {
    it('merges partial settings', () => {
      const { result } = withSetup(() => useWorkflowEditorState());

      result.updateGeneral({ name: 'Updated', enabled: false });

      expect(result.generalSettings.value.name).toBe('Updated');
      expect(result.generalSettings.value.enabled).toBe(false);
      // Untouched fields remain at default
      expect(result.generalSettings.value.description).toBe('');
    });
  });

  // ────────────────────────────────────────────────────────────────────────
  // Node operations
  // ────────────────────────────────────────────────────────────────────────
  describe('addNode', () => {
    it('appends a node to the list', () => {
      const { result } = withSetup(() => useWorkflowEditorState());
      const node = makeNode({ id: 'new1' });

      result.addNode(node);

      expect(result.nodes.value).toHaveLength(1);
      expect(result.nodes.value[0]!.id).toBe('new1');
      expect(result.nodesCount.value).toBe(1);
    });
  });

  describe('removeNode', () => {
    it('removes the node and its connected edges', () => {
      const { result } = withSetup(() => useWorkflowEditorState());
      const n1 = makeNode({ id: 'n1', type: 'core/delay' });
      const n2 = makeNode({ id: 'n2', type: 'core/delay' });
      const edge = makeEdge({ id: 'e1', source: 'n1', target: 'n2' });

      result.nodes.value = [n1, n2];
      result.edges.value = [edge];

      mockGetNodeType.mockReturnValue({ deletable: true });

      const removed = result.removeNode('n1');

      expect(removed).toContain('n1');
      expect(result.nodes.value.map(n => n.id)).not.toContain('n1');
      expect(result.edges.value).toHaveLength(0);
    });

    it('removes child nodes (e.g., text notes) along with the parent', () => {
      const { result } = withSetup(() => useWorkflowEditorState());

      const parent = makeNode({ id: 'parent', type: 'core/delay' });
      const childNote = makeNode({ id: 'note1', type: 'core/text_note', parentNodeId: 'parent' });
      const noteEdge = makeEdge({ id: 'ne1', source: 'note1', target: 'parent' });

      result.nodes.value = [parent, childNote];
      result.edges.value = [noteEdge];

      mockGetNodeType.mockReturnValue({ deletable: true });

      const removed = result.removeNode('parent');

      expect(removed).toEqual(['parent', 'note1']);
      expect(result.nodes.value).toHaveLength(0);
      expect(result.edges.value).toHaveLength(0);
    });

    it('blocks deletion of undeletable nodes (e.g., Start)', () => {
      const { result } = withSetup(() => useWorkflowEditorState());
      const startNode = makeNode({ id: '__start__', type: 'core/start' });

      result.nodes.value = [startNode];
      mockGetNodeType.mockReturnValue({ deletable: false });

      const removed = result.removeNode('__start__');

      expect(removed).toEqual([]);
      expect(result.nodes.value).toHaveLength(1);
    });

    it('returns empty array for non-existent node', () => {
      const { result } = withSetup(() => useWorkflowEditorState());
      const removed = result.removeNode('ghost');
      expect(removed).toEqual([]);
    });
  });

  describe('updateNodeConfig', () => {
    it('merges config and increments nodeConfigVersion', () => {
      const { result } = withSetup(() => useWorkflowEditorState());

      const node = makeNode({ id: 'n1', config: { delay: 1000 } });
      result.nodes.value = [node];

      mockGetNodeType.mockReturnValue({});

      const versionBefore = result.nodeConfigVersion.value;
      result.updateNodeConfig('n1', { delay: 2000, unit: 'ms' });

      expect(result.nodes.value[0]!.config).toEqual({ delay: 2000, unit: 'ms' });
      expect(result.nodeConfigVersion.value).toBe(versionBefore + 1);
    });

    it('removes orphaned edges when node has dynamic resolvers', () => {
      const { result } = withSetup(() => useWorkflowEditorState());

      const node = makeNode({ id: 'n1' });
      const orphanEdge = makeEdge({ id: 'orphan', source: 'n1', target: 'n2' });
      result.nodes.value = [node];
      result.edges.value = [orphanEdge];

      mockGetNodeType.mockReturnValue({ resolveOutputs: vi.fn() });
      mockResolveNodeHandles.mockReturnValue({ inputs: [], outputs: [{ id: 'out_new' }] });
      mockFindOrphanedEdges.mockReturnValue(['orphan']);

      result.updateNodeConfig('n1', { branches: 3 });

      expect(result.edges.value).toHaveLength(0);
    });

    it('is a no-op for non-existent node', () => {
      const { result } = withSetup(() => useWorkflowEditorState());
      result.updateNodeConfig('ghost', { foo: 'bar' });
      // Should not throw
      expect(result.nodeConfigVersion.value).toBe(0);
    });
  });

  describe('duplicateNode', () => {
    it('creates a copy offset by 40px', () => {
      const { result } = withSetup(() => useWorkflowEditorState());

      const node = makeNode({ id: 'n1', type: 'core/delay', position: { x: 100, y: 200 }, config: { delay: 5 } });
      result.nodes.value = [node];

      mockGetNodeType.mockReturnValue({ deletable: true, label: 'Delay' });

      const newId = result.duplicateNode('n1');

      expect(newId).toBeTruthy();
      expect(result.nodes.value).toHaveLength(2);

      const clone = result.nodes.value.find(n => n.id === newId);
      expect(clone).toBeDefined();
      expect(clone!.position.x).toBe(140);
      expect(clone!.position.y).toBe(240);
      expect(clone!.config).toEqual({ delay: 5 });
    });

    it('returns null for undeletable node types', () => {
      const { result } = withSetup(() => useWorkflowEditorState());

      result.nodes.value = [makeNode({ id: 'start', type: 'core/start' })];
      mockGetNodeType.mockReturnValue({ deletable: false });

      const newId = result.duplicateNode('start');
      expect(newId).toBeNull();
      expect(result.nodes.value).toHaveLength(1);
    });

    it('returns null for non-existent node', () => {
      const { result } = withSetup(() => useWorkflowEditorState());
      expect(result.duplicateNode('ghost')).toBeNull();
    });
  });

  // ────────────────────────────────────────────────────────────────────────
  // Edge operations
  // ────────────────────────────────────────────────────────────────────────
  describe('addEdge', () => {
    it('appends an edge', () => {
      const { result } = withSetup(() => useWorkflowEditorState());
      const edge = makeEdge({ id: 'e1' });

      result.addEdge(edge);

      expect(result.edges.value).toHaveLength(1);
      expect(result.edges.value[0]!.id).toBe('e1');
    });
  });

  describe('removeEdge', () => {
    it('removes an edge by ID', () => {
      const { result } = withSetup(() => useWorkflowEditorState());

      result.edges.value = [
        makeEdge({ id: 'e1' }),
        makeEdge({ id: 'e2' }),
      ];

      result.removeEdge('e1');

      expect(result.edges.value).toHaveLength(1);
      expect(result.edges.value[0]!.id).toBe('e2');
    });

    it('is a no-op if edge does not exist', () => {
      const { result } = withSetup(() => useWorkflowEditorState());
      result.edges.value = [makeEdge({ id: 'e1' })];

      result.removeEdge('ghost');

      expect(result.edges.value).toHaveLength(1);
    });
  });

  describe('updateNodes / updateEdges', () => {
    it('replaces entire nodes array', () => {
      const { result } = withSetup(() => useWorkflowEditorState());
      const newNodes = [makeNode({ id: 'a' }), makeNode({ id: 'b' })];

      result.updateNodes(newNodes);

      expect(result.nodes.value).toHaveLength(2);
    });

    it('replaces entire edges array', () => {
      const { result } = withSetup(() => useWorkflowEditorState());
      const newEdges = [makeEdge({ id: 'e1' })];

      result.updateEdges(newEdges);

      expect(result.edges.value).toHaveLength(1);
    });
  });

  // ────────────────────────────────────────────────────────────────────────
  // State variables CRUD
  // ────────────────────────────────────────────────────────────────────────
  describe('state variables (states)', () => {
    it('addState pushes a variable', () => {
      const { result } = withSetup(() => useWorkflowEditorState());
      result.addState(makeVariable({ field: 'count' }));

      expect(result.states.value).toHaveLength(1);
      expect(result.states.value[0]!.field).toBe('count');
    });

    it('updateState replaces at index', () => {
      const { result } = withSetup(() => useWorkflowEditorState());

      result.addState(makeVariable({ field: 'a' }));
      result.addState(makeVariable({ field: 'b' }));

      result.updateState(1, makeVariable({ field: 'b_updated', type: 'string' }));

      expect(result.states.value[1]!.field).toBe('b_updated');
    });

    it('updateState ignores invalid index', () => {
      const { result } = withSetup(() => useWorkflowEditorState());
      result.addState(makeVariable());

      result.updateState(5, makeVariable({ field: 'nope' }));
      result.updateState(-1, makeVariable({ field: 'nope' }));

      expect(result.states.value).toHaveLength(1);
    });

    it('removeState splices at index', () => {
      const { result } = withSetup(() => useWorkflowEditorState());
      result.addState(makeVariable({ field: 'a' }));
      result.addState(makeVariable({ field: 'b' }));
      result.addState(makeVariable({ field: 'c' }));

      result.removeState(1);

      expect(result.states.value.map(s => s.field)).toEqual(['a', 'c']);
    });

    it('removeState ignores invalid index', () => {
      const { result } = withSetup(() => useWorkflowEditorState());
      result.addState(makeVariable());

      result.removeState(10);
      result.removeState(-1);

      expect(result.states.value).toHaveLength(1);
    });

    it('moveState swaps up', () => {
      const { result } = withSetup(() => useWorkflowEditorState());
      result.addState(makeVariable({ field: 'a' }));
      result.addState(makeVariable({ field: 'b' }));

      result.moveState(1, 'up');

      expect(result.states.value.map(s => s.field)).toEqual(['b', 'a']);
    });

    it('moveState swaps down', () => {
      const { result } = withSetup(() => useWorkflowEditorState());
      result.addState(makeVariable({ field: 'a' }));
      result.addState(makeVariable({ field: 'b' }));

      result.moveState(0, 'down');

      expect(result.states.value.map(s => s.field)).toEqual(['b', 'a']);
    });

    it('moveState is a no-op at boundaries', () => {
      const { result } = withSetup(() => useWorkflowEditorState());
      result.addState(makeVariable({ field: 'a' }));
      result.addState(makeVariable({ field: 'b' }));

      result.moveState(0, 'up');
      expect(result.states.value.map(s => s.field)).toEqual(['a', 'b']);

      result.moveState(1, 'down');
      expect(result.states.value.map(s => s.field)).toEqual(['a', 'b']);
    });
  });

  // ────────────────────────────────────────────────────────────────────────
  // Capture fields CRUD
  // ────────────────────────────────────────────────────────────────────────
  describe('capture fields', () => {
    it('addCaptureField pushes a field', () => {
      const { result } = withSetup(() => useWorkflowEditorState());
      result.addCaptureField(makeCaptureField({ field: 'temp' }));

      expect(result.captureFields.value).toHaveLength(1);
    });

    it('updateCaptureField replaces at index', () => {
      const { result } = withSetup(() => useWorkflowEditorState());
      result.addCaptureField(makeCaptureField({ field: 'temp' }));
      result.updateCaptureField(0, makeCaptureField({ field: 'humidity' }));

      expect(result.captureFields.value[0]!.field).toBe('humidity');
    });

    it('updateCaptureField ignores invalid index', () => {
      const { result } = withSetup(() => useWorkflowEditorState());
      result.addCaptureField(makeCaptureField());
      result.updateCaptureField(99, makeCaptureField({ field: 'nope' }));

      expect(result.captureFields.value).toHaveLength(1);
    });

    it('removeCaptureField splices at index', () => {
      const { result } = withSetup(() => useWorkflowEditorState());
      result.addCaptureField(makeCaptureField({ field: 'a' }));
      result.addCaptureField(makeCaptureField({ field: 'b' }));

      result.removeCaptureField(0);

      expect(result.captureFields.value).toHaveLength(1);
      expect(result.captureFields.value[0]!.field).toBe('b');
    });

    it('removeCaptureField ignores invalid index', () => {
      const { result } = withSetup(() => useWorkflowEditorState());
      result.addCaptureField(makeCaptureField());
      result.removeCaptureField(-1);

      expect(result.captureFields.value).toHaveLength(1);
    });

    it('moveCaptureField swaps correctly', () => {
      const { result } = withSetup(() => useWorkflowEditorState());
      result.addCaptureField(makeCaptureField({ field: 'a' }));
      result.addCaptureField(makeCaptureField({ field: 'b' }));

      result.moveCaptureField(0, 'down');

      expect(result.captureFields.value.map(f => f.field)).toEqual(['b', 'a']);
    });

    it('moveCaptureField is a no-op at boundaries', () => {
      const { result } = withSetup(() => useWorkflowEditorState());
      result.addCaptureField(makeCaptureField({ field: 'a' }));

      result.moveCaptureField(0, 'up');
      result.moveCaptureField(0, 'down');

      expect(result.captureFields.value).toHaveLength(1);
    });
  });

  // ────────────────────────────────────────────────────────────────────────
  // External inputs CRUD
  // ────────────────────────────────────────────────────────────────────────
  describe('external inputs', () => {
    it('addExternalInput pushes a variable', () => {
      const { result } = withSetup(() => useWorkflowEditorState());
      result.addExternalInput(makeExternalInput({ field: 'deviceId' }));

      expect(result.externalInputs.value).toHaveLength(1);
    });

    it('updateExternalInput replaces at index', () => {
      const { result } = withSetup(() => useWorkflowEditorState());
      result.addExternalInput(makeExternalInput({ field: 'deviceId' }));
      result.updateExternalInput(0, makeExternalInput({ field: 'sensorId' }));

      expect(result.externalInputs.value[0]!.field).toBe('sensorId');
    });

    it('updateExternalInput ignores invalid index', () => {
      const { result } = withSetup(() => useWorkflowEditorState());
      result.addExternalInput(makeExternalInput());
      result.updateExternalInput(5, makeExternalInput({ field: 'nope' }));

      expect(result.externalInputs.value).toHaveLength(1);
    });

    it('removeExternalInput splices at index', () => {
      const { result } = withSetup(() => useWorkflowEditorState());
      result.addExternalInput(makeExternalInput({ field: 'a' }));
      result.addExternalInput(makeExternalInput({ field: 'b' }));

      result.removeExternalInput(0);

      expect(result.externalInputs.value).toHaveLength(1);
      expect(result.externalInputs.value[0]!.field).toBe('b');
    });

    it('removeExternalInput ignores invalid index', () => {
      const { result } = withSetup(() => useWorkflowEditorState());
      result.addExternalInput(makeExternalInput());
      result.removeExternalInput(-1);

      expect(result.externalInputs.value).toHaveLength(1);
    });

    it('moveExternalInput swaps correctly', () => {
      const { result } = withSetup(() => useWorkflowEditorState());
      result.addExternalInput(makeExternalInput({ field: 'a' }));
      result.addExternalInput(makeExternalInput({ field: 'b' }));

      result.moveExternalInput(1, 'up');

      expect(result.externalInputs.value.map(e => e.field)).toEqual(['b', 'a']);
    });

    it('moveExternalInput is a no-op at boundaries', () => {
      const { result } = withSetup(() => useWorkflowEditorState());
      result.addExternalInput(makeExternalInput({ field: 'x' }));

      result.moveExternalInput(0, 'up');
      result.moveExternalInput(0, 'down');

      expect(result.externalInputs.value).toHaveLength(1);
    });
  });

  // ────────────────────────────────────────────────────────────────────────
  // variablesCount computed
  // ────────────────────────────────────────────────────────────────────────
  describe('variablesCount', () => {
    it('sums states + captureFields + externalInputs', () => {
      const { result } = withSetup(() => useWorkflowEditorState());

      result.addState(makeVariable());
      result.addState(makeVariable({ field: 'second' }));
      result.addCaptureField(makeCaptureField());
      result.addExternalInput(makeExternalInput());

      expect(result.variablesCount.value).toBe(4);
    });
  });

  // ────────────────────────────────────────────────────────────────────────
  // Installed plugins
  // ────────────────────────────────────────────────────────────────────────
  describe('installed plugins', () => {
    it('addInstalledPlugin adds unique plugin IDs', () => {
      const { result } = withSetup(() => useWorkflowEditorState());

      result.addInstalledPlugin('telegram');
      result.addInstalledPlugin('slack');
      result.addInstalledPlugin('telegram'); // duplicate

      expect(result.installedPlugins.value).toEqual(['telegram', 'slack']);
    });

    it('removeInstalledPlugin filters out the ID', () => {
      const { result } = withSetup(() => useWorkflowEditorState());

      result.addInstalledPlugin('telegram');
      result.addInstalledPlugin('slack');

      result.removeInstalledPlugin('telegram');

      expect(result.installedPlugins.value).toEqual(['slack']);
    });
  });

  // ────────────────────────────────────────────────────────────────────────
  // addNoteToNode
  // ────────────────────────────────────────────────────────────────────────
  describe('addNoteToNode', () => {
    it('creates a text note node and annotation edge', () => {
      const { result } = withSetup(() => useWorkflowEditorState());

      const target = makeNode({ id: 'target1' });
      result.nodes.value = [target];

      result.addNoteToNode('target1');

      expect(result.nodes.value).toHaveLength(2);
      const noteNode = result.nodes.value.find(n => n.type === 'core/text_note');
      expect(noteNode).toBeDefined();
      expect(noteNode!.parentNodeId).toBe('target1');

      expect(result.edges.value).toHaveLength(1);
      expect(result.edges.value[0]!.source).toBe(noteNode!.id);
      expect(result.edges.value[0]!.target).toBe('target1');
    });

    it('is a no-op for non-existent target', () => {
      const { result } = withSetup(() => useWorkflowEditorState());

      result.addNoteToNode('ghost');

      expect(result.nodes.value).toHaveLength(0);
      expect(result.edges.value).toHaveLength(0);
    });
  });

  // ────────────────────────────────────────────────────────────────────────
  // updateViewport
  // ────────────────────────────────────────────────────────────────────────
  describe('updateViewport', () => {
    it('sets viewport values', () => {
      const { result } = withSetup(() => useWorkflowEditorState());

      result.updateViewport({ x: 50, y: 75, zoom: 2 });

      expect(result.viewport.value).toEqual({ x: 50, y: 75, zoom: 2 });
    });
  });

  // ────────────────────────────────────────────────────────────────────────
  // getDefaultValue
  // ────────────────────────────────────────────────────────────────────────
  describe('getDefaultValue', () => {
    it('returns correct defaults for known types', () => {
      const { result } = withSetup(() => useWorkflowEditorState());

      expect(result.getDefaultValue('string')).toBe('');
      expect(result.getDefaultValue('number')).toBe(0);
      expect(result.getDefaultValue('boolean')).toBe(false);
      expect(result.getDefaultValue('json')).toBe('{}');
    });

    it('returns empty string for unknown type', () => {
      const { result } = withSetup(() => useWorkflowEditorState());
      expect(result.getDefaultValue('unknown_type')).toBe('');
    });
  });

  // ────────────────────────────────────────────────────────────────────────
  // getCurrentWorkflow (serialization)
  // ────────────────────────────────────────────────────────────────────────
  describe('getCurrentWorkflow', () => {
    it('serializes all state into a WorkflowDefinition', () => {
      const { result } = withSetup(() => useWorkflowEditorState());

      result.updateGeneral({ name: 'Test WF', description: 'Desc', enabled: true });
      result.addState(makeVariable({ field: 'count' }));
      result.addCaptureField(makeCaptureField({ field: 'temp' }));
      result.addExternalInput(makeExternalInput({ field: 'deviceId' }));
      result.addInstalledPlugin('telegram');

      const startNode = makeNode({ id: '__start__', type: 'core/start' });
      const delayNode = makeNode({ id: 'n1', type: 'core/delay' });
      result.nodes.value = [startNode, delayNode];

      const edge = makeEdge({ id: 'e1', source: '__start__', target: 'n1', sourceHandle: 'out' });
      result.edges.value = [edge];

      result.updateViewport({ x: 10, y: 20, zoom: 1.2 });

      const wf = result.getCurrentWorkflow.value;

      expect(wf.name).toBe('Test WF');
      expect(wf.description).toBe('Desc');
      expect(wf.enabled).toBe(true);
      expect(wf.definitionVersion).toBe(1);
      expect(wf.states).toHaveLength(1);
      expect(wf.captureFields).toHaveLength(1);
      expect(wf.externalInputs).toHaveLength(1);
      expect(wf.nodes).toHaveLength(2);
      expect(wf.edges).toHaveLength(1);
      expect(wf.installedPlugins).toEqual(['telegram']);
      expect(wf.metadata.canvasViewport).toEqual({ x: 10, y: 20, zoom: 1.2 });
    });

    it('cleans edge properties — omits default pathOffset values', () => {
      const { result } = withSetup(() => useWorkflowEditorState());

      result.edges.value = [
        makeEdge({
          id: 'e1',
          source: 'a',
          target: 'b',
          sourceHandle: 'out',
          pathOffsetX: 0, // falsy — should be omitted
          pathOffsetY: 10, // truthy — should be included
        }),
      ];

      const wf = result.getCurrentWorkflow.value;
      const cleanEdge = wf.edges[0]!;

      expect(cleanEdge.sourceHandle).toBe('out');
      expect(cleanEdge.pathOffsetX).toBeUndefined();
      expect(cleanEdge.pathOffsetY).toBe(10);
    });
  });

  // ────────────────────────────────────────────────────────────────────────
  // validateAllNodes
  // ────────────────────────────────────────────────────────────────────────
  describe('validateAllNodes', () => {
    it('returns 0 and empty errors when all nodes are valid', () => {
      const { result } = withSetup(() => useWorkflowEditorState());

      const node = makeNode({ id: 'n1', type: 'core/delay' });
      result.nodes.value = [node];
      result.edges.value = [
        makeEdge({ source: 'start', target: 'n1' }),
        makeEdge({ source: 'n1', target: 'end' }),
      ];

      mockGetNodeType.mockReturnValue({
        inputs: [{ id: 'in' }],
        outputs: [{ id: 'out' }],
      });

      const count = result.validateAllNodes();

      expect(count).toBe(0);
      expect(result.nodeValidationErrors.value).toEqual({});
    });

    it('detects missing incoming connection', () => {
      const { result } = withSetup(() => useWorkflowEditorState());

      const node = makeNode({ id: 'n1', type: 'core/delay' });
      result.nodes.value = [node];
      // outgoing only
      result.edges.value = [makeEdge({ source: 'n1', target: 'n2' })];

      mockGetNodeType.mockReturnValue({
        inputs: [{ id: 'in' }],
        outputs: [{ id: 'out' }],
      });

      const count = result.validateAllNodes();

      expect(count).toBe(1);
      expect(result.nodeValidationErrors.value['n1']).toContain('noIncomingConnection');
    });

    it('detects missing outgoing connection', () => {
      const { result } = withSetup(() => useWorkflowEditorState());

      const node = makeNode({ id: 'n1', type: 'core/delay' });
      result.nodes.value = [node];
      // incoming only
      result.edges.value = [makeEdge({ source: 'start', target: 'n1' })];

      mockGetNodeType.mockReturnValue({
        inputs: [{ id: 'in' }],
        outputs: [{ id: 'out' }],
      });

      const count = result.validateAllNodes();

      expect(count).toBe(1);
      expect(result.nodeValidationErrors.value['n1']).toContain('noOutgoingConnection');
    });

    it('calls plugin validate function and reports errors', () => {
      const { result } = withSetup(() => useWorkflowEditorState());

      const node = makeNode({ id: 'n1', type: 'core/code', config: {} });
      result.nodes.value = [node];
      result.edges.value = [
        makeEdge({ source: 'start', target: 'n1' }),
        makeEdge({ source: 'n1', target: 'end' }),
      ];

      mockGetNodeType.mockReturnValue({
        inputs: [{ id: 'in' }],
        outputs: [{ id: 'out' }],
        // eslint-disable-next-line @typescript-eslint/no-unused-vars
        validate: (_config: Record<string, unknown>) => ({
          valid: false,
          errors: ['Script is required'],
        }),
      });

      const count = result.validateAllNodes();

      expect(count).toBe(1);
      expect(result.nodeValidationErrors.value['n1']).toContain('Script is required');
    });

    it('validates required properties from declarative node definitions', () => {
      const { result } = withSetup(() => useWorkflowEditorState());

      const node = makeNode({ id: 'n1', type: 'core/http', config: {} });
      result.nodes.value = [node];
      result.edges.value = [
        makeEdge({ source: 'start', target: 'n1' }),
        makeEdge({ source: 'n1', target: 'end' }),
      ];

      mockGetNodeType.mockReturnValue({
        inputs: [{ id: 'in' }],
        outputs: [{ id: 'out' }],
        properties: [
          { name: 'url', displayName: 'URL', required: true },
          { name: 'method', displayName: 'Method', required: false },
        ],
      });

      const count = result.validateAllNodes();

      expect(count).toBe(1);
      expect(result.nodeValidationErrors.value['n1']).toContain('propRequired::URL');
    });

    it('skips hidden properties based on displayOptions.show', () => {
      const { result } = withSetup(() => useWorkflowEditorState());

      const node = makeNode({
        id: 'n1',
        type: 'core/http',
        config: { mode: 'simple' },
      });
      result.nodes.value = [node];
      result.edges.value = [
        makeEdge({ source: 'start', target: 'n1' }),
        makeEdge({ source: 'n1', target: 'end' }),
      ];

      mockGetNodeType.mockReturnValue({
        inputs: [{ id: 'in' }],
        outputs: [{ id: 'out' }],
        properties: [
          {
            name: 'advancedUrl',
            displayName: 'Advanced URL',
            required: true,
            displayOptions: { show: { mode: ['advanced'] } },
          },
        ],
      });

      const count = result.validateAllNodes();

      // Property is hidden (mode=simple != advanced), so it should NOT produce an error
      expect(count).toBe(0);
    });

    it('skips annotation types from connectivity checks', () => {
      const { result } = withSetup(() => useWorkflowEditorState());

      const note = makeNode({ id: 'note1', type: 'core/text_note' });
      result.nodes.value = [note];
      result.edges.value = []; // No edges

      mockGetNodeType.mockReturnValue({
        inputs: [{ id: 'in' }],
        outputs: [],
      });

      const count = result.validateAllNodes();

      // Annotation types are exempt from connectivity checks
      expect(count).toBe(0);
    });

    it('exempts core/start from input checks and core/end from output checks', () => {
      const { result } = withSetup(() => useWorkflowEditorState());

      const startNode = makeNode({ id: 'start', type: 'core/start' });
      const endNode = makeNode({ id: 'end', type: 'core/end' });
      result.nodes.value = [startNode, endNode];
      result.edges.value = [makeEdge({ source: 'start', target: 'end' })];

      mockGetNodeType.mockImplementation((type: string) => ({
        inputs: type === 'core/start' ? [] : [{ id: 'in' }],
        outputs: type === 'core/end' ? [] : [{ id: 'out' }],
      }));

      const count = result.validateAllNodes();

      expect(count).toBe(0);
    });

    it('detects tight cycles', () => {
      const { result } = withSetup(() => useWorkflowEditorState());

      const n1 = makeNode({ id: 'n1', type: 'core/delay' });
      const n2 = makeNode({ id: 'n2', type: 'core/delay' });
      result.nodes.value = [n1, n2];
      result.edges.value = [
        makeEdge({ source: 'n1', target: 'n2' }),
        makeEdge({ source: 'n2', target: 'n1' }),
      ];

      mockGetNodeType.mockReturnValue({
        inputs: [{ id: 'in' }],
        outputs: [{ id: 'out' }],
      });
      mockDetectTightCycles.mockReturnValue(new Set(['n1', 'n2']));

      const count = result.validateAllNodes();

      expect(count).toBe(2);
      expect(result.nodeValidationErrors.value['n1']).toContain('tightCycleDetected');
      expect(result.nodeValidationErrors.value['n2']).toContain('tightCycleDetected');
    });

    it('validates goto sender/receiver pairing', () => {
      const { result } = withSetup(() => useWorkflowEditorState());

      const sender = makeNode({
        id: 'g1',
        type: 'core/goto',
        config: { role: 'sender', pairLabel: 'ErrorHandler' },
      });
      result.nodes.value = [sender];
      result.edges.value = [makeEdge({ source: 'start', target: 'g1' })];

      mockGetNodeType.mockReturnValue({
        inputs: [{ id: 'in' }],
        outputs: [],
      });

      result.validateAllNodes();

      // Sender without matching receiver
      expect(result.nodeValidationErrors.value['g1']).toEqual(
        expect.arrayContaining([expect.stringContaining('gotoSenderNeedsReceiver')]),
      );
    });
  });

  // ────────────────────────────────────────────────────────────────────────
  // clearNodeValidationErrors
  // ────────────────────────────────────────────────────────────────────────
  describe('clearNodeValidationErrors', () => {
    it('empties the errors record', () => {
      const { result } = withSetup(() => useWorkflowEditorState());

      result.nodeValidationErrors.value = { n1: ['error'] };

      result.clearNodeValidationErrors();

      expect(result.nodeValidationErrors.value).toEqual({});
    });
  });
});
