<template>
  <q-form ref="formRef" greedy>
    <!-- Security Information Section -->
    <div class="q-mb-md">
      <div class="text-subtitle1 text-weight-medium q-mb-xs">
        <q-icon name="security" color="primary" class="q-mr-xs" />
        {{ t.sections.security.value }}
      </div>
      <div class="text-body2 text-grey-7">
        {{ t.formDescriptions.security.value }}
      </div>
    </div>

    <div class="row q-col-gutter-md">
      <!-- Password -->
      <div class="col-12 col-md-6">
        <q-input
          v-model="localData.password"
          outlined
          dense
          :type="showPassword ? 'text' : 'password'"
          class="rounded-borders"
          data-testid="user-password-input"
          :label="`${t.fields.password.value} ${isEditMode ? '' : '*'}`"
          :hint="t.hints.passwordRequirements.value"
          :rules="passwordRules"
          @update:model-value="updateValue"
        >
          <template #prepend>
            <q-icon name="lock" color="primary" />
          </template>
          <template #append>
            <q-icon
              :name="showPassword ? 'visibility_off' : 'visibility'"
              class="cursor-pointer"
              @click="showPassword = !showPassword"
            />
          </template>
        </q-input>
      </div>

      <!-- Confirm Password -->
      <div class="col-12 col-md-6">
        <q-input
          v-model="confirmPassword"
          outlined
          dense
          :type="showConfirmPassword ? 'text' : 'password'"
          class="rounded-borders"
          data-testid="user-confirm-password-input"
          :label="`${t.fields.confirmPassword.value} ${isEditMode ? '' : '*'}`"
          :rules="confirmPasswordRules"
        >
          <template #prepend>
            <q-icon name="lock" color="primary" />
          </template>
          <template #append>
            <q-icon
              :name="showConfirmPassword ? 'visibility_off' : 'visibility'"
              class="cursor-pointer"
              @click="showConfirmPassword = !showConfirmPassword"
            />
          </template>
        </q-input>
      </div>

      <!-- Change Password on Next Login -->
      <div class="col-12">
        <div class="q-py-sm">
          <q-checkbox
            v-model="localData.changePasswordNextLogin"
            color="primary"
            class="q-mb-xs"
            data-testid="user-change-pwd-checkbox"
            :label="t.fields.changePasswordNextLogin.value"
            @update:model-value="updateValue"
          />
          <div class="text-caption text-grey-7 q-pl-lg">
            {{ t.hints.changePasswordNextLogin.value }}
          </div>
        </div>
      </div>

      <!-- Enabled Toggle -->
      <div class="col-12">
        <div class="q-py-sm">
          <q-checkbox
            v-model="localData.enabled"
            color="primary"
            class="q-mb-xs"
            data-testid="user-enabled-checkbox"
            :label="t.fields.enabled.value"
            @update:model-value="updateValue"
          />
          <div class="text-caption text-grey-7 q-pl-lg">
            {{ t.hints.enabled.value }}
          </div>
        </div>
      </div>
    </div>
  </q-form>
</template>

<script setup lang="ts">
defineOptions({
  name: 'Step2Security',
});

/** TYPE IMPORTS */
import type { Step2SecurityProps } from './interfaces/Step2Security.interface';
import type { QForm } from 'quasar';
import type { UserFormData } from '../../interfaces';

/** VUE IMPORTS */
import { ref, reactive, watch, computed } from 'vue';

/** COMPOSABLES */
import { useAddUserTranslations } from '@composables/i18n';

/** LOCAL IMPORTS */
import { VALIDATION } from '../../constants';

const props = withDefaults(defineProps<Step2SecurityProps>(), {
  isEditMode: false,
});

const emit = defineEmits<{
  (e: 'update:modelValue', value: Partial<UserFormData>): void;
}>();

/** COMPOSABLES & STORES */
const t = useAddUserTranslations();

/** STATE */
const formRef = ref<QForm | null>(null);
const showPassword = ref(false);
const showConfirmPassword = ref(false);
const confirmPassword = ref('');

const localData = reactive({
  password: props.modelValue.password || '',
  changePasswordNextLogin: props.modelValue.changePasswordNextLogin || false,
  enabled: props.modelValue.enabled ?? true,
});

/** COMPUTED */

/**
 * Password validation rules
 */
const passwordRules = computed(() => {
  if (props.isEditMode) {
    // In edit mode, password is optional but if provided must be valid
    return [
      (val: string) => !val || val.length >= VALIDATION.PASSWORD_MIN_LENGTH || t.validation.passwordMinLength.value,
      (val: string) => !val || val.length <= VALIDATION.PASSWORD_MAX_LENGTH || t.validation.passwordMaxLength.value,
    ];
  }
  // In create mode, password is required
  return [
    (val: string) => !!val || t.validation.passwordRequired.value,
    (val: string) => val.length >= VALIDATION.PASSWORD_MIN_LENGTH || t.validation.passwordMinLength.value,
    (val: string) => val.length <= VALIDATION.PASSWORD_MAX_LENGTH || t.validation.passwordMaxLength.value,
  ];
});

/**
 * Confirm password validation rules
 */
const confirmPasswordRules = computed(() => {
  if (props.isEditMode && !localData.password) {
    return [];
  }
  return [
    (val: string) => val === localData.password || t.validation.passwordMismatch.value,
  ];
});

/** WATCHERS */
watch(
  () => props.modelValue,
  (newVal) => {
    localData.password = newVal.password || '';
    localData.changePasswordNextLogin = newVal.changePasswordNextLogin || false;
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
