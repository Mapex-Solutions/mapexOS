<template>
  <q-drawer
    overlay
    bordered
    side="right"
    :model-value="modelValue"
    :width="450"
    @update:model-value="emit('update:modelValue', $event)"
  >
    <!-- Header with Glassmorphism -->
    <q-toolbar class="drawer-header">
      <q-icon name="account_tree" size="sm" class="q-mr-sm" color="primary" />
      <q-toolbar-title class="text-weight-medium">{{ t.title.value }}</q-toolbar-title>

      <!-- Info button for legend -->
      <q-btn
        flat
        round
        dense
        icon="info"
        color="primary"
        @click="showLegendModal = true"
      >
        <AppTooltip :content="t.legendTooltip.value" />
      </q-btn>

      <q-btn flat round dense icon="close" class="drawer-close-btn" @click="close">
        <AppTooltip :content="t.close.value" />
      </q-btn>
    </q-toolbar>

    <!-- Filters Section -->
    <div class="filters-section q-pa-md">
      <!-- Search by name with debounce -->
      <q-input
        v-model="filters.name"
        outlined
        dense
        class="modern-input"
        :label="t.filters.searchLabel.value"
        :placeholder="t.filters.searchPlaceholder.value"
        :debounce="300"
      >
        <template v-slot:prepend>
          <q-icon name="search" color="primary" />
        </template>
        <template v-slot:append v-if="filters.name">
          <q-icon
            name="clear"
            class="cursor-pointer"
            color="grey-7"
            @click="filters.name = ''"
          />
        </template>
      </q-input>

      <!-- Type filter (multi-select) -->
      <q-select
        v-model="filters.types"
        multiple
        outlined
        dense
        emit-value
        map-options
        class="q-mt-md modern-input"
        :options="typeOptions"
        :label="t.filters.typeLabel.value"
      >
        <template v-slot:prepend>
          <q-icon name="filter_list" />
        </template>
        <template v-slot:option="scope">
          <q-item v-bind="scope.itemProps">
            <q-item-section avatar>
              <q-icon
                :name="getOrganizationIcon(scope.opt.value)"
                :color="getOrganizationColor(scope.opt.value)"
              />
            </q-item-section>
            <q-item-section>
              <q-item-label>{{ scope.opt.label }}</q-item-label>
            </q-item-section>
          </q-item>
        </template>
        <template v-slot:selected-item="scope">
          <SelectableChip
            :label="scope.opt.label"
            :icon="getOrganizationIcon(scope.opt.value)"
            :color="getOrganizationColor(scope.opt.value) as any"
            dense
            size="sm"
            class="q-mt-sm q-mr-xs"
            @remove="scope.removeAtIndex(scope.index)"
          />
        </template>
      </q-select>

      <!-- Enabled filter (toggle) -->
      <q-btn-toggle
        v-model="filters.enabled"
        unelevated
        size="sm"
        class="q-mt-md full-width modern-toggle"
        toggle-color="primary"
        :options="enabledOptions"
      />
    </div>

    <q-separator />

    <!-- Organization Tree -->
    <div class="tree-container q-pa-md">
      <q-tree
        v-if="treeNodes.length > 0"
        v-model:expanded="expanded"
        ref="treeRef"
        node-key="id"
        label-key="name"
        children-key="children"
        class="modern-tree"
        :nodes="treeNodes"
      >
        <template v-slot:default-header="prop">
          <div
            class="tree-node-item row items-center full-width no-wrap"
            @click.stop="onNodeExpand(prop.node)"
            @dblclick.stop="onOrganizationClick(prop.node)"
          >
            <!-- Expand/Collapse Arrow (only for nodes with children) -->
            <q-icon
              v-if="prop.node.children && prop.node.children.length > 0"
              size="sm"
              class="q-mr-xs expand-arrow"
              :name="expanded.includes(prop.node.id) ? 'expand_more' : 'chevron_right'"
              color="primary"
            />

            <!-- Organization Type Icon -->
            <q-icon
              size="sm"
              class="q-mr-sm"
              :name="getOrganizationIcon(prop.node.type)"
              :color="getOrganizationColor(prop.node.type)"
            />

            <!-- Organization Name -->
            <div class="text-body2 ellipsis">{{ prop.node.name }}</div>

            <q-space />

            <!-- Inactive Badge -->
            <q-badge
              v-if="!prop.node.enabled"
              class="q-ml-sm"
              :color="'grey'"
              :label="t.inactive.value"
            />

            <!-- Child Count Badge -->
            <q-badge
              v-if="prop.node.childCount > 0"
              class="q-ml-xs"
              color="primary"
              :label="prop.node.childCount"
            />
          </div>
        </template>
      </q-tree>

      <!-- Empty state -->
      <div v-else-if="!loading" class="empty-state">
        <q-icon name="search_off" size="64px" color="grey-5" />
        <div class="text-subtitle1 q-mt-md text-grey-7">{{ t.empty.value }}</div>
        <div class="text-caption text-grey-5 q-mt-sm">Try adjusting your filters</div>
      </div>
    </div>

    <!-- Loading -->
    <q-inner-loading :showing="loading">
      <q-spinner-dots size="50px" color="primary" />
      <div class="text-primary q-mt-md">{{ t.loading.value }}</div>
    </q-inner-loading>

    <!-- Legend Modal -->
    <q-dialog v-model="showLegendModal">
      <q-card style="min-width: 400px; max-width: 500px;">
        <q-card-section class="row items-center q-pb-none">
          <div class="text-h6">{{ t.legend.title.value }}</div>
          <q-space />
          <q-btn v-close-popup flat round dense icon="close" />
        </q-card-section>

        <q-card-section>
          <div class="text-subtitle2 q-mb-md">{{ t.legend.description.value }}</div>

          <!-- Organization Types Legend -->
          <div class="q-mb-md">
            <div class="text-weight-medium q-mb-sm">{{ t.legend.typesTitle.value }}</div>
            <q-list dense>
              <q-item v-for="type in organizationTypes" :key="type.value">
                <q-item-section avatar>
                  <q-icon
                    size="sm"
                    :name="getOrganizationIcon(type.value)"
                    :color="getOrganizationColor(type.value)"
                  />
                </q-item-section>
                <q-item-section>
                  <q-item-label>{{ type.label }}</q-item-label>
                  <q-item-label caption>{{ type.description }}</q-item-label>
                </q-item-section>
              </q-item>
            </q-list>
          </div>

          <!-- Click Instructions -->
          <div class="q-mb-md">
            <div class="text-weight-medium q-mb-sm">{{ t.legend.actionsTitle.value }}</div>
            <q-list dense bordered class="rounded-borders">
              <q-item>
                <q-item-section avatar>
                  <q-icon name="touch_app" color="primary" />
                </q-item-section>
                <q-item-section>
                  <q-item-label>{{ t.legend.singleClick.value }}</q-item-label>
                  <q-item-label caption>{{ t.legend.singleClickDesc.value }}</q-item-label>
                </q-item-section>
              </q-item>

              <q-separator spaced />

              <q-item>
                <q-item-section avatar>
                  <q-icon name="double_arrow" color="primary" />
                </q-item-section>
                <q-item-section>
                  <q-item-label>{{ t.legend.doubleClick.value }}</q-item-label>
                  <q-item-label caption>{{ t.legend.doubleClickDesc.value }}</q-item-label>
                </q-item-section>
              </q-item>
            </q-list>
          </div>

          <!-- Status Badges -->
          <div>
            <div class="text-weight-medium q-mb-sm">{{ t.legend.statusTitle.value }}</div>
            <div class="row items-center q-gutter-sm">
              <q-badge color="grey" :label="t.inactive.value" />
              <span class="text-caption">{{ t.legend.inactiveDesc.value }}</span>
            </div>
          </div>
        </q-card-section>

        <q-card-actions class="row items-center q-px-md q-py-md">
          <q-space />
          <q-btn v-close-popup flat dense color="primary" :label="t.legend.close.value" size="sm" class="rounded-borders" />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-drawer>
</template>

<script setup lang="ts">
defineOptions({
  name: 'OrganizationTreeDrawer'
});

import type { OrganizationTreeDrawerProps, OrganizationTreeDrawerEmits } from './interfaces';
import type { OrganizationType } from '@stores/organization/types';

import { computed, ref, watch } from 'vue';
import { useRouter } from 'vue-router';
import type { QTree } from 'quasar';
import { SelectableChip } from '@components/chips';
import { AppTooltip } from '@components/tooltips';
import { getOrganizationIcon, getOrganizationColor } from 'src/utils/organization/icons';
import { useOrganizationTree } from 'src/composables/organizations/useOrganizationTree';
import { useOrganizationTreeDrawerTranslations } from 'src/composables/i18n/components/useOrganizationTreeDrawerTranslations';
import { useOrganizationStore } from '@stores/organization';

/**
 * OrganizationTreeDrawer Component
 *
 * Right-side drawer displaying hierarchical organization tree with filters.
 * Organizations have 6 levels: Vendor → Customer → Site → Building → Floor → Zone
 *
 * Features:
 * - Hierarchical q-tree display with icons (Quasar native expand/collapse)
 * - Name search filter (debounced)
 * - Type multi-select filter
 * - Enabled status toggle (all/active/inactive)
 * - Automatic expand/collapse on arrow click
 * - Single click on row to expand/collapse
 * - Double click on row to select organization context
 * - Persistent expansion state across filters
 * - 100% i18n compliance
 *
 * Usage:
 * ```vue
 * <OrganizationTreeDrawer
 *   v-model="showDrawer"
 *   @select="onOrganizationSelected"
 * />
 * ```
 */

defineProps<OrganizationTreeDrawerProps>();
const emit = defineEmits<OrganizationTreeDrawerEmits>();

const router = useRouter();
const organizationStore = useOrganizationStore();
const t = useOrganizationTreeDrawerTranslations();
const { treeNodes, loading, filters } = useOrganizationTree();

// Tree ref for programmatic control
const treeRef = ref<QTree>();

// State for legend modal
const showLegendModal = ref(false);

// Expansion state management
// Store expanded nodes persistently across filter changes
const expanded = ref<string[]>([]);

// Track if we've done initial expansion to avoid re-expansion on filter changes
const hasInitialExpanded = ref(false);

// Watch treeNodes and expand root level by default on first load only
watch(
  () => treeNodes.value,
  (newNodes) => {
    if (newNodes.length > 0 && !hasInitialExpanded.value) {
      // Expand only root level nodes (first level) on initial load
      expanded.value = newNodes.map((node) => node.id);
      hasInitialExpanded.value = true;
    }
  },
  { immediate: true }
);

/**
 * Type filter options with translated labels
 */
const typeOptions = computed<{ label: string; value: OrganizationType }[]>(() => [
  { label: t.types.vendor.value, value: 'vendor' },
  { label: t.types.customer.value, value: 'customer' },
  { label: t.types.site.value, value: 'site' },
  { label: t.types.building.value, value: 'building' },
  { label: t.types.floor.value, value: 'floor' },
  { label: t.types.zone.value, value: 'zone' },
]);

/**
 * Enabled status filter options with translated labels
 */
const enabledOptions = computed(() => [
  { label: t.filters.enabledAll.value, value: 'all' },
  { label: t.filters.enabledActive.value, value: 'active' },
  { label: t.filters.enabledInactive.value, value: 'inactive' },
]);

/**
 * Organization types with labels and descriptions for legend modal
 */
const organizationTypes = computed<{ value: OrganizationType; label: string; description: string }[]>(() => [
  {
    value: 'vendor',
    label: t.types.vendor.value,
    description: t.legend.vendorDesc.value,
  },
  {
    value: 'customer',
    label: t.types.customer.value,
    description: t.legend.customerDesc.value,
  },
  {
    value: 'site',
    label: t.types.site.value,
    description: t.legend.siteDesc.value,
  },
  {
    value: 'building',
    label: t.types.building.value,
    description: t.legend.buildingDesc.value,
  },
  {
    value: 'floor',
    label: t.types.floor.value,
    description: t.legend.floorDesc.value,
  },
  {
    value: 'zone',
    label: t.types.zone.value,
    description: t.legend.zoneDesc.value,
  },
]);

/**
 * Close drawer
 */
function close() {
  emit('update:modelValue', false);
}

/**
 * Handle node expand/collapse on single click
 * Toggles the expansion state of the node
 */
function onNodeExpand(node: any) {
  if (!treeRef.value || !node.children || node.children.length === 0) return;

  // Check if node is already expanded
  const isExpanded = treeRef.value.isExpanded(node.id);

  // Toggle expansion
  treeRef.value.setExpanded(node.id, !isExpanded);
}

/**
 * Handle organization selection on double click
 * Updates the store with selected organization and closes drawer
 * Automatically reloads the current page to refresh data for the new organization context
 */
function onOrganizationClick(node: any) {
  // Check if organization actually changed to avoid unnecessary reload
  const previousOrgId = organizationStore.selectedOrganizationId;
  const newOrgId = node.id;

  // If clicking the same organization, just close drawer without reload
  if (previousOrgId === newOrgId) {
    close();
    return;
  }

  // Update organization context in store
  organizationStore.selectOrganization(newOrgId);

  // Emit select event for parent components
  emit('select', newOrgId);

  // Close the drawer before reload to ensure clean state
  close();

  // Set flag in sessionStorage to indicate this is a manual org change (not initial load)
  // This prevents infinite reload loops
  sessionStorage.setItem('orgChangeReload', 'true');

  // Force reload of current page to refresh all data with new organization context
  // This ensures ANY page/route automatically gets fresh data when organization changes
  // Using router.go(0) provides a clean full page reload with proper loading states
  router.go(0);
}
</script>

<style lang="scss" scoped>
// Modern Drawer Header with Glassmorphism
.drawer-header {
  background: var(--mapex-header-bg);
  backdrop-filter: blur(10px);
  border-bottom: 1px solid var(--mapex-header-border);
  transition: all var(--mapex-transition-slow) ease;

  .q-toolbar__title {
    font-size: 1.1rem;
    color: var(--q-primary);
  }
}

// Close button
.drawer-close-btn {
  color: var(--mapex-text-secondary);
}

// Modern Filters Section
.filters-section {
  background: linear-gradient(135deg, rgba(var(--q-primary-rgb), 0.02) 0%, rgba(var(--q-primary-rgb), 0.05) 100%);
  border-bottom: 1px solid var(--mapex-divider);
}

// Modern Input Fields
.modern-input {
  :deep(.q-field__control) {
    border-radius: var(--mapex-radius-md);
    background: var(--mapex-surface-elevated);
    border: 1px solid var(--mapex-card-border);
    transition: all var(--mapex-transition-slow) ease;

    &:hover {
      border-color: var(--q-primary);
      box-shadow: 0 2px 8px rgba(var(--q-primary-rgb), 0.1);
    }
  }

  :deep(.q-field--focused .q-field__control) {
    border-color: var(--q-primary);
    box-shadow: 0 0 0 3px rgba(var(--q-primary-rgb), 0.1);
  }
}

// Modern Toggle Buttons
.modern-toggle {
  border-radius: var(--mapex-radius-md);
  overflow: hidden;
  box-shadow: 0 2px 4px var(--mapex-elevation-shadow);

  :deep(.q-btn) {
    flex: 1;
    border-radius: 0;
    transition: all var(--mapex-transition-base) ease;

    &:not(.q-btn--active) {
      background: var(--mapex-surface-elevated);
      color: var(--mapex-text-primary);

      &:hover {
        background: rgba(var(--q-primary-rgb), 0.05);
      }
    }
  }
}

// Tree Container with Custom Scrollbar
.tree-container {
  height: calc(100vh - 320px);
  overflow-y: auto;
  overflow-x: hidden;

  // Custom Scrollbar
  &::-webkit-scrollbar {
    width: 6px;
    background: transparent;
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

// Modern Tree Styling with Proper Hierarchy
.modern-tree {
  // Remove Quasar's default connector lines (we use custom template)
  :deep(.q-tree__node::before),
  :deep(.q-tree__node::after),
  :deep(.q-tree__node-header::before) {
    display: none !important;
  }

  // Basic spacing for nodes
  :deep(.q-tree__node) {
    padding: 2px 0;
  }

  // Ensure proper padding for the node header content
  :deep(.q-tree__node-header) {
    padding: 4px 0;
    cursor: pointer;
  }

  // Hide Quasar's default arrow (we have our own in custom template)
  :deep(.q-tree__arrow) {
    display: none !important;
  }

  // Smooth icon animations
  :deep(.q-icon) {
    transition: transform var(--mapex-transition-base) ease;
  }

  // Ensure proper indentation for child nodes
  // Each level gets indentation
  :deep(.q-tree__children) {
    padding-left: 16px;
  }
}

// Tree Node Item with Modern Hover
.tree-node-item {
  padding: 10px 12px;
  border-radius: var(--mapex-radius-md);
  margin: 2px 0;
  cursor: pointer;
  transition: all var(--mapex-transition-base) ease;
  border: 1px solid transparent;

  &:hover {
    background: rgba(var(--q-primary-rgb), 0.08);
    border-color: rgba(var(--q-primary-rgb), 0.2);
    transform: translateX(2px);
    box-shadow: 0 2px 8px rgba(var(--q-primary-rgb), 0.15);

    .q-icon {
      transform: scale(1.1);
    }
  }

  &:active {
    transform: translateX(0);
    background: rgba(var(--q-primary-rgb), 0.12);
  }

  .text-body2 {
    font-weight: 500;
    color: var(--mapex-text-primary);
  }
}

// Modern Empty State
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 20px;
  text-align: center;
  animation: fadeIn 0.4s ease;

  .q-icon {
    opacity: 0.5;
    animation: pulse 2s ease-in-out infinite;
  }
}

// Animations
@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@keyframes pulse {
  0%, 100% {
    transform: scale(1);
    opacity: 0.5;
  }
  50% {
    transform: scale(1.05);
    opacity: 0.7;
  }
}

// Chip Styling in Select
:deep(.q-chip) {
  border-radius: var(--mapex-radius-xl);
  font-weight: 500;
  box-shadow: 0 2px 4px var(--mapex-elevation-shadow);
}

// Badge Styling
:deep(.q-badge) {
  border-radius: var(--mapex-radius-lg);
  font-weight: 500;
  padding: 4px 8px;
  box-shadow: 0 2px 4px var(--mapex-elevation-shadow);
}

// Separator Styling
:deep(.q-separator) {
  background: var(--mapex-divider);
}

// Footer padding (ensure proper spacing for modal actions)
:deep(.q-card__actions) {
  padding: 16px !important;
}
</style>
