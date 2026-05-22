import { describe, it, expect, vi } from 'vitest';
import { mountWithPlugins } from '@src/test/helpers';
import EndNode from './EndNode.vue';

vi.mock('@src/components/workflow/constants', () => ({
  POSITION_OPTIONS: [
    { label: 'Top', value: 'top', icon: 'arrow_upward' },
    { label: 'Bottom', value: 'bottom', icon: 'arrow_downward' },
    { label: 'Left', value: 'left', icon: 'arrow_back' },
    { label: 'Right', value: 'right', icon: 'arrow_forward' },
  ],
}));

vi.mock('@src/composables/workflow', () => ({
  useWorkflowContext: () => ({
    getNodeType: () => ({
      type: 'core/end',
      label: 'End',
      icon: 'stop_circle',
      color: 'purple-7',
      inputs: [],
      outputs: [],
      deletable: true,
      shape: 'circle',
    }),
    updateNodeConfig: vi.fn(),
    nodes: { value: [] },
    addNoteToNode: vi.fn(),
    pushSnapshot: vi.fn(),
  }),
}));

vi.mock('@src/utils/workflow', () => ({
  resolveNodeHandles: () => ({
    inputs: [{ id: 'in-1', position: 'top', label: 'In' }],
    outputs: [],
  }),
}));

const BASE_PROPS = {
  id: 'end-1',
  data: {
    config: {},
    label: 'End',
    __nodeType: 'core/end',
  },
  selected: false,
};

describe('EndNode', () => {
  it('renders without errors', () => {
    const wrapper = mountWithPlugins(EndNode, { props: BASE_PROPS });
    expect(wrapper.exists()).toBe(true);
  });

  it('computes isError as false when terminateWithError is not set', () => {
    const wrapper = mountWithPlugins(EndNode, { props: BASE_PROPS });
    expect(wrapper.vm.isError).toBe(false);
  });

  it('computes isError as true when terminateWithError is true', () => {
    const wrapper = mountWithPlugins(EndNode, {
      props: {
        ...BASE_PROPS,
        data: {
          config: { terminateWithError: true },
          label: 'End',
          __nodeType: 'core/end',
        },
      },
    });
    expect(wrapper.vm.isError).toBe(true);
  });

  it('computes nodeIcon as check_circle for success mode', () => {
    const wrapper = mountWithPlugins(EndNode, { props: BASE_PROPS });
    expect(wrapper.vm.nodeIcon).toBe('check_circle');
  });

  it('computes nodeIcon as cancel for error mode', () => {
    const wrapper = mountWithPlugins(EndNode, {
      props: {
        ...BASE_PROPS,
        data: {
          config: { terminateWithError: true },
          label: 'End',
          __nodeType: 'core/end',
        },
      },
    });
    expect(wrapper.vm.nodeIcon).toBe('cancel');
  });

  it('computes nodeColorHex as purple for success mode', () => {
    const wrapper = mountWithPlugins(EndNode, { props: BASE_PROPS });
    expect(wrapper.vm.nodeColorHex).toBe('#7B1FA2');
  });

  it('computes nodeColorHex as red for error mode', () => {
    const wrapper = mountWithPlugins(EndNode, {
      props: {
        ...BASE_PROPS,
        data: {
          config: { terminateWithError: true },
          label: 'End',
          __nodeType: 'core/end',
        },
      },
    });
    expect(wrapper.vm.nodeColorHex).toBe('#D32F2F');
  });

  it('computes resolvedHandles from resolveNodeHandles', () => {
    const wrapper = mountWithPlugins(EndNode, { props: BASE_PROPS });
    expect(wrapper.vm.resolvedHandles.inputs).toHaveLength(1);
    expect(wrapper.vm.resolvedHandles.outputs).toHaveLength(0);
  });
});
