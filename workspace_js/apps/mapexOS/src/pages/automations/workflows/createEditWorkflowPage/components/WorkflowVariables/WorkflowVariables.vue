<script setup lang="ts">
/** TYPE IMPORTS */
import type { VariableForm } from './interfaces/WorkflowVariables.interface';

/** VUE IMPORTS */
import { ref, computed } from 'vue';

/** COMPONENTS */
import { AppTooltip } from '@components/tooltips';

/** COMPOSABLES */
import { useWorkflowEditorState } from '../../composables';
import { useCreateEditWorkflowTranslations } from '@composables/i18n/pages/automations/workflows';

/** LOCAL IMPORTS (constants and handlers ONLY) */
import { VARIABLE_TYPE_OPTIONS, DEFAULT_VALUE_BY_TYPE } from '../../constants';

/** COMPOSABLES & STORES */
const {
  states,
  addState,
  updateState,
  removeState,
  moveState,
} = useWorkflowEditorState();
const t = useCreateEditWorkflowTranslations();

/** STATE */

/**
 * Currently editing variable index (-1 = adding new)
 */
const editingIndex = ref(-1);

/**
 * Form data for add/edit
 */
const form = ref<VariableForm>({
  field: '',
  type: 'string',
  defaultValue: '',
  durable: false,
});

/**
 * Whether we are in edit mode (vs add mode)
 */
const isEditing = computed(() => editingIndex.value >= 0);

/** FUNCTIONS */

/**
 * Reset form to defaults
 *
 * @returns {void}
 */
function resetForm(): void {
  form.value = {
    field: '',
    type: 'string',
    defaultValue: '',
    durable: false,
  };
  editingIndex.value = -1;
}

/**
 * Start editing a variable
 *
 * @param {number} index - Variable index to edit
 * @returns {void}
 */
function startEdit(index: number): void {
  editingIndex.value = index;
  form.value = { ...states.value[index] } as VariableForm;
}

/**
 * Handle form submission (add or update)
 *
 * @returns {void}
 */
function handleSubmit(): void {
  if (!form.value.field.trim()) return;

  if (isEditing.value) {
    updateState(editingIndex.value, { ...form.value });
  } else {
    addState({ ...form.value });
  }

  resetForm();
}

/**
 * Handle type change — reset default value to match new type
 *
 * @param {string} newType - New variable type
 * @returns {void}
 */
function handleTypeChange(newType: string): void {
  form.value.defaultValue = DEFAULT_VALUE_BY_TYPE[newType] ?? '';
}

/**
 * Handle delete with confirmation
 *
 * @param {number} index - Variable index to delete
 * @returns {void}
 */
function handleDelete(index: number): void {
  removeState(index);
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
            {{ isEditing ? t.variables.editTitle.value : t.variables.addTitle.value }}
          </div>

          <q-input
            v-model="form.field"
            :label="t.variables.name.value"
            outlined
            dense
            class="q-mb-md"
            :hint="t.variables.nameHint.value"
            :rules="[(val: string) => !!val || t.validation.nameIsRequired.value]"
          />

          <q-select
            v-model="form.type"
            :label="t.variables.type.value"
            outlined
            dense
            :options="[...VARIABLE_TYPE_OPTIONS]"
            emit-value
            map-options
            class="q-mb-md"
            @update:model-value="handleTypeChange"
          />

          <q-input
            v-model="form.description"
            :label="t.variables.description.value"
            outlined
            dense
            class="q-mb-md"
            :hint="t.variables.descriptionHint.value"
          />

          <!-- Dynamic default value input based on type -->
          <q-input
            v-if="form.type === 'string'"
            v-model="form.defaultValue"
            :label="t.variables.defaultValue.value"
            outlined
            dense
            class="q-mb-md"
          />
          <q-input
            v-else-if="form.type === 'number'"
            v-model.number="form.defaultValue"
            :label="t.variables.defaultValue.value"
            outlined
            dense
            type="number"
            class="q-mb-md"
          />
          <q-toggle
            v-else-if="form.type === 'boolean'"
            v-model="form.defaultValue"
            :label="t.variables.defaultValue.value"
            class="q-mb-md"
          />
          <q-input
            v-else
            v-model="form.defaultValue"
            :label="t.variables.defaultValueJson.value"
            outlined
            dense
            type="textarea"
            autogrow
            class="q-mb-md"
          />

          <!-- Persistence (Ephemeral / Durable) -->
          <div class="q-mb-md">
            <div class="persistence-label">
              {{ t.variables.persistence.value }}
              <q-icon
                name="info"
                size="14px"
                color="grey-6"
                class="q-ml-xs cursor-pointer"
              >
                <AppTooltip max-width="300px">
                  <div class="q-mb-sm">
                    {{ t.variables.ephemeralTooltip.value }}
                  </div>
                  <div>
                    {{ t.variables.durableTooltip.value }}
                  </div>
                </AppTooltip>
              </q-icon>
            </div>
            <q-option-group
              v-model="form.durable"
              inline
              dense
              type="radio"
              :options="[
                { label: t.variables.ephemeralLabel.value, value: false },
                { label: t.variables.durableLabel.value, value: true },
              ]"
            />
          </div>

          <div class="row q-gutter-sm">
            <q-btn
              unelevated
              no-caps
              :color="isEditing ? 'amber-8' : 'primary'"
              :label="isEditing ? t.variables.update.value : t.variables.add.value"
              :icon="isEditing ? 'save' : 'add'"
              :disable="!form.field.trim()"
              @click="handleSubmit"
            />
            <q-btn
              v-if="isEditing"
              flat
              no-caps
              color="grey-7"
              :label="t.variables.cancel.value"
              @click="resetForm"
            />
          </div>
        </q-card-section>
      </q-card>
    </div>

    <!-- Variables List -->
    <div class="col-12 col-md-8">
      <!-- Empty state -->
      <div v-if="states.length === 0" class="empty-state">
        <q-icon name="data_object" size="48px" class="q-mb-md" />
        <p>{{ t.variables.emptyTitle.value }}</p>
        <p class="text-caption">{{ t.variables.emptyDescription.value }}</p>
      </div>

      <!-- List -->
      <q-list v-else separator bordered class="rounded-borders">
        <q-item v-for="(variable, index) in states" :key="index">
          <q-item-section>
            <q-item-label class="text-weight-medium">
              {{ variable.field }}
            </q-item-label>
            <q-item-label caption>
              <q-badge :label="variable.type" color="grey-7" class="q-mr-xs" />
              <q-badge
                v-if="variable.durable"
                label="durable"
                color="amber-8"
                class="q-mr-xs"
              />
              {{ t.variables.defaultLabel.value }} {{ variable.defaultValue }}
            </q-item-label>
            <q-item-label v-if="variable.description" caption class="q-mt-xs text-grey-6">
              {{ variable.description }}
            </q-item-label>
          </q-item-section>
          <q-item-section side>
            <div class="row q-gutter-xs">
              <q-btn flat dense round icon="arrow_upward" size="sm" :disable="index === 0" @click="moveState(index, 'up')" />
              <q-btn flat dense round icon="arrow_downward" size="sm" :disable="index === states.length - 1" @click="moveState(index, 'down')" />
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

.persistence-label {
  display: flex;
  align-items: center;
  font-size: 0.75rem;
  font-weight: 600;
  color: var(--mapex-text-secondary);
  margin-bottom: 6px;
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
