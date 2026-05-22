<script setup lang="ts">
/** TYPE IMPORTS */
import type { WorkflowDefinition } from '../../interfaces/CreateEditWorkflow.interface';
import type { editor as MonacoEditorNs } from 'monaco-editor';

defineOptions({
  name: 'JsonDebugTab',
});

/** VUE IMPORTS */
import { ref, watch, nextTick, onBeforeUnmount, computed } from 'vue';

/** COMPONENTS */
import { AppTooltip } from '@components/tooltips';

/** COMPOSABLES */
import { useWorkflowEditorState } from '../../composables';
import { useCreateEditWorkflowTranslations } from '@composables/i18n/pages/automations/workflows';

/** UTILS */
import { monaco as monacoEditor } from '@utils/monaco-setup';
import { registerMapexMonacoThemes, getMapexMonacoTheme, applyMapexMonacoTheme } from '@utils/monaco-theme';
import { notifySuccess, notifyWarning, notifyFail } from '@utils/alert/notify';

/** STORES */
import { useThemeStore } from '@stores/theme';

/** COMPOSABLES & STORES */
const { getCurrentWorkflow, setAllStates } = useWorkflowEditorState();
const themeStore = useThemeStore();
const t = useCreateEditWorkflowTranslations();

/**
 * Get formatted JSON from current workflow state
 */
const getFormattedJson = computed(() => {
  return JSON.stringify(getCurrentWorkflow.value, null, 2);
});

/** STATE */
const mode = ref<'view' | 'edit'>('view');
const monacoContainer = ref<HTMLElement | null>(null);
let monacoEditorInstance: MonacoEditorNs.IStandaloneCodeEditor | null = null;
const validationError = ref<string>('');

/** COMPUTED */

/**
 * Check if editor is read-only based on mode
 */
const isReadOnly = computed(() => mode.value === 'view');

/** WATCHERS */

/**
 * Watch for mode changes and reinitialize editor
 */
watch(mode, () => {
  void initMonacoEditor();
});

/**
 * Watch for external workflow changes in view mode
 */
watch(getFormattedJson, (newJson: string) => {
  if (mode.value === 'view' && monacoEditorInstance) {
    monacoEditorInstance.setValue(newJson);
    void monacoEditorInstance.getAction('editor.action.formatDocument')?.run();
  }
});

/**
 * Initialize on mount
 */
watch(monacoContainer, (container: HTMLElement | null) => {
  if (container) {
    void initMonacoEditor();
  }
}, { immediate: true });

/**
 * Watch theme changes to update Monaco editor
 */
watch(() => themeStore.isDark, (isDark: boolean) => {
  applyMapexMonacoTheme(isDark);
});

/** FUNCTIONS */

/**
 * Initialize Monaco Editor instance
 *
 * @returns {Promise<void>}
 */
async function initMonacoEditor(): Promise<void> {
  await nextTick();

  if (monacoEditorInstance) {
    monacoEditorInstance.dispose();
    monacoEditorInstance = null;
  }

  if (!monacoContainer.value) return;

  registerMapexMonacoThemes();

  monacoEditorInstance = monacoEditor.editor.create(monacoContainer.value, {
    value: getFormattedJson.value,
    language: 'json',
    theme: getMapexMonacoTheme(themeStore.isDark),
    readOnly: isReadOnly.value,
    automaticLayout: true,
    minimap: { enabled: false },
    fontSize: 14,
    lineNumbers: 'on',
    wordWrap: 'on',
    scrollBeyondLastLine: false,
    formatOnPaste: true,
    formatOnType: true,
  });

  void monacoEditorInstance.getAction('editor.action.formatDocument')?.run();

  if (!isReadOnly.value) {
    monacoEditorInstance.onDidChangeModelContent(() => {
      validateJson();
    });
  }
}

/**
 * Validate JSON in the editor
 *
 * @returns {void}
 */
function validateJson(): void {
  if (!monacoEditorInstance) return;

  const text = monacoEditorInstance.getValue();
  try {
    JSON.parse(text);
    validationError.value = '';
  } catch (error: any) {
    validationError.value = error.message || 'Invalid JSON';
  }
}

/**
 * Apply JSON changes to the workflow state.
 * Parses JSON, validates structure, and syncs all tabs via setAllStates.
 *
 * @returns {void}
 */
function applyJsonChanges(): void {
  if (!monacoEditorInstance) return;

  const text = monacoEditorInstance.getValue();
  try {
    const parsed = JSON.parse(text) as WorkflowDefinition;

    if (!parsed.name || typeof parsed.name !== 'string') {
      notifyWarning({ message: t.notifications.invalidWorkflowName.value });
      return;
    }

    if (!Array.isArray(parsed.nodes)) {
      notifyWarning({ message: t.notifications.invalidWorkflowNodes.value });
      return;
    }

    if (!Array.isArray(parsed.edges)) {
      notifyWarning({ message: t.notifications.invalidWorkflowEdges.value });
      return;
    }

    setAllStates(parsed);

    mode.value = 'view';
    validationError.value = '';

    notifySuccess({ message: t.notifications.workflowUpdatedFromJson.value });
  } catch (error: any) {
    notifyWarning({ message: `Invalid JSON: ${error.message}` });
  }
}

/**
 * Cancel edit mode and revert editor to current workflow state
 *
 * @returns {void}
 */
function cancelEdit(): void {
  mode.value = 'view';
  validationError.value = '';
  void initMonacoEditor();
}

/**
 * Copy JSON to clipboard
 *
 * @returns {Promise<void>}
 */
async function copyJson(): Promise<void> {
  const text = monacoEditorInstance?.getValue() ?? getFormattedJson.value;
  try {
    await navigator.clipboard.writeText(text);
    notifySuccess({ message: t.notifications.jsonCopied.value });
  } catch {
    notifyFail({ message: t.notifications.copyFailed.value });
  }
}

/** LIFECYCLE HOOKS */
onBeforeUnmount(() => {
  if (monacoEditorInstance) {
    monacoEditorInstance.dispose();
    monacoEditorInstance = null;
  }
});
</script>

<template>
  <div class="json-debug-tab q-pa-md">
    <!-- Header -->
    <div class="page-header q-mb-lg">
      <div class="row items-center justify-between">
        <div class="col-auto">
          <div class="text-h5 text-weight-medium q-mb-xs" style="color: var(--mapex-text-primary)">
            {{ t.jsonDebug.title.value }}
          </div>
          <div class="text-body2" style="color: var(--mapex-text-secondary)">
            {{ mode === 'view' ? t.jsonDebug.descriptionView.value : t.jsonDebug.descriptionEdit.value }}
          </div>
        </div>

        <div class="col-auto">
          <div class="row items-center q-gutter-sm">
            <!-- Mode Toggle -->
            <q-btn-toggle
              v-model="mode"
              toggle-color="primary"
              :options="[
                { label: t.jsonDebug.viewMode.value, value: 'view', icon: 'visibility' },
                { label: t.jsonDebug.editMode.value, value: 'edit', icon: 'edit' }
              ]"
              unelevated
              no-caps
              class="mode-toggle"
            />

            <!-- Copy Button -->
            <q-btn
              flat
              round
              dense
              icon="content_copy"
              color="primary"
              @click="copyJson"
            >
              <AppTooltip :content="t.jsonDebug.copyTooltip.value" />
            </q-btn>
          </div>
        </div>
      </div>
    </div>

    <!-- Monaco Editor -->
    <q-card flat bordered class="monaco-card">
      <div
        ref="monacoContainer"
        class="monaco-container"
      />
    </q-card>

    <!-- Validation Error -->
    <div v-if="validationError && mode === 'edit'" class="q-mt-sm text-negative text-caption">
      <q-icon name="warning" size="sm" class="q-mr-xs" />
      {{ validationError }}
    </div>

    <!-- Edit Mode Actions -->
    <div v-if="mode === 'edit'" class="q-mt-md row justify-end q-gutter-sm">
      <q-btn
        flat
        no-caps
        :label="t.jsonDebug.cancel.value"
        color="grey-7"
        @click="cancelEdit"
      />
      <q-btn
        unelevated
        no-caps
        :label="t.jsonDebug.applyChanges.value"
        color="primary"
        icon="check"
        :disable="!!validationError"
        @click="applyJsonChanges"
      />
    </div>

    <!-- Guidelines -->
    <div class="q-mt-md">
      <q-expansion-item
        class="text-primary"
        icon="help"
        :label="t.jsonDebug.guidelines.value"
      >
        <q-card>
          <q-card-section class="text-body2" style="color: var(--mapex-text-secondary)">
            <p><strong>{{ t.jsonDebug.expectedStructure.value }}</strong></p>
            <ul>
              <li><code>name</code> - {{ t.jsonDebug.structureName.value }}</li>
              <li><code>description</code> - {{ t.jsonDebug.structureDescription.value }}</li>
              <li><code>variables</code> - {{ t.jsonDebug.structureVariables.value }}</li>
              <li><code>captureFields</code> - {{ t.jsonDebug.structureCaptureFields.value }}</li>
              <li><code>externalVariables</code> - {{ t.jsonDebug.structureExternalVariables.value }}</li>
              <li><code>nodes</code> - {{ t.jsonDebug.structureNodes.value }}</li>
              <li><code>edges</code> - {{ t.jsonDebug.structureEdges.value }}</li>
            </ul>

            <p class="q-mt-md"><strong>{{ t.jsonDebug.editModeTitle.value }}</strong></p>
            <ul>
              <li>{{ t.jsonDebug.editModeHint1.value }}</li>
              <li>{{ t.jsonDebug.editModeHint2.value }}</li>
              <li>{{ t.jsonDebug.editModeHint3.value }}</li>
              <li>{{ t.jsonDebug.editModeHint4.value }}</li>
            </ul>

            <p class="q-mt-md"><strong>{{ t.jsonDebug.tipsTitle.value }}</strong></p>
            <ul>
              <li><kbd>Ctrl+F</kbd> - {{ t.jsonDebug.tipSearch.value }}</li>
              <li><kbd>Ctrl+A</kbd> - {{ t.jsonDebug.tipSelectAll.value }}</li>
              <li>{{ t.jsonDebug.tipFormat.value }}</li>
            </ul>
          </q-card-section>
        </q-card>
      </q-expansion-item>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.json-debug-tab {
  max-width: 1200px;
  text-align: left;

  .page-header {
    .text-h5 {
      font-size: 1.25rem;
      font-weight: 500;
      line-height: 1.4;
    }

    .text-body2 {
      font-size: 0.875rem;
      line-height: 1.5;
    }

    .mode-toggle {
      border-radius: var(--mapex-radius-md);
      overflow: hidden;

      :deep(.q-btn) {
        padding: 8px 16px;
      }
    }
  }
}

.monaco-card {
  border-radius: var(--mapex-radius-md) !important;
  border-color: var(--mapex-card-border) !important;
  background: var(--mapex-surface-bg) !important;
  overflow: hidden;

  .monaco-container {
    width: 100%;
    height: 600px;
    min-height: 400px;
  }
}

@media (max-width: 768px) {
  .monaco-card .monaco-container {
    height: 400px;
  }
}
</style>
