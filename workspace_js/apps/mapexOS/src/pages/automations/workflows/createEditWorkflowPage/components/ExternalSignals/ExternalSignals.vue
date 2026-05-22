<script setup lang="ts">
/** TYPE IMPORTS */
import type { ExternalSignal } from '../../interfaces/CreateEditWorkflow.interface';

/** VUE IMPORTS */
import { ref, computed } from 'vue';

/** COMPOSABLES */
import { useWorkflowEditorState } from '../../composables';
import { useCreateEditWorkflowTranslations } from '@composables/i18n/pages/automations/workflows';

/** COMPOSABLES & STORES */
const {
  externalSignals,
  addExternalSignal,
  updateExternalSignal,
  removeExternalSignal,
  moveExternalSignal,
} = useWorkflowEditorState();
const t = useCreateEditWorkflowTranslations();

/** STATE */

/**
 * Currently editing signal index (-1 = adding new)
 */
const editingIndex = ref(-1);

/**
 * Form data for add/edit
 */
const form = ref<ExternalSignal>({
  name: '',
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
  form.value = { name: '', description: '' };
  editingIndex.value = -1;
}

/**
 * Start editing a signal
 *
 * @param {number} index - Signal index to edit
 * @returns {void}
 */
function startEdit(index: number): void {
  editingIndex.value = index;
  form.value = { ...externalSignals.value[index]! };
}

/**
 * Handle form submission (add or update)
 *
 * @returns {void}
 */
function handleSubmit(): void {
  if (!form.value.name.trim()) return;

  if (isEditing.value) {
    updateExternalSignal(editingIndex.value, { ...form.value });
  } else {
    addExternalSignal({ ...form.value });
  }

  resetForm();
}

/**
 * Handle delete
 *
 * @param {number} index - Signal index to delete
 * @returns {void}
 */
function handleDelete(index: number): void {
  removeExternalSignal(index);
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
            {{ isEditing ? t.externalSignals.editTitle.value : t.externalSignals.addTitle.value }}
          </div>

          <q-input
            v-model="form.name"
            :label="t.externalSignals.name.value"
            outlined
            dense
            class="q-mb-md"
            :rules="[(val: string) => !!val || t.validation.fieldNameRequired.value]"
          />

          <q-input
            v-model="form.description"
            :label="t.externalSignals.description.value"
            outlined
            dense
            type="textarea"
            autogrow
            class="q-mb-md"
            :hint="t.externalSignals.descriptionHint.value"
          />

          <div class="row q-gutter-sm">
            <q-btn
              unelevated
              no-caps
              :color="isEditing ? 'amber-8' : 'primary'"
              :label="isEditing ? t.externalSignals.update.value : t.externalSignals.add.value"
              :icon="isEditing ? 'save' : 'add'"
              :disable="!form.name.trim()"
              @click="handleSubmit"
            />
            <q-btn
              v-if="isEditing"
              flat
              no-caps
              color="grey-7"
              :label="t.externalSignals.cancel.value"
              @click="resetForm"
            />
          </div>
        </q-card-section>
      </q-card>
    </div>

    <!-- Signals List -->
    <div class="col-12 col-md-8">
      <!-- Empty state -->
      <div v-if="externalSignals.length === 0" class="empty-state">
        <q-icon name="sensors" size="48px" class="q-mb-md" />
        <p>{{ t.externalSignals.emptyTitle.value }}</p>
        <p class="text-caption">{{ t.externalSignals.emptyDescription.value }}</p>
      </div>

      <!-- List -->
      <q-list v-else separator bordered class="rounded-borders">
        <q-item v-for="(signal, index) in externalSignals" :key="index">
          <q-item-section>
            <q-item-label class="text-weight-medium">
              {{ signal.name }}
            </q-item-label>
            <q-item-label caption>
              {{ signal.description }}
            </q-item-label>
          </q-item-section>
          <q-item-section side>
            <div class="row q-gutter-xs">
              <q-btn flat dense round icon="arrow_upward" size="sm" :disable="index === 0" @click="moveExternalSignal(index, 'up')" />
              <q-btn flat dense round icon="arrow_downward" size="sm" :disable="index === externalSignals.length - 1" @click="moveExternalSignal(index, 'down')" />
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
