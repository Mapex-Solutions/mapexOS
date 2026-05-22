<script setup lang="ts">
defineOptions({
  name: 'JsExecLogsPage'
});

/** TYPE IMPORTS */
import type { ListHeaderMenuColumn } from '@components/headers';
import type { DataRowColumn } from '@components/cards';
import type { EventsJsExecResponse } from '@mapexos/schemas';
import type { FilterField, FilterValues } from '@components/drawers';
import type { JsExecLogsPageFilters, JsExecLogsPageColumnVisibility, JsExecLogsPageCursor } from './interfaces';

/** VUE IMPORTS */
import { ref, computed, onMounted } from 'vue';

/** COMPONENTS */
import { PageHeader, ListHeaderMenu } from '@components/headers';
import { ListCardEmpty, DataRow } from '@components/cards';
import { JsonDrawer, AdvancedFiltersDrawer } from '@components/drawers';
import { AppTooltip } from '@components/tooltips';

/** COMPOSABLES */
import { useOrgChangeRefresh } from '@composables/organizations';
import { useJsExecLogsPageTranslations } from '@composables/i18n/pages/logs/jsExecLogsPage';
import { useLogger } from '@composables/useLogger';

/** UTILS */
import { notifyFail } from '@utils/alert';

/** SERVICES */
import { apis } from '@services/mapex';

/** LOCAL IMPORTS */
import {
  DEFAULT_LIMIT,
  COLUMN_VISIBILITY_DEFAULTS,
  SUCCESS_COLORS,
} from './constants';

/** TRANSLATIONS */
const t = useJsExecLogsPageTranslations();
const logger = useLogger('JsExecLogsPage');

/** Maximum number of visible filter chips */
const MAX_VISIBLE_CHIPS = 2;

/** STATE */
const loading = ref(false);
const lastUpdatedAt = ref<number | undefined>(undefined);
const logsList = ref<EventsJsExecResponse[]>([]);
const selectedEvent = ref<EventsJsExecResponse | null>(null);
const jsonDrawerOpen = ref(false);

/** Cursor Pagination */
const limit = ref(DEFAULT_LIMIT);
const cursor = ref<JsExecLogsPageCursor>({
  current: undefined,
  next: undefined,
  prev: undefined,
  hasNext: false,
  hasPrevious: false,
});

/**
 * Get today's date range (start of day to end of day in ISO format)
 */
function getTodayDateRange(): { startTime: string; endTime: string } {
  const now = new Date();
  const startOfDay = new Date(now.getFullYear(), now.getMonth(), now.getDate(), 0, 0, 0, 0);
  const endOfDay = new Date(now.getFullYear(), now.getMonth(), now.getDate(), 23, 59, 59, 999);
  return {
    startTime: startOfDay.toISOString(),
    endTime: endOfDay.toISOString(),
  };
}

/** Filters - Default date range is today */
const todayRange = getTodayDateRange();
const filters = ref<JsExecLogsPageFilters>({
  startTime: todayRange.startTime,
  endTime: todayRange.endTime,
});

/** FILTER STATE */
const showFiltersDrawer = ref(false);
const quickSearchUuid = ref('');
const quickStatusSuccess = ref<boolean | null>(null);
const advancedFilterValues = ref<FilterValues>({
  includeChildren: null,
  dateRange: { from: todayRange.startTime, to: todayRange.endTime },
  execTimeOp: null,
  execTimeValue: null,
  execTimeValueEnd: null,
});
const hasPendingAdvancedFilters = ref(false);

/** Column visibility */
const menuColumns = ref<ListHeaderMenuColumn[]>([
  { key: 'name', label: t.columns.dataSource.value, visible: COLUMN_VISIBILITY_DEFAULTS.name },
  { key: 'threadId', label: t.columns.uuid.value, visible: COLUMN_VISIBILITY_DEFAULTS.threadId },
  { key: 'success', label: t.columns.status.value, visible: COLUMN_VISIBILITY_DEFAULTS.success },
  { key: 'executionTime', label: t.columns.executionTime.value, visible: COLUMN_VISIBILITY_DEFAULTS.executionTime },
  { key: 'created', label: t.columns.timestamp.value, visible: COLUMN_VISIBILITY_DEFAULTS.created },
]);

/** COMPUTED */

/**
 * Column visibility state computed from menuColumns
 */
const columnVisibility = computed<JsExecLogsPageColumnVisibility>(() => ({
  threadId: menuColumns.value.find((col) => col.key === 'threadId')?.visible ?? true,
  name: menuColumns.value.find((col) => col.key === 'name')?.visible ?? true,
  success: menuColumns.value.find((col) => col.key === 'success')?.visible ?? true,
  executionTime: menuColumns.value.find((col) => col.key === 'executionTime')?.visible ?? true,
  created: menuColumns.value.find((col) => col.key === 'created')?.visible ?? true,
}));

/**
 * Status options for quick filter select
 */
const statusOptions = computed(() => [
  { label: t.filters.allStatus.value, value: null },
  { label: t.filters.options.success.value, value: true },
  { label: t.filters.options.failed.value, value: false },
]);

/**
 * Execution time operator options computed from translations
 */
const execTimeOperatorOptions = computed(() => [
  { label: t.execTimeOperators.lessThan.value, value: 'lt' },
  { label: t.execTimeOperators.lessThanOrEqual.value, value: 'lte' },
  { label: t.execTimeOperators.greaterThan.value, value: 'gt' },
  { label: t.execTimeOperators.greaterThanOrEqual.value, value: 'gte' },
  { label: t.execTimeOperators.between.value, value: 'between' },
]);

/**
 * Advanced filter fields configuration
 * includeChildren MUST be first and use type 'toggle' with 3 options (All/Yes/No)
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
    key: 'execTimeOp',
    type: 'select',
    label: t.filters.execTimeOperator.value,
    icon: 'timer',
    options: execTimeOperatorOptions.value,
  },
  {
    key: 'execTimeValue',
    type: 'input',
    label: t.filters.execTime.value,
    icon: 'speed',
    placeholder: t.placeholders.execTime.value,
    inputType: 'number',
  },
  {
    key: 'execTimeValueEnd',
    type: 'input',
    label: t.filters.execTimeEnd.value,
    icon: 'speed',
    placeholder: t.placeholders.execTimeEnd.value,
    inputType: 'number',
    disabled: advancedFilterValues.value.execTimeOp !== 'between',
  },
]);

/**
 * Check if any filter is active
 */
const hasActiveFilters = computed(() => {
  return !!(
    filters.value.threadId ||
    filters.value.success !== undefined ||
    filters.value.includeChildren !== undefined ||
    filters.value.execTimeOp ||
    filters.value.execTimeValue !== undefined
  );
});

/**
 * Count of advanced filters only (for badge on filter icon)
 */
const advancedFiltersCount = computed(() => {
  let count = 0;
  if (filters.value.includeChildren !== undefined) count++;
  if (filters.value.execTimeOp) count++;
  if (filters.value.execTimeValue !== undefined) count++;
  return count;
});

/**
 * Active filter chips for display
 */
const activeFilterChips = computed(() => {
  const chips: Array<{ key: string; label: string; value: string }> = [];

  if (filters.value.threadId) {
    chips.push({ key: 'threadId', label: t.filters.uuid.value, value: filters.value.threadId });
  }
  if (filters.value.success !== undefined) {
    chips.push({
      key: 'success',
      label: t.filters.status.value,
      value: filters.value.success ? t.filters.options.success.value : t.filters.options.failed.value
    });
  }
  if (filters.value.includeChildren !== undefined) {
    chips.push({
      key: 'includeChildren',
      label: t.filters.includeChildren.value,
      value: filters.value.includeChildren ? t.filters.options.yes.value : t.filters.options.no.value
    });
  }
  if (filters.value.execTimeOp && filters.value.execTimeValue !== undefined) {
    const opLabel = execTimeOperatorOptions.value.find(o => o.value === filters.value.execTimeOp)?.label || filters.value.execTimeOp;
    let valueStr = `${opLabel} ${filters.value.execTimeValue}ms`;
    if (filters.value.execTimeOp === 'between' && filters.value.execTimeValueEnd !== undefined) {
      valueStr = `${filters.value.execTimeValue}ms - ${filters.value.execTimeValueEnd}ms`;
    }
    chips.push({
      key: 'execTime',
      label: t.filters.execTime.value,
      value: valueStr
    });
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
 * Log columns configuration for DataRow
 */
const logColumns = computed<DataRowColumn[]>(() => [
  {
    key: 'icon',
    label: '',
    type: 'avatar',
    visible: 'always',
    width: 56,
    icon: () => 'code',
    color: (value: any, row: EventsJsExecResponse) => getSuccessColor(row.success),
  },
  {
    key: 'name',
    label: t.columns.dataSource.value,
    type: 'text',
    visible: 'always',
    width: 250,
    ellipsis: true,
    format: (value: string, row: EventsJsExecResponse) => value || row.threadId || t.defaults.unknown.value,
    secondaryKey: 'description',
  },
  {
    key: 'threadId',
    label: t.columns.uuid.value,
    type: 'text',
    visible: 'laptop',
    width: 180,
    ellipsis: true,
    copyable: true,
  },
  {
    key: 'success',
    label: t.columns.status.value,
    type: 'badge',
    visible: 'always',
    width: 100,
    format: (value: boolean) => value ? t.statusBadge.success.value : t.statusBadge.failed.value,
    color: (value: boolean) => getSuccessColor(value),
  },
  {
    key: 'totalExecutionTime',
    label: t.columns.execTime.value,
    type: 'badge',
    visible: 'laptop',
    width: 100,
    format: (value: number) => `${value ?? 0}ms`,
    color: (value: number) => value > 1000 ? 'orange-6' : 'grey-7',
  },
  {
    key: 'created',
    label: t.columns.timestamp.value,
    type: 'text',
    visible: 'laptop',
    width: 180,
    format: (value: string) => formatTimestamp(value),
  },
]);

/**
 * Filtered columns based on visibility settings
 */
const visibleColumns = computed(() => {
  return logColumns.value.filter((col) => {
    if (col.key === 'icon') return true;
    if (col.key === 'name') return columnVisibility.value.name;
    if (col.key === 'threadId') return columnVisibility.value.threadId;
    if (col.key === 'success') return columnVisibility.value.success;
    if (col.key === 'totalExecutionTime') return columnVisibility.value.executionTime;
    if (col.key === 'created') return columnVisibility.value.created;
    return true;
  });
});

/**
 * Selected event with parsed event field for JSON drawer
 * The event field comes as a JSON string from ClickHouse and needs to be parsed
 */
const selectedEventParsed = computed(() => {
  if (!selectedEvent.value) return null;

  const eventData = { ...selectedEvent.value } as any;

  // Parse event field if it's a string
  if (typeof eventData.event === 'string' && eventData.event) {
    try {
      eventData.event = JSON.parse(eventData.event);
    } catch (e) {
      // Keep as string if parsing fails
      logger.warn('Failed to parse event', e);
    }
  }

  return eventData;
});

/** WATCHERS */

/**
 * Re-fetch data when organization changes
 */
useOrgChangeRefresh(() => {
  resetCursor();
  void fetchData();
});

/** FUNCTIONS */

/**
 * Reset cursor state to initial values
 */
function resetCursor(): void {
  cursor.value = {
    current: undefined,
    next: undefined,
    prev: undefined,
    hasNext: false,
    hasPrevious: false,
  };
}

/**
 * Fetch JS executor events from API with current filters and cursor pagination
 *
 * @param {string} direction - Direction to fetch: 'next' (older) or 'prev' (newer)
 * @returns {Promise<void>}
 */
async function fetchData(direction: 'next' | 'prev' = 'next'): Promise<void> {
  loading.value = true;

  try {
    const queryParams: Record<string, any> = {
      limit: limit.value,
      direction,
    };

    // Set cursor based on direction
    if (direction === 'next' && cursor.value.next) {
      queryParams.cursor = cursor.value.next;
    } else if (direction === 'prev' && cursor.value.prev) {
      queryParams.cursor = cursor.value.prev;
    }

    // Add filters conditionally
    if (filters.value.threadId) {
      queryParams.threadId = filters.value.threadId;
    }
    if (typeof filters.value.success === 'boolean') {
      queryParams.success = filters.value.success;
    }
    if (filters.value.startTime) {
      queryParams.startTime = filters.value.startTime;
    }
    if (filters.value.endTime) {
      queryParams.endTime = filters.value.endTime;
    }
    if (typeof filters.value.includeChildren === 'boolean') {
      queryParams.includeChildren = filters.value.includeChildren;
    }

    // Add execution time filter
    if (filters.value.execTimeOp && filters.value.execTimeValue !== undefined) {
      queryParams.execTimeOp = filters.value.execTimeOp;
      queryParams.execTimeValue = filters.value.execTimeValue;
      // Only add end value for 'between' operator
      if (filters.value.execTimeOp === 'between' && filters.value.execTimeValueEnd !== undefined) {
        queryParams.execTimeValueEnd = filters.value.execTimeValueEnd;
      }
    }

    logger.debug('Fetching with params', queryParams);

    const response = await apis.events.events.listJsExec(queryParams);

    logger.debug('Response received', response);

    logsList.value = response.items || [];

    // Update cursor state
    cursor.value = {
      current: queryParams.cursor,
      next: response.nextCursor || undefined,
      prev: response.prevCursor || undefined,
      hasNext: response.hasNext,
      hasPrevious: response.hasPrevious,
    };

  } catch (error: any) {
    logger.error('Error fetching JS exec events', error);
    notifyFail({
      message: error?.response?.data?.message || t.messages.loadFailed.value,
      timeout: 5000,
    });
    logsList.value = [];
  } finally {
    loading.value = false;
    lastUpdatedAt.value = Date.now();
  }
}

/**
 * Refresh: resets cursor to initial and refetches from the first page with current filters.
 * Filters are preserved. Used by the ListHeaderMenu refresh button.
 */
async function refreshData(): Promise<void> {
  cursor.value = {
    current: undefined,
    next: undefined,
    prev: undefined,
    hasNext: false,
    hasPrevious: false,
  };
  await fetchData('next');
}

/**
 * Apply quick filters (search UUID + status)
 */
function applyQuickFilters(): void {
  if (quickSearchUuid.value) {
    filters.value.threadId = quickSearchUuid.value;
  } else {
    delete filters.value.threadId;
  }
  if (quickStatusSuccess.value !== null) {
    filters.value.success = quickStatusSuccess.value;
  } else {
    delete filters.value.success;
  }
  resetCursor();
  void fetchData();
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

  // Handle date range
  if (values.dateRange?.from && values.dateRange?.to) {
    filters.value.startTime = values.dateRange.from;
    filters.value.endTime = values.dateRange.to;
  } else {
    const today = getTodayDateRange();
    filters.value.startTime = today.startTime;
    filters.value.endTime = today.endTime;
  }

  // Handle execution time filter
  if (values.execTimeOp) {
    filters.value.execTimeOp = values.execTimeOp;
  } else {
    delete filters.value.execTimeOp;
  }
  if (values.execTimeValue) {
    filters.value.execTimeValue = Number(values.execTimeValue);
  } else {
    delete filters.value.execTimeValue;
  }
  if (values.execTimeValueEnd) {
    filters.value.execTimeValueEnd = Number(values.execTimeValueEnd);
  } else {
    delete filters.value.execTimeValueEnd;
  }

  // Reset pending state after apply
  hasPendingAdvancedFilters.value = false;

  resetCursor();
  void fetchData();
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
  const today = getTodayDateRange();
  advancedFilterValues.value = {
    includeChildren: null,
    dateRange: { from: today.startTime, to: today.endTime },
    execTimeOp: null,
    execTimeValue: null,
    execTimeValueEnd: null,
  };
}

/**
 * Handle advanced field change (for cascading filters like execTimeValueEnd disabled state)
 * @param {string} key - Field key that changed
 * @param {any} value - New value
 */
function handleAdvancedFieldChange(key: string, value: any): void {
  // Update local value for dynamic field states
  advancedFilterValues.value[key] = value;
}

/**
 * Remove individual filter
 * @param {string} key - Filter key to remove
 */
function removeFilter(key: string): void {
  if (key === 'threadId') {
    delete filters.value.threadId;
    quickSearchUuid.value = '';
  } else if (key === 'success') {
    delete filters.value.success;
    quickStatusSuccess.value = null;
  } else if (key === 'includeChildren') {
    delete filters.value.includeChildren;
    advancedFilterValues.value.includeChildren = null;
  } else if (key === 'execTime') {
    delete filters.value.execTimeOp;
    delete filters.value.execTimeValue;
    delete filters.value.execTimeValueEnd;
    advancedFilterValues.value.execTimeOp = null;
    advancedFilterValues.value.execTimeValue = null;
    advancedFilterValues.value.execTimeValueEnd = null;
  }

  resetCursor();
  void fetchData();
}

/**
 * Clear all filters
 */
function clearAllFilters(): void {
  // Reset filter state
  const today = getTodayDateRange();
  filters.value = {
    startTime: today.startTime,
    endTime: today.endTime,
  };

  // Reset quick filters
  quickSearchUuid.value = '';
  quickStatusSuccess.value = null;

  // Reset advanced filters
  advancedFilterValues.value = {
    includeChildren: null,
    dateRange: { from: today.startTime, to: today.endTime },
    execTimeOp: null,
    execTimeValue: null,
    execTimeValueEnd: null,
  };

  // Reset pending state
  hasPendingAdvancedFilters.value = false;

  resetCursor();
  void fetchData();
}

/**
 * Navigate to next page (older events)
 */
function goToNextPage(): void {
  if (cursor.value.hasNext) {
    void fetchData('next');
  }
}

/**
 * Navigate to previous page (newer events)
 */
function goToPrevPage(): void {
  if (cursor.value.hasPrevious) {
    void fetchData('prev');
  }
}

/**
 * Update menu columns when changed
 *
 * @param {ListHeaderMenuColumn[]} columns - Updated column visibility states
 */
function handleColumnsUpdate(columns: ListHeaderMenuColumn[]): void {
  menuColumns.value = columns;
}

/**
 * Handle limit change from header menu
 *
 * @param {number} newLimit - New limit value
 */
function handleLimitChange(newLimit: number): void {
  limit.value = newLimit;
  resetCursor();
  void fetchData();
}

/**
 * Handle card click event - opens JSON drawer
 *
 * @param {EventsJsExecResponse} item - Clicked log item
 */
function cardClick(item: EventsJsExecResponse): void {
  selectedEvent.value = item;
  jsonDrawerOpen.value = true;
}

/**
 * Get color for success status
 *
 * @param {boolean} success - Success status
 * @returns {string} Color class name
 */
function getSuccessColor(success: boolean): string {
  return success ? SUCCESS_COLORS.success : SUCCESS_COLORS.failed;
}

/**
 * Format timestamp for display
 *
 * @param {string} value - ISO timestamp string
 * @returns {string} Formatted date string
 */
function formatTimestamp(value: string): string {
  if (!value) return t.defaults.notAvailable.value;
  const date = new Date(value);
  return date.toLocaleDateString('en-US', {
    day: '2-digit',
    month: 'short',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  });
}

/** LIFECYCLE HOOKS */
onMounted(() => {
  void fetchData();
});
</script>

<template>
  <q-page class="q-pt-lg">

    <!-- Header Section -->
    <PageHeader
      icon="code"
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
          v-model="quickSearchUuid"
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
          v-model="quickStatusSuccess"
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
          <q-icon name="code" size="sm" color="primary" class="q-mr-sm"/>
          <div class="text-subtitle1 text-weight-medium text-primary">{{ t.listTitle.value }}</div>
        </div>
      </div>
      <div class="col-auto">
        <ListHeaderMenu
          :items-count="logsList.length"
          :item-label="t.itemLabel.value"
          :item-label-plural="t.itemLabelPlural.value"
          icon="code"
          :items-per-page="limit"
          :columns="menuColumns"
          :hide-total="true"
          :filtered="hasActiveFilters"
          :refreshing="loading"
          :last-updated-at="lastUpdatedAt"
          @update:items-per-page="handleLimitChange"
          @update:columns="handleColumnsUpdate"
          @refresh="refreshData"
        />
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="row justify-center q-pa-xl">
      <q-spinner color="primary" size="50px" />
    </div>

    <!-- Logs Row List -->
    <div v-else-if="logsList.length > 0" class="row">
      <div
        v-for="(log, index) in logsList"
        class="col-12 q-mb-xs"
        :key="log.threadId + '-' + log.created + '-' + index"
      >
        <DataRow
          :data="log"
          :columns="visibleColumns"
          @click="cardClick"
        />
      </div>
    </div>

    <!-- No Results -->
    <div v-else class="row q-col-gutter-lg">
      <ListCardEmpty
        :title="t.empty.title.value"
        :description="t.empty.description.value"
        icon="code"
      />
    </div>

    <!-- Cursor Pagination -->
    <div v-if="logsList.length > 0" class="row justify-center q-mt-lg q-mb-md">
      <q-btn
        flat
        color="primary"
        icon="chevron_left"
        :label="t.pagination.newer.value"
        :disable="!cursor.hasPrevious || loading"
        @click="goToPrevPage"
        class="q-mr-sm"
      />
      <q-btn
        flat
        color="primary"
        icon-right="chevron_right"
        :label="t.pagination.older.value"
        :disable="!cursor.hasNext || loading"
        @click="goToNextPage"
      />
    </div>

    <!-- JSON Drawer -->
    <JsonDrawer
      v-if="selectedEventParsed"
      v-model:show="jsonDrawerOpen"
      :jsonData="selectedEventParsed"
      :editable="false"
      :title="t.drawer.title.value"
      :subtitle="`${selectedEventParsed.success ? t.drawer.subtitleSuccess.value : t.drawer.subtitleFailed.value} • ${formatTimestamp(selectedEventParsed.created)}`"
    />

    <!-- Advanced Filters Drawer -->
    <AdvancedFiltersDrawer
      v-model="showFiltersDrawer"
      :fields="advancedFilterFields"
      :values="advancedFilterValues"
      @apply="handleAdvancedFiltersApply"
      @reset="handleAdvancedFiltersReset"
      @pending-change="handlePendingChange"
      @field-change="handleAdvancedFieldChange"
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
