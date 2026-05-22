import { ref, computed } from 'vue';
import { useRouter } from 'vue-router';
import { driver, type DriveStep, type Config, type Driver } from 'driver.js';
import 'driver.js/dist/driver.css';
import { useOnboardingTranslations } from '@composables/i18n';
import { useAuthStore } from '@stores/auth';
import { useLogger } from '@composables/useLogger';
import { apis } from '@services/mapex';

/**
 * Onboarding tour composable
 *
 * Encapsulates all Driver.js logic for the guided tour.
 * - Defines tour steps targeting layout elements by ID
 * - Provides startTour() to trigger manually or automatically
 * - Marks tour as completed via PATCH /api/v1/users/me/tour (sets startTour=false)
 * - Uses i18n translations for all texts
 *
 * @example
 * ```ts
 * const { startTour, isActive } = useOnboarding();
 *
 * // Auto-start on mount
 * onMounted(() => {
 *   if (authStore.user?.startTour) startTour();
 * });
 *
 * // Manual trigger from menu
 * function handleStartTour() { startTour(); }
 * ```
 */
export function useOnboarding() {
  const t = useOnboardingTranslations();
  const router = useRouter();
  const authStore = useAuthStore();
  const logger = useLogger('Onboarding');

  /** Whether the tour is currently active */
  const isActive = ref(false);

  /** Driver.js instance reference for external access */
  let driverInstance: Driver | null = null;

  /**
   * Tour steps configuration
   * Each step targets an element by ID in the layout.
   * The last step (userMenu) navigates to /users?tour=true on "Continue".
   */
  const tourSteps = computed<DriveStep[]>(() => [
    {
      element: '#sidebar-menu',
      popover: {
        title: t.steps.sidebar.title.value,
        description: t.steps.sidebar.description.value,
        side: 'right',
        align: 'start',
      },
    },
    {
      element: '#header-breadcrumbs',
      popover: {
        title: t.steps.breadcrumbs.title.value,
        description: t.steps.breadcrumbs.description.value,
        side: 'bottom',
        align: 'start',
      },
    },
    {
      element: '#header-org-selector',
      popover: {
        title: t.steps.orgSelector.title.value,
        description: t.steps.orgSelector.description.value,
        side: 'bottom',
        align: 'center',
      },
    },
    {
      element: '#header-lang-selector',
      popover: {
        title: t.steps.langSelector.title.value,
        description: t.steps.langSelector.description.value,
        side: 'bottom',
        align: 'center',
      },
    },
    {
      element: '#header-user-menu',
      popover: {
        title: t.steps.userMenu.title.value,
        description: t.steps.userMenu.description.value,
        side: 'bottom',
        align: 'end',
        onNextClick: () => {
          void completeTour();
          driverInstance?.destroy();
          void router.push('/users?tour=true');
        },
      },
    },
  ]);

  /**
   * Mark tour as completed by calling the backend API.
   * Updates the user's startTour flag to false and syncs with auth store.
   */
  async function completeTour(): Promise<void> {
    isActive.value = false;
    logger.debug('Tour completed');

    try {
      await apis.mapexOS?.users.disableMyTour();

      // Sync the auth store user object
      if (authStore.user) {
        authStore.user.startTour = false;
        authStore.persistTokens();
      }
    } catch (error) {
      logger.error('Failed to disable tour via API', error);
    }
  }

  /**
   * Start the onboarding tour
   * Creates a new Driver.js instance and drives through steps.
   * After the last step, navigates to /users?tour=true to continue the page tour.
   */
  function startTour(): void {
    logger.debug('Starting tour');

    const config: Config = {
      showProgress: true,
      animate: true,
      allowClose: true,
      overlayColor: 'rgba(0, 0, 0, 0.5)',
      stagePadding: 8,
      stageRadius: 8,
      popoverOffset: 12,
      nextBtnText: t.buttons.next.value,
      prevBtnText: t.buttons.previous.value,
      doneBtnText: t.buttons.continue.value,
      steps: tourSteps.value,
      onDestroyStarted: () => {
        if (isActive.value) void completeTour();
        driverInstance?.destroy();
      },
    };

    driverInstance = driver(config);
    isActive.value = true;
    driverInstance.drive();
  }

  /**
   * Check if the tour should auto-start based on backend startTour flag
   *
   * @returns {boolean} Whether the tour should start automatically
   */
  function shouldAutoStart(): boolean {
    return authStore.user?.startTour === true;
  }

  return {
    startTour,
    isActive,
    shouldAutoStart,
  };
}
