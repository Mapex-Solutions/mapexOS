import { describe, it, expect, vi } from 'vitest';
import { mountWithPlugins } from '@src/test/helpers';
import UserSelectFilter from './UserSelectFilter.vue';

vi.mock('@services/mapex', () => ({
  apis: {
    mapexOS: {
      users: {
        getById: vi.fn().mockResolvedValue({
          id: 'user-1',
          firstName: 'John',
          lastName: 'Doe',
          email: 'john@test.com',
        }),
      },
    },
  },
}));

function factory(overrides: Record<string, unknown> = {}) {
  return mountWithPlugins(UserSelectFilter, {
    props: {
      modelValue: null,
      label: 'User',
      icon: 'person',
      ...overrides,
    },
  });
}

describe('UserSelectFilter', () => {
  it('renders without errors', () => {
    const wrapper = factory();
    expect(wrapper.exists()).toBe(true);
  });

  it('defaults clearable to true', () => {
    const wrapper = factory();
    expect(wrapper.props('clearable')).toBe(true);
  });

  it('defaults disabled to false', () => {
    const wrapper = factory();
    expect(wrapper.props('disabled')).toBe(false);
  });

  it('defaults placeholder to "Click to select..."', () => {
    const wrapper = factory();
    expect(wrapper.props('placeholder')).toBe('Click to select...');
  });

  it('getUserDisplayName returns full name', () => {
    const wrapper = factory();
    expect(wrapper.vm.getUserDisplayName({ firstName: 'John', lastName: 'Doe' })).toBe('John Doe');
  });

  it('getUserDisplayName returns first name only when no last name', () => {
    const wrapper = factory();
    expect(wrapper.vm.getUserDisplayName({ firstName: 'John' })).toBe('John');
  });

  it('getUserDisplayName falls back to email', () => {
    const wrapper = factory();
    expect(wrapper.vm.getUserDisplayName({ email: 'john@test.com' })).toBe('john@test.com');
  });

  it('getUserDisplayName returns "Unknown User" when no name or email', () => {
    const wrapper = factory();
    expect(wrapper.vm.getUserDisplayName({})).toBe('Unknown User');
  });

  it('openDrawer does nothing when disabled', () => {
    const wrapper = factory({ disabled: true });
    wrapper.vm.openDrawer();
    expect(wrapper.vm.showDrawer).toBe(false);
  });

  it('openDrawer sets showDrawer to true when enabled', () => {
    const wrapper = factory();
    wrapper.vm.openDrawer();
    expect(wrapper.vm.showDrawer).toBe(true);
  });

  it('clearSelection resets selectedUser and emits null', async () => {
    const wrapper = factory();
    wrapper.vm.selectedUser = { id: 'u1', name: 'Test' };
    wrapper.vm.clearSelection();
    await wrapper.vm.$nextTick();
    expect(wrapper.vm.selectedUser).toBeNull();
    const emitted = wrapper.emitted('update:modelValue');
    expect(emitted).toBeTruthy();
    expect(emitted![0]![0]).toBeNull();
  });

  it('handleUserSelect sets selectedUser and emits user id', async () => {
    const wrapper = factory();
    wrapper.vm.handleUserSelect({ id: 'u1', firstName: 'Jane', lastName: 'Doe', email: 'jane@test.com' });
    await wrapper.vm.$nextTick();
    expect(wrapper.vm.selectedUser).toEqual({ id: 'u1', name: 'Jane Doe', email: 'jane@test.com' });
    const emitted = wrapper.emitted('update:modelValue');
    expect(emitted).toBeTruthy();
    expect(emitted![0]![0]).toBe('u1');
  });
});
