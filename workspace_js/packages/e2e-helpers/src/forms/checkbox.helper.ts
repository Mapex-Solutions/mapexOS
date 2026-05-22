import type { Page } from '@playwright/test';

/**
 * Toggle a Quasar q-checkbox identified by data-testid
 *
 * @param {Page} page - Playwright page instance
 * @param {string} testId - data-testid attribute value
 * @returns {Promise<void>}
 */
export async function toggleCheckbox(page: Page, testId: string): Promise<void> {
  await page.getByTestId(testId).click();
}

/**
 * Check whether a Quasar q-checkbox is checked
 *
 * @param {Page} page - Playwright page instance
 * @param {string} testId - data-testid attribute value
 * @returns {Promise<boolean>} Whether the checkbox is checked
 */
export async function isChecked(page: Page, testId: string): Promise<boolean> {
  const checkbox = page.getByTestId(testId);
  // Quasar checkboxes use aria-checked attribute
  const ariaChecked = await checkbox.getAttribute('aria-checked');
  if (ariaChecked !== null) {
    return ariaChecked === 'true';
  }
  // Fallback: check for the "q-checkbox--checked" class
  const classList = await checkbox.getAttribute('class');
  return classList?.includes('q-checkbox--checked') ?? false;
}

/**
 * Set a Quasar q-checkbox to a specific state
 *
 * @param {Page} page - Playwright page instance
 * @param {string} testId - data-testid attribute value
 * @param {boolean} checked - Desired checked state
 * @returns {Promise<void>}
 */
export async function setCheckbox(page: Page, testId: string, checked: boolean): Promise<void> {
  const currentState = await isChecked(page, testId);
  if (currentState !== checked) {
    await toggleCheckbox(page, testId);
  }
}
