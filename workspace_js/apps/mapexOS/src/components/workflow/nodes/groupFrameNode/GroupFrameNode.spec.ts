import { describe, it, expect, vi } from 'vitest';
import { mountWithPlugins } from '@src/test/helpers';
import GroupFrameNode from './GroupFrameNode.vue';

vi.mock('@src/composables/workflow', () => ({
  useWorkflowContext: () => ({
    updateNodeConfig: vi.fn(),
    getNodeType: () => undefined,
    nodes: { value: [] },
    addNoteToNode: vi.fn(),
    pushSnapshot: vi.fn(),
  }),
  usePluginI18n: () => ({
    t: (key: string) => key,
  }),
}));

const BASE_PROPS = {
  id: 'frame-1',
  data: {
    config: {
      title: 'My Group',
      description: 'A frame description',
      color: 'blue-grey',
      width: 400,
      height: 300,
    },
    __nodeType: 'core/group_frame',
  },
  selected: false,
};

describe('GroupFrameNode', () => {
  it('renders without errors', () => {
    const wrapper = mountWithPlugins(GroupFrameNode, { props: BASE_PROPS });
    expect(wrapper.exists()).toBe(true);
  });

  it('computes title from config', () => {
    const wrapper = mountWithPlugins(GroupFrameNode, { props: BASE_PROPS });
    expect(wrapper.vm.title).toBe('My Group');
  });

  it('defaults title to empty string when not set', () => {
    const wrapper = mountWithPlugins(GroupFrameNode, {
      props: {
        ...BASE_PROPS,
        data: { config: {}, __nodeType: 'core/group_frame' },
      },
    });
    expect(wrapper.vm.title).toBe('');
  });

  it('computes description from config', () => {
    const wrapper = mountWithPlugins(GroupFrameNode, { props: BASE_PROPS });
    expect(wrapper.vm.description).toBe('A frame description');
  });

  it('computes colorName from config', () => {
    const wrapper = mountWithPlugins(GroupFrameNode, { props: BASE_PROPS });
    expect(wrapper.vm.colorName).toBe('blue-grey');
  });

  it('computes colorHex from FRAME_COLOR_OPTIONS for known color', () => {
    const wrapper = mountWithPlugins(GroupFrameNode, { props: BASE_PROPS });
    // blue-grey should resolve to a hex from FRAME_COLOR_OPTIONS
    expect(wrapper.vm.colorHex).toMatch(/^#[0-9a-fA-F]{6}$/);
  });

  it('computes colorHex fallback for custom hex', () => {
    const wrapper = mountWithPlugins(GroupFrameNode, {
      props: {
        ...BASE_PROPS,
        data: {
          config: { ...BASE_PROPS.data.config, color: '#AABBCC' },
          __nodeType: 'core/group_frame',
        },
      },
    });
    expect(wrapper.vm.colorHex).toBe('#AABBCC');
  });

  it('computes width from config', () => {
    const wrapper = mountWithPlugins(GroupFrameNode, { props: BASE_PROPS });
    expect(wrapper.vm.width).toBe(400);
  });

  it('defaults width to 300 when not set', () => {
    const wrapper = mountWithPlugins(GroupFrameNode, {
      props: {
        ...BASE_PROPS,
        data: { config: {}, __nodeType: 'core/group_frame' },
      },
    });
    expect(wrapper.vm.width).toBe(300);
  });

  it('computes height from config', () => {
    const wrapper = mountWithPlugins(GroupFrameNode, { props: BASE_PROPS });
    expect(wrapper.vm.height).toBe(300);
  });

  it('defaults height to 200 when not set', () => {
    const wrapper = mountWithPlugins(GroupFrameNode, {
      props: {
        ...BASE_PROPS,
        data: { config: {}, __nodeType: 'core/group_frame' },
      },
    });
    expect(wrapper.vm.height).toBe(200);
  });

  it('computes frameStyle with dimensions and color', () => {
    const wrapper = mountWithPlugins(GroupFrameNode, { props: BASE_PROPS });
    expect(wrapper.vm.frameStyle).toHaveProperty('width', '400px');
    expect(wrapper.vm.frameStyle).toHaveProperty('height', '300px');
    expect(wrapper.vm.frameStyle).toHaveProperty('--frame-color');
  });

  it('applies selected class when selected is true', () => {
    const wrapper = mountWithPlugins(GroupFrameNode, {
      props: { ...BASE_PROPS, selected: true },
    });
    expect(wrapper.find('.group-frame').classes()).toContain('group-frame--selected');
  });

  it('renders title text in the header', () => {
    const wrapper = mountWithPlugins(GroupFrameNode, { props: BASE_PROPS });
    expect(wrapper.find('.group-frame__title').text()).toBe('My Group');
  });

  it('renders description when present', () => {
    const wrapper = mountWithPlugins(GroupFrameNode, { props: BASE_PROPS });
    expect(wrapper.find('.group-frame__description').exists()).toBe(true);
    expect(wrapper.find('.group-frame__description').text()).toBe('A frame description');
  });

  it('does not render description when empty', () => {
    const wrapper = mountWithPlugins(GroupFrameNode, {
      props: {
        ...BASE_PROPS,
        data: {
          config: { title: 'No Desc', color: 'blue-grey', width: 300, height: 200 },
          __nodeType: 'core/group_frame',
        },
      },
    });
    expect(wrapper.find('.group-frame__description').exists()).toBe(false);
  });
});
