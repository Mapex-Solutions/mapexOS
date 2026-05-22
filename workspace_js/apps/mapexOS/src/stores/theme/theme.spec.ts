import { describe, it, expect, beforeEach, vi } from 'vitest';
import { setActivePinia, createPinia } from 'pinia';
import { useThemeStore } from './index';

describe('ThemeStore', () => {
  beforeEach(() => {
    setActivePinia(createPinia());
    localStorage.clear();
  });

  describe('state', () => {
    it('defaults to system mode', () => {
      const store = useThemeStore();
      expect(store.mode).toBe('system');
    });

    it('defaults resolvedTheme to light', () => {
      const store = useThemeStore();
      expect(store.resolvedTheme).toBe('light');
    });

    it('reads persisted mode from localStorage', () => {
      localStorage.setItem('mapex-theme', 'dark');
      const store = useThemeStore();
      expect(store.mode).toBe('dark');
    });
  });

  describe('getters', () => {
    it('isDark returns false when resolvedTheme is light', () => {
      const store = useThemeStore();
      store.resolvedTheme = 'light';
      expect(store.isDark).toBe(false);
    });

    it('isDark returns true when resolvedTheme is dark', () => {
      const store = useThemeStore();
      store.resolvedTheme = 'dark';
      expect(store.isDark).toBe(true);
    });

    it('currentMode returns the mode', () => {
      const store = useThemeStore();
      expect(store.currentMode).toBe('system');
    });
  });

  describe('actions', () => {
    it('setTheme updates mode and persists to localStorage', () => {
      const store = useThemeStore();
      store.setTheme('dark');
      expect(store.mode).toBe('dark');
      expect(store.resolvedTheme).toBe('dark');
      expect(localStorage.getItem('mapex-theme')).toBe('dark');
    });

    it('setTheme with light sets resolvedTheme to light', () => {
      const store = useThemeStore();
      store.setTheme('light');
      expect(store.resolvedTheme).toBe('light');
    });

    it('setTheme with system resolves from matchMedia', () => {
      // Mock matchMedia to prefer dark
      Object.defineProperty(window, 'matchMedia', {
        writable: true,
        value: vi.fn().mockImplementation((query: string) => ({
          matches: query === '(prefers-color-scheme: dark)',
          media: query,
        })),
      });

      const store = useThemeStore();
      store.setTheme('system');
      expect(store.mode).toBe('system');
      expect(store.resolvedTheme).toBe('dark');
    });

    it('toggleTheme switches from light to dark', () => {
      const store = useThemeStore();
      store.setTheme('light');
      store.toggleTheme();
      expect(store.resolvedTheme).toBe('dark');
    });

    it('toggleTheme switches from dark to light', () => {
      const store = useThemeStore();
      store.setTheme('dark');
      store.toggleTheme();
      expect(store.resolvedTheme).toBe('light');
    });

    it('initTheme reads from localStorage', () => {
      localStorage.setItem('mapex-theme', 'dark');
      const store = useThemeStore();
      store.initTheme();
      expect(store.mode).toBe('dark');
      expect(store.resolvedTheme).toBe('dark');
    });

    it('initTheme defaults to system when no localStorage', () => {
      Object.defineProperty(window, 'matchMedia', {
        writable: true,
        value: vi.fn().mockImplementation((query: string) => ({
          matches: false,
          media: query,
        })),
      });

      const store = useThemeStore();
      store.initTheme();
      expect(store.mode).toBe('system');
      expect(store.resolvedTheme).toBe('light');
    });
  });
});
