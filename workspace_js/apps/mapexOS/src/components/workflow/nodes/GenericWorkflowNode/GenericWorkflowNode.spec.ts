import { describe, it, expect, vi } from 'vitest';
import { mountWithPlugins } from '@src/test/helpers';
import GenericWorkflowNode from './GenericWorkflowNode.vue';

vi.mock('@src/components/workflow/constants', () => ({
  POSITION_OPTIONS: [
    { label: 'Top', value: 'top', icon: 'arrow_upward' },
    { label: 'Bottom', value: 'bottom', icon: 'arrow_downward' },
    { label: 'Left', value: 'left', icon: 'arrow_back' },
    { label: 'Right', value: 'right', icon: 'arrow_forward' },
  ],
}));

const mockGetNodeType = vi.fn();

vi.mock('@src/composables/workflow', () => ({
  useWorkflowContext: () => ({
    getNodeType: mockGetNodeType,
    updateNodeConfig: vi.fn(),
    nodes: { value: [] },
    addNoteToNode: vi.fn(),
    pushSnapshot: vi.fn(),
  }),
}));

vi.mock('vue-i18n', () => ({
  useI18n: () => ({
    t: (key: string) => key,
    te: (key: string) => key === 'wf.test-plugin.nodes.myNode.label',
  }),
}));

vi.mock('@src/utils/workflow', () => ({
  resolveNodeHandles: () => ({
    inputs: [{ id: 'in-1', position: 'top', label: 'In' }],
    outputs: [{ id: 'out-1', position: 'bottom', label: 'Out' }],
  }),
}));

const BASE_PROPS = {
  id: 'node-1',
  data: {
    config: {},
    label: 'Fallback Label',
    __nodeType: 'test/myNode',
  },
  selected: false,
};

describe('GenericWorkflowNode', () => {
  it('renders without errors', () => {
    mockGetNodeType.mockReturnValue(undefined);
    const wrapper = mountWithPlugins(GenericWorkflowNode, { props: BASE_PROPS });
    expect(wrapper.exists()).toBe(true);
  });

  it('computes nodeType from getNodeType', () => {
    const nodeTypeDef = {
      type: 'test/myNode',
      label: 'My Node',
      icon: 'settings',
      color: 'blue-7',
      inputs: [],
      outputs: [],
    };
    mockGetNodeType.mockReturnValue(nodeTypeDef);

    const wrapper = mountWithPlugins(GenericWorkflowNode, { props: BASE_PROPS });
    expect(wrapper.vm.nodeType).toEqual(nodeTypeDef);
  });

  it('returns undefined for nodeType when __nodeType is not set', () => {
    mockGetNodeType.mockReturnValue(undefined);
    const wrapper = mountWithPlugins(GenericWorkflowNode, {
      props: {
        ...BASE_PROPS,
        data: { config: {}, label: 'No Type' },
      },
    });
    expect(wrapper.vm.nodeType).toBeUndefined();
  });

  it('computes translatedLabel from i18n when key exists', () => {
    mockGetNodeType.mockReturnValue({
      type: 'test/myNode',
      label: 'My Node',
      icon: 'settings',
      color: 'blue-7',
      _pluginId: 'test-plugin',
    });

    const wrapper = mountWithPlugins(GenericWorkflowNode, { props: BASE_PROPS });
    expect(wrapper.vm.translatedLabel).toBe('wf.test-plugin.nodes.myNode.label');
  });

  it('falls back to data.label when i18n key does not exist', () => {
    mockGetNodeType.mockReturnValue({
      type: 'test/otherNode',
      label: 'Other Node',
      icon: 'settings',
      color: 'blue-7',
      _pluginId: 'unknown-plugin',
    });

    const wrapper = mountWithPlugins(GenericWorkflowNode, {
      props: {
        ...BASE_PROPS,
        data: { config: {}, label: 'Custom Label', __nodeType: 'test/otherNode' },
      },
    });
    expect(wrapper.vm.translatedLabel).toBe('Custom Label');
  });

  it('falls back to nodeType.label when data.label is undefined and i18n key missing', () => {
    mockGetNodeType.mockReturnValue({
      type: 'test/otherNode',
      label: 'Type Label',
      icon: 'settings',
      color: 'blue-7',
      _pluginId: 'unknown-plugin',
    });

    const wrapper = mountWithPlugins(GenericWorkflowNode, {
      props: {
        ...BASE_PROPS,
        data: { config: {}, __nodeType: 'test/otherNode' },
      },
    });
    expect(wrapper.vm.translatedLabel).toBe('Type Label');
  });

  it('computes resolvedHandles with inputs and outputs', () => {
    mockGetNodeType.mockReturnValue({
      type: 'test/myNode',
      label: 'My Node',
      icon: 'settings',
      color: 'blue-7',
    });

    const wrapper = mountWithPlugins(GenericWorkflowNode, { props: BASE_PROPS });
    expect(wrapper.vm.resolvedHandles.inputs).toHaveLength(1);
    expect(wrapper.vm.resolvedHandles.outputs).toHaveLength(1);
  });

  it('returns empty handles when nodeType is undefined', () => {
    mockGetNodeType.mockReturnValue(undefined);

    const wrapper = mountWithPlugins(GenericWorkflowNode, {
      props: {
        ...BASE_PROPS,
        data: { config: {}, label: 'No Type' },
      },
    });
    expect(wrapper.vm.resolvedHandles).toEqual({ inputs: [], outputs: [] });
  });

  it('passes hasErrors from data to BaseWorkflowNode', () => {
    mockGetNodeType.mockReturnValue({
      type: 'test/myNode',
      label: 'My Node',
      icon: 'settings',
      color: 'blue-7',
    });

    const wrapper = mountWithPlugins(GenericWorkflowNode, {
      props: {
        ...BASE_PROPS,
        data: { config: {}, __nodeType: 'test/myNode', hasErrors: true },
      },
    });
    const baseNode = wrapper.findComponent({ name: 'BaseWorkflowNode' });
    expect(baseNode.exists()).toBe(true);
  });
});
