<script setup lang="ts">
defineOptions({
  name: 'GroupsListPage'
});

/** TYPE IMPORTS */
import type { ListHeaderMenuColumn } from '@components/headers';
import type { GroupResponse } from '@mapexos/schemas';
import type { FilterField, FilterValues, FilterAutocompleteOption } from '@components/drawers';
import type {
  GroupsListPageFilters,
  GroupsListPageColumnVisibility,
} from './interfaces';

/** VUE IMPORTS */
import { ref, computed, onMounted } from 'vue';
import { useRouter } from 'vue-router';

/** COMPONENTS */
import { PageHeader, ListHeaderMenu } from '@components/headers';
import { ListCardEmpty, DataRow } from '@components/cards';
import { ListPagination } from '@components/navigation';
import { GroupDetailsDrawer, AdvancedFiltersDrawer } from '@components/drawers';
import { AppTooltip } from '@components/tooltips';

/** SERVICES */
import { apis } from '@services/mapex';

/** COMPOSABLES */
import { useGroupsTranslations } from '@composables/i18n';
import { usePermissions } from '@composables/shared/usePermissions';

/** LOCAL IMPORTS */
import {
  GROUPS_LIST_PAGE_DEFAULTS,
  GROUPS_COLUMN_VISIBILITY_DEFAULTS,
  GROUPS_FILTER_DEFAULTS,
} from './constants';
import {
  fetchGroupsHandler,
  handlePageChangeHandler,
  handleItemsPerPageChangeHandler,
  handleColumnsUpdateHandler,
  editGroupHandler,
  confirmDeleteHandler,
  deleteGroupHandler,
} from './handlers';

/** COMPOSABLES & STORES */
const t = useGroupsTranslations();
const router = useRouter();
const { canCreate, canUpdate, canDelete, canRead } = usePermissions();
const canCreateGroup = canCreate('groups');
const canUpdateGroup = canUpdate('groups');
const canDeleteGroup = canDelete('groups');
const canReadGroup = canRead('groups');

/** STATE */
const groupsList = ref<GroupResponse[]>([]);
const loading = ref(false);
const lastUpdatedAt = ref<number | undefined>(undefined);
const error = ref<string | undefined>(undefined);
const itemsPerPage = ref(GROUPS_LIST_PAGE_DEFAULTS.ITEMS_PER_PAGE);
const currentPage = ref(GROUPS_LIST_PAGE_DEFAULTS.INITIAL_PAGE);
const totalPages = ref(1);
const totalItems = ref(0);
const filters = ref<GroupsListPageFilters>({ ...GROUPS_FILTER_DEFAULTS });
const columnVisibilityState = ref<GroupsListPageColumnVisibility>({ ...GROUPS_COLUMN_VISIBILITY_DEFAULTS });
const showDetailsDrawer = ref(false);
const selectedGroupId = ref<string | null>(null);

/** FILTER STATE */
const showFiltersDrawer = ref(false);
const quickSearchName = ref('');
const quickStatusEnabled = ref<boolean | null>(null);
const advancedFilterValues = ref<FilterValues>({
  includeChildren: null,
  memberId: null,
});
const hasPendingAdvancedFilters = ref(false);

/** COMPUTED */

/**
 * Status options for quick filter
 */
const statusOptions = computed(() => [
  { label: t.filters.allStatus.value, value: null },
  { label: t.filters.options.enabled.value, value: true },
  { label: t.filters.options.disabled.value, value: false },
]);

/**
 * Advanced filter fields configuration
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
    key: 'memberId',
    type: 'autocomplete',
    label: t.filters.filterByMember.value,
    icon: 'person_search',
    placeholder: t.filters.searchMemberPlaceholder.value,
    fetchOptions: fetchUserOptions,
  },
]);

/**
 * Check if any filter is active
 */
const hasActiveFilters = computed(() => {
  return !!(
    filters.value.name ||
    filters.value.enabled !== undefined ||
    filters.value.includeChildren !== undefined ||
    filters.value.memberId
  );
});

/**
 * Count of advanced filters only (for badge on filter icon)
 */
const advancedFiltersCount = computed(() => {
  let count = 0;
  if (filters.value.includeChildren !== undefined) count++;
  if (filters.value.memberId) count++;
  return count;
});

/**
 * Active filter chips for display
 */
const activeFilterChips = computed(() => {
  const chips: Array<{ key: string; label: string; value: string }> = [];

  if (filters.value.name) {
    chips.push({ key: 'name', label: t.filters.name.value, value: filters.value.name });
  }
  if (filters.value.enabled !== undefined) {
    chips.push({
      key: 'enabled',
      label: t.filters.status.value,
      value: filters.value.enabled ? t.filters.options.enabled.value : t.filters.options.disabled.value
    });
  }
  if (filters.value.includeChildren !== undefined) {
    chips.push({
      key: 'includeChildren',
      label: t.filters.includeChildren.value,
      value: filters.value.includeChildren ? t.filters.options.yes.value : t.filters.options.no.value
    });
  }
  if (filters.value.memberId && advancedFilterValues.value.memberIdLabel) {
    chips.push({ key: 'memberId', label: t.filters.member.value, value: advancedFilterValues.value.memberIdLabel as string });
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
 * Column visibility using ListHeaderMenuColumn format
 */
const menuColumns = computed(() => {
  const cols: ListHeaderMenuColumn[] = [
    { key: 'description', label: t.menuColumns.description.value, visible: columnVisibilityState.value.description },
    { key: 'membersCount', label: t.menuColumns.members.value, visible: columnVisibilityState.value.membersCount },
    { key: 'created', label: t.menuColumns.created.value, visible: columnVisibilityState.value.created },
  ];

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
    if (col.key === 'avatar' || col.key === 'name' || col.key === 'type' || col.key === 'status') {
      return true;
    }

    if (col.key === 'organizationName') {
      return filters.value.includeChildren === true && columnVisibilityState.value.organization;
    }

    if (col.key === 'description') return columnVisibilityState.value.description;
    if (col.key === 'membersCount') return columnVisibilityState.value.membersCount;
    if (col.key === 'created') return columnVisibilityState.value.created;

    return true;
  });
});

/** FUNCTIONS */

/**
 * Fetch groups from API with current filters and pagination
 * @returns {Promise<void>}
 */
async function fetchGroups(): Promise<void> {
  await fetchGroupsHandler(
    filters,
    currentPage,
    itemsPerPage,
    groupsList,
    totalPages,
    totalItems,
    loading,
    error,
  );
  lastUpdatedAt.value = Date.now();
}

/**
 * Fetch user options for autocomplete
 * @param {string} search - Search term
 * @returns {Promise<FilterAutocompleteOption[]>}
 */
async function fetchUserOptions(search: string): Promise<FilterAutocompleteOption[]> {
  const response = await apis.mapexOS.users.list({
    page: 1,
    perPage: 10,
    firstName: search,
    projection: 'id,firstName,lastName,email',
  });

  return response.items.map((user) => ({
    id: user.id || '',
    label: user.firstName
      ? `${user.firstName} ${user.lastName || ''}`.trim()
      : user.email || '',
    caption: user.email || '',
  }));
}

/**
 * Apply quick filters (search + status)
 */
function applyQuickFilters(): void {
  filters.value.name = quickSearchName.value || undefined;
  filters.value.enabled = quickStatusEnabled.value ?? undefined;
  currentPage.value = 1;

  // Update column visibility for includeChildren
  if (filters.value.includeChildren === true) {
    columnVisibilityState.value.organization = true;
  }

  void fetchGroups();
}

/**
 * Handle advanced filters apply from drawer
 * @param {FilterValues} values - Applied filter values
 */
function handleAdvancedFiltersApply(values: FilterValues): void {
  advancedFilterValues.value = values;
  filters.value.includeChildren = values.includeChildren ?? undefined;
  filters.value.memberId = values.memberId ?? undefined;
  currentPage.value = 1;

  // Reset pending state after apply
  hasPendingAdvancedFilters.value = false;

  // Update column visibility for includeChildren
  if (filters.value.includeChildren === true) {
    columnVisibilityState.value.organization = true;
  } else {
    columnVisibilityState.value.organization = false;
  }

  void fetchGroups();
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
    memberId: null,
  };
}

/**
 * Remove individual filter
 * @param {string} key - Filter key to remove
 */
function removeFilter(key: string): void {
  if (key === 'name') {
    filters.value.name = undefined;
    quickSearchName.value = '';
  } else if (key === 'enabled') {
    filters.value.enabled = undefined;
    quickStatusEnabled.value = null;
  } else if (key === 'includeChildren') {
    filters.value.includeChildren = undefined;
    advancedFilterValues.value.includeChildren = null;
    columnVisibilityState.value.organization = false;
  } else if (key === 'memberId') {
    filters.value.memberId = undefined;
    advancedFilterValues.value.memberId = null;
    advancedFilterValues.value.memberIdLabel = null;
  }

  currentPage.value = 1;
  void fetchGroups();
}

/**
 * Clear all filters
 */
function clearAllFilters(): void {
  // Reset filter state
  filters.value = { ...GROUPS_FILTER_DEFAULTS };

  // Reset quick filters
  quickSearchName.value = '';
  quickStatusEnabled.value = null;

  // Reset advanced filters
  advancedFilterValues.value = {
    includeChildren: null,
    memberId: null,
  };

  // Reset pending state
  hasPendingAdvancedFilters.value = false;

  // Reset column visibility
  columnVisibilityState.value.organization = false;

  currentPage.value = 1;
  void fetchGroups();
}

/**
 * Handle pagination navigation
 * @param {number} page - New page number
 */
function handlePageChange(page: number): void {
  handlePageChangeHandler(page, currentPage, () => void fetchGroups());
}

/**
 * Handle items per page change
 * @param {number} newValue - New items per page value
 */
function handleItemsPerPageChange(newValue: number): void {
  handleItemsPerPageChangeHandler(newValue, itemsPerPage, currentPage, () => void fetchGroups());
}

/**
 * Update menu columns when changed
 * @param {ListHeaderMenuColumn[]} columns - Updated column visibility states
 */
function handleColumnsUpdate(columns: ListHeaderMenuColumn[]): void {
  handleColumnsUpdateHandler(columns, columnVisibilityState);
}

/**
 * View group details in drawer
 * @param {any} group - Group to view
 */
function viewDetails(group: any): void {
  if (!canReadGroup.value) return;
  if (!group.id) {
    return;
  }
  selectedGroupId.value = group.id;
  showDetailsDrawer.value = true;
}

/**
 * Edit group - Navigate to edit page
 * @param {any} group - Group to edit
 */
function editGroup(group: any): void {
  if (!canUpdateGroup.value) return;
  editGroupHandler(group, t, router);
}

/**
 * Handle edit event from drawer
 * @param {string} groupId - ID of group to edit
 */
function handleDrawerEdit(groupId: string): void {
  if (!canUpdateGroup.value) return;
  void router.push(`/groups/edit/${groupId}`);
}

/**
 * Confirm delete group with dialog
 * @param {GroupResponse} group - Group to delete
 */
async function confirmDelete(group: GroupResponse): Promise<void> {
  if (!canDeleteGroup.value) return;
  await confirmDeleteHandler(group, t, deleteGroup);
}

/**
 * Delete group from API and update list
 * @param {GroupResponse} group - Group to delete
 */
async function deleteGroup(group: GroupResponse): Promise<void> {
  await deleteGroupHandler(group, groupsList, totalItems, currentPage, t, fetchGroups);
}

/** LIFECYCLE HOOKS */
onMounted(async () => await fetchGroups());
</script>

<template>
  <q-page class="q-pt-lg">

    <!-- Header Section -->
    <PageHeader
        icon="groups"
        iconColor="primary"
        :title="t.page.title.value"
        :description="t.page.description.value"
        :button="canCreateGroup ? { label: t.page.addButton.value, icon: 'add', color: 'primary', to: '/groups/add' } : undefined"
        :info="t.page.info.value"
    />

    <!-- Filters Section -->
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
          v-model="quickStatusEnabled"
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
          <q-icon name="groups" size="sm" color="primary" class="q-mr-sm"/>
          <div class="text-subtitle1 text-weight-medium text-primary">{{ t.page.listTitle }}</div>
        </div>
      </div>
      <div class="col-auto">
        <ListHeaderMenu
          icon="groups"
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
          @refresh="fetchGroups"
        />
      </div>
    </div>

    <!-- Loading Spinner -->
    <div v-if="loading" class="row justify-center q-my-lg">
      <q-spinner color="primary" size="3em" />
    </div>

    <!-- Groups Row List -->
    <div v-else class="row">
      <div
          v-for="(group, index) in groupsList"
          :key="group.id || `group-${index}`"
          class="col-12 q-mb-xs"
      >
        <DataRow
            :data="group"
            :columns="visibleColumns"
            :actions="{ showEdit: canUpdateGroup, showView: canReadGroup, showDelete: canDeleteGroup }"
            @click="viewDetails"
            @dblclick="editGroup"
            @edit="editGroup"
            @view="viewDetails"
            @delete="confirmDelete"
        />
      </div>

      <!-- No Results -->
      <div v-if="groupsList.length === 0" class="col-12">
        <ListCardEmpty
            :title="t.empty.title.value"
            :description="t.empty.description.value"
            icon="groups"
        />
      </div>
    </div>

    <!-- Pagination -->
    <ListPagination
      v-model="currentPage"
      :total-pages="totalPages"
      @change="handlePageChange"
    />

    <!-- Group Details Drawer -->
    <GroupDetailsDrawer
      v-model="showDetailsDrawer"
      :group-id="selectedGroupId"
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
