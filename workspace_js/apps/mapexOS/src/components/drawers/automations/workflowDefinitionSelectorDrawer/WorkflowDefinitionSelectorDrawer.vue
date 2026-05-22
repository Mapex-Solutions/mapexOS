<script setup lang="ts">
defineOptions({ name: 'WorkflowDefinitionSelectorDrawer' });

/** TYPE IMPORTS */
import type { DefinitionResponse } from '@mapexos/schemas';
import type {
  WorkflowDefinitionSelectorDrawerProps,
  WorkflowDefinitionSelectorDrawerEmits,
} from './interfaces';

/** VUE IMPORTS */
import { ref, computed, watch } from 'vue';

/** COMPONENTS */
import { GenericDrawer } from '@components/drawers/common/genericDrawer';

/** COMPOSABLES */
import { useCreateEditWorkflowInstanceTranslations } from '@src/composables/i18n/pages/automations/workflowInstances/createEditWorkflowInstancePage/useCreateEditWorkflowInstanceTranslations';

/** SERVICES */
import { apis } from '@services/mapex';

/** PROPS & EMITS */
const props = defineProps<WorkflowDefinitionSelectorDrawerProps>();
const emit = defineEmits<WorkflowDefinitionSelectorDrawerEmits>();

/** COMPOSABLES & STORES */
const t = useCreateEditWorkflowInstanceTranslations();

/** STATE */
const definitions = ref<DefinitionResponse[]>([]);
const isLoading = ref(false);
const searchQuery = ref('');
const statusFilter = ref<boolean | undefined>(undefined);

/** COMPUTED */

/**
 * Status options for filter select
 * @returns {Array} Status filter options
 */
const statusOptions = computed(() => [
  { label: t.fields.allStatus.value, value: undefined },
  { label: t.fields.active.value, value: true },
  { label: t.fields.inactive.value, value: false },
]);

/**
 * Footer count label
 * @returns {string} Count text
 */
const countLabel = computed(() => {
  const count = definitions.value.length;
  return `${count} ${count === 1 ? 'definition' : 'definitions'}`;
});

/** FUNCTIONS */

/**
 * Fetch workflow definitions from API
 * @returns {Promise<void>}
 */
async function fetchDefinitions(): Promise<void> {
  isLoading.value = true;
  try {
    const queryParams: Record<string, any> = { perPage: 100 };
    if (searchQuery.value) queryParams.name = searchQuery.value;
    if (typeof statusFilter.value === 'boolean') queryParams.enabled = statusFilter.value;

    const response = await apis.workflows.definition.list(queryParams);
    definitions.value = response.items || [];
  } catch {
    definitions.value = [];
  } finally {
    isLoading.value = false;
  }
}

/**
 * Handle definition selection (single-select)
 * @param {DefinitionResponse} definition - Selected definition
 * @returns {void}
 */
function selectDefinition(definition: DefinitionResponse): void {
  emit('select', definition);
  emit('update:modelValue', false);
}

/**
 * Check if a definition is the currently selected one
 * @param {DefinitionResponse} definition - Definition to check
 * @returns {boolean} True if selected
 */
function isSelected(definition: DefinitionResponse): boolean {
  return definition._id === props.selectedDefinitionId;
}

/**
 * Handle filter change — re-fetch definitions
 * @returns {void}
 */
function onFilterChange(): void {
  void fetchDefinitions();
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
    void fetchDefinitions();
  }
});
</script>

<template>
  <GenericDrawer
    :model-value="modelValue"
    :title="t.fields.definitionDrawerTitle.value"
    icon="account_tree"
    :width="500"
    @update:model-value="emit('update:modelValue', $event)"
    @close="handleCancel"
  >
    <!-- Info Banner -->
    <q-banner dense class="info-banner rounded-borders q-mb-md">
      <template #avatar>
        <q-icon name="info" size="sm" />
      </template>
      <div class="text-caption">
        {{ t.fields.definitionPlaceholder.value }}
      </div>
    </q-banner>

    <!-- Filters -->
    <div class="q-mb-md">
      <div class="section-label q-mb-md">
        <q-icon name="filter_list" size="xs" class="q-mr-xs" />
        {{ t.fields.definitionDrawerSearch.value }}
      </div>
      <div class="row q-col-gutter-md">
        <!-- Search -->
        <div class="col-12 col-sm-8">
          <q-input
            v-model="searchQuery"
            outlined
            dense
            clearable
            :label="t.fields.definitionDrawerSearch.value"
            @update:model-value="onFilterChange"
          >
            <template #prepend>
              <q-icon name="search" />
            </template>
          </q-input>
        </div>

        <!-- Status -->
        <div class="col-12 col-sm-4">
          <q-select
            v-model="statusFilter"
            outlined
            dense
            :label="t.fields.enabled.value"
            :options="statusOptions"
            emit-value
            map-options
            @update:model-value="onFilterChange"
          >
            <template #prepend>
              <q-icon name="toggle_on" />
            </template>
          </q-select>
        </div>
      </div>
    </div>

    <q-separator class="q-mb-md" />

    <!-- Results Header -->
    <div class="section-label q-mb-sm">
      <q-icon name="account_tree" size="xs" class="q-mr-xs" />
      {{ t.fields.definition.value }}
    </div>

    <!-- Loading -->
    <div v-if="isLoading" class="q-pa-md text-center">
      <q-spinner color="primary" size="3em" />
      <div class="loading-text q-mt-md">{{ t.notifications.loading.value }}</div>
    </div>

    <!-- Empty -->
    <div v-else-if="definitions.length === 0" class="q-pa-md text-center">
      <q-icon name="inbox" size="4em" class="empty-icon" />
      <div class="loading-text q-mt-md">{{ t.fields.noExternalInputs.value }}</div>
    </div>

    <!-- List -->
    <q-list v-else separator class="definition-list">
      <q-item
        v-for="def in definitions"
        :key="def._id || ''"
        clickable
        :active="isSelected(def)"
        @click="selectDefinition(def)"
      >
        <q-item-section avatar>
          <q-avatar
            :color="isSelected(def) ? 'primary' : (def.enabled ? 'primary' : 'grey-5')"
            icon="account_tree"
            text-color="white"
            size="md"
          />
        </q-item-section>

        <q-item-section>
          <q-item-label>{{ def.name }}</q-item-label>
          <q-item-label caption class="item-caption">
            v{{ def.definitionVersion }}
            <span v-if="def.status"> · {{ def.status }}</span>
            <span v-if="(def as any).externalInputs?.length"> · {{ (def as any).externalInputs.length }} inputs</span>
          </q-item-label>
        </q-item-section>

        <q-item-section side>
          <q-icon v-if="isSelected(def)" name="check_circle" color="primary" />
        </q-item-section>
      </q-item>
    </q-list>

    <!-- Footer -->
    <template #footer>
      <div class="footer-count">
        <q-icon name="account_tree" size="xs" class="q-mr-xs" />
        {{ countLabel }}
      </div>
      <q-space />
      <q-btn
        flat
        dense
        no-caps
        size="sm"
        :label="t.navigation.cancel.value"
        class="cancel-btn"
        @click="handleCancel"
      />
    </template>
  </GenericDrawer>
</template>

<style lang="scss" scoped>
.info-banner {
  background: rgba(var(--q-primary-rgb), 0.08);
  color: var(--mapex-text-primary);
}

.section-label {
  font-size: 0.7rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  color: var(--mapex-text-secondary);
  display: flex;
  align-items: center;
}

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

.definition-list {
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
