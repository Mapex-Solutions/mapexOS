<script setup lang="ts">
defineOptions({
  name: 'TelegramConfig'
});

import type { ChannelTelegramProps } from './interfaces';

import { ref, computed, watch } from 'vue';
import { PARSE_MODES } from './constants';

// Props and Emits
const props = defineProps<{
  modelValue: ChannelTelegramProps;
}>();

const emit = defineEmits<{
  (e: 'update:modelValue', value: ChannelTelegramProps): void;
}>();

// Local state
const localData = ref<ChannelTelegramProps>({ ...props.modelValue });
const showToken = ref(false);

// Computed properties
/**
 * Generates a preview message by replacing placeholders in the template
 * with example values for demonstration purposes
 */
const previewMessage = computed(() => {
  const template = localData.value.messageTemplate;
  return template
      .replace(/{{title}}/g, 'Example Title')
      .replace(/{{message}}/g, 'This is an example message to visualize how it will look in Telegram.');
});

/**
 * Returns the appropriate CSS class based on the selected parse mode
 * for styling the message preview
 */
const getParseClass = computed(() => {
  switch (localData.value.parseMode) {
    case 'Markdown':
    case 'MarkdownV2':
      return 'telegram-markdown';
    case 'HTML':
      return 'telegram-html';
    default:
      return '';
  }
});

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
  <div class="telegram-config">
    <div class="row q-col-gutter-md">
      <div class="col-12 col-md-6">
        <q-input
            v-model="localData.botName"
            outlined
            label="Bot Name *"
            hint="Full bot name, including @ if necessary"
            :rules="[val => !!val || 'Bot Name is required']"
        />
      </div>

      <div class="col-12 col-md-6">
        <q-input
            v-model="localData.botToken"
            outlined
            label="Bot Token *"
            :rules="[val => !!val || 'Bot Token is required']"
            :type="showToken ? 'text' : 'password'"
        >
          <template v-slot:append>
            <q-icon
                class="cursor-pointer"
                :name="showToken ? 'visibility_off' : 'visibility'"
                @click="showToken = !showToken"
            />
          </template>
        </q-input>
      </div>

      <div class="col-12">
        <q-select
            v-model="localData.chatNames"
            use-input
            use-chips
            multiple
            outlined
            input-debounce="0"
            new-value-mode="add-unique"
            label="Chat Names *"
            :rules="[val => val.length > 0 || 'Add at least one chat']"
        >
          <template v-slot:no-option>
            <q-item>
              <q-item-section class="text-grey">
                Type the chat name and press Enter
              </q-item-section>
            </q-item>
          </template>
        </q-select>
      </div>

      <div class="col-12 col-md-6">
        <q-select
            v-model="localData.parseMode"
            outlined
            emit-value
            label="Formatting Mode"
            :options="PARSE_MODES"
        >
          <template v-slot:option="scope">
            <q-item v-bind="scope.itemProps">
              <q-item-section>
                <q-item-label>{{ scope.opt.label || scope.opt }}</q-item-label>
                <q-item-label v-if="scope.opt.description" caption>
                  {{ scope.opt.description }}
                </q-item-label>
              </q-item-section>
            </q-item>
          </template>
        </q-select>
      </div>

      <div class="col-12 col-md-6">
        <q-toggle
            v-model="localData.disableNotification"
            label="Disable Sound Notifications"
            color="primary"
            hint="Send message silently"
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
        <q-expansion-item
            icon="visibility"
            label="Message Preview"
            header-class="text-primary"
            expand-icon-class="text-primary"
        >
          <q-card>
            <q-card-section>
              <div class="telegram-preview">
                <div class="text-subtitle2 q-mb-sm">
                  <q-avatar size="24px" font-size="16px" color="primary" text-color="white">T</q-avatar>
                  {{ localData.botName }}
                </div>
                <div class="telegram-message" :class="getParseClass">
                  {{ previewMessage }}
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
.telegram-preview {
  background: var(--mapex-surface-sunken);
  padding: 15px;
  border-radius: var(--mapex-radius-md);
}

.telegram-message {
  background: var(--mapex-surface-bg);
  padding: 10px;
  border-radius: var(--mapex-radius-md);
  max-width: 80%;
  box-shadow: var(--mapex-shadow-xs);
  white-space: pre-wrap;
}
</style>