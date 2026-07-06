<script setup lang="ts">
/** TYPE IMPORTS */
import type { MenuItem } from '../../interfaces';

/** PROPS & EMITS */
const props = defineProps<{
  item: MenuItem;
  currentRoute: string;
}>();

/** FUNCTIONS */

/**
 * Whether any descendant leaf route matches the current route, used to
 * auto-open and highlight ancestor branches at any depth.
 *
 * @param {MenuItem} node - Menu node to test
 * @returns {boolean} True when a descendant route is active
 */
function hasActiveDescendant(node: MenuItem): boolean {
  if (node.to && props.currentRoute.startsWith(node.to)) return true;
  return (node.children ?? []).some(hasActiveDescendant);
}
</script>

<template>
  <!-- Separator -->
  <q-separator v-if="item.separator" class="q-my-xs" />

  <!-- Branch: nested expansion, recurses to any depth -->
  <q-expansion-item
    v-else-if="item.children?.length"
    dense
    expand-icon="keyboard_arrow_down"
    expand-icon-class="text-grey-7"
    header-class="menu-parent-item"
    :icon="item.icon"
    :label="item.label"
    :default-opened="hasActiveDescendant(item)"
    :icon-color="hasActiveDescendant(item) ? 'primary' : 'grey-7'"
  >
    <q-list class="submenu-list">
      <MenuTreeItem
        v-for="(child, cIndex) in item.children"
        :key="cIndex"
        :item="child"
        :current-route="currentRoute"
      />
    </q-list>
  </q-expansion-item>

  <!-- Leaf -->
  <q-item
    v-else
    v-ripple
    clickable
    active-class="active-menu-item"
    class="submenu-item"
    :to="item.to"
    :active="currentRoute === item.to"
  >
    <q-item-section avatar>
      <q-icon
        :name="item.icon || 'chevron_right'"
        :color="currentRoute === item.to ? 'primary' : 'grey-7'"
      />
    </q-item-section>
    <q-item-section>{{ item.label }}</q-item-section>
  </q-item>
</template>

<style lang="scss" scoped>
// Scoped styles do not cross component boundaries, so the submenu look is
// declared here too (mirrors the sidebar) and each nesting level adds its own
// left padding, giving natural indentation by depth.
.submenu-list {
  padding-left: var(--mapex-spacing-lg) !important;
  background: var(--mapex-submenu-bg);
}

.submenu-item {
  padding: var(--mapex-spacing-sm) var(--mapex-spacing-lg);
  min-height: 40px;
  border-radius: 0 var(--mapex-radius-md) var(--mapex-radius-md) 0;
  margin: 2px 0;

  &:hover {
    background: rgba(var(--q-primary-rgb), 0.05);
  }
}

.active-menu-item {
  background: rgba(var(--q-primary-rgb), 0.1);
  border-right: 3px solid var(--q-primary);
  color: var(--q-primary);
}

:deep(.menu-parent-item) {
  border-radius: 0 var(--mapex-radius-md) var(--mapex-radius-md) 0;
  margin: 2px 0;

  &:hover {
    background: rgba(var(--q-primary-rgb), 0.05);
  }
}
</style>
