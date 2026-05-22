import { describe, it, expect, vi } from 'vitest';
import { mountWithPlugins } from '@src/test/helpers';
import OrganizationSelectorDrawer from './OrganizationSelectorDrawer.vue';

vi.mock('@composables/organizations/useOrganizationTree', () => ({
  useOrganizationTree: () => ({
    treeNodes: { value: [] },
    loading: { value: false },
    filters: { value: { name: '', types: [], enabled: 'all' } },
  }),
}));

vi.mock('@utils/organization/icons', () => ({
  getOrganizationIcon: vi.fn(() => 'domain'),
  getOrganizationColor: vi.fn(() => 'primary'),
}));

describe('OrganizationSelectorDrawer', () => {
  const defaultProps = {
    modelValue: true,
    selectedOrganizationId: null,
  };

  it('renders without errors', () => {
    const wrapper = mountWithPlugins(OrganizationSelectorDrawer, { props: defaultProps });
    expect(wrapper.exists()).toBe(true);
  });

  it('syncs showDialog with modelValue prop', () => {
    const wrapper = mountWithPlugins(OrganizationSelectorDrawer, { props: defaultProps });
    expect(wrapper.vm.showDialog).toBe(true);
  });

  it('sets showDialog to false when modelValue is false', () => {
    const wrapper = mountWithPlugins(OrganizationSelectorDrawer, {
      props: { ...defaultProps, modelValue: false },
    });
    expect(wrapper.vm.showDialog).toBe(false);
  });

  it('emits update:modelValue when showDialog changes', async () => {
    const wrapper = mountWithPlugins(OrganizationSelectorDrawer, { props: defaultProps });
    wrapper.vm.showDialog = false;
    await wrapper.vm.$nextTick();
    expect(wrapper.emitted('update:modelValue')).toBeTruthy();
    expect(wrapper.emitted('update:modelValue')![0]).toEqual([false]);
  });

  it('initializes expanded as empty array', () => {
    const wrapper = mountWithPlugins(OrganizationSelectorDrawer, { props: defaultProps });
    expect(wrapper.vm.expanded).toEqual([]);
  });

  it('computes typeOptions with all organization types', () => {
    const wrapper = mountWithPlugins(OrganizationSelectorDrawer, { props: defaultProps });
    const typeOptions = wrapper.vm.typeOptions;
    expect(typeOptions).toHaveLength(7);
    expect(typeOptions[0].label).toBe('All Types');
  });

  it('computes enabledOptions', () => {
    const wrapper = mountWithPlugins(OrganizationSelectorDrawer, { props: defaultProps });
    const enabledOptions = wrapper.vm.enabledOptions;
    expect(enabledOptions).toHaveLength(3);
    expect(enabledOptions.map((o: any) => o.value)).toEqual(['all', 'active', 'inactive']);
  });

  it('identifies selected organization via isSelected', () => {
    const wrapper = mountWithPlugins(OrganizationSelectorDrawer, {
      props: { ...defaultProps, selectedOrganizationId: 'org-1' },
    });
    const vm = wrapper.vm;
    expect(vm.isSelected({ id: 'org-1' })).toBe(true);
    expect(vm.isSelected({ id: 'org-2' })).toBe(false);
  });

  it('emits select and closes on selectOrganization', async () => {
    const wrapper = mountWithPlugins(OrganizationSelectorDrawer, { props: defaultProps });
    const vm = wrapper.vm;
    const node = { id: 'org-1', name: 'Test Org', type: 'customer', enabled: true, pathKey: '/test' };
    vm.selectOrganization(node);
    await wrapper.vm.$nextTick();
    expect(wrapper.emitted('select')).toBeTruthy();
    expect(wrapper.emitted('select')![0]![0]).toMatchObject({ id: 'org-1', name: 'Test Org' });
  });

  it('emits cancel and closes on handleCancel', async () => {
    const wrapper = mountWithPlugins(OrganizationSelectorDrawer, { props: defaultProps });
    const vm = wrapper.vm;
    vm.handleCancel();
    await wrapper.vm.$nextTick();
    expect(wrapper.emitted('cancel')).toBeTruthy();
    expect(wrapper.emitted('update:modelValue')).toBeTruthy();
  });

  it('handles ESC key when drawer is open', () => {
    const wrapper = mountWithPlugins(OrganizationSelectorDrawer, { props: defaultProps });
    window.dispatchEvent(new KeyboardEvent('keydown', { key: 'Escape' }));
    expect(wrapper.emitted('cancel')).toBeTruthy();
  });

  it('does not handle ESC key when drawer is closed', () => {
    const wrapper = mountWithPlugins(OrganizationSelectorDrawer, {
      props: { ...defaultProps, modelValue: false },
    });
    window.dispatchEvent(new KeyboardEvent('keydown', { key: 'Escape' }));
    expect(wrapper.emitted('cancel')).toBeFalsy();
  });

  it('toggles node expansion', () => {
    const wrapper = mountWithPlugins(OrganizationSelectorDrawer, { props: defaultProps });
    const vm = wrapper.vm;
    const node = { id: 'org-1' };
    vm.toggleExpand(node);
    expect(vm.expanded).toContain('org-1');
    vm.toggleExpand(node);
    expect(vm.expanded).not.toContain('org-1');
  });
});
