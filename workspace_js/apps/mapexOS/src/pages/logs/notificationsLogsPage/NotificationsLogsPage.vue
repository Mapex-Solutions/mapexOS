<script setup lang="ts">
defineOptions({
  name: 'NotificationsLogsPage'
});

/** TYPE IMPORTS */
import type { ListHeaderMenuColumn } from '@components/headers';
import type { DataRowColumn } from '@components/cards';
import type { FilterField, FilterValues } from '@components/drawers';

/** VUE IMPORTS */
import { ref, computed, onMounted } from 'vue';

/** COMPONENTS */
import { PageHeader, ListHeaderMenu } from '@components/headers';
import { ListCardEmpty, DataRow } from '@components/cards';
import { JsonDrawer, AdvancedFiltersDrawer } from '@components/drawers';
import { AppTooltip } from '@components/tooltips';

/** COMPOSABLES */
import { useNotificationsLogsPageTranslations } from '@composables/i18n/pages/logs/notificationsLogsPage';
import { useLogger } from '@composables/useLogger';

/** LOCAL IMPORTS */
import { NOTIFICATIONS_LOG_LIST_STUB } from './stubs';

/** CONSTANTS */
const MAX_VISIBLE_CHIPS = 2;

const NOTIFICATION_ICONS: Record<string, string> = {
  slack: 'mdi-slack',
  teams: 'mdi-microsoft-teams',
  email: 'mdi-email',
  push: 'mdi-bell-ring',
  telegram: 'mdi-send',
  webhook: 'mdi-webhook',
};

const NOTIFICATION_COLORS: Record<string, string> = {
  slack: 'purple-6',
  teams: 'blue-6',
  email: 'grey-6',
  push: 'orange-6',
  telegram: 'cyan-6',
  webhook: 'indigo-6',
};


/** LOCAL IMPORTS */
import type { NotificationsLogsPageFilters } from './interfaces/NotificationsLogsPage.interface';

/** COMPOSABLES & STORES */
const t = useNotificationsLogsPageTranslations();
const logger = useLogger('NotificationsLogsPage');

/** STATE */
const loading = ref(false);
const lastUpdatedAt = ref<number | undefined>(undefined);
const logsList = ref<any[]>([]);
const selectedEvent = ref<any>(null);
const jsonDrawerOpen = ref(false);

/** Pagination */
const itemsPerPage = ref(15);
const currentPage = ref(1);

/** Filters */
const filters = ref<NotificationsLogsPageFilters>({});

/** FILTER STATE - Enterprise Filter Pattern */
const showFiltersDrawer = ref(false);
const quickSearchName = ref('');
const quickStatusSuccess = ref<boolean | null>(null);
const advancedFilterValues = ref<FilterValues>({
  includeChildren: null,
  notificationType: null,
});
const hasPendingAdvancedFilters = ref(false);

/** Column visibility using ListHeaderMenuColumn format */
const menuColumns = ref<ListHeaderMenuColumn[]>([
  { key: 'notificationType', label: 'Type', visible: true },
  { key: 'tenantId', label: 'Tenant', visible: true },
  { key: 'created', label: 'Created', visible: true },
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
    key: 'notificationType',
    type: 'select',
    label: t.filters.notificationType.value,
    icon: 'notifications',
    options: [
      { label: t.notificationTypeOptions.slack.value, value: 'slack' },
      { label: t.notificationTypeOptions.teams.value, value: 'teams' },
      { label: t.notificationTypeOptions.email.value, value: 'email' },
      { label: t.notificationTypeOptions.push.value, value: 'push' },
      { label: t.notificationTypeOptions.telegram.value, value: 'telegram' },
      { label: t.notificationTypeOptions.webhook.value, value: 'webhook' },
    ],
  },
]);

/**
 * Check if any filter is active
 */
const hasActiveFilters = computed(() => {
  return !!(
    filters.value.notificationName ||
    filters.value.status !== undefined ||
    filters.value.includeChildren !== undefined ||
    filters.value.notificationType
  );
});

/**
 * Count of advanced filters only (for badge on filter icon)
 */
const advancedFiltersCount = computed(() => {
  let count = 0;
  if (filters.value.includeChildren !== undefined) count++;
  if (filters.value.notificationType) count++;
  return count;
});

/**
 * Active filter chips for display
 */
const activeFilterChips = computed(() => {
  const chips: Array<{ key: string; label: string; value: string }> = [];

  if (filters.value.notificationName) {
    chips.push({ key: 'notificationName', label: t.filters.name.value, value: filters.value.notificationName });
  }
  if (filters.value.status !== undefined) {
    chips.push({
      key: 'status',
      label: t.filters.status.value,
      value: filters.value.status ? t.filters.options.success.value : t.filters.options.failure.value
    });
  }
  if (filters.value.includeChildren !== undefined) {
    chips.push({
      key: 'includeChildren',
      label: t.filters.includeChildren.value,
      value: filters.value.includeChildren ? t.filters.options.yes.value : t.filters.options.no.value
    });
  }
  if (filters.value.notificationType) {
    const typeKey = filters.value.notificationType as keyof typeof t.notificationTypeOptions;
    const typeLabel = t.notificationTypeOptions[typeKey]?.value || filters.value.notificationType;
    chips.push({ key: 'notificationType', label: t.filters.notificationType.value, value: typeLabel });
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
  notificationType: menuColumns.value.find((col) => col.key === 'notificationType')?.visible ?? true,
  tenantId: menuColumns.value.find((col) => col.key === 'tenantId')?.visible ?? true,
  created: menuColumns.value.find((col) => col.key === 'created')?.visible ?? true,
}));

/**
 * Log columns configuration for DataRow
 */
const logColumns = computed((): DataRowColumn[] => [
  {
    key: 'icon',
    label: '',
    type: 'avatar',
    visible: 'always',
    width: 56,
    icon: (value: any, row: any) => getNotificationIcon(row.notificationType),
    color: (value: any, row: any) => getNotificationIconColor(row.notificationType),
  },
  {
    key: 'notificationName',
    label: 'Notification Name',
    type: 'text',
    visible: 'always',
    width: 250,
    ellipsis: true,
    secondaryKey: 'details',
  },
  {
    key: 'notificationType',
    label: 'Type',
    type: 'chip',
    visible: 'laptop',
    width: 130,
    format: (value: any): string => value ? String(value).toUpperCase() : 'UNKNOWN',
    color: (value: any) => getNotificationIconColor(value),
  },
  {
    key: 'tenantId',
    label: 'Tenant',
    type: 'text',
    visible: 'laptop',
    width: 160,
    ellipsis: true,
  },
  {
    key: 'status',
    label: 'Status',
    type: 'badge',
    visible: 'laptop',
    width: 110,
    format: (value: any): string => value === 'success' ? 'SUCCESS' : 'FAILURE',
    color: (value: any): string => value === 'success' ? 'green-6' : 'red-6',
  },
  {
    key: 'created',
    label: 'Created',
    type: 'text',
    visible: 'laptop',
    width: 180,
    format: (value: any): string => formatTimestamp(value),
  },
]);

/**
 * Filtered columns based on visibility settings
 */
const visibleColumns = computed(() => {
  return logColumns.value.filter((col) => {
    if (col.key === 'icon' || col.key === 'notificationName' || col.key === 'status') return true;
    if (col.key === 'notificationType') return columnVisibility.value.notificationType;
    if (col.key === 'tenantId') return columnVisibility.value.tenantId;
    if (col.key === 'created') return columnVisibility.value.created;
    return true;
  });
});

/** FUNCTIONS */

/**
 * Fetch notification logs data
 */
function fetchData(): void {
  loading.value = true;

  try {
    // TODO: Replace with actual API call when available
    // For now using stub data with client-side filtering
    let filteredData = [...NOTIFICATIONS_LOG_LIST_STUB];

    // Apply filters
    if (filters.value.notificationName) {
      const searchTerm = filters.value.notificationName.toLowerCase();
      filteredData = filteredData.filter(log =>
        log.notificationName?.toLowerCase().includes(searchTerm)
      );
    }
    if (filters.value.status !== undefined) {
      const statusValue = filters.value.status ? 'success' : 'failure';
      filteredData = filteredData.filter(log => log.status === statusValue);
    }
    if (filters.value.notificationType) {
      filteredData = filteredData.filter(log => log.notificationType === filters.value.notificationType);
    }
    // Note: includeChildren would be handled by API in real implementation

    logsList.value = filteredData;
    logger.debug('Fetched notification logs', { count: filteredData.length });
  } catch (error: any) {
    logger.error('Error fetching notification logs', error);
    logsList.value = [];
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
    filters.value.notificationName = quickSearchName.value;
  } else {
    delete filters.value.notificationName;
  }
  if (quickStatusSuccess.value !== null) {
    filters.value.status = quickStatusSuccess.value;
  } else {
    delete filters.value.status;
  }
  currentPage.value = 1;
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
  if (values.notificationType) {
    filters.value.notificationType = values.notificationType;
  } else {
    delete filters.value.notificationType;
  }

  // Reset pending state after apply
  hasPendingAdvancedFilters.value = false;

  currentPage.value = 1;
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
    notificationType: null,
  };
}

/**
 * Remove individual filter
 * @param {string} key - Filter key to remove
 */
function removeFilter(key: string): void {
  if (key === 'notificationName') {
    delete filters.value.notificationName;
    quickSearchName.value = '';
  } else if (key === 'status') {
    delete filters.value.status;
    quickStatusSuccess.value = null;
  } else if (key === 'includeChildren') {
    delete filters.value.includeChildren;
    advancedFilterValues.value.includeChildren = null;
  } else if (key === 'notificationType') {
    delete filters.value.notificationType;
    advancedFilterValues.value.notificationType = null;
  }

  currentPage.value = 1;
  void fetchData();
}

/**
 * Clear all filters
 */
function clearAllFilters(): void {
  // Reset filter state
  filters.value = {};

  // Reset quick filters
  quickSearchName.value = '';
  quickStatusSuccess.value = null;

  // Reset advanced filters
  advancedFilterValues.value = {
    includeChildren: null,
    notificationType: null,
  };

  // Reset pending state
  hasPendingAdvancedFilters.value = false;

  currentPage.value = 1;
  void fetchData();
}

/**
 * Update menu columns when changed
 * @param {ListHeaderMenuColumn[]} columns - Updated column visibility states
 */
function handleColumnsUpdate(columns: ListHeaderMenuColumn[]): void {
  menuColumns.value = columns;
}

/**
 * Handle items per page change
 * @param {number} newValue - New items per page value
 */
function handleItemsPerPageChange(newValue: number): void {
  itemsPerPage.value = newValue;
  currentPage.value = 1;
  void fetchData();
}

/**
 * Handle card click event - opens JSON drawer
 * @param {any} item - Clicked log item
 */
function cardClick(item: any): void {
  selectedEvent.value = item;
  jsonDrawerOpen.value = true;
}

/**
 * Get icon for notification type
 * @param {string} type - Notification type
 * @returns {string} Icon name
 */
function getNotificationIcon(type: string): string {
  return NOTIFICATION_ICONS[type?.toLowerCase()] || 'mdi-bell-outline';
}

/**
 * Get color for notification type
 * @param {string} type - Notification type
 * @returns {string} Color class name
 */
function getNotificationIconColor(type: string): string {
  return NOTIFICATION_COLORS[type?.toLowerCase()] || 'grey-6';
}

/**
 * Format timestamp for display
 * @param {string} value - ISO timestamp string
 * @returns {string} Formatted date string
 */
function formatTimestamp(value: string): string {
  if (!value) return 'N/A';
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
      icon="notifications"
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
          <q-icon name="notifications" size="sm" color="primary" class="q-mr-sm"/>
          <div class="text-subtitle1 text-weight-medium text-primary">{{ t.listTitle.value }}</div>
        </div>
      </div>
      <div class="col-auto">
        <ListHeaderMenu
          :item-label="t.itemLabel.value"
          :item-label-plural="t.itemLabelPlural.value"
          icon="notifications"
          :items-count="logsList.length"
          :items-per-page="itemsPerPage"
          :columns="menuColumns"
          :filtered="hasActiveFilters"
          :refreshing="loading"
          :last-updated-at="lastUpdatedAt"
          @update:items-per-page="handleItemsPerPageChange"
          @update:columns="handleColumnsUpdate"
          @refresh="fetchData"
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
        v-for="log in logsList"
        :key="log.id"
        class="col-12 q-mb-xs"
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
        icon="notifications"
      />
    </div>

    <!-- Pagination -->
    <div class="row justify-center q-mt-lg q-mb-lg">
      <q-pagination
        v-if="logsList.length > 0"
        v-model="currentPage"
        direction-links
        boundary-links
        class="rounded-borders"
        color="primary"
        active-color="primary"
        :max="Math.ceil(logsList.length / itemsPerPage) || 1"
      />
    </div>

    <!-- JSON Drawer -->
    <JsonDrawer
      v-if="selectedEvent"
      v-model:show="jsonDrawerOpen"
      :title="t.drawer.title.value"
      :jsonData="selectedEvent"
      :editable="false"
      :subtitle="`${selectedEvent.notificationType?.toUpperCase() || t.defaults.notAvailable.value} • ${formatTimestamp(selectedEvent.created)}`"
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
