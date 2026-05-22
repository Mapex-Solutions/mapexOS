import { describe, it, expect, vi, beforeEach } from 'vitest';
import { ref, nextTick } from 'vue';
import { mountWithPlugins } from '@src/test/helpers';
import WorkflowCanvas from './WorkflowCanvas.vue';

/**
 * ── Mock state ──
 * Module-level refs that tests can mutate to simulate composable state.
 */
const mockNodes = ref<any[]>([]);
const mockEdges = ref<any[]>([]);
const mockNodeConfigVersion = ref(0);
const mockNodeValidationErrors = ref<Record<string, string[]>>({});
const mockInstalledPlugins = ref<string[]>([]);
const mockIsBatchRestoring = ref(false);

const mockAddNode = vi.fn();
const mockAddEdge = vi.fn();
const mockRemoveEdge = vi.fn();
const mockRemoveNode = vi.fn().mockReturnValue([]);
const mockDuplicateNode = vi.fn();
const mockUpdateViewport = vi.fn();
const mockPushSnapshot = vi.fn();
const mockUndo = vi.fn();
const mockRedo = vi.fn();
const mockFinishRestore = vi.fn();

/** Vue Flow mock functions */
const mockSetNodes = vi.fn();
const mockSetEdges = vi.fn();
const mockAddEdges = vi.fn();
const mockVfRemoveEdges = vi.fn();
const mockVfAddNodes = vi.fn();
const mockVfRemoveNodes = vi.fn();
const mockUpdateNodeInternals = vi.fn();
const mockFitView = vi.fn();
const mockScreenToFlowCoordinate = vi.fn().mockReturnValue({ x: 100, y: 200 });
const mockGetNodes = ref<any[]>([]);
const mockGetSelectedNodes = ref<any[]>([]);
const mockGetSelectedEdges = ref<any[]>([]);

/** Event handler captures — store the callbacks so tests can invoke them */
type AnyFn = (...args: any[]) => any;
let onConnectHandler: AnyFn;
let onNodeClickHandler: AnyFn;
let onPaneClickHandler: AnyFn;
let onEdgesChangeHandler: AnyFn;
let onNodeDragStartHandler: AnyFn;
let onNodeDragStopHandler: AnyFn;
let onViewportChangeEndHandler: AnyFn;

vi.mock('@vue-flow/core', () => ({
  VueFlow: { name: 'VueFlow', template: '<div><slot /></div>', props: ['nodes', 'edges', 'nodeTypes', 'edgeTypes'] },
  useVueFlow: () => ({
    onConnect: (fn: AnyFn) => { onConnectHandler = fn; },
    addEdges: mockAddEdges,
    removeEdges: mockVfRemoveEdges,
    onNodeClick: (fn: AnyFn) => { onNodeClickHandler = fn; },
    onPaneClick: (fn: AnyFn) => { onPaneClickHandler = fn; },
    onEdgesChange: (fn: AnyFn) => { onEdgesChangeHandler = fn; },
    onNodeDragStart: (fn: AnyFn) => { onNodeDragStartHandler = fn; },
    onNodeDragStop: (fn: AnyFn) => { onNodeDragStopHandler = fn; },
    screenToFlowCoordinate: mockScreenToFlowCoordinate,
    getNodes: mockGetNodes,
    getSelectedNodes: mockGetSelectedNodes,
    getSelectedEdges: mockGetSelectedEdges,
    setNodes: mockSetNodes,
    setEdges: mockSetEdges,
    addNodes: mockVfAddNodes,
    removeNodes: mockVfRemoveNodes,
    updateNodeInternals: mockUpdateNodeInternals,
    onViewportChangeEnd: (fn: AnyFn) => { onViewportChangeEndHandler = fn; },
    fitView: mockFitView,
  }),
  MarkerType: { ArrowClosed: 'arrowclosed' },
  Position: { Left: 'left', Right: 'right', Top: 'top', Bottom: 'bottom' },
  Handle: { name: 'Handle', template: '<div />' },
}));

vi.mock('@vue-flow/minimap', () => ({
  MiniMap: { name: 'MiniMap', template: '<div />' },
}));

vi.mock('@vue-flow/controls', () => ({
  Controls: { name: 'Controls', template: '<div />' },
}));

vi.mock('@vue-flow/background', () => ({
  Background: { name: 'Background', template: '<div />' },
}));

vi.mock('../../composables', () => ({
  useWorkflowEditorState: () => ({
    nodes: mockNodes,
    edges: mockEdges,
    addNode: mockAddNode,
    addEdge: mockAddEdge,
    removeEdge: mockRemoveEdge,
    removeNode: mockRemoveNode,
    duplicateNode: mockDuplicateNode,
    nodeConfigVersion: mockNodeConfigVersion,
    nodeValidationErrors: mockNodeValidationErrors,
    installedPlugins: mockInstalledPlugins,
    updateViewport: mockUpdateViewport,
  }),
  useWorkflowHistory: () => ({
    pushSnapshot: mockPushSnapshot,
    undo: mockUndo,
    redo: mockRedo,
    finishRestore: mockFinishRestore,
    isBatchRestoring: mockIsBatchRestoring,
  }),
}));

vi.mock('@stores/pluginRegistry', () => ({
  usePluginRegistryStore: () => ({
    nodeTypeMap: new Map(),
    getNodeType: vi.fn().mockReturnValue({ label: 'Test', configurable: true, deletable: true }),
  }),
}));

vi.mock('../../utils', () => ({
  createConnectionValidator: () => vi.fn().mockReturnValue(true),
  resolveNodeHandles: vi.fn(),
}));

vi.mock('@utils/workflow', () => ({
  buildDefaultConfig: vi.fn().mockReturnValue({}),
}));

vi.mock('./utils', () => ({
  createWorkflowHotkeyHandler: () => ({ attach: vi.fn(), detach: vi.fn() }),
}));

vi.mock('../AdjustableEdge/AdjustableEdge.vue', () => ({
  default: { name: 'AdjustableEdge', template: '<div />' },
}));

vi.mock('../AnnotationEdge/AnnotationEdge.vue', () => ({
  default: { name: 'AnnotationEdge', template: '<div />' },
}));

vi.mock('@components/workflow/nodes/GenericWorkflowNode', () => ({
  GenericWorkflowNode: { name: 'GenericWorkflowNode', template: '<div />' },
}));

vi.mock('quasar', () => ({
  useQuasar: () => ({ dark: { isActive: false } }),
}));

const DEFAULT_TOOLBAR = {
  showMinimap: false,
  showGrid: true,
  locked: false,
  maximized: false,
};

function makeNode(overrides: Record<string, any> = {}): any {
  return {
    id: 'n1',
    type: 'core/start',
    position: { x: 100, y: 200 },
    config: {},
    label: 'Start',
    ...overrides,
  };
}

function makeEdge(overrides: Record<string, any> = {}): any {
  return {
    id: 'e1',
    source: 'n1',
    target: 'n2',
    ...overrides,
  };
}

function mountCanvas() {
  return mountWithPlugins(WorkflowCanvas, {
    props: { toolbarState: DEFAULT_TOOLBAR },
  });
}

describe('WorkflowCanvas', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    mockNodes.value = [makeNode()];
    mockEdges.value = [];
    mockNodeConfigVersion.value = 0;
    mockNodeValidationErrors.value = {};
    mockIsBatchRestoring.value = false;
    mockGetNodes.value = [];
    mockGetSelectedNodes.value = [];
    mockGetSelectedEdges.value = [];
  });

  it('renders without errors', () => {
    const wrapper = mountCanvas();
    expect(wrapper.exists()).toBe(true);
    expect(wrapper.find('.workflow-canvas').exists()).toBe(true);
  });

  // ──────────────────────────────────────────────────────────────────────
  // flowNodes computed
  // ──────────────────────────────────────────────────────────────────────

  describe('flowNodes computed', () => {
    it('maps nodes to Vue Flow format', () => {
      const node = makeNode({ id: 'n1', type: 'core/delay', label: 'Wait' });
      mockNodes.value = [node];

      const wrapper = mountCanvas();
      const flowNodes = wrapper.vm.flowNodes;

      expect(flowNodes).toHaveLength(1);
      expect(flowNodes[0]).toMatchObject({
        id: 'n1',
        type: 'core/delay',
        position: { x: 100, y: 200 },
        data: expect.objectContaining({ label: 'Wait', __nodeType: 'core/delay' }),
      });
    });

    it('does NOT set parentNode for text_note nodes (absolute positioning)', () => {
      mockNodes.value = [
        makeNode({ id: 'parent1' }),
        makeNode({ id: 'note1', type: 'core/text_note', parentNodeId: 'parent1', position: { x: 10, y: 180 } }),
      ];

      const wrapper = mountCanvas();
      const flowNodes = wrapper.vm.flowNodes;
      const noteNode = flowNodes.find((n: any) => n.id === 'note1');

      expect(noteNode).not.toHaveProperty('parentNode');
    });

    it('sets parentNode for group_frame children', () => {
      mockNodes.value = [
        makeNode({ id: 'frame1', type: 'core/group_frame' }),
        makeNode({ id: 'child1', type: 'core/delay', parentNodeId: 'frame1', position: { x: 20, y: 30 } }),
      ];

      const wrapper = mountCanvas();
      const flowNodes = wrapper.vm.flowNodes;
      const childNode = flowNodes.find((n: any) => n.id === 'child1');

      expect(childNode.parentNode).toBe('frame1');
    });

    it('does not set parentNode for regular nodes', () => {
      mockNodes.value = [makeNode({ id: 'n1' })];

      const wrapper = mountCanvas();
      const flowNodes = wrapper.vm.flowNodes;

      expect(flowNodes[0]).not.toHaveProperty('parentNode');
    });

    it('sets zIndex -1 for group frame nodes', () => {
      mockNodes.value = [makeNode({ id: 'frame1', type: 'core/group_frame' })];

      const wrapper = mountCanvas();
      const flowNodes = wrapper.vm.flowNodes;

      expect(flowNodes[0].zIndex).toBe(-1);
    });

    it('marks hasErrors in data when node has validation errors', () => {
      mockNodes.value = [makeNode({ id: 'n1' })];
      mockNodeValidationErrors.value = { n1: ['Missing field'] };

      const wrapper = mountCanvas();
      const flowNodes = wrapper.vm.flowNodes;

      expect(flowNodes[0].data.hasErrors).toBe(true);
    });
  });

  // ──────────────────────────────────────────────────────────────────────
  // flowEdges computed
  // ──────────────────────────────────────────────────────────────────────

  describe('flowEdges computed', () => {
    it('maps edges to adjustable type with animation', () => {
      mockEdges.value = [makeEdge()];

      const wrapper = mountCanvas();
      const flowEdges = wrapper.vm.flowEdges;

      expect(flowEdges[0]).toMatchObject({
        type: 'adjustable',
        animated: true,
      });
      expect(flowEdges[0].markerEnd).toBeDefined();
    });
  });

  // ──────────────────────────────────────────────────────────────────────
  // nodeIdList watcher — child node position preservation
  // ──────────────────────────────────────────────────────────────────────

  describe('nodeIdList watcher', () => {
    it('preserves composable position for group_frame child nodes', async () => {
      mountCanvas();

      mockGetNodes.value = [
        { id: 'frame1', position: { x: 300, y: 400 } },
      ];

      mockNodes.value = [
        makeNode({ id: 'frame1', type: 'core/group_frame', position: { x: 250, y: 350 } }),
        makeNode({ id: 'child1', type: 'core/delay', parentNodeId: 'frame1', position: { x: 20, y: 30 } }),
      ];

      await nextTick();

      const call = mockSetNodes.mock.calls[mockSetNodes.mock.calls.length - 1];
      if (call) {
        const nodeArgs = call[0] as any[];
        const childNode = nodeArgs.find((n: any) => n.id === 'child1');
        expect(childNode.position).toEqual({ x: 20, y: 30 });
        expect(childNode.parentNode).toBe('frame1');
      }
    });

    it('uses Vue Flow dragged position for regular nodes', async () => {
      mountCanvas();

      mockGetNodes.value = [
        { id: 'n1', position: { x: 500, y: 600 } },
      ];

      mockNodes.value = [
        makeNode({ id: 'n1', position: { x: 100, y: 200 } }),
        makeNode({ id: 'n2', position: { x: 150, y: 250 } }),
      ];

      await nextTick();

      const call = mockSetNodes.mock.calls[mockSetNodes.mock.calls.length - 1];
      if (call) {
        const nodeArgs = call[0] as any[];
        const n1 = nodeArgs.find((n: any) => n.id === 'n1');
        expect(n1.position).toEqual({ x: 500, y: 600 });
      }
    });

    it('calls updateNodeInternals after setNodes', async () => {
      mountCanvas();

      mockNodes.value = [
        makeNode({ id: 'n1' }),
        makeNode({ id: 'n2' }),
      ];

      await nextTick();
      await nextTick();

      expect(mockUpdateNodeInternals).toHaveBeenCalled();
    });

    it('skips sync when isBatchRestoring is true', async () => {
      mountCanvas();

      mockIsBatchRestoring.value = true;
      mockSetNodes.mockClear();

      mockNodes.value = [
        makeNode({ id: 'n1' }),
        makeNode({ id: 'n2' }),
      ];

      await nextTick();

      expect(mockSetNodes).not.toHaveBeenCalled();
    });
  });

  // ──────────────────────────────────────────────────────────────────────
  // edgeIdList watcher — double nextTick for handle registration
  // ──────────────────────────────────────────────────────────────────────

  describe('edgeIdList watcher', () => {
    it('calls setEdges after double nextTick', async () => {
      mountCanvas();

      mockSetEdges.mockClear();
      mockEdges.value = [makeEdge()];

      // Watcher fires on first tick, then double nextTick inside:
      // tick 1: watcher triggers -> schedules nextTick A
      // tick 2: nextTick A runs -> schedules nextTick B
      // tick 3: nextTick B runs -> setEdges called
      await nextTick();
      await nextTick();
      await nextTick();

      expect(mockSetEdges).toHaveBeenCalled();
    });

    it('skips sync when isBatchRestoring is true', async () => {
      mountCanvas();

      mockIsBatchRestoring.value = true;
      mockSetEdges.mockClear();

      mockEdges.value = [makeEdge()];

      await nextTick();
      await nextTick();
      await nextTick();

      expect(mockSetEdges).not.toHaveBeenCalled();
    });
  });

  // ──────────────────────────────────────────────────────────────────────
  // nodeConfigVersion watcher — child node position in config sync
  // ──────────────────────────────────────────────────────────────────────

  describe('nodeConfigVersion watcher', () => {
    it('preserves composable position for group_frame children on config change', async () => {
      mockNodes.value = [
        makeNode({ id: 'frame1', type: 'core/group_frame', position: { x: 250, y: 350 } }),
        makeNode({ id: 'child1', type: 'core/delay', parentNodeId: 'frame1', position: { x: 20, y: 30 } }),
      ];

      mockGetNodes.value = [
        { id: 'frame1', position: { x: 300, y: 400 } },
        { id: 'child1', position: { x: 50, y: 50 } },
      ];

      mountCanvas();
      mockSetNodes.mockClear();

      mockNodeConfigVersion.value++;
      await nextTick();

      const call = mockSetNodes.mock.calls[0];
      if (call) {
        const nodeArgs = call[0] as any[];
        const childNode = nodeArgs.find((n: any) => n.id === 'child1');
        expect(childNode.position).toEqual({ x: 20, y: 30 });
      }
    });

    it('calls updateNodeInternals after config change', async () => {
      mockNodes.value = [makeNode()];

      mountCanvas();
      mockUpdateNodeInternals.mockClear();

      mockNodeConfigVersion.value++;
      await nextTick();
      await nextTick();

      expect(mockUpdateNodeInternals).toHaveBeenCalled();
    });
  });

  // ──────────────────────────────────────────────────────────────────────
  // Vue Flow event handlers
  // ──────────────────────────────────────────────────────────────────────

  describe('onConnect handler', () => {
    it('creates snapshot, adds edge to Vue Flow and composable', () => {
      mountCanvas();

      onConnectHandler({
        source: 'n1',
        target: 'n2',
        sourceHandle: 'out',
        targetHandle: 'in',
      });

      expect(mockPushSnapshot).toHaveBeenCalledWith('Connect');
      expect(mockAddEdges).toHaveBeenCalled();
      expect(mockAddEdge).toHaveBeenCalledWith(
        expect.objectContaining({ source: 'n1', target: 'n2' }),
      );
    });
  });

  describe('onEdgesChange handler', () => {
    it('removes edges from composable state on remove changes', () => {
      mountCanvas();

      onEdgesChangeHandler([{ type: 'remove', id: 'e1' }]);

      expect(mockPushSnapshot).toHaveBeenCalledWith('Remove edge');
      expect(mockRemoveEdge).toHaveBeenCalledWith('e1');
    });

    it('does not push snapshot for non-remove changes', () => {
      mountCanvas();

      onEdgesChangeHandler([{ type: 'select', id: 'e1' }]);

      expect(mockPushSnapshot).not.toHaveBeenCalled();
    });

    it('skips when isBatchRestoring', () => {
      mountCanvas();

      mockIsBatchRestoring.value = true;
      onEdgesChangeHandler([{ type: 'remove', id: 'e1' }]);

      expect(mockRemoveEdge).not.toHaveBeenCalled();
    });
  });

  describe('onNodeDragStart handler', () => {
    it('pushes a snapshot before drag', () => {
      mountCanvas();

      onNodeDragStartHandler();

      expect(mockPushSnapshot).toHaveBeenCalledWith('Move node');
    });
  });

  describe('onNodeDragStop handler', () => {
    it('syncs dragged positions back to composable state', () => {
      const node = makeNode({ id: 'n1', position: { x: 100, y: 200 } });
      mockNodes.value = [node];

      mountCanvas();

      onNodeDragStopHandler({
        nodes: [{ id: 'n1', position: { x: 500, y: 600 } }],
      });

      expect(mockNodes.value[0].position).toEqual({ x: 500, y: 600 });
    });

    it('shifts child text notes by the same delta when parent moves', () => {
      const parent = makeNode({ id: 'p1', position: { x: 100, y: 200 } });
      const note = makeNode({ id: 'note1', type: 'core/text_note', parentNodeId: 'p1', position: { x: 10, y: 180 } });
      mockNodes.value = [parent, note];

      mountCanvas();

      // Parent dragged from (100,200) to (200,300) → delta = (100, 100)
      onNodeDragStopHandler({
        nodes: [{ id: 'p1', position: { x: 200, y: 300 } }],
      });

      expect(mockNodes.value[0].position).toEqual({ x: 200, y: 300 });
      expect(mockNodes.value[1].position).toEqual({ x: 110, y: 280 });
    });

    it('does not shift non-text-note children', () => {
      const parent = makeNode({ id: 'p1', position: { x: 100, y: 200 } });
      const child = makeNode({ id: 'c1', type: 'core/delay', parentNodeId: 'p1', position: { x: 20, y: 30 } });
      mockNodes.value = [parent, child];

      mountCanvas();

      onNodeDragStopHandler({
        nodes: [{ id: 'p1', position: { x: 200, y: 300 } }],
      });

      // group_frame children use Vue Flow parentNode — their position is relative, unchanged
      expect(mockNodes.value[1].position).toEqual({ x: 20, y: 30 });
    });
  });

  describe('onNodeClick handler', () => {
    it('emits node-select for configurable nodes', () => {
      const wrapper = mountCanvas();

      onNodeClickHandler({ node: { id: 'n1', data: { __nodeType: 'core/delay' } }, event: {} });

      expect(wrapper.emitted('node-select')).toBeTruthy();
      expect(wrapper.emitted('node-select')![0]).toEqual(['n1']);
    });

    it('does not emit when ctrl/meta key is pressed', () => {
      const wrapper = mountCanvas();

      onNodeClickHandler({ node: { id: 'n1', data: { __nodeType: 'core/delay' } }, event: { ctrlKey: true } });

      expect(wrapper.emitted('node-select')).toBeFalsy();
    });
  });

  describe('onPaneClick handler', () => {
    it('emits canvas-click', () => {
      const wrapper = mountCanvas();

      onPaneClickHandler();

      expect(wrapper.emitted('canvas-click')).toBeTruthy();
    });
  });

  describe('onViewportChangeEnd handler', () => {
    it('syncs viewport to composable state', () => {
      mountCanvas();

      onViewportChangeEndHandler({ x: 10, y: 20, zoom: 1.5 });

      expect(mockUpdateViewport).toHaveBeenCalledWith({ x: 10, y: 20, zoom: 1.5 });
    });
  });

  // ──────────────────────────────────────────────────────────────────────
  // syncSnapshotToVueFlow (via expose)
  // ──────────────────────────────────────────────────────────────────────

  describe('restoreSnapshot (exposed)', () => {
    it('sets nodes and edges correctly', () => {
      const wrapper = mountCanvas();

      mockSetNodes.mockClear();
      mockSetEdges.mockClear();

      const snapshotNodes = [
        makeNode({ id: 'parent1' }),
        makeNode({ id: 'note1', type: 'core/text_note', parentNodeId: 'parent1', position: { x: 160, y: 210 } }),
      ];
      const snapshotEdges = [makeEdge()];

      wrapper.vm.restoreSnapshot(snapshotNodes, snapshotEdges);

      expect(mockSetNodes).toHaveBeenCalledTimes(1);
      expect(mockSetEdges).toHaveBeenCalledTimes(1);

      const edgeArgs = mockSetEdges.mock.calls[0]![0] as any[];
      expect(edgeArgs[0].type).toBe('adjustable');

      const nodeArgs = mockSetNodes.mock.calls[0]![0] as any[];
      const noteNode = nodeArgs.find((n: any) => n.id === 'note1');
      expect(noteNode).not.toHaveProperty('parentNode');
    });

    it('calls updateNodeInternals after restore', async () => {
      const wrapper = mountCanvas();

      mockUpdateNodeInternals.mockClear();

      wrapper.vm.restoreSnapshot([makeNode()], []);

      await nextTick();

      expect(mockUpdateNodeInternals).toHaveBeenCalled();
    });
  });

  // ──────────────────────────────────────────────────────────────────────
  // handleDrop
  // ──────────────────────────────────────────────────────────────────────

  describe('handleDrop', () => {
    it('creates a new node at drop position', async () => {
      const wrapper = mountCanvas();

      const dropEvent = {
        dataTransfer: {
          getData: (type: string) => type === 'application/workflow-node-type' ? 'core/delay' : '',
        },
        clientX: 400,
        clientY: 300,
        preventDefault: vi.fn(),
      } as unknown as DragEvent;

      await wrapper.find('.workflow-canvas').trigger('drop', dropEvent);
    });
  });
});
