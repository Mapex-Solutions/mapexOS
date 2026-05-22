<script setup lang="ts">
defineOptions({
  name: 'EventTraceVisualization'
});

/** TYPE IMPORTS */
import type { EventTraceVisualizationProps, EventTraceVisualizationEmits } from './interfaces';

/** VUE IMPORTS */
import { ref, computed } from 'vue';

/** COMPOSABLES */
import { useEventTracerPageTranslations } from '@composables/i18n/pages/logs/eventTracerPage';

/** LOCAL IMPORTS */
import { DEFAULT_EXPANDED_STATE } from './constants';
import {
  STATUS_COLORS,
  STATUS_ICONS,
  TRIGGER_TYPE_ICONS,
  SOURCE_TYPE_ICONS,
  SOURCE_TYPE_COLORS,
} from '../../constants';

/** PROPS & EMITS */
const props = defineProps<EventTraceVisualizationProps>();
const emit = defineEmits<EventTraceVisualizationEmits>();

/** COMPOSABLES */
const t = useEventTracerPageTranslations();

/** STATE */
const expandedPhases = ref({ ...DEFAULT_EXPANDED_STATE });

/** COMPUTED */

/**
 * Check if trace has any data
 */
const hasData = computed(() => props.traceResult !== null);

/**
 * Overall status color
 */
const overallStatusColor = computed(() => {
  if (!props.traceResult) return STATUS_COLORS.pending;
  return props.traceResult.allSuccess ? STATUS_COLORS.success : STATUS_COLORS.failed;
});

/**
 * Overall status icon
 */
const overallStatusIcon = computed(() => {
  if (!props.traceResult) return STATUS_ICONS.pending;
  return props.traceResult.allSuccess ? STATUS_ICONS.success : STATUS_ICONS.failed;
});

/** FUNCTIONS */

/**
 * Toggle phase expansion state
 * @param {string} phase - Phase identifier
 */
function togglePhase(phase: 'ingestion' | 'routing'): void {
  expandedPhases.value[phase] = !expandedPhases.value[phase];
}

/**
 * Format duration in milliseconds to human readable
 * @param {number | null} ms - Duration in milliseconds
 * @returns {string} Formatted duration
 */
function formatDuration(ms: number | null | undefined): string {
  if (ms === null || ms === undefined) return '0ms';
  if (ms < 1000) return `${ms}ms`;
  return `${(ms / 1000).toFixed(2)}s`;
}

/**
 * Format timestamp to human readable
 * @param {string | null} timestamp - ISO timestamp
 * @returns {string} Formatted timestamp
 */
function formatTimestamp(timestamp: string | null): string {
  if (!timestamp) return '-';
  const date = new Date(timestamp);
  return date.toLocaleString('en-US', {
    day: '2-digit',
    month: 'short',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
    fractionalSecondDigits: 3,
  });
}

/**
 * Get status color for a stage
 * @param {boolean | null} success - Success status
 * @param {boolean} hasData - Whether stage has data
 * @returns {string} Color class
 */
function getStatusColor(success: boolean | null, hasData: boolean): string {
  if (!hasData) return STATUS_COLORS.pending;
  return success ? STATUS_COLORS.success : STATUS_COLORS.failed;
}

/**
 * Get trigger type icon
 * @param {string | null} type - Trigger type
 * @returns {string} Icon name
 */
function getTriggerIcon(type: string | null): string {
  return type ? TRIGGER_TYPE_ICONS[type] || 'flash_on' : 'flash_on';
}

/**
 * Get source type icon
 * @param {string | null} source - Source type
 * @returns {string} Icon name
 */
function getSourceIcon(source: string | null): string {
  return source ? SOURCE_TYPE_ICONS[source] || 'input' : 'input';
}

/**
 * Get source type color
 * @param {string | null} source - Source type
 * @returns {string} Color class
 */
function getSourceColor(source: string | null): string {
  return source ? SOURCE_TYPE_COLORS[source] || 'grey-6' : 'grey-6';
}

/**
 * Emit view details event
 * @param {string} stage - Stage identifier
 * @param {any} data - Stage data
 */
function handleViewDetails(stage: string, data: any): void {
  emit('view-details', { stage, data });
}
</script>

<template>
  <div class="event-trace-visualization">
    <!-- Loading State -->
    <div v-if="loading" class="flex flex-center q-pa-xl">
      <q-spinner color="primary" size="50px" />
      <span class="q-ml-md text-grey-7">{{ t.messages.loading?.value || 'Loading trace...' }}</span>
    </div>

    <!-- No Data State -->
    <div v-else-if="!hasData" class="text-center q-pa-xl">
      <q-icon name="search" size="64px" color="grey-5" />
      <div class="text-h6 text-grey-7 q-mt-md">{{ t.empty.searchFirst.value }}</div>
      <div class="text-body2 text-grey-6 q-mt-sm">{{ t.empty.searchFirstDescription.value }}</div>
    </div>

    <!-- Trace Visualization -->
    <div v-else-if="traceResult" class="trace-content">
      <!-- Summary Header -->
      <q-card flat bordered class="q-mb-md summary-card">
        <q-card-section class="q-pa-md">
          <div class="row items-center q-gutter-md">
            <!-- Overall Status -->
            <div class="col-auto">
              <q-avatar :color="overallStatusColor" text-color="white" size="48px">
                <q-icon :name="overallStatusIcon" />
              </q-avatar>
            </div>

            <!-- Summary Info -->
            <div class="col">
              <div class="text-subtitle1 text-weight-medium">
                {{ traceResult.allSuccess ? t.summary.allSuccess.value : t.summary.hasFailed.value }}
              </div>
              <div class="text-caption text-grey-7">
                {{ t.summary.stagesProcessed.value }}: {{ traceResult.summary.stagesWithData }}
                &bull;
                {{ t.summary.totalExecutionTime.value }}: {{ formatDuration(traceResult.totalExecutionTime) }}
              </div>
            </div>

            <!-- Quick Stats -->
            <div class="col-auto">
              <div class="row q-gutter-sm">
                <q-chip
                  dense
                  color="teal-1"
                  text-color="teal-8"
                  icon="gavel"
                >
                </q-chip>
                <q-chip
                  dense
                  color="purple-1"
                  text-color="purple-8"
                  icon="flash_on"
                >
                  {{ traceResult.summary.totalTriggers }} {{ t.stages.trigger.value }}
                </q-chip>
              </div>
            </div>
          </div>
        </q-card-section>
      </q-card>

      <!-- Phase 1: Ingestion -->
      <q-card flat bordered class="q-mb-md phase-card">
        <q-card-section
          class="phase-header cursor-pointer"
          :class="{ 'bg-blue-1': expandedPhases.ingestion }"
          @click="togglePhase('ingestion')"
        >
          <div class="row items-center no-wrap">
            <q-avatar color="blue-6" text-color="white" size="36px" class="q-mr-md">
              <q-icon name="input" />
            </q-avatar>
            <div class="col">
              <div class="text-subtitle1 text-weight-medium text-blue-8">
                {{ t.phases.ingestion.value }}
              </div>
              <div class="text-caption text-grey-7">
                {{ formatDuration(traceResult.ingestion.totalDurationMs) }}
                &bull;
                {{ traceResult.ingestion.allSuccess ? t.stageStatus.success.value : t.stageStatus.failed.value }}
              </div>
            </div>
            <q-icon
              :name="expandedPhases.ingestion ? 'expand_less' : 'expand_more'"
              size="sm"
              color="grey-6"
            />
          </div>
        </q-card-section>

        <q-slide-transition>
          <div v-show="expandedPhases.ingestion">
            <q-separator />
            <q-card-section class="q-pa-md">
              <div class="row q-col-gutter-md">
                <!-- Raw Event Stage -->
                <div class="col-12 col-md-6">
                  <q-card flat bordered class="stage-card">
                    <q-card-section class="q-pa-sm">
                      <div class="row items-center no-wrap q-mb-sm">
                        <q-avatar
                          :color="getStatusColor(traceResult.ingestion.raw.success, traceResult.ingestion.raw.hasData)"
                          text-color="white"
                          size="32px"
                        >
                          <q-icon name="terminal" size="xs" />
                        </q-avatar>
                        <div class="q-ml-sm col">
                          <div class="text-body2 text-weight-medium">{{ t.stages.raw.value }}</div>
                          <div class="text-caption text-grey-6">
                            {{ formatTimestamp(traceResult.ingestion.raw.created) }}
                          </div>
                        </div>
                        <q-chip
                          v-if="traceResult.ingestion.raw.source"
                          dense
                          :color="getSourceColor(traceResult.ingestion.raw.source)"
                          text-color="white"
                          size="sm"
                        >
                          <q-icon :name="getSourceIcon(traceResult.ingestion.raw.source)" size="xs" class="q-mr-xs" />
                          {{ traceResult.ingestion.raw.source }}
                        </q-chip>
                      </div>
                      <div v-if="traceResult.ingestion.raw.error" class="text-caption text-negative q-mt-sm">
                        <q-icon name="error" size="xs" class="q-mr-xs" />
                        {{ traceResult.ingestion.raw.error }}
                      </div>
                      <q-btn
                        v-if="traceResult.ingestion.raw.hasData"
                        flat
                        dense
                        color="primary"
                        icon="visibility"
                        :label="t.timeline.viewDetails.value"
                        size="sm"
                        class="q-mt-sm"
                        @click="handleViewDetails('raw', traceResult.ingestion.raw.data)"
                      />
                    </q-card-section>
                  </q-card>
                </div>

                <!-- JS Executor Stage -->
                <div class="col-12 col-md-6">
                  <q-card flat bordered class="stage-card">
                    <q-card-section class="q-pa-sm">
                      <div class="row items-center no-wrap q-mb-sm">
                        <q-avatar
                          :color="getStatusColor(traceResult.ingestion.jsExec.success, traceResult.ingestion.jsExec.hasData)"
                          text-color="white"
                          size="32px"
                        >
                          <q-icon name="code" size="xs" />
                        </q-avatar>
                        <div class="q-ml-sm col">
                          <div class="text-body2 text-weight-medium">{{ t.stages.jsExec.value }}</div>
                          <div class="text-caption text-grey-6">
                            {{ formatDuration(traceResult.ingestion.jsExec.durationMs) }}
                          </div>
                        </div>
                        <q-chip
                          v-if="traceResult.ingestion.jsExec.hasData"
                          dense
                          :color="traceResult.ingestion.jsExec.success ? 'positive' : 'negative'"
                          text-color="white"
                          size="sm"
                        >
                          {{ traceResult.ingestion.jsExec.success ? t.stageStatus.success.value : t.stageStatus.failed.value }}
                        </q-chip>
                      </div>
                      <div v-if="traceResult.ingestion.jsExec.failedAt" class="text-caption text-warning q-mt-sm">
                        <q-icon name="warning" size="xs" class="q-mr-xs" />
                        Failed at: {{ traceResult.ingestion.jsExec.failedAt }}
                      </div>
                      <div v-if="traceResult.ingestion.jsExec.error" class="text-caption text-negative q-mt-sm">
                        <q-icon name="error" size="xs" class="q-mr-xs" />
                        {{ traceResult.ingestion.jsExec.error }}
                      </div>
                      <q-btn
                        v-if="traceResult.ingestion.jsExec.hasData"
                        flat
                        dense
                        color="primary"
                        icon="visibility"
                        :label="t.timeline.viewDetails.value"
                        size="sm"
                        class="q-mt-sm"
                        @click="handleViewDetails('jsExec', traceResult.ingestion.jsExec.data)"
                      />
                    </q-card-section>
                  </q-card>
                </div>
              </div>
            </q-card-section>
          </div>
        </q-slide-transition>
      </q-card>

      <!-- Phase 2: Routing -->
      <q-card flat bordered class="q-mb-md phase-card">
        <q-card-section
          class="phase-header cursor-pointer"
          :class="{ 'bg-orange-1': expandedPhases.routing }"
          @click="togglePhase('routing')"
        >
          <div class="row items-center no-wrap">
            <q-avatar color="orange-6" text-color="white" size="36px" class="q-mr-md">
              <q-icon name="call_split" />
            </q-avatar>
            <div class="col">
              <div class="text-subtitle1 text-weight-medium text-orange-8">
                {{ t.phases.routing.value }}
              </div>
              <div class="text-caption text-grey-7">
                {{ traceResult.routing.destinationsCount }} {{ t.routing.destinations.value }}
                &bull;
                {{ formatDuration(traceResult.routing.totalDurationMs) }}
              </div>
            </div>
            <q-icon
              :name="expandedPhases.routing ? 'expand_less' : 'expand_more'"
              size="sm"
              color="grey-6"
            />
          </div>
        </q-card-section>

        <q-slide-transition>
          <div v-show="expandedPhases.routing">
            <q-separator />
            <q-card-section class="q-pa-md">
              <!-- Router Info -->
              <div v-if="traceResult.routing.router.hasData" class="q-mb-md router-info-section">
                <div class="row items-center q-mb-sm">
                  <q-icon name="account_tree" size="sm" color="orange-7" class="q-mr-sm" />
                  <span class="text-subtitle2 text-weight-medium text-grey-8">
                    {{ traceResult.routing.router.name || 'Route Group' }}
                  </span>
                </div>
                <div class="row q-gutter-sm">
                  <q-chip dense square color="grey-2" text-color="grey-8" size="sm" icon="layers">
                    {{ traceResult.routing.router.totalRouters }} destinations
                  </q-chip>
                  <q-chip dense square color="green-1" text-color="green-8" size="sm" icon="check_circle">
                    {{ traceResult.routing.router.matchedCount }} matched
                  </q-chip>
                  <q-chip dense square color="blue-1" text-color="blue-8" size="sm" icon="send">
                    {{ traceResult.routing.router.publishedCount }} published
                  </q-chip>
                </div>
              </div>

              <!-- Routed Destinations Grid -->
              <div class="row q-col-gutter-md">
                <!-- Direct Triggers -->
                <div
                  v-for="(trigger, idx) in traceResult.routing.directTriggers"
                  :key="`direct-trigger-${idx}`"
                  class="col-12 col-sm-6 col-md-4"
                >
                  <q-card flat bordered class="stage-card">
                    <q-card-section class="q-pa-sm">
                      <div class="row items-center no-wrap q-mb-sm">
                        <q-avatar
                          :color="getStatusColor(trigger.success, trigger.hasData)"
                          text-color="white"
                          size="32px"
                        >
                          <q-icon :name="getTriggerIcon(trigger.triggerType)" size="xs" />
                        </q-avatar>
                        <div class="q-ml-sm col">
                          <div class="text-body2 text-weight-medium ellipsis">
                            {{ trigger.triggerName || t.stages.directTrigger.value }}
                          </div>
                          <div class="text-caption text-grey-6">
                            {{ trigger.triggerType }} &bull; {{ formatDuration(trigger.durationMs) }}
                          </div>
                        </div>
                      </div>
                      <div class="row items-center q-gutter-xs q-mb-sm">
                        <q-chip
                          dense
                          :color="trigger.success ? 'positive' : 'negative'"
                          text-color="white"
                          size="sm"
                        >
                          {{ trigger.success ? 'Success' : 'Failed' }}
                        </q-chip>
                        <q-chip
                          dense
                          :color="trigger.category === 'technical' ? 'blue-6' : 'green-6'"
                          text-color="white"
                          size="sm"
                        >
                          {{ trigger.category }}
                        </q-chip>
                      </div>
                      <div v-if="trigger.error" class="text-caption text-negative q-mb-sm">
                        <q-icon name="error" size="xs" class="q-mr-xs" />
                        {{ trigger.error }}
                      </div>
                      <q-btn
                        v-if="trigger.hasData"
                        flat
                        dense
                        color="primary"
                        icon="visibility"
                        :label="t.timeline.viewDetails.value"
                        size="sm"
                        @click="handleViewDetails('directTrigger', trigger.data)"
                      />
                    </q-card-section>
                  </q-card>
                </div>

                <!-- Save Event -->
                <div v-if="traceResult.routing.saveEvent.hasData" class="col-12 col-sm-6 col-md-4">
                  <q-card flat bordered class="stage-card">
                    <q-card-section class="q-pa-sm">
                      <div class="row items-center no-wrap q-mb-sm">
                        <q-avatar
                          :color="getStatusColor(traceResult.routing.saveEvent.success, traceResult.routing.saveEvent.hasData)"
                          text-color="white"
                          size="32px"
                        >
                          <q-icon name="save" size="xs" />
                        </q-avatar>
                        <div class="q-ml-sm col">
                          <div class="text-body2 text-weight-medium">{{ t.stages.saveEvent.value }}</div>
                          <div class="text-caption text-grey-6">
                            {{ formatDuration(traceResult.routing.saveEvent.durationMs) }}
                          </div>
                        </div>
                      </div>
                      <div class="row items-center q-gutter-xs q-mb-sm">
                        <q-chip
                          dense
                          :color="traceResult.routing.saveEvent.triggered ? 'positive' : 'warning'"
                          text-color="white"
                          size="sm"
                        >
                          {{ traceResult.routing.saveEvent.triggered ? 'Matched' : 'Not Matched' }}
                        </q-chip>
                        <q-chip
                          dense
                          :color="traceResult.routing.saveEvent.success ? 'positive' : 'grey-5'"
                          text-color="white"
                          size="sm"
                        >
                          {{ traceResult.routing.saveEvent.success ? 'Published' : 'Not Published' }}
                        </q-chip>
                      </div>
                      <q-btn
                        v-if="traceResult.routing.saveEvent.data"
                        flat
                        dense
                        color="primary"
                        icon="visibility"
                        :label="t.timeline.viewDetails.value"
                        size="sm"
                        @click="handleViewDetails('saveEvent', traceResult.routing.saveEvent.data)"
                      />
                    </q-card-section>
                  </q-card>
                </div>

                <!-- Data Lake -->
                <div v-if="traceResult.routing.lakeHouse.hasData" class="col-12 col-sm-6 col-md-4">
                  <q-card flat bordered class="stage-card">
                    <q-card-section class="q-pa-sm">
                      <div class="row items-center no-wrap q-mb-sm">
                        <q-avatar
                          :color="getStatusColor(traceResult.routing.lakeHouse.success, traceResult.routing.lakeHouse.hasData)"
                          text-color="white"
                          size="32px"
                        >
                          <q-icon name="storage" size="xs" />
                        </q-avatar>
                        <div class="q-ml-sm col">
                          <div class="text-body2 text-weight-medium">{{ t.stages.lakeHouse.value }}</div>
                          <div class="text-caption text-grey-6">
                            {{ formatDuration(traceResult.routing.lakeHouse.durationMs) }}
                          </div>
                        </div>
                      </div>
                      <div class="row items-center q-gutter-xs q-mb-sm">
                        <q-chip
                          dense
                          :color="traceResult.routing.lakeHouse.triggered ? 'positive' : 'warning'"
                          text-color="white"
                          size="sm"
                        >
                          {{ traceResult.routing.lakeHouse.triggered ? 'Matched' : 'Not Matched' }}
                        </q-chip>
                        <q-chip
                          dense
                          :color="traceResult.routing.lakeHouse.success ? 'positive' : 'grey-5'"
                          text-color="white"
                          size="sm"
                        >
                          {{ traceResult.routing.lakeHouse.success ? 'Published' : 'Not Published' }}
                        </q-chip>
                      </div>
                      <q-btn
                        v-if="traceResult.routing.lakeHouse.data"
                        flat
                        dense
                        color="primary"
                        icon="visibility"
                        :label="t.timeline.viewDetails.value"
                        size="sm"
                        @click="handleViewDetails('lakeHouse', traceResult.routing.lakeHouse.data)"
                      />
                    </q-card-section>
                  </q-card>
                </div>
              </div>
            </q-card-section>
          </div>
        </q-slide-transition>
      </q-card>

    </div>
  </div>
</template>

<style scoped lang="scss">
.event-trace-visualization {
  width: 100%;
}

.summary-card {
  background: linear-gradient(135deg, var(--mapex-surface-elevated) 0%, var(--mapex-card-border) 100%);
}

.phase-card {
  overflow: hidden;
}

.phase-header {
  transition: background-color 0.2s ease;

  &:hover {
    background-color: var(--mapex-submenu-bg);
  }
}

.stage-card {
  transition: var(--mapex-transition-base);
  height: 100%;

  &:hover {
    box-shadow: var(--mapex-shadow-sm);
  }

  &--disabled {
    opacity: 0.6;
    background-color: var(--mapex-surface-bg);
  }
}

.ellipsis {
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
</style>
