<script setup lang="ts">
defineOptions({
  name: 'WorkflowSelectorDialog'
});

/** TYPE IMPORTS */
import type { WorkflowSelectorDialogProps, WorkflowSelectorDialogEmits, WorkflowSelectorItem } from './interfaces';

/** VUE IMPORTS */
import { ref, computed, watch } from 'vue';

/** COMPONENTS */
import { GenericSelectorDialog } from '@components/dialogs/common/genericSelectorDialog';

/** COMPOSABLES */
import { useTS } from '@utils/translation';

/** SERVICES */
import { apis } from '@services/mapex';

/** PROPS & EMITS */
const props = withDefaults(defineProps<WorkflowSelectorDialogProps>(), {
  selectedWorkflowId: null,
  excludeWorkflowId: null,
});

const emit = defineEmits<WorkflowSelectorDialogEmits>();

/** COMPOSABLES & STORES */
const tsTitle = useTS({ titleCase: true });
const ts = useTS({ capitalize: true });
const tsRaw = useTS({ capitalize: false });
const bp = 'components.dialogs.workflowSelector';

/** STATE */

/**
 * All workflows fetched from API
 */
const workflows = ref<WorkflowSelectorItem[]>([]);

/**
 * Loading state
 */
const loading = ref(false);

/**
 * Client-side search query (filtering done locally)
 */
const searchQuery = ref('');

/**
 * Status filter: undefined = all, true = enabled, false = disabled
 */
const statusFilter = ref<boolean | undefined>(undefined);

/** COMPUTED */

/**
 * Status filter options for dropdown
 */
const statusOptions = computed(() => [
  { label: tsTitle(`${bp}.allFilter`), value: undefined },
  { label: tsTitle(`${bp}.enabled`), value: true },
  { label: tsTitle(`${bp}.disabled`), value: false },
]);

/**
 * Filtered workflows (client-side search + status + exclude self)
 */
const filteredWorkflows = computed(() => {
  let result = workflows.value;

  // Exclude self-reference
  if (props.excludeWorkflowId) {
    result = result.filter(w => w.id !== props.excludeWorkflowId);
  }

  // Filter by status
  if (statusFilter.value !== undefined) {
    result = result.filter(w => w.enabled === statusFilter.value);
  }

  // Filter by search query
  if (searchQuery.value.trim()) {
    const q = searchQuery.value.toLowerCase();
    result = result.filter(w =>
      w.name.toLowerCase().includes(q) ||
      w.description?.toLowerCase().includes(q),
    );
  }

  return result;
});

/**
 * Pre-selected workflow IDs for highlighting
 */
const selectedIds = computed(() =>
  props.selectedWorkflowId ? [props.selectedWorkflowId] : [],
);

/** WATCHERS */

/**
 * Fetch workflows when dialog opens
 */
watch(() => props.modelValue, (isOpen) => {
  if (isOpen) {
    void fetchWorkflows();
  }
});

/** FUNCTIONS */

/**
 * Fetch workflow definitions from API
 *
 * @returns {Promise<void>}
 */
async function fetchWorkflows(): Promise<void> {
  loading.value = true;
  try {
    const response = await apis.workflows.definition.list({
      perPage: 100,
      projection: '_id,name,description,enabled',
    });

    workflows.value = (response.items || []).map((def: any): WorkflowSelectorItem => ({
      id: def._id || '',
      name: def.name || '',
      description: def.description || '',
      enabled: def.enabled ?? true,
    }));
  } catch {
    workflows.value = [];
  } finally {
    loading.value = false;
  }
}

/**
 * Handle workflow selection from GenericSelectorDialog
 *
 * @param {any[]} items - Selected items (single-select, array of 1)
 */
function handleSelect(items: any[]): void {
  const workflow = items[0] as WorkflowSelectorItem;
  if (workflow) {
    emit('select', workflow);
  }
}

/**
 * Handle search query (client-side filtering)
 *
 * @param {string} query - Search query
 */
function handleSearch(query: string): void {
  searchQuery.value = query;
}
</script>

<template>
  <GenericSelectorDialog
    :model-value="modelValue"
    :title="tsTitle(`${bp}.title`)"
    icon="hub"
    icon-color="deep-purple-6"
    :items="filteredWorkflows"
    item-key="id"
    :multi-select="false"
    :selected-ids="selectedIds"
    :loading="loading"
    :search-placeholder="tsRaw(`${bp}.searchPlaceholder`)"
    :info-banner="{ text: ts(`${bp}.infoBanner`) }"
    :empty-text="ts(`${bp}.emptyText`)"
    empty-icon="inbox"
    results-icon="hub"
    footer-icon="hub"
    :item-noun-singular="tsRaw(`${bp}.itemSingular`)"
    :item-noun-plural="tsRaw(`${bp}.itemPlural`)"
    :total-items="filteredWorkflows.length"
    :active-item-style="{ backgroundColor: 'rgba(103, 58, 183, 0.08)', borderColor: 'var(--q-deep-purple-6)' }"
    @update:model-value="emit('update:modelValue', $event)"
    @select="handleSelect"
    @cancel="emit('cancel')"
    @search="handleSearch"
  >
    <!-- Status filter -->
    <template #filters>
      <div class="col-12">
        <q-select
          v-model="statusFilter"
          outlined
          dense
          clearable
          :label="tsTitle(`${bp}.status`)"
          :options="statusOptions"
          option-label="label"
          option-value="value"
          emit-value
          map-options
        >
          <template #prepend>
            <q-icon name="toggle_on" />
          </template>
        </q-select>
      </div>
    </template>

    <!-- Item rendering -->
    <template #item="{ item }">
      <q-item-section avatar>
        <q-avatar
          icon="hub"
          color="deep-purple-6"
          text-color="white"
          size="md"
        />
      </q-item-section>
      <q-item-section>
        <q-item-label class="text-weight-medium">{{ item.name }}</q-item-label>
        <q-item-label caption lines="2">{{ item.description }}</q-item-label>
        <q-item-label caption class="q-mt-xs">
          <q-badge
            :color="item.enabled ? 'positive' : 'grey-6'"
            :label="item.enabled ? tsTitle(`${bp}.enabled`) : tsTitle(`${bp}.disabled`)"
            dense
          />
        </q-item-label>
      </q-item-section>
    </template>
  </GenericSelectorDialog>
</template>
