<script setup lang="ts">
/** TYPE IMPORTS */
import type { EdgeProps } from '@vue-flow/core';

/** VUE IMPORTS */
import { computed } from 'vue';
import { BaseEdge, getStraightPath } from '@vue-flow/core';

/** PROPS */
const props = defineProps<EdgeProps>();

/** COMPUTED */

/**
 * Compute straight path between source and target
 */
const pathData = computed(() =>
  getStraightPath({
    sourceX: props.sourceX,
    sourceY: props.sourceY,
    targetX: props.targetX,
    targetY: props.targetY,
  }),
);

/**
 * SVG path string
 */
const edgePath = computed(() => pathData.value[0]);
</script>

<template>
  <BaseEdge
    :id="id"
    :path="edgePath"
    :style="{
      stroke: 'var(--mapex-text-muted)',
      strokeWidth: 1.5,
      strokeDasharray: '4 3',
      opacity: 0.5,
    }"
  />
</template>
