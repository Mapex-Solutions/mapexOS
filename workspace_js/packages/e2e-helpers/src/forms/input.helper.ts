import type { Page, Locator } from '@playwright/test';

/**
 * Fill a Quasar q-input identified by data-testid.
 * In Quasar 2.18+, data-testid is placed directly on the native <input> element,
 * so getByTestId() returns the input itself.
 *
 * @param {Page} page - Playwright page instance
 * @param {string} testId - data-testid attribute value
 * @param {string} value - Value to fill
 * @returns {Promise<void>}
 */
export async function fillInput(page: Page, testId: string, value: string): Promise<void> {
  const input = page.getByTestId(testId);
  await input.click();
  await input.fill(value);
}

/**
 * Clear a Quasar q-input identified by data-testid
 *
 * @param {Page} page - Playwright page instance
 * @param {string} testId - data-testid attribute value
 * @returns {Promise<void>}
 */
export async function clearInput(page: Page, testId: string): Promise<void> {
  const input = page.getByTestId(testId);
  await input.click();
  await input.fill('');
}

/**
 * Get the current value of a Quasar q-input by data-testid
 *
 * @param {Page} page - Playwright page instance
 * @param {string} testId - data-testid attribute value
 * @returns {Promise<string>} Current input value
 */
export async function getInputValue(page: Page, testId: string): Promise<string> {
  return (await page.getByTestId(testId).inputValue()) || '';
}

/**
 * Get the native input locator for a Quasar q-input by data-testid.
 * In Quasar 2.18+, data-testid is on the native <input> directly.
 *
 * @param {Page} page - Playwright page instance
 * @param {string} testId - data-testid attribute value
 * @returns {Locator} Native input locator
 */
export function getNativeInput(page: Page, testId: string): Locator {
  return page.getByTestId(testId);
}

/**
 * Fill a Quasar q-input using a Locator directly
 *
 * @param {Locator} locator - Playwright locator
 * @param {string} value - Value to fill
 * @returns {Promise<void>}
 */
export async function fillInputByLocator(locator: Locator, value: string): Promise<void> {
  await locator.click();
  await locator.fill(value);
}
