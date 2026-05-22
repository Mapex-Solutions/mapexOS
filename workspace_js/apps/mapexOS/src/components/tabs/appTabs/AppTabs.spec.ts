import { describe, it, expect } from 'vitest';
import { mountWithPlugins } from '@src/test/helpers';
import AppTabs from './AppTabs.vue';
import type { AppTabItem } from './interfaces';

const MOCK_TABS: AppTabItem[] = [
  { name: 'general', label: 'General', icon: 'settings' },
  { name: 'members', label: 'Members', icon: 'people', badge: 5, badgeColor: 'red' },
  { name: 'files', label: 'Files' },
  { name: 'disabled', label: 'Disabled', disabled: true },
];

describe('AppTabs', () => {
  it('renders with required props', () => {
    const wrapper = mountWithPlugins(AppTabs, {
      props: { modelValue: 'general', tabs: MOCK_TABS },
    });

    expect(wrapper.find('.app-tabs').exists()).toBe(true);
  });

  it('renders q-tabs component', () => {
    const wrapper = mountWithPlugins(AppTabs, {
      props: { modelValue: 'general', tabs: MOCK_TABS },
    });

    expect(wrapper.find('q-tabs-stub').exists()).toBe(true);
  });

  it('applies bordered class when bordered prop is true', () => {
    const wrapper = mountWithPlugins(AppTabs, {
      props: { modelValue: 'general', tabs: MOCK_TABS, bordered: true },
    });

    expect(wrapper.find('.app-tabs--bordered').exists()).toBe(true);
  });

  it('does not apply bordered class by default', () => {
    const wrapper = mountWithPlugins(AppTabs, {
      props: { modelValue: 'general', tabs: MOCK_TABS },
    });

    expect(wrapper.find('.app-tabs--bordered').exists()).toBe(false);
  });

  it('renders separator by default (non-pill variant)', () => {
    const wrapper = mountWithPlugins(AppTabs, {
      props: { modelValue: 'general', tabs: MOCK_TABS },
    });

    expect(wrapper.find('q-separator-stub').exists()).toBe(true);
  });

  it('does not render separator when separator prop is false', () => {
    const wrapper = mountWithPlugins(AppTabs, {
      props: { modelValue: 'general', tabs: MOCK_TABS, separator: false },
    });

    expect(wrapper.find('q-separator-stub').exists()).toBe(false);
  });

  it('does not render separator for pill variant', () => {
    const wrapper = mountWithPlugins(AppTabs, {
      props: { modelValue: 'general', tabs: MOCK_TABS, variant: 'pill' },
    });

    expect(wrapper.find('q-separator-stub').exists()).toBe(false);
  });

  it('applies pill class when variant is pill', () => {
    const wrapper = mountWithPlugins(AppTabs, {
      props: { modelValue: 'general', tabs: MOCK_TABS, variant: 'pill' },
    });

    expect(wrapper.find('.app-tabs--pill').exists()).toBe(true);
  });

  it('does not apply pill class for default variant', () => {
    const wrapper = mountWithPlugins(AppTabs, {
      props: { modelValue: 'general', tabs: MOCK_TABS },
    });

    expect(wrapper.find('.app-tabs--pill').exists()).toBe(false);
  });

  it('uses transparent indicator for pill variant', () => {
    const wrapper = mountWithPlugins(AppTabs, {
      props: { modelValue: 'general', tabs: MOCK_TABS, variant: 'pill' },
    });

    expect(wrapper.find('q-tabs-stub').attributes('indicator-color')).toBe('transparent');
  });

  it('uses primary indicator for default variant', () => {
    const wrapper = mountWithPlugins(AppTabs, {
      props: { modelValue: 'general', tabs: MOCK_TABS },
    });

    expect(wrapper.find('q-tabs-stub').attributes('indicator-color')).toBe('primary');
  });

  it('passes align prop to q-tabs', () => {
    const wrapper = mountWithPlugins(AppTabs, {
      props: { modelValue: 'general', tabs: MOCK_TABS, align: 'center' },
    });

    expect(wrapper.find('q-tabs-stub').attributes('align')).toBe('center');
  });

  it('emits update:modelValue and change when activeTab changes', async () => {
    const wrapper = mountWithPlugins(AppTabs, {
      props: { modelValue: 'general', tabs: MOCK_TABS },
    });

    const qTabs = wrapper.findComponent({ name: 'q-tabs' });
    await qTabs.vm.$emit('update:modelValue', 'members');

    expect(wrapper.emitted('update:modelValue')).toBeTruthy();
    expect(wrapper.emitted('update:modelValue')![0]).toEqual(['members']);
    expect(wrapper.emitted('change')).toBeTruthy();
    expect(wrapper.emitted('change')![0]).toEqual(['members']);
  });

  it('renders default slot content', () => {
    const wrapper = mountWithPlugins(AppTabs, {
      props: { modelValue: 'general', tabs: MOCK_TABS },
      slots: { default: '<div class="panel-content">Panel</div>' },
    });

    expect(wrapper.find('.panel-content').exists()).toBe(true);
  });

  it('sets narrow-indicator to true for default variant', () => {
    const wrapper = mountWithPlugins(AppTabs, {
      props: { modelValue: 'general', tabs: MOCK_TABS },
    });

    expect(wrapper.find('q-tabs-stub').attributes('narrow-indicator')).toBe('true');
  });

  it('sets narrow-indicator to false for pill variant', () => {
    const wrapper = mountWithPlugins(AppTabs, {
      props: { modelValue: 'general', tabs: MOCK_TABS, variant: 'pill' },
    });

    expect(wrapper.find('q-tabs-stub').attributes('narrow-indicator')).toBe('false');
  });

  it('applies default-tabs class for default variant', () => {
    const wrapper = mountWithPlugins(AppTabs, {
      props: { modelValue: 'general', tabs: MOCK_TABS },
    });

    expect(wrapper.find('q-tabs-stub').attributes('class')).toContain('app-tabs__default-tabs');
  });

  it('applies pill-tabs class for pill variant', () => {
    const wrapper = mountWithPlugins(AppTabs, {
      props: { modelValue: 'general', tabs: MOCK_TABS, variant: 'pill' },
    });

    expect(wrapper.find('q-tabs-stub').attributes('class')).toContain('app-tabs__pill-tabs');
  });
});
