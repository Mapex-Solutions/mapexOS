import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { mountWithPlugins } from '@src/test/helpers';
import GroupSelectorDrawer from './GroupSelectorDrawer.vue';

vi.mock('@services/mapex', () => ({
  apis: {
    mapexOS: {
      groups: {
        list: vi.fn().mockResolvedValue({
          items: [],
          pagination: { totalPages: 1, totalItems: 0 },
        }),
      },
    },
  },
}));

vi.mock('@utils/error', () => ({
  handleApiError: vi.fn(),
}));

describe('GroupSelectorDrawer', () => {
  const defaultProps = {
    modelValue: true,
  };

  let addSpy: ReturnType<typeof vi.spyOn>;
  let removeSpy: ReturnType<typeof vi.spyOn>;

  beforeEach(() => {
    addSpy = vi.spyOn(window, 'addEventListener');
    removeSpy = vi.spyOn(window, 'removeEventListener');
  });

  afterEach(() => {
    addSpy.mockRestore();
    removeSpy.mockRestore();
  });

  it('renders without errors', () => {
    const wrapper = mountWithPlugins(GroupSelectorDrawer, { props: defaultProps });
    expect(wrapper.exists()).toBe(true);
  });

  it('computes showDialog from modelValue', () => {
    const wrapper = mountWithPlugins(GroupSelectorDrawer, { props: defaultProps });
    expect((wrapper.vm).showDialog).toBe(true);
  });

  it('sets showDialog to false when modelValue is false', () => {
    const wrapper = mountWithPlugins(GroupSelectorDrawer, {
      props: { ...defaultProps, modelValue: false },
    });
    expect((wrapper.vm).showDialog).toBe(false);
  });

  it('starts with loading as false before async fetch completes', () => {
    const wrapper = mountWithPlugins(GroupSelectorDrawer, { props: defaultProps });
    // Loading is set to true inside the async fetchGroups, but the initial ref value is false
    expect((wrapper.vm).loading).toBe(false);
  });

  it('starts with empty groups array', () => {
    const wrapper = mountWithPlugins(GroupSelectorDrawer, { props: defaultProps });
    expect((wrapper.vm).groups).toEqual([]);
  });

  it('initializes filter state correctly', () => {
    const wrapper = mountWithPlugins(GroupSelectorDrawer, { props: defaultProps });
    expect((wrapper.vm).filters.name).toBeUndefined();
    expect((wrapper.vm).filters.enabled).toBeUndefined();
    expect((wrapper.vm).filters.isTemplate).toBeUndefined();
  });

  it('computes statusOptions with 3 entries', () => {
    const wrapper = mountWithPlugins(GroupSelectorDrawer, { props: defaultProps });
    expect((wrapper.vm).statusOptions).toHaveLength(3);
  });

  it('computes templateOptions with 3 entries', () => {
    const wrapper = mountWithPlugins(GroupSelectorDrawer, { props: defaultProps });
    expect((wrapper.vm).templateOptions).toHaveLength(3);
  });

  it('defaults selectedGroupId to null', () => {
    const wrapper = mountWithPlugins(GroupSelectorDrawer, { props: defaultProps });
    expect(wrapper.props('selectedGroupId')).toBeNull();
  });

  it('emits select and closes on selectGroup', () => {
    const wrapper = mountWithPlugins(GroupSelectorDrawer, { props: defaultProps });
    const mockGroup = { id: 'g1', name: 'Group 1' };
    (wrapper.vm).selectGroup(mockGroup);
    expect(wrapper.emitted('select')).toBeTruthy();
    expect(wrapper.emitted('select')![0]).toEqual([mockGroup]);
    expect(wrapper.emitted('update:modelValue')![0]).toEqual([false]);
  });

  it('isSelected returns true when ID matches selectedGroupId', () => {
    const wrapper = mountWithPlugins(GroupSelectorDrawer, {
      props: { ...defaultProps, selectedGroupId: 'g1' },
    });
    expect((wrapper.vm).isSelected({ id: 'g1' })).toBe(true);
    expect((wrapper.vm).isSelected({ id: 'g2' })).toBe(false);
  });

  it('emits cancel on handleCancel', () => {
    const wrapper = mountWithPlugins(GroupSelectorDrawer, { props: defaultProps });
    (wrapper.vm).handleCancel();
    expect(wrapper.emitted('cancel')).toBeTruthy();
    expect(wrapper.emitted('update:modelValue')![0]).toEqual([false]);
  });

  it('getMemberCount returns correct text', () => {
    const wrapper = mountWithPlugins(GroupSelectorDrawer, { props: defaultProps });
    expect((wrapper.vm).getMemberCount({ memberCount: 1 })).toBe('1 member');
    expect((wrapper.vm).getMemberCount({ memberCount: 5 })).toBe('5 members');
    expect((wrapper.vm).getMemberCount({})).toBe('0 members');
  });

  it('registers ESC key handler on mount', () => {
    mountWithPlugins(GroupSelectorDrawer, { props: defaultProps });
    const keydownCalls = addSpy.mock.calls.filter(([event]: [string, ...unknown[]]) => event === 'keydown');
    expect(keydownCalls.length).toBeGreaterThanOrEqual(1);
  });

  it('handles ESC key when drawer is open', () => {
    const wrapper = mountWithPlugins(GroupSelectorDrawer, { props: defaultProps });
    const escEvent = new KeyboardEvent('keydown', { key: 'Escape' });
    window.dispatchEvent(escEvent);
    expect(wrapper.emitted('cancel')).toBeTruthy();
    expect(wrapper.emitted('update:modelValue')![0]).toEqual([false]);
  });
});
