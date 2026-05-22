import { describe, it, expect, vi } from 'vitest';
import { mountWithPlugins } from '@src/test/helpers';
import SubworkflowNodeConfig from './SubworkflowNodeConfig.vue';

vi.mock('@src/composables/workflow', () => ({
  useWorkflowContext: () => ({
    states: { value: [{ field: 'status', type: 'string' }, { field: 'result', type: 'any' }] },
  }),
  usePluginI18n: () => ({ t: (key: string) => key }),
}));

vi.mock('@components/tooltips', () => ({
  AppTooltip: { name: 'AppTooltip', template: '<div />' },
}));

vi.mock('@components/dialogs/common/workflowSelectorDialog', () => ({
  WorkflowSelectorDialog: { name: 'WorkflowSelectorDialog', template: '<div />' },
}));

vi.mock('@components/forms/fieldSourceSelector', () => ({
  FieldSourceSelector: { name: 'FieldSourceSelector', template: '<div />' },
}));

const BASE_CONFIG: Record<string, unknown> = {};

describe('SubworkflowNodeConfig', () => {
  // ── Initial State ──────────────────────────────────────────────────

  it('renders without errors', () => {
    const wrapper = mountWithPlugins(SubworkflowNodeConfig, { props: { config: BASE_CONFIG } });
    expect(wrapper.exists()).toBe(true);
  });

  it('starts with workflowDialogOpen as false', () => {
    const wrapper = mountWithPlugins(SubworkflowNodeConfig, { props: { config: BASE_CONFIG } });
    expect(wrapper.vm.workflowDialogOpen).toBe(false);
  });

  // ── Computed: executionMode ────────────────────────────────────────

  it('defaults executionMode to sync', () => {
    const wrapper = mountWithPlugins(SubworkflowNodeConfig, { props: { config: BASE_CONFIG } });
    expect(wrapper.vm.executionMode).toBe('sync');
  });

  it('reads executionMode from config', () => {
    const wrapper = mountWithPlugins(SubworkflowNodeConfig, {
      props: { config: { executionMode: 'async' } },
    });
    expect(wrapper.vm.executionMode).toBe('async');
  });

  // ── Computed: timeout ──────────────────────────────────────────────

  it('defaults timeout to 30 seconds', () => {
    const wrapper = mountWithPlugins(SubworkflowNodeConfig, { props: { config: BASE_CONFIG } });
    expect(wrapper.vm.timeout).toEqual({ duration: 30, unit: 'seconds' });
  });

  it('reads timeout from config', () => {
    const wrapper = mountWithPlugins(SubworkflowNodeConfig, {
      props: { config: { timeout: { duration: 60, unit: 'minutes' } } },
    });
    expect(wrapper.vm.timeout).toEqual({ duration: 60, unit: 'minutes' });
  });

  // ── Computed: inputMappings ────────────────────────────────────────

  it('defaults inputMappings to empty array', () => {
    const wrapper = mountWithPlugins(SubworkflowNodeConfig, { props: { config: BASE_CONFIG } });
    expect(wrapper.vm.inputMappings).toEqual([]);
  });

  it('reads inputMappings from config', () => {
    const mappings = [{ childVariable: 'x', source: { type: 'literal', value: '42' } }];
    const wrapper = mountWithPlugins(SubworkflowNodeConfig, {
      props: { config: { inputMappings: mappings } },
    });
    expect(wrapper.vm.inputMappings).toEqual(mappings);
  });

  // ── Computed: outputMappings ───────────────────────────────────────

  it('defaults outputMappings to empty array', () => {
    const wrapper = mountWithPlugins(SubworkflowNodeConfig, { props: { config: BASE_CONFIG } });
    expect(wrapper.vm.outputMappings).toEqual([]);
  });

  it('reads outputMappings from config', () => {
    const mappings = [{ outputKey: 'result', targetVariable: 'status' }];
    const wrapper = mountWithPlugins(SubworkflowNodeConfig, {
      props: { config: { outputMappings: mappings } },
    });
    expect(wrapper.vm.outputMappings).toEqual(mappings);
  });

  // ── Computed: variableOptions ──────────────────────────────────────

  it('computes variableOptions from workflow context states', () => {
    const wrapper = mountWithPlugins(SubworkflowNodeConfig, { props: { config: BASE_CONFIG } });
    expect(wrapper.vm.variableOptions).toEqual([
      { label: 'status', value: 'status' },
      { label: 'result', value: 'result' },
    ]);
  });

  // ── Computed: showOutputMappings ───────────────────────────────────

  it('hides output mappings when no workflow is selected', () => {
    const wrapper = mountWithPlugins(SubworkflowNodeConfig, { props: { config: BASE_CONFIG } });
    expect(wrapper.vm.showOutputMappings).toBe(false);
  });

  it('shows output mappings when workflow is selected and mode is sync', () => {
    const wrapper = mountWithPlugins(SubworkflowNodeConfig, {
      props: { config: { workflowId: 'wf-1', executionMode: 'sync' } },
    });
    expect(wrapper.vm.showOutputMappings).toBe(true);
  });

  it('hides output mappings when mode is async', () => {
    const wrapper = mountWithPlugins(SubworkflowNodeConfig, {
      props: { config: { workflowId: 'wf-1', executionMode: 'async' } },
    });
    expect(wrapper.vm.showOutputMappings).toBe(false);
  });

  // ── Computed: stateFields ──────────────────────────────────────────

  it('computes stateFields from workflow context states', () => {
    const wrapper = mountWithPlugins(SubworkflowNodeConfig, { props: { config: BASE_CONFIG } });
    expect(wrapper.vm.stateFields).toEqual([
      { name: 'status', type: 'string' },
      { name: 'result', type: 'any' },
    ]);
  });

  // ── handleWorkflowSelect ──────────────────────────────────────────

  it('emits update:config with workflow info on handleWorkflowSelect', () => {
    const wrapper = mountWithPlugins(SubworkflowNodeConfig, { props: { config: BASE_CONFIG } });
    wrapper.vm.handleWorkflowSelect({ id: 'wf-1', name: 'My Workflow' });
    const emitted = wrapper.emitted('update:config')!;
    expect(emitted[0]![0]).toMatchObject({
      workflowId: 'wf-1',
      workflowName: 'My Workflow',
    });
  });

  // ── clearWorkflow ──────────────────────────────────────────────────

  it('clears workflow and resets mappings on clearWorkflow', () => {
    const wrapper = mountWithPlugins(SubworkflowNodeConfig, {
      props: { config: { workflowId: 'wf-1', workflowName: 'W', inputMappings: [{}], outputMappings: [{}] } },
    });
    wrapper.vm.clearWorkflow();
    const emitted = wrapper.emitted('update:config')!;
    expect(emitted[0]![0]).toMatchObject({
      workflowId: undefined,
      workflowName: undefined,
      inputMappings: [],
      outputMappings: [],
    });
  });

  // ── updateExecutionMode ────────────────────────────────────────────

  it('emits update:config with new execution mode', () => {
    const wrapper = mountWithPlugins(SubworkflowNodeConfig, { props: { config: BASE_CONFIG } });
    wrapper.vm.updateExecutionMode('async');
    const emitted = wrapper.emitted('update:config')!;
    expect((emitted[0]![0] as any).executionMode).toBe('async');
  });

  it('clears output mappings when switching to async', () => {
    const config = { executionMode: 'sync', outputMappings: [{ outputKey: 'r', targetVariable: 's' }] };
    const wrapper = mountWithPlugins(SubworkflowNodeConfig, { props: { config } });
    wrapper.vm.updateExecutionMode('async');
    const emitted = wrapper.emitted('update:config')!;
    expect((emitted[0]![0] as any).outputMappings).toEqual([]);
  });

  it('preserves output mappings when switching to sync', () => {
    const config = { executionMode: 'async', outputMappings: [{ outputKey: 'r', targetVariable: 's' }] };
    const wrapper = mountWithPlugins(SubworkflowNodeConfig, { props: { config } });
    wrapper.vm.updateExecutionMode('sync');
    const emitted = wrapper.emitted('update:config')!;
    expect((emitted[0]![0] as any).outputMappings).toEqual([{ outputKey: 'r', targetVariable: 's' }]);
  });

  // ── updateTimeoutDuration ──────────────────────────────────────────

  it('emits update:config with new timeout duration', () => {
    const wrapper = mountWithPlugins(SubworkflowNodeConfig, { props: { config: BASE_CONFIG } });
    wrapper.vm.updateTimeoutDuration(60);
    const emitted = wrapper.emitted('update:config')!;
    expect((emitted[0]![0] as any).timeout.duration).toBe(60);
  });

  it('clamps timeout duration to minimum 1 for negative values', () => {
    const wrapper = mountWithPlugins(SubworkflowNodeConfig, { props: { config: BASE_CONFIG } });
    wrapper.vm.updateTimeoutDuration(-5);
    const emitted = wrapper.emitted('update:config')!;
    expect((emitted[0]![0] as any).timeout.duration).toBe(1);
  });

  it('falls back to 30 for zero (falsy) due to || guard', () => {
    const wrapper = mountWithPlugins(SubworkflowNodeConfig, { props: { config: BASE_CONFIG } });
    wrapper.vm.updateTimeoutDuration(0);
    const emitted = wrapper.emitted('update:config')!;
    expect((emitted[0]![0] as any).timeout.duration).toBe(30);
  });

  it('defaults to 30 for null timeout duration', () => {
    const wrapper = mountWithPlugins(SubworkflowNodeConfig, { props: { config: BASE_CONFIG } });
    wrapper.vm.updateTimeoutDuration(null);
    const emitted = wrapper.emitted('update:config')!;
    expect((emitted[0]![0] as any).timeout.duration).toBe(30);
  });

  it('defaults to 30 for NaN timeout duration', () => {
    const wrapper = mountWithPlugins(SubworkflowNodeConfig, { props: { config: BASE_CONFIG } });
    wrapper.vm.updateTimeoutDuration('abc');
    const emitted = wrapper.emitted('update:config')!;
    expect((emitted[0]![0] as any).timeout.duration).toBe(30);
  });

  // ── updateTimeoutUnit ──────────────────────────────────────────────

  it('emits update:config with new timeout unit', () => {
    const wrapper = mountWithPlugins(SubworkflowNodeConfig, { props: { config: BASE_CONFIG } });
    wrapper.vm.updateTimeoutUnit('minutes');
    const emitted = wrapper.emitted('update:config')!;
    expect((emitted[0]![0] as any).timeout.unit).toBe('minutes');
  });

  it('preserves timeout duration when changing unit', () => {
    const wrapper = mountWithPlugins(SubworkflowNodeConfig, {
      props: { config: { timeout: { duration: 45, unit: 'seconds' } } },
    });
    wrapper.vm.updateTimeoutUnit('hours');
    const emitted = wrapper.emitted('update:config')!;
    expect((emitted[0]![0] as any).timeout).toEqual({ duration: 45, unit: 'hours' });
  });

  // ── Input Mapping helpers ──────────────────────────────────────────

  it('adds empty input mapping on addInputMapping', () => {
    const wrapper = mountWithPlugins(SubworkflowNodeConfig, { props: { config: BASE_CONFIG } });
    wrapper.vm.addInputMapping();
    const emitted = wrapper.emitted('update:config')!;
    const mappings = (emitted[0]![0] as any).inputMappings;
    expect(mappings).toHaveLength(1);
    expect(mappings[0]).toEqual({ childVariable: '', source: { type: 'literal', value: '' } });
  });

  it('appends to existing input mappings on addInputMapping', () => {
    const existing = [{ childVariable: 'x', source: { type: 'literal', value: '1' } }];
    const wrapper = mountWithPlugins(SubworkflowNodeConfig, {
      props: { config: { inputMappings: existing } },
    });
    wrapper.vm.addInputMapping();
    const emitted = wrapper.emitted('update:config')!;
    const mappings = (emitted[0]![0] as any).inputMappings;
    expect(mappings).toHaveLength(2);
    expect(mappings[0].childVariable).toBe('x');
  });

  it('removes input mapping by index on removeInputMapping', () => {
    const existing = [
      { childVariable: 'a', source: { type: 'literal', value: '1' } },
      { childVariable: 'b', source: { type: 'literal', value: '2' } },
    ];
    const wrapper = mountWithPlugins(SubworkflowNodeConfig, {
      props: { config: { inputMappings: existing } },
    });
    wrapper.vm.removeInputMapping(0);
    const emitted = wrapper.emitted('update:config')!;
    const mappings = (emitted[0]![0] as any).inputMappings;
    expect(mappings).toHaveLength(1);
    expect(mappings[0].childVariable).toBe('b');
  });

  it('updates input child variable on updateInputChildVariable', () => {
    const existing = [{ childVariable: '', source: { type: 'literal', value: '' } }];
    const wrapper = mountWithPlugins(SubworkflowNodeConfig, {
      props: { config: { inputMappings: existing } },
    });
    wrapper.vm.updateInputChildVariable(0, 'myVar');
    const emitted = wrapper.emitted('update:config')!;
    const mappings = (emitted[0]![0] as any).inputMappings;
    expect(mappings[0].childVariable).toBe('myVar');
  });

  it('does nothing on updateInputChildVariable with invalid index', () => {
    const wrapper = mountWithPlugins(SubworkflowNodeConfig, { props: { config: BASE_CONFIG } });
    wrapper.vm.updateInputChildVariable(5, 'myVar');
    expect(wrapper.emitted('update:config')).toBeUndefined();
  });

  it('updates input source on updateInputSource', () => {
    const existing = [{ childVariable: 'x', source: { type: 'literal', value: '' } }];
    const wrapper = mountWithPlugins(SubworkflowNodeConfig, {
      props: { config: { inputMappings: existing } },
    });
    wrapper.vm.updateInputSource(0, { type: 'state', value: 'status' });
    const emitted = wrapper.emitted('update:config')!;
    const mappings = (emitted[0]![0] as any).inputMappings;
    expect(mappings[0].source).toEqual({ type: 'state', value: 'status' });
  });

  it('does nothing on updateInputSource with invalid index', () => {
    const wrapper = mountWithPlugins(SubworkflowNodeConfig, { props: { config: BASE_CONFIG } });
    wrapper.vm.updateInputSource(10, { type: 'literal', value: 'x' });
    expect(wrapper.emitted('update:config')).toBeUndefined();
  });

  // ── Output Mapping helpers ─────────────────────────────────────────

  it('adds empty output mapping on addOutputMapping', () => {
    const wrapper = mountWithPlugins(SubworkflowNodeConfig, { props: { config: BASE_CONFIG } });
    wrapper.vm.addOutputMapping();
    const emitted = wrapper.emitted('update:config')!;
    const mappings = (emitted[0]![0] as any).outputMappings;
    expect(mappings).toHaveLength(1);
    expect(mappings[0]).toEqual({ outputKey: '', targetVariable: '' });
  });

  it('removes output mapping by index on removeOutputMapping', () => {
    const existing = [
      { outputKey: 'a', targetVariable: 'status' },
      { outputKey: 'b', targetVariable: 'result' },
    ];
    const wrapper = mountWithPlugins(SubworkflowNodeConfig, {
      props: { config: { outputMappings: existing } },
    });
    wrapper.vm.removeOutputMapping(0);
    const emitted = wrapper.emitted('update:config')!;
    const mappings = (emitted[0]![0] as any).outputMappings;
    expect(mappings).toHaveLength(1);
    expect(mappings[0].outputKey).toBe('b');
  });

  it('updates output mapping outputKey on updateOutputMapping', () => {
    const existing = [{ outputKey: '', targetVariable: '' }];
    const wrapper = mountWithPlugins(SubworkflowNodeConfig, {
      props: { config: { outputMappings: existing } },
    });
    wrapper.vm.updateOutputMapping(0, 'outputKey', 'resultData');
    const emitted = wrapper.emitted('update:config')!;
    const mappings = (emitted[0]![0] as any).outputMappings;
    expect(mappings[0].outputKey).toBe('resultData');
  });

  it('updates output mapping targetVariable on updateOutputMapping', () => {
    const existing = [{ outputKey: 'r', targetVariable: '' }];
    const wrapper = mountWithPlugins(SubworkflowNodeConfig, {
      props: { config: { outputMappings: existing } },
    });
    wrapper.vm.updateOutputMapping(0, 'targetVariable', 'status');
    const emitted = wrapper.emitted('update:config')!;
    const mappings = (emitted[0]![0] as any).outputMappings;
    expect(mappings[0].targetVariable).toBe('status');
  });

  it('does nothing on updateOutputMapping with invalid index', () => {
    const wrapper = mountWithPlugins(SubworkflowNodeConfig, { props: { config: BASE_CONFIG } });
    wrapper.vm.updateOutputMapping(99, 'outputKey', 'test');
    expect(wrapper.emitted('update:config')).toBeUndefined();
  });
});
