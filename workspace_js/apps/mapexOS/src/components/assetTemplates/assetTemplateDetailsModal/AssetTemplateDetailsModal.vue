<template>
  <q-dialog v-model="isOpen">
    <q-card class="template-details-modal">
      <!-- Header -->
      <q-card-section class="row items-center template-details-modal__header">
        <q-icon name="memory" size="sm" color="primary" class="q-mr-sm" />
        <div class="text-h6 text-weight-medium ellipsis">
          {{ template?.name || t.drawer.title.value }}
        </div>
        <q-space />
        <q-btn v-close-popup flat round dense icon="close" color="grey-7">
          <AppTooltip :content="t.drawer.close.value" />
        </q-btn>
      </q-card-section>

      <!-- System-template hint -->
      <div v-if="template?.isSystem" class="q-px-md q-pt-md">
        <InfoBanner variant="warning" dense>{{ t.drawer.systemTemplateWarning.value }}</InfoBanner>
      </div>

      <!-- Tab strip (data-driven) -->
      <q-tabs
        v-model="activeTab"
        align="left"
        dense
        active-color="primary"
        indicator-color="primary"
        class="template-details-modal__tabs"
      >
        <q-tab
          v-for="tab in visibleTabs"
          :key="tab.key"
          :name="tab.key"
          :icon="tab.icon"
          :label="tab.label"
          no-caps
        />
      </q-tabs>
      <q-separator />

      <!-- Tab panels. Adding a tab: append an entry to TABS + a matching panel
           here. Next intended additions: 'downlinks' and 'digitalTwin'. -->
      <q-tab-panels v-model="activeTab" animated class="template-details-modal__body">
        <!-- Overview -->
        <q-tab-panel name="overview" class="q-gutter-y-md">
          <div class="detail-grid">
            <div class="detail-field">
              <div class="detail-field__label">{{ t.drawer.fields.status.value }}</div>
              <DetailChip
                :color="template?.enabled ? 'positive' : 'negative'"
                size="sm"
                :label="(template?.enabled ? t.status.active.value : t.status.inactive.value).toUpperCase()"
              />
            </div>
            <div class="detail-field">
              <div class="detail-field__label">{{ t.drawer.fields.isSystem.value }}</div>
              <DetailChip
                :icon="template?.isSystem ? 'lock' : 'lock_open'"
                :color="template?.isSystem ? 'orange' : 'grey'"
                size="sm"
                :label="(template?.isSystem ? t.drawer.system.yes.value : t.drawer.system.no.value).toUpperCase()"
              />
            </div>
            <div class="detail-field">
              <div class="detail-field__label">{{ t.drawer.fields.manufacturer.value }}</div>
              <DetailChip icon="factory" color="blue" size="sm" :label="template?.manufacturerName || '-'" />
            </div>
            <div class="detail-field">
              <div class="detail-field__label">{{ t.drawer.fields.model.value }}</div>
              <DetailChip icon="router" color="indigo" size="sm" :label="template?.modelName || '-'" />
            </div>
            <div class="detail-field">
              <div class="detail-field__label">{{ t.drawer.fields.version.value }}</div>
              <DetailChip icon="label" color="purple" size="sm" :label="template?.version || '-'" />
            </div>
            <div class="detail-field">
              <div class="detail-field__label">{{ t.drawer.fields.created.value }}</div>
              <div class="detail-field__value">{{ formatDate(template?.created) }}</div>
            </div>
          </div>

          <div class="detail-field">
            <div class="detail-field__label">{{ t.drawer.fields.description.value }}</div>
            <div class="detail-field__value text-grey-8">
              {{ template?.description || t.drawer.empty.description.value }}
            </div>
          </div>

          <div class="detail-field">
            <div class="detail-field__label">{{ t.drawer.fields.assetIdPath.value }}</div>
            <code class="detail-field__code">{{ template?.assetIdPath || '-' }}</code>
          </div>
        </q-tab-panel>

        <!-- Dynamic Fields -->
        <q-tab-panel name="dynamicFields">
          <DynamicFieldsTable :fields="template?.dynamicFields ?? []" />
        </q-tab-panel>

        <!-- Scripts -->
        <q-tab-panel name="scripts" class="q-gutter-y-sm">
          <div v-for="s in scriptRows" :key="s.key" class="detail-field detail-field--row">
            <div class="detail-field__label q-mb-none">{{ s.label }}</div>
            <div class="row items-center q-gutter-xs">
              <DetailChip
                :icon="hasScript(s.key) ? 'check_circle' : 'cancel'"
                :color="hasScript(s.key) ? 'green' : 'grey'"
                size="sm"
                :label="hasScript(s.key) ? t.drawer.scripts.configured.value : t.drawer.scripts.notConfigured.value"
              />
              <q-btn
                v-if="hasScript(s.key)"
                flat dense round size="sm"
                icon="visibility"
                color="primary"
                @click="viewScript(s.key, s.label)"
              >
                <AppTooltip :content="t.drawer.scriptViewer.viewScript.value" />
              </q-btn>
            </div>
          </div>
        </q-tab-panel>

        <!-- Available Fields -->
        <q-tab-panel name="availableFields">
          <AvailableFieldsList :fields="template?.availableFields ?? []" />
        </q-tab-panel>
      </q-tab-panels>

      <!-- Footer: caller supplies the actions (Edit/Duplicate, Select, ...) -->
      <q-separator />
      <q-card-actions align="right" class="template-details-modal__footer">
        <slot name="actions" :template="template" />
      </q-card-actions>
    </q-card>

    <!-- Script viewer -->
    <ScriptViewerDialog
      v-model="showScriptViewer"
      :title="currentScriptTitle"
      :script-content="currentScriptContent"
      language="javascript"
      :copy-tooltip="t.drawer.scriptViewer.copyScript.value"
      :close-tooltip="t.drawer.scriptViewer.close.value"
      :copy-success-message="t.drawer.scriptViewer.copySuccess.value"
      :copy-fail-message="t.drawer.scriptViewer.copyFail.value"
    />
  </q-dialog>
</template>

<script setup lang="ts">
defineOptions({
  name: 'AssetTemplateDetailsModal'
});

/** TYPE IMPORTS */
import type { AssetTemplateDetailsModalProps, AssetTemplateDetailsModalEmits, AssetTemplateDetailTab } from './interfaces/assetTemplateDetailsModal.interface';

/** VUE */
import { ref, computed, watch } from 'vue';
import { date } from 'quasar';

/** COMPONENTS */
import { DynamicFieldsTable, AvailableFieldsList } from '@components/assetTemplates';
import { ScriptViewerDialog } from '@components/dialogs/scriptViewer';
import { DetailChip } from '@components/chips';
import { AppTooltip } from '@components/tooltips';
import { InfoBanner } from '@components/banners';

/** COMPOSABLES */
import { useAssetTemplatesTranslations } from '@composables/i18n';

/** PROPS & EMITS */
const props = defineProps<AssetTemplateDetailsModalProps>();
const emit = defineEmits<AssetTemplateDetailsModalEmits>();

/** COMPOSABLES & STORES */
const t = useAssetTemplatesTranslations();

/** STATE */
const activeTab = ref('overview');
const showScriptViewer = ref(false);
const currentScriptTitle = ref('');
const currentScriptContent = ref('');

// The four transform scripts, in pipeline order.
const scriptRows = computed(() => [
  { key: 'scriptTest', label: t.drawer.fields.scriptTest.value },
  { key: 'scriptProcessor', label: t.drawer.fields.scriptProcessor.value },
  { key: 'scriptValidator', label: t.drawer.fields.scriptValidator.value },
  { key: 'scriptConversion', label: t.drawer.fields.scriptConversion.value },
]);

/** COMPUTED */
const isOpen = computed({
  get: () => props.modelValue,
  set: (v) => emit('update:modelValue', v),
});

const dynamicFieldsCount = computed(() => props.template?.dynamicFields?.length ?? 0);
const availableFieldsCount = computed(() => props.template?.availableFields?.length ?? 0);

// The tab catalog. Data tabs hide when their array is empty.
const visibleTabs = computed<AssetTemplateDetailTab[]>(() => {
  const tabs: AssetTemplateDetailTab[] = [
    { key: 'overview', icon: 'info', label: t.drawer.tabs.overview.value, always: true },
    { key: 'dynamicFields', icon: 'dataset', label: t.drawer.tabs.dynamicFields.value, showWhen: () => dynamicFieldsCount.value > 0 },
    { key: 'scripts', icon: 'code', label: t.drawer.tabs.scripts.value, always: true },
    { key: 'availableFields', icon: 'list', label: t.drawer.tabs.availableFields.value, showWhen: () => availableFieldsCount.value > 0 },
  ];
  return tabs.filter((tab) => tab.always || tab.showWhen?.());
});

/** WATCHERS */
// Reset to the first tab whenever a new template opens, and clamp the active
// tab if it points at a now-hidden data tab.
watch(() => [props.modelValue, props.template?.id], () => {
  if (props.modelValue) activeTab.value = 'overview';
});

/** FUNCTIONS */
function hasScript(scriptKey: string): boolean {
  return !!(props.template as Record<string, unknown> | null)?.[scriptKey];
}

function viewScript(scriptKey: string, scriptTitle: string): void {
  const content = (props.template as Record<string, unknown> | null)?.[scriptKey];
  if (typeof content !== 'string' || !content) return;
  currentScriptTitle.value = scriptTitle;
  currentScriptContent.value = content;
  showScriptViewer.value = true;
}

function formatDate(value: string | Date | undefined | null): string {
  if (!value) return '-';
  const d = typeof value === 'string' ? new Date(value) : value;
  return date.formatDate(d, 'MMM DD, YYYY HH:mm') || '-';
}
</script>

<style lang="scss" scoped>
.template-details-modal {
  width: 720px;
  max-width: 90vw;
  border-radius: var(--mapex-radius-md);

  &__header {
    padding: 16px 20px;
  }

  &__tabs {
    color: var(--mapex-text-secondary);
  }

  &__body {
    min-height: 320px;
    max-height: 60vh;
  }

  &__footer {
    padding: 12px 20px;
  }
}

.detail-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
}

.detail-field {
  display: flex;
  flex-direction: column;
  gap: 4px;

  &--row {
    flex-direction: row;
    align-items: center;
    justify-content: space-between;
    padding: 6px 0;
    border-bottom: 1px solid var(--mapex-divider);
  }

  &__label {
    font-size: var(--mapex-font-2xs);
    font-weight: var(--mapex-font-weight-semibold);
    text-transform: uppercase;
    letter-spacing: 0.6px;
    color: var(--mapex-text-secondary);
  }

  &__value {
    font-size: var(--mapex-font-sm);
    color: var(--mapex-text-primary);
    word-break: break-word;
  }

  &__code {
    font-family: 'Courier New', monospace;
    font-size: var(--mapex-font-xs);
    color: var(--mapex-text-primary);
    background: var(--mapex-surface-sunken);
    padding: 4px 8px;
    border-radius: var(--mapex-radius-sm);
    word-break: break-all;
  }
}
</style>
