<script setup lang="ts">
defineOptions({
  name: 'TeamsChannel'
});

import type { TeamsConfig } from '@components/cardBody/notificationBodyDetails/interfaces';
const props = defineProps<{ channel: TeamsConfig }>();
</script>

<template>
  <q-list dense>
    <q-item>
      <q-item-section>
        <div class="text-caption text-grey-7">Team</div>
        <div class="text-body2">{{ props.channel.teamName }}</div>
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
