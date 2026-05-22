<script setup lang="ts">
defineOptions({
  name: 'TriggerEventRow'
});

import type { RawTriggerProps } from './interfaces';
import { ref } from 'vue';
import { JsonDrawer } from '@components/drawers';
import { DetailChip } from '@components/chips';

const props = defineProps<{
  event: RawTriggerProps
}>();

defineEmits<{
  showValues: [event: RawTriggerProps];
  viewDetails: [event: RawTriggerProps];
}>();

const jsonDrawerOpen = ref(false);

const getCardClass = () => `event-card--${props.event.status}`;
const getBorderClass = () => `event-card__border--${props.event.status}`;

const getTriggerTypeColor = (type: string) => {
  const map: Record<string, string> = {
    HTTP: 'blue-6',
    Notification: 'teal-6',
    Incident: 'red-6',
    Task: 'orange-6',
    Workflow: 'indigo-6',
    MQTT: 'purple-6',
  };
  return map[type] || 'grey-6';
};

const getTriggerTypeIcon = (type: string) => {
  const icons: Record<string, string> = {
    HTTP: 'language',
    Notification: 'notifications',
    Incident: 'warning',
    Task: 'check_circle',
    Workflow: 'devices',
    MQTT: 'router',
  };
  return icons[type] || 'help';
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

<template>
  <q-card
      flat
      bordered
      class="q-card--hover hoverable cursor-pointer"
      :class="getCardClass()"
      @click.stop="jsonDrawerOpen = true"
  >
    <q-card-section class="row items-center justify-between q-pa-md">
      <!-- TRIGGER TYPE -->
      <div class="col-2 q-mb-sm">
        <div class="row items-center">
          <q-avatar
              text-color="white"
              size="40px"
              class="q-mr-sm"
              :color="getTriggerTypeColor(props.event.triggerType)"
              :icon="getTriggerTypeIcon(props.event.triggerType)"
          />
          <div class="column">
            <div class="text-caption text-grey-5 text-weight-medium">
              Trigger Type
            </div>
            <div class="text-subtitle2 text-weight-medium text-grey-9">
              {{ props.event.triggerType }}
            </div>
          </div>
        </div>
      </div>

      <!-- TENANT (optional) -->
      <div v-if="props.event.tenantName" class="col-12 col-sm-6 col-md-2 text-center q-mb-sm">
        <div class="text-caption text-grey-5 text-weight-medium">Tenant</div>
        <div class="text-body2 text-weight-medium text-grey-8">
          {{ props.event.tenantName || '—' }}
        </div>
      </div>

      <!-- TRIGGER NAME -->
      <div class="col-6 col-sm-3 col-md-2 text-center q-mb-sm">
        <div class="text-caption text-grey-5 text-weight-medium">
          Trigger Name
        </div>
        <div class="text-body2 text-weight-medium text-grey-8">
          {{ props.event.triggerName }}
        </div>
      </div>

      <!-- STATUS -->
      <div class="col-12 col-sm-4 col-md-1 text-center q-mb-sm">
        <div class="text-caption text-grey-5 text-weight-medium">Status</div>
        <DetailChip
          :label="props.event.status.toUpperCase()"
          :color="(props.event.status === 'success' ? 'green-6' : 'red-6') as any"
          size="sm"
          outline
          rounded
        />
      </div>

      <!-- CREATED -->
      <div class="col-12 col-sm-12 col-md-2 text-right q-mb-sm">
        <div class="text-caption text-grey-5 text-weight-medium">Created</div>
        <div class="text-body2 text-weight-medium text-grey-8">
          {{ formatDate(props.event.created) }}
        </div>
      </div>
    </q-card-section>

    <div class="event-card__border" :class="getBorderClass()"></div>
  </q-card>

  <JsonDrawer
      v-model:show="jsonDrawerOpen"
      :editable="false"
      title="Trigger Details"
      :jsonData="props.event"
      :subtitle="`${props.event.triggerType} • ${formatDate(props.event.created)}`"
  />
</template>

<style scoped lang="scss">
.q-avatar {
  box-shadow: var(--mapex-shadow-sm);
}
</style>
