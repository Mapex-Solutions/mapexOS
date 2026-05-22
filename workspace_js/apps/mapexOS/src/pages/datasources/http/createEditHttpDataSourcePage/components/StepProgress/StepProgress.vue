<template>
  <q-card class="stepper-card rounded-borders shadow-1">
    <q-card-section class="bg-primary text-white q-pb-md">
      <div class="text-h6 text-weight-bold q-mb-sm">
        <q-icon name="rule" size="sm" class="q-mr-xs"/>
        Progress Steps
      </div>
      <div class="text-caption">{{ t.progress.completeAllSteps.value }}</div>
    </q-card-section>

    <q-card-section class="q-pa-md">
      <div class="progress-steps">
        <div
          v-for="(st, idx) in steps"
          :key="idx"
          class="step-item"
          :class="{ active: currentStep === idx + 1, completed: currentStep > idx + 1 }"
        >
          <div class="step-icon-wrapper">
            <div class="step-icon">
              <q-icon
                :name="currentStep > idx + 1 ? 'check_circle' : st.icon"
                :color="currentStep >= idx + 1 ? 'white' : 'grey-5'"
                size="sm"
              />
            </div>
          </div>
          <div class="step-content">
            <div class="step-title">{{ st.label }}</div>
            <div class="step-description">{{ st.description }}</div>
          </div>
        </div>
      </div>

      <q-separator class="q-my-md"/>

      <div class="current-step-info">
        <div class="text-caption text-grey-6">
          <q-icon name="info" size="xs" class="q-mr-xs"/>
          All fields marked with * are required
        </div>
        <div class="text-caption text-primary q-mt-xs">
          <q-icon name="arrow_forward" size="xs" class="q-mr-xs"/>
          Current Step: {{ steps[currentStep - 1]?.label }}
        </div>
      </div>
    </q-card-section>
  </q-card>
</template>

<script setup lang="ts">
/** TYPE IMPORTS */
import type { StepProgressProps } from './interfaces/StepProgress.interface';

import { useHttpDataSourceCreateEditTranslations } from '@composables/i18n/pages/datasources/http';

defineOptions({
  name: 'StepProgress'
});

defineProps<StepProgressProps>();

const t = useHttpDataSourceCreateEditTranslations();
</script>

<style scoped>
.rounded-borders {
  border-radius: var(--mapex-radius-md);
}

.stepper-card {
  height: 100%;
}

.progress-steps {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.step-item {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  padding: 8px 0;
  position: relative;
}

.step-item:not(:last-child)::after {
  content: '';
  position: absolute;
  left: 19px;
  top: 40px;
  width: 2px;
  height: 24px;
  background-color: var(--mapex-card-border);
  z-index: 0;
}

.step-item.completed:not(:last-child)::after {
  background-color: var(--q-positive);
}

.step-icon-wrapper {
  position: relative;
  z-index: 1;
}

.step-icon {
  width: 40px;
  height: 40px;
  border-radius: var(--mapex-radius-full);
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: var(--mapex-surface-sunken);
  border: 2px solid var(--mapex-card-border);
  transition: var(--mapex-transition-slow);
}

.step-item.active .step-icon {
  background-color: var(--q-primary);
  border-color: var(--q-primary);
}

.step-item.completed .step-icon {
  background-color: var(--q-positive);
  border-color: var(--q-positive);
}

.step-content {
  flex: 1;
  padding-top: 2px;
}

.step-title {
  font-weight: 600;
  font-size: 14px;
  color: var(--mapex-text-primary);
  margin-bottom: 4px;
}

.step-item.active .step-title {
  color: var(--q-primary);
}

.step-item.completed .step-title {
  color: var(--q-primary);
}

.step-description {
  font-size: 12px;
  color: var(--mapex-text-secondary);
  line-height: 1.4;
}

.current-step-info {
  background-color: var(--mapex-surface-elevated);
  padding: 12px;
  border-radius: var(--mapex-radius-sm);
  border-left: 3px solid var(--q-primary);
}
</style>
