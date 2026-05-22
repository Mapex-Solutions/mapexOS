<script setup lang="ts">
defineOptions({
  name: 'Step1BasicInfo',
});

/** TYPE IMPORTS */
import type { Step1BasicInfoProps, Step1BasicInfoEmits } from './interfaces/Step1BasicInfo.interface';
import type { RouteGroupCreate } from '@interfaces/routing/routeGroups.interface';

/** VUE IMPORTS */
import { ref } from 'vue';
import { QForm } from 'quasar';

/** PROPS & EMITS */
const props = defineProps<Step1BasicInfoProps>();
const emit = defineEmits<Step1BasicInfoEmits>();

/** STATE */
const formRef = ref<QForm | null>(null);

/** FUNCTIONS */

/**
 * Validate the form
 *
 * @returns {Promise<boolean>} True if form is valid
 */
async function validate(): Promise<boolean> {
  if (formRef.value) {
    return await formRef.value.validate();
  }
  return false;
}

/**
 * Update a field in the form data
 *
 * @param {K} field - Field name to update
 * @param {RouteGroupCreate[K] | string | number | null} value - New value
 */
function updateField<K extends keyof RouteGroupCreate>(
  field: K,
  value: RouteGroupCreate[K] | string | number | null,
): void {
  if (value === null) return; // Ignore null values from Quasar
  emit('update:formData', {
    ...props.formData,
    [field]: value as RouteGroupCreate[K],
  });
}

/** EXPOSE */
defineExpose({
  validate,
});
</script>

<template>
  <q-form ref="formRef" greedy>
    <div class="q-mb-md">
      <div class="text-subtitle1 text-weight-medium q-mb-xs">
        <q-icon name="info" color="primary" class="q-mr-xs" />
        {{ t.createEdit.basicInfoStep.title.value }}
      </div>
      <div class="text-body2 text-grey-7">
        {{ t.createEdit.basicInfoStep.subtitle.value }}
      </div>
    </div>

    <div class="row q-col-gutter-md">
      <!-- Name -->
      <div id="route-group-field-name" class="col-xs-12 col-sm-6 col-md-6">
        <q-input
          outlined
          dense
          class="rounded-borders"
          :label="t.createEdit.basicInfoStep.fields.name.label.value + ' *'"
          :placeholder="t.createEdit.basicInfoStep.fields.name.placeholder.value"
          :hint="t.createEdit.basicInfoStep.fields.name.hint.value"
          :rules="[(val) => !!val || t.createEdit.basicInfoStep.fields.name.required.value]"
          :model-value="formData.name"
          @update:model-value="(val) => updateField('name', val)"
        >
          <template #prepend>
            <q-icon name="label" color="primary" />
          </template>
        </q-input>
      </div>

      <!-- Status -->
      <div id="route-group-field-status" class="col-xs-12 col-sm-6 col-md-6">
        <q-select
          outlined
          dense
          emit-value
          map-options
          class="rounded-borders"
          :label="t.createEdit.basicInfoStep.fields.enabled.label.value + ' *'"
          :options="statusOptions"
          :hint="t.createEdit.basicInfoStep.fields.enabled.hint.value"
          option-label="label"
          option-value="value"
          :model-value="formData.enabled"
          @update:model-value="(val) => updateField('enabled', val)"
        >
          <template #prepend>
            <q-icon name="toggle_on" color="primary" />
          </template>
        </q-select>
      </div>

      <!-- Description -->
      <div class="col-12">
        <q-input
          outlined
          dense
          class="rounded-borders"
          type="textarea"
          rows="3"
          :label="t.createEdit.basicInfoStep.fields.description.label.value"
          :placeholder="t.createEdit.basicInfoStep.fields.description.placeholder.value"
          :hint="t.createEdit.basicInfoStep.fields.description.hint.value"
          :model-value="formData.description"
          @update:model-value="(val) => updateField('description', val)"
        >
          <template #prepend>
            <q-icon name="notes" color="primary" />
          </template>
        </q-input>
      </div>

      <!-- isTemplate checkbox - Only for Vendor/Customer -->
      <div v-if="canCreateTemplate" class="col-12">
        <div class="q-py-sm">
          <q-checkbox
            :model-value="formData.isTemplate"
            @update:model-value="(val) => updateField('isTemplate', val)"
            :label="t.createEdit.basicInfoStep.fields.isTemplate.label.value"
            color="primary"
            class="q-mb-xs"
          />
          <div class="text-caption text-grey-7 q-pl-lg">
            {{ t.createEdit.basicInfoStep.fields.isTemplate.hint.value }}
          </div>
        </div>
      </div>
    </div>
  </q-form>
</template>

<style scoped lang="scss">
.rounded-borders {
  border-radius: var(--mapex-radius-md);
}
</style>
