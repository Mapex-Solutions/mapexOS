/**
 * Page Tour Composable
 *
 * Generic composable for page-level and multi-page guided tours.
 * Uses Driver.js under the hood, with support for:
 * - Per-page tours with custom steps
 * - Multi-page transitions via sessionStorage
 * - Auto-start detection via route query param (?tour=true)
 * - Reuses onboarding translations for button labels
 *
 * @example
 * ```ts
 * const { startTour, isTourMode } = usePageTour({
 *   tourId: 'users-list',
 *   steps: [
 *     { element: '#header', title: 'Header', description: 'Page header' },
 *   ],
 * });
 * ```
 */

import type { DriveStep, Config } from 'driver.js';
import type { PageTourConfig, TourSessionState } from '../interfaces';

import { ref, computed, onMounted, onUnmounted } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { driver } from 'driver.js';
import 'driver.js/dist/driver.css';

import { useOnboardingTranslations } from '@composables/i18n';
import { useLogger } from '@composables/useLogger';

/** SessionStorage key for multi-page tour state */
const TOUR_SESSION_KEY = 'mapex-tour-session';

/** TTL for tour session state (5 minutes) */
const TOUR_SESSION_TTL_MS = 5 * 60 * 1000;

/**
 * Page tour composable
 * Provides tour start/stop, auto-start detection, and multi-page transition support
 *
 * @param {PageTourConfig} config - Tour configuration with steps and optional transition
 * @returns Tour control functions and reactive state
 */
export function usePageTour(config: PageTourConfig) {
  const onboardingT = useOnboardingTranslations();
  const router = useRouter();
  const route = useRoute();
  const logger = useLogger(`PageTour:${config.tourId}`);

  /** Whether the tour is currently active */
  const isActive = ref(false);

  /** Driver.js instance reference for cleanup */
  let driverInstance: ReturnType<typeof driver> | null = null;

  /**
   * Whether the current page is in tour/demo mode
   * Detected via ?tour=true query parameter
   */
  const isTourMode = computed(() => route.query.tour === 'true');

  /**
   * Read and validate tour session state from sessionStorage
   *
   * @returns {TourSessionState | null} Valid session state or null if expired/missing
   */
  function readSessionState(): TourSessionState | null {
    try {
      const raw = sessionStorage.getItem(TOUR_SESSION_KEY);
      if (!raw) return null;

      const state: TourSessionState = JSON.parse(raw);
      const elapsed = Date.now() - state.timestamp;

      // Check TTL expiration
      if (elapsed > TOUR_SESSION_TTL_MS) {
        sessionStorage.removeItem(TOUR_SESSION_KEY);
        logger.debug('Tour session expired');
        return null;
      }

      return state;
    } catch {
      sessionStorage.removeItem(TOUR_SESSION_KEY);
      return null;
    }
  }

  /**
   * Save tour session state to sessionStorage for multi-page transition
   *
   * @param {string} flowId - Tour flow identifier
   */
  function saveSessionState(flowId: string): void {
    const state: TourSessionState = {
      flowId,
      shouldAutoStart: true,
      timestamp: Date.now(),
    };
    sessionStorage.setItem(TOUR_SESSION_KEY, JSON.stringify(state));
  }

  /**
   * Clear tour session state from sessionStorage
   */
  function clearSessionState(): void {
    sessionStorage.removeItem(TOUR_SESSION_KEY);
  }

  /**
   * Build Driver.js steps from PageTourStep config
   * Calls the steps getter to resolve translations at tour-start time
   *
   * @returns {DriveStep[]} Array of Driver.js compatible steps
   */
  function buildDriverSteps(): DriveStep[] {
    const resolvedSteps = config.steps();
    return resolvedSteps.map((step, index) => {
      const driveStep: DriveStep = {
        element: step.element,
        popover: {
          title: step.title,
          description: step.description,
          side: step.side || 'bottom',
          align: step.align || 'start',
        },
      };

      // If step has custom onNextClick handler (e.g., open/close drawer)
      if (step.onNextClick) {
        const callback = step.onNextClick;
        driveStep.popover!.onNextClick = () => {
          callback(() => driverInstance?.moveNext());
        };
      }

      // If this step triggers a multi-page transition (takes priority over step onNextClick)
      if (config.transition && index === config.transition.triggerAtStep) {
        driveStep.popover!.onNextClick = () => {
          logger.debug('Tour transition triggered', {
            from: config.tourId,
            to: config.transition!.targetRoute,
          });

          // Save session state for destination page
          saveSessionState(config.tourId);

          // Destroy current driver instance
          if (driverInstance) {
            driverInstance.destroy();
            driverInstance = null;
          }
          isActive.value = false;

          // Navigate to target route with tour query param
          void router.push(`${config.transition!.targetRoute}?tour=true`);
        };
      }

      // If step has onHighlightStarted callback
      if (step.onHighlightStarted) {
        driveStep.popover!.onPopoverRender = () => {
          step.onHighlightStarted!();
        };
      }

      return driveStep;
    });
  }

  /**
   * Start the tour
   * Creates a new Driver.js instance and drives through configured steps
   */
  function startTour(): void {
    logger.debug('Starting tour:', config.tourId);

    const steps = buildDriverSteps();

    if (steps.length === 0) {
      logger.debug('No steps configured, skipping tour');
      return;
    }

    const driverConfig: Config = {
      showProgress: true,
      animate: true,
      allowClose: true,
      overlayColor: 'rgba(0, 0, 0, 0.5)',
      stagePadding: 8,
      stageRadius: 8,
      popoverOffset: 12,
      nextBtnText: onboardingT.buttons.next.value,
      prevBtnText: onboardingT.buttons.previous.value,
      doneBtnText: onboardingT.buttons.finish.value,
      steps,
      onDestroyStarted: () => {
        logger.debug('Tour ending:', config.tourId);
        isActive.value = false;
        clearSessionState();

        if (driverInstance) {
          driverInstance.destroy();
          driverInstance = null;
        }

        // Call onTourEnd callback (e.g., navigate back to list)
        if (config.onTourEnd) {
          config.onTourEnd();
        }
      },
    };

    driverInstance = driver(driverConfig);
    isActive.value = true;
    driverInstance.drive();
  }

  /**
   * Stop the tour and clean up
   */
  function stopTour(): void {
    if (driverInstance) {
      driverInstance.destroy();
      driverInstance = null;
    }
    isActive.value = false;
    clearSessionState();
  }

  /** LIFECYCLE HOOKS */
  onMounted(() => {
    // Auto-start logic: check route query param or session state
    const shouldAutoStart =
      config.autoStart === true ||
      isTourMode.value ||
      readSessionState()?.shouldAutoStart === true;

    if (shouldAutoStart && config.steps().length > 0) {
      // Delay to allow DOM elements to render
      setTimeout(() => {
        clearSessionState();
        startTour();
      }, 500);
    }
  });

  onUnmounted(() => {
    // Cleanup on component unmount
    if (driverInstance) {
      driverInstance.destroy();
      driverInstance = null;
    }
    isActive.value = false;
  });

  return {
    startTour,
    stopTour,
    isActive,
    isTourMode,
  };
}
