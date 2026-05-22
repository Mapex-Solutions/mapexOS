<script setup lang="ts">
/** VUE IMPORTS */
import { ref } from 'vue';

/** COMPONENTS */
import ExternalVariables from '../ExternalVariables/ExternalVariables.vue';
import WorkflowVariables from '../WorkflowVariables/WorkflowVariables.vue';
import CaptureFields from '../CaptureFields/CaptureFields.vue';
import ExternalSignals from '../ExternalSignals/ExternalSignals.vue';

/** COMPOSABLES */
import { useCreateEditWorkflowTranslations } from '@composables/i18n/pages/automations/workflows';

/** STATE */

/**
 * Active sub-tab
 */
const subTab = ref<'inputs' | 'state' | 'captureFields' | 'signals'>('inputs');

/** COMPOSABLES & STORES */
const t = useCreateEditWorkflowTranslations();

/**
 * Resolve the help text for the active sub-tab
 */
const SUB_TAB_HELP: Record<string, () => string> = {
  inputs: () => t.variablesTab.inputsHelp.value,
  state: () => t.variablesTab.stateHelp.value,
  captureFields: () => t.variablesTab.captureFieldsHelp.value,
  signals: () => t.variablesTab.signalsHelp.value,
};
</script>

<template>
  <div class="variables-tab">
    <!-- Sub-tabs -->
    <div class="row items-center q-mb-lg">
      <q-tabs
        v-model="subTab"
        dense
        no-caps
        class="text-grey-8 variables-subtabs"
        active-color="primary"
        indicator-color="primary"
      >
        <q-tab name="inputs" icon="input" :label="t.variablesTab.inputsTab.value" />
        <q-tab name="signals" icon="sensors" :label="t.variablesTab.signalsTab.value" />
        <q-tab name="state" icon="data_object" :label="t.variablesTab.stateTab.value" />
        <q-tab name="captureFields" icon="analytics" :label="t.variablesTab.captureFieldsTab.value" />
      </q-tabs>
      <q-space />
      <span class="text-caption text-grey-6">
        <q-icon name="info" size="xs" class="q-mr-xs" />
        {{ SUB_TAB_HELP[subTab]?.() }}
      </span>
    </div>

    <!-- External Variables (Inputs) -->
    <div v-if="subTab === 'inputs'">
      <ExternalVariables />
    </div>

    <!-- Workflow Variables -->
    <div v-else-if="subTab === 'state'">
      <WorkflowVariables />
    </div>

    <!-- Capture Fields -->
    <div v-else-if="subTab === 'captureFields'">
      <CaptureFields />
    </div>

    <!-- External Signals -->
    <div v-else-if="subTab === 'signals'">
      <ExternalSignals />
    </div>
  </div>
</template>

<style lang="scss" scoped>
.variables-subtabs {
  background: transparent;
  border-radius: var(--mapex-radius-md);

  :deep(.q-tab) {
    min-height: 36px;
    padding: 0 16px;
    border-radius: var(--mapex-radius-sm);
    margin-right: 4px;

    &:hover:not(.q-tab--active) {
      background: var(--mapex-surface-bg);
    }
  }

  :deep(.q-tabs__content) {
    padding: 4px;
    background: var(--mapex-page-bg);
    border-radius: var(--mapex-radius-md);
  }

  :deep(.q-tab--active) {
    background: var(--mapex-surface-elevated);
    box-shadow: var(--mapex-shadow-xs);
  }

  :deep(.q-tab__indicator) {
    display: none;
  }
}
</style>
