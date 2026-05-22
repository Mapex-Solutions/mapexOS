import { describe, it, expect, vi } from 'vitest';
import { ref } from 'vue';
import { mountWithPlugins } from '@src/test/helpers';
import FieldSourceSelector from './FieldSourceSelector.vue';

/** Mock composables */
vi.mock('@src/pages/automations/workflows/createEditWorkflowPage/composables', () => ({
  useWorkflowEditorState: () => ({
    states: ref([
      { field: 'count', type: 'number' },
      { field: 'name', type: 'string' },
    ]),
    externalInputs: ref([
      { field: 'deviceId', label: 'Device ID', icon: 'devices' },
    ]),
  }),
}));

vi.mock('@stores/pluginRegistry', () => ({
  usePluginRegistryStore: () => ({
    nodeTypeMap: new Map(),
  }),
}));

vi.mock('@services/mapex', () => ({
  apis: {
    workflows: {
      credential: {
        loadOptions: vi.fn().mockResolvedValue([]),
      },
    },
  },
}));

describe('FieldSourceSelector', () => {
  const defaultProps = {
    modelValue: { type: 'literal' as const, value: '' },
    allowedTypes: ['event', 'state', 'input', 'literal', 'nodeOutput', 'loadOptions'] as const,
  };

  it('renders without errors', () => {
    const wrapper = mountWithPlugins(FieldSourceSelector, { props: { ...defaultProps } });
    expect(wrapper.exists()).toBe(true);
  });

  it('filteredTypeOptions filters by allowedTypes', () => {
    const wrapper = mountWithPlugins(FieldSourceSelector, {
      props: {
        ...defaultProps,
        allowedTypes: ['event', 'literal'],
      },
    });
    const filtered = wrapper.vm.filteredTypeOptions;
    expect(filtered).toHaveLength(2);
    expect(filtered.map((o: any) => o.value)).toEqual(['event', 'literal']);
  });

  it('filteredTypeOptions replaces loadOptions label when loadOptionsLabel is provided', () => {
    const wrapper = mountWithPlugins(FieldSourceSelector, {
      props: {
        ...defaultProps,
        allowedTypes: ['loadOptions', 'literal'],
        loadOptionsLabel: 'Search Chats',
      },
    });
    const loadOpt = wrapper.vm.filteredTypeOptions.find((o: any) => o.value === 'loadOptions');
    expect(loadOpt?.label).toBe('Search Chats');
  });

  it('showTypeSelector is true when multiple types allowed', () => {
    const wrapper = mountWithPlugins(FieldSourceSelector, {
      props: {
        ...defaultProps,
        allowedTypes: ['event', 'literal'],
      },
    });
    expect(wrapper.vm.showTypeSelector).toBe(true);
  });

  it('showTypeSelector is false when only 1 type allowed', () => {
    const wrapper = mountWithPlugins(FieldSourceSelector, {
      props: {
        ...defaultProps,
        allowedTypes: ['literal'],
      },
    });
    expect(wrapper.vm.showTypeSelector).toBe(false);
  });

  it('currentValue falls back to DEFAULT_FIELD_SOURCE_VALUE when modelValue is null', () => {
    const wrapper = mountWithPlugins(FieldSourceSelector, {
      props: {
        modelValue: null as any,
        allowedTypes: ['literal'],
      },
    });
    expect(wrapper.vm.currentValue).toEqual({ type: 'literal', value: '' });
  });

  it('updateType emits update:modelValue with new type and reset value', () => {
    const wrapper = mountWithPlugins(FieldSourceSelector, { props: { ...defaultProps } });
    wrapper.vm.updateType('event');
    expect(wrapper.emitted('update:modelValue')).toBeTruthy();
    const emitted = wrapper.emitted('update:modelValue')![0]![0] as any;
    expect(emitted.type).toBe('event');
    expect(emitted.value).toBe('');
    expect(emitted.mode).toBe('dynamic');
  });

  it('updateType for literal does not add mode', () => {
    const wrapper = mountWithPlugins(FieldSourceSelector, { props: { ...defaultProps } });
    wrapper.vm.updateType('literal');
    const emitted = wrapper.emitted('update:modelValue')![0]![0] as any;
    expect(emitted.type).toBe('literal');
    expect(emitted.mode).toBeUndefined();
  });

  it('canFetchLoadOptions is true when credentialId and loadOptionsKey are provided', () => {
    const wrapper = mountWithPlugins(FieldSourceSelector, {
      props: {
        ...defaultProps,
        credentialId: 'cred-123',
        loadOptionsKey: 'getChats',
      },
    });
    expect(wrapper.vm.canFetchLoadOptions).toBe(true);
  });

  it('canFetchLoadOptions is false when credentialId is missing', () => {
    const wrapper = mountWithPlugins(FieldSourceSelector, {
      props: {
        ...defaultProps,
        loadOptionsKey: 'getChats',
      },
    });
    expect(wrapper.vm.canFetchLoadOptions).toBe(false);
  });

  it('inputOptions contains only external inputs', () => {
    const wrapper = mountWithPlugins(FieldSourceSelector, { props: { ...defaultProps } });
    const opts = wrapper.vm.inputOptions;
    expect(opts).toHaveLength(1);
    expect(opts[0].label).toContain('input.deviceId');
    expect(opts[0].value).toBe('deviceId');
  });

  it('stateFieldNames falls back to composable states when stateFields prop is empty', () => {
    const wrapper = mountWithPlugins(FieldSourceSelector, { props: { ...defaultProps } });
    const names = wrapper.vm.stateFieldNames;
    expect(names).toHaveLength(2);
    expect(names).toContain('count');
    expect(names).toContain('name');
  });

  it('stateFieldNames uses stateFields prop when provided', () => {
    const wrapper = mountWithPlugins(FieldSourceSelector, {
      props: {
        ...defaultProps,
        stateFields: [{ name: 'custom', type: 'string' }],
      },
    });
    const names = wrapper.vm.stateFieldNames;
    expect(names).toHaveLength(1);
    expect(names[0]).toBe('custom');
  });

  it('isEventReadonly is true when type is event and mode is dynamic', () => {
    const wrapper = mountWithPlugins(FieldSourceSelector, {
      props: {
        ...defaultProps,
        modelValue: { type: 'event', value: '', mode: 'dynamic' },
      },
    });
    expect(wrapper.vm.isEventReadonly).toBe(true);
  });

  it('isEventReadonly is false when mode is manual', () => {
    const wrapper = mountWithPlugins(FieldSourceSelector, {
      props: {
        ...defaultProps,
        modelValue: { type: 'event', value: '', mode: 'manual' },
      },
    });
    expect(wrapper.vm.isEventReadonly).toBe(false);
  });

  it('stateFieldNames maps stateFields to names', () => {
    const wrapper = mountWithPlugins(FieldSourceSelector, {
      props: {
        ...defaultProps,
        stateFields: [{ name: 'foo', type: 'string' }, { name: 'bar', type: 'number' }],
      },
    });
    expect(wrapper.vm.stateFieldNames).toEqual(['foo', 'bar']);
  });

  it('handleEventClick emits openTemplateSelector when hasTemplates is false', () => {
    const wrapper = mountWithPlugins(FieldSourceSelector, {
      props: {
        ...defaultProps,
        modelValue: { type: 'event', value: '', mode: 'dynamic' },
        hasTemplates: false,
      },
    });
    wrapper.vm.handleEventClick();
    expect(wrapper.emitted('openTemplateSelector')).toBeTruthy();
  });

  it('handleEventClick emits openEventSelector when hasTemplates is true', () => {
    const wrapper = mountWithPlugins(FieldSourceSelector, {
      props: {
        ...defaultProps,
        modelValue: { type: 'event', value: '', mode: 'dynamic' },
        hasTemplates: true,
      },
    });
    wrapper.vm.handleEventClick();
    expect(wrapper.emitted('openEventSelector')).toBeTruthy();
  });
});

describe('FieldSourceSelector — literal mode (template interpolation)', () => {
  const literalProps = {
    modelValue: { type: 'literal' as const, value: '' },
    allowedTypes: ['literal'] as const,
  };

  it('renders a multi-line input (q-input type=textarea, autogrow) when type is literal', () => {
    const wrapper = mountWithPlugins(FieldSourceSelector, { props: { ...literalProps } });
    // Quasar <q-input type="textarea" autogrow> rendering varies between full DOM
    // and shallow stub in the test environment. Either:
    //  - a <textarea> element appears in the DOM (full mount), OR
    //  - the q-input stub carries autogrow / type=textarea attributes (shallow).
    const html = wrapper.html();
    const isTextarea = /<textarea/i.test(html);
    const hasAutogrow = /autogrow/i.test(html);
    const hasTypeTextarea = /type=["']?textarea["']?/i.test(html);
    expect(isTextarea || hasAutogrow || hasTypeTextarea).toBe(true);
  });

  it('emits update:modelValue with plain literal value verbatim', () => {
    const wrapper = mountWithPlugins(FieldSourceSelector, { props: { ...literalProps } });
    wrapper.vm.updateValue('hello');
    const emitted = wrapper.emitted('update:modelValue')![0]![0] as any;
    expect(emitted.type).toBe('literal');
    expect(emitted.value).toBe('hello');
  });

  it('emits update:modelValue with template string verbatim (no client-side resolution)', () => {
    const wrapper = mountWithPlugins(FieldSourceSelector, { props: { ...literalProps } });
    wrapper.vm.updateValue('Hi {{event.user.name}}');
    const emitted = wrapper.emitted('update:modelValue')![0]![0] as any;
    expect(emitted.value).toBe('Hi {{event.user.name}}');
  });

  it('exposes literalAutocompleteOpen ref initialized to false', () => {
    const wrapper = mountWithPlugins(FieldSourceSelector, { props: { ...literalProps } });
    expect(wrapper.vm.literalAutocompleteOpen).toBe(false);
  });

  it('insertNamespacePrefix replaces trailing {{ with the chosen namespace prefix', () => {
    const wrapper = mountWithPlugins(FieldSourceSelector, {
      props: {
        modelValue: { type: 'literal' as const, value: 'Hi {{' },
        allowedTypes: ['literal'] as const,
      },
    });
    wrapper.vm.insertNamespacePrefix('event');
    const emissions = wrapper.emitted('update:modelValue')!;
    const lastEmission = emissions[emissions.length - 1]![0] as any;
    expect(lastEmission.value).toBe('Hi {{event.');
    expect(lastEmission.type).toBe('literal');
  });
});
