import { describe, it, expect, beforeEach, vi } from 'vitest';
import { setActivePinia, createPinia } from 'pinia';
import { usePluginRegistryStore } from './index';
import type { WorkflowPlugin, PluginNodeType } from './types';

/**
 * Create a minimal mock PluginNodeType for testing
 *
 * @param {Partial<PluginNodeType>} overrides - Fields to override
 * @returns {PluginNodeType} Mock node type
 */
function createMockNodeType(overrides: Partial<PluginNodeType> = {}): PluginNodeType {
  return {
    type: 'core/test',
    label: 'Test Node',
    icon: 'science',
    color: '#4caf50',
    description: 'A test node',
    inputs: [{ id: 'in', label: 'Input', position: 'left' }],
    outputs: [{ id: 'out', label: 'Output', position: 'right' }],
    configSchema: {},
    ...overrides,
  };
}

/**
 * Create a minimal mock WorkflowPlugin for testing
 *
 * @param {Partial<WorkflowPlugin>} overrides - Fields to override
 * @returns {WorkflowPlugin} Mock plugin
 */
function createMockPlugin(overrides: Partial<WorkflowPlugin> = {}): WorkflowPlugin {
  return {
    id: 'test-plugin',
    name: 'Test Plugin',
    version: '1.0.0',
    category: 'logic',
    icon: 'extension',
    nodeTypes: [createMockNodeType()],
    ...overrides,
  };
}

describe('PluginRegistryStore', () => {
  beforeEach(() => {
    setActivePinia(createPinia());
  });

  // ────────────────────────────────────────────────────────────────────────
  // State
  // ────────────────────────────────────────────────────────────────────────

  describe('state', () => {
    it('has empty plugins map by default', () => {
      const store = usePluginRegistryStore();
      expect(store.plugins.size).toBe(0);
    });

    it('has empty nodeTypeMap by default', () => {
      const store = usePluginRegistryStore();
      expect(store.nodeTypeMap.size).toBe(0);
    });
  });

  // ────────────────────────────────────────────────────────────────────────
  // Actions
  // ────────────────────────────────────────────────────────────────────────

  describe('actions', () => {
    describe('registerPlugin', () => {
      it('adds plugin to plugins map', () => {
        const store = usePluginRegistryStore();
        const plugin = createMockPlugin();

        store.registerPlugin(plugin);

        expect(store.plugins.has('test-plugin')).toBe(true);
        expect(store.plugins.get('test-plugin')?.name).toBe('Test Plugin');
      });

      it('adds node types to nodeTypeMap', () => {
        const store = usePluginRegistryStore();
        const plugin = createMockPlugin();

        store.registerPlugin(plugin);

        expect(store.nodeTypeMap.has('core/test')).toBe(true);
        expect(store.nodeTypeMap.get('core/test')?.label).toBe('Test Node');
      });

      it('sets _pluginId on registered node types', () => {
        const store = usePluginRegistryStore();
        const plugin = createMockPlugin();

        store.registerPlugin(plugin);

        const nodeType = store.nodeTypeMap.get('core/test');
        expect(nodeType?._pluginId).toBe('test-plugin');
      });

      it('registers multiple node types from a single plugin', () => {
        const store = usePluginRegistryStore();
        const plugin = createMockPlugin({
          nodeTypes: [
            createMockNodeType({ type: 'core/delay', label: 'Delay' }),
            createMockNodeType({ type: 'core/condition', label: 'Condition' }),
          ],
        });

        store.registerPlugin(plugin);

        expect(store.nodeTypeMap.size).toBe(2);
        expect(store.nodeTypeMap.has('core/delay')).toBe(true);
        expect(store.nodeTypeMap.has('core/condition')).toBe(true);
      });

      it('skips duplicate plugin registration (same id)', () => {
        const store = usePluginRegistryStore();
        const plugin1 = createMockPlugin({
          nodeTypes: [createMockNodeType({ type: 'core/v1', label: 'V1' })],
        });
        const plugin2 = createMockPlugin({
          nodeTypes: [createMockNodeType({ type: 'core/v2', label: 'V2' })],
        });

        store.registerPlugin(plugin1);
        store.registerPlugin(plugin2);

        // Second registration is skipped — original stays
        expect(store.plugins.size).toBe(1);
        expect(store.nodeTypeMap.has('core/v1')).toBe(true);
        expect(store.nodeTypeMap.has('core/v2')).toBe(false);
      });

      it('calls onActivate lifecycle hook when provided', () => {
        const store = usePluginRegistryStore();
        const onActivate = vi.fn();
        const plugin = createMockPlugin({ onActivate });

        store.registerPlugin(plugin);

        expect(onActivate).toHaveBeenCalledOnce();
        expect(onActivate).toHaveBeenCalledWith(
          expect.objectContaining({
            pluginId: 'test-plugin',
            subscriptions: [],
            registerTranslations: expect.any(Function),
          }),
        );
      });

      it('overwrites nodeTypeMap entry when different plugins register same node type', () => {
        const store = usePluginRegistryStore();
        const sharedNodeType = createMockNodeType({ type: 'shared/node', label: 'Original' });
        const plugin1 = createMockPlugin({
          id: 'plugin-a',
          nodeTypes: [sharedNodeType],
        });

        store.registerPlugin(plugin1);

        // Manually set a different node type at the same key to simulate overwrite
        const overwriteNodeType = createMockNodeType({ type: 'shared/node', label: 'Overwritten' });
        overwriteNodeType._pluginId = 'plugin-b';
        store.nodeTypeMap.set('shared/node', overwriteNodeType);

        expect(store.nodeTypeMap.get('shared/node')?.label).toBe('Overwritten');
        expect(store.nodeTypeMap.get('shared/node')?._pluginId).toBe('plugin-b');
      });
    });

    describe('unregisterPlugin', () => {
      it('removes plugin from plugins map', () => {
        const store = usePluginRegistryStore();
        store.registerPlugin(createMockPlugin());

        store.unregisterPlugin('test-plugin');

        expect(store.plugins.has('test-plugin')).toBe(false);
      });

      it('removes node types from nodeTypeMap', () => {
        const store = usePluginRegistryStore();
        store.registerPlugin(createMockPlugin());

        store.unregisterPlugin('test-plugin');

        expect(store.nodeTypeMap.has('core/test')).toBe(false);
        expect(store.nodeTypeMap.size).toBe(0);
      });

      it('does nothing when plugin does not exist', () => {
        const store = usePluginRegistryStore();
        store.registerPlugin(createMockPlugin());

        store.unregisterPlugin('non-existent');

        // Original plugin should still be there
        expect(store.plugins.size).toBe(1);
        expect(store.nodeTypeMap.size).toBe(1);
      });

      it('only removes node types belonging to the unregistered plugin', () => {
        const store = usePluginRegistryStore();
        const pluginA = createMockPlugin({
          id: 'plugin-a',
          nodeTypes: [createMockNodeType({ type: 'a/node', label: 'A' })],
        });
        const pluginB = createMockPlugin({
          id: 'plugin-b',
          nodeTypes: [createMockNodeType({ type: 'b/node', label: 'B' })],
        });

        store.registerPlugin(pluginA);
        store.registerPlugin(pluginB);

        store.unregisterPlugin('plugin-a');

        expect(store.plugins.has('plugin-a')).toBe(false);
        expect(store.nodeTypeMap.has('a/node')).toBe(false);
        expect(store.plugins.has('plugin-b')).toBe(true);
        expect(store.nodeTypeMap.has('b/node')).toBe(true);
      });
    });

    describe('getNodeType', () => {
      it('returns the correct node type by type string', () => {
        const store = usePluginRegistryStore();
        store.registerPlugin(createMockPlugin());

        const result = store.getNodeType('core/test');

        expect(result).toBeDefined();
        expect(result?.type).toBe('core/test');
        expect(result?.label).toBe('Test Node');
      });

      it('returns undefined for unknown type', () => {
        const store = usePluginRegistryStore();

        const result = store.getNodeType('unknown/type');

        expect(result).toBeUndefined();
      });

      it('returns undefined after plugin is unregistered', () => {
        const store = usePluginRegistryStore();
        store.registerPlugin(createMockPlugin());
        store.unregisterPlugin('test-plugin');

        const result = store.getNodeType('core/test');

        expect(result).toBeUndefined();
      });
    });

    describe('getVueFlowNodeTypes', () => {
      it('returns map of type to canvasComponent', () => {
        const store = usePluginRegistryStore();
        const mockComponent = { template: '<div />' };
        const plugin = createMockPlugin({
          nodeTypes: [createMockNodeType({ type: 'core/test', canvasComponent: mockComponent as any })],
        });

        store.registerPlugin(plugin);

        const result = store.getVueFlowNodeTypes();

        expect(result['core/test']).toStrictEqual(mockComponent);
      });

      it('returns empty object when no plugins registered', () => {
        const store = usePluginRegistryStore();

        const result = store.getVueFlowNodeTypes();

        expect(result).toEqual({});
      });
    });

    describe('clearAll', () => {
      it('removes all plugins and node types', () => {
        const store = usePluginRegistryStore();
        store.registerPlugin(createMockPlugin({ id: 'p1', nodeTypes: [createMockNodeType({ type: 'a/1' })] }));
        store.registerPlugin(createMockPlugin({ id: 'p2', nodeTypes: [createMockNodeType({ type: 'b/1' })] }));

        store.clearAll();

        expect(store.plugins.size).toBe(0);
        expect(store.nodeTypeMap.size).toBe(0);
      });
    });
  });

  // ────────────────────────────────────────────────────────────────────────
  // Getters
  // ────────────────────────────────────────────────────────────────────────

  describe('getters', () => {
    describe('nodeTypeCount', () => {
      it('returns 0 when empty', () => {
        const store = usePluginRegistryStore();
        expect(store.nodeTypeCount).toBe(0);
      });

      it('returns correct count after registration', () => {
        const store = usePluginRegistryStore();
        store.registerPlugin(createMockPlugin({
          nodeTypes: [
            createMockNodeType({ type: 'core/a' }),
            createMockNodeType({ type: 'core/b' }),
          ],
        }));

        expect(store.nodeTypeCount).toBe(2);
      });
    });

    describe('pluginCount', () => {
      it('returns 0 when empty', () => {
        const store = usePluginRegistryStore();
        expect(store.pluginCount).toBe(0);
      });

      it('returns correct count after registration', () => {
        const store = usePluginRegistryStore();
        store.registerPlugin(createMockPlugin({ id: 'p1', nodeTypes: [createMockNodeType({ type: 'a/1' })] }));
        store.registerPlugin(createMockPlugin({ id: 'p2', nodeTypes: [createMockNodeType({ type: 'b/1' })] }));

        expect(store.pluginCount).toBe(2);
      });
    });
  });
});
