<template>
  <q-avatar
      text-color="white"
      class="flex-shrink-0"
      style="box-shadow: var(--mapex-shadow-sm);"
      :size="mobile ? '36px' : '40px'"
      :color="getColor()"
      :icon="getIcon()"
  >
    <!-- Tooltip (only on laptop/desktop >= 1024px) -->
    <AppTooltip
        v-if="getTooltip()"
        :content="getTooltip()!"
        anchor="center right"
        self="center left"
        :offset="[10, 0]"
    />
  </q-avatar>
</template>

<script setup lang="ts">
defineOptions({
  name: 'AvatarColumn'
});

import type { DataRowColumn } from '../interfaces';
import { AppTooltip } from '@components/tooltips';

const props = defineProps<{
  value: any;
  column: DataRowColumn;
  row: any;
  mobile?: boolean;
}>();

function getIcon() {
  if (typeof props.column.icon === 'function') {
    return props.column.icon(props.value, props.row);
  }
  return props.column.icon || 'person';
}

function getColor() {
  if (typeof props.column.color === 'function') {
    return props.column.color(props.value, props.row);
  }
  return props.column.color || 'primary';
}

/**
 * Get tooltip text (only shown on laptop/desktop)
 * Returns null if no tooltip configured or on mobile
 */
function getTooltip(): string | null {
  // Don't show tooltip on mobile
  if (props.mobile) return null;

  // No tooltip configured
  if (!props.column.tooltip) return null;

  // Execute function or return static string
  if (typeof props.column.tooltip === 'function') {
    return props.column.tooltip(props.value, props.row);
  }

  return props.column.tooltip;
}
</script>
