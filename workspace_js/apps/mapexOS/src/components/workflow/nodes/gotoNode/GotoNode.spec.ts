import { describe, it, expect, vi } from 'vitest';
import { mountWithPlugins } from '@src/test/helpers';
import GotoNode from './GotoNode.vue';

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
      type: 'core/goto',
      label: 'Goto',
      icon: 'near_me',
      color: 'deep-purple-6',
      inputs: [],
      outputs: [],
    }),
    updateNodeConfig: vi.fn(),
    nodes: { value: [] },
    addNoteToNode: vi.fn(),
    pushSnapshot: vi.fn(),
  }),
  usePluginI18n: () => ({
    t: (key: string) => key,
  }),
}));

vi.mock('@src/utils/workflow', () => ({
  resolveNodeHandles: () => ({
    inputs: [{ id: 'in-1', position: 'top', label: 'In' }],
    outputs: [{ id: 'out-1', position: 'bottom', label: 'Out' }],
  }),
}));

const BASE_PROPS = {
  id: 'goto-1',
  data: {
    config: { role: 'sender', pairLabel: 'ErrorHandler', pairColor: 'deep-purple-6' },
    label: 'Goto',
    __nodeType: 'core/goto',
  },
  selected: false,
};

describe('GotoNode', () => {
  it('renders without errors', () => {
    const wrapper = mountWithPlugins(GotoNode, { props: BASE_PROPS });
    expect(wrapper.exists()).toBe(true);
  });

  it('computes role from config', () => {
    const wrapper = mountWithPlugins(GotoNode, { props: BASE_PROPS });
    expect(wrapper.vm.role).toBe('sender');
  });

  it('defaults role to sender when not set', () => {
    const wrapper = mountWithPlugins(GotoNode, {
      props: {
        ...BASE_PROPS,
        data: { config: {}, label: 'Goto', __nodeType: 'core/goto' },
      },
    });
    expect(wrapper.vm.role).toBe('sender');
  });

  it('computes pairLabel from config', () => {
    const wrapper = mountWithPlugins(GotoNode, { props: BASE_PROPS });
    expect(wrapper.vm.pairLabel).toBe('ErrorHandler');
  });

  it('computes pairColor from config', () => {
    const wrapper = mountWithPlugins(GotoNode, { props: BASE_PROPS });
    expect(wrapper.vm.pairColor).toBe('deep-purple-6');
  });

  it('computes colorHex from GOTO_COLOR_OPTIONS for known color', () => {
    const wrapper = mountWithPlugins(GotoNode, { props: BASE_PROPS });
    expect(wrapper.vm.colorHex).toBe('#5e35b1');
  });

  it('computes colorHex fallback for custom hex', () => {
    const wrapper = mountWithPlugins(GotoNode, {
      props: {
        ...BASE_PROPS,
        data: {
          config: { role: 'sender', pairLabel: 'Test', pairColor: '#123456' },
          label: 'Goto',
          __nodeType: 'core/goto',
        },
      },
    });
    expect(wrapper.vm.colorHex).toBe('#123456');
  });

  it('computes colorHex default for unknown color name', () => {
    const wrapper = mountWithPlugins(GotoNode, {
      props: {
        ...BASE_PROPS,
        data: {
          config: { role: 'sender', pairLabel: 'Test', pairColor: 'unknown-color' },
          label: 'Goto',
          __nodeType: 'core/goto',
        },
      },
    });
    expect(wrapper.vm.colorHex).toBe('#5e35b1');
  });

  it('computes roleIcon as near_me for sender', () => {
    const wrapper = mountWithPlugins(GotoNode, { props: BASE_PROPS });
    expect(wrapper.vm.roleIcon).toBe('near_me');
  });

  it('computes roleIcon as place for receiver', () => {
    const wrapper = mountWithPlugins(GotoNode, {
      props: {
        ...BASE_PROPS,
        data: {
          config: { role: 'receiver', pairLabel: 'ErrorHandler', pairColor: 'deep-purple-6' },
          label: 'Goto',
          __nodeType: 'core/goto',
        },
      },
    });
    expect(wrapper.vm.roleIcon).toBe('place');
  });

  it('renders sender badge below the node', () => {
    const wrapper = mountWithPlugins(GotoNode, { props: BASE_PROPS });
    const badge = wrapper.find('.goto-node-wrapper__badge--bottom');
    expect(badge.exists()).toBe(true);
    expect(badge.text()).toBe('ErrorHandler');
  });

  it('renders receiver badge above the node', () => {
    const wrapper = mountWithPlugins(GotoNode, {
      props: {
        ...BASE_PROPS,
        data: {
          config: { role: 'receiver', pairLabel: 'ErrorHandler', pairColor: 'deep-purple-6' },
          label: 'Goto',
          __nodeType: 'core/goto',
        },
      },
    });
    const badge = wrapper.find('.goto-node-wrapper__badge--top');
    expect(badge.exists()).toBe(true);
    expect(badge.text()).toBe('ErrorHandler');
  });

  it('does not render badge when pairLabel is empty', () => {
    const wrapper = mountWithPlugins(GotoNode, {
      props: {
        ...BASE_PROPS,
        data: {
          config: { role: 'sender', pairLabel: '', pairColor: 'deep-purple-6' },
          label: 'Goto',
          __nodeType: 'core/goto',
        },
      },
    });
    expect(wrapper.find('.goto-node-wrapper__badge--bottom').exists()).toBe(false);
    expect(wrapper.find('.goto-node-wrapper__badge--top').exists()).toBe(false);
  });
});
