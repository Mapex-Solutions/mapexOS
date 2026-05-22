<script setup lang="ts">
defineOptions({
  name: 'AssetClassificationSelector'
});

import type { AssetClassificationSelectorProps, AssetClassificationSelectorEmits } from './interfaces';

import { GuidedModeSelector } from './components/guidedModeSelector';

defineProps<AssetClassificationSelectorProps>();
const emit = defineEmits<AssetClassificationSelectorEmits>();

// Forward value updates from GuidedModeSelector
function handleValueUpdate(value: any) {
  emit('update:modelValue', value);
}
</script>

<template>
  <div class="asset-classification-selector">
    <GuidedModeSelector
      :model-value="modelValue"
      :disabled="disabled"
      :required="required"
      @update:model-value="handleValueUpdate"
    />
  </div>
</template>

<style lang="scss" scoped>
.asset-classification-selector {
  width: 100%;
}
</style>
