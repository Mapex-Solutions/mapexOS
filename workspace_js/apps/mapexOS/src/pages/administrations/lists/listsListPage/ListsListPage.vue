<script setup lang="ts">
defineOptions({
	name: 'ListsListPage'
});

/** TYPE IMPORTS */
import type { ListHeaderMenuColumn } from '@components/headers';
import type { DataRowActionConfig } from '@components/cards';
import type { FilterField, FilterValues } from '@components/drawers';
import type { ListResponse } from '@mapexos/schemas';
import type {
	ListsListPageFilters,
	ListsListPageColumnVisibility,
} from './interfaces';

/** VUE IMPORTS */
import { ref, computed, onMounted } from 'vue';

/** COMPONENTS */
import { PageHeader, ListHeaderMenu } from '@components/headers';
import { ListCardEmpty, DataRow } from '@components/cards';
import { ListDrawer, AdvancedFiltersDrawer } from '@components/drawers';
import { ListPagination } from '@components/navigation';
import { AppTooltip } from '@components/tooltips';

/** COMPOSABLES */
import { useListsTranslations } from '@composables/i18n/pages/administrations/lists/useListsTranslations';
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
	COLUMN_VISIBILITY_DEFAULTS,
	LISTS_PROJECTION,
} from './constants';

/** COMPOSABLES & STORES */
const t = useListsTranslations();
const orgStore = useOrganizationStore();
const logger = useLogger('ListsListPage');
const { canUpdate, canDelete, canRead } = usePermissions();
const canUpdateList = canUpdate('lists');
const canDeleteList = canDelete('lists');
const canReadList = canRead('lists');

/** CONSTANTS */
const MAX_VISIBLE_CHIPS = 2;

/** STATE */
const listsList = ref<ListResponse[]>([]);
const loading = ref(false);
const lastUpdatedAt = ref<number | undefined>(undefined);
const error = ref<string | undefined>(undefined);
const itemsPerPage = ref(DEFAULT_ITEMS_PER_PAGE);
const currentPage = ref(1);
const totalPages = ref(1);
const totalItems = ref(0);

// Quick filters (applied on enter or select change)
const quickSearch = ref('');
const quickSource = ref<boolean | null>(null);

// Advanced filters (applied via drawer)
const showAdvancedFilters = ref(false);
const advancedFilterValues = ref<FilterValues>({
	category: 'iot',
	type: null,
	isTemplate: null,
	includeChildren: null,
});
const hasPendingAdvancedFilters = ref(false);

const columnVisibilityState = ref<ListsListPageColumnVisibility>({ ...COLUMN_VISIBILITY_DEFAULTS });
const drawerOpen = ref(false);
const selectedListId = ref<string | undefined>(undefined);

/** COMPUTED */

/**
 * Combined filters for API request
 */
const appliedFilters = computed((): ListsListPageFilters => ({
	name: quickSearch.value || undefined,
	category: advancedFilterValues.value.category || undefined,
	type: advancedFilterValues.value.type || undefined,
	isSystem: quickSource.value ?? undefined,
	isTemplate: advancedFilterValues.value.isTemplate ?? undefined,
	includeChildren: advancedFilterValues.value.includeChildren ?? undefined,
}));

/**
 * Advanced filter fields configuration
 */
const advancedFilterFields = computed((): FilterField[] => [
	{
		key: 'category',
		type: 'select',
		label: t.filters.category.value,
		icon: 'folder_special',
		options: [
			{ label: t.filters.options.iot.value, value: 'iot' },
		],
	},
	{
		key: 'type',
		type: 'select',
		label: t.filters.type.value,
		icon: 'category',
		options: [
			{ label: t.filters.options.categories.value, value: 'asset_category' },
			{ label: t.filters.options.manufacturers.value, value: 'asset_manufacturer' },
			{ label: t.filters.options.models.value, value: 'asset_model' },
		],
		disabled: !advancedFilterValues.value.category,
	},
	{
		key: 'isTemplate',
		type: 'toggle',
		label: t.filters.isTemplate.value,
		icon: 'content_copy',
		options: [
			{ label: t.filters.options.all.value, value: null },
			{ label: t.filters.options.shared.value, value: true },
			{ label: t.filters.options.local.value, value: false },
		],
	},
	{
		key: 'includeChildren',
		type: 'toggle',
		label: t.filters.includeChildren.value,
		icon: 'account_tree',
		options: [
			{ label: t.filters.options.all.value, value: null },
			{ label: t.filters.options.yes.value, value: true },
			{ label: t.filters.options.no.value, value: false },
		],
	},
]);

/**
 * Count of advanced filters only (for badge)
 */
const advancedFiltersCount = computed(() => {
	let count = 0;
	if (advancedFilterValues.value.category) count++;
	if (advancedFilterValues.value.type) count++;
	if (advancedFilterValues.value.isTemplate !== null) count++;
	if (advancedFilterValues.value.includeChildren !== null) count++;
	return count;
});

/**
 * Active filter chips for display
 */
const activeFilterChips = computed(() => {
	const chips: { key: string; label: string; value: string }[] = [];

	if (advancedFilterValues.value.category) {
		chips.push({
			key: 'category',
			label: t.filters.category.value,
			value: t.filters.options.iot.value,
		});
	}

	if (advancedFilterValues.value.type) {
		const typeLabel = advancedFilterValues.value.type === 'asset_category'
			? t.filters.options.categories.value
			: advancedFilterValues.value.type === 'asset_manufacturer'
				? t.filters.options.manufacturers.value
				: t.filters.options.models.value;
		chips.push({
			key: 'type',
			label: t.filters.type.value,
			value: typeLabel,
		});
	}

	if (advancedFilterValues.value.isTemplate !== null) {
		const label = advancedFilterValues.value.isTemplate
			? t.filters.options.shared.value
			: t.filters.options.local.value;
		chips.push({
			key: 'isTemplate',
			label: t.filters.isTemplate.value,
			value: label,
		});
	}

	if (advancedFilterValues.value.includeChildren !== null) {
		const label = advancedFilterValues.value.includeChildren
			? t.filters.options.yes.value
			: t.filters.options.no.value;
		chips.push({
			key: 'includeChildren',
			label: t.filters.includeChildren.value,
			value: label,
		});
	}

	return chips;
});

const visibleFilterChips = computed(() => activeFilterChips.value.slice(0, MAX_VISIBLE_CHIPS));
const hiddenFilterChips = computed(() => activeFilterChips.value.slice(MAX_VISIBLE_CHIPS));
const hiddenFiltersCount = computed(() => hiddenFilterChips.value.length);

/**
 * Check if any filter is active
 */
const hasActiveFilters = computed(() => {
	return activeFilterChips.value.length > 0 || quickSource.value !== null;
});

const menuColumns = computed(() => {
	const cols: ListHeaderMenuColumn[] = [
		{ key: 'category', label: t.menuColumns.category.value, visible: columnVisibilityState.value.category },
		{ key: 'type', label: t.menuColumns.type.value, visible: columnVisibilityState.value.type },
		{ key: 'source', label: t.menuColumns.source.value, visible: columnVisibilityState.value.source },
		{ key: 'scope', label: t.menuColumns.scope.value, visible: columnVisibilityState.value.scope },
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

// Filtered columns based on visibility
const visibleColumns = computed(() => {
	return t.columns.value.filter((col: any) => {
		// Always show icon and name
		if (col.key === 'icon' || col.key === 'name') {
			return true;
		}

		// Category column (always visible)
		if (col.key === 'category') {
			return columnVisibilityState.value.category;
		}

		// Type column
		if (col.key === 'type') {
			return columnVisibilityState.value.type;
		}

		// Organization column only visible when includeChildren filter is active
		if (col.key === 'organizationName') {
			return advancedFilterValues.value.includeChildren === true && columnVisibilityState.value.organization;
		}

		// Source (isSystem) column
		if (col.key === 'isSystem') {
			return columnVisibilityState.value.source;
		}

		// Scope (isTemplate) column
		if (col.key === 'isTemplate') {
			return columnVisibilityState.value.scope;
		}

		return true;
	});
});

/** FUNCTIONS */

/**
 * Fetch lists from API with current filters and pagination
 * @returns {Promise<void>}
 */
async function fetchLists(): Promise<void> {
	if (!apis.mapexOS?.lists) {
		error.value = 'Lists API not initialized';
		return;
	}

	try {
		loading.value = true;
		error.value = undefined;

		// Build query parameters from filters and pagination state
		const queryParams: Record<string, any> = {
			page: currentPage.value,
			perPage: itemsPerPage.value,
			projection: LISTS_PROJECTION,
		};

		// Add active filters to query params (only if they have values)
		if (appliedFilters.value.name) {
			queryParams.name = appliedFilters.value.name;
		}
		if (appliedFilters.value.category) {
			queryParams.category = appliedFilters.value.category;
		}
		if (appliedFilters.value.type) {
			queryParams.type = appliedFilters.value.type;
		}
		if (typeof appliedFilters.value.isSystem === 'boolean') {
			queryParams.isSystem = appliedFilters.value.isSystem;
		}
		if (typeof appliedFilters.value.isTemplate === 'boolean') {
			queryParams.isTemplate = appliedFilters.value.isTemplate;
		}
		if (typeof appliedFilters.value.includeChildren === 'boolean') {
			queryParams.includeChildren = appliedFilters.value.includeChildren;
		}

		// Clean undefined values to avoid sending "undefined" as string in URL
		const cleanedParams = cleanQueryParams(queryParams);

		const response = await apis.mapexOS.lists.list(cleanedParams);

		// Safely access response data with proper null checks
		const listsData = response?.items || [];

		// Enrich lists with organization name
		const enrichedLists = listsData.map((list: any) => {
			const organization = orgStore.flatList.find((org: any) => org.id === list.orgId);
			return {
				...list,
				organizationName: organization?.name || 'Unknown',
			};
		});

		listsList.value = enrichedLists;

		// Update pagination state from response
		if (response.pagination) {
			totalItems.value = response.pagination.totalItems || 0;
			totalPages.value = response.pagination.totalPages || 1;
		}
	} catch (err: any) {
		logger.error('Error fetching lists:', err);
		const errorMsg = err.message || t.notifications.loadFailed.value;
		error.value = errorMsg;
		notifyFail({ message: errorMsg });
	} finally {
		loading.value = false;
		lastUpdatedAt.value = Date.now();
	}
}

/**
 * Apply quick filters (search and source)
 * @returns {void}
 */
function applyQuickFilters(): void {
	currentPage.value = 1;
	void fetchLists();
}

/**
 * Handle advanced filters apply
 * @param {FilterValues} values - Applied filter values
 * @returns {void}
 */
function handleAdvancedFiltersApply(values: FilterValues): void {
	advancedFilterValues.value = {
		category: values.category || null,
		type: values.type || null,
		isTemplate: values.isTemplate,
		includeChildren: values.includeChildren,
	};

	// Auto-hide columns when includeChildren is active to prevent horizontal scroll
	if (values.includeChildren === true) {
		columnVisibilityState.value.type = false;
		columnVisibilityState.value.source = false;
		columnVisibilityState.value.scope = false;
	} else {
		// Restore columns when includeChildren is disabled
		columnVisibilityState.value.type = true;
		columnVisibilityState.value.source = true;
		columnVisibilityState.value.scope = true;
	}

	hasPendingAdvancedFilters.value = false;
	currentPage.value = 1;
	void fetchLists();
}

/**
 * Handle advanced filters reset
 * @returns {void}
 */
function handleAdvancedFiltersReset(): void {
	advancedFilterValues.value = {
		category: 'iot',
		type: null,
		isTemplate: null,
		includeChildren: null,
	};
	hasPendingAdvancedFilters.value = false;

	// Restore column visibility
	columnVisibilityState.value.type = true;
	columnVisibilityState.value.source = true;
	columnVisibilityState.value.scope = true;

	currentPage.value = 1;
	void fetchLists();
}

/**
 * Handle pending changes state
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
	if (key === 'category') {
		advancedFilterValues.value.category = null;
		// Clear type when category is cleared
		advancedFilterValues.value.type = null;
	} else if (key === 'type') {
		advancedFilterValues.value.type = null;
	} else if (key === 'isTemplate') {
		advancedFilterValues.value.isTemplate = null;
	} else if (key === 'includeChildren') {
		advancedFilterValues.value.includeChildren = null;
		// Restore column visibility
		columnVisibilityState.value.type = true;
		columnVisibilityState.value.source = true;
		columnVisibilityState.value.scope = true;
	}

	currentPage.value = 1;
	void fetchLists();
}

/**
 * Clear all filters
 * @returns {void}
 */
function clearAllFilters(): void {
	quickSearch.value = '';
	quickSource.value = null;
	advancedFilterValues.value = {
		category: 'iot',
		type: null,
		isTemplate: null,
		includeChildren: null,
	};
	hasPendingAdvancedFilters.value = false;

	// Restore column visibility
	columnVisibilityState.value.type = true;
	columnVisibilityState.value.source = true;
	columnVisibilityState.value.scope = true;

	currentPage.value = 1;
	void fetchLists();
}

/**
 * Handle pagination navigation
 * @param {number} page - New page number
 * @returns {void}
 */
function handlePageChange(page: number): void {
	currentPage.value = page;
	void fetchLists();
}

/**
 * Handle items per page change
 * @param {number} newValue - New items per page value
 * @returns {void}
 */
function handleItemsPerPageChange(newValue: number): void {
	itemsPerPage.value = newValue;
	currentPage.value = 1;
	void fetchLists();
}

/**
 * Update menu columns when changed
 * @param {ListHeaderMenuColumn[]} columns - Updated columns
 * @returns {void}
 */
function handleColumnsUpdate(columns: ListHeaderMenuColumn[]): void {
	columns.forEach(col => {
		if (col.key === 'organization') columnVisibilityState.value.organization = col.visible;
		if (col.key === 'category') columnVisibilityState.value.category = col.visible;
		if (col.key === 'type') columnVisibilityState.value.type = col.visible;
		if (col.key === 'source') columnVisibilityState.value.source = col.visible;
		if (col.key === 'scope') columnVisibilityState.value.scope = col.visible;
	});
}

/**
 * Check if user can edit/delete a list item
 * Rules:
 * - isSystem = true: Cannot edit/delete (system resource)
 * - isTemplate = true: Can only edit/delete if orgId matches current organization
 * - isTemplate = false: Can always edit/delete (local resource)
 * @param {any} list - List item to check
 * @returns {boolean} Whether the list can be modified
 */
function canModifyList(list: any): boolean {
	// System lists cannot be modified
	if (list.isSystem) {
		return false;
	}

	// Shared templates can only be modified by the owner organization
	if (list.isTemplate) {
		return list.orgId === orgStore.selectedOrganizationId;
	}

	// Local resources can always be modified
	return true;
}

/**
 * View list details
 * @param {any} list - List item to view
 * @returns {void}
 */
function viewDetails(list: any): void {
	if (!canReadList.value) return;
	selectedListId.value = list.id;
	drawerOpen.value = true;
}

/**
 * Edit list (with access control validation)
 * @param {any} list - List item to edit
 * @returns {void}
 */
function editList(list: any): void {
	if (!canUpdateList.value) return;
	if (!canModifyList(list)) {
		if (list.isSystem) {
			notifyWarning({ message: t.notifications.systemListEdit.value });
		} else if (list.isTemplate) {
			notifyWarning({ message: t.notifications.sharedListEdit.value });
		}
		return;
	}
	logger.debug('Edit list:', list);
	// TODO: Navigate to edit page
}

/**
 * Confirm delete list (with access control validation)
 * @param {any} list - List item to delete
 * @returns {Promise<void>}
 */
async function confirmDelete(list: any): Promise<void> {
	if (!canDeleteList.value) return;
	if (!canModifyList(list)) {
		if (list.isSystem) {
			notifyWarning({ message: t.notifications.systemListDelete.value });
		} else if (list.isTemplate) {
			notifyWarning({ message: t.notifications.sharedListDelete.value });
		}
		return;
	}

	const confirmed = await dialogDelete({
		title: t.dialog.confirmDelete.title.value,
		message: t.dialog.confirmDelete.message(list.name || 'this item'),
	});

	if (confirmed) {
		await deleteList(list);
	}
}

/**
 * Delete list
 * @param {any} list - List item to delete
 * @returns {Promise<void>}
 */
async function deleteList(list: any): Promise<void> {
	if (!apis.mapexOS?.lists) {
		notifyFail({ message: t.errors.apiNotInitialized.value });
		return;
	}

	if (!list.id) {
		notifyFail({ message: t.errors.idMissing.value });
		return;
	}

	try {
		await apis.mapexOS.lists.delete({ listId: list.id });

		// Update total items count after successful deletion
		totalItems.value = Math.max(0, totalItems.value - 1);

		// Remove from local list
		const index = listsList.value.findIndex(r => r.id === list.id);
		if (index !== -1) {
			listsList.value.splice(index, 1);
		}

		// If current page is now empty and not the first page, go to previous page
		if (listsList.value.length === 0 && currentPage.value > 1) {
			currentPage.value -= 1;
			await fetchLists();
		}

		notifySuccess({ message: t.notifications.deleted.value });
	} catch (err: any) {
		notifyFail({ message: err.message || t.notifications.deleteError.value });
	}
}

/**
 * Get actions configuration for a list
 * @param {any} list - List item
 * @returns {DataRowActionConfig} Actions configuration
 */
function getListActions(list: any): DataRowActionConfig {
	return {
		showEdit: canUpdateList.value && canModifyList(list),
		showView: canReadList.value,
		showDelete: canDeleteList.value && canModifyList(list),
	};
}

/** LIFECYCLE HOOKS */
onMounted(async () => await fetchLists());

// Auto-refresh lists when organization changes
useOrgChangeRefresh(async () => {
	// Reset pagination and filters when org changes
	currentPage.value = 1;
	quickSearch.value = '';
	quickSource.value = null;
	advancedFilterValues.value = {
		category: 'iot',
		type: null,
		isTemplate: null,
		includeChildren: null,
	};
	hasPendingAdvancedFilters.value = false;
	await fetchLists();
});
</script>

<template>
	<q-page class="q-pt-lg">

		<!-- Header Section -->
		<PageHeader
			icon="list"
			iconColor="primary"
			:title="t.pageHeader.title.value"
			:description="t.pageHeader.description.value"
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
					@keyup.enter="applyQuickFilters"
					@clear="applyQuickFilters"
				>
					<template #prepend>
						<q-icon name="search" color="grey-6" />
					</template>
				</q-input>
			</div>

			<!-- Source Select (System/Custom) -->
			<div class="col-auto" style="min-width: 140px;">
				<q-select
					v-model="quickSource"
					outlined
					dense
					emit-value
					map-options
					:options="[
						{ label: t.filters.allSource.value, value: null },
						{ label: t.filters.options.system.value, value: true },
						{ label: t.filters.options.custom.value, value: false },
					]"
					:label="t.filters.isSystem.value"
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
					<q-icon name="list" size="sm" color="primary" class="q-mr-sm"/>
					<div class="text-subtitle1 text-weight-medium text-primary">{{ t.listHeader.title.value }}</div>
				</div>
			</div>
			<div class="col-auto">
				<ListHeaderMenu
					icon="list"
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
					@refresh="fetchLists"
				/>
			</div>
		</div>

		<!-- Loading Spinner -->
		<div v-if="loading" class="row justify-center q-my-lg">
			<q-spinner color="primary" size="3em" />
		</div>

		<!-- Lists Row List -->
		<div v-else class="row">
			<div
				v-for="(list, index) in listsList"
				:key="list.id || `list-${index}`"
				class="col-12 q-mb-xs"
			>
				<DataRow
					:data="list"
					:columns="visibleColumns"
					:actions="getListActions(list)"
					@click="viewDetails"
					@dblclick="editList"
					@edit="editList"
					@view="viewDetails"
					@delete="confirmDelete"
				/>
			</div>

			<!-- No Results -->
			<div v-if="listsList.length === 0" class="col-12">
				<ListCardEmpty
					:title="t.empty.title.value"
					:description="t.empty.description.value"
					icon="list"
				/>
			</div>
		</div>

		<!-- Pagination -->
		<ListPagination
			v-model="currentPage"
			:total-pages="totalPages"
			@change="handlePageChange"
		/>

		<!-- List Details Drawer -->
		<ListDrawer
			v-if="selectedListId"
			v-model="drawerOpen"
			:listId="selectedListId"
		/>

		<!-- Advanced Filters Drawer -->
		<AdvancedFiltersDrawer
			v-model="showAdvancedFilters"
			:fields="advancedFilterFields"
			:values="advancedFilterValues"
			@apply="handleAdvancedFiltersApply"
			@reset="handleAdvancedFiltersReset"
			@pending-change="handlePendingChange"
		/>
	</q-page>
</template>
