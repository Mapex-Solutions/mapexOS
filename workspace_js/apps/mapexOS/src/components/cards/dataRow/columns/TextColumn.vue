<template>
  <div style="min-width: 0; display: flex; flex-direction: column; gap: 2px;">
    <!-- Primary Text -->
    <div
        :class="props.column.ellipsis ? 'ellipsis' : ''"
        :style="mobile ? 'min-width: 0; flex: 1;' : 'min-width: 0;'"
        class="text-body2 text-weight-medium"
    >
      {{ displayValue }}
      <AppTooltip v-if="props.column.ellipsis" :content="displayValue" />
    </div>

    <!-- Secondary Text (if secondaryKey exists) -->
    <div
        v-if="secondaryValue"
        class="text-caption text-grey-6 ellipsis"
        style="min-width: 0;"
    >
      {{ secondaryValue }}
      <AppTooltip :content="secondaryValue" />
    </div>
  </div>
</template>

<script setup lang="ts">
defineOptions({
  name: 'TextColumn'
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

const secondaryValue = computed(() => {
  // Support secondary as function (like format)
  if (props.column.secondary) {
    return props.column.secondary(props.value, props.row);
  }

  // Support secondaryKey as property path
  if (props.column.secondaryKey) {
    const keys = props.column.secondaryKey.split('.');
    let value = props.row;

    for (const key of keys) {
      value = value?.[key];
    }

    return value || null;
  }

  return null;
});
</script>
