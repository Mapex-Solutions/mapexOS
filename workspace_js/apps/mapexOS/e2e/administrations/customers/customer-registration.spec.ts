import { test, expect } from '../../fixtures/auth.fixture';
import { CustomerRegistrationPage } from './customer-registration.po';
import { createCustomerData } from './customer-registration.data';
import { expectToast, fillInput, getNativeInput } from '@mapexos/e2e-helpers';

test.describe('administration.customers.add', () => {
  let customerPage: CustomerRegistrationPage;

  test.beforeEach(async ({ authenticatedPage }) => {
    customerPage = new CustomerRegistrationPage(authenticatedPage);
  });

  test('happy-path: complete full customer registration flow', async ({ authenticatedPage }) => {
    const data = createCustomerData();
    await customerPage.goto();

    // Step 1: Fill basic info and advance
    await customerPage.fillStep1(data.basic);
    await customerPage.next();

    // Step 2: Fill address and advance
    await customerPage.fillStep2(data.address);
    await customerPage.next();

    // Step 3: Select access policy and advance
    await customerPage.fillStep3(data.accessPolicy);
    await customerPage.next();

    // Step 4: Verify review and save
    await customerPage.verifyReview(data);
    await customerPage.save();

    // Expect success toast and redirect to customers list
    await expectToast(authenticatedPage, 'customer created successfully', { type: 'positive' });
    await authenticatedPage.waitForURL(/\/customers$/, { timeout: 10000 });
  });

  test('navigation: Next and Previous buttons work correctly', async ({
    authenticatedPage,
  }) => {
    const data = createCustomerData();
    await customerPage.goto();

    // Start at Step 1 - Previous should not be visible
    await expect(authenticatedPage.getByTestId('wizard-previous-btn')).toBeHidden();
    await expect(authenticatedPage.getByTestId('wizard-next-btn')).toBeVisible();

    // Fill Step 1 and go to Step 2
    await customerPage.fillStep1(data.basic);
    await customerPage.next();

    // Step 2 - Previous should now be visible
    await expect(authenticatedPage.getByTestId('wizard-previous-btn')).toBeVisible();
    await expect(authenticatedPage.getByTestId('wizard-next-btn')).toBeVisible();

    // Go back to Step 1
    await customerPage.previous();

    // Verify we're back on Step 1 - check that name input still has value
    const nameInput = getNativeInput(authenticatedPage, 'customer-name-input');
    await expect(nameInput).toHaveValue(data.basic.name);

    // Go forward again through remaining steps
    await customerPage.next();
    await customerPage.fillStep2(data.address);
    await customerPage.next();

    // Step 3 - should see Previous and Next
    await expect(authenticatedPage.getByTestId('wizard-previous-btn')).toBeVisible();
    await expect(authenticatedPage.getByTestId('wizard-next-btn')).toBeVisible();

    // Select access policy and advance to Step 4 (Review)
    await customerPage.fillStep3(data.accessPolicy);
    await customerPage.next();

    // Step 4 - Next should be hidden, Save should be visible
    await expect(authenticatedPage.getByTestId('wizard-next-btn')).toBeHidden();
    await expect(authenticatedPage.getByTestId('wizard-save-btn')).toBeVisible();
  });

  test('validation.step1: block navigation when name is empty or too short', async ({
    authenticatedPage,
  }) => {
    await customerPage.goto();

    // Try to advance without filling anything - Next should be disabled
    const nextDisabled = await customerPage.isNextDisabled();
    expect(nextDisabled).toBe(true);

    // Fill a name that is too short (< 3 chars)
    await fillInput(authenticatedPage, 'customer-name-input', 'AB');
    const stillDisabled = await customerPage.isNextDisabled();
    expect(stillDisabled).toBe(true);

    // Fill a valid name (>= 3 chars) - should be enabled
    await fillInput(authenticatedPage, 'customer-name-input', 'ABC');
    const nowEnabled = await customerPage.isNextDisabled();
    expect(nowEnabled).toBe(false);
  });

  test('validation.step2: allow navigation with empty address fields', async ({
    authenticatedPage,
  }) => {
    const data = createCustomerData();
    await customerPage.goto();

    // Complete Step 1
    await customerPage.fillStep1(data.basic);
    await customerPage.next();

    // On Step 2 without filling anything - Next should still be enabled (all optional)
    const nextDisabled = await customerPage.isNextDisabled();
    expect(nextDisabled).toBe(false);
  });

  test('access-policy: card selection updates correctly', async ({
    authenticatedPage,
  }) => {
    const data = createCustomerData();
    await customerPage.goto();

    // Complete Steps 1 and 2
    await customerPage.fillStep1(data.basic);
    await customerPage.next();
    await customerPage.fillStep2(data.address);
    await customerPage.next();

    // Default should be "strict" selected
    const strictCard = authenticatedPage.getByTestId('customer-role-policy-strict');
    await expect(strictCard).toHaveClass(/option-card--selected/);

    // Click "merge" card
    await authenticatedPage.getByTestId('customer-role-policy-merge').click();
    const mergeCard = authenticatedPage.getByTestId('customer-role-policy-merge');
    await expect(mergeCard).toHaveClass(/option-card--selected/);
    await expect(strictCard).not.toHaveClass(/option-card--selected/);

    // Default scope should be "local" selected
    const localCard = authenticatedPage.getByTestId('customer-scope-local');
    await expect(localCard).toHaveClass(/option-card--selected/);

    // Click "recursive" card
    await authenticatedPage.getByTestId('customer-scope-recursive').click();
    const recursiveCard = authenticatedPage.getByTestId('customer-scope-recursive');
    await expect(recursiveCard).toHaveClass(/option-card--selected/);
    await expect(localCard).not.toHaveClass(/option-card--selected/);
  });

  test('review: all sections are visible and contain expected data', async ({
    authenticatedPage,
  }) => {
    const data = createCustomerData();
    await customerPage.goto();

    // Complete all steps
    await customerPage.fillStep1(data.basic);
    await customerPage.next();
    await customerPage.fillStep2(data.address);
    await customerPage.next();
    await customerPage.fillStep3(data.accessPolicy);
    await customerPage.next();

    // Verify all review sections are visible
    const basicSection = authenticatedPage.getByTestId('review-basic-section');
    await expect(basicSection).toBeVisible();
    await expect(basicSection).toContainText(data.basic.name);

    const addressSection = authenticatedPage.getByTestId('review-address-section');
    await expect(addressSection).toBeVisible();
    await expect(addressSection).toContainText(data.address.country!);
    await expect(addressSection).toContainText(data.address.city!);

    const accessSection = authenticatedPage.getByTestId('review-access-section');
    await expect(accessSection).toBeVisible();
    await expect(accessSection).toContainText('Strict');
    await expect(accessSection).toContainText('Local');
  });
});
