<script setup lang="ts">
defineOptions({
  name: 'OrganizationSelectorDrawer'
});

/** TYPE IMPORTS */
import type { OrganizationSelectorDrawerProps, OrganizationSelectorDrawerEmits } from './interfaces';
import type { OrganizationTreeNode, OrganizationType } from '@stores/organization/types';

/** VUE IMPORTS */
import { ref, computed, watch, onMounted, onBeforeUnmount } from 'vue';

/** COMPOSABLES */
import { useOrganizationTree } from '@composables/organizations/useOrganizationTree';

/** UTILS */
import { getOrganizationIcon, getOrganizationColor } from '@utils/organization/icons';

/** PROPS & EMITS */
const props = withDefaults(defineProps<OrganizationSelectorDrawerProps>(), {
  selectedOrganizationId: null,
});

const emit = defineEmits<OrganizationSelectorDrawerEmits>();

/** COMPOSABLES */
const { treeNodes, loading, filters } = useOrganizationTree();

/** STATE */
const expanded = ref<string[]>([]);
const hasInitialExpanded = ref(false);

/** COMPUTED */

/**
 * Drawer visibility model
 */
const showDialog = computed({
  get: () => props.modelValue,
  set: (value: boolean) => emit('update:modelValue', value),
});

/**
 * Type filter options
 */
const typeOptions = computed(() => [
  { label: 'All Types', value: undefined },
  { label: 'Vendor', value: 'vendor' as OrganizationType },
  { label: 'Customer', value: 'customer' as OrganizationType },
  { label: 'Site', value: 'site' as OrganizationType },
  { label: 'Building', value: 'building' as OrganizationType },
  { label: 'Floor', value: 'floor' as OrganizationType },
  { label: 'Zone', value: 'zone' as OrganizationType },
]);

/**
 * Enabled filter options
 */
const enabledOptions = computed(() => [
  { label: 'All', value: 'all' },
  { label: 'Active', value: 'active' },
  { label: 'Inactive', value: 'inactive' },
]);

/**
 * Flatten tree to count total organizations
 */
const totalOrganizations = computed(() => {
  function countNodes(nodes: OrganizationTreeNode[]): number {
    return nodes.reduce((acc, node) => {
      return acc + 1 + (node.children ? countNodes(node.children) : 0);
    }, 0);
  }
  return countNodes(treeNodes.value);
});

/** WATCHERS */

/**
 * Watch drawer open state
 */
watch(() => props.modelValue, (isOpen) => {
  if (isOpen) {
    // Reset filters when drawer opens
    filters.value.name = '';
    hasInitialExpanded.value = false;
  }
});

/**
 * Watch treeNodes and expand root level by default
 */
watch(
  () => treeNodes.value,
  (newNodes) => {
    if (newNodes.length > 0 && !hasInitialExpanded.value) {
      expanded.value = newNodes.map((node) => node.id);
      hasInitialExpanded.value = true;
    }
  },
  { immediate: true }
);

/** FUNCTIONS */

/**
 * Toggle node expansion
 * @param {OrganizationTreeNode} node - Node to toggle
 */
function toggleExpand(node: OrganizationTreeNode): void {
  const index = expanded.value.indexOf(node.id);
  if (index >= 0) {
    expanded.value.splice(index, 1);
  } else {
    expanded.value.push(node.id);
  }
}

/**
 * Check if node is expanded
 * @param {OrganizationTreeNode} node - Node to check
 * @returns {boolean} True if expanded
 */
function isExpanded(node: OrganizationTreeNode): boolean {
  return expanded.value.includes(node.id);
}

/**
 * Select organization and close drawer
 * @param {OrganizationTreeNode} node - Organization node to select
 */
function selectOrganization(node: OrganizationTreeNode): void {
  emit('select', {
    id: node.id,
    name: node.name,
    type: node.type,
    enabled: node.enabled,
    pathKey: node.pathKey,
  });
  showDialog.value = false;
}

/**
 * Check if organization is selected
 * @param {OrganizationTreeNode} node - Organization to check
 * @returns {boolean} True if selected
 */
function isSelected(node: OrganizationTreeNode): boolean {
  return node.id === props.selectedOrganizationId;
}

/**
 * Cancel handler
 */
function handleCancel(): void {
  emit('cancel');
  showDialog.value = false;
}

/**
 * Handle ESC key to close drawer
 * @param {KeyboardEvent} event - Keyboard event
 */
function handleEscKey(event: KeyboardEvent): void {
  if (event.key === 'Escape' && props.modelValue) {
    handleCancel();
  }
}

/**
 * Handle type filter change
 * @param {OrganizationType | undefined} type - Selected type
 */
function handleTypeFilter(type: OrganizationType | undefined): void {
  if (type) {
    filters.value.types = [type];
  } else {
    filters.value.types = [];
  }
}

/**
 * Get current type filter value
 */
const currentTypeFilter = computed(() => {
  return filters.value.types?.length === 1 ? filters.value.types[0] : undefined;
});

/** LIFECYCLE HOOKS */
onMounted(() => {
  window.addEventListener('keydown', handleEscKey);
});

onBeforeUnmount(() => {
  window.removeEventListener('keydown', handleEscKey);
});
</script>

<template>
  <q-dialog v-model="showDialog" position="right" maximized>
    <q-card style="width: 500px; max-width: 90vw; display: flex; flex-direction: column; height: 100vh;">
      <!-- Header -->
      <q-card-section class="q-pb-sm">
        <div class="row items-center">
          <q-icon name="account_tree" size="md" color="primary" class="q-mr-sm" />
          <div class="text-h6">Select Organization</div>
          <q-space />
          <q-btn icon="close" flat round dense class="rounded-borders" @click="handleCancel" />
        </div>
      </q-card-section>

      <!-- Info Banner -->
      <q-card-section class="q-pt-none q-pb-md">
        <q-banner dense class="bg-teal-1 text-teal-9 rounded-borders">
          <template #avatar>
            <q-icon name="info" color="teal-6" size="sm" />
          </template>
          <div class="text-caption">
            Select the organization where the user will have access. Click on an organization to select it.
          </div>
        </q-banner>
      </q-card-section>

      <!-- Filters -->
      <q-card-section class="q-py-md">
        <div class="text-overline text-grey-7 q-mb-md">
          <q-icon name="filter_list" size="xs" class="q-mr-xs" />
          Filters
        </div>
        <div class="row q-col-gutter-md">
          <!-- Search -->
          <div class="col-12">
            <q-input
              v-model="filters.name"
              outlined
              dense
              label="Search by name"
              placeholder="Type to search..."
              clearable
              debounce="300"
              class="rounded-borders"
            >
              <template #prepend>
                <q-icon name="search" />
              </template>
            </q-input>
          </div>

          <!-- Type and Status -->
          <div class="col-12 col-sm-6">
            <q-select
              :model-value="currentTypeFilter"
              outlined
              dense
              label="Type"
              class="rounded-borders"
              :options="typeOptions"
              emit-value
              map-options
              @update:model-value="handleTypeFilter"
            >
              <template #prepend>
                <q-icon name="category" />
              </template>
            </q-select>
          </div>

          <div class="col-12 col-sm-6">
            <q-select
              v-model="filters.enabled"
              outlined
              dense
              label="Status"
              class="rounded-borders"
              :options="enabledOptions"
              emit-value
              map-options
            >
              <template #prepend>
                <q-icon name="toggle_on" />
              </template>
            </q-select>
          </div>
        </div>
      </q-card-section>

      <q-separator />

      <!-- Results Header -->
      <q-card-section class="q-py-sm">
        <div class="text-overline text-grey-7">
          <q-icon name="account_tree" size="xs" class="q-mr-xs" />
          Organizations
        </div>
      </q-card-section>

      <!-- Organizations Tree -->
      <q-card-section class="q-pa-none" style="flex: 1; overflow: hidden;">
        <!-- Loading state -->
        <div v-if="loading" class="q-pa-md text-center">
          <q-spinner color="primary" size="3em" />
          <div class="text-grey-7 q-mt-md">Loading organizations...</div>
        </div>

        <!-- Empty state -->
        <div v-else-if="treeNodes.length === 0" class="q-pa-md text-center">
          <q-icon name="search_off" size="4em" color="grey-5" />
          <div class="text-grey-7 q-mt-md">No organizations found</div>
          <div class="text-caption text-grey-5 q-mt-sm">Try adjusting your filters</div>
        </div>

        <!-- Organizations Tree -->
        <q-scroll-area v-else style="height: 100%;">
          <div class="q-pa-md">
            <template v-for="node in treeNodes" :key="node.id">
              <div class="org-tree-node">
                <!-- Root node -->
                <div
                  class="org-tree-item rounded-borders q-mb-xs"
                  :class="{ 'org-tree-item--selected': isSelected(node) }"
                >
                  <div class="row items-center no-wrap q-pa-sm cursor-pointer" @click="selectOrganization(node)">
                    <!-- Expand arrow -->
                    <q-btn
                      v-if="node.children && node.children.length > 0"
                      flat
                      round
                      dense
                      size="sm"
                      :icon="isExpanded(node) ? 'expand_more' : 'chevron_right'"
                      color="primary"
                      class="q-mr-xs"
                      @click.stop="toggleExpand(node)"
                    />
                    <div v-else class="q-mr-md" style="width: 24px;" />

                    <!-- Icon -->
                    <q-icon
                      :name="getOrganizationIcon(node.type)"
                      :color="getOrganizationColor(node.type)"
                      size="sm"
                      class="q-mr-sm"
                    />

                    <!-- Name -->
                    <div class="col">
                      <div class="text-body2">{{ node.name }}</div>
                    </div>

                    <!-- Badges -->
                    <q-badge
                      v-if="!node.enabled"
                      color="grey"
                      label="Inactive"
                      class="q-mr-xs"
                    />
                    <q-badge
                      v-if="node.children && node.children.length > 0"
                      color="primary"
                      :label="node.children.length"
                    />
                  </div>
                </div>

                <!-- Children (recursive) -->
                <div v-if="isExpanded(node) && node.children && node.children.length > 0" class="org-tree-children">
                  <template v-for="child in node.children" :key="child.id">
                    <div class="org-tree-node">
                      <div
                        class="org-tree-item rounded-borders q-mb-xs"
                        :class="{ 'org-tree-item--selected': isSelected(child) }"
                      >
                        <div class="row items-center no-wrap q-pa-sm cursor-pointer" @click="selectOrganization(child)">
                          <q-btn
                            v-if="child.children && child.children.length > 0"
                            flat
                            round
                            dense
                            size="sm"
                            :icon="isExpanded(child) ? 'expand_more' : 'chevron_right'"
                            color="primary"
                            class="q-mr-xs"
                            @click.stop="toggleExpand(child)"
                          />
                          <div v-else class="q-mr-md" style="width: 24px;" />

                          <q-icon
                            :name="getOrganizationIcon(child.type)"
                            :color="getOrganizationColor(child.type)"
                            size="sm"
                            class="q-mr-sm"
                          />

                          <div class="col">
                            <div class="text-body2">{{ child.name }}</div>
                          </div>

                          <q-badge
                            v-if="!child.enabled"
                            color="grey"
                            label="Inactive"
                            class="q-mr-xs"
                          />
                          <q-badge
                            v-if="child.children && child.children.length > 0"
                            color="primary"
                            :label="child.children.length"
                          />
                        </div>
                      </div>

                      <!-- Level 3 children -->
                      <div v-if="isExpanded(child) && child.children && child.children.length > 0" class="org-tree-children">
                        <template v-for="grandchild in child.children" :key="grandchild.id">
                          <div class="org-tree-node">
                            <div
                              class="org-tree-item rounded-borders q-mb-xs"
                              :class="{ 'org-tree-item--selected': isSelected(grandchild) }"
                            >
                              <div class="row items-center no-wrap q-pa-sm cursor-pointer" @click="selectOrganization(grandchild)">
                                <q-btn
                                  v-if="grandchild.children && grandchild.children.length > 0"
                                  flat
                                  round
                                  dense
                                  size="sm"
                                  :icon="isExpanded(grandchild) ? 'expand_more' : 'chevron_right'"
                                  color="primary"
                                  class="q-mr-xs"
                                  @click.stop="toggleExpand(grandchild)"
                                />
                                <div v-else class="q-mr-md" style="width: 24px;" />

                                <q-icon
                                  :name="getOrganizationIcon(grandchild.type)"
                                  :color="getOrganizationColor(grandchild.type)"
                                  size="sm"
                                  class="q-mr-sm"
                                />

                                <div class="col">
                                  <div class="text-body2">{{ grandchild.name }}</div>
                                </div>

                                <q-badge
                                  v-if="!grandchild.enabled"
                                  color="grey"
                                  label="Inactive"
                                  class="q-mr-xs"
                                />
                                <q-badge
                                  v-if="grandchild.children && grandchild.children.length > 0"
                                  color="primary"
                                  :label="grandchild.children.length"
                                />
                              </div>
                            </div>

                            <!-- Level 4+ children -->
                            <div v-if="isExpanded(grandchild) && grandchild.children && grandchild.children.length > 0" class="org-tree-children">
                              <template v-for="greatgrandchild in grandchild.children" :key="greatgrandchild.id">
                                <div
                                  class="org-tree-item rounded-borders q-mb-xs"
                                  :class="{ 'org-tree-item--selected': isSelected(greatgrandchild) }"
                                >
                                  <div class="row items-center no-wrap q-pa-sm cursor-pointer" @click="selectOrganization(greatgrandchild)">
                                    <div class="q-mr-md" style="width: 24px;" />

                                    <q-icon
                                      :name="getOrganizationIcon(greatgrandchild.type)"
                                      :color="getOrganizationColor(greatgrandchild.type)"
                                      size="sm"
                                      class="q-mr-sm"
                                    />

                                    <div class="col">
                                      <div class="text-body2">{{ greatgrandchild.name }}</div>
                                    </div>

                                    <q-badge
                                      v-if="!greatgrandchild.enabled"
                                      color="grey"
                                      label="Inactive"
                                      class="q-mr-xs"
                                    />
                                  </div>
                                </div>
                              </template>
                            </div>
                          </div>
                        </template>
                      </div>
                    </div>
                  </template>
                </div>
              </div>
            </template>
          </div>
        </q-scroll-area>
      </q-card-section>

      <!-- Footer -->
      <q-separator />
      <q-card-actions class="row items-center q-px-md q-py-md">
        <div class="text-caption text-grey-7">
          <q-icon name="account_tree" size="xs" class="q-mr-xs" />
          {{ totalOrganizations }} {{ totalOrganizations === 1 ? 'organization' : 'organizations' }}
        </div>
        <q-space />
        <q-btn flat dense label="Cancel" color="grey-7" size="sm" class="rounded-borders" @click="handleCancel" />
      </q-card-actions>
    </q-card>
  </q-dialog>
</template>

<style scoped lang="scss">
.rounded-borders {
  border-radius: var(--mapex-radius-md);
}

.org-tree-children {
  padding-left: 20px;
  border-left: 1px solid var(--mapex-divider);
  margin-left: 12px;
}

.org-tree-item {
  border: 1px solid transparent;
  transition: all var(--mapex-transition-base) ease;

  &:hover {
    background-color: var(--mapex-surface-bg);
    border-color: var(--mapex-divider);
  }

  &--selected {
    background-color: var(--mapex-active-bg) !important;
    border-color: var(--q-primary) !important;
    border-left: 3px solid var(--q-primary);
  }
}

/* Better spacing for filter inputs */
:deep(.q-field--outlined .q-field__control) {
  border-radius: var(--mapex-radius-md);
}

/* Footer padding */
:deep(.q-card__actions) {
  padding: 16px !important;
}
</style>
