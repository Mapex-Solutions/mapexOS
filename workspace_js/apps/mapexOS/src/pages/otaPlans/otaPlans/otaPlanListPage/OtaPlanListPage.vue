<script setup lang="ts">
defineOptions({
  name: 'OtaPlanListPage'
});

/** TYPE IMPORTS */
import type { OTAPlanResponse, OTAPlanStatus } from '@mapexos/schemas';
import type { DataRowColumn } from '@components/cards';
import type { OtaPlanListFilters } from './interfaces/otaPlanListPage.interface';

/** VUE IMPORTS */
import { ref, computed, onMounted } from 'vue';
import { useRouter } from 'vue-router';

/** COMPONENTS */
import { PageHeader, ListHeaderMenu } from '@components/headers';
import { ListCardEmpty, DataRow } from '@components/cards';
import { ListPagination } from '@components/navigation';

/** COMPOSABLES */
import { useOtaPlanListTranslations } from '@composables/i18n/pages/otaPlans/otaPlanList/useOtaPlanListTranslations';
import { useOrgChangeRefresh } from '@composables/organizations';
import { usePermissions } from '@composables/shared/usePermissions';
import { useLogger } from '@composables/useLogger';

/** UTILS */
import { notifyFail, notifySuccess, dialogDelete } from '@utils/alert';
import { cleanQueryParams } from '@utils/query';

/** SERVICES */
import { apis } from '@services/mapex';

/** LOCAL IMPORTS */
import { DEFAULT_ITEMS_PER_PAGE, OTA_PLAN_FILTER_DEFAULTS } from './constants/otaPlanListPage.constant';
import { OTA_PLAN_STATUS_COLORS, OTA_PLAN_STATUS_ICONS } from '../constants/otaStatus.constant';

/** COMPOSABLES & STORES */
const t = useOtaPlanListTranslations();
const router = useRouter();
const logger = useLogger('OtaPlanListPage');
const { canCreate, canDelete, canRead } = usePermissions();
const canCreatePlan = canCreate('ota_plans');
const canDeletePlan = canDelete('ota_plans');
const canReadPlan = canRead('ota_plans');

/** STATE */
const plansList = ref<OTAPlanResponse[]>([]);
const loading = ref(false);
const lastUpdatedAt = ref<number | undefined>(undefined);
const itemsPerPage = ref(DEFAULT_ITEMS_PER_PAGE);
const currentPage = ref(1);
const totalPages = ref(1);
const totalItems = ref(0);
const filters = ref<OtaPlanListFilters>({ ...OTA_PLAN_FILTER_DEFAULTS });
const quickSearch = ref('');
const quickStatus = ref<OTAPlanStatus | null>(null);

/** COMPUTED */

/**
 * Status options for the quick filter select, from the wire status enum
 */
const statusOptions = computed(() => [
  { label: t.filters.clearAll.value, value: null },
  ...Object.keys(OTA_PLAN_STATUS_COLORS).map((status) => ({
    label: t.planStatus(status),
    value: status,
  })),
]);

/**
 * DataRow column definitions: name, status chip, progress, window
 */
const columns = computed((): DataRowColumn[] => [
  {
    key: 'icon',
    label: '',
    type: 'avatar',
    visible: 'always',
    width: 56,
    icon: () => 'system_update',
    color: (_value: unknown, row: Record<string, unknown>) =>
      OTA_PLAN_STATUS_COLORS[String(row.status)] || 'grey-5',
  },
  {
    key: 'name',
    label: t.columns.name.value,
    type: 'text',
    visible: 'always',
    width: 220,
    ellipsis: true,
    secondaryKey: 'description',
  },
  {
    key: 'status',
    label: t.columns.status.value,
    type: 'chip',
    visible: 'always',
    width: 140,
    format: (value: string) => t.planStatus(value),
    color: (value: string) => OTA_PLAN_STATUS_COLORS[value] || 'grey-6',
    icon: (value: string) => OTA_PLAN_STATUS_ICONS[value] || 'help',
  },
  {
    key: 'counters',
    label: t.columns.progress.value,
    type: 'text',
    visible: 'laptop',
    width: 120,
    format: (value: OTAPlanResponse['counters']) =>
      value ? `${value.succeeded}/${value.total}` : '',
  },
  {
    key: 'startAt',
    label: t.columns.startAt.value,
    type: 'date',
    visible: 'laptop',
    width: 170,
  },
  {
    key: 'maxTime',
    label: t.columns.maxTime.value,
    type: 'date',
    visible: 'laptop',
    width: 170,
  },
]);

/** Whether any quick filter is active */
const hasActiveFilters = computed(() => !!(quickSearch.value || quickStatus.value));

/** FUNCTIONS */

/**
 * Apply quick filters and refetch from page 1
 */
function applyQuickFilters(): void {
  filters.value.name = quickSearch.value || undefined;
  filters.value.status = quickStatus.value ?? undefined;
  currentPage.value = 1;
  void fetchPlans();
}

/**
 * Clear all filters and refetch
 */
function clearAllFilters(): void {
  quickSearch.value = '';
  quickStatus.value = null;
  filters.value = { ...OTA_PLAN_FILTER_DEFAULTS };
  currentPage.value = 1;
  void fetchPlans();
}

/**
 * Fetch the plans page with current filters + pagination
 */
async function fetchPlans(): Promise<void> {
  if (!apis.assets) return;
  try {
    loading.value = true;

    const queryParams = cleanQueryParams({
      page: currentPage.value,
      perPage: itemsPerPage.value,
      name: filters.value.name,
      status: filters.value.status,
    });

    const response = await apis.assets.otaPlans.listPlans(queryParams);
    plansList.value = response.items;

    if (response.pagination) {
      totalPages.value = response.pagination.totalPages || 1;
      totalItems.value = response.pagination.totalItems || 0;
    }
  } catch (err: unknown) {
    logger.error('Error fetching OTA plans:', err);
    notifyFail({ message: t.messages.loadError.value });
  } finally {
    loading.value = false;
    lastUpdatedAt.value = Date.now();
  }
}

/**
 * Pagination navigation
 */
function handlePageChange(page: number): void {
  currentPage.value = page;
  void fetchPlans();
}

/**
 * Page-size change resets to the first page
 */
function handleItemsPerPageChange(newValue: number): void {
  itemsPerPage.value = newValue;
  currentPage.value = 1;
  void fetchPlans();
}

/**
 * Navigate to the plan detail page for a row
 */
function viewDetails(plan: OTAPlanResponse): void {
  if (!canReadPlan.value || !plan.id) return;
  void router.push(`/ota_plans/${plan.id}`);
}

/**
 * Confirm and cancel a plan (the platform close routine finalizes it)
 */
async function confirmCancel(plan: OTAPlanResponse): Promise<void> {
  if (!canDeletePlan.value || !plan.id) return;

  const confirmed = await dialogDelete({
    title: t.cancelDialog.title.value,
    message: t.cancelDialog.message.value,
  });
  if (!confirmed) return;

  try {
    await apis.assets.otaPlans.deletePlan({ planId: plan.id });
    notifySuccess({ message: t.messages.cancelSuccess.value });
    await fetchPlans();
  } catch (err: unknown) {
    logger.error('Error canceling OTA plan:', err);
    notifyFail({ message: t.messages.cancelError.value });
  }
}

/** LIFECYCLE HOOKS */

onMounted(async () => {
  await fetchPlans();
});

/**
 * Refetch when the selected organization changes so the list stays scoped
 */
useOrgChangeRefresh(async () => {
  currentPage.value = 1;
  quickSearch.value = '';
  quickStatus.value = null;
  filters.value = { ...OTA_PLAN_FILTER_DEFAULTS };
  await fetchPlans();
});
</script>

<template>
  <q-page class="q-pt-lg">
    <!-- Header Section -->
    <PageHeader
      icon="system_update"
      iconColor="primary"
      :title="t.page.title.value"
      :description="t.page.description.value"
      :info="t.page.info.value"
      :button="canCreatePlan ? { label: t.page.addButton.value, icon: 'add', to: '/ota_plans/add', color: 'primary' } : undefined"
    />

    <!-- Filters Section -->
    <div class="text-caption text-grey-7 q-mb-xs">{{ t.filters.label.value }}</div>
    <div class="row items-center q-col-gutter-sm q-mb-md">
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

      <div class="col-auto" style="min-width: 170px;">
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

    <!-- Results Section -->
    <div class="row items-center q-pt-xl q-mb-md">
      <div class="col">
        <div class="row items-center">
          <q-icon name="system_update" size="sm" color="primary" class="q-mr-sm" />
          <div class="text-subtitle1 text-weight-medium text-primary">{{ t.page.listTitle.value }}</div>
        </div>
      </div>
      <div class="col-auto">
        <ListHeaderMenu
          icon="system_update"
          :items-count="totalItems"
          :item-label="t.page.itemLabel.value"
          :item-label-plural="t.page.itemLabelPlural.value"
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

    <!-- Plans Row List -->
    <div v-else class="row">
      <div
        v-for="(plan, index) in plansList"
        :key="plan.id || `plan-${index}`"
        class="col-12 q-mb-xs"
      >
        <DataRow
          :data="plan"
          :columns="columns"
          :actions="{ showView: canReadPlan, showDelete: canDeletePlan }"
          @click="viewDetails"
          @view="viewDetails"
          @delete="confirmCancel"
        />
      </div>

      <!-- No Results -->
      <div v-if="plansList.length === 0" class="col-12">
        <ListCardEmpty
          :title="t.empty.title.value"
          :description="t.empty.description.value"
          icon="system_update"
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
