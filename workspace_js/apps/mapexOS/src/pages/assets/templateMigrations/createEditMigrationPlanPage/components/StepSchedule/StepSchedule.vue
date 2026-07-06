<script setup lang="ts">
/** TYPE IMPORTS */
import type { StepScheduleProps, StepScheduleEmits, ScheduleMode } from './interfaces';

defineOptions({
  name: 'StepSchedule'
});

/** VUE IMPORTS */
import { ref, computed, watch, onMounted } from 'vue';

/** COMPOSABLES */
import { useCreateEditMigrationPlanTranslations } from '@composables/i18n';

/** UTILS */
import { date } from 'quasar';

/** PROPS & EMITS */
const props = defineProps<StepScheduleProps>();
const emit = defineEmits<StepScheduleEmits>();

/** COMPOSABLES & STORES */
const t = useCreateEditMigrationPlanTranslations();

/** CONSTANTS */
const PICKER_MASK = 'YYYY-MM-DD HH:mm';

/** STATE */
// Local execution mode. Derived from the incoming value so edits reopen correctly.
const mode = ref<ScheduleMode>(props.modelValue ? 'scheduled' : 'now');
// Working display value for the date/time picker (local mask, not ISO).
const pickerValue = ref<string>(isoToLocal(props.modelValue));

/** COMPUTED */
const modeOptions = computed(() => [
  { label: t.steps.schedule.options.now.value, value: 'now' },
  { label: t.steps.schedule.options.scheduled.value, value: 'scheduled' },
]);

const scheduledDate = computed<Date | null>(() => {
  if (!pickerValue.value) return null;
  const parsed = new Date(pickerValue.value.replace(' ', 'T'));
  return Number.isNaN(parsed.getTime()) ? null : parsed;
});

const isFuture = computed(() => !!scheduledDate.value && scheduledDate.value.getTime() > Date.now());

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
 * Emit the current schedule value (ISO or null) and its validity.
 * Run-now is always valid; a scheduled run is valid only for a future datetime.
 * @returns {void}
 */
function sync(): void {
  if (mode.value === 'now') {
    emit('update:modelValue', null);
    emit('update:valid', true);
    return;
  }

  const valid = isFuture.value;
  emit('update:modelValue', valid && scheduledDate.value ? scheduledDate.value.toISOString() : null);
  emit('update:valid', valid);
}

/** WATCHERS */
watch([mode, pickerValue], sync);

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

      <div v-if="mode === 'scheduled'" class="col-12 col-md-6">
        <q-input
          v-model="pickerValue"
          outlined
          dense
          class="rounded-borders"
          mask="####-##-## ##:##"
          :label="t.steps.schedule.field.label.value + ' *'"
          :placeholder="t.steps.schedule.field.placeholder.value"
          :hint="t.steps.schedule.field.hint.value"
          :rules="[
            (val) => !!val || t.steps.schedule.field.required.value,
            () => isFuture || t.steps.schedule.field.future.value,
          ]"
        >
          <template v-slot:prepend>
            <q-icon name="event" color="primary" />
          </template>
          <template v-slot:append>
            <q-icon name="event" class="cursor-pointer">
              <q-popup-proxy cover transition-show="scale" transition-hide="scale">
                <q-date v-model="pickerValue" :mask="PICKER_MASK">
                  <div class="row items-center justify-end">
                    <q-btn v-close-popup :label="t.steps.schedule.field.close.value" color="primary" flat />
                  </div>
                </q-date>
              </q-popup-proxy>
            </q-icon>
            <q-icon name="access_time" class="cursor-pointer q-ml-xs">
              <q-popup-proxy cover transition-show="scale" transition-hide="scale">
                <q-time v-model="pickerValue" :mask="PICKER_MASK">
                  <div class="row items-center justify-end">
                    <q-btn v-close-popup :label="t.steps.schedule.field.close.value" color="primary" flat />
                  </div>
                </q-time>
              </q-popup-proxy>
            </q-icon>
          </template>
        </q-input>
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
