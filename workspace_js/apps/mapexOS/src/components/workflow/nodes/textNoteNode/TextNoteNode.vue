<script setup lang="ts">
defineOptions({
  name: 'TextNoteNode',
});

/** TYPE IMPORTS */
import type { WorkflowNodeComponentProps } from '@src/components/workflow/interfaces';

/** VUE IMPORTS */
import { ref, computed, nextTick } from 'vue';

/** COMPOSABLES */
import { useWorkflowContext, usePluginI18n } from '@src/composables/workflow';

/** LOCAL IMPORTS */
import { NOTE_COLOR_OPTIONS } from './constants';

/** PROPS & EMITS */
const props = defineProps<WorkflowNodeComponentProps>();

/** COMPOSABLES & STORES */
const { updateNodeConfig } = useWorkflowContext();
const { t } = usePluginI18n('core-annotations');

/** STATE */

/**
 * Whether the text is being edited
 */
const isEditing = ref(false);

/**
 * Buffer for editing
 */
const editBuffer = ref('');

/**
 * Reference to the textarea for autofocus
 */
const textareaRef = ref<HTMLTextAreaElement | null>(null);

/** COMPUTED */

/**
 * Current text from node config
 */
const text = computed<string>(
  () => (props.data.config?.text as string) || '',
);

/**
 * Current background color from config
 */
const bgColor = computed<string>(
  () => (props.data.config?.color as string) || 'grey',
);

/** FUNCTIONS */

/**
 * Start editing the text
 */
function startEditing(): void {
  editBuffer.value = text.value;
  isEditing.value = true;
  void nextTick(() => textareaRef.value?.focus());
}

/**
 * Save the edited text
 */
function saveText(): void {
  isEditing.value = false;
  updateNodeConfig(props.id, { text: editBuffer.value });
}

/**
 * Handle keydown — Escape to cancel
 *
 * @param {KeyboardEvent} event - Keyboard event
 */
function handleKeydown(event: KeyboardEvent): void {
  if (event.key === 'Escape') {
    isEditing.value = false;
  }
}

/**
 * Change the note color
 *
 * @param {string} color - New color value
 */
function changeColor(color: string): void {
  updateNodeConfig(props.id, { color });
}
</script>

<template>
  <div
    class="text-note"
    :class="[
      `text-note--${bgColor}`,
      { 'text-note--selected': selected },
    ]"
    @dblclick.stop="startEditing"
  >
    <!-- Toolbar — visible on hover -->
    <div class="text-note__toolbar">
      <!-- Color picker -->
      <q-btn
        flat
        dense
        round
        size="xs"
        icon="palette"
        class="text-note__toolbar-btn"
      >
        <q-menu>
          <div class="text-note__color-grid">
            <div
              v-for="opt in NOTE_COLOR_OPTIONS"
              :key="opt.value"
              class="text-note__color-swatch"
              :class="{ 'text-note__color-swatch--active': bgColor === opt.value }"
              :style="{ background: opt.hex }"
              @click="changeColor(opt.value)"
            />
          </div>
        </q-menu>
      </q-btn>
    </div>

    <!-- Left accent strip -->
    <div class="text-note__strip" />

    <!-- Editing mode -->
    <textarea
      v-if="isEditing"
      ref="textareaRef"
      v-model="editBuffer"
      class="text-note__textarea"
      :placeholder="t('nodes.text_note.placeholder')"
      @blur="saveText"
      @keydown="handleKeydown"
    />

    <!-- Display mode -->
    <div v-else class="text-note__content">
      <div v-if="text" class="text-note__text">{{ text }}</div>
      <div v-else class="text-note__placeholder">{{ t('nodes.text_note.emptyPlaceholder') }}</div>
    </div>

  </div>
</template>

<style lang="scss" scoped>
.text-note {
  --note-accent: #78909c;

  position: relative;
  min-width: 60px;
  max-width: 160px;
  padding: 3px 6px 3px 10px;
  background: transparent;
  border: none;
  cursor: grab;
  overflow: hidden;
  opacity: 0.75;
  transition: opacity 0.2s ease;

  &:hover {
    opacity: 1;
  }

  &--selected {
    opacity: 1;
  }

  /* Color variants */
  &--amber { --note-accent: #ffb300; }
  &--blue { --note-accent: #42a5f5; }
  &--green { --note-accent: #66bb6a; }
  &--red { --note-accent: #ef5350; }
  &--purple { --note-accent: #ab47bc; }
  &--grey { --note-accent: #78909c; }

  /* Toolbar — top-right, visible on hover */
  &__toolbar {
    position: absolute;
    top: -2px;
    right: -2px;
    display: flex;
    gap: 2px;
    opacity: 0;
    transform: scale(0.8);
    transition: opacity 0.15s ease, transform 0.15s ease;
    z-index: 3;
  }

  &:hover &__toolbar {
    opacity: 1;
    transform: scale(1);
  }

  &__toolbar-btn {
    background: var(--mapex-surface-elevated) !important;
    border: 1px solid var(--mapex-card-border) !important;
    color: var(--mapex-text-secondary) !important;
    width: 18px !important;
    height: 18px !important;

    &:hover {
      color: var(--mapex-text-primary) !important;
      border-color: var(--note-accent) !important;
    }
  }

  /* Left accent strip — thin subtle line */
  &__strip {
    position: absolute;
    top: 2px;
    left: 2px;
    bottom: 2px;
    width: 2px;
    border-radius: var(--mapex-radius-full);
    background: var(--note-accent);
    opacity: 0.65;
  }

  &__content {
    user-select: none;
    width: 100%;
    overflow: hidden;
  }

  &__text {
    font-size: 9px;
    line-height: 1.1;
    color: var(--mapex-text-secondary);
    white-space: pre-wrap;
    word-break: break-word;
    font-style: italic;
  }

  &__placeholder {
    font-size: 9px;
    color: var(--mapex-text-muted);
    font-style: italic;
  }

  &__textarea {
    width: 100%;
    min-height: 16px;
    border: none;
    outline: none;
    background: transparent;
    color: var(--mapex-text-secondary);
    font-size: 9px;
    line-height: 1.1;
    font-family: inherit;
    font-style: italic;
    resize: none;
    padding: 0;

    &::placeholder {
      color: var(--mapex-text-muted);
      font-style: italic;
    }
  }

}
</style>

<!-- Unscoped: Quasar teleports q-menu to <body> -->
<style lang="scss">
.text-note__color-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 6px;
  padding: 8px;
  background: var(--mapex-surface-elevated);
  border-radius: var(--mapex-radius-md);
}

.text-note__color-swatch {
  width: 22px;
  height: 22px;
  border-radius: var(--mapex-radius-full);
  cursor: pointer;
  border: 2px solid transparent;
  transition: transform 0.15s ease, border-color 0.15s ease;

  &:hover {
    transform: scale(1.2);
  }

  &--active {
    border-color: var(--mapex-wf-text-on-accent, #fff);
    box-shadow: var(--mapex-wf-selection-ring, 0 0 0 2px rgba(255, 255, 255, 0.3));
  }
}
</style>
