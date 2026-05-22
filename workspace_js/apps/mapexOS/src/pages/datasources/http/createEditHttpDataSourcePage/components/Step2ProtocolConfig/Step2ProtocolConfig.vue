<template>
  <div class="row q-col-gutter-md">
    <!-- Info Banner -->
    <div class="col-12">
      <q-banner dense class="bg-blue-1 text-blue-9 rounded-borders">
        <template #avatar>
          <q-icon name="info" color="blue-6" />
        </template>
        <div class="text-caption">
          <strong>{{ t.protocolConfig.gatewayBannerTitle.value }}</strong> {{ t.protocolConfig.gatewayBannerBody.value }}
        </div>
      </q-banner>
    </div>

    <!-- Mode Selection (Only Push) -->
    <div class="col-12">
      <q-select
        v-model="localData.mode"
        outlined
        dense
        label="Mode *"
        class="rounded-borders"
        :options="MODE_OPTIONS"
        :rules="[(val) => !!val || 'Mode is required']"
        @update:model-value="updateValue"
      >
        <template v-slot:prepend>
          <q-icon name="upload" color="primary"/>
        </template>
      </q-select>
      <div class="text-caption text-grey-7 q-mt-xs">
        <q-icon name="info" size="xs" />
        Push mode: Devices send data to this endpoint via HTTP POST requests.
      </div>
    </div>

    <!-- Protocol Display (Always HTTP) -->
    <div class="col-12">
      <q-input
        model-value="HTTP"
        outlined
        dense
        readonly
        label="Protocol"
        class="rounded-borders"
        hint="HTTP/HTTPS protocol is fixed for this gateway type"
      >
        <template v-slot:prepend>
          <q-icon name="http" color="primary"/>
        </template>
      </q-input>
    </div>
  </div>
</template>

<script setup lang="ts">
defineOptions({
  name: 'Step2ProtocolConfig'
});

import type { StepEmits, StepProps } from '../../interfaces/httpDataSource.interface';

import { reactive, watch } from 'vue';

import { useHttpDataSourceCreateEditTranslations } from '@composables/i18n/pages/datasources/http';

import { MODE_OPTIONS } from '../../constants/httpDataSourceConstants';

const props = defineProps<StepProps>();
const emit = defineEmits<StepEmits>();

const t = useHttpDataSourceCreateEditTranslations();

const localData = reactive({
  mode: props.modelValue.mode || null,
  protocol: 'HTTP', // Always HTTP for this gateway
});

watch(() => props.modelValue, (newVal) => {
  localData.mode = newVal.mode || null;
  // Protocol is always HTTP, no need to watch
}, { deep: true });

function updateValue() {
  emit('update:modelValue', {
    ...localData,
    protocol: 'HTTP', // Ensure protocol is always HTTP
  });
}
</script>

<style scoped>
.rounded-borders {
  border-radius: var(--mapex-radius-md);
}
</style>
