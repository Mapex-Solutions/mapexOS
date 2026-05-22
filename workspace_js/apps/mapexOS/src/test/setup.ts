/**
 * Global test setup for Vitest.
 *
 * Configures:
 * - Quasar component stubs (q-btn, q-card, etc.)
 * - i18n mock (returns key as translation)
 * - API service mock
 * - Vue Router mock
 */
import { vi } from 'vitest';
import { config } from '@vue/test-utils';

/**
 * Mock Quasar components as stubs.
 * Prevents "Unknown custom element" warnings in tests.
 */
config.global.stubs = {
  'q-btn': true,
  'q-icon': true,
  'q-card': true,
  'q-card-section': true,
  'q-card-actions': true,
  'q-input': true,
  'q-select': true,
  'q-toggle': true,
  'q-checkbox': true,
  'q-dialog': true,
  'q-drawer': true,
  'q-separator': true,
  'q-space': true,
  'q-badge': true,
  'q-chip': true,
  'q-avatar': true,
  'q-banner': true,
  'q-list': true,
  'q-item': true,
  'q-item-section': true,
  'q-item-label': true,
  'q-tooltip': true,
  'q-tab': true,
  'q-tabs': true,
  'q-tab-panel': true,
  'q-tab-panels': true,
  'q-page': true,
  'q-field': true,
  'q-form': true,
  'q-spinner': true,
  'q-scroll-area': true,
  'q-menu': true,
  'q-pagination': true,
  'q-expansion-item': true,
  'q-stepper': true,
  'q-step': true,
  'q-stepper-navigation': true,
  'q-popup-proxy': true,
  'q-date': true,
  'q-time': true,
  'q-table': true,
  'q-file': true,
  'q-toolbar': true,
  'q-toolbar-title': true,
  'q-no-ssr': true,
  'q-resize-observer': true,
};

/**
 * Mock Quasar directives
 */
config.global.directives = {
  'close-popup': {},
  'ripple': {},
  'permission': {},
};

/**
 * Mock API services globally.
 * Individual tests can override with vi.mocked().
 */
vi.mock('@services/mapex', () => ({
  apis: {
    mapexOS: {
      users: { list: vi.fn(), getById: vi.fn(), create: vi.fn(), update: vi.fn(), delete: vi.fn() },
      roles: { list: vi.fn(), getById: vi.fn(), create: vi.fn(), update: vi.fn(), delete: vi.fn() },
      groups: { list: vi.fn(), getById: vi.fn(), create: vi.fn(), update: vi.fn(), delete: vi.fn() },
      organizations: { list: vi.fn(), getById: vi.fn() },
    },
    router: {
      routegroup: { list: vi.fn(), getById: vi.fn() },
    },
    workflows: {
      definition: { list: vi.fn(), getById: vi.fn(), create: vi.fn() },
      instance: { list: vi.fn(), getById: vi.fn() },
      plugin: { list: vi.fn(), getEnabled: vi.fn(), create: vi.fn(), delete: vi.fn() },
      credential: { list: vi.fn(), create: vi.fn(), test: vi.fn(), delete: vi.fn(), loadOptions: vi.fn() },
    },
  },
}));

/**
 * Mock vue-router
 */
vi.mock('vue-router', () => ({
  useRouter: () => ({
    push: vi.fn(),
    replace: vi.fn(),
    back: vi.fn(),
    currentRoute: { value: { params: {}, query: {} } },
  }),
  useRoute: () => ({
    params: {},
    query: {},
    path: '/',
    name: 'test',
  }),
}));

/**
 * Mock i18n — returns the key as the translated value
 */
vi.mock('vue-i18n', () => ({
  useI18n: () => ({
    t: (key: string) => key,
    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    te: (key: string) => false,
    locale: { value: 'en-US' },
  }),
  createI18n: () => ({
    global: {
      t: (key: string) => key,
      // eslint-disable-next-line @typescript-eslint/no-unused-vars
      te: (key: string) => false,
      locale: { value: 'en-US' },
    },
  }),
}));

/**
 * Mock monaco-theme — prevents transitive monaco-editor resolution errors.
 * ScriptEditorDialog -> monaco-theme -> monaco-setup -> monaco-editor workers.
 */
vi.mock('@utils/monaco-theme', () => ({
  registerMapexMonacoThemes: () => {},
  getMapexMonacoTheme: () => 'vs-dark',
  applyMapexMonacoTheme: () => {},
}));

/**
 * Mock boot/i18n export
 */
vi.mock('src/boot/i18n', () => ({
  i18nInstance: {
    global: {
      t: (key: string) => key,
      // eslint-disable-next-line @typescript-eslint/no-unused-vars
      te: (key: string) => false,
      locale: { value: 'en-US' },
    },
  },
}));
