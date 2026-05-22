import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { mountWithPlugins } from '@src/test/helpers';
import UserMultiSelectorDrawer from './UserMultiSelectorDrawer.vue';
import type { UserSelectorItem } from './interfaces';

/**
 * Mock API service
 */
const mockUsersList = vi.fn();

vi.mock('@services/mapex', () => ({
  apis: {
    mapexOS: {
      users: {
        list: (...args: any[]) => mockUsersList(...args),
      },
    },
  },
}));

vi.mock('@utils/error', () => ({
  handleApiError: vi.fn(),
}));

/**
 * Helper to create a mock UserSelectorItem
 */
function makeUser(overrides: Partial<UserSelectorItem> = {}): UserSelectorItem {
  return {
    id: 'user-1',
    firstName: 'John',
    lastName: 'Doe',
    email: 'john.doe@example.com',
    ...overrides,
  };
}

describe('UserMultiSelectorDrawer', () => {
  const defaultProps = {
    modelValue: true,
    excludeUserIds: [],
    selectedUserIds: [],
  };

  let addSpy: ReturnType<typeof vi.spyOn>;
  let removeSpy: ReturnType<typeof vi.spyOn>;

  beforeEach(() => {
    addSpy = vi.spyOn(window, 'addEventListener');
    removeSpy = vi.spyOn(window, 'removeEventListener');

    mockUsersList.mockResolvedValue({
      items: [
        { id: 'user-1', firstName: 'John', lastName: 'Doe', email: 'john@test.com' },
        { id: 'user-2', firstName: 'Jane', lastName: 'Smith', email: 'jane@test.com' },
      ],
      pagination: { totalPages: 1, totalItems: 2 },
    });
  });

  afterEach(() => {
    addSpy.mockRestore();
    removeSpy.mockRestore();
  });

  it('renders without errors', () => {
    const wrapper = mountWithPlugins(UserMultiSelectorDrawer, { props: defaultProps });
    expect(wrapper.exists()).toBe(true);
  });

  it('receives props correctly', () => {
    const wrapper = mountWithPlugins(UserMultiSelectorDrawer, { props: defaultProps });
    expect(wrapper.props('modelValue')).toBe(true);
  });

  describe('user selection', () => {
    it('toggleUserSelection adds user to selection', () => {
      const wrapper = mountWithPlugins(UserMultiSelectorDrawer, { props: defaultProps });
      const vm = wrapper.vm;
      const user = makeUser({ id: 'user-1' });

      vm.toggleUserSelection(user);

      expect(vm.selectedUsers).toHaveLength(1);
      expect(vm.selectedUsers[0].id).toBe('user-1');
    });

    it('toggleUserSelection removes user when already selected', () => {
      const wrapper = mountWithPlugins(UserMultiSelectorDrawer, { props: defaultProps });
      const vm = wrapper.vm;
      const user = makeUser({ id: 'user-1' });

      vm.toggleUserSelection(user);
      vm.toggleUserSelection(user);

      expect(vm.selectedUsers).toHaveLength(0);
    });

    it('isSelected returns true for selected users', () => {
      const wrapper = mountWithPlugins(UserMultiSelectorDrawer, { props: defaultProps });
      const vm = wrapper.vm;
      const user = makeUser({ id: 'user-1' });

      vm.toggleUserSelection(user);

      expect(vm.isSelected(user)).toBe(true);
    });

    it('isSelected returns false for non-selected users', () => {
      const wrapper = mountWithPlugins(UserMultiSelectorDrawer, { props: defaultProps });
      const vm = wrapper.vm;
      const user = makeUser({ id: 'user-1' });

      expect(vm.isSelected(user)).toBe(false);
    });

    it('removeFromSelection removes specific user', () => {
      const wrapper = mountWithPlugins(UserMultiSelectorDrawer, { props: defaultProps });
      const vm = wrapper.vm;

      vm.toggleUserSelection(makeUser({ id: 'user-1' }));
      vm.toggleUserSelection(makeUser({ id: 'user-2' }));

      vm.removeFromSelection(makeUser({ id: 'user-1' }));

      expect(vm.selectedUsers).toHaveLength(1);
      expect(vm.selectedUsers[0].id).toBe('user-2');
    });

    it('clearSelection removes all selected users', () => {
      const wrapper = mountWithPlugins(UserMultiSelectorDrawer, { props: defaultProps });
      const vm = wrapper.vm;

      vm.toggleUserSelection(makeUser({ id: 'user-1' }));
      vm.toggleUserSelection(makeUser({ id: 'user-2' }));

      vm.clearSelection();

      expect(vm.selectedUsers).toHaveLength(0);
    });
  });

  describe('confirm and cancel', () => {
    it('handleConfirm emits confirm with selected users and closes', () => {
      const wrapper = mountWithPlugins(UserMultiSelectorDrawer, { props: defaultProps });
      const vm = wrapper.vm;

      vm.toggleUserSelection(makeUser({ id: 'user-1' }));
      vm.handleConfirm();

      expect(wrapper.emitted('confirm')).toBeTruthy();
      expect(wrapper.emitted('confirm')![0]![0]).toHaveLength(1);
      expect(wrapper.emitted('update:modelValue')![0]).toEqual([false]);
    });

    it('handleCancel emits cancel and closes', () => {
      const wrapper = mountWithPlugins(UserMultiSelectorDrawer, { props: defaultProps });
      const vm = wrapper.vm;

      vm.handleCancel();

      expect(wrapper.emitted('cancel')).toBeTruthy();
      expect(wrapper.emitted('update:modelValue')![0]).toEqual([false]);
    });

    it('canConfirm is false when no users selected', () => {
      const wrapper = mountWithPlugins(UserMultiSelectorDrawer, { props: defaultProps });
      const vm = wrapper.vm;

      expect(vm.canConfirm).toBe(false);
    });

    it('canConfirm is true when users are selected', () => {
      const wrapper = mountWithPlugins(UserMultiSelectorDrawer, { props: defaultProps });
      const vm = wrapper.vm;

      vm.toggleUserSelection(makeUser({ id: 'user-1' }));

      expect(vm.canConfirm).toBe(true);
    });
  });

  describe('display helpers', () => {
    it('getInitials returns first letters of first and last name', () => {
      const wrapper = mountWithPlugins(UserMultiSelectorDrawer, { props: defaultProps });
      const vm = wrapper.vm;

      expect(vm.getInitials(makeUser({ firstName: 'John', lastName: 'Doe' }))).toBe('JD');
    });

    it('getInitials falls back to email when no name', () => {
      const wrapper = mountWithPlugins(UserMultiSelectorDrawer, { props: defaultProps });
      const vm = wrapper.vm;

      expect(vm.getInitials(makeUser({ firstName: '', lastName: '', email: 'test@mail.com' }))).toBe('T');
    });

    it('getDisplayName returns full name', () => {
      const wrapper = mountWithPlugins(UserMultiSelectorDrawer, { props: defaultProps });
      const vm = wrapper.vm;

      expect(vm.getDisplayName(makeUser({ firstName: 'John', lastName: 'Doe' }))).toBe('John Doe');
    });

    it('getDisplayName falls back to email when no name', () => {
      const wrapper = mountWithPlugins(UserMultiSelectorDrawer, { props: defaultProps });
      const vm = wrapper.vm;

      expect(vm.getDisplayName(makeUser({ firstName: '', lastName: '', email: 'test@mail.com' }))).toBe('test@mail.com');
    });
  });

  describe('excluded users', () => {
    it('displayUsers filters out excluded user IDs', () => {
      const wrapper = mountWithPlugins(UserMultiSelectorDrawer, {
        props: { ...defaultProps, excludeUserIds: ['user-1'] },
      });
      const vm = wrapper.vm;

      // Manually set users as if API returned them
      vm.users = [
        makeUser({ id: 'user-1' }),
        makeUser({ id: 'user-2' }),
      ];

      expect(vm.displayUsers).toHaveLength(1);
      expect(vm.displayUsers[0].id).toBe('user-2');
    });
  });

  describe('keyboard handling', () => {
    it('registers ESC key handler on mount', () => {
      mountWithPlugins(UserMultiSelectorDrawer, { props: defaultProps });
      const keydownCalls = addSpy.mock.calls.filter(([event]: [string, ...unknown[]]) => event === 'keydown');
      expect(keydownCalls.length).toBeGreaterThanOrEqual(1);
    });

    it('handles ESC key to cancel when drawer is open', () => {
      const wrapper = mountWithPlugins(UserMultiSelectorDrawer, { props: defaultProps });

      const escEvent = new KeyboardEvent('keydown', { key: 'Escape' });
      window.dispatchEvent(escEvent);

      expect(wrapper.emitted('cancel')).toBeTruthy();
    });
  });
});
