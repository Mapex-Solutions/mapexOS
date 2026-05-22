<script setup lang="ts">
defineOptions({
  name: 'EmailConfig'
});

import type { ChannelEmailProps } from './interfaces';

import { ref, computed, watch } from 'vue';
import { DetailChip } from '@components/chips';

// Props and Emits
const props = defineProps<{
  modelValue: ChannelEmailProps;
}>();

const emit = defineEmits<{
  (e: 'update:modelValue', value: ChannelEmailProps): void;
}>();

// Local state
const localData = ref<ChannelEmailProps>({ ...props.modelValue });

// Computed properties
/**
 * Generates a preview of the email template by replacing placeholders with example content
 * @returns {string} HTML string with example content replacing template placeholders
 */
const previewEmail = computed(() => {
  const template = localData.value.template;
  return template
      .replace(/{{title}}/g, 'Example Title')
      .replace(/{{message}}/g, 'This is an example message to preview how the email will look.');
});

// Methods
/**
 * Validates if the provided string is a valid email address
 * @param {string} email - The email address to validate
 * @returns {boolean} True if the email is valid, false otherwise
 */
function isValidEmail(email: string): boolean {
  const regex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
  return regex.test(email);
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
  <div class="email-config">
    <div class="row q-col-gutter-md">
      <div class="col-12 col-md-6">
        <q-input
            v-model="localData.from"
            outlined
            label="Sender *"
            hint="Sender email address"
            :rules="[
            val => !!val || 'Sender is required',
            val => isValidEmail(val) || 'Invalid email'
          ]"
        />
      </div>

      <div class="col-12 col-md-6">
        <q-input
            v-model="localData.subject"
            outlined
            label="Subject *"
            :rules="[val => !!val || 'Subject is required']"
        />
      </div>

      <div class="col-12">
        <q-select
            v-model="localData.to"
            outlined
            use-input
            use-chips
            multiple
            input-debounce="0"
            new-value-mode="add-unique"
            label="Recipients *"
            :rules="[
            val => val.length > 0 || 'Add at least one recipient',
            val => val.every((email: string) => isValidEmail(email)) || 'One or more emails are invalid'
          ]"
        >
          <template v-slot:no-option>
            <q-item>
              <q-item-section class="text-grey">
                Type an email address and press Enter
              </q-item-section>
            </q-item>
          </template>
        </q-select>
      </div>

      <div class="col-12 col-md-6">
        <q-select
            v-model="localData.cc"
            outlined
            use-input
            use-chips
            multiple
            input-debounce="0"
            new-value-mode="add-unique"
            label="CC"
            :rules="[
            val => val.every((email: string) => isValidEmail(email)) || 'One or more emails are invalid'
          ]"
        >
          <template v-slot:no-option>
            <q-item>
              <q-item-section class="text-grey">
                Type an email address and press Enter
              </q-item-section>
            </q-item>
          </template>
        </q-select>
      </div>

      <div class="col-12 col-md-6">
        <q-select
            v-model="localData.bcc"
            outlined
            use-input
            use-chips
            multiple
            input-debounce="0"
            new-value-mode="add-unique"
            label="BCC"
            :rules="[
            val => val.every((email: string) => isValidEmail(email)) || 'One or more emails are invalid'
          ]"
        >
          <template v-slot:no-option>
            <q-item>
              <q-item-section class="text-grey">
                Type an email address and press Enter
              </q-item-section>
            </q-item>
          </template>
        </q-select>
      </div>

      <div class="col-12">
        <q-input
            v-model="localData.template"
            outlined
            autogrow
            type="textarea"
            label="HTML Template *"
            hint="Use HTML with {{title}} and {{message}} as placeholders"
            :rules="[val => !!val || 'HTML template is required']"
        />
      </div>

      <div class="col-12">
        <q-toggle
            v-model="localData.attachments"
            color="primary"
            label="Allow Attachments"
        />
      </div>

      <div class="col-12">
        <q-expansion-item
            default-opened
            group="email-config"
            icon="settings"
            label="SMTP Settings"
            header-class="text-primary"
            expand-icon-class="text-primary"
        >
          <q-card>
            <q-card-section>
              <div class="row q-col-gutter-md">
                <div class="col-12 col-md-6">
                  <q-input
                      v-model="localData.smtp.host"
                      outlined
                      label="SMTP Server *"
                      :rules="[val => !!val || 'SMTP server is required']"
                  />
                </div>

                <div class="col-12 col-md-6">
                  <q-input
                      v-model.number="localData.smtp.port"
                      outlined
                      type="number"
                      label="SMTP Port *"
                      :rules="[val => !!val || 'SMTP port is required']"
                  />
                </div>

                <div class="col-12 col-md-6">
                  <q-toggle
                      v-model="localData.smtp.secure"
                      color="primary"
                      label="Secure Connection (SSL/TLS)"
                  />
                </div>

                <div class="col-12 col-md-6">
                  <q-input
                      v-model="localData.smtp.auth.user"
                      outlined
                      label="SMTP User *"
                      :rules="[val => !!val || 'SMTP user is required']"
                  />
                </div>

                <div class="col-12 col-md-6">
                  <q-input
                      v-model="localData.smtp.auth.password"
                      outlined
                      type="password"
                      label="SMTP Password *"
                      :rules="[val => !!val || 'SMTP password is required']"
                  />
                </div>
              </div>
            </q-card-section>
          </q-card>
        </q-expansion-item>
      </div>

      <div class="col-12">
        <q-expansion-item
            group="email-config"
            icon="visibility"
            label="Email Preview"
            header-class="text-primary"
            expand-icon-class="text-primary"
        >
          <q-card>
            <q-card-section>
              <div class="email-preview q-pa-md">
                <div class="text-subtitle2 q-mb-sm">From: {{ localData.from }}</div>
                <div class="text-subtitle2 q-mb-sm">Subject: {{ localData.subject }}</div>
                <div class="text-subtitle2 q-mb-sm">
                  To:
                  <DetailChip
                    v-for="email in localData.to"
                    :key="email"
                    :label="email"
                    color="primary"
                    dense
                  />
                </div>
                <q-separator class="q-my-md"/>
                <div class="email-content" v-html="previewEmail"></div>
              </div>
            </q-card-section>
          </q-card>
        </q-expansion-item>
      </div>
    </div>
  </div>
</template>

<style scoped>
.email-preview {
  background: var(--mapex-surface-bg);
  border: 1px solid var(--mapex-card-border);
  border-radius: var(--mapex-radius-xs);
}

.email-content {
  max-height: 300px;
  overflow-y: auto;
  padding: 10px;
  background: var(--mapex-surface-elevated);
  border-radius: var(--mapex-radius-xs);
}
</style>