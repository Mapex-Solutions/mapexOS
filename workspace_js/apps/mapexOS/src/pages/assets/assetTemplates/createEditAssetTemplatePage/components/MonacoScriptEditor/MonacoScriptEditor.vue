<script setup lang="ts">
defineOptions({
  name: 'MonacoScriptEditor'
});

import { ref, onMounted, onBeforeUnmount, watch } from 'vue';
import { createMonacoEditorManager } from '../../handlers';
import type { MonacoEditorManager } from '../../handlers';
import { applyMapexMonacoTheme } from '@utils/monaco-theme';
import { useThemeStore } from '@stores/theme';

// Props
interface Props {
  modelValue: string;
  title: string;
  subtitle: string;
  icon?: string;
  guidelines: {
    title: string;
    availableVariables: string;
    items: Array<{ label: string; description: string }>;
    expectedReturn: string;
    expectedReturnValue: string;
    additionalInfo?: string;
  };
  language?: 'javascript' | 'json';
  height?: string;
  errorMessage?: string;
}

// Emits
interface Emits {
  (e: 'update:modelValue', value: string): void;
}

const props = withDefaults(defineProps<Props>(), {
  icon: 'code',
  language: 'javascript',
  height: '400px',
  errorMessage: '',
});

const emit = defineEmits<Emits>();
const themeStore = useThemeStore();

// Watch theme changes to update Monaco editor
watch(() => themeStore.isDark, (isDark: boolean) => {
  applyMapexMonacoTheme(isDark);
});

// Editor container ref
const editorContainerRef = ref<HTMLElement | null>(null);

// Editor manager
let editorManager: MonacoEditorManager | null = null;

// Setup editor on mount
onMounted(() => {
  editorManager = createMonacoEditorManager(
    editorContainerRef,
    { language: props.language },
    props.modelValue,
    (value: string) => {
      emit('update:modelValue', value);
    }
  );
  editorManager.setup();
});

// Cleanup on unmount
onBeforeUnmount(() => {
  editorManager?.dispose();
  editorManager = null;
});

// Watch for external changes to modelValue
watch(() => props.modelValue, (newValue) => {
  if (editorManager && editorManager.getValue() !== newValue) {
    editorManager.setValue(newValue);
  }
});
</script>

<template>
  <div class="code-editor-section">
    <div class="q-mb-md">
      <div class="text-subtitle1 text-weight-medium q-mb-xs">
        <q-icon :name="icon" color="primary" class="q-mr-xs" />
        {{ title }}
      </div>
      <div class="text-body2 text-grey-7">
        {{ subtitle }}
      </div>
    </div>

    <q-card flat bordered class="code-editor-card">
      <div
        ref="editorContainerRef"
        class="monaco-editor-container"
        :style="{ height }"
      ></div>
    </q-card>

    <!-- Error Hint -->
    <div v-if="errorMessage" class="q-mt-sm text-negative text-caption">
      <q-icon name="warning" size="sm" class="q-mr-xs" />
      {{ errorMessage }}
    </div>

    <!-- Guidelines Expansion -->
    <div class="q-mt-md">
      <q-expansion-item
        class="text-primary"
        icon="help"
        :label="guidelines.title"
      >
        <q-card>
          <q-card-section class="text-body2">
            <p><strong>{{ guidelines.availableVariables }}</strong></p>
            <ul>
              <li v-for="(item, index) in guidelines.items" :key="index">
                <code>{{ item.label }}</code> - {{ item.description }}
              </li>
            </ul>

            <p><strong>{{ guidelines.expectedReturn }}</strong> {{ guidelines.expectedReturnValue }}</p>

            <div v-if="guidelines.additionalInfo" class="q-mt-md">
              <div v-html="guidelines.additionalInfo"></div>
            </div>
          </q-card-section>
        </q-card>
      </q-expansion-item>
    </div>
  </div>
</template>

<style scoped>
.code-editor-section .monaco-editor-container {
  border: 1px solid var(--mapex-card-border);
  border-radius: var(--mapex-radius-md);
  overflow: hidden;
}

.code-editor-card {
  border-radius: var(--mapex-radius-md) !important;
  overflow: hidden;
}
</style>
