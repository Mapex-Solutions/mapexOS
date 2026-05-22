<template>
  <q-form ref="formRef" greedy>
    <div class="q-mb-md">
      <div class="text-subtitle1 text-weight-medium q-mb-xs">
        <q-icon name="fingerprint" color="primary" class="q-mr-xs" />
        {{ t.steps.step1.title.value }}
      </div>
      <div class="text-body2 text-grey-7">
        {{ t.steps.step1.subtitle.value }}
      </div>
    </div>

    <div class="row q-col-gutter-md">
      <!-- Name - Full width -->
      <div class="col-12">
        <q-input
          v-model="localData.name"
          outlined
          dense
          class="rounded-borders"
          data-testid="asset-name-input"
          :label="t.steps.step1.fields.name.label.value + ' *'"
          :placeholder="t.steps.step1.fields.name.placeholder.value"
          :hint="t.steps.step1.fields.name.hint.value"
          :rules="[(val) => !!val || t.steps.step1.fields.name.required.value]"
          @update:model-value="updateValue"
        >
          <template #prepend>
            <q-icon name="badge" color="primary" />
          </template>
        </q-input>
      </div>

      <!-- Asset ID - Half width on tablet/desktop -->
      <div class="col-xs-12 col-sm-6">
        <q-input
          v-model="localData.assetId"
          outlined
          dense
          class="rounded-borders"
          data-testid="asset-id-input"
          :label="t.steps.step1.fields.assetId.label.value + ' *'"
          :placeholder="t.steps.step1.fields.assetId.placeholder.value"
          :hint="t.steps.step1.fields.assetId.hint.value"
          :rules="[(val) => !!val || t.steps.step1.fields.assetId.required.value]"
          @update:model-value="updateValue"
        >
          <template #prepend>
            <q-icon name="tag" color="primary" />
          </template>
        </q-input>
      </div>

      <!-- Status - Half width on tablet/desktop -->
      <div class="col-xs-12 col-sm-6">
        <q-select
          v-model="localData.enabled"
          outlined
          dense
          emit-value
          map-options
          class="rounded-borders"
          data-testid="asset-status-select"
          :label="t.steps.step1.fields.status.label.value + ' *'"
          :placeholder="t.steps.step1.fields.status.placeholder.value"
          :hint="t.steps.step1.fields.status.hint.value"
          :options="statusOptions"
          option-label="label"
          option-value="value"
          :rules="[(val) => val !== null && val !== undefined || t.steps.step1.fields.status.required.value]"
          @update:model-value="updateValue"
        >
          <template #prepend>
            <q-icon name="toggle_on" color="primary" />
          </template>
        </q-select>
      </div>

      <!-- Description - Full width -->
      <div class="col-12">
        <q-input
          v-model="localData.description"
          outlined
          dense
          type="textarea"
          rows="3"
          class="rounded-borders"
          data-testid="asset-description-input"
          :label="t.steps.step1.fields.description.label.value"
          :placeholder="t.steps.step1.fields.description.placeholder.value"
          :hint="t.steps.step1.fields.description.hint.value"
          @update:model-value="updateValue"
        >
          <template #prepend>
            <q-icon name="notes" color="primary" />
          </template>
        </q-input>
      </div>

      <!-- Debug Mode Toggle -->
      <div class="col-12">
        <q-item tag="label" class="rounded-borders debug-banner">
          <q-item-section avatar>
            <q-icon name="mdi-bug-outline" color="warning" />
          </q-item-section>
          <q-item-section>
            <q-item-label class="debug-banner__label text-weight-medium">{{ t.steps.step4.fields.debugEnabled.label.value }}</q-item-label>
            <q-item-label caption class="debug-banner__hint">{{ t.steps.step4.fields.debugEnabled.hint.value }}</q-item-label>
          </q-item-section>
          <q-item-section side>
            <q-toggle
              v-model="localData.debugEnabled"
              color="warning"
              data-testid="asset-debug-toggle"
              @update:model-value="updateValue"
            />
          </q-item-section>
        </q-item>
      </div>
    </div>
  </q-form>
</template>

<script setup lang="ts">
/** TYPE IMPORTS */
import type { Step1IdentificationProps } from './interfaces/Step1Identification.interface';

defineOptions({
  name: 'Step1Identification'
});

import type { AssetFormData, SelectOption } from '../../interfaces';
import type { QForm } from 'quasar';

import { ref, computed, reactive, watch } from 'vue';

import { useAddAssetTranslations } from '@src/composables/i18n/pages/assets/addAsset/useAddAssetTranslations';

const props = defineProps<Step1IdentificationProps>();
const emit = defineEmits<{
  (e: 'update:modelValue', value: Partial<AssetFormData>): void;
}>();

const formRef = ref<QForm | null>(null);

const t = useAddAssetTranslations();

const localData = reactive({
  name: props.modelValue.name || '',
  assetId: props.modelValue.assetId || '',
  enabled: props.modelValue.enabled ?? true,
  description: props.modelValue.description || '',
  debugEnabled: props.modelValue.debugEnabled ?? false,
});

watch(() => props.modelValue, (newVal) => {
  localData.name = newVal.name || '';
  localData.assetId = newVal.assetId || '';
  localData.enabled = newVal.enabled ?? true;
  localData.description = newVal.description || '';
  localData.debugEnabled = newVal.debugEnabled ?? false;
}, { deep: true });

function updateValue() {
  emit('update:modelValue', { ...localData });
}

const statusOptions = computed((): SelectOption[] => [
  { label: t.statusOptions.active.label.value, value: t.statusOptions.active.value },
  { label: t.statusOptions.inactive.label.value, value: t.statusOptions.inactive.value },
]);

defineExpose({
  formRef,
});
</script>

<style scoped lang="scss">
.rounded-borders {
  border-radius: var(--mapex-radius-md);
}

.debug-banner {
  background: rgba(var(--q-warning-rgb), 0.08);
  border: 1px solid rgba(var(--q-warning-rgb), 0.2);
  border-radius: var(--mapex-radius-md);

  &__label {
    color: var(--mapex-text-primary);
  }

  &__hint {
    color: var(--mapex-text-secondary) !important;
  }
}
</style>
