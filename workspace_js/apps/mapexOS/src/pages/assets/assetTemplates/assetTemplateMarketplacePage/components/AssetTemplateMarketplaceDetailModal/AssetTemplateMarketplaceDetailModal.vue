<script setup lang="ts">
defineOptions({
  name: 'AssetTemplateMarketplaceDetailModal'
});

/** TYPE IMPORTS */
import type { AssetTemplateBundle } from '@mapexos/schemas';
import type { DynamicField } from '@components/assetTemplates/dynamicFieldsTable';
import type {
  AssetTemplateMarketplaceDetailModalProps,
  AssetTemplateMarketplaceDetailModalEmits
} from './interfaces';

/** VUE IMPORTS */
import { ref, computed, watch } from 'vue';
import { useI18n } from 'vue-i18n';

/** COMPONENTS */
import { DynamicFieldsTable } from '@components/assetTemplates/dynamicFieldsTable';
import { AvailableFieldsList } from '@components/assetTemplates/availableFieldsList';
import { ScriptViewerDialog } from '@components/dialogs/scriptViewer';
import { DetailChip } from '@components/chips';

/** COMPOSABLES */
import { useAssetTemplateMarketplaceTranslations } from '@composables/i18n';
import { useLogger } from '@composables/useLogger';

/** UTILS */
import { notifySuccess, notifyFail } from '@utils/alert/notify';

/** SERVICES */
import { apis } from '@services/mapex';

/** PROPS & EMITS */
const props = defineProps<AssetTemplateMarketplaceDetailModalProps>();
const emit = defineEmits<AssetTemplateMarketplaceDetailModalEmits>();

/** COMPOSABLES & STORES */
const t = useAssetTemplateMarketplaceTranslations();
const { locale } = useI18n();
const logger = useLogger('AssetTemplateMarketplaceDetailModal');

/** STATE */
const bundle = ref<AssetTemplateBundle | null>(null);
const loading = ref(false);
const loadError = ref(false);
const installing = ref(false);
const installError = ref<string | undefined>(undefined);
const showScript = ref(false);
const scriptTitle = ref('');
const scriptContent = ref('');

/** COMPUTED */
const isOpen = computed<boolean>({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value),
});

/**
 * The template name resolved for the active locale, falling back to en-US and
 * then to any available translation, since the bundle carries a locale map.
 */
const templateName = computed<string>(() => resolveLocalized(bundle.value?.name));

/**
 * The template description resolved the same way as the name.
 */
const templateDescription = computed<string>(() => resolveLocalized(bundle.value?.description));

/**
 * The bundle dynamic fields typed as the shared table's row shape. The bundle
 * omits the optional geo path fields, which the table treats as absent.
 */
const dynamicFields = computed<DynamicField[]>(
  () => bundle.value?.dynamicFields ?? []
);

/**
 * The four template scripts with a localized label and a configured flag,
 * driving both the status indicator and the view-code action.
 */
const scripts = computed(() => {
  const b = bundle.value;
  return [
    { key: 'test', label: t.modal.scripts.test.value, content: b?.scriptTest ?? '' },
    { key: 'processor', label: t.modal.scripts.processor.value, content: b?.scriptProcessor ?? '' },
    { key: 'validator', label: t.modal.scripts.validator.value, content: b?.scriptValidator ?? '' },
    { key: 'conversion', label: t.modal.scripts.conversion.value, content: b?.scriptConversion ?? '' },
  ].map((s) => ({ ...s, configured: s.content.trim().length > 0 }));
});

/** WATCHERS */
watch(
  () => props.modelValue,
  (open) => {
    if (open && props.vendor && props.slug) {
      void fetchBundle(props.vendor, props.slug);
    }
    if (!open) {
      resetState();
    }
  }
);

/** FUNCTIONS */

/**
 * Resolve a localized `{locale: value}` map to a single string for the active
 * locale, falling back to en-US and then the first available value.
 * @param {Record<string, string> | undefined} map - The locale map to resolve.
 * @returns {string} The best-matching localized string, or an empty string.
 */
function resolveLocalized(map: Record<string, string> | undefined): string {
  if (!map) return '';
  return map[locale.value] ?? map['en-US'] ?? Object.values(map)[0] ?? '';
}

/**
 * Fetch the template bundle for the given vendor/slug and populate the modal.
 * @param {string} vendor - Template vendor key.
 * @param {string} slug - Template slug.
 */
async function fetchBundle(vendor: string, slug: string): Promise<void> {
  if (!apis.assetTemplatesMarketplace) return;

  loading.value = true;
  loadError.value = false;
  installError.value = undefined;
  bundle.value = null;

  try {
    const res = await apis.assetTemplatesMarketplace.marketplace.get({ vendor, slug });
    bundle.value = res.data;
  } catch (err) {
    logger.error('Failed to load the asset template bundle', err);
    loadError.value = true;
  } finally {
    loading.value = false;
  }
}

/**
 * Reset the transient modal state when it closes so a re-open starts clean.
 */
function resetState(): void {
  bundle.value = null;
  loadError.value = false;
  installError.value = undefined;
  showScript.value = false;
}

/**
 * Open the Monaco script viewer for a given script.
 * @param {string} title - The localized script label used as the viewer title.
 * @param {string} content - The script source to display.
 */
function openScript(title: string, content: string): void {
  scriptTitle.value = title;
  scriptContent.value = content;
  showScript.value = true;
}

/**
 * Extract a user-facing install error message, treating a checksum mismatch as
 * a distinct, explanatory case and everything else as a generic failure.
 * @param {unknown} err - The rejected install error.
 * @returns {string} A localized, user-facing error message.
 */
function resolveInstallError(err: unknown): string {
  const raw = (err as { response?: { data?: { errors?: unknown } } })?.response?.data?.errors;
  const parts: string[] = Array.isArray(raw)
    ? raw
        .map((e) =>
          typeof e === 'string'
            ? e
            : ((e as { message?: string; code?: string })?.message ??
               (e as { code?: string })?.code ??
               '')
        )
        .filter((s): s is string => Boolean(s))
    : [];
  const joined = parts.join(', ');

  if (/checksum/i.test(joined)) {
    const human = parts.filter((p) => !/^CHECKSUM_MISMATCH$/i.test(p)).join(', ');
    return human || t.install.checksumError.value;
  }

  return t.install.genericError.value;
}

/**
 * Install the previewed template into the current organization. Verifies via
 * the backend (which hard-checks the artifact sha256); on a checksum mismatch
 * the install is rejected and the reason is surfaced in the modal and a toast.
 */
async function handleInstall(): Promise<void> {
  if (!props.vendor || !props.slug || installing.value) return;

  installing.value = true;
  installError.value = undefined;
  emit('install', { vendor: props.vendor, slug: props.slug });

  try {
    await apis.assets.assetTemplate.install({ vendor: props.vendor, slug: props.slug });
    notifySuccess({ message: t.install.success.value });
    emit('installed', { vendor: props.vendor, slug: props.slug });
    isOpen.value = false;
  } catch (err) {
    const message = resolveInstallError(err);
    installError.value = message;
    notifyFail({ message });
    logger.error('Failed to install asset template from marketplace', err);
  } finally {
    installing.value = false;
  }
}

/**
 * Close the modal without installing.
 */
function close(): void {
  isOpen.value = false;
}
</script>

<template>
  <q-dialog v-model="isOpen">
    <q-card class="marketplace-detail">
      <!-- Header -->
      <q-card-section class="marketplace-detail__header">
        <div class="marketplace-detail__header-main">
          <q-icon name="widgets" size="md" color="primary" />
          <div class="marketplace-detail__titles">
            <span class="marketplace-detail__title">
              {{ templateName || t.modal.title.value }}
            </span>
            <div v-if="bundle" class="marketplace-detail__chips">
              <DetailChip
                v-if="bundle.manufacturerName"
                color="blue"
                size="sm"
                :label="bundle.manufacturerName"
              />
              <DetailChip
                v-if="bundle.modelName"
                color="grey"
                size="sm"
                :label="bundle.modelName"
              />
              <DetailChip
                v-if="bundle.version"
                color="green"
                size="sm"
                :label="bundle.version"
              />
            </div>
          </div>
        </div>
        <q-btn v-close-popup flat round dense color="grey-7" icon="close" />
      </q-card-section>

      <q-separator />

      <!-- Body -->
      <q-card-section class="marketplace-detail__body">
        <!-- Loading -->
        <div v-if="loading" class="marketplace-detail__status">
          <q-spinner color="primary" size="36px" />
          <span class="marketplace-detail__status-text">{{ t.modal.loading.value }}</span>
        </div>

        <!-- Load error -->
        <div v-else-if="loadError" class="marketplace-detail__status">
          <q-icon name="error_outline" size="36px" color="negative" />
          <span class="marketplace-detail__status-text">{{ t.modal.loadError.value }}</span>
        </div>

        <!-- Content -->
        <template v-else-if="bundle">
          <!-- Overview -->
          <section class="marketplace-detail__section">
            <h2 class="marketplace-detail__section-title">{{ t.modal.sections.overview.value }}</h2>
            <p v-if="templateDescription" class="marketplace-detail__description">
              {{ templateDescription }}
            </p>
            <dl class="marketplace-detail__grid">
              <template v-if="bundle.categoryName">
                <dt>{{ t.modal.overview.category.value }}</dt>
                <dd>{{ bundle.categoryName }}</dd>
              </template>
              <template v-if="bundle.manufacturerName">
                <dt>{{ t.modal.overview.manufacturer.value }}</dt>
                <dd>{{ bundle.manufacturerName }}</dd>
              </template>
              <template v-if="bundle.modelName">
                <dt>{{ t.modal.overview.model.value }}</dt>
                <dd>{{ bundle.modelName }}</dd>
              </template>
              <template v-if="bundle.version">
                <dt>{{ t.modal.overview.version.value }}</dt>
                <dd>{{ bundle.version }}</dd>
              </template>
              <template v-if="bundle.assetIdPath">
                <dt>{{ t.modal.overview.assetIdPath.value }}</dt>
                <dd><code>{{ bundle.assetIdPath }}</code></dd>
              </template>
            </dl>
          </section>

          <!-- Dynamic Fields -->
          <section v-if="dynamicFields.length > 0" class="marketplace-detail__section">
            <h2 class="marketplace-detail__section-title">{{ t.modal.sections.dynamicFields.value }}</h2>
            <DynamicFieldsTable :fields="dynamicFields" />
          </section>

          <!-- Scripts -->
          <section class="marketplace-detail__section">
            <h2 class="marketplace-detail__section-title">{{ t.modal.sections.scripts.value }}</h2>
            <div class="marketplace-detail__scripts">
              <div
                v-for="script in scripts"
                :key="script.key"
                class="marketplace-detail__script"
              >
                <span class="marketplace-detail__script-label">{{ script.label }}</span>
                <DetailChip
                  :color="script.configured ? 'green' : 'grey'"
                  size="sm"
                  :label="script.configured ? t.modal.scripts.configured.value : t.modal.scripts.notConfigured.value"
                />
                <q-btn
                  v-if="script.configured"
                  flat
                  dense
                  size="sm"
                  color="primary"
                  icon="code"
                  :label="t.modal.scripts.view.value"
                  @click="openScript(script.label, script.content)"
                />
              </div>
            </div>
          </section>

          <!-- Available Fields -->
          <section v-if="bundle.availableFields.length > 0" class="marketplace-detail__section">
            <h2 class="marketplace-detail__section-title">{{ t.modal.sections.availableFields.value }}</h2>
            <AvailableFieldsList :fields="bundle.availableFields" />
          </section>
        </template>
      </q-card-section>

      <!-- Install error banner -->
      <q-banner v-if="installError" dense class="marketplace-detail__error">
        <template #avatar>
          <q-icon name="gpp_bad" color="negative" />
        </template>
        {{ installError }}
      </q-banner>

      <q-separator />

      <!-- Footer -->
      <q-card-actions align="right" class="marketplace-detail__footer">
        <q-btn
          flat
          color="grey-7"
          :label="t.modal.actions.cancel.value"
          :disable="installing"
          @click="close"
        />
        <q-btn
          unelevated
          color="primary"
          icon="download"
          :label="installing ? t.install.installing.value : t.install.button.value"
          :loading="installing"
          :disable="!bundle || loading"
          @click="handleInstall"
        />
      </q-card-actions>
    </q-card>

    <!-- Monaco script viewer -->
    <ScriptViewerDialog
      v-model="showScript"
      :title="scriptTitle"
      :script-content="scriptContent"
      language="javascript"
    />
  </q-dialog>
</template>

<style lang="scss" scoped>
.marketplace-detail {
  width: 100%;
  max-width: 860px;

  &__header {
    display: flex;
    align-items: flex-start;
    justify-content: space-between;
    gap: var(--mapex-spacing-md);
    padding: var(--mapex-spacing-lg);
    background: var(--mapex-surface-sunken);
  }

  &__header-main {
    display: flex;
    align-items: flex-start;
    gap: var(--mapex-spacing-md);
  }

  &__titles {
    display: flex;
    flex-direction: column;
    gap: var(--mapex-spacing-xs);
  }

  &__title {
    font-size: 1.15rem;
    font-weight: 600;
    color: var(--mapex-text-primary);
  }

  &__chips {
    display: flex;
    flex-wrap: wrap;
    gap: var(--mapex-spacing-xs);
  }

  &__body {
    display: flex;
    flex-direction: column;
    gap: var(--mapex-spacing-lg);
    max-height: 60vh;
    overflow-y: auto;
    padding: var(--mapex-spacing-lg);
  }

  &__status {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: var(--mapex-spacing-md);
    padding: var(--mapex-spacing-xl) 0;
  }

  &__status-text {
    color: var(--mapex-text-secondary);
  }

  &__section {
    display: flex;
    flex-direction: column;
    gap: var(--mapex-spacing-sm);
  }

  &__section-title {
    margin: 0;
    font-size: 0.75rem;
    font-weight: 600;
    text-transform: uppercase;
    letter-spacing: 0.8px;
    color: var(--mapex-text-secondary);
  }

  &__description {
    margin: 0;
    color: var(--mapex-text-primary);
  }

  &__grid {
    display: grid;
    grid-template-columns: max-content 1fr;
    gap: var(--mapex-spacing-xs) var(--mapex-spacing-md);
    margin: 0;

    dt {
      color: var(--mapex-text-secondary);
      font-size: 0.85rem;
    }

    dd {
      margin: 0;
      color: var(--mapex-text-primary);
      font-size: 0.85rem;
    }
  }

  &__scripts {
    display: flex;
    flex-direction: column;
    gap: var(--mapex-spacing-xs);
  }

  &__script {
    display: flex;
    align-items: center;
    gap: var(--mapex-spacing-md);
  }

  &__script-label {
    min-width: 120px;
    color: var(--mapex-text-primary);
    font-size: 0.9rem;
  }

  &__error {
    color: var(--mapex-text-primary);
    background: var(--mapex-surface-sunken);
  }

  &__footer {
    padding: var(--mapex-spacing-md) var(--mapex-spacing-lg);
  }
}
</style>
