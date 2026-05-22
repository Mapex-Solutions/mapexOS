<template>
  <q-drawer
      v-model="localShow"
      overlay
      side="right"
      behavior="mobile"
      class="json-drawer"
      :width="$q.screen.lt.md ? $q.screen.width : 700"
  >
    <div class="column full-height">
      <!-- HEADER -->
      <div class="row items-center q-pa-md drawer-header">
        <div class="col">
          <div class="text-h6 text-primary">{{ title }}</div>
          <div v-if="subtitle" class="text-caption text-grey-8">{{ subtitle }}</div>
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
          <q-btn flat round dense icon="close" color="red-4" @click="localShow = false">
            <AppTooltip content="Close" />
          </q-btn>
        </div>
      </div>

      <!-- TABS -->
      <AppTabs v-model="tab" :tabs="drawerTabs" @change="(value: string) => onTabClick(value as 'json' | 'event')" />

      <!-- CONTENT -->
      <div class="col q-pa-md monaco-body" v-if="tab === 'json'">
        <div ref="monacoJsonContainer" class="monaco-container"/>
      </div>

      <div v-else class="col q-pa-md monaco-body">
        <!-- Monaco Editor for Event Data -->
        <div v-if="!isEmpty(props.eventData)">
          <div ref="monacoEventContainer" class="monaco-container"/>
        </div>

        <!-- Skeleton Loader -->
        <div v-else class="text-grey-6">
          <q-skeleton type="text" width="40%" class="q-mb-md"/>
          <q-skeleton type="text" width="60%" class="q-mb-md"/>
          <q-skeleton type="text" width="50%" class="q-mb-md"/>
          <q-skeleton type="text" width="70%" class="q-mb-md"/>
        </div>
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

<script setup lang="ts">
defineOptions({
  name: 'JsonAndEventDrawer'
});

import type { JsonAndEventDrawerProps, JsonAndEventDrawerEmit } from './interfaces';

import { isEmpty } from 'lodash';
import { ref, watch, nextTick, onBeforeUnmount } from 'vue';
import { useQuasar } from 'quasar';
import * as monaco from 'monaco-editor';
import { registerMapexMonacoThemes, getMapexMonacoTheme, applyMapexMonacoTheme } from '@utils/monaco-theme';
import { notifyInfo, notifySuccess, notifyFail, notifyWarning } from '@utils/alert/notify';
import { AppTooltip } from '@components/tooltips';
import { AppTabs } from '@components/tabs';
import { useThemeStore } from '@stores/theme';

const props = defineProps<JsonAndEventDrawerProps>();
const emit = defineEmits<JsonAndEventDrawerEmit>();

const $q = useQuasar();
const themeStore = useThemeStore();
const localShow = ref<boolean>(props.show);
const tab = ref<'json' | 'event'>('json');
const loadingEvent = ref(false);

/**
 * Tabs configuration
 */
const drawerTabs = [
  { name: 'json', label: 'Raw JSON', icon: 'code' },
  { name: 'event', label: 'Event Data', icon: 'event' },
];

// Monaco Editors & Containers
let monacoJsonEditor: monaco.editor.IStandaloneCodeEditor | null = null;
let monacoEventEditor: monaco.editor.IStandaloneCodeEditor | null = null;
const monacoJsonContainer = ref<HTMLElement | null>(null);
const monacoEventContainer = ref<HTMLElement | null>(null);

watch(() => props.show, val => (localShow.value = val));
watch(localShow, val => emit('update:show', val));

// Watch theme changes to update Monaco editor
watch(() => themeStore.isDark, (isDark: boolean) => {
  applyMapexMonacoTheme(isDark);
});

// Initialize either JSON or Event editor
watch([() => tab.value, localShow, () => props.jsonData, () => props.eventData], async () => {
  if (!localShow.value) return;

  await nextTick();

  registerMapexMonacoThemes();

  // JSON Tab
  if (tab.value === 'json') {
    monacoEventEditor?.dispose();
    monacoEventEditor = null;
    if (!monacoJsonContainer.value) return;
    monacoJsonEditor?.dispose();
    monacoJsonEditor = monaco.editor.create(monacoJsonContainer.value, {
      value: typeof props.jsonData === 'string' ? props.jsonData : JSON.stringify(props.jsonData, null, 2),
      language: 'json', theme: getMapexMonacoTheme(themeStore.isDark),
      readOnly: !props.editable, automaticLayout: true,
      minimap: { enabled: false }, fontSize: 14,
      lineNumbers: 'on', wordWrap: 'on',
    });
    void monacoJsonEditor.getAction('editor.action.formatDocument')?.run();
  }

  // Event Tab
  if (tab.value === 'event' && props.eventData) {
    monacoJsonEditor?.dispose();
    monacoJsonEditor = null;
    if (!monacoEventContainer.value) return;
    monacoEventEditor?.dispose();
    monacoEventEditor = monaco.editor.create(monacoEventContainer.value, {
      value: JSON.stringify(props.eventData, null, 2),
      language: 'json', theme: getMapexMonacoTheme(themeStore.isDark),
      readOnly: true, automaticLayout: true,
      minimap: { enabled: false }, fontSize: 14,
      lineNumbers: 'on', wordWrap: 'on',
    });
    void monacoEventEditor.getAction('editor.action.formatDocument')?.run();
  }
}, { immediate: true });

async function copyJson() {
  const text = monacoJsonEditor?.getValue() ?? typeof props.jsonData === 'string' ? props.jsonData : JSON.stringify(props.jsonData, null, 2);
  try {
    await navigator.clipboard.writeText(JSON.stringify(text));
    notifyInfo({ message: 'JSON copied' });
  } catch {
    notifyFail({ message: 'Copy failed' });
  }
}

function onSave() {
  if (!monacoJsonEditor) return;
  const text = monacoJsonEditor.getValue();
  try {
    const parsed = JSON.parse(text);
    emit('save', parsed);
    notifySuccess({ message: 'Saved successfully' });
  } catch {
    notifyWarning({ message: 'Invalid JSON' });
  }
}

function onTabClick(name: 'json' | 'event') {
  tab.value = name;
  if (name === 'event') {
    loadingEvent.value = true;
    const id = (props.jsonData as any).eventId;
    if (id) emit('fetch-event', id);
  }
}

watch(localShow, open => {
  if (!open) {
    monacoJsonEditor?.dispose();
    monacoEventEditor?.dispose();
    monacoJsonEditor = null;
    monacoEventEditor = null;
  }
});

onBeforeUnmount(() => {
  monacoJsonEditor?.dispose();
  monacoEventEditor?.dispose();
});
</script>

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
