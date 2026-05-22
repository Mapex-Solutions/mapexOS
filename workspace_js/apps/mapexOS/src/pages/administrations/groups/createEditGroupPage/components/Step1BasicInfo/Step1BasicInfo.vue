<template>
  <q-form ref="formRef" greedy>
    <div class="q-mb-md">
      <div class="text-subtitle1 text-weight-medium q-mb-xs">
        <q-icon name="groups" color="primary" class="q-mr-xs" />
        {{ t.sections.basicInfo.value }}
      </div>
      <div class="text-body2 text-grey-7">
        {{ t.formDescriptions.basicInfo.value }}
      </div>
    </div>

    <div class="row q-col-gutter-md">
      <!-- Name -->
      <div class="col-12 col-md-8">
        <q-input
          v-model="localData.name"
          outlined
          dense
          class="rounded-borders"
          :label="`${t.fields.name.value} *`"
          :hint="t.formDescriptions.name.value"
          :rules="[
            (val) => !!val || t.validation.nameRequired.value,
            (val) => val.length >= NAME_MIN_LENGTH || t.validation.nameMinLength.value,
            (val) => val.length <= NAME_MAX_LENGTH || t.validation.nameMaxLength.value,
          ]"
          @update:model-value="updateValue"
        >
          <template #prepend>
            <q-icon name="badge" color="primary" />
          </template>
        </q-input>
      </div>

      <!-- Status Select -->
      <div class="col-12 col-md-4">
        <q-select
          v-model="localData.enabled"
          outlined
          dense
          emit-value
          map-options
          :options="statusOptions"
          :label="`${t.fields.status.value} *`"
          :hint="t.formDescriptions.enabled.value"
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
          :rows="3"
          class="rounded-borders"
          :label="t.fields.description.value"
          :hint="t.formDescriptions.description.value"
          :rules="[
            (val) => !val || val.length <= DESCRIPTION_MAX_LENGTH || t.validation.descriptionMaxLength.value,
          ]"
          @update:model-value="updateValue"
        >
          <template #prepend>
            <q-icon name="description" color="primary" />
          </template>
        </q-input>
      </div>
    </div>
  </q-form>
</template>

<script setup lang="ts">
defineOptions({
  name: 'Step1BasicInfo',
});

/** TYPE IMPORTS */
import type { Step1BasicInfoProps } from './interfaces/Step1BasicInfo.interface';
import type { QForm } from 'quasar';
import type { GroupFormData } from '../../interfaces';

/** VUE IMPORTS */
import { ref, reactive, watch, computed } from 'vue';

/** COMPOSABLES */
import { useGroupsTranslations } from '@composables/i18n';

/** LOCAL IMPORTS */
import {
  NAME_MIN_LENGTH,
  NAME_MAX_LENGTH,
  DESCRIPTION_MAX_LENGTH,
} from '../../constants';

const props = withDefaults(defineProps<Step1BasicInfoProps>(), {
  isEditMode: false,
});

const emit = defineEmits<{
  (e: 'update:modelValue', value: Partial<GroupFormData>): void;
}>();

/** COMPOSABLES & STORES */
const t = useGroupsTranslations();

/** STATE */
const formRef = ref<QForm | null>(null);

const localData = reactive({
  name: props.modelValue.name || '',
  description: props.modelValue.description || '',
  enabled: props.modelValue.enabled ?? true,
});

/** STATUS OPTIONS */
const statusOptions = computed(() => [
  { label: t.statusOptions.active.value, value: true },
  { label: t.statusOptions.inactive.value, value: false },
]);

/** WATCHERS */
watch(
  () => props.modelValue,
  (newVal) => {
    localData.name = newVal.name || '';
    localData.description = newVal.description || '';
    localData.enabled = newVal.enabled ?? true;
  },
  { deep: true }
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

<style scoped lang="scss">
.rounded-borders {
  border-radius: var(--mapex-radius-md);
}
</style>
