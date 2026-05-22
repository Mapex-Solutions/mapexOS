<template>
  <q-form ref="formRef" greedy>
    <div class="q-mb-md">
      <div class="text-subtitle1 text-weight-medium q-mb-xs">
        <q-icon name="badge" color="primary" class="q-mr-xs" />
        {{ t.sections.basicInfo.value }}
      </div>
      <div class="text-body2 text-grey-7">
        {{ t.formDescriptions.basicInfo.value }}
      </div>
    </div>

    <div class="row q-col-gutter-md">
      <!-- Name -->
      <div class="col-12">
        <q-input
          v-model="localData.name"
          outlined
          dense
          class="rounded-borders"
          :label="`${t.fields.name.value} *`"
          :rules="[
            (val) => !!val || t.validation.nameRequired.value,
            (val) => val.length >= NAME_MIN_LENGTH || t.validation.nameMinLength.value,
            (val) => val.length <= NAME_MAX_LENGTH || t.validation.nameMaxLength.value,
          ]"
          :disable="isEditMode && isSystemRole"
          @update:model-value="updateValue"
        >
          <template #prepend>
            <q-icon name="badge" color="primary" />
          </template>
        </q-input>
      </div>

      <!-- Description -->
      <div class="col-12">
        <q-input
          v-model="localData.description"
          outlined
          dense
          type="textarea"
          autogrow
          class="rounded-borders"
          :label="t.fields.description.value"
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

      <!-- Scope Selection -->
      <div class="col-12">
        <div class="text-body2 text-weight-medium q-mb-sm">
          {{ t.fields.scope.value }} *
        </div>
        <div class="row q-col-gutter-sm">
          <div
            v-for="option in SCOPE_OPTIONS"
            :key="option.value"
            class="col-12 col-md-6"
          >
            <q-card
              flat
              bordered
              class="scope-card cursor-pointer"
              :class="{
                'scope-card--selected': localData.scope === option.value,
                'scope-card--disabled': isEditMode,
              }"
              @click="!isEditMode && selectScope(option.value)"
            >
              <q-card-section class="row items-center no-wrap q-pa-md">
                <q-icon
                  :name="option.icon"
                  size="md"
                  :color="localData.scope === option.value ? 'primary' : 'grey-6'"
                  class="q-mr-md"
                />
                <div class="col">
                  <div class="text-subtitle2 text-weight-medium">
                    {{ option.label }}
                  </div>
                  <div class="text-caption text-grey-7">
                    {{ option.description }}
                  </div>
                </div>
                <q-radio
                  v-model="localData.scope"
                  :val="option.value"
                  :disable="isEditMode"
                  color="primary"
                />
              </q-card-section>
            </q-card>
          </div>
        </div>
        <div v-if="!localData.scope && showScopeError" class="text-negative text-caption q-mt-xs">
          {{ t.validation.scopeRequired.value }}
        </div>
      </div>

      <!-- Is Template Checkbox -->
      <div class="col-12">
        <div class="q-py-sm">
          <q-checkbox
            v-model="localData.isTemplate"
            color="primary"
            class="q-mb-xs"
            :label="t.fields.isTemplate.value"
            :disable="isEditMode"
            @update:model-value="updateValue"
          />
          <div class="text-caption text-grey-7 q-pl-lg">
            {{ t.formDescriptions.isTemplate.value }}
          </div>
        </div>
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
import type { RoleFormData } from '../../interfaces';

/** VUE IMPORTS */
import { ref, reactive, watch } from 'vue';

/** COMPOSABLES */
import { useRolesTranslations } from '@composables/i18n';

/** LOCAL IMPORTS */
import {
  SCOPE_OPTIONS,
  NAME_MIN_LENGTH,
  NAME_MAX_LENGTH,
  DESCRIPTION_MAX_LENGTH,
} from '../../constants';

const props = withDefaults(defineProps<Step1BasicInfoProps>(), {
  isEditMode: false,
  isSystemRole: false,
});

const emit = defineEmits<{
  (e: 'update:modelValue', value: Partial<RoleFormData>): void;
}>();

/** COMPOSABLES & STORES */
const t = useRolesTranslations();

/** STATE */
const formRef = ref<QForm | null>(null);
const showScopeError = ref(false);

const localData = reactive({
  name: props.modelValue.name || '',
  description: props.modelValue.description || '',
  scope: props.modelValue.scope,
  isTemplate: props.modelValue.isTemplate || false,
});

/** WATCHERS */
watch(
  () => props.modelValue,
  (newVal) => {
    localData.name = newVal.name || '';
    localData.description = newVal.description || '';
    localData.scope = newVal.scope;
    localData.isTemplate = newVal.isTemplate || false;
  },
  { deep: true }
);

/** FUNCTIONS */

/**
 * Select scope option
 *
 * @param {string} value - Scope value ('global' | 'local')
 */
function selectScope(value: string): void {
  localData.scope = value as 'global' | 'local';
  showScopeError.value = false;
  updateValue();
}

/**
 * Emit updated value to parent
 */
function updateValue(): void {
  emit('update:modelValue', { ...localData });
}

/**
 * Validate the form including scope selection
 *
 * @returns {Promise<boolean>} Whether the form is valid
 */
async function validate(): Promise<boolean> {
  const formValid = await formRef.value?.validate();
  const scopeValid = !!localData.scope;

  if (!scopeValid) {
    showScopeError.value = true;
  }

  return !!formValid && scopeValid;
}

/** EXPOSE */
defineExpose({
  formRef,
  validate,
});
</script>

<style scoped lang="scss">
.rounded-borders {
  border-radius: var(--mapex-radius-md);
}

.scope-card {
  border-radius: var(--mapex-radius-md);
  transition: var(--mapex-transition-base);

  &:hover:not(.scope-card--disabled) {
    border-color: var(--q-primary);
    background-color: rgba(var(--q-primary-rgb), 0.05);
  }

  &--selected {
    border-color: var(--q-primary);
    background-color: rgba(var(--q-primary-rgb), 0.08);
  }

  &--disabled {
    opacity: 0.7;
    cursor: not-allowed;
  }
}
</style>
