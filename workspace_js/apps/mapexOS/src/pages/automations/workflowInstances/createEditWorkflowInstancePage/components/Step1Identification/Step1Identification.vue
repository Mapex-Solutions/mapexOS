<script setup lang="ts">
defineOptions({ name: 'Step1Identification' });

/** TYPE IMPORTS */
import type { WorkflowInstanceFormData } from '../../interfaces';
import type { QForm } from 'quasar';

/** VUE IMPORTS */
import { ref, computed, reactive, watch } from 'vue';

/** COMPOSABLES */
import { useCreateEditWorkflowInstanceTranslations } from '@src/composables/i18n/pages/automations/workflowInstances/createEditWorkflowInstancePage/useCreateEditWorkflowInstanceTranslations';

/** PROPS & EMITS */
const props = defineProps<{
  modelValue: WorkflowInstanceFormData;
}>();
const emit = defineEmits<{
  (e: 'update:modelValue', value: Partial<WorkflowInstanceFormData>): void;
}>();

/** COMPOSABLES & STORES */
const t = useCreateEditWorkflowInstanceTranslations();

/** STATE */
const formRef = ref<QForm | null>(null);

const localData = reactive({
  name: props.modelValue.name || '',
  description: props.modelValue.description || '',
  enabled: props.modelValue.enabled ?? true,
  isTemplate: props.modelValue.isTemplate ?? false,
  uniqueExecution: props.modelValue.uniqueExecution ?? false,
  workflowUUID: props.modelValue.workflowUUID || '',
});

/** COMPUTED */
const statusOptions = computed(() => [
  { label: 'Active', value: true },
  { label: 'Inactive', value: false },
]);

/** WATCHERS */
watch(() => props.modelValue, (newVal) => {
  localData.name = newVal.name || '';
  localData.description = newVal.description || '';
  localData.enabled = newVal.enabled ?? true;
  localData.isTemplate = newVal.isTemplate ?? false;
  localData.uniqueExecution = newVal.uniqueExecution ?? false;
  localData.workflowUUID = newVal.workflowUUID || '';
}, { deep: true });

/** FUNCTIONS */

/**
 * Emit partial update to parent
 * @returns {void}
 */
function updateValue(): void {
  emit('update:modelValue', { ...localData });
}

/**
 * Handle uniqueExecution toggle — clear UUID when disabled
 * @param {boolean} val - New value
 * @returns {void}
 */
function handleUniqueChange(val: boolean): void {
  if (!val) {
    localData.workflowUUID = '';
  }
  updateValue();
}

/**
 * Generate a new UUID v4
 * @returns {void}
 */
function generateUUID(): void {
  localData.workflowUUID = crypto.randomUUID();
  updateValue();
}

defineExpose({ formRef });
</script>

<template>
  <q-form ref="formRef" greedy>
    <div class="row q-col-gutter-md">
      <!-- Name -->
      <div class="col-12 col-sm-6">
        <q-input
          v-model="localData.name"
          outlined
          dense
          class="rounded-borders"
          data-testid="instance-name-input"
          :label="t.fields.name.value + ' *'"
          :placeholder="t.fields.namePlaceholder.value"
          :rules="[(val: string) => !!val || t.fields.requiredField.value]"
          @update:model-value="updateValue"
        >
          <template #prepend>
            <q-icon name="badge" color="primary" />
          </template>
        </q-input>
      </div>

      <!-- Status -->
      <div class="col-12 col-sm-6">
        <q-select
          v-model="localData.enabled"
          outlined
          dense
          emit-value
          map-options
          class="rounded-borders"
          data-testid="instance-status-select"
          :label="t.fields.enabled.value + ' *'"
          :options="statusOptions"
          option-label="label"
          option-value="value"
          @update:model-value="updateValue"
        >
          <template #prepend>
            <q-icon name="toggle_on" color="primary" />
          </template>
        </q-select>
      </div>

      <!-- Description -->
      <div class="col-12">
        <q-input
          v-model="localData.description"
          outlined
          dense
          type="textarea"
          rows="3"
          class="rounded-borders"
          data-testid="instance-description-input"
          :label="t.fields.description.value"
          :placeholder="t.fields.descriptionPlaceholder.value"
          @update:model-value="updateValue"
        >
          <template #prepend>
            <q-icon name="notes" color="primary" />
          </template>
        </q-input>
      </div>

      <!-- Shared Template checkbox -->
      <div class="col-12 col-sm-6">
        <q-checkbox
          v-model="localData.isTemplate"
          dense
          class="q-mt-sm"
          data-testid="instance-template-checkbox"
          :label="t.fields.isTemplate.value"
          @update:model-value="updateValue"
        />
        <div class="text-caption text-secondary q-ml-lg q-pl-xs">
          {{ t.fields.isTemplateHint.value }}
        </div>
      </div>

      <!-- Unique Execution checkbox -->
      <div class="col-12 col-sm-6">
        <q-checkbox
          v-model="localData.uniqueExecution"
          dense
          class="q-mt-sm"
          data-testid="instance-unique-checkbox"
          :label="t.fields.uniqueExecution.value"
          @update:model-value="handleUniqueChange"
        />
        <div class="text-caption text-secondary q-ml-lg q-pl-xs">
          {{ t.fields.uniqueExecutionHint.value }}
        </div>
      </div>

      <!-- Workflow UUID (visible only when uniqueExecution is true) -->
      <div v-if="localData.uniqueExecution" class="col-12">
        <q-input
          v-model="localData.workflowUUID"
          outlined
          dense
          class="rounded-borders"
          data-testid="instance-uuid-input"
          :label="t.fields.workflowUUID.value + ' *'"
          :placeholder="t.fields.workflowUUIDPlaceholder.value"
          :hint="t.fields.workflowUUIDHint.value"
          :rules="[(val: string) => !!val || t.fields.requiredField.value]"
          @update:model-value="updateValue"
        >
          <template #prepend>
            <q-icon name="fingerprint" color="primary" />
          </template>
          <template #append>
            <q-btn
              flat
              dense
              no-caps
              size="sm"
              color="primary"
              :label="t.fields.generateUUID.value"
              icon="casino"
              @click="generateUUID"
            />
          </template>
        </q-input>
      </div>
    </div>
  </q-form>
</template>

<style scoped lang="scss">
.rounded-borders {
  border-radius: var(--mapex-radius-md);
}

.text-secondary {
  color: var(--mapex-text-secondary);
}
</style>
