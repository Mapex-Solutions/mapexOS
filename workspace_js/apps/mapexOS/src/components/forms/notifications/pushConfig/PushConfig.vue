<script setup lang="ts">
defineOptions({
  name: 'PushConfig'
});

import type { ChannelPushProps } from './interfaces';

import { ref, watch } from 'vue';
import { SERVICE_PROVIDERS, PRIORITY_OPTIONS, SOUND_OPTIONS } from './constants';

// Props and Emits
const props = defineProps<{
  modelValue: ChannelPushProps;
}>();

const emit = defineEmits<{
  (e: 'update:modelValue', value: ChannelPushProps): void;
}>();

// Local state
const localData = ref<ChannelPushProps>({ ...props.modelValue });
const showApiKey = ref(false);

// Methods
/**
 * Gets the human-readable label for a service provider based on its value
 * @param {string} value - The service provider value to lookup
 * @returns {string} The provider label or the original value if not found
 */
function getProviderLabel(value: string): string {
  const provider = SERVICE_PROVIDERS.find(p => p.value === value);
  return provider ? provider.label : value;
}

/**
 * Gets the human-readable label for a priority level based on its value
 * @param {string} value - The priority value to lookup
 * @returns {string} The priority label or the original value if not found
 */
function getPriorityLabel(value: string): string {
  const priority = PRIORITY_OPTIONS.find(p => p.value === value);
  return priority ? priority.label : value;
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
  <div class="push-config">
    <div class="row q-col-gutter-md">
      <div class="col-12 col-md-6">
        <q-input
            v-model="localData.appName"
            outlined
            label="Application Name *"
            :rules="[val => !!val || 'Application name is required']"
        />
      </div>

      <div class="col-12 col-md-6">
        <q-input
            v-model.number="localData.deviceCount"
            outlined
            readonly
            type="number"
            label="Device Count"
            hint="This value is automatically updated by the system"
        />
      </div>

      <div class="col-12 col-md-6">
        <q-select
            v-model="localData.serviceProvider"
            outlined
            emit-value
            map-options
            label="Service Provider *"
            :options="SERVICE_PROVIDERS"
            :rules="[val => !!val || 'Service provider is required']"
        >
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

      <div class="col-12 col-md-6">
        <q-select
            v-model="localData.priority"
            outlined
            emit-value
            map-options
            label="Priority"
            :options="PRIORITY_OPTIONS"
        >
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

      <div class="col-12">
        <q-input
            v-model="localData.apiKey"
            outlined
            label="API Key *"
            :rules="[val => !!val || 'API Key is required']"
            :type="showApiKey ? 'text' : 'password'"
        >
          <template v-slot:append>
            <q-icon
                class="cursor-pointer"
                :name="showApiKey ? 'visibility_off' : 'visibility'"
                @click="showApiKey = !showApiKey"
            />
          </template>
        </q-input>
      </div>

      <div class="col-12 col-md-6">
        <q-input
            v-model.number="localData.ttl"
            outlined
            type="number"
            label="TTL (seconds)"
            hint="Notification time to live in seconds"
        />
      </div>

      <div class="col-12 col-md-6">
        <q-input
            v-model="localData.clickAction"
            outlined
            label="Click Action"
            hint="Action to execute when user clicks the notification"
        />
      </div>

      <div class="col-12 col-md-6">
        <q-toggle
            v-model="localData.badge"
            label="Show Badge"
            color="primary"
            hint="Display notification counter on app icon"
        />
      </div>

      <div class="col-12 col-md-6">
        <q-select
            v-model="localData.sound"
            outlined
            emit-value
            label="Sound"
            :options="SOUND_OPTIONS"
        />
      </div>

      <div class="col-12">
        <q-card flat bordered class="bg-grey-1">
          <q-card-section>
            <div class="text-subtitle2">Additional Information</div>
            <p class="text-caption q-mt-sm">
              This notification channel will send push messages to
              <strong>{{ localData.deviceCount }}</strong> registered devices through the
              <strong>{{ getProviderLabel(localData.serviceProvider) }}</strong> provider.
            </p>
            <p class="text-caption">
              Notifications will have <strong>{{ getPriorityLabel(localData.priority) }}</strong>
              priority and will expire after <strong>{{ localData.ttl }}</strong> seconds if not delivered.
            </p>
          </q-card-section>
        </q-card>
      </div>
    </div>
  </div>
</template>