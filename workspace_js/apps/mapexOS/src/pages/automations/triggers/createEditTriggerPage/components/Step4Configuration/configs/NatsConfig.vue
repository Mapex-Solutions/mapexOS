<script setup lang="ts">
defineOptions({
  name: 'NatsConfig'
});

/** TYPE IMPORTS */
import type { Trigger } from '../../../interfaces';

/** VUE IMPORTS */
import { ref, watch } from 'vue';

/** COMPOSABLES */
import { useCreateEditTriggerTranslations } from '@composables/i18n/pages/automations/triggers/createEditTrigger';

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
const config = ref({
  server: props.modelValue.config?.server || '',
  subject: props.modelValue.config?.subject || '',
  username: props.modelValue.config?.username || '',
  password: props.modelValue.config?.password || '',
  token: props.modelValue.config?.token || '',
  message: props.modelValue.config?.message || {},
  useTLS: props.modelValue.config?.useTLS || false,
});

const messageJson = ref(JSON.stringify(config.value.message, null, 2));

/** WATCHERS */
watch(config, (newConfig) => {
  emit('update:modelValue', { ...props.modelValue, config: newConfig });
}, { deep: true });

watch(messageJson, (newMessageJson) => {
  try {
    config.value.message = JSON.parse(newMessageJson);
  } catch {
    // Invalid JSON
  }
});

/** FUNCTIONS */
function formatJson(): void {
  try {
    const parsed = JSON.parse(messageJson.value);
    messageJson.value = JSON.stringify(parsed, null, 2);
  } catch {
    // Invalid JSON
  }
}
</script>

<template>
  <div class="nats-config">
    <div class="row q-col-gutter-md">
      <div class="col-12">
        <q-input
          v-model="config.server"
          outlined
          dense
          :label="translations.step4Configs.nats.serverLabel.value"
          :placeholder="translations.step4Configs.nats.serverPlaceholder.value"
          :hint="translations.step4Configs.nats.serverHint.value"
          :rules="[(val: string) => !!val || translations.step4Configs.nats.serverRequired.value]"
        >
          <template v-slot:prepend>
            <q-icon name="dns" />
          </template>
        </q-input>
      </div>

      <div class="col-12">
        <q-input
          v-model="config.subject"
          outlined
          dense
          :label="translations.step4Configs.nats.subjectLabel.value"
          :placeholder="translations.step4Configs.nats.subjectPlaceholder.value"
          :hint="translations.step4Configs.nats.subjectHint.value"
          :rules="[(val: string) => !!val || translations.step4Configs.nats.subjectRequired.value]"
        >
          <template v-slot:prepend>
            <q-icon name="topic" />
          </template>
        </q-input>
      </div>

      <div class="col-12">
        <q-toggle
          v-model="config.useTLS"
          :label="translations.step4Configs.nats.useTlsLabel.value"
          color="primary"
        />
      </div>

      <div class="col-12 col-md-6">
        <q-input
          v-model="config.username"
          outlined
          dense
          :label="translations.step4Configs.nats.usernameLabel.value"
          :placeholder="translations.step4Configs.nats.usernamePlaceholder.value"
          :hint="translations.step4Configs.nats.usernameHint.value"
        >
          <template v-slot:prepend>
            <q-icon name="person" />
          </template>
        </q-input>
      </div>

      <div class="col-12 col-md-6">
        <q-input
          v-model="config.password"
          outlined
          dense
          type="password"
          :label="translations.step4Configs.nats.passwordLabel.value"
          :placeholder="translations.step4Configs.nats.passwordPlaceholder.value"
          :hint="translations.step4Configs.nats.passwordHint.value"
        >
          <template v-slot:prepend>
            <q-icon name="lock" />
          </template>
        </q-input>
      </div>

      <div class="col-12">
        <q-input
          v-model="config.token"
          outlined
          dense
          :label="translations.step4Configs.nats.tokenLabel.value"
          :placeholder="translations.step4Configs.nats.tokenPlaceholder.value"
          :hint="translations.step4Configs.nats.tokenHint.value"
        >
          <template v-slot:prepend>
            <q-icon name="vpn_key" />
          </template>
        </q-input>
      </div>

      <div class="col-12">
        <div class="row items-center q-mb-sm">
          <div class="col">
            <div class="text-subtitle2 text-weight-medium text-dark">{{ translations.step4Configs.nats.messageTitle.value }}</div>
          </div>
          <div class="col-auto">
            <q-btn
              flat
              dense
              size="sm"
              icon="format_align_left"
              color="primary"
              :label="translations.step4Configs.nats.formatJsonButton.value"
              @click="formatJson"
            />
          </div>
        </div>

        <q-input
          v-model="messageJson"
          outlined
          dense
          type="textarea"
          :placeholder="translations.step4Configs.nats.messagePlaceholder.value"
          :hint="translations.step4Configs.nats.messageHint.value"
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
.nats-config {
  .json-editor {
    font-family: 'Courier New', monospace;
  }
}
</style>
