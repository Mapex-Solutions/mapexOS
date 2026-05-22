import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { mountWithPlugins } from '@src/test/helpers';
import WorkflowDetailsDrawer from './WorkflowDetailsDrawer.vue';

vi.mock('quasar', () => ({
  date: {
    formatDate: vi.fn(() => 'Jan 01, 2024 12:00'),
  },
}));

vi.mock('@composables/i18n', () => ({
  useWorkflowListPageTranslations: () => new Proxy({}, {
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

vi.mock('src/pages/automations/workflows/workflowListPage/constants', () => ({
  DUMMY_WORKFLOWS: [
    {
      id: 'wf-1',
      name: 'Test Workflow',
      description: 'A test workflow',
      enabled: true,
      isTemplate: false,
      definitionVersion: 1,
      nodesCount: 5,
      edgesCount: 4,
      timezone: 'UTC',
      created: '2024-01-01T00:00:00Z',
      updated: '2024-01-01T00:00:00Z',
    },
  ],
}));

describe('WorkflowDetailsDrawer', () => {
  const defaultProps = {
    modelValue: true,
    workflowId: 'wf-1',
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
    const wrapper = mountWithPlugins(WorkflowDetailsDrawer, { props: defaultProps });
    expect(wrapper.exists()).toBe(true);
  });

  it('starts with loading state', () => {
    const wrapper = mountWithPlugins(WorkflowDetailsDrawer, { props: defaultProps });
    expect((wrapper.vm).loading).toBe(true);
  });

  it('starts with workflow as null', () => {
    const wrapper = mountWithPlugins(WorkflowDetailsDrawer, { props: defaultProps });
    expect((wrapper.vm).workflow).toBeNull();
  });

  it('registers ESC key handler on mount', () => {
    mountWithPlugins(WorkflowDetailsDrawer, { props: defaultProps });
    const keydownCalls = addSpy.mock.calls.filter(([event]: [string, ...unknown[]]) => event === 'keydown');
    expect(keydownCalls.length).toBeGreaterThanOrEqual(1);
  });

  it('handles ESC key when drawer is open', () => {
    const wrapper = mountWithPlugins(WorkflowDetailsDrawer, { props: defaultProps });
    const escEvent = new KeyboardEvent('keydown', { key: 'Escape' });
    window.dispatchEvent(escEvent);
    expect(wrapper.emitted('update:modelValue')).toBeTruthy();
    expect(wrapper.emitted('update:modelValue')![0]).toEqual([false]);
  });

  it('ignores ESC key when drawer is closed', () => {
    const wrapper = mountWithPlugins(WorkflowDetailsDrawer, {
      props: { ...defaultProps, modelValue: false },
    });
    const escEvent = new KeyboardEvent('keydown', { key: 'Escape' });
    window.dispatchEvent(escEvent);
    expect(wrapper.emitted('update:modelValue')).toBeFalsy();
  });

  it('emits update:modelValue(false) on close', () => {
    const wrapper = mountWithPlugins(WorkflowDetailsDrawer, { props: defaultProps });
    (wrapper.vm).close();
    expect(wrapper.emitted('update:modelValue')![0]).toEqual([false]);
  });

  it('does not emit edit when workflow is null', () => {
    const wrapper = mountWithPlugins(WorkflowDetailsDrawer, { props: defaultProps });
    (wrapper.vm).handleEdit();
    expect(wrapper.emitted('edit')).toBeFalsy();
  });
});
