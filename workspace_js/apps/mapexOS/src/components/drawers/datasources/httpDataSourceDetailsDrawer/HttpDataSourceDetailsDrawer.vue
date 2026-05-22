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
      <q-icon name="settings_input_antenna" size="sm" color="primary" class="q-mr-sm" />
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
          <q-spinner size="3em" color="primary" class="q-mb-md" />
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

        <!-- Data Source Data -->
        <div v-else-if="dataSource" class="q-px-md q-py-lg">

          <!-- Basic Information Section -->
          <div class="section q-mb-md">
            <div class="section-header">
              <q-icon name="info" color="primary" size="sm" class="q-mr-sm" />
              <span class="text-subtitle1 text-weight-medium">{{ t.drawer.sections.basicInfo.value }}</span>
            </div>
            <q-separator class="q-my-sm" />

            <div class="field-row">
              <div class="field-label">{{ t.drawer.fields.name.value }}</div>
              <div class="field-value text-weight-medium">{{ dataSource?.name || '-' }}</div>
            </div>

            <div class="field-row">
              <div class="field-label">{{ t.drawer.fields.status.value }}</div>
              <div class="field-value">
                <DetailChip
                  :color="dataSource?.enabled ? 'positive' : 'negative'"
                  size="sm"
                  :label="dataSource?.enabled ? t.drawer.status.active.value.toUpperCase() : t.drawer.status.inactive.value.toUpperCase()"
                />
              </div>
            </div>

            <div class="field-row">
              <div class="field-label">{{ t.drawer.fields.description.value }}</div>
              <div class="field-value text-grey-8">
                {{ dataSource?.description || t.drawer.empty.description.value }}
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

            <!-- Mode & Protocol (2 columns) -->
            <div class="row q-col-gutter-sm">
              <div class="col-6">
                <div class="field-row">
                  <div class="field-label">{{ t.drawer.fields.mode.value }}</div>
                  <div class="field-value">
                    <DetailChip
                      :color="dataSource?.mode === 'pull' ? 'green' : 'orange'"
                      size="sm"
                      :label="dataSource?.mode === 'pull' ? t.drawer.modes.pull.value.toUpperCase() : t.drawer.modes.push.value.toUpperCase()"
                    />
                  </div>
                </div>
              </div>
              <div class="col-6">
                <div class="field-row">
                  <div class="field-label">{{ t.drawer.fields.protocol.value }}</div>
                  <div class="field-value">
                    <DetailChip
                      color="default"
                      size="sm"
                      :label="dataSource?.protocol?.toUpperCase() || 'HTTP'"
                    />
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- Authentication Section -->
          <div class="section q-mb-md">
            <div class="section-header">
              <q-icon name="lock" color="primary" size="sm" class="q-mr-sm" />
              <span class="text-subtitle1 text-weight-medium">{{ t.drawer.sections.authentication.value }}</span>
            </div>
            <q-separator class="q-my-sm" />

            <div class="field-row">
              <div class="field-label">{{ t.drawer.fields.authType.value }}</div>
              <div class="field-value">
                <DetailChip
                  color="blue"
                  size="sm"
                  :label="getAuthTypeLabel(dataSource?.auth?.type)"
                />
              </div>
            </div>
          </div>

          <!-- Asset Binding Section -->
          <div class="section q-mb-md">
            <div class="section-header">
              <q-icon name="link" color="primary" size="sm" class="q-mr-sm" />
              <span class="text-subtitle1 text-weight-medium">{{ t.drawer.sections.assetBinding.value }}</span>
            </div>
            <q-separator class="q-my-sm" />

            <div class="field-row">
              <div class="field-label">{{ t.drawer.fields.assetBindType.value }}</div>
              <div class="field-value">
                <DetailChip
                  :color="dataSource?.assetBind?.type === 'fixedAssetId' ? 'blue' : 'purple'"
                  size="sm"
                  :label="dataSource?.assetBind?.type === 'fixedAssetId' ? 'FIXED' : 'DYNAMIC'"
                />
              </div>
            </div>

            <div v-if="dataSource?.assetBind?.type === 'fixedAssetId'" class="field-row">
              <div class="field-label">{{ t.drawer.fields.fixedAssetId.value }}</div>
              <div class="field-value">
                <DetailChip
                  color="default"
                  size="sm"
                  :label="dataSource?.assetBind?.data?.assetId || '-'"
                />
              </div>
            </div>

            <div v-if="dataSource?.assetBind?.type === 'uuidField'" class="field-row">
              <div class="field-label">{{ t.drawer.fields.uuidField.value }}</div>
              <div class="field-value">
                <DetailChip
                  color="default"
                  size="sm"
                  :label="dataSource?.assetBind?.data?.uuidField?.join('.') || '-'"
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
                  <div class="field-value">{{ formatDate(dataSource?.created) }}</div>
                </div>
              </div>
              <div class="col-6">
                <div class="field-row">
                  <div class="field-label">{{ t.drawer.fields.updated.value }}</div>
                  <div class="field-value">{{ formatDate(dataSource?.updated) }}</div>
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
        :disable="!dataSource"
        @click="handleEdit"
      />
    </div>
  </q-drawer>
</template>

<script setup lang="ts">
defineOptions({
  name: 'HttpDataSourceDetailsDrawer'
});

/** TYPE IMPORTS */
import type { DataSourceResponse } from '@mapexos/schemas';
import type { HttpDataSourceDetailsDrawerProps, HttpDataSourceDetailsDrawerEmits } from './interfaces';

/** VUE IMPORTS */
import { ref, watch, onMounted, onBeforeUnmount } from 'vue';
import { date } from 'quasar';

/** COMPONENTS */
import { DetailChip } from '@components/chips';
import { AppTooltip } from '@components/tooltips';

/** SERVICES */
import { apis } from '@services/mapex';

/** COMPOSABLES */
import { useHttpDataSourcesTranslations } from '@composables/i18n';
import { useLogger } from '@composables/useLogger';

/** UTILS */
import { notifyFail } from '@utils/alert';

const props = defineProps<HttpDataSourceDetailsDrawerProps>();
const emit = defineEmits<HttpDataSourceDetailsDrawerEmits>();

// Translations
const t = useHttpDataSourcesTranslations();
const logger = useLogger('HttpDataSourceDetailsDrawer');

// State
const dataSource = ref<DataSourceResponse | null>(null);
const loading = ref(false);
const error = ref(false);

// Handle ESC key to close drawer
function handleEscKey(event: KeyboardEvent) {
  if (event.key === 'Escape' && props.modelValue) {
    close();
  }
}

onMounted(() => {
  window.addEventListener('keydown', handleEscKey);
});

onBeforeUnmount(() => {
  window.removeEventListener('keydown', handleEscKey);
});

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
  if (!dataSource.value?.id) return;
  emit('edit', dataSource.value.id);
  close();
}

// Format date using Quasar date utils
function formatDate(dateValue?: string): string {
  if (!dateValue) return '-';

  try {
    const dateObj = new Date(dateValue);
    return date.formatDate(dateObj, 'MMM DD, YYYY HH:mm');
  } catch {
    return '-';
  }
}

// Get auth type label
function getAuthTypeLabel(authType?: string): string {
  if (!authType) return t.drawer.authTypes.none.value.toUpperCase();

  const typeMap: Record<string, any> = {
    'none': t.drawer.authTypes.none,
    'apiKey': t.drawer.authTypes.apiKey,
    'ip_whitelist': t.drawer.authTypes.ipWhitelist,
    'jwt': t.drawer.authTypes.jwt,
    'oauth2': t.drawer.authTypes.oauth2,
  };

  return typeMap[authType]?.value?.toUpperCase() || authType.toUpperCase();
}

// Fetch data source details
async function fetchDataSource() {
  if (!props.dataSourceId || !apis.httpGateway) {
    return;
  }

  loading.value = true;
  error.value = false;
  dataSource.value = null;

  try {
    const response = await apis.httpGateway.datasource.getById({
      dataSourceId: props.dataSourceId,
    });

    dataSource.value = response;
  } catch (err: any) {
    logger.error('Error fetching data source:', err);
    error.value = true;
    notifyFail({ message: t.drawer.error.value });
  } finally {
    loading.value = false;
  }
}

// Watch for changes in dataSourceId
watch(() => props.dataSourceId, (newId) => {
  if (newId && props.modelValue) {
    void fetchDataSource();
  }
}, { immediate: true });

// Watch for drawer opening
watch(() => props.modelValue, (isOpen) => {
  if (isOpen && props.dataSourceId) {
    void fetchDataSource();
  } else if (!isOpen) {
    // Reset state when drawer closes
    dataSource.value = null;
    error.value = false;
  }
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
