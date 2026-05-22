<template>
  <q-form ref="formRef" greedy>
    <!-- Auth Config Section (V1: Locked to Internal) -->
    <div class="q-mb-lg">
      <div class="q-mb-md">
        <div class="text-subtitle1 text-weight-medium q-mb-xs">
          <q-icon name="security" color="primary" class="q-mr-xs" />
          {{ t.sections.authConfig.value }}
        </div>
        <div class="text-body2 text-grey-7">
          {{ t.formDescriptions.authConfig.value }}
        </div>
      </div>

      <div class="row q-col-gutter-md">
        <!-- Auth Provider Selection (V1: locked to internal) -->
        <div class="col-12">
          <div class="text-body2 text-weight-medium q-mb-sm">
            {{ t.fields.authProvider.value }}
          </div>
          <div class="row q-col-gutter-sm">
            <div
              v-for="option in AUTH_PROVIDER_OPTIONS"
              :key="option.value"
              class="col-12 col-md-6"
            >
              <q-card
                flat
                bordered
                class="option-card"
                :class="{
                  'option-card--selected': option.value === 'internal',
                  'option-card--disabled': option.disabled,
                }"
              >
                <q-card-section class="row items-center no-wrap q-pa-md">
                  <q-icon
                    :name="option.icon"
                    size="md"
                    :color="option.value === 'internal' ? 'primary' : 'grey-4'"
                    class="q-mr-md"
                  />
                  <div class="col">
                    <div class="text-subtitle2 text-weight-medium" :class="{ 'text-grey-5': option.disabled }">
                      {{ option.label }}
                    </div>
                    <div class="text-caption text-grey-7">
                      {{ option.description }}
                    </div>
                  </div>
                  <q-radio
                    :model-value="'internal'"
                    :val="option.value"
                    color="primary"
                    :disable="true"
                  />
                </q-card-section>
              </q-card>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Access Policy Section -->
    <div>
      <div class="q-mb-md">
        <div class="text-subtitle1 text-weight-medium q-mb-xs">
          <q-icon name="policy" color="primary" class="q-mr-xs" />
          {{ t.sections.accessPolicy.value }}
        </div>
        <div class="text-body2 text-grey-7">
          {{ t.formDescriptions.accessPolicy.value }}
        </div>
      </div>

      <div class="row q-col-gutter-md">
        <!-- Role Policy Selection -->
        <div class="col-12">
          <div class="text-body2 text-weight-medium q-mb-sm">
            {{ t.fields.rolePolicy.value }} *
          </div>
          <div class="row q-col-gutter-sm">
            <div
              v-for="option in ROLE_POLICY_OPTIONS"
              :key="option.value"
              class="col-12 col-md-6"
            >
              <q-card
                flat
                bordered
                class="option-card cursor-pointer"
                :class="{
                  'option-card--selected': localData.accessPolicy.rolePolicy === option.value,
                }"
                :data-testid="`customer-role-policy-${option.value}`"
                @click="selectRolePolicy(option.value)"
              >
                <q-card-section class="row items-center no-wrap q-pa-md">
                  <q-icon
                    :name="option.icon"
                    size="md"
                    :color="localData.accessPolicy.rolePolicy === option.value ? 'primary' : 'grey-6'"
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
                    v-model="localData.accessPolicy.rolePolicy"
                    :val="option.value"
                    color="primary"
                  />
                </q-card-section>
              </q-card>
            </div>
          </div>
        </div>

        <!-- Default Scope Selection -->
        <div class="col-12">
          <div class="text-body2 text-weight-medium q-mb-sm">
            {{ t.fields.defaultScope.value }} *
          </div>
          <div class="row q-col-gutter-sm">
            <div
              v-for="option in DEFAULT_SCOPE_OPTIONS"
              :key="option.value"
              class="col-12 col-md-6"
            >
              <q-card
                flat
                bordered
                class="option-card cursor-pointer"
                :class="{
                  'option-card--selected': localData.accessPolicy.defaultScope === option.value,
                }"
                :data-testid="`customer-scope-${option.value}`"
                @click="selectDefaultScope(option.value)"
              >
                <q-card-section class="row items-center no-wrap q-pa-md">
                  <q-icon
                    :name="option.icon"
                    size="md"
                    :color="localData.accessPolicy.defaultScope === option.value ? 'primary' : 'grey-6'"
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
                    v-model="localData.accessPolicy.defaultScope"
                    :val="option.value"
                    color="primary"
                  />
                </q-card-section>
              </q-card>
            </div>
          </div>
        </div>
      </div>
    </div>
  </q-form>
</template>

<script setup lang="ts">
defineOptions({
  name: 'Step3AccessPolicy',
});

/** TYPE IMPORTS */
import type { Step3AccessPolicyProps } from './interfaces/Step3AccessPolicy.interface';
import type { QForm } from 'quasar';
import type {
  OrganizationFormData,
  RolePolicyType,
  DefaultScopeType,
} from '../../interfaces';

/** VUE IMPORTS */
import { ref, reactive, watch } from 'vue';

/** COMPOSABLES */
import { useAddCustomerTranslations } from '@composables/i18n';

/** LOCAL IMPORTS */
import {
  AUTH_PROVIDER_OPTIONS,
  ROLE_POLICY_OPTIONS,
  DEFAULT_SCOPE_OPTIONS,
} from '../../constants';

const props = defineProps<Step3AccessPolicyProps>();

const emit = defineEmits<{
  (e: 'update:modelValue', value: Partial<OrganizationFormData>): void;
}>();

/** COMPOSABLES & STORES */
const t = useAddCustomerTranslations();

/** STATE */
const formRef = ref<QForm | null>(null);

const localData = reactive({
  accessPolicy: {
    rolePolicy: props.modelValue.accessPolicy?.rolePolicy || 'strict',
    defaultScope: props.modelValue.accessPolicy?.defaultScope || 'local',
  },
});

/** WATCHERS */
watch(
  () => props.modelValue,
  (newVal) => {
    localData.accessPolicy.rolePolicy = newVal.accessPolicy?.rolePolicy || 'strict';
    localData.accessPolicy.defaultScope = newVal.accessPolicy?.defaultScope || 'local';
  },
  { deep: true },
);

/** FUNCTIONS */

/**
 * Select role policy
 *
 * @param {RolePolicyType} value - Role policy value
 */
function selectRolePolicy(value: RolePolicyType): void {
  localData.accessPolicy.rolePolicy = value;
  updateValue();
}

/**
 * Select default scope
 *
 * @param {DefaultScopeType} value - Default scope value
 */
function selectDefaultScope(value: DefaultScopeType): void {
  localData.accessPolicy.defaultScope = value;
  updateValue();
}

/**
 * Emit updated value to parent
 * V1: AuthConfig is always overridden to internal in the handler
 */
function updateValue(): void {
  emit('update:modelValue', {
    accessPolicy: { ...localData.accessPolicy },
  });
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

.option-card {
  border-radius: var(--mapex-radius-md);
  transition: var(--mapex-transition-base);

  &:hover:not(.option-card--disabled) {
    border-color: var(--q-primary);
    background-color: rgba(var(--q-primary-rgb), 0.05);
  }

  &--selected {
    border-color: var(--q-primary);
    background-color: rgba(var(--q-primary-rgb), 0.08);
  }

  &--disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }
}
</style>
