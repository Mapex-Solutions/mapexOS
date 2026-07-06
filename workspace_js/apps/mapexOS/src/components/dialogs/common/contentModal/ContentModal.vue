<template>
  <q-dialog v-model="isOpen" @hide="handleClose">
    <q-card class="content-modal" style="min-width: 640px; max-width: 860px;">
      <!-- Header -->
      <q-card-section class="row items-center q-pb-none">
        <q-icon
          v-if="icon"
          :name="icon"
          size="md"
          color="primary"
          class="q-mr-md"
        />
        <div class="text-h6 text-weight-medium">{{ title }}</div>
        <q-space />
        <q-btn
          v-close-popup
          flat
          round
          dense
          icon="close"
          color="grey-7"
        />
      </q-card-section>

      <q-separator class="q-my-sm" />

      <!-- Content -->
      <q-card-section class="q-pt-md">
        <slot />
      </q-card-section>

      <!-- Footer Actions (optional) -->
      <q-card-actions v-if="$slots.actions" align="right" class="q-px-md q-pb-md">
        <slot name="actions" />
      </q-card-actions>
    </q-card>
  </q-dialog>
</template>

<script setup lang="ts">
defineOptions({
  name: 'ContentModal'
});

/** TYPE IMPORTS */
import type { ContentModalProps, ContentModalEmits } from './interfaces';

/** VUE IMPORTS */
import { ref, watch } from 'vue';

/** PROPS & EMITS */
const props = defineProps<ContentModalProps>();

const emit = defineEmits<ContentModalEmits>();

/** STATE */
const isOpen = ref(props.modelValue);

/** WATCHERS */
watch(
  () => props.modelValue,
  (newValue) => {
    isOpen.value = newValue;
  }
);

watch(isOpen, (newValue) => {
  emit('update:modelValue', newValue);
});

/** FUNCTIONS */

/**
 * Emit a closed state when the dialog is hidden
 * @returns {void}
 */
function handleClose(): void {
  emit('update:modelValue', false);
}
</script>

<style lang="scss" scoped>
.content-modal {
  border-radius: var(--mapex-radius-md);

  .q-card__section {
    padding: 20px 24px;
  }
}
</style>
