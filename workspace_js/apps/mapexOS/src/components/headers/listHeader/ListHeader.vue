<script setup lang="ts">
defineOptions({
  name: 'ListHeader'
});

import type { ListHeadersProps } from './interfaces';

// Props will be available directly as variables (no 'props.')
const props = defineProps<ListHeadersProps>();
const { title, icon , button } = props;
</script>

<template>
  <div class="row items-center q-mb-md">
    <div class="col">
      <div class="row items-center">
        <q-icon size="sm" class="q-mr-sm" color="primary" :name="icon" />
        <div class="text-subtitle1 text-weight-medium text-primary">{{ title }}</div>
      </div>
    </div>
    <div class="col-auto">
      <q-btn
        outline
        rounded
        class="q-py-xs q-px-md"
        :label="button.label"
        :icon="button.icon"
        :color="button.color || 'primary'"
        :to="button.to"
      />
    </div>
  </div>
</template>
