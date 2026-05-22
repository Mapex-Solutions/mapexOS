<script setup lang="ts">
defineOptions({
  name: 'DlqListView'
});

/** TYPE IMPORTS */
import type { EventsDLQResponse } from '@mapexos/schemas';

/** VUE IMPORTS */
import { ref } from 'vue';

/** COMPOSABLES */
import { useDlqLogsPageTranslations } from '@composables/i18n/pages/logs/dlqLogsPage';

/** CONSTANTS */
// Service type visual config — data-driven colors (not design tokens)
const SERVICE_CONFIG: Record<string, { icon: string; color: string }> = {
  workflow: { icon: 'account_tree', color: '#7E57C2' },
  triggers: { icon: 'flash_on', color: '#EF6C00' },
  router: { icon: 'route', color: '#1E88E5' },
  events: { icon: 'event_note', color: '#00897B' },
  assets: { icon: 'devices', color: '#43A047' },
  'mapex-iam': { icon: 'admin_panel_settings', color: '#5C6BC0' },
  'http-gateway': { icon: 'http', color: '#00ACC1' },
  'js-executor': { icon: 'code', color: '#F9A825' },
  'js-workflow-executor': { icon: 'terminal', color: '#AB47BC' },
};

const DEFAULT_CONFIG = { icon: 'dns', color: 'var(--mapex-text-muted)' };

/** PROPS & EMITS */
defineProps<{
  entries: EventsDLQResponse[];
  selectedId: string | null;
  hasMore: boolean;
}>();

const emit = defineEmits<{
  select: [entry: EventsDLQResponse];
  loadMore: [done: (stop?: boolean) => void];
}>();

/** COMPOSABLES & STORES */
const t = useDlqLogsPageTranslations();

/** STATE */
const scrollTargetRef = ref<HTMLElement | null>(null);

/** FUNCTIONS */

/**
 * Get config for a service type
 * @param {string} serviceType - The service type string
 * @returns {{ icon: string; color: string }}
 */
function getServiceConfig(serviceType: string): { icon: string; color: string } {
  return SERVICE_CONFIG[serviceType?.toLowerCase()] || DEFAULT_CONFIG;
}

/**
 * Format timestamp to short relative display
 * @param {string} value - ISO timestamp
 * @returns {string} Formatted time string
 */
function formatTime(value: string): string {
  if (!value) return '';
  const date = new Date(value);
  const now = new Date();
  const diffMs = now.getTime() - date.getTime();
  const diffMins = Math.floor(diffMs / 60000);
  const diffHours = Math.floor(diffMs / 3600000);
  const diffDays = Math.floor(diffMs / 86400000);

  if (diffMins < 1) return t.time.justNow.value;
  if (diffMins < 60) return `${diffMins}m`;
  if (diffHours < 24) return `${diffHours}h`;
  if (diffDays < 7) return `${diffDays}d`;
  return date.toLocaleDateString('en-US', { day: '2-digit', month: 'short' });
}

/**
 * Truncate error message for preview
 * @param {string} error - Full error message
 * @returns {string} Truncated error
 */
function truncateError(error: string): string {
  if (!error) return t.list.unknownError.value;
  return error.length > 90 ? error.substring(0, 90) + '...' : error;
}

/**
 * Handle infinite scroll load event from Quasar
 * @param {number} _index - Page index (unused)
 * @param {Function} done - Callback to signal completion
 */
function onLoad(_index: number, done: (stop?: boolean) => void): void {
  emit('loadMore', done);
}
</script>

<template>
  <div ref="scrollTargetRef" class="dlq-list">
    <q-infinite-scroll
      :scroll-target="scrollTargetRef || undefined"
      :offset="250"
      @load="onLoad"
    >
      <!-- Entries -->
      <div
        v-for="entry in entries"
        :key="entry.id"
        class="dlq-list__row"
        :class="{ 'dlq-list__row--selected': selectedId === entry.id }"
        @click="emit('select', entry)"
      >
        <!-- Avatar -->
        <div
          class="dlq-list__avatar"
          :style="{ background: getServiceConfig(entry.serviceType).color }"
        >
          <q-icon :name="getServiceConfig(entry.serviceType).icon" size="16px" color="white" />
          <span v-if="entry.errorCount > 1" class="dlq-list__avatar-badge">{{ entry.errorCount }}</span>
        </div>

        <!-- Content -->
        <div class="dlq-list__content">
          <div class="dlq-list__service-name">{{ entry.serviceName }}</div>
          <div class="dlq-list__error-preview">{{ truncateError(entry.lastError) }}</div>
        </div>

        <!-- Time -->
        <div class="dlq-list__time">{{ formatTime(entry.created) }}</div>
      </div>

      <!-- Loading Spinner -->
      <template #loading>
        <div class="row justify-center q-pa-md">
          <q-spinner color="primary" size="24px" />
        </div>
      </template>
    </q-infinite-scroll>

    <!-- End of list -->
    <div v-if="!hasMore && entries.length > 0" class="dlq-list__end">
      {{ entries.length }} {{ t.itemLabelPlural.value }}
    </div>
  </div>
</template>

<style lang="scss" scoped>
.dlq-list {
  flex: 1;
  min-height: 0;
  overflow-y: auto;
  padding-top: var(--mapex-spacing-md);

  &__row {
    display: flex;
    align-items: center;
    gap: var(--mapex-spacing-md);
    padding: var(--mapex-spacing-md) var(--mapex-spacing-lg);
    background: var(--mapex-surface-bg);
    border: 1px solid var(--mapex-card-border);
    border-radius: var(--mapex-radius-xs);
    margin: 0 var(--mapex-spacing-md) 8px var(--mapex-spacing-md);
    cursor: pointer;
    transition: var(--mapex-transition-base);

    &:hover {
      background-color: var(--mapex-surface-elevated);
      box-shadow: 0 2px 6px var(--mapex-elevation-shadow);
    }

    &--selected {
      background: var(--mapex-active-bg) !important;
    }
  }

  &__avatar {
    width: 36px;
    height: 36px;
    border-radius: var(--mapex-radius-md);
    display: flex;
    align-items: center;
    justify-content: center;
    flex-shrink: 0;
    position: relative;
  }

  &__avatar-badge {
    position: absolute;
    top: -4px;
    right: -4px;
    min-width: 16px;
    height: 16px;
    border-radius: 8px;
    background: var(--mapex-danger);
    color: var(--mapex-surface-bg);
    font-size: 9px;
    font-weight: var(--mapex-font-weight-bold);
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 0 3px;
  }

  &__content {
    flex: 1;
    min-width: 0;
  }

  &__service-name {
    font-size: var(--mapex-font-sm);
    font-weight: var(--mapex-font-weight-medium);
    color: var(--mapex-text-primary);
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  &__error-preview {
    font-size: var(--mapex-font-xs);
    color: var(--mapex-text-muted);
    margin-top: 2px;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  &__time {
    font-size: var(--mapex-font-xs);
    color: var(--mapex-text-muted);
    flex-shrink: 0;
    font-variant-numeric: tabular-nums;
  }

  &__end {
    text-align: center;
    padding: var(--mapex-spacing-sm);
    font-size: var(--mapex-font-xs);
    color: var(--mapex-text-muted);
  }
}
</style>
