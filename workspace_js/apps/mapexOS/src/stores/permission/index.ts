import { defineStore } from 'pinia';
import { state } from './state';
import { getters } from './getters';
import { actions } from './actions';

export const usePermissionStore = defineStore('permission', {
  state,
  getters,
  actions,
});

export type PermissionStore = ReturnType<typeof usePermissionStore>;
