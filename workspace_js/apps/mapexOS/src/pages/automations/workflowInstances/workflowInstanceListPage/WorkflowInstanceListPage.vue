<script setup lang="ts">
defineOptions({
  name: 'WorkflowInstanceListPage',
});

/** TYPE IMPORTS */
import type { ListHeaderMenuColumn } from '@components/headers';
import type { DataRowColumn } from '@components/cards';
import type {
  WorkflowInstanceListPageFilters,
  WorkflowInstanceListPageColumnVisibility,
  WorkflowInstanceListItem,
} from './interfaces/workflowInstanceListPage.interface';

/** VUE IMPORTS */
import { ref, computed, onMounted } from 'vue';
import { useRouter } from 'vue-router';

/** COMPONENTS */
import { PageHeader, ListHeaderMenu } from '@components/headers';
import { ListCardEmpty, DataRow } from '@components/cards';
import { ListPagination } from '@components/navigation';
import { AppTooltip } from '@components/tooltips';
import { WorkflowInstanceDetailsDrawer } from '@components/drawers/automations/workflowInstanceDetailsDrawer';
import { ExecuteResultDialog } from './components';

/** COMPOSABLES */
import { useWorkflowInstanceListPageTranslations } from '@composables/i18n';
import { useOrgChangeRefresh } from '@composables/organizations';
import { usePermissions } from '@composables/shared/usePermissions';

/** UTILS */
import { notifySuccess, dialogDelete } from '@utils/alert';
import { handleApiError } from '@utils/error';

/** SERVICES */
import { apis } from '@services/mapex';

/** LOCAL IMPORTS */
import {
  DEFAULT_ITEMS_PER_PAGE,
  INSTANCE_COLUMN_VISIBILITY_DEFAULTS,
  INSTANCE_FILTER_DEFAULTS,
  LIST_PROJECTION,
} from './constants';

/** COMPOSABLES & STORES */
const router = useRouter();
const t = useWorkflowInstanceListPageTranslations();
const { canCreate, canUpdate, canDelete, canRead } = usePermissions();
const canCreateInstance = canCreate('workflows.instances');
const canUpdateInstance = canUpdate('workflows.instances');
const canDeleteInstance = canDelete('workflows.instances');
const canReadInstance = canRead('workflows.instances');

/** STATE */
const instancesList = ref<WorkflowInstanceListItem[]>([]);
const loading = ref(false);
const lastUpdatedAt = ref<number | undefined>(undefined);
const showDetailsDrawer = ref(false);
const selectedInstanceId = ref<string | null>(null);
const itemsPerPage = ref(DEFAULT_ITEMS_PER_PAGE);
const currentPage = ref(1);
const totalPages = ref(1);
const totalItems = ref(0);
const filters = ref<WorkflowInstanceListPageFilters>({ ...INSTANCE_FILTER_DEFAULTS });
const columnVisibilityState = ref<WorkflowInstanceListPageColumnVisibility>({ ...INSTANCE_COLUMN_VISIBILITY_DEFAULTS });

// Quick filter state (inline)
const quickSearch = ref('');
const quickStatus = ref<boolean | null>(null);

// Advanced filter state (drawer)
const advancedFilterValues = ref<Record<string, any>>({
  uniqueExecution: null,
});

/** COMPUTED */

/**
 * Status options for quick filter select
 */
const statusOptions = computed(() => [
  { label: t.filters.allStatus.value, value: null },
  { label: t.filters.options.enabled.value, value: true },
  { label: t.filters.options.disabled.value, value: false },
]);

/**
 * Check if any filters are active
 */
const hasActiveFilters = computed(() => {
  return !!(
    quickSearch.value ||
    quickStatus.value !== null ||
    advancedFilterValues.value.uniqueExecution !== null
  );
});

/**
 * Active filter chips for visual feedback
 */
const activeFilterChips = computed(() => {
  const chips: Array<{ key: string; label: string; value: string }> = [];

  if (quickSearch.value) {
    chips.push({
      key: 'name',
      label: t.filters.instanceName.value,
      value: quickSearch.value,
    });
  }

  if (quickStatus.value !== null) {
    chips.push({
      key: 'status',
      label: t.filters.status.value,
      value: quickStatus.value ? t.filters.options.enabled.value : t.filters.options.disabled.value,
    });
  }

  if (advancedFilterValues.value.uniqueExecution !== null) {
    chips.push({
      key: 'uniqueExecution',
      label: t.filters.uniqueExecution.value,
      value: advancedFilterValues.value.uniqueExecution ? t.filters.options.yes.value : t.filters.options.no.value,
    });
  }

  return chips;
});

/** Maximum number of visible filter chips */
const MAX_VISIBLE_CHIPS = 2;

const visibleFilterChips = computed(() => activeFilterChips.value.slice(0, MAX_VISIBLE_CHIPS));
const hiddenFilterChips = computed(() => activeFilterChips.value.slice(MAX_VISIBLE_CHIPS));
const hiddenFiltersCount = computed(() => hiddenFilterChips.value.length);

/**
 * Column visibility using ListHeaderMenuColumn format
 */
const menuColumns = computed((): ListHeaderMenuColumn[] => [
  { key: 'definitionName', label: t.menuColumns.definitionName.value, visible: columnVisibilityState.value.definitionName },
  { key: 'inputsCount', label: t.menuColumns.inputsCount.value, visible: columnVisibilityState.value.inputsCount },
  { key: 'uniqueExecution', label: t.menuColumns.uniqueExecution.value, visible: columnVisibilityState.value.uniqueExecution },
]);

/**
 * Filtered columns based on visibility state
 */
const visibleColumns = computed(() => {
  return t.columns.value.filter((col: DataRowColumn) => {
    if (col.key === 'icon' || col.key === 'name') return true;
    if (col.key === 'definitionName') return columnVisibilityState.value.definitionName;
    if (col.key === 'inputsCount') return columnVisibilityState.value.inputsCount;
    if (col.key === 'uniqueExecution') return columnVisibilityState.value.uniqueExecution;
    return true;
  });
});

/** FUNCTIONS */

/**
 * Apply quick filters and fetch data
 * @returns {void}
 */
function applyQuickFilters(): void {
  filters.value.name = quickSearch.value || undefined;
  filters.value.status = quickStatus.value ?? undefined;
  currentPage.value = 1;
  void fetchInstances();
}

/**
 * Remove individual filter
 * @param {string} key - Filter key to remove
 * @returns {void}
 */
function removeFilter(key: string): void {
  if (key === 'name') {
    quickSearch.value = '';
    filters.value.name = undefined;
  } else if (key === 'status') {
    quickStatus.value = null;
    filters.value.status = undefined;
  } else if (key === 'uniqueExecution') {
    advancedFilterValues.value.uniqueExecution = null;
  }

  currentPage.value = 1;
  void fetchInstances();
}

/**
 * Clear all filters
 * @returns {void}
 */
function clearAllFilters(): void {
  quickSearch.value = '';
  quickStatus.value = null;
  advancedFilterValues.value = { uniqueExecution: null };
  filters.value = { ...INSTANCE_FILTER_DEFAULTS };
  currentPage.value = 1;
  void fetchInstances();
}

/**
 * Fetch workflow instances from API with server-side filtering and pagination
 * @returns {Promise<void>}
 */
async function fetchInstances(): Promise<void> {
  try {
    loading.value = true;

    const queryParams: Record<string, any> = {
      page: currentPage.value,
      perPage: itemsPerPage.value,
      projection: LIST_PROJECTION,
    };

    // Quick filters
    if (filters.value.name) queryParams.name = filters.value.name;
    if (typeof filters.value.status === 'boolean') queryParams.enabled = filters.value.status;

    // Advanced filters
    if (advancedFilterValues.value.uniqueExecution !== null) queryParams.uniqueExecution = advancedFilterValues.value.uniqueExecution;

    const response = await apis.workflows.instance.list(queryParams);

    // Map API response to WorkflowInstanceListItem
    instancesList.value = (response.items || []).map((inst: any): WorkflowInstanceListItem => ({
      id: inst._id || '',
      name: inst.name || '',
      description: inst.description || '',
      enabled: inst.enabled ?? true,
      definitionName: inst.definitionName || '—',
      inputsCount: inst.externalInputs ? Object.keys(inst.externalInputs).length : 0,
      uniqueExecution: inst.uniqueExecution ?? false,
      workflowUUID: inst.workflowUUID || '',
    }));

    // Update pagination from response
    if (response.pagination) {
      totalItems.value = response.pagination.totalItems || 0;
      totalPages.value = response.pagination.totalPages || 1;
    } else {
      totalItems.value = instancesList.value.length;
      totalPages.value = 1;
    }
  } catch (error: any) {
    handleApiError(error, { defaultMessage: 'Failed to load instances' });
    instancesList.value = [];
  } finally {
    loading.value = false;
    lastUpdatedAt.value = Date.now();
  }
}

/**
 * Handle pagination navigation
 * @param {number} page - New page number
 * @returns {void}
 */
function handlePageChange(page: number): void {
  currentPage.value = page;
  void fetchInstances();
}

/**
 * Handle items per page change
 * @param {number} newValue - New items per page value
 * @returns {void}
 */
function handleItemsPerPageChange(newValue: number): void {
  itemsPerPage.value = newValue;
  currentPage.value = 1;
  void fetchInstances();
}

/**
 * Update menu columns when changed
 * @param {ListHeaderMenuColumn[]} columns - Updated column visibility states
 * @returns {void}
 */
function handleColumnsUpdate(columns: ListHeaderMenuColumn[]): void {
  columns.forEach(col => {
    if (col.key === 'definitionName') columnVisibilityState.value.definitionName = col.visible;
    if (col.key === 'inputsCount') columnVisibilityState.value.inputsCount = col.visible;
    if (col.key === 'uniqueExecution') columnVisibilityState.value.uniqueExecution = col.visible;
  });
}

/**
 * View instance details in drawer
 * @param {WorkflowInstanceListItem} instance - Instance to view
 * @returns {void}
 */
function viewDetails(instance: WorkflowInstanceListItem): void {
  if (!canReadInstance.value) return;
  selectedInstanceId.value = instance.id;
  showDetailsDrawer.value = true;
}

/**
 * Handle edit from drawer
 * @param {string} instanceId - Instance ID to edit
 * @returns {void}
 */
function handleDrawerEdit(instanceId: string): void {
  if (!canUpdateInstance.value) return;
  void router.push(`/workflow_instances/${instanceId}`);
}

/**
 * Navigate to instance editor
 * @param {WorkflowInstanceListItem} instance - Instance to edit
 * @returns {void}
 */
function editInstance(instance: WorkflowInstanceListItem): void {
  if (!canUpdateInstance.value) return;
  void router.push(`/workflow_instances/${instance.id}`);
}

/**
 * Confirm delete operation with dialog
 * @param {WorkflowInstanceListItem} instance - Instance to delete
 * @returns {Promise<void>}
 */
async function confirmDelete(instance: WorkflowInstanceListItem): Promise<void> {
  if (!canDeleteInstance.value) return;
  const instanceName = instance.name || 'this instance';
  const confirmed = await dialogDelete({
    title: t.dialog.deleteTitle.value,
    message: t.messages.confirmDelete(instanceName),
  });

  if (confirmed) {
    await deleteInstance(instance);
  }
}

/**
 * Delete instance via API
 * @param {WorkflowInstanceListItem} instance - Instance to delete
 * @returns {Promise<void>}
 */
async function deleteInstance(instance: WorkflowInstanceListItem): Promise<void> {
  try {
    await apis.workflows.instance.delete({ instanceId: instance.id });
    notifySuccess({ message: t.messages.deletedSuccessfully.value });

    if (instancesList.value.length <= 1 && currentPage.value > 1) {
      currentPage.value -= 1;
    }
    await fetchInstances();
  } catch (error: any) {
    handleApiError(error, { defaultMessage: 'Failed to delete instance' });
  }
}

/** EXECUTE STATE */
const executeLoading = ref(false);
const executeDialogOpen = ref(false);
const executeResult = ref<{ workflowUUID: string; status: string; errorInfo?: any } | null>(null);
const executeError = ref<string | null>(null);

/**
 * Handle custom DataRow actions
 * @param {string} actionKey - The action key
 * @param {WorkflowInstanceListItem} row - The instance row data
 */
function handleCustomAction(actionKey: string, row: WorkflowInstanceListItem): void {
  if (actionKey === 'execute') {
    void executeInstance(row.id);
  }
}

/**
 * Execute a workflow instance via HTTP API
 * @param {string} instanceId - The instance ID to execute
 * @returns {Promise<void>}
 */
async function executeInstance(instanceId: string): Promise<void> {
  executeLoading.value = true;
  executeError.value = null;
  executeResult.value = null;
  executeDialogOpen.value = true;

  try {
    const result = await apis.workflows.instance.execute(
      { instanceId },
      {},
    );
    executeResult.value = result;
  } catch (error: any) {
    executeError.value = error?.response?.data?.message || error?.message || 'Execution failed';
  } finally {
    executeLoading.value = false;
  }
}

/**
 * Navigate to instance creation page
 * @returns {void}
 */
function navigateToCreate(): void {
  void router.push('/workflow_instances/add');
}

/** LIFECYCLE HOOKS */
onMounted(async () => {
  await fetchInstances();
});

/**
 * Auto-refresh when organization changes
 * Resets pagination, filters and refetches data
 */
useOrgChangeRefresh(async () => {
  currentPage.value = 1;
  quickSearch.value = '';
  quickStatus.value = null;
  advancedFilterValues.value = { uniqueExecution: null };
  filters.value = { ...INSTANCE_FILTER_DEFAULTS };
  await fetchInstances();
});
</script>

<template>
  <q-page class="q-pt-lg">

    <!-- Header Section -->
    <PageHeader
      icon="play_circle"
      icon-color="teal-7"
      :title="t.page.title.value"
      :description="t.page.description.value"
      :button="canCreateInstance ? { label: t.page.addButton.value, icon: 'add', onClick: navigateToCreate, color: 'teal-7' } : undefined"
      :info="t.page.info.value"
    />

    <!-- Filters Section -->
    <div class="text-caption text-grey-7 q-mb-xs">{{ t.filters.label.value }}</div>
    <div class="row items-center q-col-gutter-sm q-mb-md">
      <!-- Search Input -->
      <div class="col">
        <q-input
          v-model="quickSearch"
          outlined
          dense
          clearable
          :placeholder="t.filters.searchPlaceholder.value"
          class="filter-input"
          @keyup.enter="applyQuickFilters"
          @clear="applyQuickFilters"
        >
          <template #prepend>
            <q-icon name="search" color="grey-6" />
          </template>
        </q-input>
      </div>

      <!-- Status Select -->
      <div class="col-auto" style="min-width: 140px;">
        <q-select
          v-model="quickStatus"
          outlined
          dense
          emit-value
          map-options
          :options="statusOptions"
          :label="t.filters.status.value"
          class="filter-input"
          @update:model-value="applyQuickFilters"
        />
      </div>

      <!-- Workflow Select -->
    </div>

    <!-- Active Filter Chips (with limit) -->
    <div v-if="hasActiveFilters" class="row items-center q-mb-md q-gutter-xs">
      <!-- Visible Chips (max 2) -->
      <q-chip
        v-for="chip in visibleFilterChips"
        :key="chip.key"
        removable
        dense
        outline
        color="teal-7"
        size="sm"
        @remove="removeFilter(chip.key)"
      >
        <span class="text-weight-medium">{{ chip.label }}:</span>&nbsp;{{ chip.value }}
      </q-chip>

      <!-- +N Badge for hidden filters -->
      <q-badge
        v-if="hiddenFiltersCount > 0"
        color="teal-7"
        outline
        class="q-pa-xs cursor-pointer"
      >
        +{{ hiddenFiltersCount }}
        <AppTooltip>
          <div v-for="chip in hiddenFilterChips" :key="chip.key" class="q-mb-xs">
            <span class="text-weight-medium">{{ chip.label }}:</span> {{ chip.value }}
          </div>
        </AppTooltip>
      </q-badge>

      <!-- Clear All Button -->
      <q-btn
        flat
        dense
        size="sm"
        color="grey-7"
        icon="filter_alt_off"
        :label="t.filters.clearAll.value"
        no-caps
        class="q-ml-sm"
        @click="clearAllFilters"
      />
    </div>

    <!-- Results Section -->
    <div class="row items-center q-pt-xl q-mb-md">
      <div class="col">
        <div class="row items-center">
          <q-icon name="play_circle" size="sm" color="teal-7" class="q-mr-sm" />
          <div class="text-subtitle1 text-weight-medium" style="color: var(--mapex-text-primary);">
            {{ t.page.listTitle.value }}
          </div>
        </div>
      </div>
      <div class="col-auto">
        <ListHeaderMenu
          icon="play_circle"
          :items-count="totalItems"
          :item-label="t.page.itemLabel.value"
          :item-label-plural="t.page.itemLabelPlural.value"
          :items-per-page="itemsPerPage"
          :columns="menuColumns"
          :filtered="hasActiveFilters"
          :refreshing="loading"
          :last-updated-at="lastUpdatedAt"
          @update:items-per-page="handleItemsPerPageChange"
          @update:columns="handleColumnsUpdate"
          @refresh="fetchInstances"
        />
      </div>
    </div>

    <!-- Loading Spinner -->
    <div v-if="loading" class="row justify-center q-my-lg">
      <q-spinner color="teal-7" size="3em" />
    </div>

    <!-- Instance Row List -->
    <div v-else class="row">
      <div
        v-for="(instance, index) in instancesList"
        :key="instance.id || `instance-${index}`"
        class="col-12 q-mb-xs"
      >
        <DataRow
          :data="instance"
          :columns="visibleColumns"
          :actions="{
            showEdit: canUpdateInstance,
            showView: canReadInstance,
            showDelete: canDeleteInstance,
            customActions: [
              {
                key: 'execute',
                icon: 'play_arrow',
                label: 'Run',
                color: 'green-7',
                condition: (row: any) => row.enabled,
              }
            ]
          }"
          @click="viewDetails"
          @dblclick="editInstance"
          @edit="editInstance"
          @view="viewDetails"
          @delete="confirmDelete"
          @action="handleCustomAction"
        />
      </div>

      <!-- No Results -->
      <div v-if="instancesList.length === 0" class="col-12">
        <ListCardEmpty
          :title="t.empty.title.value"
          :description="t.empty.description.value"
          icon="play_circle"
        />
      </div>
    </div>

    <!-- Pagination -->
    <ListPagination
      v-model="currentPage"
      :total-pages="totalPages"
      @change="handlePageChange"
    />

    <!-- Instance Details Drawer -->
    <WorkflowInstanceDetailsDrawer
      v-model="showDetailsDrawer"
      :instance-id="selectedInstanceId"
      @edit="handleDrawerEdit"
    />

    <!-- Execute Result Dialog -->
    <ExecuteResultDialog
      v-model="executeDialogOpen"
      :loading="executeLoading"
      :result="executeResult"
      :error="executeError"
    />
  </q-page>
</template>

<style lang="scss" scoped>
.filter-input {
  :deep(.q-field__control) {
    border-radius: var(--mapex-radius-md);
  }
}
</style>
