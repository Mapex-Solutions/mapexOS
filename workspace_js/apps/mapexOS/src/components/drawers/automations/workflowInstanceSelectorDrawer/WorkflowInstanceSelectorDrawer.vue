<script setup lang="ts">
defineOptions({ name: 'WorkflowInstanceSelectorDrawer' });

/** TYPE IMPORTS */
import type { InstanceResponse } from '@mapexos/schemas';
import type {
  WorkflowInstanceSelectorDrawerProps,
  WorkflowInstanceSelectorDrawerEmits,
} from './interfaces';

/** VUE IMPORTS */
import { ref, watch } from 'vue';

/** COMPONENTS */
import { GenericDrawer } from '@components/drawers/common/genericDrawer';
import { DetailChip } from '@components/chips';

/** SERVICES */
import { apis } from '@services/mapex';

/** PROPS & EMITS */
const props = defineProps<WorkflowInstanceSelectorDrawerProps>();
const emit = defineEmits<WorkflowInstanceSelectorDrawerEmits>();

/** STATE */
const instances = ref<InstanceResponse[]>([]);
const isLoading = ref(false);
const searchQuery = ref('');

/** FUNCTIONS */

/**
 * Fetch workflow instances from API
 * @returns {Promise<void>}
 */
async function fetchInstances(): Promise<void> {
  isLoading.value = true;
  try {
    const queryParams: Record<string, any> = {
      perPage: 100,
      enabled: true,
    };
    if (searchQuery.value) queryParams.name = searchQuery.value;

    const response = await apis.workflows.instance.list(queryParams);
    instances.value = response.items || [];
  } catch {
    instances.value = [];
  } finally {
    isLoading.value = false;
  }
}

/**
 * Handle instance selection (single-select)
 * @param {InstanceResponse} instance - Selected instance
 * @returns {void}
 */
function selectInstance(instance: InstanceResponse): void {
  emit('select', instance);
  emit('update:modelValue', false);
}

/**
 * Check if an instance is the currently selected one
 * @param {InstanceResponse} instance - Instance to check
 * @returns {boolean} True if selected
 */
function isSelected(instance: InstanceResponse): boolean {
  return instance._id === props.selectedInstanceId;
}

/**
 * Handle filter change — re-fetch
 * @returns {void}
 */
function onFilterChange(): void {
  void fetchInstances();
}

/**
 * Handle cancel action
 * @returns {void}
 */
function handleCancel(): void {
  emit('cancel');
  emit('update:modelValue', false);
}

/** WATCHERS */
watch(() => props.modelValue, (isOpen) => {
  if (isOpen) {
    void fetchInstances();
  }
});
</script>

<template>
  <GenericDrawer
    :model-value="modelValue"
    title="Select Workflow Instance"
    icon="play_circle"
    icon-color="teal-7"
    :width="500"
    @update:model-value="emit('update:modelValue', $event)"
    @close="handleCancel"
  >
    <!-- Filters -->
    <div class="q-mb-md">
      <q-input
        v-model="searchQuery"
        outlined
        dense
        clearable
        placeholder="Search instances..."
        @update:model-value="onFilterChange"
      >
        <template #prepend>
          <q-icon name="search" />
        </template>
      </q-input>
    </div>

    <q-separator class="q-mb-md" />

    <!-- Loading -->
    <div v-if="isLoading" class="q-pa-md text-center">
      <q-spinner color="primary" size="3em" />
      <div class="loading-text q-mt-md">Loading instances...</div>
    </div>

    <!-- Empty -->
    <div v-else-if="instances.length === 0" class="q-pa-md text-center">
      <q-icon name="inbox" size="4em" class="empty-icon" />
      <div class="loading-text q-mt-md">No instances found</div>
    </div>

    <!-- List -->
    <q-list v-else separator class="instance-list">
      <q-item
        v-for="inst in instances"
        :key="inst._id || ''"
        clickable
        :active="isSelected(inst)"
        @click="selectInstance(inst)"
      >
        <q-item-section avatar>
          <q-avatar
            :color="isSelected(inst) ? 'primary' : (inst.enabled ? 'teal-7' : 'grey-5')"
            icon="play_circle"
            text-color="white"
            size="md"
          />
        </q-item-section>

        <q-item-section>
          <q-item-label>{{ inst.name }}</q-item-label>
          <q-item-label caption class="item-caption">
            {{ inst.definitionName || '—' }}
            <span v-if="inst.uniqueExecution"> · UNIQUE</span>
          </q-item-label>
        </q-item-section>

        <q-item-section side>
          <div class="row q-gutter-xs items-center">
            <DetailChip
              v-if="inst.uniqueExecution"
              dense
              size="xs"
              color="warning"
              label="U"
            />
            <q-icon v-if="isSelected(inst)" name="check_circle" color="primary" />
          </div>
        </q-item-section>
      </q-item>
    </q-list>

    <!-- Footer -->
    <template #footer>
      <div class="footer-count">
        <q-icon name="play_circle" size="xs" class="q-mr-xs" />
        {{ instances.length }} {{ instances.length === 1 ? 'instance' : 'instances' }}
      </div>
      <q-space />
      <q-btn flat dense no-caps size="sm" label="Cancel" class="cancel-btn" @click="handleCancel" />
    </template>
  </GenericDrawer>
</template>

<style lang="scss" scoped>
.loading-text {
  color: var(--mapex-text-secondary);
}

.empty-icon {
  color: var(--mapex-text-muted);
}

.item-caption {
  color: var(--mapex-text-secondary);
}

.footer-count {
  font-size: 0.75rem;
  color: var(--mapex-text-secondary);
  display: flex;
  align-items: center;
}

.cancel-btn {
  color: var(--mapex-text-secondary);
}

.instance-list {
  :deep(.q-item) {
    transition: var(--mapex-transition-base);
  }

  :deep(.q-item:hover) {
    background-color: var(--mapex-hover-overlay);
  }

  :deep(.q-item.q-item--active) {
    background-color: rgba(var(--q-primary-rgb), 0.08) !important;
    border-left: 3px solid var(--q-primary);
  }
}

:deep(.q-field--outlined .q-field__control) {
  border-radius: var(--mapex-radius-md);
}
</style>
