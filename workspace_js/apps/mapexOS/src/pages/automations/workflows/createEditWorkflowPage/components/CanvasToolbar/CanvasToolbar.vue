<script setup lang="ts">
/** TYPE IMPORTS */
import type { CanvasToolbarState } from '../../interfaces/CreateEditWorkflow.interface';

/** VUE IMPORTS */
import { ref, computed } from 'vue';

/** COMPONENTS */
import { InfoModal } from '@components/dialogs/infoModal';
import { AppTooltip } from '@components/tooltips';

/** COMPOSABLES */
import { useCreateEditWorkflowTranslations } from '@composables/i18n/pages/automations/workflows';
import { useWorkflowHistory } from '../../composables';

/** LOCAL IMPORTS */
import { HOTKEY_DEFINITIONS } from './constants';

/** PROPS & EMITS */
const model = defineModel<CanvasToolbarState>({ required: true });
const emit = defineEmits<{
  (e: 'auto-organize'): void;
  (e: 'undo'): void;
  (e: 'redo'): void;
}>();

/** COMPOSABLES & STORES */
const { canUndo, canRedo } = useWorkflowHistory();
const t = useCreateEditWorkflowTranslations();

/** STATE */
const showHotkeysModal = ref(false);

/** COMPUTED */

/**
 * Hotkey items built from i18n translations
 */
const hotkeyItems = computed(() =>
  HOTKEY_DEFINITIONS.map(def => ({
    icon: def.icon,
    color: def.color,
    title: def.title,
    text: t.canvasToolbar[def.i18nKey].value,
  })),
);

/** FUNCTIONS */

/**
 * Toggle node movement lock
 *
 * @returns {void}
 */
function toggleLock(): void {
  model.value = { ...model.value, locked: !model.value.locked };
}

/**
 * Toggle maximized (fullscreen) state
 *
 * @returns {void}
 */
function toggleMaximize(): void {
  model.value = { ...model.value, maximized: !model.value.maximized };
}
</script>

<template>
  <q-toolbar class="canvas-toolbar">
    <!-- Auto-organize -->
    <q-btn flat dense icon="account_tree" size="sm" @click="emit('auto-organize')">
      <AppTooltip :content="t.canvasToolbar.autoOrganize.value" />
    </q-btn>

    <q-separator vertical class="q-mx-xs" />

    <!-- Lock -->
    <q-btn
      flat
      dense
      :icon="model.locked ? 'lock' : 'lock_open'"
      size="sm"
      :color="model.locked ? 'amber-8' : undefined"
      @click="toggleLock"
    >
      <AppTooltip :content="model.locked ? t.canvasToolbar.unlockCanvas.value : t.canvasToolbar.lockCanvas.value" />
    </q-btn>

    <q-separator vertical class="q-mx-xs" />

    <!-- Undo / Redo -->
    <q-btn
      flat
      dense
      icon="undo"
      size="sm"
      :disable="!canUndo"
      @click="emit('undo')"
    >
      <AppTooltip :content="t.canvasToolbar.undo.value" />
    </q-btn>
    <q-btn
      flat
      dense
      icon="redo"
      size="sm"
      :disable="!canRedo"
      @click="emit('redo')"
    >
      <AppTooltip :content="t.canvasToolbar.redo.value" />
    </q-btn>

    <q-space />

    <!-- Canvas toggles -->
    <q-toggle v-model="model.showMinimap" :label="t.canvasToolbar.minimap.value" dense size="sm" />
    <q-toggle v-model="model.showGrid" :label="t.canvasToolbar.grid.value" dense size="sm" class="q-ml-sm" />

    <q-separator vertical class="q-mx-xs" />

    <!-- Hotkeys info -->
    <q-btn flat dense icon="keyboard" size="sm" @click="showHotkeysModal = true">
      <AppTooltip :content="t.canvasToolbar.keyboardShortcuts.value" />
    </q-btn>

    <!-- Maximize -->
    <q-btn
      flat
      dense
      :icon="model.maximized ? 'fullscreen_exit' : 'fullscreen'"
      size="sm"
      :color="model.maximized ? 'primary' : undefined"
      @click="toggleMaximize"
    >
      <AppTooltip :content="model.maximized ? t.canvasToolbar.exitFullscreen.value : t.canvasToolbar.fullscreen.value" />
    </q-btn>

    <!-- Hotkeys modal -->
    <InfoModal
      v-model="showHotkeysModal"
      icon="keyboard"
      :title="t.canvasToolbar.shortcutsTitle.value"
      :description="t.canvasToolbar.shortcutsDescription.value"
      :items="hotkeyItems"
      :close-label="t.canvasToolbar.close.value"
    />
  </q-toolbar>
</template>

<style lang="scss" scoped>
.canvas-toolbar {
  min-height: 40px;
  padding: 4px 8px;
  background: var(--mapex-surface-bg);
  border-bottom: 1px solid var(--mapex-card-border);
}
</style>
