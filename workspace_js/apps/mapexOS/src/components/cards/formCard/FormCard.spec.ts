import { describe, it, expect } from 'vitest';
import { mountWithPlugins } from '@src/test/helpers';
import FormCard from './FormCard.vue';

const baseHeader = {
  icon: 'settings',
  title: 'Test Form',
  description: 'A test form card',
};

function factory(overrides: Record<string, unknown> = {}) {
  return mountWithPlugins(FormCard, {
    props: {
      header: baseHeader,
      navigation: {
        currentStep: 1,
        totalSteps: 3,
        showPreviousButton: true,
        showNextButton: true,
        showSaveButton: true,
      },
      ...overrides,
    },
  });
}

describe('FormCard', () => {
  it('renders without errors', () => {
    const wrapper = factory();
    expect(wrapper.exists()).toBe(true);
  });

  it('shouldShowPreviousButton is false on step 1', () => {
    const wrapper = factory();
    expect(wrapper.vm.shouldShowPreviousButton).toBe(false);
  });

  it('shouldShowPreviousButton is true on step 2', () => {
    const wrapper = factory({
      navigation: {
        currentStep: 2,
        totalSteps: 3,
        showPreviousButton: true,
        showNextButton: true,
        showSaveButton: true,
      },
    });
    expect(wrapper.vm.shouldShowPreviousButton).toBe(true);
  });

  it('shouldShowPreviousButton is false when showPreviousButton is disabled', () => {
    const wrapper = factory({
      navigation: {
        currentStep: 2,
        totalSteps: 3,
        showPreviousButton: false,
        showNextButton: true,
        showSaveButton: true,
      },
    });
    expect(wrapper.vm.shouldShowPreviousButton).toBe(false);
  });

  it('shouldShowNextButton is true when not on last step', () => {
    const wrapper = factory();
    expect(wrapper.vm.shouldShowNextButton).toBe(true);
  });

  it('shouldShowNextButton is false on last step', () => {
    const wrapper = factory({
      navigation: {
        currentStep: 3,
        totalSteps: 3,
        showPreviousButton: true,
        showNextButton: true,
        showSaveButton: true,
      },
    });
    expect(wrapper.vm.shouldShowNextButton).toBe(false);
  });

  it('shouldShowSaveButton is false when not on last step', () => {
    const wrapper = factory();
    expect(wrapper.vm.shouldShowSaveButton).toBe(false);
  });

  it('shouldShowSaveButton is true on last step', () => {
    const wrapper = factory({
      navigation: {
        currentStep: 3,
        totalSteps: 3,
        showPreviousButton: true,
        showNextButton: true,
        showSaveButton: true,
      },
    });
    expect(wrapper.vm.shouldShowSaveButton).toBe(true);
  });

  it('shouldShowSaveButton is false when showSaveButton is disabled', () => {
    const wrapper = factory({
      navigation: {
        currentStep: 3,
        totalSteps: 3,
        showPreviousButton: true,
        showNextButton: true,
        showSaveButton: false,
      },
    });
    expect(wrapper.vm.shouldShowSaveButton).toBe(false);
  });

  it('emits previous with decremented step', () => {
    const wrapper = factory({
      navigation: {
        currentStep: 2,
        totalSteps: 3,
        showPreviousButton: true,
        showNextButton: true,
        showSaveButton: true,
      },
    });
    wrapper.vm.handlePreviousStep();
    expect(wrapper.emitted('previous')).toBeTruthy();
    expect(wrapper.emitted('previous')![0]).toEqual([1]);
  });

  it('emits next with incremented step', () => {
    const wrapper = factory();
    wrapper.vm.handleNextStep();
    expect(wrapper.emitted('next')).toBeTruthy();
    expect(wrapper.emitted('next')![0]).toEqual([2]);
  });

  it('emits save event', () => {
    const wrapper = factory();
    wrapper.vm.handleSave();
    expect(wrapper.emitted('save')).toBeTruthy();
  });
});
