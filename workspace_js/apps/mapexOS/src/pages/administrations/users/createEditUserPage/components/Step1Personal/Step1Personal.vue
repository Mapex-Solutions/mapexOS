<template>
  <q-form ref="formRef" greedy>
    <!-- Personal Information Section -->
    <div class="q-mb-md">
      <div class="text-subtitle1 text-weight-medium q-mb-xs">
        <q-icon name="person" color="primary" class="q-mr-xs" />
        {{ t.sections.personalInfo.value }}
      </div>
      <div class="text-body2 text-grey-7">
        {{ t.formDescriptions.personal.value }}
      </div>
    </div>

    <div class="row q-col-gutter-md">
      <!-- First Name -->
      <div class="col-12 col-md-6">
        <q-input
          v-model="localData.firstName"
          outlined
          dense
          class="rounded-borders"
          data-testid="user-firstname-input"
          :label="`${t.fields.firstName.value} *`"
          :rules="[
            (val) => !!val || t.validation.firstNameRequired.value,
            (val) => val.length >= VALIDATION.FIRST_NAME_MIN_LENGTH || t.validation.firstNameMinLength.value,
            (val) => val.length <= VALIDATION.FIRST_NAME_MAX_LENGTH || t.validation.firstNameMaxLength.value,
          ]"
          @update:model-value="updateValue"
        >
          <template #prepend>
            <q-icon name="person" color="primary" />
          </template>
        </q-input>
      </div>

      <!-- Last Name -->
      <div class="col-12 col-md-6">
        <q-input
          v-model="localData.lastName"
          outlined
          dense
          class="rounded-borders"
          data-testid="user-lastname-input"
          :label="`${t.fields.lastName.value} *`"
          :rules="[
            (val) => !!val || t.validation.lastNameRequired.value,
            (val) => val.length >= VALIDATION.LAST_NAME_MIN_LENGTH || t.validation.lastNameMinLength.value,
            (val) => val.length <= VALIDATION.LAST_NAME_MAX_LENGTH || t.validation.lastNameMaxLength.value,
          ]"
          @update:model-value="updateValue"
        >
          <template #prepend>
            <q-icon name="person" color="primary" />
          </template>
        </q-input>
      </div>

      <!-- Email -->
      <div class="col-12">
        <q-input
          v-model="localData.email"
          outlined
          dense
          type="email"
          class="rounded-borders"
          data-testid="user-email-input"
          :label="`${t.fields.email.value} *`"
          :rules="[
            (val) => !!val || t.validation.emailRequired.value,
            (val) => EMAIL_VALIDATION_REGEX.test(val) || t.validation.emailInvalid.value,
            (val) => val.length <= VALIDATION.EMAIL_MAX_LENGTH || t.validation.emailMaxLength.value,
          ]"
          :disable="isEditMode"
          @update:model-value="updateValue"
        >
          <template #prepend>
            <q-icon name="email" color="primary" />
          </template>
        </q-input>
      </div>

      <!-- Phone -->
      <div class="col-12 col-md-6">
        <q-input
          v-model="localData.phone"
          outlined
          dense
          type="tel"
          class="rounded-borders"
          data-testid="user-phone-input"
          :label="t.fields.phone.value"
          :hint="t.hints.phoneFormat.value"
          @update:model-value="updateValue"
        >
          <template #prepend>
            <q-icon name="phone" color="primary" />
          </template>
        </q-input>
      </div>

      <!-- Job Title -->
      <div class="col-12 col-md-6">
        <q-input
          v-model="localData.jobTitle"
          outlined
          dense
          class="rounded-borders"
          data-testid="user-jobtitle-input"
          :label="t.fields.jobTitle.value"
          :rules="[
            (val) => !val || val.length <= VALIDATION.JOB_TITLE_MAX_LENGTH || t.validation.jobTitleMaxLength.value,
          ]"
          @update:model-value="updateValue"
        >
          <template #prepend>
            <q-icon name="work" color="primary" />
          </template>
        </q-input>
      </div>
    </div>
  </q-form>
</template>

<script setup lang="ts">
defineOptions({
  name: 'Step1Personal',
});

/** TYPE IMPORTS */
import type { Step1PersonalProps } from './interfaces/Step1Personal.interface';
import type { QForm } from 'quasar';
import type { UserFormData } from '../../interfaces';

/** VUE IMPORTS */
import { ref, reactive, watch } from 'vue';

/** COMPOSABLES */
import { useAddUserTranslations } from '@composables/i18n';

/** LOCAL IMPORTS */
import { VALIDATION, EMAIL_VALIDATION_REGEX } from '../../constants';

const props = withDefaults(defineProps<Step1PersonalProps>(), {
  isEditMode: false,
});

const emit = defineEmits<{
  (e: 'update:modelValue', value: Partial<UserFormData>): void;
}>();

/** COMPOSABLES & STORES */
const t = useAddUserTranslations();

/** STATE */
const formRef = ref<QForm | null>(null);

const localData = reactive({
  firstName: props.modelValue.firstName || '',
  lastName: props.modelValue.lastName || '',
  email: props.modelValue.email || '',
  phone: props.modelValue.phone || '',
  jobTitle: props.modelValue.jobTitle || '',
});

/** WATCHERS */
watch(
  () => props.modelValue,
  (newVal) => {
    localData.firstName = newVal.firstName || '';
    localData.lastName = newVal.lastName || '';
    localData.email = newVal.email || '';
    localData.phone = newVal.phone || '';
    localData.jobTitle = newVal.jobTitle || '';
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
