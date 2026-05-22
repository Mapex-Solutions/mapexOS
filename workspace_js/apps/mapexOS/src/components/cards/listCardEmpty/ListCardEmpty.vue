<script setup lang="ts">
defineOptions({
  name: 'ListCardEmpty'
});

import type { ListCardEmptyProps } from './interfaces';

const props = defineProps<ListCardEmptyProps>();

// Emit declaration for click event
const emit = defineEmits<{
  (e: 'button-click'): void;
}>();

// Emit event when button is clicked
const handleClick = () => {
  emit('button-click');
};
</script>

<template>
  <div class="col-12">
    <q-card class="empty-state-card shadow-1 rounded-borders">
      <q-card-section class="column items-center q-py-xl">
        <div class="empty-state-icon-container q-mb-md">
          <q-icon color="primary" size="60px" class="empty-state-icon" :name="props.icon"/>
        </div>
        <div class="text-h5 text-weight-medium text-primary q-mb-sm">{{ props.title }}</div>
        <div class="text-body1 text-grey-7 text-center q-mb-lg">
          {{ props.description }}
        </div>
        <div class="row q-gutter-md">
          <q-btn
              class="rounded-borders"
              color="primary"
              :ripple="false"
              :label="props?.buttonLabel || 'Reset filter'"
              :icon="props?.buttonIcon || 'refresh'"
              @click="handleClick"
          />
        </div>
      </q-card-section>
    </q-card>
  </div>
</template>

<style scoped lang="scss">
.empty-state-card {
  margin-top: 25px;
  border: 1px solid var(--mapex-card-border);
  background: var(--mapex-surface-bg);
  max-height: 350px;
}

.empty-state-icon-container {
  background-color: rgba(var(--mapex-primary-rgb), 0.1);
  border-radius: var(--mapex-radius-full);
  width: 110px;
  height: 110px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.empty-state-icon {
  color: var(--mapex-primary) !important;
}
</style>