import { describe, it, expect, vi } from 'vitest';
import { mountWithPlugins } from '@src/test/helpers';
import AssetClassificationSelector from './AssetClassificationSelector.vue';

vi.mock('./components/guidedModeSelector', () => ({
  GuidedModeSelector: { name: 'GuidedModeSelector', template: '<div />' },
}));

describe('AssetClassificationSelector', () => {
  it('renders without errors', () => {
    const wrapper = mountWithPlugins(AssetClassificationSelector, {
      props: { modelValue: undefined },
    });
    expect(wrapper.exists()).toBe(true);
  });

  it('renders with a modelValue', () => {
    const wrapper = mountWithPlugins(AssetClassificationSelector, {
      props: {
        modelValue: {
          categoryId: 'cat-1',
          manufacturerId: 'man-1',
          modelId: 'mod-1',
          version: '1.0',
        },
      },
    });
    expect(wrapper.exists()).toBe(true);
  });

  it('emits update:modelValue on handleValueUpdate', () => {
    const wrapper = mountWithPlugins(AssetClassificationSelector, {
      props: { modelValue: undefined },
    });
    const value = {
      categoryId: 'cat-1',
      manufacturerId: 'man-1',
      modelId: 'mod-1',
      version: '2.0',
    };
    wrapper.vm.handleValueUpdate(value);
    const emitted = wrapper.emitted('update:modelValue')!;
    expect(emitted[0]![0]).toEqual(value);
  });

  it('renders with disabled prop', () => {
    const wrapper = mountWithPlugins(AssetClassificationSelector, {
      props: { modelValue: undefined, disabled: true },
    });
    expect(wrapper.exists()).toBe(true);
  });

  it('renders with required prop', () => {
    const wrapper = mountWithPlugins(AssetClassificationSelector, {
      props: { modelValue: undefined, required: true },
    });
    expect(wrapper.exists()).toBe(true);
  });
});
