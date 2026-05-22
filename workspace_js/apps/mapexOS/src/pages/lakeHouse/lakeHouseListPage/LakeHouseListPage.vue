<script setup lang="ts">
defineOptions({
  name: 'LakeHouseListPage'
});

/** TYPE IMPORTS */
import type { ListHeaderMenuColumn } from '@components/headers';
import type { DataRowColumn } from '@components/cards';

/** VUE IMPORTS */
import { ref, computed, watch } from 'vue';

/** COMPONENTS */
import { PageHeader, ListHeaderMenu } from '@components/headers';
import { ListCardEmpty, DataRow } from '@components/cards';
import { ListPagination } from '@components/navigation';
import { AdvancedFiltersDrawer } from '@components/drawers';
import { AppTooltip } from '@components/tooltips';

/** COMPOSABLES */
import { useLakeHouseTranslations } from '@composables/i18n/pages/lakeHouse/useLakeHouseTranslations';
import { usePermissions } from '@composables/shared/usePermissions';
import { useLogger } from '@composables/useLogger';

/** LOCAL IMPORTS */
import { LAKE_HOUSE_LIST_STUB } from './stubs';

/** CONSTANTS */
const MAX_VISIBLE_CHIPS = 2;

/** COMPOSABLES & STORES */
const translations = useLakeHouseTranslations();
const logger = useLogger('LakeHouseListPage');
const { canUpdate, canRead } = usePermissions();
const canUpdateRetention = canUpdate('retention');
const canReadRetention = canRead('retention');

/** STATE */
const lakeHousesList = ref<any[]>([]);
const channels = ref<any[]>(LAKE_HOUSE_LIST_STUB);
const loading = ref(false);
const lastUpdatedAt = ref<number | undefined>(undefined);
lakeHousesList.value = [...channels.value];

// Pagination
const itemsPerPage = ref(15);
const currentPage = ref(1);
const totalPages = ref(1);

// Quick filters state
const quickSearch = ref('');
const quickStatus = ref<boolean | null>(null);

// Advanced filters state
const showAdvancedFilters = ref(false);
const advancedFilterValues = ref<Record<string, any>>({
  includeChildren: null,
  type: null,
  region: '',
});
const hasPendingAdvancedFilters = ref(false);

// Column visibility
const menuColumns = ref<ListHeaderMenuColumn[]>(translations.menuColumns.value);

/** COMPUTED */

/**
 * Status toggle options for quick filter
 */
const statusOptions = computed(() => [
  { label: translations.filters.allStatus.value, value: null },
  { label: translations.filters.options.active.value, value: true },
  { label: translations.filters.options.inactive.value, value: false },
]);

/**
 * Count of active advanced filters
 */
const advancedFiltersCount = computed(() => {
  let count = 0;
  if (advancedFilterValues.value.includeChildren !== null) count++;
  if (advancedFilterValues.value.type !== null) count++;
  if (advancedFilterValues.value.region) count++;
  return count;
});

/**
 * Active filter chips for display
 */
const activeFilterChips = computed(() => {
  const chips: Array<{ key: string; label: string; value: string }> = [];

  if (quickSearch.value) {
    chips.push({
      key: 'name',
      label: translations.filters.name.value,
      value: quickSearch.value,
    });
  }

  if (quickStatus.value !== null) {
    const statusLabel = quickStatus.value
      ? translations.filters.options.active.value
      : translations.filters.options.inactive.value;
    chips.push({
      key: 'status',
      label: translations.filters.status.value,
      value: statusLabel,
    });
  }

  if (advancedFilterValues.value.includeChildren !== null) {
    const childLabel = advancedFilterValues.value.includeChildren
      ? translations.filters.options.yes.value
      : translations.filters.options.no.value;
    chips.push({
      key: 'includeChildren',
      label: translations.filters.includeChildren.value,
      value: childLabel,
    });
  }

  if (advancedFilterValues.value.type !== null) {
    const typeLabels: Record<string, () => string> = {
      'aws-s3': () => translations.filters.options.awsS3.value,
      'azure-blob': () => translations.filters.options.azureBlob.value,
      'gcp-storage': () => translations.filters.options.gcpStorage.value,
      'minio': () => translations.filters.options.minio.value,
    };
    const typeLabel = typeLabels[advancedFilterValues.value.type]?.() || advancedFilterValues.value.type;
    chips.push({
      key: 'type',
      label: translations.filters.type.value,
      value: typeLabel,
    });
  }

  if (advancedFilterValues.value.region) {
    chips.push({
      key: 'region',
      label: translations.filters.region.value,
      value: advancedFilterValues.value.region,
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
 * Legacy computed for backward compatibility with visibleColumns
 */
const columnVisibility = computed(() => ({
  bucket: menuColumns.value.find(col => col.key === 'bucket')?.visible ?? true,
  region: menuColumns.value.find(col => col.key === 'region')?.visible ?? true,
  maxSize: menuColumns.value.find(col => col.key === 'maxSize')?.visible ?? true,
  frequency: menuColumns.value.find(col => col.key === 'frequency')?.visible ?? true,
}));

/**
 * Data lake columns with icon functions
 */
const lakeHouseColumns = computed(() => {
  return [
    {
      key: 'icon',
      label: '',
      type: 'avatar',
      visible: 'always',
      width: 56,
      icon: (value: any, row: any) => getLakeHouseIcon(row.type),
      color: (value: any, row: any) => getLakeHouseIconColor(row.type),
    },
    {
      key: 'name',
      label: translations.columns.value.find((col: DataRowColumn) => col.key === 'name')?.label || 'Name',
      type: 'text',
      visible: 'always',
      width: 250,
      ellipsis: true,
      secondaryKey: 'description',
    },
    {
      key: 'bucket',
      label: translations.columns.value.find((col: DataRowColumn) => col.key === 'bucket')?.label || 'Bucket',
      type: 'text',
      visible: 'laptop',
      width: 180,
      format: (value: any, row: any) => row.credentials?.bucket || translations.status.notAvailable.value,
      ellipsis: true,
    },
    {
      key: 'region',
      label: translations.columns.value.find((col: DataRowColumn) => col.key === 'region')?.label || 'Region',
      type: 'chip',
      visible: 'laptop',
      width: 140,
      format: (value: any, row: any) => row.credentials?.region || translations.status.notAvailable.value,
      color: 'blue-6',
    },
    {
      key: 'maxSize',
      label: translations.columns.value.find((col: DataRowColumn) => col.key === 'maxSize')?.label || 'Max Size',
      type: 'chip',
      visible: 'laptop',
      width: 120,
      format: (value: any, row: any) => `${row.pathConfig?.maxFileSize || 0} MB`,
      color: 'purple-6',
    },
    {
      key: 'frequency',
      label: translations.columns.value.find((col: DataRowColumn) => col.key === 'frequency')?.label || 'Frequency',
      type: 'text',
      visible: 'laptop',
      width: 180,
      format: (value: any, row: any) => {
        const freq = row.frequency;
        return freq ? `${freq.interval} ${freq.type} at ${freq.time}` : translations.status.notAvailable.value;
      },
    },
    {
      key: 'status',
      label: translations.columns.value.find((col: DataRowColumn) => col.key === 'status')?.label || 'Status',
      type: 'badge',
      visible: 'always',
      width: 100,
      format: (value: any) => value ? translations.status.active.value.toUpperCase() : translations.status.inactive.value.toUpperCase(),
      color: (value: any) => value ? 'green-6' : 'red-6',
    },
  ] as DataRowColumn[];
});

/**
 * Filtered columns based on visibility
 */
const visibleColumns = computed(() => {
  return lakeHouseColumns.value.filter((col: any) => {
    if (col.key === 'icon' || col.key === 'name' || col.key === 'status') {
      return true;
    }
    if (col.key === 'bucket') return columnVisibility.value.bucket;
    if (col.key === 'region') return columnVisibility.value.region;
    if (col.key === 'maxSize') return columnVisibility.value.maxSize;
    if (col.key === 'frequency') return columnVisibility.value.frequency;
    return true;
  });
});

/**
 * Filtered data lakes based on all filters
 */
const filteredLakeHouses = computed(() => {
  let result = [...channels.value];

  // Quick search filter
  if (quickSearch.value) {
    const search = quickSearch.value.toLowerCase();
    result = result.filter(dl =>
      dl.name?.toLowerCase().includes(search) ||
      dl.description?.toLowerCase().includes(search)
    );
  }

  // Quick status filter
  if (quickStatus.value !== null) {
    result = result.filter(dl => dl.status === quickStatus.value);
  }

  // Advanced: includeChildren filter
  // Note: This would be used by the API in real implementation
  // For frontend filtering, this is a placeholder

  // Advanced: type filter
  if (advancedFilterValues.value.type !== null) {
    result = result.filter(dl => dl.type === advancedFilterValues.value.type);
  }

  // Advanced: region filter
  if (advancedFilterValues.value.region) {
    const regionSearch = advancedFilterValues.value.region.toLowerCase();
    result = result.filter(dl =>
      dl.credentials?.region?.toLowerCase().includes(regionSearch)
    );
  }

  return result;
});

/** WATCHERS */

/**
 * Update filtered list when filteredLakeHouses changes
 */
watch(filteredLakeHouses, (newList) => {
  lakeHousesList.value = newList;
}, { immediate: true });

/** FUNCTIONS */

/**
 * Refresh lake house list (re-apply filters with current state)
 * @returns {void}
 */
function fetchLakeHouses(): void {
  loading.value = true;
  lakeHousesList.value = filteredLakeHouses.value;
  loading.value = false;
  lastUpdatedAt.value = Date.now();
}

/**
 * Apply quick filters immediately
 * @returns {void}
 */
function applyQuickFilters(): void {
  currentPage.value = 1;
  lakeHousesList.value = filteredLakeHouses.value;
}

/**
 * Handle advanced filters apply from drawer
 * @param {Record<string, any>} appliedFilters - Applied filter values
 * @returns {void}
 */
function handleAdvancedFiltersApply(appliedFilters: Record<string, any>): void {
  advancedFilterValues.value = { ...appliedFilters };
  hasPendingAdvancedFilters.value = false;
  currentPage.value = 1;
  showAdvancedFilters.value = false;
  lakeHousesList.value = filteredLakeHouses.value;
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
  if (key === 'name') {
    quickSearch.value = '';
  } else if (key === 'status') {
    quickStatus.value = null;
  } else if (key === 'includeChildren') {
    advancedFilterValues.value.includeChildren = null;
  } else if (key === 'type') {
    advancedFilterValues.value.type = null;
  } else if (key === 'region') {
    advancedFilterValues.value.region = '';
  }
  currentPage.value = 1;
  lakeHousesList.value = filteredLakeHouses.value;
}

/**
 * Clear all filters
 * @returns {void}
 */
function clearAllFilters(): void {
  quickSearch.value = '';
  quickStatus.value = null;
  advancedFilterValues.value = {
    includeChildren: null,
    type: null,
    region: '',
  };
  hasPendingAdvancedFilters.value = false;
  currentPage.value = 1;
  lakeHousesList.value = filteredLakeHouses.value;
}

/**
 * Update menu columns when changed
 * @param {ListHeaderMenuColumn[]} columns - Updated columns
 * @returns {void}
 */
function handleColumnsUpdate(columns: ListHeaderMenuColumn[]): void {
  menuColumns.value = columns;
}

/**
 * Handle card click
 * @param {any} item - Clicked item
 * @returns {void}
 */
function cardClick(item: any): void {
  logger.debug('Card clicked:', item);
}

/**
 * Edit data lake
 * @param {any} channel - Data lake to edit
 * @returns {void}
 */
function editNotification(channel: any): void {
  if (!canUpdateRetention.value) return;
  logger.debug('Edit channel:', channel);
}

/**
 * View data lake details
 * @param {any} channel - Data lake to view
 * @returns {void}
 */
function viewDetails(channel: any): void {
  if (!canReadRetention.value) return;
  logger.debug('View channel details:', channel);
}

/**
 * Confirm delete data lake — retention has no delete permission, always blocked
 */
function confirmDelete(): void {
  return;
}

/**
 * Handle page change
 * @param {number} page - New page number
 * @returns {void}
 */
function handlePageChange(page: number): void {
  currentPage.value = page;
}

/**
 * Returns the icon name for a given data lake type
 * @param {string} type - Data lake type
 * @returns {string}
 */
function getLakeHouseIcon(type: string): string {
  switch (type) {
    case 'aws-s3':
      return 'mdi-aws';
    case 'azure-blob':
      return 'mdi-microsoft-azure';
    case 'gcp-storage':
      return 'mdi-google-cloud';
    case 'minio':
      return 'mdi-database';
    default:
      return 'mdi-database';
  }
}

/**
 * Returns the icon color for a given data lake type
 * @param {string} type - Data lake type
 * @returns {string}
 */
function getLakeHouseIconColor(type: string): string {
  switch (type) {
    case 'aws-s3':
      return 'orange-6';
    case 'azure-blob':
      return 'blue-6';
    case 'gcp-storage':
      return 'red-6';
    case 'minio':
      return 'purple-6';
    default:
      return 'purple-6';
  }
}
</script>

<template>
  <q-page class="q-pt-lg">

    <!-- Header Section -->
    <PageHeader
        icon="cloud_upload"
        iconColor="primary"
        :title="translations.page.title.value"
        :description="translations.page.description.value"
        :button="canUpdateRetention ? { label: translations.page.button.add.value, icon: 'add', to: '/lakehouse/add', color: 'primary' } : undefined"
    />

    <!-- Filters Section -->
    <div class="text-caption text-grey-7 q-mb-xs">{{ translations.filters.label.value }}</div>
    <div class="row items-center q-col-gutter-sm q-mb-md">
      <!-- Search Input -->
      <div class="col">
        <q-input
          v-model="quickSearch"
          outlined
          dense
          clearable
          :placeholder="translations.filters.searchPlaceholder.value"
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
          :label="translations.filters.status.value"
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
          <AppTooltip :content="hasPendingAdvancedFilters
            ? translations.filters.pendingFilters.value
            : translations.filters.advancedFilters.value"
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
        :label="translations.filters.clearAll.value"
        no-caps
        class="q-ml-sm"
        @click="clearAllFilters"
      />
    </div>

    <!-- Results Section -->
    <div class="row items-center q-pt-xl q-mb-md">
      <div class="col">
        <div class="row items-center">
          <q-icon name="cloud_upload" size="sm" color="primary" class="q-mr-sm"/>
          <div class="text-subtitle1 text-weight-medium text-primary">{{ translations.page.listTitle.value }}</div>
        </div>
      </div>
      <div class="col-auto">
        <ListHeaderMenu
          :items-count="lakeHousesList.length"
          :item-label="translations.menuLabels.singular.value"
          :item-label-plural="translations.menuLabels.plural.value"
          icon="cloud_upload"
          :items-per-page="itemsPerPage"
          :columns="menuColumns"
          :refreshing="loading"
          :last-updated-at="lastUpdatedAt"
          @update:items-per-page="itemsPerPage = $event"
          @update:columns="handleColumnsUpdate"
          @refresh="fetchLakeHouses"
        />
      </div>
    </div>

    <!-- Data Lakes Row List -->
    <div class="row">
      <div
          v-for="datalake in lakeHousesList"
          :key="datalake.id"
          class="col-12 q-mb-xs"
      >
        <DataRow
            :data="datalake"
            :columns="visibleColumns"
            :actions="{ showEdit: canUpdateRetention, showView: canReadRetention, showDelete: false }"
            @click="cardClick"
            @dblclick="editNotification"
            @edit="editNotification"
            @view="viewDetails"
            @delete="confirmDelete"
        />
      </div>
    </div>

    <!-- No Results -->
    <div v-if="lakeHousesList.length === 0" class="row q-col-gutter-lg">
      <ListCardEmpty
          :title="translations.empty.title.value"
          :description="translations.empty.description.value"
          icon="cloud_upload"
      />
    </div>

    <!-- Pagination -->
    <ListPagination
      v-model="currentPage"
      :total-pages="totalPages"
      @change="handlePageChange"
    />

    <!-- Advanced Filters Drawer -->
    <AdvancedFiltersDrawer
      v-model="showAdvancedFilters"
      :title="translations.filters.advancedFilters.value"
      :fields="translations.advancedFilters.value"
      :values="advancedFilterValues"
      @apply="handleAdvancedFiltersApply"
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
