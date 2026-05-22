<script setup lang="ts">
defineOptions({
  name: 'TriggerLogsPage'
});

/** TYPE IMPORTS */
import type { ListHeaderMenuColumn } from '@components/headers';
import type { DataRowColumn } from '@components/cards';
import type { EventsTriggerResponse } from '@mapexos/schemas';
import type { FilterField, FilterValues } from '@components/drawers';

/** VUE IMPORTS */
import { ref, computed, onMounted } from 'vue';

/** COMPONENTS */
import { PageHeader, ListHeaderMenu } from '@components/headers';
import { ListCardEmpty, DataRow } from '@components/cards';
import { JsonDrawer, AdvancedFiltersDrawer } from '@components/drawers';
import { AppTooltip } from '@components/tooltips';

/** COMPOSABLES */
import { useOrgChangeRefresh } from '@composables/organizations';
import { useTriggerLogsPageTranslations } from '@composables/i18n/pages/logs/triggerLogsPage';
import { useLogger } from '@composables/useLogger';

/** UTILS */
import { notifyFail } from '@utils/alert';

/** SERVICES */
import { apis } from '@services/mapex';

/** CONSTANTS */
const DEFAULT_LIMIT = 15;
const MAX_VISIBLE_CHIPS = 2;

const STATUS_COLORS = {
  success: 'green-6',
  failure: 'red-6',
};

const TRIGGER_TYPE_ICONS: Record<string, string> = {
  http: 'http',
  mqtt: 'sensors',
  rabbitmq: 'message',
  nats: 'hub',
  websocket: 'cable',
  email: 'email',
  teams: 'groups',
  slack: 'chat',
};

const TRIGGER_TYPE_COLORS: Record<string, string> = {
  http: 'blue-6',
  mqtt: 'teal-6',
  rabbitmq: 'orange-6',
  nats: 'purple-6',
  websocket: 'cyan-6',
  email: 'amber-8',
  teams: 'indigo-6',
  slack: 'green-7',
};

const SOURCE_COLORS: Record<string, string> = {
  router: 'primary',
};

const DURATION_THRESHOLDS = {
  fast: 100,
  normal: 500,
};

const DURATION_COLORS = {
  fast: 'green-6',
  normal: 'orange-6',
  slow: 'red-6',
};

/** LOCAL IMPORTS */
import type { TriggerLogsPageCursor, TriggerLogsPageFilters } from './interfaces/TriggerLogsPage.interface';

/** COMPOSABLES & STORES */
const t = useTriggerLogsPageTranslations();
const logger = useLogger('TriggerLogsPage');

/** STATE */
const loading = ref(false);
const lastUpdatedAt = ref<number | undefined>(undefined);
const logsList = ref<EventsTriggerResponse[]>([]);
const selectedEvent = ref<EventsTriggerResponse | null>(null);
const jsonDrawerOpen = ref(false);

/** Cursor Pagination */
const limit = ref(DEFAULT_LIMIT);
const cursor = ref<TriggerLogsPageCursor>({
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
const filters = ref<TriggerLogsPageFilters>({
  startTime: todayRange.startTime,
  endTime: todayRange.endTime,
});

/** FILTER STATE - Enterprise Filter Pattern */
const showFiltersDrawer = ref(false);
const quickSearchTriggerId = ref('');
const quickStatusSuccess = ref<boolean | null>(null);
const advancedFilterValues = ref<FilterValues>({
  includeChildren: null,
  triggerType: null,
  category: null,
  source: null,
});
const hasPendingAdvancedFilters = ref(false);

/** Column visibility using ListHeaderMenuColumn format */
const menuColumns = ref<ListHeaderMenuColumn[]>([
  { key: 'triggerType', label: 'Trigger Type', visible: true },
  { key: 'category', label: 'Category', visible: true },
  { key: 'source', label: 'Source', visible: true },
  { key: 'duration', label: 'Duration', visible: true },
  { key: 'created', label: 'Timestamp', visible: true },
]);

/** COMPUTED */

/**
 * Status options for quick filter
 */
const statusOptions = computed(() => [
  { label: t.filters.allStatus.value, value: null },
  { label: t.filters.options.success.value, value: true },
  { label: t.filters.options.failure.value, value: false },
]);

/**
 * Advanced filter fields configuration
 * includeChildren is FIRST as per Enterprise Filter Pattern
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
    key: 'triggerType',
    type: 'select',
    label: t.filters.triggerType.value,
    icon: 'flash_on',
    options: [
      { label: t.filters.allStatus.value, value: null },
      { label: t.triggerTypeOptions.http.value, value: 'http' },
      { label: t.triggerTypeOptions.mqtt.value, value: 'mqtt' },
      { label: t.triggerTypeOptions.rabbitmq.value, value: 'rabbitmq' },
      { label: t.triggerTypeOptions.nats.value, value: 'nats' },
      { label: t.triggerTypeOptions.websocket.value, value: 'websocket' },
      { label: t.triggerTypeOptions.email.value, value: 'email' },
      { label: t.triggerTypeOptions.teams.value, value: 'teams' },
      { label: t.triggerTypeOptions.slack.value, value: 'slack' },
    ],
  },
  {
    key: 'category',
    type: 'select',
    label: t.filters.category.value,
    icon: 'category',
    options: [
      { label: t.filters.allStatus.value, value: null },
      { label: t.categoryOptions.technical.value, value: 'technical' },
      { label: t.categoryOptions.communication.value, value: 'communication' },
    ],
  },
  {
    key: 'source',
    type: 'select',
    label: t.filters.source.value,
    icon: 'source',
    options: [
      { label: t.filters.allStatus.value, value: null },
      { label: t.sourceOptions.router.value, value: 'router' },
    ],
  },
]);

/**
 * Check if any filter is active
 */
const hasActiveFilters = computed(() => {
  return !!(
    filters.value.triggerId ||
    filters.value.success !== undefined ||
    filters.value.includeChildren !== undefined ||
    filters.value.triggerType ||
    filters.value.category ||
    filters.value.source
  );
});

/**
 * Count of advanced filters only (for badge on filter icon)
 */
const advancedFiltersCount = computed(() => {
  let count = 0;
  if (filters.value.includeChildren !== undefined) count++;
  if (filters.value.triggerType) count++;
  if (filters.value.category) count++;
  if (filters.value.source) count++;
  return count;
});

/**
 * Active filter chips for display
 */
const activeFilterChips = computed(() => {
  const chips: Array<{ key: string; label: string; value: string }> = [];

  if (filters.value.triggerId) {
    chips.push({ key: 'triggerId', label: t.filters.triggerId.value, value: filters.value.triggerId });
  }
  if (filters.value.success !== undefined) {
    chips.push({
      key: 'success',
      label: t.filters.status.value,
      value: filters.value.success ? t.filters.options.success.value : t.filters.options.failure.value
    });
  }
  if (filters.value.includeChildren !== undefined) {
    chips.push({
      key: 'includeChildren',
      label: t.filters.includeChildren.value,
      value: filters.value.includeChildren ? t.filters.options.yes.value : t.filters.options.no.value
    });
  }
  if (filters.value.triggerType) {
    const typeLabel = t.triggerTypeOptions[filters.value.triggerType as keyof typeof t.triggerTypeOptions]?.value || filters.value.triggerType;
    chips.push({ key: 'triggerType', label: t.filters.triggerType.value, value: typeLabel });
  }
  if (filters.value.category) {
    const categoryLabel = t.categoryOptions[filters.value.category as keyof typeof t.categoryOptions]?.value || filters.value.category;
    chips.push({ key: 'category', label: t.filters.category.value, value: categoryLabel });
  }
  if (filters.value.source) {
    const srcLabel = t.sourceOptions[filters.value.source as keyof typeof t.sourceOptions]?.value || filters.value.source;
    chips.push({ key: 'source', label: t.filters.source.value, value: srcLabel });
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
 * Column visibility state computed from menuColumns
 */
const columnVisibility = computed(() => ({
  triggerType: menuColumns.value.find((col) => col.key === 'triggerType')?.visible ?? true,
  category: menuColumns.value.find((col) => col.key === 'category')?.visible ?? true,
  source: menuColumns.value.find((col) => col.key === 'source')?.visible ?? true,
  duration: menuColumns.value.find((col) => col.key === 'duration')?.visible ?? true,
  created: menuColumns.value.find((col) => col.key === 'created')?.visible ?? true,
}));

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
    icon: (value: any, row: EventsTriggerResponse) => getTriggerTypeIcon(row.triggerType),
    color: (value: any, row: EventsTriggerResponse) => row.success ? STATUS_COLORS.success : STATUS_COLORS.failure,
  },
  {
    key: 'triggerName',
    label: 'Trigger Name',
    type: 'text',
    visible: 'always',
    width: 250,
    ellipsis: true,
    format: (value: string, row: EventsTriggerResponse) => value || row.triggerId || 'Unknown',
    secondaryKey: 'error',
  },
  {
    key: 'triggerType',
    label: 'Type',
    type: 'chip',
    visible: 'laptop',
    width: 100,
    format: (value: string) => value?.toUpperCase() || 'N/A',
    color: (value: string) => getTriggerTypeColor(value),
  },
  {
    key: 'category',
    label: 'Category',
    type: 'chip',
    visible: 'laptop',
    width: 130,
    format: (value: string) => value ? value.charAt(0).toUpperCase() + value.slice(1) : 'N/A',
    color: () => 'grey-7',
  },
  {
    key: 'source',
    label: 'Source',
    type: 'badge',
    visible: 'laptop',
    width: 110,
    color: (value: string) => getSourceColor(value),
  },
  {
    key: 'success',
    label: 'Status',
    type: 'badge',
    visible: 'laptop',
    width: 100,
    format: (value: boolean) => value ? 'SUCCESS' : 'FAILURE',
    color: (value: boolean) => value ? STATUS_COLORS.success : STATUS_COLORS.failure,
  },
  {
    key: 'durationMs',
    label: 'Duration',
    type: 'badge',
    visible: 'laptop',
    width: 80,
    format: (value: number) => `${value ?? 0}ms`,
    color: (value: number) => getDurationColor(value),
  },
  {
    key: 'created',
    label: 'Timestamp',
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
    if (col.key === 'icon' || col.key === 'triggerName' || col.key === 'success') return true;
    if (col.key === 'triggerType') return columnVisibility.value.triggerType;
    if (col.key === 'category') return columnVisibility.value.category;
    if (col.key === 'source') return columnVisibility.value.source;
    if (col.key === 'durationMs') return columnVisibility.value.duration;
    if (col.key === 'created') return columnVisibility.value.created;
    return true;
  });
});

/**
 * Selected event with parsed fields for JSON drawer
 */
const selectedEventParsed = computed(() => {
  if (!selectedEvent.value) return null;

  const eventData = { ...selectedEvent.value } as any;

  // Parse JSON string fields
  const jsonFields = ['requestData', 'responseData'];
  for (const field of jsonFields) {
    if (typeof eventData[field] === 'string' && eventData[field]) {
      try {
        eventData[field] = JSON.parse(eventData[field]);
      } catch (e) {
        logger.warn(`Failed to parse ${field}`, e);
      }
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
 * Fetch trigger events from API with current filters and cursor pagination
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
    if (filters.value.triggerId) {
      queryParams.triggerId = filters.value.triggerId;
    }
    if (filters.value.triggerType) {
      queryParams.triggerType = filters.value.triggerType;
    }
    if (filters.value.category) {
      queryParams.category = filters.value.category;
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

    const response = await apis.events.events.listTrigger(queryParams);

    logger.debug('Response received', response);

    logsList.value = response.items || [];

    // Update cursor state - only include defined values
    const newCursor: TriggerLogsPageCursor = {
      hasNext: response.hasNext,
      hasPrevious: response.hasPrevious,
    };
    if (queryParams.cursor) newCursor.current = queryParams.cursor;
    if (response.nextCursor) newCursor.next = response.nextCursor;
    if (response.prevCursor) newCursor.prev = response.prevCursor;
    cursor.value = newCursor;

  } catch (error: any) {
    logger.error('Error fetching trigger events', error);
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
  if (quickSearchTriggerId.value) {
    filters.value.triggerId = quickSearchTriggerId.value;
  } else {
    delete filters.value.triggerId;
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
  if (values.triggerType) {
    filters.value.triggerType = values.triggerType;
  } else {
    delete filters.value.triggerType;
  }
  if (values.category) {
    filters.value.category = values.category;
  } else {
    delete filters.value.category;
  }
  if (values.source) {
    filters.value.source = values.source;
  } else {
    delete filters.value.source;
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
    triggerType: null,
    category: null,
    source: null,
  };
}

/**
 * Remove individual filter
 * @param {string} key - Filter key to remove
 */
function removeFilter(key: string): void {
  if (key === 'triggerId') {
    delete filters.value.triggerId;
    quickSearchTriggerId.value = '';
  } else if (key === 'success') {
    delete filters.value.success;
    quickStatusSuccess.value = null;
  } else if (key === 'includeChildren') {
    delete filters.value.includeChildren;
    advancedFilterValues.value.includeChildren = null;
  } else if (key === 'triggerType') {
    delete filters.value.triggerType;
    advancedFilterValues.value.triggerType = null;
  } else if (key === 'category') {
    delete filters.value.category;
    advancedFilterValues.value.category = null;
  } else if (key === 'source') {
    delete filters.value.source;
    advancedFilterValues.value.source = null;
  }

  resetCursor();
  void fetchData();
}

/**
 * Clear all filters
 */
function clearAllFilters(): void {
  // Reset filter state (keep date range)
  filters.value = {
    startTime: todayRange.startTime,
    endTime: todayRange.endTime,
  };

  // Reset quick filters
  quickSearchTriggerId.value = '';
  quickStatusSuccess.value = null;

  // Reset advanced filters
  advancedFilterValues.value = {
    includeChildren: null,
    triggerType: null,
    category: null,
    source: null,
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
 * @param {EventsTriggerResponse} item - Clicked log item
 */
function cardClick(item: EventsTriggerResponse): void {
  selectedEvent.value = item;
  jsonDrawerOpen.value = true;
}

/**
 * Get icon for trigger type
 *
 * @param {string} type - Trigger type
 * @returns {string} Icon name
 */
function getTriggerTypeIcon(type: string): string {
  return TRIGGER_TYPE_ICONS[type?.toLowerCase()] || 'flash_on';
}

/**
 * Get color for trigger type chip
 *
 * @param {string} type - Trigger type
 * @returns {string} Color class name
 */
function getTriggerTypeColor(type: string): string {
  return TRIGGER_TYPE_COLORS[type?.toLowerCase()] || 'grey-6';
}

/**
 * Get color for source badge
 *
 * @returns {string} Color class name
 */
function getSourceColor(source: string): string {
  return SOURCE_COLORS[source?.toLowerCase()] || 'grey-7';
}

/**
 * Get color for duration badge based on milliseconds
 *
 * @param {number} durationMs - Duration in milliseconds
 * @returns {string} Color class name
 */
function getDurationColor(durationMs: number): string {
  if (durationMs < DURATION_THRESHOLDS.fast) return DURATION_COLORS.fast;
  if (durationMs < DURATION_THRESHOLDS.normal) return DURATION_COLORS.normal;
  return DURATION_COLORS.slow;
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
      icon="flash_on"
      iconColor="primary"
      :title="t.pageHeader.title.value"
      :description="t.pageHeader.description.value"
    />

    <!-- Filters Section - Enterprise Filter Pattern -->
    <div class="text-caption text-grey-7 q-mb-xs">{{ t.filters.label.value }}</div>
    <div class="row items-center q-col-gutter-sm q-mb-md">
      <!-- Search Input -->
      <div class="col">
        <q-input
          v-model="quickSearchTriggerId"
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
          <q-icon name="flash_on" size="sm" color="primary" class="q-mr-sm"/>
          <div class="text-subtitle1 text-weight-medium text-primary">{{ t.listTitle.value }}</div>
        </div>
      </div>
      <div class="col-auto">
        <ListHeaderMenu
          :item-label="t.itemLabel.value"
          :item-label-plural="t.itemLabelPlural.value"
          icon="flash_on"
          :items-count="logsList.length"
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
        :key="log.triggerId + '-' + log.created + '-' + index"
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
        icon="flash_on"
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
