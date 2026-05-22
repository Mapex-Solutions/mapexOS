<script setup lang="ts">
defineOptions({
  name: 'EmailChannel'
});

import type { EmailConfig } from '@components/cardBody/notificationBodyDetails/interfaces';

import { computed } from 'vue';
import { AppTooltip } from '@components/tooltips';
const props = defineProps<{ channel: EmailConfig }>();

const displayList  = computed(() => props.channel.to.slice(0, 2));
const extraList    = computed(() => props.channel.to.slice(2));
const extraCount   = computed(() => extraList.value.length);
</script>

<template>
  <q-list dense>
    <q-item>
      <q-item-section>
        <div class="text-caption text-grey-7">From</div>
        <div class="text-body2">{{ props.channel.from }}</div>
      </q-item-section>
    </q-item>

    <q-item>
      <q-item-section>
        <div class="text-caption text-grey-7">To</div>
        <div class="row items-center">
          <span
            v-for="(addr, i) in displayList"
            :key="i"
            class="text-body2 q-mr-xs"
          >
            {{ addr }}
          </span>
          <q-chip
            v-if="extraCount"
            dense
            size="sm"
            class="bg-grey-3 text-grey-8"
          >
            +{{ extraCount }}
            <AppTooltip :content="extraList.join(', ')" />
          </q-chip>
        </div>
      </q-item-section>
    </q-item>
  </q-list>
</template>
