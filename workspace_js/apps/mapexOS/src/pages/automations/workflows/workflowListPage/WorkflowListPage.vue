<script setup lang="ts">
defineOptions({
  name: 'WorkflowListPage',
});

/** TYPE IMPORTS */
import type { ListHeaderMenuColumn } from '@components/headers';
import type { FilterField } from '@components/drawers';
import type { DataRowColumn } from '@components/cards';
import type {
  WorkflowListPageFilters,
  WorkflowListPageColumnVisibility,
  WorkflowListItem,
} from './interfaces/workflowListPage.interface';

/** VUE IMPORTS */
import { ref, computed, onMounted } from 'vue';
import { useRouter } from 'vue-router';

/** COMPONENTS */
import { PageHeader, ListHeaderMenu } from '@components/headers';
import { ListCardEmpty, DataRow } from '@components/cards';
import { WorkflowDetailsDrawer, AdvancedFiltersDrawer } from '@components/drawers';
import { ListPagination } from '@components/navigation';
import { AppTooltip } from '@components/tooltips';

/** COMPOSABLES */
import { useWorkflowListPageTranslations } from '@composables/i18n';
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
  WORKFLOW_COLUMN_VISIBILITY_DEFAULTS,
  WORKFLOW_FILTER_DEFAULTS,
} from './constants';

/** COMPOSABLES & STORES */
const t = useWorkflowListPageTranslations();
const router = useRouter();
const { canCreate, canUpdate, canDelete, canRead } = usePermissions();
const canCreateWorkflow = canCreate('workflows');
const canUpdateWorkflow = canUpdate('workflows');
const canDeleteWorkflow = canDelete('workflows');
const canReadWorkflow = canRead('workflows');

/** STATE */
const workflowsList = ref<WorkflowListItem[]>([]);
const loading = ref(false);
const lastUpdatedAt = ref<number | undefined>(undefined);
const showDetailsDrawer = ref(false);
const showFiltersDrawer = ref(false);
const selectedWorkflowId = ref<string | null>(null);
const itemsPerPage = ref(DEFAULT_ITEMS_PER_PAGE);
const currentPage = ref(1);
const totalPages = ref(1);
const totalItems = ref(0);
const filters = ref<WorkflowListPageFilters>({ ...WORKFLOW_FILTER_DEFAULTS });
const columnVisibilityState = ref<WorkflowListPageColumnVisibility>({ ...WORKFLOW_COLUMN_VISIBILITY_DEFAULTS });

// Quick filter state (inline)
const quickSearch = ref('');
const quickStatus = ref<boolean | null>(null);

// Advanced filter state (drawer)
const advancedFilterValues = ref<Record<string, any>>({
  isTemplate: null,
  version: null,
  health: null,
  nodesCount: null,
  pluginType: null,
});

// Pending changes state
const hasPendingAdvancedFilters = ref(false);

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
 * Advanced filter fields for the drawer
 */
const advancedFilterFields = computed((): FilterField[] => [
  {
    key: 'isTemplate',
    label: t.filters.isTemplate.value,
    type: 'toggle',
    icon: 'content_copy',
    options: [
      { label: t.filters.allStatus.value, value: null },
      { label: t.filters.options.yes.value, value: true },
      { label: t.filters.options.no.value, value: false },
    ],
  },
  {
    key: 'version',
    label: t.filters.version.value,
    type: 'input',
    icon: 'tag',
    placeholder: t.filters.versionPlaceholder.value,
    inputType: 'number',
  },
  {
    key: 'health',
    label: t.filters.health.value,
    type: 'select',
    icon: 'monitor_heart',
    placeholder: t.filters.health.value,
    options: [
      { label: t.filters.allStatus.value, value: null },
      { label: t.filters.options.valid.value, value: 'valid' },
      { label: t.filters.options.pluginMissing.value, value: 'plugin_missing' },
      { label: t.filters.options.invalid.value, value: 'invalid' },
    ],
  },
  {
    key: 'nodesCount',
    label: t.filters.nodesCount.value,
    type: 'input',
    icon: 'hub',
    placeholder: t.filters.nodesCountPlaceholder.value,
    inputType: 'number',
  },
  {
    key: 'pluginType',
    label: t.filters.pluginType.value,
    type: 'input',
    icon: 'extension',
    placeholder: t.filters.pluginTypePlaceholder.value,
  },
]);

/**
 * Check if any filters are active (quick or advanced)
 */
const hasActiveFilters = computed(() => {
  return !!(
    quickSearch.value ||
    quickStatus.value !== null ||
    advancedFilterValues.value.isTemplate !== null ||
    advancedFilterValues.value.version !== null ||
    advancedFilterValues.value.health !== null ||
    advancedFilterValues.value.nodesCount !== null ||
    advancedFilterValues.value.pluginType !== null
  );
});

/**
 * Count of active advanced filters (for badge)
 */
const advancedFiltersCount = computed(() => {
  let count = 0;
  if (advancedFilterValues.value.isTemplate !== null) count++;
  if (advancedFilterValues.value.version !== null) count++;
  if (advancedFilterValues.value.health !== null) count++;
  if (advancedFilterValues.value.nodesCount !== null) count++;
  if (advancedFilterValues.value.pluginType !== null) count++;
  return count;
});

/**
 * Active filter chips for visual feedback
 */
const activeFilterChips = computed(() => {
  const chips: Array<{ key: string; label: string; value: string }> = [];

  if (quickSearch.value) {
    chips.push({
      key: 'name',
      label: t.filters.workflowName.value,
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

  if (advancedFilterValues.value.isTemplate !== null) {
    chips.push({
      key: 'isTemplate',
      label: t.filters.isTemplate.value,
      value: advancedFilterValues.value.isTemplate ? t.filters.options.yes.value : t.filters.options.no.value,
    });
  }

  if (advancedFilterValues.value.version !== null) {
    chips.push({
      key: 'version',
      label: t.filters.version.value,
      value: `v${advancedFilterValues.value.version}`,
    });
  }

  if (advancedFilterValues.value.health !== null) {
    chips.push({
      key: 'health',
      label: t.filters.health.value,
      value: String(advancedFilterValues.value.health),
    });
  }

  if (advancedFilterValues.value.nodesCount !== null) {
    chips.push({
      key: 'nodesCount',
      label: t.filters.nodesCount.value,
      value: `≥ ${advancedFilterValues.value.nodesCount}`,
    });
  }

  if (advancedFilterValues.value.pluginType !== null) {
    chips.push({
      key: 'pluginType',
      label: t.filters.pluginType.value,
      value: String(advancedFilterValues.value.pluginType),
    });
  }

  return chips;
});

/** Maximum number of visible filter chips */
const MAX_VISIBLE_CHIPS = 2;

/**
 * Chips that are visible (limited to MAX_VISIBLE_CHIPS)
 */
const visibleFilterChips = computed(() => {
  return activeFilterChips.value.slice(0, MAX_VISIBLE_CHIPS);
});

/**
 * Chips that are hidden (beyond the limit)
 */
const hiddenFilterChips = computed(() => {
  return activeFilterChips.value.slice(MAX_VISIBLE_CHIPS);
});

/**
 * Count of hidden filters
 */
const hiddenFiltersCount = computed(() => {
  return hiddenFilterChips.value.length;
});

/**
 * Column visibility using ListHeaderMenuColumn format with reactive translations
 */
const menuColumns = computed((): ListHeaderMenuColumn[] => [
  { key: 'version', label: t.menuColumns.version.value, visible: columnVisibilityState.value.version },
  { key: 'nodesCount', label: t.menuColumns.nodesCount.value, visible: columnVisibilityState.value.nodesCount },
  { key: 'pluginsCount', label: t.menuColumns.pluginsCount.value, visible: columnVisibilityState.value.pluginsCount },
]);

/**
 * Filtered columns based on visibility state
 */
const visibleColumns = computed(() => {
  return t.columns.value.filter((col: DataRowColumn) => {
    // Always show avatar, name, and status (health)
    if (col.key === 'icon' || col.key === 'name' || col.key === 'status') {
      return true;
    }

    // Filter based on columnVisibility
    if (col.key === 'definitionVersion') return columnVisibilityState.value.version;
    if (col.key === 'nodesCount') return columnVisibilityState.value.nodesCount;
    if (col.key === 'pluginsCount') return columnVisibilityState.value.pluginsCount;

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
  void fetchWorkflows();
}

/**
 * Handle advanced filters apply from drawer
 * @param {Record<string, any>} appliedFilters - Applied filter values
 * @returns {void}
 */
function handleAdvancedFiltersApply(appliedFilters: Record<string, any>): void {
  advancedFilterValues.value = { ...appliedFilters };
  hasPendingAdvancedFilters.value = false;
  currentPage.value = 1;
  showFiltersDrawer.value = false;
  void fetchWorkflows();
}

/**
 * Handle pending state change from advanced filters drawer
 * @param {boolean} hasPending - Whether there are pending changes
 * @returns {void}
 */
function handlePendingChange(hasPending: boolean): void {
  hasPendingAdvancedFilters.value = hasPending;
}

/**
 * Handle advanced filters reset from drawer
 * @returns {void}
 */
function handleAdvancedFiltersReset(): void {
  advancedFilterValues.value = { isTemplate: null, version: null, health: null, nodesCount: null, pluginType: null };
  currentPage.value = 1;
  void fetchWorkflows();
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
  } else if (key === 'isTemplate') {
    advancedFilterValues.value.isTemplate = null;
  } else if (key === 'version') {
    advancedFilterValues.value.version = null;
  } else if (key === 'health') {
    advancedFilterValues.value.health = null;
  } else if (key === 'nodesCount') {
    advancedFilterValues.value.nodesCount = null;
  } else if (key === 'pluginType') {
    advancedFilterValues.value.pluginType = null;
  }

  currentPage.value = 1;
  void fetchWorkflows();
}

/**
 * Clear all filters (quick and advanced)
 * @returns {void}
 */
function clearAllFilters(): void {
  quickSearch.value = '';
  quickStatus.value = null;
  advancedFilterValues.value = { isTemplate: null, version: null, health: null, nodesCount: null, pluginType: null };
  filters.value = { ...WORKFLOW_FILTER_DEFAULTS };
  hasPendingAdvancedFilters.value = false;
  currentPage.value = 1;
  void fetchWorkflows();
}

/**
 * Fetch workflows from API with server-side filtering and pagination
 * @returns {Promise<void>}
 */
async function fetchWorkflows(): Promise<void> {
  try {
    loading.value = true;

    const queryParams: Record<string, any> = {
      page: currentPage.value,
      perPage: itemsPerPage.value,
    };

    // Quick filters
    if (filters.value.name) queryParams.name = filters.value.name;
    if (typeof filters.value.status === 'boolean') queryParams.enabled = filters.value.status;

    // Advanced filters (server-side)
    if (advancedFilterValues.value.isTemplate !== null) queryParams.isTemplate = advancedFilterValues.value.isTemplate;
    if (advancedFilterValues.value.version !== null) queryParams.definitionVersion = Number(advancedFilterValues.value.version);
    if (advancedFilterValues.value.health !== null) queryParams.status = advancedFilterValues.value.health;

    const response = await apis.workflows.definition.list(queryParams);

    // Map API response to WorkflowListItem
    workflowsList.value = (response.items || []).map((def: any): WorkflowListItem => ({
      id: def._id || '',
      name: def.name || '',
      description: def.description || '',
      enabled: def.enabled ?? true,
      isTemplate: def.isTemplate ?? false,
      definitionVersion: def.definitionVersion || 1,
      nodesCount: def.nodes?.length || 0,
      edgesCount: def.edges?.length || 0,
      timezone: def.timezone?.value || 'UTC',
      created: def.created || '',
      updated: def.updated || '',
      status: def.status || 'valid',
      missingPlugins: def.missingPlugins || [],
      pluginsCount: def.installedPlugins?.length || 0,
    }));

    // Update pagination from response
    if (response.pagination) {
      totalItems.value = response.pagination.totalItems || 0;
      totalPages.value = response.pagination.totalPages || 1;
    } else {
      totalItems.value = workflowsList.value.length;
      totalPages.value = 1;
    }
  } catch (error: any) {
    handleApiError(error, { defaultMessage: 'Failed to load workflows' });
    workflowsList.value = [];
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
  void fetchWorkflows();
}

/**
 * Handle items per page change
 * @param {number} newValue - New items per page value
 * @returns {void}
 */
function handleItemsPerPageChange(newValue: number): void {
  itemsPerPage.value = newValue;
  currentPage.value = 1;
  void fetchWorkflows();
}

/**
 * Update menu columns when changed
 * @param {ListHeaderMenuColumn[]} columns - Updated column visibility states
 * @returns {void}
 */
function handleColumnsUpdate(columns: ListHeaderMenuColumn[]): void {
  columns.forEach(col => {
    if (col.key === 'version') columnVisibilityState.value.version = col.visible;
    if (col.key === 'nodesCount') columnVisibilityState.value.nodesCount = col.visible;
    if (col.key === 'pluginsCount') columnVisibilityState.value.pluginsCount = col.visible;
  });
}

/**
 * View workflow details in drawer
 * @param {WorkflowListItem} workflow - Workflow to view
 * @returns {void}
 */
function viewDetails(workflow: WorkflowListItem): void {
  if (!canReadWorkflow.value) return;
  selectedWorkflowId.value = workflow.id;
  showDetailsDrawer.value = true;
}

/**
 * Navigate to workflow editor
 * @param {WorkflowListItem} workflow - Workflow to edit
 * @returns {void}
 */
function editWorkflow(workflow: WorkflowListItem): void {
  if (!canUpdateWorkflow.value) return;
  void router.push(`/workflows/edit/${workflow.id}`);
}

/**
 * Handle edit event from drawer
 * @param {string} workflowId - ID of workflow to edit
 * @returns {void}
 */
function handleDrawerEdit(workflowId: string): void {
  if (!canUpdateWorkflow.value) return;
  void router.push(`/workflows/edit/${workflowId}`);
}

/**
 * Confirm delete operation with dialog
 * @param {WorkflowListItem} workflow - Workflow to delete
 * @returns {Promise<void>}
 */
async function confirmDelete(workflow: WorkflowListItem): Promise<void> {
  if (!canDeleteWorkflow.value) return;
  const workflowName = workflow.name || 'this workflow';
  const confirmed = await dialogDelete({
    title: t.dialog.deleteTitle.value,
    message: t.messages.confirmDelete(workflowName),
  });

  if (confirmed) {
    await deleteWorkflow(workflow);
  }
}

/**
 * Delete workflow via API
 * @param {WorkflowListItem} workflow - Workflow to delete
 * @returns {Promise<void>}
 */
async function deleteWorkflow(workflow: WorkflowListItem): Promise<void> {
  try {
    await apis.workflows.definition.delete({ workflowId: workflow.id });

    notifySuccess({ message: t.messages.deletedSuccessfully.value });

    // If current page would be empty, go to previous page
    if (workflowsList.value.length <= 1 && currentPage.value > 1) {
      currentPage.value -= 1;
    }

    await fetchWorkflows();
  } catch (error: any) {
    handleApiError(error, { defaultMessage: 'Failed to delete workflow' });
  }
}

/**
 * Navigate to workflow creation page
 * @returns {void}
 */
function navigateToCreate(): void {
  void router.push('/workflows/add');
}

/** LIFECYCLE HOOKS */
onMounted(async () => {
  await fetchWorkflows();
});

/**
 * Auto-refresh when organization changes
 * Resets pagination, filters and refetches data
 */
useOrgChangeRefresh(async () => {
  currentPage.value = 1;
  quickSearch.value = '';
  quickStatus.value = null;
  advancedFilterValues.value = { isTemplate: null, version: null, health: null, nodesCount: null, pluginType: null };
  filters.value = { ...WORKFLOW_FILTER_DEFAULTS };
  await fetchWorkflows();
});
</script>

<template>
  <q-page class="q-pt-lg">

    <!-- Header Section -->
    <PageHeader
      icon="account_tree"
      icon-color="primary"
      :title="t.page.title.value"
      :description="t.page.description.value"
      :button="canCreateWorkflow ? { label: t.page.addButton.value, icon: 'add', onClick: navigateToCreate, color: 'primary' } : undefined"
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

      <!-- Filter Icon Button -->
      <div class="col-auto">
        <q-btn
          round
          flat
          icon="tune"
          color="grey-7"
          @click="showFiltersDrawer = true"
        >
          <q-badge
            v-if="advancedFiltersCount > 0 || hasPendingAdvancedFilters"
            :color="hasPendingAdvancedFilters ? 'warning' : 'primary'"
            floating
            rounded
            :label="advancedFiltersCount || '!'"
          />
          <AppTooltip :content="hasPendingAdvancedFilters
            ? t.filters.pendingFilters.value
            : t.filters.advancedFilters.value"
          />
        </q-btn>
      </div>
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
        color="primary"
        size="sm"
        @remove="removeFilter(chip.key)"
      >
        <span class="text-weight-medium">{{ chip.label }}:</span>&nbsp;{{ chip.value }}
      </q-chip>

      <!-- +N Badge for hidden filters -->
      <q-badge
        v-if="hiddenFiltersCount > 0"
        color="primary"
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
          <q-icon name="account_tree" size="sm" color="primary" class="q-mr-sm" />
          <div class="text-subtitle1 text-weight-medium text-primary">{{ t.page.listTitle.value }}</div>
        </div>
      </div>
      <div class="col-auto">
        <ListHeaderMenu
          icon="account_tree"
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
          @refresh="fetchWorkflows"
        />
      </div>
    </div>

    <!-- Loading Spinner -->
    <div v-if="loading" class="row justify-center q-my-lg">
      <q-spinner color="primary" size="3em" />
    </div>

    <!-- Workflow Row List -->
    <div v-else class="row">
      <div
        v-for="(workflow, index) in workflowsList"
        :key="workflow.id || `workflow-${index}`"
        class="col-12 q-mb-xs"
      >
        <DataRow
          :data="workflow"
          :columns="visibleColumns"
          :actions="{ showEdit: canUpdateWorkflow, showView: canReadWorkflow, showDelete: canDeleteWorkflow }"
          @click="viewDetails"
          @dblclick="editWorkflow"
          @edit="editWorkflow"
          @view="viewDetails"
          @delete="confirmDelete"
        />
      </div>

      <!-- No Results -->
      <div v-if="workflowsList.length === 0" class="col-12">
        <ListCardEmpty
          :title="t.empty.title.value"
          :description="t.empty.description.value"
          icon="account_tree"
        />
      </div>
    </div>

    <!-- Pagination -->
    <ListPagination
      v-model="currentPage"
      :total-pages="totalPages"
      @change="handlePageChange"
    />

    <!-- Workflow Details Drawer -->
    <WorkflowDetailsDrawer
      v-model="showDetailsDrawer"
      :workflow-id="selectedWorkflowId"
      @edit="handleDrawerEdit"
    />

    <!-- Advanced Filters Drawer -->
    <AdvancedFiltersDrawer
      v-model="showFiltersDrawer"
      :fields="advancedFilterFields"
      :values="advancedFilterValues"
      @apply="handleAdvancedFiltersApply"
      @reset="handleAdvancedFiltersReset"
      @pending-change="handlePendingChange"
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
