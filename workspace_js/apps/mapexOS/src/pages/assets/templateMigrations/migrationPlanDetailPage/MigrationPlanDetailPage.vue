<script setup lang="ts">
defineOptions({
  name: 'MigrationPlanDetailPage'
});

/** TYPE IMPORTS */
import type { QTableColumn } from 'quasar';
import type {
  MigrationPlanResponse,
  MigrationExecutionResponse,
  MigrationExecutionStatus,
} from '@mapexos/schemas';
import type { DetailChipColor } from '@components/chips';
import type { MigrationExecutionStatusOption } from './interfaces';

/** VUE IMPORTS */
import { ref, computed, onMounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';

/** COMPONENTS */
import { PageHeader } from '@components/headers';
import { ListPagination } from '@components/navigation';
import { BaseButton } from '@components/buttons';
import { DetailChip } from '@components/chips';
import { ListCardEmpty } from '@components/cards';
import { AppTooltip } from '@components/tooltips';

/** COMPOSABLES */
import { useMigrationPlanDetailTranslations } from '@composables/i18n';
import { useLogger } from '@composables/useLogger';

/** UTILS */
import { date } from 'quasar';
import { notifyFail, notifySuccess, dialogDelete } from '@utils/alert';
import { cleanQueryParams } from '@utils/query';

/** SERVICES */
import { apis } from '@services/mapex';

/** LOCAL IMPORTS */
import {
  DEFAULT_EXECUTIONS_PER_PAGE,
  MIGRATION_EDITABLE_STATUSES,
  EXECUTION_STATUS_VALUES,
  EXECUTION_STATUS_COLORS,
  EXECUTION_STATUS_ICONS,
  PLAN_STATUS_DETAIL_COLORS,
} from './constants';
import {
  MIGRATION_STATUS_ICONS,
  MIGRATION_SCHEDULE_MASK,
  MIGRATION_ROUTE_BASE,
} from '../migrationPlansListPage/constants';

/** COMPOSABLES & STORES */
const t = useMigrationPlanDetailTranslations();
const route = useRoute();
const router = useRouter();
const logger = useLogger('MigrationPlanDetailPage');

/** STATE */
const plan = ref<MigrationPlanResponse | null>(null);
const loading = ref(false);
const cancelling = ref(false);

const executions = ref<MigrationExecutionResponse[]>([]);
const executionsLoading = ref(false);
const execStatusFilter = ref<MigrationExecutionStatus | null>(null);
const currentPage = ref(1);
const totalPages = ref(1);
const totalItems = ref(0);

/** Resolved asset names keyed by asset id (falls back to the id) */
const assetNames = ref<Map<string, string>>(new Map());

/** COMPUTED */

/** Plan id from the route */
const planId = computed(() => route.params.id as string);

/** Whether edit/cancel actions are allowed for the current plan status */
const canModify = computed(() =>
  !!plan.value && MIGRATION_EDITABLE_STATUSES.includes(plan.value.status)
);

/** DetailChip color for the current plan status */
const statusDetailColor = computed<DetailChipColor>(() =>
  plan.value ? PLAN_STATUS_DETAIL_COLORS[plan.value.status] ?? 'default' : 'default'
);

/** Status chip icon for the current plan */
const statusIcon = computed(() =>
  plan.value ? MIGRATION_STATUS_ICONS[plan.value.status] ?? 'help' : 'help'
);

/** Progress fraction (0..1) for the linear progress bar */
const progressValue = computed(() => {
  const pct = Number(plan.value?.progressPct) || 0;
  return Math.min(Math.max(pct / 100, 0), 1);
});

/** Rounded progress percentage for display */
const progressPct = computed(() => Math.round(Number(plan.value?.progressPct) || 0));

/** Options for the executions status filter select */
const executionStatusOptions = computed((): MigrationExecutionStatusOption[] => [
  { label: t.executions.allStatus.value, value: null },
  ...EXECUTION_STATUS_VALUES.map(status => ({
    label: t.executionStatusLabel(status),
    value: status,
  })),
]);

/** Executions table column definitions */
const executionColumns = computed((): QTableColumn<MigrationExecutionResponse>[] => [
  { name: 'assetId', label: t.executions.columns.assetId.value, field: 'assetId', align: 'left' },
  { name: 'status', label: t.executions.columns.status.value, field: 'status', align: 'center' },
  { name: 'error', label: t.executions.columns.error.value, field: 'error', align: 'left' },
  { name: 'attempts', label: t.executions.columns.attempts.value, field: 'attempts', align: 'center' },
  {
    name: 'updated',
    label: t.executions.columns.updated.value,
    field: 'updated',
    align: 'left',
    format: (val: string) => formatDateTime(val),
  },
]);

/** FUNCTIONS */

/**
 * Format an ISO timestamp for display, tolerating empty values.
 * @param value - ISO timestamp
 * @returns Formatted date/time or an em dash
 */
function formatDateTime(value: string | undefined): string {
  if (!value) return '—';
  return date.formatDate(value, MIGRATION_SCHEDULE_MASK);
}

/**
 * Format the plan schedule, showing the immediate label when the plan is not
 * scheduled for the future relative to its creation time.
 * @returns Formatted schedule or the immediate label
 */
function formatSchedule(): string {
  const scheduleAt = plan.value?.scheduleAt;
  const created = plan.value?.created;
  if (!scheduleAt) return t.header.immediate.value;

  const scheduleTime = new Date(scheduleAt).getTime();
  const createdTime = created ? new Date(created).getTime() : Date.now();
  if (Number.isNaN(scheduleTime) || scheduleTime <= createdTime) {
    return t.header.immediate.value;
  }
  return date.formatDate(scheduleAt, MIGRATION_SCHEDULE_MASK);
}

/** DetailChip color for an execution status */
function executionColor(status: MigrationExecutionStatus): DetailChipColor {
  return EXECUTION_STATUS_COLORS[status] ?? 'grey';
}

/** Chip icon for an execution status */
function executionIcon(status: MigrationExecutionStatus): string {
  return EXECUTION_STATUS_ICONS[status] ?? 'help';
}

/** Resolved asset name for an id, falling back to the id itself */
function assetName(id: string): string {
  return assetNames.value.get(id) || id;
}

/**
 * Fetch and cache asset names for any uncached ids in the given executions.
 * @param rows - Executions whose asset ids should be resolved
 * @returns {Promise<void>}
 */
async function resolveAssetNames(rows: MigrationExecutionResponse[]): Promise<void> {
  if (!apis.assets) return;

  const missing = [
    ...new Set(rows.map(row => row.assetId).filter(id => id && !assetNames.value.has(id))),
  ];
  if (missing.length === 0) return;

  const resolved = await Promise.all(
    missing.map(async id => {
      try {
        const asset = await apis.assets.asset.getById({ assetId: id });
        return [id, asset?.name || id] as const;
      } catch {
        return [id, id] as const;
      }
    })
  );

  const next = new Map(assetNames.value);
  for (const [id, name] of resolved) next.set(id, name);
  assetNames.value = next;
}

/**
 * Fetch the migration plan from the API.
 * @returns {Promise<void>}
 */
async function fetchPlan(): Promise<void> {
  if (!apis.assets) {
    notifyFail({ message: t.errors.apiNotInitialized.value });
    return;
  }
  if (!planId.value) {
    notifyFail({ message: t.page.notFound.value });
    void router.back();
    return;
  }

  try {
    loading.value = true;
    plan.value = await apis.assets.migration.get({ id: planId.value });
  } catch (err: any) {
    logger.error('Error fetching migration plan:', err);
    notifyFail({ message: t.page.loadError.value });
  } finally {
    loading.value = false;
  }
}

/**
 * Fetch the plan executions with the current filter and pagination.
 * @returns {Promise<void>}
 */
async function fetchExecutions(): Promise<void> {
  if (!apis.assets || !planId.value) return;

  try {
    executionsLoading.value = true;

    const query = cleanQueryParams({
      page: currentPage.value,
      perPage: DEFAULT_EXECUTIONS_PER_PAGE,
      status: execStatusFilter.value ?? undefined,
    });

    const response = await apis.assets.migration.executions({ id: planId.value }, query);
    executions.value = response?.items || [];
    await resolveAssetNames(executions.value);

    if (response?.pagination) {
      totalItems.value = response.pagination.totalItems || 0;
      totalPages.value = response.pagination.totalPages || 1;
    }
  } catch (err: any) {
    logger.error('Error fetching migration executions:', err);
    notifyFail({ message: t.page.loadError.value });
  } finally {
    executionsLoading.value = false;
  }
}

/**
 * Refresh both the plan and its executions (used by polling and post-action refresh).
 * @returns {Promise<void>}
 */
async function refreshAll(): Promise<void> {
  await Promise.all([fetchPlan(), fetchExecutions()]);
}

/**
 * Apply the executions status filter and refetch from the first page.
 * @returns {void}
 */
function applyExecutionFilter(): void {
  currentPage.value = 1;
  void fetchExecutions();
}

/**
 * Header refresh button: clear the status filter back to ALL and reload everything.
 * @returns {Promise<void>}
 */
async function resetAndRefresh(): Promise<void> {
  execStatusFilter.value = null;
  currentPage.value = 1;
  await refreshAll();
}

/**
 * Handle executions pagination navigation.
 * @param {number} page - New page number
 * @returns {void}
 */
function handlePageChange(page: number): void {
  currentPage.value = page;
  void fetchExecutions();
}

/**
 * Navigate to the create/edit page in edit mode for this plan.
 * @returns {void}
 */
function goToEdit(): void {
  if (!planId.value) return;
  void router.push(`${MIGRATION_ROUTE_BASE}/edit/${planId.value}`);
}

/**
 * Confirm and cancel the plan, then refresh on success.
 * @returns {Promise<void>}
 */
async function confirmCancel(): Promise<void> {
  const confirmed = await dialogDelete({
    title: t.cancelDialog.title.value,
    message: t.cancelDialog.message.value,
    ok: { label: t.cancelDialog.confirm.value, color: 'negative' },
  });
  if (!confirmed) return;

  if (!apis.assets || !planId.value) return;

  try {
    cancelling.value = true;
    await apis.assets.migration.cancel({ id: planId.value });
    notifySuccess({ message: t.cancelDialog.success.value });
    await refreshAll();
  } catch (err: any) {
    logger.error('Error cancelling migration plan:', err);
    notifyFail({ message: t.cancelDialog.error.value });
  } finally {
    cancelling.value = false;
  }
}

/** LIFECYCLE HOOKS */

onMounted(async () => {
  await refreshAll();
});
</script>

<template>
  <q-page class="q-pa-lg">
    <!-- Header Section -->
    <PageHeader
      icon="swap_horiz"
      iconColor="primary"
      :title="plan?.name || ''"
      :description="plan?.description || ''"
      :button="{ label: t.page.back.value, icon: 'arrow_back', flat: true, to: MIGRATION_ROUTE_BASE }"
    />

    <!-- Loading Spinner -->
    <div v-if="loading && !plan" class="row justify-center q-my-xl">
      <q-spinner color="primary" size="3em" />
    </div>

    <div v-else-if="plan" class="row q-col-gutter-lg q-mb-lg">
      <!-- LEFT: plan information panel -->
      <div class="col-12 col-md-4">
        <q-card class="rounded-borders">
          <!-- Panel header -->
          <q-card-section class="bg-grey-1 q-pb-md">
            <div class="text-h6 text-weight-bold text-primary">
              <q-icon size="sm" name="swap_horiz" color="primary" class="q-mr-xs" />
              {{ t.sections.overview.value }}
            </div>
            <div class="text-caption text-grey-7">{{ t.sections.overviewSubtitle.value }}</div>
          </q-card-section>

          <!-- Panel body -->
          <q-card-section class="q-pa-lg">
            <!-- Section: Status -->
            <div class="info-section">
              <div class="section-title">
                <q-icon name="flag" size="xs" class="q-mr-xs" />{{ t.sections.status.value }}
              </div>
              <div class="info-field q-mb-none">
                <div class="info-label">{{ t.fields.status.value }}</div>
                <DetailChip
                  dense
                  size="sm"
                  :color="statusDetailColor"
                  :icon="statusIcon"
                  :label="t.planStatusLabel(plan.status)"
                />
              </div>
            </div>

            <!-- Section: Templates -->
            <div class="info-section">
              <div class="section-title">
                <q-icon name="dashboard_customize" size="xs" class="q-mr-xs" />{{ t.sections.templates.value }}
              </div>
              <div class="info-field">
                <div class="info-label">{{ t.fields.fromTemplate.value }}</div>
                <div class="info-value">{{ plan.fromTemplateName || plan.fromTemplateId }}</div>
              </div>
              <div class="info-field q-mb-none">
                <div class="info-label">{{ t.fields.toTemplate.value }}</div>
                <div class="info-value">{{ plan.toTemplateName || plan.toTemplateId }}</div>
              </div>
            </div>

            <!-- Section: Schedule -->
            <div class="info-section">
              <div class="section-title">
                <q-icon name="schedule" size="xs" class="q-mr-xs" />{{ t.sections.schedule.value }}
              </div>
              <div class="row q-col-gutter-md">
                <div class="col-6">
                  <div class="info-field q-mb-none">
                    <div class="info-label">{{ t.header.schedule.value }}</div>
                    <div class="info-value">{{ formatSchedule() }}</div>
                  </div>
                </div>
                <div class="col-6">
                  <div class="info-field q-mb-none">
                    <div class="info-label">{{ t.header.created.value }}</div>
                    <div class="info-value">{{ formatDateTime(plan.created) }}</div>
                  </div>
                </div>
              </div>
            </div>

            <q-separator class="q-my-md" />

            <!-- Progress + counters -->
            <div>
              <div class="row items-center q-mb-sm">
                <span class="section-title q-mb-none">{{ t.progress.title.value }}</span>
                <q-space />
                <span class="text-body2 text-weight-medium text-primary">{{ progressPct }}%</span>
              </div>
              <q-linear-progress
                :value="progressValue"
                rounded
                size="10px"
                color="primary"
                track-color="grey-3"
              />
              <div class="row q-gutter-x-md q-mt-sm">
                <div class="text-caption text-grey-7">
                  {{ t.progress.total.value }}: <span class="text-body2 text-weight-medium">{{ plan.total }}</span>
                </div>
                <div class="text-caption text-grey-7">
                  {{ t.progress.migrated.value }}: <span class="text-body2 text-positive text-weight-medium">{{ plan.migrated }}</span>
                </div>
                <div class="text-caption text-grey-7">
                  {{ t.progress.failed.value }}: <span class="text-body2 text-negative text-weight-medium">{{ plan.failed }}</span>
                </div>
              </div>
            </div>
          </q-card-section>
        </q-card>
      </div>

      <!-- RIGHT: executions panel -->
      <div class="col-12 col-md-8">
        <q-card class="rounded-borders">
          <!-- Panel header -->
          <q-card-section class="bg-grey-1 q-pb-md">
            <div class="row items-center">
              <div class="col">
                <div class="text-h6 text-weight-bold text-primary">
                  <q-icon size="sm" name="format_list_bulleted" color="primary" class="q-mr-xs" />
                  {{ t.executions.title.value }}
                </div>
                <div class="text-caption text-grey-7">
                  {{ totalItems }} {{ totalItems === 1 ? t.executions.itemLabel.value : t.executions.itemLabelPlural.value }}
                </div>
              </div>
              <div class="col-auto" style="min-width: 200px;">
                <q-select
                  v-model="execStatusFilter"
                  outlined
                  dense
                  emit-value
                  map-options
                  :options="executionStatusOptions"
                  :label="t.executions.filterStatus.value"
                  class="filter-input"
                  @update:model-value="applyExecutionFilter"
                />
              </div>
              <div class="col-auto">
                <q-btn
                  flat
                  round
                  dense
                  color="primary"
                  icon="refresh"
                  :loading="executionsLoading || loading"
                  @click="refreshAll"
                >
                  <AppTooltip :content="t.actions.refresh.value" />
                </q-btn>
              </div>
            </div>
          </q-card-section>

          <q-separator />

          <!-- Panel body -->
          <q-card-section class="q-pa-lg">
            <!-- Loading (initial) -->
            <div v-if="executionsLoading && executions.length === 0" class="row justify-center q-py-xl">
              <q-spinner color="primary" size="2.5em" />
            </div>

            <!-- Empty: no table chrome, just the standard empty card -->
            <ListCardEmpty
              v-else-if="executions.length === 0"
              icon="devices"
              :title="t.executions.emptyTitle.value"
              :description="t.executions.empty.value"
              :button-label="t.executions.resetFilter.value"
              button-icon="filter_alt_off"
              @button-click="resetAndRefresh"
            />

            <!-- Table + pagination -->
            <template v-else>
              <q-table
                flat
                :rows="executions"
                :columns="executionColumns"
                row-key="id"
                :loading="executionsLoading"
                hide-pagination
                :rows-per-page-options="[0]"
              >
                <template #body-cell-assetId="props">
                  <q-td :props="props">
                    {{ assetName(props.row.assetId) }}
                  </q-td>
                </template>

                <template #body-cell-status="props">
                  <q-td :props="props" class="text-center">
                    <DetailChip
                      dense
                      size="sm"
                      :color="executionColor(props.row.status)"
                      :icon="executionIcon(props.row.status)"
                      :label="t.executionStatusLabel(props.row.status)"
                    />
                  </q-td>
                </template>

                <template #body-cell-error="props">
                  <q-td :props="props">
                    <span v-if="props.row.status === 'failed' && props.row.error" class="text-negative">
                      {{ props.row.error }}
                    </span>
                    <span v-else class="text-grey-6">—</span>
                  </q-td>
                </template>
              </q-table>

              <ListPagination
                v-model="currentPage"
                :total-pages="totalPages"
                @change="handlePageChange"
              />
            </template>

            <q-separator class="q-my-md" />

            <!-- Actions -->
            <div class="row justify-end q-gutter-sm">
              <BaseButton
                outline
                no-caps
                color="primary"
                icon="edit"
                :label="t.actions.edit.value"
                :disable="!canModify"
                @click="goToEdit"
              />
              <BaseButton
                outline
                no-caps
                color="negative"
                icon="cancel"
                :label="t.actions.cancel.value"
                :disable="!canModify"
                :loading="cancelling"
                @click="confirmCancel"
              />
            </div>
          </q-card-section>
        </q-card>
      </div>
    </div>

    <!-- Not Found -->
    <div v-else class="row justify-center q-my-xl text-grey-7">
      {{ t.page.notFound.value }}
    </div>
  </q-page>
</template>

<style lang="scss" scoped>
.info-section {
  margin-bottom: var(--mapex-spacing-lg);

  &:last-child {
    margin-bottom: 0;
  }
}

.section-title {
  display: flex;
  align-items: center;
  font-size: 0.7rem;
  font-weight: 600;
  letter-spacing: 0.05em;
  text-transform: uppercase;
  color: var(--mapex-text-secondary);
  margin-bottom: var(--mapex-spacing-sm);
}

.info-field {
  margin-bottom: var(--mapex-spacing-md);
}

.info-label {
  font-size: 0.75rem;
  color: var(--mapex-text-secondary);
  margin-bottom: 4px;
}

.info-value {
  font-size: 0.95rem;
  font-weight: 400;
  color: var(--mapex-text-primary);
  word-break: break-word;
}

.filter-input {
  :deep(.q-field__control) {
    border-radius: var(--mapex-radius-md);
  }
}
</style>
