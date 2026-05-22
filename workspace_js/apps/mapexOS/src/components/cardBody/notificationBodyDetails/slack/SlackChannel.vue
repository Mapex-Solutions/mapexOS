<script setup lang="ts">
defineOptions({
  name: 'SlackChannel'
});

import type { SlackConfig } from '@components/cardBody/notificationBodyDetails/interfaces';
const props = defineProps<{ channel: SlackConfig }>();
</script>

<template>
  <q-list dense>
    <q-item>
      <q-item-section>
        <div class="text-caption text-grey-7">Workspace</div>
        <div class="text-body2">{{ props.channel.workspace }}</div>
      </q-item-section>
    </q-item>
    <q-item>
      <q-item-section>
        <div class="text-caption text-grey-7">Channel</div>
        <div class="text-body2">{{ props.channel.channelName }}</div>
      </q-item-section>
    </q-item>
  </q-list>
</template>
