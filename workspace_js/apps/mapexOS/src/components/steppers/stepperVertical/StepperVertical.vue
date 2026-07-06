<script setup lang="ts">
defineOptions({
  name: 'StepperVertical'
});

/** TYPE IMPORTS */
import type { StepperVerticalProps, StepperVerticalItem, StepperRenderRow } from './interfaces';

/** VUE IMPORTS */
import { computed } from 'vue';

/** PROPS & EMITS */
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

const emit = defineEmits<{
  (e: 'step-click', stepNumber: number): void
}>();

/** COMPUTED */

/**
 * Flattens the (possibly grouped) steps into ordered render rows. Leaves are
 * numbered 1..N in document order; a grouped item yields a header row followed
 * by its child leaf rows. With no children this is the flat list, unchanged.
 */
const renderRows = computed<StepperRenderRow[]>(() => {
  const rows: StepperRenderRow[] = [];
  let leaf = 0;
  for (const item of props.steps) {
    if (item.children && item.children.length) {
      const childNumbers: number[] = [];
      const childRows: StepperRenderRow[] = [];
      for (const child of item.children) {
        leaf += 1;
        childNumbers.push(leaf);
        childRows.push({ kind: 'leaf', item: child, leafNumber: leaf, indented: true, childLeafNumbers: [] });
      }
      rows.push({ kind: 'group', item, leafNumber: 0, indented: false, childLeafNumbers: childNumbers });
      rows.push(...childRows);
    } else {
      leaf += 1;
      rows.push({ kind: 'leaf', item, leafNumber: leaf, indented: false, childLeafNumbers: [] });
    }
  }
  return rows;
});

/** True when at least one step is a group (enables the grouped layout). */
const hasGroups = computed(() => props.steps.some((s) => !!s.children?.length));

/** FUNCTIONS */

/**
 * True in edit mode / free navigation, where every step is reachable.
 * @returns {boolean} Whether all steps are navigable.
 */
function isEditNav(): boolean {
  return props.allowStepNavigation || props.mode === 'editing';
}

/**
 * Builds optional per-leaf DOM attributes (id for tours/CSS selectors).
 * @param {number} leafNumber - The 1-based leaf number.
 * @returns {Record<string, string>} Attributes to bind.
 */
function getStepAttrs(leafNumber: number): Record<string, string> {
  const attrs: Record<string, string> = {};
  if (props.stepIdPrefix) {
    attrs.id = `${props.stepIdPrefix}-${leafNumber}`;
  }
  return attrs;
}

/**
 * Whether a leaf should render as active.
 * @param {number} n - The leaf number.
 * @returns {boolean} Active state.
 */
function isLeafActive(n: number): boolean {
  return isEditNav() ? true : props.currentStep === n;
}

/**
 * Whether a leaf should render as completed.
 * @param {number} n - The leaf number.
 * @returns {boolean} Completed state.
 */
function isLeafCompleted(n: number): boolean {
  return isEditNav() ? true : props.currentStep > n;
}

/**
 * Whether a leaf can be clicked to navigate.
 * @param {number} n - The leaf number.
 * @returns {boolean} Clickable state.
 */
function isLeafClickable(n: number): boolean {
  return props.allowStepNavigation || props.mode === 'editing' || n < props.currentStep;
}

/**
 * Icon for a leaf (check when completed in creating mode).
 * @param {StepperVerticalItem} step - The leaf item.
 * @param {number} n - The leaf number.
 * @returns {string} Icon name.
 */
function getLeafIcon(step: StepperVerticalItem, n: number): string {
  if (isEditNav()) {
    return step.icon;
  }
  return props.currentStep > n ? 'check_circle' : step.icon;
}

/**
 * CSS state class for a leaf icon.
 * @param {number} n - The leaf number.
 * @returns {string} State class.
 */
function getLeafIconStateClass(n: number): string {
  if (isEditNav()) {
    return props.currentStep === n ? 'step-icon--active' : 'step-icon--completed';
  }
  if (props.currentStep === n) return 'step-icon--active';
  if (props.currentStep > n) return 'step-icon--completed';
  return 'step-icon--pending';
}

/**
 * A group is active when the current leaf is one of its children.
 * @param {number[]} children - Child leaf numbers.
 * @returns {boolean} Active state.
 */
function isGroupActive(children: number[]): boolean {
  return children.includes(props.currentStep);
}

/**
 * A group is completed when the current leaf is past all of its children.
 * @param {number[]} children - Child leaf numbers.
 * @returns {boolean} Completed state.
 */
function isGroupCompleted(children: number[]): boolean {
  return children.length > 0 && props.currentStep > Math.max(...children);
}

/**
 * Title of the current active leaf.
 * @returns {string} The current step title.
 */
function getCurrentStepLabel(): string {
  const row = renderRows.value.find((r) => r.kind === 'leaf' && r.leafNumber === props.currentStep);
  return row ? row.item.title : '';
}

/**
 * Emits a navigation request for the clicked leaf when allowed.
 * @param {number} leafNumber - The 1-based leaf number.
 */
function handleStepClick(leafNumber: number): void {
  if (props.allowStepNavigation || props.mode === 'editing' || leafNumber < props.currentStep) {
    emit('step-click', leafNumber);
  }
}
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
      <div class="progress-steps" :class="{ 'progress-steps--grouped': hasGroups }">
        <template v-for="(row, idx) in renderRows" :key="idx">
          <!-- Group header (visual only, not navigable) -->
          <div
            v-if="row.kind === 'group'"
            class="step-group"
            :class="{
              'step-group--active': isGroupActive(row.childLeafNumbers),
              'step-group--completed': isGroupCompleted(row.childLeafNumbers),
            }"
          >
            <q-icon size="xs" :name="row.item.icon" class="step-group__icon" />
            <span class="step-group__title">{{ row.item.title }}</span>
          </div>

          <!-- Leaf step -->
          <div
            v-else
            v-bind="getStepAttrs(row.leafNumber)"
            class="step-item"
            :class="{
              active: isLeafActive(row.leafNumber),
              completed: isLeafCompleted(row.leafNumber),
              clickable: isLeafClickable(row.leafNumber),
              'step-item--indented': row.indented,
            }"
            @click="handleStepClick(row.leafNumber)"
          >
            <div class="step-icon-wrapper">
              <div class="step-icon" :class="getLeafIconStateClass(row.leafNumber)">
                <q-icon
                    size="sm"
                    :name="getLeafIcon(row.item, row.leafNumber)"
                />
              </div>
            </div>
            <div class="step-content">
              <div class="step-title">{{ row.item.title }}</div>
              <div class="step-description">{{ row.item.description }}</div>
            </div>
          </div>
        </template>
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

// Connector line between steps (flat mode only)
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

// Grouped layout uses indentation instead of connector lines
.progress-steps--grouped {
  gap: 10px;

  .step-item::after {
    display: none;
  }
}

// Group header — a lightweight, non-interactive context label
.step-group {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 4px 0 0;
  color: var(--mapex-text-secondary);
}

.step-group__icon {
  color: var(--mapex-text-muted);
}

.step-group__title {
  font-size: var(--mapex-font-xs);
  font-weight: var(--mapex-font-weight-bold);
  text-transform: uppercase;
  letter-spacing: 0.04em;
}

.step-group--active {
  .step-group__title,
  .step-group__icon {
    color: var(--mapex-primary);
  }
}

.step-item--indented {
  padding-left: 16px;
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
  transition: var(--mapex-transition-slow);

  background-color: var(--mapex-surface-elevated);
  border: 2px solid var(--mapex-card-border);
  color: var(--mapex-text-muted);
}

.step-icon--active {
  background-color: var(--q-primary);
  border-color: var(--q-primary);
  color: white;
}

.step-icon--completed {
  background-color: rgba(var(--mapex-primary-rgb), 0.15);
  border-color: rgba(var(--mapex-primary-rgb), 0.3);
  color: var(--mapex-primary);
}

.step-icon--pending {
  // Uses the base .step-icon styles (muted)
}

.step-content {
  flex: 1;
  padding-top: 2px;
}

.step-title {
  font-weight: var(--mapex-font-weight-semibold);
  font-size: var(--mapex-font-md);
  color: var(--mapex-text-primary);
  margin-bottom: 4px;
}

.step-description {
  font-size: var(--mapex-font-xs);
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
