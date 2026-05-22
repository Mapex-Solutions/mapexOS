<script setup lang="ts">
/** TYPE IMPORTS */
import type { WorkflowGeneralSettings } from '../../interfaces/CreateEditWorkflow.interface';

/** VUE IMPORTS */
import { ref, computed, watch, nextTick } from 'vue';

/** COMPOSABLES */
import { useWorkflowEditorState } from '../../composables';
import { useCreateEditWorkflowTranslations } from '@composables/i18n/pages/automations/workflows';

/** LOCAL IMPORTS (constants and handlers ONLY) */
import { TIMEZONE_TYPE_OPTIONS, IANA_TIMEZONE_OPTIONS } from './constants';

/** COMPOSABLES & STORES */
const { generalSettings, updateGeneral, externalInputs } = useWorkflowEditorState();
const t = useCreateEditWorkflowTranslations();

/** STATE */

/**
 * Local copy of general settings for form binding
 */
const form = ref<WorkflowGeneralSettings>({ ...generalSettings.value });

/** STATE */

/**
 * Filtered IANA timezone options (for search/filter in q-select)
 */
const filteredTimezones = ref(IANA_TIMEZONE_OPTIONS);

/** COMPUTED */

/**
 * Input variable path options for timezone value when type is 'variable'.
 * Shows input.{field} paths from the Data tab — the workflow's external contract.
 */
const inputPathOptions = computed(() =>
  externalInputs.value.map((input) => ({
    label: `input.${input.field}`,
    value: `input.${input.field}`,
    caption: input.label,
    icon: input.icon || 'input',
  })),
);

/** FUNCTIONS */

/**
 * Filter IANA timezone options based on user input
 *
 * @param {string} val - Current filter text
 * @param {(callbackFn: () => void) => void} update - Quasar filter update callback
 */
function filterTimezones(val: string, update: (fn: () => void) => void): void {
  update(() => {
    const needle = val.toLowerCase();
    filteredTimezones.value = needle
      ? IANA_TIMEZONE_OPTIONS.filter((tz) => tz.label.toLowerCase().includes(needle))
      : IANA_TIMEZONE_OPTIONS;
  });
}

/** WATCHERS */

/**
 * Flag to prevent bidirectional watcher ping-pong.
 * When true, the form→composable watcher skips its update.
 */
let syncing = false;

/**
 * Sync external changes back to local form
 */
watch(
  () => generalSettings.value,
  (newVal) => {
    syncing = true;
    form.value = { ...newVal };
    void nextTick(() => { syncing = false; });
  },
  { deep: true },
);

/**
 * Push local form changes to composable
 */
watch(
  () => form.value,
  (newVal) => {
    if (syncing) return;
    updateGeneral(newVal);
  },
  { deep: true },
);
</script>

<template>
  <div class="general-tab">
    <div class="row q-col-gutter-md">
      <!-- ── Left Column: Core Settings ── -->
      <div class="col-12 col-md-7">
        <!-- Basic Information -->
        <q-card flat bordered class="q-mb-md">
          <q-card-section>
            <div class="text-subtitle1 text-weight-medium q-mb-md">{{ t.generalTab.basicInfo.value }}</div>

            <div class="row q-col-gutter-sm">
              <div class="col-12 col-sm-8">
                <q-input
                  v-model="form.name"
                  :label="t.generalTab.name.value"
                  outlined
                  dense
                  :rules="[(val: string) => !!val || t.validation.nameIsRequired.value]"
                />
              </div>
              <div class="col-12 col-sm-4">
                <q-select
                  v-model="form.enabled"
                  :label="t.generalTab.status.value"
                  outlined
                  dense
                  :options="[
                    { label: t.generalTab.enabled.value, value: true },
                    { label: t.generalTab.disabled.value, value: false },
                  ]"
                  emit-value
                  map-options
                />
              </div>
              <div class="col-12">
                <q-input
                  v-model="form.description"
                  :label="t.generalTab.description.value"
                  outlined
                  dense
                  type="textarea"
                  autogrow
                />
              </div>
            </div>

            <q-separator class="q-my-md" />

            <div class="row q-gutter-md">
              <q-checkbox
                v-model="form.sharedWithChildren"
                :label="t.generalTab.sharedWithChildren.value"
                dense
              />
            </div>
          </q-card-section>
        </q-card>

      </div>

      <!-- ── Right Column: Execution Settings ── -->
      <div class="col-12 col-md-5">
        <!-- Timezone -->
        <q-card flat bordered class="q-mb-md">
          <q-card-section>
            <div class="text-subtitle1 text-weight-medium q-mb-md">
              <q-icon name="schedule" color="blue-7" class="q-mr-xs" />
              {{ t.generalTab.timezone.value }}
            </div>

            <q-select
              v-model="form.timezone.type"
              :label="t.generalTab.timezoneType.value"
              outlined
              dense
              :options="[...TIMEZONE_TYPE_OPTIONS]"
              emit-value
              map-options
              class="q-mb-md"
            />

            <!-- Literal: searchable IANA timezone select -->
            <q-select
              v-if="form.timezone.type === 'literal'"
              v-model="form.timezone.value"
              :label="t.generalTab.timezoneValue.value"
              outlined
              dense
              use-input
              input-debounce="100"
              emit-value
              map-options
              option-value="value"
              option-label="label"
              :options="filteredTimezones"
              :hint="t.generalTab.timezoneHintIana.value"
              class="q-mb-md"
              @filter="filterTimezones"
            >
              <template #option="scope">
                <q-item v-bind="scope.itemProps">
                  <q-item-section avatar>
                    <q-icon name="schedule" size="20px" color="grey-7" />
                  </q-item-section>
                  <q-item-section>
                    <q-item-label>{{ scope.opt.label }}</q-item-label>
                    <q-item-label caption>{{ scope.opt.region }}</q-item-label>
                  </q-item-section>
                </q-item>
              </template>
            </q-select>

            <!-- Variable: select from Data/Inputs contract -->
            <q-select
              v-else
              v-model="form.timezone.value"
              :label="t.generalTab.timezoneValue.value"
              outlined
              dense
              emit-value
              map-options
              option-value="value"
              option-label="label"
              :options="inputPathOptions"
              :hint="t.generalTab.timezoneHintVariable.value"
              class="q-mb-md"
            >
              <template #option="scope">
                <q-item v-bind="scope.itemProps">
                  <q-item-section avatar>
                    <q-icon :name="scope.opt.icon" size="20px" color="cyan-6" />
                  </q-item-section>
                  <q-item-section>
                    <q-item-label>{{ scope.opt.label }}</q-item-label>
                    <q-item-label caption>{{ scope.opt.caption }}</q-item-label>
                  </q-item-section>
                </q-item>
              </template>
              <template #no-option>
                <q-item>
                  <q-item-section class="text-grey-6 text-caption">
                    {{ t.externalInputs.emptyTitle.value }}
                  </q-item-section>
                </q-item>
              </template>
            </q-select>

            <div class="text-caption text-grey-6">
              <q-icon name="info" size="xs" class="q-mr-xs" />
              {{ t.generalTab.timezoneInfo.value }}
            </div>
          </q-card-section>
        </q-card>
      </div>
    </div>
  </div>
</template>
