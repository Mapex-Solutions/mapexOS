<template>
  <!-- Invisible backdrop for click outside detection -->
  <Teleport to="body">
    <div
      v-if="modelValue"
      class="drawer-backdrop"
      :style="{ right: `${width}px` }"
      @click="handleClose"
    />
  </Teleport>

  <q-drawer
    overlay
    bordered
    side="right"
    :model-value="modelValue"
    :width="width"
    @update:model-value="emit('update:modelValue', $event)"
    @keydown.esc="handleClose"
  >
    <!-- Header -->
    <q-toolbar class="drawer-header">
      <q-icon
        v-if="icon"
        :name="icon"
        size="sm"
        :color="iconColor"
        class="q-mr-sm"
      />
      <q-toolbar-title class="text-weight-medium">{{ title }}</q-toolbar-title>

      <q-btn
        flat
        round
        dense
        icon="close"
        class="drawer-close-btn"
        @click="handleClose"
      >
        <AppTooltip :content="closeTooltip" />
      </q-btn>
    </q-toolbar>

    <q-separator />

    <!-- Scrollable Content -->
    <div class="drawer-content">
      <q-scroll-area class="fit">
        <div class="q-pa-md">
          <slot />
        </div>
      </q-scroll-area>
    </div>

    <!-- Footer (optional) -->
    <template v-if="$slots.footer">
      <q-separator />
      <div class="drawer-footer">
        <slot name="footer" />
      </div>
    </template>
  </q-drawer>
</template>

<script setup lang="ts">
defineOptions({
  name: 'GenericDrawer'
});

/** TYPE IMPORTS */
import type { GenericDrawerProps, GenericDrawerEmits } from './interfaces';

/** VUE IMPORTS */
import { onMounted, onBeforeUnmount } from 'vue';

/** COMPONENTS */
import { AppTooltip } from '@components/tooltips';

/** PROPS & EMITS */
const props = withDefaults(defineProps<GenericDrawerProps>(), {
  icon: 'info',
  iconColor: 'primary',
  width: 380,
  closeTooltip: 'Close',
});
const emit = defineEmits<GenericDrawerEmits>();

/** FUNCTIONS */

/**
 * Handle ESC key to close drawer
 * @param {KeyboardEvent} event - Keyboard event
 */
function handleEscKey(event: KeyboardEvent): void {
  if (event.key === 'Escape' && props.modelValue) {
    handleClose();
  }
}

/**
 * Close drawer and emit events
 */
function handleClose(): void {
  emit('update:modelValue', false);
  emit('close');
}

/** LIFECYCLE HOOKS */
onMounted(() => {
  window.addEventListener('keydown', handleEscKey);
});

onBeforeUnmount(() => {
  window.removeEventListener('keydown', handleEscKey);
});
</script>

<style lang="scss" scoped>
// Flex layout for drawer content
:deep(.q-drawer__content) {
  display: flex;
  flex-direction: column;
  height: 100%;
}

// Drawer Header - Fixed at top
.drawer-header {
  flex-shrink: 0;
  background: var(--mapex-header-bg);
  backdrop-filter: blur(10px);
  border-bottom: 1px solid var(--mapex-header-border);

  .q-toolbar__title {
    font-size: 1.1rem;
    color: var(--q-primary);
  }
}

// Close button - uses theme-aware muted color
.drawer-close-btn {
  color: var(--mapex-text-secondary);
}

// Backdrop (teleported to body, needs :global) - transparent, just for click detection
:global(.drawer-backdrop) {
  position: fixed;
  top: 0;
  left: 0;
  bottom: 0;
  background: transparent;
  z-index: 5999; // Below q-drawer (6000)
  cursor: default;
}

// Drawer Content - Scrollable middle section
.drawer-content {
  flex: 1;
  min-height: 0; // Important for flex children with overflow
  overflow: hidden;

  :deep(.q-scrollarea__content) {
    width: 100%;
    max-width: 100%;
    overflow-x: hidden;
  }
}

// Drawer Footer - Fixed at bottom
.drawer-footer {
  flex-shrink: 0;
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 16px;
  background: var(--mapex-header-bg);
  backdrop-filter: blur(10px);
  border-top: 1px solid var(--mapex-header-border);
  box-shadow: 0 -2px 8px var(--mapex-elevation-shadow);
}

// Custom Scrollbar
:deep(.q-scrollarea__content) {
  &::-webkit-scrollbar {
    width: 6px;
  }

  &::-webkit-scrollbar-track {
    background: transparent;
    border-radius: var(--mapex-radius-lg);
  }

  &::-webkit-scrollbar-thumb {
    background: rgba(var(--q-primary-rgb), 0.3);
    border-radius: var(--mapex-radius-lg);
    transition: background var(--mapex-transition-base) ease;

    &:hover {
      background: rgba(var(--q-primary-rgb), 0.5);
    }
  }
}
</style>
