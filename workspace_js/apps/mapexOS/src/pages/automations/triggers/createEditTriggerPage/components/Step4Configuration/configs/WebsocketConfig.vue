<script setup lang="ts">
defineOptions({
  name: 'WebsocketConfig'
});

/** TYPE IMPORTS */
import type { Trigger } from '../../../interfaces';

/** VUE IMPORTS */
import { ref, watch } from 'vue';

/** COMPOSABLES */
import { useCreateEditTriggerTranslations } from '@composables/i18n/pages/automations/triggers/createEditTrigger';

/** PROPS & EMITS */
const props = defineProps<{ modelValue: Trigger }>();
const emit = defineEmits<{ 'update:modelValue': [value: Trigger] }>();

/** COMPOSABLES & STORES */
const translations = useCreateEditTriggerTranslations();

/** STATE */
const config = ref({
  url: props.modelValue.config?.url || '',
  message: props.modelValue.config?.message || {},
  headers: props.modelValue.config?.headers || {},
});

const messageJson = ref(JSON.stringify(config.value.message, null, 2));
const headers = ref<Array<{ key: string; value: string }>>(
  Object.entries(config.value.headers).map(([key, value]) => ({ key, value: String(value) }))
);

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

watch(headers, (newHeaders) => {
  const headersObj: Record<string, string> = {};
  newHeaders.forEach((h) => {
    if (h.key && h.value) headersObj[h.key] = h.value;
  });
  config.value.headers = headersObj;
}, { deep: true });

/** FUNCTIONS */
function formatJson(): void {
  try {
    const parsed = JSON.parse(messageJson.value);
    messageJson.value = JSON.stringify(parsed, null, 2);
  } catch {
    // Invalid JSON
  }
}

function addHeader(): void {
  headers.value.push({ key: '', value: '' });
}

function removeHeader(index: number): void {
  headers.value.splice(index, 1);
}
</script>

<template>
  <div class="websocket-config">
    <div class="row q-col-gutter-md">
      <div class="col-12">
        <q-input
          v-model="config.url"
          outlined
          dense
          :label="translations.step4Configs.websocket.urlLabel.value"
          :placeholder="translations.step4Configs.websocket.urlPlaceholder.value"
          :hint="translations.step4Configs.websocket.urlHint.value"
          :rules="[(val: string) => !!val || translations.step4Configs.websocket.urlRequired.value]"
        >
          <template v-slot:prepend>
            <q-icon name="link" />
          </template>
        </q-input>
      </div>

      <div class="col-12">
        <div class="row items-center q-mb-sm">
          <div class="col">
            <div class="text-subtitle2 text-weight-medium text-dark">{{ translations.step4Configs.websocket.headersTitle.value }}</div>
          </div>
          <div class="col-auto">
            <q-btn flat dense size="sm" icon="add" color="primary" :label="translations.step4Configs.websocket.addHeaderButton.value" @click="addHeader" />
          </div>
        </div>

        <div v-if="headers.length > 0" class="headers-list">
          <div v-for="(header, index) in headers" :key="index" class="row q-col-gutter-sm q-mb-sm items-center">
            <div class="col-5">
              <q-input v-model="header.key" outlined dense :placeholder="translations.step4Configs.websocket.headerNamePlaceholder.value" />
            </div>
            <div class="col">
              <q-input v-model="header.value" outlined dense :placeholder="translations.step4Configs.websocket.headerValuePlaceholder.value" />
            </div>
            <div class="col-auto">
              <q-btn flat dense round size="sm" icon="delete" color="negative" @click="removeHeader(index)" />
            </div>
          </div>
        </div>
      </div>

      <div class="col-12">
        <div class="row items-center q-mb-sm">
          <div class="col">
            <div class="text-subtitle2 text-weight-medium text-dark">{{ translations.step4Configs.websocket.messageTitle.value }}</div>
          </div>
          <div class="col-auto">
            <q-btn flat dense size="sm" icon="format_align_left" color="primary" :label="translations.step4Configs.websocket.formatJsonButton.value" @click="formatJson" />
          </div>
        </div>

        <q-input
          v-model="messageJson"
          outlined
          dense
          type="textarea"
          :placeholder="translations.step4Configs.websocket.messagePlaceholder.value"
          :hint="translations.step4Configs.websocket.messageHint.value"
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
.websocket-config {
  .headers-list {
    border: 1px solid var(--mapex-card-border);
    border-radius: var(--mapex-radius-xs);
    padding: 12px;
  }

  .json-editor {
    font-family: 'Courier New', monospace;
  }
}
</style>
