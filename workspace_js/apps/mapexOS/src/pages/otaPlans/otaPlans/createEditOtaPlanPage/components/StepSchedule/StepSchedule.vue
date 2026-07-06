<script setup lang="ts">
/** TYPE IMPORTS */
import type { OTARolloutConfig } from '@mapexos/schemas';
import type { StepScheduleProps, StepScheduleEmits, ScheduleMode } from './interfaces/stepSchedule.interface';

defineOptions({
  name: 'StepSchedule'
});

/** VUE IMPORTS */
import { ref, computed, watch, onMounted } from 'vue';

/** COMPOSABLES */
import { useCreateEditOtaPlanTranslations } from '@composables/i18n/pages/otaPlans/createEditOtaPlan/useCreateEditOtaPlanTranslations';

/** UTILS */
import { date } from 'quasar';

/** PROPS & EMITS */
const props = defineProps<StepScheduleProps>();
const emit = defineEmits<StepScheduleEmits>();

/** COMPOSABLES & STORES */
const t = useCreateEditOtaPlanTranslations();

/** CONSTANTS */
const PICKER_MASK = 'YYYY-MM-DD HH:mm';

/** STATE */
// Local execution mode. Derived from the incoming value so edits reopen correctly.
const mode = ref<ScheduleMode>(props.startAt ? 'scheduled' : 'now');
// Working display values for the pickers (local mask, not ISO).
const startPicker = ref<string>(isoToLocal(props.startAt));
const maxPicker = ref<string>(isoToLocal(props.maxTime));

/** COMPUTED */
const modeOptions = computed<{ label: string; value: ScheduleMode }[]>(() => [
  { label: t.steps.schedule.options.now.value, value: 'now' },
  { label: t.steps.schedule.options.scheduled.value, value: 'scheduled' },
]);

const startDate = computed<Date | null>(() => parsePicker(startPicker.value));
const maxDate = computed<Date | null>(() => parsePicker(maxPicker.value));

const isStartFuture = computed(() => !!startDate.value && startDate.value.getTime() > Date.now());

/** The deadline must come after the effective start (now for run-now mode) */
const isMaxAfterStart = computed(() => {
  if (!maxDate.value) return false;
  const startMs = mode.value === 'scheduled' && startDate.value ? startDate.value.getTime() : Date.now();
  return maxDate.value.getTime() > startMs;
});

/** FUNCTIONS */

/**
 * Convert an ISO datetime string to the local picker mask.
 * @param iso - ISO datetime string or null
 * @returns Masked local datetime, or empty string when absent/invalid
 */
function isoToLocal(iso: string | null): string {
  if (!iso) return '';
  const parsed = new Date(iso);
  if (Number.isNaN(parsed.getTime())) return '';
  return date.formatDate(parsed, PICKER_MASK);
}

/**
 * Parse a masked picker value into a Date.
 * @param value - Masked local datetime
 * @returns Date, or null when empty/invalid
 */
function parsePicker(value: string): Date | null {
  if (!value) return null;
  const parsed = new Date(value.replace(' ', 'T'));
  return Number.isNaN(parsed.getTime()) ? null : parsed;
}

/**
 * Forward a pacing/abort change.
 */
function updateRollout(partial: Partial<OTARolloutConfig>): void {
  emit('update:rolloutConfig', { ...props.rolloutConfig, ...partial });
}

/**
 * Emit the current schedule slice (ISO or null) and its validity.
 * Run-now is valid with any future deadline; a scheduled run also requires a
 * future start, with the deadline after it.
 */
function sync(): void {
  const startValid = mode.value === 'now' || isStartFuture.value;
  const valid = startValid && isMaxAfterStart.value && props.rolloutConfig.ratePerMinute > 0;

  emit('update:startAt', mode.value === 'scheduled' && isStartFuture.value && startDate.value ? startDate.value.toISOString() : null);
  emit('update:maxTime', isMaxAfterStart.value && maxDate.value ? maxDate.value.toISOString() : null);
  emit('update:valid', valid);
}

/** WATCHERS */
watch([mode, startPicker, maxPicker, () => props.rolloutConfig.ratePerMinute], sync);

/** LIFECYCLE HOOKS */
onMounted(sync);
</script>

<template>
  <div>
    <div class="q-mb-md">
      <div class="text-subtitle1 text-weight-medium q-mb-xs">
        <q-icon name="schedule" color="primary" class="q-mr-xs" />
        {{ t.steps.schedule.title.value }}
      </div>
      <div class="text-body2 text-grey-7">
        {{ t.steps.schedule.subtitle.value }}
      </div>
    </div>

    <div class="row q-col-gutter-md">
      <!-- Execution mode -->
      <div class="col-12">
        <div class="text-caption text-grey-7 q-mb-xs">{{ t.steps.schedule.modeLabel.value }}</div>
        <q-option-group
          v-model="mode"
          type="radio"
          color="primary"
          inline
          :options="modeOptions"
        />
      </div>

      <!-- Start (only when scheduled) -->
      <div v-if="mode === 'scheduled'" class="col-12 col-md-6">
        <q-input
          v-model="startPicker"
          outlined
          dense
          class="rounded-borders"
          mask="####-##-## ##:##"
          :label="t.steps.schedule.field.label.value + ' *'"
          :placeholder="t.steps.schedule.field.placeholder.value"
          :hint="t.steps.schedule.field.hint.value"
          :rules="[
            (val) => !!val || t.steps.schedule.field.required.value,
            () => isStartFuture || t.steps.schedule.field.future.value,
          ]"
        >
          <template v-slot:prepend>
            <q-icon name="event" color="primary" />
          </template>
          <template v-slot:append>
            <q-icon name="event" class="cursor-pointer">
              <q-popup-proxy cover transition-show="scale" transition-hide="scale">
                <q-date v-model="startPicker" :mask="PICKER_MASK">
                  <div class="row items-center justify-end">
                    <q-btn v-close-popup :label="t.steps.schedule.field.close.value" color="primary" flat />
                  </div>
                </q-date>
              </q-popup-proxy>
            </q-icon>
            <q-icon name="access_time" class="cursor-pointer q-ml-xs">
              <q-popup-proxy cover transition-show="scale" transition-hide="scale">
                <q-time v-model="startPicker" :mask="PICKER_MASK">
                  <div class="row items-center justify-end">
                    <q-btn v-close-popup :label="t.steps.schedule.field.close.value" color="primary" flat />
                  </div>
                </q-time>
              </q-popup-proxy>
            </q-icon>
          </template>
        </q-input>
      </div>

      <!-- Deadline (always required) -->
      <div class="col-12 col-md-6">
        <q-input
          v-model="maxPicker"
          outlined
          dense
          class="rounded-borders"
          mask="####-##-## ##:##"
          :label="t.steps.schedule.maxTime.label.value + ' *'"
          :placeholder="t.steps.schedule.maxTime.placeholder.value"
          :hint="t.steps.schedule.maxTime.hint.value"
          :rules="[
            (val) => !!val || t.steps.schedule.maxTime.required.value,
            () => isMaxAfterStart || t.steps.schedule.maxTime.afterStart.value,
          ]"
        >
          <template v-slot:prepend>
            <q-icon name="event_busy" color="primary" />
          </template>
          <template v-slot:append>
            <q-icon name="event" class="cursor-pointer">
              <q-popup-proxy cover transition-show="scale" transition-hide="scale">
                <q-date v-model="maxPicker" :mask="PICKER_MASK">
                  <div class="row items-center justify-end">
                    <q-btn v-close-popup :label="t.steps.schedule.field.close.value" color="primary" flat />
                  </div>
                </q-date>
              </q-popup-proxy>
            </q-icon>
            <q-icon name="access_time" class="cursor-pointer q-ml-xs">
              <q-popup-proxy cover transition-show="scale" transition-hide="scale">
                <q-time v-model="maxPicker" :mask="PICKER_MASK">
                  <div class="row items-center justify-end">
                    <q-btn v-close-popup :label="t.steps.schedule.field.close.value" color="primary" flat />
                  </div>
                </q-time>
              </q-popup-proxy>
            </q-icon>
          </template>
        </q-input>
      </div>

      <!-- Pacing + abort -->
      <div class="col-12">
        <div class="text-subtitle2 q-mb-sm">
          <q-icon name="speed" color="primary" class="q-mr-xs" />
          {{ t.steps.schedule.pacing.title.value }}
        </div>
        <div class="row q-col-gutter-md">
          <div class="col-12 col-md-4">
            <q-input
              :model-value="rolloutConfig.ratePerMinute"
              outlined
              dense
              type="number"
              min="1"
              class="rounded-borders"
              :label="t.steps.schedule.pacing.ratePerMinute.value"
              @update:model-value="(v) => updateRollout({ ratePerMinute: Number(v) || 0 })"
            />
          </div>
          <div class="col-12 col-md-4">
            <q-input
              :model-value="rolloutConfig.abortThresholdPct"
              outlined
              dense
              type="number"
              min="0"
              max="100"
              class="rounded-borders"
              :label="t.steps.schedule.pacing.abortThresholdPct.value"
              :hint="t.steps.schedule.pacing.abortThresholdHint.value"
              @update:model-value="(v) => updateRollout({ abortThresholdPct: Number(v) || 0 })"
            />
          </div>
          <div class="col-12 col-md-4">
            <q-input
              :model-value="rolloutConfig.abortMinExecuted"
              outlined
              dense
              type="number"
              min="0"
              class="rounded-borders"
              :label="t.steps.schedule.pacing.abortMinExecuted.value"
              @update:model-value="(v) => updateRollout({ abortMinExecuted: Number(v) || 0 })"
            />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.rounded-borders {
  border-radius: var(--mapex-radius-md);
}

.cursor-pointer {
  cursor: pointer;
}
</style>
