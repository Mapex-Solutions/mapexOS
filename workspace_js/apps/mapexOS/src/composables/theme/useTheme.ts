import { watch, onMounted } from 'vue';
import { useQuasar } from 'quasar';
import { useThemeStore } from '@stores/theme';
import { storeToRefs } from 'pinia';

/**
 * Composable that bridges the Pinia theme store with Quasar's dark mode API.
 * Must be called inside setup() of a component with access to $q.
 *
 * @returns {{ isDark: Ref<boolean>, currentMode: Ref<ThemeMode>, setTheme: (mode: ThemeMode) => void, toggleTheme: () => void }} Theme reactive state and actions
 */
export function useTheme() {
  const $q = useQuasar();
  const themeStore = useThemeStore();
  const { isDark, currentMode } = storeToRefs(themeStore);

  /**
   * Apply resolved theme to Quasar's dark mode
   *
   * @returns {void}
   */
  function applyTheme(): void {
    $q.dark.set(themeStore.isDark);
  }

  // Watch store changes and sync to Quasar
  watch(() => themeStore.resolvedTheme, () => {
    applyTheme();
  });

  // Initialize theme and listen for system preference changes
  onMounted(() => {
    themeStore.initTheme();
    applyTheme();

    window.matchMedia('(prefers-color-scheme: dark)')
      .addEventListener('change', () => {
        if (themeStore.currentMode === 'system') {
          themeStore.setTheme('system');
          applyTheme();
        }
      });
  });

  return {
    isDark,
    currentMode,
    setTheme: themeStore.setTheme.bind(themeStore),
    toggleTheme: () => {
      themeStore.toggleTheme();
      applyTheme();
    },
  };
}
