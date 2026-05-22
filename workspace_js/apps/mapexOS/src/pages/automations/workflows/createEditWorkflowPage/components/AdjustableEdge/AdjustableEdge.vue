<script setup lang="ts">
/** TYPE IMPORTS */
import type { EdgeProps } from '@vue-flow/core';

/** VUE IMPORTS */
import { computed, ref, onUnmounted } from 'vue';
import { BaseEdge, EdgeLabelRenderer, getSmoothStepPath, useVueFlow } from '@vue-flow/core';

/** COMPOSABLES */
import { useWorkflowEditorState } from '../../composables';

/** PROPS */
const props = defineProps<EdgeProps>();

/** COMPOSABLES & STORES */
const { edges } = useWorkflowEditorState();
const { screenToFlowCoordinate } = useVueFlow();

/** STATE */
const isDragging = ref(false);
const dragStartFlow = ref<{ x: number; y: number } | null>(null);
const dragStartOffset = ref<{ x: number; y: number }>({ x: 0, y: 0 });

/** COMPUTED */

/**
 * Read current edge data from composable state
 */
const edgeData = computed(() => edges.value.find(e => e.id === props.id));

/**
 * Current X offset from natural center
 */
const currentOffsetX = computed(() => edgeData.value?.pathOffsetX ?? 0);

/**
 * Current Y offset from natural center
 */
const currentOffsetY = computed(() => edgeData.value?.pathOffsetY ?? 0);

/**
 * Shifted center X (natural midpoint + user offset)
 */
const shiftedCenterX = computed(() =>
  (props.sourceX + props.targetX) / 2 + currentOffsetX.value,
);

/**
 * Shifted center Y (natural midpoint + user offset)
 */
const shiftedCenterY = computed(() =>
  (props.sourceY + props.targetY) / 2 + currentOffsetY.value,
);

/**
 * Compute the smooth step path with the shifted center
 */
const pathData = computed(() =>
  getSmoothStepPath({
    sourceX: props.sourceX,
    sourceY: props.sourceY,
    sourcePosition: props.sourcePosition,
    targetX: props.targetX,
    targetY: props.targetY,
    targetPosition: props.targetPosition,
    borderRadius: 8,
    centerX: shiftedCenterX.value,
    centerY: shiftedCenterY.value,
  }),
);

/**
 * SVG path string
 */
const edgePath = computed(() => pathData.value[0]);

/**
 * Label X position (follows shifted center)
 */
const labelX = computed(() => pathData.value[1]);

/**
 * Label Y position (follows shifted center)
 */
const labelY = computed(() => pathData.value[2]);

/**
 * Resolved marker start (ensures string, never undefined)
 */
const resolvedMarkerStart = computed(() => props.markerStart || '');

/**
 * Resolved marker end (ensures string, never undefined)
 */
const resolvedMarkerEnd = computed(() => props.markerEnd || '');

/** FUNCTIONS */

/**
 * Start drag interaction on the handle.
 * Records starting mouse position (flow coords) and current offset.
 *
 * @param {MouseEvent} event - The mousedown event
 */
function onDragStart(event: MouseEvent): void {
  event.stopPropagation();
  event.preventDefault();
  isDragging.value = true;
  dragStartFlow.value = screenToFlowCoordinate({ x: event.clientX, y: event.clientY });
  dragStartOffset.value = { x: currentOffsetX.value, y: currentOffsetY.value };

  document.addEventListener('mousemove', onDragMove);
  document.addEventListener('mouseup', onDragEnd);
}

/**
 * Handle mouse move during drag.
 * Calculates delta in flow coordinates and updates the edge offset.
 *
 * @param {MouseEvent} event - The mousemove event
 */
function onDragMove(event: MouseEvent): void {
  if (!isDragging.value || !dragStartFlow.value) return;

  const currentFlow = screenToFlowCoordinate({ x: event.clientX, y: event.clientY });
  const deltaX = currentFlow.x - dragStartFlow.value.x;
  const deltaY = currentFlow.y - dragStartFlow.value.y;

  updateEdgeOffset(
    dragStartOffset.value.x + deltaX,
    dragStartOffset.value.y + deltaY,
  );
}

/**
 * End drag interaction. Clean up document listeners.
 */
function onDragEnd(): void {
  isDragging.value = false;
  dragStartFlow.value = null;
  document.removeEventListener('mousemove', onDragMove);
  document.removeEventListener('mouseup', onDragEnd);
}

/**
 * Persist offset to the composable edge state.
 *
 * @param {number} offsetX - New X offset
 * @param {number} offsetY - New Y offset
 */
function updateEdgeOffset(offsetX: number, offsetY: number): void {
  const edge = edges.value.find(e => e.id === props.id);
  if (!edge) return;
  edge.pathOffsetX = offsetX;
  edge.pathOffsetY = offsetY;
}

/**
 * Reset edge offset to zero on double-click
 */
function onResetOffset(): void {
  updateEdgeOffset(0, 0);
}

/** LIFECYCLE HOOKS */
onUnmounted(() => {
  document.removeEventListener('mousemove', onDragMove);
  document.removeEventListener('mouseup', onDragEnd);
});
</script>

<template>
  <!-- SVG path rendered by BaseEdge (handles markers, interaction area, animation) -->
  <BaseEdge
    :id="id"
    :path="edgePath"
    :label="label"
    :label-x="labelX"
    :label-y="labelY"
    :style="style"
    :marker-start="resolvedMarkerStart"
    :marker-end="resolvedMarkerEnd"
  />

  <!-- Draggable handle at midpoint (HTML overlay in flow coordinate space) -->
  <EdgeLabelRenderer>
    <div
      class="adjustable-edge-handle"
      :class="{
        'adjustable-edge-handle--dragging': isDragging,
        'adjustable-edge-handle--visible': selected,
      }"
      :style="{
        position: 'absolute',
        transform: `translate(-50%, -50%) translate(${labelX}px, ${labelY}px)`,
        pointerEvents: 'all',
      }"
      @mousedown="onDragStart"
      @dblclick="onResetOffset"
    >
      <div class="adjustable-edge-handle__dot" />
    </div>
  </EdgeLabelRenderer>
</template>

<style lang="scss" scoped>
.adjustable-edge-handle {
  width: 20px;
  height: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: grab;
  opacity: 0;
  transition: opacity 0.15s ease;
  z-index: 10;

  &:hover,
  &--visible,
  &--dragging {
    opacity: 1;
  }

  &--dragging {
    cursor: grabbing;
  }

  &__dot {
    width: 8px;
    height: 8px;
    border-radius: 50%;
    background: var(--wf-edge-color, var(--mapex-primary));
    border: 1.5px solid var(--wf-handle-border, #fff);
    box-shadow: 0 0 4px var(--wf-edge-glow, rgba(0, 0, 0, 0.2));
    transition: transform 0.15s;
  }

  &:hover &__dot {
    transform: scale(1.3);
  }
}
</style>
