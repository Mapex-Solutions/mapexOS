<script setup lang="ts">
defineOptions({
  name: 'GotoNode',
});

/** TYPE IMPORTS */
import type { PluginNodeType, WorkflowNodeComponentProps } from '@src/components/workflow/interfaces';

/** VUE IMPORTS */
import { computed } from 'vue';

/** COMPONENTS */
import BaseWorkflowNode from '../BaseWorkflowNode/BaseWorkflowNode.vue';

/** COMPOSABLES */
import { useWorkflowContext, usePluginI18n } from '@src/composables/workflow';

/** UTILS */
import { resolveNodeHandles } from '@src/utils/workflow';

/** LOCAL IMPORTS */
import { GOTO_COLOR_OPTIONS } from './constants';

/** PROPS & EMITS */
const props = defineProps<WorkflowNodeComponentProps>();

/** COMPOSABLES & STORES */
const { getNodeType, nodes } = useWorkflowContext();
const { t } = usePluginI18n('core-flow-control');

/** COMPUTED */

/**
 * Resolved node type definition from plugin registry
 */
const nodeType = computed<PluginNodeType | undefined>(() => {
  if (!props.data.__nodeType) return undefined;
  return getNodeType(props.data.__nodeType);
});

/**
 * Resolved handles (dynamic via resolveInputs/resolveOutputs)
 */
const resolvedHandles = computed(() => {
  if (!nodeType.value) return { inputs: [], outputs: [] };
  return resolveNodeHandles(nodeType.value, props.data.config);
});

/**
 * Current role from node config
 */
const role = computed<'sender' | 'receiver'>(
  () => (props.data.config?.role as 'sender' | 'receiver') || 'sender',
);

/**
 * Current pair label from config
 */
const pairLabel = computed<string>(
  () => (props.data.config?.pairLabel as string) || '',
);

/**
 * Effective pair color — receivers inherit from matched sender dynamically
 */
const pairColor = computed<string>(() => {
  const own = (props.data.config?.pairColor as string) || 'deep-purple-6';
  if (role.value === 'receiver' && pairLabel.value) {
    const sender = nodes.value.find(n =>
      n.type === 'core/goto' &&
      n.id !== props.id &&
      (n.config?.role as string) === 'sender' &&
      (n.config?.pairLabel as string) === pairLabel.value,
    );
    if (sender) return (sender.config?.pairColor as string) || 'deep-purple-6';
  }
  return own;
});

/**
 * Hex color resolved from GOTO_COLOR_OPTIONS or custom hex
 */
const colorHex = computed<string>(() => {
  const opt = GOTO_COLOR_OPTIONS.find(o => o.value === pairColor.value);
  if (opt) return opt.hex;
  return pairColor.value.startsWith('#') ? pairColor.value : '#5e35b1';
});

/**
 * Icon based on role — correlated navigation pair
 */
const roleIcon = computed<string>(() =>
  role.value === 'sender' ? 'near_me' : 'place',
);
</script>

<template>
  <div class="goto-node-wrapper">
    <!-- Receiver: badge above -->
    <div
      v-if="pairLabel && role === 'receiver'"
      class="goto-node-wrapper__badge goto-node-wrapper__badge--top"
      :style="{ background: colorHex }"
    >
      {{ pairLabel }}
    </div>

    <BaseWorkflowNode
      :id="id"
      :label="data.label ?? t('nodes.goto.label')"
      :icon="roleIcon"
      :color="pairColor"
      :color-hex="colorHex"
      :selected="selected ?? false"
      :inputs="resolvedHandles.inputs"
      :outputs="resolvedHandles.outputs"
      :has-errors="!!data.hasErrors"
    />

    <!-- Sender: badge below -->
    <div
      v-if="pairLabel && role === 'sender'"
      class="goto-node-wrapper__badge goto-node-wrapper__badge--bottom"
      :style="{ background: colorHex }"
    >
      {{ pairLabel }}
    </div>
  </div>
</template>

<style lang="scss" scoped>
.goto-node-wrapper {
  display: flex;
  flex-direction: column;
  align-items: center;

  &__badge {
    font-size: 8px;
    font-weight: 700;
    letter-spacing: 0.5px;
    color: var(--mapex-wf-text-on-accent, #fff);
    padding: 1px 6px;
    border-radius: var(--mapex-radius-sm);
    white-space: nowrap;
    max-width: 80px;
    overflow: hidden;
    text-overflow: ellipsis;
    text-align: center;
    line-height: 1.3;

    &--top {
      margin-bottom: 2px;
    }

    &--bottom {
      margin-top: 2px;
    }
  }
}
</style>
