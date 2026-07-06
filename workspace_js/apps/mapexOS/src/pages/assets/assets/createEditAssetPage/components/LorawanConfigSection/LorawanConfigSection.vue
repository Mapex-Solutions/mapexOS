<script setup lang="ts">
defineOptions({
  name: 'LorawanConfigSection'
});

/** TYPE IMPORTS */
import type { LorawanFormConfig } from '../../interfaces';
import type { FrequencyPlanOption } from '../../constants/frequencyPlans.constant';

/** VUE IMPORTS */
import { reactive, computed, watch } from 'vue';

/** COMPOSABLES */
import { useTS } from '@utils/translation';

/** LOCAL IMPORTS (constants and handlers ONLY) */
import { LORAWAN_FREQUENCY_PLANS } from '../../constants/frequencyPlans.constant';
import { LORAWAN_KIND_DEVICE, LORAWAN_KIND_GATEWAY, CERT_TTL_UNITS } from '../../interfaces/createEditAsset.interface';

/** PROPS & EMITS */
const props = defineProps<{
  modelValue: LorawanFormConfig;
  /** Asset UUID from Step 1 — for both kinds this IS the LoRaWAN EUI (the
   *  contract keys the asset by its EUI); shown disabled so it cannot diverge. */
  assetUUID?: string;
}>();

const emit = defineEmits<{
  'update:modelValue': [value: Partial<LorawanFormConfig>];
}>();

/** COMPOSABLES & STORES */
const ts = useTS({ capitalize: true });
const bp = 'pages.assets.assets.lorawanConfig';

/** STATE — local mirror, re-emitted on every change so the parent owns the truth */
const localData = reactive<LorawanFormConfig>({
  ...props.modelValue,
  gateway: { ...props.modelValue.gateway },
});

/** WATCHERS */
watch(
  () => props.modelValue,
  (newVal) => {
    Object.assign(localData, newVal);
    localData.gateway = { ...newVal.gateway };
  },
  { deep: true }
);

/** COMPUTED — kind / activation / version gating */

/** A LoRaWAN asset is either a gateway (radio infrastructure) or an end-device
 *  (sensor). The kind drives the entire field set and the backend cold-config
 *  path (the LNS rejects a gateway connect whose projection kind != gateway). */
const isDevice = computed<boolean>(() => localData.kind === LORAWAN_KIND_DEVICE);
const isGateway = computed<boolean>(() => localData.kind === LORAWAN_KIND_GATEWAY);

/** OTAA derives session keys at join (appKey/nwkKey); ABP ships fixed session
 *  keys (devAddr/nwkSKey/appSKey). Only one key set is collected at a time. */
const isOtaa = computed<boolean>(() => localData.activation === 'otaa');

/** NwkKey only exists on LoRaWAN 1.1 — earlier MAC versions derive it from AppKey. */
const isMac11 = computed<boolean>(() => localData.macVersion === '1.1');

/** The asset UUID is the EUI for both kinds; the hint clarifies which. */
const euiHint = computed<string>(() =>
  isGateway.value ? ts(`${bp}.eui.hintGateway`) : ts(`${bp}.eui.hintDevice`)
);

/** OPTIONS */

const classOptions = computed(() => [
  { label: ts(`${bp}.device.classOptions.a`), value: 'A' },
  { label: ts(`${bp}.device.classOptions.b`), value: 'B' },
  { label: ts(`${bp}.device.classOptions.c`), value: 'C' },
]);

const macVersionOptions = ['1.0.2', '1.0.3', '1.0.4', '1.1'].map((v) => ({ label: v, value: v }));

const activationOptions = computed(() => [
  { label: ts(`${bp}.device.activationOptions.otaa`), value: 'otaa' },
  { label: ts(`${bp}.device.activationOptions.abp`), value: 'abp' },
]);

const authModeOptions = computed(() => [
  { label: ts(`${bp}.gateway.authModeOptions.cert`), value: 'cert' },
  { label: ts(`${bp}.gateway.authModeOptions.key`), value: 'key' },
  { label: ts(`${bp}.gateway.authModeOptions.eui`), value: 'eui' },
]);

const certTtlUnitOptions = computed(() =>
  CERT_TTL_UNITS.map((u) => ({ label: ts(`${bp}.gateway.certTtlUnits.${u}`), value: u }))
);

/** A device's region and a gateway's frequency plan both come from the plans
 *  the LNS embeds. The label is the human-readable name (never the raw id);
 *  the id is the submitted value. Gateways only accept gateway-flagged plans. */
const planToOption = (p: FrequencyPlanOption) => ({ label: p.name, value: p.id });
const devicePlanOptions = computed(() => LORAWAN_FREQUENCY_PLANS.map(planToOption));
const gatewayPlanOptions = computed(() =>
  LORAWAN_FREQUENCY_PLANS.filter((p) => p.gateways).map(planToOption)
);

/** FUNCTIONS */

/**
 * Re-emits the full local config to the parent so the wizard's form state and
 * Step 5 review reflect every edit.
 */
function emitUpdate(): void {
  emit('update:modelValue', { ...localData, gateway: { ...localData.gateway } });
}

/**
 * Generates `bytes` of cryptographically-random data as an uppercase hex
 * string (used to mint OTAA/ABP secret keys client-side).
 * @param {number} bytes - Number of random bytes to produce.
 * @returns {string} Uppercase hex, length `bytes * 2`.
 */
function randomHex(bytes: number): string {
  const buf = new Uint8Array(bytes);
  crypto.getRandomValues(buf);
  return Array.from(buf, (b) => b.toString(16).padStart(2, '0')).join('').toUpperCase();
}

/**
 * Fills a 16-byte secret key field with a fresh random value and re-emits.
 * The operator can still paste an existing key instead of generating one.
 * @param {'appKey' | 'nwkKey' | 'nwkSKey' | 'appSKey'} field - Key to fill.
 */
function generateKey(field: 'appKey' | 'nwkKey' | 'nwkSKey' | 'appSKey'): void {
  localData[field] = randomHex(16);
  emitUpdate();
}

/**
 * Fills the gateway Basics Station token with a fresh 32-byte random value and
 * re-emits. The operator can still paste an existing token instead.
 */
function generateGatewayKey(): void {
  localData.gateway.apiKey = randomHex(32);
  emitUpdate();
}
</script>

<template>
  <div class="lorawan-config">
    <!-- EUI — the asset UUID from Step 1 IS the EUI for both kinds (disabled) -->
    <div class="lorawan-config__section">
      <q-input
        :model-value="assetUUID || '—'"
        :label="ts(`${bp}.eui.label`)"
        :hint="euiHint"
        outlined
        dense
        disable
      />
    </div>

    <!-- ====================== DEVICE (sensor) ====================== -->
    <template v-if="isDevice">
      <q-separator class="lorawan-config__divider" />
      <div class="lorawan-config__group">
        <div class="lorawan-config__header">{{ ts(`${bp}.sections.identity`) }}</div>
        <div class="row q-col-gutter-md">
          <div class="col-12 col-md-6">
            <q-select
              v-model="localData.region"
              :options="devicePlanOptions"
              :label="ts(`${bp}.device.region`)"
              outlined
              dense
              emit-value
              map-options
              use-input
              input-debounce="0"
              @update:model-value="emitUpdate"
            />
          </div>
          <div v-if="isOtaa" class="col-12 col-md-6">
            <q-input
              v-model="localData.joinEui"
              :label="ts(`${bp}.device.joinEui`)"
              outlined
              dense
              hint="8 bytes (16 hex)"
              @update:model-value="emitUpdate"
            />
          </div>
        </div>
      </div>

      <q-separator class="lorawan-config__divider" />
      <div class="lorawan-config__group">
        <div class="lorawan-config__header">{{ ts(`${bp}.sections.profile`) }}</div>
        <div class="row q-col-gutter-md">
          <div class="col-12 col-md-4">
            <q-select
              v-model="localData.class"
              :options="classOptions"
              :label="ts(`${bp}.device.class`)"
              outlined
              dense
              emit-value
              map-options
              @update:model-value="emitUpdate"
            />
          </div>
          <div class="col-12 col-md-4">
            <q-select
              v-model="localData.macVersion"
              :options="macVersionOptions"
              :label="ts(`${bp}.device.macVersion`)"
              outlined
              dense
              emit-value
              map-options
              @update:model-value="emitUpdate"
            />
          </div>
          <div class="col-12 col-md-4">
            <q-input
              v-model="localData.phyVersion"
              :label="ts(`${bp}.device.phyVersion`)"
              outlined
              dense
              @update:model-value="emitUpdate"
            />
          </div>
        </div>
      </div>

      <q-separator class="lorawan-config__divider" />
      <div class="lorawan-config__group">
        <div class="lorawan-config__header">{{ ts(`${bp}.sections.activation`) }}</div>
        <div class="row q-col-gutter-md">
          <div class="col-12 col-md-6">
            <q-select
              v-model="localData.activation"
              :options="activationOptions"
              :label="ts(`${bp}.device.activation`)"
              outlined
              dense
              emit-value
              map-options
              @update:model-value="emitUpdate"
            />
          </div>
        </div>
        <div class="row q-col-gutter-md q-mt-xs">
          <!-- OTAA root keys -->
          <template v-if="isOtaa">
            <div class="col-12 col-md-6">
              <q-input
                v-model="localData.appKey"
                :label="ts(`${bp}.device.appKey`)"
                outlined
                dense
                hint="16 bytes (32 hex)"
                autocomplete="off"
                @update:model-value="emitUpdate"
              >
                <template #append>
                  <q-btn
                    flat
                    dense
                    round
                    icon="autorenew"
                    color="primary"
                    @click="generateKey('appKey')"
                  >
                    <q-tooltip>{{ ts(`${bp}.device.generate`) }}</q-tooltip>
                  </q-btn>
                </template>
              </q-input>
            </div>
            <div v-if="isMac11" class="col-12 col-md-6">
              <q-input
                v-model="localData.nwkKey"
                :label="ts(`${bp}.device.nwkKey`)"
                outlined
                dense
                hint="16 bytes (32 hex)"
                autocomplete="off"
                @update:model-value="emitUpdate"
              >
                <template #append>
                  <q-btn
                    flat
                    dense
                    round
                    icon="autorenew"
                    color="primary"
                    @click="generateKey('nwkKey')"
                  >
                    <q-tooltip>{{ ts(`${bp}.device.generate`) }}</q-tooltip>
                  </q-btn>
                </template>
              </q-input>
            </div>
          </template>

          <!-- ABP fixed session keys -->
          <template v-else>
            <div class="col-12 col-md-4">
              <q-input
                v-model="localData.devAddr"
                :label="ts(`${bp}.device.devAddr`)"
                outlined
                dense
                hint="4 bytes (8 hex)"
                @update:model-value="emitUpdate"
              />
            </div>
            <div class="col-12 col-md-4">
              <q-input
                v-model="localData.nwkSKey"
                :label="ts(`${bp}.device.nwkSKey`)"
                outlined
                dense
                hint="16 bytes (32 hex)"
                autocomplete="off"
                @update:model-value="emitUpdate"
              >
                <template #append>
                  <q-btn
                    flat
                    dense
                    round
                    icon="autorenew"
                    color="primary"
                    @click="generateKey('nwkSKey')"
                  >
                    <q-tooltip>{{ ts(`${bp}.device.generate`) }}</q-tooltip>
                  </q-btn>
                </template>
              </q-input>
            </div>
            <div class="col-12 col-md-4">
              <q-input
                v-model="localData.appSKey"
                :label="ts(`${bp}.device.appSKey`)"
                outlined
                dense
                hint="16 bytes (32 hex)"
                autocomplete="off"
                @update:model-value="emitUpdate"
              >
                <template #append>
                  <q-btn
                    flat
                    dense
                    round
                    icon="autorenew"
                    color="primary"
                    @click="generateKey('appSKey')"
                  >
                    <q-tooltip>{{ ts(`${bp}.device.generate`) }}</q-tooltip>
                  </q-btn>
                </template>
              </q-input>
            </div>
          </template>
        </div>
      </div>
    </template>

    <!-- ====================== GATEWAY ====================== -->
    <template v-if="isGateway">
      <q-separator class="lorawan-config__divider" />
      <div class="lorawan-config__group">
        <div class="lorawan-config__header">{{ ts(`${bp}.sections.frequency`) }}</div>
        <div class="row q-col-gutter-md">
          <div class="col-12">
            <q-select
              v-model="localData.gateway.frequencyPlanId"
              :options="gatewayPlanOptions"
              :label="ts(`${bp}.gateway.frequencyPlanId`)"
              outlined
              dense
              emit-value
              map-options
              use-input
              input-debounce="0"
              @update:model-value="emitUpdate"
            />
          </div>
          <div class="col-12">
            <q-select
              v-model="localData.gateway.frequencyPlanIds"
              :options="gatewayPlanOptions"
              :label="ts(`${bp}.gateway.frequencyPlanIds`)"
              :hint="ts(`${bp}.gateway.frequencyPlanIdsHint`)"
              outlined
              dense
              multiple
              emit-value
              map-options
              use-chips
              use-input
              input-debounce="0"
              @update:model-value="emitUpdate"
            />
          </div>
        </div>
      </div>

      <q-separator class="lorawan-config__divider" />
      <div class="lorawan-config__group">
        <div class="lorawan-config__header">{{ ts(`${bp}.sections.authentication`) }}</div>
        <div class="row q-col-gutter-md">
          <div class="col-12 col-md-6">
            <q-select
              v-model="localData.gateway.authMode"
              :options="authModeOptions"
              :label="ts(`${bp}.gateway.authMode`)"
              outlined
              dense
              emit-value
              map-options
              @update:model-value="emitUpdate"
            />
          </div>
          <template v-if="localData.gateway.authMode === 'cert' && localData.gateway.certTTL">
            <div class="col-6 col-md-3">
              <q-input
                v-model.number="localData.gateway.certTTL.value"
                type="number"
                :label="ts(`${bp}.gateway.certTtlValue`)"
                outlined
                dense
                min="1"
                @update:model-value="emitUpdate"
              />
            </div>
            <div class="col-6 col-md-3">
              <q-select
                v-model="localData.gateway.certTTL.unit"
                :options="certTtlUnitOptions"
                :label="ts(`${bp}.gateway.certTtlUnit`)"
                outlined
                dense
                emit-value
                map-options
                @update:model-value="emitUpdate"
              />
            </div>
          </template>
          <div v-if="localData.gateway.authMode === 'key'" class="col-12 col-md-6">
            <q-input
              v-model="localData.gateway.apiKey"
              :label="ts(`${bp}.gateway.apiKey`)"
              outlined
              dense
              hint="32 bytes (64 hex)"
              autocomplete="off"
              @update:model-value="emitUpdate"
            >
              <template #append>
                <q-btn
                  flat
                  dense
                  round
                  icon="autorenew"
                  color="primary"
                  @click="generateGatewayKey"
                >
                  <q-tooltip>{{ ts(`${bp}.device.generate`) }}</q-tooltip>
                </q-btn>
              </template>
            </q-input>
          </div>
        </div>
      </div>
    </template>
  </div>
</template>

<style scoped lang="scss">
.lorawan-config {
  display: flex;
  flex-direction: column;
  gap: var(--mapex-spacing-lg);

  &__section {
    display: flex;
    flex-direction: column;
    gap: var(--mapex-spacing-sm);
  }

  &__group {
    display: flex;
    flex-direction: column;
    gap: var(--mapex-spacing-sm);
  }

  &__description {
    font-size: var(--mapex-font-sm);
    color: var(--mapex-text-secondary);
  }

  &__header {
    font-weight: var(--mapex-font-weight-semibold);
    font-size: var(--mapex-font-md);
    color: var(--mapex-text-primary);
  }

  &__divider {
    margin-top: var(--mapex-spacing-sm);
  }
}
</style>
