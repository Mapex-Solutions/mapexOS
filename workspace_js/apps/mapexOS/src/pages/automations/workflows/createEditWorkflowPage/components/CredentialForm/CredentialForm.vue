<script setup lang="ts">
/** TYPE IMPORTS */
import type { CredentialFieldDefinition } from '@components/workflow/interfaces';
import type { CredentialFormProps, CredentialFormEmits } from './interfaces/CredentialForm.interface';

/** VUE IMPORTS */
import { ref, computed } from 'vue';

/** COMPOSABLES */
import { useCreateEditWorkflowTranslations } from '@composables/i18n/pages/automations/workflows';

/** PROPS & EMITS */
const props = defineProps<CredentialFormProps>();
const emit = defineEmits<CredentialFormEmits>();

/** COMPOSABLES & STORES */
const t = useCreateEditWorkflowTranslations();

/** STATE */

/**
 * Credential name input
 */
const name = ref(props.initialName);

/**
 * Form data keyed by field name
 */
const formData = ref<Record<string, unknown>>(buildInitialData());

/**
 * Password visibility toggle per secret field
 */
const secretVisible = ref<Record<string, boolean>>({});

/** COMPUTED */

/**
 * Whether the form is valid (name is required, required fields filled)
 */
const isValid = computed(() => {
  if (!name.value.trim()) return false;
  for (const field of props.credentialDef.fields) {
    if (field.required && !formData.value[field.name]) return false;
  }
  return true;
});

/** FUNCTIONS */

/**
 * Build initial form data from credential field defaults
 *
 * @returns {Record<string, unknown>} Initial form values
 */
function buildInitialData(): Record<string, unknown> {
  const data: Record<string, unknown> = {};
  for (const field of props.credentialDef.fields) {
    data[field.name] = field.default ?? '';
  }
  return data;
}

/**
 * Get the input type for a credential field
 *
 * @param {CredentialFieldDefinition} field - Field definition
 * @returns {'text' | 'password' | 'number'} HTML input type compatible with q-input
 */
function getInputType(field: CredentialFieldDefinition): 'text' | 'password' | 'number' {
  if (field.isSecret && !secretVisible.value[field.name]) return 'password';
  if (field.type === 'number') return 'number';
  return 'text';
}

/**
 * Toggle password visibility for a secret field
 *
 * @param {string} fieldName - Field name to toggle
 */
function toggleSecret(fieldName: string): void {
  secretVisible.value[fieldName] = !secretVisible.value[fieldName];
}

/**
 * Handle form submission
 */
function handleSave(): void {
  if (!isValid.value) return;

  // Only send fields that have a non-empty value (skip placeholder secrets on edit)
  const data: Record<string, unknown> = {};
  for (const field of props.credentialDef.fields) {
    const val = formData.value[field.name];
    if (val !== '' && val !== undefined) {
      data[field.name] = val;
    }
  }

  emit('save', { name: name.value.trim(), data });
}
</script>

<template>
  <div class="credential-form">
    <!-- Credential Name -->
    <q-input
      v-model="name"
      outlined
      dense
      class="q-mb-md"
      :label="t.credentials.credentialName.value"
      :hint="t.credentials.credentialNameHint.value"
    />

    <!-- Dynamic Fields from credentialDef.fields -->
    <template v-for="field in credentialDef.fields" :key="field.name">
      <!-- Options type → q-select -->
      <q-select
        v-if="field.type === 'options' && field.options"
        v-model="formData[field.name]"
        outlined
        dense
        emit-value
        map-options
        class="q-mb-md"
        :options="field.options"
        :label="field.displayName"
        :hint="field.hint"
      />

      <!-- Boolean type → q-toggle -->
      <q-toggle
        v-else-if="field.type === 'boolean'"
        v-model="formData[field.name]"
        class="q-mb-md"
        :label="field.displayName"
      />

      <!-- String / Number / Secret → q-input -->
      <q-input
        v-else
        :model-value="String(formData[field.name] ?? '')"
        outlined
        dense
        class="q-mb-md"
        :type="getInputType(field)"
        :label="field.displayName + (field.required ? ' *' : '')"
        :hint="field.hint"
        :placeholder="props.isEdit && field.isSecret ? t.credentials.secretPlaceholder.value : undefined"
        @update:model-value="formData[field.name] = field.type === 'number' ? Number($event) : $event"
      >
        <template v-if="field.isSecret" #append>
          <q-icon
            :name="secretVisible[field.name] ? 'visibility_off' : 'visibility'"
            class="cursor-pointer"
            @click="toggleSecret(field.name)"
          />
        </template>
      </q-input>
    </template>

    <!-- Actions -->
    <div class="row justify-end q-gutter-sm q-mt-md">
      <q-btn
        flat
        no-caps
        color="grey-7"
        :label="t.credentials.cancel.value"
        @click="emit('cancel')"
      />
      <q-btn
        unelevated
        no-caps
        color="primary"
        :label="t.credentials.save.value"
        :disable="!isValid"
        @click="handleSave"
      />
    </div>
  </div>
</template>
