import { describe, it, expect, vi } from 'vitest';
import { mountWithPlugins } from '@src/test/helpers';
import EventFieldInput from './EventFieldInput.vue';

/** Mock composables */
vi.mock('@composables/i18n/components/forms/useEventFieldInputTranslations', () => ({
  useEventFieldInputTranslations: () => {
    const handler: ProxyHandler<Record<string, unknown>> = {
      get(_target, prop) {
        if (prop === 'value') return String(prop);
        return new Proxy({ value: String(prop) }, handler);
      },
    };
    return new Proxy({}, handler);
  },
}));

vi.mock('@components/tooltips', () => ({
  AppTooltip: { template: '<span />' },
}));

describe('EventFieldInput', () => {
  const defaultFieldValue = { type: 'literal' as const, value: '', mode: 'dynamic' as const };

  it('renders without errors', () => {
    const wrapper = mountWithPlugins(EventFieldInput, {
      props: { modelValue: defaultFieldValue },
    });
    expect(wrapper.exists()).toBe(true);
  });

  it('fieldValue falls back to default when modelValue is undefined', () => {
    const wrapper = mountWithPlugins(EventFieldInput, {
      props: { modelValue: undefined },
    });
    expect(wrapper.vm.fieldValue.type).toBe('literal');
    expect(wrapper.vm.fieldValue.value).toBe('');
  });

  it('fieldTypeOptions has 4 options', () => {
    const wrapper = mountWithPlugins(EventFieldInput, {
      props: { modelValue: defaultFieldValue },
    });
    expect(wrapper.vm.fieldTypeOptions).toHaveLength(4);
    const values = wrapper.vm.fieldTypeOptions.map((o: any) => o.value);
    expect(values).toEqual(['event', 'state', 'variable', 'literal']);
  });

  it('fieldInputMode returns current type', () => {
    const wrapper = mountWithPlugins(EventFieldInput, {
      props: { modelValue: { type: 'event', value: 'foo', mode: 'dynamic' } },
    });
    expect(wrapper.vm.fieldInputMode).toBe('event');
  });

  it('currentModeOption returns matching option for literal', () => {
    const wrapper = mountWithPlugins(EventFieldInput, {
      props: { modelValue: defaultFieldValue },
    });
    expect(wrapper.vm.currentModeOption.value).toBe('literal');
    expect(wrapper.vm.currentModeOption.icon).toBe('format_quote');
  });

  it('isReadonly is true when type=event and mode=dynamic', () => {
    const wrapper = mountWithPlugins(EventFieldInput, {
      props: { modelValue: { type: 'event', value: '', mode: 'dynamic' } },
    });
    expect(wrapper.vm.isReadonly).toBe(true);
  });

  it('isReadonly is false when type=event and mode=manual', () => {
    const wrapper = mountWithPlugins(EventFieldInput, {
      props: { modelValue: { type: 'event', value: '', mode: 'manual' } },
    });
    expect(wrapper.vm.isReadonly).toBe(false);
  });

  it('isReadonly is false when type is not event', () => {
    const wrapper = mountWithPlugins(EventFieldInput, {
      props: { modelValue: defaultFieldValue },
    });
    expect(wrapper.vm.isReadonly).toBe(false);
  });

  it('currentMode returns null when type is not event', () => {
    const wrapper = mountWithPlugins(EventFieldInput, {
      props: { modelValue: defaultFieldValue },
    });
    expect(wrapper.vm.currentMode).toBeNull();
  });

  it('currentMode returns dynamic by default for event type', () => {
    const wrapper = mountWithPlugins(EventFieldInput, {
      props: { modelValue: { type: 'event', value: '' } },
    });
    expect(wrapper.vm.currentMode).toBe('dynamic');
  });

  it('availableStateNames maps stateFields to names', () => {
    const wrapper = mountWithPlugins(EventFieldInput, {
      props: {
        modelValue: defaultFieldValue,
        stateFields: [{ name: 'alpha', type: 'string' }, { name: 'beta', type: 'number' }],
      },
    });
    expect(wrapper.vm.availableStateNames).toEqual(['alpha', 'beta']);
  });

  it('inputValue reflects modelValue value', () => {
    const wrapper = mountWithPlugins(EventFieldInput, {
      props: { modelValue: { type: 'literal', value: 'hello' } },
    });
    expect(wrapper.vm.inputValue).toBe('hello');
  });

  it('handleInputClick emits openTemplateSelector when no templates', () => {
    const wrapper = mountWithPlugins(EventFieldInput, {
      props: {
        modelValue: { type: 'event', value: '', mode: 'dynamic' },
        hasTemplates: false,
      },
    });
    wrapper.vm.handleInputClick();
    expect(wrapper.emitted('openTemplateSelector')).toBeTruthy();
  });

  it('handleInputClick emits openEventSelector when templates available', () => {
    const wrapper = mountWithPlugins(EventFieldInput, {
      props: {
        modelValue: { type: 'event', value: '', mode: 'dynamic' },
        hasTemplates: true,
      },
    });
    wrapper.vm.handleInputClick();
    expect(wrapper.emitted('openEventSelector')).toBeTruthy();
  });

  it('handleInputClick does nothing in manual mode', () => {
    const wrapper = mountWithPlugins(EventFieldInput, {
      props: {
        modelValue: { type: 'event', value: '', mode: 'manual' },
        hasTemplates: true,
      },
    });
    wrapper.vm.handleInputClick();
    expect(wrapper.emitted('openEventSelector')).toBeFalsy();
    expect(wrapper.emitted('openTemplateSelector')).toBeFalsy();
  });

  it('switchToManualMode emits update:modelValue with mode manual', () => {
    const wrapper = mountWithPlugins(EventFieldInput, {
      props: { modelValue: { type: 'event', value: 'x', mode: 'dynamic' } },
    });
    wrapper.vm.switchToManualMode();
    const emitted = wrapper.emitted('update:modelValue')![0]![0] as any;
    expect(emitted.mode).toBe('manual');
  });

  it('switchToDynamicMode emits update:modelValue with mode dynamic', () => {
    const wrapper = mountWithPlugins(EventFieldInput, {
      props: { modelValue: { type: 'event', value: 'x', mode: 'manual' } },
    });
    wrapper.vm.switchToDynamicMode();
    const emitted = wrapper.emitted('update:modelValue')![0]![0] as any;
    expect(emitted.mode).toBe('dynamic');
  });
});
