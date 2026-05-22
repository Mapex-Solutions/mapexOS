import { test as base } from '@playwright/test';
import type { Page } from '@playwright/test';
import { login } from '@mapexos/e2e-helpers';

/**
 * Custom test fixture that provides a pre-authenticated page
 * Opens the login page, fills credentials via UI, and waits for redirect
 */
export const test = base.extend<{ authenticatedPage: Page }>({
  authenticatedPage: async ({ page }, use) => {
    await login(page, 'vendor@mapex.global', 'Mapex@123');
    await use(page);
  },
});

export { expect } from '@playwright/test';
