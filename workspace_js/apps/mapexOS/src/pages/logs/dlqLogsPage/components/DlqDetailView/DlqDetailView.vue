<script setup lang="ts">
defineOptions({
  name: 'DlqDetailView'
});

/** TYPE IMPORTS */
import type { EventsDLQResponse } from '@mapexos/schemas';

/** VUE IMPORTS */
import { ref, watch, nextTick, onBeforeUnmount, computed } from 'vue';

/** COMPONENTS */
import { AppTooltip } from '@components/tooltips';
import { AppTabs } from '@components/tabs';

/** COMPOSABLES */
import { useDlqLogsPageTranslations } from '@composables/i18n/pages/logs/dlqLogsPage';

/** UTILS */
import * as monaco from 'monaco-editor';
import { registerMapexMonacoThemes, getMapexMonacoTheme, applyMapexMonacoTheme } from '@utils/monaco-theme';
import { notifyInfo, notifyFail } from '@utils/alert/notify';

/** STORES */
import { useThemeStore } from '@stores/theme';

/** PROPS & EMITS */
const props = defineProps<{
  entry: EventsDLQResponse;
}>();

const emit = defineEmits<{
  close: [];
}>();

/** COMPOSABLES & STORES */
const t = useDlqLogsPageTranslations();
const themeStore = useThemeStore();

/** STATE */
const activeTab = ref('payload');
const monacoContainer = ref<HTMLElement | null>(null);
let monacoEditor: monaco.editor.IStandaloneCodeEditor | null = null;

/** COMPUTED */

/**
 * Tab items for AppTabs with i18n labels
 */
const tabs = computed(() => [
  { name: 'payload', label: t.detail.tabs.payload.value, icon: 'data_object' },
  { name: 'headers', label: t.detail.tabs.headers.value, icon: 'list_alt' },
  { name: 'error', label: t.detail.tabs.error.value, icon: 'bug_report' },
]);

/**
 * Parse originalData as formatted JSON
 * @returns {string} Pretty-printed JSON or raw string
 */
const payloadContent = computed(() => {
  try {
    const parsed = JSON.parse(props.entry.originalData);
    return JSON.stringify(parsed, null, 2);
  } catch {
    return props.entry.originalData || t.detail.fallback.empty.value;
  }
});

/**
 * Parse originalHeaders as formatted JSON
 * @returns {string} Pretty-printed JSON or raw string
 */
const headersContent = computed(() => {
  if (!props.entry.originalHeaders) return t.detail.fallback.noHeaders.value;
  try {
    const parsed = JSON.parse(props.entry.originalHeaders);
    return JSON.stringify(parsed, null, 2);
  } catch {
    return props.entry.originalHeaders;
  }
});

/**
 * Get the content for the active tab
 * @returns {string} Content to display in Monaco
 */
const activeContent = computed(() => {
  switch (activeTab.value) {
    case 'payload': return payloadContent.value;
    case 'headers': return headersContent.value;
    case 'error': return props.entry.lastError || t.detail.fallback.noErrorMessage.value;
    default: return payloadContent.value;
  }
});

/**
 * Get the Monaco language for the active tab
 * @returns {string} Language identifier
 */
const activeLanguage = computed(() => {
  if (activeTab.value === 'error') return 'plaintext';
  return 'json';
});

/** FUNCTIONS */

/**
 * Format a full timestamp
 * @param {string} value - ISO timestamp
 * @returns {string} Formatted date string
 */
function formatFullTimestamp(value: string): string {
  if (!value) return t.defaults.notAvailable.value;
  const date = new Date(value);
  return date.toLocaleDateString('en-US', {
    day: '2-digit',
    month: 'short',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
  });
}

/**
 * Calculate duration between two timestamps
 * @param {string} start - Start ISO timestamp
 * @param {string} end - End ISO timestamp
 * @returns {string} Human-readable duration
 */
function formatDuration(start: string, end: string): string {
  if (!start || !end) return t.defaults.notAvailable.value;
  const diffMs = new Date(end).getTime() - new Date(start).getTime();
  const diffSecs = Math.floor(diffMs / 1000);
  const diffMins = Math.floor(diffSecs / 60);
  const diffHours = Math.floor(diffMins / 60);
  const diffDays = Math.floor(diffHours / 24);

  if (diffDays > 0) return `${diffDays}d ${diffHours % 24}h`;
  if (diffHours > 0) return `${diffHours}h ${diffMins % 60}m`;
  if (diffMins > 0) return `${diffMins}m ${diffSecs % 60}s`;
  return `${diffSecs}s`;
}

/**
 * Copy current Monaco content to clipboard
 * @returns {Promise<void>}
 */
async function copyContent(): Promise<void> {
  const text = monacoEditor?.getValue() ?? activeContent.value;
  try {
    await navigator.clipboard.writeText(text);
    notifyInfo({ message: t.messages.copied.value });
  } catch {
    notifyFail({ message: t.messages.copyFailed.value });
  }
}

/**
 * Create or update the Monaco editor instance
 */
function setupMonaco(): void {
  if (!monacoContainer.value) return;

  if (monacoEditor) {
    monacoEditor.dispose();
    monacoEditor = null;
  }

  registerMapexMonacoThemes();

  monacoEditor = monaco.editor.create(monacoContainer.value, {
    value: activeContent.value,
    language: activeLanguage.value,
    theme: getMapexMonacoTheme(themeStore.isDark),
    readOnly: true,
    automaticLayout: true,
    minimap: { enabled: false },
    fontSize: 13,
    lineNumbers: 'on',
    wordWrap: 'on',
    scrollBeyondLastLine: false,
    renderLineHighlight: 'none',
    folding: true,
  });

  void monacoEditor.getAction('editor.action.formatDocument')?.run();
}

/** WATCHERS */

/** Watch entry + tab changes to re-create editor */
watch(
  [() => props.entry.id, activeTab],
  async () => {
    await nextTick();
    setupMonaco();
  },
  { immediate: true },
);

/** Watch theme changes */
watch(() => themeStore.isDark, (isDark: boolean) => {
  applyMapexMonacoTheme(isDark);
});

/** LIFECYCLE HOOKS */
onBeforeUnmount(() => {
  if (monacoEditor) {
    monacoEditor.dispose();
    monacoEditor = null;
  }
});
</script>

<template>
  <div class="dlq-detail column full-height">

    <!-- ═══ Header Bar ═══ -->
    <div class="dlq-detail__header">
      <div class="row items-center q-pa-md">
        <q-btn
          flat
          round
          dense
          icon="arrow_back"
          class="dlq-detail__back-btn q-mr-sm"
          @click="emit('close')"
        >
          <AppTooltip :content="t.detail.backToList.value" />
        </q-btn>
        <div class="col">
          <div class="dlq-detail__title">{{ entry.serviceName }}</div>
          <div class="dlq-detail__subtitle">{{ entry.eventType }} &middot; {{ entry.consumerName }}</div>
        </div>
        <q-btn
          flat
          round
          dense
          icon="content_copy"
          class="dlq-detail__action-btn"
          @click="copyContent"
        >
          <AppTooltip :content="t.detail.copyContent.value" />
        </q-btn>
      </div>
    </div>

    <!-- ═══ Info Card ═══ -->
    <div class="dlq-detail__info-card q-ma-md">
      <!-- Error message -->
      <div class="dlq-detail__error-row">
        <q-icon name="error" size="sm" class="dlq-detail__error-icon" />
        <span class="dlq-detail__error-text">{{ entry.lastError || t.detail.unknownError.value }}</span>
      </div>

      <q-separator class="dlq-detail__separator" />

      <!-- Stats row -->
      <div class="dlq-detail__stats-row">
        <div class="dlq-detail__stat">
          <span class="dlq-detail__stat-label">{{ t.detail.errors.value }}</span>
          <span class="dlq-detail__stat-value dlq-detail__stat-value--danger">{{ entry.errorCount }}</span>
        </div>
        <div class="dlq-detail__stat">
          <span class="dlq-detail__stat-label">{{ t.detail.deliveries.value }}</span>
          <span class="dlq-detail__stat-value">{{ entry.totalDeliveries }}</span>
        </div>
        <div class="dlq-detail__stat">
          <span class="dlq-detail__stat-label">{{ t.detail.duration.value }}</span>
          <span class="dlq-detail__stat-value">{{ formatDuration(entry.firstDelivery, entry.lastDelivery) }}</span>
        </div>
        <div class="dlq-detail__stat">
          <span class="dlq-detail__stat-label">{{ t.detail.retention.value }}</span>
          <span class="dlq-detail__stat-value">{{ entry.retentionDays }}d</span>
        </div>
      </div>

      <q-separator class="dlq-detail__separator" />

      <!-- Context row -->
      <div class="dlq-detail__context-row">
        <div class="dlq-detail__context-item">
          <q-icon name="schedule" size="xs" class="dlq-detail__context-icon" />
          <span>{{ formatFullTimestamp(entry.firstDelivery) }}</span>
          <span class="dlq-detail__context-arrow">→</span>
          <span>{{ formatFullTimestamp(entry.lastDelivery) }}</span>
        </div>
        <div class="dlq-detail__context-badges">
          <q-badge outline class="dlq-detail__badge dlq-detail__badge--primary" :label="entry.originalStream" />
          <q-badge outline class="dlq-detail__badge" :label="entry.originalSubject" />
          <q-badge outline class="dlq-detail__badge" :label="entry.serviceType" />
        </div>
      </div>
    </div>

    <!-- ═══ Tabs + Editor ═══ -->
    <div class="q-mx-md">
      <AppTabs v-model="activeTab" :tabs="tabs" :separator="false" />
    </div>

    <div class="col q-pa-none dlq-detail__editor-body">
      <div ref="monacoContainer" class="dlq-detail__monaco-container" />
    </div>

    <!-- ═══ Footer ═══ -->
    <div class="dlq-detail__footer">
      <q-icon name="fingerprint" size="xs" class="q-mr-2xs" />
      ID: {{ entry.id }}
      <q-space />
      {{ t.detail.searchHint.value }}
    </div>
  </div>
</template>

<style lang="scss" scoped>
.dlq-detail {
  background: var(--mapex-page-bg);

  // ── Header ──
  &__header {
    background: var(--mapex-surface-elevated);
    border-bottom: 1px solid var(--mapex-card-border);
    flex-shrink: 0;
  }

  &__back-btn {
    color: var(--mapex-text-secondary);
  }

  &__title {
    font-size: var(--mapex-font-lg);
    font-weight: var(--mapex-font-weight-bold);
    color: var(--mapex-text-primary);
  }

  &__subtitle {
    font-size: var(--mapex-font-xs);
    color: var(--mapex-text-muted);
  }

  &__action-btn {
    color: var(--mapex-text-secondary);
  }

  // ── Info Card ──
  &__info-card {
    background: var(--mapex-surface-elevated);
    border: 1px solid var(--mapex-card-border);
    border-radius: var(--mapex-radius-lg);
    box-shadow: var(--mapex-shadow-xs);
    flex-shrink: 0;
  }

  &__separator {
    background: var(--mapex-divider);
  }

  // Error row
  &__error-row {
    display: flex;
    align-items: flex-start;
    gap: var(--mapex-spacing-sm);
    padding: var(--mapex-spacing-md) var(--mapex-spacing-lg);
  }

  &__error-icon {
    color: var(--mapex-danger);
    flex-shrink: 0;
    margin-top: 2px;
  }

  &__error-text {
    font-size: var(--mapex-font-sm);
    color: var(--mapex-text-primary);
    line-height: var(--mapex-line-height-base);
    word-break: break-word;
  }

  // Stats row
  &__stats-row {
    display: flex;
    align-items: center;
    gap: var(--mapex-spacing-2xl);
    padding: var(--mapex-spacing-sm) var(--mapex-spacing-lg);
  }

  &__stat {
    display: flex;
    align-items: center;
    gap: var(--mapex-spacing-xs);
  }

  &__stat-label {
    font-size: var(--mapex-font-xs);
    color: var(--mapex-text-muted);
  }

  &__stat-value {
    font-size: var(--mapex-font-sm);
    font-weight: var(--mapex-font-weight-bold);
    color: var(--mapex-text-primary);

    &--danger {
      color: var(--mapex-danger);
    }
  }

  // Context row
  &__context-row {
    display: flex;
    align-items: center;
    justify-content: space-between;
    flex-wrap: wrap;
    gap: var(--mapex-spacing-sm);
    padding: var(--mapex-spacing-sm) var(--mapex-spacing-lg);
  }

  &__context-item {
    display: flex;
    align-items: center;
    gap: var(--mapex-spacing-xs);
    font-size: var(--mapex-font-xs);
    color: var(--mapex-text-muted);
  }

  &__context-icon {
    color: var(--mapex-text-muted);
  }

  &__context-arrow {
    color: var(--mapex-text-muted);
    margin: 0 var(--mapex-spacing-2xs);
  }

  &__context-badges {
    display: flex;
    gap: var(--mapex-spacing-xs);
  }

  &__badge {
    border-color: var(--mapex-card-border);
    color: var(--mapex-text-secondary);
    font-size: var(--mapex-font-2xs);

    &--primary {
      border-color: var(--mapex-active-border);
      color: var(--mapex-primary);
    }
  }


  // ── Editor ──
  &__editor-body {
    background-color: var(--mapex-surface-bg);
  }

  &__monaco-container {
    width: 100%;
    height: 100%;
    min-height: 200px;
  }

  // ── Footer ──
  &__footer {
    display: flex;
    align-items: center;
    padding: var(--mapex-spacing-sm) var(--mapex-spacing-lg);
    font-size: var(--mapex-font-xs);
    color: var(--mapex-text-muted);
    background: var(--mapex-surface-elevated);
    border-top: 1px solid var(--mapex-card-border);
    flex-shrink: 0;
  }
}
</style>
