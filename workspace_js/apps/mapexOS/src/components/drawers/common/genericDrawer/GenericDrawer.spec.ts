import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { mountWithPlugins } from '@src/test/helpers';
import GenericDrawer from './GenericDrawer.vue';

describe('GenericDrawer', () => {
  const defaultProps = {
    modelValue: true,
    title: 'Test Drawer',
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
    const wrapper = mountWithPlugins(GenericDrawer, { props: defaultProps });
    expect(wrapper.exists()).toBe(true);
  });

  it('receives modelValue prop correctly', () => {
    const wrapper = mountWithPlugins(GenericDrawer, { props: defaultProps });
    expect(wrapper.props('modelValue')).toBe(true);
  });

  it('defaults icon to "info"', () => {
    const wrapper = mountWithPlugins(GenericDrawer, { props: defaultProps });
    expect(wrapper.props('icon')).toBe('info');
  });

  it('defaults iconColor to "primary"', () => {
    const wrapper = mountWithPlugins(GenericDrawer, { props: defaultProps });
    expect(wrapper.props('iconColor')).toBe('primary');
  });

  it('defaults width to 380', () => {
    const wrapper = mountWithPlugins(GenericDrawer, { props: defaultProps });
    expect(wrapper.props('width')).toBe(380);
  });

  it('defaults closeTooltip to "Close"', () => {
    const wrapper = mountWithPlugins(GenericDrawer, { props: defaultProps });
    expect(wrapper.props('closeTooltip')).toBe('Close');
  });

  it('emits update:modelValue(false) and close on handleClose', () => {
    const wrapper = mountWithPlugins(GenericDrawer, { props: defaultProps });
    (wrapper.vm).handleClose();
    expect(wrapper.emitted('update:modelValue')).toBeTruthy();
    expect(wrapper.emitted('update:modelValue')![0]).toEqual([false]);
    expect(wrapper.emitted('close')).toBeTruthy();
  });

  it('registers ESC key handler on mount', () => {
    mountWithPlugins(GenericDrawer, { props: defaultProps });
    const keydownCalls = addSpy.mock.calls.filter(([event]: [string, ...unknown[]]) => event === 'keydown');
    expect(keydownCalls.length).toBeGreaterThanOrEqual(1);
  });

  it('handles ESC key when drawer is open', () => {
    const wrapper = mountWithPlugins(GenericDrawer, { props: defaultProps });
    const escEvent = new KeyboardEvent('keydown', { key: 'Escape' });
    window.dispatchEvent(escEvent);
    expect(wrapper.emitted('update:modelValue')).toBeTruthy();
    expect(wrapper.emitted('update:modelValue')![0]).toEqual([false]);
  });

  it('ignores ESC key when drawer is closed', () => {
    const wrapper = mountWithPlugins(GenericDrawer, {
      props: { ...defaultProps, modelValue: false },
    });
    const escEvent = new KeyboardEvent('keydown', { key: 'Escape' });
    window.dispatchEvent(escEvent);
    expect(wrapper.emitted('update:modelValue')).toBeFalsy();
  });

  it('accepts custom width', () => {
    const wrapper = mountWithPlugins(GenericDrawer, {
      props: { ...defaultProps, width: 500 },
    });
    expect(wrapper.props('width')).toBe(500);
  });
});
