<script setup lang="ts">
defineOptions({
  name: 'ChannelBaseInfo'
});

import type { ChannelType, ChannelBaseProps } from './interfaces';

import { ref, computed, watch } from 'vue';
import { CHANNEL_TYPES } from './constants';

// Props and Emits
const props = defineProps<{
  modelValue: ChannelBaseProps;
}>();

const emit = defineEmits<{
  (e: 'update:modelValue', value: ChannelBaseProps): void;
}>();

// Local state
const localData = ref<ChannelBaseProps>({ ...props.modelValue });

// Computed properties
const isActive = computed({
  get: () => localData.value.status === 'Active',
  set: (value) => {
    localData.value.status = value ? 'Active' : 'Inactive';
    emit('update:modelValue', localData.value);
  },
});

// Methods
/**
 * Handles the change event when channel type is updated
 * Updates the icon based on the selected type and emits the change to parent component
 * @param {ChannelType} type - The selected channel type
 */
function handleTypeChange(type: ChannelType) {
  // Updates the icon when type changes
  const typeOption = CHANNEL_TYPES.find(t => t.value === type);
  if (typeOption) {
    localData.value.icon = typeOption.icon;
  }

  // Emits the event to parent component
  emit('update:modelValue', localData.value);
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
  <div class="row q-col-gutter-md">
    <div class="col-12">
      <q-select
          v-model="localData.type"
          outlined
          emit-value
          map-options
          label="Channel Type *"
          option-label="label"
          option-value="value"
          :options="CHANNEL_TYPES"
          :rules="[val => !!val || 'Notification type is required']"
          @update:model-value="handleTypeChange"
      >
        <template v-slot:option="scope">
          <q-item v-bind="scope.itemProps">
            <q-item-section avatar>
              <q-icon :name="scope.opt.icon"/>
            </q-item-section>
            <q-item-section>
              <q-item-label>{{ scope.opt.label }}</q-item-label>
              <q-item-label caption>{{ scope.opt.description }}</q-item-label>
            </q-item-section>
          </q-item>
        </template>
      </q-select>
    </div>

    <div class="col-12 col-md-8">
      <q-input
          v-model="localData.name"
          outlined
          label="Channel Name *"
          :rules="[val => !!val || 'Name is required']"
      />
    </div>

    <div class="col-12 col-md-4">
      <q-toggle
          v-model="isActive"
          color="positive"
          label="Active Status"
      />
    </div>

    <div class="col-12">
      <q-input
          v-model="localData.description"
          outlined
          autogrow
          type="textarea"
          label="Description"
      />
    </div>
  </div>
</template>