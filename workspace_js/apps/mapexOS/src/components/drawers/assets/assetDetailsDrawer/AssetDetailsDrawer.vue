<template>
  <!-- Invisible backdrop for click outside detection -->
  <Teleport to="body">
    <div
      v-if="modelValue"
      class="drawer-backdrop"
      @click="close"
    />
  </Teleport>

  <q-drawer
    overlay
    bordered
    side="right"
    :model-value="modelValue"
    :width="450"
    @update:model-value="emit('update:modelValue', $event)"
    @keydown.esc="close"
  >
    <!-- Header -->
    <q-toolbar class="drawer-header">
      <q-icon name="devices" size="sm" class="q-mr-sm" color="primary" />
      <q-toolbar-title class="text-weight-medium">{{ t.drawer.title.value }}</q-toolbar-title>

      <q-btn flat round dense icon="close" class="drawer-close-btn" @click="close">
        <AppTooltip :content="t.drawer.close.value" />
      </q-btn>
    </q-toolbar>

    <q-separator />

    <!-- Content -->
    <div class="drawer-content">
      <q-scroll-area class="fit">
        <!-- Loading State -->
        <div v-if="loading" class="q-pa-lg text-center">
          <q-spinner size="3em" class="q-mb-md" color="primary" />
          <div class="text-grey-7">{{ t.drawer.loading.value }}</div>
        </div>

        <!-- Error State -->
        <div v-else-if="error" class="q-pa-lg">
          <q-banner rounded class="bg-negative text-white">
            <template #avatar>
              <q-icon name="error" color="white" />
            </template>
            {{ t.drawer.error.value }}
          </q-banner>
        </div>

        <!-- Asset Data -->
        <div v-else-if="asset" class="q-px-md q-py-lg">

          <!-- Basic Information Section -->
          <div class="section q-mb-md">
            <div class="section-header">
              <q-icon name="info" color="primary" size="sm" class="q-mr-sm" />
              <span class="text-subtitle1 text-weight-medium">{{ t.drawer.sections.basicInfo.value }}</span>
            </div>
            <q-separator class="q-my-sm" />

            <!-- Name (full width) -->
            <div class="field-row">
              <div class="field-label">{{ t.drawer.fields.name.value }}</div>
              <div class="field-value text-weight-medium">{{ asset?.name || '-' }}</div>
            </div>

            <!-- Description (full width) -->
            <div class="field-row">
              <div class="field-label">{{ t.drawer.fields.description.value }}</div>
              <div class="field-value text-grey-8">
                {{ asset?.description || t.drawer.empty.description.value }}
              </div>
            </div>

            <!-- UUID & Status (2 columns) -->
            <div class="row q-col-gutter-sm">
              <div class="col-8">
                <div class="field-row">
                  <div class="field-label">{{ t.drawer.fields.uuid.value }}</div>
                  <div class="field-value uuid-field-value">
                    <DetailChip
                      class="uuid-chip"
                      icon="fingerprint"
                      color="default"
                      size="sm"
                      :label="asset?.assetUUID || '-'"
                    />
                    <AppTooltip v-if="asset?.assetUUID" :content="asset.assetUUID" />
                  </div>
                </div>
              </div>
              <div class="col-4">
                <div class="field-row">
                  <div class="field-label">{{ t.drawer.fields.status.value }}</div>
                  <div class="field-value">
                    <DetailChip
                      :color="asset?.enabled ? 'positive' : 'negative'"
                      size="sm"
                      :label="asset?.enabled ? t.status.active.value.toUpperCase() : t.status.inactive.value.toUpperCase()"
                    />
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- Configuration Section -->
          <div class="section q-mb-md">
            <div class="section-header">
              <q-icon name="settings" color="primary" size="sm" class="q-mr-sm" />
              <span class="text-subtitle1 text-weight-medium">{{ t.drawer.sections.configuration.value }}</span>
            </div>
            <q-separator class="q-my-sm" />

            <!-- Template (full width) -->
            <div class="field-row q-mb-md">
              <div class="field-label">{{ t.drawer.fields.template.value }}</div>
              <div class="field-value">{{ asset?.assetTemplateName || '-' }}</div>
            </div>

            <!-- Category & Manufacturer (2 columns) -->
            <div class="row q-col-gutter-sm">
              <div class="col-6">
                <div class="field-row">
                  <div class="field-label">{{ t.drawer.fields.category.value }}</div>
                  <div class="field-value">{{ asset?.categoryName || '-' }}</div>
                </div>
              </div>
              <div class="col-6">
                <div class="field-row">
                  <div class="field-label">{{ t.drawer.fields.manufacturer.value }}</div>
                  <div class="field-value">{{ asset?.manufacturerName || '-' }}</div>
                </div>
              </div>
            </div>

            <!-- Model & Version (2 columns) -->
            <div class="row q-col-gutter-sm">
              <div class="col-6">
                <div class="field-row">
                  <div class="field-label">{{ t.drawer.fields.model.value }}</div>
                  <div class="field-value">{{ asset?.modelName || '-' }}</div>
                </div>
              </div>
              <div class="col-6">
                <div class="field-row">
                  <div class="field-label">{{ t.drawer.fields.version.value }}</div>
                  <div class="field-value">{{ asset?.version || '-' }}</div>
                </div>
              </div>
            </div>
          </div>

          <!-- Protocol Section -->
          <div class="section q-mb-md">
            <div class="section-header">
              <q-icon name="cable" color="primary" size="sm" class="q-mr-sm" />
              <span class="text-subtitle1 text-weight-medium">{{ t.drawer.sections.protocol.value }}</span>
            </div>
            <q-separator class="q-my-sm" />

            <div class="field-row">
              <div class="field-label">{{ t.drawer.fields.protocolType.value }}</div>
              <div class="field-value">
                <DetailChip
                  icon="cable"
                  :color="getProtocolColorName(asset?.protocol?.type)"
                  size="sm"
                  :label="asset?.protocol?.type?.toUpperCase() || '-'"
                />
              </div>
            </div>

            <!-- MQTT Configuration -->
            <template v-if="asset?.protocol?.type === 'mqtt' && asset?.protocol?.mqtt">
              <!-- Client ID (full width) -->
              <div class="field-row q-mb-md">
                <div class="field-label">{{ t.drawer.fields.clientId.value }}</div>
                <div class="field-value">
                  <DetailChip
                    color="default"
                    size="sm"
                    :label="asset.protocol.mqtt.clientId"
                  />
                </div>
              </div>

              <!-- Username -->
              <div class="row q-col-gutter-sm">
                <div class="col-12">
                  <div class="field-row">
                    <div class="field-label">{{ t.drawer.fields.username.value }}</div>
                    <div class="field-value">{{ asset.protocol.mqtt.username }}</div>
                  </div>
                </div>
              </div>

            </template>
          </div>

          <!-- Auth Section — MQTT only. The asset declares ONE auth
               mode (password XOR cert) at create/update time and the
               broker enforces mutual exclusion at CONNECT; the drawer
               mirrors that contract and renders only the relevant
               section so operators never see the legacy "both rows"
               state from earlier builds. -->
          <div v-if="asset?.protocol?.type === 'mqtt'" class="section q-mb-md">
            <div class="section-header">
              <q-icon name="security" color="primary" size="sm" class="q-mr-sm" />
              <span class="text-subtitle1 text-weight-medium">{{ t.drawer.sections.auth.value }}</span>
            </div>
            <q-separator class="q-my-sm" />

            <!-- Password row — status only; mutations go through the
                 wizard (Edit asset) to keep credential entry on a
                 single audited form. Rendered only when the asset
                 declared password mode. -->
            <div v-if="mqttAuthType === 'password'" class="field-row">
              <div class="field-label">{{ t.drawer.auth.password.label.value }}</div>
              <div class="field-value">
                <DetailChip
                  :color="hasPasswordHash ? 'positive' : 'grey'"
                  size="sm"
                  :label="hasPasswordHash ? t.drawer.auth.password.set.value : t.drawer.auth.password.notSet.value"
                />
                <div class="text-caption text-grey-7 q-mt-xs">
                  {{ t.drawer.auth.password.hint.value }}
                </div>
              </div>
            </div>

            <!-- Certificate row — full lifecycle inline. Rendered only
                 when the asset declared cert mode. -->
            <div v-if="mqttAuthType === 'cert'" class="field-row">
              <div class="field-label">{{ t.drawer.auth.certificate.label.value }}</div>
              <div class="field-value">
                <template v-if="asset?.currentCert">
                  <div class="row items-center q-gutter-sm q-mb-sm">
                    <DetailChip
                      color="positive"
                      size="sm"
                      :label="t.drawer.auth.certificate.active.value"
                    />
                    <span class="text-caption text-grey-7">{{ certExpiryLabel }}</span>
                  </div>

                  <div class="cert-meta">
                    <div class="row">
                      <span class="label">{{ t.drawer.auth.certificate.fields.serial.value }}</span>
                      <code>{{ asset.currentCert.serial }}</code>
                    </div>
                    <div class="row">
                      <span class="label">{{ t.drawer.auth.certificate.fields.fingerprint.value }}</span>
                      <code class="fp">{{ asset.currentCert.fingerprint }}</code>
                    </div>
                    <div class="row">
                      <span class="label">{{ t.drawer.auth.certificate.fields.subjectCN.value }}</span>
                      <span>{{ asset.currentCert.subjectCN }}</span>
                    </div>
                    <div class="row">
                      <span class="label">{{ t.drawer.auth.certificate.fields.issued.value }}</span>
                      <span>{{ formatDate(asset.currentCert.issuedAt) }}</span>
                    </div>
                    <div class="row">
                      <span class="label">{{ t.drawer.auth.certificate.fields.expires.value }}</span>
                      <span>{{ formatDate(asset.currentCert.expiresAt) }}</span>
                    </div>
                  </div>

                  <div class="row q-gutter-sm q-mt-sm">
                    <q-btn
                      dense
                      outline
                      no-caps
                      color="primary"
                      icon="refresh"
                      :label="t.drawer.auth.certificate.actions.regenerate.value"
                      :loading="certActionLoading"
                      @click="openGenerateDialog(true)"
                    />
                    <q-btn
                      dense
                      flat
                      no-caps
                      color="negative"
                      icon="block"
                      :label="t.drawer.auth.certificate.actions.revoke.value"
                      :loading="certActionLoading"
                      @click="confirmRevokeCert"
                    />
                  </div>
                </template>

                <template v-else>
                  <DetailChip
                    color="grey"
                    size="sm"
                    :label="t.drawer.auth.certificate.noActive.value"
                  />
                  <div class="q-mt-sm">
                    <q-btn
                      unelevated
                      rounded
                      no-caps
                      color="primary"
                      icon="badge"
                      :label="t.drawer.auth.certificate.actions.generate.value"
                      :loading="certActionLoading"
                      class="q-px-md"
                      @click="openGenerateDialog(false)"
                    />
                  </div>
                </template>
              </div>
            </div>

            <!-- Revoked list — collapsible to keep the drawer compact
                 when there is no history worth surfacing. -->
            <q-expansion-item
              v-if="revokedCerts.length > 0"
              :label="t.drawer.auth.revoked.title.value + ' (' + revokedCerts.length + ')'"
              header-class="text-grey-7 text-caption q-px-none"
              dense
            >
              <div class="text-caption text-grey-7 q-mb-sm">
                {{ t.drawer.auth.revoked.retentionNotice.value }}
              </div>
              <q-list dense bordered class="rounded-borders">
                <q-item v-for="row in revokedCerts" :key="row.serial" class="q-py-sm">
                  <q-item-section>
                    <q-item-label class="text-caption">
                      <code>{{ row.serial }}</code>
                    </q-item-label>
                    <q-item-label caption>
                      {{ row.reason }} - {{ formatDate(row.revokedAt) }}
                    </q-item-label>
                  </q-item-section>
                </q-item>
              </q-list>
            </q-expansion-item>
          </div>

          <!-- Location Section -->
          <div class="section q-mb-md" v-if="hasLocation">
            <div class="section-header">
              <q-icon name="location_on" color="primary" size="sm" class="q-mr-sm" />
              <span class="text-subtitle1 text-weight-medium">{{ t.drawer.sections.location.value }}</span>
            </div>
            <q-separator class="q-my-sm" />

            <!-- Latitude & Longitude (2 columns) -->
            <div class="row q-col-gutter-sm">
              <div class="col-6">
                <div class="field-row">
                  <div class="field-label">{{ t.drawer.fields.latitude.value }}</div>
                  <div class="field-value">{{ asset?.latitude || '-' }}</div>
                </div>
              </div>
              <div class="col-6">
                <div class="field-row">
                  <div class="field-label">{{ t.drawer.fields.longitude.value }}</div>
                  <div class="field-value">{{ asset?.longitude || '-' }}</div>
                </div>
              </div>
            </div>
          </div>

          <!-- Routing Section -->
          <div class="section q-mb-md">
            <div class="section-header">
              <q-icon name="alt_route" color="primary" size="sm" class="q-mr-sm" />
              <span class="text-subtitle1 text-weight-medium">{{ t.drawer.sections.routing.value }}</span>
            </div>
            <q-separator class="q-my-sm" />

            <div class="field-row">
              <div class="field-label">{{ t.drawer.fields.routeGroups.value }}</div>
              <div class="field-value">
                <div v-if="asset?.routeGroupNames && asset.routeGroupNames.length > 0" class="row q-gutter-xs">
                  <DetailChip
                    v-for="(routeName, index) in asset.routeGroupNames"
                    :key="index"
                    icon="route"
                    color="blue"
                    size="sm"
                    :label="routeName"
                  />
                </div>
                <div v-else class="text-grey-6">
                  {{ t.drawer.empty.routeGroups.value }}
                </div>
              </div>
            </div>
          </div>

          <!-- Health Monitoring Section -->
          <div v-if="asset?.healthMonitor?.enabled" class="section q-mb-md">
            <div class="section-header">
              <q-icon name="monitor_heart" color="primary" size="sm" class="q-mr-sm" />
              <span class="text-subtitle1 text-weight-medium">{{ t.drawer.sections.healthMonitoring.value }}</span>
            </div>
            <q-separator class="q-my-sm" />

            <div class="field-row">
              <div class="field-label">{{ t.drawer.fields.healthStatus.value }}</div>
              <div class="field-value">
                <q-badge
                  :color="asset.healthStatus === 'online' ? 'green' : asset.healthStatus === 'offline' ? 'red' : 'grey-5'"
                  :label="asset.healthStatus || 'unknown'"
                />
              </div>
            </div>

            <div v-if="asset.lastSeenAt" class="field-row">
              <div class="field-label">{{ t.drawer.fields.lastSeen.value }}</div>
              <div class="field-value text-caption">{{ new Date(asset.lastSeenAt).toLocaleString() }}</div>
            </div>

            <div class="field-row">
              <div class="field-label">{{ t.drawer.fields.threshold.value }}</div>
              <div class="field-value">{{ asset.healthMonitor.thresholdMinutes }} min</div>
            </div>

            <div class="field-row">
              <div class="field-label">{{ t.drawer.fields.requiredMisses.value }}</div>
              <div class="field-value">{{ asset.healthMonitor.requiredMisses }}</div>
            </div>
          </div>

          <!-- Organization Section -->
          <div class="section q-mb-md">
            <div class="section-header">
              <q-icon name="domain" color="primary" size="sm" class="q-mr-sm" />
              <span class="text-subtitle1 text-weight-medium">{{ t.drawer.sections.organization.value }}</span>
            </div>
            <q-separator class="q-my-sm" />

            <div class="field-row">
              <div class="field-label">{{ t.drawer.fields.organization.value }}</div>
              <div class="field-value">
                <DetailChip
                  icon="domain"
                  color="indigo"
                  size="sm"
                  :label="(asset as any)?.organizationName || asset?.orgId || '-'"
                />
              </div>
            </div>

            <div class="field-row">
              <div class="field-label">{{ t.drawer.fields.pathKey.value }}</div>
              <div class="field-value text-grey-8">{{ asset?.pathKey || '-' }}</div>
            </div>

            <div class="field-row" v-if="asset?.customerId">
              <div class="field-label">{{ t.drawer.fields.customer.value }}</div>
              <div class="field-value">
                <DetailChip
                  icon="business"
                  color="blue"
                  size="sm"
                  :label="(asset as any)?.customerName || asset.customerId"
                />
              </div>
            </div>
          </div>

          <!-- Timestamps Section -->
          <div class="section">
            <div class="section-header">
              <q-icon name="schedule" color="primary" size="sm" class="q-mr-sm" />
              <span class="text-subtitle1 text-weight-medium">{{ t.drawer.sections.timestamps.value }}</span>
            </div>
            <q-separator class="q-my-sm" />

            <!-- Created & Updated (2 columns) -->
            <div class="row q-col-gutter-sm">
              <div class="col-6">
                <div class="field-row">
                  <div class="field-label">{{ t.drawer.fields.created.value }}</div>
                  <div class="field-value">{{ formatDate(asset?.created) }}</div>
                </div>
              </div>
              <div class="col-6">
                <div class="field-row">
                  <div class="field-label">{{ t.drawer.fields.updated.value }}</div>
                  <div class="field-value">{{ formatDate(asset?.updated) }}</div>
                </div>
              </div>
            </div>
          </div>

        </div>
      </q-scroll-area>
    </div>

    <!-- Footer Actions -->
    <q-separator />
    <div class="drawer-footer">
      <q-space />
      <q-btn
        unelevated
        icon="edit"
        color="primary"
        :label="t.drawer.edit.value"
        :disable="!asset"
        @click="handleEdit"
      />
    </div>
  </q-drawer>

  <!-- Cert lifecycle confirmations are surfaced through the platform's
       Quasar Dialog plugin helpers (`dialogWarning` / `dialogDelete`).
       Keeping the drawer free of inline `<q-dialog>` blocks avoids
       z-index / stacking issues against the drawer overlay and keeps
       the styling consistent with the rest of the app. -->
</template>

<script setup lang="ts">
defineOptions({
  name: 'AssetDetailsDrawer'
});

/** TYPE IMPORTS (ALL types first, grouped) */
import type { AssetDetailsDrawerProps, AssetDetailsDrawerEmits } from './interfaces/assetDetailsDrawer.interface';
import type { AssetResponse, RevokedCertResponse } from '@mapexos/schemas';

/** VUE IMPORTS */
import { ref, watch, computed, onMounted, onBeforeUnmount } from 'vue';
import { date } from 'quasar';

/** COMPONENTS */
import { DetailChip } from '@components/chips';
import { AppTooltip } from '@components/tooltips';

/** COMPOSABLES */
import { useAssetsTranslations } from '@composables/i18n';
import { useLogger } from '@composables/useLogger';

/** UTILS */
import { notifyFail, notifySuccess, dialogWarning, dialogDelete } from '@utils/alert';
import { downloadCertZip, decodeBase64ToBytes } from '@utils/zipDownload';

/** SERVICES */
import { apis } from '@services/mapex';

/** STORES */
import { useOrganizationStore } from '@stores/organization';

/** PROPS & EMITS */
const props = defineProps<AssetDetailsDrawerProps>();

const emit = defineEmits<AssetDetailsDrawerEmits>();

/** COMPOSABLES & STORES */
const t = useAssetsTranslations();
const logger = useLogger('AssetDetailsDrawer');
const orgStore = useOrganizationStore();

/** STATE */
const asset = ref<AssetResponse | null>(null);
const loading = ref(false);
const error = ref(false);

const revokedCerts = ref<RevokedCertResponse[]>([]);
const certActionLoading = ref(false);

/** COMPUTED */
const hasLocation = computed(() => {
  return asset.value?.latitude !== undefined || asset.value?.longitude !== undefined;
});

/**
 * Declared MQTT auth mode for the asset. Drives mutually-exclusive
 * rendering of the password vs. certificate row — the broker enforces
 * the same exclusion at CONNECT, so the drawer mirrors that contract
 * instead of showing both rows as the legacy build did.
 */
const mqttAuthType = computed<'password' | 'cert' | undefined>(() => {
  return (asset.value?.protocol as { mqtt?: { authType?: 'password' | 'cert' } } | undefined)
    ?.mqtt?.authType;
});

/**
 * Whether the asset has a stored bcrypt hash for password-mode CONNECTs.
 * The contract uses `passwordHash` on the read-model side; on the
 * request-side the field is `password` (plaintext, write-only). We rely
 * on the read-model field here.
 */
const hasPasswordHash = computed(() => {
  const hash = (asset.value?.protocol as { mqtt?: { passwordHash?: string } } | undefined)
    ?.mqtt?.passwordHash;
  return typeof hash === 'string' && hash.length > 0;
});

/**
 * Human-readable expiry summary for the active cert. Renders "expires
 * in N days" when in the future, "expired" once past the date — the
 * drawer is mostly observational, so a friendly relative line beats
 * a second full date string next to "Expires".
 */
const certExpiryLabel = computed(() => {
  const exp = asset.value?.currentCert?.expiresAt;
  if (!exp) return '';
  const expDate = typeof exp === 'string' ? new Date(exp) : exp;
  const diffMs = expDate.getTime() - Date.now();
  if (diffMs <= 0) return t.drawer.auth.certificate.expired.value;
  const days = Math.ceil(diffMs / (1000 * 60 * 60 * 24));
  return `${t.drawer.auth.certificate.active.value} - ${days}d`;
});

/** WATCHERS */
watch(() => props.assetId, (newAssetId) => {
  if (newAssetId && props.modelValue) {
    void fetchAssetDetails(newAssetId);
  }
}, { immediate: true });

watch(() => props.modelValue, (isOpen) => {
  if (isOpen && props.assetId) {
    void fetchAssetDetails(props.assetId);
  } else if (!isOpen) {
    asset.value = null;
    error.value = false;
    revokedCerts.value = [];
  }
});

/** FUNCTIONS */

/**
 * Fetch asset details by ID from API
 * @param {string} assetId - Asset ID
 * @returns {Promise<void>}
 */
async function fetchAssetDetails(assetId: string): Promise<void> {
  if (!apis.assets) {
    error.value = true;
    notifyFail({ message: 'Assets API not initialized' });
    return;
  }

  loading.value = true;
  error.value = false;
  asset.value = null;

  try {
    const response = await apis.assets.asset.getById({ assetId });

    const organization = orgStore.flatList.find((org: any) => org.id === response.orgId);
    const customer = response.customerId
      ? orgStore.flatList.find((org: any) => org.id === response.customerId)
      : null;

    asset.value = {
      ...response,
      organizationName: organization?.name || 'Unknown',
      customerName: customer?.name || undefined,
    } as any;

    if (response.protocol?.type === 'mqtt' && response.assetUUID) {
      void fetchRevokedCerts(response.assetUUID);
    }
  } catch (err: any) {
    logger.error('Error fetching asset details:', err);
    error.value = true;
    notifyFail({ message: t.drawer.error.value });
  } finally {
    loading.value = false;
  }
}

/**
 * Fetch revoked certificates for the asset. Failures are surfaced
 * to the logs only — the revoked list is auxiliary information, so
 * a backend hiccup must not block the drawer from rendering the rest
 * of the asset.
 */
async function fetchRevokedCerts(assetUUID: string): Promise<void> {
  try {
    revokedCerts.value = await apis.assets.mqttcerts.listRevoked(assetUUID);
  } catch (err: any) {
    logger.warn('Failed to load revoked certs', err);
    revokedCerts.value = [];
  }
}

/**
 * Open the cert generation confirmation using the platform's standard
 * warning dialog (Quasar Dialog plugin via `dialogWarning`). The same
 * modal surface is used by every destructive / cautionary action in
 * the app — keeps the visual + interaction contract consistent and
 * sidesteps the z-index issues a local `<q-dialog>` had against the
 * drawer overlay. `replaces=true` appends the extra warning about
 * revoking the active cert at issue time.
 */
async function openGenerateDialog(replaces: boolean): Promise<void> {
  if (!asset.value?.assetUUID || certActionLoading.value) return;
  const warning = t.drawer.auth.certificate.dialog.warning.value;
  const replaceWarning = replaces ? `<br><br>${t.drawer.auth.certificate.dialog.replaceWarning.value}` : '';
  const confirmed = await dialogWarning({
    title: t.drawer.auth.certificate.dialog.title.value,
    message: `${warning}${replaceWarning}`,
    html: true,
    ok: { label: t.drawer.auth.certificate.dialog.generateButton.value, color: 'primary', unelevated: true, rounded: true, noCaps: true },
    cancel: { label: t.drawer.auth.certificate.dialog.skipButton.value, flat: true, noCaps: true },
  });
  if (confirmed) {
    await issueCert(replaces);
  }
}

/**
 * Call POST /mqtt_certs (force=true when replacing the active cert so
 * the backend revokes the incumbent atomically), trigger the zip
 * download, refresh the drawer to surface the freshly-issued cert.
 */
async function issueCert(replaces: boolean): Promise<void> {
  if (!asset.value?.assetUUID || !asset.value?.id) return;
  const assetUUID = asset.value.assetUUID;
  const assetIdLocal = asset.value.id;
  const filenameBase = asset.value.name || assetUUID;
  certActionLoading.value = true;
  try {
    const res = await apis.assets.mqttcerts.issueCert(assetUUID, replaces);
    await downloadCertZip({
      filename: `${filenameBase}-mqtt-cert.zip`,
      certPEM: decodeBase64ToBytes(res.certPEM),
      keyPEM: decodeBase64ToBytes(res.keyPEM),
      caChainPEM: decodeBase64ToBytes(res.caChainPEM),
    });
    notifySuccess({ message: t.drawer.auth.certificate.issueSuccess.value });
    await fetchAssetDetails(assetIdLocal);
  } catch (err: any) {
    logger.error('Failed to issue cert', err);
    notifyFail({ message: t.drawer.auth.certificate.issueFailed.value });
  } finally {
    certActionLoading.value = false;
  }
}

/**
 * Revoke flow — confirm with the operator first via the platform's
 * delete-styled dialog (the device loses MQTT access immediately, so
 * the negative-color OK button matches the destructive intent), then
 * DELETE /mqtt_certs/:serial and refresh both the asset and the
 * revoked list.
 */
async function confirmRevokeCert(): Promise<void> {
  if (!asset.value?.id || !asset.value?.currentCert?.serial || certActionLoading.value) return;
  const confirmed = await dialogDelete({
    title: t.drawer.auth.certificate.revokeConfirmTitle.value,
    message: t.drawer.auth.certificate.revokeConfirmBody.value,
    ok: { label: t.drawer.auth.certificate.actions.revoke.value, color: 'negative', unelevated: true, noCaps: true },
    cancel: { label: t.drawer.close.value, flat: true, noCaps: true },
  });
  if (confirmed) {
    await revokeCert();
  }
}

async function revokeCert(): Promise<void> {
  if (!asset.value?.id || !asset.value?.currentCert?.serial) return;
  const assetIdLocal = asset.value.id;
  const serial = asset.value.currentCert.serial;
  certActionLoading.value = true;
  try {
    await apis.assets.mqttcerts.revokeCert(serial);
    notifySuccess({ message: t.drawer.auth.certificate.revokeSuccess.value });
    await fetchAssetDetails(assetIdLocal);
  } catch (err: any) {
    logger.error('Failed to revoke cert', err);
    notifyFail({ message: t.drawer.auth.certificate.revokeFailed.value });
  } finally {
    certActionLoading.value = false;
  }
}

/**
 * Get protocol color name (DetailChip format)
 * @param {string} protocolType - Protocol type
 * @returns {string} Color name compatible with DetailChip
 */
function getProtocolColorName(protocolType?: string): 'purple' | 'blue' | 'orange' | 'grey' {
  const colors: Record<string, 'purple' | 'blue' | 'orange' | 'grey'> = {
    'mqtt': 'purple',
    'http': 'blue',
    'lorawan': 'orange',
  };
  return colors[protocolType?.toLowerCase() || ''] || 'grey';
}

/**
 * Format date using Quasar date utils
 * @param {any} dateValue - Date value to format
 * @returns {string} Formatted date
 */
function formatDate(dateValue: any): string {
  if (!dateValue) return '-';

  try {
    const dateObj = typeof dateValue === 'string' ? new Date(dateValue) : dateValue;
    return date.formatDate(dateObj, 'MMM DD, YYYY HH:mm');
  } catch {
    return '-';
  }
}

/**
 * Close drawer
 */
function close(): void {
  emit('update:modelValue', false);
}

/**
 * Handle edit action
 */
function handleEdit(): void {
  if (!asset.value?.id) return;
  emit('edit', asset.value.id);
  close();
}

/**
 * Handle ESC key to close drawer
 * @param {KeyboardEvent} event - Keyboard event
 */
function handleEscKey(event: KeyboardEvent): void {
  if (event.key === 'Escape' && props.modelValue) {
    close();
  }
}

/** LIFECYCLE HOOKS */
onMounted(() => {
  window.addEventListener('keydown', handleEscKey);
});

onBeforeUnmount(() => {
  window.removeEventListener('keydown', handleEscKey);
});
</script>

<style lang="scss" scoped>
// Flex layout for drawer content
:deep(.q-drawer__content) {
  display: flex;
  flex-direction: column;
  height: 100%;
}

// Drawer Header
.drawer-header {
  flex-shrink: 0;
  background: var(--mapex-header-bg);
  backdrop-filter: blur(10px);
  border-bottom: 1px solid var(--mapex-header-border);

  .q-toolbar__title {
    font-size: 1.1rem;
    color: var(--q-primary);
  }
}

// Close button
.drawer-close-btn {
  color: var(--mapex-text-secondary);
}

// Backdrop (teleported to body, needs :global) - transparent, just for click detection
:global(.drawer-backdrop) {
  position: fixed;
  top: 0;
  left: 0;
  right: 450px; // Leave space for drawer (450px width)
  bottom: 0;
  background: transparent;
  z-index: 5999; // Below q-drawer (6000)
  cursor: default;
}

// Drawer Content
.drawer-content {
  flex: 1;
  min-height: 0; // Important for flex children with overflow
  overflow: hidden;

  :deep(.q-scrollarea__content) {
    width: 100%;
    max-width: 100%;
    overflow-x: hidden;
  }
}

// Drawer Footer - Fixed at bottom
.drawer-footer {
  flex-shrink: 0;
  display: flex;
  align-items: center;
  padding: 12px 16px;
  background: var(--mapex-header-bg);
  backdrop-filter: blur(10px);
  border-top: 1px solid var(--mapex-header-border);
  box-shadow: 0 -2px 8px var(--mapex-elevation-shadow);
}

// Section Styling
.section {
  .section-header {
    display: flex;
    align-items: center;
    color: var(--q-primary);
    margin-bottom: 8px;
  }
}

// Field Row Styling
.field-row {
  display: flex;
  flex-direction: column;
  padding: 10px 0;
  border-bottom: 1px solid var(--mapex-divider);

  &:last-child {
    border-bottom: none;
  }

  .field-label {
    font-size: 0.7rem;
    font-weight: 600;
    text-transform: uppercase;
    color: var(--mapex-text-secondary);
    margin-bottom: 4px;
    letter-spacing: 0.8px;
  }

  .field-value {
    font-size: 0.9rem;
    color: var(--mapex-text-primary);
    word-break: break-word;
    line-height: 1.4;
  }
}

// Post-regenerate credential dialog wrapper
.regenerate-credential-dialog {
  background: var(--mapex-surface-primary);
  padding: var(--mapex-space-md);
  border-radius: var(--mapex-radius-md);
  min-width: 480px;
  max-width: 720px;
}

// Cert metadata block under the Auth section. Compact two-column
// rows so serial/fingerprint stay readable without overflowing the
// 450px drawer width.
.cert-meta {
  background: var(--mapex-surface-secondary);
  border-radius: var(--mapex-radius-sm);
  padding: var(--mapex-space-sm);
  font-size: 0.85em;

  .row {
    display: flex;
    gap: var(--mapex-space-sm);
    padding: 2px 0;
    align-items: baseline;

    .label {
      color: var(--mapex-text-secondary);
      min-width: 90px;
      font-size: 0.85em;
      text-transform: uppercase;
      letter-spacing: 0.4px;
    }

    code {
      background: var(--mapex-surface-primary);
      padding: 2px 6px;
      border-radius: var(--mapex-radius-sm);
      word-break: break-all;
    }

    .fp {
      font-size: 0.8em;
    }
  }
}

.rounded-borders {
  border-radius: var(--mapex-radius-sm);
}

// UUID chip ellipsis — long asset UUIDs must not overflow into the Status column.
.uuid-field-value {
  min-width: 0;
  overflow: hidden;
}

:deep(.uuid-chip) {
  max-width: 100%;

  .q-chip__content {
    overflow: hidden;
  }

  .chip-label {
    overflow: hidden;
    text-overflow: ellipsis;
    min-width: 0;
  }
}

.revoke-confirm-dialog {
  background: var(--mapex-surface-primary);
  min-width: 360px;
  max-width: 480px;
  border-radius: var(--mapex-radius-md);
}

// Custom Scrollbar
:deep(.q-scrollarea__content) {
  &::-webkit-scrollbar {
    width: 6px;
  }

  &::-webkit-scrollbar-track {
    background: transparent;
    border-radius: var(--mapex-radius-lg);
  }

  &::-webkit-scrollbar-thumb {
    background: rgba(var(--q-primary-rgb), 0.3);
    border-radius: var(--mapex-radius-lg);
    transition: background var(--mapex-transition-base) ease;

    &:hover {
      background: rgba(var(--q-primary-rgb), 0.5);
    }
  }
}
</style>
