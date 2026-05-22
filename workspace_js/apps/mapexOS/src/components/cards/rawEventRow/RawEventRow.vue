<template>
  <q-card
      flat
      bordered
      class="q-card--hover hoverable cursor-pointer"
      :class="getCardClass()"
      @click.stop="jsonDrawerOpen = true"
  >
    <!-- WRAP is the default -->
    <q-card-section class="row items-center justify-between q-pa-md">

      <!-- ASSET -->
      <div class="col-4 q-mb-sm">
        <div class="row items-center">
          <q-avatar
              text-color="white"
              size="40px"
              class="q-mr-sm"
              :color="getAssetIconColor()"
              :icon="event.asset.icon"
          />
          <div class="column">
            <div class="text-caption text-grey-5 text-weight-medium">Asset</div>
            <div class="text-subtitle2 text-weight-medium text-grey-9">
              {{ event.asset.name }}
            </div>
            <div class="text-caption text-grey-6">
              {{ event.asset.description }}
            </div>
          </div>
        </div>
      </div>

      <!-- TENANT ID -->
      <div class="col-12 col-sm-6 col-md-2 text-center q-mb-sm">
        <div class="text-caption text-grey-5 text-weight-medium">Tenant</div>
        <div class="text-body2 text-weight-medium text-grey-8">
          {{ 'Mapex ABC 123 456' }}
        </div>
      </div>

      <!-- EVENT TYPE -->
      <div class="col-12 col-sm-6 col-md-2 text-center q-mb-sm">
        <div class="text-caption text-grey-5 text-weight-medium">Event Type</div>
        <div class="text-body2 text-weight-medium text-grey-8">
          {{ event.type }}
        </div>
      </div>

      <!-- STATUS -->
      <div class="col-6 col-sm-3 col-md-1 text-center q-mb-sm">
        <div class="text-caption text-grey-5 text-weight-medium">Status</div>
        <DetailChip
          :label="event.status.toUpperCase()"
          :color="getStatusColor() as any"
          size="sm"
          outline
          rounded
        />
      </div>

      <!-- PROTOCOL -->
      <div class="col-6 col-sm-3 col-md-1 text-center q-mb-sm">
        <div class="text-caption text-grey-5 text-weight-medium">Protocol</div>
        <DetailChip
          :label="event.protocol"
          :color="getProtocolColor() as any"
          size="sm"
          outline
          rounded
        />
      </div>

      <!-- CREATED -->
      <div class="col-12 col-sm-12 col-md-2 text-right q-mb-sm">
        <div class="text-caption text-grey-5 text-weight-medium">Created</div>
        <div class="text-body2 text-weight-medium text-grey-8">
          {{ formatDate(event.created) }}
        </div>
      </div>

    </q-card-section>
    <div class="event-card__border" :class="getBorderClass()"></div>
  </q-card>

  <!-- JSON Drawer -->
  <JsonDrawer
    v-model:show="jsonDrawerOpen"
    :editable="false"
    title="Audit Log Details"
    :jsonData="props.event"
    :subtitle="`${props.event.type} • ${formatDate(props.event.created)}`"
  />
</template>

<script setup lang="ts">
defineOptions({
  name: 'RawEventRow'
});

import type { RawEventProps } from './interfaces';

import { ref } from 'vue';
import { JsonDrawer } from '@components/drawers';
import { DetailChip } from '@components/chips';

const props = defineProps<{
  event: RawEventProps
}>();

defineEmits<{
  showValues: [event: RawEventProps]
  viewDetails: [event: RawEventProps]
}>();

const jsonDrawerOpen = ref(false);

const getCardClass = () => {
  return `event-card--${props.event.status}`;
};

const getBorderClass = () => {
  return `event-card__border--${props.event.status}`;
};

const getStatusColor = () => {
  const colors = {
    high: 'red-6',
    medium: 'orange-6',
    low: 'green-6',
  };
  return colors[props.event.status];
};

const getAssetIconColor = () => {
  const colors = {
    high: 'red-5',
    medium: 'orange-5',
    low: 'green-5',
  };
  return colors[props.event.status];
};

const getProtocolColor = () => {
  const protocolColors: Record<string, string> = {
    'MQTT': 'purple-6',
    'HTTP': 'blue-6',
    'TCP': 'teal-6',
    'UDP': 'indigo-6',
  };
  return protocolColors[props.event.protocol] || 'grey-6';
};

const formatDate = (dateString: string) => {
  const date = new Date(dateString);
  return date.toLocaleDateString('en-US', {
    day: '2-digit',
    month: 'short',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  });
};

</script>

<style scoped lang="scss">
.q-avatar {
  box-shadow: var(--mapex-shadow-sm);
}
</style>