import { test, expect } from '../../fixtures/auth.fixture';
import { UserRegistrationPage } from './user-registration.po';
import { createUserData } from './user-registration.data';
import { expectToast, fillInput, getNativeInput } from '@mapexos/e2e-helpers';

test.describe('administration.users.add', () => {
  let userPage: UserRegistrationPage;

  test.beforeEach(async ({ authenticatedPage }) => {
    userPage = new UserRegistrationPage(authenticatedPage);
  });

  test('happy-path: complete full registration flow', async ({ authenticatedPage }) => {
    const data = createUserData();
    await userPage.goto();

    // Step 1: Fill personal info and advance
    await userPage.fillStep1(data.personal);
    await userPage.next();

    // Step 2: Fill security info and advance
    await userPage.fillStep2(data.security);
    await userPage.next();

    // Step 3: Select group access and advance
    await userPage.fillStep3WithGroup();
    await userPage.next();

    // Step 4: Verify review and save
    await userPage.verifyReview(data);
    await userPage.save();

    // Expect success toast and redirect to users list
    await expectToast(authenticatedPage, 'success', { type: 'positive' });
    await authenticatedPage.waitForURL(/\/users$/, { timeout: 10000 });
  });

  test('navigation: Next and Previous buttons work correctly', async ({
    authenticatedPage,
  }) => {
    const data = createUserData();
    await userPage.goto();

    // Start at Step 1 - Previous should not be visible
    await expect(authenticatedPage.getByTestId('wizard-previous-btn')).toBeHidden();
    await expect(authenticatedPage.getByTestId('wizard-next-btn')).toBeVisible();

    // Fill Step 1 and go to Step 2
    await userPage.fillStep1(data.personal);
    await userPage.next();

    // Step 2 - Previous should now be visible
    await expect(authenticatedPage.getByTestId('wizard-previous-btn')).toBeVisible();
    await expect(authenticatedPage.getByTestId('wizard-next-btn')).toBeVisible();

    // Go back to Step 1
    await userPage.previous();

    // Verify we're back on Step 1 - check that firstname input still has value
    const firstNameInput = getNativeInput(authenticatedPage, 'user-firstname-input');
    await expect(firstNameInput).toHaveValue(data.personal.firstName);

    // Go forward again
    await userPage.next();
    // Fill security and advance to Step 3
    await userPage.fillStep2(data.security);
    await userPage.next();

    // Step 3 - should see Previous and Next
    await expect(authenticatedPage.getByTestId('wizard-previous-btn')).toBeVisible();

    // Select a group and advance to Step 4 (Review)
    await userPage.fillStep3WithGroup();
    await userPage.next();

    // Step 4 - Next should be hidden, Save should be visible
    await expect(authenticatedPage.getByTestId('wizard-next-btn')).toBeHidden();
    await expect(authenticatedPage.getByTestId('wizard-save-btn')).toBeVisible();
  });

  test('validation.step1: block navigation when required fields are empty', async ({
    authenticatedPage,
  }) => {
    await userPage.goto();

    // Try to advance without filling anything
    const nextDisabled = await userPage.isNextDisabled();
    expect(nextDisabled).toBe(true);

    // Fill only first name - should still be disabled
    await fillInput(authenticatedPage, 'user-firstname-input', 'Test');
    const stillDisabled = await userPage.isNextDisabled();
    expect(stillDisabled).toBe(true);

    // Fill last name too - still need email
    await fillInput(authenticatedPage, 'user-lastname-input', 'User');
    const stillDisabled2 = await userPage.isNextDisabled();
    expect(stillDisabled2).toBe(true);

    // Fill email - now should be enabled
    await fillInput(authenticatedPage, 'user-email-input', 'test@mapex.test');
    const nowEnabled = await userPage.isNextDisabled();
    expect(nowEnabled).toBe(false);
  });

  test('validation.step2: block navigation when password is empty or too short', async ({
    authenticatedPage,
  }) => {
    const data = createUserData();
    await userPage.goto();

    // Complete Step 1
    await userPage.fillStep1(data.personal);
    await userPage.next();

    // On Step 2 without filling anything - Next should be disabled
    const nextDisabled = await userPage.isNextDisabled();
    expect(nextDisabled).toBe(true);

    // Fill a short password (< 8 chars) - should still be disabled
    await fillInput(authenticatedPage, 'user-password-input', 'Short1');
    const stillDisabled = await userPage.isNextDisabled();
    expect(stillDisabled).toBe(true);

    // Fill a valid password - should be enabled
    await fillInput(authenticatedPage, 'user-password-input', 'TestP@ss123');
    await fillInput(authenticatedPage, 'user-confirm-password-input', 'TestP@ss123');
    const nowEnabled = await userPage.isNextDisabled();
    expect(nowEnabled).toBe(false);
  });

  test('validation.step3: block navigation when no group is selected', async ({
    authenticatedPage,
  }) => {
    const data = createUserData();
    await userPage.goto();

    // Complete Step 1
    await userPage.fillStep1(data.personal);
    await userPage.next();

    // Complete Step 2
    await userPage.fillStep2(data.security);
    await userPage.next();

    // On Step 3 without selecting a group - Next should be disabled
    const nextDisabled = await userPage.isNextDisabled();
    expect(nextDisabled).toBe(true);
  });

  test('review: edit button navigates back to the correct step', async ({
    authenticatedPage,
  }) => {
    const data = createUserData();
    await userPage.goto();

    // Complete all steps to reach review
    await userPage.fillStep1(data.personal);
    await userPage.next();
    await userPage.fillStep2(data.security);
    await userPage.next();
    await userPage.fillStep3WithGroup();
    await userPage.next();

    // Verify we're on review - all sections visible
    await expect(authenticatedPage.getByTestId('review-personal-section')).toBeVisible();
    await expect(authenticatedPage.getByTestId('review-security-section')).toBeVisible();
    await expect(authenticatedPage.getByTestId('review-access-section')).toBeVisible();

    // Click edit on personal section -> should go to Step 1
    await userPage.clickEditOnReview('personal');

    // Verify Step 1 fields are visible and populated
    const firstNameInput = authenticatedPage.getByTestId('user-firstname-input');
    await expect(firstNameInput).toBeVisible();
    await expect(firstNameInput).toHaveValue(data.personal.firstName);
  });

  test('tour: pre-fill form with demo data', async ({ authenticatedPage }) => {
    await userPage.goto({ tour: true });

    // Wait for the Driver.js tour overlay to appear
    const tourPopover = authenticatedPage.locator('.driver-popover');
    await tourPopover.waitFor({ state: 'visible', timeout: 10000 });

    // Close the tour overlay so form elements are accessible
    await authenticatedPage.locator('.driver-popover-close-btn').click();
    await tourPopover.waitFor({ state: 'hidden', timeout: 5000 });

    // Verify Step 1 is pre-filled with demo data
    const firstNameInput = getNativeInput(authenticatedPage, 'user-firstname-input');
    await expect(firstNameInput).toHaveValue('John', { timeout: 10000 });

    const lastNameInput = getNativeInput(authenticatedPage, 'user-lastname-input');
    await expect(lastNameInput).toHaveValue('Doe');

    const emailInput = getNativeInput(authenticatedPage, 'user-email-input');
    await expect(emailInput).toHaveValue('john.doe@example.com');
  });
});
