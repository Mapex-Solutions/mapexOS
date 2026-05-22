<script setup lang="ts">
defineOptions({
  name: 'EventStorePage'
});

/** TYPE IMPORTS */
import type { ListHeaderMenuColumn } from '@components/headers';
import type { DataRowActionConfig } from '@components/cards';
import type {
  EventStorePageFilters,
  EventStorePageCursor,
  EventStoreItem,
} from './interfaces';
import type { DynamicFiltersResult } from '@components/drawers';
import type { PageTourStep } from '@composables/tour';

/** VUE IMPORTS */
import { ref, computed, onMounted } from 'vue';

/** COMPONENTS */
import { PageHeader, ListHeaderMenu } from '@components/headers';
import { ListCardEmpty, DataRow } from '@components/cards';
import { JsonDrawer, AdvancedFiltersDrawer, DynamicFiltersDrawer } from '@components/drawers';
import { AppTooltip } from '@components/tooltips';

/** COMPOSABLES */
import { useEventStorePageTranslations } from '@composables/i18n/pages/events/eventStorePage';
import { useOrgChangeRefresh } from '@composables/organizations';
import { usePageTour } from '@composables/tour';
import { useLogger } from '@composables/useLogger';

/** UTILS */
import { notifyFail } from '@utils/alert';

/** SERVICES */
import { apis } from '@services/mapex';

/** LOCAL IMPORTS */
import { DEFAULT_LIMIT, COLUMN_VISIBILITY_DEFAULTS, EVENT_STORE_TOUR_STEPS } from './constants';

/** CONSTANTS */
const MAX_VISIBLE_CHIPS = 2;

/** COMPOSABLES & STORES */
const t = useEventStorePageTranslations();
const logger = useLogger('EventStorePage');

/**
 * Build tour steps with resolved translations and drawer callbacks
 * Pattern: header → searchInput → advancedFiltersBtn → advancedFiltersOpen → dynamicFiltersBtn → dynamicFiltersOpen → results
 *
 * @returns {PageTourStep[]} Tour steps with resolved text and callbacks
 */
function buildTourSteps(): PageTourStep[] {
  return EVENT_STORE_TOUR_STEPS.map((step) => {
    const key = step.translationKey as keyof typeof t.tour;
    const translation = t.tour[key];
    const result: PageTourStep = {
      element: step.element,
      title: translation.title.value,
      description: translation.description.value,
    };
    if (step.side) result.side = step.side;
    if (step.align) result.align = step.align;

    // Advanced filters button: open drawer on Next click
    if (step.translationKey === 'advancedFiltersBtn') {
      result.onNextClick = (moveNext) => {
        showAdvancedFilters.value = true;
        setTimeout(moveNext, 400);
      };
    }

    // Advanced filters open: close drawer on Next click
    if (step.translationKey === 'advancedFiltersOpen') {
      result.onNextClick = (moveNext) => {
        showAdvancedFilters.value = false;
        setTimeout(moveNext, 300);
      };
    }

    // Dynamic filters button: open drawer with demo data on Next click
    if (step.translationKey === 'dynamicFiltersBtn') {
      result.onHighlightStarted = () => {
        showAdvancedFilters.value = false;
      };
      result.onNextClick = (moveNext) => {
        dynamicFiltersDemoMode.value = true;
        showDynamicFilters.value = true;
        setTimeout(moveNext, 400);
      };
    }

    // Dynamic filters open: close drawer and clear demo on Next click
    if (step.translationKey === 'dynamicFiltersOpen') {
      result.onNextClick = (moveNext) => {
        showDynamicFilters.value = false;
        dynamicFiltersDemoMode.value = false;
        setTimeout(moveNext, 300);
      };
    }

    // Results step: ensure drawers are closed and demo mode is off
    if (step.translationKey === 'results') {
      result.onHighlightStarted = () => {
        showAdvancedFilters.value = false;
        showDynamicFilters.value = false;
        dynamicFiltersDemoMode.value = false;
      };
    }

    return result;
  });
}

/** PAGE TOUR */
const { startTour } = usePageTour({
  tourId: 'event-store',
  steps: buildTourSteps,
});

/** STATE */
const loading = ref(false);
const lastUpdatedAt = ref<number | undefined>(undefined);
const eventsList = ref<EventStoreItem[]>([]);
const selectedEvent = ref<EventStoreItem | null>(null);
const jsonDrawerOpen = ref(false);

// Quick filters state
const quickSearchThreadId = ref('');

// Advanced filters state
const showAdvancedFilters = ref(false);
const advancedFilterValues = ref<Record<string, any>>({
  includeChildren: null,
  startDate: '',
  endDate: '',
  threadId: '',
  source: '',
});
const hasPendingAdvancedFilters = ref(false);

/** Dynamic filters state */
const showDynamicFilters = ref(false);
const dynamicFiltersResult = ref<DynamicFiltersResult | null>(null);
const hasPendingDynamicFilters = ref(false);
const dynamicFiltersDemoMode = ref(false);

/** Cursor Pagination */
const limit = ref(DEFAULT_LIMIT);
const cursor = ref<EventStorePageCursor>({
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

/** Fixed Filters - Default date range is today */
const todayRange = getTodayDateRange();
const filters = ref<EventStorePageFilters>({
  threadId: undefined,
  assetId: undefined,
  assetTemplateId: undefined,
  eventType: undefined,
  source: undefined,
  startTime: todayRange.startTime,
  endTime: todayRange.endTime,
});

/** Column visibility */
const menuColumns = ref<ListHeaderMenuColumn[]>([
  { key: 'threadId', label: t.menuColumns.threadId.value, visible: COLUMN_VISIBILITY_DEFAULTS.threadId },
  { key: 'source', label: t.menuColumns.source.value, visible: COLUMN_VISIBILITY_DEFAULTS.source },
  { key: 'created', label: t.menuColumns.created.value, visible: COLUMN_VISIBILITY_DEFAULTS.created },
]);

/** Row actions configuration */
const rowActionsConfig = computed<DataRowActionConfig>(() => ({
  showEdit: false,
  showView: true,
  showDelete: false,
}));

/** COMPUTED */

/**
 * Column visibility state computed from menuColumns
 */
const columnVisibility = computed(() => ({
  threadId: menuColumns.value.find((col) => col.key === 'threadId')?.visible ?? true,
  source: menuColumns.value.find((col) => col.key === 'source')?.visible ?? true,
  created: menuColumns.value.find((col) => col.key === 'created')?.visible ?? true,
}));

/**
 * Whether any filter is actively applied (advanced, dynamic, or quick search)
 */
const isFiltered = computed(() =>
  advancedFiltersCount.value > 0 ||
  !!dynamicFiltersResult.value ||
  !!quickSearchThreadId.value,
);

/**
 * Count of active advanced filters
 */
const advancedFiltersCount = computed(() => {
  let count = 0;
  if (advancedFilterValues.value.includeChildren !== null) count++;
  if (advancedFilterValues.value.startDate) count++;
  if (advancedFilterValues.value.endDate) count++;
  if (advancedFilterValues.value.threadId) count++;
  if (advancedFilterValues.value.source) count++;
  return count;
});

/**
 * Active filter chips for display
 */
const activeFilterChips = computed(() => {
  const chips: Array<{ key: string; label: string; value: string }> = [];

  if (quickSearchThreadId.value) {
    chips.push({
      key: 'threadId',
      label: t.filters.threadId.value,
      value: quickSearchThreadId.value,
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

  if (advancedFilterValues.value.startDate) {
    chips.push({
      key: 'startDate',
      label: t.filters.startDate.value,
      value: advancedFilterValues.value.startDate,
    });
  }

  if (advancedFilterValues.value.endDate) {
    chips.push({
      key: 'endDate',
      label: t.filters.endDate.value,
      value: advancedFilterValues.value.endDate,
    });
  }

  if (advancedFilterValues.value.threadId) {
    chips.push({
      key: 'threadId',
      label: t.filters.threadId.value,
      value: advancedFilterValues.value.threadId,
    });
  }

  if (dynamicFiltersResult.value) {
    chips.push({
      key: 'dynamicFilter',
      label: t.filters.assetTemplate.value,
      value: `${dynamicFiltersResult.value.sourceName} → ${dynamicFiltersResult.value.templateName}`,
    });
  }

  if (advancedFilterValues.value.source) {
    chips.push({
      key: 'source',
      label: t.filters.source.value,
      value: advancedFilterValues.value.source,
    });
  }

  return chips;
});

/**
 * Visible filter chips (limited to MAX_VISIBLE_CHIPS)
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
 * Filtered columns based on visibility settings
 */
const visibleColumns = computed(() => {
  return t.columns.value.filter((col: any) => {
    if (col.key === 'icon') return true;
    if (col.key === 'threadId') return columnVisibility.value.threadId;
    if (col.key === 'source') return columnVisibility.value.source;
    if (col.key === 'created') return columnVisibility.value.created;
    return true;
  });
});

/**
 * Parsed event data interface for JSON drawer display
 */
interface ParsedEventData extends Omit<EventStoreItem, 'payload' | 'metadata'> {
  payload: Record<string, unknown> | string;
  metadata?: Record<string, unknown> | string;
}

/**
 * Selected event with parsed payload for JSON drawer
 */
const selectedEventParsed = computed<ParsedEventData | null>(() => {
  if (!selectedEvent.value) return null;

  const eventData: ParsedEventData = { ...selectedEvent.value };

  // Parse payload field if it's a string
  if (typeof eventData.payload === 'string' && eventData.payload) {
    try {
      eventData.payload = JSON.parse(eventData.payload);
    } catch (e) {
      logger.warn('Failed to parse payload', e);
    }
  }

  // Parse metadata field if it's a string
  if (typeof eventData.metadata === 'string' && eventData.metadata) {
    try {
      eventData.metadata = JSON.parse(eventData.metadata);
    } catch (e) {
      logger.warn('Failed to parse metadata', e);
    }
  }

  return eventData;
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
 * Apply quick filters
 * @returns {void}
 */
function applyQuickFilters(): void {
  filters.value.threadId = quickSearchThreadId.value || undefined;
  resetCursor();
  void fetchEvents();
}

/**
 * Build evaFilters array from dynamic filters result for backend POST
 * Only includes fields that have a value set.
 * Maps frontend types to EVA bucket names and stringifies values.
 *
 * @returns {Array | undefined} EvaFilter array or undefined if no dynamic filters
 */
function buildEvaFilters(): Array<{ fieldId: number; bucket: string; operator: string; value: string; endValue?: string }> | undefined {
  if (!dynamicFiltersResult.value) return undefined;

  const evaFilters: Array<{ fieldId: number; bucket: string; operator: string; value: string; endValue?: string }> = [];

  /** Map frontend field type to backend EVA bucket name */
  const typeToBucket: Record<string, string> = {
    number: 'number',
    string: 'string',
    boolean: 'bool',
    date: 'date',
  };

  for (const field of dynamicFiltersResult.value.fields) {
    // Skip fields without a value
    if (field.value === undefined || field.value === null || field.value === '') continue;

    const bucket = typeToBucket[field.type] || 'string';
    const filter: { fieldId: number; bucket: string; operator: string; value: string; endValue?: string } = {
      fieldId: field.fieldId,
      bucket,
      operator: field.type === 'boolean' ? 'eq' : field.operator,
      value: String(field.value),
    };

    // Add endValue for "between" operator
    if (field.operator === 'between' && field.endValue !== undefined && field.endValue !== null && field.endValue !== '') {
      filter.endValue = String(field.endValue);
    }

    // For "like" operator on strings, append % for starts-with
    if (field.operator === 'like' && bucket === 'string') {
      filter.value = filter.value + '%';
    }

    evaFilters.push(filter);
  }

  return evaFilters.length > 0 ? evaFilters : undefined;
}

/**
 * Fetch events from API with current filters and cursor pagination
 * Uses POST /store/query to support EVA dynamic field filters
 *
 * @param {string} direction - Direction to fetch: 'next' (older) or 'prev' (newer)
 * @returns {Promise<void>}
 */
async function fetchEvents(direction: 'next' | 'prev' = 'next'): Promise<void> {
  loading.value = true;

  try {
    const queryBody: Record<string, any> = {
      limit: limit.value,
      direction,
    };

    // Set cursor based on direction
    if (direction === 'next' && cursor.value.next) {
      queryBody.cursor = cursor.value.next;
    } else if (direction === 'prev' && cursor.value.prev) {
      queryBody.cursor = cursor.value.prev;
    }

    // Add filters conditionally
    if (filters.value.threadId) {
      queryBody.threadId = filters.value.threadId;
    }
    if (filters.value.assetId) {
      queryBody.assetId = filters.value.assetId;
    }
    if (filters.value.assetTemplateId) {
      queryBody.assetTemplateId = filters.value.assetTemplateId;
    }
    if (filters.value.eventType) {
      queryBody.eventType = filters.value.eventType;
    }
    if (filters.value.source) {
      queryBody.source = filters.value.source;
    }
    if (filters.value.startTime) {
      queryBody.startTime = filters.value.startTime;
    }
    if (filters.value.endTime) {
      queryBody.endTime = filters.value.endTime;
    }

    // Add EVA dynamic field filters
    const evaFilters = buildEvaFilters();
    if (evaFilters) {
      queryBody.evaFilters = evaFilters;
    }

    logger.debug('Fetching with params', queryBody);

    const response = await apis.events.events.queryStore(queryBody);

    logger.debug('Response received', response);

    eventsList.value = (response.items || []) as EventStoreItem[];

    // Update cursor state
    const newCursor: EventStorePageCursor = {
      hasNext: response.hasNext,
      hasPrevious: response.hasPrevious,
      current: queryBody.cursor,
      next: response.nextCursor ?? undefined,
      prev: response.prevCursor ?? undefined,
    };
    cursor.value = newCursor;

  } catch (error: any) {
    logger.error('Error fetching events', error);
    notifyFail({
      message: error?.response?.data?.message || t.notifications.loadFailed.value,
      timeout: 5000,
    });
    eventsList.value = [];
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
  await fetchEvents('next');
}

/**
 * Handle advanced filters apply from drawer
 * @param {Record<string, any>} appliedFilters - Applied filter values
 * @returns {void}
 */
function handleAdvancedFiltersApply(appliedFilters: Record<string, any>): void {
  advancedFilterValues.value = { ...appliedFilters };

  // Map dates to startTime/endTime
  if (appliedFilters.startDate) {
    filters.value.startTime = new Date(appliedFilters.startDate).toISOString();
  } else {
    const today = getTodayDateRange();
    filters.value.startTime = today.startTime;
  }

  if (appliedFilters.endDate) {
    filters.value.endTime = new Date(appliedFilters.endDate + 'T23:59:59.999Z').toISOString();
  } else {
    const today = getTodayDateRange();
    filters.value.endTime = today.endTime;
  }

  // Map other filters
  filters.value.threadId = appliedFilters.threadId || undefined;
  filters.value.assetTemplateId = dynamicFiltersResult.value?.assetTemplateId || undefined;
  filters.value.source = appliedFilters.source || undefined;

  hasPendingAdvancedFilters.value = false;
  showAdvancedFilters.value = false;
  resetCursor();
  void fetchEvents();
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
  if (key === 'threadId') {
    quickSearchThreadId.value = '';
    advancedFilterValues.value.threadId = '';
    filters.value.threadId = undefined;
  } else if (key === 'includeChildren') {
    advancedFilterValues.value.includeChildren = null;
  } else if (key === 'startDate') {
    advancedFilterValues.value.startDate = '';
    const today = getTodayDateRange();
    filters.value.startTime = today.startTime;
  } else if (key === 'endDate') {
    advancedFilterValues.value.endDate = '';
    const today = getTodayDateRange();
    filters.value.endTime = today.endTime;
  } else if (key === 'dynamicFilter') {
    dynamicFiltersResult.value = null;
    filters.value.assetTemplateId = undefined;
  } else if (key === 'source') {
    advancedFilterValues.value.source = '';
    filters.value.source = undefined;
  }
  resetCursor();
  void fetchEvents();
}

/**
 * Clear all filters
 * @returns {void}
 */
function clearAllFilters(): void {
  quickSearchThreadId.value = '';
  dynamicFiltersResult.value = null;
  advancedFilterValues.value = {
    includeChildren: null,
    startDate: '',
    endDate: '',
    threadId: '',
    source: '',
  };
  const today = getTodayDateRange();
  filters.value = {
    threadId: undefined,
    assetId: undefined,
    assetTemplateId: undefined,
    eventType: undefined,
    source: undefined,
    startTime: today.startTime,
    endTime: today.endTime,
  };
  hasPendingAdvancedFilters.value = false;
  resetCursor();
  void fetchEvents();
}

/**
 * Handle view event details
 *
 * @param {EventStoreItem} event - Event to view
 */
function handleViewEvent(event: EventStoreItem): void {
  selectedEvent.value = event;
  jsonDrawerOpen.value = true;
}

/**
 * Handle column visibility change
 *
 * @param {ListHeaderMenuColumn[]} updatedColumns - Updated columns
 */
function handleColumnChange(updatedColumns: ListHeaderMenuColumn[]): void {
  menuColumns.value = updatedColumns;
}

/**
 * Navigate to next page
 */
function handleNextPage(): void {
  if (cursor.value.hasNext) {
    void fetchEvents('next');
  }
}

/**
 * Navigate to previous page
 */
function handlePrevPage(): void {
  if (cursor.value.hasPrevious) {
    void fetchEvents('prev');
  }
}

/**
 * Handle dynamic filters apply from DynamicFiltersDrawer
 * Saves the result, sets assetTemplateId filter, and refetches events
 *
 * @param {DynamicFiltersResult} result - Applied dynamic filters result
 */
function handleDynamicApply(result: DynamicFiltersResult): void {
  dynamicFiltersResult.value = result;
  filters.value.assetTemplateId = result.assetTemplateId;
  hasPendingDynamicFilters.value = false;
  resetCursor();
  void fetchEvents();
}

/**
 * Handle dynamic filters reset from DynamicFiltersDrawer
 * Clears the result and assetTemplateId filter, then refetches events
 */
function handleDynamicReset(): void {
  dynamicFiltersResult.value = null;
  filters.value.assetTemplateId = undefined;
  hasPendingDynamicFilters.value = false;
  resetCursor();
  void fetchEvents();
}

/**
 * Handle pending change in dynamic filters
 *
 * @param {boolean} hasPending - Whether there are pending dynamic filters
 */
function handleDynamicPendingChange(hasPending: boolean): void {
  hasPendingDynamicFilters.value = hasPending;
}

/** LIFECYCLE */

onMounted(() => {
  void fetchEvents();
});

/** Organization change refresh */
useOrgChangeRefresh(() => {
  quickSearchThreadId.value = '';
  dynamicFiltersResult.value = null;
  advancedFilterValues.value = {
    includeChildren: null,
    startDate: '',
    endDate: '',
    threadId: '',
    source: '',
  };
  const today = getTodayDateRange();
  filters.value = {
    threadId: undefined,
    assetId: undefined,
    assetTemplateId: undefined,
    eventType: undefined,
    source: undefined,
    startTime: today.startTime,
    endTime: today.endTime,
  };
  hasPendingAdvancedFilters.value = false;
  hasPendingDynamicFilters.value = false;
  resetCursor();
  void fetchEvents();
});
</script>

<template>
  <q-page padding>
    <!-- Page Header -->
    <div id="page-header-section">
      <PageHeader
        icon="storage"
        iconColor="primary"
        :title="t.pageHeader.title.value"
        :description="t.pageHeader.description.value"
        :tour="{ enabled: true }"
        @start-tour="startTour"
      />
    </div>

    <!-- Filters Section -->
    <div class="text-caption text-grey-7 q-mb-xs">{{ t.filters.label.value }}</div>
    <div class="row items-center q-col-gutter-sm q-mb-md">
      <!-- Search Input -->
      <div id="filter-search-input" class="col">
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

      <!-- Filter Icon Button -->
      <div class="col-auto">
        <q-btn
          id="advanced-filters-btn"
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

      <!-- Dynamic Filters Button -->
      <div class="col-auto">
        <q-btn
          id="dynamic-filters-btn"
          round
          flat
          icon="schema"
          color="grey-7"
          @click="showDynamicFilters = true"
        >
          <q-badge
            v-if="dynamicFiltersResult || hasPendingDynamicFilters"
            :color="hasPendingDynamicFilters ? 'warning' : 'secondary'"
            floating
            rounded
            :label="hasPendingDynamicFilters ? '!' : '1'"
          />
          <AppTooltip :content="hasPendingDynamicFilters
            ? t.filters.pendingFilters.value
            : t.filters.dynamicFilters.value"
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

    <!-- List Header Menu -->
    <div class="row items-center q-pt-xl q-mb-md">
      <div class="col">
        <div class="row items-center">
          <q-icon name="storage" size="sm" color="primary" class="q-mr-sm"/>
          <div class="text-subtitle1 text-weight-medium text-primary">{{ t.listHeader.title.value }}</div>
        </div>
      </div>
      <div class="col-auto">
        <ListHeaderMenu
          icon="storage"
          :items-count="eventsList.length"
          :item-label="t.listHeader.itemLabel.value"
          :item-label-plural="t.listHeader.itemLabelPlural.value"
          :items-per-page="limit"
          :filtered="isFiltered"
          :columns="menuColumns"
          :refreshing="loading"
          :last-updated-at="lastUpdatedAt"
          @update:columns="handleColumnChange"
          @refresh="refreshData"
        />
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="q-pa-xl text-center">
      <q-spinner color="primary" size="3em" />
    </div>

    <!-- Empty State -->
    <ListCardEmpty
      v-else-if="eventsList.length === 0"
      :title="t.empty.title.value"
      :description="t.empty.description.value"
      icon="event_busy"
    />

    <!-- Events List -->
    <div v-else id="results-section" class="q-gutter-sm">
      <DataRow
        v-for="(event, index) in eventsList"
        :key="`event-${index}`"
        :columns="visibleColumns"
        :data="event"
        :actions="rowActionsConfig"
        @click="handleViewEvent(event)"
        @view="handleViewEvent(event)"
      />
    </div>

    <!-- Pagination -->
    <div v-if="eventsList.length > 0" class="row justify-center q-mt-lg q-gutter-sm">
      <q-btn
        flat
        icon="chevron_left"
        :label="t.pagination.previous.value"
        :disable="!cursor.hasPrevious"
        @click="handlePrevPage"
      />
      <q-btn
        flat
        icon-right="chevron_right"
        :label="t.pagination.next.value"
        :disable="!cursor.hasNext"
        @click="handleNextPage"
      />
    </div>

    <!-- JSON Drawer for Event Details -->
    <JsonDrawer
      :show="jsonDrawerOpen"
      :title="t.drawer.title.value"
      :json-data="selectedEventParsed || {}"
      @update:show="jsonDrawerOpen = $event"
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

    <!-- Dynamic Filters Drawer -->
    <DynamicFiltersDrawer
      v-model="showDynamicFilters"
      :demo="dynamicFiltersDemoMode"
      @apply="handleDynamicApply"
      @reset="handleDynamicReset"
      @pending-change="handleDynamicPendingChange"
    />
  </q-page>
</template>

<style scoped lang="scss">
.filter-input {
  :deep(.q-field__control) {
    border-radius: var(--mapex-radius-md);
  }
}

.rounded-borders {
  border-radius: var(--mapex-radius-md);
}
</style>
