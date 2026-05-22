import { markRaw } from 'vue';
import type { WorkflowPlugin } from '../interfaces';
import { i18nInstance } from 'src/boot/i18n';

import GenericWorkflowNode from '../nodes/GenericWorkflowNode/GenericWorkflowNode.vue';
import TextNoteNode from '../nodes/textNoteNode/TextNoteNode.vue';
import GotoNode from '../nodes/gotoNode/GotoNode.vue';
import GroupFrameNode from '../nodes/groupFrameNode/GroupFrameNode.vue';
import EndNode from '../nodes/endNode/EndNode.vue';

import { enUS, ptBR } from '../i18n';

import CodeNodeConfig from '../nodes/codeNode/configs/CodeNodeConfig.vue';
import ConditionNodeConfig from '../nodes/conditionNode/configs/ConditionNodeConfig.vue';
import TriggerEventNodeConfig from '../nodes/triggerEventNode/configs/TriggerEventNodeConfig.vue';
import WaitForNodeConfig from '../nodes/waitForNode/configs/WaitForNodeConfig.vue';
import WaitSignalNodeConfig from '../nodes/waitSignalNode/configs/WaitSignalNodeConfig.vue';
import EndNodeConfig from '../nodes/endNode/configs/EndNodeConfig.vue';
import LoopNodeConfig from '../nodes/loopNode/configs/LoopNodeConfig.vue';
import SubworkflowNodeConfig from '../nodes/subWorkflowNode/configs/SubworkflowNodeConfig.vue';
import SwitchNodeConfig from '../nodes/switchNode/configs/SwitchNodeConfig.vue';
import SetStateNodeConfig from '../nodes/setStateNode/configs/SetStateNodeConfig.vue';
import GotoNodeConfig from '../nodes/gotoNode/configs/GotoNodeConfig.vue';
import GroupFrameNodeConfig from '../nodes/groupFrameNode/configs/GroupFrameNodeConfig.vue';

import type { SwitchCase } from '../interfaces';

import {
  validateTriggerEvent,
  validateCondition,
  validateSetState,
  validateCode,
  validateSwitch,
  validateSubworkflow,
  validateLoop,
  validateEnd,
  validateGoto,
  validateWaitSignal,
  validateWaitFor,
} from '../validators';

/**
 * Start node type identifier
 */
export const START_NODE_TYPE = 'core/start';

/**
 * Default start node ID (auto-created on every workflow)
 */
export const START_NODE_ID = '__start__';

/**
 * Core Triggers plugin — entry point node types
 */
export const CORE_TRIGGERS_PLUGIN: WorkflowPlugin = {
  id: 'core-triggers',
  name: 'Core Triggers',
  version: '1.0.0',
  category: 'triggers',
  icon: 'bolt',
  onActivate(context) {
    context.registerTranslations('en-US', enUS.coreTriggers);
    context.registerTranslations('pt-BR', ptBR.coreTriggers);
  },
  nodeTypes: [
    {
      type: START_NODE_TYPE,
      label: 'Start',
      icon: 'play_arrow',
      color: 'amber-8',
      description: 'Workflow entry point',
      inputs: [],
      outputs: [{ id: 'out', label: 'Out', position: 'bottom' }],
      configSchema: {},
      canvasComponent: markRaw(GenericWorkflowNode),
      defaults: {},
      deletable: false,
      catalogHidden: true,
      shape: 'circle',
    },
    {
      type: 'core/trigger_event',
      label: 'Trigger Event',
      icon: 'bolt',
      color: 'amber-8',
      description: 'Entry point — triggered by a configured trigger from the platform',
      inputs: [{ id: 'in', label: 'In', position: 'top' }],
      outputs: [
        { id: 'success', label: 'Success', position: 'bottom', color: '#4caf50' },
        { id: 'error', label: 'Error', position: 'bottom', color: '#ef5350' },
      ],
      timeout: { duration: 30, unit: 'seconds', enableOutput: false },
      errorHandler: { enabled: false, maxAttempts: 3, initialInterval: 5, intervalUnit: 'seconds', backoffMultiplier: 2.0 },
      configSchema: {},
      canvasComponent: markRaw(GenericWorkflowNode),
      configComponent: markRaw(TriggerEventNodeConfig),
      defaults: {},
      validate: validateTriggerEvent,
      availableOutputs: [
        { path: '*', description: 'Downstream triggered output (if any)' },
      ],
    },
  ],
};

/**
 * Core Logic plugin — conditions and evaluation node types
 */
export const CORE_LOGIC_PLUGIN: WorkflowPlugin = {
  id: 'core-logic',
  name: 'Core Logic',
  version: '1.0.0',
  category: 'logic',
  icon: 'rule',
  onActivate(context) {
    context.registerTranslations('en-US', enUS.coreLogic);
    context.registerTranslations('pt-BR', ptBR.coreLogic);
  },
  nodeTypes: [
    {
      type: 'core/condition',
      label: 'Conditions',
      icon: 'fact_check',
      color: 'blue-7',
      description: 'Evaluate conditions and branch based on result',
      inputs: [{ id: 'in', label: 'In', position: 'top' }],
      outputs: [
        { id: 'true', label: 'True', position: 'bottom', color: '#4caf50' },
        { id: 'false', label: 'False', position: 'bottom', color: '#ef5350' },
      ],
      configSchema: {},
      canvasComponent: markRaw(GenericWorkflowNode),
      configComponent: markRaw(ConditionNodeConfig),
      validate: validateCondition,
      defaults: {
        logic: 'AND',
        items: [],
        selectedTemplateIds: [],
      },
    },
  ],
};

/**
 * Core Data plugin — data manipulation and logging node types
 */
export const CORE_DATA_PLUGIN: WorkflowPlugin = {
  id: 'core-data',
  name: 'Core Data',
  version: '1.0.0',
  category: 'state',
  icon: 'edit_note',
  onActivate(context) {
    context.registerTranslations('en-US', enUS.coreData);
    context.registerTranslations('pt-BR', ptBR.coreData);
  },
  nodeTypes: [
    {
      type: 'core/log',
      label: 'Log',
      icon: 'article',
      color: 'teal-7',
      description: 'Emit observability event',
      inputs: [{ id: 'in', label: 'In', position: 'top' }],
      outputs: [{ id: 'out', label: 'Out', position: 'bottom' }],
      configSchema: {},
      canvasComponent: markRaw(GenericWorkflowNode),
      properties: [
        {
          name: 'message',
          displayName: 'Message',
          type: 'string',
          default: '',
        },
        {
          name: 'level',
          displayName: 'Level',
          type: 'options',
          default: 'info',
          options: [
            { label: 'Info', value: 'info' },
            { label: 'Warn', value: 'warn' },
            { label: 'Error', value: 'error' },
            { label: 'Debug', value: 'debug' },
          ],
        },
      ],
      defaults: { message: '', level: 'info' },
    },
    {
      type: 'core/set_state',
      label: 'Set State',
      icon: 'edit_note',
      color: 'teal-7',
      description: 'Set, increment, decrement or remove a state variable',
      inputs: [{ id: 'in', label: 'In', position: 'top' }],
      outputs: [{ id: 'out', label: 'Out', position: 'bottom' }],
      configSchema: {},
      canvasComponent: markRaw(GenericWorkflowNode),
      configComponent: markRaw(SetStateNodeConfig),
      validate: validateSetState,
      defaults: {
        operation: 'set',
        targetField: '',
        valueSource: { type: 'literal', value: '' },
        selectedTemplateIds: [],
      },
    },
    {
      type: 'core/code',
      label: 'Code',
      icon: 'code',
      color: 'teal-7',
      description: 'Execute a JavaScript snippet',
      inputs: [{ id: 'in', label: 'In', position: 'top' }],
      outputs: [
        { id: 'success', label: 'Success', position: 'bottom', color: '#4caf50' },
        { id: 'error', label: 'Error', position: 'bottom', color: '#ef5350' },
      ],
      timeout: { duration: 30, unit: 'seconds', enableOutput: false },
      errorHandler: { enabled: false, maxAttempts: 3, initialInterval: 5, intervalUnit: 'seconds', backoffMultiplier: 2.0 },
      configSchema: {},
      canvasComponent: markRaw(GenericWorkflowNode),
      configComponent: markRaw(CodeNodeConfig),
      validate: validateCode,
      defaults: {
        script: '// Access: state, event, inputs, nodes\n\nreturn {};',
        timeout: 5000,
      },
      availableOutputs: [
        { path: '*', description: 'Depends on return object of script' },
      ],
    },
  ],
};

/**
 * Core Flow Control plugin — branching, merging, subworkflow node types
 */
export const CORE_FLOW_CONTROL_PLUGIN: WorkflowPlugin = {
  id: 'core-flow-control',
  name: 'Core Flow Control',
  version: '1.0.0',
  category: 'flow_control',
  icon: 'device_hub',
  onActivate(context) {
    context.registerTranslations('en-US', enUS.coreFlowControl);
    context.registerTranslations('pt-BR', ptBR.coreFlowControl);
  },
  nodeTypes: [
    {
      type: 'core/fanout',
      label: 'Fanout',
      icon: 'account_tree',
      color: 'purple-7',
      description: 'Parallel execution fork',
      inputs: [{ id: 'in', label: 'In', position: 'top' }],
      outputs: [
        { id: 'out_1', label: 'Out 1', position: 'bottom' },
        { id: 'out_2', label: 'Out 2', position: 'bottom' },
      ],
      configSchema: {},
      canvasComponent: markRaw(GenericWorkflowNode),
      properties: [
        {
          name: 'branches',
          displayName: 'Number of Branches',
          type: 'number',
          default: 2,
        },
        {
          name: 'mode',
          displayName: 'Mode',
          type: 'options',
          default: 'waitAll',
          options: [
            { label: 'Wait All', value: 'waitAll' },
            { label: 'First Completed', value: 'firstCompleted' },
          ],
        },
      ],
      defaults: { branches: 2, mode: 'waitAll' },
      resolveOutputs: (config) => {
        const count = Math.max(1, (config.branches as number) || 2);
        return Array.from({ length: count }, (_, i) => ({
          id: `out_${i + 1}`,
          label: `Out ${i + 1}`,
          position: 'bottom' as const,
        }));
      },
    },
    {
      type: 'core/sequence',
      label: 'Sequence',
      icon: 'format_list_numbered',
      color: 'purple-7',
      description: 'Sequential execution — runs each step in order',
      inputs: [{ id: 'in', label: 'In', position: 'top' }],
      outputs: [
        { id: 'step_1', label: 'Step 1', position: 'bottom' },
        { id: 'step_2', label: 'Step 2', position: 'bottom' },
      ],
      configSchema: {},
      canvasComponent: markRaw(GenericWorkflowNode),
      properties: [
        {
          name: 'steps',
          displayName: 'Number of Steps',
          type: 'number',
          default: 2,
        },
      ],
      defaults: { steps: 2 },
      resolveOutputs: (config) => {
        const count = Math.max(1, (config.steps as number) || 2);
        return Array.from({ length: count }, (_, i) => ({
          id: `step_${i + 1}`,
          label: `Step ${i + 1}`,
          position: 'bottom' as const,
        }));
      },
    },
    {
      type: 'core/merge',
      label: 'Merge',
      icon: 'call_merge',
      color: 'purple-7',
      description: 'Join parallel branches',
      inputs: [
        { id: 'in_1', label: 'In 1', position: 'top' },
        { id: 'in_2', label: 'In 2', position: 'top' },
      ],
      outputs: [{ id: 'out', label: 'Out', position: 'bottom' }],
      configSchema: {},
      canvasComponent: markRaw(GenericWorkflowNode),
      properties: [
        {
          name: 'branches',
          displayName: 'Number of Branches',
          type: 'number',
          default: 2,
        },
        {
          name: 'strategy',
          displayName: 'Merge Strategy',
          type: 'options',
          default: 'all',
          options: [
            { label: 'All', value: 'all' },
            { label: 'Any', value: 'any' },
            { label: 'First', value: 'first' },
          ],
        },
      ],
      defaults: { branches: 2, strategy: 'all' },
      resolveInputs: (config) => {
        const count = Math.max(1, (config.branches as number) || 2);
        return Array.from({ length: count }, (_, i) => ({
          id: `in_${i + 1}`,
          label: `In ${i + 1}`,
          position: 'top' as const,
        }));
      },
    },
    {
      type: 'core/switch',
      label: 'Switch',
      icon: 'alt_route',
      color: 'purple-7',
      description: 'Route to paths based on matching conditions',
      inputs: [{ id: 'in', label: 'In', position: 'top' }],
      outputs: [{ id: 'default', label: 'Default', position: 'bottom' }],
      configSchema: {},
      canvasComponent: markRaw(GenericWorkflowNode),
      configComponent: markRaw(SwitchNodeConfig),
      validate: validateSwitch,
      defaults: { cases: [], matchMode: 'first', selectedTemplateIds: [] },
      resolveOutputs: (config) => {
        const cases = (config.cases as SwitchCase[]) || [];
        const outputs = cases.map((c, i) => ({
          id: c.id,
          label: `Case ${i + 1}`,
          position: 'bottom' as const,
        }));
        outputs.push({ id: 'default', label: 'Default', position: 'bottom' });
        return outputs;
      },
    },
    {
      type: 'core/subworkflow',
      label: 'Subworkflow',
      icon: 'hub',
      color: 'deep-purple-6',
      description: 'Execute another workflow as a child process',
      inputs: [{ id: 'in', label: 'In', position: 'top' }],
      outputs: [
        { id: 'success', label: 'Success', position: 'bottom', color: '#4caf50' },
        { id: 'error', label: 'Error', position: 'bottom', color: '#ef5350' },
      ],
      timeout: { duration: 3600, unit: 'seconds', enableOutput: false },
      errorHandler: { enabled: false, maxAttempts: 3, initialInterval: 5, intervalUnit: 'seconds', backoffMultiplier: 2.0 },
      configSchema: {},
      canvasComponent: markRaw(GenericWorkflowNode),
      configComponent: markRaw(SubworkflowNodeConfig),
      validate: validateSubworkflow,
      defaults: {
        workflowId: undefined,
        workflowName: undefined,
        executionMode: 'sync',
        timeout: { duration: 30, unit: 'seconds' },
        inputMappings: [],
        outputMappings: [],
      },
      availableOutputs: [
        { path: '*', description: 'Depends on child workflow outputs' },
      ],
    },
    {
      type: 'core/loop',
      label: 'Loop',
      icon: 'loop',
      color: 'purple-7',
      description: 'Iterate over a list',
      inputs: [{ id: 'in', label: 'In', position: 'top' }],
      outputs: [
        { id: 'body', label: 'Body', position: 'bottom', color: '#2196f3' },
        { id: 'done', label: 'Done', position: 'bottom', color: '#4caf50' },
        { id: 'error', label: 'Error', position: 'bottom', color: '#ef5350' },
      ],
      errorHandler: { enabled: false, maxAttempts: 3, initialInterval: 5, intervalUnit: 'seconds', backoffMultiplier: 2.0 },
      configSchema: {},
      canvasComponent: markRaw(GenericWorkflowNode),
      configComponent: markRaw(LoopNodeConfig),
      validate: validateLoop,
      defaults: {
        source: { type: 'state', value: '' },
      },
      availableOutputs: [
        { path: 'item', description: 'Current iteration element' },
        { path: 'index', description: '0-based position' },
      ],
    },
    {
      type: 'core/end',
      label: 'End',
      icon: 'check_circle',
      color: 'purple-7',
      description: 'Workflow termination point with optional error mode',
      inputs: [{ id: 'in', label: 'In', position: 'top' }],
      outputs: [],
      configSchema: {},
      canvasComponent: markRaw(EndNode),
      configComponent: markRaw(EndNodeConfig),
      validate: validateEnd,
      defaults: {
        terminateWithError: false,
        errorCode: '',
        errorMessage: { type: 'literal', value: '' },
      },
      shape: 'circle',
    },
    {
      type: 'core/goto',
      label: 'Goto',
      icon: 'place',
      color: 'deep-purple-6',
      description: 'Virtual portal — connect workflow sections without edges',
      inputs: [{ id: 'in', label: 'In', position: 'top' }],
      outputs: [{ id: 'out', label: 'Out', position: 'bottom' }],
      configSchema: {},
      canvasComponent: markRaw(GotoNode),
      configComponent: markRaw(GotoNodeConfig),
      validate: validateGoto,
      defaults: { role: 'sender', pairLabel: '', pairColor: 'deep-purple-6' },
      resolveInputs: (config) =>
        config.role === 'sender'
          ? [{ id: 'in', label: 'In', position: 'top' as const }]
          : [],
      resolveOutputs: (config) =>
        config.role === 'receiver'
          ? [{ id: 'out', label: 'Out', position: 'bottom' as const }]
          : [],
    },
  ],
};

/**
 * Core Timers plugin — delay, wait signal, and wait for node types
 */
export const CORE_TIMERS_PLUGIN: WorkflowPlugin = {
  id: 'core-timers',
  name: 'Core Timers',
  version: '1.0.0',
  category: 'timers',
  icon: 'hourglass_empty',
  onActivate(context) {
    context.registerTranslations('en-US', enUS.coreTimers);
    context.registerTranslations('pt-BR', ptBR.coreTimers);
  },
  nodeTypes: [
    {
      type: 'core/delay',
      label: 'Delay',
      icon: 'hourglass_empty',
      color: 'orange-7',
      description: 'Fixed wait (seconds to years)',
      inputs: [{ id: 'in', label: 'In', position: 'top' }],
      outputs: [{ id: 'out', label: 'Out', position: 'bottom' }],
      configSchema: {},
      canvasComponent: markRaw(GenericWorkflowNode),
      properties: [
        {
          name: 'duration',
          displayName: 'Duration',
          type: 'number',
          default: 30,
        },
        {
          name: 'unit',
          displayName: 'Unit',
          type: 'options',
          default: 'seconds',
          options: [
            { label: 'Seconds', value: 'seconds' },
            { label: 'Minutes', value: 'minutes' },
            { label: 'Hours', value: 'hours' },
            { label: 'Days', value: 'days' },
            { label: 'Months (30 days)', value: 'months' },
            { label: 'Years (365 days)', value: 'years' },
          ],
        },
      ],
      defaults: { duration: 30, unit: 'seconds' },
    },
    {
      type: 'core/wait_signal',
      label: 'Wait Signal',
      icon: 'notifications_active',
      color: 'orange-7',
      description: 'Wait for external signal with timeout',
      inputs: [{ id: 'in', label: 'In', position: 'top' }],
      outputs: [
        { id: 'received', label: 'Received', position: 'bottom', color: '#4caf50' },
      ],
      timeout: { duration: 86400, unit: 'seconds', enableOutput: false },
      configSchema: {},
      canvasComponent: markRaw(GenericWorkflowNode),
      configComponent: markRaw(WaitSignalNodeConfig),
      validate: validateWaitSignal,
      defaults: { signalName: '', mappings: [] },
      availableOutputs: [
        { path: '*', description: 'Signal payload received from publisher' },
      ],
    },
    {
      type: 'core/wait_for',
      label: 'Wait For',
      icon: 'pending_actions',
      color: 'orange-7',
      description: 'Wait until a state variable condition is met',
      inputs: [{ id: 'in', label: 'In', position: 'top' }],
      outputs: [
        { id: 'matched', label: 'Matched', position: 'bottom', color: '#4caf50' },
      ],
      timeout: { duration: 86400, unit: 'seconds', enableOutput: false },
      configSchema: {},
      canvasComponent: markRaw(GenericWorkflowNode),
      configComponent: markRaw(WaitForNodeConfig),
      validate: validateWaitFor,
      defaults: {
        field: '',
        operator: 'equals',
        compareTo: { source: 'literal', value: '' },
        interval: '30s',
      },
      availableOutputs: [
        { path: '*', description: 'Resolved condition payload (state snapshot when matched)' },
      ],
    },
  ],
};

/**
 * Core Annotations plugin — non-functional visual elements for documenting workflows
 */
export const CORE_ANNOTATIONS_PLUGIN: WorkflowPlugin = {
  id: 'core-annotations',
  name: 'Core Annotations',
  version: '1.0.0',
  category: 'annotations',
  icon: 'sticky_note_2',
  onActivate(context) {
    context.registerTranslations('en-US', enUS.coreAnnotations);
    context.registerTranslations('pt-BR', ptBR.coreAnnotations);
  },
  nodeTypes: [
    {
      type: 'core/text_note',
      label: 'Text Note',
      icon: 'sticky_note_2',
      color: 'amber-7',
      description: 'Add a text note to document your workflow',
      inputs: [],
      outputs: [],
      configSchema: {},
      canvasComponent: markRaw(TextNoteNode),
      defaults: { text: '', color: 'grey' },
      catalogHidden: true,
      configurable: false,
    },
    {
      type: 'core/group_frame',
      label: 'Group Frame',
      icon: 'dashboard',
      color: 'blue-grey-6',
      description: 'Visual container to organize workflow sections',
      inputs: [],
      outputs: [],
      configSchema: {},
      canvasComponent: markRaw(GroupFrameNode),
      configComponent: markRaw(GroupFrameNodeConfig),
      defaults: {
        title: '',
        description: '',
        color: 'blue-grey',
        width: 300,
        height: 200,
      },
      configurable: true,
    },
  ],
};

/**
 * All core plugins grouped for registration
 */
export const CORE_PLUGINS: WorkflowPlugin[] = [
  CORE_TRIGGERS_PLUGIN,
  CORE_LOGIC_PLUGIN,
  CORE_DATA_PLUGIN,
  CORE_FLOW_CONTROL_PLUGIN,
  CORE_TIMERS_PLUGIN,
  CORE_ANNOTATIONS_PLUGIN,
];

/**
 * Register all core plugins in the plugin registry
 *
 * @param {(plugin: WorkflowPlugin) => void} registerFn - Function to register a plugin
 */
export function bootWorkflowPlugins(registerFn: (plugin: WorkflowPlugin) => void): void {
  CORE_PLUGINS.forEach(registerFn);

  // Register shared workflow translations (used by BaseWorkflowNode, execution viewer, etc.)
  if (i18nInstance) {
    i18nInstance.global.mergeLocaleMessage('en-US', { wf: { common: enUS.shared } });
    i18nInstance.global.mergeLocaleMessage('pt-BR', { wf: { common: ptBR.shared } });
  }
}
