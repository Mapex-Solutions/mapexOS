<script setup lang="ts">
defineOptions({
  name: 'ConditionalRoutingToggle',
});

/** TYPE IMPORTS */
import type { ConditionalRoutingToggleProps, ConditionalRoutingToggleEmits } from './interfaces/ConditionalRoutingToggle.interface';

/** PROPS & EMITS */
defineProps<ConditionalRoutingToggleProps>();
const emit = defineEmits<ConditionalRoutingToggleEmits>();
</script>

<template>
  <q-card
    flat
    bordered
    class="q-mb-md"
  >
    <q-card-section class="row items-center q-py-md">
      <div class="col">
        <div class="text-body1 text-weight-medium text-grey-9">
          {{ t.createEdit.routersStep.conditionalRouting.label.value }}
        </div>
        <div class="text-caption text-grey-6 q-mt-xs">
          {{ t.createEdit.routersStep.conditionalRouting.hint.value }}
        </div>
      </div>
      <div class="col-auto">
        <q-toggle
          color="secondary"
          size="md"
          checked-icon="check"
          unchecked-icon="close"
          :model-value="modelValue"
          @update:model-value="emit('update:modelValue', $event)"
        />
      </div>
    </q-card-section>
  </q-card>
</template>
