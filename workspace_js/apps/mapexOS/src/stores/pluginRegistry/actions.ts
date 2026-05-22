import type { PluginRegistryState, WorkflowPlugin, PluginNodeType } from './types';
import type { PluginActivationContext } from '@src/components/workflow/interfaces';
import { i18nInstance } from 'src/boot/i18n';

export const actions = {
  /**
   * Register a plugin and all its node types
   *
   * @param {WorkflowPlugin} plugin - Plugin to register
   * @returns {void}
   */
  registerPlugin(plugin: WorkflowPlugin): void {
    const store = this as unknown as PluginRegistryState & typeof actions;

    if (store.plugins.has(plugin.id)) {
      console.warn(`[PluginRegistry] Plugin "${plugin.id}" already registered, skipping`);
      return;
    }

    store.plugins.set(plugin.id, plugin);

    for (const nodeType of plugin.nodeTypes) {
      if (store.nodeTypeMap.has(nodeType.type)) {
        console.warn(`[PluginRegistry] Node type "${nodeType.type}" already registered, overwriting`);
      }
      nodeType._pluginId = plugin.id;
      store.nodeTypeMap.set(nodeType.type, nodeType);
    }

    // Call plugin lifecycle hook with activation context
    if (plugin.onActivate) {
      const context: PluginActivationContext = {
        pluginId: plugin.id,
        subscriptions: [],
        registerTranslations: (locale, messages) => {
          if (!i18nInstance) {
            console.warn('[PluginRegistry] i18n instance not available, translations skipped');
            return;
          }
          i18nInstance.global.mergeLocaleMessage(locale, {
            wf: { [plugin.id]: messages },
          });
        },
      };
      plugin.onActivate(context);
    }
  },

  /**
   * Unregister a plugin and remove its node types
   *
   * @param {string} pluginId - Plugin ID to remove
   * @returns {void}
   */
  unregisterPlugin(pluginId: string): void {
    const store = this as unknown as PluginRegistryState & typeof actions;

    const plugin = store.plugins.get(pluginId);
    if (!plugin) return;

    for (const nodeType of plugin.nodeTypes) {
      store.nodeTypeMap.delete(nodeType.type);
    }

    store.plugins.delete(pluginId);
  },

  /**
   * Get a specific node type by its type string
   *
   * @param {string} type - Node type string (e.g., 'core/delay')
   * @returns {PluginNodeType | undefined} Node type definition
   */
  getNodeType(type: string): PluginNodeType | undefined {
    const store = this as unknown as PluginRegistryState & typeof actions;
    return store.nodeTypeMap.get(type);
  },

  /**
   * Get all registered node types as a Vue Flow nodeTypes map
   * Maps type string → canvas component for Vue Flow
   *
   * @returns {Record<string, any>} Vue Flow nodeTypes map
   */
  getVueFlowNodeTypes(): Record<string, any> {
    const store = this as unknown as PluginRegistryState & typeof actions;
    const result: Record<string, any> = {};

    for (const [type, nodeType] of store.nodeTypeMap.entries()) {
      result[type] = nodeType.canvasComponent;
    }

    return result;
  },

  /**
   * Clear all plugins and node types
   *
   * @returns {void}
   */
  clearAll(): void {
    const store = this as unknown as PluginRegistryState & typeof actions;
    store.plugins.clear();
    store.nodeTypeMap.clear();
  },
};
