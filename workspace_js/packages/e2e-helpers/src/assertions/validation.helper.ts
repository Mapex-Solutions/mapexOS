import type { Page } from '@playwright/test';
import { expect } from '@playwright/test';

/**
 * Assert that a Quasar validation error message is visible
 *
 * @param {Page} page - Playwright page instance
 * @param {string} errorText - Expected validation error text
 * @param {number} timeout - Max time to wait (ms)
 * @returns {Promise<void>}
 */
export async function expectValidationError(
  page: Page,
  errorText: string,
  timeout: number = 5000,
): Promise<void> {
  // Quasar renders validation errors in .q-field__messages elements
  const errorMessage = page.locator('.q-field__messages').filter({ hasText: errorText });
  await expect(errorMessage).toBeVisible({ timeout });
}

/**
 * Assert that a specific field (by data-testid) has a validation error
 *
 * @param {Page} page - Playwright page instance
 * @param {string} testId - data-testid of the field
 * @param {string} errorText - Expected error text (optional, just checks for error state if omitted)
 * @returns {Promise<void>}
 */
export async function expectFieldError(
  page: Page,
  testId: string,
  errorText?: string,
): Promise<void> {
  const field = page.getByTestId(testId);

  // Check that the field has the error class
  await expect(field.locator('.q-field--error')).toBeVisible();

  if (errorText) {
    const message = field.locator('.q-field__messages').filter({ hasText: errorText });
    await expect(message).toBeVisible();
  }
}

/**
 * Assert that no validation errors are visible on the page
 *
 * @param {Page} page - Playwright page instance
 * @returns {Promise<void>}
 */
export async function expectNoValidationErrors(page: Page): Promise<void> {
  const errors = page.locator('.q-field--error');
  await expect(errors).toHaveCount(0);
}
