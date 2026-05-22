<script setup lang="ts">
/** TYPE IMPORTS */
import type { PluginNodeType, WorkflowNodeComponentProps } from '@src/components/workflow/interfaces';

/** VUE IMPORTS */
import { computed } from 'vue';
import { useI18n } from 'vue-i18n';

/** COMPONENTS */
import BaseWorkflowNode from '../BaseWorkflowNode/BaseWorkflowNode.vue';

/** COMPOSABLES */
import { useWorkflowContext } from '@src/composables/workflow';

/** UTILS */
import { resolveNodeHandles } from '@src/utils/workflow';

/** PROPS & EMITS */
const props = defineProps<WorkflowNodeComponentProps>();

/** COMPOSABLES & STORES */
const { getNodeType, nodes } = useWorkflowContext();
const i18n = useI18n();

/** COMPUTED */

/**
 * Resolved node type definition from plugin registry
 */
const nodeType = computed<PluginNodeType | undefined>(() => {
  if (!props.data.__nodeType) return undefined;
  return getNodeType(props.data.__nodeType);
});

/**
 * Translated node label resolved from plugin i18n with fallback
 */
const translatedLabel = computed(() => {
  const nt = nodeType.value;
  if (!nt?._pluginId) return props.data.label ?? nt?.label ?? '';

  const shortName = nt.type.split('/').pop() || '';
  const key = `wf.${nt._pluginId}.nodes.${shortName}.label`;
  if (!i18n.te(key)) return props.data.label ?? nt.label ?? '';
  return String(i18n.t(key));
});

/**
 * Resolved handles (dynamic or static) based on node type and config
 */
/**
 * Node timeout from the workflow node instance (for dynamic timeout handle)
 */
const nodeTimeout = computed(() => {
  const node = nodes.value.find(n => n.id === props.id);
  return node?.timeout;
});

const resolvedHandles = computed(() => {
  if (!nodeType.value) return { inputs: [], outputs: [] };
  return resolveNodeHandles(nodeType.value, props.data.config, nodeTimeout.value);
});

</script>

<template>
  <BaseWorkflowNode
    :id="id"
    :label="translatedLabel"
    :icon="nodeType?.icon || 'help'"
    :color="nodeType?.color || 'grey-7'"
    :selected="selected ?? false"
    :inputs="resolvedHandles.inputs"
    :outputs="resolvedHandles.outputs"
    :deletable="nodeType?.deletable ?? true"
    :shape="nodeType?.shape ?? 'square'"
    :has-errors="!!data.hasErrors"
  />
</template>
