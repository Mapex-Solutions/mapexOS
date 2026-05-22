<script setup lang="ts">
/** TYPE IMPORTS */
import type { PluginCatalogItemProps } from './interfaces/PluginCatalogItem.interface';

/** COMPONENTS */
import { AppTooltip } from '@components/tooltips';

/** PROPS & EMITS */
const props = defineProps<PluginCatalogItemProps>();

/** FUNCTIONS */

/**
 * Handle drag start — set node type in dataTransfer
 *
 * @param {DragEvent} event - HTML5 drag event
 * @returns {void}
 */
function handleDragStart(event: DragEvent): void {
  if (!event.dataTransfer) return;

  event.dataTransfer.setData('application/workflow-node-type', props.nodeType.type);
  event.dataTransfer.effectAllowed = 'move';
}
</script>

<template>
  <div
    class="catalog-item"
    :class="{ 'catalog-item--collapsed': collapsed }"
    draggable="true"
    @dragstart="handleDragStart"
  >
    <q-icon
      :name="nodeType.icon"
      size="18px"
      :color="nodeType.color"
      class="catalog-item__icon"
    />
    <template v-if="!collapsed">
      <span class="catalog-item__label">{{ nodeType.label }}</span>
    </template>

    <!-- Tooltip (always show on hover) -->
    <AppTooltip anchor="center right" self="center left" :offset="[8, 0]">
      <div class="text-weight-medium">{{ nodeType.label }}</div>
      <div class="text-caption">{{ nodeType.description }}</div>
      <div class="text-caption text-grey-5 q-mt-xs">
        {{ nodeType.type }}
      </div>
    </AppTooltip>
  </div>
</template>

<style lang="scss" scoped>
.catalog-item {
  display: flex;
  align-items: center;
  padding: 6px 12px 6px 24px;
  cursor: grab;
  user-select: none;
  transition: var(--mapex-transition-fast);
  border-radius: var(--mapex-radius-sm);
  margin: 0 4px;

  &:hover {
    background: var(--mapex-surface-elevated);
  }

  &:active {
    cursor: grabbing;
    opacity: 0.7;
  }

  &--collapsed {
    justify-content: center;
    padding: 8px;
  }

  &__icon {
    flex-shrink: 0;
  }

  &__label {
    margin-left: 8px;
    font-size: 13px;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }
}
</style>
