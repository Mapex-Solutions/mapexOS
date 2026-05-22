<script setup lang="ts">
defineOptions({
  name: 'WebhookConfig'
});

import type { ChannelWebhookProps } from './interfaces';
import { ref, watch, onMounted } from 'vue';

// Props and Emits
const props = defineProps<{
  modelValue: ChannelWebhookProps;
}>();

const emit = defineEmits<{
  (e: 'update:modelValue', value: ChannelWebhookProps): void;
}>();

// Local state
const localData = ref<ChannelWebhookProps>({ ...props.modelValue });
const headerKeys = ref<string[]>([]);
const headerValues = ref<string[]>([]);
const headerIndices = ref<number[]>([]);
const httpMethods = ['GET', 'POST', 'PUT', 'DELETE', 'PATCH'];

// State for webhook testing
const testTitle = ref('Test Title');
const testMessage = ref('This is a test message for the webhook');
const testResult = ref(null);
const testing = ref(false);

// Initialize headers arrays
onMounted(() => {
  initializeHeaders();
});

/**
 * Initializes the headers arrays from the local data
 * If no headers exist, adds an empty one
 */
function initializeHeaders() {
  headerKeys.value = [];
  headerValues.value = [];
  headerIndices.value = [];

  let index = 0;
  for (const key in localData.value.headers) {
    headerKeys.value[index] = key;
    headerValues.value[index] = localData.value.headers[key] as string;
    headerIndices.value.push(index);
    index++;
  }

  // If no headers exist, add an empty one
  if (index === 0) {
    addHeader();
  }
}

/**
 * Updates the headers object when a header key is changed
 * Rebuilds the entire headers object with new keys and values
 */
function updateHeaderKey() {
  const newHeaders: Record<string, string> = {};

  // Recreate the headers object with the new keys and values
  headerIndices.value.forEach((idx) => {
    if (headerKeys.value[idx]) {
      newHeaders[headerKeys.value[idx]] = headerValues.value[idx] || '';
    }
  });

  localData.value.headers = newHeaders;
  emit('update:modelValue', localData.value);
}

/**
 * Updates a specific header value
 * @param index - The index of the header to update
 */
function updateHeaderValue(index: number) {
  const key = headerKeys.value[index];
  if (key) {
    localData.value.headers[key] = headerValues.value[index] || '';
    emit('update:modelValue', localData.value);
  }
}

/**
 * Adds a new empty header row
 */
function addHeader() {
  const newIndex = headerIndices.value.length > 0
      ? Math.max(...headerIndices.value) + 1
      : 0;

  headerIndices.value.push(newIndex);
  headerKeys.value[newIndex] = '';
  headerValues.value[newIndex] = '';
}

/**
 * Removes a header row and updates the headers object
 * @param index - The index of the header to remove
 */
function removeHeader(index: number) {
  const key = headerKeys.value[index];
  if (key && localData.value.headers[key]) {
    const newHeaders: Record<string, string> = {};
    for (const k in localData.value.headers) {
      if (k !== key) {
        newHeaders[k] = localData.value.headers[k] as string;
      }
    }
    localData.value.headers = newHeaders;
  }

  const indexPosition = headerIndices.value.indexOf(index);
  if (indexPosition !== -1) {
    headerIndices.value.splice(indexPosition, 1);
  }

  emit('update:modelValue', localData.value);
}

/**
 * Validates if a given string is a valid URL
 * @param url - The URL string to validate
 * @returns true if valid URL, false otherwise
 */
function isValidUrl(url: string): boolean {
  try {
    new URL(url);
    return true;
  } catch {
    return false;
  }
}

/**
 * Validates if a given string is valid JSON
 * @param str - The JSON string to validate
 * @returns true if valid JSON, false otherwise
 */
function isValidJson(str: string): boolean {
  try {
    JSON.parse(str);
    return true;
  } catch {
    return false;
  }
}

/**
 * Tests the webhook configuration with sample data
 * Simulates a webhook call and displays the result
 */
async function testWebhook() {
  testing.value = true;
  try {
    // Here you would make an API call to test the webhook
    // For now, we'll simulate a response
    await new Promise(resolve => setTimeout(resolve, 1500));

    const currentTimestamp = new Date().toISOString();
    const samplePayload = localData.value.payload
        .replace(/{{title}}/g, testTitle.value)
        .replace(/{{message}}/g, testMessage.value)
        .replace(/{{timestamp}}/g, currentTimestamp);

    testResult.value = {
      success: true,
      request: {
        method: localData.value.method,
        url: localData.value.url,
        headers: localData.value.headers,
        payload: JSON.parse(samplePayload),
      },
      response: {
        status: 200,
        statusText: 'OK',
        data: {
          received: true,
          message: 'Webhook received successfully',
        },
      },
    } as any;
  } catch (error: any) {
    testResult.value = {
      success: false,
      error: error instanceof Error ? error.message : 'Unknown error',
    } as any;
  } finally {
    testing.value = false;
  }
}

// Watchers
/**
 * Watches for changes in props.modelValue and updates local data
 */
watch(() => props.modelValue, (newValue) => {
  localData.value = { ...newValue };
  initializeHeaders();
}, { deep: true });

/**
 * Watches for changes in local data and emits updates to parent
 */
watch(() => localData.value, (newValue) => {
  emit('update:modelValue', newValue);
}, { deep: true });
</script>

<template>
  <div class="webhook-config">
    <div class="row q-col-gutter-md">
      <div class="col-12 col-md-6">
        <q-input
            v-model="localData.name"
            outlined
            label="Webhook Name *"
            :rules="[val => !!val || 'Webhook Name is required']"
        />
      </div>

      <div class="col-12 col-md-6">
        <q-select
            v-model="localData.method"
            outlined
            emit-value
            label="HTTP Method *"
            :options="httpMethods"
            :rules="[val => !!val || 'HTTP Method is required']"
        >
          <template v-slot:option="scope">
            <q-item v-bind="scope.itemProps">
              <q-item-section>
                <q-item-label>{{ scope.opt }}</q-item-label>
              </q-item-section>
            </q-item>
          </template>
        </q-select>
      </div>

      <div class="col-12">
        <q-input
            v-model="localData.url"
            outlined
            label="Endpoint URL *"
            :rules="[
            val => !!val || 'Endpoint URL is required',
            val => isValidUrl(val) || 'Invalid URL'
          ]"
        />
      </div>

      <div class="col-12">
        <q-expansion-item
            default-opened
            icon="code"
            label="Headers"
            header-class="text-primary"
            expand-icon-class="text-primary"
        >
          <q-card>
            <q-card-section>
              <div class="headers-container">
                <div v-for="(index) in headerIndices" :key="index" class="row q-col-gutter-sm q-mb-sm">
                  <div class="col-5">
                    <q-input
                        v-model="headerKeys[index]"
                        outlined
                        dense
                        label="Key"
                        @update:model-value="updateHeaderKey()"
                    />
                  </div>
                  <div class="col-6">
                    <q-input
                        v-model="headerValues[index]"
                        outlined
                        dense
                        label="Value"
                        @update:model-value="updateHeaderValue(index)"
                    />
                  </div>
                  <div class="col-1 flex items-center">
                    <q-btn
                        flat
                        round
                        dense
                        icon="delete"
                        color="negative"
                        @click="removeHeader(index)"
                    />
                  </div>
                </div>

                <q-btn
                    flat
                    class="q-mt-sm"
                    label="Add Header"
                    color="primary"
                    @click="addHeader"
                />
              </div>
            </q-card-section>
          </q-card>
        </q-expansion-item>
      </div>

      <div class="col-12">
        <q-input
            v-model="localData.payload"
            outlined
            autogrow
            type="textarea"
            label="Payload Template *"
            hint="Use JSON format with {{title}}, {{message}} and {{timestamp}} as placeholders"
            :rules="[
            val => !!val || 'Payload Template is required',
            val => isValidJson(val) || 'Invalid JSON'
          ]"
        />
      </div>

      <div class="col-12 col-md-6">
        <q-input
            v-model.number="localData.timeout"
            outlined
            type="number"
            label="Timeout (ms)"
            hint="Maximum time to wait for response"
        />
      </div>

      <div class="col-12 col-md-6">
        <q-input
            v-model.number="localData.retryCount"
            outlined
            type="number"
            label="Retry Attempts"
            hint="Number of attempts in case of failure"
        />
      </div>

      <div class="col-12">
        <q-expansion-item
            icon="send"
            label="Webhook Test"
            header-class="text-primary"
            expand-icon-class="text-primary"
        >
          <q-card>
            <q-card-section>
              <div class="row q-col-gutter-md">
                <div class="col-12">
                  <q-input
                      v-model="testTitle"
                      outlined
                      dense
                      label="Test Title"
                  />
                </div>
                <div class="col-12">
                  <q-input
                      v-model="testMessage"
                      outlined
                      dense
                      autogrow
                      type="textarea"
                      label="Test Message"
                  />
                </div>
                <div class="col-12">
                  <q-btn
                      label="Test Webhook"
                      color="primary"
                      :loading="testing"
                      @click="testWebhook"
                  />
                </div>

                <div v-if="testResult" class="col-12">
                  <q-card flat bordered class="bg-grey-1">
                    <q-card-section>
                      <div class="text-subtitle2">Test Result</div>
                      <pre class="test-result">{{ JSON.stringify(testResult, null, 2) }}</pre>
                    </q-card-section>
                  </q-card>
                </div>
              </div>
            </q-card-section>
          </q-card>
        </q-expansion-item>
      </div>
    </div>
  </div>
</template>

<style scoped>
.headers-container {
  max-height: 300px;
  overflow-y: auto;
}

.test-result {
  background: var(--mapex-surface-sunken);
  padding: 10px;
  border-radius: var(--mapex-radius-xs);
  max-height: 200px;
  overflow-y: auto;
  font-family: monospace;
  font-size: 12px;
}
</style>