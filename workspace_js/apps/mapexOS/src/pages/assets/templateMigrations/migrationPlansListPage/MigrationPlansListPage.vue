<script setup lang="ts">
defineOptions({
  name: 'MigrationPlansListPage'
});

/** TYPE IMPORTS */
import type { DataRowColumn, DataRowActionConfig } from '@components/cards';
import type { MigrationPlanResponse, MigrationPlanStatus } from '@mapexos/schemas';
import type { MigrationPlansListPageFilters, MigrationStatusOption } from './interfaces';

/** VUE IMPORTS */
import { ref, computed, onMounted } from 'vue';
import { useRouter } from 'vue-router';

/** COMPONENTS */
import { PageHeader, ListHeaderMenu } from '@components/headers';
import { ListCardEmpty, DataRow } from '@components/cards';
import { ListPagination } from '@components/navigation';

/** COMPOSABLES */
import { useMigrationPlansTranslations } from '@composables/i18n';
import { useOrgChangeRefresh } from '@composables/organizations';
import { usePermissions } from '@composables/shared/usePermissions';
import { useLogger } from '@composables/useLogger';

/** UTILS */
import { date } from 'quasar';
import { notifyFail, notifySuccess, dialogDelete } from '@utils/alert';
import { cleanQueryParams } from '@utils/query';

/** SERVICES */
import { apis } from '@services/mapex';

/** LOCAL IMPORTS */
import {
  DEFAULT_ITEMS_PER_PAGE,
  MIGRATION_FILTER_DEFAULTS,
  MIGRATION_STATUS_VALUES,
  MIGRATION_STATUS_COLORS,
  MIGRATION_STATUS_ICONS,
  MIGRATION_EDITABLE_STATUSES,
  MIGRATION_SCHEDULE_MASK,
  MIGRATION_ROUTE_BASE,
} from './constants';

/** COMPOSABLES & STORES */
const t = useMigrationPlansTranslations();
const router = useRouter();
const logger = useLogger('MigrationPlansListPage');
const { canCreate, canRead, canUpdate, canDelete } = usePermissions();
const canCreatePlan = canCreate('templatemigrations');
const canReadPlan = canRead('templatemigrations');
const canUpdatePlan = canUpdate('templatemigrations');
const canDeletePlan = canDelete('templatemigrations');

/** STATE */
const plansList = ref<MigrationPlanResponse[]>([]);
const loading = ref(false);
const error = ref<string | undefined>(undefined);
const lastUpdatedAt = ref<number | undefined>(undefined);
const itemsPerPage = ref(DEFAULT_ITEMS_PER_PAGE);
const currentPage = ref(1);
const totalPages = ref(1);
const totalItems = ref(0);
const filters = ref<MigrationPlansListPageFilters>({ ...MIGRATION_FILTER_DEFAULTS });

const quickSearch = ref('');
const quickStatus = ref<MigrationPlanStatus | null>(null);

/** COMPUTED */

/**
 * Status options for the quick filter select
 */
const statusOptions = computed((): MigrationStatusOption[] => [
  { label: t.filters.allStatus.value, value: null },
  ...MIGRATION_STATUS_VALUES.map(status => ({
    label: t.statusLabel(status),
    value: status,
  })),
]);

/**
 * Whether any quick filter is active
 */
const hasActiveFilters = computed(() => !!quickSearch.value || quickStatus.value !== null);

/**
 * DataRow column definitions for the migration plans list
 */
const columns = computed((): DataRowColumn[] => [
  {
    key: 'icon',
    label: '',
    type: 'avatar',
    visible: 'always',
    width: 56,
    icon: (_value: any, row: any) => MIGRATION_STATUS_ICONS[row.status as MigrationPlanStatus] ?? 'swap_horiz',
    color: (_value: any, row: any) => MIGRATION_STATUS_COLORS[row.status as MigrationPlanStatus] ?? 'primary',
    tooltip: (_value: any, row: any) => t.statusLabel(row.status),
  },
  {
    key: 'name',
    label: t.columns.name.value,
    type: 'text',
    visible: 'always',
    width: 240,
    ellipsis: true,
    secondaryKey: 'description',
  },
  {
    key: 'fromTemplateName',
    label: t.columns.route.value,
    type: 'text',
    visible: 'laptop',
    width: 260,
    ellipsis: true,
    format: (value: any, row: any) => value || row.fromTemplateId,
    secondaryKey: 'toTemplateName',
  },
  {
    key: 'scheduleAt',
    label: t.columns.schedule.value,
    type: 'text',
    visible: 'laptop',
    width: 180,
    format: (value: any, row: any) => formatSchedule(value, row.created),
  },
  {
    key: 'status',
    label: t.columns.status.value,
    type: 'chip',
    visible: 'always',
    width: 160,
    align: 'center',
    format: (value: any) => t.statusLabel(value).toUpperCase(),
    color: (value: any) => MIGRATION_STATUS_COLORS[value as MigrationPlanStatus] ?? 'grey-6',
    icon: (value: any) => MIGRATION_STATUS_ICONS[value as MigrationPlanStatus] ?? 'help',
  },
  {
    key: 'progressPct',
    label: t.columns.progress.value,
    type: 'text',
    visible: 'laptop',
    width: 120,
    align: 'center',
    format: (value: any) => `${Math.round(Number(value) || 0)}%`,
  },
]);

/**
 * Row action menu config. The DataRow built-in Edit/View/Delete are disabled
 * (they are unconditional and English-only); the menu is driven entirely by
 * these i18n custom actions. View is always available; Edit and Remove appear
 * only while the plan is still editable (scheduled/pending), mirroring the
 * backend guard.
 */
const rowActions = computed((): DataRowActionConfig => ({
  showEdit: false,
  showView: false,
  showDelete: false,
  customActions: [
    ...(canReadPlan.value
      ? [{
          key: 'view',
          label: t.actions.view.value,
          description: t.actions.viewHint.value,
          icon: 'visibility',
          color: 'blue-7',
        }]
      : []),
    ...(canUpdatePlan.value
      ? [{
          key: 'edit',
          label: t.actions.edit.value,
          description: t.actions.editHint.value,
          icon: 'edit',
          color: 'primary',
          condition: (row: any) => MIGRATION_EDITABLE_STATUSES.includes(row.status),
        }]
      : []),
    ...(canDeletePlan.value
      ? [{
          key: 'remove',
          label: t.actions.remove.value,
          description: t.actions.removeHint.value,
          icon: 'cancel',
          color: 'negative',
          condition: (row: any) => MIGRATION_EDITABLE_STATUSES.includes(row.status),
        }]
      : []),
  ],
}));

/** FUNCTIONS */

/**
 * Format the schedule cell: shows the localized "immediate" label when the
 * plan is not scheduled for the future relative to its creation time.
 * @param scheduleAt - ISO schedule timestamp
 * @param created - ISO creation timestamp
 * @returns Formatted date/time or the immediate label
 */
function formatSchedule(scheduleAt: string | undefined, created: string | undefined): string {
  if (!scheduleAt) {
    return t.schedule.immediate.value;
  }

  const scheduleTime = new Date(scheduleAt).getTime();
  const createdTime = created ? new Date(created).getTime() : Date.now();

  if (Number.isNaN(scheduleTime) || scheduleTime <= createdTime) {
    return t.schedule.immediate.value;
  }

  return date.formatDate(scheduleAt, MIGRATION_SCHEDULE_MASK);
}

/**
 * Apply the quick filters and refetch from the first page
 * @returns {void}
 */
function applyQuickFilters(): void {
  filters.value.name = quickSearch.value || undefined;
  filters.value.status = quickStatus.value ?? undefined;
  currentPage.value = 1;
  void fetchPlans();
}

/**
 * Clear all quick filters and refetch
 * @returns {void}
 */
function clearAllFilters(): void {
  quickSearch.value = '';
  quickStatus.value = null;
  filters.value = { ...MIGRATION_FILTER_DEFAULTS };
  currentPage.value = 1;
  void fetchPlans();
}

/**
 * Fetch migration plans from the API with current filters and pagination
 * @returns {Promise<void>}
 */
async function fetchPlans(): Promise<void> {
  if (!apis.assets) {
    error.value = t.errors.apiNotInitialized.value;
    return;
  }

  try {
    loading.value = true;
    error.value = undefined;

    const queryParams: Record<string, any> = {
      page: currentPage.value,
      perPage: itemsPerPage.value,
    };

    if (filters.value.name) {
      queryParams.name = filters.value.name;
    }
    if (filters.value.status) {
      queryParams.status = filters.value.status;
    }

    const cleanedParams = cleanQueryParams(queryParams);
    const response = await apis.assets.migration.list(cleanedParams);

    plansList.value = response?.items || [];

    if (response.pagination) {
      totalItems.value = response.pagination.totalItems || 0;
      totalPages.value = response.pagination.totalPages || 1;
    }
  } catch (err: any) {
    logger.error('Error fetching migration plans:', err);
    const errorMsg = err.message || t.notifications.loadError.value;
    error.value = errorMsg;
    notifyFail({ message: t.notifications.loadError.value });
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
  void fetchPlans();
}

/**
 * Handle items-per-page change from the list header menu
 * @param {number} newValue - New items per page value
 * @returns {void}
 */
function handleItemsPerPageChange(newValue: number): void {
  itemsPerPage.value = newValue;
  currentPage.value = 1;
  void fetchPlans();
}

/**
 * Navigate to the migration plan detail page
 * @param {MigrationPlanResponse} plan - Selected plan
 * @returns {void}
 */
function viewDetails(plan: MigrationPlanResponse): void {
  if (!plan.id) return;
  void router.push(`${MIGRATION_ROUTE_BASE}/${plan.id}`);
}

/**
 * Dispatch a row action-menu event to its handler.
 * @param {string} key - Action key
 * @param {MigrationPlanResponse} plan - Row plan
 * @returns {void}
 */
function handleRowAction(key: string, plan: MigrationPlanResponse): void {
  if (key === 'view') {
    viewDetails(plan);
  } else if (key === 'edit') {
    if (plan.id) void router.push(`${MIGRATION_ROUTE_BASE}/edit/${plan.id}`);
  } else if (key === 'remove') {
    void confirmRemovePlan(plan);
  }
}

/**
 * Confirm and remove (cancel) a scheduled plan, then refresh the list.
 * @param {MigrationPlanResponse} plan - Plan to remove
 * @returns {Promise<void>}
 */
async function confirmRemovePlan(plan: MigrationPlanResponse): Promise<void> {
  if (!plan.id) return;

  const confirmed = await dialogDelete({
    title: t.cancelDialog.title.value,
    message: t.cancelDialog.message.value,
    ok: { label: t.cancelDialog.confirm.value, color: 'negative' },
  });
  if (!confirmed) return;

  if (!apis.assets) {
    notifyFail({ message: t.errors.apiNotInitialized.value });
    return;
  }

  try {
    await apis.assets.migration.cancel({ id: plan.id });
    notifySuccess({ message: t.cancelDialog.success.value });
    await fetchPlans();
  } catch (err: any) {
    logger.error('Error removing migration plan:', err);
    notifyFail({ message: t.cancelDialog.error.value });
  }
}

/** LIFECYCLE HOOKS */

onMounted(async () => {
  await fetchPlans();
});

/**
 * Refetch when the active organization changes
 */
useOrgChangeRefresh(async () => {
  currentPage.value = 1;
  quickSearch.value = '';
  quickStatus.value = null;
  filters.value = { ...MIGRATION_FILTER_DEFAULTS };
  await fetchPlans();
});
</script>

<template>
  <q-page class="q-pt-lg">

    <!-- Header Section -->
    <PageHeader
      icon="swap_horiz"
      iconColor="primary"
      :title="t.pageHeader.title.value"
      :description="t.pageHeader.description.value"
      :button="canCreatePlan ? { label: t.pageHeader.button.value, icon: 'add', to: `${MIGRATION_ROUTE_BASE}/add`, color: 'primary' } : undefined"
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
      <div class="col-auto" style="min-width: 180px;">
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

      <!-- Clear All Button -->
      <div v-if="hasActiveFilters" class="col-auto">
        <q-btn
          flat
          dense
          size="sm"
          color="grey-7"
          icon="filter_alt_off"
          :label="t.filters.clearAll.value"
          no-caps
          @click="clearAllFilters"
        />
      </div>
    </div>

    <!-- Results Header -->
    <div class="row items-center q-pt-xl q-mb-md">
      <div class="col">
        <div class="row items-center">
          <q-icon name="swap_horiz" size="sm" color="primary" class="q-mr-sm"/>
          <div class="text-subtitle1 text-weight-medium text-primary">{{ t.listHeader.title.value }}</div>
        </div>
      </div>
      <div class="col-auto">
        <ListHeaderMenu
          icon="swap_horiz"
          :items-count="totalItems"
          :item-label="t.listHeader.itemLabel.value"
          :item-label-plural="t.listHeader.itemLabelPlural.value"
          :items-per-page="itemsPerPage"
          :filtered="hasActiveFilters"
          :refreshing="loading"
          :last-updated-at="lastUpdatedAt"
          @update:items-per-page="handleItemsPerPageChange"
          @refresh="fetchPlans"
        />
      </div>
    </div>

    <!-- Loading Spinner -->
    <div v-if="loading" class="row justify-center q-my-lg">
      <q-spinner color="primary" size="3em" />
    </div>

    <!-- Migration Plans Row List -->
    <div v-else class="row">
      <div
        v-for="(plan, index) in plansList"
        :key="plan.id || `plan-${index}`"
        class="col-12 q-mb-xs"
      >
        <DataRow
          :data="plan"
          :columns="columns"
          :show-actions="true"
          :actions="rowActions"
          @click="viewDetails"
          @action="handleRowAction"
        />
      </div>

      <!-- Empty State -->
      <div v-if="plansList.length === 0" class="col-12">
        <ListCardEmpty
          :title="t.empty.title.value"
          :description="t.empty.description.value"
          icon="swap_horiz"
        />
      </div>
    </div>

    <!-- Pagination -->
    <ListPagination
      v-model="currentPage"
      :total-pages="totalPages"
      @change="handlePageChange"
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
