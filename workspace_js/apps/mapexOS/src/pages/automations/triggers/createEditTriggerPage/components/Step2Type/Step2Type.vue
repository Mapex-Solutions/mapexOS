<script setup lang="ts">
defineOptions({
  name: 'Step2Type'
});

/** TYPE IMPORTS */
import type { Trigger, TriggerTypeOption, TriggerType } from '../../interfaces';

/** VUE IMPORTS */
import { ref, watch, computed } from 'vue';

/** COMPONENTS */

/** COMPOSABLES */
import { useCreateEditTriggerTranslations } from '@composables/i18n/pages/automations/triggers/createEditTrigger';

/** UTILS */

/** SERVICES */

/** STORES */

/** LOCAL IMPORTS */
import { TRIGGER_TYPE_OPTIONS } from '../../constants';

/** PROPS & EMITS */
const props = defineProps<{
  modelValue: Trigger;
}>();

const emit = defineEmits<{
  'update:modelValue': [value: Trigger];
  'type-selected': [type: TriggerType];
}>();

/** COMPOSABLES & STORES */
const t = useCreateEditTriggerTranslations();

/** STATE */

/**
 * Currently selected trigger type option
 */
const selectedType = ref(null as TriggerTypeOption | null);

/** COMPUTED */

/**
 * Filter trigger types based on selected category
 */
const filteredTriggerTypes = computed<TriggerTypeOption[]>(() => {
  return TRIGGER_TYPE_OPTIONS.filter(
    (type) => type.category === props.modelValue.category
  );
});

/** WATCHERS */

/**
 * Watch for type selection changes and emit updates
 */
watch(selectedType, (newType) => {
  if (newType) {
    const updatedTrigger: Trigger = {
      ...props.modelValue,
      triggerType: newType.value,
    };
    emit('update:modelValue', updatedTrigger);
    emit('type-selected', newType.value);
  }
});

/** FUNCTIONS */

/**
 * Handle trigger type card click
 * @param {TriggerTypeOption} type - Selected trigger type option
 * @returns {void}
 */
function selectType(type: TriggerTypeOption): void {
  selectedType.value = type;
}

/** LIFECYCLE HOOKS */
</script>

<template>
  <div class="step2-type">
    <div class="text-body1 text-grey-8 q-mb-lg">
      {{ modelValue.category === 'technical' ? t.steps.step2.introTechnical.value : t.steps.step2.introCommunication.value }}
    </div>

    <div class="row q-col-gutter-md">
      <div
        v-for="type in filteredTriggerTypes"
        :key="type.value"
        class="col-12 col-md-6"
      >
        <q-card
          flat
          bordered
          :class="[
            'type-card cursor-pointer transition-all',
            selectedType?.value === type.value ? 'selected' : ''
          ]"
          @click="selectType(type)"
        >
          <q-card-section class="q-pa-md">
            <div class="row items-center no-wrap">
              <div class="col-auto q-mr-md">
                <q-avatar size="56px" color="primary-1" text-color="primary">
                  <q-icon :name="type.icon" size="32px" />
                </q-avatar>
              </div>
              <div class="col">
                <div class="text-h6 text-weight-medium text-dark">{{ type.label }}</div>
                <div class="text-body2 text-grey-7">{{ type.description }}</div>
              </div>
              <div v-if="selectedType?.value === type.value" class="col-auto">
                <q-icon name="check_circle" color="primary" size="28px" />
              </div>
            </div>
          </q-card-section>
        </q-card>
      </div>
    </div>

    <!-- Info box -->
    <q-banner rounded class="bg-blue-1 text-blue-9 q-mt-lg">
      <template v-slot:avatar>
        <q-icon name="info" color="blue-7" />
      </template>
      <div class="text-body2">
        <strong>{{ t.steps.step2.note.prefix.value }}:</strong> {{ t.steps.step2.note.text.value }}
      </div>
    </q-banner>
  </div>
</template>

<style lang="scss" scoped>
.step2-type {
  .type-card {
    border: 2px solid transparent;
    transition: var(--mapex-transition-slow);

    &:hover {
      border-color: var(--mapex-active-border);
      transform: translateY(-2px);
      box-shadow: var(--mapex-shadow-md);
    }

    &.selected {
      border-color: var(--mapex-active-border);
      background-color: var(--mapex-active-bg);
    }
  }
}
</style>
