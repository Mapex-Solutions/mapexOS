<template>
  <!-- Invisible backdrop for click outside detection -->
  <Teleport to="body">
    <div
      v-if="modelValue"
      class="drawer-backdrop"
      @click="close"
    />
  </Teleport>

  <q-drawer
    overlay
    bordered
    side="right"
    :model-value="modelValue"
    :width="450"
    @update:model-value="emit('update:modelValue', $event)"
    @keydown.esc="close"
  >
    <!-- Header -->
    <q-toolbar class="drawer-header">
      <q-icon name="flash_on" size="sm" class="q-mr-sm" color="primary" />
      <q-toolbar-title class="text-weight-medium">Trigger Details</q-toolbar-title>

      <q-btn flat round dense icon="close" class="drawer-close-btn" @click="close">
        <AppTooltip content="Close" />
      </q-btn>
    </q-toolbar>

    <q-separator />

    <!-- Content -->
    <div class="drawer-content">
      <q-scroll-area class="fit">
        <!-- Loading State -->
        <div v-if="loading" class="q-pa-lg text-center">
          <q-spinner size="3em" class="q-mb-md" color="primary" />
          <div class="text-grey-7">Loading trigger details...</div>
        </div>

        <!-- Error State -->
        <div v-else-if="error" class="q-pa-lg">
          <q-banner rounded class="bg-negative text-white">
            <template #avatar>
              <q-icon name="error" color="white" />
            </template>
            Failed to load trigger details
          </q-banner>
        </div>

        <!-- Trigger Data -->
        <div v-else-if="trigger" class="q-px-md q-py-lg">

          <!-- System Trigger Warning -->
          <div v-if="trigger.isSystem" class="q-mb-md">
            <q-banner rounded class="bg-warning text-white">
              <template #avatar>
                <q-icon name="lock" color="white" />
              </template>
              This is a system trigger and cannot be modified
            </q-banner>
          </div>

          <!-- Basic Information Section -->
          <div class="section q-mb-md">
            <div class="section-header">
              <q-icon name="info" color="primary" size="sm" class="q-mr-sm" />
              <span class="text-subtitle1 text-weight-medium">Basic Information</span>
            </div>
            <q-separator class="q-my-sm" />

            <!-- Name (full width) -->
            <div class="field-row q-mb-md">
              <div class="field-label">Name</div>
              <div class="field-value text-weight-medium">{{ trigger?.name || '-' }}</div>
            </div>

            <!-- Category & Type (2 columns) -->
            <div class="row q-col-gutter-sm q-mb-md">
              <div class="col-6">
                <div class="field-row">
                  <div class="field-label">Category</div>
                  <div class="field-value">
                    <DetailChip
                      :icon="getCategoryIcon(trigger?.category)"
                      :color="getCategoryColor(trigger?.category)"
                      size="sm"
                      :label="trigger?.category?.toUpperCase() || 'UNKNOWN'"
                    />
                  </div>
                </div>
              </div>
              <div class="col-6">
                <div class="field-row">
                  <div class="field-label">Type</div>
                  <div class="field-value">
                    <DetailChip
                      :icon="getTypeIcon(trigger?.triggerType)"
                      color="blue"
                      size="sm"
                      :label="trigger?.triggerType?.toUpperCase() || 'UNKNOWN'"
                    />
                  </div>
                </div>
              </div>
            </div>

            <!-- Status & Template (2 columns) -->
            <div class="row q-col-gutter-sm q-mb-md">
              <div class="col-6">
                <div class="field-row">
                  <div class="field-label">Status</div>
                  <div class="field-value">
                    <DetailChip
                      :icon="trigger?.enabled ? 'check_circle' : 'cancel'"
                      :color="trigger?.enabled ? 'green' : 'red'"
                      size="sm"
                      :label="trigger?.enabled ? 'ACTIVE' : 'INACTIVE'"
                    />
                  </div>
                </div>
              </div>
              <div class="col-6">
                <div class="field-row">
                  <div class="field-label">Is Template</div>
                  <div class="field-value">
                    <DetailChip
                      :color="trigger?.isTemplate ? 'indigo' : 'grey'"
                      size="sm"
                      :label="trigger?.isTemplate ? 'YES' : 'NO'"
                    />
                  </div>
                </div>
              </div>
            </div>

            <!-- Description (full width) -->
            <div class="field-row">
              <div class="field-label">Description</div>
              <div class="field-value text-grey-8">
                {{ trigger?.description || 'No description provided' }}
              </div>
            </div>
          </div>

          <!-- Configuration Section -->
          <div class="section q-mb-md">
            <div class="section-header">
              <q-icon name="settings" color="primary" size="sm" class="q-mr-sm" />
              <span class="text-subtitle1 text-weight-medium">Configuration</span>
            </div>
            <q-separator class="q-my-sm" />

            <!-- Dynamic config display based on trigger type -->
            <div v-if="configEntries.length > 0">
              <div v-for="(entry, index) in configEntries" :key="index" class="field-row">
                <div class="field-label">{{ formatConfigKey(entry.key) }}</div>
                <div class="field-value">{{ formatConfigValue(entry.value) }}</div>
              </div>
            </div>
            <div v-else class="text-grey-6 text-center q-py-md">
              No configuration available
            </div>
          </div>

          <!-- Timestamps Section -->
          <div class="section">
            <div class="section-header">
              <q-icon name="schedule" color="primary" size="sm" class="q-mr-sm" />
              <span class="text-subtitle1 text-weight-medium">Timestamps</span>
            </div>
            <q-separator class="q-my-sm" />

            <!-- Created & Updated (2 columns) -->
            <div class="row q-col-gutter-sm">
              <div class="col-6">
                <div class="field-row">
                  <div class="field-label">Created</div>
                  <div class="field-value">{{ formatDate(trigger?.created) }}</div>
                </div>
              </div>
              <div class="col-6">
                <div class="field-row">
                  <div class="field-label">Updated</div>
                  <div class="field-value">{{ formatDate(trigger?.updated) }}</div>
                </div>
              </div>
            </div>
          </div>

        </div>
      </q-scroll-area>
    </div>

    <!-- Footer Actions -->
    <q-separator />
    <div class="drawer-footer">
      <q-space />
      <q-btn
        unelevated
        icon="edit"
        color="primary"
        label="Edit"
        :disable="!trigger || isSystemTrigger"
        @click="handleEdit"
      >
        <AppTooltip v-if="isSystemTrigger" content="System triggers cannot be edited" />
      </q-btn>
    </div>
  </q-drawer>
</template>

<script setup lang="ts">
defineOptions({
  name: 'TriggerDetailDrawer'
});

/** TYPE IMPORTS */
import type { TriggerDetailDrawerProps, TriggerDetailDrawerEmits } from './interfaces';
import type { TriggerResponse } from '@mapexos/schemas';

/** VUE IMPORTS */
import { ref, watch, computed, onMounted, onBeforeUnmount } from 'vue';
import { date } from 'quasar';

/** COMPONENTS */
import { DetailChip } from '@components/chips';
import { AppTooltip } from '@components/tooltips';

/** COMPOSABLES */
import { useCommonErrors } from '@composables/i18n';
import { useLogger } from '@composables/useLogger';

/** UTILS */
import { notifyFail } from '@utils/alert';

/** SERVICES */
import { apis } from '@services/mapex';

const errors = useCommonErrors();
const logger = useLogger('TriggerDetailDrawer');

/** PROPS & EMITS */
const props = defineProps<TriggerDetailDrawerProps>();

const emit = defineEmits<TriggerDetailDrawerEmits>();

/** STATE */
const trigger = ref<TriggerResponse | null>(null);
const loading = ref(false);
const error = ref(false);

/** COMPUTED */
const isSystemTrigger = computed(() => {
  return trigger.value?.isSystem === true;
});

/**
 * Extract config entries for display
 */
const configEntries = computed(() => {
  if (!trigger.value?.config) return [];

  const config = trigger.value.config as Record<string, unknown>;
  const entries: { key: string; value: unknown }[] = [];

  // If config is wrapped by trigger type, extract inner config
  const triggerType = trigger.value.triggerType;
  const innerConfig = triggerType && config[triggerType]
    ? config[triggerType] as Record<string, unknown>
    : config;

  for (const [key, value] of Object.entries(innerConfig)) {
    entries.push({ key, value });
  }

  return entries;
});

/** WATCHERS */
watch(() => props.triggerId, (newTriggerId) => {
  if (newTriggerId && props.modelValue) {
    void fetchTriggerDetails(newTriggerId);
  }
}, { immediate: true });

watch(() => props.modelValue, (isOpen) => {
  if (isOpen && props.triggerId) {
    void fetchTriggerDetails(props.triggerId);
  } else if (!isOpen) {
    trigger.value = null;
    error.value = false;
  }
});

/** FUNCTIONS */

/**
 * Fetch trigger details by ID from API
 * @param {string} triggerId - Trigger ID
 * @returns {Promise<void>}
 */
async function fetchTriggerDetails(triggerId: string): Promise<void> {
  if (!apis.triggers) {
    error.value = true;
    notifyFail({ message: errors.apiNotInitialized.value });
    return;
  }

  loading.value = true;
  error.value = false;
  trigger.value = null;

  try {
    const response = await apis.triggers.trigger.getById({ triggerId });
    trigger.value = response;
  } catch (err: unknown) {
    logger.error('Error fetching trigger details:', err);
    error.value = true;
    notifyFail({ message: errors.loadFailed.value });
  } finally {
    loading.value = false;
  }
}

/**
 * Get icon for category
 * @param {string} category - Trigger category
 * @returns {string} Icon name
 */
function getCategoryIcon(category?: string): string {
  switch (category) {
    case 'technical': return 'dns';
    case 'communication': return 'chat';
    default: return 'category';
  }
}

/**
 * Get color for category
 * @param {string} category - Trigger category
 * @returns {string} Color name
 */
function getCategoryColor(category?: string): 'purple' | 'teal' | 'grey' {
  switch (category) {
    case 'technical': return 'purple';
    case 'communication': return 'teal';
    default: return 'grey';
  }
}

/**
 * Get icon for trigger type
 * @param {string} type - Trigger type
 * @returns {string} Icon name
 */
function getTypeIcon(type?: string): string {
  switch (type?.toLowerCase()) {
    case 'http': return 'http';
    case 'mqtt': return 'wifi_tethering';
    case 'rabbitmq': return 'memory';
    case 'nats': return 'cloud';
    case 'websocket': return 'swap_vert';
    case 'email': return 'email';
    case 'teams': return 'groups';
    case 'slack': return 'tag';
    default: return 'flash_on';
  }
}

/**
 * Format config key for display
 * @param {string} key - Config key
 * @returns {string} Formatted key
 */
function formatConfigKey(key: string): string {
  return key
    .replace(/([A-Z])/g, ' $1')
    .replace(/^./, str => str.toUpperCase())
    .trim();
}

/**
 * Format config value for display
 * @param {unknown} value - Config value
 * @returns {string} Formatted value
 */
function formatConfigValue(value: unknown): string {
  if (value === null || value === undefined) return '-';
  if (typeof value === 'string') return value;
  if (typeof value === 'number') return value.toString();
  if (typeof value === 'boolean') return value ? 'Yes' : 'No';
  if (typeof value === 'bigint') return value.toString();
  if (typeof value === 'object') return JSON.stringify(value, null, 2);
  return '-';
}

/**
 * Format date using Quasar date utils
 * @param {unknown} dateValue - Date value to format
 * @returns {string} Formatted date
 */
function formatDate(dateValue: unknown): string {
  if (!dateValue) return '-';

  try {
    const dateObj = typeof dateValue === 'string' ? new Date(dateValue) : dateValue as Date;
    return date.formatDate(dateObj, 'MMM DD, YYYY HH:mm');
  } catch {
    return '-';
  }
}

/**
 * Close drawer
 */
function close(): void {
  emit('update:modelValue', false);
}

/**
 * Handle edit action
 */
function handleEdit(): void {
  if (isSystemTrigger.value) return;
  if (!trigger.value?.id) return;
  emit('edit', trigger.value.id);
  close();
}

/**
 * Handle ESC key to close drawer
 * @param {KeyboardEvent} event - Keyboard event
 */
function handleEscKey(event: KeyboardEvent): void {
  if (event.key === 'Escape' && props.modelValue) {
    close();
  }
}

/** LIFECYCLE HOOKS */
onMounted(() => {
  window.addEventListener('keydown', handleEscKey);
});

onBeforeUnmount(() => {
  window.removeEventListener('keydown', handleEscKey);
});
</script>

<style lang="scss" scoped>
// Flex layout for drawer content
:deep(.q-drawer__content) {
  display: flex;
  flex-direction: column;
  height: 100%;
}

// Drawer Header
.drawer-header {
  flex-shrink: 0;
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

// Backdrop (teleported to body, needs :global) - transparent, just for click detection
:global(.drawer-backdrop) {
  position: fixed;
  top: 0;
  left: 0;
  right: 450px; // Leave space for drawer (450px width)
  bottom: 0;
  background: transparent;
  z-index: 5999; // Below q-drawer (6000)
  cursor: default;
}

// Drawer Content
.drawer-content {
  flex: 1;
  min-height: 0; // Important for flex children with overflow
  overflow: hidden;

  :deep(.q-scrollarea__content) {
    width: 100%;
    max-width: 100%;
    overflow-x: hidden;
  }
}

// Drawer Footer - Fixed at bottom
.drawer-footer {
  flex-shrink: 0;
  display: flex;
  align-items: center;
  padding: 12px 16px;
  background: var(--mapex-header-bg);
  backdrop-filter: blur(10px);
  border-top: 1px solid var(--mapex-header-border);
  box-shadow: 0 -2px 8px var(--mapex-elevation-shadow);
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
