import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { mountWithPlugins } from '@src/test/helpers';
import GroupDetailsDrawer from './GroupDetailsDrawer.vue';

vi.mock('quasar', () => ({
  date: {
    formatDate: vi.fn(() => 'Jan 01, 2024 12:00'),
  },
}));

vi.mock('vue-router', () => ({
  useRouter: () => ({
    push: vi.fn(),
  }),
}));

vi.mock('@composables/i18n', () => ({
  useGroupsTranslations: () => new Proxy({}, {
    get: (_t: any, prop: string) => {
      if (prop === 'value') return prop;
      return new Proxy({ value: String(prop) }, {
        get: (_t2: any, p2: string) => {
          if (p2 === 'value') return String(prop);
          return new Proxy({ value: String(p2) }, {
            get: (_t3: any, p3: string) => {
              if (p3 === 'value') return String(p2);
              return { value: String(p3) };
            },
          });
        },
      });
    },
  }),
}));

vi.mock('@composables/useLogger', () => ({
  useLogger: () => ({
    debug: vi.fn(),
    info: vi.fn(),
    warn: vi.fn(),
    error: vi.fn(),
  }),
}));

vi.mock('@services/mapex', () => ({
  apis: {
    mapexOS: {
      groups: {
        getById: vi.fn().mockResolvedValue({
          id: 'grp-1',
          name: 'Test Group',
          enabled: true,
          orgId: 'org-1',
          membersCount: 0,
        }),
        getMembers: vi.fn().mockResolvedValue({ items: [] }),
      },
      roles: {
        list: vi.fn().mockResolvedValue({ items: [] }),
      },
      users: {
        getById: vi.fn(),
      },
    },
  },
}));

vi.mock('@stores/organization', () => ({
  useOrganizationStore: () => ({
    flatList: [{ id: 'org-1', name: 'Test Org' }],
  }),
}));

vi.mock('@utils/alert', () => ({
  notifyFail: vi.fn(),
}));

describe('GroupDetailsDrawer', () => {
  const defaultProps = {
    modelValue: true,
    groupId: 'grp-1',
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
    const wrapper = mountWithPlugins(GroupDetailsDrawer, { props: defaultProps });
    expect(wrapper.exists()).toBe(true);
  });

  it('starts with loading state', () => {
    const wrapper = mountWithPlugins(GroupDetailsDrawer, { props: defaultProps });
    expect((wrapper.vm).loading).toBe(true);
  });

  it('starts with group as null', () => {
    const wrapper = mountWithPlugins(GroupDetailsDrawer, { props: defaultProps });
    expect((wrapper.vm).group).toBeNull();
  });

  it('starts with empty membersList', () => {
    const wrapper = mountWithPlugins(GroupDetailsDrawer, { props: defaultProps });
    expect((wrapper.vm).membersList).toEqual([]);
  });

  it('starts with empty rolesList', () => {
    const wrapper = mountWithPlugins(GroupDetailsDrawer, { props: defaultProps });
    expect((wrapper.vm).rolesList).toEqual([]);
  });

  it('computes displayedRoles as empty when rolesList is empty', () => {
    const wrapper = mountWithPlugins(GroupDetailsDrawer, { props: defaultProps });
    expect((wrapper.vm).displayedRoles).toEqual([]);
  });

  it('computes displayedMembers as empty when membersList is empty', () => {
    const wrapper = mountWithPlugins(GroupDetailsDrawer, { props: defaultProps });
    expect((wrapper.vm).displayedMembers).toEqual([]);
  });

  it('registers ESC key handler on mount', () => {
    mountWithPlugins(GroupDetailsDrawer, { props: defaultProps });
    const keydownCalls = addSpy.mock.calls.filter(([event]: [string, ...unknown[]]) => event === 'keydown');
    expect(keydownCalls.length).toBeGreaterThanOrEqual(1);
  });

  it('handles ESC key when drawer is open', () => {
    const wrapper = mountWithPlugins(GroupDetailsDrawer, { props: defaultProps });
    const escEvent = new KeyboardEvent('keydown', { key: 'Escape' });
    window.dispatchEvent(escEvent);
    expect(wrapper.emitted('update:modelValue')).toBeTruthy();
    expect(wrapper.emitted('update:modelValue')![0]).toEqual([false]);
  });

  it('ignores ESC key when drawer is closed', () => {
    const wrapper = mountWithPlugins(GroupDetailsDrawer, {
      props: { ...defaultProps, modelValue: false },
    });
    const escEvent = new KeyboardEvent('keydown', { key: 'Escape' });
    window.dispatchEvent(escEvent);
    expect(wrapper.emitted('update:modelValue')).toBeFalsy();
  });

  it('emits update:modelValue(false) on close', () => {
    const wrapper = mountWithPlugins(GroupDetailsDrawer, { props: defaultProps });
    (wrapper.vm).close();
    expect(wrapper.emitted('update:modelValue')![0]).toEqual([false]);
  });

  it('does not emit edit when group is null', () => {
    const wrapper = mountWithPlugins(GroupDetailsDrawer, { props: defaultProps });
    (wrapper.vm).handleEdit();
    expect(wrapper.emitted('edit')).toBeFalsy();
  });

  it('getInitials returns correct initials', () => {
    const wrapper = mountWithPlugins(GroupDetailsDrawer, { props: defaultProps });
    expect((wrapper.vm).getInitials('John Doe')).toBe('JD');
    expect((wrapper.vm).getInitials('Jane')).toBe('JA');
    expect((wrapper.vm).getInitials('')).toBe('?');
  });
});
