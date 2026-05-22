<script setup lang="ts">
defineOptions({
  name: 'HealthMonitoringSection'
});

/** TYPE IMPORTS */
import type { HealthMonitorFormConfig } from '../../interfaces';
import type { RouteGroupResponse } from '@mapexos/schemas';

/** VUE IMPORTS */
import { ref, reactive, watch, computed } from 'vue';
import { copyToClipboard } from 'quasar';

/** COMPONENTS */
import { BaseButton } from '@components/buttons';
import { SelectableChip } from '@components/chips';
import { RouteGroupSelectorDrawer } from '@components/drawers';
import { AppTooltip } from '@components/tooltips';

/** COMPOSABLES */
import { useTS } from '@utils/translation';

/** UTILS */
import { notifySuccess, notifyFail } from '@utils/alert/notify';

/** PROPS & EMITS */
const props = defineProps<{
  modelValue: HealthMonitorFormConfig;
  /** Asset UUID — when provided, the explicit-mode HTTP curl example shows
   *  the real value in the body; otherwise falls back to <assetUUID>. */
  assetUUID?: string;
  /** Asset protocol — MQTT-protocol assets are always monitored via broker
   *  presence advisories (CONNECT/DISCONNECT), so the enable toggle and
   *  the heartbeat-mode selector are both hidden. Operators still pick
   *  thresholds and the online/offline route groups. */
  protocol?: string;
}>();

const emit = defineEmits<{
  'update:modelValue': [value: Partial<HealthMonitorFormConfig>];
}>();

/** COMPOSABLES & STORES */
const ts = useTS({ capitalize: true });
const bp = 'pages.assets.assets.healthMonitoring';

/**
 * Router kinds permitted for healthStatus event source.
 * Mirrors validateHealthMonitorConfig (asset_helpers.go) and
 * HealthStatusAllowedRouterKinds (route_types.go).
 */
const HEALTH_ALLOWED_KINDS = ['trigger', 'workflow'] as const;

/** COMPUTED — protocol gating */

/**
 * MQTT-protocol assets are health-monitored exclusively via broker
 * presence (CONNECT / DISCONNECT advisories). The operator does NOT
 * pick a heartbeat mode and cannot disable monitoring — both fields
 * are forced (enabled=true, heartbeatMode='explicit') and the inputs
 * are hidden from the UI.
 */
const isMqtt = computed<boolean>(() => (props.protocol ?? '').toUpperCase() === 'MQTT');

/** STATE */
const showOfflineDrawer = ref(false);
const showOnlineDrawer = ref(false);

const localData = reactive({
  enabled: isMqtt.value ? true : props.modelValue.enabled,
  thresholdMinutes: props.modelValue.thresholdMinutes,
  requiredMisses: props.modelValue.requiredMisses,
  heartbeatMode: isMqtt.value ? 'explicit' : (props.modelValue.heartbeatMode ?? 'implicit'),
  offlineRouteGroupIds: props.modelValue.offlineRouteGroupIds || [],
  onlineRouteGroupIds: props.modelValue.onlineRouteGroupIds || [],
});

const selectedOfflineRouteGroups = ref<RouteGroupResponse[]>(
  props.modelValue.selectedOfflineRouteGroups || []
);
const selectedOnlineRouteGroups = ref<RouteGroupResponse[]>(
  props.modelValue.selectedOnlineRouteGroups || []
);

/** WATCHERS */
watch(() => props.modelValue, (newVal) => {
  localData.enabled = isMqtt.value ? true : newVal.enabled;
  localData.thresholdMinutes = newVal.thresholdMinutes;
  localData.requiredMisses = newVal.requiredMisses;
  localData.heartbeatMode = isMqtt.value ? 'explicit' : (newVal.heartbeatMode ?? 'implicit');
  localData.offlineRouteGroupIds = newVal.offlineRouteGroupIds || [];
  localData.onlineRouteGroupIds = newVal.onlineRouteGroupIds || [];
  if (newVal.selectedOfflineRouteGroups?.length) {
    selectedOfflineRouteGroups.value = newVal.selectedOfflineRouteGroups;
  }
  if (newVal.selectedOnlineRouteGroups?.length) {
    selectedOnlineRouteGroups.value = newVal.selectedOnlineRouteGroups;
  }
}, { deep: true });

// React to protocol changes (Step4 lets the operator switch MQTT↔HTTP).
// Switching INTO MQTT forces enabled=true + heartbeatMode='explicit'
// and re-emits so the parent's form state reflects the locked values.
watch(isMqtt, (nowMqtt) => {
  if (nowMqtt) {
    localData.enabled = true;
    localData.heartbeatMode = 'explicit';
    emitUpdate();
  }
});

/** FUNCTIONS */

/**
 * Emits updated form data to parent.
 */
function emitUpdate(): void {
  emit('update:modelValue', {
    ...localData,
    selectedOfflineRouteGroups: selectedOfflineRouteGroups.value,
    selectedOnlineRouteGroups: selectedOnlineRouteGroups.value,
  });
}

/**
 * Handles offline RouteGroup selection from drawer.
 * @param {RouteGroupResponse[]} routeGroups - Selected route groups
 */
function handleOfflineSelect(routeGroups: RouteGroupResponse[]): void {
  selectedOfflineRouteGroups.value = routeGroups;
  localData.offlineRouteGroupIds = routeGroups.map(rg => rg.id).filter(Boolean) as string[];
  emitUpdate();
}

/**
 * Handles online RouteGroup selection from drawer.
 * @param {RouteGroupResponse[]} routeGroups - Selected route groups
 */
function handleOnlineSelect(routeGroups: RouteGroupResponse[]): void {
  selectedOnlineRouteGroups.value = routeGroups;
  localData.onlineRouteGroupIds = routeGroups.map(rg => rg.id).filter(Boolean) as string[];
  emitUpdate();
}

/**
 * Removes a route group from the offline list.
 * @param {string} id - Route group ID to remove
 */
function removeOfflineRouteGroup(id: string): void {
  selectedOfflineRouteGroups.value = selectedOfflineRouteGroups.value.filter(rg => rg.id !== id);
  localData.offlineRouteGroupIds = selectedOfflineRouteGroups.value.map(rg => rg.id).filter(Boolean) as string[];
  emitUpdate();
}

/**
 * Removes a route group from the online list.
 * @param {string} id - Route group ID to remove
 */
function removeOnlineRouteGroup(id: string): void {
  selectedOnlineRouteGroups.value = selectedOnlineRouteGroups.value.filter(rg => rg.id !== id);
  localData.onlineRouteGroupIds = selectedOnlineRouteGroups.value.map(rg => rg.id).filter(Boolean) as string[];
  emitUpdate();
}

/** COMPUTED — heartbeat banner content */

/**
 * Curl example targeting the heartbeat HTTP endpoint. The body carries
 * the assetUUID; orgId and pathKey are derived server-side from the
 * resolved DataSource (never from the body).
 */
const httpCurlExample = computed<string>(() => {
  const uuid = props.assetUUID || '<assetUUID>';
  return `curl -X POST 'http://localhost:8080/api/v1/heartbeat?ds=<dataSourceId>' \\\n  -H 'x-api-key: <your-key>' \\\n  -H 'Content-Type: application/json' \\\n  -d '{"assetUUID":"${uuid}"}'`;
});

/**
 * Should the "monitor only" hint banner be shown? Per the model, it only
 * applies when BOTH route group lists are empty AND the form is in
 * monitoring mode (mode is irrelevant — the banner is about routing, not
 * heartbeat source). Keeping the gate on emptiness avoids the false
 * "monitor only" message that previously rendered even when route groups
 * were configured.
 */
const showMonitorOnlyHint = computed<boolean>(() => {
  return (
    localData.offlineRouteGroupIds.length === 0 &&
    localData.onlineRouteGroupIds.length === 0
  );
});

/**
 * heartbeatMode toggle options for q-btn-toggle.
 */
const heartbeatModeOptions = computed(() => [
  { label: ts(`${bp}.heartbeatMode.optionImplicit`), value: 'implicit' },
  { label: ts(`${bp}.heartbeatMode.optionExplicit`), value: 'explicit' },
]);

/** FUNCTIONS — copy helpers */

/**
 * Copies a string to the system clipboard and surfaces feedback via the
 * shared Notify util — success/fail toasts follow the platform pattern
 * (see StandardizedPayloadHelpModal for the canonical reference).
 *
 * @param {string} text - Plain-text payload to copy
 */
async function copyText(text: string): Promise<void> {
  try {
    await copyToClipboard(text);
    notifySuccess({
      message: ts(`${bp}.heartbeatMode.copySuccess`),
      timeout: 2000,
    });
  } catch {
    notifyFail({
      message: ts(`${bp}.heartbeatMode.copyFail`),
      timeout: 2000,
    });
  }
}
</script>

<template>
  <q-card flat bordered class="q-mt-md">
    <q-card-section>
      <div class="row items-center q-mb-sm">
        <q-icon name="health_and_safety" color="primary" size="sm" class="q-mr-sm" />
        <div class="text-subtitle1 text-weight-medium">{{ ts(`${bp}.title`) }}</div>
        <q-space />
        <q-toggle
          v-if="!isMqtt"
          v-model="localData.enabled"
          :label="localData.enabled ? ts(`${bp}.enabled`) : ts(`${bp}.disabled`)"
          color="primary"
          @update:model-value="emitUpdate"
        />
        <q-chip
          v-else
          color="positive"
          text-color="white"
          icon="check_circle"
          dense
          :label="ts(`${bp}.mqttAlwaysOn`)"
        />
      </div>

      <!-- MQTT presence banner: explains that MQTT assets are health-
           monitored exclusively via broker presence advisories, so the
           operator only picks thresholds and route groups. -->
      <q-banner
        v-if="isMqtt"
        class="mqtt-presence-banner q-mb-md"
        rounded
        dense
      >
        <template v-slot:avatar>
          <q-icon name="hub" color="primary" />
        </template>
        <div class="text-weight-medium q-mb-xs">{{ ts(`${bp}.mqttBanner.title`) }}</div>
        <div class="text-caption">{{ ts(`${bp}.mqttBanner.body`) }}</div>
      </q-banner>

      <q-slide-transition>
        <div v-if="localData.enabled">
          <q-separator class="q-mb-lg" />

          <!-- ============================================================ -->
          <!-- SECTION 1: Heartbeat source — HIDDEN for MQTT (broker only)  -->
          <!-- ============================================================ -->
          <div v-if="!isMqtt" class="health-section">
            <div class="health-section__label">{{ ts(`${bp}.sections.heartbeatSource`) }}</div>
            <div class="health-section__description">{{ ts(`${bp}.sections.heartbeatSourceDescription`) }}</div>

            <q-btn-toggle
              v-model="localData.heartbeatMode"
              :options="heartbeatModeOptions"
              toggle-color="primary"
              color="white"
              text-color="primary"
              unelevated
              spread
              no-caps
              @update:model-value="emitUpdate"
            />

            <!-- Implicit mode: same card treatment as explicit, just no channels -->
            <div v-if="localData.heartbeatMode === 'implicit'" class="mode-card q-mt-md">
              <div class="mode-card__header">
                <q-icon name="autorenew" color="primary" size="20px" />
                <div class="mode-card__heading">
                  <div class="mode-card__title">{{ ts(`${bp}.heartbeatMode.implicitTitle`) }}</div>
                  <div class="mode-card__subtitle">{{ ts(`${bp}.heartbeatMode.implicitHint`) }}</div>
                </div>
              </div>
            </div>

            <!-- Explicit mode: dedicated card with breathing space + channels -->
            <div v-else class="mode-card q-mt-md">
              <div class="mode-card__header">
                <q-icon name="wifi_tethering" color="primary" size="20px" />
                <div class="mode-card__heading">
                  <div class="mode-card__title">{{ ts(`${bp}.heartbeatMode.explicitTitle`) }}</div>
                  <div class="mode-card__subtitle">{{ ts(`${bp}.heartbeatMode.explicitHint`) }}</div>
                </div>
              </div>

              <div class="channel-card" role="group">
                <div class="channel-card__head">
                  <span class="channel-card__chip channel-card__chip--http">HTTP</span>
                  <span class="channel-card__title">{{ ts(`${bp}.heartbeatMode.httpCardTitle`) }}</span>
                  <q-space />
                  <BaseButton
                    icon="content_copy"
                    flat
                    round
                    dense
                    size="sm"
                    color="primary"
                    :aria-label="ts(`${bp}.heartbeatMode.copyAria`)"
                    @click="copyText(httpCurlExample)"
                  >
                    <AppTooltip :content="ts(`${bp}.heartbeatMode.copyAria`)" />
                  </BaseButton>
                </div>
                <div class="channel-card__body">
                  <code>{{ httpCurlExample }}</code>
                </div>
              </div>
            </div>
          </div>

          <!-- ============================================================ -->
          <!-- SECTION 2: Detection thresholds                              -->
          <!-- ============================================================ -->
          <div class="health-section">
            <div class="health-section__label">{{ ts(`${bp}.sections.detectionThresholds`) }}</div>
            <div class="health-section__description">{{ ts(`${bp}.sections.detectionThresholdsDescription`) }}</div>

            <div class="row q-col-gutter-md">
              <div class="col-12 col-sm-6">
                <q-input
                  v-model.number="localData.thresholdMinutes"
                  type="number"
                  :label="ts(`${bp}.threshold`)"
                  :hint="ts(`${bp}.thresholdHint`)"
                  :min="10"
                  outlined
                  dense
                  :suffix="ts(`${bp}.thresholdUnit`)"
                  @update:model-value="emitUpdate"
                />
              </div>
              <div class="col-12 col-sm-6">
                <q-input
                  v-model.number="localData.requiredMisses"
                  type="number"
                  :label="ts(`${bp}.requiredMisses`)"
                  :hint="ts(`${bp}.requiredMissesHint`)"
                  :min="1"
                  outlined
                  dense
                  @update:model-value="emitUpdate"
                />
              </div>
            </div>
          </div>

          <!-- ============================================================ -->
          <!-- SECTION 3: Transition routing                                -->
          <!-- ============================================================ -->
          <div class="health-section">
            <div class="health-section__label">{{ ts(`${bp}.sections.transitionRouting`) }}</div>
            <div class="health-section__description">{{ ts(`${bp}.sections.transitionRoutingDescription`) }}</div>

            <q-banner v-if="showMonitorOnlyHint" class="bg-grey-2 text-grey-9 q-mb-md text-caption" rounded dense>
              <template v-slot:avatar>
                <q-icon name="info" color="primary" />
              </template>
              {{ ts(`${bp}.monitorOnlyHint`) }}
            </q-banner>

            <div class="row q-col-gutter-md">
            <!-- Offline Route Groups -->
            <div class="col-12 col-sm-6">
              <q-card flat bordered class="health-route-card health-route-card--offline full-height">
                <q-card-section class="q-pb-none">
                  <div class="row items-center q-mb-xs">
                    <q-icon name="mdi-wifi-off" color="negative" size="xs" class="q-mr-sm" />
                    <div class="text-subtitle2 text-weight-bold">{{ ts(`${bp}.offlineRouteGroups`) }}</div>
                  </div>
                  <div class="text-caption text-grey">{{ ts(`${bp}.offlineRouteGroupsHint`) }}</div>
                </q-card-section>

                <q-card-section>
                  <div v-if="selectedOfflineRouteGroups.length" class="row q-gutter-sm q-mb-sm">
                    <SelectableChip
                      v-for="rg in selectedOfflineRouteGroups"
                      :key="rg.id ?? ''"
                      :label="rg.name"
                      icon="alt_route"
                      color="negative"
                      removable
                      @remove="removeOfflineRouteGroup(rg.id!)"
                    />
                  </div>
                  <div v-else class="text-caption text-grey-6 text-italic q-mb-sm">
                    {{ ts(`${bp}.offlineRouteGroupsEmpty`) }}
                  </div>

                  <q-btn
                    outline
                    dense
                    no-caps
                    icon="add"
                    color="negative"
                    size="sm"
                    :label="ts(`${bp}.addRouteGroup`)"
                    class="full-width"
                    @click="showOfflineDrawer = true"
                  />
                </q-card-section>
              </q-card>
            </div>

            <!-- Online Route Groups -->
            <div class="col-12 col-sm-6">
              <q-card flat bordered class="health-route-card health-route-card--online full-height">
                <q-card-section class="q-pb-none">
                  <div class="row items-center q-mb-xs">
                    <q-icon name="mdi-wifi-check" color="positive" size="xs" class="q-mr-sm" />
                    <div class="text-subtitle2 text-weight-bold">{{ ts(`${bp}.onlineRouteGroups`) }}</div>
                  </div>
                  <div class="text-caption text-grey">{{ ts(`${bp}.onlineRouteGroupsHint`) }}</div>
                </q-card-section>

                <q-card-section>
                  <div v-if="selectedOnlineRouteGroups.length" class="row q-gutter-sm q-mb-sm">
                    <SelectableChip
                      v-for="rg in selectedOnlineRouteGroups"
                      :key="rg.id ?? ''"
                      :label="rg.name"
                      icon="alt_route"
                      color="positive"
                      removable
                      @remove="removeOnlineRouteGroup(rg.id!)"
                    />
                  </div>
                  <div v-else class="text-caption text-grey-6 text-italic q-mb-sm">
                    {{ ts(`${bp}.onlineRouteGroupsEmpty`) }}
                  </div>

                  <q-btn
                    outline
                    dense
                    no-caps
                    icon="add"
                    color="positive"
                    size="sm"
                    :label="ts(`${bp}.addRouteGroup`)"
                    class="full-width"
                    @click="showOnlineDrawer = true"
                  />
                </q-card-section>
              </q-card>
            </div>
            </div>
          </div>
        </div>
      </q-slide-transition>
    </q-card-section>

    <!-- Route Group Selector Drawers -->
    <RouteGroupSelectorDrawer
      v-model="showOfflineDrawer"
      :multi-select="true"
      :allowed-router-kinds="[...HEALTH_ALLOWED_KINDS]"
      @select="handleOfflineSelect"
    />
    <RouteGroupSelectorDrawer
      v-model="showOnlineDrawer"
      :multi-select="true"
      :allowed-router-kinds="[...HEALTH_ALLOWED_KINDS]"
      @select="handleOnlineSelect"
    />
  </q-card>
</template>

<style lang="scss" scoped>
.health-route-card {
  border-radius: var(--mapex-radius-md);
  transition: var(--mapex-transition-base);

  &--offline {
    border-left: 3px solid var(--q-negative);
  }

  &--online {
    border-left: 3px solid var(--q-positive);
  }
}

.health-section {
  margin-bottom: 28px;

  &:last-child {
    margin-bottom: 0;
  }

  &__label {
    font-size: 0.65rem;
    font-weight: 700;
    letter-spacing: 0.5px;
    color: var(--mapex-text-secondary);
    text-transform: uppercase;
    margin-bottom: 4px;
  }

  &__description {
    font-size: 0.8rem;
    color: var(--mapex-text-secondary);
    margin-bottom: 12px;
  }
}

.heartbeat-mode-banner {
  background: var(--mapex-surface-elevated);
  color: var(--mapex-text-primary);
  border: 1px solid var(--mapex-card-border);
  border-left: 3px solid var(--q-primary);
}

.mqtt-presence-banner {
  background: var(--mapex-surface-elevated);
  color: var(--mapex-text-primary);
  border: 1px solid var(--mapex-card-border);
  border-left: 3px solid var(--q-primary);
}

/* Mode card — shared frame for both 'implicit' and 'explicit' heartbeat
 * choices so they have identical breathing room. Explicit adds channel
 * sub-cards inside; implicit has only the header. */
.mode-card {
  background: var(--mapex-surface-elevated);
  border: 1px solid var(--mapex-card-border);
  border-left: 3px solid var(--q-primary);
  border-radius: var(--mapex-radius-md);
  padding: 16px 18px;

  &__header {
    display: flex;
    align-items: flex-start;
    gap: 12px;
  }

  /* When channel-cards follow the header, give the header bottom margin so
   * the children don't snug up against it. Implicit has no children, so no
   * extra margin needed. */
  & .channel-card:first-of-type {
    margin-top: 16px;
  }

  &__heading {
    flex: 1 1 auto;
    min-width: 0;
  }

  &__title {
    font-size: 0.875rem;
    font-weight: 600;
    color: var(--mapex-text-primary);
    line-height: 1.3;
  }

  &__subtitle {
    font-size: 0.75rem;
    color: var(--mapex-text-secondary);
    margin-top: 2px;
    line-height: 1.4;
  }
}

/* Per-channel card — protocol chip + copy button up top, code block below
 * with full width for horizontal scroll without elbow-jostling the button. */
.channel-card {
  background: var(--mapex-surface);
  border: 1px solid var(--mapex-card-border);
  border-radius: var(--mapex-radius-sm);
  padding: 10px 12px;

  & + & {
    margin-top: 10px;
  }

  &__head {
    display: flex;
    align-items: center;
    gap: 8px;
    margin-bottom: 8px;
  }

  &__chip {
    font-size: 0.65rem;
    font-weight: 700;
    padding: 3px 9px;
    border-radius: var(--mapex-radius-sm);
    letter-spacing: 0.6px;
    line-height: 1.2;
    color: var(--mapex-text-on-primary, white);

    &--mqtt {
      background: #6750A4; /* purple-ish, distinguishable */
    }
    &--http {
      background: var(--q-primary);
    }
  }

  &__title {
    font-size: 0.8rem;
    font-weight: 500;
    color: var(--mapex-text-primary);
  }

  &__body {
    background: var(--mapex-surface-elevated);
    border-radius: var(--mapex-radius-sm);
    padding: 8px 10px;
    overflow-x: auto;
    overflow-y: hidden;
    white-space: nowrap;

    code {
      font-family: ui-monospace, SFMono-Regular, Menlo, monospace;
      font-size: 0.75rem;
      color: var(--mapex-text-primary);
      background: transparent;
    }

    &::-webkit-scrollbar {
      height: 4px;
    }
    &::-webkit-scrollbar-thumb {
      background: var(--mapex-card-border);
      border-radius: 2px;
    }

    /* Variant for prose hints (e.g. MQTT presence "nothing to configure" message)
     * — wraps long text instead of forcing a horizontal scroll. */
    &--text {
      white-space: normal;
      font-size: 0.78rem;
      color: var(--mapex-text-secondary);
      line-height: 1.45;
    }
  }
}

.text-mono {
  font-family: ui-monospace, SFMono-Regular, Menlo, monospace;
  font-size: 0.75rem;
  color: var(--mapex-text-secondary);
}
</style>
