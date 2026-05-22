<script setup lang="ts">
defineOptions({
  name: 'TriggerListPage'
});

/** TYPE IMPORTS (ALL types first, grouped) */
import type { ListHeaderMenuColumn } from '@components/headers';
import type { DataRowColumn } from '@components/cards';
import type { TriggerResponse } from '@mapexos/schemas';
import type { TriggerListPageFilters } from './handlers';

/** VUE IMPORTS */
import { ref, computed, onMounted } from 'vue';
import { useRouter } from 'vue-router';

/** COMPONENTS */
import { AdvancedFiltersDrawer } from '@components/drawers/common';
import { PageHeader, ListHeaderMenu } from '@components/headers';
import { ListCardEmpty, DataRow } from '@components/cards';
import { ListPagination } from '@components/navigation';
import { TriggerDetailDrawer } from '@components/drawers/triggers';
import { AppTooltip } from '@components/tooltips';

/** UTILS */
import { notifySuccess, dialogDelete } from '@utils/alert';

/** SERVICES */
import { apis } from '@services/mapex';

/** COMPOSABLES */
import { useTriggersTranslations } from '@composables/i18n/pages/automations/triggers/useTriggersTranslations';
import { useLogger } from '@composables/useLogger';
import { usePermissions } from '@composables/shared/usePermissions';

/** LOCAL IMPORTS (constants and handlers ONLY - NO types here!) */
import { fetchTriggersHandler, handlePageChangeHandler } from './handlers';

/** CONSTANTS */
const MAX_VISIBLE_CHIPS = 2;

/** COMPOSABLES & STORES */
const router = useRouter();
const translations = useTriggersTranslations();
const logger = useLogger('TriggerListPage');
const { canCreate, canUpdate, canDelete, canRead } = usePermissions();
const canCreateTrigger = canCreate('triggers');
const canUpdateTrigger = canUpdate('triggers');
const canDeleteTrigger = canDelete('triggers');
const canReadTrigger = canRead('triggers');

/** STATE */
const loading = ref(false);
const lastUpdatedAt = ref<number | undefined>(undefined);
const triggersList = ref<TriggerResponse[]>([]);

// Quick filters (inline)
const quickSearch = ref('');
const quickStatus = ref<boolean | null>(null);

// Advanced filters (drawer)
const showAdvancedFilters = ref(false);
const advancedFilterValues = ref<Record<string, any>>({
  includeChildren: null,
  category: null,
  triggerType: null,
});
const hasPendingAdvancedFilters = ref(false);

// Applied filters (sent to API)
const filters = ref<TriggerListPageFilters>({
  name: undefined,
  status: undefined,
  includeChildren: undefined,
  category: undefined,
  triggerType: undefined,
});

// Pagination
const itemsPerPage = ref(15);
const currentPage = ref(1);
const totalPages = ref(0);
const totalItems = ref(0);

// Column visibility using ListHeaderMenuColumn format - initialized with translations
const menuColumns = ref<ListHeaderMenuColumn[]>(translations.menuColumns.value);

// Detail drawer state
const showDetailDrawer = ref(false);
const selectedTriggerId = ref<string | null>(null);

/** COMPUTED */

/**
 * Map menuColumns visibility to column keys
 * Used to hide/show columns based on user preferences
 */
const columnVisibility = computed(() => ({
  triggerType: menuColumns.value.find(col => col.key === 'triggerType')?.visible ?? true,
  category: menuColumns.value.find(col => col.key === 'category')?.visible ?? true,
}));

/**
 * Complete columns configuration with icon, formatting, and visibility
 */
const triggerColumns = computed((): DataRowColumn[] => {
  return [
    {
      key: 'icon',
      label: '',
      type: 'avatar',
      visible: 'always',
      width: 56,
      icon: () => 'flash_on',
      color: (value: unknown, row: TriggerResponse) => row.enabled ? 'primary' : 'grey-5',
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
      key: 'triggerType',
      label: translations.columns.value.find((col: DataRowColumn) => col.key === 'triggerType')?.label || 'Type',
      type: 'chip',
      visible: 'laptop',
      width: 120,
      format: (value: unknown) => typeof value === 'string' ? value.toUpperCase() : 'UNKNOWN',
      color: 'blue-6',
    },
    {
      key: 'category',
      label: translations.columns.value.find((col: DataRowColumn) => col.key === 'category')?.label || 'Category',
      type: 'chip',
      visible: 'laptop',
      width: 150,
      format: (value: unknown) => typeof value === 'string' ? value.toUpperCase() : 'UNKNOWN',
      color: 'purple-6',
    },
    {
      key: 'enabled',
      label: translations.columns.value.find((col: DataRowColumn) => col.key === 'status')?.label || 'Status',
      type: 'badge',
      visible: 'always',
      width: 100,
      format: (value: unknown) => value ? translations.filters.options.active.value : translations.filters.options.inactive.value,
      color: (value: unknown) => value ? 'green-6' : 'red-6',
    },
  ];
});

/**
 * Filter columns based on user visibility preferences
 * Always show icon, name, and status - hide others based on menuColumns
 */
const visibleColumns = computed((): DataRowColumn[] => {
  return triggerColumns.value.filter((col: DataRowColumn) => {
    // Always show avatar, name, and status
    if (col.key === 'icon' || col.key === 'name' || col.key === 'status') {
      return true;
    }

    // Filter based on columnVisibility
    if (col.key === 'triggerType') return columnVisibility.value.triggerType;
    if (col.key === 'category') return columnVisibility.value.category;

    return true;
  });
});

/**
 * Build active filter chips from current filter state
 */
const activeFilterChips = computed(() => {
  const chips: Array<{ key: string; label: string; value: string }> = [];

  // Quick filters
  if (quickSearch.value) {
    chips.push({ key: 'name', label: translations.filters.label.value, value: quickSearch.value });
  }
  if (quickStatus.value !== null) {
    chips.push({
      key: 'status',
      label: translations.filters.allStatus.value,
      value: quickStatus.value ? translations.filters.options.active.value : translations.filters.options.inactive.value,
    });
  }

  // Advanced filters
  if (filters.value.includeChildren !== undefined && filters.value.includeChildren !== null) {
    chips.push({
      key: 'includeChildren',
      label: translations.filters.includeChildren.value,
      value: filters.value.includeChildren ? translations.filters.options.yes.value : translations.filters.options.no.value,
    });
  }
  if (filters.value.category) {
    const categoryValue = filters.value.category === 'technical'
      ? translations.filters.options.technical.value
      : translations.filters.options.communication.value;
    chips.push({ key: 'category', label: translations.filters.category.value, value: categoryValue });
  }
  if (filters.value.triggerType) {
    chips.push({ key: 'triggerType', label: translations.filters.triggerType.value, value: filters.value.triggerType.toUpperCase() });
  }

  return chips;
});

/**
 * Count of advanced filters currently applied
 */
const advancedFiltersCount = computed(() => {
  let count = 0;
  if (filters.value.includeChildren !== undefined && filters.value.includeChildren !== null) count++;
  if (filters.value.category) count++;
  if (filters.value.triggerType) count++;
  return count;
});

/**
 * Visible filter chips (limited by MAX_VISIBLE_CHIPS)
 */
const visibleFilterChips = computed(() => {
  return activeFilterChips.value.slice(0, MAX_VISIBLE_CHIPS);
});

/**
 * Hidden filter chips (shown in tooltip)
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

/** FUNCTIONS */

/**
 * Apply quick filters to the main filters and fetch data
 */
function applyQuickFilters(): void {
  filters.value.name = quickSearch.value || undefined;
  filters.value.status = quickStatus.value ?? undefined;
  currentPage.value = 1;
  void fetchTriggers();
}

/**
 * Handle advanced filters apply from drawer
 * @param {Record<string, any>} appliedFilters - Filters from drawer
 */
function handleAdvancedFiltersApply(appliedFilters: Record<string, any>): void {
  // Update advanced filter values
  advancedFilterValues.value = { ...appliedFilters };

  // Apply to main filters
  filters.value.includeChildren = appliedFilters.includeChildren ?? undefined;
  filters.value.category = appliedFilters.category ?? undefined;
  filters.value.triggerType = appliedFilters.triggerType ?? undefined;

  // Reset pending state
  hasPendingAdvancedFilters.value = false;

  // Reset pagination and fetch
  currentPage.value = 1;
  void fetchTriggers();
}

/**
 * Handle pending change event from advanced filters drawer
 * @param {boolean} hasPending - Whether there are pending changes
 */
function handlePendingChange(hasPending: boolean): void {
  hasPendingAdvancedFilters.value = hasPending;
}

/**
 * Remove a specific filter by key
 * @param {string} key - Filter key to remove
 */
function removeFilter(key: string): void {
  switch (key) {
    case 'name':
      quickSearch.value = '';
      break;
    case 'status':
      quickStatus.value = null;
      break;
    case 'includeChildren':
      advancedFilterValues.value.includeChildren = null;
      filters.value.includeChildren = undefined;
      break;
    case 'category':
      advancedFilterValues.value.category = null;
      filters.value.category = undefined;
      break;
    case 'triggerType':
      advancedFilterValues.value.triggerType = null;
      filters.value.triggerType = undefined;
      break;
  }
  currentPage.value = 1;
  void fetchTriggers();
}

/**
 * Clear all filters
 */
function clearAllFilters(): void {
  // Clear quick filters
  quickSearch.value = '';
  quickStatus.value = null;

  // Clear advanced filters
  advancedFilterValues.value = {
    includeChildren: null,
    category: null,
    triggerType: null,
  };
  filters.value = {
    name: undefined,
    status: undefined,
    includeChildren: undefined,
    category: undefined,
    triggerType: undefined,
  };

  hasPendingAdvancedFilters.value = false;
  currentPage.value = 1;
  void fetchTriggers();
}

/**
 * Fetch triggers from API with current filters and pagination
 * @returns {Promise<void>}
 */
async function fetchTriggers(): Promise<void> {
  loading.value = true;

  try {
    const result = await fetchTriggersHandler(
      filters.value,
      currentPage.value,
      itemsPerPage.value
    );

    triggersList.value = result.triggers;
    totalPages.value = result.totalPages;
    totalItems.value = result.totalItems;
  } finally {
    loading.value = false;
    lastUpdatedAt.value = Date.now();
  }
}

/**
 * Handle page change event from pagination component
 * Updates current page and refetches data
 *
 * @param {number} page - New page number (1-indexed)
 */
function handlePageChange(page: number): void {
  handlePageChangeHandler(page, currentPage, () => {
    void fetchTriggers();
  });
}

/**
 * Update menu columns when changed
 * @param {ListHeaderMenuColumn[]} columns - Updated columns
 */
function handleColumnsUpdate(columns: ListHeaderMenuColumn[]): void {
  menuColumns.value = columns;
}

/**
 * Handle single click - open detail drawer
 * @param {TriggerResponse} item - Clicked trigger
 */
function cardClick(item: TriggerResponse): void {
  if (!canReadTrigger.value) return;
  if (!item.id) return;
  selectedTriggerId.value = item.id;
  showDetailDrawer.value = true;
}

/**
 * Handle double click or edit action - navigate to edit page
 * @param {TriggerResponse} trigger - Trigger to edit
 */
function editTrigger(trigger: TriggerResponse): void {
  if (!canUpdateTrigger.value) return;
  if (!trigger.id) return;
  void router.push(`/triggers/edit/${trigger.id}`);
}

/**
 * Handle view details action from menu - open drawer
 * @param {TriggerResponse} trigger - Trigger to view
 */
function viewDetails(trigger: TriggerResponse): void {
  if (!canReadTrigger.value) return;
  if (!trigger.id) return;
  selectedTriggerId.value = trigger.id;
  showDetailDrawer.value = true;
}

/**
 * Handle edit from drawer
 * @param {string} triggerId - Trigger ID to edit
 */
function handleDrawerEdit(triggerId: string): void {
  if (!canUpdateTrigger.value) return;
  void router.push(`/triggers/edit/${triggerId}`);
}

/**
 * Confirm deletion of trigger
 * @param {TriggerResponse} trigger - Trigger to delete
 */
async function confirmDelete(trigger: TriggerResponse): Promise<void> {
  if (!canDeleteTrigger.value) return;
  const confirmed = await dialogDelete({
    title: translations.dialog.confirmDelete.title.value,
    message: translations.dialog.confirmDelete.message(trigger.name || 'this trigger'),
  });

  if (confirmed) {
    await deleteTrigger(trigger);
  }
}

/**
 * Delete trigger via API
 * @param {TriggerResponse} trigger - Trigger to delete
 */
async function deleteTrigger(trigger: TriggerResponse): Promise<void> {
  if (!trigger.id) return;

  try {
    await apis.triggers.trigger.delete({ triggerId: trigger.id.toString() });
    notifySuccess({ message: translations.notifications.deleteSuccess.value });

    // Refetch data after deletion
    void fetchTriggers();
  } catch (error: unknown) {
    logger.error('Error deleting trigger:', error);
  }
}

/** LIFECYCLE HOOKS */
onMounted(() => {
  void fetchTriggers();
});

</script>

<template>
  <q-page class="q-pt-lg">

    <!-- Header Section -->
    <PageHeader
        icon="flash_on"
        iconColor="primary"
        :title="translations.page.title.value"
        :description="translations.page.description.value"
        :button="canCreateTrigger ? { label: translations.page.button.add.value, icon: 'add', to: '/triggers/add', color: 'primary' } : undefined"
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
          :options="[
            { label: translations.filters.allStatus.value, value: null },
            { label: translations.filters.options.active.value, value: true },
            { label: translations.filters.options.inactive.value, value: false },
          ]"
          :label="translations.filters.allStatus.value"
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
      <!-- Visible Chips -->
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
          <q-icon name="flash_on" size="sm" color="primary" class="q-mr-sm"/>
          <div class="text-subtitle1 text-weight-medium text-primary">{{ translations.page.listTitle.value }}</div>
        </div>
      </div>
      <div class="col-auto">
        <ListHeaderMenu
          :items-count="totalItems"
          :item-label="translations.menuLabels.singular.value"
          :item-label-plural="translations.menuLabels.plural.value"
          icon="flash_on"
          :items-per-page="itemsPerPage"
          :columns="menuColumns"
          :refreshing="loading"
          :last-updated-at="lastUpdatedAt"
          @update:items-per-page="itemsPerPage = $event"
          @update:columns="handleColumnsUpdate"
          @refresh="fetchTriggers"
        />
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="flex flex-center q-pa-xl">
      <q-spinner color="primary" size="50px" />
    </div>

    <!-- Triggers Row List -->
    <div v-else-if="triggersList.length > 0" class="row">
      <div
          v-for="trigger in triggersList"
          :key="trigger.id?.toString() || `trigger-${Math.random()}`"
          class="col-12 q-mb-xs"
      >
        <DataRow
            :data="trigger"
            :columns="visibleColumns"
            :actions="{ showEdit: canUpdateTrigger, showView: canReadTrigger, showDelete: canDeleteTrigger }"
            @click="cardClick"
            @dblclick="editTrigger"
            @edit="editTrigger"
            @view="viewDetails"
            @delete="confirmDelete"
        />
      </div>
    </div>

    <!-- No Results -->
    <div v-else class="row q-col-gutter-lg">
      <ListCardEmpty
          :title="translations.empty.title.value"
          :description="translations.empty.description.value"
          icon="flash_on"
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
      :fields="translations.filterItems.value"
      :values="advancedFilterValues"
      @apply="handleAdvancedFiltersApply"
      @pending-change="handlePendingChange"
    />

    <!-- Trigger Detail Drawer -->
    <TriggerDetailDrawer
      v-model="showDetailDrawer"
      :trigger-id="selectedTriggerId"
      @edit="handleDrawerEdit"
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
