<script setup lang="ts">
defineOptions({
  name: 'HttpConfig'
});

/** TYPE IMPORTS */
import type { Trigger } from '../../../interfaces';

/** VUE IMPORTS */
import { ref, watch } from 'vue';

/** COMPONENTS */

/** COMPOSABLES */
import { useLogger } from '@composables/useLogger';
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
const logger = useLogger('HttpConfig');
const translations = useCreateEditTriggerTranslations();

/** STATE */

/**
 * Local config state
 */
const config = ref({
  endpoint: props.modelValue.config?.endpoint || '',
  method: props.modelValue.config?.method || 'POST',
  headers: props.modelValue.config?.headers || {},
  body: props.modelValue.config?.body || {},
  timeout: props.modelValue.config?.timeout || 30000,
});

/**
 * HTTP method options
 */
const methodOptions = ['GET', 'POST', 'PUT', 'PATCH', 'DELETE'];

/**
 * Headers in editable format
 */
const headers = ref<Array<{ key: string; value: string }>>(
  Object.entries(config.value.headers).map(([key, value]) => ({ key, value: String(value) }))
);

/**
 * Body JSON editor
 */
const bodyJson = ref(
  Object.keys(config.value.body).length > 0
    ? JSON.stringify(config.value.body, null, 2)
    : '{\n  \n}'
);

/**
 * JSON validation error
 */
const bodyJsonError = ref<string | null>(null);

/** COMPUTED */

/** WATCHERS */

/**
 * Watch config changes and emit updates
 */
watch(
  config,
  (newConfig) => {
    logger.debug('Emitting config update:', newConfig);
    emit('update:modelValue', {
      ...props.modelValue,
      config: newConfig,
    });
  },
  { deep: true }
);

/**
 * Watch headers array and update config
 */
watch(
  headers,
  (newHeaders) => {
    const headersObj: Record<string, string> = {};
    newHeaders.forEach((h) => {
      if (h.key && h.value) {
        headersObj[h.key] = h.value;
      }
    });
    config.value.headers = headersObj;
  },
  { deep: true }
);

/**
 * Watch body JSON and update config
 */
watch(bodyJson, (newBodyJson) => {
  try {
    const parsed = JSON.parse(newBodyJson);
    logger.debug('Parsed body JSON:', parsed);
    config.value.body = parsed;
    bodyJsonError.value = null;
  } catch (error) {
    // Store error but keep the text for user to fix
    bodyJsonError.value = error instanceof Error ? error.message : translations.step4Configs.http.invalidJson.value;
    logger.debug('Invalid JSON:', bodyJsonError.value);
    // Don't update config.value.body if JSON is invalid
  }
});

/** FUNCTIONS */

/**
 * Add a new header row
 * @returns {void}
 */
function addHeader(): void {
  headers.value.push({ key: '', value: '' });
}

/**
 * Remove a header row
 * @param {number} index - Index of header to remove
 * @returns {void}
 */
function removeHeader(index: number): void {
  headers.value.splice(index, 1);
}

/**
 * Format JSON body
 * @returns {void}
 */
function formatJson(): void {
  try {
    const parsed = JSON.parse(bodyJson.value);
    bodyJson.value = JSON.stringify(parsed, null, 2);
  } catch {
    // Invalid JSON, show error
  }
}

/** LIFECYCLE HOOKS */
</script>

<template>
  <div class="http-config">
    <div class="row q-col-gutter-md">
      <!-- Endpoint -->
      <div class="col-12">
        <q-input
          v-model="config.endpoint"
          outlined
          dense
          :label="translations.step4Configs.http.endpointLabel.value"
          :placeholder="translations.step4Configs.http.endpointPlaceholder.value"
          :hint="translations.step4Configs.http.endpointHint.value"
          :rules="[(val: string) => !!val || translations.step4Configs.http.endpointRequired.value]"
        >
          <template v-slot:prepend>
            <q-icon name="link" />
          </template>
        </q-input>
      </div>

      <!-- Method and Timeout -->
      <div class="col-12 col-md-6">
        <q-select
          v-model="config.method"
          outlined
          dense
          :label="translations.step4Configs.http.methodLabel.value"
          :options="methodOptions"
          emit-value
          map-options
        >
          <template v-slot:prepend>
            <q-icon name="http" />
          </template>
        </q-select>
      </div>

      <div class="col-12 col-md-6">
        <q-input
          v-model.number="config.timeout"
          outlined
          dense
          type="number"
          :label="translations.step4Configs.http.timeoutLabel.value"
          :placeholder="translations.step4Configs.http.timeoutPlaceholder.value"
          :hint="translations.step4Configs.http.timeoutHint.value"
        >
          <template v-slot:prepend>
            <q-icon name="schedule" />
          </template>
        </q-input>
      </div>

      <!-- Headers Section -->
      <div class="col-12">
        <div class="row items-center q-mb-sm">
          <div class="col">
            <div class="text-subtitle2 text-weight-medium text-dark">{{ translations.step4Configs.http.headersTitle.value }}</div>
          </div>
          <div class="col-auto">
            <q-btn
              flat
              dense
              size="sm"
              icon="add"
              color="primary"
              :label="translations.step4Configs.http.addHeaderButton.value"
              @click="addHeader"
            />
          </div>
        </div>

        <div v-if="headers.length > 0" class="headers-list">
          <div
            v-for="(header, index) in headers"
            :key="index"
            class="row q-col-gutter-sm q-mb-sm items-center"
          >
            <div class="col-5">
              <q-input
                v-model="header.key"
                outlined
                dense
                :placeholder="translations.step4Configs.http.headerNamePlaceholder.value"
              />
            </div>
            <div class="col">
              <q-input
                v-model="header.value"
                outlined
                dense
                :placeholder="translations.step4Configs.http.headerValuePlaceholder.value"
              />
            </div>
            <div class="col-auto">
              <q-btn
                flat
                dense
                round
                size="sm"
                icon="delete"
                color="negative"
                @click="removeHeader(index)"
              />
            </div>
          </div>
        </div>

        <q-banner v-else rounded class="bg-grey-2 text-grey-7">
          <div class="text-body2">{{ translations.step4Configs.http.noHeadersMessage.value }}</div>
        </q-banner>
      </div>

      <!-- Body Section -->
      <div class="col-12">
        <div class="row items-center q-mb-sm">
          <div class="col">
            <div class="text-subtitle2 text-weight-medium text-dark">{{ translations.step4Configs.http.bodyTitle.value }}</div>
          </div>
          <div class="col-auto">
            <q-btn
              flat
              dense
              size="sm"
              icon="format_align_left"
              color="primary"
              :label="translations.step4Configs.http.formatJsonButton.value"
              @click="formatJson"
            />
          </div>
        </div>

        <q-input
          v-model="bodyJson"
          outlined
          dense
          type="textarea"
          :placeholder="translations.step4Configs.http.bodyPlaceholder.value"
          :hint="bodyJsonError || translations.step4Configs.http.bodyHint.value"
          :error="!!bodyJsonError"
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
.http-config {
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
