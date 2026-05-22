<script setup lang="ts">
defineOptions({ name: 'ExecutionNodeWrapper', inheritAttrs: false });

import { computed, ref, onMounted, nextTick } from 'vue';
import { GenericWorkflowNode } from '@src/components/workflow/nodes/GenericWorkflowNode';
import { usePluginRegistryStore } from '@stores/pluginRegistry';

const pluginRegistry = usePluginRegistryStore();

const props = defineProps<{
  id: string;
  data: { config: Record<string, unknown>; label?: string; __nodeType?: string; __execStatus?: string; hasErrors?: boolean };
  type: string;
}>();

// Aligned with runtime-emitted node statuses:
// completed / waiting / retrying / timeout / error / cancelled + "" (notReached)
const STATUS_ICONS: Record<string, string> = {
  completed: 'check',
  waiting:   'hourglass_top',
  retrying:  'replay',
  timeout:   'timer_off',
  error:     'close',
  cancelled: 'block',
  '':        'remove',
};

const wrapperRef = ref<HTMLElement | null>(null);
const badgeOffset = ref({ top: '-6px', right: '-6px' });
const execStatus = computed(() => (props.data?.__execStatus as string) || '');
const statusIcon = computed(() => STATUS_ICONS[execStatus.value] || '');

/** Resolve canvas component from plugin registry, fallback to GenericWorkflowNode */
const nodeComponent = computed(() => {
  const nodeType = pluginRegistry.nodeTypeMap.get(props.type);
  return nodeType?.canvasComponent ?? GenericWorkflowNode;
});

/**
 * Calculate badge position relative to the icon square
 */
function updateBadgePosition(): void {
  if (!wrapperRef.value) return;
  const icon = wrapperRef.value.querySelector('.wf-node__icon');
  if (!icon) return;
  const wrapperRect = wrapperRef.value.getBoundingClientRect();
  const iconRect = icon.getBoundingClientRect();
  const top = iconRect.top - wrapperRect.top - 6;
  const right = wrapperRect.right - iconRect.right - 6;
  badgeOffset.value = { top: `${top}px`, right: `${right}px` };
}

onMounted(() => {
  void nextTick(() => updateBadgePosition());
});
</script>

<template>
  <div ref="wrapperRef" style="position: relative">
    <component :is="nodeComponent" v-bind="$attrs" :id="id" :data="data" :type="type" />
    <div
      v-if="statusIcon"
      class="exec-status-badge"
      :class="`exec-status-badge--${execStatus || 'off'}`"
      :style="{ top: badgeOffset.top, right: badgeOffset.right }"
    >
      <q-icon :name="statusIcon" size="11px" />
    </div>
  </div>
</template>

<style scoped>
.exec-status-badge {
  position: absolute;
  width: 18px;
  height: 18px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 5;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.3);
  pointer-events: none;
}

.exec-status-badge--completed { background: #4caf50; color: #fff; }
.exec-status-badge--waiting   { background: #ff9800; color: #fff; }
.exec-status-badge--retrying  { background: #ffc107; color: #fff; }
.exec-status-badge--timeout   { background: #fdd835; color: #5d4037; }
.exec-status-badge--error     { background: #f44336; color: #fff; }
.exec-status-badge--cancelled { background: #616161; color: #fff; }
.exec-status-badge--off       { background: #9e9e9e; color: #f5f5f5; opacity: 0.6; }
</style>
