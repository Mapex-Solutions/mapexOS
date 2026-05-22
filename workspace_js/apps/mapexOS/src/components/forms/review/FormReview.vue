<script setup lang="ts">
defineOptions({
  name: 'FormReview'
});

/** TYPE IMPORTS */
import type { FormReviewProps, FormReviewEmits, ReviewFieldDef } from './interfaces';

/** COMPONENTS */
import { DetailChip } from '@components/chips';

/** VUE IMPORTS */

/** PROPS & EMITS */
withDefaults(defineProps<FormReviewProps>(), {
  editMode: true,
  showSuccessBanner: true,
  successMessage: 'Ready to save! Click "Save" to confirm your changes.',
  description: 'Review your configuration before saving. You can edit any section by clicking the edit button.',
});

const emit = defineEmits<FormReviewEmits>();

/** FUNCTIONS */

/**
 * Format date and time based on specified format
 * @param {string | number | Date} value - Date value to format
 * @param {string} format - Format type: 'date', 'time', or 'datetime'
 * @returns {string} Formatted date string
 */
function formatDateTime(value: string | number | Date, format?: string): string {
  if (!value) return '—';

  try {
    const date = new Date(value);
    if (isNaN(date.getTime())) return String(value);

    if (format === 'date') {
      return date.toLocaleDateString('en-US', { year: 'numeric', month: 'short', day: 'numeric' });
    } else if (format === 'time') {
      return date.toLocaleTimeString('en-US', { hour: '2-digit', minute: '2-digit' });
    } else {
      return date.toLocaleString('en-US', { year: 'numeric', month: 'short', day: 'numeric', hour: '2-digit', minute: '2-digit' });
    }
  } catch {
    return String(value);
  }
}

/**
 * Convert value to string safely, handling objects
 * @param {unknown} value - Value to convert
 * @returns {string} String representation
 */
function safeStringify(value: unknown): string {
  if (value === null || value === undefined) return '';
  if (typeof value === 'string') return value;
  if (typeof value === 'number') return value.toString();
  if (typeof value === 'boolean') return value.toString();
  if (typeof value === 'bigint') return value.toString();
  if (typeof value === 'symbol') return value.toString();
  if (typeof value === 'function') return '[Function]';
  if (typeof value === 'object') {
    try {
      return JSON.stringify(value);
    } catch {
      return '[Object]';
    }
  }
  return '[Unknown]';
}

/**
 * Get badge/chip color based on field configuration
 * @param {ReviewFieldDef} field - Field definition
 * @returns {string} Color string
 */
function getColor(field: ReviewFieldDef): string {
  if (field.badgeColors) {
    if (typeof field.badgeColors === 'string') {
      return field.badgeColors;
    } else if (typeof field.badgeColors === 'object' && field.value !== null && field.value !== undefined) {
      const key = safeStringify(field.value);
      return field.badgeColors[key] || 'primary';
    }
  }
  return 'primary';
}

/**
 * Handle edit button click
 * @param {number} stepNumber - Step number to navigate to
 */
function handleEditSection(stepNumber: number): void {
  emit('editSection', stepNumber);
}

/**
 * Get display value for text fields
 * @param {ReviewFieldDef} field - Field definition
 * @returns {string} Display value
 */
function getTextValue(field: ReviewFieldDef): string {
  if (field.value === null || field.value === undefined || field.value === '') {
    return '—';
  }
  return safeStringify(field.value);
}

/**
 * Get boolean display value
 * @param {unknown} value - Boolean value
 * @returns {string} 'Active' or 'Inactive'
 */
function getBooleanLabel(value: unknown): string {
  return value ? 'Active' : 'Inactive';
}

/**
 * Get boolean badge color
 * @param {unknown} value - Boolean value
 * @returns {string} Color string
 */
function getBooleanColor(value: unknown): string {
  return value ? 'green-6' : 'grey-6';
}

/**
 * Format JSON for display
 * @param {unknown} value - JSON value
 * @returns {string} Formatted JSON string
 */
function formatJson(value: unknown): string {
  if (value === null || value === undefined) return '—';
  try {
    return JSON.stringify(value, null, 2);
  } catch {
    return safeStringify(value);
  }
}
</script>

<template>
  <div class="form-review">
    <!-- Description -->
    <div v-if="description" class="text-body1 text-grey-8 q-mb-lg">
      {{ description }}
    </div>

    <!-- Sections -->
    <template v-for="(section, index) in sections" :key="index">
      <q-card flat bordered class="q-mb-md" :data-testid="section.testId">
        <q-card-section>
          <!-- Section Header -->
          <div class="row items-center q-mb-md">
            <div class="col">
              <div class="row items-center">
                <q-icon
                  :name="section.icon.name"
                  size="sm"
                  :color="section.icon.color || 'primary'"
                  class="q-mr-sm"
                />
                <span class="text-h6 text-weight-medium text-dark">{{ section.label }}</span>
              </div>
            </div>
            <div v-if="editMode" class="col-auto">
              <q-btn
                flat
                dense
                size="sm"
                icon="edit"
                color="primary"
                label="Edit"
                :data-testid="section.testId ? `${section.testId}-edit-btn` : undefined"
                @click="handleEditSection(section.stepNumber)"
              />
            </div>
          </div>

          <!-- Section Fields -->
          <div class="row q-col-gutter-md">
            <template v-for="(field, fieldIndex) in section.fields" :key="fieldIndex">
              <div :class="`col-12 col-md-${field.colSize || 6}`">
                <div class="review-field">
                  <div class="field-label">{{ field.label }}</div>

                  <!-- Text Field -->
                  <div v-if="field.type === 'text'" class="field-value">
                    {{ getTextValue(field) }}
                  </div>

                  <!-- Badge Field -->
                  <div v-else-if="field.type === 'badge'" class="field-value">
                    <q-badge :color="getColor(field)">
                      {{ field.value ?? '—' }}
                    </q-badge>
                  </div>

                  <!-- Chip Field -->
                  <div v-else-if="field.type === 'chip'" class="field-value">
                    <DetailChip
                      :value="String(field.value ?? '—')"
                      :icon="field.icon"
                      :color="getColor(field) as any"
                      size="sm"
                    />
                  </div>

                  <!-- Boolean Field -->
                  <div v-else-if="field.type === 'boolean'" class="field-value">
                    <q-badge :color="getBooleanColor(field.value)">
                      {{ getBooleanLabel(field.value) }}
                    </q-badge>
                  </div>

                  <!-- DateTime Field -->
                  <div v-else-if="field.type === 'datetime'" class="field-value">
                    <q-icon name="schedule" size="xs" class="q-mr-xs text-grey-6" />
                    {{ formatDateTime(field.value as string | number | Date, field.format) }}
                  </div>

                  <!-- JSON Field -->
                  <div v-else-if="field.type === 'json'" class="field-value">
                    <div class="json-preview">
                      <pre class="json-content">{{ formatJson(field.value) }}</pre>
                    </div>
                  </div>

                  <!-- Fallback -->
                  <div v-else class="field-value">
                    {{ getTextValue(field) }}
                  </div>
                </div>
              </div>
            </template>
          </div>
        </q-card-section>
      </q-card>
    </template>

    <!-- Empty State -->
    <div v-if="!sections.length" class="empty-state text-center q-pa-xl">
      <q-icon name="info" size="64px" color="grey-5" class="q-mb-md" />
      <div class="text-h6 text-grey-6 q-mb-sm">No sections to review</div>
      <div class="text-body2 text-grey-5">
        Complete the previous steps to view data here.
      </div>
    </div>

    <!-- Success Banner -->
    <q-banner v-if="showSuccessBanner && sections.length" rounded class="bg-green-1 text-green-9 q-mt-lg">
      <template #avatar>
        <q-icon name="check_circle" color="green-7" />
      </template>
      <div class="text-body2">
        <strong>{{ successMessage }}</strong>
      </div>
    </q-banner>
  </div>
</template>

<style lang="scss" scoped>
.form-review {
  .review-field {
    .field-label {
      font-size: 0.75rem;
      font-weight: 500;
      color: var(--mapex-text-secondary);
      text-transform: uppercase;
      margin-bottom: 4px;
      letter-spacing: 0.5px;
    }

    .field-value {
      font-size: 0.95rem;
      color: var(--mapex-text-primary);
      min-height: 24px;
    }
  }

  .json-preview {
    background-color: var(--mapex-surface-elevated);
    border-radius: var(--mapex-radius-xs);
    padding: 12px;
    border: 1px solid var(--mapex-card-border);
    max-height: 200px;
    overflow: auto;

    .json-content {
      margin: 0;
      font-family: 'Courier New', monospace;
      font-size: 0.85rem;
      color: var(--mapex-text-primary);
      white-space: pre-wrap;
      word-wrap: break-word;
    }
  }

  .empty-state {
    background: var(--mapex-surface-elevated);
    border-radius: var(--mapex-radius-lg);
    border: 2px dashed var(--mapex-card-border);
  }
}
</style>
