<script setup lang="ts">
defineOptions({
  name: 'RolesListPage'
});

/** TYPE IMPORTS */
import type { ListHeaderMenuColumn } from '@components/headers';
import type { FilterField } from '@components/drawers';
import type { RoleResponse, RoleQuery } from '@mapexos/schemas';
import type {
  RolesListPageFilters,
  RolesListPageColumnVisibility,
} from './interfaces';

/** VUE IMPORTS */
import { ref, computed, onMounted } from 'vue';
import { useRouter } from 'vue-router';

/** COMPONENTS */
import { PageHeader, ListHeaderMenu } from '@components/headers';
import { ListCardEmpty, DataRow } from '@components/cards';
import { RoleDetailsDrawer, AdvancedFiltersDrawer } from '@components/drawers';
import { ListPagination } from '@components/navigation';
import { AppTooltip } from '@components/tooltips';

/** COMPOSABLES */
import { useRolesTranslations } from '@composables/i18n';
import { usePermissions } from '@composables/shared/usePermissions';
import { useLogger } from '@composables/useLogger';

/** UTILS */
import { notifySuccess, notifyFail, notifyWarning, dialogDelete } from '@utils/alert';

/** SERVICES */
import { apis } from '@services/mapex';

/** STORES */
import { useOrganizationStore } from '@stores/organization';

/** LOCAL IMPORTS */
import {
  DEFAULT_ITEMS_PER_PAGE,
  COLUMN_VISIBILITY_DEFAULTS,
  FILTER_DEFAULTS,
  ROLES_PROJECTION,
} from './constants';

/** COMPOSABLES & STORES */
const t = useRolesTranslations();
const orgStore = useOrganizationStore();
const logger = useLogger('RolesListPage');
const router = useRouter();
const { canCreate, canUpdate, canDelete, canRead } = usePermissions();
const canCreateRole = canCreate('roles');
const canUpdateRole = canUpdate('roles');
const canDeleteRole = canDelete('roles');
const canReadRole = canRead('roles');

/** STATE */
const rolesList = ref<RoleResponse[]>([]);
const loading = ref(false);
const lastUpdatedAt = ref<number | undefined>(undefined);
const error = ref<string | undefined>(undefined);
const itemsPerPage = ref(DEFAULT_ITEMS_PER_PAGE);
const currentPage = ref(1);
const totalPages = ref(1);
const totalItems = ref(0);
const filters = ref<RolesListPageFilters>({ ...FILTER_DEFAULTS });
const columnVisibilityState = ref<RolesListPageColumnVisibility>({ ...COLUMN_VISIBILITY_DEFAULTS });
const showDetailsDrawer = ref(false);
const showFiltersDrawer = ref(false);
const selectedRoleId = ref<string | null>(null);

// Quick filter state (inline)
const quickSearch = ref('');
const quickIsSystem = ref<boolean | null>(null);

// Advanced filter state (drawer)
const advancedFilterValues = ref<Record<string, any>>({
  includeChildren: null,
  scope: null,
  permission: null,
  isTemplate: null,
});

// Pending changes state
const hasPendingAdvancedFilters = ref(false);

/** COMPUTED */

/**
 * Status options for quick filter select (isSystem)
 */
const isSystemOptions = computed(() => [
  { label: t.filters.allStatus.value, value: null },
  { label: t.filters.options.system.value, value: true },
  { label: t.filters.options.custom.value, value: false },
]);

/**
 * Advanced filter fields for the drawer
 */
const advancedFilterFields = computed((): FilterField[] => {
  const fields: FilterField[] = [
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
      key: 'scope',
      label: t.filters.scope.value,
      type: 'toggle',
      icon: 'public',
      options: [
        { label: t.filters.allStatus.value, value: null },
        { label: t.filters.options.global.value, value: 'global' },
        { label: t.filters.options.local.value, value: 'local' },
      ],
    },
    {
      key: 'permission',
      label: t.filters.permission.value,
      type: 'input',
      icon: 'vpn_key',
      placeholder: t.filters.filterByPermission.value,
    },
  ];

  // Add isTemplate filter only for Customer and Site organizations
  if (orgStore.isCustomer || orgStore.isSite) {
    fields.push({
      key: 'isTemplate',
      label: t.filters.isTemplate.value,
      type: 'toggle',
      icon: 'content_copy',
      options: [
        { label: t.filters.allStatus.value, value: null },
        { label: t.filters.options.templates.value, value: true },
        { label: t.filters.options.local.value, value: false },
      ],
    });
  }

  return fields;
});

/**
 * Check if any filters are active (quick or advanced)
 */
const hasActiveFilters = computed(() => {
  return !!(
    quickSearch.value ||
    quickIsSystem.value !== null ||
    advancedFilterValues.value.includeChildren !== null ||
    advancedFilterValues.value.scope ||
    advancedFilterValues.value.permission ||
    advancedFilterValues.value.isTemplate !== null
  );
});

/**
 * Count of active advanced filters (for badge)
 */
const advancedFiltersCount = computed(() => {
  let count = 0;
  if (advancedFilterValues.value.includeChildren !== null) count++;
  if (advancedFilterValues.value.scope) count++;
  if (advancedFilterValues.value.permission) count++;
  if (advancedFilterValues.value.isTemplate !== null) count++;
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
      label: t.filters.name.value,
      value: quickSearch.value,
    });
  }

  if (quickIsSystem.value !== null) {
    chips.push({
      key: 'isSystem',
      label: t.filters.isSystem.value,
      value: quickIsSystem.value ? t.filters.options.system.value : t.filters.options.custom.value,
    });
  }

  if (advancedFilterValues.value.includeChildren !== null) {
    chips.push({
      key: 'includeChildren',
      label: t.filters.includeChildren.value,
      value: advancedFilterValues.value.includeChildren ? t.filters.options.yes.value : t.filters.options.no.value,
    });
  }

  if (advancedFilterValues.value.scope) {
    chips.push({
      key: 'scope',
      label: t.filters.scope.value,
      value: advancedFilterValues.value.scope === 'global' ? t.filters.options.global.value : t.filters.options.local.value,
    });
  }

  if (advancedFilterValues.value.permission) {
    chips.push({
      key: 'permission',
      label: t.filters.permission.value,
      value: advancedFilterValues.value.permission,
    });
  }

  if (advancedFilterValues.value.isTemplate !== null) {
    chips.push({
      key: 'isTemplate',
      label: t.filters.isTemplate.value,
      value: advancedFilterValues.value.isTemplate ? t.filters.options.templates.value : t.filters.options.local.value,
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
    { key: 'description', label: t.menuColumns.description.value, visible: columnVisibilityState.value.description },
    { key: 'permissions', label: t.menuColumns.permissions.value, visible: columnVisibilityState.value.permissions },
    { key: 'scope', label: t.menuColumns.scope.value, visible: columnVisibilityState.value.scope },
    { key: 'isTemplate', label: t.menuColumns.templateSource.value, visible: columnVisibilityState.value.isTemplate },
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
    // Always show avatar, name, and type
    if (col.key === 'avatar' || col.key === 'name' || col.key === 'type') {
      return true;
    }

    // Organization column only visible when includeChildren filter is active
    if (col.key === 'organizationName') {
      return advancedFilterValues.value.includeChildren === true && columnVisibilityState.value.organization;
    }

    // Filter based on columnVisibility
    if (col.key === 'description') return columnVisibilityState.value.description;
    if (col.key === 'permissions') return columnVisibilityState.value.permissions;
    if (col.key === 'scope') return columnVisibilityState.value.scope;
    if (col.key === 'isTemplate') return columnVisibilityState.value.isTemplate;
    if (col.key === 'created') return columnVisibilityState.value.created;

    return true;
  });
});

/** FUNCTIONS */

/**
 * Fetch roles from API with current filters and pagination
 * @returns {Promise<void>}
 */
async function fetchRoles(): Promise<void> {
  if (!apis.mapexOS?.roles) {
    error.value = 'Roles API not initialized';
    return;
  }

  try {
    loading.value = true;
    error.value = undefined;

    // Build query parameters from filters and pagination state
    const queryParams: Record<string, any> = {
      page: currentPage.value,
      perPage: itemsPerPage.value,
      projection: ROLES_PROJECTION,
    };

    // Add active filters conditionally (only if they have values)
    if (filters.value.name) {
      queryParams.name = filters.value.name;
    }
    if (typeof filters.value.isSystem === 'boolean') {
      queryParams.isSystem = filters.value.isSystem;
    }
    if (filters.value.scope) {
      queryParams.scope = filters.value.scope;
    }
    if (filters.value.permission) {
      queryParams.permission = filters.value.permission;
    }
    if (typeof filters.value.includeChildren === 'boolean') {
      queryParams.includeChildren = filters.value.includeChildren;
    }
    if (typeof filters.value.isTemplate === 'boolean') {
      queryParams.isTemplate = filters.value.isTemplate;
    }

    const response = await apis.mapexOS.roles.list(queryParams as RoleQuery);

    // Enrich roles with organization name
    const enrichedRoles = (response.items || []).map((role: any) => {
      const organization = orgStore.flatList.find((org: any) => org.id === role.orgId);
      return {
        ...role,
        organizationName: organization?.name || 'Unknown',
      };
    });

    rolesList.value = enrichedRoles;

    // Update pagination state from response
    if (response.pagination) {
      totalPages.value = response.pagination.totalPages || 1;
      totalItems.value = response.pagination.totalItems || 0;
    }
  } catch (err: any) {
    logger.error('Error fetching roles:', err);
    const errorMsg = err.message || 'Failed to fetch roles';
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
  filters.value.isSystem = quickIsSystem.value ?? undefined;
  currentPage.value = 1;
  void fetchRoles();
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
  filters.value.scope = appliedFilters.scope || undefined;
  filters.value.permission = appliedFilters.permission || undefined;
  filters.value.isTemplate = appliedFilters.isTemplate;

  // Reset pending state after apply
  hasPendingAdvancedFilters.value = false;

  // Auto-hide columns to prevent horizontal scroll when includeChildren is active
  if (appliedFilters.includeChildren === true) {
    columnVisibilityState.value.description = false;
    columnVisibilityState.value.permissions = false;
    columnVisibilityState.value.created = false;
    columnVisibilityState.value.isTemplate = false;
  } else {
    columnVisibilityState.value.description = true;
    columnVisibilityState.value.permissions = true;
    columnVisibilityState.value.created = true;
    columnVisibilityState.value.isTemplate = true;
  }

  currentPage.value = 1;
  showFiltersDrawer.value = false;
  void fetchRoles();
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
    scope: null,
    permission: null,
    isTemplate: null,
  };

  filters.value.includeChildren = undefined;
  filters.value.scope = undefined;
  filters.value.permission = undefined;
  filters.value.isTemplate = undefined;

  // Restore column visibility
  columnVisibilityState.value.description = true;
  columnVisibilityState.value.permissions = true;
  columnVisibilityState.value.created = true;
  columnVisibilityState.value.isTemplate = true;

  currentPage.value = 1;
  void fetchRoles();
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
  } else if (key === 'isSystem') {
    quickIsSystem.value = null;
    filters.value.isSystem = undefined;
  } else if (key === 'includeChildren') {
    advancedFilterValues.value.includeChildren = null;
    filters.value.includeChildren = undefined;
    // Restore column visibility
    columnVisibilityState.value.description = true;
    columnVisibilityState.value.permissions = true;
    columnVisibilityState.value.created = true;
    columnVisibilityState.value.isTemplate = true;
  } else if (key === 'scope') {
    advancedFilterValues.value.scope = null;
    filters.value.scope = undefined;
  } else if (key === 'permission') {
    advancedFilterValues.value.permission = null;
    filters.value.permission = undefined;
  } else if (key === 'isTemplate') {
    advancedFilterValues.value.isTemplate = null;
    filters.value.isTemplate = undefined;
  }

  currentPage.value = 1;
  void fetchRoles();
}

/**
 * Clear all filters (quick and advanced)
 * @returns {void}
 */
function clearAllFilters(): void {
  quickSearch.value = '';
  quickIsSystem.value = null;
  advancedFilterValues.value = {
    includeChildren: null,
    scope: null,
    permission: null,
    isTemplate: null,
  };
  filters.value = { ...FILTER_DEFAULTS };

  // Reset pending state
  hasPendingAdvancedFilters.value = false;

  // Restore column visibility
  columnVisibilityState.value.description = true;
  columnVisibilityState.value.permissions = true;
  columnVisibilityState.value.created = true;
  columnVisibilityState.value.isTemplate = true;

  currentPage.value = 1;
  void fetchRoles();
}

/**
 * Handle pagination navigation
 * @param {number} page - New page number
 * @returns {void}
 */
function handlePageChange(page: number): void {
  currentPage.value = page;
  void fetchRoles();
}

/**
 * Handle items per page change
 * @param {number} newValue - New items per page value
 * @returns {void}
 */
function handleItemsPerPageChange(newValue: number): void {
  itemsPerPage.value = newValue;
  currentPage.value = 1;
  void fetchRoles();
}

/**
 * Update menu columns when changed
 * @param {ListHeaderMenuColumn[]} columns - Updated column visibility states
 * @returns {void}
 */
function handleColumnsUpdate(columns: ListHeaderMenuColumn[]): void {
  columns.forEach(col => {
    if (col.key === 'organization') columnVisibilityState.value.organization = col.visible;
    if (col.key === 'description') columnVisibilityState.value.description = col.visible;
    if (col.key === 'permissions') columnVisibilityState.value.permissions = col.visible;
    if (col.key === 'scope') columnVisibilityState.value.scope = col.visible;
    if (col.key === 'isTemplate') columnVisibilityState.value.isTemplate = col.visible;
    if (col.key === 'created') columnVisibilityState.value.created = col.visible;
  });
}

/**
 * Check if user can edit/delete a role
 * @param {any} role - Role to check
 * @returns {boolean} True if user can modify the role
 */
function canModifyRole(role: any): boolean {
  // System roles cannot be modified
  if (role.isSystem) {
    return false;
  }

  // Shared templates can only be modified by the owner organization
  if (role.isTemplate) {
    return role.orgId === orgStore.selectedOrganizationId;
  }

  // Local resources can always be modified
  return true;
}

/**
 * View role details in drawer
 * @param {RoleResponse} role - Role to view
 * @returns {void}
 */
function viewDetails(role: RoleResponse): void {
  if (!canReadRole.value) return;
  if (!role.id) {
    notifyFail({ message: t.errors.idMissing.value });
    return;
  }
  selectedRoleId.value = role.id;
  showDetailsDrawer.value = true;
}

/**
 * Edit role - Navigate to edit page
 * @param {any} role - Role to edit
 * @returns {void}
 */
function editRole(role: any): void {
  if (!canUpdateRole.value) return;
  if (!canModifyRole(role)) {
    if (role.isSystem) {
      notifyWarning({ message: t.notifications.systemEdit.value });
    } else if (role.isTemplate) {
      notifyWarning({ message: t.notifications.sharedEdit.value });
    }
    return;
  }

  if (!role.id) {
    notifyFail({ message: t.errors.idMissing.value });
    return;
  }

  logger.debug('Navigating to edit role:', role.id);
  void router.push(`/roles/edit/${role.id}`);
}

/**
 * Handle edit event from drawer - Navigate to edit page
 * @param {string} roleId - ID of role to edit
 * @returns {void}
 */
function handleDrawerEdit(roleId: string): void {
  if (!canUpdateRole.value) return;
  logger.debug('Navigating to edit role from drawer:', roleId);
  void router.push(`/roles/edit/${roleId}`);
}

/**
 * Confirm delete operation with dialog
 * @param {RoleResponse} role - Role to delete
 * @returns {Promise<void>}
 */
async function confirmDelete(role: RoleResponse): Promise<void> {
  if (!canDeleteRole.value) return;
  if (!canModifyRole(role)) {
    if (role.isSystem) {
      notifyWarning({ message: t.notifications.systemDelete.value });
    } else if (role.isTemplate) {
      notifyWarning({ message: t.notifications.sharedDelete.value });
    }
    return;
  }

  const roleName = role.name || 'this role';
  const confirmed = await dialogDelete({
    title: t.dialog.deleteTitle.value,
    message: t.messages.confirmDelete(roleName),
  });

  if (confirmed) {
    await deleteRole(role);
  }
}

/**
 * Delete role from API and update list
 * @param {RoleResponse} role - Role to delete
 * @returns {Promise<void>}
 */
async function deleteRole(role: RoleResponse): Promise<void> {
  if (!apis.mapexOS?.roles) {
    notifyFail({ message: t.errors.apiNotInitialized.value });
    return;
  }

  if (!role.id) {
    notifyFail({ message: t.errors.idMissing.value });
    return;
  }

  try {
    await apis.mapexOS.roles.delete({ roleId: role.id });

    // Update total items count after successful deletion
    totalItems.value = Math.max(0, totalItems.value - 1);

    // Remove from local list
    const index = rolesList.value.findIndex(r => r.id === role.id);
    if (index !== -1) {
      rolesList.value.splice(index, 1);
    }

    // If current page is now empty and not the first page, go to previous page
    if (rolesList.value.length === 0 && currentPage.value > 1) {
      currentPage.value -= 1;
      await fetchRoles();
    }

    notifySuccess({ message: t.messages.deletedSuccessfully.value });
  } catch (err: any) {
    notifyFail({ message: err.message || 'Failed to delete role' });
  }
}

/** LIFECYCLE HOOKS */
onMounted(async () => await fetchRoles());

</script>

<template>
  <q-page class="q-pt-lg">

    <!-- Header Section -->
    <PageHeader
        icon="admin_panel_settings"
        iconColor="primary"
        :title="t.page.title.value"
        :description="t.page.description.value"
        :button="canCreateRole ? { label: t.page.addButton.value, icon: 'add', color: 'primary', to: '/roles/add' } : undefined"
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

      <!-- Type Select (isSystem) -->
      <div class="col-auto" style="min-width: 140px;">
        <q-select
          v-model="quickIsSystem"
          outlined
          dense
          emit-value
          map-options
          :options="isSystemOptions"
          :label="t.filters.isSystem.value"
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
          <q-icon name="admin_panel_settings" size="sm" color="primary" class="q-mr-sm"/>
          <div class="text-subtitle1 text-weight-medium text-primary">{{ t.page.listTitle.value }}</div>
        </div>
      </div>
      <div class="col-auto">
        <ListHeaderMenu
          icon="admin_panel_settings"
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
          @refresh="fetchRoles"
        />
      </div>
    </div>

    <!-- Loading Spinner -->
    <div v-if="loading" class="row justify-center q-my-lg">
      <q-spinner color="primary" size="3em" />
    </div>

    <!-- Roles Row List -->
    <div v-else class="row">
      <div
          v-for="(role, index) in rolesList"
          :key="role.id || `role-${index}`"
          class="col-12 q-mb-xs"
      >
        <DataRow
            :data="role"
            :columns="visibleColumns"
            :actions="{
              showEdit: canUpdateRole && canModifyRole(role),
              showView: canReadRole,
              showDelete: canDeleteRole && canModifyRole(role),
            }"
            @click="viewDetails"
            @dblclick="editRole"
            @edit="editRole"
            @view="viewDetails"
            @delete="confirmDelete"
        />
      </div>

      <!-- No Results -->
      <div v-if="rolesList.length === 0" class="col-12">
        <ListCardEmpty
            :title="t.empty.title.value"
            :description="t.empty.description.value"
            icon="admin_panel_settings"
        />
      </div>
    </div>

    <!-- Pagination -->
    <ListPagination
      v-model="currentPage"
      :total-pages="totalPages"
      @change="handlePageChange"
    />

    <!-- Role Details Drawer -->
    <RoleDetailsDrawer
      v-model="showDetailsDrawer"
      :role-id="selectedRoleId"
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
