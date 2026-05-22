<script setup lang="ts">
defineOptions({
  name: 'DlqLogsPage'
});

/** TYPE IMPORTS */
import type { EventsDLQResponse, EventsDLQServiceCount } from '@mapexos/schemas';
import type { DlqServiceTypeGroup } from './interfaces';

/** VUE IMPORTS */
import { ref, computed, onMounted, watch } from 'vue';

/** COMPONENTS */
import { DlqSidebar, DlqListView, DlqDetailView } from './components';
import { ListCardEmpty } from '@components/cards';

/** COMPOSABLES */
import { useOrgChangeRefresh } from '@composables/organizations';
import { useDlqLogsPageTranslations } from '@composables/i18n/pages/logs/dlqLogsPage';
import { useLogger } from '@composables/useLogger';

/** UTILS */
import { notifyFail } from '@utils/alert';

/** SERVICES */
import { apis } from '@services/mapex';

/** CONSTANTS */
const DEFAULT_LIMIT = 20;
const FIXED_SERVICE_TYPES: Record<string, string> = {
  workflow: 'account_tree',
  triggers: 'flash_on',
  router: 'route',
  events: 'event_note',
  assets: 'devices',
  'mapex-iam': 'admin_panel_settings',
  'http-gateway': 'http',
  'js-executor': 'code',
  'js-workflow-executor': 'terminal',
};

/** COMPOSABLES & STORES */
const t = useDlqLogsPageTranslations();
const logger = useLogger('DlqLogsPage');

/** STATE */
const initialLoading = ref(true);
const allEntries = ref<EventsDLQResponse[]>([]);
const selectedEntry = ref<EventsDLQResponse | null>(null);
const showDetail = ref(false);
const activeServiceType = ref<string | null>(null);
const countsFromApi = ref<EventsDLQServiceCount[]>([]);
const totalFromApi = ref(0);
const datePreset = ref<'today' | '7d' | '30d' | 'all'>('today');
const searchQuery = ref('');
const scrollKey = ref(0);

/** Cursor */
const nextCursor = ref<string | undefined>(undefined);
const hasMore = ref(true);

/** COMPUTED */

/**
 * Date presets with i18n labels
 */
const datePresets = computed(() => [
  { value: 'today' as const, label: t.datePresets.today.value },
  { value: '7d' as const, label: '7d' },
  { value: '30d' as const, label: '30d' },
  { value: 'all' as const, label: t.datePresets.all.value },
]);

/**
 * Get date range based on preset
 * @returns {{ startTime?: string; endTime?: string }}
 */
const dateRange = computed(() => {
  const now = new Date();
  const endOfDay = new Date(now.getFullYear(), now.getMonth(), now.getDate(), 23, 59, 59, 999);

  switch (datePreset.value) {
    case 'today': {
      const startOfDay = new Date(now.getFullYear(), now.getMonth(), now.getDate(), 0, 0, 0, 0);
      return { startTime: startOfDay.toISOString(), endTime: endOfDay.toISOString() };
    }
    case '7d': {
      const start = new Date(now);
      start.setDate(start.getDate() - 7);
      start.setHours(0, 0, 0, 0);
      return { startTime: start.toISOString(), endTime: endOfDay.toISOString() };
    }
    case '30d': {
      const start = new Date(now);
      start.setDate(start.getDate() - 30);
      start.setHours(0, 0, 0, 0);
      return { startTime: start.toISOString(), endTime: endOfDay.toISOString() };
    }
    case 'all':
    default:
      return {};
  }
});

/**
 * Build fixed service type groups with counts from API
 */
const serviceTypeGroups = computed<DlqServiceTypeGroup[]>(() => {
  const countMap = new Map<string, number>();
  for (const c of countsFromApi.value) {
    countMap.set(c.serviceType, c.count);
  }
  return Object.entries(FIXED_SERVICE_TYPES).map(([serviceType, icon]) => ({
    serviceType,
    icon,
    count: countMap.get(serviceType) || 0,
  }));
});

/**
 * Header title based on active filter
 */
const headerTitle = computed(() => {
  return activeServiceType.value || t.sidebar.allFailures.value;
});

/**
 * Header badge count
 */
const headerCount = computed(() => {
  if (activeServiceType.value) {
    const group = countsFromApi.value.find(c => c.serviceType === activeServiceType.value);
    return group?.count || allEntries.value.length;
  }
  return totalFromApi.value;
});

/** WATCHERS */

/**
 * Re-fetch data when organization changes
 */
useOrgChangeRefresh(() => {
  void resetAndReload();
});

/**
 * When date preset changes, reload everything
 */
watch(datePreset, () => {
  void resetAndReload();
});

/** FUNCTIONS */

/**
 * Reset all state and reload
 * @returns {Promise<void>}
 */
async function resetAndReload(): Promise<void> {
  selectedEntry.value = null;
  showDetail.value = false;
  allEntries.value = [];
  nextCursor.value = undefined;
  hasMore.value = true;
  scrollKey.value++;
  await fetchCounts();
  initialLoading.value = false;
}

/**
 * Fetch DLQ counts from API grouped by service type
 * @returns {Promise<void>}
 */
async function fetchCounts(): Promise<void> {
  try {
    const queryParams: Record<string, any> = {};
    if (dateRange.value.startTime) queryParams.startTime = dateRange.value.startTime;
    if (dateRange.value.endTime) queryParams.endTime = dateRange.value.endTime;

    const response = await apis.events.events.getDLQCounts(queryParams);
    countsFromApi.value = response.counts || [];
    totalFromApi.value = response.total || 0;
  } catch (error: any) {
    logger.error('Error fetching DLQ counts', error);
  }
}

/**
 * Fetch next page — called by q-infinite-scroll via DlqListView
 * Handles both first page and subsequent pages.
 * @param {Function} done - Quasar callback: done() to continue, done(true) to stop
 * @returns {Promise<void>}
 */
async function fetchPage(done: (stop?: boolean) => void): Promise<void> {
  try {
    const params: Record<string, any> = {
      limit: DEFAULT_LIMIT,
      direction: 'next',
    };

    if (nextCursor.value) {
      params.cursor = normalizeCursorDate(nextCursor.value);
    }
    if (activeServiceType.value) params.serviceType = activeServiceType.value;
    if (searchQuery.value) params.lastError = searchQuery.value;
    if (dateRange.value.startTime) params.startTime = dateRange.value.startTime;
    if (dateRange.value.endTime) params.endTime = dateRange.value.endTime;

    const response = await apis.events.events.listDLQ(params);
    const newItems = response.items || [];

    allEntries.value = [...allEntries.value, ...newItems];
    nextCursor.value = response.nextCursor ? String(response.nextCursor) : undefined;
    hasMore.value = response.hasNext;

    done(!response.hasNext);
  } catch (error: any) {
    logger.error('Error fetching DLQ events', error);
    notifyFail({
      message: error?.response?.data?.message || t.messages.loadFailed.value,
      timeout: 5000,
    });
    done(true);
  }
}

/**
 * Normalize cursor date to ISO 8601 with exactly 3ms digits.
 * Go/ClickHouse can return microsecond precision (6+ digits) which fails Zod validation.
 * @param {string} cursor - Raw cursor string from API
 * @returns {string} Normalized ISO string (YYYY-MM-DDTHH:mm:ss.SSSZ)
 */
function normalizeCursorDate(cursor: string): string {
  const d = new Date(cursor);
  if (isNaN(d.getTime())) return cursor;
  return d.toISOString();
}

/**
 * Handle service type change from sidebar
 * @param {string | null} serviceType - Selected service type or null for all
 */
function handleServiceTypeChange(serviceType: string | null): void {
  closeDetail();
  activeServiceType.value = serviceType;
  allEntries.value = [];
  nextCursor.value = undefined;
  hasMore.value = true;
  scrollKey.value++;
}

/**
 * Handle search input with debounce
 * @param {string} value - Search input value
 */
let searchTimer: ReturnType<typeof setTimeout> | null = null;
function handleSearchInput(value: string): void {
  searchQuery.value = value;
  if (searchTimer) clearTimeout(searchTimer);
  searchTimer = setTimeout(() => {
    allEntries.value = [];
    nextCursor.value = undefined;
    hasMore.value = true;
    scrollKey.value++;
  }, 400);
}

/**
 * Handle entry selection from list
 * @param {EventsDLQResponse} entry - Selected DLQ entry
 */
function handleSelect(entry: EventsDLQResponse): void {
  selectedEntry.value = entry;
  showDetail.value = true;
}

/**
 * Close detail view — back to list
 */
function closeDetail(): void {
  selectedEntry.value = null;
  showDetail.value = false;
}

/**
 * Handle reload — refresh everything
 */
function handleReload(): void {
  void resetAndReload();
}

/** LIFECYCLE HOOKS */
onMounted(async () => {
  await fetchCounts();
  initialLoading.value = false;
});
</script>

<template>
  <q-page class="dlq-page">
    <div class="dlq-page__layout row no-wrap">

      <!-- Sidebar -->
      <DlqSidebar
        :active-service-type="activeServiceType"
        :service-type-groups="serviceTypeGroups"
        :total-count="totalFromApi"
        @update:active-service-type="handleServiceTypeChange"
      />

      <!-- Content Area -->
      <div class="col dlq-page__content">

        <!-- Detail View (full-width, replaces list) -->
        <DlqDetailView
          v-if="showDetail && selectedEntry"
          :entry="selectedEntry"
          @close="closeDetail"
        />

        <!-- List View -->
        <template v-else>
          <!-- Header Bar -->
          <div class="dlq-page__header">
            <div class="dlq-page__header-left">
              <q-icon name="report_problem" size="sm" class="dlq-page__header-icon" />
              <span class="dlq-page__header-title">{{ headerTitle }}</span>
              <q-badge rounded :label="headerCount" class="dlq-page__header-badge" />
            </div>

            <div class="dlq-page__header-center">
              <div class="dlq-page__date-pills">
                <button
                  v-for="preset in datePresets"
                  :key="preset.value"
                  class="dlq-page__date-pill"
                  :class="{ 'dlq-page__date-pill--active': datePreset === preset.value }"
                  @click="datePreset = preset.value"
                >
                  {{ preset.label }}
                </button>
              </div>
            </div>

            <div class="dlq-page__header-right">
              <q-input
                :model-value="searchQuery"
                outlined
                dense
                rounded
                clearable
                :placeholder="t.sidebar.searchErrorPlaceholder.value"
                class="dlq-page__search"
                @update:model-value="handleSearchInput($event as string)"
                @clear="handleSearchInput('')"
              >
                <template #prepend>
                  <q-icon name="search" size="xs" />
                </template>
              </q-input>

              <q-btn
                flat
                round
                dense
                icon="refresh"
                class="dlq-page__refresh-btn"
                @click="handleReload"
              >
                <q-tooltip>{{ t.actions.refresh.value }}</q-tooltip>
              </q-btn>
            </div>
          </div>

          <!-- Message List Card -->
          <q-card flat bordered class="dlq-page__list-card">
            <!-- Initial Loading -->
            <div v-if="initialLoading" class="dlq-page__loading column items-center justify-center">
              <q-spinner color="primary" size="50px" />
            </div>

            <!-- Empty State -->
            <div v-else-if="allEntries.length === 0 && !hasMore" class="dlq-page__empty-wrapper">
              <ListCardEmpty
                :title="t.empty.title.value"
                :description="t.empty.description.value"
                icon="error_outline"
              />
            </div>

            <!-- Entries List -->
            <DlqListView
              v-else
              :key="scrollKey"
              :entries="allEntries"
              :selected-id="null"
              :has-more="hasMore"
              @select="handleSelect"
              @load-more="(done: (stop?: boolean) => void) => fetchPage(done)"
            />
          </q-card>
        </template>

      </div>

    </div>
  </q-page>
</template>

<style lang="scss" scoped>
.dlq-page {
  &__layout {
    background: var(--mapex-page-bg);
    height: calc(100vh - 80px);
  }

  &__content {
    display: flex;
    flex-direction: column;
    overflow: hidden;
  }

  &__header {
    display: flex;
    align-items: center;
    gap: var(--mapex-spacing-lg);
    padding: var(--mapex-spacing-md) var(--mapex-spacing-lg);
    background: var(--mapex-surface-elevated);
    border-bottom: 1px solid var(--mapex-card-border);
    flex-shrink: 0;
  }

  &__header-left {
    display: flex;
    align-items: center;
    gap: var(--mapex-spacing-sm);
    flex-shrink: 0;
  }

  &__header-center {
    flex: 1;
    display: flex;
    justify-content: center;
  }

  &__header-right {
    display: flex;
    align-items: center;
    gap: var(--mapex-spacing-sm);
    flex-shrink: 0;
  }

  &__header-icon {
    color: var(--mapex-danger);
  }

  &__header-title {
    font-size: var(--mapex-font-lg);
    font-weight: var(--mapex-font-weight-bold);
    color: var(--mapex-text-primary);
  }

  &__header-badge {
    background: var(--mapex-text-muted);
    color: var(--mapex-surface-bg);
  }

  &__date-pills {
    display: flex;
    background: var(--mapex-surface-sunken);
    border-radius: var(--mapex-radius-sm);
    padding: 2px;
    gap: 2px;
  }

  &__date-pill {
    all: unset;
    cursor: pointer;
    padding: 4px 14px;
    font-size: var(--mapex-font-xs);
    font-weight: var(--mapex-font-weight-medium);
    color: var(--mapex-text-muted);
    border-radius: var(--mapex-radius-xs);
    transition: var(--mapex-transition-fast);
    white-space: nowrap;

    &:hover:not(&--active) {
      color: var(--mapex-text-primary);
      background: var(--mapex-surface-highlight);
    }

    &--active {
      background: var(--mapex-surface-elevated);
      color: var(--mapex-primary);
      box-shadow: var(--mapex-shadow-xs);
    }
  }

  &__search {
    width: 200px;

    :deep(.q-field__control) {
      height: 32px;
    }
  }

  &__refresh-btn {
    color: var(--mapex-text-secondary);
  }

  &__list-card {
    flex: 1;
    min-height: 0;
    display: flex;
    flex-direction: column;
    overflow: hidden;
    border-radius: 0;
    border-left: none;
    border-right: none;
    border-bottom: none;
    background: var(--mapex-surface-bg);
  }

  &__loading {
    flex: 1;
    min-height: 400px;
  }

  &__empty-wrapper {
    flex: 1;
    display: flex;
    align-items: center;
    justify-content: center;
    padding: var(--mapex-spacing-lg);

    :deep(.empty-state-card) {
      max-width: 520px;
      margin-top: 0;
    }
  }
}
</style>
