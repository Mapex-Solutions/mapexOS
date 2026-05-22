<script setup lang="ts">
defineOptions({
  name: 'LakeHouseFrequency'
});

import { computed } from 'vue';

import type { LakeHouseConfigProps } from '@components/forms/lakeHouse';
import type { LakeHouseFrequency } from './interfaces';

import { DEFAULT_FREQUENCY_CONFIG } from './constants';

// this ref is automatically tied to `modelValue` + `update:modelValue`
const modelRef = defineModel<LakeHouseConfigProps>({
  default: () => ({
    frequency: DEFAULT_FREQUENCY_CONFIG,
  }),
});

const frequencyOptions = [
  { value: 'minute', label: 'Minute', icon: 'schedule' },
  { value: 'hour',   label: 'Hour',   icon: 'access_time' },
  { value: 'day',    label: 'Day',    icon: 'today' },
  { value: 'week',   label: 'Week',   icon: 'date_range' },
  { value: 'month',  label: 'Month',  icon: 'calendar_month' },
  { value: 'year',   label: 'Year',   icon: 'calendar_today' },
] as const;

const weekdayOptions = [
  { value: 'monday',    label: 'Monday',    short: 'Mon' },
  { value: 'tuesday',   label: 'Tuesday',   short: 'Tue' },
  { value: 'wednesday', label: 'Wednesday', short: 'Wed' },
  { value: 'thursday',  label: 'Thursday',  short: 'Thu' },
  { value: 'friday',    label: 'Friday',    short: 'Fri' },
  { value: 'saturday',  label: 'Saturday',  short: 'Sat' },
  { value: 'sunday',    label: 'Sunday',    short: 'Sun' },
];

// Computed properties
const showInterval = computed(() =>
  ['minute', 'hour', 'day'].includes(modelRef.value.frequency.type),
);

const showTime = computed(() =>
  ['day', 'week', 'month', 'year'].includes(modelRef.value.frequency.type),
);

const showWeekdays = computed(() =>
  modelRef.value.frequency.type === 'week',
);

const showDayOfMonth = computed(() =>
  modelRef.value.frequency.type === 'month',
);

const intervalSuffix = computed(() => {
  const type = modelRef.value.frequency.type;
  const interval = modelRef.value.frequency.interval || 1;

  switch (type) {
    case 'minute': {
      return interval === 1 ? 'minute' : 'minutes';
    }
    case 'hour': {
      return interval === 1 ? 'hour' : 'hours';
    }
    case 'day': {
      return interval === 1 ? 'day' : 'days';
    }
    default: {
      return '';
    }
  }
});

const frequencyDescription = computed(() => {
  const freq = modelRef.value.frequency;
  const interval = freq.interval || 1;

  switch (freq.type) {
    case 'minute': {
      return interval === 1
        ? 'Every minute'
        : `Every ${interval} minutes`;
    }

    case 'hour': {
      return interval === 1
        ? 'Every hour'
        : `Every ${interval} hours`;
    }

    case 'day': {
      const timeStr = freq.time ? ` at ${freq.time}` : '';
      return interval === 1
        ? `Daily${timeStr}`
        : `Every ${interval} days${timeStr}`;
    }

    case 'week': {
      const weekdaysStr = freq.weekdays?.length
        ? ` (${freq.weekdays.map((d: any) => weekdayOptions.find(w => w.value === d)?.short).join(', ')})`
        : '';
      const timeStr = freq.time ? ` at ${freq.time}` : '';
      return `Weekly${weekdaysStr}${timeStr}`;
    }

    case 'month': {
      const dayStr = freq.dayOfMonth ? ` on day ${freq.dayOfMonth}` : '';
      const timeStr = freq.time ? ` at ${freq.time}` : '';
      return `Monthly${dayStr}${timeStr}`;
    }

    case 'year': {
      const timeStr = freq.time ? ` at ${freq.time}` : '';
      return `Yearly${timeStr}`;
    }

    default: {
      return 'Frequency not defined';
    }
  }
});

const nextExecution = computed(() => {
  // This would calculate the next execution time based on current settings
  // For now, just return a placeholder
  const now = new Date();
  const next = new Date(now.getTime() + 24 * 60 * 60 * 1000); // Tomorrow
  return next.toLocaleString('en-US');
});

// Methods
function selectFrequency(type: LakeHouseFrequency['type']) {
  modelRef.value.frequency = {
    type,
    interval: ['minute', 'hour', 'day'].includes(type) ? 1 : undefined,
    time: ['day', 'week', 'month', 'year'].includes(type) ? '09:00' : undefined,
    weekdays: type === 'week' ? [] : undefined,
    dayOfMonth: type === 'month' ? 1 : undefined,
  };
}

function toggleWeekday(day: string) {
  if (!modelRef.value.frequency.weekdays) {
    modelRef.value.frequency.weekdays = [];
  }

  const weekdays = modelRef.value.frequency.weekdays;
  const index = weekdays.indexOf(day);

  if (index > -1) {
    weekdays.splice(index, 1);
  } else {
    weekdays.push(day);
  }
}

function isValidTime(time: string): boolean {
  if (!time) return true; // Optional field
  const timeRegex = /^([01]?[0-9]|2[0-3]):[0-5][0-9]$/;
  return timeRegex.test(time);
}
</script>

<template>
  <div class="frequency-config">
    <div class="row q-col-gutter-md">
      <!-- Frequency Type -->
      <div class="col-12">
        <div class="row q-col-gutter-sm">
          <div
            v-for="freq in frequencyOptions"
            :key="freq.value"
            class="col-6 col-sm-4 col-md-2"
          >
            <q-card
              flat
              bordered
              class="frequency-card cursor-pointer transition-all text-center"
              :class="modelRef.frequency.type === freq.value
                  ? 'frequency-card--selected'
                  : 'frequency-card--default'"
              @click="selectFrequency(freq.value)"
            >
              <q-card-section class="q-py-md">
                <q-icon size="24px" class="q-mb-xs" :name="freq.icon" />
                <div class="text-body2 text-weight-medium">{{ freq.label }}</div>
              </q-card-section>
            </q-card>
          </div>
        </div>
      </div>

      <!-- Custom Interval (for applicable frequencies) -->
      <div v-if="showInterval" class="col-12 col-md-6">
        <q-input
          v-model.number="modelRef.frequency.interval"
          outlined
          type="number"
          min="1"
          label="Interval"
          hint="Set a custom interval"
          :suffix="intervalSuffix"
          :rules="[val => val > 0 || 'Interval must be greater than 0']"
        />
      </div>

      <!-- Time Selection -->
      <div v-if="showTime" class="col-12 col-md-6">
        <q-input
          v-model="modelRef.frequency.time"
          outlined
          mask="##:##"
          placeholder="14:30"
          label="Time"
          hint="Execution time (24h format)"
          :rules="[val => isValidTime(val) || 'Invalid time format (HH:mm)']"
        >
          <template v-slot:append>
            <q-icon name="access_time" class="cursor-pointer">
              <q-popup-proxy cover transition-show="scale" transition-hide="scale">
                <q-time v-model="modelRef.frequency.time" format24h>
                  <div class="row items-center justify-end">
                    <q-btn v-close-popup flat label="Close" color="primary"/>
                  </div>
                </q-time>
              </q-popup-proxy>
            </q-icon>
          </template>
        </q-input>
      </div>

      <!-- Weekdays Selection (for weekly frequency) -->
      <div v-if="showWeekdays" class="col-12">
        <div class="text-subtitle2 q-mb-md">Weekdays</div>
        <div class="row q-col-gutter-xs">
          <div
            v-for="day in weekdayOptions"
            :key="day.value"
            class="col-auto"
          >
            <q-btn
              round
              size="sm"
              class="weekday-btn"
              :label="day.short"
              :color="modelRef.frequency.weekdays?.includes(day.value) ? 'primary' : 'grey-4'"
              :text-color="modelRef.frequency.weekdays?.includes(day.value) ? 'white' : 'grey-8'"
              @click="toggleWeekday(day.value)"
            />
          </div>
        </div>
        <div class="text-caption text-grey-6 q-mt-xs">
          Select weekdays for export
        </div>
      </div>

      <!-- Day of Month (for monthly frequency) -->
      <div v-if="showDayOfMonth" class="col-12 col-md-6">
        <q-input
          v-model.number="modelRef.frequency.dayOfMonth"
          outlined
          type="number"
          min="1"
          max="31"
          label="Day of Month"
          hint="Day of month for execution (1-31)"
          :rules="[
            val => val >= 1 && val <= 31 || 'Day must be between 1 and 31'
          ]"
        />
      </div>
    </div>

    <!-- Frequency Preview -->
    <div class="q-mt-lg">
      <q-card flat bordered class="bg-blue-1">
        <q-card-section>
          <div class="text-subtitle2 text-primary q-mb-sm">
            <q-icon name="schedule" class="q-mr-xs"/>
            Frequency Summary
          </div>
          <div class="text-body2 text-grey-8">
            {{ frequencyDescription }}
          </div>
          <div v-if="nextExecution" class="text-caption text-grey-6 q-mt-xs">
            Next estimated execution: {{ nextExecution }}
          </div>
        </q-card-section>
      </q-card>
    </div>
  </div>
</template>

<style scoped lang="scss">
.frequency-card {
  min-height: 80px;
  border-radius: var(--mapex-radius-md);

  &--default {
    border: 2px solid var(--mapex-card-border);

    &:hover {
      border-color: var(--mapex-card-hover-border);
      box-shadow: var(--mapex-shadow-sm);
    }
  }

  &--selected {
    border: 2px solid var(--q-primary);
    background-color: rgba(var(--q-primary-rgb), 0.05);
  }
}

.weekday-btn {
  width: 32px;
  height: 32px;
  min-width: 32px;
}

.transition-all {
  transition: var(--mapex-transition-base);
}
</style>
