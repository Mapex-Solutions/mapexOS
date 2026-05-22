<script setup lang="ts">
defineOptions({
  name: 'WorkflowDetailsDrawer'
});

/** TYPE IMPORTS */
import type { WorkflowDetailsDrawerProps, WorkflowDetailsDrawerEmits } from './interfaces/workflowDetailsDrawer.interface';
import type { WorkflowListItem } from 'src/pages/automations/workflows/workflowListPage/interfaces/workflowListPage.interface';

/** VUE IMPORTS */
import { ref, watch, onMounted, onBeforeUnmount } from 'vue';
import { date } from 'quasar';

/** COMPONENTS */
import { DetailChip } from '@components/chips';
import { AppTooltip } from '@components/tooltips';

/** COMPOSABLES */
import { useWorkflowListPageTranslations } from '@composables/i18n';

/** SERVICES */
import { apis } from '@services/mapex';

/** PROPS & EMITS */
const props = defineProps<WorkflowDetailsDrawerProps>();
const emit = defineEmits<WorkflowDetailsDrawerEmits>();

/** COMPOSABLES & STORES */
const t = useWorkflowListPageTranslations();

/** STATE */
const workflow = ref<WorkflowListItem | null>(null);
const loading = ref(false);
const error = ref(false);

/** WATCHERS */
watch(() => props.workflowId, (newId) => {
  if (newId && props.modelValue) {
    void fetchWorkflowDetails(newId);
  }
}, { immediate: true });

watch(() => props.modelValue, (isOpen) => {
  if (isOpen && props.workflowId) {
    void fetchWorkflowDetails(props.workflowId);
  } else if (!isOpen) {
    workflow.value = null;
    error.value = false;
  }
});

/** FUNCTIONS */

/**
 * Fetch workflow details by ID from API
 *
 * @param {string} workflowId - Workflow ID to fetch
 * @returns {Promise<void>}
 */
async function fetchWorkflowDetails(workflowId: string): Promise<void> {
  loading.value = true;
  error.value = false;
  workflow.value = null;

  try {
    const def = await apis.workflows.definition.getById({ workflowId });

    workflow.value = {
      id: def._id || '',
      name: def.name || '',
      description: def.description || '',
      enabled: def.enabled ?? true,
      isTemplate: def.isTemplate ?? false,
      definitionVersion: def.definitionVersion || 1,
      nodesCount: def.nodes?.length || 0,
      edgesCount: def.edges?.length || 0,
      timezone: (def.timezone as any)?.value || 'UTC',
      created: def.created || '',
      updated: def.updated || '',
      status: (def.status as WorkflowListItem['status']) || 'valid',
      missingPlugins: def.missingPlugins || [],
      pluginsCount: def.installedPlugins?.length || 0,
    };
  } catch {
    error.value = true;
  } finally {
    loading.value = false;
  }
}

/**
 * Format date using Quasar date utils
 *
 * @param {string} dateValue - ISO date string to format
 * @returns {string} Formatted date
 */
function formatDate(dateValue?: string): string {
  if (!dateValue) return '-';

  try {
    const dateObj = new Date(dateValue);
    return date.formatDate(dateObj, 'MMM DD, YYYY HH:mm');
  } catch {
    return '-';
  }
}

/**
 * Close drawer
 *
 * @returns {void}
 */
function close(): void {
  emit('update:modelValue', false);
}

/**
 * Handle edit action — emit workflow ID and close
 *
 * @returns {void}
 */
function handleEdit(): void {
  if (!workflow.value?.id) return;
  emit('edit', workflow.value.id);
  close();
}

/**
 * Handle ESC key to close drawer
 *
 * @param {KeyboardEvent} event - Keyboard event
 * @returns {void}
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
      <q-icon name="account_tree" size="sm" class="q-mr-sm" color="primary" />
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

        <!-- Workflow Data -->
        <div v-else-if="workflow" class="q-px-md q-py-lg">

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
              <div class="field-value text-weight-medium">{{ workflow.name || '-' }}</div>
            </div>

            <!-- Description (full width) -->
            <div class="field-row">
              <div class="field-label">{{ t.drawer.fields.description.value }}</div>
              <div class="field-value text-grey-8">
                {{ workflow.description || t.drawer.values.noDescription.value }}
              </div>
            </div>

            <!-- Status & Template (2 columns) -->
            <div class="row q-col-gutter-sm">
              <div class="col-6">
                <div class="field-row">
                  <div class="field-label">{{ t.drawer.fields.status.value }}</div>
                  <div class="field-value">
                    <DetailChip
                      :color="workflow.enabled ? 'positive' : 'negative'"
                      size="sm"
                      :label="workflow.enabled ? t.drawer.values.enabled.value.toUpperCase() : t.drawer.values.disabled.value.toUpperCase()"
                    />
                  </div>
                </div>
              </div>
              <div class="col-6">
                <div class="field-row">
                  <div class="field-label">{{ t.drawer.fields.isTemplate.value }}</div>
                  <div class="field-value">
                    <DetailChip
                      :color="workflow.isTemplate ? 'blue' : 'grey'"
                      size="sm"
                      :label="workflow.isTemplate ? t.drawer.values.yes.value.toUpperCase() : t.drawer.values.no.value.toUpperCase()"
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

            <!-- Version (full width) -->
            <div class="field-row">
              <div class="field-label">{{ t.drawer.fields.version.value }}</div>
              <div class="field-value">
                <DetailChip
                  color="blue"
                  size="sm"
                  :label="`v${workflow.definitionVersion}`"
                />
              </div>
            </div>

            <!-- Nodes & Edges (2 columns) -->
            <div class="row q-col-gutter-sm">
              <div class="col-6">
                <div class="field-row">
                  <div class="field-label">{{ t.drawer.fields.nodesCount.value }}</div>
                  <div class="field-value">{{ workflow.nodesCount }}</div>
                </div>
              </div>
              <div class="col-6">
                <div class="field-row">
                  <div class="field-label">{{ t.drawer.fields.edgesCount.value }}</div>
                  <div class="field-value">{{ workflow.edgesCount }}</div>
                </div>
              </div>
            </div>

            <!-- Timezone (full width) -->
            <div class="field-row">
              <div class="field-label">{{ t.drawer.fields.timezone.value }}</div>
              <div class="field-value">
                <DetailChip
                  icon="schedule"
                  color="default"
                  size="sm"
                  :label="workflow.timezone"
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
                  <div class="field-value">{{ formatDate(workflow.created) }}</div>
                </div>
              </div>
              <div class="col-6">
                <div class="field-row">
                  <div class="field-label">{{ t.drawer.fields.updated.value }}</div>
                  <div class="field-value">{{ formatDate(workflow.updated) }}</div>
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
        :disable="!workflow"
        @click="handleEdit"
      />
    </div>
  </q-drawer>
</template>

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
  box-shadow: var(--mapex-shadow-sm);
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
    transition: background 0.2s ease;

    &:hover {
      background: rgba(var(--q-primary-rgb), 0.5);
    }
  }
}
</style>
