import { describe, it, expect, vi } from 'vitest';
import { mountWithPlugins } from '@src/test/helpers';
import JsonAndEventDrawer from './JsonAndEventDrawer.vue';

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

vi.mock('lodash', () => ({
  isEmpty: (val: any) => !val || (typeof val === 'object' && Object.keys(val).length === 0),
}));

describe('JsonAndEventDrawer', () => {
  const defaultProps = {
    show: true,
    title: 'JSON & Event Viewer',
    jsonData: { eventId: 'evt-1', key: 'value' },
  };

  it('renders without errors', () => {
    const wrapper = mountWithPlugins(JsonAndEventDrawer, { props: defaultProps });
    expect(wrapper.exists()).toBe(true);
  });

  it('syncs localShow with show prop', () => {
    const wrapper = mountWithPlugins(JsonAndEventDrawer, { props: defaultProps });
    expect((wrapper.vm).localShow).toBe(true);
  });

  it('initializes tab to "json"', () => {
    const wrapper = mountWithPlugins(JsonAndEventDrawer, { props: defaultProps });
    expect((wrapper.vm).tab).toBe('json');
  });

  it('emits update:show when localShow changes', async () => {
    const wrapper = mountWithPlugins(JsonAndEventDrawer, { props: defaultProps });
    (wrapper.vm).localShow = false;
    await wrapper.vm.$nextTick();
    expect(wrapper.emitted('update:show')).toBeTruthy();
    expect(wrapper.emitted('update:show')![0]).toEqual([false]);
  });

  it('sets localShow to false when show is false', () => {
    const wrapper = mountWithPlugins(JsonAndEventDrawer, {
      props: { ...defaultProps, show: false },
    });
    expect((wrapper.vm).localShow).toBe(false);
  });

  it('emits fetch-event on tab click to "event"', () => {
    const wrapper = mountWithPlugins(JsonAndEventDrawer, { props: defaultProps });
    (wrapper.vm).onTabClick('event');
    expect(wrapper.emitted('fetch-event')).toBeTruthy();
    expect(wrapper.emitted('fetch-event')![0]).toEqual(['evt-1']);
  });

  it('has drawerTabs with json and event entries', () => {
    const wrapper = mountWithPlugins(JsonAndEventDrawer, { props: defaultProps });
    const tabs = (wrapper.vm).drawerTabs;
    expect(tabs).toHaveLength(2);
    expect(tabs[0].name).toBe('json');
    expect(tabs[1].name).toBe('event');
  });
});
