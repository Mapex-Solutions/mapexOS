import type { Page } from '@playwright/test';
import { expect } from '@playwright/test';
import { fillInput, getNativeInput, clickNext, clickPrevious, clickSave } from '@mapexos/e2e-helpers';
import type { UserRegistrationData, UserPersonalData, UserSecurityData } from './user-registration.data';

/**
 * Page Object for the User Registration Wizard
 * Encapsulates all interactions with the /users/add page
 */
export class UserRegistrationPage {
  constructor(private readonly page: Page) {}

  /**
   * Navigate to the user registration page
   *
   * @param {object} options - Navigation options
   * @param {boolean} options.tour - Whether to enable tour mode
   * @returns {Promise<void>}
   */
  async goto(options?: { tour?: boolean }): Promise<void> {
    const url = options?.tour ? '/users/add?tour=true' : '/users/add';
    await this.page.goto(url);
    await this.page.waitForLoadState('networkidle');
  }

  /**
   * Fill Step 1 - Personal Information
   *
   * @param {UserPersonalData} data - Personal info to fill
   * @returns {Promise<void>}
   */
  async fillStep1(data: UserPersonalData): Promise<void> {
    await fillInput(this.page, 'user-firstname-input', data.firstName);
    await fillInput(this.page, 'user-lastname-input', data.lastName);
    await fillInput(this.page, 'user-email-input', data.email);

    if (data.phone) {
      await fillInput(this.page, 'user-phone-input', data.phone);
    }
    if (data.jobTitle) {
      await fillInput(this.page, 'user-jobtitle-input', data.jobTitle);
    }
  }

  /**
   * Fill Step 2 - Security Settings
   *
   * @param {UserSecurityData} data - Security info to fill
   * @returns {Promise<void>}
   */
  async fillStep2(data: UserSecurityData): Promise<void> {
    await fillInput(this.page, 'user-password-input', data.password);
    await fillInput(this.page, 'user-confirm-password-input', data.confirmPassword);

    if (data.changePasswordNextLogin) {
      await this.page.getByTestId('user-change-pwd-checkbox').click();
    }
  }

  /**
   * Fill Step 3 - Access Configuration (group mode)
   * Selects group access type + opens drawer and picks first group
   *
   * @returns {Promise<void>}
   */
  async fillStep3WithGroup(): Promise<void> {
    // Group access type is default, ensure it's selected
    await this.page.getByTestId('user-access-type-group').click();

    // Click the group selector to open the drawer
    await this.page.getByTestId('user-group-select-btn').click();

    // Wait for the drawer to open and groups to load
    const drawer = this.page.locator('.q-dialog');
    await expect(drawer).toBeVisible();

    // Wait for loading spinner to disappear
    await this.page.locator('.q-dialog .q-spinner').waitFor({ state: 'hidden', timeout: 10000 }).catch(() => {});

    // Click the first group in the list
    const firstGroup = this.page.locator('.q-dialog .q-item').first();
    await expect(firstGroup).toBeVisible({ timeout: 10000 });
    await firstGroup.click();

    // Wait for drawer to close
    await expect(drawer).toBeHidden({ timeout: 5000 });
  }

  /**
   * Verify the review step displays the expected data
   *
   * @param {UserRegistrationData} data - Expected data to verify
   * @returns {Promise<void>}
   */
  async verifyReview(data: UserRegistrationData): Promise<void> {
    const personalSection = this.page.getByTestId('review-personal-section');
    await expect(personalSection).toBeVisible();
    await expect(personalSection).toContainText(data.personal.firstName);
    await expect(personalSection).toContainText(data.personal.lastName);
    await expect(personalSection).toContainText(data.personal.email);

    const securitySection = this.page.getByTestId('review-security-section');
    await expect(securitySection).toBeVisible();

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
   * @param {'personal' | 'security' | 'access'} section - Section to edit
   * @returns {Promise<void>}
   */
  async clickEditOnReview(section: 'personal' | 'security' | 'access'): Promise<void> {
    await this.page.getByTestId(`review-edit-${section}-btn`).click();
  }
}
