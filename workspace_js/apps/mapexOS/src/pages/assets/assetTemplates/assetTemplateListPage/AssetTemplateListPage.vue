<script setup lang="ts">
defineOptions({
  name: 'AssetTemplateListPage'
});

/** TYPE IMPORTS */
import type { ListHeaderMenuColumn } from '@components/headers';
import type { DataRowActionConfig } from '@components/cards';
import type { FilterField } from '@components/drawers';
import type {
  AssetTemplateListPageFilters,
  AssetTemplateListPageColumnVisibility,
  DynamicFilterOptions,
  EnrichedAssetTemplate,
} from './interfaces';

/** VUE IMPORTS */
import { ref, computed, onMounted } from 'vue';
import { useRouter } from 'vue-router';

/** COMPONENTS */
import { PageHeader, ListHeaderMenu } from '@components/headers';
import { ListCardEmpty, DataRow } from '@components/cards';
import { AssetTemplateDetailsDrawer, AdvancedFiltersDrawer } from '@components/drawers';
import { ListPagination } from '@components/navigation';
import { AppTooltip } from '@components/tooltips';

/** COMPOSABLES */
import { useAssetTemplatesTranslations } from '@composables/i18n';
import { useOrgChangeRefresh } from '@composables/organizations';
import { usePermissions } from '@composables/shared/usePermissions';
import { useLogger } from '@composables/useLogger';

/** UTILS */
import { notifyFail, notifySuccess, notifyWarning, dialogDelete } from '@utils/alert';
import { cleanQueryParams } from '@utils/query';

/** SERVICES */
import { apis } from '@services/mapex';

/** STORES */
import { useOrganizationStore } from '@stores/organization';

/** LOCAL IMPORTS */
import {
  DEFAULT_ITEMS_PER_PAGE,
  ASSET_TEMPLATE_COLUMN_VISIBILITY_DEFAULTS,
  ASSET_TEMPLATE_FILTER_DEFAULTS,
  ASSET_TEMPLATE_PROJECTION,
  LIST_TYPE,
  CASCADING_FILTER_DEFAULTS,
} from './constants';

/** COMPOSABLES & STORES */
const t = useAssetTemplatesTranslations();
const orgStore = useOrganizationStore();
const router = useRouter();
const logger = useLogger('AssetTemplateListPage');
const { canCreate, canUpdate, canDelete, canRead } = usePermissions();
const canCreateTemplate = canCreate('assettemplates');
const canUpdateTemplate = canUpdate('assettemplates');
const canDeleteTemplate = canDelete('assettemplates');
const canReadTemplate = canRead('assettemplates');

/** STATE */
const assetTemplatesList = ref<EnrichedAssetTemplate[]>([]);
const loading = ref(false);
const lastUpdatedAt = ref<number | undefined>(undefined);
const error = ref<string | undefined>(undefined);
const showDetailsDrawer = ref(false);
const showFiltersDrawer = ref(false);
const selectedTemplateId = ref<string | null>(null);
const itemsPerPage = ref(DEFAULT_ITEMS_PER_PAGE);
const currentPage = ref(1);
const totalPages = ref(1);
const totalItems = ref(0);
const hasNext = ref(false);
const hasPrev = ref(false);
const filters = ref<AssetTemplateListPageFilters>({ ...ASSET_TEMPLATE_FILTER_DEFAULTS });
const columnVisibilityState = ref<AssetTemplateListPageColumnVisibility>({ ...ASSET_TEMPLATE_COLUMN_VISIBILITY_DEFAULTS });

// Quick filter state (inline)
const quickSearch = ref('');
const quickStatus = ref<boolean | null>(null);

// Advanced filter state (drawer)
const advancedFilterValues = ref<Record<string, any>>({
  includeChildren: null,
  isSystem: null,
  isTemplate: null,
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
  { label: t.quickFilters.allStatus.value, value: null },
  { label: t.quickFilters.options.enabled.value, value: true },
  { label: t.quickFilters.options.disabled.value, value: false },
]);

/**
 * Advanced filter fields for the drawer
 */
const advancedFilterFields = computed((): FilterField[] => {
  const fields: FilterField[] = [
    {
      key: 'includeChildren',
      label: t.quickFilters.includeChildren.value,
      type: 'toggle',
      icon: 'account_tree',
      options: [
        { label: t.quickFilters.allStatus.value, value: null },
        { label: t.quickFilters.options.yes.value, value: true },
        { label: t.quickFilters.options.no.value, value: false },
      ],
    },
    {
      key: 'isSystem',
      label: t.quickFilters.isSystem.value,
      type: 'toggle',
      icon: 'lock',
      options: [
        { label: t.quickFilters.allStatus.value, value: null },
        { label: t.quickFilters.options.system.value, value: true },
        { label: t.quickFilters.options.custom.value, value: false },
      ],
    },
  ];

  // Add isTemplate filter only for Customer and Site organizations
  if (orgStore.isCustomer || orgStore.isSite) {
    fields.push({
      key: 'isTemplate',
      label: t.quickFilters.isTemplate.value,
      type: 'toggle',
      icon: 'content_copy',
      options: [
        { label: t.quickFilters.allStatus.value, value: null },
        { label: t.quickFilters.options.shared.value, value: true },
        { label: t.quickFilters.options.local.value, value: false },
      ],
    });
  }

  // Cascading filters
  fields.push(
    {
      key: 'categoryId',
      label: t.quickFilters.category.value,
      type: 'select',
      icon: 'category',
      options: categoryOptions.value,
      loading: loadingCategories.value,
      placeholder: t.quickFilters.filterByCategory.value,
    },
    {
      key: 'manufacturerId',
      label: t.quickFilters.manufacturer.value,
      type: 'select',
      icon: 'factory',
      options: manufacturerOptions.value,
      loading: loadingManufacturers.value,
      disabled: !advancedFilterValues.value.categoryId,
      placeholder: t.quickFilters.filterByManufacturer.value,
    },
    {
      key: 'modelId',
      label: t.quickFilters.model.value,
      type: 'select',
      icon: 'memory',
      options: modelOptions.value,
      loading: loadingModels.value,
      disabled: !advancedFilterValues.value.manufacturerId,
      placeholder: t.quickFilters.filterByModel.value,
    }
  );

  return fields;
});

/**
 * Check if any filters are active (quick or advanced)
 */
const hasActiveFilters = computed(() => {
  return !!(
    quickSearch.value ||
    quickStatus.value !== null ||
    advancedFilterValues.value.includeChildren !== null ||
    advancedFilterValues.value.isSystem !== null ||
    advancedFilterValues.value.isTemplate !== null ||
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
  if (advancedFilterValues.value.isSystem !== null) count++;
  if (advancedFilterValues.value.isTemplate !== null) count++;
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
      label: t.quickFilters.name.value,
      value: quickSearch.value,
    });
  }

  if (quickStatus.value !== null) {
    chips.push({
      key: 'status',
      label: t.quickFilters.status.value,
      value: quickStatus.value ? t.quickFilters.options.enabled.value : t.quickFilters.options.disabled.value,
    });
  }

  if (advancedFilterValues.value.includeChildren !== null) {
    chips.push({
      key: 'includeChildren',
      label: t.quickFilters.includeChildren.value,
      value: advancedFilterValues.value.includeChildren ? t.quickFilters.options.yes.value : t.quickFilters.options.no.value,
    });
  }

  if (advancedFilterValues.value.isSystem !== null) {
    chips.push({
      key: 'isSystem',
      label: t.quickFilters.isSystem.value,
      value: advancedFilterValues.value.isSystem ? t.quickFilters.options.system.value : t.quickFilters.options.custom.value,
    });
  }

  if (advancedFilterValues.value.isTemplate !== null) {
    chips.push({
      key: 'isTemplate',
      label: t.quickFilters.isTemplate.value,
      value: advancedFilterValues.value.isTemplate ? t.quickFilters.options.shared.value : t.quickFilters.options.local.value,
    });
  }

  if (advancedFilterValues.value.categoryId) {
    const category = categoryOptions.value.find(c => c.value === advancedFilterValues.value.categoryId);
    chips.push({
      key: 'categoryId',
      label: t.quickFilters.category.value,
      value: category?.label || advancedFilterValues.value.categoryId,
    });
  }

  if (advancedFilterValues.value.manufacturerId) {
    const manufacturer = manufacturerOptions.value.find(m => m.value === advancedFilterValues.value.manufacturerId);
    chips.push({
      key: 'manufacturerId',
      label: t.quickFilters.manufacturer.value,
      value: manufacturer?.label || advancedFilterValues.value.manufacturerId,
    });
  }

  if (advancedFilterValues.value.modelId) {
    const model = modelOptions.value.find(m => m.value === advancedFilterValues.value.modelId);
    chips.push({
      key: 'modelId',
      label: t.quickFilters.model.value,
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
    { key: 'manufacturerModel', label: t.menuColumns.manufacturerModel.value, visible: columnVisibilityState.value.manufacturerModel },
    { key: 'version', label: t.menuColumns.version.value, visible: columnVisibilityState.value.version },
    { key: 'isSystem', label: t.menuColumns.templateType.value, visible: columnVisibilityState.value.isSystem },
    { key: 'isTemplate', label: t.menuColumns.templateSource.value, visible: columnVisibilityState.value.isTemplate },
  ];

  // Only show organization toggle when includeChildren is active
  if (advancedFilterValues.value.includeChildren === true) {
    cols.unshift({
      key: 'organization',
      label: t.menuColumns.organization.value,
      visible: columnVisibilityState.value.organization
    });
  }

  return cols;
});

/**
 * Filtered columns based on visibility
 */
const visibleColumns = computed(() => {
  return t.columns.value.filter((col: any) => {
    // Always show icon and name
    if (col.key === 'icon' || col.key === 'name') {
      return true;
    }

    // Organization column only visible when includeChildren filter is active
    if (col.key === 'organizationName') {
      return advancedFilterValues.value.includeChildren === true && columnVisibilityState.value.organization;
    }

    // Filter based on columnVisibility
    if (col.key === 'manufacturer') return columnVisibilityState.value.manufacturerModel;
    if (col.key === 'version') return columnVisibilityState.value.version;
    if (col.key === 'isSystem') return columnVisibilityState.value.isSystem;
    if (col.key === 'isTemplate') return columnVisibilityState.value.isTemplate;

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
  void fetchAssetTemplates();
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
  filters.value.isSystem = appliedFilters.isSystem;
  filters.value.isTemplate = appliedFilters.isTemplate;
  filters.value.categoryId = appliedFilters.categoryId || undefined;
  filters.value.manufacturerId = appliedFilters.manufacturerId || undefined;
  filters.value.modelId = appliedFilters.modelId || undefined;

  // Reset pending state after apply
  hasPendingAdvancedFilters.value = false;

  // Auto-hide less important columns when includeChildren is active
  if (appliedFilters.includeChildren === true) {
    columnVisibilityState.value.isSystem = false;
    columnVisibilityState.value.isTemplate = false;
  } else {
    columnVisibilityState.value.isSystem = true;
    columnVisibilityState.value.isTemplate = true;
  }

  currentPage.value = 1;
  showFiltersDrawer.value = false;
  void fetchAssetTemplates();
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
    isSystem: null,
    isTemplate: null,
    categoryId: null,
    manufacturerId: null,
    modelId: null,
  };

  // Reset cascading filter options
  manufacturerOptions.value = [];
  modelOptions.value = [];

  filters.value.includeChildren = undefined;
  filters.value.isSystem = undefined;
  filters.value.isTemplate = undefined;
  filters.value.categoryId = undefined;
  filters.value.manufacturerId = undefined;
  filters.value.modelId = undefined;

  // Restore column visibility
  columnVisibilityState.value.isSystem = true;
  columnVisibilityState.value.isTemplate = true;

  currentPage.value = 1;
  void fetchAssetTemplates();
}

/**
 * Handle field change in advanced filters (for cascading logic)
 * @param {string} key - Field key that changed
 * @param {any} value - New value
 * @returns {void}
 */
function handleAdvancedFieldChange(key: string, value: any): void {
  logger.debug(`Advanced field change: ${key} = ${value}`);

  // Update the field value in advancedFilterValues (required for cascading to work)
  advancedFilterValues.value[key] = value;

  // When Template Source is set (Shared or Local), force Template Type to Custom
  // System templates are never Shared/Local — they are global
  if (key === 'isTemplate' && value !== null) {
    advancedFilterValues.value.isSystem = false;
  }

  // When Template Type is set to System, reset Template Source (not applicable)
  if (key === 'isSystem' && value === true) {
    advancedFilterValues.value.isTemplate = null;
  }

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
    // Restore column visibility
    columnVisibilityState.value.isSystem = true;
    columnVisibilityState.value.isTemplate = true;
  } else if (key === 'isSystem') {
    advancedFilterValues.value.isSystem = null;
    filters.value.isSystem = undefined;
  } else if (key === 'isTemplate') {
    advancedFilterValues.value.isTemplate = null;
    filters.value.isTemplate = undefined;
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
  void fetchAssetTemplates();
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
    isSystem: null,
    isTemplate: null,
    categoryId: null,
    manufacturerId: null,
    modelId: null,
  };
  manufacturerOptions.value = [];
  modelOptions.value = [];
  filters.value = { ...ASSET_TEMPLATE_FILTER_DEFAULTS };

  // Reset pending state
  hasPendingAdvancedFilters.value = false;

  // Restore column visibility
  columnVisibilityState.value.isSystem = true;
  columnVisibilityState.value.isTemplate = true;

  currentPage.value = 1;
  void fetchAssetTemplates();
}

/**
 * Fetch asset templates from API with current filters and pagination
 * @returns {Promise<void>}
 */
async function fetchAssetTemplates(): Promise<void> {
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
      projection: ASSET_TEMPLATE_PROJECTION,
    };

    // Add active filters to query params (only if they have values)
    if (filters.value.name) {
      queryParams.name = filters.value.name;
    }
    if (typeof filters.value.status === 'boolean') {
      queryParams.enabled = filters.value.status;
    }
    if (typeof filters.value.isSystem === 'boolean') {
      queryParams.isSystem = filters.value.isSystem;
    }
    if (typeof filters.value.isTemplate === 'boolean') {
      queryParams.isTemplate = filters.value.isTemplate;
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

    const response = await apis.assets.assetTemplate.list(cleanedParams);

    // Safely access response data with proper null checks
    // API returns PaginatedResponse<AssetTemplateResponse>
    const templatesData = response?.items || [];

    // Enrich templates with computed fields and organization name
    const enrichedTemplates = templatesData.map((template: any) => {
      const organization = orgStore.flatList.find((org: any) => org.id === template.orgId);
      return {
        ...template,
        organizationName: organization?.name || 'Unknown',
        // Compute script flags for display
        hasPreprocessor: !!template.scriptProcessor,
        hasValidation: !!template.scriptValidator,
        hasConversion: !!template.scriptConversion,
        // Map backend fields to frontend display fields (for column compatibility)
        manufacturer: template.manufacturerName,
        deviceModel: template.modelName,
      };
    });

    assetTemplatesList.value = enrichedTemplates;

    // Update pagination state from response
    if (response.pagination) {
      totalItems.value = response.pagination.totalItems || 0;
      totalPages.value = response.pagination.totalPages || 1;
      // Calculate hasNext and hasPrev based on current page and total pages
      hasNext.value = currentPage.value < totalPages.value;
      hasPrev.value = currentPage.value > 1;
    }
  } catch (err: any) {
    logger.error('Error fetching asset templates:', err);
    const errorMsg = err.message || 'Failed to fetch asset templates';
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
  void fetchAssetTemplates();
}

/**
 * Handle items per page change
 * @param {number} newValue - New items per page value
 * @returns {void}
 */
function handleItemsPerPageChange(newValue: number): void {
  itemsPerPage.value = newValue;
  currentPage.value = 1; // Reset to first page
  void fetchAssetTemplates();
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
    if (col.key === 'organization') columnVisibilityState.value.organization = col.visible;
    if (col.key === 'manufacturerModel') columnVisibilityState.value.manufacturerModel = col.visible;
    if (col.key === 'version') columnVisibilityState.value.version = col.visible;
    if (col.key === 'isSystem') columnVisibilityState.value.isSystem = col.visible;
    if (col.key === 'isTemplate') columnVisibilityState.value.isTemplate = col.visible;
  });
}

/**
 * Check if user can edit/delete a template
 * Rules:
 * - isSystem = true: Cannot edit/delete (system resource)
 * - isTemplate = true: Can only edit/delete if orgId matches current organization
 * - isTemplate = false: Can always edit/delete (local resource)
 * @param {EnrichedAssetTemplate} template - Template to check modification permissions
 * @returns {boolean} True if user can modify the template
 */
function canModifyTemplate(template: EnrichedAssetTemplate): boolean {
  // System templates cannot be modified
  if (template.isSystem) {
    return false;
  }

  // Shared templates can only be modified by the owner organization
  if (template.isTemplate) {
    return template.orgId === orgStore.selectedOrganizationId;
  }

  // Local resources can always be modified
  return true;
}

/**
 * View template details in drawer
 * IMPORTANT: Pass only the template ID, not the full template object
 * The drawer will fetch complete template data using the API
 * @param {EnrichedAssetTemplate} template - Template to view
 * @returns {void}
 */
function viewDetails(template: EnrichedAssetTemplate): void {
  if (!canReadTemplate.value) return;
  if (!template.id) {
    notifyFail({ message: t.errors.idMissing.value });
    return;
  }
  selectedTemplateId.value = template.id;
  showDetailsDrawer.value = true;
}

/**
 * Edit template (with access control validation)
 * @param {EnrichedAssetTemplate} template - Template to edit
 * @returns {void}
 */
function editTemplate(template: EnrichedAssetTemplate): void {
  if (!canUpdateTemplate.value) return;
  if (!canModifyTemplate(template)) {
    if (template.isSystem) {
      notifyWarning({ message: t.notifications.systemTemplateEdit.value });
    } else if (template.isTemplate) {
      notifyWarning({ message: t.notifications.sharedTemplateEdit.value });
    }
    return;
  }

  // Navigate to edit page
  void router.push(`/assets_template/edit/${template.id}`);
}

/**
 * Handle edit from drawer
 * @param {string} templateId - ID of template to edit
 * @returns {void}
 */
function handleDrawerEdit(templateId: string): void {
  if (!canUpdateTemplate.value) return;
  // Close drawer before navigating
  showDetailsDrawer.value = false;

  // Navigate to edit page
  void router.push(`/assets_template/edit/${templateId}`);
}

/**
 * Handle duplicate from drawer
 * @param {string} templateId - ID of template to duplicate
 * @returns {void}
 */
function handleDrawerDuplicate(templateId: string): void {
  logger.debug('Duplicate template from drawer:', templateId);
  // TODO: Implement duplicate functionality
  // router.push(`/assets_template/duplicate/${templateId}`);
}

/**
 * Confirm delete template (with access control validation)
 * @param {EnrichedAssetTemplate} template - Template to delete
 * @returns {Promise<void>}
 */
async function confirmDelete(template: EnrichedAssetTemplate): Promise<void> {
  if (!canDeleteTemplate.value) return;
  if (!canModifyTemplate(template)) {
    if (template.isSystem) {
      notifyWarning({ message: t.notifications.systemTemplateDelete.value });
    } else if (template.isTemplate) {
      notifyWarning({ message: t.notifications.sharedTemplateDelete.value });
    }
    return;
  }

  const confirmed = await dialogDelete({
    title: t.dialog.confirmDelete.title.value,
    message: t.dialog.confirmDelete.message(template.name || ''),
  });

  if (confirmed) {
    await deleteTemplate(template);
  }
}

/**
 * Delete template from API and update list
 * @param {EnrichedAssetTemplate} template - Template to delete
 * @returns {Promise<void>}
 */
async function deleteTemplate(template: EnrichedAssetTemplate): Promise<void> {
  if (!apis.assets) {
    notifyFail({ message: t.errors.apiNotInitialized.value });
    return;
  }

  if (!template.id) {
    notifyFail({ message: t.errors.idMissing.value });
    return;
  }

  try {
    await apis.assets.assetTemplate.delete({ assetTemplateId: template.id });

    // Update total items count after successful deletion
    totalItems.value = Math.max(0, totalItems.value - 1);

    // Remove from local list
    const index = assetTemplatesList.value.findIndex(r => r.id === template.id);
    if (index !== -1) {
      assetTemplatesList.value.splice(index, 1);
    }

    // If current page is now empty and not the first page, go to previous page
    if (assetTemplatesList.value.length === 0 && currentPage.value > 1) {
      currentPage.value -= 1;
      await fetchAssetTemplates();
    }

    notifySuccess({ message: t.notifications.deleted.value });
  } catch (err: any) {
    notifyFail({ message: err.message || t.notifications.deleteError.value });
  }
}

/**
 * Get actions configuration for a template
 * Combines RBAC permission checks with data-level access control
 * @param {EnrichedAssetTemplate} template - Template to get actions for
 * @returns {DataRowActionConfig} Action configuration object
 */
function getTemplateActions(template: EnrichedAssetTemplate): DataRowActionConfig {
  return {
    showEdit: canUpdateTemplate.value && canModifyTemplate(template),
    showView: canReadTemplate.value,
    showDelete: canDeleteTemplate.value && canModifyTemplate(template),
  };
}

/** WATCHERS */

/** LIFECYCLE HOOKS */

onMounted(async () => {
  await fetchCategories();
  await fetchAssetTemplates();
});

/**
 * Auto-refresh asset templates when organization changes
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
    isSystem: null,
    isTemplate: null,
    categoryId: null,
    manufacturerId: null,
    modelId: null,
  };
  manufacturerOptions.value = [];
  modelOptions.value = [];
  filters.value = { ...ASSET_TEMPLATE_FILTER_DEFAULTS };

  // Restore column visibility
  columnVisibilityState.value.isSystem = true;
  columnVisibilityState.value.isTemplate = true;

  await fetchAssetTemplates();
});

</script>

<template>
  <q-page class="q-pt-lg">

    <!-- Header Section -->
    <PageHeader
        icon="memory"
        iconColor="primary"
        :title="t.pageHeader.title.value"
        :description="t.pageHeader.description.value"
        :button="canCreateTemplate ? { label: t.pageHeader.button.value, icon: 'add', to: '/assets_template/add', color: 'primary' } : undefined"
        :info="t.pageHeader.info.value"
    />

    <!-- Filters Section -->
    <div class="text-caption text-grey-7 q-mb-xs">{{ t.quickFilters.label.value }}</div>
    <div class="row items-center q-col-gutter-sm q-mb-md">
      <!-- Search Input -->
      <div class="col">
        <q-input
          v-model="quickSearch"
          outlined
          dense
          clearable
          :placeholder="t.quickFilters.searchPlaceholder.value"
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
          :label="t.quickFilters.status.value"
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
            ? t.quickFilters.pendingFilters.value
            : t.quickFilters.advancedFilters.value"
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
        :label="t.quickFilters.clearAll.value"
        no-caps
        class="q-ml-sm"
        @click="clearAllFilters"
      />
    </div>

    <!-- Results Section -->
    <div class="row items-center q-pt-xl q-mb-md">
      <div class="col">
        <div class="row items-center">
          <q-icon name="memory" size="sm" color="primary" class="q-mr-sm"/>
          <div class="text-subtitle1 text-weight-medium text-primary">{{ t.listHeader.title.value }}</div>
        </div>
      </div>
      <div class="col-auto">
        <ListHeaderMenu
          icon="memory"
          :items-count="totalItems"
          :item-label="t.listHeader.itemLabel.value"
          :item-label-plural="t.listHeader.itemLabelPlural.value"
          :items-per-page="itemsPerPage"
          :columns="menuColumns"
          :filtered="hasActiveFilters"
          :refreshing="loading"
          :last-updated-at="lastUpdatedAt"
          @update:items-per-page="handleItemsPerPageChange"
          @update:columns="handleColumnsUpdate"
          @refresh="fetchAssetTemplates"
        />
      </div>
    </div>

    <!-- Loading Spinner -->
    <div v-if="loading" class="row justify-center q-my-lg">
      <q-spinner color="primary" size="3em" />
    </div>

    <!-- Asset Templates Row List -->
    <div v-else class="row">
      <div
          v-for="(assetTemplate, index) in assetTemplatesList"
          :key="assetTemplate.id || `template-${index}`"
          class="col-12 q-mb-xs"
      >
        <DataRow
            :data="assetTemplate"
            :columns="visibleColumns"
            :actions="getTemplateActions(assetTemplate)"
            @click="viewDetails"
            @dblclick="editTemplate"
            @edit="editTemplate"
            @view="viewDetails"
            @delete="confirmDelete"
        />
      </div>

      <!-- No Results -->
      <div v-if="assetTemplatesList.length === 0" class="col-12">
        <ListCardEmpty
            :title="t.empty.title.value"
            :description="t.empty.description.value"
            icon="memory"
        />
      </div>
    </div>

    <!-- Pagination -->
    <ListPagination
      v-model="currentPage"
      :total-pages="totalPages"
      @change="handlePageChange"
    />

    <!-- Asset Template Details Drawer -->
    <AssetTemplateDetailsDrawer
      v-model="showDetailsDrawer"
      :template-id="selectedTemplateId"
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
