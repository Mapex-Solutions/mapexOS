<script setup lang="ts">
defineOptions({
  name: 'AuditLogsPage'
});

/** TYPE IMPORTS */
import type { ListHeaderMenuColumn } from '@components/headers';
import type { DataRowColumn } from '@components/cards';
import type { FilterField, FilterValues } from '@components/drawers';
import type {
  AuditLogEntry,
  AuditLogsPageFilters,
  AuditLogsPageColumnVisibility,
} from './interfaces';

/** VUE IMPORTS */
import { ref, computed } from 'vue';

/** COMPONENTS */
import { PageHeader, ListHeaderMenu } from '@components/headers';
import { ListCardEmpty, DataRow } from '@components/cards';
import { JsonDrawer, AdvancedFiltersDrawer } from '@components/drawers';
import { AppTooltip } from '@components/tooltips';

/** COMPOSABLES */
import { useLogger } from '@composables/useLogger';
import { useAuditLogsPageTranslations } from '@composables/i18n';

/** LOCAL IMPORTS */
import { AUDIT_LOG_LIST_STUB } from './stubs';
import {
  DEFAULT_ITEMS_PER_PAGE,
  INITIAL_PAGE,
  MAX_VISIBLE_CHIPS,
  COLUMN_VISIBILITY_DEFAULTS,
  FILTER_DEFAULTS,
  ACTION_COLORS,
  TYPE_ICONS,
  DEFAULT_TYPE_ICON,
  DEFAULT_ACTION_COLOR,
} from './constants';

const logger = useLogger('AuditLogsPage');

/** COMPOSABLES & STORES */
const t = useAuditLogsPageTranslations();

/** STATE */
const logsList = ref<AuditLogEntry[]>([...AUDIT_LOG_LIST_STUB]);
const selectedEvent = ref<AuditLogEntry | null>(null);
const jsonDrawerOpen = ref(false);
const loading = ref(false);
const lastUpdatedAt = ref<number | undefined>(undefined);
const itemsPerPage = ref(DEFAULT_ITEMS_PER_PAGE);
const currentPage = ref(INITIAL_PAGE);
const totalPages = ref(1);
const totalItems = ref(AUDIT_LOG_LIST_STUB.length);
const filters = ref<AuditLogsPageFilters>({ ...FILTER_DEFAULTS });
const columnVisibilityState = ref<AuditLogsPageColumnVisibility>({ ...COLUMN_VISIBILITY_DEFAULTS });

/** FILTER STATE */
const showFiltersDrawer = ref(false);
const quickSearchName = ref('');
const quickStatusEnabled = ref<string | null>(null);
const advancedFilterValues = ref<FilterValues>({
  includeChildren: null,
  action: null,
  resourceType: null,
  dateRange: null,
});
const hasPendingAdvancedFilters = ref(false);

/** COMPUTED */

/**
 * Status options for quick filter
 */
const statusOptions = computed(() => [
  { label: t.filters.allStatus.value, value: null },
  { label: t.filters.options.success.value, value: 'success' },
  { label: t.filters.options.failure.value, value: 'failure' },
]);

/**
 * Action options for advanced filter
 */
const actionOptions = computed(() => [
  { label: t.filters.allStatus.value, value: null },
  { label: t.actionOptions.create.value, value: 'Create' },
  { label: t.actionOptions.edit.value, value: 'Edit' },
  { label: t.actionOptions.delete.value, value: 'Delete' },
]);

/**
 * Resource type options for advanced filter
 */
const resourceTypeOptions = computed(() => [
  { label: t.filters.allStatus.value, value: null },
  { label: t.resourceTypeOptions.userLog.value, value: 'userLog' },
  { label: t.resourceTypeOptions.dataSource.value, value: 'dataSource' },
  { label: t.resourceTypeOptions.assets.value, value: 'assets' },
  { label: t.resourceTypeOptions.payloadHandler.value, value: 'payloadHandler' },
  { label: t.resourceTypeOptions.triggers.value, value: 'triggers' },
  { label: t.resourceTypeOptions.users.value, value: 'users' },
  { label: t.resourceTypeOptions.customers.value, value: 'customers' },
]);

/**
 * Advanced filter fields configuration
 */
const advancedFilterFields = computed((): FilterField[] => [
  {
    key: 'includeChildren',
    type: 'toggle',
    label: t.filters.includeChildrenOrgs.value,
    icon: 'account_tree',
    options: [
      { label: t.filters.allStatus.value, value: null },
      { label: t.filters.options.yes.value, value: true },
      { label: t.filters.options.no.value, value: false },
    ],
  },
  {
    key: 'action',
    type: 'select',
    label: t.filters.action.value,
    icon: 'bolt',
    options: actionOptions.value,
  },
  {
    key: 'resourceType',
    type: 'select',
    label: t.filters.resourceType.value,
    icon: 'category',
    options: resourceTypeOptions.value,
  },
]);

/**
 * Check if any filter is active
 */
const hasActiveFilters = computed(() => {
  return !!(
    filters.value.search ||
    filters.value.status ||
    filters.value.includeChildren !== undefined ||
    filters.value.action ||
    filters.value.resourceType
  );
});

/**
 * Count of advanced filters only (for badge on filter icon)
 */
const advancedFiltersCount = computed(() => {
  let count = 0;
  if (filters.value.includeChildren !== undefined) count++;
  if (filters.value.action) count++;
  if (filters.value.resourceType) count++;
  return count;
});

/**
 * Active filter chips for display
 */
const activeFilterChips = computed(() => {
  const chips: Array<{ key: string; label: string; value: string }> = [];

  if (filters.value.search) {
    chips.push({ key: 'search', label: t.filters.actor.value, value: filters.value.search });
  }
  if (filters.value.status) {
    chips.push({
      key: 'status',
      label: t.filters.status.value,
      value: filters.value.status === 'success' ? t.filters.options.success.value : t.filters.options.failure.value
    });
  }
  if (filters.value.includeChildren !== undefined) {
    chips.push({
      key: 'includeChildren',
      label: t.filters.includeChildren.value,
      value: filters.value.includeChildren ? t.filters.options.yes.value : t.filters.options.no.value
    });
  }
  if (filters.value.action) {
    chips.push({ key: 'action', label: t.filters.action.value, value: filters.value.action });
  }
  if (filters.value.resourceType) {
    const resourceLabel = resourceTypeOptions.value.find(o => o.value === filters.value.resourceType)?.label || filters.value.resourceType;
    chips.push({ key: 'resourceType', label: t.filters.resourceType.value, value: resourceLabel });
  }

  return chips;
});

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
 * Column visibility using ListHeaderMenuColumn format
 */
const menuColumns = computed((): ListHeaderMenuColumn[] => [
  { key: 'action', label: t.menuColumns.action.value, visible: columnVisibilityState.value.action },
  { key: 'resource', label: t.menuColumns.resource.value, visible: columnVisibilityState.value.resource },
  { key: 'created', label: t.menuColumns.timestamp.value, visible: columnVisibilityState.value.created },
]);

/**
 * Get icon name for resource type
 * @param {string} type - Resource type
 * @returns {string} Icon name
 */
function getTypeIcon(type: string): string {
  return TYPE_ICONS[type] || DEFAULT_TYPE_ICON;
}

/**
 * Get color for action type
 * @param {string} action - Action type
 * @returns {string} Color class name
 */
function getActionColor(action: string): string {
  return ACTION_COLORS[action] || DEFAULT_ACTION_COLOR;
}

/**
 * Log columns configuration
 */
const logColumns = ref<DataRowColumn[]>([
  {
    key: 'icon',
    label: '',
    type: 'avatar',
    visible: 'always',
    width: 56,
    icon: (_value: any, row: any) => getTypeIcon(row.type),
    color: (_value: any, row: any) => row.status === 'success' ? 'primary' : 'red-5',
  },
  {
    key: 'actor',
    label: 'Actor',
    type: 'text',
    visible: 'always',
    width: 200,
    ellipsis: true,
    secondaryKey: 'details',
  },
  {
    key: 'action',
    label: 'Action',
    type: 'chip',
    visible: 'laptop',
    width: 100,
    color: (value: any) => getActionColor(value),
  },
  {
    key: 'resource',
    label: 'Resource',
    type: 'text',
    visible: 'laptop',
    width: 180,
    ellipsis: true,
  },
  {
    key: 'status',
    label: 'Status',
    type: 'badge',
    visible: 'laptop',
    width: 110,
    format: (value: any) => value ? value.toUpperCase() : 'UNKNOWN',
    color: (value: any) => value === 'success' ? 'green-6' : 'red-6',
  },
  {
    key: 'created',
    label: 'Timestamp',
    type: 'text',
    visible: 'laptop',
    width: 180,
    format: (value: any) => {
      if (!value) return 'N/A';
      const date = new Date(value);
      return date.toLocaleDateString('en-US', {
        day: '2-digit',
        month: 'short',
        year: 'numeric',
        hour: '2-digit',
        minute: '2-digit',
      });
    },
  },
]);

/**
 * Filtered columns based on visibility
 */
const visibleColumns = computed(() => {
  return logColumns.value.filter((col: any) => {
    // Always show avatar, actor, and status
    if (col.key === 'icon' || col.key === 'actor' || col.key === 'status') {
      return true;
    }

    // Filter based on columnVisibility
    if (col.key === 'action') return columnVisibilityState.value.action;
    if (col.key === 'resource') return columnVisibilityState.value.resource;
    if (col.key === 'created') return columnVisibilityState.value.created;

    return true;
  });
});

/** FUNCTIONS */

/**
 * Fetch logs (mock implementation)
 * @returns {void}
 */
function fetchLogs(): void {
  loading.value = true;
  try {
    // Mock implementation - filter the stub data based on filters
    let filteredLogs = [...AUDIT_LOG_LIST_STUB];

    if (filters.value.search) {
      const searchTerm = filters.value.search.toLowerCase();
      filteredLogs = filteredLogs.filter(log =>
        log.actor.toLowerCase().includes(searchTerm) ||
        log.resource.toLowerCase().includes(searchTerm)
      );
    }

    if (filters.value.status) {
      filteredLogs = filteredLogs.filter(log => log.status === filters.value.status);
    }

    if (filters.value.action) {
      filteredLogs = filteredLogs.filter(log => log.action === filters.value.action);
    }

    if (filters.value.resourceType) {
      filteredLogs = filteredLogs.filter(log => log.type === filters.value.resourceType);
    }

    logsList.value = filteredLogs;
    totalItems.value = filteredLogs.length;
    totalPages.value = Math.ceil(filteredLogs.length / itemsPerPage.value) || 1;
  } catch (err) {
    logger.error('Failed to fetch logs:', err);
  } finally {
    loading.value = false;
    lastUpdatedAt.value = Date.now();
  }
}

/**
 * Apply quick filters (search + status)
 */
function applyQuickFilters(): void {
  if (quickSearchName.value) {
    filters.value.search = quickSearchName.value;
  } else {
    delete filters.value.search;
  }
  if (quickStatusEnabled.value !== null) {
    filters.value.status = quickStatusEnabled.value;
  } else {
    delete filters.value.status;
  }
  currentPage.value = 1;
  fetchLogs();
}

/**
 * Handle advanced filters apply from drawer
 * @param {FilterValues} values - Applied filter values
 */
function handleAdvancedFiltersApply(values: FilterValues): void {
  advancedFilterValues.value = values;
  if (values.includeChildren !== null && values.includeChildren !== undefined) {
    filters.value.includeChildren = values.includeChildren;
  } else {
    delete filters.value.includeChildren;
  }
  if (values.action) {
    filters.value.action = values.action;
  } else {
    delete filters.value.action;
  }
  if (values.resourceType) {
    filters.value.resourceType = values.resourceType;
  } else {
    delete filters.value.resourceType;
  }
  currentPage.value = 1;

  // Reset pending state after apply
  hasPendingAdvancedFilters.value = false;

  fetchLogs();
}

/**
 * Handle pending state change from advanced filters drawer
 * @param {boolean} hasPending - Whether there are pending changes
 */
function handlePendingChange(hasPending: boolean): void {
  hasPendingAdvancedFilters.value = hasPending;
}

/**
 * Handle advanced filters reset from drawer
 */
function handleAdvancedFiltersReset(): void {
  advancedFilterValues.value = {
    includeChildren: null,
    action: null,
    resourceType: null,
    dateRange: null,
  };
}

/**
 * Remove individual filter
 * @param {string} key - Filter key to remove
 */
function removeFilter(key: string): void {
  if (key === 'search') {
    delete filters.value.search;
    quickSearchName.value = '';
  } else if (key === 'status') {
    delete filters.value.status;
    quickStatusEnabled.value = null;
  } else if (key === 'includeChildren') {
    delete filters.value.includeChildren;
    advancedFilterValues.value.includeChildren = null;
  } else if (key === 'action') {
    delete filters.value.action;
    advancedFilterValues.value.action = null;
  } else if (key === 'resourceType') {
    delete filters.value.resourceType;
    advancedFilterValues.value.resourceType = null;
  }

  currentPage.value = 1;
  fetchLogs();
}

/**
 * Clear all filters
 */
function clearAllFilters(): void {
  // Reset filter state
  filters.value = { ...FILTER_DEFAULTS };

  // Reset quick filters
  quickSearchName.value = '';
  quickStatusEnabled.value = null;

  // Reset advanced filters
  advancedFilterValues.value = {
    includeChildren: null,
    action: null,
    resourceType: null,
    dateRange: null,
  };

  // Reset pending state
  hasPendingAdvancedFilters.value = false;

  currentPage.value = 1;
  fetchLogs();
}

/**
 * Update menu columns when changed
 * @param {ListHeaderMenuColumn[]} columns - Updated column visibility states
 * @returns {void}
 */
function handleColumnsUpdate(columns: ListHeaderMenuColumn[]): void {
  columns.forEach((col) => {
    if (col.key === 'action') columnVisibilityState.value.action = col.visible;
    if (col.key === 'resource') columnVisibilityState.value.resource = col.visible;
    if (col.key === 'created') columnVisibilityState.value.created = col.visible;
  });
}

/**
 * Handle card click event - opens JSON drawer
 * @param {AuditLogEntry} item - Clicked audit log item
 * @returns {void}
 */
function cardClick(item: AuditLogEntry): void {
  selectedEvent.value = item;
  jsonDrawerOpen.value = true;
}

/**
 * Edit log entry (placeholder)
 * @param {AuditLogEntry} log - Log entry to edit
 * @returns {void}
 */
function editLog(log: AuditLogEntry): void {
  logger.debug('Edit log:', log);
}

/**
 * View log entry details (placeholder)
 * @param {AuditLogEntry} log - Log entry to view
 * @returns {void}
 */
function viewDetails(log: AuditLogEntry): void {
  logger.debug('View log details:', log);
}

/**
 * Handle pagination navigation
 * @param {number} page - New page number
 */
function handlePageChange(page: number): void {
  currentPage.value = page;
  fetchLogs();
}

/**
 * Handle items per page change
 * @param {number} newValue - New items per page value
 */
function handleItemsPerPageChange(newValue: number): void {
  itemsPerPage.value = newValue;
  currentPage.value = 1;
  fetchLogs();
}
</script>

<template>
  <q-page class="q-pt-lg">

    <!-- Header Section -->
    <PageHeader
        icon="fact_check"
        iconColor="primary"
        :title="t.pageHeader.title.value"
        :description="t.pageHeader.description.value"
    />

    <!-- Filters Section -->
    <div class="text-caption text-grey-7 q-mb-xs">{{ t.filters.label.value }}</div>
    <div class="row items-center q-col-gutter-sm q-mb-md">
      <!-- Search Input -->
      <div class="col">
        <q-input
          v-model="quickSearchName"
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
          v-model="quickStatusEnabled"
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
          <q-icon name="fact_check" size="sm" color="primary" class="q-mr-sm"/>
          <div class="text-subtitle1 text-weight-medium text-primary">{{ t.listTitle.value }}</div>
        </div>
      </div>
      <div class="col-auto">
        <ListHeaderMenu
          icon="fact_check"
          :item-label="t.itemLabel.value"
          :item-label-plural="t.itemLabelPlural.value"
          :items-count="totalItems"
          :items-per-page="itemsPerPage"
          :columns="menuColumns"
          :filtered="hasActiveFilters"
          :refreshing="loading"
          :last-updated-at="lastUpdatedAt"
          @update:items-per-page="handleItemsPerPageChange"
          @update:columns="handleColumnsUpdate"
          @refresh="fetchLogs"
        />
      </div>
    </div>

    <!-- Loading Spinner -->
    <div v-if="loading" class="row justify-center q-my-lg">
      <q-spinner color="primary" size="3em" />
    </div>

    <!-- Logs Row List -->
    <div v-else class="row">
      <div
          v-for="(log, index) in logsList"
          :key="log.id || `log-${index}`"
          class="col-12 q-mb-xs"
      >
        <DataRow
            :data="log"
            :columns="visibleColumns"
            @click="cardClick"
            @dblclick="editLog"
            @edit="editLog"
            @view="viewDetails"
        />
      </div>

      <!-- No Results -->
      <div v-if="logsList.length === 0" class="col-12">
        <ListCardEmpty
            :title="t.empty.title.value"
            :description="t.empty.description.value"
            icon="fact_check"
        />
      </div>
    </div>

    <!-- Pagination -->
    <div class="row justify-center q-mt-lg q-mb-lg">
      <q-pagination
          v-if="logsList.length"
          v-model="currentPage"
          direction-links
          boundary-links
          class="rounded-borders"
          color="primary"
          active-color="primary"
          :max="totalPages"
          @update:model-value="handlePageChange"
      />
    </div>

    <!-- JSON Drawer -->
    <JsonDrawer
      v-if="selectedEvent"
      v-model:show="jsonDrawerOpen"
      :title="t.drawer.title.value"
      :jsonData="selectedEvent"
      :editable="false"
      :subtitle="`${selectedEvent.type} • ${selectedEvent.created ? new Date(selectedEvent.created).toLocaleDateString('en-US', { day: '2-digit', month: 'short', year: 'numeric', hour: '2-digit', minute: '2-digit' }) : 'N/A'}`"
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
