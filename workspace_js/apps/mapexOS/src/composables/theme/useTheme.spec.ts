import { describe, it, expect, beforeEach, vi } from 'vitest';
import { setActivePinia, createPinia } from 'pinia';
import { useThemeStore } from '@stores/theme';

/** Mock Quasar's useQuasar */
const mockDarkSet = vi.fn();
vi.mock('quasar', () => ({
  useQuasar: () => ({
    dark: {
      set: mockDarkSet,
      isActive: false,
    },
  }),
}));

/** Import after mocks are set up */
import { useTheme } from './useTheme';

describe('useTheme', () => {
  beforeEach(() => {
    setActivePinia(createPinia());
    localStorage.clear();
    mockDarkSet.mockClear();
  });

  it('returns isDark, currentMode, setTheme, and toggleTheme', () => {
    const result = useTheme();

    expect(result).toHaveProperty('isDark');
    expect(result).toHaveProperty('currentMode');
    expect(result).toHaveProperty('setTheme');
    expect(result).toHaveProperty('toggleTheme');
  });

  it('isDark reflects the theme store getter', () => {
    const store = useThemeStore();
    store.resolvedTheme = 'dark';

    const { isDark } = useTheme();
    expect(isDark.value).toBe(true);
  });

  it('currentMode reflects the theme store getter', () => {
    useThemeStore();

    const { currentMode } = useTheme();
    expect(currentMode.value).toBe('system');
  });

  it('setTheme delegates to themeStore.setTheme', () => {
    const store = useThemeStore();
    const spy = vi.spyOn(store, 'setTheme');

    const { setTheme } = useTheme();
    setTheme('dark');

    expect(spy).toHaveBeenCalledWith('dark');
  });

  it('toggleTheme calls store.toggleTheme and applies to Quasar dark mode', () => {
    const store = useThemeStore();
    store.setTheme('light');
    mockDarkSet.mockClear();

    const { toggleTheme } = useTheme();
    toggleTheme();

    // After toggle from light, resolved should be dark
    expect(store.resolvedTheme).toBe('dark');
    expect(mockDarkSet).toHaveBeenCalledWith(true);
  });

  it('toggleTheme from dark goes to light', () => {
    const store = useThemeStore();
    store.setTheme('dark');
    mockDarkSet.mockClear();

    const { toggleTheme } = useTheme();
    toggleTheme();

    expect(store.resolvedTheme).toBe('light');
    expect(mockDarkSet).toHaveBeenCalledWith(false);
  });
});
