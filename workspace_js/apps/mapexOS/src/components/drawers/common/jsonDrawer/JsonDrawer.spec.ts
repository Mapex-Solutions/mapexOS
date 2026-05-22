import { describe, it, expect, vi } from 'vitest';
import { mountWithPlugins } from '@src/test/helpers';
import JsonDrawer from './JsonDrawer.vue';

vi.mock('monaco-editor', () => ({
  editor: {
    create: vi.fn(() => ({
      dispose: vi.fn(),
      getValue: vi.fn(() => '{}'),
      getAction: vi.fn(() => ({ run: vi.fn() })),
    })),
  },
}));

vi.mock('@utils/monaco-theme', () => ({
  registerMapexMonacoThemes: vi.fn(),
  getMapexMonacoTheme: vi.fn(() => 'mapex-dark'),
  applyMapexMonacoTheme: vi.fn(),
}));

vi.mock('@utils/alert/notify', () => ({
  notifyInfo: vi.fn(),
  notifySuccess: vi.fn(),
  notifyFail: vi.fn(),
  notifyWarning: vi.fn(),
}));

vi.mock('quasar', () => ({
  useQuasar: () => ({ screen: { lt: { md: false }, width: 1024 } }),
}));

vi.mock('@stores/theme', () => ({
  useThemeStore: () => ({ isDark: false }),
}));

describe('JsonDrawer', () => {
  const defaultProps = {
    show: true,
    title: 'JSON Viewer',
    jsonData: { key: 'value' },
  };

  it('renders without errors', () => {
    const wrapper = mountWithPlugins(JsonDrawer, { props: defaultProps });
    expect(wrapper.exists()).toBe(true);
  });

  it('syncs localShow with show prop', () => {
    const wrapper = mountWithPlugins(JsonDrawer, { props: defaultProps });
    expect((wrapper.vm).localShow).toBe(true);
  });

  it('sets localShow to false when show is false', () => {
    const wrapper = mountWithPlugins(JsonDrawer, {
      props: { ...defaultProps, show: false },
    });
    expect((wrapper.vm).localShow).toBe(false);
  });

  it('emits update:show when localShow changes', async () => {
    const wrapper = mountWithPlugins(JsonDrawer, { props: defaultProps });
    (wrapper.vm).localShow = false;
    await wrapper.vm.$nextTick();
    expect(wrapper.emitted('update:show')).toBeTruthy();
    expect(wrapper.emitted('update:show')![0]).toEqual([false]);
  });

  it('updates localShow when show prop changes', async () => {
    const wrapper = mountWithPlugins(JsonDrawer, { props: defaultProps });
    expect((wrapper.vm).localShow).toBe(true);
    await wrapper.setProps({ show: false });
    expect((wrapper.vm).localShow).toBe(false);
  });

  it('handles string jsonData', () => {
    const wrapper = mountWithPlugins(JsonDrawer, {
      props: { ...defaultProps, jsonData: '{"raw":"json"}' },
    });
    expect(wrapper.exists()).toBe(true);
  });
});
