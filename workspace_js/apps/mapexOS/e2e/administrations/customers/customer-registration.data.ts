/**
 * Test data factories for customer registration E2E tests
 */

/**
 * Customer basic info test data
 */
export interface CustomerBasicData {
  name: string;
  phone?: string;
  enabled?: boolean;
}

/**
 * Customer address test data
 */
export interface CustomerAddressData {
  country?: string;
  state?: string;
  city?: string;
  zipCode?: string;
}

/**
 * Customer access policy test data
 */
export interface CustomerAccessPolicyData {
  rolePolicy: 'strict' | 'merge';
  defaultScope: 'local' | 'recursive';
}

/**
 * Full customer registration test data
 */
export interface CustomerRegistrationData {
  basic: CustomerBasicData;
  address: CustomerAddressData;
  accessPolicy: CustomerAccessPolicyData;
}

/**
 * Generate a valid customer registration dataset
 *
 * @param {Partial<CustomerRegistrationData>} overrides - Fields to override
 * @returns {CustomerRegistrationData} Complete test data
 */
export function createCustomerData(overrides?: Partial<CustomerRegistrationData>): CustomerRegistrationData {
  const timestamp = Date.now();
  return {
    basic: {
      name: `Test Customer ${timestamp}`,
      phone: '+5511999999999',
      enabled: true,
      ...overrides?.basic,
    },
    address: {
      country: 'Brazil',
      state: 'São Paulo',
      city: 'Campinas',
      zipCode: '13000-000',
      ...overrides?.address,
    },
    accessPolicy: {
      rolePolicy: 'strict',
      defaultScope: 'local',
      ...overrides?.accessPolicy,
    },
  };
}

/**
 * Generate minimal customer data (only required fields)
 *
 * @returns {CustomerRegistrationData} Minimal valid data
 */
export function createMinimalCustomerData(): CustomerRegistrationData {
  const timestamp = Date.now();
  return {
    basic: {
      name: `Min Customer ${timestamp}`,
    },
    address: {},
    accessPolicy: {
      rolePolicy: 'strict',
      defaultScope: 'local',
    },
  };
}
