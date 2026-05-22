<script setup lang="ts">
defineOptions({
  name: 'AssetRawLogsPage'
});

/** TYPE IMPORTS */
import type { ListHeaderMenuColumn } from '@components/headers';
import type { DataRowColumn, DataRowActionConfig } from '@components/cards';
import type { EventsRawResponse } from '@mapexos/schemas';
import type { AssetRawLogsPageFilters, AssetRawLogsPageColumnVisibility, AssetRawLogsPageCursor } from './interfaces';

/** VUE IMPORTS */
import { ref, computed, onMounted } from 'vue';
import { useRouter } from 'vue-router';

/** COMPONENTS */
import { PageHeader, ListHeaderMenu } from '@components/headers';
import { ListCardEmpty, DataRow } from '@components/cards';
import { JsonDrawer, AdvancedFiltersDrawer } from '@components/drawers';
import { AppTooltip } from '@components/tooltips';

/** COMPOSABLES */
import { useOrgChangeRefresh } from '@composables/organizations';
import { useAssetRawLogsPageTranslations } from '@composables/i18n/pages/logs/assetRawLogsPage';
import { useLogger } from '@composables/useLogger';

/** UTILS */
import { notifyFail } from '@utils/alert';

/** SERVICES */
import { apis } from '@services/mapex';

/** LOCAL IMPORTS */
import {
  DEFAULT_LIMIT,
  COLUMN_VISIBILITY_DEFAULTS,
  SOURCE_COLORS,
  SUCCESS_COLORS,
  DEFAULT_COLOR,
} from './constants';

/** CONSTANTS */
const MAX_VISIBLE_CHIPS = 2;

/** COMPOSABLES */
const t = useAssetRawLogsPageTranslations();
const logger = useLogger('AssetRawLogsPage');
const router = useRouter();

/** STATE */
const loading = ref(false);
const lastUpdatedAt = ref<number | undefined>(undefined);
const logsList = ref<EventsRawResponse[]>([]);
const selectedEvent = ref<EventsRawResponse | null>(null);
const jsonDrawerOpen = ref(false);

/** Cursor Pagination */
const limit = ref(DEFAULT_LIMIT);
const cursor = ref<AssetRawLogsPageCursor>({
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
const filters = ref<AssetRawLogsPageFilters>({
  startTime: todayRange.startTime,
  endTime: todayRange.endTime,
});

/** Quick filters state */
const quickSearchUuid = ref('');
const quickStatus = ref<boolean | null>(null);

/** Advanced filters state */
const showAdvancedFilters = ref(false);
const advancedFilterValues = ref<Record<string, any>>({
  includeChildren: null,
  startDate: todayRange.startTime.split('T')[0],
  endDate: todayRange.endTime.split('T')[0],
  source: null,
  threadId: '',
});
const hasPendingAdvancedFilters = ref(false);

/** Column visibility */
const menuColumns = ref<ListHeaderMenuColumn[]>([
  { key: 'name', label: t.columns.dataSource.value, visible: COLUMN_VISIBILITY_DEFAULTS.name },
  { key: 'threadId', label: t.columns.uuid.value, visible: COLUMN_VISIBILITY_DEFAULTS.threadId },
  { key: 'source', label: t.columns.source.value, visible: COLUMN_VISIBILITY_DEFAULTS.source },
  { key: 'created', label: t.columns.timestamp.value, visible: COLUMN_VISIBILITY_DEFAULTS.created },
]);

/** Row actions configuration */
const rowActionsConfig = computed<DataRowActionConfig>(() => ({
  showEdit: false,
  showView: true,
  showDelete: false,
  customActions: [
    {
      key: 'trackEvent',
      label: t.actions.trackEvent.value,
      icon: 'account_tree',
      color: 'secondary',
    },
  ],
}));

/** COMPUTED */

/**
 * Status select options for quick filter
 */
const statusOptions = computed(() => [
  { label: t.filters.allStatus.value, value: null },
  { label: t.filters.options.success.value, value: true },
  { label: t.filters.options.failed.value, value: false },
]);

/**
 * Count of active advanced filters
 */
const advancedFiltersCount = computed(() => {
  let count = 0;
  if (advancedFilterValues.value.includeChildren !== null) count++;
  if (advancedFilterValues.value.source) count++;
  if (advancedFilterValues.value.threadId) count++;
  return count;
});

/**
 * Active filter chips for display
 */
const activeFilterChips = computed(() => {
  const chips: Array<{ key: string; label: string; value: string }> = [];

  if (quickSearchUuid.value) {
    chips.push({
      key: 'quickSearchUuid',
      label: t.filters.uuid.value,
      value: quickSearchUuid.value,
    });
  }

  if (quickStatus.value !== null) {
    const statusLabel = quickStatus.value
      ? t.filters.options.success.value
      : t.filters.options.failed.value;
    chips.push({
      key: 'success',
      label: t.filters.status.value,
      value: statusLabel,
    });
  }

  if (advancedFilterValues.value.includeChildren !== null) {
    const childrenLabel = advancedFilterValues.value.includeChildren
      ? t.filters.options.yes.value
      : t.filters.options.no.value;
    chips.push({
      key: 'includeChildren',
      label: t.filters.includeChildrenOrgs.value,
      value: childrenLabel,
    });
  }

  if (advancedFilterValues.value.source) {
    const sourceLabel = formatSourceLabel(advancedFilterValues.value.source);
    chips.push({
      key: 'source',
      label: t.filters.source.value,
      value: sourceLabel,
    });
  }

  if (advancedFilterValues.value.threadId) {
    chips.push({
      key: 'threadId',
      label: t.filters.uuid.value,
      value: advancedFilterValues.value.threadId,
    });
  }

  return chips;
});

/**
 * Visible filter chips (limited by MAX_VISIBLE_CHIPS)
 */
const visibleFilterChips = computed(() => activeFilterChips.value.slice(0, MAX_VISIBLE_CHIPS));

/**
 * Hidden filter chips (beyond MAX_VISIBLE_CHIPS)
 */
const hiddenFilterChips = computed(() => activeFilterChips.value.slice(MAX_VISIBLE_CHIPS));

/**
 * Count of hidden filters
 */
const hiddenFiltersCount = computed(() => hiddenFilterChips.value.length);

/**
 * Column visibility state computed from menuColumns
 */
const columnVisibility = computed<AssetRawLogsPageColumnVisibility>(() => ({
  name: menuColumns.value.find((col) => col.key === 'name')?.visible ?? true,
  threadId: menuColumns.value.find((col) => col.key === 'threadId')?.visible ?? true,
  source: menuColumns.value.find((col) => col.key === 'source')?.visible ?? true,
  created: menuColumns.value.find((col) => col.key === 'created')?.visible ?? true,
}));

/**
 * Source options computed from translations
 */
const sourceOptions = computed(() => [
  { label: t.sourceOptions.httpGateway.value, value: 'http_gateway' },
  { label: t.sourceOptions.mqttGateway.value, value: 'mqtt_gateway' },
  { label: t.sourceOptions.lorawanGateway.value, value: 'lorawan_gateway', disable: true },
]);

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
    icon: () => 'terminal',
    color: (value: any, row: EventsRawResponse) => getSuccessColor(row.success),
  },
  {
    key: 'name',
    label: t.columns.dataSource.value,
    type: 'text',
    visible: 'always',
    width: 250,
    ellipsis: true,
    format: (value: string, row: EventsRawResponse) => value || row.threadId || t.defaults.unknown.value,
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
    key: 'source',
    label: t.columns.source.value,
    type: 'chip',
    visible: 'laptop',
    width: 150,
    format: (value: string) => formatSourceLabel(value),
    color: (value: string) => getSourceColor(value),
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
    if (col.key === 'source') return columnVisibility.value.source;
    if (col.key === 'created') return columnVisibility.value.created;
    return true;
  });
});

/** WATCHERS */

/**
 * Re-fetch data when organization changes
 */
useOrgChangeRefresh(() => {
  resetCursor();
  quickSearchUuid.value = '';
  quickStatus.value = null;
  const resetRange = getTodayDateRange();
  advancedFilterValues.value = {
    includeChildren: null,
    startDate: resetRange.startTime.split('T')[0],
    endDate: resetRange.endTime.split('T')[0],
    source: null,
    threadId: '',
  };
  filters.value = {
    startTime: resetRange.startTime,
    endTime: resetRange.endTime,
  };
  hasPendingAdvancedFilters.value = false;
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
 * Fetch raw events from API with current filters and cursor pagination
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
    if (filters.value.source) {
      queryParams.source = filters.value.source;
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

    const response = await apis.events.events.listRaw(queryParams);

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
    logger.error('Error fetching raw events', {
      message: error?.message,
      details: error?.details,
      error,
    });
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
 * Apply quick filters immediately
 * @returns {void}
 */
function applyQuickFilters(): void {
  if (quickSearchUuid.value) {
    filters.value.threadId = quickSearchUuid.value;
  } else {
    delete filters.value.threadId;
  }
  if (quickStatus.value !== null) {
    filters.value.success = quickStatus.value;
  } else {
    delete filters.value.success;
  }
  resetCursor();
  void fetchData();
}

/**
 * Handle advanced filters apply from drawer
 * @param {Record<string, any>} appliedFilters - Applied filter values
 * @returns {void}
 */
function handleAdvancedFiltersApply(appliedFilters: Record<string, any>): void {
  advancedFilterValues.value = { ...appliedFilters };
  if (appliedFilters.includeChildren !== null && appliedFilters.includeChildren !== undefined) {
    filters.value.includeChildren = appliedFilters.includeChildren;
  } else {
    delete filters.value.includeChildren;
  }
  // Handle date filters - convert from date string to ISO datetime
  if (appliedFilters.startDate) {
    const startDate = new Date(appliedFilters.startDate);
    startDate.setHours(0, 0, 0, 0);
    filters.value.startTime = startDate.toISOString();
  } else {
    // Default to today's start
    const today = getTodayDateRange();
    filters.value.startTime = today.startTime;
  }
  if (appliedFilters.endDate) {
    const endDate = new Date(appliedFilters.endDate);
    endDate.setHours(23, 59, 59, 999);
    filters.value.endTime = endDate.toISOString();
  } else {
    // Default to today's end
    const today = getTodayDateRange();
    filters.value.endTime = today.endTime;
  }
  if (appliedFilters.source) {
    filters.value.source = appliedFilters.source;
  } else {
    delete filters.value.source;
  }
  if (appliedFilters.threadId || quickSearchUuid.value) {
    filters.value.threadId = appliedFilters.threadId || quickSearchUuid.value;
  } else {
    delete filters.value.threadId;
  }
  hasPendingAdvancedFilters.value = false;
  resetCursor();
  showAdvancedFilters.value = false;
  void fetchData();
}

/**
 * Handle pending change in advanced filters
 * @param {boolean} hasPending - Whether there are pending changes
 * @returns {void}
 */
function handlePendingChange(hasPending: boolean): void {
  hasPendingAdvancedFilters.value = hasPending;
}

/**
 * Remove a specific filter
 * @param {string} key - Filter key to remove
 * @returns {void}
 */
function removeFilter(key: string): void {
  if (key === 'quickSearchUuid') {
    quickSearchUuid.value = '';
    if (advancedFilterValues.value.threadId) {
      filters.value.threadId = advancedFilterValues.value.threadId;
    } else {
      delete filters.value.threadId;
    }
  } else if (key === 'success') {
    quickStatus.value = null;
    delete filters.value.success;
  } else if (key === 'includeChildren') {
    advancedFilterValues.value.includeChildren = null;
    delete filters.value.includeChildren;
  } else if (key === 'source') {
    advancedFilterValues.value.source = null;
    delete filters.value.source;
  } else if (key === 'threadId') {
    advancedFilterValues.value.threadId = '';
    if (quickSearchUuid.value) {
      filters.value.threadId = quickSearchUuid.value;
    } else {
      delete filters.value.threadId;
    }
  }
  resetCursor();
  void fetchData();
}

/**
 * Clear all filters
 * @returns {void}
 */
function clearAllFilters(): void {
  quickSearchUuid.value = '';
  quickStatus.value = null;
  const resetRange = getTodayDateRange();
  advancedFilterValues.value = {
    includeChildren: null,
    startDate: resetRange.startTime.split('T')[0],
    endDate: resetRange.endTime.split('T')[0],
    source: null,
    threadId: '',
  };
  filters.value = {
    startTime: resetRange.startTime,
    endTime: resetRange.endTime,
  };
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
 * @param {EventsRawResponse} item - Clicked log item
 */
function cardClick(item: EventsRawResponse): void {
  selectedEvent.value = item;
  jsonDrawerOpen.value = true;
}

/**
 * Handle view action - opens JSON drawer
 *
 * @param {EventsRawResponse} item - Log item to view
 */
function handleView(item: EventsRawResponse): void {
  selectedEvent.value = item;
  jsonDrawerOpen.value = true;
}

/**
 * Handle custom action event from DataRow
 *
 * @param {string} actionKey - The key of the action that was triggered
 * @param {EventsRawResponse} item - The log item associated with the action
 */
function handleAction(actionKey: string, item: EventsRawResponse): void {
  if (actionKey === 'trackEvent') {
    // Navigate to Event Tracer page with the eventTrackerId
    void router.push({
      path: '/logs/event_tracer',
      query: { eventTrackerId: item.eventTrackerId },
    });
  }
}

/**
 * Get color for source type
 *
 * @param {string} source - Source type
 * @returns {string} Color class name
 */
function getSourceColor(source: string): string {
  return SOURCE_COLORS[source] || DEFAULT_COLOR;
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
 * Format source label for display
 *
 * @param {string} source - Source type (e.g., 'http_gateway')
 * @returns {string} Formatted label (e.g., 'HTTP Gateway')
 */
function formatSourceLabel(source: string): string {
  const option = sourceOptions.value.find((opt) => opt.value === source);
  return option?.label || source;
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
      icon="terminal"
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
          @click="showAdvancedFilters = true"
        >
          <q-badge
            v-if="advancedFiltersCount > 0 || hasPendingAdvancedFilters"
            :color="hasPendingAdvancedFilters ? 'warning' : 'primary'"
            floating
            rounded
            :label="hasPendingAdvancedFilters ? '!' : advancedFiltersCount"
          />
          <AppTooltip :content="hasPendingAdvancedFilters
            ? t.filters.pendingFilters.value
            : t.filters.advancedFilters.value"
          />
        </q-btn>
      </div>
    </div>

    <!-- Active Filter Chips -->
    <div v-if="activeFilterChips.length > 0" class="row items-center q-mb-md q-gutter-xs">
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
          <q-icon name="terminal" size="sm" color="primary" class="q-mr-sm"/>
          <div class="text-subtitle1 text-weight-medium text-primary">{{ t.listTitle.value }}</div>
        </div>
      </div>
      <div class="col-auto">
        <ListHeaderMenu
          :items-count="logsList.length"
          :item-label="t.itemLabel.value"
          :item-label-plural="t.itemLabelPlural.value"
          icon="terminal"
          :items-per-page="limit"
          :columns="menuColumns"
          :hide-total="true"
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
        :key="log.threadId + '-' + index"
      >
        <DataRow
          :data="log"
          :columns="visibleColumns"
          :actions="rowActionsConfig"
          @click="cardClick"
          @view="handleView"
          @action="handleAction"
        />
      </div>
    </div>

    <!-- No Results -->
    <div v-else class="row q-col-gutter-lg">
      <ListCardEmpty
        :title="t.empty.title.value"
        :description="t.empty.description.value"
        icon="terminal"
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
      v-if="selectedEvent"
      v-model:show="jsonDrawerOpen"
      :jsonData="selectedEvent"
      :editable="false"
      :title="t.drawer.title.value"
      :subtitle="`${formatSourceLabel(selectedEvent.source)} • ${formatTimestamp(selectedEvent.created)}`"
    />

    <!-- Advanced Filters Drawer -->
    <AdvancedFiltersDrawer
      v-model="showAdvancedFilters"
      :title="t.filters.advancedFilters.value"
      :fields="t.advancedFilters.value"
      :values="advancedFilterValues"
      @apply="handleAdvancedFiltersApply"
      @pending-change="handlePendingChange"
    />
  </q-page>
</template>

<style scoped lang="scss">
.filter-input {
  :deep(.q-field__control) {
    border-radius: var(--mapex-radius-md);
  }
}
</style>
