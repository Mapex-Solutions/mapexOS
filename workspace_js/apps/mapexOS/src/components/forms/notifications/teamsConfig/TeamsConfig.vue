<script setup lang="ts">
defineOptions({
  name: 'TeamsConfig'
});

import type { ChannelTeamsProps } from './interfaces';

import { ref, computed, watch } from 'vue';

// Props and Emits
const props = defineProps<{
  modelValue: ChannelTeamsProps;
}>();

const emit = defineEmits<{
  (e: 'update:modelValue', value: ChannelTeamsProps): void;
}>();

// Local state
const localData = ref<ChannelTeamsProps>({ ...props.modelValue });

// Computed properties
/**
 * Generates a preview message by replacing placeholders in the template
 * with example values for demonstration purposes
 */
const previewMessage = computed(() => {
  const template = localData.value.messageTemplate;
  return template
      .replace(/{{title}}/g, 'Example Title')
      .replace(/{{message}}/g, 'This is an example message to visualize how it will look in Microsoft Teams.');
});

/**
 * Formats the preview message as JSON with proper indentation
 * Falls back to plain text if JSON parsing fails
 */
const formattedPreview = computed(() => {
  try {
    const jsonObj = JSON.parse(previewMessage.value);
    return JSON.stringify(jsonObj, null, 2);
  } catch {
    return previewMessage.value;
  }
});

// Methods
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

// Watchers
/**
 * Watches for changes in props.modelValue and updates local data
 */
watch(() => props.modelValue, (newValue) => {
  localData.value = { ...newValue };
}, { deep: true });

/**
 * Watches for changes in local data and emits updates to parent
 */
watch(() => localData.value, (newValue) => {
  emit('update:modelValue', newValue);
}, { deep: true });
</script>

<template>
  <div class="teams-config">
    <div class="row q-col-gutter-md">
      <div class="col-12 col-md-6">
        <q-input
            v-model="localData.teamName"
            outlined
            label="Team Name *"
            :rules="[val => !!val || 'Team Name is required']"
        />
      </div>

      <div class="col-12 col-md-6">
        <q-toggle
            v-model="localData.adaptiveCard"
            color="primary"
            label="Use Adaptive Card"
        />
      </div>

      <div class="col-12">
        <q-select
            v-model="localData.channelsName"
            outlined
            use-input
            use-chips
            multiple
            input-debounce="0"
            new-value-mode="add-unique"
            label="Channels *"
            :rules="[val => val.length > 0 || 'Select at least one channel']"
        >
          <template v-slot:no-option>
            <q-item>
              <q-item-section class="text-grey">
                Type to add a new channel
              </q-item-section>
            </q-item>
          </template>
        </q-select>
      </div>

      <div class="col-12">
        <q-input
            v-model="localData.webhookUrl"
            outlined
            label="Webhook URL *"
            :rules="[
            val => !!val || 'Webhook URL is required',
            val => isValidUrl(val) || 'Invalid URL'
          ]"
        />
      </div>

      <div class="col-12">
        <q-input
            v-model="localData.messageTemplate"
            outlined
            autogrow
            type="textarea"
            label="Message Template (JSON) *"
            hint="Use JSON format with {{title}} and {{message}} as placeholders"
            :rules="[
            val => !!val || 'Message template is required',
            val => isValidJson(val) || 'Invalid JSON'
          ]"
        />
      </div>

      <div class="col-12">
        <q-card flat bordered class="bg-grey-1">
          <q-card-section>
            <div class="text-subtitle2">Message Example</div>
            <div class="q-mt-sm message-preview">
              <pre>{{ formattedPreview }}</pre>
            </div>
          </q-card-section>
        </q-card>
      </div>
    </div>
  </div>
</template>

<style scoped>
.message-preview {
  max-height: 200px;
  overflow-y: auto;
  background: var(--mapex-surface-sunken);
  padding: 10px;
  border-radius: var(--mapex-radius-xs);
  font-family: monospace;
}
</style>