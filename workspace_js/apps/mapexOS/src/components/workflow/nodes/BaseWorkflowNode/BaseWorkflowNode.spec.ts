import { describe, it, expect, vi } from 'vitest';
import { mountWithPlugins } from '@src/test/helpers';
import BaseWorkflowNode from './BaseWorkflowNode.vue';

vi.mock('@src/composables/workflow', () => ({
  useWorkflowContext: () => ({
    updateNodeConfig: vi.fn(),
    nodes: { value: [] },
    addNoteToNode: vi.fn(),
    pushSnapshot: vi.fn(),
  }),
}));

const BASE_PROPS = {
  id: 'node-1',
  icon: 'play_arrow',
  color: 'green-7',
  label: 'Start',
  selected: false,
  hasErrors: false,
  inputs: [],
  outputs: [],
};

describe('BaseWorkflowNode', () => {
  it('renders without errors', () => {
    const wrapper = mountWithPlugins(BaseWorkflowNode, { props: BASE_PROPS });
    expect(wrapper.exists()).toBe(true);
  });

  it('renders label from props', () => {
    const wrapper = mountWithPlugins(BaseWorkflowNode, {
      props: { ...BASE_PROPS, label: 'My Node' },
    });
    expect(wrapper.find('.wf-node__label').text()).toBe('My Node');
  });

  it('applies selected class when selected is true', () => {
    const wrapper = mountWithPlugins(BaseWorkflowNode, {
      props: { ...BASE_PROPS, selected: true },
    });
    expect(wrapper.find('.wf-node').classes()).toContain('wf-node--selected');
  });

  it('applies circle class when shape is circle', () => {
    const wrapper = mountWithPlugins(BaseWorkflowNode, {
      props: { ...BASE_PROPS, shape: 'circle' },
    });
    expect(wrapper.find('.wf-node').classes()).toContain('wf-node--circle');
  });

  it('applies error class when hasErrors is true', () => {
    const wrapper = mountWithPlugins(BaseWorkflowNode, {
      props: { ...BASE_PROPS, hasErrors: true },
    });
    expect(wrapper.find('.wf-node').classes()).toContain('wf-node--error');
  });

  it('does not apply error class when hasErrors is false', () => {
    const wrapper = mountWithPlugins(BaseWorkflowNode, {
      props: { ...BASE_PROPS, hasErrors: false },
    });
    expect(wrapper.find('.wf-node').classes()).not.toContain('wf-node--error');
  });

  it('computes rootStyle with --node-color when colorHex is provided', () => {
    const wrapper = mountWithPlugins(BaseWorkflowNode, {
      props: { ...BASE_PROPS, colorHex: '#FF0000' },
    });
    expect(wrapper.vm.rootStyle).toHaveProperty('--node-color', '#FF0000');
  });

  it('computes rootStyle without --node-color when colorHex is not provided', () => {
    const wrapper = mountWithPlugins(BaseWorkflowNode, {
      props: BASE_PROPS,
    });
    expect(wrapper.vm.rootStyle).not.toHaveProperty('--node-color');
  });

  it('computes maxHandleCount as max of inputs and outputs', () => {
    const wrapper = mountWithPlugins(BaseWorkflowNode, {
      props: {
        ...BASE_PROPS,
        inputs: [{ id: 'in-1', position: 'top', label: 'In' }],
        outputs: [
          { id: 'out-1', position: 'bottom', label: 'Out1' },
          { id: 'out-2', position: 'bottom', label: 'Out2' },
        ],
      },
    });
    expect(wrapper.vm.maxHandleCount).toBe(2);
  });

  it('computes nodeMinWidth as 44px for single handle', () => {
    const wrapper = mountWithPlugins(BaseWorkflowNode, {
      props: {
        ...BASE_PROPS,
        inputs: [{ id: 'in-1', position: 'top', label: 'In' }],
        outputs: [],
      },
    });
    expect(wrapper.vm.nodeMinWidth).toBe('44px');
  });

  it('computes nodeMinWidth based on handle count for multiple handles', () => {
    const wrapper = mountWithPlugins(BaseWorkflowNode, {
      props: {
        ...BASE_PROPS,
        outputs: [
          { id: 'out-1', position: 'bottom', label: 'Out1' },
          { id: 'out-2', position: 'bottom', label: 'Out2' },
          { id: 'out-3', position: 'bottom', label: 'Out3' },
        ],
      },
    });
    expect(wrapper.vm.nodeMinWidth).toBe('150px');
  });

  it('renders Handle stubs for each input', () => {
    const wrapper = mountWithPlugins(BaseWorkflowNode, {
      props: {
        ...BASE_PROPS,
        inputs: [
          { id: 'in-1', position: 'top', label: 'Input 1' },
          { id: 'in-2', position: 'top', label: 'Input 2' },
        ],
      },
    });
    const handles = wrapper.findAll('handle-stub');
    expect(handles.length).toBeGreaterThanOrEqual(2);
  });

  it('renders Handle stubs for each output', () => {
    const wrapper = mountWithPlugins(BaseWorkflowNode, {
      props: {
        ...BASE_PROPS,
        outputs: [
          { id: 'out-1', position: 'bottom', label: 'Out 1' },
        ],
      },
    });
    const handles = wrapper.findAll('handle-stub');
    expect(handles.length).toBeGreaterThanOrEqual(1);
  });
});
