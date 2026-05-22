import { describe, it, expect } from 'vitest';
import { mountWithPlugins } from '@src/test/helpers';
import StepperVertical from './StepperVertical.vue';
import type { StepperVerticalProps, StepperVerticalItem } from './interfaces';

/** Stub that renders its default slot as a plain div */
const SlotStub = { template: '<div><slot /></div>' };

const quasarStubs = {
  'q-card': SlotStub,
  'q-card-section': SlotStub,
  'q-icon': { template: '<i />' },
  'q-separator': { template: '<hr />' },
};

const mockSteps: StepperVerticalItem[] = [
  { title: 'Basic Info', description: 'Enter basic details', icon: 'info' },
  { title: 'Configuration', description: 'Set up config', icon: 'settings' },
  { title: 'Review', description: 'Review and confirm', icon: 'check' },
];

const makeProps = (overrides: Partial<StepperVerticalProps> = {}): StepperVerticalProps => ({
  steps: mockSteps,
  ...overrides,
});

describe('StepperVertical', () => {
  describe('rendering', () => {
    it('should render with minimal props (steps only)', () => {
      const wrapper = mountWithPlugins(StepperVertical, {
        props: makeProps(),
        stubs: quasarStubs,
      });
      expect(wrapper.exists()).toBe(true);
    });

    it('should render default title and subtitle', () => {
      const wrapper = mountWithPlugins(StepperVertical, {
        props: makeProps(),
        stubs: quasarStubs,
      });
      expect(wrapper.text()).toContain('Configuration Steps');
      expect(wrapper.text()).toContain('Complete all steps');
    });

    it('should render custom title and subtitle', () => {
      const wrapper = mountWithPlugins(StepperVertical, {
        props: makeProps({ title: 'Setup Wizard', subtitle: 'Follow the steps below' }),
        stubs: quasarStubs,
      });
      expect(wrapper.text()).toContain('Setup Wizard');
      expect(wrapper.text()).toContain('Follow the steps below');
    });

    it('should render all step titles and descriptions', () => {
      const wrapper = mountWithPlugins(StepperVertical, {
        props: makeProps(),
        stubs: quasarStubs,
      });
      expect(wrapper.text()).toContain('Basic Info');
      expect(wrapper.text()).toContain('Enter basic details');
      expect(wrapper.text()).toContain('Configuration');
      expect(wrapper.text()).toContain('Set up config');
      expect(wrapper.text()).toContain('Review');
      expect(wrapper.text()).toContain('Review and confirm');
    });

    it('should render the correct number of step items', () => {
      const wrapper = mountWithPlugins(StepperVertical, {
        props: makeProps(),
        stubs: quasarStubs,
      });
      const stepItems = wrapper.findAll('.step-item');
      expect(stepItems).toHaveLength(3);
    });

    it('should render default info text', () => {
      const wrapper = mountWithPlugins(StepperVertical, {
        props: makeProps(),
        stubs: quasarStubs,
      });
      expect(wrapper.text()).toContain('All fields marked with * are required');
    });

    it('should render custom info text', () => {
      const wrapper = mountWithPlugins(StepperVertical, {
        props: makeProps({ infoText: 'Custom info' }),
        stubs: quasarStubs,
      });
      expect(wrapper.text()).toContain('Custom info');
    });

    it('should apply fullHeight class when fullHeight is true', () => {
      const wrapper = mountWithPlugins(StepperVertical, {
        props: makeProps({ fullHeight: true }),
        stubs: quasarStubs,
      });
      expect(wrapper.find('.stepper-card').exists()).toBe(true);
    });

    it('should not apply fullHeight class when fullHeight is false', () => {
      const wrapper = mountWithPlugins(StepperVertical, {
        props: makeProps({ fullHeight: false }),
        stubs: quasarStubs,
      });
      expect(wrapper.find('.stepper-card').exists()).toBe(false);
    });
  });

  describe('step ID prefix', () => {
    it('should generate step IDs when stepIdPrefix is provided', () => {
      const wrapper = mountWithPlugins(StepperVertical, {
        props: makeProps({ stepIdPrefix: 'wizard-step' }),
      });
      expect(wrapper.vm.getStepAttrs(0)).toEqual({ id: 'wizard-step-1' });
      expect(wrapper.vm.getStepAttrs(1)).toEqual({ id: 'wizard-step-2' });
      expect(wrapper.vm.getStepAttrs(2)).toEqual({ id: 'wizard-step-3' });
    });

    it('should return empty attrs when stepIdPrefix is not provided', () => {
      const wrapper = mountWithPlugins(StepperVertical, { props: makeProps() });
      expect(wrapper.vm.getStepAttrs(0)).toEqual({});
    });
  });

  describe('isActive (creating mode)', () => {
    it('should mark only the current step as active', () => {
      const wrapper = mountWithPlugins(StepperVertical, {
        props: makeProps({ currentStep: 2 }),
      });
      expect(wrapper.vm.isActive(0)).toBe(false);
      expect(wrapper.vm.isActive(1)).toBe(true);
      expect(wrapper.vm.isActive(2)).toBe(false);
    });
  });

  describe('isActive (editing mode)', () => {
    it('should mark all steps as active in editing mode', () => {
      const wrapper = mountWithPlugins(StepperVertical, {
        props: makeProps({ mode: 'editing', currentStep: 1 }),
      });
      expect(wrapper.vm.isActive(0)).toBe(true);
      expect(wrapper.vm.isActive(1)).toBe(true);
      expect(wrapper.vm.isActive(2)).toBe(true);
    });

    it('should mark all steps as active when allowStepNavigation is true', () => {
      const wrapper = mountWithPlugins(StepperVertical, {
        props: makeProps({ allowStepNavigation: true, currentStep: 1 }),
      });
      expect(wrapper.vm.isActive(0)).toBe(true);
      expect(wrapper.vm.isActive(1)).toBe(true);
      expect(wrapper.vm.isActive(2)).toBe(true);
    });
  });

  describe('isCompleted', () => {
    it('should mark previous steps as completed in creating mode', () => {
      const wrapper = mountWithPlugins(StepperVertical, {
        props: makeProps({ currentStep: 3 }),
      });
      expect(wrapper.vm.isCompleted(0)).toBe(true);
      expect(wrapper.vm.isCompleted(1)).toBe(true);
      expect(wrapper.vm.isCompleted(2)).toBe(false);
    });

    it('should mark all steps as completed in editing mode', () => {
      const wrapper = mountWithPlugins(StepperVertical, {
        props: makeProps({ mode: 'editing', currentStep: 1 }),
      });
      expect(wrapper.vm.isCompleted(0)).toBe(true);
      expect(wrapper.vm.isCompleted(1)).toBe(true);
      expect(wrapper.vm.isCompleted(2)).toBe(true);
    });
  });

  describe('getCurrentStepLabel', () => {
    it('should return the title of the current step', () => {
      const wrapper = mountWithPlugins(StepperVertical, {
        props: makeProps({ currentStep: 2 }),
      });
      expect(wrapper.vm.getCurrentStepLabel()).toBe('Configuration');
    });

    it('should return empty string for out-of-bounds step', () => {
      const wrapper = mountWithPlugins(StepperVertical, {
        props: makeProps({ currentStep: 10 }),
      });
      expect(wrapper.vm.getCurrentStepLabel()).toBe('');
    });

    it('should return first step label by default (currentStep=1)', () => {
      const wrapper = mountWithPlugins(StepperVertical, { props: makeProps() });
      expect(wrapper.vm.getCurrentStepLabel()).toBe('Basic Info');
    });
  });

  describe('getStepIcon', () => {
    it('should return check_circle for completed steps in creating mode', () => {
      const wrapper = mountWithPlugins(StepperVertical, {
        props: makeProps({ currentStep: 3 }),
      });
      expect(wrapper.vm.getStepIcon(mockSteps[0], 0)).toBe('check_circle');
      expect(wrapper.vm.getStepIcon(mockSteps[1], 1)).toBe('check_circle');
    });

    it('should return the step own icon for the current step in creating mode', () => {
      const wrapper = mountWithPlugins(StepperVertical, {
        props: makeProps({ currentStep: 2 }),
      });
      expect(wrapper.vm.getStepIcon(mockSteps[1], 1)).toBe('settings');
    });

    it('should return the step own icon for future steps in creating mode', () => {
      const wrapper = mountWithPlugins(StepperVertical, {
        props: makeProps({ currentStep: 1 }),
      });
      expect(wrapper.vm.getStepIcon(mockSteps[2], 2)).toBe('check');
    });

    it('should always return step own icon in editing mode', () => {
      const wrapper = mountWithPlugins(StepperVertical, {
        props: makeProps({ mode: 'editing', currentStep: 3 }),
      });
      expect(wrapper.vm.getStepIcon(mockSteps[0], 0)).toBe('info');
      expect(wrapper.vm.getStepIcon(mockSteps[1], 1)).toBe('settings');
      expect(wrapper.vm.getStepIcon(mockSteps[2], 2)).toBe('check');
    });
  });

  describe('getStepIconStateClass', () => {
    it('should return "step-icon--active" for current step in creating mode', () => {
      const wrapper = mountWithPlugins(StepperVertical, {
        props: makeProps({ currentStep: 2 }),
      });
      expect(wrapper.vm.getStepIconStateClass(1)).toBe('step-icon--active');
    });

    it('should return "step-icon--completed" for past steps in creating mode', () => {
      const wrapper = mountWithPlugins(StepperVertical, {
        props: makeProps({ currentStep: 3 }),
      });
      expect(wrapper.vm.getStepIconStateClass(0)).toBe('step-icon--completed');
    });

    it('should return "step-icon--pending" for future steps in creating mode', () => {
      const wrapper = mountWithPlugins(StepperVertical, {
        props: makeProps({ currentStep: 1 }),
      });
      expect(wrapper.vm.getStepIconStateClass(2)).toBe('step-icon--pending');
    });

    it('should return "step-icon--active" for current step in editing mode', () => {
      const wrapper = mountWithPlugins(StepperVertical, {
        props: makeProps({ mode: 'editing', currentStep: 2 }),
      });
      expect(wrapper.vm.getStepIconStateClass(1)).toBe('step-icon--active');
    });

    it('should return "step-icon--completed" for non-current steps in editing mode', () => {
      const wrapper = mountWithPlugins(StepperVertical, {
        props: makeProps({ mode: 'editing', currentStep: 2 }),
      });
      expect(wrapper.vm.getStepIconStateClass(0)).toBe('step-icon--completed');
      expect(wrapper.vm.getStepIconStateClass(2)).toBe('step-icon--completed');
    });
  });

  describe('handleStepClick (emits)', () => {
    it('should emit "step-click" in editing mode', async () => {
      const wrapper = mountWithPlugins(StepperVertical, {
        props: makeProps({ mode: 'editing', currentStep: 1 }),
        stubs: quasarStubs,
      });
      const stepItems = wrapper.findAll('.step-item');
      await stepItems[2]!.trigger('click');
      expect(wrapper.emitted('step-click')).toBeTruthy();
      expect(wrapper.emitted('step-click')![0]).toEqual([3]);
    });

    it('should emit "step-click" when allowStepNavigation is true', async () => {
      const wrapper = mountWithPlugins(StepperVertical, {
        props: makeProps({ allowStepNavigation: true, currentStep: 1 }),
        stubs: quasarStubs,
      });
      const stepItems = wrapper.findAll('.step-item');
      await stepItems[1]!.trigger('click');
      expect(wrapper.emitted('step-click')![0]).toEqual([2]);
    });

    it('should emit "step-click" when clicking a previous step in creating mode', async () => {
      const wrapper = mountWithPlugins(StepperVertical, {
        props: makeProps({ currentStep: 3 }),
        stubs: quasarStubs,
      });
      const stepItems = wrapper.findAll('.step-item');
      await stepItems[0]!.trigger('click');
      expect(wrapper.emitted('step-click')![0]).toEqual([1]);
    });

    it('should NOT emit "step-click" when clicking a future step in creating mode', async () => {
      const wrapper = mountWithPlugins(StepperVertical, {
        props: makeProps({ currentStep: 1 }),
        stubs: quasarStubs,
      });
      const stepItems = wrapper.findAll('.step-item');
      await stepItems[2]!.trigger('click');
      expect(wrapper.emitted('step-click')).toBeFalsy();
    });

    it('should NOT emit "step-click" when clicking the current step in creating mode', async () => {
      const wrapper = mountWithPlugins(StepperVertical, {
        props: makeProps({ currentStep: 2 }),
        stubs: quasarStubs,
      });
      const stepItems = wrapper.findAll('.step-item');
      await stepItems[1]!.trigger('click');
      expect(wrapper.emitted('step-click')).toBeFalsy();
    });
  });

  describe('clickable CSS class', () => {
    it('should add clickable class to all steps in editing mode', () => {
      const wrapper = mountWithPlugins(StepperVertical, {
        props: makeProps({ mode: 'editing', currentStep: 1 }),
        stubs: quasarStubs,
      });
      const stepItems = wrapper.findAll('.step-item');
      stepItems.forEach((item) => {
        expect(item.classes()).toContain('clickable');
      });
    });

    it('should add clickable class only to previous steps in creating mode', () => {
      const wrapper = mountWithPlugins(StepperVertical, {
        props: makeProps({ currentStep: 2 }),
        stubs: quasarStubs,
      });
      const stepItems = wrapper.findAll('.step-item');
      expect(stepItems[0]!.classes()).toContain('clickable');
      expect(stepItems[1]!.classes()).not.toContain('clickable');
      expect(stepItems[2]!.classes()).not.toContain('clickable');
    });
  });
});
