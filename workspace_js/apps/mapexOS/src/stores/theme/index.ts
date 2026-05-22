import { defineStore } from 'pinia';
import { state } from './state';
import { getters } from './getters';
import { actions } from './actions';

export const useThemeStore = defineStore('theme', {
  state,
  getters,
  actions,
});

export type ThemeStore = ReturnType<typeof useThemeStore>;
