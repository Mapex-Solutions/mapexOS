<script setup lang="ts">
/** TYPE IMPORTS */
import type { Step1BasicInfoProps, Step1BasicInfoEmits } from './interfaces/Step1BasicInfo.interface';

defineOptions({
  name: 'Step1BasicInfo'
});

import type { AssetTemplateData } from '../../interfaces';
import type { AssetClassification } from '@components/forms/assetClassificationSelector';
import { computed } from 'vue';
import { AssetClassificationSelector } from '@components/forms/assetClassificationSelector';
import { useAddAssetTemplateTranslations } from '@src/composables/i18n/pages/assets/addAssetTemplate/useAddAssetTemplateTranslations';

// Props
// Emits
const props = defineProps<Step1BasicInfoProps>();
const emit = defineEmits<Step1BasicInfoEmits>();

// i18n
const t = useAddAssetTemplateTranslations();

// Computed v-model for data binding
const data = computed({
  get: () => props.modelValue,
  set: (value: AssetTemplateData) => emit('update:modelValue', value)
});

// Status options
const statusOptions = computed(() => [
  { label: t.statusOptions.active.label.value, value: t.statusOptions.active.value },
  { label: t.statusOptions.inactive.label.value, value: t.statusOptions.inactive.value },
]);

// Computed property for asset classification
const classification = computed<AssetClassification | undefined>({
  get() {
    const assetData = data.value;
    if (assetData.categoryId && assetData.manufacturerId && assetData.modelId && assetData.version) {
      return {
        categoryId: assetData.categoryId,
        categoryName: assetData.categoryName,
        manufacturerId: assetData.manufacturerId,
        manufacturerName: assetData.manufacturerName,
        modelId: assetData.modelId,
        modelName: assetData.modelName,
        version: assetData.version,
      };
    }
    return undefined;
  },
  set(value: AssetClassification | undefined) {
    const updated = { ...data.value };

    if (value) {
      updated.categoryId = value.categoryId;
      updated.categoryName = value.categoryName;
      updated.manufacturerId = value.manufacturerId;
      updated.manufacturerName = value.manufacturerName;
      updated.modelId = value.modelId;
      updated.modelName = value.modelName;
      updated.version = value.version;
    } else {
      updated.categoryId = undefined;
      updated.categoryName = undefined;
      updated.manufacturerId = undefined;
      updated.manufacturerName = undefined;
      updated.modelId = undefined;
      updated.modelName = undefined;
      updated.version = undefined;
    }

    emit('update:modelValue', updated);
  }
});
</script>

<template>
  <div>
    <div class="q-mb-md">
      <div class="text-subtitle1 text-weight-medium q-mb-xs">
        <q-icon name="info" color="primary" class="q-mr-xs" />
        {{ t.steps.step1.title.value }}
      </div>
      <div class="text-body2 text-grey-7">
        {{ t.steps.step1.subtitle.value }}
      </div>
    </div>

    <div class="row q-col-gutter-md">
      <div class="col-12 col-md-6">
        <q-input
          v-model="data.name"
          outlined
          dense
          class="rounded-borders"
          :label="t.steps.step1.fields.name.label.value + ' *'"
          :placeholder="t.steps.step1.fields.name.placeholder.value"
          :hint="t.steps.step1.fields.name.hint.value"
          :rules="[(val) => !!val || t.steps.step1.fields.name.required.value]"
        >
          <template v-slot:prepend>
            <q-icon name="label" color="primary" />
          </template>
        </q-input>
      </div>

      <div class="col-12 col-md-6">
        <q-select
          v-model="data.enabled"
          outlined
          dense
          emit-value
          map-options
          class="rounded-borders"
          option-label="label"
          option-value="value"
          :label="t.steps.step1.fields.status.label.value + ' *'"
          :placeholder="t.steps.step1.fields.status.placeholder.value"
          :hint="t.steps.step1.fields.status.hint.value"
          :options="statusOptions"
          :rules="[(val) => val !== null && val !== undefined || t.steps.step1.fields.status.required.value]"
        >
          <template v-slot:prepend>
            <q-icon name="toggle_on" color="primary" />
          </template>
        </q-select>
      </div>

      <div class="col-12">
        <q-input
          v-model="data.description"
          outlined
          dense
          type="textarea"
          rows="3"
          class="rounded-borders"
          :label="t.steps.step1.fields.description.label.value"
          :placeholder="t.steps.step1.fields.description.placeholder.value"
          :hint="t.steps.step1.fields.description.hint.value"
        >
          <template v-slot:prepend>
            <q-icon name="notes" color="primary" />
          </template>
        </q-input>
      </div>

      <!-- Asset Classification Selector -->
      <div class="col-12">
        <AssetClassificationSelector
          v-model="classification"
          :required="true"
        />
      </div>

      <!-- isTemplate checkbox - Only for Vendor/Customer -->
      <div v-if="canCreateTemplate" class="col-12">
        <div class="q-py-sm">
          <q-checkbox
            v-model="data.isTemplate"
            color="primary"
            class="q-mb-xs"
            :label="t.steps.step1.fields.isTemplate.label.value"
          />
          <div class="text-caption text-grey-7 q-pl-lg">
            {{ t.steps.step1.fields.isTemplate.hint.value }}
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.rounded-borders {
  border-radius: var(--mapex-radius-md);
}
</style>
