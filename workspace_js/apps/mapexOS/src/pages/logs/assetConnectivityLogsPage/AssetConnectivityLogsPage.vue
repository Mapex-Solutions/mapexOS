<script setup lang="ts">
defineOptions({
  name: 'AssetConnectivityLogsPage'
});

/** TYPE IMPORTS */
import type { ListHeaderMenuColumn } from '@components/headers';
import type { DataRowColumn, DataRowActionConfig } from '@components/cards';
import type { AssetConnectivityEvent } from '@mapexos/schemas';
import type {
  AssetConnectivityLogsPageFilters,
  AssetConnectivityLogsPageColumnVisibility,
  AssetConnectivityLogsPageCursor,
} from './interfaces';

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
import { useAssetConnectivityLogsPageTranslations } from '@composables/i18n/pages/logs/assetConnectivityLogsPage';
import { useLogger } from '@composables/useLogger';

/** UTILS */
import { notifyFail } from '@utils/alert';

/** SERVICES */
import { apis } from '@services/mapex';

/** LOCAL IMPORTS */
import {
  DEFAULT_LIMIT,
  COLUMN_VISIBILITY_DEFAULTS,
  EVENT_TYPE_COLORS,
  DEFAULT_COLOR,
} from './constants';

/** CONSTANTS */
const MAX_VISIBLE_CHIPS = 2;

/** COMPOSABLES */
const t = useAssetConnectivityLogsPageTranslations();
const logger = useLogger('AssetConnectivityLogsPage');
const router = useRouter();

/** STATE */
const loading = ref(false);
const lastUpdatedAt = ref<number | undefined>(undefined);
const logsList = ref<AssetConnectivityEvent[]>([]);
const selectedEvent = ref<AssetConnectivityEvent | null>(null);
const jsonDrawerOpen = ref(false);

/** Cursor Pagination */
const limit = ref(DEFAULT_LIMIT);
const cursor = ref<AssetConnectivityLogsPageCursor>({
  current: undefined,
  next: undefined,
  prev: undefined,
  hasNext: false,
  hasPrevious: false,
});

/**
 * Get today's date range (start of day to end of day in ISO format)
 */
function getTodayDateRange(): { from: string; to: string } {
  const now = new Date();
  const startOfDay = new Date(now.getFullYear(), now.getMonth(), now.getDate(), 0, 0, 0, 0);
  const endOfDay = new Date(now.getFullYear(), now.getMonth(), now.getDate(), 23, 59, 59, 999);
  return {
    from: startOfDay.toISOString(),
    to: endOfDay.toISOString(),
  };
}

/** Filters - Default date range is today */
const todayRange = getTodayDateRange();
const filters = ref<AssetConnectivityLogsPageFilters>({
  from: todayRange.from,
  to: todayRange.to,
});

/** Quick filters state */
const quickSearchAssetUUID = ref('');
const quickEventType = ref<'offline' | 'online' | null>(null);

/** Advanced filters state */
const showAdvancedFilters = ref(false);
const advancedFilterValues = ref<Record<string, any>>({
  includeChildren: null,
  startDate: todayRange.from.split('T')[0],
  endDate: todayRange.to.split('T')[0],
  eventType: null,
  assetUUID: '',
});
const hasPendingAdvancedFilters = ref(false);

/** Column visibility */
const menuColumns = ref<ListHeaderMenuColumn[]>([
  { key: 'asset', label: t.columns.asset.value, visible: COLUMN_VISIBILITY_DEFAULTS.asset },
  { key: 'assetUUID', label: t.columns.assetUUID.value, visible: COLUMN_VISIBILITY_DEFAULTS.assetUUID },
  { key: 'eventType', label: t.columns.eventType.value, visible: COLUMN_VISIBILITY_DEFAULTS.eventType },
  { key: 'lastSeenAt', label: t.columns.lastSeenAt.value, visible: COLUMN_VISIBILITY_DEFAULTS.lastSeenAt },
  { key: 'missCount', label: t.columns.missCount.value, visible: COLUMN_VISIBILITY_DEFAULTS.missCount },
  { key: 'thresholdMinutes', label: t.columns.thresholdMinutes.value, visible: COLUMN_VISIBILITY_DEFAULTS.thresholdMinutes },
  { key: 'created', label: t.columns.timestamp.value, visible: COLUMN_VISIBILITY_DEFAULTS.created },
]);

/** Row actions configuration */
const rowActionsConfig = computed<DataRowActionConfig>(() => ({
  showEdit: false,
  showView: true,
  showDelete: false,
  customActions: [
    {
      key: 'viewAsset',
      label: t.actions.viewAsset.value,
      icon: 'inventory_2',
      color: 'secondary',
    },
  ],
}));

/** COMPUTED */

/**
 * Event type select options for quick filter
 */
const eventTypeOptions = computed(() => [
  { label: t.filters.allStatus.value, value: null },
  { label: t.filters.options.offline.value, value: 'offline' },
  { label: t.filters.options.online.value, value: 'online' },
]);

/**
 * Count of active advanced filters
 */
const advancedFiltersCount = computed(() => {
  let count = 0;
  if (advancedFilterValues.value.includeChildren !== null) count++;
  if (advancedFilterValues.value.eventType) count++;
  if (advancedFilterValues.value.assetUUID) count++;
  return count;
});

/**
 * Active filter chips for display
 */
const activeFilterChips = computed(() => {
  const chips: Array<{ key: string; label: string; value: string }> = [];

  if (quickSearchAssetUUID.value) {
    chips.push({
      key: 'quickSearchAssetUUID',
      label: t.filters.assetUUID.value,
      value: quickSearchAssetUUID.value,
    });
  }

  if (quickEventType.value !== null) {
    const typeLabel = quickEventType.value === 'offline'
      ? t.filters.options.offline.value
      : t.filters.options.online.value;
    chips.push({
      key: 'eventType',
      label: t.filters.eventType.value,
      value: typeLabel,
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

  if (advancedFilterValues.value.assetUUID) {
    chips.push({
      key: 'assetUUID',
      label: t.filters.assetUUID.value,
      value: advancedFilterValues.value.assetUUID,
    });
  }

  return chips;
});

const visibleFilterChips = computed(() => activeFilterChips.value.slice(0, MAX_VISIBLE_CHIPS));
const hiddenFilterChips = computed(() => activeFilterChips.value.slice(MAX_VISIBLE_CHIPS));
const hiddenFiltersCount = computed(() => hiddenFilterChips.value.length);

/**
 * Column visibility state computed from menuColumns
 */
const columnVisibility = computed<AssetConnectivityLogsPageColumnVisibility>(() => ({
  asset: menuColumns.value.find((col) => col.key === 'asset')?.visible ?? true,
  assetUUID: menuColumns.value.find((col) => col.key === 'assetUUID')?.visible ?? true,
  eventType: menuColumns.value.find((col) => col.key === 'eventType')?.visible ?? true,
  lastSeenAt: menuColumns.value.find((col) => col.key === 'lastSeenAt')?.visible ?? true,
  missCount: menuColumns.value.find((col) => col.key === 'missCount')?.visible ?? true,
  thresholdMinutes: menuColumns.value.find((col) => col.key === 'thresholdMinutes')?.visible ?? false,
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
    icon: (_: any, row: AssetConnectivityEvent) => row.eventType === 'offline' ? 'wifi_off' : 'wifi',
    color: (_: any, row: AssetConnectivityEvent) => getEventTypeColor(row.eventType),
  },
  {
    key: 'asset',
    label: t.columns.asset.value,
    type: 'text',
    visible: 'always',
    width: 250,
    ellipsis: true,
    format: (_: any, row: AssetConnectivityEvent) => row.assetName || row.assetUUID || t.defaults.unknown.value,
    secondaryKey: 'assetUUID',
  },
  {
    key: 'assetUUID',
    label: t.columns.assetUUID.value,
    type: 'text',
    visible: 'laptop',
    width: 180,
    ellipsis: true,
    copyable: true,
  },
  {
    key: 'eventType',
    label: t.columns.eventType.value,
    type: 'badge',
    visible: 'always',
    width: 100,
    format: (value: string) => value === 'offline' ? t.statusBadge.offline.value : t.statusBadge.online.value,
    color: (value: string) => getEventTypeColor(value),
  },
  {
    key: 'lastSeenAt',
    label: t.columns.lastSeenAt.value,
    type: 'text',
    visible: 'laptop',
    width: 180,
    format: (value: string) => value ? formatTimestamp(value) : t.defaults.notAvailable.value,
  },
  {
    key: 'missCount',
    label: t.columns.missCount.value,
    type: 'text',
    visible: 'laptop',
    width: 140,
    format: (value: number) => value != null ? String(value) : t.defaults.notAvailable.value,
  },
  {
    key: 'thresholdMinutes',
    label: t.columns.thresholdMinutes.value,
    type: 'text',
    visible: 'laptop',
    width: 140,
    format: (value: number) => value != null ? String(value) : t.defaults.notAvailable.value,
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
    if (col.key === 'asset') return columnVisibility.value.asset;
    if (col.key === 'assetUUID') return columnVisibility.value.assetUUID;
    if (col.key === 'eventType') return columnVisibility.value.eventType;
    if (col.key === 'lastSeenAt') return columnVisibility.value.lastSeenAt;
    if (col.key === 'missCount') return columnVisibility.value.missCount;
    if (col.key === 'thresholdMinutes') return columnVisibility.value.thresholdMinutes;
    if (col.key === 'created') return columnVisibility.value.created;
    return true;
  });
});

/** WATCHERS */

useOrgChangeRefresh(() => {
  resetCursor();
  quickSearchAssetUUID.value = '';
  quickEventType.value = null;
  const resetRange = getTodayDateRange();
  advancedFilterValues.value = {
    includeChildren: null,
    startDate: resetRange.from.split('T')[0],
    endDate: resetRange.to.split('T')[0],
    eventType: null,
    assetUUID: '',
  };
  filters.value = {
    from: resetRange.from,
    to: resetRange.to,
  };
  hasPendingAdvancedFilters.value = false;
  void fetchData();
});

/** FUNCTIONS */

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
 * Fetch connectivity events from API with current filters + cursor pagination.
 */
async function fetchData(direction: 'next' | 'prev' = 'next'): Promise<void> {
  loading.value = true;

  try {
    const queryParams: Record<string, any> = {
      limit: limit.value,
      direction,
    };

    if (direction === 'next' && cursor.value.next) {
      queryParams.cursor = cursor.value.next;
    } else if (direction === 'prev' && cursor.value.prev) {
      queryParams.cursor = cursor.value.prev;
    }

    if (filters.value.assetUUID) {
      queryParams.assetUUID = filters.value.assetUUID;
    }
    if (filters.value.eventType) {
      queryParams.eventType = filters.value.eventType;
    }
    if (filters.value.from) {
      queryParams.from = filters.value.from;
    }
    if (filters.value.to) {
      queryParams.to = filters.value.to;
    }
    if (typeof filters.value.includeChildren === 'boolean') {
      queryParams.includeChildren = filters.value.includeChildren;
    }

    logger.debug('Fetching with params', queryParams);

    const response = await apis.events.events.listConnectivityHistory(queryParams);

    logger.debug('Response received', response);

    logsList.value = response.items || [];

    cursor.value = {
      current: queryParams.cursor,
      next: response.nextCursor ? String(response.nextCursor) : undefined,
      prev: response.prevCursor ? String(response.prevCursor) : undefined,
      hasNext: response.hasNext,
      hasPrevious: response.hasPrevious,
    };

  } catch (error: any) {
    logger.error('Error fetching connectivity events', {
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

async function refreshData(): Promise<void> {
  resetCursor();
  await fetchData('next');
}

function applyQuickFilters(): void {
  if (quickSearchAssetUUID.value) {
    filters.value.assetUUID = quickSearchAssetUUID.value;
  } else {
    delete filters.value.assetUUID;
  }
  if (quickEventType.value !== null) {
    filters.value.eventType = quickEventType.value;
  } else {
    delete filters.value.eventType;
  }
  resetCursor();
  void fetchData();
}

function handleAdvancedFiltersApply(appliedFilters: Record<string, any>): void {
  advancedFilterValues.value = { ...appliedFilters };
  if (appliedFilters.includeChildren !== null && appliedFilters.includeChildren !== undefined) {
    filters.value.includeChildren = appliedFilters.includeChildren;
  } else {
    delete filters.value.includeChildren;
  }
  if (appliedFilters.startDate) {
    const startDate = new Date(appliedFilters.startDate);
    startDate.setHours(0, 0, 0, 0);
    filters.value.from = startDate.toISOString();
  } else {
    const today = getTodayDateRange();
    filters.value.from = today.from;
  }
  if (appliedFilters.endDate) {
    const endDate = new Date(appliedFilters.endDate);
    endDate.setHours(23, 59, 59, 999);
    filters.value.to = endDate.toISOString();
  } else {
    const today = getTodayDateRange();
    filters.value.to = today.to;
  }
  if (appliedFilters.eventType) {
    filters.value.eventType = appliedFilters.eventType;
  } else if (quickEventType.value) {
    filters.value.eventType = quickEventType.value;
  } else {
    delete filters.value.eventType;
  }
  if (appliedFilters.assetUUID || quickSearchAssetUUID.value) {
    filters.value.assetUUID = appliedFilters.assetUUID || quickSearchAssetUUID.value;
  } else {
    delete filters.value.assetUUID;
  }
  hasPendingAdvancedFilters.value = false;
  resetCursor();
  showAdvancedFilters.value = false;
  void fetchData();
}

function handlePendingChange(hasPending: boolean): void {
  hasPendingAdvancedFilters.value = hasPending;
}

function removeFilter(key: string): void {
  if (key === 'quickSearchAssetUUID') {
    quickSearchAssetUUID.value = '';
    if (advancedFilterValues.value.assetUUID) {
      filters.value.assetUUID = advancedFilterValues.value.assetUUID;
    } else {
      delete filters.value.assetUUID;
    }
  } else if (key === 'eventType') {
    quickEventType.value = null;
    advancedFilterValues.value.eventType = null;
    delete filters.value.eventType;
  } else if (key === 'includeChildren') {
    advancedFilterValues.value.includeChildren = null;
    delete filters.value.includeChildren;
  } else if (key === 'assetUUID') {
    advancedFilterValues.value.assetUUID = '';
    if (quickSearchAssetUUID.value) {
      filters.value.assetUUID = quickSearchAssetUUID.value;
    } else {
      delete filters.value.assetUUID;
    }
  }
  resetCursor();
  void fetchData();
}

function clearAllFilters(): void {
  quickSearchAssetUUID.value = '';
  quickEventType.value = null;
  const resetRange = getTodayDateRange();
  advancedFilterValues.value = {
    includeChildren: null,
    startDate: resetRange.from.split('T')[0],
    endDate: resetRange.to.split('T')[0],
    eventType: null,
    assetUUID: '',
  };
  filters.value = {
    from: resetRange.from,
    to: resetRange.to,
  };
  hasPendingAdvancedFilters.value = false;
  resetCursor();
  void fetchData();
}

function goToNextPage(): void {
  if (cursor.value.hasNext) {
    void fetchData('next');
  }
}

function goToPrevPage(): void {
  if (cursor.value.hasPrevious) {
    void fetchData('prev');
  }
}

function handleColumnsUpdate(columns: ListHeaderMenuColumn[]): void {
  menuColumns.value = columns;
}

function handleLimitChange(newLimit: number): void {
  limit.value = newLimit;
  resetCursor();
  void fetchData();
}

function cardClick(item: AssetConnectivityEvent): void {
  selectedEvent.value = item;
  jsonDrawerOpen.value = true;
}

function handleView(item: AssetConnectivityEvent): void {
  selectedEvent.value = item;
  jsonDrawerOpen.value = true;
}

function handleAction(actionKey: string, item: AssetConnectivityEvent): void {
  if (actionKey === 'viewAsset') {
    // Navigate to the Asset list, filtered by UUID (deep-link via query string).
    void router.push({
      path: '/assets',
      query: { assetUUID: item.assetUUID },
    });
  }
}

function getEventTypeColor(eventType: string): string {
  return EVENT_TYPE_COLORS[eventType] || DEFAULT_COLOR;
}

function formatTimestamp(value: string | Date): string {
  if (!value) return t.defaults.notAvailable.value;
  const date = typeof value === 'string' ? new Date(value) : value;
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
      icon="wifi"
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
          v-model="quickSearchAssetUUID"
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

      <!-- Event Type Select -->
      <div class="col-auto" style="min-width: 160px;">
        <q-select
          v-model="quickEventType"
          outlined
          dense
          emit-value
          map-options
          :options="eventTypeOptions"
          :label="t.filters.eventType.value"
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
          <q-icon name="wifi" size="sm" color="primary" class="q-mr-sm"/>
          <div class="text-subtitle1 text-weight-medium text-primary">{{ t.listTitle.value }}</div>
        </div>
      </div>
      <div class="col-auto">
        <ListHeaderMenu
          :items-count="logsList.length"
          :item-label="t.itemLabel.value"
          :item-label-plural="t.itemLabelPlural.value"
          icon="wifi"
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
        :key="log.eventId + '-' + index"
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
        icon="wifi"
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
      :subtitle="`${selectedEvent.assetName || selectedEvent.assetUUID} • ${formatTimestamp(selectedEvent.created)}`"
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
