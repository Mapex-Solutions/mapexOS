import { describe, it, expect, vi } from 'vitest';
import { mountWithPlugins } from '@src/test/helpers';
import TextNoteNode from './TextNoteNode.vue';

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
  id: 'note-1',
  data: {
    config: {
      text: 'This is a note',
      color: 'amber',
    },
    __nodeType: 'core/text_note',
  },
  selected: false,
};

describe('TextNoteNode', () => {
  it('renders without errors', () => {
    const wrapper = mountWithPlugins(TextNoteNode, { props: BASE_PROPS });
    expect(wrapper.exists()).toBe(true);
  });

  it('computes text from config', () => {
    const wrapper = mountWithPlugins(TextNoteNode, { props: BASE_PROPS });
    expect(wrapper.vm.text).toBe('This is a note');
  });

  it('defaults text to empty string when not set', () => {
    const wrapper = mountWithPlugins(TextNoteNode, {
      props: {
        ...BASE_PROPS,
        data: { config: {}, __nodeType: 'core/text_note' },
      },
    });
    expect(wrapper.vm.text).toBe('');
  });

  it('computes bgColor from config', () => {
    const wrapper = mountWithPlugins(TextNoteNode, { props: BASE_PROPS });
    expect(wrapper.vm.bgColor).toBe('amber');
  });

  it('defaults bgColor to grey when not set', () => {
    const wrapper = mountWithPlugins(TextNoteNode, {
      props: {
        ...BASE_PROPS,
        data: { config: {}, __nodeType: 'core/text_note' },
      },
    });
    expect(wrapper.vm.bgColor).toBe('grey');
  });

  it('applies color variant class based on bgColor', () => {
    const wrapper = mountWithPlugins(TextNoteNode, { props: BASE_PROPS });
    expect(wrapper.find('.text-note').classes()).toContain('text-note--amber');
  });

  it('applies selected class when selected is true', () => {
    const wrapper = mountWithPlugins(TextNoteNode, {
      props: { ...BASE_PROPS, selected: true },
    });
    expect(wrapper.find('.text-note').classes()).toContain('text-note--selected');
  });

  it('displays note text when text is set', () => {
    const wrapper = mountWithPlugins(TextNoteNode, { props: BASE_PROPS });
    expect(wrapper.find('.text-note__text').exists()).toBe(true);
    expect(wrapper.find('.text-note__text').text()).toBe('This is a note');
  });

  it('displays placeholder when text is empty', () => {
    const wrapper = mountWithPlugins(TextNoteNode, {
      props: {
        ...BASE_PROPS,
        data: { config: { color: 'grey' }, __nodeType: 'core/text_note' },
      },
    });
    expect(wrapper.find('.text-note__placeholder').exists()).toBe(true);
  });

  it('starts in non-editing mode', () => {
    const wrapper = mountWithPlugins(TextNoteNode, { props: BASE_PROPS });
    expect(wrapper.vm.isEditing).toBe(false);
  });

  it('does not show textarea in non-editing mode', () => {
    const wrapper = mountWithPlugins(TextNoteNode, { props: BASE_PROPS });
    expect(wrapper.find('.text-note__textarea').exists()).toBe(false);
  });
});
