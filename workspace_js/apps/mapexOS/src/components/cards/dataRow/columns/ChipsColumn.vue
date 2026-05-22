<template>
  <div :class="containerClass">
    <!-- Visible chips (up to maxVisible) -->
    <q-chip
      v-for="(item, index) in visibleItems"
      :key="index"
      dense
      outline
      size="sm"
      text-color="white"
      class="text-weight-medium q-mr-xs"
      style="width: fit-content; max-width: 100%;"
      :color="getColor()"
    >
      {{ item }}
      <AppTooltip :content="item" />
    </q-chip>

    <!-- "See more" chip when there are hidden items -->
    <q-chip
      v-if="hiddenCount > 0"
      dense
      size="sm"
      color="grey-6"
      text-color="white"
      class="text-weight-medium cursor-pointer"
      clickable
      @click.stop="showAllDialog = true"
    >
      +{{ hiddenCount }}
      <AppTooltip :content="`${hiddenCount} more items`" />
    </q-chip>

    <!-- Dialog to show all items -->
    <q-dialog v-model="showAllDialog">
      <q-card style="min-width: 300px; max-width: 500px;">
        <q-card-section class="row items-center">
          <div class="text-h6">{{ column.label || 'All Items' }}</div>
          <q-space />
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>

        <q-separator />

        <q-card-section class="q-pt-md">
          <div class="row q-gutter-sm">
            <q-chip
              v-for="(item, index) in allItems"
              :key="index"
              dense
              outline
              size="sm"
              text-color="white"
              class="text-weight-medium"
              :color="getColor()"
            >
              {{ item }}
            </q-chip>
          </div>
        </q-card-section>
      </q-card>
    </q-dialog>
  </div>
</template>

<script setup lang="ts">
defineOptions({
  name: 'ChipsColumn'
});

import { ref, computed } from 'vue';
import type { DataRowColumn } from '../interfaces';
import { AppTooltip } from '@components/tooltips';

/** PROPS */
const props = withDefaults(defineProps<{
  value: any;
  column: DataRowColumn;
  row: any;
  mobile?: boolean;
}>(), {
  mobile: false,
});

/** STATE */
const showAllDialog = ref(false);

/** CONSTANTS */
const MAX_VISIBLE = 3;

/** COMPUTED */

/**
 * All items as array
 */
const allItems = computed(() => {
  if (Array.isArray(props.value)) {
    return props.value;
  }
  if (typeof props.value === 'string') {
    return [props.value];
  }
  return [];
});

/**
 * Visible items (limited to MAX_VISIBLE)
 */
const visibleItems = computed(() => {
  return allItems.value.slice(0, MAX_VISIBLE);
});

/**
 * Count of hidden items
 */
const hiddenCount = computed(() => {
  return Math.max(0, allItems.value.length - MAX_VISIBLE);
});

/**
 * Container alignment class
 */
const containerClass = computed(() => {
  const classes = ['flex', 'items-center', 'flex-wrap'];

  if (props.column.align === 'center') {
    classes.push('justify-center');
  } else if (props.column.align === 'right') {
    classes.push('justify-end');
  } else {
    classes.push('justify-start');
  }

  return classes.join(' ');
});

/** FUNCTIONS */

/**
 * Get chip color from column config
 * @returns {string} Color name
 */
function getColor(): string {
  if (typeof props.column.color === 'function') {
    return props.column.color(props.value, props.row);
  }
  return (props.column.color as string) || 'primary';
}
</script>
