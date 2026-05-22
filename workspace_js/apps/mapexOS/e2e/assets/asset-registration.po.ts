import type { Page } from '@playwright/test';
import { expect } from '@playwright/test';
import {
  fillInput,
  getNativeInput,
  selectOption,
  clickNext,
  clickPrevious,
  clickSave,
} from '@mapexos/e2e-helpers';
import type {
  AssetRegistrationData,
  AssetIdentificationData,
  AssetConnectivityData,
} from './asset-registration.data';

/**
 * Page Object for the Asset Registration Wizard
 * Encapsulates all interactions with the /assets/add page (5-step wizard)
 */
export class AssetRegistrationPage {
  constructor(private readonly page: Page) {}

  /**
   * Navigate to the asset registration page
   *
   * @returns {Promise<void>}
   */
  async goto(): Promise<void> {
    await this.page.goto('/assets/add');
    await this.page.waitForLoadState('networkidle');
  }

  /**
   * Fill Step 1 - Identification
   *
   * @param {AssetIdentificationData} data - Identification data to fill
   * @returns {Promise<void>}
   */
  async fillStep1(data: AssetIdentificationData): Promise<void> {
    await fillInput(this.page, 'asset-name-input', data.name);
    await fillInput(this.page, 'asset-id-input', data.assetId);

    if (data.description) {
      await fillInput(this.page, 'asset-description-input', data.description);
    }

    // Status is "Active" (true) by default; change only if explicitly false
    if (data.enabled === false) {
      await selectOption(this.page, 'asset-status-select', 'Inactive');
    }

    // Debug is off by default; toggle only if explicitly true
    if (data.debugEnabled === true) {
      await this.page.getByTestId('asset-debug-toggle').click();
    }
  }

  /**
   * Fill Step 2 - Select Asset Template via drawer (single-select mode)
   * Opens the template selector drawer and picks the first available template.
   * In single-select mode, clicking an item auto-confirms and closes the drawer.
   *
   * @returns {Promise<void>}
   */
  async fillStep2WithTemplate(): Promise<void> {
    // Click the template input to open the drawer
    await this.page.getByTestId('asset-template-input').click();

    // Wait for the dialog to appear
    const dialog = this.page.locator('.q-dialog');
    await dialog.waitFor({ state: 'visible' });
    await this.page.waitForLoadState('networkidle');

    // Wait for loading spinner to disappear and list items to render
    const list = dialog.locator('.q-list');
    await list.waitFor({ state: 'visible', timeout: 15000 });

    // Click the first template item in the list (single-select auto-closes the drawer)
    const firstItem = list.locator('.q-item').first();
    await firstItem.waitFor({ state: 'visible', timeout: 10000 });
    await firstItem.click();

    // Wait for the dialog to close (single-select auto-confirms)
    await dialog.waitFor({ state: 'hidden', timeout: 5000 });

    // Verify template card is visible (selection was successful)
    await expect(this.page.getByTestId('asset-template-card')).toBeVisible();
  }

  /**
   * Fill Step 3 - Select Route Groups via drawer (multi-select mode)
   * Opens the route group selector drawer, picks the first available route group,
   * and clicks "Confirm Selection".
   *
   * @returns {Promise<void>}
   */
  async fillStep3WithRouteGroups(): Promise<void> {
    // Click the "Select Route Groups" button
    await this.page.getByTestId('asset-routegroup-select-btn').click();

    // Wait for the dialog to appear
    const dialog = this.page.locator('.q-dialog');
    await dialog.waitFor({ state: 'visible' });
    await this.page.waitForLoadState('networkidle');

    // Wait for loading spinner to disappear and list items to render
    const list = dialog.locator('.q-list');
    await list.waitFor({ state: 'visible', timeout: 15000 });

    // Click the first route group item (multi-select: toggles checkbox)
    const firstItem = list.locator('.q-item').first();
    await firstItem.waitFor({ state: 'visible', timeout: 10000 });
    await firstItem.click();

    // Click "Confirm Selection" button (multi-select mode has confirm button)
    await dialog.getByRole('button', { name: 'Confirm Selection' }).click();

    // Wait for the dialog to close
    await dialog.waitFor({ state: 'hidden', timeout: 5000 });

    // Verify route group count card is visible (selection was successful)
    await expect(this.page.getByTestId('asset-routegroup-count')).toBeVisible();
  }

  /**
   * Fill Step 4 - Connectivity
   * Selects protocol and fills conditional fields (MQTT) and optional lat/lng
   *
   * @param {AssetConnectivityData} data - Connectivity data to fill
   * @returns {Promise<void>}
   */
  async fillStep4(data: AssetConnectivityData): Promise<void> {
    // Protocol is HTTP by default; only change if different
    if (data.protocol !== 'HTTP') {
      await selectOption(this.page, 'asset-protocol-select', data.protocol);
    }

    // Fill MQTT fields when protocol is MQTT
    if (data.protocol === 'MQTT') {
      if (data.mqttUsername) {
        await fillInput(this.page, 'asset-mqtt-username-input', data.mqttUsername);
      }
      if (data.mqttClientId) {
        await fillInput(this.page, 'asset-mqtt-clientid-input', data.mqttClientId);
      }
    }

    // Fill optional location fields
    if (data.latitude !== undefined) {
      await fillInput(this.page, 'asset-latitude-input', data.latitude.toString());
    }
    if (data.longitude !== undefined) {
      await fillInput(this.page, 'asset-longitude-input', data.longitude.toString());
    }
  }

  /**
   * Verify the review step displays the expected data
   *
   * @param {AssetRegistrationData} data - Expected data to verify
   * @returns {Promise<void>}
   */
  async verifyReview(data: AssetRegistrationData): Promise<void> {
    const identSection = this.page.getByTestId('review-identification-section');
    await expect(identSection).toBeVisible();
    await expect(identSection).toContainText(data.identification.name);
    await expect(identSection).toContainText(data.identification.assetId);

    const templateSection = this.page.getByTestId('review-template-section');
    await expect(templateSection).toBeVisible();

    const routeGroupsSection = this.page.getByTestId('review-routegroups-section');
    await expect(routeGroupsSection).toBeVisible();

    const connectivitySection = this.page.getByTestId('review-connectivity-section');
    await expect(connectivitySection).toBeVisible();
    await expect(connectivitySection).toContainText(data.connectivity.protocol);
  }

  /**
   * Click Next to advance to the next step
   *
   * @returns {Promise<void>}
   */
  async next(): Promise<void> {
    await clickNext(this.page);
  }

  /**
   * Click Previous to go back to the previous step
   *
   * @returns {Promise<void>}
   */
  async previous(): Promise<void> {
    await clickPrevious(this.page);
  }

  /**
   * Click Save to submit the form
   *
   * @returns {Promise<void>}
   */
  async save(): Promise<void> {
    await clickSave(this.page);
  }

  /**
   * Check if the Next button is disabled
   *
   * @returns {Promise<boolean>}
   */
  async isNextDisabled(): Promise<boolean> {
    return this.page.getByTestId('wizard-next-btn').isDisabled();
  }

  /**
   * Check if the Save button is disabled
   *
   * @returns {Promise<boolean>}
   */
  async isSaveDisabled(): Promise<boolean> {
    return this.page.getByTestId('wizard-save-btn').isDisabled();
  }

  /**
   * Get a native input locator inside a Quasar q-input by testid
   *
   * @param {string} testId - data-testid of the q-input wrapper
   * @returns {Locator}
   */
  getNativeInput(testId: string) {
    return getNativeInput(this.page, testId);
  }

  /**
   * Click edit button on a review section to navigate back
   *
   * @param {'identification' | 'template' | 'routegroups' | 'connectivity'} section - Section to edit
   * @returns {Promise<void>}
   */
  async clickEditOnReview(
    section: 'identification' | 'template' | 'routegroups' | 'connectivity',
  ): Promise<void> {
    await this.page.getByTestId(`review-${section}-section-edit-btn`).click();
  }
}
