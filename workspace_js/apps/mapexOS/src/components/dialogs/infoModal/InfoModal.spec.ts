import { describe, it, expect } from 'vitest';
import { mountWithPlugins } from '@src/test/helpers';
import InfoModal from './InfoModal.vue';

describe('InfoModal', () => {
  const defaultProps = {
    modelValue: true,
    title: 'Test Title',
    description: 'Test description text',
  };

  it('renders without errors', () => {
    const wrapper = mountWithPlugins(InfoModal, { props: defaultProps });
    expect(wrapper.exists()).toBe(true);
  });

  it('syncs isOpen ref with modelValue prop', () => {
    const wrapper = mountWithPlugins(InfoModal, { props: defaultProps });
    expect(wrapper.vm.isOpen).toBe(true);
  });

  it('sets isOpen to false when modelValue is false', () => {
    const wrapper = mountWithPlugins(InfoModal, {
      props: { ...defaultProps, modelValue: false },
    });
    expect(wrapper.vm.isOpen).toBe(false);
  });

  it('emits update:modelValue when isOpen changes', async () => {
    const wrapper = mountWithPlugins(InfoModal, { props: defaultProps });
    wrapper.vm.isOpen = false;
    await wrapper.vm.$nextTick();
    expect(wrapper.emitted('update:modelValue')).toBeTruthy();
    expect(wrapper.emitted('update:modelValue')![0]).toEqual([false]);
  });

  it('emits update:modelValue(false) on handleClose', () => {
    const wrapper = mountWithPlugins(InfoModal, { props: defaultProps });
    (wrapper.vm).handleClose();
    expect(wrapper.emitted('update:modelValue')).toBeTruthy();
    expect(wrapper.emitted('update:modelValue')![0]).toEqual([false]);
  });

  it('defaults showActions to true', () => {
    const wrapper = mountWithPlugins(InfoModal, { props: defaultProps });
    expect(wrapper.props('showActions')).toBe(true);
  });

  it('updates isOpen when modelValue prop changes', async () => {
    const wrapper = mountWithPlugins(InfoModal, { props: defaultProps });
    expect(wrapper.vm.isOpen).toBe(true);
    await wrapper.setProps({ modelValue: false });
    expect(wrapper.vm.isOpen).toBe(false);
  });
});
