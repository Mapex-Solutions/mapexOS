import type { Directive, DirectiveBinding } from 'vue';
import { usePermissionStore } from '@stores/permission';
import { watch } from 'vue';

/**
 * Vue directive for permission-based element visibility.
 * Hides elements (display: none) when the user lacks required permissions.
 *
 * Uses display:none instead of removing from DOM to avoid layout thrashing.
 *
 * @example
 * ```vue
 * <!-- Single permission -->
 * <q-btn v-permission="'assets.create'" label="Add" />
 *
 * <!-- Any of (default modifier) -->
 * <q-btn v-permission:any="['assets.create', 'assets.update']" label="Save" />
 *
 * <!-- All of -->
 * <q-btn v-permission:all="['assets.create', 'assets.update']" label="Full Edit" />
 * ```
 */

/**
 * Update element visibility based on permission check result
 *
 * @param {HTMLElement} el - DOM element
 * @param {DirectiveBinding} binding - Vue directive binding
 */
function updateVisibility(el: HTMLElement, binding: DirectiveBinding): void {
  const store = usePermissionStore();
  const mode = binding.arg || 'any';
  const value = binding.value;

  if (!value) {
    el.style.display = '';
    return;
  }

  let hasAccess = false;

  if (typeof value === 'string') {
    hasAccess = store.hasPermission(value);
  } else if (Array.isArray(value)) {
    if (mode === 'all') {
      hasAccess = store.hasAllPermissions(value);
    } else {
      hasAccess = store.hasAnyPermission(value);
    }
  }

  el.style.display = hasAccess ? '' : 'none';
}

export const vPermission: Directive = {
  mounted(el: HTMLElement, binding: DirectiveBinding) {
    // Initial check
    updateVisibility(el, binding);

    // Watch for permission store changes (reactive)
    const store = usePermissionStore();
    const stopWatch = watch(
      () => [store.permissions, store.version],
      () => updateVisibility(el, binding),
      { deep: true },
    );

    // Store cleanup function on element
    (el as any).__permissionCleanup = stopWatch;
  },

  updated(el: HTMLElement, binding: DirectiveBinding) {
    updateVisibility(el, binding);
  },

  unmounted(el: HTMLElement) {
    // Clean up watcher
    const cleanup = (el as any).__permissionCleanup;
    if (typeof cleanup === 'function') {
      cleanup();
    }
  },
};
