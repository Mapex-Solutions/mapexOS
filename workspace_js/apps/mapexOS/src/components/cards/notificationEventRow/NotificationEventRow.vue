<script setup lang="ts">
defineOptions({
  name: 'NotificationEventRow'
});

import type { RawNotificationProps } from './interfaces';
import { ref } from 'vue';
import { JsonDrawer } from '@components/drawers';
import { DetailChip } from '@components/chips';

const props = defineProps<{
  event: RawNotificationProps
}>();

const jsonDrawerOpen = ref(false);

const getCardClass = () => `event-card--${props.event.status}`;
const getBorderClass = () => `event-card__border--${props.event.status}`;

const getNotificationTypeColor = (type: string) => {
  const map: Record<string, string> = {
    slack: 'purple-6',
    teams: 'blue-6',
    email: 'grey-6',
    push: 'orange-6',
    telegram: 'cyan-6',
    webhook: 'indigo-6',
  };
  // normalize just in case
  const key = type.toLowerCase();
  return map[key] || 'grey-6';
};

const getNotificationTypeIcon = (type: string) => {
  const icons: Record<string, string> = {
    slack: 'mdi-slack',
    teams: 'mdi-microsoft-teams',
    email: 'mdi-email',
    push: 'mdi-bell-ring',
    telegram: 'mdi-send',
    webhook: 'mdi-webhook',
  };
  const key = type.toLowerCase();
  return icons[key] || 'mdi-bell-outline';
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
      <!-- NOTIFICATION TYPE -->
      <div class="col-2 q-mb-sm">
        <div class="row items-center">
          <q-avatar
              text-color="white"
              size="40px"
              class="q-mr-sm"
              :color="getNotificationTypeColor(props.event.notificationType)"
              :icon="getNotificationTypeIcon(props.event.notificationType)"
          />
          <div class="column">
            <div class="text-caption text-grey-5 text-weight-medium">
              Notification Type
            </div>
            <div class="text-subtitle2 text-weight-medium text-grey-9 text-capitalize">
              {{ props.event.notificationType }}
            </div>
          </div>
        </div>
      </div>

      <!-- TENANT (optional) -->
      <div v-if="props.event.tenantId" class="col-12 col-sm-6 col-md-2 text-center q-mb-sm">
        <div class="text-caption text-grey-5 text-weight-medium">Tenant</div>
        <div class="text-body2 text-weight-medium text-grey-8">
          {{ props.event.tenantId || '—' }}
        </div>
      </div>

      <!-- NOTIFICATION NAME -->
      <div class="col-6 col-sm-3 col-md-2 text-center q-mb-sm">
        <div class="text-caption text-grey-5 text-weight-medium">
          Notification Name
        </div>
        <div class="text-body2 text-weight-medium text-grey-8">
          {{ props.event.notificationName }}
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
      :subtitle="`${props.event.notificationType} • ${formatDate(props.event.created)}`"
  />
</template>

<style scoped lang="scss">
.q-avatar {
  box-shadow: var(--mapex-shadow-sm);
}
</style>
