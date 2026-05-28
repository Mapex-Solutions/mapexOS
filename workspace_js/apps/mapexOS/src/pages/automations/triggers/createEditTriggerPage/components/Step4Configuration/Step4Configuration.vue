<script setup lang="ts">
defineOptions({
  name: 'Step4Configuration'
});

/** TYPE IMPORTS */
import type { Trigger } from '../../interfaces';
import type { QForm } from 'quasar';

/** VUE IMPORTS */
import { ref, computed } from 'vue';

/** COMPONENTS */
import HttpConfig from './configs/HttpConfig.vue';
import MqttConfig from './configs/MqttConfig.vue';
import RabbitmqConfig from './configs/RabbitmqConfig.vue';
import NatsConfig from './configs/NatsConfig.vue';
import WebsocketConfig from './configs/WebsocketConfig.vue';
import EmailConfig from './configs/EmailConfig.vue';
import TeamsConfig from './configs/TeamsConfig.vue';
import SlackConfig from './configs/SlackConfig.vue';

/** COMPOSABLES */
import { useCreateEditTriggerTranslations } from '@composables/i18n/pages/automations/triggers/createEditTrigger';

/** UTILS */

/** SERVICES */

/** STORES */

/** LOCAL IMPORTS */

/** PROPS & EMITS */
const props = defineProps<{
  modelValue: Trigger;
}>();

const emit = defineEmits<{
  'update:modelValue': [value: Trigger];
}>();

/** COMPOSABLES & STORES */
const t = useCreateEditTriggerTranslations();

/** CONSTANTS */
const PLACEHOLDER_SAMPLE = '{{placeholderName}}';

/** STATE */

/**
 * Form reference for validation
 */
const formRef = ref<QForm | null>(null);

/** COMPUTED */

/**
 * Get the appropriate config component based on trigger type
 */
const configComponent = computed(() => {
  const componentMap: Record<string, any> = {
    http: HttpConfig,
    mqtt: MqttConfig,
    rabbitmq: RabbitmqConfig,
    nats: NatsConfig,
    websocket: WebsocketConfig,
    email: EmailConfig,
    teams: TeamsConfig,
    slack: SlackConfig,
  };

  return componentMap[props.modelValue.triggerType] || null;
});

/** WATCHERS */

/** FUNCTIONS */

/**
 * Update trigger configuration
 * @param {Trigger} updatedTrigger - Updated trigger data
 * @returns {void}
 */
function updateTrigger(updatedTrigger: Trigger): void {
  emit('update:modelValue', updatedTrigger);
}

/**
 * Validate form
 * @returns {Promise<boolean>} Validation result
 */
async function validate(): Promise<boolean> {
  if (!formRef.value) return false;
  return await formRef.value.validate();
}

/** LIFECYCLE HOOKS */

/** EXPOSE */
defineExpose({
  formRef,
  validate,
});
</script>

<template>
  <div class="step4-configuration">
    <div class="text-body1 text-grey-8 q-mb-lg">
      Configure the specific settings for your {{ modelValue.triggerType.toUpperCase() }} trigger. Use placeholders like <code v-text="'{{placeholderName}}'"></code> for dynamic values.
    </div>

    <q-form ref="formRef">
      <!-- Dynamic Config Component -->
      <component
        v-if="configComponent"
        :is="configComponent"
        :model-value="modelValue"
        @update:model-value="updateTrigger"
      />

      <!-- Fallback if no config component -->
      <q-banner v-else rounded class="bg-orange-1 text-orange-9">
        <template v-slot:avatar>
          <q-icon name="warning" color="orange-7" />
        </template>
        <div class="text-body2">
          Configuration component for <strong>{{ modelValue.triggerType }}</strong> is not yet implemented.
        </div>
      </q-banner>

      <!-- Info box about placeholders -->
      <q-banner rounded class="bg-blue-1 text-blue-9 q-mt-lg">
        <template v-slot:avatar>
          <q-icon name="info" color="blue-7" />
        </template>
        <div class="text-body2">
          <strong>{{ t.steps.step4.placeholderSyntaxTitle.value }}</strong>
          {{ t.steps.step4.placeholderSyntaxDescription(PLACEHOLDER_SAMPLE) }}
        </div>
      </q-banner>
    </q-form>
  </div>
</template>

<style lang="scss" scoped>
.step4-configuration {
  code {
    background-color: var(--mapex-submenu-bg);
    padding: 2px 6px;
    border-radius: var(--mapex-radius-xs);
    font-family: 'Courier New', monospace;
    font-size: 0.9em;
  }
}
</style>
