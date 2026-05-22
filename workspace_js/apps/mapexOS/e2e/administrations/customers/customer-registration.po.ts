import type { Page } from '@playwright/test';
import { expect } from '@playwright/test';
import { fillInput, getNativeInput, clickNext, clickPrevious, clickSave } from '@mapexos/e2e-helpers';
import type { CustomerRegistrationData, CustomerBasicData, CustomerAddressData, CustomerAccessPolicyData } from './customer-registration.data';

/**
 * Page Object for the Customer Registration Wizard
 * Encapsulates all interactions with the /customers/add page
 */
export class CustomerRegistrationPage {
  constructor(private readonly page: Page) {}

  /**
   * Navigate to the customer registration page via the customers list.
   * Goes to /customers first and clicks the "Add" button so the URL
   * includes the required ?parentId query parameter.
   *
   * @returns {Promise<void>}
   */
  async goto(): Promise<void> {
    await this.page.goto('/customers');
    await this.page.waitForLoadState('networkidle');

    // Click the "Add" button which links to /customers/add?parentId=<orgId>
    const addBtn = this.page.locator('a[href*="/customers/add"]');
    await addBtn.click();
    await this.page.waitForLoadState('networkidle');
  }

  /**
   * Fill Step 1 - Basic Information
   *
   * @param {CustomerBasicData} data - Basic info to fill
   * @returns {Promise<void>}
   */
  async fillStep1(data: CustomerBasicData): Promise<void> {
    await fillInput(this.page, 'customer-name-input', data.name);

    if (data.phone) {
      await fillInput(this.page, 'customer-phone-input', data.phone);
    }

    // Enabled is true by default; toggle only if explicitly false
    if (data.enabled === false) {
      await this.page.getByTestId('customer-enabled-checkbox').click();
    }
  }

  /**
   * Fill Step 2 - Address Information
   *
   * @param {CustomerAddressData} data - Address info to fill
   * @returns {Promise<void>}
   */
  async fillStep2(data: CustomerAddressData): Promise<void> {
    if (data.country) {
      await fillInput(this.page, 'customer-country-input', data.country);
    }
    if (data.state) {
      await fillInput(this.page, 'customer-state-input', data.state);
    }
    if (data.city) {
      await fillInput(this.page, 'customer-city-input', data.city);
    }
    if (data.zipCode) {
      await fillInput(this.page, 'customer-zipcode-input', data.zipCode);
    }
  }

  /**
   * Fill Step 3 - Access Policy (card-based selection)
   *
   * @param {CustomerAccessPolicyData} data - Access policy selections
   * @returns {Promise<void>}
   */
  async fillStep3(data: CustomerAccessPolicyData): Promise<void> {
    await this.page.getByTestId(`customer-role-policy-${data.rolePolicy}`).click();
    await this.page.getByTestId(`customer-scope-${data.defaultScope}`).click();
  }

  /**
   * Verify the review step displays the expected data
   *
   * @param {CustomerRegistrationData} data - Expected data to verify
   * @returns {Promise<void>}
   */
  async verifyReview(data: CustomerRegistrationData): Promise<void> {
    const basicSection = this.page.getByTestId('review-basic-section');
    await expect(basicSection).toBeVisible();
    await expect(basicSection).toContainText(data.basic.name);

    const addressSection = this.page.getByTestId('review-address-section');
    await expect(addressSection).toBeVisible();

    const accessSection = this.page.getByTestId('review-access-section');
    await expect(accessSection).toBeVisible();
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
   * @param {'basic' | 'address' | 'access'} section - Section to edit
   * @returns {Promise<void>}
   */
  async clickEditOnReview(section: 'basic' | 'address' | 'access'): Promise<void> {
    await this.page.getByTestId(`review-${section}-section-edit-btn`).click();
  }
}
