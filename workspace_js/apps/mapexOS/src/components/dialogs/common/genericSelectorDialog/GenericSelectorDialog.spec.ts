import { describe, it, expect, vi } from 'vitest';
import { mountWithPlugins } from '@src/test/helpers';
import GenericSelectorDialog from './GenericSelectorDialog.vue';

vi.mock('@utils/translation', () => ({
  useTS: () => (key: string) => key,
}));

vi.mock('./constants', () => ({
  DEFAULT_DIALOG_WIDTH: 600,
  DEFAULT_SEARCH_DEBOUNCE_MS: 300,
  DEFAULT_SCROLL_THRESHOLD: 0.9,
  DEFAULT_ITEM_KEY: 'id',
  DEFAULT_EMPTY_ICON: 'inbox',
  DEFAULT_EMPTY_TEXT: 'No items',
  DEFAULT_LOADING_TEXT: 'Loading...',
  DEFAULT_CONFIRM_LABEL: 'Confirm',
  DEFAULT_CANCEL_LABEL: 'Cancel',
  DEFAULT_ITEM_NOUN_SINGULAR: 'item',
  DEFAULT_ITEM_NOUN_PLURAL: 'items',
  DEFAULT_INFO_BANNER: { icon: 'info', bgClass: 'bg-teal-1', textClass: 'text-teal-9', iconColor: 'teal-6' },
  DEFAULT_ACTIVE_ITEM_STYLE: { backgroundColor: 'rgba(0,150,136,0.08)', borderColor: 'var(--q-teal-6)' },
}));

const BASE_PROPS = {
  modelValue: true,
  title: 'Select Item',
  items: [
    { id: '1', name: 'Item A' },
    { id: '2', name: 'Item B' },
  ],
};

describe('GenericSelectorDialog', () => {
  it('renders without errors', () => {
    const wrapper = mountWithPlugins(GenericSelectorDialog, { props: BASE_PROPS });
    expect(wrapper.exists()).toBe(true);
  });

  it('computes showDialog from modelValue', () => {
    const wrapper = mountWithPlugins(GenericSelectorDialog, { props: BASE_PROPS });
    expect(wrapper.vm.showDialog).toBe(true);
  });

  it('computes canConfirm as false with no selection', () => {
    const wrapper = mountWithPlugins(GenericSelectorDialog, { props: BASE_PROPS });
    expect(wrapper.vm.canConfirm).toBe(false);
  });

  it('computes footerCountText correctly', () => {
    const wrapper = mountWithPlugins(GenericSelectorDialog, { props: BASE_PROPS });
    expect(wrapper.vm.footerCountText).toBe('2 items');
  });

  it('computes footerCountText for single item', () => {
    const wrapper = mountWithPlugins(GenericSelectorDialog, {
      props: { ...BASE_PROPS, items: [{ id: '1', name: 'Item A' }] },
    });
    expect(wrapper.vm.footerCountText).toBe('1 item');
  });

  it('checks isSelected returns false for unselected item', () => {
    const wrapper = mountWithPlugins(GenericSelectorDialog, { props: BASE_PROPS });
    expect(wrapper.vm.isSelected({ id: '1' })).toBe(false);
  });

  it('handles single-select toggleItem by emitting select and closing', () => {
    const wrapper = mountWithPlugins(GenericSelectorDialog, { props: BASE_PROPS });
    wrapper.vm.toggleItem({ id: '1', name: 'Item A' });
    const emitted = wrapper.emitted('select')!;
    expect(emitted[0]![0]).toEqual([{ id: '1', name: 'Item A' }]);
  });

  it('handles multi-select toggle adding/removing items', () => {
    const wrapper = mountWithPlugins(GenericSelectorDialog, {
      props: { ...BASE_PROPS, multiSelect: true },
    });
    wrapper.vm.toggleItem({ id: '1', name: 'Item A' });
    expect(wrapper.vm.selectedItems).toHaveLength(1);
    wrapper.vm.toggleItem({ id: '1', name: 'Item A' });
    expect(wrapper.vm.selectedItems).toHaveLength(0);
  });

  it('confirms multi-select selection', () => {
    const wrapper = mountWithPlugins(GenericSelectorDialog, {
      props: { ...BASE_PROPS, multiSelect: true },
    });
    wrapper.vm.toggleItem({ id: '1', name: 'Item A' });
    wrapper.vm.confirmSelection();
    const emitted = wrapper.emitted('select')!;
    expect(emitted[0]![0]).toHaveLength(1);
  });

  it('emits cancel on handleCancel', () => {
    const wrapper = mountWithPlugins(GenericSelectorDialog, { props: BASE_PROPS });
    wrapper.vm.handleCancel();
    expect(wrapper.emitted('cancel')).toBeTruthy();
  });

  it('handles search input and emits search', () => {
    vi.useFakeTimers();
    const wrapper = mountWithPlugins(GenericSelectorDialog, { props: BASE_PROPS });
    wrapper.vm.handleSearchInput('test');
    expect(wrapper.vm.searchQuery).toBe('test');
    vi.advanceTimersByTime(400);
    const emitted = wrapper.emitted('search')!;
    expect(emitted[0]![0]).toBe('test');
    vi.useRealTimers();
  });

  it('handles search clear bypassing debounce', () => {
    const wrapper = mountWithPlugins(GenericSelectorDialog, { props: BASE_PROPS });
    wrapper.vm.handleSearchClear();
    expect(wrapper.vm.searchQuery).toBe('');
    const emitted = wrapper.emitted('search')!;
    expect(emitted[0]![0]).toBe('');
  });

  it('computes resolvedBanner as null when no infoBanner', () => {
    const wrapper = mountWithPlugins(GenericSelectorDialog, { props: BASE_PROPS });
    expect(wrapper.vm.resolvedBanner).toBeNull();
  });

  it('computes resolvedBanner with defaults when infoBanner provided', () => {
    const wrapper = mountWithPlugins(GenericSelectorDialog, {
      props: { ...BASE_PROPS, infoBanner: { text: 'Select items' } },
    });
    expect(wrapper.vm.resolvedBanner).toBeTruthy();
    expect(wrapper.vm.resolvedBanner!.text).toBe('Select items');
  });
});
