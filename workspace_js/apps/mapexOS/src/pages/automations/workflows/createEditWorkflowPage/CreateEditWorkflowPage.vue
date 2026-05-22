<script setup lang="ts">
defineOptions({
  name: 'CreateEditWorkflowPage',
});

/** VUE IMPORTS */
import { ref, computed, watch, onMounted, onBeforeUnmount } from 'vue';
import { useRouter, useRoute, onBeforeRouteLeave } from 'vue-router';

/** COMPONENTS */
import { GeneralTab } from './components/GeneralTab';
import { VariablesTab } from './components/VariablesTab';
import { WorkflowTab } from './components/WorkflowTab';
import { JsonDebugTab } from './components/JsonDebugTab';
import { PluginsTab } from './components/PluginsTab';
import { AppTabs } from '@components/tabs';
import { PageHeader } from '@components/headers';
import { AppTooltip } from '@components/tooltips';

/** COMPOSABLES */
import { useWorkflowEditorState, useWorkflowHistory } from './composables';
import { useLogger } from '@composables/useLogger';
import { useCreateEditWorkflowTranslations } from '@composables/i18n/pages/automations/workflows';

/** UTILS */
import { notifySuccess, notifyFail, dialogConfirm } from '@utils/alert';

/** SERVICES */
import { apis } from '@services/mapex';

/** STORES */
import { usePluginRegistryStore } from '@stores/pluginRegistry';

/** LOCAL IMPORTS (constants and handlers ONLY) */
import { bootWorkflowPlugins } from '@src/components/workflow/constants';
import { bootMarketplacePlugins } from './utils/manifestLoader';

/** COMPOSABLES & STORES */
const router = useRouter();
const route = useRoute();
const logger = useLogger('CreateEditWorkflowPage');
usePluginRegistryStore();
const {
  variablesCount,
  nodesCount,
  getCurrentWorkflow,
  setAllStates,
  resetAllStates,
  nodeValidationErrors,
  validateAllNodes,
  definitionStatus,
  missingPlugins,
} = useWorkflowEditorState();
const { clearHistory } = useWorkflowHistory();
const t = useCreateEditWorkflowTranslations();

/** EDIT MODE DETECTION */
const isEditMode = ref(!!route.params.id);
const workflowId = ref(route.params.id as string | undefined);

/** STATE */

/**
 * Loading state for fetching workflow data in edit mode
 */
const isLoading = ref(false);

/**
 * Currently selected tab
 */
const tab = ref('general');

/**
 * Save operation loading state
 */
const saving = ref(false);

/**
 * Track if there are unsaved changes
 */
const hasUnsavedChanges = ref(false);

/**
 * Validation errors
 */
const validationErrors = ref<string[]>([]);

/** COMPUTED */

/**
 * Available tabs configuration
 * Badges show dynamic counts for variables and nodes
 */
const tabs = computed(() => [
  { id: 'tab-general', name: 'general', label: t.tabs.general.value, icon: 'settings' },
  {
    id: 'tab-data',
    name: 'data',
    label: t.tabs.data.value,
    icon: 'database',
    badge: variablesCount.value > 0 ? variablesCount.value : undefined,
    badgeColor: 'primary',
  },
  {
    id: 'tab-workflow',
    name: 'workflow',
    label: t.tabs.workflow.value,
    icon: 'account_tree',
    badge: nodesCount.value > 0 ? nodesCount.value : undefined,
    badgeColor: 'teal-7',
  },
  {
    id: 'tab-plugins',
    name: 'plugins',
    label: t.tabs.plugins.value,
    icon: 'extension',
  },
  {
    id: 'tab-debug',
    name: 'debug',
    label: t.tabs.jsonDebug.value,
    icon: 'bug_report',
  },
]);

/**
 * Dynamic page title based on mode
 */
const pageTitle = computed(() =>
  isEditMode.value ? t.page.titleEdit.value : t.page.title.value,
);

/**
 * Dynamic page description based on mode
 */
const pageDescription = computed(() =>
  isEditMode.value
    ? t.page.descriptionEdit.value
    : t.page.description.value,
);

/**
 * Dynamic save button label based on mode
 */
const saveButtonLabel = computed(() =>
  isEditMode.value ? t.buttons.update.value : t.buttons.save.value,
);

/**
 * Total validation error count (global + per-node)
 */
const totalErrorCount = computed(() =>
  validationErrors.value.length + Object.values(nodeValidationErrors.value).flat().length,
);

/** WATCHERS */

/**
 * Watch for changes and mark as unsaved.
 * Uses a serialized snapshot to avoid deep watcher pitfalls with Vue Flow sync.
 */
const workflowSnapshot = computed(() => JSON.stringify(getCurrentWorkflow.value));

watch(workflowSnapshot, () => {
  if (!isLoading.value) {
    hasUnsavedChanges.value = true;
  }
  if (validationErrors.value.length > 0 || Object.keys(nodeValidationErrors.value).length > 0) {
    validateWorkflow();
  }
});

/** FUNCTIONS */

/**
 * Load workflow data from API in EDIT mode
 *
 * @returns {Promise<void>}
 */
async function loadWorkflowData(): Promise<void> {
  if (!isEditMode.value || !workflowId.value) return;

  isLoading.value = true;
  try {
    const response = await apis.workflows.definition.getById({ workflowId: workflowId.value });
    setAllStates(response as Parameters<typeof setAllStates>[0]);
  } catch (error) {
    logger.error('Failed to load workflow:', error);
    notifyFail({ message: t.notifications.loadFailed.value });
    void router.push('/workflows');
  } finally {
    isLoading.value = false;
  }
}

/**
 * Validate the current workflow
 *
 * @returns {boolean} True if valid, false otherwise
 */
function validateWorkflow(): boolean {
  const errors: string[] = [];
  const workflow = getCurrentWorkflow.value;

  if (!workflow.name || workflow.name.trim() === '') {
    errors.push(t.validation.nameRequired.value);
  }

  if (workflow.nodes.length === 0) {
    errors.push(t.validation.nodesRequired.value);
  }

  // Per-node validation
  const nodesWithErrors = validateAllNodes();
  if (nodesWithErrors > 0) {
    errors.push(t.resolveValidationError(`nodesHaveErrors::${nodesWithErrors}`));
  }

  validationErrors.value = errors;
  return errors.length === 0;
}

/**
 * Save the workflow to the backend
 *
 * @returns {Promise<boolean>} True if saved successfully
 */
async function saveWorkflow(): Promise<boolean> {
  saving.value = true;

  try {
    const workflow = getCurrentWorkflow.value;

    if (isEditMode.value && workflowId.value) {
      await apis.workflows.definition.update(
        { workflowId: workflowId.value },
        workflow as unknown as Parameters<typeof apis.workflows.definition.update>[1],
      );
    } else {
      const created = await apis.workflows.definition.create(
        workflow as unknown as Parameters<typeof apis.workflows.definition.create>[0],
      );
      if (created._id) {
        workflowId.value = created._id;
        isEditMode.value = true;
      }
    }

    notifySuccess({
      message: isEditMode.value ? t.notifications.updatedSuccess.value : t.notifications.savedSuccess.value,
    });

    hasUnsavedChanges.value = false;
    return true;
  } catch (error) {
    logger.error('Failed to save workflow:', error);
    notifyFail({ message: t.notifications.saveFailed.value });
    return false;
  } finally {
    saving.value = false;
  }
}

/**
 * Save and close — returns to workflows list
 *
 * @returns {Promise<void>}
 */
async function saveAndClose(): Promise<void> {
  if (!validateWorkflow()) {
    notifyFail({ message: t.validation.fixErrors.value });
    return;
  }

  const saved = await saveWorkflow();
  if (saved) {
    void router.push('/workflows');
  }
}

/**
 * Handle cancel action
 * Prompts user if there are unsaved changes
 *
 * @returns {Promise<void>}
 */
async function handleCancel(): Promise<void> {
  if (hasUnsavedChanges.value) {
    const confirmed = await dialogConfirm({
      title: t.dialogs.unsavedChangesTitle.value,
      message: t.dialogs.unsavedChangesCancel.value,
    });

    if (confirmed) {
      void router.push('/workflows');
    }
  } else {
    void router.push('/workflows');
  }
}

/**
 * Handle keyboard shortcuts
 *
 * @param {KeyboardEvent} e - Keyboard event
 * @returns {void}
 */
function handleKeydown(e: KeyboardEvent): void {
  if ((e.ctrlKey || e.metaKey) && e.key === 's') {
    e.preventDefault();
    void saveAndClose();
  }
}

/** LIFECYCLE HOOKS */

onMounted(() => {
  // Boot workflow plugins (registers core node types in catalog)
  const registry = usePluginRegistryStore();
  bootWorkflowPlugins((plugin) => registry.registerPlugin(plugin));

  // Boot marketplace plugins (async, falls back to core-only on failure)
  void bootMarketplacePlugins((plugin) => registry.registerPlugin(plugin));

  if (isEditMode.value) {
    void loadWorkflowData();
  } else {
    resetAllStates();
  }

  // Clear undo/redo history after initial state is set
  clearHistory();

  window.addEventListener('keydown', handleKeydown);
});

onBeforeUnmount(() => {
  window.removeEventListener('keydown', handleKeydown);
});

onBeforeRouteLeave((_to, _from, next) => {
  if (hasUnsavedChanges.value) {
    void dialogConfirm({
      title: t.dialogs.unsavedChangesTitle.value,
      message: t.dialogs.unsavedChangesLeave.value,
    }).then((confirmed) => {
      next(confirmed ? undefined : false);
    });
  } else {
    next();
  }
});
</script>

<template>
  <q-page class="q-pt-lg">
    <!-- Loading Spinner for EDIT mode data fetch -->
    <div v-if="isLoading" class="flex flex-center q-pa-xl">
      <div class="column items-center">
        <q-spinner color="primary" size="50px" />
        <span class="q-mt-md text-grey-7">{{ t.page.loading.value }}</span>
      </div>
    </div>

    <!-- Main Content -->
    <div v-else>
      <!-- Header Section -->
      <div id="workflow-builder-header">
        <PageHeader
          icon="account_tree"
          icon-color="primary"
          :title="pageTitle"
          :description="pageDescription"
          :button="{
            label: t.page.back.value,
            icon: 'arrow_back',
            flat: true,
            onClick: handleCancel,
          }"
        >
          <template #actions>
            <!-- Validate Button -->
            <q-btn
              flat
              rounded
              no-caps
              class="q-px-md q-mr-sm"
              color="grey-7"
              icon="check_circle"
              :label="t.buttons.validate.value"
              @click="validateWorkflow"
            >
              <AppTooltip :content="t.tooltips.validate.value" />
            </q-btn>

            <!-- Save Button -->
            <div class="relative-position">
              <q-btn
                unelevated
                rounded
                no-caps
                class="q-px-md"
                color="primary"
                icon="save"
                :label="saveButtonLabel"
                :loading="saving"
                :disable="saving"
                @click="saveAndClose"
                :ripple="false"
              >
                <AppTooltip :content="totalErrorCount > 0 ? t.tooltips.saveDisabled.value : t.tooltips.saveReady.value" />
              </q-btn>

              <!-- Validation Error Badge -->
              <q-badge
                v-if="totalErrorCount > 0"
                color="negative"
                floating
                rounded
                style="top: -4px; right: -4px;"
              >
                {{ totalErrorCount }}
              </q-badge>

              <!-- Unsaved Changes Badge -->
              <q-badge
                v-else-if="hasUnsavedChanges"
                color="warning"
                text-color="dark"
                floating
                rounded
                style="top: -4px; right: -4px; font-size: 8px; padding: 2px 4px;"
              >
                !
                <AppTooltip :content="t.tooltips.unsavedChanges.value" />
              </q-badge>
            </div>
          </template>
        </PageHeader>
      </div>

      <!-- Plugin Missing Warning -->
      <q-banner
        v-if="definitionStatus === 'plugin_missing' && missingPlugins.length > 0"
        dense
        rounded
        class="q-mx-md q-mt-sm workflow-status-banner"
      >
        <template #avatar>
          <q-icon name="warning" color="warning" />
        </template>
        <span class="text-body2">
          {{ t.statusBanner.pluginMissing.value }}
        </span>
        <span class="text-caption text-weight-medium q-ml-xs">
          {{ missingPlugins.join(', ') }}
        </span>
      </q-banner>

      <!-- Tabs -->
      <div class="workflow-builder">
        <div id="workflow-builder-tabs">
          <AppTabs v-model="tab" :tabs="tabs" :separator="false" />
        </div>

        <q-tab-panels v-model="tab" animated keep-alive>
          <!-- GENERAL -->
          <q-tab-panel name="general">
            <GeneralTab />
          </q-tab-panel>

          <!-- DATA (Inputs / State / Capture) -->
          <q-tab-panel name="data">
            <VariablesTab />
          </q-tab-panel>

          <!-- WORKFLOW (Canvas) -->
          <q-tab-panel name="workflow" class="workflow-tab-panel">
            <WorkflowTab />
          </q-tab-panel>

          <!-- PLUGINS -->
          <q-tab-panel name="plugins">
            <PluginsTab />
          </q-tab-panel>

          <!-- JSON DEBUG -->
          <q-tab-panel name="debug">
            <JsonDebugTab />
          </q-tab-panel>
        </q-tab-panels>
      </div>
    </div>
  </q-page>
</template>

<style lang="scss" scoped>
.workflow-status-banner {
  background: color-mix(in srgb, var(--q-warning) 12%, var(--mapex-surface-bg));
  border: 1px solid color-mix(in srgb, var(--q-warning) 30%, transparent);
}

.workflow-builder {
  max-width: 100%;
  margin: 0 auto;
}

.workflow-tab-panel {
  padding: 0 !important;
}

.q-tab {
  min-height: 48px;
  padding: 0 24px;
}

:deep() {
  .q-tab__content {
    min-width: unset;
  }
}
</style>
