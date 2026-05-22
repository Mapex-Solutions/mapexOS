<script setup lang="ts">
defineOptions({ name: 'ExecutionNode' });

import { computed } from 'vue';
import { Handle, Position } from '@vue-flow/core';

const props = defineProps<{
  id: string;
  data: {
    label?: string;
    __nodeType?: string;
    __status?: string;
    __icon?: string;
    __color?: string;
  };
}>();

const label = computed(() => props.data?.label || props.id);
const nodeType = computed(() => props.data?.__nodeType || '');
const icon = computed(() => props.data?.__icon || 'memory');
const color = computed(() => props.data?.__color || '#9e9e9e');

const shortType = computed(() => {
  const t = nodeType.value;
  return t.includes('/') ? t.split('/')[1] : t;
});
</script>

<template>
  <div class="exec-node">
    <Handle type="target" :position="Position.Top" class="exec-node__handle" />

    <div class="exec-node__body">
      <q-icon :name="icon" size="18px" :style="{ color }" class="q-mr-xs" />
      <div class="exec-node__info">
        <div class="exec-node__label">{{ label }}</div>
        <div class="exec-node__type">{{ shortType }}</div>
      </div>
    </div>

    <Handle type="source" :position="Position.Bottom" class="exec-node__handle" />
  </div>
</template>

<style scoped>
.exec-node {
  min-width: 120px;
  max-width: 220px;
}

.exec-node__body {
  display: flex;
  align-items: center;
  padding: 8px 12px;
  background: var(--mapex-surface-primary, #1e1e2e);
  border-radius: 8px;
  border: inherit;
}

.exec-node__info {
  overflow: hidden;
}

.exec-node__label {
  font-size: 12px;
  font-weight: 600;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  color: var(--mapex-text-primary, #fff);
}

.exec-node__type {
  font-size: 10px;
  color: var(--mapex-text-secondary, #999);
}

.exec-node__handle {
  width: 6px;
  height: 6px;
  background: var(--mapex-border-secondary, #555);
  border: none;
}
</style>
