import { describe, it, expect, vi } from 'vitest';
import { mountWithPlugins } from '@src/test/helpers';
import TriggerSelectorDialog from './TriggerSelectorDialog.vue';

vi.mock('@utils/translation', () => ({
  useTS: () => (key: string) => key,
}));

vi.mock('@components/dialogs/common/genericSelectorDialog', () => ({
  GenericSelectorDialog: { name: 'GenericSelectorDialog', template: '<div />' },
}));

vi.mock('@services/mapex', () => ({
  apis: {
    triggers: {
      trigger: {
        list: vi.fn().mockResolvedValue({ items: [] }),
      },
    },
  },
}));

vi.mock('@utils/error', () => ({
  handleApiError: vi.fn(),
}));

const BASE_PROPS = {
  modelValue: false,
};

describe('TriggerSelectorDialog', () => {
  it('renders without errors', () => {
    const wrapper = mountWithPlugins(TriggerSelectorDialog, { props: BASE_PROPS });
    expect(wrapper.exists()).toBe(true);
  });

  it('starts with empty triggers', () => {
    const wrapper = mountWithPlugins(TriggerSelectorDialog, { props: BASE_PROPS });
    expect(wrapper.vm.triggers).toEqual([]);
  });

  it('starts with loading as false', () => {
    const wrapper = mountWithPlugins(TriggerSelectorDialog, { props: BASE_PROPS });
    expect(wrapper.vm.loading).toBe(false);
  });

  it('computes selectedIds from selectedTriggerId', () => {
    const wrapper = mountWithPlugins(TriggerSelectorDialog, {
      props: { ...BASE_PROPS, selectedTriggerId: 'trig-1' },
    });
    expect(wrapper.vm.selectedIds).toEqual(['trig-1']);
  });

  it('computes selectedIds as empty when no selectedTriggerId', () => {
    const wrapper = mountWithPlugins(TriggerSelectorDialog, { props: BASE_PROPS });
    expect(wrapper.vm.selectedIds).toEqual([]);
  });

  it('computes filteredTriggers with empty search', () => {
    const wrapper = mountWithPlugins(TriggerSelectorDialog, { props: BASE_PROPS });
    expect(wrapper.vm.filteredTriggers).toEqual([]);
  });

  it('handles search query update', () => {
    const wrapper = mountWithPlugins(TriggerSelectorDialog, { props: BASE_PROPS });
    wrapper.vm.handleSearch('email');
    expect(wrapper.vm.searchQuery).toBe('email');
  });

  it('handles select by emitting trigger', () => {
    const wrapper = mountWithPlugins(TriggerSelectorDialog, { props: BASE_PROPS });
    const trigger = { id: 'trig-1', name: 'Email Trigger', category: 'email' };
    wrapper.vm.handleSelect([trigger]);
    const emitted = wrapper.emitted('select')!;
    expect(emitted[0]![0]).toEqual(trigger);
  });

  it('getCategoryIcon returns correct icon for known categories', () => {
    const wrapper = mountWithPlugins(TriggerSelectorDialog, { props: BASE_PROPS });
    expect(wrapper.vm.getCategoryIcon('email')).toBe('email');
    expect(wrapper.vm.getCategoryIcon('slack')).toBe('chat');
    expect(wrapper.vm.getCategoryIcon('http')).toBe('http');
  });

  it('getCategoryIcon returns default for unknown category', () => {
    const wrapper = mountWithPlugins(TriggerSelectorDialog, { props: BASE_PROPS });
    expect(wrapper.vm.getCategoryIcon('unknown')).toBe('notifications');
  });

  it('getCategoryColor returns correct color for known categories', () => {
    const wrapper = mountWithPlugins(TriggerSelectorDialog, { props: BASE_PROPS });
    expect(wrapper.vm.getCategoryColor('email')).toBe('blue');
    expect(wrapper.vm.getCategoryColor('mqtt')).toBe('orange');
  });

  it('getCategoryColor returns primary for unknown category', () => {
    const wrapper = mountWithPlugins(TriggerSelectorDialog, { props: BASE_PROPS });
    expect(wrapper.vm.getCategoryColor('unknown')).toBe('primary');
  });

  it('computes categoryOptions with all filter as first entry', () => {
    const wrapper = mountWithPlugins(TriggerSelectorDialog, { props: BASE_PROPS });
    expect(wrapper.vm.categoryOptions).toHaveLength(1);
    expect(wrapper.vm.categoryOptions[0].value).toBeUndefined();
  });
});
