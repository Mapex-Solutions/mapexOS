<script setup lang="ts">
defineOptions({
  name: 'SlackConfig'
});

import type { ChannelSlackProps } from './interfaces';

import { ref, computed, watch } from 'vue';

// Props and Emits
const props = defineProps<{
  modelValue: ChannelSlackProps;
}>();

const emit = defineEmits<{
  (e: 'update:modelValue', value: ChannelSlackProps): void;
}>();

// Local state
const localData = ref<ChannelSlackProps>({ ...props.modelValue });

// Computed properties
/**
 * Generates a preview of the message template by replacing placeholders with example content
 * @returns {string} Message template with example content replacing placeholders
 */
const previewMessage = computed(() => {
  const template = localData.value.messageTemplate;
  return template
      .replace(/{{title}}/g, 'Example Title')
      .replace(/{{message}}/g, 'This is an example message to preview how it will look in Slack.');
});

// Methods
/**
 * Validates if the provided string is a valid URL
 * @param {string} url - The URL string to validate
 * @returns {boolean} True if the URL is valid, false otherwise
 */
function isValidUrl(url: string): boolean {
  try {
    new URL(url);
    return true;
  } catch {
    return false;
  }
}

// Watchers
/**
 * Watches for changes in props.modelValue and updates local data accordingly
 * Deep watcher to catch nested property changes
 */
watch(() => props.modelValue, (newValue) => {
  localData.value = { ...newValue };
}, { deep: true });

/**
 * Watches for changes in local data and emits updates to parent component
 * Deep watcher to catch nested property changes
 */
watch(() => localData.value, (newValue) => {
  emit('update:modelValue', newValue);
}, { deep: true });
</script>

<template>
  <div class="slack-config">
    <div class="row q-col-gutter-md">
      <div class="col-12 col-md-6">
        <q-input
            v-model="localData.workspace"
            outlined
            label="Workspace *"
            :rules="[val => !!val || 'Workspace is required']"
        />
      </div>

      <div class="col-12 col-md-6">
        <q-input
            v-model="localData.botName"
            outlined
            label="Bot Name"
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
            label="Message Template *"
            hint="Use {{title}} and {{message}} as placeholders"
            :rules="[val => !!val || 'Message template is required']"
        />
      </div>

      <div class="col-12">
        <q-card flat bordered class="bg-grey-1">
          <q-card-section>
            <div class="text-subtitle2">Message Example</div>
            <div class="q-mt-sm">
              {{ previewMessage }}
            </div>
          </q-card-section>
        </q-card>
      </div>
    </div>
  </div>
</template>