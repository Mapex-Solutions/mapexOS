<script setup lang="ts">
/** TYPE IMPORTS */
import type { PluginCatalogGroupProps } from './interfaces/PluginCatalogGroup.interface';
import type { PluginNodeType } from '../../interfaces/CreateEditWorkflow.interface';

/** VUE IMPORTS */
import { ref } from 'vue';

/** COMPONENTS */
import { AppTooltip } from '@components/tooltips';
import PluginCatalogItem from '../PluginCatalogItem/PluginCatalogItem.vue';

/** UTILS */
import { buildDefaultConfig } from '@utils/workflow';

/** COMPOSABLES */
import { useWorkflowEditorState } from '../../composables';

/** STORES */
import { usePluginRegistryStore } from '@stores/pluginRegistry';

/** PROPS & EMITS */
defineProps<PluginCatalogGroupProps>();

/** COMPOSABLES & STORES */
const { addNode, viewport } = useWorkflowEditorState();
const pluginRegistry = usePluginRegistryStore();

/** STATE */

/**
 * Whether this group accordion is expanded
 */
const expanded = ref(true);

/**
 * Counter for unique node IDs from click-to-add
 */
const clickCounter = ref(0);

/** FUNCTIONS */

/**
 * Handle drag start from the popup menu item
 *
 * @param {DragEvent} event - HTML5 drag event
 * @param {string} nodeTypeStr - Node type identifier
 * @returns {void}
 */
function handleMenuDragStart(event: DragEvent, nodeTypeStr: string): void {
  if (!event.dataTransfer) return;
  event.dataTransfer.setData('application/workflow-node-type', nodeTypeStr);
  event.dataTransfer.effectAllowed = 'move';
}

/**
 * Add a node to the canvas on click.
 * Places the node in the visible area using the current viewport position.
 *
 * @param {PluginNodeType} nodeType - Node type to add
 * @returns {void}
 */
function handleMenuClick(nodeType: PluginNodeType): void {
  const offset = ++clickCounter.value * 30;
  const { x, y, zoom } = viewport.value;

  // Convert viewport top-left to flow coordinates + padding + stacking offset
  const position = {
    x: Math.round((-x + 120 + offset) / zoom),
    y: Math.round((-y + 80 + offset) / zoom),
  };

  const newId = `n_${nodeType.type.replace('/', '_')}_${Date.now()}`;

  const operationProp = nodeType.properties?.find(p => p.name === 'operation');
  const defaultOperation = operationProp?.default as string | undefined;

  addNode({
    id: newId,
    type: nodeType.type,
    position,
    config: buildDefaultConfig(nodeType, defaultOperation),
    label: nodeType.label,
  });
}

/**
 * Add a node from expanded catalog item click
 *
 * @param {string} nodeTypeStr - Node type identifier
 * @returns {void}
 */
function handleItemClick(nodeTypeStr: string): void {
  const nodeType = pluginRegistry.getNodeType(nodeTypeStr);
  if (!nodeType) return;
  handleMenuClick(nodeType);
}
</script>

<template>
  <div class="catalog-group">
    <!-- COLLAPSED: category icon with popup menu -->
    <template v-if="collapsed">
      <div class="catalog-group__header catalog-group__header--collapsed">
        <q-icon :name="group.icon" size="18px" />
        <AppTooltip :content="group.label" anchor="center right" self="center left" :offset="[4, 0]" />

        <!-- Popup menu with node items -->
        <q-menu
          anchor="top right"
          self="top left"
          :offset="[4, 0]"
          class="catalog-group__menu"
        >
          <div class="catalog-group__menu-header">
            <q-icon :name="group.icon" size="16px" class="q-mr-sm" />
            <span class="text-caption text-weight-medium text-uppercase">
              {{ group.label }}
            </span>
          </div>

          <q-list dense class="catalog-group__menu-list">
            <q-item
              v-for="nodeType in group.nodeTypes"
              :key="nodeType.type"
              v-close-popup
              clickable
              draggable="true"
              class="catalog-group__menu-item"
              @click="handleMenuClick(nodeType)"
              @dragstart="handleMenuDragStart($event, nodeType.type)"
            >
              <q-item-section side>
                <q-icon :name="nodeType.icon" :color="nodeType.color" size="18px" />
              </q-item-section>
              <q-item-section>
                <q-item-label>{{ nodeType.label }}</q-item-label>
                <q-item-label caption>{{ nodeType.description }}</q-item-label>
              </q-item-section>
            </q-item>
          </q-list>
        </q-menu>
      </div>
    </template>

    <!-- EXPANDED: accordion with items inline -->
    <template v-else>
      <div
        class="catalog-group__header"
        @click="expanded = !expanded"
      >
        <q-icon :name="group.icon" size="18px" class="q-mr-sm" />
        <span class="text-caption text-weight-medium text-uppercase">
          {{ group.label }}
        </span>
        <q-space />
        <q-icon
          :name="expanded ? 'expand_less' : 'expand_more'"
          size="18px"
        />
      </div>

      <q-slide-transition>
        <div v-show="expanded" class="catalog-group__items">
          <PluginCatalogItem
            v-for="nodeType in group.nodeTypes"
            :key="nodeType.type"
            :node-type="nodeType"
            :collapsed="false"
            @click="handleItemClick(nodeType.type)"
          />
        </div>
      </q-slide-transition>
    </template>
  </div>
</template>

<style lang="scss" scoped>
.catalog-group {
  &__header {
    display: flex;
    align-items: center;
    padding: 8px 12px;
    cursor: pointer;
    user-select: none;
    color: var(--mapex-text-secondary);
    transition: var(--mapex-transition-fast);

    &:hover {
      background: var(--mapex-surface-elevated);
    }

    &--collapsed {
      justify-content: center;
      padding: 10px 8px;
    }
  }

  &__items {
    padding: 0 0 4px;
  }

  &__menu {
    background: var(--mapex-surface-elevated) !important;
    border: 1px solid var(--mapex-card-border) !important;
    border-radius: var(--mapex-radius-md) !important;
    box-shadow: var(--mapex-shadow-md) !important;
    min-width: 220px !important;
    max-width: 280px !important;
  }

  &__menu-header {
    display: flex;
    align-items: center;
    padding: 8px 12px 4px;
    color: var(--mapex-text-secondary);
    border-bottom: 1px solid var(--mapex-card-border);
    margin-bottom: 4px;
  }

  &__menu-list {
    padding: 4px 0;
  }

  &__menu-item {
    cursor: grab;

    &:active {
      cursor: grabbing;
      opacity: 0.7;
    }
  }
}
</style>
