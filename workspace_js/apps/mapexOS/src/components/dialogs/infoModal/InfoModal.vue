<template>
  <q-dialog v-model="isOpen" @hide="handleClose">
    <q-card class="info-modal" style="min-width: 500px; max-width: 600px;">
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
        <!-- Main description -->
        <div class="text-body1 text-grey-8 q-mb-md" style="line-height: 1.6;">
          {{ description }}
        </div>

        <!-- Features/Items list (optional) -->
        <div v-if="items && items.length > 0" class="q-mb-md">
          <div
            v-for="(item, index) in items"
            :key="index"
            class="row items-start q-mb-sm"
          >
            <q-icon
              :name="item.icon || 'check_circle'"
              :color="item.color || 'primary'"
              size="sm"
              class="q-mr-sm q-mt-xs"
            />
            <div>
              <div v-if="item.title" class="text-weight-medium text-grey-9">
                {{ item.title }}
              </div>
              <div class="text-grey-7">{{ item.text }}</div>
            </div>
          </div>
        </div>

        <!-- Documentation link -->
        <div v-if="docsUrl" class="q-mt-lg">
          <q-btn
            flat
            dense
            no-caps
            icon="description"
            color="primary"
            :label="docsLabel || 'View Documentation'"
            :href="docsUrl"
            target="_blank"
            class="q-px-sm"
          >
            <q-icon name="open_in_new" size="xs" class="q-ml-xs" />
          </q-btn>
        </div>
      </q-card-section>

      <!-- Footer Actions (optional) -->
      <q-card-actions v-if="showActions" align="right" class="q-px-md q-pb-md">
        <q-btn
          v-close-popup
          flat
          no-caps
          :label="closeLabel || 'Got it'"
          color="primary"
          class="q-px-lg"
        />
      </q-card-actions>
    </q-card>
  </q-dialog>
</template>

<script setup lang="ts">
defineOptions({
  name: 'InfoModal'
});

import type { InfoModalProps, InfoModalEmits } from './interfaces';

import { ref, watch } from 'vue';

const props = withDefaults(defineProps<InfoModalProps>(), {
  showActions: true,
});

const emit = defineEmits<InfoModalEmits>();

const isOpen = ref(props.modelValue);

watch(
  () => props.modelValue,
  (newValue) => {
    isOpen.value = newValue;
  }
);

watch(isOpen, (newValue) => {
  emit('update:modelValue', newValue);
});

function handleClose() {
  emit('update:modelValue', false);
}
</script>

<style lang="scss" scoped>
.info-modal {
  border-radius: var(--mapex-radius-md);

  .q-card__section {
    padding: 20px 24px;
  }
}
</style>
