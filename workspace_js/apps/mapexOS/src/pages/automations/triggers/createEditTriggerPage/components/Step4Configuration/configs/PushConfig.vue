<script setup lang="ts">
defineOptions({
  name: 'PushConfig'
});

import type { Trigger } from '../../../interfaces';
import { ref, watch } from 'vue';

import { useCreateEditTriggerTranslations } from '@composables/i18n/pages/automations/triggers/createEditTrigger';

const props = defineProps<{ modelValue: Trigger }>();
const emit = defineEmits<{ 'update:modelValue': [value: Trigger] }>();

const t = useCreateEditTriggerTranslations();

const config = ref(props.modelValue.config || {});

watch(config, (newConfig) => {
  emit('update:modelValue', { ...props.modelValue, config: newConfig });
}, { deep: true });
</script>

<template>
  <div class="text-center q-pa-lg">
    <q-icon name="construction" size="3rem" color="warning" />
    <div class="text-h6 q-mt-md">{{ t.steps.step4.comingSoon.value }}</div>
    <div class="text-body2 text-grey-7">{{ t.steps.step4.comingSoonDetail.value }}</div>
  </div>
</template>
