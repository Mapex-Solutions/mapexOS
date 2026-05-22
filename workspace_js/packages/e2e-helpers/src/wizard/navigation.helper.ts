import type { Page } from '@playwright/test';
import { clickNext } from './stepper.helper';

/**
 * Complete the current step and advance to the next one
 * Clicks "Next" and waits for the step transition
 *
 * @param {Page} page - Playwright page instance
 * @param {number} expectedStep - The step number expected after clicking next
 * @returns {Promise<void>}
 */
export async function completeStep(page: Page, expectedStep: number): Promise<void> {
  await clickNext(page);
  // Wait for step indicator to update
  await page.waitForTimeout(300);
}

/**
 * Navigate to a specific step by clicking on the stepper
 * Only works when step navigation is allowed (edit mode)
 *
 * @param {Page} page - Playwright page instance
 * @param {number} stepNumber - 1-based step number to navigate to
 * @param {string} stepIdPrefix - Prefix used for step IDs (default: 'step')
 * @returns {Promise<void>}
 */
export async function navigateToStep(
  page: Page,
  stepNumber: number,
  stepIdPrefix: string = 'step',
): Promise<void> {
  const stepElement = page.locator(`#${stepIdPrefix}-${stepNumber}`);
  await stepElement.click();
  await page.waitForTimeout(300);
}

/**
 * Complete multiple wizard steps sequentially
 * Useful for quickly advancing through a wizard
 *
 * @param {Page} page - Playwright page instance
 * @param {number} fromStep - Starting step (1-based)
 * @param {number} toStep - Target step (1-based)
 * @returns {Promise<void>}
 */
export async function advanceToStep(
  page: Page,
  fromStep: number,
  toStep: number,
): Promise<void> {
  for (let step = fromStep; step < toStep; step++) {
    await completeStep(page, step + 1);
  }
}
