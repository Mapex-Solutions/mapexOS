<script setup lang="ts">
defineOptions({
  name: 'AccessAuditPage'
});

/** TYPE IMPORTS */
import type { ListHeaderMenuColumn } from '@components/headers';
import type { FilterField } from '@components/drawers';
import type {
  AccessAuditPageFilters,
  AccessAuditPageColumnVisibility,
  EnrichedMembership,
} from './interfaces';

/** VUE IMPORTS */
import { ref, computed, onMounted } from 'vue';
import { useRouter } from 'vue-router';

/** COMPONENTS */
import { PageHeader, ListHeaderMenu } from '@components/headers';
import { ListCardEmpty, DataRow } from '@components/cards';
import { AdvancedFiltersDrawer } from '@components/drawers';
import { ListPagination } from '@components/navigation';
import { AppTooltip } from '@components/tooltips';

/** COMPOSABLES */
import { useAccessAuditTranslations } from '@composables/i18n';
import { useLogger } from '@composables/useLogger';

/** UTILS */
import { notifyFail } from '@utils/alert';
import { cleanQueryParams } from '@utils/query';

/** SERVICES */
import { apis } from '@services/mapex';

/** STORES */
import { useOrganizationStore } from '@stores/organization';

/** LOCAL IMPORTS */
import {
  DEFAULT_ITEMS_PER_PAGE,
  COLUMN_VISIBILITY_DEFAULTS,
  FILTER_DEFAULTS,
  MEMBERSHIPS_PROJECTION,
} from './constants';

/** COMPOSABLES & STORES */
const t = useAccessAuditTranslations();
const router = useRouter();
const orgStore = useOrganizationStore();
const logger = useLogger('AccessAuditPage');

/** STATE - Caches for name resolution */
const rolesCache = ref<Map<string, string>>(new Map());
const usersCache = ref<Map<string, string>>(new Map());
const groupsCache = ref<Map<string, string>>(new Map());

/** STATE */
const membershipsList = ref<EnrichedMembership[]>([]);
const loading = ref(false);
const lastUpdatedAt = ref<number | undefined>(undefined);
const error = ref<string | undefined>(undefined);
const showFiltersDrawer = ref(false);
const itemsPerPage = ref(DEFAULT_ITEMS_PER_PAGE);
const currentPage = ref(1);
const totalPages = ref(1);
const totalItems = ref(0);
const filters = ref<AccessAuditPageFilters>({ ...FILTER_DEFAULTS });
const columnVisibilityState = ref<AccessAuditPageColumnVisibility>({ ...COLUMN_VISIBILITY_DEFAULTS });

// Quick filter state (inline)
const quickSearch = ref('');
const quickAssigneeType = ref<string | null>(null);

// Advanced filter state (drawer)
const advancedFilterValues = ref<Record<string, any>>({
  scope: null,
  includeChildren: null,
  enabled: null,
});

// Pending changes state
const hasPendingAdvancedFilters = ref(false);

/** COMPUTED */

/**
 * Assignee type options for quick filter select
 */
const assigneeTypeOptions = computed(() => [
  { label: t.filters.allStatus.value, value: null },
  { label: t.filters.options.user.value, value: 'user' },
  { label: t.filters.options.group.value, value: 'group' },
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
    key: 'scope',
    label: t.filters.scope.value,
    type: 'toggle',
    icon: 'domain',
    options: [
      { label: t.filters.allStatus.value, value: null },
      { label: t.filters.options.local.value, value: 'local' },
      { label: t.filters.options.recursive.value, value: 'recursive' },
    ],
  },
  {
    key: 'enabled',
    label: t.filters.status.value,
    type: 'toggle',
    icon: 'toggle_on',
    options: [
      { label: t.filters.allStatus.value, value: null },
      { label: t.filters.options.enabled.value, value: true },
      { label: t.filters.options.disabled.value, value: false },
    ],
  },
]);

/**
 * Check if any filters are active (quick or advanced)
 */
const hasActiveFilters = computed(() => {
  return !!(
    quickSearch.value ||
    quickAssigneeType.value !== null ||
    advancedFilterValues.value.scope !== null ||
    advancedFilterValues.value.includeChildren !== null ||
    advancedFilterValues.value.enabled !== null
  );
});

/**
 * Count of active advanced filters (for badge)
 */
const advancedFiltersCount = computed(() => {
  let count = 0;
  if (advancedFilterValues.value.scope !== null) count++;
  if (advancedFilterValues.value.includeChildren !== null) count++;
  if (advancedFilterValues.value.enabled !== null) count++;
  return count;
});

/**
 * Active filter chips for visual feedback
 */
const activeFilterChips = computed(() => {
  const chips: Array<{ key: string; label: string; value: string }> = [];

  if (quickSearch.value) {
    chips.push({
      key: 'search',
      label: t.filters.assigneeType.value,
      value: quickSearch.value,
    });
  }

  if (quickAssigneeType.value !== null) {
    chips.push({
      key: 'assigneeType',
      label: t.filters.assigneeType.value,
      value: quickAssigneeType.value === 'user'
        ? t.filters.options.user.value
        : t.filters.options.group.value,
    });
  }

  if (advancedFilterValues.value.scope !== null) {
    chips.push({
      key: 'scope',
      label: t.filters.scope.value,
      value: advancedFilterValues.value.scope === 'local'
        ? t.filters.options.local.value
        : t.filters.options.recursive.value,
    });
  }

  if (advancedFilterValues.value.includeChildren !== null) {
    chips.push({
      key: 'includeChildren',
      label: t.filters.includeChildren.value,
      value: advancedFilterValues.value.includeChildren
        ? t.filters.options.yes.value
        : t.filters.options.no.value,
    });
  }

  if (advancedFilterValues.value.enabled !== null) {
    chips.push({
      key: 'enabled',
      label: t.filters.status.value,
      value: advancedFilterValues.value.enabled
        ? t.filters.options.enabled.value
        : t.filters.options.disabled.value,
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
    { key: 'organization', label: t.menuColumns.organization.value, visible: columnVisibilityState.value.organization },
    { key: 'roles', label: t.menuColumns.roles.value, visible: columnVisibilityState.value.roles },
    { key: 'scope', label: t.menuColumns.scope.value, visible: columnVisibilityState.value.scope },
  ];

  return cols;
});

/**
 * Filtered columns based on visibility
 */
const visibleColumns = computed(() => {
  return t.columns.value.filter((col: any) => {
    // Always show avatar, name (assignee)
    if (col.key === 'avatar' || col.key === 'assigneeName') {
      return true;
    }

    // Filter based on columnVisibility
    if (col.key === 'orgName') return columnVisibilityState.value.organization;
    if (col.key === 'roleNames') return columnVisibilityState.value.roles;
    if (col.key === 'scope') return columnVisibilityState.value.scope;

    return true;
  });
});

/** FUNCTIONS */

/**
 * Apply quick filters and fetch data
 * @returns {void}
 */
function applyQuickFilters(): void {
  filters.value.assigneeType = quickAssigneeType.value ?? undefined;
  currentPage.value = 1;
  void fetchMemberships();
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
  filters.value.scope = appliedFilters.scope || undefined;
  filters.value.includeChildren = appliedFilters.includeChildren;
  filters.value.enabled = appliedFilters.enabled;

  // Reset pending state after apply
  hasPendingAdvancedFilters.value = false;

  currentPage.value = 1;
  showFiltersDrawer.value = false;
  void fetchMemberships();
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
    scope: null,
    includeChildren: null,
    enabled: null,
  };

  filters.value.scope = undefined;
  filters.value.includeChildren = undefined;
  filters.value.enabled = undefined;

  currentPage.value = 1;
  void fetchMemberships();
}

/**
 * Remove individual filter
 * @param {string} key - Filter key to remove
 * @returns {void}
 */
function removeFilter(key: string): void {
  if (key === 'search') {
    quickSearch.value = '';
  } else if (key === 'assigneeType') {
    quickAssigneeType.value = null;
    filters.value.assigneeType = undefined;
  } else if (key === 'scope') {
    advancedFilterValues.value.scope = null;
    filters.value.scope = undefined;
  } else if (key === 'includeChildren') {
    advancedFilterValues.value.includeChildren = null;
    filters.value.includeChildren = undefined;
  } else if (key === 'enabled') {
    advancedFilterValues.value.enabled = null;
    filters.value.enabled = undefined;
  }

  currentPage.value = 1;
  void fetchMemberships();
}

/**
 * Clear all filters (quick and advanced)
 * @returns {void}
 */
function clearAllFilters(): void {
  quickSearch.value = '';
  quickAssigneeType.value = null;
  advancedFilterValues.value = {
    scope: null,
    includeChildren: null,
    enabled: null,
  };
  filters.value = { ...FILTER_DEFAULTS };

  // Reset pending state
  hasPendingAdvancedFilters.value = false;

  currentPage.value = 1;
  void fetchMemberships();
}

/**
 * Fetch memberships from API with current filters and pagination
 * @returns {Promise<void>}
 */
async function fetchMemberships(): Promise<void> {
  if (!apis.mapexOS?.memberships) {
    error.value = 'Memberships API not initialized';
    return;
  }

  try {
    loading.value = true;
    error.value = undefined;

    // Build query parameters from filters and pagination state
    const queryParams: Record<string, any> = {
      page: currentPage.value,
      perPage: itemsPerPage.value,
      projection: MEMBERSHIPS_PROJECTION,
    };

    // Add active filters conditionally (only if they have values)
    if (filters.value.assigneeType) {
      queryParams.assigneeType = filters.value.assigneeType;
    }
    if (filters.value.assigneeId) {
      queryParams.assigneeId = filters.value.assigneeId;
    }
    if (filters.value.roleId) {
      queryParams.roleId = filters.value.roleId;
    }
    if (filters.value.scope) {
      queryParams.scope = filters.value.scope;
    }
    if (typeof filters.value.enabled === 'boolean') {
      queryParams.enabled = filters.value.enabled;
    }
    if (typeof filters.value.includeChildren === 'boolean') {
      queryParams.includeChildren = filters.value.includeChildren;
    }

    // Clean undefined values to avoid sending "undefined" as string in URL
    const cleanedParams = cleanQueryParams(queryParams);

    const response = await apis.mapexOS.memberships.list(cleanedParams);

    // Enrich memberships with organization and role names
    const enrichedMemberships = enrichMemberships(response.items || []);

    // Client-side search filter by assignee name
    if (quickSearch.value) {
      const searchLower = quickSearch.value.toLowerCase();
      membershipsList.value = enrichedMemberships.filter(m =>
        m.assigneeName.toLowerCase().includes(searchLower)
      );
    } else {
      membershipsList.value = enrichedMemberships;
    }

    // Update pagination state from response
    if (response.pagination) {
      totalPages.value = response.pagination.totalPages || 1;
      totalItems.value = response.pagination.totalItems || 0;
    }
  } catch (err: any) {
    logger.error('Error fetching memberships:', err);
    const errorMsg = err.message || 'Failed to fetch memberships';
    error.value = errorMsg;
    notifyFail({ message: errorMsg });
  } finally {
    loading.value = false;
    lastUpdatedAt.value = Date.now();
  }
}

/**
 * Fetch roles and cache them for name resolution
 * @returns {Promise<void>}
 */
async function fetchRolesForCache(): Promise<void> {
  if (!apis.mapexOS?.roles || rolesCache.value.size > 0) {
    return;
  }

  try {
    const response = await apis.mapexOS.roles.list({ perPage: 100 });
    (response.items || []).forEach((role: any) => {
      if (role.id && role.name) {
        rolesCache.value.set(role.id, role.name);
      }
    });
  } catch (err) {
    logger.warn('Failed to fetch roles for cache:', err);
  }
}

/**
 * Fetch users and cache them for name resolution
 * @returns {Promise<void>}
 */
async function fetchUsersForCache(): Promise<void> {
  if (!apis.mapexOS?.users || usersCache.value.size > 0) {
    return;
  }

  try {
    const response = await apis.mapexOS.users.list({ perPage: 100, projection: 'id,firstName,lastName,email' });
    (response.items || []).forEach((user: any) => {
      if (user.id) {
        const firstName = user.firstName || '';
        const lastName = user.lastName || '';
        const fullName = `${firstName} ${lastName}`.trim();
        usersCache.value.set(user.id, fullName || user.email || user.id);
      }
    });
  } catch (err) {
    logger.warn('Failed to fetch users for cache:', err);
  }
}

/**
 * Fetch groups and cache them for name resolution
 * @returns {Promise<void>}
 */
async function fetchGroupsForCache(): Promise<void> {
  if (!apis.mapexOS?.groups || groupsCache.value.size > 0) {
    return;
  }

  try {
    const response = await apis.mapexOS.groups.list({ perPage: 100, projection: 'id,name' });
    (response.items || []).forEach((group: any) => {
      if (group.id && group.name) {
        groupsCache.value.set(group.id, group.name);
      }
    });
  } catch (err) {
    logger.warn('Failed to fetch groups for cache:', err);
  }
}

/**
 * Enrich memberships with organization and role names
 * @param {any[]} memberships - Raw membership data
 * @returns {EnrichedMembership[]} Enriched memberships
 */
function enrichMemberships(memberships: any[]): EnrichedMembership[] {
  return memberships.map((membership: any) => {
    // Find organization
    const organization = orgStore.flatList.find((org: any) => org.id === membership.orgId);

    // Map role IDs to role names using cache
    const roleNames = (membership.roleIds || []).map((roleId: string) => {
      return rolesCache.value.get(roleId) || roleId;
    });

    // Get assignee name from appropriate cache based on type
    let assigneeName = membership.assigneeName || '';
    if (!assigneeName && membership.assigneeId) {
      if (membership.assigneeType === 'user') {
        assigneeName = usersCache.value.get(membership.assigneeId) || membership.assigneeId;
      } else if (membership.assigneeType === 'group') {
        assigneeName = groupsCache.value.get(membership.assigneeId) || membership.assigneeId;
      }
    }

    return {
      id: membership.id || '',
      assigneeType: membership.assigneeType || '',
      assigneeId: membership.assigneeId || '',
      assigneeName: assigneeName || membership.assigneeId || '',
      orgId: membership.orgId || '',
      orgName: organization?.name || 'Unknown',
      orgPathKey: membership.orgPathKey || '',
      roleIds: membership.roleIds || [],
      roleNames,
      scope: membership.scope || '',
      enabled: membership.enabled ?? true,
      created: membership.created || '',
    };
  });
}

/**
 * Handle pagination navigation
 * @param {number} page - New page number
 * @returns {void}
 */
function handlePageChange(page: number): void {
  currentPage.value = page;
  void fetchMemberships();
}

/**
 * Handle items per page change
 * @param {number} newValue - New items per page value
 * @returns {void}
 */
function handleItemsPerPageChange(newValue: number): void {
  itemsPerPage.value = newValue;
  currentPage.value = 1; // Reset to first page
  void fetchMemberships();
}

/**
 * Update menu columns when changed
 * @param {ListHeaderMenuColumn[]} columns - Updated column visibility states
 * @returns {void}
 */
function handleColumnsUpdate(columns: ListHeaderMenuColumn[]): void {
  columns.forEach(col => {
    if (col.key === 'organization') columnVisibilityState.value.organization = col.visible;
    if (col.key === 'roles') columnVisibilityState.value.roles = col.visible;
    if (col.key === 'scope') columnVisibilityState.value.scope = col.visible;
    if (col.key === 'enabled') columnVisibilityState.value.enabled = col.visible;
  });
}

/**
 * Navigate to user or group detail page based on assignee type
 * @param {EnrichedMembership} membership - Membership to view assignee
 * @returns {void}
 */
function viewAssignee(membership: EnrichedMembership): void {
  if (!membership.assigneeId) {
    notifyFail({ message: t.errors.assigneeIdMissing.value });
    return;
  }

  if (membership.assigneeType === 'user') {
    void router.push(`/users/detail/${membership.assigneeId}`);
  } else if (membership.assigneeType === 'group') {
    void router.push(`/groups/detail/${membership.assigneeId}`);
  }
}

/** LIFECYCLE HOOKS */
onMounted(async () => {
  // Fetch all caches in parallel for name resolution
  await Promise.all([
    fetchRolesForCache(),
    fetchUsersForCache(),
    fetchGroupsForCache(),
  ]);
  await fetchMemberships();
});
</script>

<template>
  <q-page class="q-pt-lg">

    <!-- Header Section -->
    <PageHeader
        icon="security"
        iconColor="primary"
        :title="t.page.title.value"
        :description="t.page.description.value"
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

      <!-- Assignee Type Select -->
      <div class="col-auto" style="min-width: 140px;">
        <q-select
          v-model="quickAssigneeType"
          outlined
          dense
          emit-value
          map-options
          :options="assigneeTypeOptions"
          :label="t.filters.assigneeType.value"
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
          <q-icon name="security" size="sm" color="primary" class="q-mr-sm"/>
          <div class="text-subtitle1 text-weight-medium text-primary">{{ t.page.listTitle }}</div>
        </div>
      </div>
      <div class="col-auto">
        <ListHeaderMenu
          icon="security"
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
          @refresh="fetchMemberships"
        />
      </div>
    </div>

    <!-- Loading Spinner -->
    <div v-if="loading" class="row justify-center q-my-lg">
      <q-spinner color="primary" size="3em" />
    </div>

    <!-- Memberships Row List -->
    <div v-else class="row">
      <div
          v-for="(membership, index) in membershipsList"
          :key="membership.id || `membership-${index}`"
          class="col-12 q-mb-xs"
      >
        <DataRow
            :data="membership"
            :columns="visibleColumns"
            :show-actions="false"
            @click="viewAssignee"
            @dblclick="viewAssignee"
        />
      </div>

      <!-- No Results -->
      <div v-if="membershipsList.length === 0" class="col-12">
        <ListCardEmpty
            :title="t.empty.title.value"
            :description="t.empty.description.value"
            icon="security"
        />
      </div>
    </div>

    <!-- Pagination -->
    <ListPagination
      v-model="currentPage"
      :total-pages="totalPages"
      @change="handlePageChange"
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
