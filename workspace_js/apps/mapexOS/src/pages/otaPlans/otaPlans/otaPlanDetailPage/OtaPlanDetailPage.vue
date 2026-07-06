<script setup lang="ts">
defineOptions({
  name: 'OtaPlanDetailPage'
});

/** TYPE IMPORTS */
import type { QTableColumn } from 'quasar';
import type { OTAPlanResponse, OTAExecutionResponse, OTAExecutionState } from '@mapexos/schemas';
import type { DetailChipColor } from '@components/chips';
import type { OtaExecutionStateOption } from './interfaces/otaPlanDetailPage.interface';

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
import { useOtaPlanDetailTranslations } from '@composables/i18n/pages/otaPlans/otaPlanDetail/useOtaPlanDetailTranslations';
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
  OTA_ROUTE_BASE,
  OTA_DATETIME_MASK,
} from './constants/otaPlanDetailPage.constant';
import {
  OTA_PLAN_STATUS_CHIP_COLORS,
  OTA_PLAN_STATUS_ICONS,
  OTA_EXECUTION_STATE_COLORS,
  OTA_EXECUTION_STATE_ICONS,
  OTA_PLAN_CANCELABLE_STATUSES,
} from '../constants/otaStatus.constant';

/** COMPOSABLES & STORES */
const t = useOtaPlanDetailTranslations();
const route = useRoute();
const router = useRouter();
const logger = useLogger('OtaPlanDetailPage');

/** STATE */
const plan = ref<OTAPlanResponse | null>(null);
const loading = ref(false);
const cancelling = ref(false);
const downloadingFirmware = ref(false);

const executions = ref<OTAExecutionResponse[]>([]);
const executionsLoading = ref(false);
const stateFilter = ref<OTAExecutionState | null>(null);
const currentPage = ref(1);
const totalPages = ref(1);
const totalItems = ref(0);

/** Resolved names keyed by id (fall back to the id) */
const assetNames = ref<Map<string, string>>(new Map());
const templateNames = ref<Map<string, string>>(new Map());

/** COMPUTED */

/** Plan id from the route */
const planId = computed(() => route.params.id as string);

/** Whether the plan can still be canceled */
const canCancel = computed(() =>
  !!plan.value && OTA_PLAN_CANCELABLE_STATUSES.includes(plan.value.status)
);

/** DetailChip color + icon for the plan status */
const statusColor = computed<DetailChipColor>(() =>
  plan.value ? OTA_PLAN_STATUS_CHIP_COLORS[plan.value.status] ?? 'grey' : 'grey'
);
const statusIcon = computed(() =>
  plan.value ? OTA_PLAN_STATUS_ICONS[plan.value.status] ?? 'help' : 'help'
);

/** Whether the firmware artifact can be downloaded (present and not purged) */
const firmwareDownloadable = computed(() =>
  !!plan.value?.firmware && plan.value.firmware.status !== 'PURGED'
);

/** Progress fraction (0..1) from the counters */
const progressValue = computed(() => {
  const c = plan.value?.counters;
  if (!c || c.total === 0) return 0;
  return Math.min(Math.max(c.succeeded / c.total, 0), 1);
});
const progressPct = computed(() => Math.round(progressValue.value * 100));

/** Resolved source/target template display names */
const sourceTemplateName = computed(() =>
  plan.value ? templateNames.value.get(plan.value.sourceTemplateId) || plan.value.sourceTemplateId : ''
);
const targetTemplateName = computed(() =>
  plan.value ? templateNames.value.get(plan.value.targetTemplateId) || plan.value.targetTemplateId : ''
);

/** Options for the executions state filter */
const executionStateOptions = computed((): OtaExecutionStateOption[] => [
  { label: t.executions.allStatus.value, value: null },
  ...(Object.keys(OTA_EXECUTION_STATE_ICONS) as OTAExecutionState[]).map((state) => ({
    label: t.executionState(state),
    value: state,
  })),
]);

/** Executions table columns */
const executionColumns = computed((): QTableColumn<OTAExecutionResponse>[] => [
  { name: 'assetId', label: t.executions.columns.device.value, field: 'assetId', align: 'left' },
  { name: 'state', label: t.executions.columns.state.value, field: 'state', align: 'center' },
  { name: 'percentage', label: t.executions.columns.progress.value, field: 'percentage', align: 'left' },
  { name: 'attempts', label: t.executions.columns.attempts.value, field: 'attempts', align: 'center' },
  { name: 'error', label: t.executions.columns.error.value, field: 'error', align: 'left' },
  {
    name: 'updated',
    label: t.executions.columns.updated.value,
    field: 'updated',
    align: 'left',
    format: (val: string) => formatDateTime(val),
  },
]);

/** FUNCTIONS */

/** Format an ISO timestamp, tolerating empty values */
function formatDateTime(value: string | undefined): string {
  if (!value) return '—';
  return date.formatDate(value, OTA_DATETIME_MASK);
}

/** Human-readable byte size */
function formatSize(bytes: number | undefined): string {
  if (!bytes) return '—';
  const units = ['B', 'KB', 'MB', 'GB'];
  const i = Math.min(Math.floor(Math.log(bytes) / Math.log(1024)), units.length - 1);
  return `${(bytes / 1024 ** i).toFixed(1)} ${units[i]}`;
}

/** Format the schedule, showing the immediate label for a run-now plan */
function formatSchedule(): string {
  const startAt = plan.value?.startAt;
  const created = plan.value?.created;
  if (!startAt) return t.header.immediate.value;
  const startTime = new Date(startAt).getTime();
  const createdTime = created ? new Date(created).getTime() : Date.now();
  if (Number.isNaN(startTime) || startTime <= createdTime + 1000) {
    return t.header.immediate.value;
  }
  return date.formatDate(startAt, OTA_DATETIME_MASK);
}

/** DetailChip color + icon for an execution state */
function stateColor(state: string): DetailChipColor {
  return OTA_EXECUTION_STATE_COLORS[state] ?? 'grey';
}
function stateIcon(state: string): string {
  return OTA_EXECUTION_STATE_ICONS[state] ?? 'help';
}

/** Resolved asset name for an id, falling back to the id */
function assetName(id: string): string {
  return assetNames.value.get(id) || id;
}

/** Resolve + cache the source/target template names for the current plan */
async function resolveTemplateNames(): Promise<void> {
  if (!apis.assets || !plan.value) return;
  const ids = [plan.value.sourceTemplateId, plan.value.targetTemplateId].filter(
    (id) => id && !templateNames.value.has(id)
  );
  if (ids.length === 0) return;

  const resolved = await Promise.all(
    ids.map(async (id) => {
      try {
        const tpl = await apis.assets.assetTemplate.getById({ assetTemplateId: id });
        return [id, tpl?.name || id] as const;
      } catch {
        return [id, id] as const;
      }
    })
  );
  const next = new Map(templateNames.value);
  for (const [id, name] of resolved) next.set(id, name);
  templateNames.value = next;
}

/** Resolve + cache asset names for any uncached execution ids */
async function resolveAssetNames(rows: OTAExecutionResponse[]): Promise<void> {
  if (!apis.assets) return;
  const missing = [...new Set(rows.map((r) => r.assetId).filter((id) => id && !assetNames.value.has(id)))];
  if (missing.length === 0) return;

  const resolved = await Promise.all(
    missing.map(async (id) => {
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

/** Fetch the plan (with the denormalized firmware block) */
async function fetchPlan(): Promise<void> {
  if (!apis.assets) {
    notifyFail({ message: t.page.apiNotInitialized.value });
    return;
  }
  if (!planId.value) {
    notifyFail({ message: t.page.notFound.value });
    void router.back();
    return;
  }
  try {
    loading.value = true;
    plan.value = await apis.assets.otaPlans.getPlan({ planId: planId.value });
    await resolveTemplateNames();
  } catch (err: unknown) {
    logger.error('Error fetching OTA plan:', err);
    notifyFail({ message: t.page.loadError.value });
  } finally {
    loading.value = false;
  }
}

/** Fetch the plan executions with the current filter + pagination */
async function fetchExecutions(): Promise<void> {
  if (!apis.assets || !planId.value) return;
  try {
    executionsLoading.value = true;
    const query = cleanQueryParams({
      page: currentPage.value,
      perPage: DEFAULT_EXECUTIONS_PER_PAGE,
      state: stateFilter.value ?? undefined,
    });
    const response = await apis.assets.otaPlans.listExecutions({ planId: planId.value }, query);
    executions.value = response?.items || [];
    await resolveAssetNames(executions.value);
    if (response?.pagination) {
      totalItems.value = response.pagination.totalItems || 0;
      totalPages.value = response.pagination.totalPages || 1;
    }
  } catch (err: unknown) {
    logger.error('Error fetching OTA executions:', err);
    notifyFail({ message: t.page.loadError.value });
  } finally {
    executionsLoading.value = false;
  }
}

/** Refresh both plan + executions (polling and post-action) */
async function refreshAll(): Promise<void> {
  await Promise.all([fetchPlan(), fetchExecutions()]);
}

/** Apply the state filter and refetch from the first page */
function applyStateFilter(): void {
  currentPage.value = 1;
  void fetchExecutions();
}

/** Header refresh button: clear the state filter back to ALL and reload everything */
async function resetAndRefresh(): Promise<void> {
  stateFilter.value = null;
  currentPage.value = 1;
  await refreshAll();
}

/** Executions pagination */
function handlePageChange(page: number): void {
  currentPage.value = page;
  void fetchExecutions();
}

/** Confirm and cancel the plan, then refresh */
async function confirmCancel(): Promise<void> {
  const confirmed = await dialogDelete({
    title: t.cancelDialog.title.value,
    message: t.cancelDialog.message.value,
    ok: { label: t.cancelDialog.confirm.value, color: 'negative' },
  });
  if (!confirmed || !apis.assets || !planId.value) return;
  try {
    cancelling.value = true;
    await apis.assets.otaPlans.deletePlan({ planId: planId.value });
    notifySuccess({ message: t.cancelDialog.success.value });
    await refreshAll();
  } catch (err: unknown) {
    logger.error('Error cancelling OTA plan:', err);
    notifyFail({ message: t.cancelDialog.error.value });
  } finally {
    cancelling.value = false;
  }
}

/**
 * Mint a presigned download URL for the firmware and open it. The backend
 * returns 404 when the artifact is no longer in storage — surfaced as a
 * "firmware not found" notice.
 */
async function downloadFirmware(): Promise<void> {
  if (!apis.assets || !planId.value) return;
  try {
    downloadingFirmware.value = true;
    const res = await apis.assets.otaPlans.downloadFirmware({ planId: planId.value });
    const link = document.createElement('a');
    link.href = res.url;
    link.target = '_blank';
    link.rel = 'noopener';
    document.body.appendChild(link);
    link.click();
    link.remove();
  } catch (err: unknown) {
    logger.error('Error downloading firmware:', err);
    const status = (err as { response?: { status?: number } })?.response?.status;
    notifyFail({ message: status === 404 ? t.firmware.notFound.value : t.firmware.downloadError.value });
  } finally {
    downloadingFirmware.value = false;
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
      icon="system_update"
      iconColor="primary"
      :title="plan?.name || ''"
      :description="plan?.description || ''"
      :button="{ label: t.page.back.value, icon: 'arrow_back', flat: true, to: OTA_ROUTE_BASE }"
    />

    <!-- Loading Spinner -->
    <div v-if="loading && !plan" class="row justify-center q-my-xl">
      <q-spinner color="primary" size="3em" />
    </div>

    <div v-else-if="plan" class="row q-col-gutter-lg q-mb-lg">
      <!-- LEFT: plan info panel -->
      <div class="col-12 col-md-4">
        <q-card class="rounded-borders">
          <!-- Panel header -->
          <q-card-section class="bg-grey-1 q-pb-md">
            <div class="text-h6 text-weight-bold text-primary">
              <q-icon size="sm" name="system_update" color="primary" class="q-mr-xs" />
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
              <div class="info-field">
                <div class="info-label">{{ t.fields.status.value }}</div>
                <DetailChip
                  dense
                  size="sm"
                  :color="statusColor"
                  :icon="statusIcon"
                  :label="t.planStatus(plan.status)"
                />
              </div>
            </div>

            <!-- Section: Firmware -->
            <div class="info-section">
              <div class="row items-center justify-between q-mb-sm">
                <div class="section-title q-mb-none">
                  <q-icon name="memory" size="xs" class="q-mr-xs" />{{ t.firmware.title.value }}
                </div>
                <q-btn
                  v-if="firmwareDownloadable"
                  dense
                  flat
                  no-caps
                  size="sm"
                  color="primary"
                  icon="download"
                  :label="t.firmware.download.value"
                  :loading="downloadingFirmware"
                  @click="downloadFirmware"
                />
              </div>
              <div v-if="plan.firmware" class="firmware-box">
                <div class="fw-row">
                  <span class="fw-key">{{ t.firmware.version.value }}</span>
                  <span class="fw-val">{{ plan.firmware.version }}</span>
                </div>
                <div class="fw-row">
                  <span class="fw-key">{{ t.firmware.filename.value }}</span>
                  <span class="fw-val">{{ plan.firmware.filename }}</span>
                </div>
                <div class="fw-row">
                  <span class="fw-key">{{ t.firmware.size.value }}</span>
                  <span class="fw-val">{{ formatSize(plan.firmware.size) }}</span>
                </div>
                <div class="fw-row">
                  <span class="fw-key">{{ t.firmware.checksum.value }}</span>
                  <span class="fw-val fw-mono">{{ plan.firmware.checksum }}</span>
                </div>
                <div class="fw-row">
                  <span class="fw-key">{{ t.firmware.artifactStatus.value }}</span>
                  <span class="fw-val">{{ plan.firmware.status }}</span>
                </div>
              </div>
              <div v-else class="info-value text-grey-6">{{ t.firmware.empty.value }}</div>
            </div>

            <!-- Section: Templates -->
            <div class="info-section">
              <div class="section-title">
                <q-icon name="dashboard_customize" size="xs" class="q-mr-xs" />{{ t.sections.templates.value }}
              </div>
              <div class="info-field">
                <div class="info-label">{{ t.fields.sourceTemplate.value }}</div>
                <div class="info-value">{{ sourceTemplateName }}</div>
              </div>
              <div class="info-field q-mb-none">
                <div class="info-label">{{ t.fields.targetTemplate.value }}</div>
                <div class="info-value">{{ targetTemplateName }}</div>
              </div>
            </div>

            <!-- Section: Schedule -->
            <div class="info-section q-mb-none">
              <div class="section-title">
                <q-icon name="schedule" size="xs" class="q-mr-xs" />{{ t.sections.schedule.value }}
              </div>
              <div class="info-field">
                <div class="info-label">{{ t.fields.schedule.value }}</div>
                <div class="info-value">{{ formatSchedule() }}</div>
              </div>
              <div class="row q-col-gutter-md">
                <div class="col-6">
                  <div class="info-field q-mb-none">
                    <div class="info-label">{{ t.fields.deadline.value }}</div>
                    <div class="info-value">{{ formatDateTime(plan.maxTime) }}</div>
                  </div>
                </div>
                <div class="col-6">
                  <div class="info-field q-mb-none">
                    <div class="info-label">{{ t.fields.created.value }}</div>
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
                  {{ t.progress.total.value }}: <span class="text-body2 text-weight-medium">{{ plan.counters.total }}</span>
                </div>
                <div class="text-caption text-grey-7">
                  {{ t.progress.succeeded.value }}: <span class="text-body2 text-positive text-weight-medium">{{ plan.counters.succeeded }}</span>
                </div>
                <div class="text-caption text-grey-7">
                  {{ t.progress.failed.value }}: <span class="text-body2 text-negative text-weight-medium">{{ plan.counters.failed }}</span>
                </div>
                <div class="text-caption text-grey-7">
                  {{ t.progress.timedOut.value }}: <span class="text-body2 text-weight-medium">{{ plan.counters.timedOut }}</span>
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
                  <q-icon size="sm" name="devices" color="primary" class="q-mr-xs" />
                  {{ t.executions.title.value }}
                </div>
                <div class="text-caption text-grey-7">
                  {{ totalItems }} {{ totalItems === 1 ? t.executions.itemLabel.value : t.executions.itemLabelPlural.value }}
                </div>
              </div>
              <div class="col-auto" style="min-width: 200px;">
                <q-select
                  v-model="stateFilter"
                  outlined
                  dense
                  emit-value
                  map-options
                  :options="executionStateOptions"
                  :label="t.executions.filterStatus.value"
                  class="filter-input"
                  @update:model-value="applyStateFilter"
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
                  <q-td :props="props">{{ assetName(props.row.assetId) }}</q-td>
                </template>

                <template #body-cell-state="props">
                  <q-td :props="props" class="text-center">
                    <DetailChip
                      dense
                      size="sm"
                      :color="stateColor(props.row.state)"
                      :icon="stateIcon(props.row.state)"
                      :label="t.executionState(props.row.state)"
                    />
                  </q-td>
                </template>

                <template #body-cell-percentage="props">
                  <q-td :props="props">
                    <q-linear-progress
                      v-if="props.row.state === 'DOWNLOADING' || props.row.state === 'UPDATING'"
                      rounded
                      size="8px"
                      color="orange-6"
                      track-color="grey-3"
                      :value="(props.row.percentage || 0) / 100"
                    />
                    <span v-else class="text-grey-7">{{ props.row.percentage || 0 }}%</span>
                  </q-td>
                </template>

                <template #body-cell-error="props">
                  <q-td :props="props">
                    <span v-if="props.row.error" class="text-negative">{{ props.row.error }}</span>
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
            <div class="row justify-end">
              <BaseButton
                outline
                no-caps
                color="negative"
                icon="cancel"
                :label="t.actions.cancel.value"
                :disable="!canCancel"
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

.firmware-box {
  border: 1px solid var(--mapex-card-border);
  border-radius: var(--mapex-radius-md);
  padding: var(--mapex-spacing-sm);
}

.fw-row {
  display: flex;
  justify-content: space-between;
  gap: var(--mapex-spacing-sm);
  padding: 2px 0;
}

.fw-key {
  font-size: 0.75rem;
  color: var(--mapex-text-secondary);
}

.fw-val {
  font-size: 0.85rem;
  color: var(--mapex-text-primary);
  text-align: right;
  word-break: break-all;
}

.fw-mono {
  font-family: monospace;
}

.filter-input {
  :deep(.q-field__control) {
    border-radius: var(--mapex-radius-md);
  }
}
</style>
