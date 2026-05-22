<template>
  <q-form ref="formRef" greedy>
    <div class="q-mb-md">
      <div class="text-subtitle1 text-weight-medium q-mb-xs">
        <q-icon name="location_on" color="primary" class="q-mr-xs" />
        {{ t.sections.address.value }}
      </div>
      <div class="text-body2 text-grey-7">
        {{ t.formDescriptions.address.value }}
      </div>
    </div>

    <div class="row q-col-gutter-md">
      <!-- Country -->
      <div class="col-12 col-md-6">
        <q-input
          v-model="localData.country"
          outlined
          dense
          class="rounded-borders"
          data-testid="customer-country-input"
          :label="t.fields.country.value"
          :rules="[
            (val) => !val || val.length <= COUNTRY_MAX_LENGTH || t.validation.countryMaxLength.value,
          ]"
          @update:model-value="updateValue"
        >
          <template #prepend>
            <q-icon name="public" color="primary" />
          </template>
        </q-input>
      </div>

      <!-- State -->
      <div class="col-12 col-md-6">
        <q-input
          v-model="localData.state"
          outlined
          dense
          class="rounded-borders"
          data-testid="customer-state-input"
          :label="t.fields.state.value"
          :rules="[
            (val) => !val || val.length <= STATE_MAX_LENGTH || t.validation.stateMaxLength.value,
          ]"
          @update:model-value="updateValue"
        >
          <template #prepend>
            <q-icon name="map" color="primary" />
          </template>
        </q-input>
      </div>

      <!-- City -->
      <div class="col-12 col-md-6">
        <q-input
          v-model="localData.city"
          outlined
          dense
          class="rounded-borders"
          data-testid="customer-city-input"
          :label="t.fields.city.value"
          :rules="[
            (val) => !val || val.length <= CITY_MAX_LENGTH || t.validation.cityMaxLength.value,
          ]"
          @update:model-value="updateValue"
        >
          <template #prepend>
            <q-icon name="location_city" color="primary" />
          </template>
        </q-input>
      </div>

      <!-- Zip Code -->
      <div class="col-12 col-md-6">
        <q-input
          v-model="localData.zipCode"
          outlined
          dense
          class="rounded-borders"
          data-testid="customer-zipcode-input"
          :label="t.fields.zipCode.value"
          :rules="[
            (val) => !val || val.length <= ZIPCODE_MAX_LENGTH || t.validation.zipCodeMaxLength.value,
          ]"
          @update:model-value="updateValue"
        >
          <template #prepend>
            <q-icon name="markunread_mailbox" color="primary" />
          </template>
        </q-input>
      </div>
    </div>
  </q-form>
</template>

<script setup lang="ts">
defineOptions({
  name: 'Step2Address',
});

/** TYPE IMPORTS */
import type { Step2AddressProps } from './interfaces/Step2Address.interface';
import type { QForm } from 'quasar';
import type { OrganizationFormData, Address } from '../../interfaces';

/** VUE IMPORTS */
import { ref, reactive, watch } from 'vue';

/** COMPOSABLES */
import { useAddCustomerTranslations } from '@composables/i18n';

/** LOCAL IMPORTS */
import {
  CITY_MAX_LENGTH,
  STATE_MAX_LENGTH,
  COUNTRY_MAX_LENGTH,
  ZIPCODE_MAX_LENGTH,
} from '../../constants';

const props = defineProps<Step2AddressProps>();

const emit = defineEmits<{
  (e: 'update:modelValue', value: Partial<OrganizationFormData>): void;
}>();

/** COMPOSABLES & STORES */
const t = useAddCustomerTranslations();

/** STATE */
const formRef = ref<QForm | null>(null);

const localData = reactive<Address>({
  city: props.modelValue.address?.city || '',
  state: props.modelValue.address?.state || '',
  country: props.modelValue.address?.country || '',
  zipCode: props.modelValue.address?.zipCode || '',
});

/** WATCHERS */
watch(
  () => props.modelValue,
  (newVal) => {
    localData.city = newVal.address?.city || '';
    localData.state = newVal.address?.state || '';
    localData.country = newVal.address?.country || '';
    localData.zipCode = newVal.address?.zipCode || '';
  },
  { deep: true }
);

/** FUNCTIONS */

/**
 * Emit updated value to parent
 */
function updateValue(): void {
  emit('update:modelValue', { address: { ...localData } });
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
