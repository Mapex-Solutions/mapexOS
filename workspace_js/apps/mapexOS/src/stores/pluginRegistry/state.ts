import type { PluginRegistryState } from './types';

/**
 * Plugin registry initial state factory
 *
 * @returns {PluginRegistryState} Initial state
 */
export const state = (): PluginRegistryState => ({
  plugins: new Map(),
  nodeTypeMap: new Map(),
});
