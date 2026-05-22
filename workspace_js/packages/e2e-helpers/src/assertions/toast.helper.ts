import type { Page } from '@playwright/test';
import { expect } from '@playwright/test';

/**
 * Assert that a Quasar toast notification appears with the expected text
 *
 * @param {Page} page - Playwright page instance
 * @param {string} text - Expected text content in the toast
 * @param {object} options - Additional options
 * @param {'positive' | 'negative' | 'warning' | 'info'} options.type - Toast type to match
 * @param {number} options.timeout - Max time to wait for the toast (ms)
 * @returns {Promise<void>}
 */
export async function expectToast(
  page: Page,
  text: string,
  options: {
    type?: 'positive' | 'negative' | 'warning' | 'info';
    timeout?: number;
  } = {},
): Promise<void> {
  const { type, timeout = 10000 } = options;

  // Quasar renders toasts as .q-notification elements
  let selector = '.q-notification';
  if (type) {
    selector += `.bg-${type}`;
  }

  const toast = page.locator(selector).filter({ hasText: text });
  await expect(toast).toBeVisible({ timeout });
}

/**
 * Wait for a toast to disappear
 *
 * @param {Page} page - Playwright page instance
 * @param {string} text - Text content of the toast to wait for
 * @param {number} timeout - Max time to wait (ms)
 * @returns {Promise<void>}
 */
export async function waitForToastDismiss(
  page: Page,
  text: string,
  timeout: number = 10000,
): Promise<void> {
  const toast = page.locator('.q-notification').filter({ hasText: text });
  await expect(toast).toBeHidden({ timeout });
}
