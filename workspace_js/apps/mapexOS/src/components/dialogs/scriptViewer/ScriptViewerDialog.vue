<script setup lang="ts">
defineOptions({
  name: 'ScriptViewerDialog'
});

/**
 * ScriptViewerDialog - Script Viewer Modal Component
 *
 * NEW MAPEXOS MODAL PATTERN (2025):
 * ================================
 * Header Design:
 * - q-card-section with class "modal-header"
 * - Background: #fafafa (light grey)
 * - Padding: 20px 24px (vertical horizontal)
 * - Border-bottom: 1px solid #e0e0e0
 * - Title: text-h5 (1.25rem), text-grey-9, font-weight-medium
 * - Icon: size="md", color="primary"
 * - Action buttons: flat, round, dense, color="grey-7"
 * - Separator after header
 *
 * This pattern should be used for ALL new modals in MapexOS.
 */
import type { ScriptViewerDialogProps, ScriptViewerDialogEmits } from './interfaces';

import { ref, watch, onBeforeUnmount, nextTick } from 'vue';
import * as monaco from 'monaco-editor';
import { registerMapexMonacoThemes, getMapexMonacoTheme, applyMapexMonacoTheme } from '@utils/monaco-theme';
import { useThemeStore } from '@stores/theme';

import { AppTooltip } from '@components/tooltips';
import { notifyInfo, notifyFail } from '@utils/alert/notify';
import { useLogger } from '@composables/useLogger';

const logger = useLogger('ScriptViewerDialog');
const themeStore = useThemeStore();

const props = withDefaults(defineProps<ScriptViewerDialogProps>(), {
  language: 'javascript',
  copyTooltip: 'Copy script',
  closeTooltip: 'Close',
  copySuccessMessage: 'Script copied to clipboard',
  copyFailMessage: 'Failed to copy script'
});

const emit = defineEmits<ScriptViewerDialogEmits>();

// Watch theme changes to update Monaco editor
watch(() => themeStore.isDark, (isDark: boolean) => {
  applyMapexMonacoTheme(isDark);
});

// Editor container ref
const editorContainerRef = ref<HTMLElement | null>(null);
let editorInstance: monaco.editor.IStandaloneCodeEditor | null = null;

// Setup Monaco Editor when dialog opens
watch(() => props.modelValue, async (isOpen) => {
  if (isOpen) {
    await nextTick();
    // Add a small delay to ensure DOM is fully rendered
    setTimeout(() => {
      setupEditor();
    }, 100);
  } else {
    disposeEditor();
  }
});

// Monaco Editor options for read-only viewing
const editorOptions: monaco.editor.IStandaloneEditorConstructionOptions = {
  language: props.language,
  automaticLayout: true,
  minimap: { enabled: false },
  scrollBeyondLastLine: false,
  wordWrap: 'on',
  fontSize: 13,
  lineNumbers: 'on',
  roundedSelection: false,
  cursorStyle: 'line',
  readOnly: true, // IMPORTANT: Read-only mode
  domReadOnly: true,
  contextmenu: false,
  scrollbar: {
    vertical: 'auto',
    horizontal: 'auto',
    verticalScrollbarSize: 10,
    horizontalScrollbarSize: 10,
  },
};

// Setup Monaco Editor
function setupEditor() {
  try {
    if (!editorContainerRef.value) {
      logger.warn('Editor container ref is null');
      return;
    }

    // Dispose previous instance if exists
    if (editorInstance) {
      editorInstance.dispose();
      editorInstance = null;
    }

    // Create new read-only editor instance
    registerMapexMonacoThemes();
    editorInstance = monaco.editor.create(editorContainerRef.value, {
      ...editorOptions,
      theme: getMapexMonacoTheme(themeStore.isDark),
      value: props.scriptContent || '// No script content',
      language: props.language,
    });
  } catch (error) {
    logger.error('Error setting up Monaco Editor', error);
  }
}

// Dispose Monaco Editor
function disposeEditor() {
  if (editorInstance) {
    editorInstance.dispose();
    editorInstance = null;
  }
}

// Cleanup on unmount
onBeforeUnmount(() => {
  disposeEditor();
});

// Copy script to clipboard
async function copyToClipboard() {
  try {
    await navigator.clipboard.writeText(props.scriptContent);
    notifyInfo({
      message: props.copySuccessMessage,
      timeout: 2000,
    });
  } catch (error) {
    logger.error('Failed to copy to clipboard:', error);
    notifyFail({
      message: props.copyFailMessage,
      timeout: 2000,
    });
  }
}
</script>

<template>
  <q-dialog :model-value="modelValue" @update:model-value="emit('update:modelValue', $event)">
    <q-card class="script-viewer-card" style="min-width: 800px; max-width: 90vw; min-height: 600px; max-height: 90vh;">
      <!-- Header - New MapexOS Modal Pattern -->
      <q-card-section class="modal-header">
        <div class="row items-center">
          <q-icon name="code" size="md" color="primary" class="q-mr-sm" />
          <div class="text-h5 text-weight-medium text-grey-9">{{ title }}</div>
          <q-space />

          <!-- Copy Button -->
          <q-btn
            flat
            round
            dense
            icon="content_copy"
            color="grey-7"
            @click="copyToClipboard"
          >
            <AppTooltip :content="copyTooltip" />
          </q-btn>

          <!-- Close Button -->
          <q-btn
            flat
            round
            dense
            icon="close"
            color="grey-7"
            v-close-popup
          >
            <AppTooltip :content="closeTooltip" />
          </q-btn>
        </div>
      </q-card-section>

      <q-separator />

      <!-- Monaco Editor Container -->
      <q-card-section class="q-pa-none full-height">
        <div
          ref="editorContainerRef"
          class="monaco-editor-container"
        ></div>
      </q-card-section>
    </q-card>
  </q-dialog>
</template>

<style scoped>
.script-viewer-card {
  display: flex;
  flex-direction: column;
}

/* New MapexOS Modal Header Pattern */
.modal-header {
  padding: 20px 24px !important;
  background: var(--mapex-surface-bg);
  border-bottom: 1px solid var(--mapex-card-border);
}

.modal-header .text-h5 {
  font-size: 1.25rem;
  line-height: 1.4;
}

.monaco-editor-container {
  width: 100%;
  height: 600px;
}

.full-height {
  flex: 1;
  overflow: hidden;
}
</style>
