import { describe, it, expect, vi } from 'vitest';
import { mountWithPlugins } from '@src/test/helpers';
import TriggerSelectorDrawer from './TriggerSelectorDrawer.vue';

// Mock API service
vi.mock('@services/mapex', () => ({
  apis: {
    triggers: {
      trigger: {
        list: vi.fn().mockResolvedValue({
          items: [
            { id: 't1', name: 'HTTP Trigger', description: 'HTTP endpoint', category: 'http' },
            { id: 't2', name: 'Email Trigger', description: 'Email receiver', category: 'email' },
            { id: 't3', name: 'Slack Trigger', description: 'Slack event', category: 'slack' },
          ],
        }),
      },
    },
  },
}));

// Mock error handler
vi.mock('@utils/error', () => ({
  handleApiError: vi.fn(),
}));

describe('TriggerSelectorDrawer', () => {
  const defaultProps = {
    modelValue: false,
    selectedTriggerId: null,
  };

  it('renders without errors', () => {
    const wrapper = mountWithPlugins(TriggerSelectorDrawer, {
      props: defaultProps,
    });
    expect(wrapper.exists()).toBe(true);
  });

  it('computes showDrawer from modelValue', () => {
    const wrapper = mountWithPlugins(TriggerSelectorDrawer, {
      props: { ...defaultProps, modelValue: true },
    });
    expect(wrapper.vm.showDrawer).toBe(true);
  });

  it('emits update:modelValue when showDrawer is set', () => {
    const wrapper = mountWithPlugins(TriggerSelectorDrawer, {
      props: defaultProps,
    });
    wrapper.vm.showDrawer = true;
    expect(wrapper.emitted('update:modelValue')).toBeTruthy();
    expect(wrapper.emitted('update:modelValue')![0]).toEqual([true]);
  });

  it('returns correct icon from getCategoryIcon', () => {
    const wrapper = mountWithPlugins(TriggerSelectorDrawer, {
      props: defaultProps,
    });
    expect(wrapper.vm.getCategoryIcon('email')).toBe('email');
    expect(wrapper.vm.getCategoryIcon('slack')).toBe('chat');
    expect(wrapper.vm.getCategoryIcon('http')).toBe('http');
    expect(wrapper.vm.getCategoryIcon('mqtt')).toBe('router');
    expect(wrapper.vm.getCategoryIcon('unknown')).toBe('notifications');
  });

  it('returns correct color from getCategoryColor', () => {
    const wrapper = mountWithPlugins(TriggerSelectorDrawer, {
      props: defaultProps,
    });
    expect(wrapper.vm.getCategoryColor('email')).toBe('blue');
    expect(wrapper.vm.getCategoryColor('slack')).toBe('purple');
    expect(wrapper.vm.getCategoryColor('teams')).toBe('indigo');
    expect(wrapper.vm.getCategoryColor('unknown')).toBe('primary');
  });

  it('isSelected returns true when trigger id matches', () => {
    const wrapper = mountWithPlugins(TriggerSelectorDrawer, {
      props: { ...defaultProps, selectedTriggerId: 't1' },
    });
    expect(wrapper.vm.isSelected({ id: 't1', name: 'Test' } as any)).toBe(true);
  });

  it('isSelected returns false when trigger id does not match', () => {
    const wrapper = mountWithPlugins(TriggerSelectorDrawer, {
      props: { ...defaultProps, selectedTriggerId: 't1' },
    });
    expect(wrapper.vm.isSelected({ id: 't2', name: 'Test' } as any)).toBe(false);
  });

  it('close emits update:modelValue false', () => {
    const wrapper = mountWithPlugins(TriggerSelectorDrawer, {
      props: defaultProps,
    });
    wrapper.vm.close();
    expect(wrapper.emitted('update:modelValue')![0]).toEqual([false]);
  });

  it('handleTriggerSelect emits select and closes', () => {
    const wrapper = mountWithPlugins(TriggerSelectorDrawer, {
      props: defaultProps,
    });
    const trigger = { id: 't1', name: 'HTTP Trigger' } as any;
    wrapper.vm.handleTriggerSelect(trigger);
    expect(wrapper.emitted('select')![0]).toEqual([trigger]);
    expect(wrapper.emitted('update:modelValue')).toBeTruthy();
  });

  it('fetches triggers and populates list', async () => {
    const wrapper = mountWithPlugins(TriggerSelectorDrawer, {
      props: { ...defaultProps, modelValue: true },
    });
    await wrapper.vm.fetchTriggers();
    expect(wrapper.vm.triggers).toHaveLength(3);
  });

  it('computes categories from loaded triggers', async () => {
    const wrapper = mountWithPlugins(TriggerSelectorDrawer, {
      props: { ...defaultProps, modelValue: true },
    });
    await wrapper.vm.fetchTriggers();
    expect(wrapper.vm.categories).toContain('all');
    expect(wrapper.vm.categories).toContain('http');
    expect(wrapper.vm.categories).toContain('email');
  });

  it('filteredTriggers filters by category', async () => {
    const wrapper = mountWithPlugins(TriggerSelectorDrawer, {
      props: { ...defaultProps, modelValue: true },
    });
    await wrapper.vm.fetchTriggers();
    wrapper.vm.selectedCategory = 'http';
    expect(wrapper.vm.filteredTriggers).toHaveLength(1);
    expect(wrapper.vm.filteredTriggers[0].name).toBe('HTTP Trigger');
  });

  it('filteredTriggers filters by search query', async () => {
    const wrapper = mountWithPlugins(TriggerSelectorDrawer, {
      props: { ...defaultProps, modelValue: true },
    });
    await wrapper.vm.fetchTriggers();
    wrapper.vm.searchQuery = 'email';
    expect(wrapper.vm.filteredTriggers.length).toBeGreaterThanOrEqual(1);
    expect(wrapper.vm.filteredTriggers.some((t: any) => t.name.toLowerCase().includes('email'))).toBe(true);
  });
});
