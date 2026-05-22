<script setup lang="ts">
/** TYPE IMPORTS */
import type { Step6TestPayloadProps, Step6TestPayloadEmits } from './interfaces/Step6TestPayload.interface';

defineOptions({
  name: 'Step6TestPayload'
});

import { ref, computed, onMounted, onBeforeUnmount, watch } from 'vue';
import { AppTooltip } from '@components/tooltips';
import { useAddAssetTemplateTranslations } from '@src/composables/i18n/pages/assets/addAssetTemplate/useAddAssetTemplateTranslations';
import { createMonacoEditorManager } from '../../handlers';
import type { MonacoEditorManager } from '../../handlers';
import { DEFAULT_ASSET_TEMPLATE_DATA } from '../../constants';
import { applyMapexMonacoTheme } from '@utils/monaco-theme';
import { useThemeStore } from '@stores/theme';

// Props
// Emits
const props = defineProps<Step6TestPayloadProps>();
const emit = defineEmits<Step6TestPayloadEmits>();

// i18n
const t = useAddAssetTemplateTranslations();
const themeStore = useThemeStore();

// Watch theme changes to update Monaco editor
watch(() => themeStore.isDark, (isDark: boolean) => {
  applyMapexMonacoTheme(isDark);
});

// Editor container ref
const scriptTestEditorRef = ref<HTMLElement | null>(null);

// Editor manager
let scriptTestEditorManager: MonacoEditorManager | null = null;

// Computed v-model for scriptTest (MANDATORY field from backend)
const scriptTest = computed({
  get: () => props.modelValue.scriptTest || '',
  set: (value: string) => {
    const updated = { ...props.modelValue, scriptTest: value };
    emit('update:modelValue', updated);
  }
});

// Load example script test
function loadExampleScriptTest() {
  if (scriptTestEditorManager) {
    scriptTestEditorManager.setValue(DEFAULT_ASSET_TEMPLATE_DATA.scriptTest || '');
  }
}

// Setup editor on mount
onMounted(() => {
  // Setup scriptTest editor (JSON - test payload)
  scriptTestEditorManager = createMonacoEditorManager(
    scriptTestEditorRef,
    { language: 'json' },
    scriptTest.value,
    (value: string) => {
      scriptTest.value = value;
    }
  );
  scriptTestEditorManager.setup();
});

// Cleanup on unmount
onBeforeUnmount(() => {
  scriptTestEditorManager?.dispose();
  scriptTestEditorManager = null;
});

// Watch for external changes
watch(() => props.modelValue.scriptTest, (newValue) => {
  if (scriptTestEditorManager && scriptTestEditorManager.getValue() !== newValue) {
    scriptTestEditorManager.setValue(newValue || '');
  }
});
</script>

<template>
  <div>
    <div class="q-mb-md">
      <div class="text-subtitle1 text-weight-medium q-mb-xs">
        <q-icon name="science" color="primary" class="q-mr-xs" />
        {{ t.steps.step6.title.value }}
      </div>
      <div class="text-body2 text-grey-7">
        {{ t.steps.step6.subtitle.value }}
      </div>
    </div>

    <!-- Test Payload Editor (MANDATORY) -->
    <div class="row items-center q-mb-sm">
      <div class="text-subtitle2 text-weight-medium">
        <q-icon name="code" color="secondary" class="q-mr-xs" />
        Test Payload
      </div>
      <q-space />
      <q-btn
        flat
        dense
        size="sm"
        color="secondary"
        icon="file_copy"
        label="LOAD EXAMPLE"
        @click="loadExampleScriptTest"
      >
        <AppTooltip content="Load example test payload" />
      </q-btn>
    </div>
    <q-card flat bordered>
      <div
        ref="scriptTestEditorRef"
        class="monaco-editor-container"
        style="height: 400px;"
      ></div>
    </q-card>
    <div class="text-caption text-grey-7 q-mt-xs">
      JSON payload example used to test the conversion script
    </div>

    <div class="q-mt-md">
      <q-banner dense rounded class="bg-blue-1 text-primary">
        <template v-slot:avatar>
          <q-icon name="info" color="primary" size="sm" />
        </template>
        <div class="text-body2">
          {{ t.steps.step6.banner.info.value }}
        </div>
      </q-banner>
    </div>
  </div>
</template>

<style scoped>
.monaco-editor-container {
  border: 1px solid var(--mapex-card-border);
  border-radius: var(--mapex-radius-md);
  overflow: hidden;
}
</style>
