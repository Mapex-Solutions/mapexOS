<script setup lang="ts">
defineOptions({
  name: 'StepperVertical'
});

import type { StepperVerticalProps, StepperVerticalItem } from './interfaces';

// Define props with defaults
const props = withDefaults(defineProps<StepperVerticalProps>(), {
  currentStep: 1,
  title: 'Configuration Steps',
  subtitle: 'Complete all steps',
  headerIcon: 'timeline',
  infoText: 'All fields marked with * are required',
  currentStepLabel: 'Current Step',
  mode: 'creating',
  allowStepNavigation: false,
});

/**
 * Generates step attributes including optional ID
 * @param index Zero-based index of the step
 * @returns Object with step attributes
 */
const getStepAttrs = (index: number): Record<string, string> => {
  const attrs: Record<string, string> = {};
  if (props.stepIdPrefix) {
    attrs.id = `${props.stepIdPrefix}-${index + 1}`;
  }
  return attrs;
};

// Define emits
const emit = defineEmits<{
  (e: 'step-click', stepNumber: number): void
}>();

/**
 * Determines if a step should be shown as active
 * @param index Zero-based index of the step
 * @returns Whether the step should be shown as active
 */
const isActive = (index: number): boolean => {
  // If navigation is allowed (edit mode), all steps are active
  if (props.allowStepNavigation || props.mode === 'editing') {
    return true;
  }
  return props.currentStep === index + 1;
};

/**
 * Determines if a step should be shown as completed
 * @param index Zero-based index of the step
 * @returns Whether the step should be shown as completed
 */
const isCompleted = (index: number): boolean => {
  // If navigation is allowed (edit mode), all steps are completed
  if (props.allowStepNavigation || props.mode === 'editing') {
    return true;
  }
  return props.currentStep > index + 1;
};

/**
 * Gets the label of the current active step
 * @returns The label of the current step
 */
const getCurrentStepLabel = (): string => {
  const step = props.steps[props.currentStep - 1];
  return step ? step.title : '';
};

/**
 * Determines the icon to display for a step
 * @param step The step item
 * @param index Zero-based index of the step
 * @returns Icon name to display
 */
const getStepIcon = (step: StepperVerticalItem, index: number): string => {
  // If navigation is allowed (edit mode), always use the step's icon
  if (props.allowStepNavigation || props.mode === 'editing') {
    return step.icon;
  }

  // In creating mode, use check_circle for completed steps
  if (props.currentStep > index + 1) {
    return 'check_circle';
  }
  return step.icon;
};

/**
 * Returns the CSS state class for a step icon.
 * All visual styling is handled via CSS classes + CSS variables,
 * so dark mode works automatically without JS changes.
 *
 * @param {number} index - Zero-based index of the step
 * @returns {string} CSS class name for the step icon state
 */
const getStepIconStateClass = (index: number): string => {
  if (props.allowStepNavigation || props.mode === 'editing') {
    // Edit mode: current step = active, others = completed (muted)
    return props.currentStep === index + 1
      ? 'step-icon--active'
      : 'step-icon--completed';
  }

  // Creating mode
  if (props.currentStep === index + 1) return 'step-icon--active';
  if (props.currentStep > index + 1) return 'step-icon--completed';
  return 'step-icon--pending';
};

/**
 * Handles click on a step item
 * @param stepNumber 1-based index of the clicked step
 */
const handleStepClick = (stepNumber: number): void => {
  // Allow navigation if:
  // 1. allowStepNavigation is true (can click any step)
  // 2. mode is 'editing' (backward compatibility)
  // 3. stepNumber < currentStep (can always go back to previous steps)
  if (props.allowStepNavigation || props.mode === 'editing' || stepNumber < props.currentStep) {
    emit('step-click', stepNumber);
  }
};
</script>

<template>
  <q-card class="rounded-borders shadow-2" :class="props?.fullHeight ? 'stepper-card' : ''">
    <q-card-section class="bg-primary text-white q-pb-md">
      <div class="text-h6 text-weight-bold q-mb-sm">
        <q-icon size="sm" class="q-mr-xs" :name="headerIcon" />
        {{ props.title }}
      </div>
      <div class="text-caption">{{ props.subtitle }}</div>
    </q-card-section>

    <q-card-section class="q-pa-md">
      <div class="progress-steps">
        <div
            v-for="(st, idx) in steps"
            :key="idx"
            v-bind="getStepAttrs(idx)"
            class="step-item"
            :class="{
              active: isActive(idx),
              completed: isCompleted(idx),
              clickable: props.allowStepNavigation || props.mode === 'editing' || idx < props.currentStep - 1
            }"
            @click="handleStepClick(idx + 1)"
        >
          <div class="step-icon-wrapper">
            <div class="step-icon" :class="getStepIconStateClass(idx)">
              <q-icon
                  size="sm"
                  :name="getStepIcon(st, idx)"
              />
            </div>
          </div>
          <div class="step-content">
            <div class="step-title">{{ st.title }}</div>
            <div class="step-description">{{ st.description }}</div>
          </div>
        </div>
      </div>

      <q-separator class="q-my-md" />

      <div class="current-step-info">
        <div class="text-caption text-grey-6">
          <q-icon size="xs" class="q-mr-xs" name="info" />
          {{ infoText }}
        </div>
        <div class="text-caption text-primary q-mt-xs">
          <q-icon size="xs" class="q-mr-xs" name="arrow_forward" />
          {{ currentStepLabel }}: {{ getCurrentStepLabel() }}
        </div>
      </div>
    </q-card-section>
  </q-card>
</template>

<style scoped lang="scss">
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

.step-item.clickable {
  cursor: pointer;
}

.step-item.clickable:hover .step-title {
  text-decoration: underline;
}

// Connector line between steps
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

.step-icon-wrapper {
  position: relative;
  z-index: 1;
}

// ── Step icon base ─────────────────────────────────────────
.step-icon {
  width: 40px;
  height: 40px;
  border-radius: var(--mapex-radius-full);
  display: flex;
  align-items: center;
  justify-content: center;
  transition: var(--mapex-transition-slow);

  // Default (pending) - muted appearance
  background-color: var(--mapex-surface-elevated);
  border: 2px solid var(--mapex-card-border);
  color: var(--mapex-text-muted);
}

// ── Active step: solid primary, white icon ─────────────────
.step-icon--active {
  background-color: var(--q-primary);
  border-color: var(--q-primary);
  color: white;
}

// ── Completed step: subtle primary bg, primary icon ────────
.step-icon--completed {
  background-color: rgba(var(--mapex-primary-rgb), 0.15);
  border-color: rgba(var(--mapex-primary-rgb), 0.3);
  color: var(--mapex-primary);
}

// ── Pending step: inherits base (muted) ────────────────────
.step-icon--pending {
  // Uses the base .step-icon styles (muted)
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
