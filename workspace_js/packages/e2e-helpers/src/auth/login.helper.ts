import type { Page, BrowserContext } from '@playwright/test';

/**
 * Default authentication credentials
 */
const DEFAULT_CREDENTIALS = {
  email: 'admin@mapex.global',
  password: 'mapex123',
};

/**
 * Login via the UI form
 *
 * @param {Page} page - Playwright page instance
 * @param {string} email - User email
 * @param {string} password - User password
 * @returns {Promise<void>}
 */
export async function login(
  page: Page,
  email: string = DEFAULT_CREDENTIALS.email,
  password: string = DEFAULT_CREDENTIALS.password,
): Promise<void> {
  await page.goto('/');
  await page.getByTestId('login-email-input').fill(email);
  await page.getByTestId('login-password-input').fill(password);
  await page.getByTestId('login-submit-btn').click();
  await page.waitForURL(/\/(home|dashboard)/);
}

/**
 * Login by injecting auth tokens directly into localStorage
 * Faster than UI login - use for tests that don't test login itself
 *
 * @param {BrowserContext} context - Playwright browser context
 * @param {string} baseURL - Application base URL
 * @param {string} email - User email
 * @param {string} password - User password
 * @returns {Promise<Page>} Page with auth tokens set
 */
export async function loginViaStorage(
  context: BrowserContext,
  baseURL: string,
  email: string = DEFAULT_CREDENTIALS.email,
  password: string = DEFAULT_CREDENTIALS.password,
): Promise<Page> {
  // Create a page and perform login through the UI once
  const page = await context.newPage();
  await page.goto(baseURL);

  // Fill login form
  await page.getByTestId('login-email-input').fill(email);
  await page.getByTestId('login-password-input').fill(password);
  await page.getByTestId('login-submit-btn').click();

  // Wait for redirect after successful login
  await page.waitForURL(/\/(home|dashboard)/, { timeout: 15000 });

  return page;
}

/**
 * Login via API and set tokens in storage
 * Most efficient method - avoids UI interaction entirely
 *
 * @param {Page} page - Playwright page instance
 * @param {string} baseURL - Application base URL
 * @param {string} email - User email
 * @param {string} password - User password
 * @returns {Promise<void>}
 */
export async function loginViaApi(
  page: Page,
  baseURL: string,
  email: string = DEFAULT_CREDENTIALS.email,
  password: string = DEFAULT_CREDENTIALS.password,
): Promise<void> {
  // Navigate to app first to set localStorage domain
  await page.goto(baseURL);

  // Call auth API directly
  const response = await page.request.post(`${baseURL}/api/v1/auth/login`, {
    data: { email, password },
  });

  if (response.ok()) {
    const data = await response.json();

    // Inject tokens into localStorage
    await page.evaluate((authData) => {
      localStorage.setItem('auth_tokens', JSON.stringify(authData));
    }, data);
  }
}
