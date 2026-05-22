<template>
  <div class="row q-col-gutter-md">
    <!-- Working Hours -->
    <div class="col-12">
      <div class="row items-center q-mb-sm">
        <q-icon name="schedule" color="primary" class="q-mr-xs" />
        <div class="text-subtitle2 text-weight-medium">{{ t.workingHours.title.value }}</div>
      </div>
      <q-toggle
        v-model="localData.enableWorkingHours"
        :label="t.workingHours.enableWorkingHours.value"
        color="primary"
        @update:model-value="updateValue"
      />
    </div>

    <template v-if="localData.enableWorkingHours">
      <div class="col-12">
        <q-checkbox
          v-model="localData.daysOfWeek"
          val="Monday"
          :label="t.workingHours.daysOfWeek.monday.value"
          color="primary"
          @update:model-value="updateValue"
        />
        <q-checkbox
          v-model="localData.daysOfWeek"
          val="Tuesday"
          :label="t.workingHours.daysOfWeek.tuesday.value"
          color="primary"
          @update:model-value="updateValue"
        />
        <q-checkbox
          v-model="localData.daysOfWeek"
          val="Wednesday"
          :label="t.workingHours.daysOfWeek.wednesday.value"
          color="primary"
          @update:model-value="updateValue"
        />
        <q-checkbox
          v-model="localData.daysOfWeek"
          val="Thursday"
          :label="t.workingHours.daysOfWeek.thursday.value"
          color="primary"
          @update:model-value="updateValue"
        />
        <q-checkbox
          v-model="localData.daysOfWeek"
          val="Friday"
          :label="t.workingHours.daysOfWeek.friday.value"
          color="primary"
          @update:model-value="updateValue"
        />
        <q-checkbox
          v-model="localData.daysOfWeek"
          val="Saturday"
          :label="t.workingHours.daysOfWeek.saturday.value"
          color="primary"
          @update:model-value="updateValue"
        />
        <q-checkbox
          v-model="localData.daysOfWeek"
          val="Sunday"
          :label="t.workingHours.daysOfWeek.sunday.value"
          color="primary"
          @update:model-value="updateValue"
        />
      </div>

      <!-- Time Intervals -->
      <div class="col-12">
        <div class="text-subtitle2 text-weight-medium text-grey-7 q-mb-sm">
          <q-icon name="access_time" size="sm" class="q-mr-xs" />
          {{ t.workingHours.timeIntervals.value }}
        </div>
      </div>

      <div
        v-for="(interval, index) in localData.timeIntervals"
        :key="index"
        class="col-12"
      >
        <q-card flat bordered class="rounded-borders q-pa-md">
          <div class="row q-col-gutter-md items-center">
            <div class="col-12 col-sm-5">
              <q-input
                v-model="interval.startTime"
                outlined
                dense
                :label="t.workingHours.startTime.value"
                placeholder="09:00"
                mask="##:##"
                class="rounded-borders"
                :rules="[
                  (val: any) => !!val || t.workingHours.startTimeRequired.value,
                  (val: any) => /^([0-1]?[0-9]|2[0-3]):[0-5][0-9]$/.test(val) || t.workingHours.invalidTimeFormat.value
                ]"
                @update:model-value="updateValue"
              >
                <template #prepend>
                  <q-icon name="schedule" color="primary" />
                </template>
              </q-input>
            </div>
            <div class="col-12 col-sm-5">
              <q-input
                v-model="interval.endTime"
                outlined
                dense
                :label="t.workingHours.endTime.value"
                placeholder="17:00"
                mask="##:##"
                class="rounded-borders"
                :rules="[
                  (val: any) => !!val || t.workingHours.endTimeRequired.value,
                  (val: any) => /^([0-1]?[0-9]|2[0-3]):[0-5][0-9]$/.test(val) || t.workingHours.invalidTimeFormat.value
                ]"
                @update:model-value="updateValue"
              >
                <template #prepend>
                  <q-icon name="schedule" color="primary" />
                </template>
              </q-input>
            </div>
            <div class="col-12 col-sm-2 text-center">
              <q-btn
                v-if="localData.timeIntervals.length > 1"
                flat
                dense
                round
                icon="delete"
                color="negative"
                size="sm"
                @click="removeInterval(index)"
              >
                <AppTooltip :content="t.workingHours.removeInterval.value" />
              </q-btn>
            </div>
          </div>
        </q-card>
      </div>

      <div class="col-12">
        <q-btn
          dense
          flat
          icon="add_circle"
          :label="t.workingHours.addInterval.value"
          color="primary"
          :ripple="false"
          @click="addInterval"
        />
      </div>

      <div class="col-12">
        <q-input
          v-model="localData.timezone"
          outlined
          dense
          class="rounded-borders"
          :label="`${t.workingHours.timezone.value} *`"
          :placeholder="t.workingHours.timezonePlaceholder.value"
          :rules="[(val: any) => !!val || t.workingHours.timezoneRequired.value]"
          @update:model-value="updateValue"
        />
      </div>
    </template>

    <!-- Rate Limit -->
    <div class="col-12">
      <div class="row items-center q-mb-sm">
        <q-icon name="speed" color="primary" class="q-mr-xs" />
        <div class="text-subtitle2 text-weight-medium">{{ t.rateLimit.title.value }}</div>
      </div>
      <q-toggle
        v-model="localData.enableRateLimit"
        :label="t.rateLimit.enableRateLimit.value"
        color="primary"
        @update:model-value="updateValue"
      />
    </div>

    <template v-if="localData.enableRateLimit">
      <div class="col-12 col-sm-6">
        <q-select
          v-model="localData.rateLimitType"
          outlined
          dense
          class="rounded-borders"
          :label="`${t.rateLimit.limitType.value} *`"
          :options="RATE_LIMIT_TYPE_OPTIONS"
          :rules="[(val: any) => !!val || t.rateLimit.limitTypeRequired.value]"
          @update:model-value="updateValue"
        >
          <template #prepend>
            <q-icon name="category" color="primary" />
          </template>
        </q-select>
      </div>

      <div class="col-12 col-sm-6">
        <q-input
          v-model="localData.rateLimitValue"
          outlined
          dense
          type="number"
          class="rounded-borders"
          :label="`${t.rateLimit.value.value} *`"
          :rules="[(val: any) => val >= 0 || t.rateLimit.valueRequired.value]"
          @update:model-value="updateValue"
        >
          <template #prepend>
            <q-icon name="speed" color="primary" />
          </template>
        </q-input>
      </div>

      <div class="col-12 col-sm-6">
        <q-input
          v-model="localData.burstCapacity"
          outlined
          dense
          type="number"
          class="rounded-borders"
          :label="t.rateLimit.burstCapacity.value"
          @update:model-value="updateValue"
        >
          <template #prepend>
            <q-icon name="trending_up" color="primary" />
          </template>
        </q-input>
      </div>

      <div class="col-12 col-sm-6">
        <q-select
          v-model="localData.actionOnExceed"
          outlined
          dense
          class="rounded-borders"
          :label="`${t.rateLimit.actionOnExceed.value} *`"
          :options="RATE_LIMIT_ACTION_OPTIONS"
          :rules="[(val: any) => !!val || t.rateLimit.actionRequired.value]"
          @update:model-value="updateValue"
        />
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
defineOptions({
  name: 'Step3WorkingHours'
});

/** TYPE IMPORTS */
import type { StepEmits, StepProps } from '../../interfaces/httpDataSource.interface';

/** VUE IMPORTS */
import { reactive, watch } from 'vue';

/** COMPONENTS */
import { AppTooltip } from '@components/tooltips';

/** COMPOSABLES */
import { useHttpDataSourceCreateEditTranslations } from '@composables/i18n/pages/datasources/http';

/** LOCAL IMPORTS */
import { RATE_LIMIT_ACTION_OPTIONS, RATE_LIMIT_TYPE_OPTIONS } from '../../constants/httpDataSourceConstants';

/** PROPS & EMITS */
const props = defineProps<StepProps>();
const emit = defineEmits<StepEmits>();

/** COMPOSABLES & STORES */
const t = useHttpDataSourceCreateEditTranslations();

const localData = reactive({
  enableWorkingHours: props.modelValue.enableWorkingHours || false,
  daysOfWeek: props.modelValue.daysOfWeek || [],
  timeIntervals: props.modelValue.timeIntervals || [{ startTime: '09:00', endTime: '17:00' }],
  timezone: props.modelValue.timezone || '',
  enableRateLimit: props.modelValue.enableRateLimit || false,
  rateLimitType: props.modelValue.rateLimitType || null,
  rateLimitValue: props.modelValue.rateLimitValue || 100,
  burstCapacity: props.modelValue.burstCapacity || 200,
  actionOnExceed: props.modelValue.actionOnExceed || null,
});

watch(() => props.modelValue, (newVal) => {
  localData.enableWorkingHours = newVal.enableWorkingHours || false;
  localData.daysOfWeek = newVal.daysOfWeek || [];
  localData.timeIntervals = newVal.timeIntervals || [{ startTime: '09:00', endTime: '17:00' }];
  localData.timezone = newVal.timezone || '';
  localData.enableRateLimit = newVal.enableRateLimit || false;
  localData.rateLimitType = newVal.rateLimitType || null;
  localData.rateLimitValue = newVal.rateLimitValue || 100;
  localData.burstCapacity = newVal.burstCapacity || 200;
  localData.actionOnExceed = newVal.actionOnExceed || null;
}, { deep: true, immediate: true });

/**
 * Add a new time interval to the working hours configuration
 * Adds a default interval from 09:00 to 17:00 and updates parent
 * @returns {void}
 */
function addInterval(): void {
  localData.timeIntervals.push({ startTime: '09:00', endTime: '17:00' });
  updateValue();
}

/**
 * Remove a time interval from the working hours configuration
 * Prevents removal if only one interval remains and updates parent
 * @param {number} index - Index of the interval to remove
 * @returns {void}
 */
function removeInterval(index: number): void {
  if (localData.timeIntervals.length > 1) {
    localData.timeIntervals.splice(index, 1);
    updateValue();
  }
}

/**
 * Emit updated values to parent component
 * Merges local form data with existing model value
 * @returns {void}
 */
function updateValue(): void {
  emit('update:modelValue', {
    ...props.modelValue,
    ...localData
  });
}
</script>

<style scoped>
.rounded-borders {
  border-radius: var(--mapex-radius-md);
}
</style>
