<script setup lang="ts">
defineOptions({ name: 'WorkflowConfig' });

/** TYPE IMPORTS */
import type { WorkflowConfigProps, WorkflowConfigEmits } from './interfaces/WorkflowConfig.interface';
import type { InstanceResponse } from '@mapexos/schemas';

/** VUE IMPORTS */
import { ref, computed } from 'vue';

/** COMPONENTS */
import { WorkflowInstanceSelectorDrawer } from '@components/drawers/automations/workflowInstanceSelectorDrawer';

/** PROPS & EMITS */
const props = defineProps<WorkflowConfigProps>();
const emit = defineEmits<WorkflowConfigEmits>();

/** STATE */
const instanceDrawerOpen = ref(false);
const selectedInstanceName = ref('');
const selectedUniqueExecution = ref(false);
const selectedWorkflowUUID = ref('');
const selectedExternalSignals = ref<string[]>([]);

/** COMPUTED */

/**
 * Whether signal selection is visible (signal + signalOrStart modes)
 */
const showSignalSelect = computed(() =>
  props.workflow.mode === 'signal' || props.workflow.mode === 'signalOrStart',
);

/**
 * Whether UUID field is visible (only when instance has uniqueExecution)
 */
const showUUIDField = computed(() => selectedUniqueExecution.value);

/**
 * Current instanceId from data
 */
const currentInstanceId = computed(() =>
  (props.workflow.data?.instanceId as string) || '',
);

/**
 * Current signalName from data
 */
const currentSignalName = computed(() =>
  (props.workflow.data?.signalName as string) || '',
);

/** FUNCTIONS */

/**
 * Update mode and reset data
 * @param {string} newMode - New delivery mode
 * @returns {void}
 */
function updateMode(newMode: string): void {
  emit('update:workflow', {
    ...props.workflow,
    mode: newMode as 'newInstance' | 'signal' | 'signalOrStart',
    data: {},
  });
  selectedInstanceName.value = '';
  selectedUniqueExecution.value = false;
  selectedWorkflowUUID.value = '';
  selectedExternalSignals.value = [];
}

/**
 * Handle instance selection from drawer.
 * Populates instanceId, UUID (if unique), and signal options.
 * @param {InstanceResponse} instance - Selected instance
 * @returns {void}
 */
function onInstanceSelected(instance: InstanceResponse): void {
  const instanceId = instance._id || '';
  selectedInstanceName.value = instance.name || instanceId;
  selectedUniqueExecution.value = instance.uniqueExecution ?? false;
  selectedWorkflowUUID.value = instance.workflowUUID || '';
  selectedExternalSignals.value = (instance as any).externalSignals || [];

  const newData: Record<string, any> = { instanceId };
  if (selectedUniqueExecution.value && selectedWorkflowUUID.value) {
    newData.workflowUUID = selectedWorkflowUUID.value;
  }

  emit('update:workflow', { ...props.workflow, data: newData });
}

/**
 * Update signal name in data
 * @param {string} signalName - Selected signal name
 * @returns {void}
 */
function updateSignalName(signalName: string): void {
  const newData = { ...props.workflow.data, signalName };
  emit('update:workflow', { ...props.workflow, data: newData });
}


</script>

<template>
  <!-- Mode -->
  <div class="col-12">
    <q-select
      :model-value="workflow.mode"
      outlined
      dense
      emit-value
      map-options
      :options="[
        { label: t.createEdit.routersStep.routerCard.workflow.modeOptions.newInstance.value, value: 'newInstance', description: t.createEdit.routersStep.routerCard.workflow.modeDescriptions.newInstance.value },
        { label: t.createEdit.routersStep.routerCard.workflow.modeOptions.signal.value, value: 'signal', description: t.createEdit.routersStep.routerCard.workflow.modeDescriptions.signal.value },
        { label: t.createEdit.routersStep.routerCard.workflow.modeOptions.signalOrStart.value, value: 'signalOrStart', description: t.createEdit.routersStep.routerCard.workflow.modeDescriptions.signalOrStart.value },
      ]"
      :hint="t.createEdit.routersStep.routerCard.workflow.modeHint.value"
      @update:model-value="(val: string) => updateMode(val)"
    >
      <template #prepend>
        <q-icon name="account_tree" size="xs" class="icon-primary" />
      </template>
      <template #option="scope">
        <q-item v-bind="scope.itemProps">
          <q-item-section>
            <q-item-label>{{ scope.opt.label }}</q-item-label>
            <q-item-label caption>{{ scope.opt.description }}</q-item-label>
          </q-item-section>
        </q-item>
      </template>
    </q-select>
  </div>

  <!-- Workflow Instance Selector -->
  <div class="col-12">
    <q-input
      :model-value="selectedInstanceName || currentInstanceId"
      outlined
      dense
      readonly
      class="cursor-pointer"
      :label="t.createEdit.routersStep.routerCard.workflow.workflowId.value + ' *'"
      :hint="t.createEdit.routersStep.routerCard.workflow.workflowIdHint.value"
      @click="instanceDrawerOpen = true"
    >
      <template #prepend>
        <q-icon name="play_circle" size="xs" class="icon-primary" />
      </template>
      <template #append>
        <q-icon name="chevron_right" class="icon-muted" />
      </template>
    </q-input>
  </div>

  <!-- Workflow UUID (read-only, only if uniqueExecution) -->
  <div v-if="showUUIDField" class="col-12">
    <q-input
      :model-value="selectedWorkflowUUID"
      outlined
      dense
      readonly
      :label="t.createEdit.routersStep.routerCard.workflow.workflowUUID.value"
      :hint="t.createEdit.routersStep.routerCard.workflow.workflowUUIDHint.value"
    >
      <template #prepend>
        <q-icon name="fingerprint" size="xs" class="icon-primary" />
      </template>
    </q-input>
  </div>

  <!-- Signal Name (signal + signalOrStart) -->
  <div v-if="showSignalSelect" class="col-12">
    <q-select
      v-if="selectedExternalSignals.length > 0"
      :model-value="currentSignalName"
      outlined
      dense
      emit-value
      :options="selectedExternalSignals"
      :label="t.createEdit.routersStep.routerCard.workflow.signalName.value + ' *'"
      :hint="t.createEdit.routersStep.routerCard.workflow.signalNameHint.value"
      @update:model-value="(val: string) => updateSignalName(val)"
    >
      <template #prepend>
        <q-icon name="bolt" size="xs" class="icon-primary" />
      </template>
    </q-select>
    <q-input
      v-else
      :model-value="currentSignalName"
      outlined
      dense
      :label="t.createEdit.routersStep.routerCard.workflow.signalName.value + ' *'"
      :hint="t.createEdit.routersStep.routerCard.workflow.signalNameHint.value"
      @update:model-value="(val: string | number | null) => updateSignalName(String(val ?? ''))"
    >
      <template #prepend>
        <q-icon name="bolt" size="xs" class="icon-primary" />
      </template>
    </q-input>
  </div>

  <!-- Instance Selector Drawer -->
  <WorkflowInstanceSelectorDrawer
    v-model="instanceDrawerOpen"
    :selected-instance-id="currentInstanceId"
    @select="onInstanceSelected"
    @cancel="instanceDrawerOpen = false"
  />
</template>

<style lang="scss" scoped>
.icon-primary {
  color: var(--q-primary);
}

.icon-muted {
  color: var(--mapex-text-muted);
}
</style>
