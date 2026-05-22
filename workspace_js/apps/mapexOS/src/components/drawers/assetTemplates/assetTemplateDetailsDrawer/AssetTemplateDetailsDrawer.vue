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
      <q-icon name="memory" size="sm" class="q-mr-sm" color="primary" />
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
          <q-spinner size="3em" class="q-mb-md" color="primary" />
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

        <!-- Template Data -->
        <div v-else-if="template" class="q-px-md q-py-lg">

          <!-- System Template Warning -->
          <div v-if="template.isSystem" class="q-mb-md">
            <q-banner rounded class="bg-warning text-white">
              <template #avatar>
                <q-icon name="lock" color="white" />
              </template>
              {{ t.drawer.systemTemplateWarning.value }}
            </q-banner>
          </div>

          <!-- Basic Information Section -->
          <div class="section q-mb-md">
            <div class="section-header">
              <q-icon name="info" color="primary" size="sm" class="q-mr-sm" />
              <span class="text-subtitle1 text-weight-medium">{{ t.drawer.sections.basicInfo.value }}</span>
            </div>
            <q-separator class="q-my-sm" />

            <!-- Name (full width) -->
            <div class="field-row q-mb-md">
              <div class="field-label">{{ t.drawer.fields.name.value }}</div>
              <div class="field-value text-weight-medium">{{ template?.name || '-' }}</div>
            </div>

            <!-- Status & Is System (2 columns) -->
            <div class="row q-col-gutter-sm q-mb-md">
              <div class="col-6">
                <div class="field-row">
                  <div class="field-label">{{ t.drawer.fields.status.value }}</div>
                  <div class="field-value">
                    <DetailChip
                      :color="template?.enabled ? 'positive' : 'negative'"
                      size="sm"
                      :label="template?.enabled ? t.status.active.value.toUpperCase() : t.status.inactive.value.toUpperCase()"
                    />
                  </div>
                </div>
              </div>
              <div class="col-6">
                <div class="field-row">
                  <div class="field-label">{{ t.drawer.fields.isSystem.value }}</div>
                  <div class="field-value">
                    <DetailChip
                      :icon="template?.isSystem ? 'lock' : 'lock_open'"
                      :color="template?.isSystem ? 'orange' : 'grey'"
                      size="sm"
                      :label="template?.isSystem ? t.drawer.system.yes.value.toUpperCase() : t.drawer.system.no.value.toUpperCase()"
                    />
                  </div>
                </div>
              </div>
            </div>

            <!-- Description (full width) -->
            <div class="field-row">
              <div class="field-label">{{ t.drawer.fields.description.value }}</div>
              <div class="field-value text-grey-8">
                {{ template?.description || t.drawer.empty.description.value }}
              </div>
            </div>
          </div>

          <!-- Configuration Section -->
          <div class="section q-mb-md">
            <div class="section-header">
              <q-icon name="settings" color="primary" size="sm" class="q-mr-sm" />
              <span class="text-subtitle1 text-weight-medium">{{ t.drawer.sections.configuration.value }}</span>
            </div>
            <q-separator class="q-my-sm" />

            <!-- Manufacturer & Model (2 columns) -->
            <div class="row q-col-gutter-sm q-mb-md">
              <div class="col-6">
                <div class="field-row">
                  <div class="field-label">{{ t.drawer.fields.manufacturer.value }}</div>
                  <div class="field-value">
                    <DetailChip
                      icon="factory"
                      color="blue"
                      size="sm"
                      :label="template?.manufacturerName || '-'"
                    />
                  </div>
                </div>
              </div>
              <div class="col-6">
                <div class="field-row">
                  <div class="field-label">{{ t.drawer.fields.model.value }}</div>
                  <div class="field-value">
                    <DetailChip
                      icon="router"
                      color="indigo"
                      size="sm"
                      :label="template?.modelName || '-'"
                    />
                  </div>
                </div>
              </div>
            </div>

            <!-- Version (full width) -->
            <div class="field-row q-mb-md">
              <div class="field-label">{{ t.drawer.fields.version.value }}</div>
              <div class="field-value">
                <DetailChip
                  icon="label"
                  color="purple"
                  size="sm"
                  :label="template?.version || '-'"
                />
              </div>
            </div>

            <!-- Asset ID Path (full width) -->
            <div class="field-row">
              <div class="field-label">{{ t.drawer.fields.assetIdPath.value }}</div>
              <div class="field-value text-grey-8">
                <code class="bg-grey-2 q-pa-xs rounded-borders">{{ template?.assetIdPath || '-' }}</code>
              </div>
            </div>
          </div>

          <!-- Scripts Section -->
          <div class="section q-mb-md">
            <div class="section-header">
              <q-icon name="code" color="primary" size="sm" class="q-mr-sm" />
              <span class="text-subtitle1 text-weight-medium">{{ t.drawer.sections.scripts.value }}</span>
            </div>
            <q-separator class="q-my-sm" />

            <!-- Script Test & Script Processor (2 columns) -->
            <div class="row q-col-gutter-sm q-mb-md">
              <div class="col-6">
                <div class="field-row">
                  <div class="field-label">{{ t.drawer.fields.scriptTest.value }}</div>
                  <div class="field-value row items-center q-gutter-xs">
                    <DetailChip
                      :icon="hasScript('scriptTest') ? 'check_circle' : 'cancel'"
                      :color="hasScript('scriptTest') ? 'green' : 'grey'"
                      size="sm"
                      :label="hasScript('scriptTest') ? t.drawer.scripts.configured.value : t.drawer.scripts.notConfigured.value"
                    />
                    <q-btn
                      v-if="hasScript('scriptTest')"
                      flat
                      dense
                      round
                      size="sm"
                      icon="visibility"
                      color="primary"
                      @click="viewScript('scriptTest', t.drawer.fields.scriptTest.value)"
                    >
                      <AppTooltip :content="t.drawer.scriptViewer.viewScript.value" />
                    </q-btn>
                  </div>
                </div>
              </div>
              <div class="col-6">
                <div class="field-row">
                  <div class="field-label">{{ t.drawer.fields.scriptProcessor.value }}</div>
                  <div class="field-value row items-center q-gutter-xs">
                    <DetailChip
                      :icon="hasScript('scriptProcessor') ? 'check_circle' : 'cancel'"
                      :color="hasScript('scriptProcessor') ? 'green' : 'grey'"
                      size="sm"
                      :label="hasScript('scriptProcessor') ? t.drawer.scripts.configured.value : t.drawer.scripts.notConfigured.value"
                    />
                    <q-btn
                      v-if="hasScript('scriptProcessor')"
                      flat
                      dense
                      round
                      size="sm"
                      icon="visibility"
                      color="primary"
                      @click="viewScript('scriptProcessor', t.drawer.fields.scriptProcessor.value)"
                    >
                      <AppTooltip :content="t.drawer.scriptViewer.viewScript.value" />
                    </q-btn>
                  </div>
                </div>
              </div>
            </div>

            <!-- Script Validator & Script Conversion (2 columns) -->
            <div class="row q-col-gutter-sm q-mb-md">
              <div class="col-6">
                <div class="field-row">
                  <div class="field-label">{{ t.drawer.fields.scriptValidator.value }}</div>
                  <div class="field-value row items-center q-gutter-xs">
                    <DetailChip
                      :icon="hasScript('scriptValidator') ? 'check_circle' : 'cancel'"
                      :color="hasScript('scriptValidator') ? 'green' : 'grey'"
                      size="sm"
                      :label="hasScript('scriptValidator') ? t.drawer.scripts.configured.value : t.drawer.scripts.notConfigured.value"
                    />
                    <q-btn
                      v-if="hasScript('scriptValidator')"
                      flat
                      dense
                      round
                      size="sm"
                      icon="visibility"
                      color="primary"
                      @click="viewScript('scriptValidator', t.drawer.fields.scriptValidator.value)"
                    >
                      <AppTooltip :content="t.drawer.scriptViewer.viewScript.value" />
                    </q-btn>
                  </div>
                </div>
              </div>
              <div class="col-6">
                <div class="field-row">
                  <div class="field-label">{{ t.drawer.fields.scriptConversion.value }}</div>
                  <div class="field-value row items-center q-gutter-xs">
                    <DetailChip
                      :icon="hasScript('scriptConversion') ? 'check_circle' : 'cancel'"
                      :color="hasScript('scriptConversion') ? 'green' : 'grey'"
                      size="sm"
                      :label="hasScript('scriptConversion') ? t.drawer.scripts.configured.value : t.drawer.scripts.notConfigured.value"
                    />
                    <q-btn
                      v-if="hasScript('scriptConversion')"
                      flat
                      dense
                      round
                      size="sm"
                      icon="visibility"
                      color="primary"
                      @click="viewScript('scriptConversion', t.drawer.fields.scriptConversion.value)"
                    >
                      <AppTooltip :content="t.drawer.scriptViewer.viewScript.value" />
                    </q-btn>
                  </div>
                </div>
              </div>
            </div>

            <!-- Script Summary (full width) -->
            <div class="field-row">
              <div class="field-label">{{ t.drawer.fields.scriptsSummary.value }}</div>
              <div class="field-value">
                <DetailChip
                  icon="code"
                  :color="getScriptSummaryColorName()"
                  size="sm"
                  :label="getScriptSummary()"
                />
              </div>
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
                  <div class="field-value">{{ formatDate(template?.created) }}</div>
                </div>
              </div>
              <div class="col-6">
                <div class="field-row">
                  <div class="field-label">{{ t.drawer.fields.updated.value }}</div>
                  <div class="field-value">{{ formatDate(template?.updated) }}</div>
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
        :label="t.drawer.edit.value"
        :disable="!template || isSystemTemplate"
        @click="handleEdit"
      >
        <AppTooltip v-if="isSystemTemplate" :content="t.drawer.systemTemplateTooltip.value" />
      </q-btn>
    </div>
  </q-drawer>

  <!-- Script Viewer Dialog -->
  <ScriptViewerDialog
    v-model="showScriptViewer"
    :title="currentScriptTitle"
    :script-content="currentScriptContent"
    :language="currentScriptLanguage"
    :copy-tooltip="t.drawer.scriptViewer.copyScript.value"
    :close-tooltip="t.drawer.scriptViewer.close.value"
    :copy-success-message="t.drawer.scriptViewer.copySuccess.value"
    :copy-fail-message="t.drawer.scriptViewer.copyFail.value"
  />
</template>

<script setup lang="ts">
defineOptions({
  name: 'AssetTemplateDetailsDrawer'
});

/** TYPE IMPORTS (ALL types first, grouped) */
import type { AssetTemplateDetailsDrawerProps, AssetTemplateDetailsDrawerEmits } from './interfaces/assetTemplateDetailsDrawer.interface';
import type { AssetTemplateResponse } from '@mapexos/schemas';

/** VUE IMPORTS */
import { ref, watch, computed, onMounted, onBeforeUnmount } from 'vue';
import { date } from 'quasar';

/** COMPONENTS */
import { ScriptViewerDialog } from '@components/dialogs/scriptViewer';
import { DetailChip } from '@components/chips';
import { AppTooltip } from '@components/tooltips';

/** COMPOSABLES */
import { useAssetTemplatesTranslations } from '@composables/i18n';
import { useLogger } from '@composables/useLogger';

/** UTILS */
import { notifyFail } from '@utils/alert';

/** SERVICES */
import { apis } from '@services/mapex';

/** PROPS & EMITS */
const props = defineProps<AssetTemplateDetailsDrawerProps>();

const emit = defineEmits<AssetTemplateDetailsDrawerEmits>();

/** COMPOSABLES & STORES */
const t = useAssetTemplatesTranslations();
const logger = useLogger('AssetTemplateDetailsDrawer');

/** STATE */
const template = ref<AssetTemplateResponse | null>(null);
const loading = ref(false);
const error = ref(false);

const showScriptViewer = ref(false);
const currentScriptTitle = ref('');
const currentScriptContent = ref('');
const currentScriptLanguage = ref<'javascript' | 'json'>('javascript');

/** COMPUTED */
const isSystemTemplate = computed(() => {
  return template.value?.isSystem === true;
});

/** WATCHERS */
watch(() => props.templateId, (newTemplateId) => {
  if (newTemplateId && props.modelValue) {
    void fetchTemplateDetails(newTemplateId);
  }
}, { immediate: true });

watch(() => props.modelValue, (isOpen) => {
  if (isOpen && props.templateId) {
    void fetchTemplateDetails(props.templateId);
  } else if (!isOpen) {
    template.value = null;
    error.value = false;
  }
});

/** FUNCTIONS */

/**
 * Fetch asset template details by ID from API
 * @param {string} templateId - Template ID
 * @returns {Promise<void>}
 */
async function fetchTemplateDetails(templateId: string): Promise<void> {
  if (!apis.assets) {
    error.value = true;
    notifyFail({ message: 'Assets API not initialized' });
    return;
  }

  loading.value = true;
  error.value = false;
  template.value = null;

  try {
    const response = await apis.assets.assetTemplate.getById({ assetTemplateId: templateId });
    template.value = response;
  } catch (err: any) {
    logger.error('Error fetching asset template details:', err);
    error.value = true;
    notifyFail({ message: t.drawer.error.value });
  } finally {
    loading.value = false;
  }
}

/**
 * Check if a script is configured
 * @param {string} scriptKey - Script key to check
 * @returns {boolean} True if script exists
 */
function hasScript(scriptKey: string): boolean {
  return !!(template.value as any)?.[scriptKey];
}

/**
 * Get script summary count
 * @returns {string} Summary text
 */
function getScriptSummary(): string {
  let count = 0;
  if (hasScript('scriptTest')) count++;
  if (hasScript('scriptProcessor')) count++;
  if (hasScript('scriptValidator')) count++;
  if (hasScript('scriptConversion')) count++;
  return `${count}/4 ${t.drawer.scripts.configured.value.toLowerCase()}`;
}

/**
 * Get script summary color name (DetailChip format)
 * @returns {string} Color name compatible with DetailChip
 */
function getScriptSummaryColorName(): 'green' | 'orange' | 'grey' {
  let count = 0;
  if (hasScript('scriptTest')) count++;
  if (hasScript('scriptProcessor')) count++;
  if (hasScript('scriptValidator')) count++;
  if (hasScript('scriptConversion')) count++;

  if (count === 4) return 'green';
  if (count >= 2) return 'orange';
  return 'grey';
}

/**
 * Format date using Quasar date utils
 * @param {any} dateValue - Date value to format
 * @returns {string} Formatted date
 */
function formatDate(dateValue: any): string {
  if (!dateValue) return '-';

  try {
    const dateObj = typeof dateValue === 'string' ? new Date(dateValue) : dateValue;
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
  if (isSystemTemplate.value) return;
  if (!template.value?.id) return;
  emit('edit', template.value.id);
  close();
}

/**
 * Open script viewer
 * @param {string} scriptKey - Script key
 * @param {string} scriptTitle - Script title
 * @param {'javascript' | 'json'} language - Script language
 */
function viewScript(scriptKey: string, scriptTitle: string, language: 'javascript' | 'json' = 'javascript'): void {
  const scriptContent = (template.value as any)?.[scriptKey] || '';
  if (!scriptContent) return;

  currentScriptTitle.value = scriptTitle;
  currentScriptContent.value = scriptContent;
  currentScriptLanguage.value = language;
  showScriptViewer.value = true;
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

    code {
      font-family: 'Courier New', monospace;
      font-size: 0.85rem;
    }
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
