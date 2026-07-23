<script setup lang="ts">
defineOptions({
  name: 'AssetTemplateDetailsDrawer'
});

/** TYPE IMPORTS */
import type { AssetTemplateDetailsDrawerProps, AssetTemplateDetailsDrawerEmits } from './interfaces/assetTemplateDetailsDrawer.interface';
import type { AssetTemplateResponse } from '@mapexos/schemas';

/** VUE IMPORTS */
import { ref, watch, computed, onMounted, onBeforeUnmount } from 'vue';
import { date } from 'quasar';

/** COMPONENTS */
import { ScriptViewerDialog } from '@components/dialogs/scriptViewer';
import { ContentModal } from '@components/dialogs/common';
import { DynamicFieldsTable, AvailableFieldsList } from '@components/assetTemplates';
import { DetailChip } from '@components/chips';
import { AppTooltip } from '@components/tooltips';
import { InfoBanner } from '@components/banners';
import { BaseButton } from '@components/buttons';

/** COMPOSABLES */
import { useAssetTemplatesTranslations, useCommonErrors } from '@composables/i18n';
import { useLogger } from '@composables/useLogger';

/** UTILS */
import { notifyFail, notifySuccess } from '@utils/alert';

/** SERVICES */
import { apis } from '@services/mapex';

/** PROPS & EMITS */
const props = defineProps<AssetTemplateDetailsDrawerProps>();
const emit = defineEmits<AssetTemplateDetailsDrawerEmits>();

/** COMPOSABLES & STORES */
const t = useAssetTemplatesTranslations();
const errors = useCommonErrors();
const logger = useLogger('AssetTemplateDetailsDrawer');

/** STATE */
const template = ref<AssetTemplateResponse | null>(null);
const loading = ref(false);
const error = ref(false);

const showScriptViewer = ref(false);
const currentScriptTitle = ref('');
const currentScriptContent = ref('');
const currentScriptLanguage = ref<'javascript' | 'json'>('javascript');

const showFieldsModal = ref(false);
const fieldsModalKind = ref<'dynamic' | 'available' | null>(null);

const cloning = ref(false);

/** COMPUTED */
const isSystemTemplate = computed(() => template.value?.isSystem === true);

// Marketplace templates are read-only; the user clones them to a local copy to edit.
const isMarketplaceTemplate = computed(() => template.value?.isMarketplace === true);

// Editing is blocked for both platform-owned system templates and marketplace templates.
const isReadonlyTemplate = computed(() => isSystemTemplate.value || isMarketplaceTemplate.value);

const dynamicFieldsCount = computed(() => template.value?.dynamicFields?.length ?? 0);

const availableFieldsCount = computed(() => template.value?.availableFields?.length ?? 0);

const fieldsModalTitle = computed(() =>
  fieldsModalKind.value === 'dynamic'
    ? t.drawer.fieldsModal.dynamicTitle.value
    : t.drawer.fieldsModal.availableTitle.value
);

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
 * Fetch asset template details by ID from the API.
 * @param {string} templateId - Template ID.
 * @returns {Promise<void>}
 */
async function fetchTemplateDetails(templateId: string): Promise<void> {
  if (!apis.assets) {
    error.value = true;
    notifyFail({ message: errors.apiNotInitialized.value });
    return;
  }

  loading.value = true;
  error.value = false;
  template.value = null;

  try {
    template.value = await apis.assets.assetTemplate.getById({ assetTemplateId: templateId });
  } catch (err) {
    logger.error('Error fetching asset template details:', err);
    error.value = true;
    notifyFail({ message: t.drawer.error.value });
  } finally {
    loading.value = false;
  }
}

/**
 * Read a dynamically-keyed field (a script body) off the template. The source
 * is accepted as `unknown` and narrowed here so the index access is typed as
 * `unknown` rather than an implicit any.
 * @param {unknown} source - The template to read from.
 * @param {string} key - The field key.
 * @returns {unknown} The field value, or undefined when absent.
 */
function templateField(source: unknown, key: string): unknown {
  return source && typeof source === 'object'
    ? (source as Record<string, unknown>)[key]
    : undefined;
}

/**
 * Whether a script field is configured on the template.
 * @param {string} scriptKey - Script key to check.
 * @returns {boolean} True when the script exists.
 */
function hasScript(scriptKey: string): boolean {
  return !!templateField(template.value, scriptKey);
}

/**
 * Human summary of how many of the four scripts are configured.
 * @returns {string} Summary text, e.g. "3/4 configured".
 */
function getScriptSummary(): string {
  return `${configuredScriptCount()}/4 ${t.drawer.scripts.configured.value.toLowerCase()}`;
}

/**
 * DetailChip color name reflecting how complete the script configuration is.
 * @returns {'green' | 'orange' | 'grey'} A DetailChip-compatible color name.
 */
function getScriptSummaryColorName(): 'green' | 'orange' | 'grey' {
  const count = configuredScriptCount();
  if (count === 4) return 'green';
  if (count >= 2) return 'orange';
  return 'grey';
}

/**
 * Count how many of the four template scripts are configured.
 * @returns {number} The configured-script count (0-4).
 */
function configuredScriptCount(): number {
  return ['scriptProcessor', 'scriptValidator', 'scriptConversion', 'scriptTest']
    .filter((key) => hasScript(key)).length;
}

/**
 * Format a timestamp for display.
 * @param {string | Date | null | undefined} value - Date value to format.
 * @returns {string} Formatted date, or "-" when absent/invalid.
 */
function formatDate(value: string | Date | null | undefined): string {
  if (!value) return '-';
  try {
    const parsed = typeof value === 'string' ? new Date(value) : value;
    return date.formatDate(parsed, 'MMM DD, YYYY HH:mm');
  } catch {
    return '-';
  }
}

/**
 * Close the drawer.
 */
function close(): void {
  emit('update:modelValue', false);
}

/**
 * Emit the edit request for the current template, unless it is read-only
 * (system or marketplace).
 */
function handleEdit(): void {
  if (isReadonlyTemplate.value) return;
  if (!template.value?.id) return;
  emit('edit', template.value.id);
  close();
}

/**
 * Clone a marketplace template into a new local, editable copy owned by the org,
 * then hand the new template id to the parent so it can open it for editing.
 */
async function handleClone(): Promise<void> {
  if (!template.value?.id || cloning.value) return;

  cloning.value = true;
  try {
    // Prefix the copy's name with the localized "Clone -" so it is easy to spot,
    // set in the single clone request (no follow-up rename call).
    const created = await apis.assets.assetTemplate.clone(
      { assetTemplateId: template.value.id },
      { name: `${t.actions.clonePrefix.value} - ${template.value.name}` },
    );
    notifySuccess({ message: t.drawer.clone.success.value });
    if (created?.id) emit('cloned', created.id);
    close();
  } catch (err) {
    logger.error('Failed to clone asset template', err);
    notifyFail({ message: t.drawer.clone.error.value });
  } finally {
    cloning.value = false;
  }
}

/**
 * Open the script viewer for a template script.
 * @param {string} scriptKey - Script key.
 * @param {string} scriptTitle - Dialog title.
 * @param {'javascript' | 'json'} language - Syntax highlighting language.
 */
function viewScript(scriptKey: string, scriptTitle: string, language: 'javascript' | 'json' = 'javascript'): void {
  const scriptContent = templateField(template.value, scriptKey);
  if (typeof scriptContent !== 'string' || !scriptContent) return;

  currentScriptTitle.value = scriptTitle;
  currentScriptContent.value = scriptContent;
  currentScriptLanguage.value = language;
  showScriptViewer.value = true;
}

/**
 * Open the fields modal showing the template's dynamic fields.
 */
function viewDynamicFields(): void {
  fieldsModalKind.value = 'dynamic';
  showFieldsModal.value = true;
}

/**
 * Open the fields modal showing the template's available fields.
 */
function viewAvailableFields(): void {
  fieldsModalKind.value = 'available';
  showFieldsModal.value = true;
}

/**
 * Close the drawer on the Escape key while it is open.
 * @param {KeyboardEvent} event - Keyboard event.
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

<template>
  <!-- Invisible backdrop for click-outside detection -->
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
          <div class="drawer-hint">{{ t.drawer.loading.value }}</div>
        </div>

        <!-- Error State -->
        <div v-else-if="error" class="q-pa-lg">
          <InfoBanner variant="danger">
            {{ t.drawer.error.value }}
          </InfoBanner>
        </div>

        <!-- Template Data -->
        <div v-else-if="template" class="q-px-md q-py-lg">

          <!-- System Template Warning -->
          <div v-if="template.isSystem" class="q-mb-md">
            <InfoBanner variant="warning">
              {{ t.drawer.systemTemplateWarning.value }}
            </InfoBanner>
          </div>

          <!-- Setup context: identity + classification + asset id path -->
          <div class="section q-mb-md">
            <div class="section-header">
              <q-icon name="mdi-cog-outline" color="primary" size="sm" class="q-mr-sm" />
              <span class="text-subtitle1 text-weight-medium">{{ t.drawer.contexts.setup.value }}</span>
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
            <div class="field-row q-mb-md">
              <div class="field-label">{{ t.drawer.fields.description.value }}</div>
              <div class="field-value">
                {{ template?.description || t.drawer.empty.description.value }}
              </div>
            </div>

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
              <div class="field-value">
                <code class="drawer-code">{{ template?.assetIdPath || '-' }}</code>
              </div>
            </div>
          </div>

          <!-- Uplink context: the ingestion scripts + test payload -->
          <div class="section q-mb-md">
            <div class="section-header">
              <q-icon name="mdi-upload-network-outline" color="primary" size="sm" class="q-mr-sm" />
              <span class="text-subtitle1 text-weight-medium">{{ t.drawer.contexts.uplink.value }}</span>
            </div>
            <q-separator class="q-my-sm" />

            <!-- Preprocessor & Validation (2 columns) -->
            <div class="row q-col-gutter-sm q-mb-md">
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
            </div>

            <!-- Conversion & Test Payload (2 columns) -->
            <div class="row q-col-gutter-sm q-mb-md">
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

          <!-- Retrieval context: the stored/queryable fields -->
          <div class="section q-mb-md">
            <div class="section-header">
              <q-icon name="mdi-database-search-outline" color="primary" size="sm" class="q-mr-sm" />
              <span class="text-subtitle1 text-weight-medium">{{ t.drawer.contexts.retrieval.value }}</span>
            </div>
            <q-separator class="q-my-sm" />

            <!-- Dynamic Fields -->
            <div class="field-row">
              <div class="field-label">{{ t.drawer.fields.dynamicFields.value }}</div>
              <div class="field-value row items-center q-gutter-xs">
                <DetailChip
                  icon="dataset"
                  :color="dynamicFieldsCount > 0 ? 'blue' : 'grey'"
                  size="sm"
                  :label="String(dynamicFieldsCount)"
                />
                <q-btn
                  v-if="dynamicFieldsCount > 0"
                  flat
                  dense
                  round
                  size="sm"
                  icon="visibility"
                  color="primary"
                  @click="viewDynamicFields()"
                >
                  <AppTooltip :content="t.drawer.viewFields.value" />
                </q-btn>
              </div>
            </div>

            <!-- Available Fields -->
            <div class="field-row">
              <div class="field-label">{{ t.drawer.fields.availableFields.value }}</div>
              <div class="field-value row items-center q-gutter-xs">
                <DetailChip
                  icon="list"
                  :color="availableFieldsCount > 0 ? 'blue' : 'grey'"
                  size="sm"
                  :label="String(availableFieldsCount)"
                />
                <q-btn
                  v-if="availableFieldsCount > 0"
                  flat
                  dense
                  round
                  size="sm"
                  icon="visibility"
                  color="primary"
                  @click="viewAvailableFields()"
                >
                  <AppTooltip :content="t.drawer.viewFields.value" />
                </q-btn>
              </div>
            </div>
          </div>

          <!-- Finalization context: audit timestamps -->
          <div class="section">
            <div class="section-header">
              <q-icon name="mdi-clipboard-check" color="primary" size="sm" class="q-mr-sm" />
              <span class="text-subtitle1 text-weight-medium">{{ t.drawer.contexts.finalization.value }}</span>
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

      <!-- Clone a marketplace template into a local, editable copy -->
      <BaseButton
        v-if="isMarketplaceTemplate"
        outline
        no-caps
        icon="content_copy"
        color="primary"
        class="drawer-footer__clone"
        :label="t.drawer.clone.button.value"
        :loading="cloning"
        :disable="!template"
        @click="handleClone"
      />

      <!-- Marketplace templates are read-only, so Edit is hidden (Clone replaces it) -->
      <BaseButton
        v-if="!isMarketplaceTemplate"
        unelevated
        icon="edit"
        color="primary"
        :label="t.drawer.edit.value"
        :disable="!template || isSystemTemplate"
        @click="handleEdit"
      >
        <AppTooltip v-if="isSystemTemplate" :content="t.drawer.systemTemplateTooltip.value" />
      </BaseButton>
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

  <!-- Fields Content Modal -->
  <ContentModal v-model="showFieldsModal" :title="fieldsModalTitle" icon="dataset">
    <DynamicFieldsTable
      v-if="fieldsModalKind === 'dynamic'"
      :fields="template?.dynamicFields ?? []"
    />
    <AvailableFieldsList
      v-else-if="fieldsModalKind === 'available'"
      :fields="template?.availableFields ?? []"
    />
  </ContentModal>
</template>

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

// Secondary hint text (loading state)
.drawer-hint {
  color: var(--mapex-text-secondary);
}

// Drawer Footer - Fixed at bottom
.drawer-footer {
  flex-shrink: 0;
  display: flex;
  align-items: center;
  gap: var(--mapex-spacing-sm);
  padding: 12px 16px;
  background: var(--mapex-header-bg);
  backdrop-filter: blur(10px);
  border-top: 1px solid var(--mapex-header-border);
  box-shadow: 0 -2px 8px var(--mapex-elevation-shadow);
}

// Section (context) Styling
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

// Inline code (asset id path)
.drawer-code {
  display: inline-block;
  background: var(--mapex-surface-sunken);
  padding: var(--mapex-spacing-xs);
  border-radius: var(--mapex-radius-sm);
  font-family: ui-monospace, 'SFMono-Regular', 'Menlo', monospace;
  font-size: 0.85rem;
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
