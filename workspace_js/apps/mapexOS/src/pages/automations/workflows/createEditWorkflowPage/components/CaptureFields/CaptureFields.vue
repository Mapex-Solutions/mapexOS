<script setup lang="ts">
/** TYPE IMPORTS */
import type { CaptureField } from '../../interfaces/CreateEditWorkflow.interface';

/** VUE IMPORTS */
import { ref, computed } from 'vue';

/** COMPOSABLES */
import { useWorkflowEditorState } from '../../composables';
import { useCreateEditWorkflowTranslations } from '@composables/i18n/pages/automations/workflows';

/** LOCAL IMPORTS (constants and handlers ONLY) */
import { VARIABLE_TYPE_OPTIONS } from '../../constants';

/** COMPOSABLES & STORES */
const {
  captureFields,
  addCaptureField,
  updateCaptureField,
  removeCaptureField,
  moveCaptureField,
} = useWorkflowEditorState();
const t = useCreateEditWorkflowTranslations();

/** STATE */

/**
 * Currently editing field index (-1 = adding new)
 */
const editingIndex = ref(-1);

/**
 * Form data for add/edit
 */
const form = ref<CaptureField>({
  field: '',
  type: 'string',
  description: '',
});

/**
 * Whether we are in edit mode
 */
const isEditing = computed(() => editingIndex.value >= 0);

/** FUNCTIONS */

/**
 * Reset form to defaults
 *
 * @returns {void}
 */
function resetForm(): void {
  form.value = { field: '', type: 'string', description: '' };
  editingIndex.value = -1;
}

/**
 * Start editing a capture field
 *
 * @param {number} index - Field index to edit
 * @returns {void}
 */
function startEdit(index: number): void {
  editingIndex.value = index;
  form.value = { ...captureFields.value[index]! };
}

/**
 * Handle form submission (add or update)
 *
 * @returns {void}
 */
function handleSubmit(): void {
  if (!form.value.field.trim()) return;

  if (isEditing.value) {
    updateCaptureField(editingIndex.value, { ...form.value });
  } else {
    addCaptureField({ ...form.value });
  }

  resetForm();
}

/**
 * Handle delete
 *
 * @param {number} index - Field index to delete
 * @returns {void}
 */
function handleDelete(index: number): void {
  removeCaptureField(index);
  if (editingIndex.value === index) {
    resetForm();
  }
}
</script>

<template>
  <div class="row q-col-gutter-md">
    <!-- Sidebar: Add/Edit Form -->
    <div class="col-12 col-md-4">
      <q-card flat bordered class="sticky-sidebar">
        <q-card-section>
          <div class="text-subtitle2 text-weight-medium q-mb-md">
            {{ isEditing ? t.captureFields.editTitle.value : t.captureFields.addTitle.value }}
          </div>

          <q-input
            v-model="form.field"
            :label="t.captureFields.fieldName.value"
            outlined
            dense
            class="q-mb-md"
            :rules="[(val: string) => !!val || t.validation.fieldNameRequired.value]"
          />

          <q-select
            v-model="form.type"
            :label="t.captureFields.type.value"
            outlined
            dense
            :options="[...VARIABLE_TYPE_OPTIONS]"
            emit-value
            map-options
            class="q-mb-md"
          />

          <q-input
            v-model="form.description"
            :label="t.captureFields.description.value"
            outlined
            dense
            type="textarea"
            autogrow
            class="q-mb-md"
            :hint="t.captureFields.descriptionHint.value"
          />

          <div class="row q-gutter-sm">
            <q-btn
              unelevated
              no-caps
              :color="isEditing ? 'amber-8' : 'primary'"
              :label="isEditing ? t.captureFields.update.value : t.captureFields.add.value"
              :icon="isEditing ? 'save' : 'add'"
              :disable="!form.field.trim()"
              @click="handleSubmit"
            />
            <q-btn
              v-if="isEditing"
              flat
              no-caps
              color="grey-7"
              :label="t.captureFields.cancel.value"
              @click="resetForm"
            />
          </div>
        </q-card-section>
      </q-card>
    </div>

    <!-- Capture Fields List -->
    <div class="col-12 col-md-8">
      <!-- Empty state -->
      <div v-if="captureFields.length === 0" class="empty-state">
        <q-icon name="analytics" size="48px" class="q-mb-md" />
        <p>{{ t.captureFields.emptyTitle.value }}</p>
        <p class="text-caption">{{ t.captureFields.emptyDescription.value }}</p>
      </div>

      <!-- List -->
      <q-list v-else separator bordered class="rounded-borders">
        <q-item v-for="(field, index) in captureFields" :key="index">
          <q-item-section>
            <q-item-label class="text-weight-medium">
              {{ field.field }}
            </q-item-label>
            <q-item-label caption>
              <q-badge :label="field.type" color="grey-7" class="q-mr-xs" />
              {{ field.description }}
            </q-item-label>
          </q-item-section>
          <q-item-section side>
            <div class="row q-gutter-xs">
              <q-btn flat dense round icon="arrow_upward" size="sm" :disable="index === 0" @click="moveCaptureField(index, 'up')" />
              <q-btn flat dense round icon="arrow_downward" size="sm" :disable="index === captureFields.length - 1" @click="moveCaptureField(index, 'down')" />
              <q-btn flat dense round icon="edit" size="sm" color="primary" @click="startEdit(index)" />
              <q-btn flat dense round icon="delete" size="sm" color="negative" @click="handleDelete(index)" />
            </div>
          </q-item-section>
        </q-item>
      </q-list>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.sticky-sidebar {
  position: sticky;
  top: 80px;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 300px;
  height: 100%;
  text-align: center;
  color: var(--mapex-text-secondary);
}
</style>
