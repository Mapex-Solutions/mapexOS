<script setup lang="ts">
defineOptions({
  name: 'CardBodyDetails'
});

import type { CarbBodyDetailsProps } from './interfaces';

import { AppTooltip } from '@components/tooltips';

const props = defineProps<CarbBodyDetailsProps>();
</script>

<template>
  <q-card-section class="q-pa-md rounded-borders q-mb-md" :class="props?.container?.color || 'bg-grey-2'">
    <!-- Section Title -->
    <div class="row">
      <div class="col-12 text-subtitle2 text-weight-medium q-mb-sm">
        {{ props.title }}
        <q-icon
            v-if="props.tenantName"
            size="18px"
            class="text-grey-6 q-ml-xs"
            name="info"
        >
          <AppTooltip :content="props.tenantName" />
        </q-icon>
      </div>

    </div>

    <!-- Fields Grid -->
    <div class="row q-col-gutter-sm">
      <div
          v-for="(field, index) in props.items"
          :key="index"
          :class="`col-${field.cols ?? 12}`"
      >
        <!-- Icon Field (icon + value?) -->
        <div v-if="field.type === 'icon'">
          <div class="text-caption text-grey-7">{{ field.name }}</div>
          <q-item dense class="rounded-borders q-pa-sm">
            <q-item-section class="q-pr-xs">
              <q-icon
                  size="sm"
                  :name="field.icon"
                  :color="field.iconColor ?? 'grey-7'"
              />
            </q-item-section>
            <q-item-section v-if="field.value">
              <div class="text-h6" :class="field.color ? `text-${field.color}-10` : 'text-grey-8'">
                {{ String(field.value) }}
                <AppTooltip v-if="field.tooltip" :content="field.tooltip" />
              </div>
            </q-item-section>
          </q-item>
        </div>

        <!-- Mini cards with description + icon and value -->
        <q-card v-else-if="field.type === 'card'" flat bordered class="rounded-borders q-pa-sm" :class="field.color">
          <div class="text-subtitle2 text-teal-9 text-center">
            {{ field.name }}
          </div>
          <div class="row items-center justify-center q-mt-xs">
            <q-icon size="23px" class="q-mr-xs" :name="field.icon" :color="field.iconColor" />
            <span class="text-h6 text-teal-10">
            {{ String(field.value) }}
          </span>
          </div>
        </q-card>

        <!-- Simple Field (name + chip value) -->
        <div v-else-if="field.type === 'chip'">
          <div class="text-caption text-grey-7">{{ field.name }}</div>
          <q-chip
              dense
              outline
              text-color="white"
              :size="field.size || 'md'"
              :color="field.color || 'grey-7'"
          >
            {{ String(field.value) }}
            <AppTooltip v-if="field.tooltip" :content="field.tooltip" />
          </q-chip>

        </div>

        <!-- Simple Field (name + value) -->
        <div v-else-if="field.type === 'iconsGroup'">
          <div class="text-caption text-grey-7">{{ field.name }}</div>
          <q-icon
              v-for="(iconItem, index) in field?.icons || []"
              size="sm"
              :key="index"
              :name="iconItem.icon"
              :color="iconItem.color ?? 'grey-7'"
          >
            <AppTooltip v-if="iconItem.tooltip" :content="iconItem.tooltip" />
          </q-icon>
        </div>

        <!-- Simple Field (name + value) -->
        <div v-else>
          <div class="text-caption text-grey-7">{{ field.name }}</div>
          <div class="text-body2">
            {{ String(field.value) }}
            <AppTooltip v-if="field.tooltip" :content="field.tooltip" />
          </div>
        </div>
      </div>
    </div>
  </q-card-section>
</template>