<template>
  <q-chip
      dense
      outline
      size="sm"
      text-color="white"
      class="text-weight-medium"
      style="width: fit-content; max-width: 100%;"
      :color="getColor()"
      :label="displayValue"
  >
    <AppTooltip :content="displayValue" />
  </q-chip>
</template>

<script setup lang="ts">
defineOptions({
  name: 'BadgeColumn'
});

import { computed } from 'vue';
import type { DataRowColumn } from '../interfaces';
import { AppTooltip } from '@components/tooltips';

const props = defineProps<{
  value: any;
  column: DataRowColumn;
  row: any;
  mobile?: boolean;
}>();

const displayValue = computed(() => {
  if (props.column.format) {
    return props.column.format(props.value, props.row);
  }
  return props.value || 'N/A';
});

function getColor() {
  if (typeof props.column.color === 'function') {
    return props.column.color(props.value, props.row);
  }
  return props.column.color || 'grey';
}
</script>
