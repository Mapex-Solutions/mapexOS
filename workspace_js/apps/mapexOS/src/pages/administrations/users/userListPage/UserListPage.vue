<script setup lang="ts">
defineOptions({
  name: 'UserListPage'
});

/** TYPE IMPORTS */
import type { ListHeaderMenuColumn } from '@components/headers';
import type { UserResponse, UserQuery } from '@mapexos/schemas';
import type { FilterField, FilterValues } from '@components/drawers';
import type { PageTourStep } from '@composables/tour';
import type {
  UserListPageFilters,
  UserListPageColumnVisibility,
} from './interfaces';

/** VUE IMPORTS */
import { ref, computed, onMounted } from 'vue';
import { useRouter } from 'vue-router';

/** COMPONENTS */
import { PageHeader, ListHeaderMenu } from '@components/headers';
import { ListCardEmpty, DataRow } from '@components/cards';
import { UserDetailsDrawer, AdvancedFiltersDrawer } from '@components/drawers';
import { ListPagination } from '@components/navigation';
import { AppTooltip } from '@components/tooltips';

/** COMPOSABLES */
import { useUsersTranslations } from '@composables/i18n';
import { usePageTour } from '@composables/tour';
import { usePermissions } from '@composables/shared/usePermissions';
import { useLogger } from '@composables/useLogger';

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
  USERS_PROJECTION,
  USER_LIST_TOUR_STEPS,
  USER_LIST_TOUR_TRANSITION,
  TOUR_BUTTON_STEP,
  ROW_ACTIONS_STEP,
  DEMO_USER,
} from './constants';

/** COMPOSABLES & STORES */
const t = useUsersTranslations();
const router = useRouter();
const orgStore = useOrganizationStore();
const logger = useLogger('UserListPage');
const { canCreate, canUpdate, canDelete, canRead } = usePermissions();
const canCreateUser = canCreate('users');
const canUpdateUser = canUpdate('users');
const canDeleteUser = canDelete('users');
const canReadUser = canRead('users');

/** Whether this page has tour button enabled in PageHeader */
const hasTourButton = false;

/** Demo row state for tour */
const showDemoRow = ref(false);

/**
 * Open action menu on demo row for tour demonstration
 * Clicks the more_vert button to open the q-menu
 */
function openDemoRowActionMenu(): void {
  const demoRow = document.querySelector('[data-tour-demo-row="true"]');
  if (demoRow) {
    const menuButton = demoRow.querySelector('.data-row-actions-cell button');
    if (menuButton) {
      (menuButton as HTMLElement).click();
    }
  }
}

/**
 * Close any open action menus and hide demo row
 */
function cleanupDemoRow(): void {
  showDemoRow.value = false;
  // Close any open menus by clicking outside
  const overlay = document.querySelector('.q-menu__backdrop');
  if (overlay) {
    (overlay as HTMLElement).click();
  }
}

/**
 * Build tour steps with resolved translations and drawer callbacks
 * If hasTourButton is true, prepends a step showing the tour button
 * Pattern: [tourButton] → header → filters → advancedFiltersBtn → advancedFiltersOpen → results → rowActions → addNew
 *
 * @returns {PageTourStep[]} Tour steps with resolved text and callbacks
 */
function buildTourSteps(): PageTourStep[] {
  // Insert rowActions step after results
  const stepsWithRowActions = [...USER_LIST_TOUR_STEPS];
  const resultsIndex = stepsWithRowActions.findIndex(s => s.translationKey === 'results');
  if (resultsIndex !== -1) {
    stepsWithRowActions.splice(resultsIndex + 1, 0, ROW_ACTIONS_STEP);
  }

  // Build base steps from constants
  const allStepDefs = hasTourButton
    ? [TOUR_BUTTON_STEP, ...stepsWithRowActions]
    : stepsWithRowActions;

  return allStepDefs.map((step) => {
    const key = step.translationKey as keyof typeof t.tour;
    const translation = t.tour[key];
    const result: PageTourStep = {
      element: step.element,
      title: translation.title.value,
      description: translation.description.value,
    };
    if (step.side) result.side = step.side;
    if (step.align) result.align = step.align;

    // Advanced filters button: open drawer on Next click
    if (step.translationKey === 'advancedFiltersBtn') {
      result.onNextClick = (moveNext) => {
        showFiltersDrawer.value = true;
        setTimeout(moveNext, 400);
      };
    }

    // Advanced filters open: close drawer on Next click
    if (step.translationKey === 'advancedFiltersOpen') {
      result.onNextClick = (moveNext) => {
        showFiltersDrawer.value = false;
        setTimeout(moveNext, 300);
      };
    }

    // Results step: ensure drawer is closed, show demo row and prepare for next step
    if (step.translationKey === 'results') {
      result.onHighlightStarted = () => {
        showFiltersDrawer.value = false;
        showDemoRow.value = true;
      };
      result.onNextClick = (moveNext) => {
        // Open the action menu before moving to next step
        setTimeout(() => {
          openDemoRowActionMenu();
          setTimeout(moveNext, 300);
        }, 100);
      };
    }

    // Row actions step: target the open menu, cleanup on next
    if (step.translationKey === 'rowActions') {
      result.onNextClick = (moveNext) => {
        cleanupDemoRow();
        setTimeout(moveNext, 200);
      };
    }

    return result;
  });
}

/**
 * Compute transition config with adjusted step index
 * When tourButton step is prepended, all step indices shift by 1
 * Note: USER_LIST_TOUR_TRANSITION.triggerAtStep already accounts for rowActions step
 */
function getTransitionConfig() {
  if (!hasTourButton) return USER_LIST_TOUR_TRANSITION;

  return {
    ...USER_LIST_TOUR_TRANSITION,
    triggerAtStep: USER_LIST_TOUR_TRANSITION.triggerAtStep + 1,
  };
}

/** PAGE TOUR */
usePageTour({
  tourId: 'users-list',
  steps: buildTourSteps,
  transition: getTransitionConfig(),
});

/** STATE */
const usersList = ref<UserResponse[]>([]);
const loading = ref(false);
const lastUpdatedAt = ref<number | undefined>(undefined);
const error = ref<string | undefined>(undefined);
const showDetailsDrawer = ref(false);
const selectedUserId = ref<string | null>(null);
const itemsPerPage = ref(DEFAULT_ITEMS_PER_PAGE);
const currentPage = ref(1);
const totalPages = ref(1);
const totalItems = ref(0);
const filters = ref<UserListPageFilters>({ ...FILTER_DEFAULTS });
const columnVisibilityState = ref<UserListPageColumnVisibility>({ ...COLUMN_VISIBILITY_DEFAULTS });

/** FILTER STATE */
const showFiltersDrawer = ref(false);
const quickSearchName = ref('');
const quickStatusEnabled = ref<boolean | null>(null);
const advancedFilterValues = ref<FilterValues>({
  includeChildren: null,
  email: null,
  firstName: null,
  lastName: null,
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
    key: 'email',
    type: 'input',
    label: t.filters.filterByEmail.value,
    icon: 'email',
    placeholder: t.filters.email.value,
  },
  {
    key: 'firstName',
    type: 'input',
    label: t.filters.filterByFirstName.value,
    icon: 'person',
    placeholder: t.filters.firstName.value,
  },
  {
    key: 'lastName',
    type: 'input',
    label: t.filters.filterByLastName.value,
    icon: 'person_outline',
    placeholder: t.filters.lastName.value,
  },
]);

/**
 * Check if any filter is active
 */
const hasActiveFilters = computed(() => {
  return !!(
    filters.value.email ||
    filters.value.firstName ||
    filters.value.lastName ||
    filters.value.enabled !== undefined ||
    filters.value.includeChildren !== undefined
  );
});

/**
 * Count of advanced filters only (for badge on filter icon)
 */
const advancedFiltersCount = computed(() => {
  let count = 0;
  if (filters.value.includeChildren !== undefined) count++;
  if (filters.value.email) count++;
  if (filters.value.firstName) count++;
  if (filters.value.lastName) count++;
  return count;
});

/**
 * Active filter chips for display
 */
const activeFilterChips = computed(() => {
  const chips: Array<{ key: string; label: string; value: string }> = [];

  if (quickSearchName.value) {
    chips.push({ key: 'search', label: t.filters.email.value, value: quickSearchName.value });
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
  if (filters.value.email) {
    chips.push({ key: 'email', label: t.filters.email.value, value: filters.value.email });
  }
  if (filters.value.firstName) {
    chips.push({ key: 'firstName', label: t.filters.firstName.value, value: filters.value.firstName });
  }
  if (filters.value.lastName) {
    chips.push({ key: 'lastName', label: t.filters.lastName.value, value: filters.value.lastName });
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
    { key: 'email', label: t.menuColumns.email.value, visible: columnVisibilityState.value.email },
    { key: 'jobTitle', label: t.menuColumns.jobTitle.value, visible: columnVisibilityState.value.jobTitle },
    { key: 'groups', label: t.menuColumns.groups.value, visible: columnVisibilityState.value.groups },
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
    // Always show avatar, name, and status
    if (col.key === 'avatar' || col.key === 'name' || col.key === 'status') {
      return true;
    }

    // Organization column only visible when includeChildren filter is active
    if (col.key === 'organizationName') {
      return filters.value.includeChildren === true && columnVisibilityState.value.organization;
    }

    // Filter based on columnVisibility
    if (col.key === 'email') return columnVisibilityState.value.email;
    if (col.key === 'jobTitle') return columnVisibilityState.value.jobTitle;
    if (col.key === 'groupsCount') return columnVisibilityState.value.groups;

    return true;
  });
});

/** FUNCTIONS */

/**
 * Fetch users from API with current filters and pagination
 * @returns {Promise<void>}
 */
async function fetchUsers(): Promise<void> {
  if (!apis.mapexOS?.users) {
    error.value = 'Users API not initialized';
    return;
  }

  try {
    loading.value = true;
    error.value = undefined;

    // Build query parameters from filters and pagination state
    const queryParams: Record<string, any> = {
      page: currentPage.value,
      perPage: itemsPerPage.value,
      projection: USERS_PROJECTION,
    };

    // Add active filters conditionally (only if they have values)
    if (filters.value.email) {
      queryParams.email = filters.value.email;
    }
    if (filters.value.firstName) {
      queryParams.firstName = filters.value.firstName;
    }
    if (filters.value.lastName) {
      queryParams.lastName = filters.value.lastName;
    }
    if (typeof filters.value.enabled === 'boolean') {
      queryParams.enabled = filters.value.enabled;
    }
    if (typeof filters.value.includeChildren === 'boolean') {
      queryParams.includeChildren = filters.value.includeChildren;
    }

    const response = await apis.mapexOS.users.list(queryParams as UserQuery);

    // Enrich users with organization name
    const enrichedUsers = (response.items || []).map((user: any) => {
      const organization = orgStore.flatList.find((org: any) => org.id === user.orgId);
      return {
        ...user,
        organizationName: organization?.name || 'Unknown',
      };
    });

    usersList.value = enrichedUsers;

    // Update pagination state from response
    if (response.pagination) {
      totalPages.value = response.pagination.totalPages || 1;
      totalItems.value = response.pagination.totalItems || 0;
    }
  } catch (err: any) {
    logger.error('Error fetching users:', err);
    const errorMsg = err.message || 'Failed to fetch users';
    error.value = errorMsg;
    notifyFail({ message: errorMsg });
  } finally {
    loading.value = false;
    lastUpdatedAt.value = Date.now();
  }
}

/**
 * Apply quick filters (search + status)
 */
function applyQuickFilters(): void {
  // Quick search maps to email (most common search field)
  filters.value.email = quickSearchName.value || undefined;
  filters.value.enabled = quickStatusEnabled.value ?? undefined;
  currentPage.value = 1;

  // Update column visibility for includeChildren
  if (filters.value.includeChildren === true) {
    columnVisibilityState.value.organization = true;
  }

  void fetchUsers();
}

/**
 * Handle advanced filters apply from drawer
 * @param {FilterValues} values - Applied filter values
 */
function handleAdvancedFiltersApply(values: FilterValues): void {
  advancedFilterValues.value = values;
  filters.value.includeChildren = values.includeChildren ?? undefined;
  filters.value.email = values.email ?? quickSearchName.value ?? undefined;
  filters.value.firstName = values.firstName ?? undefined;
  filters.value.lastName = values.lastName ?? undefined;
  currentPage.value = 1;

  // Reset pending state after apply
  hasPendingAdvancedFilters.value = false;

  // Update column visibility for includeChildren
  if (filters.value.includeChildren === true) {
    columnVisibilityState.value.organization = true;
  } else {
    columnVisibilityState.value.organization = false;
  }

  void fetchUsers();
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
    email: null,
    firstName: null,
    lastName: null,
  };
}

/**
 * Remove individual filter
 * @param {string} key - Filter key to remove
 */
function removeFilter(key: string): void {
  if (key === 'search') {
    quickSearchName.value = '';
    filters.value.email = advancedFilterValues.value.email ?? undefined;
  } else if (key === 'enabled') {
    filters.value.enabled = undefined;
    quickStatusEnabled.value = null;
  } else if (key === 'includeChildren') {
    filters.value.includeChildren = undefined;
    advancedFilterValues.value.includeChildren = null;
    columnVisibilityState.value.organization = false;
  } else if (key === 'email') {
    filters.value.email = quickSearchName.value || undefined;
    advancedFilterValues.value.email = null;
  } else if (key === 'firstName') {
    filters.value.firstName = undefined;
    advancedFilterValues.value.firstName = null;
  } else if (key === 'lastName') {
    filters.value.lastName = undefined;
    advancedFilterValues.value.lastName = null;
  }

  currentPage.value = 1;
  void fetchUsers();
}

/**
 * Clear all filters
 */
function clearAllFilters(): void {
  // Reset filter state
  filters.value = { ...FILTER_DEFAULTS };

  // Reset quick filters
  quickSearchName.value = '';
  quickStatusEnabled.value = null;

  // Reset advanced filters
  advancedFilterValues.value = {
    includeChildren: null,
    email: null,
    firstName: null,
    lastName: null,
  };

  // Reset pending state
  hasPendingAdvancedFilters.value = false;

  // Reset column visibility
  columnVisibilityState.value.organization = false;

  currentPage.value = 1;
  void fetchUsers();
}

/**
 * Handle pagination navigation
 * @param {number} page - New page number
 * @returns {void}
 */
function handlePageChange(page: number): void {
  currentPage.value = page;
  void fetchUsers();
}

/**
 * Handle items per page change
 * @param {number} newValue - New items per page value
 * @returns {void}
 */
function handleItemsPerPageChange(newValue: number): void {
  itemsPerPage.value = newValue;
  currentPage.value = 1; // Reset to first page
  void fetchUsers();
}

/**
 * Update menu columns when changed
 * @param {ListHeaderMenuColumn[]} columns - Updated column visibility states
 * @returns {void}
 */
function handleColumnsUpdate(columns: ListHeaderMenuColumn[]): void {
  columns.forEach(col => {
    if (col.key === 'organization') columnVisibilityState.value.organization = col.visible;
    if (col.key === 'email') columnVisibilityState.value.email = col.visible;
    if (col.key === 'jobTitle') columnVisibilityState.value.jobTitle = col.visible;
    if (col.key === 'groups') columnVisibilityState.value.groups = col.visible;
  });
}

/**
 * View user details in drawer (quick view)
 *
 * @param {UserResponse} user - User to view
 * @returns {void}
 */
function viewDetailsDrawer(user: UserResponse): void {
  if (!canReadUser.value) return;
  if (!user.id) {
    notifyFail({ message: t.errors.idMissing.value });
    return;
  }
  selectedUserId.value = user.id;
  showDetailsDrawer.value = true;
}

/**
 * Edit user - navigates to edit form
 *
 * @param {UserResponse} user - User to edit
 * @returns {void}
 */
function editUser(user: UserResponse): void {
  if (!canUpdateUser.value) return;
  if (!user.id) {
    notifyFail({ message: t.errors.idMissing.value });
    return;
  }
  logger.debug('Navigate to user edit:', user.id);
  void router.push(`/users/edit/${user.id}`);
}

/**
 * Handle edit event from drawer
 * @param {string} userId - ID of user to edit
 * @returns {void}
 */
function handleDrawerEdit(userId: string): void {
  if (!canUpdateUser.value) return;
  logger.debug('Edit user from drawer:', userId);
  void router.push(`/users/edit/${userId}`);
}

/**
 * Confirm delete operation with dialog
 * @param {UserResponse} user - User to delete
 * @returns {Promise<void>}
 */
async function confirmDelete(user: UserResponse): Promise<void> {
  if (!canDeleteUser.value) return;
  const userName = `${user.firstName || ''} ${user.lastName || ''}`.trim() || user.email || 'this user';
  const confirmed = await dialogDelete({
    title: t.dialog.deleteTitle.value,
    message: t.messages.confirmDelete(userName),
  });

  if (confirmed) {
    await deleteUser(user);
  }
}

/**
 * Delete user from API and update list
 * @param {UserResponse} user - User to delete
 * @returns {Promise<void>}
 */
async function deleteUser(user: UserResponse): Promise<void> {
  if (!apis.mapexOS?.users) {
    notifyFail({ message: t.errors.apiNotInitialized.value });
    return;
  }

  if (!user.id) {
    notifyFail({ message: t.errors.idMissing.value });
    return;
  }

  try {
    await apis.mapexOS.users.delete({ userId: user.id });

    // Update total items count after successful deletion
    totalItems.value = Math.max(0, totalItems.value - 1);

    // Remove from local list
    const index = usersList.value.findIndex(r => r.id === user.id);
    if (index !== -1) {
      usersList.value.splice(index, 1);
    }

    // If current page is now empty and not the first page, go to previous page
    if (usersList.value.length === 0 && currentPage.value > 1) {
      currentPage.value -= 1;
      await fetchUsers();
    }

    notifySuccess({ message: t.messages.deletedSuccessfully.value });
  } catch (err: any) {
    notifyFail({ message: err.message || 'Failed to delete user' });
  }
}

/** LIFECYCLE HOOKS */
onMounted(async () => await fetchUsers());

</script>

<template>
  <q-page class="q-pt-lg">

    <!-- Header Section -->
    <div id="page-header-section">
      <PageHeader
          icon="person"
          iconColor="primary"
          :title="t.page.title.value"
          :description="t.page.description.value"
          :button="canCreateUser ? { label: t.page.addButton.value, icon: 'add', to: '/users/add', color: 'primary', id: 'add-user-btn' } : undefined"
      />
    </div>

    <!-- Filters Section -->
    <div id="filter-section">
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
          id="advanced-filters-btn"
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

    </div>

    <!-- Results Section -->
    <div id="results-section" class="row items-center q-pt-xl q-mb-md">
      <div class="col">
        <div class="row items-center">
          <q-icon name="person" size="sm" color="primary" class="q-mr-sm"/>
          <div class="text-subtitle1 text-weight-medium text-primary">{{ t.page.listTitle }}</div>
        </div>
      </div>
      <div class="col-auto">
        <ListHeaderMenu
          icon="person"
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
          @refresh="fetchUsers"
        />
      </div>
    </div>

    <!-- Loading Spinner -->
    <div v-if="loading" class="row justify-center q-my-lg">
      <q-spinner color="primary" size="3em" />
    </div>

    <!-- Users Row List -->
    <div v-else class="row">
      <!-- Demo Row for Tour -->
      <div
        v-if="showDemoRow"
        class="col-12 q-mb-xs demo-row-highlight"
        data-tour-demo-row="true"
      >
        <DataRow
          :data="DEMO_USER"
          :columns="visibleColumns"
          @click="() => {}"
          @dblclick="() => {}"
          @edit="() => {}"
          @view="() => {}"
          @delete="() => {}"
        />
      </div>

      <!-- Real Users -->
      <div
          v-for="(user, index) in usersList"
          :key="user.id || `user-${index}`"
          class="col-12 q-mb-xs"
      >
        <DataRow
            :data="user"
            :columns="visibleColumns"
            :actions="{ showEdit: canUpdateUser, showView: canReadUser, showDelete: canDeleteUser }"
            @click="viewDetailsDrawer"
            @dblclick="editUser"
            @edit="editUser"
            @view="viewDetailsDrawer"
            @delete="confirmDelete"
        />
      </div>

      <!-- No Results -->
      <div v-if="usersList.length === 0" class="col-12">
        <ListCardEmpty
            :title="t.empty.title.value"
            :description="t.empty.description.value"
            icon="person"
        />
      </div>
    </div>

    <!-- Pagination -->
    <ListPagination
      v-model="currentPage"
      :total-pages="totalPages"
      @change="handlePageChange"
    />

    <!-- User Details Drawer -->
    <UserDetailsDrawer
      v-model="showDetailsDrawer"
      :user-id="selectedUserId"
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

// Demo row highlight for tour
.demo-row-highlight {
  position: relative;

  &::before {
    content: 'DEMO';
    position: absolute;
    top: 8px;
    left: 8px;
    z-index: 10;
    background: $primary;
    color: white;
    font-size: 10px;
    font-weight: 600;
    padding: 2px 6px;
    border-radius: var(--mapex-radius-xs);
    letter-spacing: 0.5px;
  }

  :deep(.data-row-card) {
    border: 2px solid $primary;
    background: rgba($primary, 0.02);
  }
}
</style>
