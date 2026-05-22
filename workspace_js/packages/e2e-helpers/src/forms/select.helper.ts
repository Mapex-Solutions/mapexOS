import type { Page } from '@playwright/test';

/**
 * Select an option from a Quasar q-select dropdown by data-testid
 *
 * @param {Page} page - Playwright page instance
 * @param {string} testId - data-testid attribute value of the q-select
 * @param {string} optionText - Visible text of the option to select
 * @returns {Promise<void>}
 */
export async function selectOption(
  page: Page,
  testId: string,
  optionText: string,
): Promise<void> {
  const select = page.getByTestId(testId);
  await select.click();
  // Quasar renders dropdown options in a q-menu portal
  await page.locator('.q-menu .q-item').filter({ hasText: optionText }).click();
}

/**
 * Select multiple options from a Quasar q-select (multiple mode)
 *
 * @param {Page} page - Playwright page instance
 * @param {string} testId - data-testid attribute value of the q-select
 * @param {string[]} optionTexts - Array of visible option texts to select
 * @returns {Promise<void>}
 */
export async function selectMultipleOptions(
  page: Page,
  testId: string,
  optionTexts: string[],
): Promise<void> {
  const select = page.getByTestId(testId);
  await select.click();

  for (const text of optionTexts) {
    await page.locator('.q-menu .q-item').filter({ hasText: text }).click();
  }

  // Close the dropdown by pressing Escape
  await page.keyboard.press('Escape');
}

/**
 * Clear selection of a Quasar q-select
 *
 * @param {Page} page - Playwright page instance
 * @param {string} testId - data-testid attribute value of the q-select
 * @returns {Promise<void>}
 */
export async function clearSelect(page: Page, testId: string): Promise<void> {
  const select = page.getByTestId(testId);
  const clearBtn = select.locator('.q-field__append .q-icon').filter({ hasText: 'cancel' });
  if (await clearBtn.isVisible()) {
    await clearBtn.click();
  }
}
