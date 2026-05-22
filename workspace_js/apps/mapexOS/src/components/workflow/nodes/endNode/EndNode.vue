<script setup lang="ts">
defineOptions({
  name: 'EndNode',
});

/** TYPE IMPORTS */
import type { PluginNodeType, WorkflowNodeComponentProps } from '@src/components/workflow/interfaces';

/** VUE IMPORTS */
import { computed } from 'vue';

/** COMPONENTS */
import BaseWorkflowNode from '../BaseWorkflowNode/BaseWorkflowNode.vue';

/** COMPOSABLES */
import { useWorkflowContext } from '@src/composables/workflow';

/** UTILS */
import { resolveNodeHandles } from '@src/utils/workflow';

/** LOCAL IMPORTS */
import { END_NODE_VARIANTS } from './constants';

/** PROPS & EMITS */
const props = defineProps<WorkflowNodeComponentProps>();

/** COMPOSABLES & STORES */
const { getNodeType } = useWorkflowContext();

/** COMPUTED */

/**
 * Resolved node type definition from plugin registry
 */
const nodeType = computed<PluginNodeType | undefined>(() => {
  if (!props.data.__nodeType) return undefined;
  return getNodeType(props.data.__nodeType);
});

/**
 * Resolved handles (dynamic or static) based on node type and config
 */
const resolvedHandles = computed(() => {
  if (!nodeType.value) return { inputs: [], outputs: [] };
  return resolveNodeHandles(nodeType.value, props.data.config);
});

/**
 * Whether the node is in error termination mode
 */
const isError = computed(() => props.data.config?.terminateWithError === true);

/**
 * Dynamic icon based on termination mode
 */
const nodeIcon = computed(() =>
  isError.value ? END_NODE_VARIANTS.error.icon : END_NODE_VARIANTS.success.icon,
);

/**
 * Dynamic hex color based on termination mode
 */
const nodeColorHex = computed(() =>
  isError.value ? END_NODE_VARIANTS.error.hex : END_NODE_VARIANTS.success.hex,
);
</script>

<template>
  <BaseWorkflowNode
    :id="id"
    :label="data.label ?? nodeType?.label ?? ''"
    :icon="nodeIcon"
    :color="isError ? 'red-7' : 'purple-7'"
    :color-hex="nodeColorHex"
    :selected="selected ?? false"
    :inputs="resolvedHandles.inputs"
    :outputs="resolvedHandles.outputs"
    :deletable="nodeType?.deletable ?? true"
    shape="circle"
    :has-errors="!!data.hasErrors"
  />
</template>
