<script setup lang="ts">
defineOptions({
  name: 'Step3BasicInfo'
});

/** TYPE IMPORTS */
import type { Trigger } from '../../interfaces';
import type { QForm } from 'quasar';

/** VUE IMPORTS */
import { ref, computed } from 'vue';

/** COMPONENTS */

/** COMPOSABLES */
import { useCreateEditTriggerTranslations } from '@composables/i18n/pages/automations/triggers/createEditTrigger';

/** UTILS */

/** SERVICES */

/** STORES */
import { useOrganizationStore } from '@stores/organization';

/** LOCAL IMPORTS */

/** PROPS & EMITS */
const props = defineProps<{
  modelValue: Trigger;
}>();

const emit = defineEmits<{
  'update:modelValue': [value: Trigger];
}>();

/** COMPOSABLES & STORES */
const translations = useCreateEditTriggerTranslations();
const orgStore = useOrganizationStore();

/** STATE */

/**
 * Form reference for validation
 */
const formRef = ref<QForm | null>(null);

/**
 * Local copies of modelValue fields to avoid prop mutation
 */
const localName = ref(props.modelValue.name);
const localDescription = ref(props.modelValue.description);
const localEnabled = ref(props.modelValue.enabled);
const localIsTemplate = ref(props.modelValue.isTemplate ?? false);

/** COMPUTED */

/**
 * Check if user can create shared templates
 * Only Vendor and Customer organizations can create shared templates
 */
const canCreateTemplate = ref(orgStore.isVendor || orgStore.isCustomer);

/**
 * Translated options for the status select
 */
const statusOptions = computed(() => [
  { label: translations.statusOptions.active.label.value, value: true },
  { label: translations.statusOptions.inactive.label.value, value: false },
]);

/** WATCHERS */

/** FUNCTIONS */

/**
 * Update trigger with new values
 * @param {Partial<Trigger>} updates - Partial updates to apply
 * @returns {void}
 */
function updateTrigger(updates: Partial<Trigger>): void {
  const updated: Trigger = { ...props.modelValue, ...updates };
  emit('update:modelValue', updated);
}

/**
 * Validate form
 * @returns {Promise<boolean>} Validation result
 */
async function validate(): Promise<boolean> {
  if (!formRef.value) return false;
  return await formRef.value.validate();
}

/** LIFECYCLE HOOKS */

/** EXPOSE */
defineExpose({
  formRef,
  validate,
});
</script>

<template>
  <div class="step3-basic-info">
    <div class="text-body1 text-grey-8 q-mb-lg">
      {{ translations.steps.step3.intro.value }}
    </div>

    <q-form ref="formRef">
      <div class="row q-col-gutter-md">
        <!-- Name (col-12 col-sm-9) -->
        <div class="col-12 col-sm-9">
          <q-input
            v-model="localName"
            outlined
            dense
            :label="translations.steps.step3.fields.nameLabel.value"
            :placeholder="translations.steps.step3.fields.namePlaceholder.value"
            :hint="translations.steps.step3.fields.nameHint.value"
            :rules="[
              (val: string) => !!val || translations.steps.step3.fields.nameRequired.value,
              (val: string) => val.length >= 3 || translations.steps.step3.fields.nameMinLength.value
            ]"
            @update:model-value="(val) => updateTrigger({ name: String(val || '') })"
          >
            <template v-slot:prepend>
              <q-icon name="label" />
            </template>
          </q-input>
        </div>

        <!-- Status (col-12 col-sm-3) -->
        <div class="col-12 col-sm-3">
          <q-select
            v-model="localEnabled"
            outlined
            dense
            emit-value
            map-options
            :label="translations.steps.step3.fields.statusLabel.value"
            :options="statusOptions"
            @update:model-value="(val) => updateTrigger({ enabled: val })"
          >
            <template v-slot:prepend>
              <q-icon name="toggle_on" />
            </template>
          </q-select>
        </div>

        <!-- Description (col-12) -->
        <div class="col-12">
          <q-input
            v-model="localDescription"
            outlined
            dense
            type="textarea"
            :label="translations.steps.step3.fields.descriptionLabel.value"
            :placeholder="translations.steps.step3.fields.descriptionPlaceholder.value"
            :hint="translations.steps.step3.fields.descriptionHint.value"
            rows="3"
            @update:model-value="(val) => {
              const updated: Trigger = { ...modelValue };
              if (val) {
                updated.description = String(val);
              } else {
                delete updated.description;
              }
              updateTrigger(updated);
            }"
          >
            <template v-slot:prepend>
              <q-icon name="description" />
            </template>
          </q-input>
        </div>

        <!-- isTemplate checkbox - Only for Vendor/Customer -->
        <div v-if="canCreateTemplate" class="col-12">
          <div class="q-py-sm">
            <q-checkbox
              v-model="localIsTemplate"
              color="primary"
              class="q-mb-xs"
              :label="translations.steps.step3.fields.sharedTemplateLabel.value"
              @update:model-value="(val) => updateTrigger({ isTemplate: val })"
            />
            <div class="text-caption text-grey-7 q-pl-lg">
              {{ translations.steps.step3.fields.sharedTemplateHint.value }}
            </div>
          </div>
        </div>
      </div>
    </q-form>

    <!-- Info box -->
    <q-banner rounded class="bg-blue-1 text-blue-9 q-mt-lg">
      <template v-slot:avatar>
        <q-icon name="info" color="blue-7" />
      </template>
      <div class="text-body2">
        <strong>{{ translations.steps.step3.tip.prefix.value }}</strong> {{ translations.steps.step3.tip.text.value }}
      </div>
    </q-banner>
  </div>
</template>

<style lang="scss" scoped>
.step3-basic-info {
  // Component-specific styles
}
</style>
