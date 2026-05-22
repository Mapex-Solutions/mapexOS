<script setup lang="ts">
defineOptions({
  name: 'AssetsListPage'
});

/** TYPE IMPORTS */
import type { ListHeaderMenuColumn } from '@components/headers';
import type { FilterField } from '@components/drawers';
import type {
  AssetsListPageFilters,
  AssetsListPageColumnVisibility,
  DynamicFilterOptions,
  EnrichedAsset,
} from './interfaces';

/** VUE IMPORTS */
import { ref, computed, onMounted } from 'vue';
import { useRouter } from 'vue-router';

/** COMPONENTS */
import { PageHeader, ListHeaderMenu } from '@components/headers';
import { ListCardEmpty, DataRow } from '@components/cards';
import { AssetDetailsDrawer, AdvancedFiltersDrawer } from '@components/drawers';
import { ListPagination } from '@components/navigation';
import { AppTooltip } from '@components/tooltips';

/** COMPOSABLES */
import { useAssetsTranslations } from '@composables/i18n';
import { useOrgChangeRefresh } from '@composables/organizations';
import { usePermissions } from '@composables/shared/usePermissions';
import { useLogger } from '@composables/useLogger';

/** UTILS */
import { notifyFail, notifySuccess, dialogDelete } from '@utils/alert';
import { cleanQueryParams } from '@utils/query';

/** SERVICES */
import { apis } from '@services/mapex';

/** STORES */
import { useOrganizationStore } from '@stores/organization';

/** LOCAL IMPORTS */
import {
  DEFAULT_ITEMS_PER_PAGE,
  ASSETS_COLUMN_VISIBILITY_DEFAULTS,
  ASSETS_FILTER_DEFAULTS,
  ASSETS_PROJECTION,
  LIST_TYPE,
  CASCADING_FILTER_DEFAULTS,
} from './constants';

/** COMPOSABLES & STORES */
const t = useAssetsTranslations();
const orgStore = useOrganizationStore();
const router = useRouter();
const logger = useLogger('AssetsListPage');
const { canCreate, canUpdate, canDelete, canRead } = usePermissions();
const canCreateAsset = canCreate('assets');
const canUpdateAsset = canUpdate('assets');
const canDeleteAsset = canDelete('assets');
const canReadAsset = canRead('assets');

/** STATE */
const assetsList = ref<EnrichedAsset[]>([]);
const loading = ref(false);
const lastUpdatedAt = ref<number | undefined>(undefined);
const error = ref<string | undefined>(undefined);
const showDetailsDrawer = ref(false);
const showFiltersDrawer = ref(false);
const selectedAssetId = ref<string | null>(null);
const itemsPerPage = ref(DEFAULT_ITEMS_PER_PAGE);
const currentPage = ref(1);
const totalPages = ref(1);
const totalItems = ref(0);
const hasNext = ref(false);
const hasPrev = ref(false);
const filters = ref<AssetsListPageFilters>({ ...ASSETS_FILTER_DEFAULTS });
const columnVisibilityState = ref<AssetsListPageColumnVisibility>({ ...ASSETS_COLUMN_VISIBILITY_DEFAULTS });

// Quick filter state (inline)
const quickSearch = ref('');
const quickStatus = ref<boolean | null>(null);

// Advanced filter state (drawer)
const advancedFilterValues = ref<Record<string, any>>({
  includeChildren: null,
  assetUUID: null,
  categoryId: null,
  manufacturerId: null,
  modelId: null,
});

// Cascading filter options
const categoryOptions = ref<DynamicFilterOptions[]>([]);
const manufacturerOptions = ref<DynamicFilterOptions[]>([]);
const modelOptions = ref<DynamicFilterOptions[]>([]);
const loadingCategories = ref(false);
const loadingManufacturers = ref(false);
const loadingModels = ref(false);

// Pending changes state
const hasPendingAdvancedFilters = ref(false);

/** COMPUTED */

/**
 * Status options for quick filter select
 */
const statusOptions = computed(() => [
  { label: t.filters.allStatus.value, value: null },
  { label: t.filters.options.enabled.value, value: true },
  { label: t.filters.options.disabled.value, value: false },
]);

/**
 * Advanced filter fields for the drawer
 */
const advancedFilterFields = computed((): FilterField[] => [
  {
    key: 'includeChildren',
    label: t.filters.includeChildren.value,
    type: 'toggle',
    icon: 'account_tree',
    options: [
      { label: t.filters.allStatus.value, value: null },
      { label: t.filters.options.yes.value, value: true },
      { label: t.filters.options.no.value, value: false },
    ],
  },
  {
    key: 'assetUUID',
    label: t.filters.assetUUID.value,
    type: 'input',
    icon: 'fingerprint',
    placeholder: t.filters.filterByUUID.value,
  },
  {
    key: 'categoryId',
    label: t.filters.category.value,
    type: 'select',
    icon: 'category',
    options: categoryOptions.value,
    loading: loadingCategories.value,
    placeholder: t.filters.filterByCategory.value,
  },
  {
    key: 'manufacturerId',
    label: t.filters.manufacturer.value,
    type: 'select',
    icon: 'factory',
    options: manufacturerOptions.value,
    loading: loadingManufacturers.value,
    disabled: !advancedFilterValues.value.categoryId,
    placeholder: t.filters.filterByManufacturer.value,
  },
  {
    key: 'modelId',
    label: t.filters.model.value,
    type: 'select',
    icon: 'memory',
    options: modelOptions.value,
    loading: loadingModels.value,
    disabled: !advancedFilterValues.value.manufacturerId,
    placeholder: t.filters.filterByModel.value,
  },
]);

/**
 * Check if any filters are active (quick or advanced)
 */
const hasActiveFilters = computed(() => {
  return !!(
    quickSearch.value ||
    quickStatus.value !== null ||
    advancedFilterValues.value.includeChildren !== null ||
    advancedFilterValues.value.assetUUID ||
    advancedFilterValues.value.categoryId ||
    advancedFilterValues.value.manufacturerId ||
    advancedFilterValues.value.modelId
  );
});

/**
 * Count of active advanced filters (for badge)
 */
const advancedFiltersCount = computed(() => {
  let count = 0;
  if (advancedFilterValues.value.includeChildren !== null) count++;
  if (advancedFilterValues.value.assetUUID) count++;
  if (advancedFilterValues.value.categoryId) count++;
  if (advancedFilterValues.value.manufacturerId) count++;
  if (advancedFilterValues.value.modelId) count++;
  return count;
});

/**
 * Active filter chips for visual feedback
 */
const activeFilterChips = computed(() => {
  const chips: Array<{ key: string; label: string; value: string }> = [];

  if (quickSearch.value) {
    chips.push({
      key: 'name',
      label: t.filters.assetName.value,
      value: quickSearch.value,
    });
  }

  if (quickStatus.value !== null) {
    chips.push({
      key: 'status',
      label: t.filters.status.value,
      value: quickStatus.value ? t.filters.options.enabled.value : t.filters.options.disabled.value,
    });
  }

  if (advancedFilterValues.value.includeChildren !== null) {
    chips.push({
      key: 'includeChildren',
      label: t.filters.includeChildren.value,
      value: advancedFilterValues.value.includeChildren ? t.filters.options.yes.value : t.filters.options.no.value,
    });
  }

  if (advancedFilterValues.value.assetUUID) {
    chips.push({
      key: 'assetUUID',
      label: t.filters.assetUUID.value,
      value: advancedFilterValues.value.assetUUID,
    });
  }

  if (advancedFilterValues.value.categoryId) {
    const category = categoryOptions.value.find(c => c.value === advancedFilterValues.value.categoryId);
    chips.push({
      key: 'categoryId',
      label: t.filters.category.value,
      value: category?.label || advancedFilterValues.value.categoryId,
    });
  }

  if (advancedFilterValues.value.manufacturerId) {
    const manufacturer = manufacturerOptions.value.find(m => m.value === advancedFilterValues.value.manufacturerId);
    chips.push({
      key: 'manufacturerId',
      label: t.filters.manufacturer.value,
      value: manufacturer?.label || advancedFilterValues.value.manufacturerId,
    });
  }

  if (advancedFilterValues.value.modelId) {
    const model = modelOptions.value.find(m => m.value === advancedFilterValues.value.modelId);
    chips.push({
      key: 'modelId',
      label: t.filters.model.value,
      value: model?.label || advancedFilterValues.value.modelId,
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

/**
 * Column visibility using ListHeaderMenuColumn format with reactive translations
 */
const menuColumns = computed(() => {
  const cols: ListHeaderMenuColumn[] = [
    { key: 'uuid', label: t.menuColumns.uuid.value, visible: columnVisibilityState.value.uuid },
    { key: 'manufacturerModel', label: t.menuColumns.manufacturerModel.value, visible: columnVisibilityState.value.manufacturerModel },
    { key: 'protocol', label: t.menuColumns.protocol.value, visible: columnVisibilityState.value.protocol },
    { key: 'debugEnabled', label: t.menuColumns.debug.value, visible: columnVisibilityState.value.debugEnabled },
    { key: 'healthStatus', label: t.menuColumns.status.value, visible: columnVisibilityState.value.healthStatus ?? true },
  ];

  // Only show organization toggle when includeChildren is active
  if (advancedFilterValues.value.includeChildren === true) {
    cols.unshift({ key: 'organization', label: t.menuColumns.organization.value, visible: columnVisibilityState.value.organization });
  }

  return cols;
});

/**
 * Filtered columns based on visibility and includeChildren filter
 */
const visibleColumns = computed(() => {
  return t.columns.value.filter(col => {
    // Always show avatar, name, and status
    if (col.key === 'icon' || col.key === 'name' || col.key === 'status') {
      return true;
    }

    // Filter based on columnVisibility
    if (col.key === 'assetUUID') return columnVisibilityState.value.uuid;
    if (col.key === 'manufacturerName') return columnVisibilityState.value.manufacturerModel;
    if (col.key === 'protocol.type') return columnVisibilityState.value.protocol;
    if (col.key === 'debugEnabled') return columnVisibilityState.value.debugEnabled;
    if (col.key === 'healthStatus') return columnVisibilityState.value.healthStatus ?? true;

    // Only show organization column when includeChildren is active
    if (col.key === 'organizationName') {
      return advancedFilterValues.value.includeChildren === true && columnVisibilityState.value.organization;
    }

    return true;
  });
});

/** FUNCTIONS */

/**
 * Apply quick filters and fetch data
 * @returns {void}
 */
function applyQuickFilters(): void {
  filters.value.name = quickSearch.value || undefined;
  filters.value.status = quickStatus.value ?? undefined;
  currentPage.value = 1;
  void fetchAssets();
}

/**
 * Handle advanced filters apply from drawer
 * @param {Record<string, any>} appliedFilters - Applied filter values
 * @returns {void}
 */
function handleAdvancedFiltersApply(appliedFilters: Record<string, any>): void {
  logger.debug('Advanced filters applied:', appliedFilters);

  advancedFilterValues.value = { ...appliedFilters };

  // Map advanced filter values to API filter format
  filters.value.includeChildren = appliedFilters.includeChildren;
  filters.value.assetUUID = appliedFilters.assetUUID || undefined;
  filters.value.categoryId = appliedFilters.categoryId || undefined;
  filters.value.manufacturerId = appliedFilters.manufacturerId || undefined;
  filters.value.modelId = appliedFilters.modelId || undefined;

  // Reset pending state after apply
  hasPendingAdvancedFilters.value = false;

  currentPage.value = 1;
  showFiltersDrawer.value = false;
  void fetchAssets();
}

/**
 * Handle pending state change from advanced filters drawer
 * @param {boolean} hasPending - Whether there are pending changes
 * @returns {void}
 */
function handlePendingChange(hasPending: boolean): void {
  hasPendingAdvancedFilters.value = hasPending;
}

/**
 * Handle advanced filters reset from drawer
 * @returns {void}
 */
function handleAdvancedFiltersReset(): void {
  advancedFilterValues.value = {
    includeChildren: null,
    assetUUID: null,
    categoryId: null,
    manufacturerId: null,
    modelId: null,
  };

  // Reset cascading filter options
  manufacturerOptions.value = [];
  modelOptions.value = [];

  filters.value.includeChildren = undefined;
  filters.value.assetUUID = undefined;
  filters.value.categoryId = undefined;
  filters.value.manufacturerId = undefined;
  filters.value.modelId = undefined;

  currentPage.value = 1;
  void fetchAssets();
}

/**
 * Handle field change in advanced filters (for cascading logic)
 * Updates the field value and handles cascading dependencies
 * @param {string} key - Field key that changed
 * @param {any} value - New value
 * @returns {void}
 */
function handleAdvancedFieldChange(key: string, value: any): void {
  logger.debug(`Advanced field change: ${key} = ${value}`);

  // Update the field value in advancedFilterValues (required for cascading to work)
  advancedFilterValues.value[key] = value;

  if (key === 'categoryId') {
    // Reset dependent fields
    advancedFilterValues.value.manufacturerId = null;
    advancedFilterValues.value.modelId = null;
    manufacturerOptions.value = [];
    modelOptions.value = [];

    if (value) {
      void fetchManufacturers(value);
    }
  } else if (key === 'manufacturerId') {
    // Reset dependent field
    advancedFilterValues.value.modelId = null;
    modelOptions.value = [];

    if (value) {
      void fetchModels(value);
    }
  }
}

/**
 * Remove individual filter
 * @param {string} key - Filter key to remove
 * @returns {void}
 */
function removeFilter(key: string): void {
  if (key === 'name') {
    quickSearch.value = '';
    filters.value.name = undefined;
  } else if (key === 'status') {
    quickStatus.value = null;
    filters.value.status = undefined;
  } else if (key === 'includeChildren') {
    advancedFilterValues.value.includeChildren = null;
    filters.value.includeChildren = undefined;
  } else if (key === 'assetUUID') {
    advancedFilterValues.value.assetUUID = null;
    filters.value.assetUUID = undefined;
  } else if (key === 'categoryId') {
    advancedFilterValues.value.categoryId = null;
    advancedFilterValues.value.manufacturerId = null;
    advancedFilterValues.value.modelId = null;
    manufacturerOptions.value = [];
    modelOptions.value = [];
    filters.value.categoryId = undefined;
    filters.value.manufacturerId = undefined;
    filters.value.modelId = undefined;
  } else if (key === 'manufacturerId') {
    advancedFilterValues.value.manufacturerId = null;
    advancedFilterValues.value.modelId = null;
    modelOptions.value = [];
    filters.value.manufacturerId = undefined;
    filters.value.modelId = undefined;
  } else if (key === 'modelId') {
    advancedFilterValues.value.modelId = null;
    filters.value.modelId = undefined;
  }

  currentPage.value = 1;
  void fetchAssets();
}

/**
 * Clear all filters (quick and advanced)
 * @returns {void}
 */
function clearAllFilters(): void {
  quickSearch.value = '';
  quickStatus.value = null;
  advancedFilterValues.value = {
    includeChildren: null,
    assetUUID: null,
    categoryId: null,
    manufacturerId: null,
    modelId: null,
  };
  manufacturerOptions.value = [];
  modelOptions.value = [];
  filters.value = { ...ASSETS_FILTER_DEFAULTS };

  // Reset pending state
  hasPendingAdvancedFilters.value = false;

  currentPage.value = 1;
  void fetchAssets();
}

/**
 * Fetch assets from API with current filters and pagination
 * @returns {Promise<void>}
 */
async function fetchAssets(): Promise<void> {
  if (!apis.assets) {
    error.value = 'Assets API not initialized';
    return;
  }

  try {
    loading.value = true;
    error.value = undefined;

    // Build query parameters from filters and pagination state
    const queryParams: Record<string, any> = {
      page: currentPage.value,
      perPage: itemsPerPage.value,
      projection: ASSETS_PROJECTION,
    };

    // Add active filters to query params (only if they have values)
    if (filters.value.name) {
      queryParams.name = filters.value.name;
    }
    if (filters.value.assetUUID) {
      queryParams.assetUUID = filters.value.assetUUID;
    }
    if (typeof filters.value.status === 'boolean') {
      queryParams.enabled = filters.value.status;
    }
    if (filters.value.categoryId) {
      queryParams.categoryId = filters.value.categoryId;
    }
    if (filters.value.manufacturerId) {
      queryParams.manufacturerId = filters.value.manufacturerId;
    }
    if (filters.value.modelId) {
      queryParams.modelId = filters.value.modelId;
    }
    if (typeof filters.value.includeChildren === 'boolean') {
      queryParams.includeChildren = filters.value.includeChildren;
    }

    // Clean undefined values to avoid sending "undefined" as string in URL
    const cleanedParams = cleanQueryParams(queryParams);

    const response = await apis.assets.asset.list(cleanedParams);

    // Enrich assets with organization name when includeChildren is active
    const enrichedAssets = response.items.map((asset: any) => {
      // Find organization name from store's flatList using orgId
      const organization = orgStore.flatList.find(org => org.id === asset.orgId);
      return {
        ...asset,
        organizationName: organization?.name || 'Unknown',
      };
    });

    assetsList.value = enrichedAssets;

    // Update pagination state from response
    if (response.pagination) {
      totalPages.value = response.pagination.totalPages || 1;
      totalItems.value = response.pagination.totalItems || 0;
      // Calculate hasNext and hasPrev based on current page and total pages
      hasNext.value = currentPage.value < totalPages.value;
      hasPrev.value = currentPage.value > 1;
    }
  } catch (err: any) {
    logger.error('Error fetching assets:', err);
    const errorMsg = err.message || 'Failed to fetch assets';
    error.value = errorMsg;
    notifyFail({ message: errorMsg });
  } finally {
    loading.value = false;
    lastUpdatedAt.value = Date.now();
  }
}

/**
 * Handle pagination navigation
 * @param {number} page - New page number
 * @returns {void}
 */
function handlePageChange(page: number): void {
  currentPage.value = page;
  void fetchAssets();
}

/**
 * Handle items per page change
 * @param {number} newValue - New items per page value
 * @returns {void}
 */
function handleItemsPerPageChange(newValue: number): void {
  itemsPerPage.value = newValue;
  currentPage.value = 1; // Reset to first page
  void fetchAssets();
}

/**
 * Fetch categories from API
 * @returns {Promise<void>}
 */
async function fetchCategories(): Promise<void> {
  if (!apis.mapexOS?.lists) {
    return;
  }

  try {
    loadingCategories.value = true;
    const response = await apis.mapexOS.lists.list({
      type: LIST_TYPE.ASSET_CATEGORY,
      page: CASCADING_FILTER_DEFAULTS.page,
      perPage: CASCADING_FILTER_DEFAULTS.perPage,
    });

    categoryOptions.value = response.items.map((item: any) => ({
      label: item.name,
      value: item.id,
    }));
  } catch (err: any) {
    logger.error('Error fetching categories:', err);
  } finally {
    loadingCategories.value = false;
  }
}

/**
 * Fetch manufacturers based on selected category
 * @param {string} categoryId - Category ID to fetch manufacturers for
 * @returns {Promise<void>}
 */
async function fetchManufacturers(categoryId: string): Promise<void> {
  if (!apis.mapexOS?.lists || !categoryId) {
    manufacturerOptions.value = [];
    return;
  }

  try {
    loadingManufacturers.value = true;
    const response = await apis.mapexOS.lists.list({
      type: LIST_TYPE.ASSET_MANUFACTURER,
      parentId: categoryId,
      page: CASCADING_FILTER_DEFAULTS.page,
      perPage: CASCADING_FILTER_DEFAULTS.perPage,
    });

    manufacturerOptions.value = response.items.map((item: any) => ({
      label: item.name,
      value: item.id,
    }));
  } catch (err: any) {
    logger.error('Error fetching manufacturers:', err);
    manufacturerOptions.value = [];
  } finally {
    loadingManufacturers.value = false;
  }
}

/**
 * Fetch models based on selected manufacturer
 * @param {string} manufacturerId - Manufacturer ID to fetch models for
 * @returns {Promise<void>}
 */
async function fetchModels(manufacturerId: string): Promise<void> {
  if (!apis.mapexOS?.lists || !manufacturerId) {
    modelOptions.value = [];
    return;
  }

  try {
    loadingModels.value = true;
    const response = await apis.mapexOS.lists.list({
      type: LIST_TYPE.ASSET_MODEL,
      parentId: manufacturerId,
      page: CASCADING_FILTER_DEFAULTS.page,
      perPage: CASCADING_FILTER_DEFAULTS.perPage,
    });

    modelOptions.value = response.items.map((item: any) => ({
      label: item.name,
      value: item.id,
    }));
  } catch (err: any) {
    logger.error('Error fetching models:', err);
    modelOptions.value = [];
  } finally {
    loadingModels.value = false;
  }
}

/**
 * Update menu columns when changed
 * @param {ListHeaderMenuColumn[]} columns - Updated column visibility states
 * @returns {void}
 */
function handleColumnsUpdate(columns: ListHeaderMenuColumn[]): void {
  columns.forEach(col => {
    if (col.key === 'uuid') columnVisibilityState.value.uuid = col.visible;
    if (col.key === 'manufacturerModel') columnVisibilityState.value.manufacturerModel = col.visible;
    if (col.key === 'protocol') columnVisibilityState.value.protocol = col.visible;
    if (col.key === 'debugEnabled') columnVisibilityState.value.debugEnabled = col.visible;
    if (col.key === 'healthStatus') columnVisibilityState.value.healthStatus = col.visible;
    if (col.key === 'organization') columnVisibilityState.value.organization = col.visible;
  });
}

/**
 * View asset details in drawer
 * IMPORTANT: Pass only the asset ID, not the full asset object
 * The drawer will fetch complete asset data using the API
 * @param {EnrichedAsset} asset - Asset to view
 * @returns {void}
 */
function viewDetails(asset: EnrichedAsset): void {
  if (!canReadAsset.value) return;
  if (!asset.id) {
    notifyFail({ message: t.errors.idMissing.value });
    return;
  }
  selectedAssetId.value = asset.id;
  showDetailsDrawer.value = true;
}

/**
 * Edit asset (navigate to edit page)
 * @param {EnrichedAsset} asset - Asset to edit
 * @returns {void}
 */
function editRule(asset: EnrichedAsset): void {
  if (!canUpdateAsset.value) return;
  void router.push(`/assets/edit/${asset.id}`);
}

/**
 * Handle edit event from drawer
 * @param {string} assetId - ID of asset to edit
 * @returns {void}
 */
function handleDrawerEdit(assetId: string): void {
  if (!canUpdateAsset.value) return;
  void router.push(`/assets/edit/${assetId}`);
}

/**
 * Handle duplicate event from drawer
 * @param {string} assetId - ID of asset to duplicate
 * @returns {void}
 */
function handleDrawerDuplicate(assetId: string): void {
  logger.debug('Duplicate asset from drawer:', assetId);
  // TODO: Navigate to duplicate page or implement duplicate functionality
  // router.push(`/assets/duplicate/${assetId}`);
}

/**
 * Confirm delete operation with dialog
 * @param {EnrichedAsset} asset - Asset to delete
 * @returns {Promise<void>}
 */
async function confirmDelete(asset: EnrichedAsset): Promise<void> {
  if (!canDeleteAsset.value) return;
  const assetName = asset.name || 'this asset';
  const confirmed = await dialogDelete({
    title: t.dialog.deleteTitle.value,
    message: t.messages.confirmDelete(assetName),
  });

  if (confirmed) {
    await deleteRule(asset);
  }
}

/**
 * Delete asset from API and update list
 * @param {EnrichedAsset} asset - Asset to delete
 * @returns {Promise<void>}
 */
async function deleteRule(asset: EnrichedAsset): Promise<void> {
  if (!orgStore.selectedOrganizationId) {
    notifyFail({ message: t.errors.noOrganization.value });
    return;
  }

  if (!apis.assets) {
    notifyFail({ message: t.errors.apiNotInitialized.value });
    return;
  }

  if (!asset.id) {
    notifyFail({ message: t.errors.idMissing.value });
    return;
  }

  try {
    // Call delete API with correct assetId field from asset.id
    await apis.assets.asset.delete({ assetId: asset.id });

    // Update total items count after successful deletion
    totalItems.value = Math.max(0, totalItems.value - 1);

    // Remove from local list
    const index = assetsList.value.findIndex(r => r.id === asset.id);
    if (index !== -1) {
      assetsList.value.splice(index, 1);
    }

    // If current page is now empty and not the first page, go to previous page
    if (assetsList.value.length === 0 && currentPage.value > 1) {
      currentPage.value -= 1;
      await fetchAssets();
    }

    notifySuccess({ message: t.messages.deletedSuccessfully.value });
  } catch (err: any) {
    notifyFail({ message: err.message || 'Failed to delete asset' });
  }
}

/** LIFECYCLE HOOKS */

onMounted(async () => {
  await fetchCategories();
  await fetchAssets();
});

/**
 * Auto-refresh assets when organization changes
 * This ensures data stays consistent with selected organization context
 * Only triggers on actual org changes, not on initial mount (to avoid double-fetching)
 */
useOrgChangeRefresh(async () => {
  // Reset pagination and filters when org changes
  currentPage.value = 1;
  quickSearch.value = '';
  quickStatus.value = null;
  advancedFilterValues.value = {
    includeChildren: null,
    assetUUID: null,
    categoryId: null,
    manufacturerId: null,
    modelId: null,
  };
  manufacturerOptions.value = [];
  modelOptions.value = [];
  filters.value = { ...ASSETS_FILTER_DEFAULTS };
  await fetchAssets();
});

</script>

<template>
  <q-page class="q-pt-lg">

    <!-- Header Section -->
    <PageHeader
        icon="devices"
        iconColor="primary"
        :title="t.page.title.value"
        :description="t.page.description.value"
        :button="canCreateAsset ? { label: t.page.addButton.value, icon: 'add', to: '/assets/add', color: 'primary' } : undefined"
        :info="t.page.info.value"
    />

    <!-- Filters Section -->
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
          <q-icon name="devices" size="sm" color="primary" class="q-mr-sm"/>
          <div class="text-subtitle1 text-weight-medium text-primary">{{ t.page.listTitle.value }}</div>
        </div>
      </div>
      <div class="col-auto">
        <ListHeaderMenu
          icon="devices"
          :items-count="totalItems"
          :item-label="t.page.itemLabel.value"
          :item-label-plural="t.page.itemLabelPlural.value"
          :items-per-page="itemsPerPage"
          :columns="menuColumns"
          :filtered="hasActiveFilters"
          :refreshing="loading"
          :last-updated-at="lastUpdatedAt"
          @update:items-per-page="handleItemsPerPageChange"
          @update:columns="handleColumnsUpdate"
          @refresh="fetchAssets"
        />
      </div>
    </div>

    <!-- Loading Spinner -->
    <div v-if="loading" class="row justify-center q-my-lg">
      <q-spinner color="primary" size="3em" />
    </div>

    <!-- Assets Row List -->
    <div v-else class="row">
      <div
          v-for="(asset, index) in assetsList"
          :key="asset.id || `asset-${index}`"
          class="col-12 q-mb-xs"
      >
        <DataRow
            :data="asset"
            :columns="visibleColumns"
            :actions="{ showEdit: canUpdateAsset, showView: canReadAsset, showDelete: canDeleteAsset }"
            @click="viewDetails"
            @dblclick="editRule"
            @edit="editRule"
            @view="viewDetails"
            @delete="confirmDelete"
        />
      </div>

      <!-- No Results -->
      <div v-if="assetsList.length === 0" class="col-12">
        <ListCardEmpty
            :title="t.empty.title.value"
            :description="t.empty.description.value"
            icon="devices"
        />
      </div>
    </div>

    <!-- Pagination -->
    <ListPagination
      v-model="currentPage"
      :total-pages="totalPages"
      @change="handlePageChange"
    />

    <!-- Asset Details Drawer -->
    <AssetDetailsDrawer
      v-model="showDetailsDrawer"
      :asset-id="selectedAssetId"
      @edit="handleDrawerEdit"
      @duplicate="handleDrawerDuplicate"
    />

    <!-- Advanced Filters Drawer -->
    <AdvancedFiltersDrawer
      v-model="showFiltersDrawer"
      :fields="advancedFilterFields"
      :values="advancedFilterValues"
      @apply="handleAdvancedFiltersApply"
      @reset="handleAdvancedFiltersReset"
      @field-change="handleAdvancedFieldChange"
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
