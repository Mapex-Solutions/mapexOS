<script setup lang="ts">
defineOptions({
  name: 'RouterLogsPage'
});

/** TYPE IMPORTS */
import type { ListHeaderMenuColumn } from '@components/headers';
import type { DataRowColumn } from '@components/cards';
import type { FilterField, FilterValues } from '@components/drawers';
import type { EventsRouterResponse } from '@mapexos/schemas';
import type { RouterLogsPageFilters, RouterLogsPageColumnVisibility, RouterLogsPageCursor } from './interfaces';

/** VUE IMPORTS */
import { ref, computed, onMounted } from 'vue';

/** COMPONENTS */
import { PageHeader, ListHeaderMenu } from '@components/headers';
import { ListCardEmpty, DataRow } from '@components/cards';
import { JsonDrawer, AdvancedFiltersDrawer } from '@components/drawers';
import { AppTooltip } from '@components/tooltips';

/** COMPOSABLES */
import { useOrgChangeRefresh } from '@composables/organizations';
import { useRouterLogsPageTranslations } from '@composables/i18n/pages/logs/routerLogsPage';
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
  COUNT_COLORS,
} from './constants';

/** TRANSLATIONS */
const t = useRouterLogsPageTranslations();
const logger = useLogger('RouterLogsPage');

/** STATE */
const loading = ref(false);
const lastUpdatedAt = ref<number | undefined>(undefined);
const logsList = ref<EventsRouterResponse[]>([]);
const selectedEvent = ref<EventsRouterResponse | null>(null);
const jsonDrawerOpen = ref(false);

/** Cursor Pagination */
const limit = ref(DEFAULT_LIMIT);
const cursor = ref<RouterLogsPageCursor>({
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
const filters = ref<RouterLogsPageFilters>({
  startTime: todayRange.startTime,
  endTime: todayRange.endTime,
});

/** FILTER STATE */
const showFiltersDrawer = ref(false);
const quickSearchThreadId = ref('');
const quickStatusSuccess = ref<boolean | null>(null);
const advancedFilterValues = ref<FilterValues>({
  includeChildren: null,
  assetId: null,
  routerId: null,
});
const hasPendingAdvancedFilters = ref(false);

/** Column visibility */
const menuColumns = ref<ListHeaderMenuColumn[]>([
  { key: 'name', label: t.columns.routeGroup.value, visible: COLUMN_VISIBILITY_DEFAULTS.name },
  { key: 'threadId', label: t.columns.uuid.value, visible: COLUMN_VISIBILITY_DEFAULTS.threadId },
  { key: 'success', label: t.columns.status.value, visible: COLUMN_VISIBILITY_DEFAULTS.success },
  { key: 'routersCount', label: t.columns.totalRouters.value, visible: COLUMN_VISIBILITY_DEFAULTS.routersCount },
  { key: 'publishedCount', label: t.columns.publishedCount.value, visible: COLUMN_VISIBILITY_DEFAULTS.publishedCount },
  { key: 'created', label: t.columns.timestamp.value, visible: COLUMN_VISIBILITY_DEFAULTS.created },
]);

/** COMPUTED */

/**
 * Column visibility state computed from menuColumns
 */
const columnVisibility = computed<RouterLogsPageColumnVisibility>(() => ({
  threadId: menuColumns.value.find((col) => col.key === 'threadId')?.visible ?? true,
  name: menuColumns.value.find((col) => col.key === 'name')?.visible ?? true,
  success: menuColumns.value.find((col) => col.key === 'success')?.visible ?? true,
  routersCount: menuColumns.value.find((col) => col.key === 'routersCount')?.visible ?? true,
  publishedCount: menuColumns.value.find((col) => col.key === 'publishedCount')?.visible ?? true,
  created: menuColumns.value.find((col) => col.key === 'created')?.visible ?? true,
}));

/**
 * Status options for quick filter
 */
const statusOptions = computed(() => [
  { label: t.filters.allStatus.value, value: null },
  { label: t.statusOptions.success.value, value: true },
  { label: t.statusOptions.failed.value, value: false },
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
    key: 'assetId',
    type: 'input',
    label: t.filters.assetId.value,
    icon: 'memory',
    placeholder: t.filters.assetIdPlaceholder.value,
  },
  {
    key: 'routerId',
    type: 'input',
    label: t.filters.routerId.value,
    icon: 'route',
    placeholder: t.filters.routerIdPlaceholder.value,
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
    filters.value.assetId ||
    filters.value.routerId
  );
});

/**
 * Count of advanced filters only (for badge on filter icon)
 */
const advancedFiltersCount = computed(() => {
  let count = 0;
  if (filters.value.includeChildren !== undefined) count++;
  if (filters.value.assetId) count++;
  if (filters.value.routerId) count++;
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
      value: filters.value.success ? t.statusOptions.success.value : t.statusOptions.failed.value
    });
  }
  if (filters.value.includeChildren !== undefined) {
    chips.push({
      key: 'includeChildren',
      label: t.filters.includeChildren.value,
      value: filters.value.includeChildren ? t.filters.options.yes.value : t.filters.options.no.value
    });
  }
  if (filters.value.assetId) {
    chips.push({ key: 'assetId', label: t.filters.assetId.value, value: filters.value.assetId });
  }
  if (filters.value.routerId) {
    chips.push({ key: 'routerId', label: t.filters.routerId.value, value: filters.value.routerId });
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
 * Log columns configuration for DataRow
 */
const logColumns = computed<DataRowColumn[]>(() => [
  {
    key: 'icon',
    label: '',
    type: 'avatar',
    visible: 'always',
    width: 56,
    icon: () => 'route',
    color: (value: any, row: EventsRouterResponse) => getSuccessColor(row.success),
  },
  {
    key: 'name',
    label: t.columns.routeGroup.value,
    type: 'text',
    visible: 'always',
    width: 250,
    ellipsis: true,
    format: (value: string, row: EventsRouterResponse) => value || row.routerId || t.defaults.unknown.value,
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
    key: 'totalRouters',
    label: t.columns.totalRouters.value,
    type: 'badge',
    visible: 'laptop',
    width: 80,
    format: (value: number) => `${value ?? 0}`,
    color: () => 'grey-7',
  },
  {
    key: 'matchedCount',
    label: t.columns.matchedCount.value,
    type: 'badge',
    visible: 'laptop',
    width: 80,
    format: (value: number, row: EventsRouterResponse) => `${value ?? 0}/${row.totalRouters ?? 0}`,
    color: (value: number, row: EventsRouterResponse) => getCountColor(value, row.totalRouters),
  },
  {
    key: 'publishedCount',
    label: t.columns.publishedCount.value,
    type: 'badge',
    visible: 'laptop',
    width: 80,
    format: (value: number, row: EventsRouterResponse) => `${value ?? 0}/${row.matchedCount ?? 0}`,
    color: (value: number, row: EventsRouterResponse) => getCountColor(value, row.matchedCount),
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
    if (col.key === 'totalRouters' || col.key === 'matchedCount') return columnVisibility.value.routersCount;
    if (col.key === 'publishedCount') return columnVisibility.value.publishedCount;
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
    hasNext: false,
    hasPrevious: false,
  };
}

/**
 * Fetch router events from API with current filters and cursor pagination
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
    if (filters.value.assetId) {
      queryParams.assetId = filters.value.assetId;
    }
    if (filters.value.routerId) {
      queryParams.routerId = filters.value.routerId;
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

    logger.debug('Fetching with params', queryParams);

    const response = await apis.events.events.listRouter(queryParams);

    logger.debug('Response received', response);

    logsList.value = response.items || [];

    // Update cursor state - only include defined values
    const newCursor: RouterLogsPageCursor = {
      hasNext: response.hasNext,
      hasPrevious: response.hasPrevious,
    };
    if (queryParams.cursor) newCursor.current = queryParams.cursor;
    if (response.nextCursor) newCursor.next = response.nextCursor;
    if (response.prevCursor) newCursor.prev = response.prevCursor;
    cursor.value = newCursor;

  } catch (error: any) {
    logger.error('Error fetching router events', error);
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
    hasNext: false,
    hasPrevious: false,
  };
  await fetchData('next');
}

/**
 * Apply quick filters (search + status)
 */
function applyQuickFilters(): void {
  if (quickSearchThreadId.value) {
    filters.value.threadId = quickSearchThreadId.value;
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
  if (values.assetId) {
    filters.value.assetId = values.assetId;
  } else {
    delete filters.value.assetId;
  }
  if (values.routerId) {
    filters.value.routerId = values.routerId;
  } else {
    delete filters.value.routerId;
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
  advancedFilterValues.value = {
    includeChildren: null,
    assetId: null,
    routerId: null,
  };
}

/**
 * Remove individual filter
 * @param {string} key - Filter key to remove
 */
function removeFilter(key: string): void {
  if (key === 'threadId') {
    delete filters.value.threadId;
    quickSearchThreadId.value = '';
  } else if (key === 'success') {
    delete filters.value.success;
    quickStatusSuccess.value = null;
  } else if (key === 'includeChildren') {
    delete filters.value.includeChildren;
    advancedFilterValues.value.includeChildren = null;
  } else if (key === 'assetId') {
    delete filters.value.assetId;
    advancedFilterValues.value.assetId = null;
  } else if (key === 'routerId') {
    delete filters.value.routerId;
    advancedFilterValues.value.routerId = null;
  }

  resetCursor();
  void fetchData();
}

/**
 * Clear all filters
 */
function clearAllFilters(): void {
  // Reset filter state (keep date range defaults)
  const today = getTodayDateRange();
  filters.value = {
    startTime: today.startTime,
    endTime: today.endTime,
  };

  // Reset quick filters
  quickSearchThreadId.value = '';
  quickStatusSuccess.value = null;

  // Reset advanced filters
  advancedFilterValues.value = {
    includeChildren: null,
    assetId: null,
    routerId: null,
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
 * @param {EventsRouterResponse} item - Clicked log item
 */
function cardClick(item: EventsRouterResponse): void {
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
 * Get color for count badges based on value vs total
 *
 * @param {number} value - Current count value
 * @param {number} total - Total count
 * @returns {string} Color class name
 */
function getCountColor(value: number, total: number): string {
  if (!value || value === 0) return COUNT_COLORS.zero;
  if (value < total) return COUNT_COLORS.partial;
  return COUNT_COLORS.full;
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
      icon="route"
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
          v-model="quickSearchThreadId"
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
          <q-icon name="route" size="sm" color="primary" class="q-mr-sm"/>
          <div class="text-subtitle1 text-weight-medium text-primary">{{ t.listTitle.value }}</div>
        </div>
      </div>
      <div class="col-auto">
        <ListHeaderMenu
          :items-count="logsList.length"
          :item-label="t.itemLabel.value"
          :item-label-plural="t.itemLabelPlural.value"
          icon="route"
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
        icon="route"
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
