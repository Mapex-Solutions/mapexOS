import { describe, it, expect, vi } from 'vitest';
import { mountWithPlugins, createMockTranslations } from '@src/test/helpers';
import StandardizedPayloadHelpModal from './StandardizedPayloadHelpModal.vue';

vi.mock('quasar', () => ({
  copyToClipboard: vi.fn().mockResolvedValue(undefined),
}));

vi.mock('@components/buttons', () => ({
  BaseButton: { name: 'BaseButton', template: '<button><slot /></button>' },
}));

vi.mock('@components/chips', () => ({
  DetailChip: { name: 'DetailChip', template: '<span />' },
}));

vi.mock('@components/tooltips', () => ({
  AppTooltip: { name: 'AppTooltip', template: '<div />' },
}));

vi.mock('@utils/alert/notify', () => ({
  notifySuccess: vi.fn(),
  notifyFail: vi.fn(),
}));

vi.mock('@src/composables/i18n/components/dialogs/useStandardizedPayloadHelpTranslations', () => ({
  useStandardizedPayloadHelpTranslations: () => createMockTranslations(),
}));

const BASE_PROPS = {
  modelValue: true,
};

describe('StandardizedPayloadHelpModal', () => {
  it('renders without errors', () => {
    const wrapper = mountWithPlugins(StandardizedPayloadHelpModal, { props: BASE_PROPS });
    expect(wrapper.exists()).toBe(true);
  });

  it('computes isOpen from modelValue', () => {
    const wrapper = mountWithPlugins(StandardizedPayloadHelpModal, { props: BASE_PROPS });
    expect(wrapper.vm.isOpen).toBe(true);
  });

  it('computes isOpen as false when modelValue is false', () => {
    const wrapper = mountWithPlugins(StandardizedPayloadHelpModal, {
      props: { modelValue: false },
    });
    expect(wrapper.vm.isOpen).toBe(false);
  });

  it('emits update:modelValue on closeModal', () => {
    const wrapper = mountWithPlugins(StandardizedPayloadHelpModal, { props: BASE_PROPS });
    wrapper.vm.closeModal();
    const emitted = wrapper.emitted('update:modelValue')!;
    expect(emitted[0]![0]).toBe(false);
  });

  it('calls copyToClipboard on copyCode', async () => {
    const { copyToClipboard } = await import('quasar');
    const wrapper = mountWithPlugins(StandardizedPayloadHelpModal, { props: BASE_PROPS });
    await wrapper.vm.copyCode('test code');
    expect(copyToClipboard).toHaveBeenCalledWith('test code');
  });
});
