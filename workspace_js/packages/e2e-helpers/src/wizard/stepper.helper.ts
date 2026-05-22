import type { Page } from '@playwright/test';

/**
 * Click the "Next" button in a FormCard wizard
 *
 * @param {Page} page - Playwright page instance
 * @returns {Promise<void>}
 */
export async function clickNext(page: Page): Promise<void> {
  const btn = page.getByTestId('wizard-next-btn');
  await btn.waitFor({ state: 'visible' });
  await btn.click();
}

/**
 * Click the "Previous" button in a FormCard wizard
 *
 * @param {Page} page - Playwright page instance
 * @returns {Promise<void>}
 */
export async function clickPrevious(page: Page): Promise<void> {
  const btn = page.getByTestId('wizard-previous-btn');
  await btn.waitFor({ state: 'visible' });
  await btn.click();
}

/**
 * Click the "Save" button in a FormCard wizard
 *
 * @param {Page} page - Playwright page instance
 * @returns {Promise<void>}
 */
export async function clickSave(page: Page): Promise<void> {
  const btn = page.getByTestId('wizard-save-btn');
  await btn.waitFor({ state: 'visible' });
  await btn.click();
}

/**
 * Check if the Next button is disabled
 *
 * @param {Page} page - Playwright page instance
 * @returns {Promise<boolean>} Whether the Next button is disabled
 */
export async function isNextDisabled(page: Page): Promise<boolean> {
  const btn = page.getByTestId('wizard-next-btn');
  return btn.isDisabled();
}

/**
 * Check if the Save button is disabled
 *
 * @param {Page} page - Playwright page instance
 * @returns {Promise<boolean>} Whether the Save button is disabled
 */
export async function isSaveDisabled(page: Page): Promise<boolean> {
  const btn = page.getByTestId('wizard-save-btn');
  return btn.isDisabled();
}
