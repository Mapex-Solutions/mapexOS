import type { Page } from '@playwright/test';
import { expect } from '@playwright/test';

/**
 * Wait for all Quasar loading indicators to disappear
 *
 * @param {Page} page - Playwright page instance
 * @param {number} timeout - Max time to wait (ms)
 * @returns {Promise<void>}
 */
export async function waitForLoading(page: Page, timeout: number = 15000): Promise<void> {
  // Wait for Quasar's q-loading overlay to disappear
  const loading = page.locator('.q-loading');
  await expect(loading).toBeHidden({ timeout });
}

/**
 * Wait for Quasar spinner components to disappear
 *
 * @param {Page} page - Playwright page instance
 * @param {number} timeout - Max time to wait (ms)
 * @returns {Promise<void>}
 */
export async function waitForSpinner(page: Page, timeout: number = 15000): Promise<void> {
  const spinner = page.locator('.q-spinner');
  await expect(spinner).toBeHidden({ timeout });
}

/**
 * Wait for the page to be fully loaded (no loading bars or spinners)
 *
 * @param {Page} page - Playwright page instance
 * @param {number} timeout - Max time to wait (ms)
 * @returns {Promise<void>}
 */
export async function waitForPageReady(page: Page, timeout: number = 15000): Promise<void> {
  // Wait for Quasar's loading bar to finish
  const loadingBar = page.locator('.q-loading-bar');
  await expect(loadingBar).toBeHidden({ timeout });

  // Also wait for any inner content loaders
  const innerLoader = page.locator('.q-inner-loading');
  const count = await innerLoader.count();
  if (count > 0) {
    await expect(innerLoader.first()).toBeHidden({ timeout });
  }
}
