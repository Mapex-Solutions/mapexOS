/**
 * Test data factories for user registration E2E tests
 */

/**
 * User personal info test data
 */
export interface UserPersonalData {
  firstName: string;
  lastName: string;
  email: string;
  phone?: string;
  jobTitle?: string;
}

/**
 * User security test data
 */
export interface UserSecurityData {
  password: string;
  confirmPassword: string;
  changePasswordNextLogin?: boolean;
  enabled?: boolean;
}

/**
 * Full user registration test data
 */
export interface UserRegistrationData {
  personal: UserPersonalData;
  security: UserSecurityData;
  accessType: 'group' | 'direct';
}

/**
 * Generate a valid user registration dataset
 *
 * @param {Partial<UserRegistrationData>} overrides - Fields to override
 * @returns {UserRegistrationData} Complete test data
 */
export function createUserData(overrides?: Partial<UserRegistrationData>): UserRegistrationData {
  const timestamp = Date.now();
  return {
    personal: {
      firstName: 'Test',
      lastName: 'User',
      email: `test.user.${timestamp}@mapex.test`,
      phone: '+5511999999999',
      jobTitle: 'QA Engineer',
      ...overrides?.personal,
    },
    security: {
      password: 'TestP@ss123',
      confirmPassword: 'TestP@ss123',
      changePasswordNextLogin: false,
      enabled: true,
      ...overrides?.security,
    },
    accessType: overrides?.accessType ?? 'group',
  };
}

/**
 * Generate user data with mismatched passwords (for validation tests)
 *
 * @returns {UserRegistrationData} Data with password mismatch
 */
export function createUserDataWithPasswordMismatch(): UserRegistrationData {
  return createUserData({
    security: {
      password: 'TestP@ss123',
      confirmPassword: 'DifferentP@ss456',
    },
  });
}

/**
 * Generate minimal user data (only required fields)
 *
 * @returns {UserRegistrationData} Minimal valid data
 */
export function createMinimalUserData(): UserRegistrationData {
  const timestamp = Date.now();
  return {
    personal: {
      firstName: 'Min',
      lastName: 'User',
      email: `min.user.${timestamp}@mapex.test`,
    },
    security: {
      password: 'MinP@ss123',
      confirmPassword: 'MinP@ss123',
    },
    accessType: 'group',
  };
}
