<script setup lang="ts">
defineOptions({
  name: 'CodeNodeConfig',
});

/** TYPE IMPORTS */
import type {
  NodeConfigComponentProps,
  NodeConfigComponentEmits,
} from '@src/components/workflow/interfaces';
import type { ScriptEditorGuideline } from '@components/dialogs/scriptEditor';

/** VUE IMPORTS */
import { ref, computed } from 'vue';

/** COMPONENTS */
import { ScriptEditorDialog } from '@components/dialogs/scriptEditor';

/** COMPOSABLES */
import { usePluginI18n } from '@src/composables/workflow';

/** PROPS & EMITS */
const props = defineProps<NodeConfigComponentProps>();
const emit = defineEmits<NodeConfigComponentEmits>();

/** COMPOSABLES & STORES */
const { t } = usePluginI18n('core-data');

/** STATE */

/**
 * Whether the script editor dialog is open
 */
const editorOpen = ref(false);

/** COMPUTED */

/**
 * Current script content from config
 */
const script = computed<string>(
  () => (props.config.script as string) ?? '// Access: state, event, inputs, nodes\n\nreturn {};',
);

/**
 * Execution timeout in milliseconds
 */
const timeout = computed<number>(
  () => (props.config.timeout as number) ?? 5000,
);

/**
 * Preview of the script (first few lines)
 */
const scriptPreview = computed(() => {
  const lines = script.value.split('\n');
  const preview = lines.slice(0, 8).join('\n');
  return lines.length > 8 ? preview + '\n...' : preview;
});

/**
 * Line count for display
 */
const lineCount = computed(() => script.value.split('\n').length);

/**
 * Guidelines for the script editor dialog
 */
const editorGuidelines = computed<ScriptEditorGuideline[]>(() => [
  { code: 'state', description: 'workflow state' },
  { code: 'event', description: 'trigger payload' },
  { code: 'inputs', description: 'external inputs' },
  { code: 'nodes', description: 'previous outputs' },
]);

/** FUNCTIONS */

/**
 * Emit config update with partial merge
 *
 * @param {Record<string, unknown>} partial - Partial config to merge
 */
function emitUpdate(partial: Record<string, unknown>): void {
  emit('update:config', { ...props.config, ...partial });
}

/**
 * Handle script content change from editor dialog
 *
 * @param {string} value - Updated script content
 */
function handleScriptChange(value: string): void {
  emitUpdate({ script: value });
}

/**
 * Update execution timeout
 *
 * @param {number} value - Timeout in milliseconds
 */
function updateTimeout(value: number): void {
  emitUpdate({ timeout: Math.max(100, value) });
}
</script>

<template>
  <div class="code-config">
    <!-- Script Section -->
    <div class="code-config__section">
      <div class="code-config__section-header">
        <div class="code-config__section-label">{{ t('nodes.code.config.scriptSection') }}</div>
        <q-badge color="grey-7" :label="`${lineCount} ${t('nodes.code.config.linesBadge')}`" />
      </div>

      <!-- Preview (read-only) -->
      <div class="code-config__preview" @click="editorOpen = true">
        <pre class="code-config__preview-code">{{ scriptPreview }}</pre>
        <div class="code-config__preview-overlay">
          <q-icon name="open_in_full" size="sm" />
        </div>
      </div>

      <!-- Open Editor button -->
      <q-btn
        outline
        no-caps
        dense
        color="primary"
        icon="code"
        :label="t('nodes.code.config.openEditor')"
        class="full-width q-mt-sm"
        @click="editorOpen = true"
      />
    </div>

    <!-- Timeout Section -->
    <div class="code-config__section">
      <div class="code-config__section-label">{{ t('nodes.code.config.timeoutSection') }}</div>

      <q-input
        :model-value="timeout"
        outlined
        dense
        type="number"
        suffix="ms"
        :hint="t('nodes.code.config.timeoutHint')"
        :rules="[(val: number) => val >= 100 || t('nodes.code.config.timeoutMin')]"
        @update:model-value="(val: string | number | null) => updateTimeout(Number(val ?? 5000))"
      >
        <template #prepend>
          <q-icon name="timer" color="grey-6" size="xs" />
        </template>
      </q-input>
    </div>

    <!-- Guidelines -->
    <div class="code-config__hint">
      <q-icon name="info" color="blue-6" size="xs" class="q-mr-sm" />
      <span class="text-caption" style="color: var(--mapex-text-secondary);">
        Available: <code>state</code>, <code>event</code>, <code>inputs</code>, <code>nodes</code>.
      </span>
    </div>

    <!-- Return object banner -->
    <q-banner dense rounded class="q-mt-sm" style="background: var(--mapex-wf-tint-2); border: 1px solid var(--mapex-wf-tint-border);">
      <template #avatar>
        <q-icon name="lightbulb" color="amber-8" size="xs" />
      </template>
      <span class="text-caption" style="color: var(--mapex-text-secondary);">
        To use the output in other nodes, <code>return</code> must be an <strong>object</strong>.
        <br />
        <code>return { message: 'hello' }</code> &rarr; accessible as <code>nodes.&lt;id&gt;.message</code>
      </span>
    </q-banner>

    <!-- Script Editor Dialog (reusable component) -->
    <ScriptEditorDialog
      v-model="editorOpen"
      :title="t('nodes.code.config.scriptEditorTitle')"
      :script-content="script"
      language="javascript"
      :guidelines="editorGuidelines"
      @update:script-content="handleScriptChange"
    >
      <template #guidelines-right>
        <span>
          <q-icon name="warning" color="amber-8" size="xs" class="q-mr-xs" />
          {{ t('nodes.code.config.sandboxedHint') }}
        </span>
      </template>
    </ScriptEditorDialog>
  </div>
</template>

<style lang="scss" scoped>
.code-config {
  &__section {
    margin-bottom: 16px;
  }

  &__section-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 6px;
  }

  &__section-label {
    font-size: 0.65rem;
    font-weight: 700;
    letter-spacing: 0.5px;
    color: var(--mapex-text-secondary);
    text-transform: uppercase;
  }

  &__preview {
    position: relative;
    border: 1px solid var(--mapex-card-border);
    border-radius: var(--mapex-radius-md);
    background: var(--mapex-surface-elevated);
    cursor: pointer;
    overflow: hidden;
    transition: border-color var(--mapex-transition-fast);

    &:hover {
      border-color: var(--mapex-primary);
    }

    &:hover &-overlay {
      opacity: 1;
    }
  }

  &__preview-code {
    margin: 0;
    padding: 10px 12px;
    font-family: 'Roboto Mono', monospace;
    font-size: 0.7rem;
    line-height: 1.5;
    color: var(--mapex-text-secondary);
    white-space: pre;
    overflow: hidden;
    max-height: 150px;
  }

  &__preview-overlay {
    position: absolute;
    inset: 0;
    display: flex;
    align-items: center;
    justify-content: center;
    background: var(--mapex-hover-overlay);
    color: var(--mapex-text-primary);
    opacity: 0;
    transition: opacity var(--mapex-transition-fast);
  }

  &__hint {
    display: flex;
    align-items: flex-start;
    padding: 8px 12px;
    border-radius: var(--mapex-radius-md);
    background: var(--mapex-wf-tint-2);
    border: 1px solid var(--mapex-wf-tint-border);

    code {
      font-size: 0.7rem;
      padding: 1px 4px;
      border-radius: var(--mapex-radius-xs);
      background: var(--mapex-surface-elevated);
    }
  }
}
</style>
