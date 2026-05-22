<script setup lang="ts">
defineOptions({ name: 'WorkflowInstanceDetailsDrawer' });

/** TYPE IMPORTS */
import type {
  WorkflowInstanceDetailsDrawerProps,
  WorkflowInstanceDetailsDrawerEmits,
} from './interfaces';

/** VUE IMPORTS */
import { ref, watch } from 'vue';

/** COMPONENTS */
import { GenericDrawer } from '@components/drawers/common/genericDrawer';
import { DetailChip } from '@components/chips';

/** COMPOSABLES */
import { useWorkflowInstanceListPageTranslations } from '@composables/i18n';

/** SERVICES */
import { apis } from '@services/mapex';

/** PROPS & EMITS */
const props = defineProps<WorkflowInstanceDetailsDrawerProps>();
const emit = defineEmits<WorkflowInstanceDetailsDrawerEmits>();

/** COMPOSABLES & STORES */
const t = useWorkflowInstanceListPageTranslations();

/** STATE */
const instance = ref<Record<string, any> | null>(null);
const loading = ref(false);
const error = ref(false);

/** WATCHERS */
watch(() => props.instanceId, (newId) => {
  if (newId && props.modelValue) {
    void fetchDetails(newId);
  }
}, { immediate: true });

watch(() => props.modelValue, (isOpen) => {
  if (isOpen && props.instanceId) {
    void fetchDetails(props.instanceId);
  } else if (!isOpen) {
    instance.value = null;
    error.value = false;
  }
});

/** FUNCTIONS */

/**
 * Fetch instance details by ID from API
 * @param {string} instanceId - Instance ID to fetch
 * @returns {Promise<void>}
 */
async function fetchDetails(instanceId: string): Promise<void> {
  loading.value = true;
  error.value = false;
  instance.value = null;

  try {
    const data = await apis.workflows.instance.getById({ instanceId });
    instance.value = data;
  } catch {
    error.value = true;
  } finally {
    loading.value = false;
  }
}

/**
 * Handle edit action — emit instance ID and close drawer
 * @returns {void}
 */
function handleEdit(): void {
  if (!instance.value?._id) return;
  emit('edit', instance.value._id);
  emit('update:modelValue', false);
}

/**
 * Get external inputs as array of key-value pairs for display
 * @returns {Array<{ key: string; value: string }>}
 */
function getExternalInputsList(): Array<{ key: string; value: string }> {
  const inputs = instance.value?.externalInputs;
  if (!inputs || typeof inputs !== 'object') return [];
  return Object.entries(inputs).map(([key, val]) => ({
    key,
    value: val === null || val === undefined
      ? '—'
      : typeof val === 'object'
        ? JSON.stringify(val)
        : String(val as string | number | boolean),
  }));
}
</script>

<template>
  <GenericDrawer
    :model-value="modelValue"
    :title="t.drawer.title.value"
    icon="play_circle"
    icon-color="teal-7"
    :width="450"
    @update:model-value="emit('update:modelValue', $event)"
    @close="emit('update:modelValue', false)"
  >
    <!-- Loading -->
    <div v-if="loading" class="q-pa-lg text-center">
      <q-spinner size="3em" class="q-mb-md" color="primary" />
      <div class="drawer-loading-text">{{ t.drawer.loading.value }}</div>
    </div>

    <!-- Error -->
    <div v-else-if="error" class="q-pa-lg">
      <q-banner rounded class="bg-negative text-white">
        <template #avatar>
          <q-icon name="error" color="white" />
        </template>
        {{ t.drawer.error.value }}
      </q-banner>
    </div>

    <!-- Content -->
    <div v-else-if="instance" class="q-px-md q-py-lg">

      <!-- Section: Basic Info -->
      <div class="section q-mb-md">
        <div class="section-header">
          <q-icon name="info" color="primary" size="sm" class="q-mr-sm" />
          <span class="text-subtitle1 text-weight-medium">{{ t.drawer.sections.basicInfo.value }}</span>
        </div>
        <q-separator class="q-my-sm" />

        <div class="field-row">
          <div class="field-label">{{ t.drawer.fields.name.value }}</div>
          <div class="field-value text-weight-medium">{{ instance.name || '—' }}</div>
        </div>

        <div class="field-row">
          <div class="field-label">{{ t.drawer.fields.description.value }}</div>
          <div class="field-value field-value--secondary">{{ instance.description || t.drawer.values.noDescription.value }}</div>
        </div>

        <div class="row q-col-gutter-sm">
          <div class="col-6">
            <div class="field-row">
              <div class="field-label">{{ t.drawer.fields.status.value }}</div>
              <div class="field-value">
                <DetailChip
                  dense
                  size="sm"
                  :color="instance.enabled ? 'positive' : 'negative'"
                  :label="instance.enabled ? t.drawer.values.enabled.value.toUpperCase() : t.drawer.values.disabled.value.toUpperCase()"
                />
              </div>
            </div>
          </div>
          <div class="col-6">
            <div class="field-row">
              <div class="field-label">{{ t.drawer.fields.unique.value }}</div>
              <div class="field-value">
                <DetailChip
                  v-if="instance.uniqueExecution"
                  dense
                  size="sm"
                  color="warning"
                  :label="t.drawer.values.unique.value.toUpperCase()"
                />
                <span v-else class="field-value--muted">{{ t.drawer.values.notApplicable.value }}</span>
              </div>
            </div>
          </div>
        </div>

        <div v-if="instance.workflowUUID" class="field-row">
          <div class="field-label">{{ t.drawer.fields.instanceUUID.value }}</div>
          <div class="field-value field-value--mono">{{ instance.workflowUUID }}</div>
        </div>
      </div>

      <!-- Section: Definition -->
      <div class="section q-mb-md">
        <div class="section-header">
          <q-icon name="account_tree" color="primary" size="sm" class="q-mr-sm" />
          <span class="text-subtitle1 text-weight-medium">{{ t.drawer.sections.workflow.value }}</span>
        </div>
        <q-separator class="q-my-sm" />

        <div class="field-row">
          <div class="field-label">{{ t.drawer.fields.definition.value }}</div>
          <div class="field-value">{{ instance.definitionName || '—' }}</div>
        </div>

        <div class="row q-col-gutter-sm">
          <div class="col-6">
            <div class="field-row">
              <div class="field-label">{{ t.drawer.fields.version.value }}</div>
              <div class="field-value">
                <DetailChip
                  v-if="instance.definitionVersion"
                  dense
                  size="sm"
                  color="blue"
                  :label="`v${instance.definitionVersion}`"
                />
                <span v-else>—</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Section: External Inputs -->
      <div class="section q-mb-md">
        <div class="section-header">
          <q-icon name="input" color="primary" size="sm" class="q-mr-sm" />
          <span class="text-subtitle1 text-weight-medium">{{ t.drawer.sections.inputs.value }}</span>
        </div>
        <q-separator class="q-my-sm" />

        <template v-if="getExternalInputsList().length > 0">
          <div v-for="input in getExternalInputsList()" :key="input.key" class="field-row">
            <div class="field-label">{{ input.key }}</div>
            <div class="field-value">{{ input.value }}</div>
          </div>
        </template>
        <div v-else class="field-value--muted">{{ t.drawer.fields.noInputs.value }}</div>
      </div>

      <!-- Section: Timestamps -->
      <div class="section">
        <div class="section-header">
          <q-icon name="schedule" color="primary" size="sm" class="q-mr-sm" />
          <span class="text-subtitle1 text-weight-medium">{{ t.drawer.sections.timestamps.value }}</span>
        </div>
        <q-separator class="q-my-sm" />

        <div class="row q-col-gutter-sm">
          <div class="col-6">
            <div class="field-row">
              <div class="field-label">{{ t.drawer.fields.created.value }}</div>
              <div class="field-value">{{ instance.created || '—' }}</div>
            </div>
          </div>
          <div class="col-6">
            <div class="field-row">
              <div class="field-label">{{ t.drawer.fields.updated.value }}</div>
              <div class="field-value">{{ instance.updated || '—' }}</div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Footer -->
    <template #footer>
      <q-space />
      <q-btn
        unelevated
        icon="edit"
        color="primary"
        :label="t.drawer.edit.value"
        :disable="!instance"
        @click="handleEdit"
      />
    </template>
  </GenericDrawer>
</template>

<style lang="scss" scoped>
.section {
  .section-header {
    display: flex;
    align-items: center;
    color: var(--q-primary);
    margin-bottom: 8px;
  }
}

.field-row {
  display: flex;
  flex-direction: column;
  padding: 10px 0;
  border-bottom: 1px solid var(--mapex-divider);

  &:last-child {
    border-bottom: none;
  }
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

  &--secondary {
    color: var(--mapex-text-secondary);
  }

  &--muted {
    color: var(--mapex-text-muted);
    font-size: 0.85rem;
  }

  &--mono {
    font-family: monospace;
    font-size: 0.8rem;
    word-break: break-all;
    color: var(--mapex-text-secondary);
  }
}

.drawer-loading-text {
  color: var(--mapex-text-secondary);
}
</style>
