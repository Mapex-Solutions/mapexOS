import { describe, it, expect } from 'vitest';
import { mountWithPlugins } from '@src/test/helpers';
import CardBodyDetails from './CardBodyDetails.vue';

function factory(overrides: Record<string, unknown> = {}) {
  return mountWithPlugins(CardBodyDetails, {
    props: {
      title: 'Details Section',
      items: [
        { name: 'Field 1', value: 'Value 1', type: 'text' },
        { name: 'Field 2', value: 'Value 2', type: 'chip', color: 'blue-5' },
      ],
      ...overrides,
    },
  });
}

describe('CardBodyDetails', () => {
  it('renders without errors', () => {
    const wrapper = factory();
    expect(wrapper.exists()).toBe(true);
  });

  it('receives title prop correctly', () => {
    const wrapper = factory();
    expect(wrapper.props('title')).toBe('Details Section');
  });

  it('receives items prop correctly', () => {
    const wrapper = factory();
    expect(wrapper.props('items')).toHaveLength(2);
  });

  it('accepts optional tenantName prop', () => {
    const wrapper = factory({ tenantName: 'Tenant ABC' });
    expect(wrapper.props('tenantName')).toBe('Tenant ABC');
  });

  it('accepts optional container prop with color', () => {
    const wrapper = factory({ container: { color: 'bg-blue-1' } });
    expect(wrapper.props('container')).toEqual({ color: 'bg-blue-1' });
  });

  it('renders with icon type items', () => {
    const wrapper = factory({
      items: [
        { name: 'Status', type: 'icon', icon: 'check', iconColor: 'green-5', value: 'OK' },
      ],
    });
    expect(wrapper.exists()).toBe(true);
  });

  it('renders with card type items', () => {
    const wrapper = factory({
      items: [
        { name: 'Count', type: 'card', icon: 'devices', iconColor: 'teal-5', value: 42, color: 'bg-teal-1' },
      ],
    });
    expect(wrapper.exists()).toBe(true);
  });

  it('renders with iconsGroup type items', () => {
    const wrapper = factory({
      items: [
        {
          name: 'Protocols',
          type: 'iconsGroup',
          icons: [
            { icon: 'wifi', color: 'blue-5', tooltip: 'WiFi' },
            { icon: 'bluetooth', color: 'indigo-5' },
          ],
        },
      ],
    });
    expect(wrapper.exists()).toBe(true);
  });
});
