<script setup lang="ts">
defineOptions({
  name: 'HttpDataSourcesListPage'
});

/** TYPE IMPORTS */
import type { ListHeaderMenuColumn } from '@components/headers';
import type { DataSourceResponse } from '@mapexos/schemas';
import type {
  HttpDataSourcesListPageFilters,
  HttpDataSourcesListPageColumnVisibility,
  EnrichedDataSource,
} from './interfaces';

/** VUE IMPORTS */
import { ref, computed, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { copyToClipboard } from 'quasar';

/** COMPONENTS */
import { PageHeader, ListHeaderMenu } from '@components/headers';
import { ListCardEmpty, DataRow } from '@components/cards';
import { HttpDataSourceDetailsDrawer, AdvancedFiltersDrawer } from '@components/drawers';
import { ListPagination } from '@components/navigation';
import { AppTooltip } from '@components/tooltips';

/** COMPOSABLES */
import { useHttpDataSourcesTranslations } from '@composables/i18n';
import { usePermissions } from '@composables/shared/usePermissions';

/** UTILS */
import { notifySuccess } from '@utils/alert';

/** LOCAL IMPORTS */
import {
  HTTP_DATASOURCES_LIST_PAGE_DEFAULTS,
  HTTP_DATASOURCES_COLUMN_VISIBILITY_DEFAULTS,
  HTTP_DATASOURCES_FILTER_DEFAULTS,
} from './constants';
import {
  fetchDataSourcesHandler,
  handlePageChangeHandler,
  handleItemsPerPageChangeHandler,
  handleColumnsUpdateHandler,
  viewDetailsHandler,
  editDataSourceHandler,
  confirmDeleteHandler,
  deleteDataSourceHandler,
} from './handlers';

/** CONSTANTS */
const MAX_VISIBLE_CHIPS = 2;

/** COMPOSABLES & STORES */
const router = useRouter();
const t = useHttpDataSourcesTranslations();
const { canCreate, canUpdate, canDelete, canRead } = usePermissions();
const canCreateDS = canCreate('datasources');
const canUpdateDS = canUpdate('datasources');
const canDeleteDS = canDelete('datasources');
const canReadDS = canRead('datasources');

/** STATE */
const dataSourcesList = ref<EnrichedDataSource[]>([]);
const loading = ref(false);
const lastUpdatedAt = ref<number | undefined>(undefined);
const error = ref<string | undefined>(undefined);
const itemsPerPage = ref(HTTP_DATASOURCES_LIST_PAGE_DEFAULTS.ITEMS_PER_PAGE);
const currentPage = ref(HTTP_DATASOURCES_LIST_PAGE_DEFAULTS.INITIAL_PAGE);
const totalPages = ref(1);
const totalItems = ref(0);
const filters = ref<HttpDataSourcesListPageFilters>({ ...HTTP_DATASOURCES_FILTER_DEFAULTS });
const columnVisibilityState = ref<HttpDataSourcesListPageColumnVisibility>({ ...HTTP_DATASOURCES_COLUMN_VISIBILITY_DEFAULTS });
const detailsDrawerOpen = ref(false);
const selectedDataSourceId = ref<string | undefined>(undefined);

// Quick filters state
const quickSearch = ref('');
const quickStatus = ref<boolean | null>(null);

// Advanced filters state
const showAdvancedFilters = ref(false);
const advancedFilterValues = ref<Record<string, any>>({
  includeChildren: null,
  mode: undefined,
  auth: undefined,
  assetBind: undefined,
});
const hasPendingAdvancedFilters = ref(false);

/** COMPUTED */

/**
 * Status toggle options for quick filter
 */
const statusOptions = computed(() => [
  { label: t.filters.allStatus.value, value: null },
  { label: t.filters.options.active.value, value: true },
  { label: t.filters.options.inactive.value, value: false },
]);

/**
 * Active filter chips for display
 */
const activeFilterChips = computed(() => {
  const chips: Array<{ key: string; label: string; value: any }> = [];

  if (quickSearch.value) {
    chips.push({
      key: 'name',
      label: `${t.filters.name.value}: ${quickSearch.value}`,
      value: quickSearch.value,
    });
  }

  if (quickStatus.value !== null) {
    const statusLabel = quickStatus.value
      ? t.filters.options.active.value
      : t.filters.options.inactive.value;
    chips.push({
      key: 'enabled',
      label: `${t.filters.status.value}: ${statusLabel}`,
      value: quickStatus.value,
    });
  }

  if (advancedFilterValues.value.mode) {
    const modeLabel = advancedFilterValues.value.mode === 'pull'
      ? t.filters.options.pull.value
      : t.filters.options.push.value;
    chips.push({
      key: 'mode',
      label: `${t.filters.mode.value}: ${modeLabel}`,
      value: advancedFilterValues.value.mode,
    });
  }

  if (advancedFilterValues.value.includeChildren !== null) {
    const includeLabel = advancedFilterValues.value.includeChildren
      ? t.filters.options.yes.value
      : t.filters.options.no.value;
    chips.push({
      key: 'includeChildren',
      label: `${t.filters.includeChildren.value}: ${includeLabel}`,
      value: advancedFilterValues.value.includeChildren,
    });
  }

  if (advancedFilterValues.value.auth) {
    const authLabels: Record<string, () => string> = {
      'none': () => t.filters.options.none.value,
      'apiKey': () => t.filters.options.apiKey.value,
      'ip_whitelist': () => t.filters.options.ipWhitelist.value,
      'jwt': () => t.filters.options.jwt.value,
      'oauth2': () => t.filters.options.oauth2.value,
    };
    const authLabel = authLabels[advancedFilterValues.value.auth]?.() || advancedFilterValues.value.auth;
    chips.push({
      key: 'auth',
      label: `${t.filters.authType.value}: ${authLabel}`,
      value: advancedFilterValues.value.auth,
    });
  }

  if (advancedFilterValues.value.assetBind) {
    const bindLabel = advancedFilterValues.value.assetBind === 'fixedAssetId'
      ? t.filters.options.fixed.value
      : t.filters.options.dynamic.value;
    chips.push({
      key: 'assetBind',
      label: `${t.filters.assetBindType.value}: ${bindLabel}`,
      value: advancedFilterValues.value.assetBind,
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
 * Count of active advanced filters
 */
const advancedFiltersCount = computed(() => {
  let count = 0;
  if (advancedFilterValues.value.includeChildren !== null) count++;
  if (advancedFilterValues.value.mode) count++;
  if (advancedFilterValues.value.auth) count++;
  if (advancedFilterValues.value.assetBind) count++;
  return count;
});

/**
 * Column visibility using ListHeaderMenuColumn format with reactive translations
 */
const menuColumns = computed(() => {
  const cols: ListHeaderMenuColumn[] = [
    { key: 'assetBind', label: t.menuColumns.assetBind.value, visible: columnVisibilityState.value.assetBind },
    { key: 'auth', label: t.menuColumns.auth.value, visible: columnVisibilityState.value.auth },
    { key: 'mode', label: t.menuColumns.mode.value, visible: columnVisibilityState.value.mode },
  ];

  // Only show organization toggle when includeChildren is active
  if (filters.value.includeChildren === true) {
    cols.unshift({
      key: 'organization',
      label: t.menuColumns.organization.value,
      visible: columnVisibilityState.value.organization
    });
  }

  return cols;
});

/**
 * Custom actions configuration for DataRow
 */
const rowActions = computed(() => ({
  showEdit: canUpdateDS.value,
  showView: canReadDS.value,
  showDelete: canDeleteDS.value,
  customActions: [
    {
      key: 'copyEndpoint',
      label: t.actions.copyEndpoint.value,
      icon: 'content_copy',
      color: 'grey-7',
    },
  ],
}));

/**
 * Filtered columns based on visibility
 */
const visibleColumns = computed(() => {
  return t.columns.value.filter((col: any) => {
    if (col.key === 'icon' || col.key === 'name') {
      return true;
    }
    if (col.key === 'organizationName') {
      return filters.value.includeChildren === true && columnVisibilityState.value.organization;
    }
    if (col.key === 'assetBind.type') return columnVisibilityState.value.assetBind;
    if (col.key === 'auth.type') return columnVisibilityState.value.auth;
    if (col.key === 'mode') return columnVisibilityState.value.mode;
    return true;
  });
});

/** FUNCTIONS */

/**
 * Fetch data sources from API with current filters and pagination
 * @returns {Promise<void>}
 */
async function fetchDataSources(): Promise<void> {
  await fetchDataSourcesHandler(
    filters,
    currentPage,
    itemsPerPage,
    dataSourcesList,
    totalPages,
    totalItems,
    loading,
    error,
  );
  lastUpdatedAt.value = Date.now();
}

/**
 * Apply quick filters immediately
 * @returns {void}
 */
function applyQuickFilters(): void {
  filters.value.name = quickSearch.value || undefined;
  filters.value.enabled = quickStatus.value ?? undefined;
  currentPage.value = 1;
  void fetchDataSources();
}

/**
 * Handle advanced filters apply from drawer
 * @param {Record<string, any>} appliedFilters - Applied filter values
 * @returns {void}
 */
function handleAdvancedFiltersApply(appliedFilters: Record<string, any>): void {
  advancedFilterValues.value = { ...appliedFilters };
  filters.value.includeChildren = appliedFilters.includeChildren ?? undefined;
  filters.value.mode = appliedFilters.mode;
  filters.value.auth = appliedFilters.auth;
  filters.value.assetBind = appliedFilters.assetBind;
  hasPendingAdvancedFilters.value = false;
  currentPage.value = 1;
  showAdvancedFilters.value = false;
  void fetchDataSources();
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
  } else if (key === 'enabled') {
    quickStatus.value = null;
  } else if (key === 'mode') {
    advancedFilterValues.value.mode = undefined;
    filters.value.mode = undefined;
  } else if (key === 'includeChildren') {
    advancedFilterValues.value.includeChildren = null;
    filters.value.includeChildren = undefined;
  } else if (key === 'auth') {
    advancedFilterValues.value.auth = undefined;
    filters.value.auth = undefined;
  } else if (key === 'assetBind') {
    advancedFilterValues.value.assetBind = undefined;
    filters.value.assetBind = undefined;
  }
  currentPage.value = 1;
  void fetchDataSources();
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
    mode: undefined,
    auth: undefined,
    assetBind: undefined,
  };
  filters.value = { ...HTTP_DATASOURCES_FILTER_DEFAULTS };
  hasPendingAdvancedFilters.value = false;
  currentPage.value = 1;
  void fetchDataSources();
}

/**
 * Handle pagination navigation
 * @param {number} page - New page number
 * @returns {void}
 */
function handlePageChange(page: number): void {
  handlePageChangeHandler(page, currentPage, () => void fetchDataSources());
}

/**
 * Handle items per page change
 * @param {number} newValue - New items per page value
 * @returns {void}
 */
function handleItemsPerPageChange(newValue: number): void {
  handleItemsPerPageChangeHandler(newValue, itemsPerPage, currentPage, () => void fetchDataSources());
}

/**
 * Update menu columns when changed
 * @param {ListHeaderMenuColumn[]} columns - Updated column visibility states
 * @returns {void}
 */
function handleColumnsUpdate(columns: ListHeaderMenuColumn[]): void {
  handleColumnsUpdateHandler(columns, columnVisibilityState);
}

/**
 * View data source details
 * @param {DataSourceResponse} dataSource - Data source to view
 * @returns {void}
 */
function viewDetails(dataSource: DataSourceResponse): void {
  if (!canReadDS.value) return;
  viewDetailsHandler(dataSource, selectedDataSourceId, detailsDrawerOpen);
}

/**
 * Edit data source
 * @param {any} dataSource - Data source to edit
 * @returns {void}
 */
function editDataSource(dataSource: any): void {
  if (!canUpdateDS.value) return;
  editDataSourceHandler(dataSource, router);
}

/**
 * Confirm delete data source with dialog
 * @param {DataSourceResponse} dataSource - Data source to delete
 * @returns {Promise<void>}
 */
async function confirmDelete(dataSource: DataSourceResponse): Promise<void> {
  if (!canDeleteDS.value) return;
  await confirmDeleteHandler(dataSource, t, deleteDataSource);
}

/**
 * Delete data source from API and update list
 * @param {DataSourceResponse} dataSource - Data source to delete
 * @returns {Promise<void>}
 */
async function deleteDataSource(dataSource: DataSourceResponse): Promise<void> {
  await deleteDataSourceHandler(dataSource, dataSourcesList, totalItems, currentPage, t, fetchDataSources);
}

/**
 * Handle edit action from drawer
 * @param {string} dataSourceId - Data source ID to edit
 * @returns {void}
 */
function handleEditFromDrawer(dataSourceId: string): void {
  if (!canUpdateDS.value) return;
  void router.push(`/data_sources/http/${dataSourceId}/edit`);
}

/**
 * Handle custom action from DataRow menu
 * @param {string} actionKey - Action key
 * @param {DataSourceResponse} dataSource - Data source
 * @returns {void}
 */
function handleCustomAction(actionKey: string, dataSource: DataSourceResponse): void {
  if (actionKey === 'copyEndpoint') {
    void copyEndpointUrl(dataSource);
  }
}

/**
 * Copy endpoint URL to clipboard
 * @param {DataSourceResponse} dataSource - Data source to copy URL for
 * @returns {Promise<void>}
 */
async function copyEndpointUrl(dataSource: DataSourceResponse): Promise<void> {
  const endpointUrl = `{{your_endpoint}}/api/v1/events?ds=${dataSource.id}`;

  try {
    await copyToClipboard(endpointUrl);
    notifySuccess({ message: t.notifications.endpointCopied.value });
  } catch {
    const textArea = document.createElement('textarea');
    textArea.value = endpointUrl;
    document.body.appendChild(textArea);
    textArea.select();
    document.execCommand('copy');
    document.body.removeChild(textArea);
    notifySuccess({ message: t.notifications.endpointCopied.value });
  }
}

/** LIFECYCLE HOOKS */
onMounted(async () => await fetchDataSources());
</script>

<template>
  <q-page class="q-pt-lg">

    <!-- Header Section -->
    <PageHeader
        icon="settings_input_antenna"
        iconColor="primary"
        :title="t.pageHeader.title.value"
        :description="t.pageHeader.description.value"
        :button="canCreateDS ? { label: t.pageHeader.button.value, icon: 'add', to: '/data_sources/http/add', color: 'primary' } : undefined"
        :info="t.pageHeader.info.value"
    />

    <!-- Quick Filters Section -->
    <div class="text-caption text-grey-7 q-mb-xs">{{ t.filters.label.value }}</div>
    <div class="row items-center q-col-gutter-sm q-mb-md">
      <!-- Search Input -->
      <div class="col">
        <q-input
          v-model="quickSearch"
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

      <!-- Advanced Filters Button -->
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
        <span class="text-weight-medium">{{ chip.label }}</span>
      </q-chip>

      <!-- Hidden Chips Badge -->
      <q-badge
        v-if="hiddenFiltersCount > 0"
        color="primary"
        outline
        class="q-pa-xs cursor-pointer"
      >
        +{{ hiddenFiltersCount }}
        <AppTooltip>
          <div v-for="chip in hiddenFilterChips" :key="chip.key" class="q-mb-xs">
            {{ chip.label }}
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
          <q-icon name="settings_input_antenna" size="sm" color="primary" class="q-mr-sm"/>
          <div class="text-subtitle1 text-weight-medium text-primary">{{ t.list.title.value }}</div>
        </div>
      </div>
      <div class="col-auto">
        <ListHeaderMenu
          icon="settings_input_antenna"
          :items-count="totalItems"
          :item-label="t.headerMenu.itemLabel.value"
          :item-label-plural="t.headerMenu.itemLabelPlural.value"
          :items-per-page="itemsPerPage"
          :columns="menuColumns"
          :refreshing="loading"
          :last-updated-at="lastUpdatedAt"
          @update:items-per-page="handleItemsPerPageChange"
          @update:columns="handleColumnsUpdate"
          @refresh="fetchDataSources"
        />
      </div>
    </div>

    <!-- Loading Spinner -->
    <div v-if="loading" class="row justify-center q-my-lg">
      <q-spinner color="primary" size="3em" />
    </div>

    <!-- Data Sources Row List -->
    <div v-else class="row">
      <div
          v-for="(dataSource, index) in dataSourcesList"
          :key="dataSource.id || `datasource-${index}`"
          class="col-12 q-mb-xs"
      >
        <DataRow
            :data="dataSource"
            :columns="visibleColumns"
            :actions="rowActions"
            @click="viewDetails"
            @dblclick="editDataSource"
            @edit="editDataSource"
            @view="viewDetails"
            @delete="confirmDelete"
            @action="handleCustomAction"
        />
      </div>

      <!-- No Results -->
      <div v-if="dataSourcesList.length === 0" class="col-12">
        <ListCardEmpty
            :title="t.empty.title.value"
            :description="t.empty.description.value"
            icon="settings_input_antenna"
        />
      </div>
    </div>

    <!-- Pagination -->
    <ListPagination
      v-model="currentPage"
      :total-pages="totalPages"
      @change="handlePageChange"
    />

    <!-- Details Drawer -->
    <HttpDataSourceDetailsDrawer
        v-if="selectedDataSourceId"
        v-model="detailsDrawerOpen"
        :data-source-id="selectedDataSourceId"
        @edit="handleEditFromDrawer"
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

<style lang="scss" scoped>
.filter-input {
  :deep(.q-field__control) {
    border-radius: var(--mapex-radius-md);
  }
}
</style>
