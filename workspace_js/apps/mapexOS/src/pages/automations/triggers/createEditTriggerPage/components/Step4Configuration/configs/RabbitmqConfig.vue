<script setup lang="ts">
defineOptions({
  name: 'RabbitmqConfig'
});

/** TYPE IMPORTS */
import type { Trigger } from '../../../interfaces';

/** VUE IMPORTS */
import { computed, ref, watch } from 'vue';

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
  host: props.modelValue.config?.host || '',
  port: props.modelValue.config?.port || 5672,
  vhost: props.modelValue.config?.vhost || '/',
  username: props.modelValue.config?.username || '',
  password: props.modelValue.config?.password || '',
  publishMode: props.modelValue.config?.publishMode || 'exchange',
  exchange: props.modelValue.config?.exchange || '',
  exchangeType: props.modelValue.config?.exchangeType || 'direct',
  routingKey: props.modelValue.config?.routingKey || '',
  queue: props.modelValue.config?.queue || '',
  message: props.modelValue.config?.message || {},
  useTLS: props.modelValue.config?.useTLS || false,
});

/**
 * Publish mode options
 */
const publishModeOptions = computed(() => [
  {
    label: translations.step4Configs.rabbitmq.publishModeExchange.value,
    value: 'exchange',
    description: translations.step4Configs.rabbitmq.publishModeExchangeDescription.value,
  },
  {
    label: translations.step4Configs.rabbitmq.publishModeQueue.value,
    value: 'queue',
    description: translations.step4Configs.rabbitmq.publishModeQueueDescription.value,
  },
]);

/**
 * Exchange type options
 */
const exchangeTypeOptions = computed(() => [
  { label: translations.step4Configs.rabbitmq.exchangeTypeDirect.value, value: 'direct' },
  { label: translations.step4Configs.rabbitmq.exchangeTypeFanout.value, value: 'fanout' },
  { label: translations.step4Configs.rabbitmq.exchangeTypeTopic.value, value: 'topic' },
  { label: translations.step4Configs.rabbitmq.exchangeTypeHeaders.value, value: 'headers' },
]);

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
  <div class="rabbitmq-config">
    <div class="row q-col-gutter-md">
      <!-- Host -->
      <div class="col-12 col-md-8">
        <q-input
          v-model="config.host"
          outlined
          dense
          :label="translations.step4Configs.rabbitmq.hostLabel.value"
          :placeholder="translations.step4Configs.rabbitmq.hostPlaceholder.value"
          :hint="translations.step4Configs.rabbitmq.hostHint.value"
          :rules="[(val: string) => !!val || translations.step4Configs.rabbitmq.hostRequired.value]"
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
          :label="translations.step4Configs.rabbitmq.portLabel.value"
          :placeholder="translations.step4Configs.rabbitmq.portPlaceholder.value"
          :hint="translations.step4Configs.rabbitmq.portHint.value"
          :rules="[(val: number) => !!val || translations.step4Configs.rabbitmq.portRequired.value]"
        >
          <template v-slot:prepend>
            <q-icon name="settings_ethernet" />
          </template>
        </q-input>
      </div>

      <!-- Virtual Host -->
      <div class="col-12 col-md-6">
        <q-input
          v-model="config.vhost"
          outlined
          dense
          :label="translations.step4Configs.rabbitmq.vhostLabel.value"
          :placeholder="translations.step4Configs.rabbitmq.vhostPlaceholder.value"
          :hint="translations.step4Configs.rabbitmq.vhostHint.value"
        >
          <template v-slot:prepend>
            <q-icon name="apartment" />
          </template>
        </q-input>
      </div>

      <!-- Use TLS -->
      <div class="col-12 col-md-6 flex items-center">
        <q-toggle
          v-model="config.useTLS"
          :label="translations.step4Configs.rabbitmq.useTlsLabel.value"
          color="primary"
        />
      </div>

      <!-- Username -->
      <div class="col-12 col-md-6">
        <q-input
          v-model="config.username"
          outlined
          dense
          :label="translations.step4Configs.rabbitmq.usernameLabel.value"
          :placeholder="translations.step4Configs.rabbitmq.usernamePlaceholder.value"
          :hint="translations.step4Configs.rabbitmq.usernameHint.value"
          :rules="[(val: string) => !!val || translations.step4Configs.rabbitmq.usernameRequired.value]"
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
          :label="translations.step4Configs.rabbitmq.passwordLabel.value"
          :placeholder="translations.step4Configs.rabbitmq.passwordPlaceholder.value"
          :hint="translations.step4Configs.rabbitmq.passwordHint.value"
          :rules="[(val: string) => !!val || translations.step4Configs.rabbitmq.passwordRequired.value]"
        >
          <template v-slot:prepend>
            <q-icon name="lock" />
          </template>
        </q-input>
      </div>

      <!-- Publish Mode Selector -->
      <div class="col-12">
        <q-separator class="q-my-md" />
        <div class="text-subtitle2 text-weight-medium text-dark q-mb-md">
          <q-icon name="publish" class="q-mr-sm" />
          {{ translations.step4Configs.rabbitmq.publishingTitle.value }}
        </div>
      </div>

      <div class="col-12">
        <q-select
          v-model="config.publishMode"
          outlined
          dense
          :label="translations.step4Configs.rabbitmq.publishModeLabel.value"
          :options="publishModeOptions"
          emit-value
          map-options
          option-label="label"
          option-value="value"
        >
          <template v-slot:prepend>
            <q-icon name="settings" />
          </template>
          <template v-slot:option="scope">
            <q-item v-bind="scope.itemProps">
              <q-item-section>
                <q-item-label>{{ scope.opt.label }}</q-item-label>
                <q-item-label caption>{{ scope.opt.description }}</q-item-label>
              </q-item-section>
            </q-item>
          </template>
        </q-select>
      </div>

      <!-- EXCHANGE MODE FIELDS -->
      <template v-if="config.publishMode === 'exchange'">
        <!-- Exchange -->
        <div class="col-12 col-md-8">
          <q-input
            v-model="config.exchange"
            outlined
            dense
            :label="translations.step4Configs.rabbitmq.exchangeLabel.value"
            :placeholder="translations.step4Configs.rabbitmq.exchangePlaceholder.value"
            :hint="translations.step4Configs.rabbitmq.exchangeHint.value"
            :rules="[(val: string) => !!val || translations.step4Configs.rabbitmq.exchangeRequired.value]"
          >
            <template v-slot:prepend>
              <q-icon name="swap_horiz" />
            </template>
          </q-input>
        </div>

        <!-- Exchange Type -->
        <div class="col-12 col-md-4">
          <q-select
            v-model="config.exchangeType"
            outlined
            dense
            :label="translations.step4Configs.rabbitmq.exchangeTypeLabel.value"
            :options="exchangeTypeOptions"
            emit-value
            map-options
          >
            <template v-slot:prepend>
              <q-icon name="category" />
            </template>
          </q-select>
        </div>

        <!-- Routing Key -->
        <div class="col-12">
          <q-input
            v-model="config.routingKey"
            outlined
            dense
            :label="translations.step4Configs.rabbitmq.routingKeyLabel.value"
            :placeholder="translations.step4Configs.rabbitmq.routingKeyPlaceholder.value"
            :hint="translations.step4Configs.rabbitmq.routingKeyHint.value"
          >
            <template v-slot:prepend>
              <q-icon name="route" />
            </template>
          </q-input>
        </div>
      </template>

      <!-- DIRECT QUEUE MODE FIELDS -->
      <template v-if="config.publishMode === 'queue'">
        <!-- Queue Name -->
        <div class="col-12">
          <q-input
            v-model="config.queue"
            outlined
            dense
            :label="translations.step4Configs.rabbitmq.queueLabel.value"
            :placeholder="translations.step4Configs.rabbitmq.queuePlaceholder.value"
            :hint="translations.step4Configs.rabbitmq.queueHint.value"
            :rules="[(val: string) => !!val || translations.step4Configs.rabbitmq.queueRequired.value]"
          >
            <template v-slot:prepend>
              <q-icon name="inbox" />
            </template>
          </q-input>
        </div>

        <!-- Info Banner -->
        <div class="col-12">
          <q-banner rounded class="bg-blue-1 text-blue-9">
            <template v-slot:avatar>
              <q-icon name="info" color="blue-7" />
            </template>
            <div class="text-body2">
              <strong>{{ translations.step4Configs.rabbitmq.directQueueBannerPrefix.value }}</strong> {{ translations.step4Configs.rabbitmq.directQueueBannerText.value }}
            </div>
          </q-banner>
        </div>
      </template>

      <!-- Message Section -->
      <div class="col-12">
        <div class="row items-center q-mb-sm">
          <div class="col">
            <div class="text-subtitle2 text-weight-medium text-dark">{{ translations.step4Configs.rabbitmq.messageTitle.value }}</div>
          </div>
          <div class="col-auto">
            <q-btn
              flat
              dense
              size="sm"
              icon="format_align_left"
              color="primary"
              :label="translations.step4Configs.rabbitmq.formatJsonButton.value"
              @click="formatJson"
            />
          </div>
        </div>

        <q-input
          v-model="messageJson"
          outlined
          dense
          type="textarea"
          :placeholder="translations.step4Configs.rabbitmq.messagePlaceholder.value"
          :hint="translations.step4Configs.rabbitmq.messageHint.value"
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
.rabbitmq-config {
  .json-editor {
    font-family: 'Courier New', monospace;
  }
}
</style>
