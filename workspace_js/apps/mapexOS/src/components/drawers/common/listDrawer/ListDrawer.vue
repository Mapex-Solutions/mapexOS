<script setup lang="ts">
defineOptions({
  name: 'ListDrawer'
});

import type { ListDrawerProps, ListDrawerEmits } from './interfaces';
import type { ListResponse } from '@mapexos/schemas';

import { ref, watch, onMounted, onBeforeUnmount } from 'vue';
import { date } from 'quasar';

import { DetailChip } from '@components/chips';
import { AppTooltip } from '@components/tooltips';
import { apis } from '@services/mapex';
import { useListDrawerTranslations } from '@composables/i18n';
import { useLogger } from '@composables/useLogger';
import { notifyFail } from '@utils/alert';

const props = defineProps<ListDrawerProps>();
const emit = defineEmits<ListDrawerEmits>();

const t = useListDrawerTranslations();
const logger = useLogger('ListDrawer');

// State
const data = ref<ListResponse | null>(null);
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

// Fetch list details
async function fetchListDetails() {
  if (!props.listId || !apis.mapexOS?.lists) {
    return;
  }

  loading.value = true;
  error.value = false;
  data.value = null;

  try {
    const response = await apis.mapexOS?.lists.getById({ listId: props.listId });
    data.value = response;
  } catch (err: any) {
    logger.error('Error fetching list details:', err);
    error.value = true;
    notifyFail({ message: t.drawer.error.value });
  } finally {
    loading.value = false;
  }
}

// Watch for listId changes
watch(() => props.listId, (newId) => {
  if (newId && props.modelValue) {
    void fetchListDetails();
  }
}, { immediate: true });

// Watch for drawer open/close
watch(() => props.modelValue, (isOpen) => {
  if (isOpen && props.listId) {
    void fetchListDetails();
  } else if (!isOpen) {
    data.value = null;
    error.value = false;
  }
});

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

function close() {
  emit('update:modelValue', false);
}
</script>

<template>
  <q-drawer
    overlay
    bordered
    side="right"
    :model-value="modelValue"
    :width="450"
    @update:model-value="emit('update:modelValue', $event)"
  >
    <!-- Header -->
    <q-toolbar class="drawer-header">
      <q-icon name="list" size="sm" color="primary" class="q-mr-sm" />
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

        <!-- Data Content -->
        <div v-else-if="data" class="q-px-md q-py-lg">
          <!-- Basic Information Section -->
          <div class="section q-mb-lg">
            <div class="section-header q-mb-md">
              <q-icon name="info" size="sm" class="q-mr-xs" />
              <span class="text-subtitle2 text-weight-medium">{{ t.drawer.sections.basicInfo.value }}</span>
            </div>

            <!-- Name (full width) -->
            <div class="field-row">
              <div class="field-label">{{ t.drawer.fields.name.value }}</div>
              <div class="field-value">{{ data.name }}</div>
            </div>

            <!-- Value (full width) -->
            <div class="field-row">
              <div class="field-label">{{ t.drawer.fields.value.value }}</div>
              <div class="field-value text-grey-8">{{ data.value }}</div>
            </div>

            <!-- Category & Type (2 columns) -->
            <div class="row q-col-gutter-sm">
              <div class="col-12">
                <div class="field-row">
                  <div class="field-label">{{ t.drawer.fields.type.value }}</div>
                  <div class="field-value">
                    <DetailChip
                      :value="data.type"
                      color="purple"
                      size="sm"
                      dense
                    />
                  </div>
                </div>
              </div>
            </div>

            <!-- Is System & Is Template (2 columns) -->
            <div class="row q-col-gutter-sm">
              <div class="col-6">
                <div class="field-row">
                  <div class="field-label">{{ t.drawer.fields.isSystem.value }}</div>
                  <div class="field-value">
                    <q-badge :color="data.isSystem ? 'green-6' : 'grey-5'">
                      {{ data.isSystem ? t.drawer.values.yes.value : t.drawer.values.no.value }}
                    </q-badge>
                  </div>
                </div>
              </div>
              <div class="col-6">
                <div class="field-row">
                  <div class="field-label">{{ t.drawer.fields.isTemplate.value }}</div>
                  <div class="field-value">
                    <q-badge :color="data.isTemplate ? 'blue-6' : 'grey-5'">
                      {{ data.isTemplate ? t.drawer.values.yes.value : t.drawer.values.no.value }}
                    </q-badge>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- IDs Section -->
          <div class="section q-mb-lg">
            <div class="section-header q-mb-md">
              <q-icon name="key" size="sm" class="q-mr-xs" />
              <span class="text-subtitle2 text-weight-medium">{{ t.drawer.sections.ids.value }}</span>
            </div>

            <!-- ID (full width) -->
            <div class="field-row">
              <div class="field-label">{{ t.drawer.fields.id.value }}</div>
              <div class="field-value text-grey-8" style="font-family: monospace; font-size: 0.85rem;">
                {{ data.id || '-' }}
              </div>
            </div>

            <!-- Parent Type & Parent Name (2 columns) -->
            <div class="row q-col-gutter-sm">
              <div class="col-6">
                <div class="field-row">
                  <div class="field-label">{{ t.drawer.fields.parentType.value }}</div>
                  <div class="field-value">
                    <DetailChip
                      v-if="data.parentType"
                      :value="data.parentType"
                      color="indigo"
                      size="sm"
                      dense
                    />
                    <span v-else class="text-grey-5">{{ t.drawer.empty.parentType.value }}</span>
                  </div>
                </div>
              </div>
              <div class="col-6">
                <div class="field-row">
                  <div class="field-label">{{ t.drawer.fields.parentName.value }}</div>
                  <div class="field-value">
                    {{ data.parentName || t.drawer.empty.parentName.value }}
                  </div>
                </div>
              </div>
            </div>

            <!-- Organization ID (full width) -->
            <div class="field-row">
              <div class="field-label">{{ t.drawer.fields.orgId.value }}</div>
              <div class="field-value text-grey-8" style="font-family: monospace; font-size: 0.85rem;">
                {{ data.orgId || t.drawer.empty.orgId.value }}
              </div>
            </div>
          </div>

          <!-- Metadata Section -->
          <div v-if="data.metadata && Object.keys(data.metadata).length > 0" class="section q-mb-lg">
            <div class="section-header q-mb-md">
              <q-icon name="data_object" size="sm" class="q-mr-xs" />
              <span class="text-subtitle2 text-weight-medium">{{ t.drawer.sections.metadata.value }}</span>
            </div>

            <div class="field-row">
              <div class="field-value">
                <pre class="metadata-pre">{{ JSON.stringify(data.metadata, null, 2) }}</pre>
              </div>
            </div>
          </div>

          <!-- Timestamps Section -->
          <div class="section">
            <div class="section-header q-mb-md">
              <q-icon name="schedule" size="sm" class="q-mr-xs" />
              <span class="text-subtitle2 text-weight-medium">{{ t.drawer.sections.timestamps.value }}</span>
            </div>

            <!-- Created & Updated (2 columns) -->
            <div class="row q-col-gutter-sm">
              <div class="col-6">
                <div class="field-row">
                  <div class="field-label">{{ t.drawer.fields.created.value }}</div>
                  <div class="field-value">{{ formatDate(data.created) }}</div>
                </div>
              </div>
              <div class="col-6">
                <div class="field-row">
                  <div class="field-label">{{ t.drawer.fields.updated.value }}</div>
                  <div class="field-value">{{ formatDate(data.updated) }}</div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </q-scroll-area>
    </div>
  </q-drawer>
</template>

<style lang="scss" scoped>
// Drawer Header
.drawer-header {
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

// Drawer Content
.drawer-content {
  height: calc(100vh - 64px);
  overflow: hidden;
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

// Metadata Pre
.metadata-pre {
  background: var(--mapex-surface-bg);
  border: 1px solid var(--mapex-divider);
  border-radius: var(--mapex-radius-sm);
  padding: 12px;
  font-size: 0.8rem;
  font-family: 'Courier New', monospace;
  color: var(--mapex-text-secondary);
  overflow-x: auto;
  margin: 0;
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
