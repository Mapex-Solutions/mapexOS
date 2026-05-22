import { describe, it, expect } from 'vitest';
import { mountWithPlugins } from '@src/test/helpers';
import AvailableFieldsList from './AvailableFieldsList.vue';

describe('AvailableFieldsList', () => {
  const defaultProps = {
    fields: ['payload.temperature', 'payload.humidity', 'payload.timestamp'],
  };

  it('renders without errors', () => {
    const wrapper = mountWithPlugins(AvailableFieldsList, { props: defaultProps });
    expect(wrapper.exists()).toBe(true);
  });

  it('receives fields prop', () => {
    const wrapper = mountWithPlugins(AvailableFieldsList, { props: defaultProps });
    expect(wrapper.props('fields')).toEqual([
      'payload.temperature',
      'payload.humidity',
      'payload.timestamp',
    ]);
  });

  it('defaults loading to false', () => {
    const wrapper = mountWithPlugins(AvailableFieldsList, { props: defaultProps });
    expect(wrapper.props('loading')).toBe(false);
  });

  it('defaults maxHeight to 300', () => {
    const wrapper = mountWithPlugins(AvailableFieldsList, { props: defaultProps });
    expect(wrapper.props('maxHeight')).toBe(300);
  });

  it('shouldScroll is false when fields count <= threshold', () => {
    const wrapper = mountWithPlugins(AvailableFieldsList, {
      props: { fields: ['a', 'b', 'c'] },
    });
    expect(wrapper.vm.shouldScroll).toBe(false);
  });

  it('shouldScroll is true when fields count > threshold (10)', () => {
    const fields = Array.from({ length: 15 }, (_, i) => `field.${i}`);
    const wrapper = mountWithPlugins(AvailableFieldsList, {
      props: { fields },
    });
    expect(wrapper.vm.shouldScroll).toBe(true);
  });

  it('scrollAreaHeight is "auto" when shouldScroll is false', () => {
    const wrapper = mountWithPlugins(AvailableFieldsList, {
      props: { fields: ['a', 'b'] },
    });
    expect(wrapper.vm.scrollAreaHeight).toBe('auto');
  });

  it('scrollAreaHeight uses maxHeight when shouldScroll is true', () => {
    const fields = Array.from({ length: 15 }, (_, i) => `field.${i}`);
    const wrapper = mountWithPlugins(AvailableFieldsList, {
      props: { fields, maxHeight: 400 },
    });
    expect(wrapper.vm.scrollAreaHeight).toBe('400px');
  });

  it('emits field-click when handleFieldClick is called', () => {
    const wrapper = mountWithPlugins(AvailableFieldsList, { props: defaultProps });
    (wrapper.vm).handleFieldClick('payload.temperature');
    expect(wrapper.emitted('field-click')).toBeTruthy();
    expect(wrapper.emitted('field-click')![0]).toEqual(['payload.temperature']);
  });

  it('renders with empty fields array', () => {
    const wrapper = mountWithPlugins(AvailableFieldsList, {
      props: { fields: [] },
    });
    expect(wrapper.exists()).toBe(true);
    expect(wrapper.vm.shouldScroll).toBe(false);
  });
});
