<script setup lang="ts">
/** TYPE IMPORTS */
import type { IconSectionNavProps, IconSectionNavEmits } from './interfaces/IconSectionNav.interface';

/** VUE IMPORTS */
import { computed } from 'vue';

/** LOCAL IMPORTS (constants and handlers ONLY) */
import { ICON_SECTION_NAV_DEFAULT_WIDTH } from './constants';

/** PROPS & EMITS */
const props = withDefaults(defineProps<IconSectionNavProps>(), {
  width: ICON_SECTION_NAV_DEFAULT_WIDTH,
});

const emit = defineEmits<IconSectionNavEmits>();

/** COMPUTED */

/**
 * Inline style for the rail width
 */
const railStyle = computed(() => ({
  width: `${props.width}px`,
  minWidth: `${props.width}px`,
}));

/** FUNCTIONS */

/**
 * Handle item click — emit new active section
 *
 * @param {string} name - Section name that was clicked
 */
function handleItemClick(name: string): void {
  emit('update:modelValue', name);
}

/**
 * Check if a section is currently active
 *
 * @param {string} name - Section name to check
 * @returns {boolean} Whether the section is active
 */
function isActive(name: string): boolean {
  return props.modelValue === name;
}
</script>

<template>
  <nav
    class="icon-section-nav"
    :style="railStyle"
    role="tablist"
    aria-orientation="vertical"
  >
    <button
      v-for="item in items"
      :key="item.name"
      class="icon-section-nav__item"
      :class="{ 'icon-section-nav__item--active': isActive(item.name) }"
      role="tab"
      :aria-selected="isActive(item.name)"
      :aria-label="item.tooltip"
      @click="handleItemClick(item.name)"
    >
      <q-icon
        :name="item.icon"
        size="20px"
        :color="isActive(item.name) ? 'primary' : 'grey-6'"
      />
      <q-badge
        v-if="item.badge"
        :color="item.badgeColor ?? 'primary'"
        floating
        rounded
        class="icon-section-nav__badge"
      />
      <q-tooltip
        anchor="center right"
        self="center left"
        :offset="[8, 0]"
      >
        {{ item.tooltip }}
      </q-tooltip>
    </button>
  </nav>
</template>

<style lang="scss" scoped>
.icon-section-nav {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: var(--mapex-spacing-xs) 0;
  border-right: 1px solid var(--mapex-card-border);
  background: var(--mapex-surface-bg);
  flex-shrink: 0;

  &__item {
    position: relative;
    display: flex;
    align-items: center;
    justify-content: center;
    width: 100%;
    height: 40px;
    border: none;
    background: transparent;
    cursor: pointer;
    transition: var(--mapex-transition-fast);
    border-left: 3px solid transparent;
    padding: 0;

    &:hover {
      background: rgba(var(--mapex-primary-rgb), 0.05);
    }

    &--active {
      background: var(--mapex-active-bg);
      border-left-color: var(--q-primary);
    }
  }

  &__badge {
    position: absolute;
    top: 6px;
    right: 6px;
    min-width: 8px;
    min-height: 8px;
    width: 8px;
    height: 8px;
    padding: 0;
  }
}
</style>
