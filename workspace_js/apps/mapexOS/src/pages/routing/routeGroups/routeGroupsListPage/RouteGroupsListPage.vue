<script setup lang="ts">
defineOptions({
  name: 'RouteGroupsListPage'
});

/** TYPE IMPORTS */
import type { ListHeaderMenuColumn } from '@components/headers';
import type { DataRowActionConfig } from '@components/cards';
import type {
  RouteGroupsListPageFilters,
  RouteGroupsListPageColumnVisibility,
  EnrichedRouteGroup,
} from './interfaces';

/** VUE IMPORTS */
import { ref, computed, onMounted } from 'vue';
import { useRouter } from 'vue-router';

/** COMPONENTS */
import { PageHeader, ListHeaderMenu } from '@components/headers';
import { ListCardEmpty, DataRow } from '@components/cards';
import { ListPagination } from '@components/navigation';
import { RouteGroupDetailsDrawer, AdvancedFiltersDrawer } from '@components/drawers';
import { AppTooltip } from '@components/tooltips';

/** COMPOSABLES */
import { useRouteGroupsTranslations } from '@composables/i18n/pages/routing/routeGroups/useRouteGroupsTranslations';
import { useOrgChangeRefresh } from '@composables/organizations';
import { usePermissions } from '@composables/shared/usePermissions';

/** STORES */
import { useOrganizationStore } from '@stores/organization';

/** LOCAL IMPORTS */
import {
  DEFAULT_ITEMS_PER_PAGE,
  COLUMN_VISIBILITY_DEFAULTS,
  FILTER_DEFAULTS,
} from './constants';
import {
  fetchRouteGroupsHandler,
  handlePageChangeHandler,
  handleItemsPerPageChangeHandler,
  handleColumnsUpdateHandler,
  editRouteGroupHandler,
  confirmDeleteHandler,
  deleteRouteGroupHandler,
} from './handlers';

/** CONSTANTS */
const MAX_VISIBLE_CHIPS = 2;

/** COMPOSABLES & STORES */
const t = useRouteGroupsTranslations();
const orgStore = useOrganizationStore();
const router = useRouter();
const { canCreate, canUpdate, canDelete, canRead } = usePermissions();
const canCreateRG = canCreate('routegroups');
const canUpdateRG = canUpdate('routegroups');
const canDeleteRG = canDelete('routegroups');
const canReadRG = canRead('routegroups');

/** STATE */
const routeGroupsList = ref<EnrichedRouteGroup[]>([]);
const loading = ref(false);
const error = ref<string | undefined>(undefined);
const itemsPerPage = ref(DEFAULT_ITEMS_PER_PAGE);
const currentPage = ref(1);
const totalPages = ref(1);
const totalItems = ref(0);
const hasNext = ref(false);
const hasPrev = ref(false);
const filters = ref<RouteGroupsListPageFilters>({ ...FILTER_DEFAULTS });
const columnVisibilityState = ref<RouteGroupsListPageColumnVisibility>({ ...COLUMN_VISIBILITY_DEFAULTS });
const lastUpdatedAt = ref<number | undefined>(undefined);

// Quick filters state
const quickSearch = ref('');
const quickStatus = ref<boolean | null>(null);

// Advanced filters state
const showAdvancedFilters = ref(false);
const advancedFilterValues = ref<Record<string, any>>({
  includeChildren: null,
  isTemplate: undefined,
});
const hasPendingAdvancedFilters = ref(false);

// Drawer state
const showDetailsDrawer = ref(false);
const selectedRouteGroupId = ref<string | undefined>(undefined);

/** COMPUTED */

/**
 * Status select options for quick filter
 */
const statusOptions = computed(() => [
  { label: t.filters.allStatus.value, value: null },
  { label: t.filters.options.active.value, value: true },
  { label: t.filters.options.inactive.value, value: false },
]);

/**
 * Count of active advanced filters
 */
const advancedFiltersCount = computed(() => {
  let count = 0;
  if (advancedFilterValues.value.includeChildren !== null) count++;
  if (advancedFilterValues.value.isTemplate !== undefined) count++;
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
      label: t.filters.name.value,
      value: quickSearch.value,
    });
  }

  if (quickStatus.value !== null) {
    const statusLabel = quickStatus.value
      ? t.filters.options.active.value
      : t.filters.options.inactive.value;
    chips.push({
      key: 'enabled',
      label: t.filters.enabled.value,
      value: statusLabel,
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

  if (advancedFilterValues.value.isTemplate !== undefined) {
    const templateLabel = advancedFilterValues.value.isTemplate
      ? t.filters.options.templates.value
      : t.filters.options.local.value;
    chips.push({
      key: 'isTemplate',
      label: t.filters.isTemplate.value,
      value: templateLabel,
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
 * Column visibility using ListHeaderMenuColumn format with reactive translations
 */
const menuColumns = computed(() => {
  const cols: ListHeaderMenuColumn[] = [
    { key: 'routers', label: t.menuColumns.routers.value, visible: columnVisibilityState.value.routers },
    { key: 'isTemplate', label: t.menuColumns.templateSource.value, visible: columnVisibilityState.value.isTemplate },
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
 * Filtered columns based on visibility
 */
const visibleColumns = computed(() => {
  return t.columns.value.filter((col: any) => {
    // Always show icon, name, and status
    if (col.key === 'icon' || col.key === 'name' || col.key === 'status') {
      return true;
    }

    // Organization column only visible when includeChildren filter is active
    if (col.key === 'organizationName') {
      return filters.value.includeChildren === true && columnVisibilityState.value.organization;
    }

    // Filter based on columnVisibility
    if (col.key === 'routersCount') return columnVisibilityState.value.routers;
    if (col.key === 'isTemplate') return columnVisibilityState.value.isTemplate;

    return true;
  });
});

/** FUNCTIONS */

/**
 * Fetch route groups from API with current filters and pagination
 * @returns {Promise<void>}
 */
async function fetchRouteGroups(): Promise<void> {
  await fetchRouteGroupsHandler(
    filters,
    currentPage,
    itemsPerPage,
    routeGroupsList,
    totalPages,
    totalItems,
    hasNext,
    hasPrev,
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
  void fetchRouteGroups();
}

/**
 * Handle advanced filters apply from drawer
 * @param {Record<string, any>} appliedFilters - Applied filter values
 * @returns {void}
 */
function handleAdvancedFiltersApply(appliedFilters: Record<string, any>): void {
  advancedFilterValues.value = { ...appliedFilters };
  filters.value.includeChildren = appliedFilters.includeChildren ?? undefined;
  filters.value.isTemplate = appliedFilters.isTemplate;
  hasPendingAdvancedFilters.value = false;
  currentPage.value = 1;
  showAdvancedFilters.value = false;
  void fetchRouteGroups();
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
  } else if (key === 'includeChildren') {
    advancedFilterValues.value.includeChildren = null;
    filters.value.includeChildren = undefined;
  } else if (key === 'isTemplate') {
    advancedFilterValues.value.isTemplate = undefined;
    filters.value.isTemplate = undefined;
  }
  currentPage.value = 1;
  void fetchRouteGroups();
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
    isTemplate: undefined,
  };
  filters.value = { ...FILTER_DEFAULTS };
  hasPendingAdvancedFilters.value = false;
  currentPage.value = 1;
  void fetchRouteGroups();
}

/**
 * Handle pagination navigation
 * @param {number} page - New page number
 * @returns {void}
 */
function handlePageChange(page: number): void {
  handlePageChangeHandler(page, currentPage, () => void fetchRouteGroups());
}

/**
 * Handle items per page change
 * @param {number} newValue - New items per page value
 * @returns {void}
 */
function handleItemsPerPageChange(newValue: number): void {
  handleItemsPerPageChangeHandler(newValue, itemsPerPage, currentPage, () => void fetchRouteGroups());
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
 * View route group details in drawer
 * @param {EnrichedRouteGroup} routeGroup - Route group to view
 * @returns {void}
 */
function viewDetails(routeGroup: EnrichedRouteGroup): void {
  if (!canReadRG.value) return;
  if (!routeGroup.id) {
    return;
  }
  selectedRouteGroupId.value = routeGroup.id;
  showDetailsDrawer.value = true;
}

/**
 * Edit route group
 * @param {EnrichedRouteGroup} routeGroup - Route group to edit
 * @returns {void}
 */
function editRouteGroup(routeGroup: EnrichedRouteGroup): void {
  if (!canUpdateRG.value) return;
  editRouteGroupHandler(routeGroup, orgStore.selectedOrganizationId, router, t);
}

/**
 * Confirm delete route group with dialog
 * @param {EnrichedRouteGroup} routeGroup - Route group to delete
 * @returns {Promise<void>}
 */
async function confirmDelete(routeGroup: EnrichedRouteGroup): Promise<void> {
  if (!canDeleteRG.value) return;
  await confirmDeleteHandler(routeGroup, orgStore.selectedOrganizationId, t, deleteRouteGroup);
}

/**
 * Delete route group from API and update list
 * @param {EnrichedRouteGroup} routeGroup - Route group to delete
 * @returns {Promise<void>}
 */
async function deleteRouteGroup(routeGroup: EnrichedRouteGroup): Promise<void> {
  await deleteRouteGroupHandler(routeGroup, routeGroupsList, totalItems, currentPage, t, fetchRouteGroups);
}

/**
 * Navigate to add route group page
 * @returns {void}
 */
function goToAddRouteGroup(): void {
  void router.push('/routing/route_groups/add');
}

/**
 * Get actions configuration for a route group
 * @param {EnrichedRouteGroup} _routeGroup - Route group (unused)
 * @returns {DataRowActionConfig}
 */
// eslint-disable-next-line @typescript-eslint/no-unused-vars
function getRouteGroupActions(_routeGroup: EnrichedRouteGroup): DataRowActionConfig {
  return {
    showEdit: canUpdateRG.value,
    showView: canReadRG.value,
    showDelete: canDeleteRG.value,
  };
}

/** LIFECYCLE HOOKS */
onMounted(async () => await fetchRouteGroups());

/**
 * Auto-refresh route groups when organization changes
 */
useOrgChangeRefresh(async () => {
  currentPage.value = 1;
  quickSearch.value = '';
  quickStatus.value = null;
  advancedFilterValues.value = {
    includeChildren: null,
    isTemplate: undefined,
  };
  filters.value = { ...FILTER_DEFAULTS };
  hasPendingAdvancedFilters.value = false;
  await fetchRouteGroups();
});
</script>

<template>
  <q-page class="q-pt-lg">

    <!-- Header Section -->
    <PageHeader
        icon="route"
        iconColor="primary"
        :title="t.pageHeader.title.value"
        :description="t.pageHeader.description.value"
        :button="canCreateRG ? { label: t.pageHeader.button.value, icon: 'add', onClick: goToAddRouteGroup, color: 'primary' } : undefined"
        :info="t.pageHeader.info.value"
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
          :label="t.filters.enabled.value"
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
          <q-icon name="route" size="sm" color="primary" class="q-mr-sm"/>
          <div class="text-subtitle1 text-weight-medium text-primary">{{ t.listHeader.title.value }}</div>
        </div>
      </div>
      <div class="col-auto">
        <ListHeaderMenu
          icon="route"
          :items-count="totalItems"
          :item-label="t.listHeader.itemLabel.value"
          :item-label-plural="t.listHeader.itemLabelPlural.value"
          :items-per-page="itemsPerPage"
          :columns="menuColumns"
          :refreshing="loading"
          :last-updated-at="lastUpdatedAt"
          @update:items-per-page="handleItemsPerPageChange"
          @update:columns="handleColumnsUpdate"
          @refresh="fetchRouteGroups"
        />
      </div>
    </div>

    <!-- Loading Spinner -->
    <div v-if="loading" class="row justify-center q-my-lg">
      <q-spinner color="primary" size="3em" />
    </div>

    <!-- Route Groups Row List -->
    <div v-else class="row">
      <div
          v-for="routeGroup in routeGroupsList"
          :key="routeGroup.id || Math.random()"
          class="col-12 q-mb-xs"
      >
        <DataRow
            :data="routeGroup"
            :columns="visibleColumns"
            :actions="getRouteGroupActions(routeGroup)"
            @click="viewDetails"
            @dblclick="editRouteGroup"
            @edit="editRouteGroup"
            @view="viewDetails"
            @delete="confirmDelete"
        />
      </div>

      <!-- No Results -->
      <div v-if="routeGroupsList.length === 0" class="col-12">
        <ListCardEmpty
            :title="t.empty.title.value"
            :description="t.empty.description.value"
            icon="route"
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
    <RouteGroupDetailsDrawer
      v-model="showDetailsDrawer"
      :route-group-id="selectedRouteGroupId"
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
