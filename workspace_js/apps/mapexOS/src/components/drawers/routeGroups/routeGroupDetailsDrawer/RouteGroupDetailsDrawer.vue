<template>
  <q-drawer
    overlay
    bordered
    side="right"
    :model-value="modelValue"
    :width="450"
    @update:model-value="emit('update:modelValue', $event)"
  >
    <!-- Header -->
    <q-toolbar class="drawer-header">
      <q-icon name="route" size="sm" color="primary" class="q-mr-sm" />
      <q-toolbar-title class="text-weight-medium">{{ t.drawer.title.value }}</q-toolbar-title>

      <q-btn flat round dense icon="close" class="drawer-close-btn" @click="close">
        <AppTooltip :content="t.drawer.close.value" />
      </q-btn>
    </q-toolbar>

    <q-separator />

    <!-- Content -->
    <div class="drawer-content">
      <q-scroll-area class="fit">
        <!-- Loading State -->
        <div v-if="loading" class="q-pa-lg text-center">
          <q-spinner size="3em" color="primary" class="q-mb-md" />
          <div class="text-grey-7">{{ t.drawer.loading.value }}</div>
        </div>

        <!-- Error State -->
        <div v-else-if="error" class="q-pa-lg">
          <q-banner rounded class="bg-negative text-white">
            <template #avatar>
              <q-icon name="error" color="white" />
            </template>
            {{ t.drawer.error.value }}
          </q-banner>
        </div>

        <!-- Route Group Data -->
        <div v-else-if="routeGroup" class="q-px-md q-py-lg">

          <!-- Basic Information Section -->
          <div class="section q-mb-md">
            <div class="section-header">
              <q-icon name="info" color="primary" size="sm" class="q-mr-sm" />
              <span class="text-subtitle1 text-weight-medium">{{ t.drawer.sections.basicInfo.value }}</span>
            </div>
            <q-separator class="q-my-sm" />

            <div class="field-row">
              <div class="field-label">{{ t.drawer.fields.name.value }}</div>
              <div class="field-value text-weight-medium">{{ routeGroup?.name || '-' }}</div>
            </div>

            <div class="field-row">
              <div class="field-label">{{ t.drawer.fields.enabled.value }}</div>
              <div class="field-value">
                <DetailChip
                  :color="routeGroup?.enabled ? 'positive' : 'negative'"
                  size="sm"
                  :label="routeGroup?.enabled ? t.status.active.value.toUpperCase() : t.status.inactive.value.toUpperCase()"
                />
              </div>
            </div>

            <div class="field-row">
              <div class="field-label">{{ t.drawer.fields.description.value }}</div>
              <div class="field-value text-grey-8">
                {{ routeGroup?.description || t.drawer.empty.description.value }}
              </div>
            </div>
          </div>

          <!-- Routers Section -->
          <div class="section q-mb-md">
            <div class="section-header">
              <q-icon name="alt_route" color="primary" size="sm" class="q-mr-sm" />
              <span class="text-subtitle1 text-weight-medium">{{ t.drawer.sections.routers.value }}</span>
              <q-badge
                v-if="routeGroup?.routers?.length"
                color="primary"
                class="q-ml-sm"
                :label="routeGroup.routers.length"
              />
            </div>
            <q-separator class="q-my-sm" />

            <!-- No Routers -->
            <div v-if="!routeGroup?.routers?.length" class="text-grey-6 text-caption q-py-sm">
              {{ t.drawer.empty.routers.value }}
            </div>

            <!-- Router Cards -->
            <div v-else class="q-gutter-sm q-mt-sm">
              <q-card
                v-for="(router, index) in routeGroup.routers"
                :key="index"
                flat
                bordered
                class="router-card"
              >
                <q-card-section class="q-py-sm q-px-md">
                  <div class="row items-center">
                    <q-icon
                      :name="getRouterIcon(router.kind)"
                      :color="getRouterColor(router.kind)"
                      size="sm"
                      class="q-mr-sm"
                    />
                    <div class="col">
                      <div class="text-body2 text-weight-medium">
                        {{ getRouterKindLabel(router.kind) }}
                      </div>
                      <div class="text-caption text-grey-6">
                        {{ getRouterDestination(router) }}
                      </div>
                    </div>
                    <q-badge
                      v-if="router.match"
                      color="orange"
                      text-color="white"
                      class="q-ml-sm"
                    >
                      <q-icon name="filter_alt" size="xs" class="q-mr-xs" />
                      {{ router.match.rules?.length || 0 }} rules
                    </q-badge>
                  </div>

                  <!-- Match Rules Preview -->
                  <div v-if="router.match?.rules?.length" class="q-mt-sm">
                    <div class="text-caption text-grey-7">
                      <q-icon name="settings_ethernet" size="xs" class="q-mr-xs" />
                      Policy: <strong>{{ router.match.policy === 'all' ? 'ALL (AND)' : 'ANY (OR)' }}</strong>
                    </div>
                  </div>
                </q-card-section>
              </q-card>
            </div>
          </div>

          <!-- Timestamps Section -->
          <div class="section">
            <div class="section-header">
              <q-icon name="schedule" color="primary" size="sm" class="q-mr-sm" />
              <span class="text-subtitle1 text-weight-medium">{{ t.drawer.sections.timestamps.value }}</span>
            </div>
            <q-separator class="q-my-sm" />

            <!-- Created & Updated (2 columns) -->
            <div class="row q-col-gutter-sm">
              <div class="col-6">
                <div class="field-row">
                  <div class="field-label">{{ t.drawer.fields.created.value }}</div>
                  <div class="field-value">{{ formatDate(routeGroup?.created) }}</div>
                </div>
              </div>
              <div class="col-6">
                <div class="field-row">
                  <div class="field-label">{{ t.drawer.fields.updated.value }}</div>
                  <div class="field-value">{{ formatDate(routeGroup?.updated) }}</div>
                </div>
              </div>
            </div>
          </div>

        </div>
      </q-scroll-area>
    </div>
  </q-drawer>
</template>

<script setup lang="ts">
defineOptions({
  name: 'RouteGroupDetailsDrawer'
});

/** TYPE IMPORTS */
import type { RouteGroupResponse as SchemaRouteGroupResponse } from '@mapexos/schemas';
import type { RouteGroupDetailsDrawerProps, RouteGroupDetailsDrawerEmits, RouterWithMatch } from './interfaces/routeGroupDetailsDrawer.interface';

// RouteGroup response with proper types
type RouteGroupResponse = Omit<SchemaRouteGroupResponse, 'routers'> & {
  routers?: RouterWithMatch[];
};

/** VUE IMPORTS */
import { ref, watch, onMounted, onBeforeUnmount } from 'vue';
import { date } from 'quasar';

/** COMPONENTS */
import { DetailChip } from '@components/chips';
import { AppTooltip } from '@components/tooltips';

/** SERVICES */
import { apis } from '@services/mapex';

/** COMPOSABLES */
import { useRouteGroupsTranslations } from '@composables/i18n';
import { useLogger } from '@composables/useLogger';

/** UTILS */
import { notifyFail } from '@utils/alert';

/** PROPS & EMITS */
const props = defineProps<RouteGroupDetailsDrawerProps>();
const emit = defineEmits<RouteGroupDetailsDrawerEmits>();

/** COMPOSABLES */
const t = useRouteGroupsTranslations();
const logger = useLogger('RouteGroupDetailsDrawer');

/** STATE */
const routeGroup = ref<RouteGroupResponse | null>(null);
const loading = ref(false);
const error = ref(false);

/** FUNCTIONS */

/**
 * Handle ESC key to close drawer
 * @param {KeyboardEvent} event - Keyboard event
 */
function handleEscKey(event: KeyboardEvent): void {
  if (event.key === 'Escape' && props.modelValue) {
    close();
  }
}

/**
 * Close the drawer
 */
function close(): void {
  emit('update:modelValue', false);
}

/**
 * Format date using Quasar date utils
 * @param {string | undefined} dateValue - Date string to format
 * @returns {string} Formatted date string
 */
function formatDate(dateValue?: string): string {
  if (!dateValue) return '-';

  try {
    const dateObj = new Date(dateValue);
    return date.formatDate(dateObj, 'MMM DD, YYYY HH:mm');
  } catch {
    return '-';
  }
}

/**
 * Get router icon based on kind
 * @param {string} kind - Router kind
 * @returns {string} Icon name
 */
function getRouterIcon(kind: string): string {
  const iconMap: Record<string, string> = {
    'lake_house': 'storage',
    'notification': 'notifications',
    'save_event': 'save',
  };
  return iconMap[kind] || 'route';
}

/**
 * Get router color based on kind
 * @param {string} kind - Router kind
 * @returns {string} Color name
 */
function getRouterColor(kind: string): string {
  const colorMap: Record<string, string> = {
    'lake_house': 'purple-6',
    'notification': 'orange-6',
    'save_event': 'green-6',
  };
  return colorMap[kind] || 'grey-6';
}

/**
 * Get router kind label
 * @param {string} kind - Router kind
 * @returns {string} Translated label
 */
function getRouterKindLabel(kind: string): string {
  const kindKey = kind as keyof typeof t.routerKinds;
  return t.routerKinds[kindKey]?.label?.value || kind;
}

/**
 * Get router destination info
 * @param {RouterWithMatch} router - Router object
 * @returns {string} Destination description
 */
function getRouterDestination(router: RouterWithMatch): string {
  switch (router.kind) {
    case 'lake_house':
      return router.lakeHouse?.lakeHouseId || 'No data lake configured';
    case 'notification':
      return router.notification?.notificationId || 'No notification configured';
    case 'save_event':
      return 'Saves event to storage';
    default:
      return 'Unknown destination';
  }
}

/**
 * Fetch route group details from API
 */
async function fetchRouteGroup(): Promise<void> {
  if (!props.routeGroupId || !apis.router) {
    return;
  }

  loading.value = true;
  error.value = false;
  routeGroup.value = null;

  try {
    const response = await apis.router.routegroup.getById({
      routeGroupId: props.routeGroupId,
    });

    // Cast to local type with match field support
    routeGroup.value = response as RouteGroupResponse;
  } catch (err: any) {
    logger.error('Error fetching route group:', err);
    error.value = true;
    notifyFail({ message: t.drawer.error.value });
  } finally {
    loading.value = false;
  }
}

/** LIFECYCLE HOOKS */
onMounted(() => {
  window.addEventListener('keydown', handleEscKey);
});

onBeforeUnmount(() => {
  window.removeEventListener('keydown', handleEscKey);
});

/** WATCHERS */

// Watch for changes in routeGroupId
watch(() => props.routeGroupId, (newId) => {
  if (newId && props.modelValue) {
    void fetchRouteGroup();
  }
}, { immediate: true });

// Watch for drawer opening
watch(() => props.modelValue, (isOpen) => {
  if (isOpen && props.routeGroupId) {
    void fetchRouteGroup();
  } else if (!isOpen) {
    // Reset state when drawer closes
    routeGroup.value = null;
    error.value = false;
  }
});
</script>

<style lang="scss" scoped>
// Drawer Header
.drawer-header {
  background: var(--mapex-header-bg);
  backdrop-filter: blur(10px);
  border-bottom: 1px solid var(--mapex-header-border);

  .q-toolbar__title {
    font-size: 1.1rem;
    color: var(--q-primary);
  }
}

// Close button
.drawer-close-btn {
  color: var(--mapex-text-secondary);
}

// Drawer Content
.drawer-content {
  height: calc(100vh - 64px); // Subtract header height
  overflow: hidden;
}

// Section Styling
.section {
  .section-header {
    display: flex;
    align-items: center;
    color: var(--q-primary);
    margin-bottom: 8px;
  }
}

// Field Row Styling
.field-row {
  display: flex;
  flex-direction: column;
  padding: 10px 0;
  border-bottom: 1px solid var(--mapex-divider);

  &:last-child {
    border-bottom: none;
  }

  .field-label {
    font-size: 0.7rem;
    font-weight: 600;
    text-transform: uppercase;
    color: var(--mapex-text-secondary);
    margin-bottom: 4px;
    letter-spacing: 0.8px;
  }

  .field-value {
    font-size: 0.9rem;
    color: var(--mapex-text-primary);
    word-break: break-word;
    line-height: 1.4;
  }
}

// Router Card
.router-card {
  border-radius: var(--mapex-radius-md);
  background: var(--mapex-surface-bg);
}

// Custom Scrollbar
:deep(.q-scrollarea__content) {
  &::-webkit-scrollbar {
    width: 6px;
  }

  &::-webkit-scrollbar-track {
    background: transparent;
    border-radius: var(--mapex-radius-lg);
  }

  &::-webkit-scrollbar-thumb {
    background: rgba(var(--q-primary-rgb), 0.3);
    border-radius: var(--mapex-radius-lg);
    transition: background var(--mapex-transition-base) ease;

    &:hover {
      background: rgba(var(--q-primary-rgb), 0.5);
    }
  }
}
</style>
