<script setup lang="ts">
defineOptions({
  name: 'AuditEventRow'
});

import type { AuditLogProps } from './interfaces';

import { ref } from 'vue';
import { JsonDrawer } from '@components/drawers';
import { DetailChip } from '@components/chips';

const props = defineProps<{ event: AuditLogProps }>();
const jsonDrawerOpen = ref(false);

// --- Icon & color per type ---
const typeMap: Record<AuditLogProps['type'], { icon: string; color: string }> = {
  userLog:        { icon: 'person',                  color: 'blue-5' },
  dataSource:     { icon: 'settings_input_antenna',  color: 'teal-5' },
  assets:         { icon: 'devices',                 color: 'indigo-5' },
  payloadHandler: { icon: 'memory',                  color: 'amber-5' },
  businessRule:   { icon: 'gavel',                   color: 'deep-orange-5' },
  triggers:       { icon: 'flash_on',                color: 'purple-5' },
  ruleTemplate:   { icon: 'rule',                    color: 'lime-5' },
  users:          { icon: 'group',                   color: 'green-5' },
  customers:      { icon: 'domain',                  color: 'cyan-5' }
};

// --- Human-readable labels for each type ---
const typeLabels: Record<AuditLogProps['type'], string> = {
  userLog:        'User Log',
  dataSource:     'Data Source',
  assets:         'Assets',
  payloadHandler: 'Payload Handler',
  businessRule:   'Business Rule',
  triggers:       'Triggers',
  ruleTemplate:   'Rule Template',
  users:          'Users',
  customers:      'Customers'
};

function getTypeLabel(type: AuditLogProps['type']): string {
  return typeLabels[type] || type;
}

// --- Action → Color map ---
const actionColorMap: Record<AuditLogProps['action'], string> = {
  Create: 'green-6',
  Update: 'blue-6',
  Edit:   'orange-6',
  Delete: 'red-6'
};

// --- Status → Color map ---
const statusColorMap: Record<AuditLogProps['status'], string> = {
  success: 'green-6',
  failure: 'red-6',
  warning: 'orange-6',
  info:    'blue-6'
};

// Card classes by action
const getCardClass   = () => `event-card--${props.event.action.toLowerCase()}`;
const getBorderClass = () => `event-card__border--${props.event.action.toLowerCase()}`;
const getActionColor = () => actionColorMap[props.event.action];

function formatDate(iso: string): string {
  return new Date(iso).toLocaleString('default', {
    day: '2-digit', month: 'short',
    year: 'numeric', hour: '2-digit', minute: '2-digit'
  });
}
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

      <!-- TYPE -->
      <div class="col-12 col-sm-4 col-md-2 q-mb-sm">
        <div class="row items-center">
          <q-avatar
            text-color="white"
            size="40px"
            class="q-mr-sm"
            :color="typeMap[event.type]?.color || 'grey-5'"
            :icon="typeMap[event.type]?.icon || 'assignment'"
          />
          <div class="column">
            <div class="text-caption text-grey-5 text-weight-medium">Type</div>
            <div class="text-subtitle2 text-weight-medium text-grey-9">
              {{ getTypeLabel(event.type) }}
            </div>
          </div>
        </div>
      </div>

      <!-- ACTOR -->
      <div class="col-12 col-sm-6 col-md-2 text-center q-mb-sm">
        <div class="text-caption text-grey-5 text-weight-medium">Actor</div>
        <div class="text-body2 text-weight-medium text-grey-8">
          {{ event.actor }}
        </div>
      </div>

      <!-- ACTION -->
      <div class="col-6 col-sm-6 col-md-2 text-center q-mb-sm">
        <div class="text-caption text-grey-5 text-weight-medium">Action</div>
        <DetailChip
          :label="event.action.toUpperCase()"
          :color="getActionColor() as any"
          size="sm"
          outline
          rounded
        />
      </div>

      <!-- RESOURCE -->
      <div class="col-12 col-sm-6 col-md-2 text-center q-mb-sm">
        <div class="text-caption text-grey-5 text-weight-medium">Resource</div>
        <div class="text-body2 text-weight-medium text-grey-8">
          {{ event.resource }}
        </div>
      </div>

      <!-- STATUS -->
      <div class="col-6 col-sm-3 col-md-1 text-center q-mb-sm">
        <div class="text-caption text-grey-5 text-weight-medium">Status</div>
        <DetailChip
          :label="event.status.toUpperCase()"
          :color="(statusColorMap[event.status] || 'grey-6') as any"
          size="sm"
          outline
          rounded
        />
      </div>

      <!-- TIMESTAMP -->
      <div class="col-12 col-sm-12 col-md-3 text-right q-mb-sm">
        <div class="text-caption text-grey-5 text-weight-medium">When</div>
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
    :jsonData="event"
    :subtitle="`${getTypeLabel(event.type)} • ${formatDate(event.created)}`"
  />
</template>

<style scoped lang="scss">
.q-avatar {
  box-shadow: var(--mapex-shadow-sm);
}
</style>
