<script setup lang="ts">
/** TYPE IMPORTS */
import type { WorkflowExecutionsCursor, WorkflowExecutionsFilters, WorkflowExecutionItem } from './interfaces';
import type { ListHeaderMenuColumn } from '@components/headers';
import type { DataRowColumn } from '@components/cards';
import type { FilterField, FilterAutocompleteOption } from '@components/drawers';

/** LIBRARY IMPORTS */
import { ref, computed, onMounted } from 'vue';

/** COMPONENT IMPORTS */
import { PageHeader, ListHeaderMenu } from '@components/headers';
import { ListCardEmpty, DataRow } from '@components/cards';
import { AdvancedFiltersDrawer } from '@components/drawers';
import { ExecutionDetailModal } from './components/ExecutionDetailModal';

/** COMPOSABLES */
import { useOrgChangeRefresh } from '@composables/organizations';
import { useWorkflowExecutionsPageTranslations } from '@composables/i18n/pages/logs/workflowExecutionsPage';
import { useLogger } from '@composables/useLogger';

/** SERVICES */
import { apis } from '@services/mapex';

/** CONSTANT IMPORTS */
import {
  DEFAULT_LIMIT,
  MAX_VISIBLE_CHIPS,
  STATUS_COLORS,
  STATUS_ICONS,
  DURATION_THRESHOLDS,
  DURATION_COLORS,
  HOT_STATUSES,
  COLD_STATUSES,
} from './constants';

defineOptions({
  name: 'WorkflowExecutionsPage',
});

/** STATE */
const logger = useLogger('WorkflowExecutionsPage');
const t = useWorkflowExecutionsPageTranslations();

const isLoading = ref(false);
const lastUpdatedAt = ref<number | undefined>(undefined);
const executionsList = ref<WorkflowExecutionItem[]>([]);
const selectedExecution = ref<WorkflowExecutionItem | null>(null);
const showDetailModal = ref(false);
const showAdvancedFilters = ref(false);
const hasPendingAdvancedFilters = ref(false);
const limit = ref(DEFAULT_LIMIT);

const cursor = ref<WorkflowExecutionsCursor>({
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
  return { startTime: startOfDay.toISOString(), endTime: endOfDay.toISOString() };
}

const todayRange = getTodayDateRange();
const filters = ref<WorkflowExecutionsFilters>({
  startTime: todayRange.startTime,
  endTime: todayRange.endTime,
});

const quickSearchValue = ref('');
const quickStatusValue = ref<string | null>(null);

/** COMPUTED — Status options for quick filter */
const statusOptions = computed(() => [
  { label: t.statusOptions.all.value, value: null },
  { label: t.statusOptions.running.value, value: 'running' },
  { label: t.statusOptions.waiting.value, value: 'waiting' },
  { label: t.statusOptions.completed.value, value: 'completed' },
  { label: t.statusOptions.failed.value, value: 'failed' },
  { label: t.statusOptions.cancelled.value, value: 'cancelled' },
]);

/** COMPUTED — Advanced filter fields */
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
    key: 'startDate',
    type: 'input',
    label: t.filters.startDate.value,
    icon: 'event',
    inputType: 'date',
  },
  {
    key: 'endDate',
    type: 'input',
    label: t.filters.endDate.value,
    icon: 'event',
    inputType: 'date',
  },
  {
    key: 'instanceId',
    type: 'input',
    label: t.filters.instanceId.value,
    icon: 'tag',
    placeholder: t.filters.instanceIdPlaceholder.value,
  },
  {
    key: 'definitionId',
    type: 'autocomplete',
    label: t.filters.workflow.value,
    icon: 'account_tree',
    placeholder: t.filters.workflowPlaceholder.value,
    fetchOptions: fetchWorkflowDefinitions,
  },
]);

const advancedFilterValues = ref<Record<string, unknown>>({
  startDate: todayRange.startTime.split('T')[0],
  endDate: todayRange.endTime.split('T')[0],
  includeChildren: null,
  instanceId: null,
  definitionId: null,
});

/** COMPUTED — Active filter chips */
const activeFilterChips = computed(() => {
  const chips: Array<{ key: string; label: string; value: string }> = [];
  if (filters.value.includeChildren !== undefined) {
    chips.push({ key: 'includeChildren', label: t.filters.includeChildren.value, value: filters.value.includeChildren ? t.filters.options.yes.value : t.filters.options.no.value });
  }
  if (filters.value.instanceId) {
    chips.push({ key: 'instanceId', label: t.filters.instanceId.value, value: filters.value.instanceId });
  }
  if (filters.value.definitionId) {
    const defLabel = (advancedFilterValues.value.definitionIdLabel as string) || filters.value.definitionId;
    chips.push({ key: 'definitionId', label: t.filters.workflow.value, value: defLabel });
  }
  return chips;
});



const advancedFiltersCount = computed(() => {
  let count = 0;
  if (filters.value.includeChildren !== undefined) count++;
  if (filters.value.instanceId) count++;
  if (filters.value.definitionId) count++;
  return count;
});

const visibleFilterChips = computed(() => activeFilterChips.value.slice(0, MAX_VISIBLE_CHIPS));
const hiddenFilterChips = computed(() => activeFilterChips.value.slice(MAX_VISIBLE_CHIPS));
const hiddenFiltersCount = computed(() => hiddenFilterChips.value.length);

/** COMPUTED — Column definitions */
const menuColumns = ref<ListHeaderMenuColumn[]>([
  { key: 'instanceName', label: t.columns.instanceName.value, visible: true },
  { key: 'triggerSource', label: t.columns.triggerSource.value, visible: true },
  { key: 'definitionName', label: t.columns.definitionName.value, visible: true },
  { key: 'status', label: t.columns.status.value, visible: true },
  { key: 'duration', label: t.columns.duration.value, visible: true },
  { key: 'timestamp', label: t.columns.timestamp.value, visible: true },
]);

const logColumns = computed<DataRowColumn[]>(() => [
  {
    key: 'status',
    label: t.columns.status.value,
    type: 'badge',
    visible: 'always',
    width: 120,
    format: (val: string) => getStatusLabel(val),
    color: (val: string) => STATUS_COLORS[val] || 'grey-6',
    icon: (val: string) => STATUS_ICONS[val] || 'help',
  },
  {
    key: 'instanceName',
    label: t.columns.instanceName.value,
    type: 'text',
    visible: 'always',
    width: 200,
    ellipsis: true,
  },
  {
    key: 'triggerSource',
    label: t.columns.triggerSource.value,
    type: 'chip',
    visible: 'laptop',
    width: 110,
    format: (val: string) => {
      const labels: Record<string, string> = { workflow: 'Workflow', subworkflow: 'Subworkflow', http: 'HTTP' };
      return labels[val] || val || '-';
    },
    color: (val: string) => {
      const colors: Record<string, string> = { workflow: 'teal-3', subworkflow: 'deep-purple-3', http: 'blue-3' };
      return colors[val] || 'grey-4';
    },
  },
  {
    key: 'definitionName',
    label: t.columns.definitionName.value,
    type: 'text',
    visible: 'laptop',
    width: 180,
    ellipsis: true,
  },
  {
    key: 'durationMs',
    label: t.columns.duration.value,
    type: 'text',
    visible: 'laptop',
    width: 100,
    format: (val: number, row: WorkflowExecutionItem) => {
      if (val > 0) return formatDuration(val);
      if (row.created) return formatDuration(Date.now() - new Date(row.created).getTime());
      return '-';
    },
    color: (val: number, row: WorkflowExecutionItem) => {
      const ms = val > 0 ? val : (row.created ? Date.now() - new Date(row.created).getTime() : 0);
      return getDurationColor(ms);
    },
  },
  {
    key: 'created',
    label: t.columns.timestamp.value,
    type: 'text',
    visible: 'laptop',
    width: 160,
    format: (val: string) => formatTimestamp(val),
  },
]);

/** FUNCTIONS — Data fetching */
async function fetchData(direction: 'next' | 'prev' = 'next'): Promise<void> {
  isLoading.value = true;

  try {
  const status = quickStatusValue.value;
  const isHot = status ? HOT_STATUSES.includes(status as typeof HOT_STATUSES[number]) : false;
  const isCold = status ? COLD_STATUSES.includes(status as typeof COLD_STATUSES[number]) : false;

  let coldItems: WorkflowExecutionItem[] = [];
  let hotItems: WorkflowExecutionItem[] = [];

  // COLD source (ClickHouse via Events Service) — terminal statuses ONLY
  if (!status || isCold) {
    try {
      const coldParams: Record<string, unknown> = {
        limit: limit.value,
        direction,
      };
      if (direction === 'next' && cursor.value.next) coldParams.cursor = cursor.value.next;
      if (direction === 'prev' && cursor.value.prev) coldParams.cursor = cursor.value.prev;
      // When "All", force cold to only fetch terminal statuses
      if (status) {
        coldParams.status = status;
      }
      if (filters.value.instanceId) coldParams.instanceId = filters.value.instanceId;
      if (filters.value.definitionId) coldParams.definitionId = filters.value.definitionId;
      if (filters.value.startTime) coldParams.startTime = filters.value.startTime;
      if (filters.value.endTime) coldParams.endTime = filters.value.endTime;
      if (typeof filters.value.includeChildren === 'boolean') coldParams.includeChildren = filters.value.includeChildren;

      const coldResponse = await apis.events.events.listWorkflow(coldParams);
      coldItems = (coldResponse.items || []).map((item: Record<string, unknown>) => ({
        id: (item.executionId as string) || '',
        workflowUUID: item.workflowUUID as string,
        instanceId: item.instanceId as string,
        definitionId: item.definitionId as string,
        workflowName: item.workflowName as string,
        instanceName: (item.instanceName as string) || (item.workflowName as string) || '',
        definitionName: (item.definitionName as string) || '',
        status: item.status as string,
        durationMs: item.durationMs as number,
        errorMessage: item.errorMessage as string,
        executionPath: item.executionPath as string,
        nodeOutputs: item.nodeOutputs as string,
        errorInfo: item.errorInfo as string,
        eventPayload: item.eventPayload as string,
        triggerSource: (item.triggerSource as string) || '',
        parentExecutionId: item.parentExecutionId as string,
        depth: (item.depth as number) || 0,
        created: item.created as string,
        finished: item.finished as string,
        source: 'cold' as const,
      }));

      cursor.value = {
        hasNext: coldResponse.hasNext,
        hasPrevious: coldResponse.hasPrevious,
        next: coldResponse.nextCursor as string | undefined,
        prev: coldResponse.prevCursor as string | undefined,
      };
    } catch (err) {
      logger.error('Cold source (Events) failed', err);
    }
  }

  // HOT source (MongoDB via Workflow Service) — ONLY running + waiting
  // Note: hot API does NOT support startTime/endTime — date filtering is cold-only
  if (!status || isHot) {
    try {
      const hotParams: Record<string, unknown> = {
        page: 1,
        perPage: limit.value,
      };
      hotParams.status = status || HOT_STATUSES.join(',');
      if (filters.value.instanceId) hotParams.instanceId = filters.value.instanceId;
      if (filters.value.definitionId) hotParams.definitionId = filters.value.definitionId;
      if (typeof filters.value.includeChildren === 'boolean') hotParams.includeChildren = filters.value.includeChildren;

      const hotResponse = await apis.workflows.execution.list(hotParams);
      const hotRawItems = (hotResponse as Record<string, unknown>).items as Array<Record<string, unknown>> || [];
      hotItems = hotRawItems.map((item) => ({
        id: (item._id as string) || (item.ID as string) || '',
        workflowUUID: (item.workflowUUID as string) || (item.WorkflowUUID as string) || '',
        instanceId: (item.instanceId as string) || (item.InstanceID as string) || '',
        definitionId: (item.definitionId as string) || (item.DefinitionID as string) || '',
        workflowName: (item.workflowName as string) || (item.WorkflowName as string) || '',
        instanceName: (item.instanceName as string) || (item.InstanceName as string) || (item.workflowName as string) || (item.WorkflowName as string) || '',
        definitionName: (item.definitionName as string) || (item.DefinitionName as string) || '',
        status: (item.status as string) || (item.Status as string) || '',
        durationMs: (item.durationMs as number) || (item.DurationMs as number) || 0,
        errorMessage: ((item.errorInfo || item.ErrorInfo) as Record<string, unknown>)?.message as string,
        executionPath: (item.executionPath || item.ExecutionPath) ? JSON.stringify(item.executionPath || item.ExecutionPath) : undefined,
        nodeOutputs: (item.nodeOutputs || item.NodeOutputs) ? JSON.stringify(item.nodeOutputs || item.NodeOutputs) : undefined,
        errorInfo: (item.errorInfo || item.ErrorInfo) ? JSON.stringify(item.errorInfo || item.ErrorInfo) : undefined,
        eventPayload: (item.eventPayload || item.EventPayload) ? JSON.stringify(item.eventPayload || item.EventPayload) : undefined,
        triggerSource: (item.triggerSource as string) || (item.TriggerSource as string) || '',
        parentExecutionId: (item.parentExecutionId as string) || (item.ParentExecutionID as string) || '',
        depth: (item.depth as number) || (item.Depth as number) || 0,
        created: (item.created as string) || (item.Created as string) || (item.startedAt as string) || (item.StartedAt as string) || '',
        finished: (item.completedAt as string) || (item.CompletedAt as string) || '',
        source: 'hot' as const,
      }));
    } catch (err) {
      logger.warn('Hot source unavailable', err);
    }
  }

  // Merge hot (running/waiting) + cold (terminal) — no overlap, no dedup needed
  const merged = [...hotItems, ...coldItems];
  merged.sort((a, b) => new Date(b.created).getTime() - new Date(a.created).getTime());

  executionsList.value = merged;
  } finally {
    isLoading.value = false;
    lastUpdatedAt.value = Date.now();
  }
}

/**
 * Refresh: resets cursor to initial and refetches from the first page with current filters.
 * Filters are preserved. Used by the ListHeaderMenu refresh button.
 */
async function refreshData(): Promise<void> {
  cursor.value = { hasNext: false, hasPrevious: false };
  await fetchData('next');
}

/** FUNCTIONS — Filters */
function applyQuickFilters(): void {
  cursor.value = { hasNext: false, hasPrevious: false };
  void fetchData();
}

function handleAdvancedFiltersApply(values: Record<string, unknown>): void {
  // Convert date strings to ISO datetime
  if (values.startDate) {
    const d = new Date(values.startDate as string);
    d.setHours(0, 0, 0, 0);
    filters.value.startTime = d.toISOString();
  } else {
    const today = getTodayDateRange();
    filters.value.startTime = today.startTime;
  }
  if (values.endDate) {
    const d = new Date(values.endDate as string);
    d.setHours(23, 59, 59, 999);
    filters.value.endTime = d.toISOString();
  } else {
    const today = getTodayDateRange();
    filters.value.endTime = today.endTime;
  }
  filters.value.includeChildren = values.includeChildren as boolean | undefined;
  filters.value.instanceId = (values.instanceId as string) || undefined;
  filters.value.definitionId = (values.definitionId as string) || undefined;
  advancedFilterValues.value = { ...values };
  showAdvancedFilters.value = false;
  cursor.value = { hasNext: false, hasPrevious: false };
  void fetchData();
}

function handleAdvancedFiltersReset(): void {
  const today = getTodayDateRange();
  filters.value = { startTime: today.startTime, endTime: today.endTime };
  advancedFilterValues.value = {
    startDate: today.startTime.split('T')[0],
    endDate: today.endTime.split('T')[0],
    includeChildren: null,
    instanceId: null,
    definitionId: null,
  };
  showAdvancedFilters.value = false;
  cursor.value = { hasNext: false, hasPrevious: false };
  void fetchData();
}

function handlePendingChange(hasPending: boolean): void {
  hasPendingAdvancedFilters.value = hasPending;
}

function removeFilter(key: string): void {
  if (key === 'includeChildren') {
    filters.value.includeChildren = undefined;
    advancedFilterValues.value.includeChildren = null;
  } else if (key === 'instanceId') {
    filters.value.instanceId = undefined;
    advancedFilterValues.value.instanceId = null;
  } else if (key === 'definitionId') {
    filters.value.definitionId = undefined;
    advancedFilterValues.value.definitionId = null;
  }
  cursor.value = { hasNext: false, hasPrevious: false };
  void fetchData();
}

function clearAllFilters(): void {
  const today = getTodayDateRange();
  filters.value = { startTime: today.startTime, endTime: today.endTime };
  quickSearchValue.value = '';
  quickStatusValue.value = null;
  advancedFilterValues.value = {
    startDate: today.startTime.split('T')[0],
    endDate: today.endTime.split('T')[0],
    includeChildren: null,
    instanceId: null,
    definitionId: null,
  };
  cursor.value = { hasNext: false, hasPrevious: false };
  void fetchData();
}

/** FUNCTIONS — Pagination */
function goToNextPage(): void {
  if (cursor.value.hasNext) void fetchData('next');
}

function goToPrevPage(): void {
  if (cursor.value.hasPrevious) void fetchData('prev');
}

function handleLimitChange(newLimit: number): void {
  limit.value = newLimit;
  cursor.value = { hasNext: false, hasPrevious: false };
  void fetchData();
}

function handleColumnsUpdate(cols: ListHeaderMenuColumn[]): void {
  menuColumns.value = cols;
}

/** FUNCTIONS — Autocomplete */

/**
 * Fetch workflow definitions for autocomplete filter
 * @param {string} search - Search term
 * @returns {Promise<FilterAutocompleteOption[]>}
 */
async function fetchWorkflowDefinitions(search: string): Promise<FilterAutocompleteOption[]> {
  const response = await apis.workflows.definition.list({
    page: 1,
    perPage: 10,
    name: search,
  });

  return (response.items || []).map((def: Record<string, unknown>) => ({
    id: (def._id as string) || (def.id as string) || '',
    label: (def.name as string) || '',
  }));
}

/** FUNCTIONS — Row click */
function cardClick(item: WorkflowExecutionItem): void {
  selectedExecution.value = item;
  showDetailModal.value = true;
}

/** FUNCTIONS — Formatting */
function getStatusLabel(status: string): string {
  const badgeMap: Record<string, string> = {
    running: t.statusBadge.running.value,
    waiting: t.statusBadge.waiting.value,
    created: t.statusBadge.created.value,
    completed: t.statusBadge.completed.value,
    failed: t.statusBadge.failed.value,
    cancelled: t.statusBadge.cancelled.value,
  };
  return badgeMap[status] || status.toUpperCase();
}

function formatDuration(ms: number): string {
  if (!ms || ms <= 0) return '-';
  if (ms < 1000) return `${ms}ms`;
  if (ms < 60000) return `${(ms / 1000).toFixed(1)}s`;
  return `${(ms / 60000).toFixed(1)}m`;
}

function getDurationColor(ms: number): string {
  if (ms < DURATION_THRESHOLDS.fast) return DURATION_COLORS.fast;
  if (ms < DURATION_THRESHOLDS.normal) return DURATION_COLORS.normal;
  return DURATION_COLORS.slow;
}

function formatTimestamp(ts: string): string {
  if (!ts) return '-';
  try {
    return new Date(ts).toLocaleString();
  } catch {
    return ts;
  }
}

/** LIFECYCLE */
useOrgChangeRefresh(() => {
  cursor.value = { hasNext: false, hasPrevious: false };
  void fetchData();
});

onMounted(() => {
  void fetchData();
});
</script>

<template>
  <q-page :class="showDetailModal ? 'q-pa-none full-detail' : 'q-pa-lg'">
    <!-- ═══ Detail View (replaces list when an execution is selected) ═══ -->
    <div v-if="showDetailModal && selectedExecution" key="detail" class="execution-detail-wrapper">
      <ExecutionDetailModal
        :execution="selectedExecution"
        @close="showDetailModal = false"
      />
    </div>

    <!-- ═══ List View ═══ -->
    <div v-else key="list">
      <!-- Page Header -->
      <PageHeader
        :title="t.pageHeader.title.value"
        :description="t.pageHeader.description.value"
        icon="timeline"
      />

      <!-- Filters Section -->
      <div class="text-caption text-grey-7 q-mb-xs">{{ t.filters.label.value }}</div>
      <div class="row items-center q-col-gutter-sm q-mb-md">
        <!-- Search Input -->
        <div class="col">
          <q-input
            v-model="quickSearchValue"
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
        <div class="col-auto" style="min-width: 140px">
          <q-select
            v-model="quickStatusValue"
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
              :label="advancedFiltersCount || '!'"
            />
            <q-tooltip>{{ hasPendingAdvancedFilters ? t.filters.pendingFilters.value : t.filters.advancedFilters.value }}</q-tooltip>
          </q-btn>
        </div>
      </div>

      <!-- Active filter chips -->
      <div v-if="activeFilterChips.length > 0" class="row q-gutter-xs q-mb-md items-center">
        <q-chip
          v-for="chip in visibleFilterChips"
          :key="chip.key"
          removable
          dense
          color="primary"
          text-color="white"
          @remove="removeFilter(chip.key)"
        >
          {{ chip.label }}: {{ chip.value }}
        </q-chip>
        <q-badge v-if="hiddenFiltersCount > 0" color="grey-6" class="cursor-pointer">
          +{{ hiddenFiltersCount }}
          <q-tooltip>
            <div v-for="chip in hiddenFilterChips" :key="chip.key">
              {{ chip.label }}: {{ chip.value }}
            </div>
          </q-tooltip>
        </q-badge>
        <q-btn flat dense size="sm" :label="t.filters.clearAll.value" color="negative" @click="clearAllFilters" />
      </div>

      <!-- Results Header -->
      <div class="row items-center justify-between q-pt-xl q-mb-sm">
        <div class="row items-center q-gutter-xs">
          <q-icon name="timeline" size="20px" color="primary" />
          <span class="text-subtitle2">{{ t.listTitle.value }}</span>
        </div>
        <ListHeaderMenu
          :columns="menuColumns"
          :items-count="executionsList.length"
          :item-label="'execution'"
          :items-per-page="limit"
          :refreshing="isLoading"
          :last-updated-at="lastUpdatedAt"
          @update:columns="handleColumnsUpdate"
          @update:limit="handleLimitChange"
          @refresh="refreshData"
        />
      </div>

      <!-- Loading -->
      <div v-if="isLoading" class="row justify-center q-pa-xl">
        <q-spinner color="primary" size="40px" />
      </div>

      <!-- Data List -->
      <template v-else-if="executionsList.length > 0">
        <DataRow
          v-for="item in executionsList"
          :key="item.id"
          :data="item"
          :columns="logColumns"
          :actions="{ showEdit: false, showDelete: false }"
          @click="cardClick(item)"
          @dblclick="cardClick(item)"
          @view="cardClick(item)"
        />

        <!-- Cursor Pagination -->
        <div class="row justify-center q-mt-md q-gutter-sm">
          <q-btn
            flat
            :disable="!cursor.hasPrevious"
            :label="t.pagination.newer.value"
            icon="chevron_left"
            @click="goToPrevPage"
          />
          <q-btn
            flat
            :disable="!cursor.hasNext"
            :label="t.pagination.older.value"
            icon-right="chevron_right"
            @click="goToNextPage"
          />
        </div>
      </template>

      <!-- Empty state -->
      <ListCardEmpty
        v-else
        :title="t.empty.title.value"
        :description="t.empty.description.value"
        icon="timeline"
      />

      <!-- Advanced Filters Drawer -->
      <AdvancedFiltersDrawer
        v-model="showAdvancedFilters"
        :fields="advancedFilterFields"
        :values="advancedFilterValues"
        @apply="handleAdvancedFiltersApply"
        @reset="handleAdvancedFiltersReset"
        @pending-change="handlePendingChange"
      />
    </div>
  </q-page>
</template>

<style lang="scss" scoped>
.full-detail {
  display: flex;
  flex-direction: column;

  &.q-page {
    min-height: 0 !important;
    height: 100%;
  }
}

.execution-detail-wrapper {
  display: flex;
  flex-direction: column;
  flex: 1;
  min-height: 0;
}

.filter-input {
  min-width: 0;
}
</style>
