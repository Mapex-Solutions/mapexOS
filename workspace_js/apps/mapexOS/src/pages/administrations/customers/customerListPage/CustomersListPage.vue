<script setup lang="ts">
defineOptions({
  name: 'CustomersListPage'
});

/** TYPE IMPORTS */
import type { ListHeaderMenuColumn } from '@components/headers';
import type { FilterField } from '@components/drawers';
import type { OrganizationResponse, OrganizationQuery } from '@mapexos/schemas';
import type {
  CustomersListPageFilters,
  CustomersListPageColumnVisibility,
} from './interfaces';

/** VUE IMPORTS */
import { ref, computed, onMounted } from 'vue';
import { useRouter } from 'vue-router';

/** COMPONENTS */
import { PageHeader, ListHeaderMenu } from '@components/headers';
import { ListCardEmpty, DataRow } from '@components/cards';
import { CustomerDetailsDrawer, AdvancedFiltersDrawer } from '@components/drawers';
import { ListPagination } from '@components/navigation';
import { AppTooltip } from '@components/tooltips';

/** COMPOSABLES */
import { useCustomersTranslations } from '@composables/i18n';
import { useLogger } from '@composables/useLogger';
import { usePermissions } from '@composables/shared/usePermissions';

/** UTILS */
import { notifySuccess, notifyFail, dialogDelete } from '@utils/alert';

/** SERVICES */
import { apis } from '@services/mapex';

/** STORES */
import { useOrganizationStore } from '@stores/organization';

/** LOCAL IMPORTS */
import {
  DEFAULT_ITEMS_PER_PAGE,
  COLUMN_VISIBILITY_DEFAULTS,
  FILTER_DEFAULTS,
  CUSTOMERS_PROJECTION,
} from './constants';

/** COMPOSABLES & STORES */
const t = useCustomersTranslations();
const orgStore = useOrganizationStore();
const logger = useLogger('CustomersListPage');
const router = useRouter();
const { canCreate, canUpdate, canDelete, canRead } = usePermissions();
const canCreateOrg = canCreate('organizations');
const canUpdateOrg = canUpdate('organizations');
const canDeleteOrg = canDelete('organizations');
const canReadOrg = canRead('organizations');

/** STATE */
const customersList = ref<OrganizationResponse[]>([]);
const loading = ref(false);
const lastUpdatedAt = ref<number | undefined>(undefined);
const error = ref<string | undefined>(undefined);
const itemsPerPage = ref(DEFAULT_ITEMS_PER_PAGE);
const currentPage = ref(1);
const totalPages = ref(1);
const totalItems = ref(0);
const filters = ref<CustomersListPageFilters>({ ...FILTER_DEFAULTS });
const columnVisibilityState = ref<CustomersListPageColumnVisibility>({ ...COLUMN_VISIBILITY_DEFAULTS });
const showDetailsDrawer = ref(false);
const showFiltersDrawer = ref(false);
const selectedCustomerId = ref<string | null>(null);

// Quick filter state (inline)
const quickSearch = ref('');
const quickStatus = ref<boolean | null>(null);

// Advanced filter state (drawer)
const advancedFilterValues = ref<Record<string, any>>({
  includeChildren: null,
  type: null,
});

// Pending changes state
const hasPendingAdvancedFilters = ref(false);

/** COMPUTED */

/**
 * Route for the add organization button
 * Parent org is resolved from the organization store (no query params needed)
 */
const addOrganizationRoute = '/customers/add';

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
    label: t.filters.includeChildrenOrgs.value,
    type: 'toggle',
    icon: 'account_tree',
    options: [
      { label: t.filters.allStatus.value, value: null },
      { label: t.filters.options.yes.value, value: true },
      { label: t.filters.options.no.value, value: false },
    ],
  },
  {
    key: 'type',
    label: t.filters.type.value,
    type: 'select',
    icon: 'category',
    placeholder: t.filters.type.value,
    options: [
      { label: t.filters.options.customer.value, value: 'customer' },
      { label: t.filters.options.site.value, value: 'site' },
      { label: t.filters.options.building.value, value: 'building' },
      { label: t.filters.options.floor.value, value: 'floor' },
      { label: t.filters.options.zone.value, value: 'zone' },
    ],
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
    advancedFilterValues.value.type
  );
});

/**
 * Count of active advanced filters (for badge)
 */
const advancedFiltersCount = computed(() => {
  let count = 0;
  if (advancedFilterValues.value.includeChildren !== null) count++;
  if (advancedFilterValues.value.type) count++;
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
      label: t.filters.organizationName.value,
      value: quickSearch.value,
    });
  }

  if (quickStatus.value !== null) {
    chips.push({
      key: 'enabled',
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

  if (advancedFilterValues.value.type) {
    const typeLabels: Record<string, string> = {
      customer: t.filters.options.customer.value,
      site: t.filters.options.site.value,
      building: t.filters.options.building.value,
      floor: t.filters.options.floor.value,
      zone: t.filters.options.zone.value,
    };
    chips.push({
      key: 'type',
      label: t.filters.type.value,
      value: typeLabels[advancedFilterValues.value.type] || advancedFilterValues.value.type,
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
    { key: 'address', label: t.menuColumns.address.value, visible: columnVisibilityState.value.address },
    { key: 'created', label: t.menuColumns.created.value, visible: columnVisibilityState.value.created },
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
    // Always show icon, name, and status
    if (col.key === 'icon' || col.key === 'name' || col.key === 'status') {
      return true;
    }

    // Organization column only visible when includeChildren filter is active
    if (col.key === 'organizationName') {
      return advancedFilterValues.value.includeChildren === true && columnVisibilityState.value.organization;
    }

    // Filter based on columnVisibility
    if (col.key === 'address.city') return columnVisibilityState.value.address;
    if (col.key === 'created') return columnVisibilityState.value.created;

    return true;
  });
});

/** FUNCTIONS */

/**
 * Fetch customers from API with current filters and pagination
 * @returns {Promise<void>}
 */
async function fetchCustomers(): Promise<void> {
  if (!apis.mapexOS?.organizations) {
    error.value = 'Organizations API not initialized';
    return;
  }

  try {
    loading.value = true;
    error.value = undefined;

    // Build query parameters
    const queryParams: Record<string, any> = {
      page: currentPage.value,
      perPage: itemsPerPage.value,
      projection: CUSTOMERS_PROJECTION,
    };

    // Add filters conditionally (only if they have values)
    if (filters.value.name) {
      queryParams.name = filters.value.name;
    }
    if (typeof filters.value.enabled === 'boolean') {
      queryParams.enabled = filters.value.enabled;
    }
    if (typeof filters.value.includeChildren === 'boolean') {
      queryParams.includeChildren = filters.value.includeChildren;
    }
    if (filters.value.type) {
      queryParams.type = filters.value.type;
    }

    const response = await apis.mapexOS.organizations.list(queryParams as OrganizationQuery);

    // Filter out current organization (exclude logged-in org from list)
    const filteredItems = response.items.filter((customer: any) => customer.id !== orgStore.selectedOrganizationId);

    // Enrich customers with organization name when includeChildren is active
    if (filters.value.includeChildren) {
      const enrichedCustomers = filteredItems.map((customer: any) => {
        const organization = orgStore.flatList.find(org => org.id === customer.customerId);
        return {
          ...customer,
          organizationName: organization?.name || 'Unknown',
        };
      });
      customersList.value = enrichedCustomers;
    } else {
      customersList.value = filteredItems;
    }

    // Update pagination state from response
    if (response.pagination) {
      totalPages.value = response.pagination.totalPages || 1;
      totalItems.value = response.pagination.totalItems || 0;
    }
  } catch (err: any) {
    logger.error('Error fetching customers:', err);
    const errorMsg = err.message || 'Failed to fetch customers';
    error.value = errorMsg;
    notifyFail({ message: errorMsg });
  } finally {
    loading.value = false;
    lastUpdatedAt.value = Date.now();
  }
}

/**
 * Apply quick filters and fetch data
 * @returns {void}
 */
function applyQuickFilters(): void {
  filters.value.name = quickSearch.value || undefined;
  filters.value.enabled = quickStatus.value ?? undefined;
  currentPage.value = 1;
  void fetchCustomers();
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
  filters.value.type = appliedFilters.type || undefined;

  // Reset pending state after apply
  hasPendingAdvancedFilters.value = false;

  currentPage.value = 1;
  showFiltersDrawer.value = false;
  void fetchCustomers();
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
    type: null,
  };

  filters.value.includeChildren = undefined;
  filters.value.type = undefined;

  currentPage.value = 1;
  void fetchCustomers();
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
  } else if (key === 'enabled') {
    quickStatus.value = null;
    filters.value.enabled = undefined;
  } else if (key === 'includeChildren') {
    advancedFilterValues.value.includeChildren = null;
    filters.value.includeChildren = undefined;
  } else if (key === 'type') {
    advancedFilterValues.value.type = null;
    filters.value.type = undefined;
  }

  currentPage.value = 1;
  void fetchCustomers();
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
    type: null,
  };
  filters.value = { ...FILTER_DEFAULTS };

  // Reset pending state
  hasPendingAdvancedFilters.value = false;

  currentPage.value = 1;
  void fetchCustomers();
}

/**
 * Handle pagination navigation
 * @param {number} page - New page number
 * @returns {void}
 */
function handlePageChange(page: number): void {
  currentPage.value = page;
  void fetchCustomers();
}

/**
 * Handle items per page change
 * @param {number} newValue - New items per page value
 * @returns {void}
 */
function handleItemsPerPageChange(newValue: number): void {
  itemsPerPage.value = newValue;
  currentPage.value = 1;
  void fetchCustomers();
}

/**
 * Update menu columns when changed
 * @param {ListHeaderMenuColumn[]} columns - Updated columns array
 * @returns {void}
 */
function handleColumnsUpdate(columns: ListHeaderMenuColumn[]): void {
  columns.forEach(col => {
    if (col.key === 'organization') columnVisibilityState.value.organization = col.visible;
    if (col.key === 'address') columnVisibilityState.value.address = col.visible;
    if (col.key === 'created') columnVisibilityState.value.created = col.visible;
  });
}

/**
 * View customer details in drawer
 * @param {OrganizationResponse} customer - Customer to view
 * @returns {void}
 */
function viewDetails(customer: OrganizationResponse): void {
  if (!canReadOrg.value) return;
  if (!customer.id) {
    notifyFail({ message: t.errors.idMissing.value });
    return;
  }
  selectedCustomerId.value = customer.id;
  showDetailsDrawer.value = true;
}

/**
 * Edit customer (navigate to edit page)
 * @param {OrganizationResponse} customer - Customer to edit
 * @returns {void}
 */
function editCustomer(customer: OrganizationResponse): void {
  if (!canUpdateOrg.value) return;
  if (!customer.id) {
    notifyFail({ message: t.errors.idMissing.value });
    return;
  }
  void router.push(`/customers/edit/${customer.id}`);
}

/**
 * Handle edit event from drawer
 * @param {string} customerId - ID of customer to edit
 * @returns {void}
 */
function handleDrawerEdit(customerId: string): void {
  if (!canUpdateOrg.value) return;
  showDetailsDrawer.value = false;
  void router.push(`/customers/edit/${customerId}`);
}

/**
 * Confirm delete operation with dialog
 * @param {OrganizationResponse} customer - Customer to delete
 * @returns {Promise<void>}
 */
async function confirmDelete(customer: OrganizationResponse): Promise<void> {
  if (!canDeleteOrg.value) return;
  const customerName = customer.name || 'this customer';
  const confirmed = await dialogDelete({
    title: t.dialog.deleteTitle.value,
    message: t.messages.confirmDelete(customerName),
  });

  if (confirmed) {
    await deleteCustomer(customer);
  }
}

/**
 * Delete customer from API and update list
 * @param {OrganizationResponse} customer - Customer to delete
 * @returns {Promise<void>}
 */
async function deleteCustomer(customer: OrganizationResponse): Promise<void> {
  if (!apis.mapexOS?.organizations) {
    notifyFail({ message: t.errors.apiNotInitialized.value });
    return;
  }

  if (!customer.id) {
    notifyFail({ message: t.errors.idMissing.value });
    return;
  }

  try {
    await apis.mapexOS.organizations.delete({ organizationId: customer.id });

    // Update total items count after successful deletion
    totalItems.value = Math.max(0, totalItems.value - 1);

    // Remove from local list
    const index = customersList.value.findIndex(r => r.id === customer.id);
    if (index !== -1) {
      customersList.value.splice(index, 1);
    }

    // If current page is now empty and not the first page, go to previous page
    if (customersList.value.length === 0 && currentPage.value > 1) {
      currentPage.value -= 1;
      await fetchCustomers();
    }

    notifySuccess({ message: t.messages.deletedSuccessfully.value });
  } catch (err: any) {
    notifyFail({ message: err.message || 'Failed to delete customer' });
  }
}

/** LIFECYCLE HOOKS */
onMounted(async () => await fetchCustomers());
</script>

<template>
  <q-page class="q-pt-lg">

    <!-- Header Section -->
    <PageHeader
        icon="domain"
        iconColor="primary"
        :title="t.page.title.value"
        :description="t.page.description.value"
        :button="canCreateOrg ? { label: t.page.addButton.value, icon: 'add', to: addOrganizationRoute, color: 'primary' } : undefined"
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
          <q-icon name="domain" size="sm" color="primary" class="q-mr-sm"/>
          <div class="text-subtitle1 text-weight-medium text-primary">{{ t.page.listTitle.value }}</div>
        </div>
      </div>
      <div class="col-auto">
        <ListHeaderMenu
          icon="domain"
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
          @refresh="fetchCustomers"
        />
      </div>
    </div>

    <!-- Loading Spinner -->
    <div v-if="loading" class="row justify-center q-my-lg">
      <q-spinner color="primary" size="3em" />
    </div>

    <!-- Customers Row List -->
    <div v-else class="row">
      <div
          v-for="(customer, index) in customersList"
          :key="customer.id || `customer-${index}`"
          class="col-12 q-mb-xs"
      >
        <DataRow
            :data="customer"
            :columns="visibleColumns"
            :actions="{ showEdit: canUpdateOrg, showView: canReadOrg, showDelete: canDeleteOrg }"
            @click="viewDetails"
            @dblclick="editCustomer"
            @edit="editCustomer"
            @view="viewDetails"
            @delete="confirmDelete"
        />
      </div>

      <!-- No Results -->
      <div v-if="customersList.length === 0" class="col-12">
        <ListCardEmpty
            :title="t.empty.title.value"
            :description="t.empty.description.value"
            icon="domain"
        />
      </div>
    </div>

    <!-- Pagination -->
    <ListPagination
      v-model="currentPage"
      :total-pages="totalPages"
      @change="handlePageChange"
    />

    <!-- Customer Details Drawer -->
    <CustomerDetailsDrawer
      v-model="showDetailsDrawer"
      :customer-id="selectedCustomerId"
      @edit="handleDrawerEdit"
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
