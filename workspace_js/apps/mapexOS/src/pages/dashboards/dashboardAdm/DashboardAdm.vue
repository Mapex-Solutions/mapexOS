<template>
  <q-page class="q-pa-lg dashboard-background">
    <!-- Header Section -->
    <div class="row items-center q-mb-lg">
      <div class="col-12">
        <div class="row items-center">
          <div class="q-mr-md">
            <q-icon name="dashboard" size="xl" color="primary" />
          </div>
          <div>
            <div class="text-h4 text-weight-bold text-primary">{{ t.header.title.value }}</div>
            <div class="text-subtitle1 text-grey-7">
              {{ t.header.subtitlePrefix.value }}
              <span class="text-weight-bold">{{ orgStore.selectedOrganizationName ?? t.header.notAvailable.value }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div class="row q-col-gutter-lg">
      <!-- Primary KPI Cards -->
      <div class="col-12">
        <div class="row q-col-gutter-lg">
          <div
            v-for="kpi in primaryKpis"
            :key="kpi.label"
            class="col-12 col-sm-6 col-md-3"
          >
            <q-card flat class="kpi-card kpi-card--clickable" @click="navigateTo(kpi.route)">
              <q-card-section class="q-pa-lg">
                <div class="row items-center">
                  <div class="col">
                    <template v-if="loading">
                      <q-skeleton type="text" width="60px" height="40px" class="q-mb-sm" />
                      <q-skeleton type="text" width="80px" />
                    </template>
                    <template v-else>
                      <div class="kpi-value" :class="kpi.color">{{ kpi.value }}</div>
                      <div class="kpi-label">{{ kpi.label }}</div>
                    </template>
                  </div>
                  <div class="kpi-icon-container" :class="kpi.bgColor">
                    <q-icon size="32px" :name="kpi.icon" :class="kpi.color" />
                  </div>
                </div>
              </q-card-section>
            </q-card>
          </div>
        </div>
      </div>

      <!-- Secondary KPI Cards -->
      <div class="col-12">
        <div class="row q-col-gutter-lg">
          <div
            v-for="kpi in secondaryKpis"
            :key="kpi.label"
            class="col-12 col-sm-6 col-md-3"
          >
            <q-card flat class="kpi-card kpi-card--clickable" @click="navigateTo(kpi.route)">
              <q-card-section class="q-pa-lg">
                <div class="row items-center">
                  <div class="col">
                    <template v-if="loading">
                      <q-skeleton type="text" width="60px" height="40px" class="q-mb-sm" />
                      <q-skeleton type="text" width="80px" />
                    </template>
                    <template v-else>
                      <div class="kpi-value" :class="kpi.color">{{ kpi.value }}</div>
                      <div class="kpi-label">{{ kpi.label }}</div>
                    </template>
                  </div>
                  <div class="kpi-icon-container" :class="kpi.bgColor">
                    <q-icon size="32px" :name="kpi.icon" :class="kpi.color" />
                  </div>
                </div>
              </q-card-section>
            </q-card>
          </div>
        </div>
      </div>

      <!-- Quick Actions + Platform Info -->
      <div class="col-12 col-md-8 column">
        <q-card flat class="actions-card full-height">
          <q-card-section class="q-pa-lg">
            <div class="section-title q-mb-lg">{{ t.sections.quickActions.value }}</div>
            <div class="row q-col-gutter-md">
              <div
                v-for="action in quickActions"
                :key="action.label"
                class="col-6 col-sm-4"
              >
                <div class="action-btn" @click="navigateTo(action.route)">
                  <div class="action-icon-container" :class="action.bgColor">
                    <q-icon size="24px" :name="action.icon" :class="action.color" />
                  </div>
                  <div class="action-label">{{ action.label }}</div>
                </div>
              </div>
            </div>
          </q-card-section>
        </q-card>
      </div>

      <div class="col-12 col-md-4 column">
        <q-card flat class="platform-card full-height">
          <q-card-section class="q-pa-lg">
            <div class="section-title q-mb-lg">{{ t.sections.platformInfo.value }}</div>

            <div class="platform-info-list">
              <div class="platform-info-item">
                <div class="platform-info-label">{{ t.platformInfo.organization.value }}</div>
                <div class="platform-info-value">
                  {{ orgStore.selectedOrganizationName ?? t.platformInfo.notAvailable.value }}
                </div>
              </div>

              <q-separator class="q-my-md" />

              <div class="platform-info-item">
                <div class="platform-info-label">{{ t.platformInfo.type.value }}</div>
                <div class="platform-info-value">
                  <DetailChip
                    :label="organizationType"
                    :color="organizationTypeColor"
                    size="sm"
                  />
                </div>
              </div>

              <q-separator class="q-my-md" />

              <div class="platform-info-item">
                <div class="platform-info-label">{{ t.platformInfo.totalOrganizations.value }}</div>
                <div class="platform-info-value text-weight-bold">
                  {{ orgStore.totalCount }}
                </div>
              </div>
            </div>
          </q-card-section>
        </q-card>
      </div>
    </div>
  </q-page>
</template>

<script setup lang="ts">
defineOptions({
  name: 'DashboardAdm',
});

/** TYPE IMPORTS */
import type { DetailChipColor } from '@components/chips/DetailChip';

/** VUE IMPORTS */
import { ref, computed, onMounted } from 'vue';
import { useRouter } from 'vue-router';

/** COMPONENTS */
import { DetailChip } from '@components/chips/DetailChip';

/** COMPOSABLES */
import { useLogger } from '@composables/useLogger';
import { useOrgChangeRefresh } from '@composables/organizations';
import { useDashboardAdmTranslations } from '@composables/i18n/pages/dashboards/dashboardAdm';

/** SERVICES */
import { apis } from '@services/mapex';

/** STORES */
import { useOrganizationStore } from '@stores/organization';

/** COMPOSABLES & STORES */
const logger = useLogger('DashboardAdm');
const router = useRouter();
const orgStore = useOrganizationStore();
const t = useDashboardAdmTranslations();

/** STATE */
const loading = ref(false);

/**
 * Raw KPI values fetched from the counter APIs.
 *
 * Stored separately from the label/icon metadata so that locale changes do
 * not reset the values and so that `fetchDashboardData` can mutate them
 * without breaking the computed pipeline below.
 */
const primaryKpiValues = ref<number[]>([0, 0, 0, 0]);
const secondaryKpiValues = ref<number[]>([0, 0, 0, 0]);

/**
 * Primary KPI cards
 *
 * Labels are reactive computed properties so they respond to locale changes.
 */
const primaryKpis = computed(() => [
  { label: t.primaryKpis.assets.value, value: primaryKpiValues.value[0] ?? 0, icon: 'devices', color: 'text-blue-6', bgColor: 'bg-blue-1', route: '/assets' },
  { label: t.primaryKpis.users.value, value: primaryKpiValues.value[1] ?? 0, icon: 'people', color: 'text-teal-6', bgColor: 'bg-teal-1', route: '/users' },
  { label: t.primaryKpis.triggers.value, value: primaryKpiValues.value[2] ?? 0, icon: 'flash_on', color: 'text-orange-6', bgColor: 'bg-orange-1', route: '/triggers' },
  { label: t.primaryKpis.wfDefinitions.value, value: primaryKpiValues.value[3] ?? 0, icon: 'account_tree', color: 'text-purple-6', bgColor: 'bg-purple-1', route: '/workflows' },
]);

/**
 * Secondary KPI cards
 */
const secondaryKpis = computed(() => [
  { label: t.secondaryKpis.assetTemplates.value, value: secondaryKpiValues.value[0] ?? 0, icon: 'dashboard_customize', color: 'text-indigo-6', bgColor: 'bg-indigo-1', route: '/assets_template' },
  { label: t.secondaryKpis.groups.value, value: secondaryKpiValues.value[1] ?? 0, icon: 'group_work', color: 'text-cyan-6', bgColor: 'bg-cyan-1', route: '/groups' },
  { label: t.secondaryKpis.routeGroups.value, value: secondaryKpiValues.value[2] ?? 0, icon: 'alt_route', color: 'text-amber-8', bgColor: 'bg-amber-1', route: '/routing/route_groups' },
  { label: t.secondaryKpis.wfInstances.value, value: secondaryKpiValues.value[3] ?? 0, icon: 'play_circle', color: 'text-deep-purple-6', bgColor: 'bg-deep-purple-1', route: '/workflow_instances' },
]);

/**
 * Quick action shortcuts
 */
const quickActions = computed(() => [
  { label: t.quickActions.newAsset.value, route: '/assets/add', icon: 'add_circle_outline', color: 'text-blue-6', bgColor: 'bg-blue-1' },
  { label: t.quickActions.newUser.value, route: '/users/add', icon: 'person_add', color: 'text-teal-6', bgColor: 'bg-teal-1' },
  { label: t.quickActions.newTemplate.value, route: '/assets_template/add', icon: 'dashboard_customize', color: 'text-indigo-6', bgColor: 'bg-indigo-1' },
  { label: t.quickActions.newRouteGroup.value, route: '/routing/route_groups/add', icon: 'alt_route', color: 'text-amber-8', bgColor: 'bg-amber-1' },
  { label: t.quickActions.viewLogs.value, route: '/logs/raw_logs', icon: 'receipt_long', color: 'text-grey-7', bgColor: 'bg-grey-2' },
]);

/** COMPUTED */
const organizationType = computed(() => {
  const org = orgStore.selectedOrganization;
  if (!org) return t.platformInfo.notAvailable.value;
  return org.type.charAt(0).toUpperCase() + org.type.slice(1);
});

const organizationTypeColor = computed<DetailChipColor>(() => {
  const org = orgStore.selectedOrganization;
  if (!org) return 'grey';
  const colorMap: Record<string, DetailChipColor> = {
    vendor: 'primary',
    customer: 'teal',
    site: 'amber',
    building: 'orange',
    floor: 'cyan',
    zone: 'indigo',
  };
  return colorMap[org.type] ?? 'grey';
});

/** FUNCTIONS */

/**
 * Fetch dashboard KPI data from all counter APIs in parallel
 *
 * Uses dedicated /counter endpoints with Redis cache (6h TTL) instead of list() calls.
 * Promise.allSettled ensures one failing API does not break the entire dashboard.
 */
async function fetchDashboardData(): Promise<void> {
  loading.value = true;
  try {
    const [
      assets,
      users,
      triggers,
      definitions,
      templates,
      groups,
      routeGroups,
      instances,
    ] = await Promise.allSettled([
      apis.assets?.asset?.counter(),
      apis.mapexOS?.users?.counter(),
      apis.triggers?.trigger?.counter(),
      apis.workflows?.definition?.counter(),
      apis.assets?.assetTemplate?.counter(),
      apis.mapexOS?.groups?.counter(),
      apis.router?.routegroup?.counter(),
      apis.workflows?.instance?.counter(),
    ]);

    const primaryTotals = [assets, users, triggers, definitions].map(extractTotal);
    const secondaryTotals = [templates, groups, routeGroups, instances].map(extractTotal);

    primaryKpiValues.value = primaryTotals;
    secondaryKpiValues.value = secondaryTotals;

    logger.debug('Dashboard data loaded successfully');
  } catch (error) {
    logger.error('Failed to fetch dashboard data', error);
  } finally {
    loading.value = false;
  }
}

/**
 * Extract count from a settled counter promise result
 *
 * @param {PromiseSettledResult<any>} result - The settled promise result
 * @returns {number} The entity count, or 0 if the promise was rejected
 */
function extractTotal(result: PromiseSettledResult<any>): number {
  if (result.status === 'fulfilled') {
    return result.value?.count ?? 0;
  }
  return 0;
}

/**
 * Navigate to a route using vue-router
 *
 * @param {string} route - The route path to navigate to
 */
function navigateTo(route: string): void {
  void router.push(route);
}

/** LIFECYCLE HOOKS */
onMounted(() => {
  void fetchDashboardData();
});

useOrgChangeRefresh(() => {
  void fetchDashboardData();
});
</script>

<style lang="scss" scoped>
// Dashboard Background
.dashboard-background {
  background: transparent;
  min-height: 100vh;
}

// KPI Cards
.kpi-card {
  border-radius: var(--mapex-radius-xl);
  box-shadow: 0 2px 12px var(--mapex-elevation-shadow);
  border: 1px solid var(--mapex-card-border);
  background: var(--mapex-surface-elevated);
  transition: var(--mapex-transition-slow);
  overflow: hidden;
}

.kpi-card--clickable {
  cursor: pointer;
}

.kpi-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 32px var(--mapex-hover-shadow);
}

.kpi-value {
  font-size: 2.5rem;
  font-weight: 700;
  line-height: 1;
  margin-bottom: 8px;
}

.kpi-label {
  font-size: 0.875rem;
  font-weight: 500;
  color: var(--mapex-text-secondary);
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.kpi-icon-container {
  width: 64px;
  height: 64px;
  border-radius: var(--mapex-radius-xl);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

// Actions Card
.actions-card {
  border-radius: var(--mapex-radius-xl);
  box-shadow: 0 2px 12px var(--mapex-elevation-shadow);
  border: 1px solid var(--mapex-card-border);
  background: var(--mapex-surface-elevated);
}

.action-btn {
  padding: 20px;
  border-radius: var(--mapex-radius-lg);
  border: 1px solid var(--mapex-card-border);
  background: var(--mapex-surface-bg);
  cursor: pointer;
  transition: var(--mapex-transition-slow);
  text-align: center;
  height: 100px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}

.action-btn:hover {
  border-color: var(--mapex-card-hover-border);
  box-shadow: 0 8px 24px var(--mapex-hover-shadow);
  transform: translateY(-2px);
}

.action-icon-container {
  width: 48px;
  height: 48px;
  border-radius: var(--mapex-radius-lg);
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 8px;
}

.action-label {
  font-size: 0.875rem;
  font-weight: 600;
  color: var(--mapex-text-primary);
}

// Platform Info Card
.platform-card {
  border-radius: var(--mapex-radius-xl);
  box-shadow: 0 2px 12px var(--mapex-elevation-shadow);
  border: 1px solid var(--mapex-card-border);
  background: var(--mapex-surface-elevated);
}

.platform-info-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.platform-info-label {
  font-size: 0.875rem;
  font-weight: 500;
  color: var(--mapex-text-secondary);
}

.platform-info-value {
  font-size: 0.9375rem;
  color: var(--mapex-text-primary);
}

// Common Styles
.section-title {
  font-size: 1.25rem;
  font-weight: 700;
  color: var(--mapex-text-primary);
}

// Responsive Design
@media (max-width: 768px) {
  .kpi-value {
    font-size: 2rem;
  }

  .kpi-icon-container {
    width: 48px;
    height: 48px;
  }

  .action-btn {
    height: 80px;
    padding: 16px;
  }

  .action-icon-container {
    width: 36px;
    height: 36px;
  }
}
</style>
