import { describe, it, expect, vi, beforeEach } from 'vitest';
import { setActivePinia, createPinia } from 'pinia';

/**
 * Mock driver.js
 */
const mockDrive = vi.fn();
const mockDestroy = vi.fn();
const mockDriver = vi.fn(() => ({
  drive: mockDrive,
  destroy: mockDestroy,
}));

vi.mock('driver.js', () => ({
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  driver: (_arg: unknown) => mockDriver(),
}));

vi.mock('driver.js/dist/driver.css', () => ({}));

/**
 * Mock vue-router
 */
const mockPush = vi.fn();

vi.mock('vue-router', () => ({
  useRouter: () => ({ push: mockPush }),
  useRoute: () => ({ query: {} }),
}));

/**
 * Mock composables
 */
vi.mock('@composables/i18n', () => ({
  useOnboardingTranslations: () => ({
    steps: {
      sidebar: { title: { value: 'Sidebar' }, description: { value: 'Nav menu' } },
      breadcrumbs: { title: { value: 'Breadcrumbs' }, description: { value: 'Current path' } },
      orgSelector: { title: { value: 'Org Selector' }, description: { value: 'Switch org' } },
      langSelector: { title: { value: 'Language' }, description: { value: 'Switch lang' } },
      userMenu: { title: { value: 'User Menu' }, description: { value: 'Profile settings' } },
    },
    buttons: {
      next: { value: 'Next' },
      previous: { value: 'Previous' },
      continue: { value: 'Continue' },
      finish: { value: 'Finish' },
    },
  }),
}));

vi.mock('@composables/useLogger', () => ({
  useLogger: () => ({
    debug: vi.fn(),
    error: vi.fn(),
    info: vi.fn(),
    warn: vi.fn(),
  }),
}));

/**
 * Mock API
 */
const mockDisableMyTour = vi.fn().mockResolvedValue(undefined);

vi.mock('@services/mapex', () => ({
  apis: {
    mapexOS: {
      users: {
        disableMyTour: () => mockDisableMyTour(),
      },
    },
  },
}));

describe('useOnboarding', () => {
  beforeEach(() => {
    setActivePinia(createPinia());
    mockDrive.mockReset();
    mockDestroy.mockReset();
    mockDriver.mockClear();
    mockPush.mockReset();
    mockDisableMyTour.mockReset().mockResolvedValue(undefined);
  });

  async function importFresh() {
    const mod = await import('./useOnboarding');
    return mod.useOnboarding;
  }

  it('returns startTour, isActive, shouldAutoStart', async () => {
    const useOnboarding = await importFresh();
    const result = useOnboarding();

    expect(result).toHaveProperty('startTour');
    expect(result).toHaveProperty('isActive');
    expect(result).toHaveProperty('shouldAutoStart');
  });

  it('isActive starts as false', async () => {
    const useOnboarding = await importFresh();
    const { isActive } = useOnboarding();

    expect(isActive.value).toBe(false);
  });

  describe('startTour', () => {
    it('creates driver instance and calls drive()', async () => {
      const useOnboarding = await importFresh();
      const { startTour, isActive } = useOnboarding();

      startTour();

      expect(mockDriver).toHaveBeenCalled();
      expect(mockDrive).toHaveBeenCalled();
      expect(isActive.value).toBe(true);
    });

    it('passes 5 tour steps to driver config', async () => {
      const useOnboarding = await importFresh();
      const { startTour } = useOnboarding();

      startTour();

      const config = (mockDriver.mock.calls as unknown as [unknown[]])[0][0] as Record<string, unknown>;
      expect(config.steps).toHaveLength(5);
    });

    it('configures correct button labels', async () => {
      const useOnboarding = await importFresh();
      const { startTour } = useOnboarding();

      startTour();

      const config = (mockDriver.mock.calls as unknown as [unknown[]])[0][0] as Record<string, unknown>;
      expect(config.nextBtnText).toBe('Next');
      expect(config.prevBtnText).toBe('Previous');
      expect(config.doneBtnText).toBe('Continue');
    });

    it('targets correct DOM elements', async () => {
      const useOnboarding = await importFresh();
      const { startTour } = useOnboarding();

      startTour();

      const config = (mockDriver.mock.calls as unknown as [unknown[]])[0][0] as Record<string, unknown>;
      const elements = (config.steps as { element: string }[]).map((s) => s.element);

      expect(elements).toContain('#sidebar-menu');
      expect(elements).toContain('#header-breadcrumbs');
      expect(elements).toContain('#header-org-selector');
      expect(elements).toContain('#header-lang-selector');
      expect(elements).toContain('#header-user-menu');
    });
  });

  describe('shouldAutoStart', () => {
    it('returns true when user.startTour is true', async () => {
      const useOnboarding = await importFresh();

      // Set auth store user
      const { useAuthStore } = await import('@stores/auth');
      const authStore = useAuthStore();
      authStore.user = { startTour: true } as any;

      const { shouldAutoStart } = useOnboarding();

      expect(shouldAutoStart()).toBe(true);
    });

    it('returns false when user.startTour is false', async () => {
      const useOnboarding = await importFresh();

      const { useAuthStore } = await import('@stores/auth');
      const authStore = useAuthStore();
      authStore.user = { startTour: false } as any;

      const { shouldAutoStart } = useOnboarding();

      expect(shouldAutoStart()).toBe(false);
    });

    it('returns false when no user exists', async () => {
      const useOnboarding = await importFresh();

      const { useAuthStore } = await import('@stores/auth');
      const authStore = useAuthStore();
      authStore.user = null as any;

      const { shouldAutoStart } = useOnboarding();

      expect(shouldAutoStart()).toBe(false);
    });
  });
});
