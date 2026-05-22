<script setup lang="ts">
defineOptions({
  name: 'JsonDrawer'
});

import type { JsonDrawerProps, JsonDrawerEmit } from './interfaces';

import { ref, watch, nextTick, onBeforeUnmount } from 'vue';
import { useQuasar } from 'quasar';

import * as monaco from 'monaco-editor';
import { registerMapexMonacoThemes, getMapexMonacoTheme, applyMapexMonacoTheme } from '@utils/monaco-theme';
import { notifyInfo, notifySuccess, notifyFail, notifyWarning } from '@utils/alert/notify';
import { AppTooltip } from '@components/tooltips';

import { useThemeStore } from '@stores/theme';

// --- Define Props & Emits ---
const props = defineProps<JsonDrawerProps>();
const emit = defineEmits<JsonDrawerEmit>();

const $q = useQuasar();
const themeStore = useThemeStore();

// Local copy of show for v-model
const localShow = ref<boolean>(props.show);

watch(() => props.show, val => (localShow.value = val));
watch(localShow, val => emit('update:show', val));

// Watch theme changes to update Monaco editor
watch(() => themeStore.isDark, (isDark: boolean) => {
  applyMapexMonacoTheme(isDark);
});

// Monaco editor ref
const monacoContainer = ref<HTMLElement | null>(null);
let monacoEditor: monaco.editor.IStandaloneCodeEditor | null = null;

// Return a pretty-printed JSON string
function getJsonString(): string {
  return typeof props.jsonData === 'string'
      ? props.jsonData
      : JSON.stringify(props.jsonData, null, 2);
}

// Initialize or re-initialize Monaco when needed
watch(
    [localShow, () => props.jsonData, () => props.editable],
    async ([isOpen]) => {
      if (!isOpen) return;

      await nextTick();
      if (monacoEditor) {
        monacoEditor.dispose();
        monacoEditor = null;
      }
      if (!monacoContainer.value) return;

      registerMapexMonacoThemes();

      monacoEditor = monaco.editor.create(monacoContainer.value, {
        value: getJsonString(),
        language: 'json',
        theme: getMapexMonacoTheme(themeStore.isDark),
        readOnly: !props.editable,
        automaticLayout: true,
        minimap: { enabled: false },
        fontSize: 14,
        lineNumbers: 'on',
        wordWrap: 'on',
      });

      // auto-format
      void monacoEditor.getAction('editor.action.formatDocument')?.run();
    },
    { immediate: true },
);

// Copy current JSON text to clipboard
async function copyJson() {
  const text = monacoEditor?.getValue() ?? getJsonString();
  try {
    await navigator.clipboard.writeText(text);
    notifyInfo({ message: 'JSON copied' });
  } catch {
    notifyFail({ message: 'Copy failed' });
  }
}

// Emit save with parsed JSON if valid
function onSave() {
  if (!monacoEditor) return;
  const text = monacoEditor.getValue();
  try {
    const parsed = JSON.parse(text);
    emit('save', parsed);
    notifySuccess({ message: 'Saved successfully' });
  } catch {
    notifyWarning({ message: 'Invalid JSON' });
  }
}

// Cleanup editor on close/unmount
watch(localShow, open => {
  if (!open && monacoEditor) {
    monacoEditor.dispose();
    monacoEditor = null;
  }
});
onBeforeUnmount(() => {
  if (monacoEditor) {
    monacoEditor.dispose();
    monacoEditor = null;
  }
});
</script>

<template>
  <q-drawer
      v-model="localShow"
      overlay
      side="right"
      behavior="mobile"
      class="json-drawer"
      :width="$q.screen.lt.md ? $q.screen.width : 800"
  >
    <div class="column full-height">

      <!-- HEADER -->
      <div class="row items-center q-pa-md drawer-header">
        <div class="col">
          <div class="text-h6 text-primary">{{ title }}</div>
          <div v-if="subtitle" class="text-caption text-grey-8">
            {{ subtitle }}
          </div>
        </div>
        <div class="col-auto">
          <q-btn
              v-if="editable"
              flat
              round
              dense
              icon="save"
              class="q-mr-sm"
              color="primary"
              @click="onSave"
          >
            <AppTooltip content="Save JSON" />
          </q-btn>
          <q-btn
              flat
              round
              dense
              icon="content_copy"
              class="q-mr-sm"
              color="primary"
              @click="copyJson"
          >
            <AppTooltip content="Copy JSON" />
          </q-btn>
          <q-btn
              flat
              round
              dense
              icon="close"
              color="red-4"
              @click="localShow = false"
          >
            <AppTooltip content="Close" />
          </q-btn>
        </div>
      </div>

      <!-- MONACO EDITOR -->
      <div class="col q-pa-none monaco-body">
        <div ref="monacoContainer" class="monaco-container"/>
      </div>

      <!-- FOOTER NOTE -->
      <div class="row items-center q-pa-md drawer-footer">
        <div class="col">
          <div class="text-caption text-grey-8">
            Use Ctrl+F to search • Ctrl+A to select all
            <span v-if="editable"> • Editing enabled</span>
          </div>
        </div>
      </div>

    </div>
  </q-drawer>
</template>

<style scoped lang="scss">
.json-drawer {
  .drawer-header {
    background: var(--mapex-header-bg);
    border-bottom: 1px solid var(--mapex-header-border);
  }

  .drawer-footer {
    background: var(--mapex-header-bg);
    border-top: 1px solid var(--mapex-header-border);
  }

  .monaco-body {
    background-color: var(--mapex-surface-bg);
  }

  .monaco-container {
    width: 100%;
    height: 100%;
    min-height: 400px;
  }
}
</style>
