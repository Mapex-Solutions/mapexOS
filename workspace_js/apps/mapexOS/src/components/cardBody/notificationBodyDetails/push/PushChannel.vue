<script setup lang="ts">
defineOptions({
  name: 'PushChannel'
});

import type { PushConfig } from '@components/cardBody/notificationBodyDetails/interfaces';
const props = defineProps<{ channel: PushConfig }>();
</script>

<template>
  <q-list dense>
    <q-item>
      <q-item-section>
        <div class="text-caption text-grey-7">App</div>
        <div class="text-body2">{{ props.channel.appName }}</div>
      </q-item-section>
    </q-item>
    <q-item v-if="props.channel.deviceCount != null">
      <q-item-section>
        <div class="text-caption text-grey-7">Devices</div>
        <div class="text-body2">{{ props.channel.deviceCount }}</div>
      </q-item-section>
    </q-item>
  </q-list>
</template>
