<template>
  <q-form ref="formRef" greedy>
    <div class="q-mb-md">
      <div class="text-subtitle1 text-weight-medium q-mb-xs">
        <q-icon name="wifi" color="primary" class="q-mr-xs" />
        {{ t.steps.step4.title.value }}
      </div>
      <div class="text-body2 text-grey-7">
        {{ t.steps.step4.subtitle.value }}
      </div>
    </div>

    <q-banner rounded class="bg-blue-1 text-blue-9 q-mb-md">
      <template #avatar>
        <q-icon name="info" color="blue" />
      </template>
      {{ t.steps.step4.banner.info.value }}
    </q-banner>

    <div class="row q-col-gutter-md">
      <div class="col-12">
        <q-select
          v-model="localData.protocol"
          outlined
          dense
          emit-value
          map-options
          class="rounded-borders"
          data-testid="asset-protocol-select"
          :label="t.steps.step4.fields.protocol.label.value + ' *'"
          :placeholder="t.steps.step4.fields.protocol.placeholder.value"
          :hint="t.steps.step4.fields.protocol.hint.value"
          :options="protocolOptions"
          option-label="label"
          option-value="value"
          option-disable="disable"
          :rules="[(val) => !!val || t.steps.step4.fields.protocol.required.value]"
          @update:model-value="handleProtocolChange"
        >
          <template #prepend>
            <q-icon name="router" color="primary" />
          </template>
        </q-select>
      </div>

      <!-- HTTP Info Banner -->
      <div v-if="localData.protocol === 'HTTP'" class="col-12">
        <q-banner rounded class="bg-grey-2 text-grey-8">
          <template #avatar>
            <q-icon name="http" color="grey-7" />
          </template>
          {{ t.steps.step4.banner.httpInfo.value }}
        </q-banner>
      </div>

      <!-- MQTT Configuration Fields -->
      <template v-if="localData.protocol === 'MQTT'">
        <div class="col-12">
          <q-banner rounded class="bg-green-1 text-green-9 q-mb-sm">
            <template #avatar>
              <q-icon name="mdi-lan-connect" color="green" />
            </template>
            {{ t.steps.step4.banner.mqttInfo.value }}
          </q-banner>
        </div>

        <!-- Client ID — derived from assetUUID, readonly. The broker
             uses the operator-chosen asset UUID as the canonical
             MQTT client identity, so we surface it here as a chip
             rather than a free-form input. -->
        <div class="col-12 col-md-6">
          <q-input
            :model-value="localData.mqttConfig.clientId"
            outlined
            dense
            readonly
            class="rounded-borders"
            data-testid="asset-mqtt-clientid-input"
            :label="t.steps.step4.fields.mqttClientId.label.value"
            :hint="t.steps.step4.fields.mqttClientId.hintDerived.value"
          >
            <template #prepend>
              <q-icon name="mdi-identifier" color="primary" />
            </template>
            <template #append>
              <q-icon name="lock" size="xs" color="grey-6" />
            </template>
          </q-input>
        </div>

        <!-- Auth Type — operator-chosen credential mode the broker
             enforces at CONNECT. The fields below render based on the
             selection. -->
        <div class="col-12 col-md-6">
          <q-select
            v-model="localData.mqttConfig.authType"
            outlined
            dense
            emit-value
            map-options
            class="rounded-borders"
            data-testid="asset-mqtt-authtype-select"
            :label="t.steps.step4.fields.authType.label.value + ' *'"
            :hint="t.steps.step4.fields.authType.hint.value"
            :options="authTypeOptions"
            option-label="label"
            option-value="value"
            :rules="[(val) => !!val || t.steps.step4.fields.authType.required.value]"
            @update:model-value="onAuthTypeChange"
          >
            <template #prepend>
              <q-icon name="security" color="primary" />
            </template>
          </q-select>
        </div>

        <!-- Password mode fields. Username is the bare assetUUID —
             globally unique and the broker plugin uses it directly as
             the cache lookup key. Password is the only field the
             operator types here. -->
        <template v-if="localData.mqttConfig.authType === MQTT_AUTH_TYPE_PASSWORD">
          <div class="col-12 col-md-6">
            <q-input
              :model-value="localData.mqttConfig.username"
              outlined
              dense
              readonly
              class="rounded-borders"
              data-testid="asset-mqtt-username-input"
              :label="t.steps.step4.fields.mqttUsername.label.value"
              :hint="t.steps.step4.fields.mqttUsername.hintDerived.value"
            >
              <template #prepend>
                <q-icon name="person" color="primary" />
              </template>
              <template #append>
                <q-icon name="lock" size="xs" color="grey-6" />
              </template>
            </q-input>
          </div>

          <div class="col-12 col-md-6">
            <q-input
              v-model="localData.mqttConfig.password"
              outlined
              dense
              :type="passwordVisible ? 'text' : 'password'"
              class="rounded-borders"
              data-testid="asset-mqtt-password-input"
              :label="t.steps.step4.fields.mqttPassword.label.value + (props.isEditMode ? '' : ' *')"
              :placeholder="props.isEditMode
                ? t.steps.step4.fields.mqttPassword.placeholderEdit.value
                : t.steps.step4.fields.mqttPassword.placeholder.value"
              :hint="props.isEditMode
                ? t.steps.step4.fields.mqttPassword.hintEdit.value
                : t.steps.step4.fields.mqttPassword.hint.value"
              :rules="mqttPasswordRules"
              autocomplete="new-password"
              @update:model-value="updateValue"
            >
              <template #prepend>
                <q-icon name="vpn_key" color="primary" />
              </template>
              <template #append>
                <q-icon
                  :name="passwordVisible ? 'visibility_off' : 'visibility'"
                  class="cursor-pointer"
                  @click="passwordVisible = !passwordVisible"
                />
              </template>
            </q-input>
          </div>
        </template>

        <!-- Cert mode banner. The cert is issued in the post-save
             dialog (POST /api/v1/mqtt_certs returns the PEM + key once;
             zip download starts automatically) so the operator does
             not enter any field here. -->
        <template v-if="localData.mqttConfig.authType === MQTT_AUTH_TYPE_CERT">
          <div class="col-12">
            <q-banner rounded class="cert-mode-banner">
              <template #avatar>
                <q-icon name="badge" color="primary" />
              </template>
              <div class="text-weight-medium q-mb-xs">
                {{ t.steps.step4.fields.authType.certBannerTitle.value }}
              </div>
              <div class="text-caption">
                {{ t.steps.step4.fields.authType.certBannerBody.value }}
              </div>
            </q-banner>
          </div>

          <!-- Cert TTL — operator-declared validity window. Backend
               clamps to [1 day, 10 years] and persists on the asset
               so each subsequent IssueCert uses the same TTL. -->
          <div class="col-12 col-md-3">
            <q-input
              v-model.number="certTTLValueModel"
              outlined
              dense
              type="number"
              min="1"
              class="rounded-borders"
              data-testid="asset-cert-ttl-value-input"
              :label="t.steps.step4.fields.certTTL.value.label.value + ' *'"
              :hint="t.steps.step4.fields.certTTL.value.hint.value"
              :rules="certTTLValueRules"
              @update:model-value="updateValue"
            >
              <template #prepend>
                <q-icon name="schedule" color="primary" />
              </template>
            </q-input>
          </div>
          <div class="col-12 col-md-3">
            <q-select
              v-model="certTTLUnitModel"
              outlined
              dense
              emit-value
              map-options
              class="rounded-borders"
              data-testid="asset-cert-ttl-unit-select"
              :label="t.steps.step4.fields.certTTL.unit.label.value + ' *'"
              :hint="t.steps.step4.fields.certTTL.unit.hint.value"
              :options="certTTLUnitOptions"
              option-label="label"
              option-value="value"
              :rules="[(val) => !!val || t.steps.step4.fields.certTTL.unit.required.value]"
              @update:model-value="updateValue"
            >
              <template #prepend>
                <q-icon name="timelapse" color="primary" />
              </template>
            </q-select>
          </div>
        </template>
      </template>

      <!-- Location Fields -->
      <div class="col-12 col-md-6">
        <q-input
          v-model.number="localData.latitude"
          outlined
          dense
          type="number"
          class="rounded-borders"
          data-testid="asset-latitude-input"
          :label="t.steps.step4.fields.latitude.label.value"
          :placeholder="t.steps.step4.fields.latitude.placeholder.value"
          :hint="t.steps.step4.fields.latitude.hint.value"
          @update:model-value="updateValue"
        >
          <template #prepend>
            <q-icon name="place" color="primary" />
          </template>
        </q-input>
      </div>
      <div class="col-12 col-md-6">
        <q-input
          v-model.number="localData.longitude"
          outlined
          dense
          type="number"
          class="rounded-borders"
          data-testid="asset-longitude-input"
          :label="t.steps.step4.fields.longitude.label.value"
          :placeholder="t.steps.step4.fields.longitude.placeholder.value"
          :hint="t.steps.step4.fields.longitude.hint.value"
          @update:model-value="updateValue"
        >
          <template #prepend>
            <q-icon name="place" color="primary" />
          </template>
        </q-input>
      </div>
    </div>
  </q-form>
</template>

<script setup lang="ts">
defineOptions({
  name: 'Step4Connectivity'
});

/** TYPE IMPORTS */
import type { Step4ConnectivityProps } from './interfaces/Step4Connectivity.interface';
import type { AssetFormData, SelectOption } from '../../interfaces';
import type { QForm } from 'quasar';

/** VUE IMPORTS */
import { ref, computed, reactive, watch } from 'vue';

/** COMPOSABLES */
import { useAddAssetTranslations } from '@src/composables/i18n/pages/assets/addAsset/useAddAssetTranslations';

/** LOCAL IMPORTS */
import { INITIAL_MQTT_CONFIG } from '../../constants';
import { MQTT_AUTH_TYPE_PASSWORD, MQTT_AUTH_TYPE_CERT, CERT_TTL_UNITS } from '../../interfaces/createEditAsset.interface';
import type { CertTTLUnit } from '../../interfaces/createEditAsset.interface';

const props = defineProps<Step4ConnectivityProps>();
const emit = defineEmits<{
  (e: 'update:modelValue', value: Partial<AssetFormData>): void;
}>();

/** REFS */
const formRef = ref<QForm | null>(null);
const passwordVisible = ref(false);

/** COMPOSABLES */
const t = useAddAssetTranslations();

/** STATE */
const localData = reactive({
  protocol: props.modelValue.protocol || 'HTTP',
  latitude: props.modelValue.latitude || null,
  longitude: props.modelValue.longitude || null,
  mqttConfig: { ...INITIAL_MQTT_CONFIG, ...props.modelValue.mqttConfig },
});

// Derive the canonical MQTT identity from the operator-chosen asset
// UUID (Step 1). Username equals the bare assetUUID — globally unique
// and used as-is by the broker plugin's cache lookup. Operators
// program the device with the value they see here.
function syncDerivedIdentity(): void {
  const assetUUID = props.modelValue.assetId ?? '';
  localData.mqttConfig.clientId = assetUUID;
  localData.mqttConfig.username = assetUUID;
}
syncDerivedIdentity();

/** WATCHERS */
watch(() => props.modelValue, (newVal) => {
  localData.protocol = newVal.protocol || 'HTTP';
  localData.latitude = newVal.latitude || null;
  localData.longitude = newVal.longitude || null;
  localData.mqttConfig = { ...INITIAL_MQTT_CONFIG, ...newVal.mqttConfig };
  syncDerivedIdentity();
}, { deep: true });

// Asset UUID changes upstream (Step 1) — the wizard owns the source of
// truth, so re-derive client/username whenever the operator edits it.
watch(() => props.modelValue.assetId, () => {
  syncDerivedIdentity();
  updateValue();
});

/** COMPUTED */
const protocolOptions = computed((): SelectOption[] => [
  { label: t.steps.step4.protocolOptions.http.value, value: 'HTTP', disable: false },
  { label: t.steps.step4.protocolOptions.mqtt.value, value: 'MQTT', disable: false },
  { label: t.steps.step4.protocolOptions.lorawan.value, value: 'LoRaWAN', disable: true },
]);

const authTypeOptions = computed(() => [
  { label: t.steps.step4.fields.authType.optionPassword.value, value: MQTT_AUTH_TYPE_PASSWORD },
  { label: t.steps.step4.fields.authType.optionCert.value, value: MQTT_AUTH_TYPE_CERT },
]);

// Cap on the total day count the platform accepts for a device cert.
// Mirrors the backend constant (MaxDeviceCertTTLDays = 3650 = 10 years)
// so the form catches bad combos before the request leaves the browser.
const CERT_TTL_MAX_DAYS = 3650;
const CERT_TTL_UNIT_TO_DAYS: Record<CertTTLUnit, number> = {
  day: 1,
  week: 7,
  month: 30,
  year: 365,
};

/**
 * Cert TTL options rendered in the dropdown. Labels are i18n-bound so
 * the same enum surfaces as "dia / semana / mês / ano" in pt-BR.
 */
const certTTLUnitOptions = computed(() => CERT_TTL_UNITS.map((unit) => ({
  label: t.steps.step4.fields.certTTL.units[unit].value,
  value: unit,
})));

/**
 * Two-way binding helpers for the Value+Unit pair. Persist on
 * `mqttConfig.certTTL` (an optional struct on the wire) and default
 * to `{1, year}` whenever the asset arrives without an explicit TTL.
 */
const certTTLValueModel = computed<number>({
  get: () => localData.mqttConfig.certTTL?.value ?? 1,
  set: (v: number) => {
    const current = localData.mqttConfig.certTTL ?? { value: 1, unit: 'year' as CertTTLUnit };
    localData.mqttConfig.certTTL = { ...current, value: v };
  },
});

const certTTLUnitModel = computed<CertTTLUnit>({
  get: () => localData.mqttConfig.certTTL?.unit ?? 'year',
  set: (u: CertTTLUnit) => {
    const current = localData.mqttConfig.certTTL ?? { value: 1, unit: 'year' as CertTTLUnit };
    localData.mqttConfig.certTTL = { ...current, unit: u };
  },
});

/**
 * Cert TTL rules. Required + min(1) on the value, plus a combined
 * upper bound on (value * unit-to-days) so the form blocks 11 year /
 * 4000 day attempts at the input layer instead of waiting for a 422
 * from the signer.
 */
const certTTLValueRules = computed(() => [
  (val: number | null) =>
    (typeof val === 'number' && val >= 1) || t.steps.step4.fields.certTTL.value.required.value,
  (val: number | null) => {
    const unit = certTTLUnitModel.value;
    const totalDays = (val ?? 0) * CERT_TTL_UNIT_TO_DAYS[unit];
    return totalDays <= CERT_TTL_MAX_DAYS || t.steps.step4.fields.certTTL.value.maxDays.value;
  },
]);

/**
 * Password rules. Required on create in password-mode (operator must
 * provision an initial credential); on edit it stays optional so a
 * blank value keeps the existing bcrypt hash. Non-empty values must
 * still clear the platform's bcrypt minimum.
 */
const mqttPasswordRules = computed(() => {
  const minLen = (val: string) =>
    !val || val.length >= 8 || t.steps.step4.fields.mqttPassword.minLength.value;
  if (localData.mqttConfig.authType !== MQTT_AUTH_TYPE_PASSWORD) {
    return [];
  }
  if (props.isEditMode) {
    return [minLen];
  }
  return [
    (val: string) => !!val || t.steps.step4.fields.mqttPassword.required.value,
    minLen,
  ];
});

/** FUNCTIONS */

/**
 * Handle protocol change event. Resets MQTT-specific state when
 * switching away from MQTT so the create payload does not carry stale
 * username / clientId / password for an HTTP-protocol asset. When
 * switching INTO MQTT re-derive client/username from the asset UUID.
 */
function handleProtocolChange(): void {
  if (localData.protocol !== 'MQTT' && localData.protocol !== 'LoRaWAN') {
    localData.mqttConfig = { ...INITIAL_MQTT_CONFIG };
  } else if (localData.protocol === 'MQTT') {
    syncDerivedIdentity();
  }
  updateValue();
}

/**
 * Handle auth-type change. Clears the password field when switching
 * to cert-mode so a stale plaintext value does not flow into the
 * payload (cert-mode rejects `password` on the backend validator).
 */
function onAuthTypeChange(): void {
  if (localData.mqttConfig.authType === MQTT_AUTH_TYPE_CERT) {
    localData.mqttConfig.password = '';
  }
  updateValue();
}

/**
 * Emit updated values to parent component. Client/username are
 * platform-derived so the parent always receives the canonical
 * lookup-key format the broker expects.
 */
function updateValue(): void {
  emit('update:modelValue', {
    protocol: localData.protocol,
    latitude: localData.latitude,
    longitude: localData.longitude,
    mqttConfig: localData.mqttConfig,
  });
}

defineExpose({
  formRef,
});
</script>

<style scoped lang="scss">
.rounded-borders {
  border-radius: var(--mapex-radius-md);
}

.cert-mode-banner {
  background: var(--mapex-surface-info-soft, #e3f2fd);
  color: var(--mapex-text-primary);
  border: 1px solid var(--mapex-border-info, #90caf9);
  border-radius: var(--mapex-radius-md);
}
</style>
