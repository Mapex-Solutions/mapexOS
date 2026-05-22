import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { mountWithPlugins } from '@src/test/helpers';
import CustomerDetailsDrawer from './CustomerDetailsDrawer.vue';

vi.mock('quasar', () => ({
  date: {
    formatDate: vi.fn(() => 'Jan 01, 2024 12:00'),
  },
}));

vi.mock('@composables/i18n', () => ({
  useCustomersTranslations: () => new Proxy({}, {
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
      organizations: {
        getById: vi.fn().mockResolvedValue({
          id: 'cust-1',
          name: 'Test Customer',
          enabled: true,
          type: 'customer',
        }),
      },
    },
  },
}));

vi.mock('@utils/alert', () => ({
  notifyFail: vi.fn(),
}));

describe('CustomerDetailsDrawer', () => {
  const defaultProps = {
    modelValue: true,
    customerId: 'cust-1',
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
    const wrapper = mountWithPlugins(CustomerDetailsDrawer, { props: defaultProps });
    expect(wrapper.exists()).toBe(true);
  });

  it('starts with loading state', () => {
    const wrapper = mountWithPlugins(CustomerDetailsDrawer, { props: defaultProps });
    expect((wrapper.vm).loading).toBe(true);
  });

  it('starts with customer as null', () => {
    const wrapper = mountWithPlugins(CustomerDetailsDrawer, { props: defaultProps });
    expect((wrapper.vm).customer).toBeNull();
  });

  it('computes hasAddress as false when customer is null', () => {
    const wrapper = mountWithPlugins(CustomerDetailsDrawer, { props: defaultProps });
    expect((wrapper.vm).hasAddress).toBe(false);
  });

  it('registers ESC key handler on mount', () => {
    mountWithPlugins(CustomerDetailsDrawer, { props: defaultProps });
    const keydownCalls = addSpy.mock.calls.filter(([event]: [string, ...unknown[]]) => event === 'keydown');
    expect(keydownCalls.length).toBeGreaterThanOrEqual(1);
  });

  it('handles ESC key when drawer is open', () => {
    const wrapper = mountWithPlugins(CustomerDetailsDrawer, { props: defaultProps });
    const escEvent = new KeyboardEvent('keydown', { key: 'Escape' });
    window.dispatchEvent(escEvent);
    expect(wrapper.emitted('update:modelValue')).toBeTruthy();
    expect(wrapper.emitted('update:modelValue')![0]).toEqual([false]);
  });

  it('ignores ESC key when drawer is closed', () => {
    const wrapper = mountWithPlugins(CustomerDetailsDrawer, {
      props: { ...defaultProps, modelValue: false },
    });
    const escEvent = new KeyboardEvent('keydown', { key: 'Escape' });
    window.dispatchEvent(escEvent);
    expect(wrapper.emitted('update:modelValue')).toBeFalsy();
  });

  it('emits update:modelValue(false) on close', () => {
    const wrapper = mountWithPlugins(CustomerDetailsDrawer, { props: defaultProps });
    (wrapper.vm).close();
    expect(wrapper.emitted('update:modelValue')![0]).toEqual([false]);
  });

  it('does not emit edit when customer is null', () => {
    const wrapper = mountWithPlugins(CustomerDetailsDrawer, { props: defaultProps });
    (wrapper.vm).handleEdit();
    expect(wrapper.emitted('edit')).toBeFalsy();
  });

  it('getTypeIcon returns correct icon for customer type', () => {
    const wrapper = mountWithPlugins(CustomerDetailsDrawer, { props: defaultProps });
    expect((wrapper.vm).getTypeIcon('customer')).toBe('domain');
    expect((wrapper.vm).getTypeIcon('site')).toBe('location_on');
    expect((wrapper.vm).getTypeIcon('building')).toBe('apartment');
    expect((wrapper.vm).getTypeIcon('unknown')).toBe('domain');
  });

  it('getTypeColorName returns correct color for customer type', () => {
    const wrapper = mountWithPlugins(CustomerDetailsDrawer, { props: defaultProps });
    expect((wrapper.vm).getTypeColorName('customer')).toBe('green');
    expect((wrapper.vm).getTypeColorName('site')).toBe('orange');
    expect((wrapper.vm).getTypeColorName('unknown')).toBe('primary');
  });
});
