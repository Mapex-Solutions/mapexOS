import { test, expect } from '../fixtures/auth.fixture';
import { AssetRegistrationPage } from './asset-registration.po';
import { createAssetData } from './asset-registration.data';
import { expectToast, fillInput, getNativeInput } from '@mapexos/e2e-helpers';

test.describe('assets.add', () => {
  let assetPage: AssetRegistrationPage;

  test.beforeEach(async ({ authenticatedPage }) => {
    assetPage = new AssetRegistrationPage(authenticatedPage);
  });

  test('happy-path: complete full asset registration flow', async ({ authenticatedPage }) => {
    const data = createAssetData();
    await assetPage.goto();

    // Step 1: Fill identification and advance
    await assetPage.fillStep1(data.identification);
    await assetPage.next();

    // Step 2: Select asset template via drawer and advance
    await assetPage.fillStep2WithTemplate();
    await assetPage.next();

    // Step 3: Select route groups via drawer and advance
    await assetPage.fillStep3WithRouteGroups();
    await assetPage.next();

    // Step 4: Fill connectivity and advance
    await assetPage.fillStep4(data.connectivity);
    await assetPage.next();

    // Step 5: Verify review and save
    await assetPage.verifyReview(data);
    await assetPage.save();

    // Expect success toast and redirect to assets list
    await expectToast(authenticatedPage, 'asset created successfully', { type: 'positive' });
    await authenticatedPage.waitForURL(/\/assets$/, { timeout: 10000 });
  });

  test('navigation: Next and Previous buttons work correctly', async ({
    authenticatedPage,
  }) => {
    const data = createAssetData();
    await assetPage.goto();

    // Step 1 - Previous should not be visible, Next should be visible
    await expect(authenticatedPage.getByTestId('wizard-previous-btn')).toBeHidden();
    await expect(authenticatedPage.getByTestId('wizard-next-btn')).toBeVisible();

    // Fill Step 1 and go to Step 2
    await assetPage.fillStep1(data.identification);
    await assetPage.next();

    // Step 2 - Previous should now be visible
    await expect(authenticatedPage.getByTestId('wizard-previous-btn')).toBeVisible();
    await expect(authenticatedPage.getByTestId('wizard-next-btn')).toBeVisible();

    // Go back to Step 1 - verify data persists
    await assetPage.previous();
    const nameInput = getNativeInput(authenticatedPage, 'asset-name-input');
    await expect(nameInput).toHaveValue(data.identification.name);

    // Advance through all steps to reach Step 5 (Review)
    await assetPage.next();
    await assetPage.fillStep2WithTemplate();
    await assetPage.next();
    await assetPage.fillStep3WithRouteGroups();
    await assetPage.next();
    await assetPage.fillStep4(data.connectivity);
    await assetPage.next();

    // Step 5 (Review) - Next should be hidden, Save should be visible
    await expect(authenticatedPage.getByTestId('wizard-next-btn')).toBeHidden();
    await expect(authenticatedPage.getByTestId('wizard-save-btn')).toBeVisible();
  });

  test('validation.step1: block navigation when required fields are empty', async ({
    authenticatedPage,
  }) => {
    await assetPage.goto();

    // Step 1 uses click-time validation: button is enabled but clicking
    // triggers QForm.validate() which blocks navigation if fields are invalid.

    // Click Next with empty form - should stay on Step 1 (name input still visible)
    await assetPage.next();
    await expect(authenticatedPage.getByTestId('asset-name-input')).toBeVisible();

    // Fill only name and click Next - should stay on Step 1 (assetId still required)
    await fillInput(authenticatedPage, 'asset-name-input', 'Test Asset');
    await assetPage.next();
    await expect(authenticatedPage.getByTestId('asset-name-input')).toBeVisible();

    // Fill assetId too and click Next - should advance to Step 2 (template input visible)
    await fillInput(authenticatedPage, 'asset-id-input', 'A001');
    await assetPage.next();
    await expect(authenticatedPage.getByTestId('asset-template-input')).toBeVisible();
  });

  test('validation.step2: block navigation when no template is selected', async ({
    authenticatedPage,
  }) => {
    const data = createAssetData();
    await assetPage.goto();

    // Complete Step 1
    await assetPage.fillStep1(data.identification);
    await assetPage.next();

    // Step 2 without selecting template - Next should be disabled
    const nextDisabled = await assetPage.isNextDisabled();
    expect(nextDisabled).toBe(true);
  });

  test('validation.step3: block navigation when no route group is selected', async ({
    authenticatedPage,
  }) => {
    const data = createAssetData();
    await assetPage.goto();

    // Complete Steps 1 and 2
    await assetPage.fillStep1(data.identification);
    await assetPage.next();
    await assetPage.fillStep2WithTemplate();
    await assetPage.next();

    // Step 3 without selecting route group - Next should be disabled
    const nextDisabled = await assetPage.isNextDisabled();
    expect(nextDisabled).toBe(true);
  });

  test('review: all sections are visible and contain expected data', async ({
    authenticatedPage,
  }) => {
    const data = createAssetData();
    await assetPage.goto();

    // Complete all steps
    await assetPage.fillStep1(data.identification);
    await assetPage.next();
    await assetPage.fillStep2WithTemplate();
    await assetPage.next();
    await assetPage.fillStep3WithRouteGroups();
    await assetPage.next();
    await assetPage.fillStep4(data.connectivity);
    await assetPage.next();

    // Verify all 4 review sections are visible with correct data
    const identSection = authenticatedPage.getByTestId('review-identification-section');
    await expect(identSection).toBeVisible();
    await expect(identSection).toContainText(data.identification.name);
    await expect(identSection).toContainText(data.identification.assetId);

    const templateSection = authenticatedPage.getByTestId('review-template-section');
    await expect(templateSection).toBeVisible();

    const routeGroupsSection = authenticatedPage.getByTestId('review-routegroups-section');
    await expect(routeGroupsSection).toBeVisible();

    const connectivitySection = authenticatedPage.getByTestId('review-connectivity-section');
    await expect(connectivitySection).toBeVisible();
    await expect(connectivitySection).toContainText(data.connectivity.protocol);
  });
});
