import { describe, it, expect } from 'vitest';
import { mountWithPlugins } from '@src/test/helpers';
import IconSectionNav from './IconSectionNav.vue';

describe('IconSectionNav', () => {
  const mockItems = [
    { name: 'general', icon: 'settings', tooltip: 'General Settings' },
    { name: 'security', icon: 'lock', tooltip: 'Security' },
    { name: 'notifications', icon: 'notifications', tooltip: 'Notifications', badge: true, badgeColor: 'red' },
  ];

  const defaultProps = {
    modelValue: 'general',
    items: mockItems,
  };

  it('renders without errors', () => {
    const wrapper = mountWithPlugins(IconSectionNav, { props: defaultProps });
    expect(wrapper.exists()).toBe(true);
  });

  it('receives items prop', () => {
    const wrapper = mountWithPlugins(IconSectionNav, { props: defaultProps });
    expect(wrapper.props('items')).toEqual(mockItems);
    expect(wrapper.props('items')).toHaveLength(3);
  });

  it('receives modelValue prop', () => {
    const wrapper = mountWithPlugins(IconSectionNav, { props: defaultProps });
    expect(wrapper.props('modelValue')).toBe('general');
  });

  it('defaults width to 40', () => {
    const wrapper = mountWithPlugins(IconSectionNav, { props: defaultProps });
    expect(wrapper.props('width')).toBe(40);
  });

  it('computes railStyle with correct width', () => {
    const wrapper = mountWithPlugins(IconSectionNav, { props: defaultProps });
    expect(wrapper.vm.railStyle).toEqual({ width: '40px', minWidth: '40px' });
  });

  it('computes railStyle with custom width', () => {
    const wrapper = mountWithPlugins(IconSectionNav, {
      props: { ...defaultProps, width: 60 },
    });
    expect(wrapper.vm.railStyle).toEqual({ width: '60px', minWidth: '60px' });
  });

  it('isActive returns true for current modelValue', () => {
    const wrapper = mountWithPlugins(IconSectionNav, { props: defaultProps });
    expect((wrapper.vm).isActive('general')).toBe(true);
    expect((wrapper.vm).isActive('security')).toBe(false);
  });

  it('emits update:modelValue on handleItemClick', () => {
    const wrapper = mountWithPlugins(IconSectionNav, { props: defaultProps });
    (wrapper.vm).handleItemClick('security');
    expect(wrapper.emitted('update:modelValue')).toBeTruthy();
    expect(wrapper.emitted('update:modelValue')![0]).toEqual(['security']);
  });

  it('renders nav element with role tablist', () => {
    const wrapper = mountWithPlugins(IconSectionNav, { props: defaultProps });
    const nav = wrapper.find('nav');
    expect(nav.exists()).toBe(true);
    expect(nav.attributes('role')).toBe('tablist');
  });

  it('renders a button for each item', () => {
    const wrapper = mountWithPlugins(IconSectionNav, { props: defaultProps });
    const buttons = wrapper.findAll('button');
    expect(buttons).toHaveLength(3);
  });

  it('applies active class to the current item button', () => {
    const wrapper = mountWithPlugins(IconSectionNav, { props: defaultProps });
    const buttons = wrapper.findAll('button');
    expect(buttons[0]!.classes()).toContain('icon-section-nav__item--active');
    expect(buttons[1]!.classes()).not.toContain('icon-section-nav__item--active');
  });

  it('sets aria-selected on active item', () => {
    const wrapper = mountWithPlugins(IconSectionNav, { props: defaultProps });
    const buttons = wrapper.findAll('button');
    expect(buttons[0]!.attributes('aria-selected')).toBe('true');
    expect(buttons[1]!.attributes('aria-selected')).toBe('false');
  });
});
