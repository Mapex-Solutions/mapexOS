import { describe, it, expect, vi } from 'vitest';
import { mountWithPlugins } from '@src/test/helpers';
import ListHeaderMenu from './ListHeaderMenu.vue';

vi.mock('@composables/i18n/components/headers', () => ({
  useListHeaderMenuTranslations: () => ({
    itemsPerPage: { value: 'Items per page' },
    visibleColumns: { value: 'Visible columns' },
    filtered: { value: '(filtered)' },
  }),
}));

function factory(overrides: Record<string, unknown> = {}) {
  return mountWithPlugins(ListHeaderMenu, {
    props: {
      itemsCount: 50,
      itemLabel: 'Asset',
      itemsPerPage: 10,
      ...overrides,
    },
  });
}

describe('ListHeaderMenu', () => {
  it('renders without errors', () => {
    const wrapper = factory();
    expect(wrapper.exists()).toBe(true);
  });

  it('buttonLabel shows "X of Y ITEMS" when perPage < total', () => {
    const wrapper = factory({ itemsCount: 50, itemsPerPage: 10, itemLabel: 'Asset' });
    expect(wrapper.vm.buttonLabel).toBe('10 OF 50 ASSETS');
  });

  it('buttonLabel shows just count when perPage >= total', () => {
    const wrapper = factory({ itemsCount: 5, itemsPerPage: 10, itemLabel: 'Asset' });
    expect(wrapper.vm.buttonLabel).toBe('5 ASSETS');
  });

  it('buttonLabel uses singular label when count is 1', () => {
    const wrapper = factory({ itemsCount: 1, itemsPerPage: 10, itemLabel: 'Asset' });
    expect(wrapper.vm.buttonLabel).toBe('1 ASSET');
  });

  it('buttonLabel uses itemLabelPlural when provided', () => {
    const wrapper = factory({ itemsCount: 3, itemsPerPage: 10, itemLabel: 'Child', itemLabelPlural: 'Children' });
    expect(wrapper.vm.buttonLabel).toBe('3 CHILDREN');
  });

  it('buttonLabel appends filtered suffix when filtered is true', () => {
    const wrapper = factory({ itemsCount: 5, itemsPerPage: 10, itemLabel: 'Asset', filtered: true });
    expect(wrapper.vm.buttonLabel).toContain('(filtered)');
  });

  it('toggleColumn updates column visibility and emits', async () => {
    const columns = [
      { key: 'name', label: 'Name', visible: true },
      { key: 'status', label: 'Status', visible: true },
    ];
    const wrapper = factory({ columns });
    wrapper.vm.toggleColumn('status', false);
    await wrapper.vm.$nextTick();

    expect(wrapper.vm.localColumns.find((c: any) => c.key === 'status')?.visible).toBe(false);
    const emitted = wrapper.emitted('update:columns');
    expect(emitted).toBeTruthy();
  });

  it('toggleColumn does nothing for non-existent key', () => {
    const columns = [{ key: 'name', label: 'Name', visible: true }];
    const wrapper = factory({ columns });
    wrapper.vm.toggleColumn('nonexistent', false);
    expect(wrapper.vm.localColumns[0].visible).toBe(true);
  });

  it('defaults showItemsPerPage to true', () => {
    const wrapper = factory();
    expect(wrapper.props('showItemsPerPage')).toBe(true);
  });

  it('defaults showColumnVisibility to true', () => {
    const wrapper = factory();
    expect(wrapper.props('showColumnVisibility')).toBe(true);
  });

  it('defaults filtered to false', () => {
    const wrapper = factory();
    expect(wrapper.props('filtered')).toBe(false);
  });
});
