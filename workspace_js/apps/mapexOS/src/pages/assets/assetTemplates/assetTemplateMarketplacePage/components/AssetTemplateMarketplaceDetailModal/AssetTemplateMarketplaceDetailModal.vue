<script setup lang="ts">
defineOptions({
  name: 'AssetTemplateMarketplaceDetailModal'
});

/** TYPE IMPORTS */
import type { AssetTemplateBundle } from '@mapexos/schemas';
import type { DynamicField } from '@components/assetTemplates/dynamicFieldsTable';
import type { AppTabItem } from '@components/tabs';
import type {
  AssetTemplateMarketplaceDetailModalProps,
  AssetTemplateMarketplaceDetailModalEmits
} from './interfaces';

/** VUE IMPORTS */
import { ref, computed, watch } from 'vue';
import { useI18n } from 'vue-i18n';

/** COMPONENTS */
import { AppTabs } from '@components/tabs';
import { InfoBanner } from '@components/banners';
import { DynamicFieldsTable } from '@components/assetTemplates/dynamicFieldsTable';
import { ScriptViewerDialog } from '@components/dialogs/scriptViewer';
import { AppTooltip } from '@components/tooltips';

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
const shareWithChildren = ref(false);
const showScript = ref(false);
const scriptTitle = ref('');
const scriptContent = ref('');
const activeTab = ref('setup');

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
 * Whether the template exposes any queryable (dynamic) fields.
 */
const hasRetrieval = computed(() => dynamicFields.value.length > 0);

/**
 * Tab strip mirroring the add-template wizard's context grouping: Setup (identity),
 * Uplink (the payload scripts) and Retrieval (the queryable fields). Retrieval
 * shows only when the template carries fields.
 */
const tabs = computed<AppTabItem[]>(() => {
  const list: AppTabItem[] = [
    { id: 'tab-setup', name: 'setup', label: t.modal.sections.setup.value, icon: 'mdi-cog-outline' },
    { id: 'tab-uplink', name: 'uplink', label: t.modal.sections.uplink.value, icon: 'mdi-upload-network-outline' },
  ];
  if (hasRetrieval.value) {
    list.push({
      id: 'tab-retrieval',
      name: 'retrieval',
      label: t.modal.sections.retrieval.value,
      icon: 'mdi-database-search-outline',
      badge: dynamicFields.value.length || undefined,
    });
  }
  return list;
});

/**
 * The specification tiles shown on the General tab, filtered to the values the
 * bundle actually carries.
 */
const specs = computed(() => {
  const b = bundle.value;
  if (!b) return [];
  // The model name doubles as the modal title, so the tile shows only the short
  // model token (the part before " - ") to avoid repeating the whole sentence.
  const modelToken = b.modelName?.split(' - ')[0]?.trim() ?? '';
  return [
    { icon: 'category', label: t.modal.overview.category.value, value: b.categoryName },
    { icon: 'factory', label: t.modal.overview.manufacturer.value, value: b.manufacturerName },
    { icon: 'developer_board', label: t.modal.overview.model.value, value: modelToken },
    { icon: 'sell', label: t.modal.overview.version.value, value: b.version },
  ].filter((s) => Boolean(s.value));
});

/**
 * The four pipeline scripts, each with an icon, a one-line role description and
 * a configured flag, driving the script cards and the view-code action.
 */
const scripts = computed(() => {
  const b = bundle.value;
  return [
    { key: 'test', icon: 'science', label: t.modal.scripts.test.value, description: t.modal.scripts.descriptions.test.value, content: b?.scriptTest ?? '' },
    { key: 'processor', icon: 'tune', label: t.modal.scripts.processor.value, description: t.modal.scripts.descriptions.processor.value, content: b?.scriptProcessor ?? '' },
    { key: 'validator', icon: 'verified', label: t.modal.scripts.validator.value, description: t.modal.scripts.descriptions.validator.value, content: b?.scriptValidator ?? '' },
    { key: 'conversion', icon: 'sync_alt', label: t.modal.scripts.conversion.value, description: t.modal.scripts.descriptions.conversion.value, content: b?.scriptConversion ?? '' },
  ].map((s) => ({ ...s, configured: s.content.trim().length > 0 }));
});

/** WATCHERS */
watch(
  () => props.modelValue,
  (open) => {
    if (open && props.vendor && props.slug) {
      activeTab.value = 'setup';
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
    // The bundle endpoint serves the raw bundle (not a {status,errors,data}
    // envelope), so the client returns it directly.
    bundle.value = await apis.assetTemplatesMarketplace.marketplace.get({ vendor, slug });
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
    await apis.assets.assetTemplate.install(
      { vendor: props.vendor, slug: props.slug },
      { shareWithChildren: shareWithChildren.value },
    );
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

        <!-- Content, organized into tabs -->
        <template v-else-if="bundle">
          <AppTabs v-model="activeTab" :tabs="tabs" />

          <q-tab-panels v-model="activeTab" animated class="marketplace-detail__panels">
            <!-- Setup: identity + specifications + asset id path -->
            <q-tab-panel name="setup" class="marketplace-detail__panel">
              <header class="panel-head">
                <q-icon name="mdi-cog-outline" size="20px" class="panel-head__icon" />
                <div class="panel-head__text">
                  <h3 class="panel-head__title">{{ t.modal.sections.setup.value }}</h3>
                  <p class="panel-head__subtitle">{{ t.modal.subtitles.setup.value }}</p>
                </div>
              </header>

              <div v-if="templateDescription" class="marketplace-detail__block">
                <span class="marketplace-detail__block-label">
                  {{ t.modal.overview.description.value }}
                </span>
                <p class="marketplace-detail__description">{{ templateDescription }}</p>
              </div>

              <div class="marketplace-detail__block">
                <span class="marketplace-detail__block-label">
                  {{ t.modal.overview.specifications.value }}
                </span>
                <div class="marketplace-detail__specs">
                  <div
                    v-for="spec in specs"
                    :key="spec.label"
                    class="spec-tile"
                  >
                    <q-icon :name="spec.icon" size="18px" class="spec-tile__icon" />
                    <div class="spec-tile__text">
                      <span class="spec-tile__label">{{ spec.label }}</span>
                      <span class="spec-tile__value">{{ spec.value }}</span>
                    </div>
                  </div>
                </div>
              </div>

              <div v-if="bundle.assetIdPath" class="marketplace-detail__block">
                <span class="marketplace-detail__block-label">
                  {{ t.modal.overview.assetIdPath.value }}
                </span>
                <code class="marketplace-detail__code">{{ bundle.assetIdPath }}</code>
              </div>
            </q-tab-panel>

            <!-- Uplink: the payload processing pipeline -->
            <q-tab-panel name="uplink" class="marketplace-detail__panel">
              <header class="panel-head">
                <q-icon name="mdi-upload-network-outline" size="20px" class="panel-head__icon" />
                <div class="panel-head__text">
                  <h3 class="panel-head__title">{{ t.modal.sections.uplink.value }}</h3>
                  <p class="panel-head__subtitle">{{ t.modal.subtitles.uplink.value }}</p>
                </div>
              </header>
              <div class="marketplace-detail__scripts">
                <div
                  v-for="script in scripts"
                  :key="script.key"
                  class="script-card"
                  :class="{ 'script-card--off': !script.configured }"
                >
                  <div class="script-card__head">
                    <q-icon :name="script.icon" size="20px" class="script-card__icon" />
                    <span class="script-card__title">{{ script.label }}</span>
                    <q-icon
                      :name="script.configured ? 'check_circle' : 'remove_circle_outline'"
                      size="16px"
                      :color="script.configured ? 'positive' : 'grey-5'"
                      class="script-card__status"
                    >
                      <AppTooltip
                        :text="script.configured ? t.modal.scripts.configured.value : t.modal.scripts.notConfigured.value"
                      />
                    </q-icon>
                  </div>
                  <span class="script-card__desc">{{ script.description }}</span>
                  <q-btn
                    v-if="script.configured"
                    flat
                    dense
                    no-caps
                    size="sm"
                    color="primary"
                    icon="code"
                    :label="t.modal.scripts.view.value"
                    class="script-card__view"
                    @click="openScript(script.label, script.content)"
                  />
                </div>
              </div>
            </q-tab-panel>

            <!-- Retrieval: the queryable (dynamic) fields -->
            <q-tab-panel v-if="hasRetrieval" name="retrieval" class="marketplace-detail__panel">
              <header class="panel-head">
                <q-icon name="mdi-database-search-outline" size="20px" class="panel-head__icon" />
                <div class="panel-head__text">
                  <h3 class="panel-head__title">{{ t.modal.sections.retrieval.value }}</h3>
                  <p class="panel-head__subtitle">{{ t.modal.subtitles.retrieval.value }}</p>
                </div>
              </header>
              <DynamicFieldsTable :fields="dynamicFields" />
            </q-tab-panel>
          </q-tab-panels>
        </template>
      </q-card-section>

      <!-- Install error banner -->
      <InfoBanner v-if="installError" variant="danger" icon="gpp_bad" dense class="q-mx-md q-mb-sm">
        {{ installError }}
      </InfoBanner>

      <q-separator />

      <!-- Footer -->
      <q-card-actions align="right" class="marketplace-detail__footer">
        <q-toggle
          v-model="shareWithChildren"
          :label="t.install.shareWithChildren.value"
          :disable="installing"
          color="primary"
          class="q-mr-auto"
        >
          <AppTooltip :text="t.install.shareWithChildrenHint.value" />
        </q-toggle>
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

  // A single fixed height for the whole body so switching tabs never resizes
  // the modal; the active panel scrolls internally instead.
  &__body {
    display: flex;
    flex-direction: column;
    height: 60vh;
    padding: 0;
    overflow: hidden;
  }

  &__status {
    display: flex;
    flex: 1;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: var(--mapex-spacing-md);
  }

  &__status-text {
    color: var(--mapex-text-secondary);
  }

  // The panels fill the remaining body height; each panel scrolls on its own.
  &__panels {
    flex: 1;
    min-height: 0;

    :deep(.q-panel-parent),
    :deep(.q-panel) {
      height: 100%;
    }
  }

  &__panel {
    height: 100%;
    overflow-y: auto;
    display: flex;
    flex-direction: column;
    gap: var(--mapex-spacing-xl);
    padding: var(--mapex-spacing-lg);
  }

  &__description {
    margin: 0;
    color: var(--mapex-text-secondary);
    line-height: 1.55;
  }

  &__block {
    display: flex;
    flex-direction: column;
    gap: var(--mapex-spacing-sm);
  }

  &__block-label {
    font-size: 0.7rem;
    font-weight: 600;
    text-transform: uppercase;
    letter-spacing: 0.7px;
    color: var(--mapex-text-secondary);
  }

  &__specs {
    display: grid;
    grid-template-columns: repeat(2, minmax(0, 1fr));
    gap: var(--mapex-spacing-sm);
  }

  &__code {
    align-self: flex-start;
    max-width: 100%;
    overflow-x: auto;
    font-family: 'Courier New', monospace;
    font-size: 0.8rem;
    color: var(--mapex-text-primary);
    background: var(--mapex-surface-sunken);
    border: 1px solid var(--mapex-divider);
    border-radius: var(--mapex-radius-sm, 6px);
    padding: 6px 10px;
  }

  &__scripts {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
    gap: var(--mapex-spacing-md);
  }

  &__footer {
    padding: var(--mapex-spacing-md) var(--mapex-spacing-lg);
  }
}

/* Standard header shown at the top of every tab panel: title + explanation. */
.panel-head {
  display: flex;
  align-items: flex-start;
  gap: var(--mapex-spacing-sm);
  padding-bottom: var(--mapex-spacing-md);
  border-bottom: 1px solid var(--mapex-divider);

  &__icon {
    color: var(--mapex-primary);
    margin-top: 2px;
    flex-shrink: 0;
  }

  &__title {
    margin: 0;
    font-size: 1rem;
    font-weight: 600;
    color: var(--mapex-text-primary);
    text-transform: capitalize;
    line-height: 1.2;
  }

  &__subtitle {
    margin: 3px 0 0;
    font-size: 0.8rem;
    color: var(--mapex-text-secondary);
    line-height: 1.4;
  }
}

/* A labelled spec value on the General tab. */
.spec-tile {
  display: flex;
  align-items: center;
  gap: var(--mapex-spacing-sm);
  padding: var(--mapex-spacing-sm) var(--mapex-spacing-md);
  background: var(--mapex-surface-sunken);
  border: 1px solid var(--mapex-divider);
  border-radius: var(--mapex-radius-md);

  &__icon {
    color: var(--mapex-primary);
    flex-shrink: 0;
  }

  &__text {
    display: flex;
    flex-direction: column;
    min-width: 0;
  }

  &__label {
    font-size: 0.65rem;
    font-weight: 600;
    text-transform: uppercase;
    letter-spacing: 0.5px;
    color: var(--mapex-text-secondary);
  }

  &__value {
    font-size: 0.9rem;
    font-weight: 500;
    color: var(--mapex-text-primary);
    word-break: break-word;
  }
}

/* One step of the processing pipeline on the Scripts tab. */
.script-card {
  display: flex;
  flex-direction: column;
  gap: var(--mapex-spacing-xs);
  padding: var(--mapex-spacing-md);
  background: var(--mapex-surface-elevated);
  border: 1px solid var(--mapex-divider);
  border-radius: var(--mapex-radius-md);
  transition: border-color var(--mapex-transition-base), box-shadow var(--mapex-transition-base);

  &:hover {
    border-color: var(--mapex-primary);
    box-shadow: var(--mapex-shadow-sm, 0 2px 8px rgba(0, 0, 0, 0.08));
  }

  &--off {
    opacity: 0.7;

    &:hover {
      border-color: var(--mapex-divider);
      box-shadow: none;
    }
  }

  &__head {
    display: flex;
    align-items: center;
    gap: var(--mapex-spacing-sm);
  }

  &__icon {
    color: var(--mapex-primary);
  }

  &__title {
    font-size: 0.95rem;
    font-weight: 600;
    color: var(--mapex-text-primary);
    text-transform: capitalize;
  }

  &__status {
    margin-left: auto;
  }

  &__desc {
    font-size: 0.8rem;
    color: var(--mapex-text-secondary);
    line-height: 1.4;
  }

  &__view {
    align-self: flex-start;
    margin-top: var(--mapex-spacing-xs);
  }
}
</style>
