import { describe, it, expect, beforeEach, vi } from 'vitest';
import { ref } from 'vue';

/** Capture lifecycle callbacks so we can trigger them manually */
const mountedCallbacks: (() => void)[] = [];
const unmountedCallbacks: (() => void)[] = [];

vi.mock('vue', async () => {
  // eslint-disable-next-line @typescript-eslint/consistent-type-imports
  const actual = await vi.importActual<typeof import('vue')>('vue');
  return {
    ...actual,
    onMounted: (cb: () => void) => mountedCallbacks.push(cb),
    onUnmounted: (cb: () => void) => unmountedCallbacks.push(cb),
  };
});

import { useStepperNavigation } from './useStepperNavigation';

/**
 * Helper to dispatch a KeyboardEvent on document with an optional target override.
 */
function pressKey(key: string, target?: EventTarget): void {
  const event = new KeyboardEvent('keydown', { key, bubbles: true });
  if (target) {
    Object.defineProperty(event, 'target', { value: target });
  }
  document.dispatchEvent(event);
}

describe('useStepperNavigation', () => {
  let currentStep: ReturnType<typeof ref<number>>;
  let changeStep: ReturnType<typeof vi.fn>;

  beforeEach(() => {
    currentStep = ref(1);
    changeStep = vi.fn((step: number) => {
      currentStep.value = step;
    });
    mountedCallbacks.length = 0;
    unmountedCallbacks.length = 0;
  });

  function setup(totalSteps = 5) {
    useStepperNavigation({
      // eslint-disable-next-line @typescript-eslint/consistent-type-imports
      currentStep: currentStep as import('vue').Ref<number>,
      totalSteps,
      changeStep: changeStep as unknown as (step: number) => void,
    });

    // Simulate component mount
    mountedCallbacks.forEach((cb) => cb());

    return {
      unmount: () => unmountedCallbacks.forEach((cb) => cb()),
    };
  }

  describe('ArrowRight navigation', () => {
    it('advances to next step when not at last step', () => {
      const { unmount } = setup(5);

      currentStep.value = 2;
      pressKey('ArrowRight');

      expect(changeStep).toHaveBeenCalledWith(3);
      unmount();
    });

    it('does not advance past the last step', () => {
      const { unmount } = setup(3);

      currentStep.value = 3;
      pressKey('ArrowRight');

      expect(changeStep).not.toHaveBeenCalled();
      unmount();
    });
  });

  describe('ArrowLeft navigation', () => {
    it('goes to previous step when not at first step', () => {
      const { unmount } = setup(5);

      currentStep.value = 3;
      pressKey('ArrowLeft');

      expect(changeStep).toHaveBeenCalledWith(2);
      unmount();
    });

    it('does not go below step 1', () => {
      const { unmount } = setup(5);

      currentStep.value = 1;
      pressKey('ArrowLeft');

      expect(changeStep).not.toHaveBeenCalled();
      unmount();
    });
  });

  describe('input focus guard', () => {
    it('ignores keydown when target is an input element', () => {
      const { unmount } = setup(5);

      currentStep.value = 2;
      const input = document.createElement('input');
      pressKey('ArrowRight', input);

      expect(changeStep).not.toHaveBeenCalled();
      unmount();
    });

    it('ignores keydown when target is a textarea', () => {
      const { unmount } = setup(5);

      currentStep.value = 2;
      const textarea = document.createElement('textarea');
      pressKey('ArrowRight', textarea);

      expect(changeStep).not.toHaveBeenCalled();
      unmount();
    });

    it('ignores keydown when target is a select element', () => {
      const { unmount } = setup(5);

      currentStep.value = 2;
      const select = document.createElement('select');
      pressKey('ArrowRight', select);

      expect(changeStep).not.toHaveBeenCalled();
      unmount();
    });
  });

  describe('irrelevant keys', () => {
    it('does not react to unrelated keys', () => {
      const { unmount } = setup(5);

      currentStep.value = 2;
      pressKey('Enter');
      pressKey('Space');
      pressKey('ArrowUp');
      pressKey('ArrowDown');

      expect(changeStep).not.toHaveBeenCalled();
      unmount();
    });
  });

  describe('cleanup', () => {
    it('removes event listener on unmount', () => {
      const { unmount } = setup(5);

      unmount();
      currentStep.value = 2;
      pressKey('ArrowRight');

      expect(changeStep).not.toHaveBeenCalled();
    });
  });
});
