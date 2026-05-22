import { defineStore } from 'pinia';
import { state } from './state';
import { getters } from './getters';
import { actions } from './actions';

export const usePluginRegistryStore = defineStore('pluginRegistry', {
  state,
  getters,
  actions,
});

export type PluginRegistryStore = ReturnType<typeof usePluginRegistryStore>;
