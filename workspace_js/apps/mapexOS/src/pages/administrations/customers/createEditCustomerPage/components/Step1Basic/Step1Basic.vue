<template>
  <q-form ref="formRef" greedy>
    <div class="q-mb-md">
      <div class="text-subtitle1 text-weight-medium q-mb-xs">
        <q-icon name="business" color="primary" class="q-mr-xs" />
        {{ t.sections.basicInfo.value }}
      </div>
      <div class="text-body2 text-grey-7">
        {{ t.formDescriptions.basic.value }}
      </div>
    </div>

    <div class="row q-col-gutter-md">
      <!-- Name -->
      <div class="col-12 col-md-6">
        <q-input
          v-model="localData.name"
          outlined
          dense
          class="rounded-borders"
          data-testid="customer-name-input"
          :label="`${t.fields.name.value} *`"
          :rules="[
            (val) => !!val || t.validation.nameRequired.value,
            (val) => val.length >= NAME_MIN_LENGTH || t.validation.nameMinLength.value,
            (val) => val.length <= NAME_MAX_LENGTH || t.validation.nameMaxLength.value,
          ]"
          @update:model-value="updateValue"
        >
          <template #prepend>
            <q-icon name="business" color="primary" />
          </template>
        </q-input>
      </div>

      <!-- Status Select -->
      <div class="col-12 col-md-6">
        <q-select
          v-model="localData.enabled"
          outlined
          dense
          emit-value
          map-options
          class="rounded-borders"
          option-label="label"
          option-value="value"
          data-testid="customer-enabled-select"
          :label="`${t.fields.enabled.value} *`"
          :hint="t.hints.enabled.value"
          :options="statusOptions"
          @update:model-value="updateValue"
        >
          <template #prepend>
            <q-icon name="toggle_on" color="primary" />
          </template>
        </q-select>
      </div>

      <!-- Phone (only for types that support it) -->
      <div v-if="hasPhone" class="col-12">
        <q-input
          v-model="localData.phone"
          outlined
          dense
          type="tel"
          class="rounded-borders"
          data-testid="customer-phone-input"
          :label="t.fields.phone.value"
          :hint="t.hints.phoneFormat.value"
          @update:model-value="updateValue"
        >
          <template #prepend>
            <q-icon name="phone" color="primary" />
          </template>
        </q-input>
      </div>
    </div>
  </q-form>
</template>

<script setup lang="ts">
defineOptions({
  name: 'Step1Basic',
});

/** TYPE IMPORTS */
import type { Step1BasicProps } from './interfaces/Step1Basic.interface';
import type { QForm } from 'quasar';
import type { OrganizationFormData } from '../../interfaces';

/** VUE IMPORTS */
import { ref, reactive, computed, watch } from 'vue';

/** COMPOSABLES */
import { useAddCustomerTranslations } from '@composables/i18n';

/** LOCAL IMPORTS */
import { NAME_MIN_LENGTH, NAME_MAX_LENGTH } from '../../constants';

const props = withDefaults(defineProps<Step1BasicProps>(), {
  hasPhone: true,
});

const emit = defineEmits<{
  (e: 'update:modelValue', value: Partial<OrganizationFormData>): void;
}>();

/** COMPOSABLES & STORES */
const t = useAddCustomerTranslations();

/** STATE */
const formRef = ref<QForm | null>(null);

/** COMPUTED */
const statusOptions = computed(() => [
  { label: t.status.enabled.value, value: true },
  { label: t.status.disabled.value, value: false },
]);

const localData = reactive({
  name: props.modelValue.name || '',
  phone: props.modelValue.phone || '',
  enabled: props.modelValue.enabled ?? true,
});

/** WATCHERS */
watch(
  () => props.modelValue,
  (newVal) => {
    localData.name = newVal.name || '';
    localData.phone = newVal.phone || '';
    localData.enabled = newVal.enabled ?? true;
  },
  { deep: true },
);

/** FUNCTIONS */

/**
 * Emit updated value to parent
 */
function updateValue(): void {
  emit('update:modelValue', { ...localData });
}

/** EXPOSE */
defineExpose({
  formRef,
});
</script>

<style scoped>
.rounded-borders {
  border-radius: var(--mapex-radius-md);
}
</style>
