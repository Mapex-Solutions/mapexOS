import { describe, it, expect, vi } from 'vitest';
import { nextTick } from 'vue';
import { mountWithPlugins } from '@src/test/helpers';
import FontIconDialog from './FontIconDialog.vue';

// Mock the icons.json with a small subset
vi.mock('./fonts/icons.json', () => ({
  default: {
    icons: [
      { name: 'home', categories: ['navigation'] },
      { name: 'search', categories: ['action'] },
      { name: 'settings', categories: ['action'] },
      { name: 'person', categories: ['social'] },
      { name: 'star', categories: ['toggle'] },
      { name: 'favorite', categories: ['action'] },
      { name: 'delete', categories: ['action'] },
      { name: 'add', categories: ['navigation'] },
    ],
  },
}));

// Mock constants
vi.mock('./constants', () => ({
  CATEGORY_ICON_MAP: {
    navigation: 'explore',
    action: 'touch_app',
    social: 'people',
    toggle: 'toggle_on',
  } as Record<string, string>,
  CATEGORY_DESCRIPTION_MAP: {
    navigation: 'Navigation icons',
    action: 'Action icons',
    social: 'Social icons',
    toggle: 'Toggle icons',
  } as Record<string, string>,
}));

describe('FontIconDialog', () => {
  const defaultProps = {
    modelValue: '',
    show: true,
  };

  it('renders without errors', () => {
    const wrapper = mountWithPlugins(FontIconDialog, { props: defaultProps });
    expect(wrapper.exists()).toBe(true);
  });

  it('initializes allIcons from JSON on mount', () => {
    const wrapper = mountWithPlugins(FontIconDialog, { props: defaultProps });
    expect(wrapper.vm.allIcons).toHaveLength(8);
  });

  it('computes unique sorted categories', () => {
    const wrapper = mountWithPlugins(FontIconDialog, { props: defaultProps });
    expect(wrapper.vm.categories).toEqual(['action', 'navigation', 'social', 'toggle']);
  });

  it('filters icons by search term via filteredIcons computed', async () => {
    const wrapper = mountWithPlugins(FontIconDialog, { props: defaultProps });
    // First clear category so filter takes precedence
    wrapper.vm.selectedCategory = '';
    await nextTick();
    // Now set filter
    wrapper.vm.filter = 'home';
    await nextTick();
    expect(wrapper.vm.filteredIcons).toHaveLength(1);
    expect(wrapper.vm.filteredIcons[0].name).toBe('home');
  });

  it('filters icons by selected category when no search term', async () => {
    const wrapper = mountWithPlugins(FontIconDialog, { props: defaultProps });
    wrapper.vm.filter = '';
    wrapper.vm.selectedCategory = 'navigation';
    await nextTick();
    expect(wrapper.vm.filteredIcons.every((i: any) => i.categories.includes('navigation'))).toBe(true);
    expect(wrapper.vm.filteredIcons).toHaveLength(2);
  });

  it('returns all icons when no filter and no category selected', async () => {
    const wrapper = mountWithPlugins(FontIconDialog, { props: defaultProps });
    wrapper.vm.selectedCategory = '';
    await nextTick();
    expect(wrapper.vm.filteredIcons).toHaveLength(8);
  });

  it('search filter is case-insensitive', async () => {
    const wrapper = mountWithPlugins(FontIconDialog, { props: defaultProps });
    wrapper.vm.selectedCategory = '';
    await nextTick();
    wrapper.vm.filter = 'HOME';
    await nextTick();
    expect(wrapper.vm.filteredIcons).toHaveLength(1);
    expect(wrapper.vm.filteredIcons[0].name).toBe('home');
  });

  it('clears filter when category changes', async () => {
    const wrapper = mountWithPlugins(FontIconDialog, { props: defaultProps });
    // Ensure we start with no category so filter isn't cleared
    wrapper.vm.selectedCategory = '';
    await nextTick();
    wrapper.vm.filter = 'test';
    await nextTick();
    expect(wrapper.vm.filter).toBe('test');
    // Now change category — watcher should clear filter
    wrapper.vm.selectedCategory = 'social';
    await nextTick();
    expect(wrapper.vm.filter).toBe('');
  });

  it('syncs show prop with computed show', () => {
    const wrapper = mountWithPlugins(FontIconDialog, { props: defaultProps });
    expect(wrapper.vm.show).toBe(true);
  });

  it('emits update:show when show computed is set', async () => {
    const wrapper = mountWithPlugins(FontIconDialog, { props: defaultProps });
    wrapper.vm.show = false;
    await nextTick();
    expect(wrapper.emitted('update:show')).toBeTruthy();
    expect(wrapper.emitted('update:show')![0]).toEqual([false]);
  });

  it('selects an icon via selectIcon', () => {
    const wrapper = mountWithPlugins(FontIconDialog, { props: defaultProps });
    (wrapper.vm).selectIcon('star');
    expect(wrapper.vm.selectedIcon).toBe('star');
  });

  it('syncs selectedIcon from modelValue prop', () => {
    const wrapper = mountWithPlugins(FontIconDialog, {
      props: { ...defaultProps, modelValue: 'home' },
    });
    expect(wrapper.vm.selectedIcon).toBe('home');
  });
});
