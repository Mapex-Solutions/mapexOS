<template>
  <code
      :class="props.column.ellipsis ? 'ellipsis' : ''"
      class="text-body2 text-weight-medium text-grey-8"
      style="font-family: 'Courier New', monospace; display: block; min-width: 0;"
  >
    {{ displayValue }}
    <AppTooltip v-if="props.column.ellipsis" :content="displayValue" />
  </code>
</template>

<script setup lang="ts">
defineOptions({
  name: 'CodeColumn'
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
  return props.value || '-';
});
</script>
