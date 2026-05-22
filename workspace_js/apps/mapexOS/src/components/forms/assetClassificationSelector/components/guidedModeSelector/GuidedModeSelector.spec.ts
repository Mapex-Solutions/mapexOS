import { describe, it, expect, vi } from 'vitest';
import { mountWithPlugins, createMockTranslations } from '@src/test/helpers';
import GuidedModeSelector from './GuidedModeSelector.vue';

vi.mock('@composables/i18n/components/forms/assetClassificationSelector/useAssetClassificationSelectorTranslations', () => ({
  useAssetClassificationSelectorTranslations: () => createMockTranslations(),
}));

vi.mock('@composables/useLogger', () => ({
  useLogger: () => ({
    error: vi.fn(),
    warn: vi.fn(),
    info: vi.fn(),
    debug: vi.fn(),
  }),
}));

vi.mock('@services/mapex', () => ({
  apis: {
    mapexOS: {
      lists: {
        list: vi.fn().mockResolvedValue({ items: [], pagination: { totalPages: 1 } }),
      },
    },
  },
}));

vi.mock('@utils/alert', () => ({
  notifyFail: vi.fn(),
}));

describe('GuidedModeSelector', () => {
  it('renders without errors', () => {
    const wrapper = mountWithPlugins(GuidedModeSelector, {
      props: { modelValue: undefined },
    });
    expect(wrapper.exists()).toBe(true);
  });

  it('starts with no selections', () => {
    const wrapper = mountWithPlugins(GuidedModeSelector, {
      props: { modelValue: undefined },
    });
    expect(wrapper.vm.selectedCategoryId).toBeUndefined();
    expect(wrapper.vm.selectedManufacturerId).toBeUndefined();
    expect(wrapper.vm.selectedModelId).toBeUndefined();
    expect(wrapper.vm.version).toBe('');
  });

  it('initializes from modelValue', () => {
    const wrapper = mountWithPlugins(GuidedModeSelector, {
      props: {
        modelValue: {
          categoryId: 'cat-1',
          manufacturerId: 'man-1',
          modelId: 'mod-1',
          version: '1.0',
          categoryName: 'Cat',
          manufacturerName: 'Man',
          modelName: 'Mod',
        },
      },
    });
    expect(wrapper.vm.selectedCategoryId).toBe('cat-1');
    expect(wrapper.vm.selectedManufacturerId).toBe('man-1');
    expect(wrapper.vm.selectedModelId).toBe('mod-1');
    expect(wrapper.vm.version).toBe('1.0');
  });

  it('computes manufacturerDisabled as true when no category', () => {
    const wrapper = mountWithPlugins(GuidedModeSelector, {
      props: { modelValue: undefined },
    });
    expect(wrapper.vm.manufacturerDisabled).toBe(true);
  });

  it('computes modelDisabled as true when no manufacturer', () => {
    const wrapper = mountWithPlugins(GuidedModeSelector, {
      props: { modelValue: undefined },
    });
    expect(wrapper.vm.modelDisabled).toBe(true);
  });

  it('handleCategoryChange resets downstream selections', () => {
    const wrapper = mountWithPlugins(GuidedModeSelector, {
      props: {
        modelValue: {
          categoryId: 'cat-1',
          manufacturerId: 'man-1',
          modelId: 'mod-1',
          version: '1.0',
        },
      },
    });
    wrapper.vm.handleCategoryChange('cat-2');
    expect(wrapper.vm.selectedCategoryId).toBe('cat-2');
    expect(wrapper.vm.selectedManufacturerId).toBeUndefined();
    expect(wrapper.vm.selectedModelId).toBeUndefined();
    expect(wrapper.vm.version).toBe('');
  });

  it('handleManufacturerChange resets model and version', () => {
    const wrapper = mountWithPlugins(GuidedModeSelector, {
      props: {
        modelValue: {
          categoryId: 'cat-1',
          manufacturerId: 'man-1',
          modelId: 'mod-1',
          version: '1.0',
        },
      },
    });
    wrapper.vm.handleManufacturerChange('man-2');
    expect(wrapper.vm.selectedManufacturerId).toBe('man-2');
    expect(wrapper.vm.selectedModelId).toBeUndefined();
    expect(wrapper.vm.version).toBe('');
  });

  it('handleModelChange updates selectedModelId', () => {
    const wrapper = mountWithPlugins(GuidedModeSelector, {
      props: { modelValue: undefined },
    });
    wrapper.vm.handleModelChange('mod-1');
    expect(wrapper.vm.selectedModelId).toBe('mod-1');
  });

  it('handleCategoryChange with null clears category', () => {
    const wrapper = mountWithPlugins(GuidedModeSelector, {
      props: { modelValue: undefined },
    });
    wrapper.vm.handleCategoryChange(null);
    expect(wrapper.vm.selectedCategoryId).toBeUndefined();
  });

  it('starts with empty lists', () => {
    const wrapper = mountWithPlugins(GuidedModeSelector, {
      props: { modelValue: undefined },
    });
    expect(wrapper.vm.manufacturers).toEqual([]);
    expect(wrapper.vm.models).toEqual([]);
  });
});
