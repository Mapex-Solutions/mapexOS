<script setup lang="ts">
defineOptions({
  name: 'ScriptEditorDialog',
});

/**
 * ScriptEditorDialog - Editable Script Modal Component
 *
 * Follows the ScriptViewerDialog pattern but with editable Monaco editor.
 * Uses the NEW MAPEXOS MODAL PATTERN (2025) for header design.
 */

/** TYPE IMPORTS */
import type { ScriptEditorDialogProps, ScriptEditorDialogEmits } from './interfaces';

/** VUE IMPORTS */
import { ref, computed, watch, onBeforeUnmount, nextTick } from 'vue';

/** VUE IMPORTS */
import * as monaco from 'monaco-editor';

/** COMPONENTS */
import { AppTooltip } from '@components/tooltips';

/** COMPOSABLES */
import { useTS } from '@utils/translation';

/** UTILS */
import { registerMapexMonacoThemes, getMapexMonacoTheme, applyMapexMonacoTheme } from '@utils/monaco-theme';

/** STORES */
import { useThemeStore } from '@stores/theme';

/** COMPOSABLES & STORES */
const tsTitle = useTS({ titleCase: true });

/** PROPS & EMITS */
const props = withDefaults(defineProps<ScriptEditorDialogProps>(), {
  language: 'javascript',
  closeTooltip: '',
  guidelines: () => [],
});

/** COMPUTED */
const resolvedCloseTooltip = computed(() =>
  props.closeTooltip || tsTitle('components.dialogs.scriptEditor.closeTooltip'),
);

const emit = defineEmits<ScriptEditorDialogEmits>();

/** COMPOSABLES & STORES */
const themeStore = useThemeStore();

/** STATE */

/**
 * DOM reference for Monaco editor container
 */
const editorContainerRef = ref<HTMLElement | null>(null);

/**
 * Monaco editor instance
 */
let editorInstance: monaco.editor.IStandaloneCodeEditor | null = null;

/** WATCHERS */

/**
 * Watch theme changes to update Monaco editor
 */
watch(() => themeStore.isDark, (isDark: boolean) => {
  applyMapexMonacoTheme(isDark);
});

/**
 * Setup editor when dialog opens, dispose when it closes
 */
watch(() => props.modelValue, async (isOpen) => {
  if (isOpen) {
    await nextTick();
    setTimeout(() => setupEditor(), 100);
  } else {
    disposeEditor();
  }
});

/** FUNCTIONS */

/**
 * Setup Monaco editor inside the dialog
 */
function setupEditor(): void {
  if (!editorContainerRef.value) return;

  if (editorInstance) {
    editorInstance.dispose();
    editorInstance = null;
  }

  registerMapexMonacoThemes();
  editorInstance = monaco.editor.create(editorContainerRef.value, {
    value: props.scriptContent || '',
    language: props.language,
    theme: getMapexMonacoTheme(themeStore.isDark),
    automaticLayout: true,
    minimap: { enabled: false },
    scrollBeyondLastLine: false,
    wordWrap: 'on',
    fontSize: 13,
    lineNumbers: 'on',
    roundedSelection: false,
    cursorStyle: 'line',
    formatOnPaste: true,
    formatOnType: true,
    scrollbar: {
      vertical: 'auto',
      horizontal: 'auto',
      verticalScrollbarSize: 10,
      horizontalScrollbarSize: 10,
    },
  });

  editorInstance.onDidChangeModelContent(() => {
    const value = editorInstance?.getValue() ?? '';
    emit('update:scriptContent', value);
  });
}

/**
 * Dispose Monaco editor instance
 */
function disposeEditor(): void {
  if (editorInstance) {
    editorInstance.dispose();
    editorInstance = null;
  }
}

/** LIFECYCLE HOOKS */

onBeforeUnmount(() => {
  disposeEditor();
});
</script>

<template>
  <q-dialog :model-value="modelValue" @update:model-value="emit('update:modelValue', $event)">
    <q-card class="script-editor-card" style="min-width: 800px; max-width: 90vw; min-height: 600px; max-height: 90vh;">
      <!-- Header - New MapexOS Modal Pattern -->
      <q-card-section class="modal-header">
        <div class="row items-center">
          <q-icon name="code" size="md" color="primary" class="q-mr-sm" />
          <div class="text-h5 text-weight-medium" style="color: var(--mapex-text-primary);">
            {{ title }}
          </div>
          <q-space />

          <!-- Close Button -->
          <q-btn
            flat
            round
            dense
            icon="close"
            color="grey-7"
            v-close-popup
          >
            <AppTooltip :content="resolvedCloseTooltip" />
          </q-btn>
        </div>
      </q-card-section>

      <q-separator />

      <!-- Guidelines bar (optional) -->
      <template v-if="guidelines.length > 0">
        <q-card-section class="script-editor-card__guidelines q-py-sm">
          <div class="row items-center q-gutter-md text-caption" style="color: var(--mapex-text-secondary);">
            <span v-for="g in guidelines" :key="g.code">
              <code>{{ g.code }}</code> — {{ g.description }}
            </span>
            <q-space />
            <slot name="guidelines-right" />
          </div>
        </q-card-section>

        <q-separator />
      </template>

      <!-- Monaco Editor Container -->
      <q-card-section class="q-pa-none script-editor-card__editor-section">
        <div
          ref="editorContainerRef"
          class="script-editor-card__editor"
        />
      </q-card-section>
    </q-card>
  </q-dialog>
</template>

<style lang="scss" scoped>
.script-editor-card {
  display: flex;
  flex-direction: column;

  &__guidelines {
    background: var(--mapex-surface-elevated);

    code {
      font-size: 0.7rem;
      padding: 1px 4px;
      border-radius: var(--mapex-radius-xs);
      background: var(--mapex-surface-bg);
    }
  }

  &__editor-section {
    flex: 1;
    overflow: hidden;
  }

  &__editor {
    width: 100%;
    height: 600px;
  }
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
</style>
