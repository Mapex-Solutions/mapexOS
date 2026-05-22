import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { mountWithPlugins } from '@src/test/helpers';
import AssetDetailsDrawer from './AssetDetailsDrawer.vue';

vi.mock('quasar', () => ({
  date: {
    formatDate: vi.fn(() => 'Jan 01, 2024 12:00'),
  },
}));

// Deep proxy: returns a self-replicating proxy at any nesting depth,
// terminating with `value` to satisfy the i18n composable contract
// `(t.path.to.key.value -> string)`. Lets the test mount the drawer
// without registering vue-i18n or maintaining a fixture per level.
function deepI18nProxy(label = ''): any {
  return new Proxy({ value: label }, {
    get: (target: any, prop: string) => {
      if (prop === 'value') return target.value;
      return deepI18nProxy(String(prop));
    },
  });
}

vi.mock('@composables/i18n', () => ({
  useAssetsTranslations: () => deepI18nProxy(''),
}));

vi.mock('@composables/useLogger', () => ({
  useLogger: () => ({
    debug: vi.fn(),
    info: vi.fn(),
    warn: vi.fn(),
    error: vi.fn(),
  }),
}));

vi.mock('@utils/alert', () => ({
  notifyFail: vi.fn(),
  notifySuccess: vi.fn(),
}));

vi.mock('@utils/zipDownload', () => ({
  downloadCertZip: vi.fn(),
  decodeBase64ToBytes: vi.fn(() => new Uint8Array()),
}));

vi.mock('@services/mapex', () => ({
  apis: {
    assets: {
      asset: {
        getById: vi.fn().mockResolvedValue({
          id: 'asset-1',
          name: 'Test Asset',
          enabled: true,
          orgId: 'org-1',
        }),
      },
      mqttcerts: {
        issueCert: vi.fn(),
        revokeCert: vi.fn(),
        listRevoked: vi.fn().mockResolvedValue([]),
      },
    },
  },
}));

vi.mock('@stores/organization', () => ({
  useOrganizationStore: () => ({
    flatList: [{ id: 'org-1', name: 'Test Org' }],
  }),
}));

describe('AssetDetailsDrawer', () => {
  const defaultProps = {
    modelValue: true,
    assetId: 'asset-1',
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
    const wrapper = mountWithPlugins(AssetDetailsDrawer, { props: defaultProps });
    expect(wrapper.exists()).toBe(true);
  });

  it('starts with loading state', () => {
    const wrapper = mountWithPlugins(AssetDetailsDrawer, { props: defaultProps });
    expect((wrapper.vm).loading).toBe(true);
  });

  it('starts with asset as null', () => {
    const wrapper = mountWithPlugins(AssetDetailsDrawer, { props: defaultProps });
    // Initially null before fetch completes
    expect((wrapper.vm).asset).toBeNull();
  });

  it('computes hasLocation based on asset data', () => {
    const wrapper = mountWithPlugins(AssetDetailsDrawer, { props: defaultProps });
    // With no asset loaded, hasLocation should be false
    expect((wrapper.vm).hasLocation).toBe(false);
  });

  it('registers ESC key handler on mount', () => {
    mountWithPlugins(AssetDetailsDrawer, { props: defaultProps });
    const keydownCalls = addSpy.mock.calls.filter(([event]: [string, ...unknown[]]) => event === 'keydown');
    expect(keydownCalls.length).toBeGreaterThanOrEqual(1);
  });

  it('handles ESC key when drawer is open', () => {
    const wrapper = mountWithPlugins(AssetDetailsDrawer, { props: defaultProps });
    const escEvent = new KeyboardEvent('keydown', { key: 'Escape' });
    window.dispatchEvent(escEvent);
    expect(wrapper.emitted('update:modelValue')).toBeTruthy();
    expect(wrapper.emitted('update:modelValue')![0]).toEqual([false]);
  });

  it('ignores ESC key when drawer is closed', () => {
    const wrapper = mountWithPlugins(AssetDetailsDrawer, {
      props: { ...defaultProps, modelValue: false },
    });
    const escEvent = new KeyboardEvent('keydown', { key: 'Escape' });
    window.dispatchEvent(escEvent);
    expect(wrapper.emitted('update:modelValue')).toBeFalsy();
  });

  it('emits update:modelValue(false) on close', () => {
    const wrapper = mountWithPlugins(AssetDetailsDrawer, { props: defaultProps });
    (wrapper.vm).close();
    expect(wrapper.emitted('update:modelValue')![0]).toEqual([false]);
  });

  it('does not emit edit when asset is null', () => {
    const wrapper = mountWithPlugins(AssetDetailsDrawer, { props: defaultProps });
    (wrapper.vm).handleEdit();
    expect(wrapper.emitted('edit')).toBeFalsy();
  });
});
