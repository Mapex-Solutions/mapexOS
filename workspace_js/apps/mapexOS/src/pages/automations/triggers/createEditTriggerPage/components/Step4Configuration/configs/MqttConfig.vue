<script setup lang="ts">
defineOptions({
  name: 'MqttConfig'
});

/** TYPE IMPORTS */
import type { Trigger } from '../../../interfaces';

/** VUE IMPORTS */
import { ref, watch } from 'vue';

/** COMPONENTS */

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
const translations = useCreateEditTriggerTranslations();

/** STATE */

/**
 * Local config state
 */
const config = ref({
  broker: props.modelValue.config?.broker || '',
  port: props.modelValue.config?.port || 1883,
  topic: props.modelValue.config?.topic || '',
  qos: props.modelValue.config?.qos || 0,
  username: props.modelValue.config?.username || '',
  password: props.modelValue.config?.password || '',
  clientId: props.modelValue.config?.clientId || '',
  message: props.modelValue.config?.message || {},
  useTLS: props.modelValue.config?.useTLS || false,
});

/**
 * QoS options
 */
const qosOptions = [
  { label: '0 - At most once', value: 0 },
  { label: '1 - At least once', value: 1 },
  { label: '2 - Exactly once', value: 2 },
];

/**
 * Message JSON editor
 */
const messageJson = ref(JSON.stringify(config.value.message, null, 2));

/** COMPUTED */

/** WATCHERS */

/**
 * Watch config changes and emit updates
 */
watch(
  config,
  (newConfig) => {
    emit('update:modelValue', {
      ...props.modelValue,
      config: newConfig,
    });
  },
  { deep: true }
);

/**
 * Watch message JSON and update config
 */
watch(messageJson, (newMessageJson) => {
  try {
    config.value.message = JSON.parse(newMessageJson);
  } catch {
    // Invalid JSON, don't update
  }
});

/** FUNCTIONS */

/**
 * Format JSON message
 * @returns {void}
 */
function formatJson(): void {
  try {
    const parsed = JSON.parse(messageJson.value);
    messageJson.value = JSON.stringify(parsed, null, 2);
  } catch {
    // Invalid JSON, show error
  }
}

/** LIFECYCLE HOOKS */
</script>

<template>
  <div class="mqtt-config">
    <div class="row q-col-gutter-md">
      <!-- Broker URL -->
      <div class="col-12 col-md-8">
        <q-input
          v-model="config.broker"
          outlined
          dense
          :label="translations.step4Configs.mqtt.brokerLabel.value"
          :placeholder="translations.step4Configs.mqtt.brokerPlaceholder.value"
          :hint="translations.step4Configs.mqtt.brokerHint.value"
          :rules="[(val: string) => !!val || translations.step4Configs.mqtt.brokerRequired.value]"
        >
          <template v-slot:prepend>
            <q-icon name="dns" />
          </template>
        </q-input>
      </div>

      <!-- Port -->
      <div class="col-12 col-md-4">
        <q-input
          v-model.number="config.port"
          outlined
          dense
          type="number"
          :label="translations.step4Configs.mqtt.portLabel.value"
          :placeholder="translations.step4Configs.mqtt.portPlaceholder.value"
          :hint="translations.step4Configs.mqtt.portHint.value"
          :rules="[(val: number) => !!val || translations.step4Configs.mqtt.portRequired.value]"
        >
          <template v-slot:prepend>
            <q-icon name="settings_ethernet" />
          </template>
        </q-input>
      </div>

      <!-- Topic -->
      <div class="col-12">
        <q-input
          v-model="config.topic"
          outlined
          dense
          :label="translations.step4Configs.mqtt.topicLabel.value"
          :placeholder="translations.step4Configs.mqtt.topicPlaceholder.value"
          :hint="translations.step4Configs.mqtt.topicHint.value"
          :rules="[(val: string) => !!val || translations.step4Configs.mqtt.topicRequired.value]"
        >
          <template v-slot:prepend>
            <q-icon name="topic" />
          </template>
        </q-input>
      </div>

      <!-- QoS -->
      <div class="col-12 col-md-4">
        <q-select
          v-model="config.qos"
          outlined
          dense
          :label="translations.step4Configs.mqtt.qosLabel.value"
          :options="qosOptions"
          emit-value
          map-options
        >
          <template v-slot:prepend>
            <q-icon name="priority_high" />
          </template>
        </q-select>
      </div>

      <!-- Client ID -->
      <div class="col-12 col-md-8">
        <q-input
          v-model="config.clientId"
          outlined
          dense
          :label="translations.step4Configs.mqtt.clientIdLabel.value"
          :placeholder="translations.step4Configs.mqtt.clientIdPlaceholder.value"
          :hint="translations.step4Configs.mqtt.clientIdHint.value"
        >
          <template v-slot:prepend>
            <q-icon name="fingerprint" />
          </template>
        </q-input>
      </div>

      <!-- Use TLS -->
      <div class="col-12">
        <q-toggle
          v-model="config.useTLS"
          :label="translations.step4Configs.mqtt.useTlsLabel.value"
          color="primary"
        />
      </div>

      <!-- Username -->
      <div class="col-12 col-md-6">
        <q-input
          v-model="config.username"
          outlined
          dense
          :label="translations.step4Configs.mqtt.usernameLabel.value"
          :placeholder="translations.step4Configs.mqtt.usernamePlaceholder.value"
          :hint="translations.step4Configs.mqtt.usernameHint.value"
        >
          <template v-slot:prepend>
            <q-icon name="person" />
          </template>
        </q-input>
      </div>

      <!-- Password -->
      <div class="col-12 col-md-6">
        <q-input
          v-model="config.password"
          outlined
          dense
          type="password"
          :label="translations.step4Configs.mqtt.passwordLabel.value"
          :placeholder="translations.step4Configs.mqtt.passwordPlaceholder.value"
          :hint="translations.step4Configs.mqtt.passwordHint.value"
        >
          <template v-slot:prepend>
            <q-icon name="lock" />
          </template>
        </q-input>
      </div>

      <!-- Message Section -->
      <div class="col-12">
        <div class="row items-center q-mb-sm">
          <div class="col">
            <div class="text-subtitle2 text-weight-medium text-dark">{{ translations.step4Configs.mqtt.messageTitle.value }}</div>
          </div>
          <div class="col-auto">
            <q-btn
              flat
              dense
              size="sm"
              icon="format_align_left"
              color="primary"
              :label="translations.step4Configs.mqtt.formatJsonButton.value"
              @click="formatJson"
            />
          </div>
        </div>

        <q-input
          v-model="messageJson"
          outlined
          dense
          type="textarea"
          :placeholder="translations.step4Configs.mqtt.messagePlaceholder.value"
          :hint="translations.step4Configs.mqtt.messageHint.value"
          rows="10"
          class="json-editor"
        >
          <template v-slot:prepend>
            <q-icon name="code" />
          </template>
        </q-input>
      </div>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.mqtt-config {
  .json-editor {
    font-family: 'Courier New', monospace;
  }
}
</style>
