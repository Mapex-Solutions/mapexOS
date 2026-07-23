<template>
  <div :class="containerClass">
    <q-icon :name="iconName" :color="iconColor" size="22px">
      <AppTooltip v-if="tooltipContent" :content="tooltipContent" />
    </q-icon>
  </div>
</template>

<script setup lang="ts">
defineOptions({
  name: 'IconColumn'
});

/** TYPE IMPORTS */
import type { DataRowColumn } from '../interfaces';

/** VUE IMPORTS */
import { computed } from 'vue';

/** COMPONENTS */
import { AppTooltip } from '@components/tooltips';

/** PROPS & EMITS */
const props = defineProps<{
  value: any;
  column: DataRowColumn;
  row: any;
  mobile?: boolean;
}>();

/** COMPUTED */

/**
 * Resolved icon name — from the column's icon (string or (value, row) fn).
 */
const iconName = computed<string>(() => {
  if (typeof props.column.icon === 'function') {
    return props.column.icon(props.value, props.row);
  }
  return props.column.icon || 'help_outline';
});

/**
 * Resolved icon color — from the column's color (string or (value, row) fn).
 */
const iconColor = computed<string>(() => {
  if (typeof props.column.color === 'function') {
    return props.column.color(props.value, props.row);
  }
  return props.column.color || 'grey-6';
});

/**
 * Resolved tooltip content — rendered through the shared AppTooltip.
 */
const tooltipContent = computed<string>(() => {
  if (typeof props.column.tooltip === 'function') {
    return props.column.tooltip(props.value, props.row);
  }
  return props.column.tooltip || '';
});

/**
 * Alignment container, mirroring the other column renderers.
 */
const containerClass = computed<string>(() => {
  const classes = ['flex', 'items-center'];
  if (props.column.align === 'center') {
    classes.push('justify-center');
  } else if (props.column.align === 'right') {
    classes.push('justify-end');
  } else {
    classes.push('justify-start');
  }
  return classes.join(' ');
});
</script>
