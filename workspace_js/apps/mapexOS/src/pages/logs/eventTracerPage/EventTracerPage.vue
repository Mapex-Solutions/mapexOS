<script setup lang="ts">
defineOptions({
  name: 'EventTracerPage'
});

/** TYPE IMPORTS */
import type { FilterField, FilterValues } from '@components/drawers';
import type {
  EventsRawResponse,
  EventsJsExecResponse,
  EventsRouterResponse,
  EventsTriggerResponse,
} from '@mapexos/schemas';
import type {
  EventTraceResult,
  IngestionPhase,
  RoutingPhase,
  RawEventStageData,
  JsExecStageData,
  RouterStageData,
  DirectTriggerStageData,
  RouterResult,
  TraceSummary,
} from './interfaces';

/** VUE IMPORTS */
import { ref, computed, onMounted } from 'vue';
import { useRoute } from 'vue-router';

/** COMPONENTS */
import { PageHeader } from '@components/headers';
import { JsonDrawer, AdvancedFiltersDrawer } from '@components/drawers';
import { AppTooltip } from '@components/tooltips';
import { EventTraceVisualization } from './components';

/** COMPOSABLES */
import { useEventTracerPageTranslations } from '@composables/i18n/pages/logs/eventTracerPage';
import { useLogger } from '@composables/useLogger';

/** UTILS */
import { notifyFail } from '@utils/alert';

/** SERVICES */
import { apis } from '@services/mapex';

/** LOCAL IMPORTS */
import { MIN_SEARCH_LENGTH, MAX_ITEMS_PER_STAGE } from './constants';

/** COMPOSABLES */
const t = useEventTracerPageTranslations();
const logger = useLogger('EventTracerPage');
const route = useRoute();

/** STATE - Search */
const loading = ref(false);
const hasSearched = ref(false);
const loadingMessage = ref('');

/** LOCAL IMPORTS */
import type { EventTracerFilters } from './interfaces/eventTracerPage.interface';

/** STATE - Filters */
const filters = ref<EventTracerFilters>({});

/** STATE - Quick Filters */
const quickSearchEventTrackerId = ref('');

/** STATE - Advanced Filters */
const showFiltersDrawer = ref(false);
const advancedFilterValues = ref<FilterValues>({
  includeChildren: null,
});
const hasPendingAdvancedFilters = ref(false);

/** STATE - Trace Result */
const traceResult = ref<EventTraceResult | null>(null);

/** STATE - JSON Drawer */
const jsonDrawerOpen = ref(false);
const jsonDrawerData = ref<any>(null);
const jsonDrawerTitle = ref('');
const jsonDrawerSubtitle = ref('');

/** COMPUTED */

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
]);

/**
 * Check if any filter is active
 */
const hasActiveFilters = computed(() => {
  return !!(
    filters.value.eventTrackerId ||
    filters.value.includeChildren !== undefined
  );
});

/**
 * Count of advanced filters only (for badge on filter icon)
 */
const advancedFiltersCount = computed(() => {
  let count = 0;
  if (filters.value.includeChildren !== undefined) count++;
  return count;
});

/**
 * Active filter chips for display
 */
const activeFilterChips = computed(() => {
  const chips: Array<{ key: string; label: string; value: string }> = [];

  if (filters.value.eventTrackerId) {
    chips.push({
      key: 'eventTrackerId',
      label: t.filters.eventTrackerId.value,
      value: filters.value.eventTrackerId
    });
  }
  if (filters.value.includeChildren !== undefined) {
    chips.push({
      key: 'includeChildren',
      label: t.filters.includeChildren.value,
      value: filters.value.includeChildren ? t.filters.options.yes.value : t.filters.options.no.value
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

/** FUNCTIONS */

/**
 * Get today's date range
 * @returns {{ startTime: string; endTime: string }} Today's date range in ISO format
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

/**
 * Apply quick filters (search)
 */
function applyQuickFilters(): void {
  const eventTrackerId = quickSearchEventTrackerId.value?.trim();

  if (!eventTrackerId || eventTrackerId.length < MIN_SEARCH_LENGTH) {
    if (eventTrackerId && eventTrackerId.length > 0) {
      notifyFail({ message: t.messages.minLength.value });
    }
    return;
  }

  filters.value.eventTrackerId = eventTrackerId;

  // Use today's date range as default
  const today = getTodayDateRange();
  filters.value.startTime = today.startTime;
  filters.value.endTime = today.endTime;

  void fetchTraceData();
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

  // Reset pending state after apply
  hasPendingAdvancedFilters.value = false;

  // If we have an eventTrackerId, refetch data
  if (filters.value.eventTrackerId) {
    void fetchTraceData();
  }
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
  };
}

/**
 * Remove individual filter
 * @param {string} key - Filter key to remove
 */
function removeFilter(key: string): void {
  if (key === 'eventTrackerId') {
    delete filters.value.eventTrackerId;
    quickSearchEventTrackerId.value = '';
  } else if (key === 'includeChildren') {
    delete filters.value.includeChildren;
    advancedFilterValues.value.includeChildren = null;
  }

  // If we still have an eventTrackerId, refetch data
  if (filters.value.eventTrackerId) {
    void fetchTraceData();
  }
}

/**
 * Clear all filters
 */
function clearAllFilters(): void {
  // Reset filter state
  filters.value = {};

  // Reset quick filters
  quickSearchEventTrackerId.value = '';

  // Reset advanced filters
  advancedFilterValues.value = {
    includeChildren: null,
  };

  // Reset pending state
  hasPendingAdvancedFilters.value = false;

  // Reset search state
  hasSearched.value = false;
  traceResult.value = null;
}

/**
 * Fetch all trace data for the event tracker ID
 * @returns {Promise<void>}
 */
async function fetchTraceData(): Promise<void> {
  if (!filters.value.eventTrackerId) return;

  loading.value = true;
  hasSearched.value = true;
  traceResult.value = null;

  try {
    const eventTrackerId = filters.value.eventTrackerId;
    const startTime = filters.value.startTime;
    const endTime = filters.value.endTime;

    // Fetch all stages in parallel
    loadingMessage.value = t.messages.loading.value;

    const [rawEvents, jsExecEvents, routerEvents, triggerEvents] = await Promise.all([
      fetchRawEvents(eventTrackerId, startTime, endTime),
      fetchJsExecEvents(eventTrackerId, startTime, endTime),
      fetchRouterEvents(eventTrackerId, startTime, endTime),
      fetchTriggerEvents(eventTrackerId, startTime, endTime),
    ]);

    // Build the trace result
    traceResult.value = buildTraceResult(
      eventTrackerId,
      rawEvents,
      jsExecEvents,
      routerEvents,
      triggerEvents,
    );

    logger.debug('Trace result built', traceResult.value);

  } catch (error: any) {
    logger.error('Error fetching trace data', error);
    notifyFail({
      message: error?.response?.data?.message || t.messages.loadFailed.value,
      timeout: 5000,
    });
  } finally {
    loading.value = false;
    loadingMessage.value = '';
  }
}

/**
 * Fetch raw events by event tracker ID
 */
async function fetchRawEvents(
  eventTrackerId: string,
  startTime?: string,
  endTime?: string,
): Promise<EventsRawResponse[]> {
  try {
    const response = await apis.events.events.listRaw({
      eventTrackerId,
      startTime,
      endTime,
      limit: MAX_ITEMS_PER_STAGE,
    });
    return response.items || [];
  } catch (error) {
    logger.error('Error fetching raw events', error);
    return [];
  }
}

/**
 * Fetch JS executor events by event tracker ID
 */
async function fetchJsExecEvents(
  eventTrackerId: string,
  startTime?: string,
  endTime?: string,
): Promise<EventsJsExecResponse[]> {
  try {
    const response = await apis.events.events.listJsExec({
      eventTrackerId,
      startTime,
      endTime,
      limit: MAX_ITEMS_PER_STAGE,
    });
    return response.items || [];
  } catch (error) {
    logger.error('Error fetching JS exec events', error);
    return [];
  }
}

/**
 * Fetch router events by event tracker ID
 */
async function fetchRouterEvents(
  eventTrackerId: string,
  startTime?: string,
  endTime?: string,
): Promise<EventsRouterResponse[]> {
  try {
    const response = await apis.events.events.listRouter({
      eventTrackerId,
      startTime,
      endTime,
      limit: MAX_ITEMS_PER_STAGE,
    });
    return response.items || [];
  } catch (error) {
    logger.error('Error fetching router events', error);
    return [];
  }
}

/**
 * Fetch trigger events by event tracker ID
 */
async function fetchTriggerEvents(
  eventTrackerId: string,
  startTime?: string,
  endTime?: string,
): Promise<EventsTriggerResponse[]> {
  try {
    const response = await apis.events.events.listTrigger({
      eventTrackerId,
      startTime,
      endTime,
      limit: MAX_ITEMS_PER_STAGE,
    });
    return response.items || [];
  } catch (error) {
    logger.error('Error fetching trigger events', error);
    return [];
  }
}

/**
 * Build the complete trace result from all fetched data
 */
function buildTraceResult(
  eventTrackerId: string,
  rawEvents: EventsRawResponse[],
  jsExecEvents: EventsJsExecResponse[],
  routerEvents: EventsRouterResponse[],
  triggerEvents: EventsTriggerResponse[],
): EventTraceResult {
  // Build ingestion phase
  const ingestion = buildIngestionPhase(rawEvents, jsExecEvents);

  // Build routing phase
  const routing = buildRoutingPhase(routerEvents, triggerEvents);

  // Calculate totals
  const totalExecutionTime =
    ingestion.totalDurationMs +
    routing.totalDurationMs;

  // Build summary
  const summary = buildTraceSummary(ingestion, routing);

  const allSuccess = ingestion.allSuccess && routing.allSuccess;

  return {
    eventTrackerId,
    threadId: rawEvents[0]?.threadId || jsExecEvents[0]?.threadId || '',
    created: rawEvents[0]?.created || jsExecEvents[0]?.created || new Date().toISOString(),
    ingestion,
    routing,
    totalExecutionTime,
    allSuccess,
    summary,
  };
}

/**
 * Build the ingestion phase from raw and JS exec events
 */
function buildIngestionPhase(
  rawEvents: EventsRawResponse[],
  jsExecEvents: EventsJsExecResponse[],
): IngestionPhase {
  const rawEvent = rawEvents[0];
  const jsExecEvent = jsExecEvents[0];

  const raw: RawEventStageData = {
    hasData: !!rawEvent,
    success: rawEvent?.success ?? null,
    created: rawEvent?.created ?? null,
    durationMs: null, // Raw events don't have duration
    error: rawEvent?.error ?? null,
    data: rawEvent ?? null,
    source: rawEvent?.source ?? null,
    threadId: rawEvent?.threadId ?? null,
  };

  const jsExec: JsExecStageData = {
    hasData: !!jsExecEvent,
    success: jsExecEvent?.success ?? null,
    created: jsExecEvent?.created ?? null,
    durationMs: jsExecEvent?.totalExecutionTime ?? null,
    error: jsExecEvent?.error ?? null,
    data: jsExecEvent ?? null,
    failedAt: jsExecEvent?.failedAt ?? null,
  };

  const totalDurationMs = (jsExec.durationMs || 0);
  const allSuccess = (!raw.hasData || raw.success === true) && (!jsExec.hasData || jsExec.success === true);

  return {
    raw,
    jsExec,
    totalDurationMs,
    allSuccess,
  };
}

/**
 * Build the routing phase from router and trigger events
 */
function buildRoutingPhase(
  routerEvents: EventsRouterResponse[],
  triggerEvents: EventsTriggerResponse[],
): RoutingPhase {
  const routerEvent = routerEvents[0];

  // Parse routers array from event JSON
  let parsedRouters: RouterResult[] | null = null;
  if (routerEvent?.event) {
    try {
      parsedRouters = JSON.parse(routerEvent.event);
    } catch (e) {
      logger.warn('Failed to parse router event JSON', e);
    }
  }

  const router: RouterStageData = {
    hasData: !!routerEvent,
    success: routerEvent?.success ?? null,
    created: routerEvent?.created ?? null,
    durationMs: null,
    error: routerEvent?.error ?? null,
    data: routerEvent ?? null,
    routerId: routerEvent?.routerId ?? null,
    name: routerEvent?.name ?? null,
    totalRouters: routerEvent?.totalRouters ?? 0,
    matchedCount: routerEvent?.matchedCount ?? 0,
    publishedCount: routerEvent?.publishedCount ?? 0,
    routers: parsedRouters,
  };

  // Check for special router types (save_event, lake_house)
  const saveEventRouter = parsedRouters?.find(r => r.kind === 'save_event');
  const lakeHouseRouter = parsedRouters?.find(r => r.kind === 'lake_house');

  // Separate direct triggers (from router) and rule engine triggers
  const directTriggers: DirectTriggerStageData[] = triggerEvents
    .filter(t => t.source === 'router')
    .map(t => ({
      hasData: true,
      success: t.success,
      created: t.created,
      durationMs: t.durationMs,
      error: t.error ?? null,
      data: t,
      triggerId: t.triggerId,
      triggerName: t.triggerName,
      triggerType: t.triggerType,
      category: t.category,
    }));

  // Calculate totals
  const destinationsCount = directTriggers.length +
    (saveEventRouter ? 1 : 0) + (lakeHouseRouter ? 1 : 0);

  const totalDurationMs =
    directTriggers.reduce((acc, t) => acc + (t.durationMs || 0), 0);

  const allSuccess =
    directTriggers.every(t => t.success);

  return {
    router,
    directTriggers,
    saveEvent: {
      hasData: !!saveEventRouter,
      success: saveEventRouter?.published ?? null,
      created: routerEvent?.created ?? null,
      durationMs: null,
      error: null,
      triggered: saveEventRouter?.matched ?? false,
      data: saveEventRouter ?? null,
    },
    lakeHouse: {
      hasData: !!lakeHouseRouter,
      success: lakeHouseRouter?.published ?? null,
      created: routerEvent?.created ?? null,
      durationMs: null,
      error: null,
      triggered: lakeHouseRouter?.matched ?? false,
      data: lakeHouseRouter ?? null,
    },
    destinationsCount,
    totalDurationMs,
    allSuccess,
  };
}

/**
 * Build trace summary statistics
 */
function buildTraceSummary(
  ingestion: IngestionPhase,
  routing: RoutingPhase,
): TraceSummary {
  let stagesWithData = 0;
  let stagesSucceeded = 0;
  let stagesFailed = 0;
  let firstFailure: string | null = null;

  // Count ingestion stages
  if (ingestion.raw.hasData) {
    stagesWithData++;
    if (ingestion.raw.success) stagesSucceeded++;
    else if (ingestion.raw.success === false) {
      stagesFailed++;
      if (!firstFailure) firstFailure = ingestion.raw.error || 'Raw event failed';
    }
  }
  if (ingestion.jsExec.hasData) {
    stagesWithData++;
    if (ingestion.jsExec.success) stagesSucceeded++;
    else if (ingestion.jsExec.success === false) {
      stagesFailed++;
      if (!firstFailure) firstFailure = ingestion.jsExec.error || 'JS Executor failed';
    }
  }

  // Count routing stages
  if (routing.router.hasData) stagesWithData++;
  stagesWithData += routing.directTriggers.length;

  routing.directTriggers.forEach(t => {
    if (t.success) stagesSucceeded++;
    else if (t.success === false) {
      stagesFailed++;
      if (!firstFailure) firstFailure = t.error || 'Direct trigger failed';
    }
  });

  const totalTriggers = routing.directTriggers.length;

  return {
    stagesWithData,
    stagesSucceeded,
    stagesFailed,
    totalTriggers,
    firstFailure,
  };
}

/**
 * Fields that contain JSON strings and should be parsed
 */
const JSON_STRING_FIELDS = [
  'event',
  'evaluationTree',
  'conditionLogs',
  'actionsToDispatch',
  'requestData',
  'responseData',
  'metadata',
  'stateChanges',
];

/**
 * Parse JSON string fields in an object
 * @param {Record<string, any>} data - Object with potential JSON string fields
 * @returns {Record<string, any>} Object with parsed JSON fields
 */
function parseJsonFields(data: Record<string, any> | null): Record<string, any> | null {
  if (!data) return null;

  const result = { ...data };

  for (const field of JSON_STRING_FIELDS) {
    if (result[field] && typeof result[field] === 'string') {
      try {
        result[field] = JSON.parse(result[field]);
      } catch {
        // Keep as string if parsing fails
      }
    }
  }

  return result;
}

/**
 * Handle view details event from visualization
 * @param {{ stage: string; data: any }} event - Stage and data to view
 */
function handleViewDetails(event: { stage: string; data: any }): void {
  jsonDrawerData.value = parseJsonFields(event.data);
  jsonDrawerTitle.value = getStageTitle(event.stage);
  jsonDrawerSubtitle.value = formatTimestamp(event.data?.created);
  jsonDrawerOpen.value = true;
}

/**
 * Get stage title for drawer
 * @param {string} stage - Stage identifier
 * @returns {string} Stage title
 */
function getStageTitle(stage: string): string {
  const titles: Record<string, string> = {
    raw: t.stages.raw.value,
    jsExec: t.stages.jsExec.value,
    router: t.stages.router.value,
    directTrigger: t.stages.directTrigger.value,
    saveEvent: t.stages.saveEvent.value,
    lakeHouse: t.stages.lakeHouse.value,
  };
  return titles[stage] || stage;
}

/**
 * Format timestamp for display
 * @param {string | null} value - ISO timestamp string
 * @returns {string} Formatted date string
 */
function formatTimestamp(value: string | null): string {
  if (!value) return '-';
  const date = new Date(value);
  return date.toLocaleDateString('en-US', {
    day: '2-digit',
    month: 'short',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
  });
}

/** LIFECYCLE */
onMounted(() => {
  // Check if eventTrackerId was passed as query param (e.g., from "Track this event" action)
  const queryEventTrackerId = route.query.eventTrackerId as string | undefined;

  if (queryEventTrackerId && queryEventTrackerId.length >= MIN_SEARCH_LENGTH) {
    // Set the filter and auto-fetch
    filters.value.eventTrackerId = queryEventTrackerId;
    quickSearchEventTrackerId.value = queryEventTrackerId;

    // Use today's date range as default
    const today = getTodayDateRange();
    filters.value.startTime = today.startTime;
    filters.value.endTime = today.endTime;

    // Auto-fetch trace data
    void fetchTraceData();
  }
});
</script>

<template>
  <q-page class="q-pt-lg">

    <!-- Header Section -->
    <PageHeader
      icon="account_tree"
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
          v-model="quickSearchEventTrackerId"
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

    <!-- Loading State -->
    <div v-if="loading" class="row justify-center q-pa-xl">
      <div class="column items-center">
        <q-spinner color="primary" size="50px" />
        <div class="text-grey-7 q-mt-md">{{ loadingMessage || t.messages.loading.value }}</div>
      </div>
    </div>

    <!-- Trace Visualization -->
    <EventTraceVisualization
      v-else-if="hasSearched"
      :trace-result="traceResult"
      :loading="loading"
      @view-details="handleViewDetails"
    />

    <!-- Empty State - No Search Yet -->
    <q-card v-else flat bordered>
      <q-card-section class="text-center q-pa-xl">
        <q-icon name="account_tree" size="64px" color="grey-5" />
        <div class="text-h6 text-grey-7 q-mt-md">{{ t.empty.searchFirst.value }}</div>
        <div class="text-body2 text-grey-6 q-mt-sm">{{ t.empty.searchFirstDescription.value }}</div>
      </q-card-section>
    </q-card>

    <!-- JSON Drawer -->
    <JsonDrawer
      v-if="jsonDrawerData"
      v-model:show="jsonDrawerOpen"
      :jsonData="jsonDrawerData"
      :editable="false"
      :title="jsonDrawerTitle"
      :subtitle="jsonDrawerSubtitle"
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

<style scoped lang="scss">
.filter-input {
  :deep(.q-field__control) {
    border-radius: var(--mapex-radius-md);
  }
}
</style>
