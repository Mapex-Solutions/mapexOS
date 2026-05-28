<script setup lang="ts">
defineOptions({
  name: 'MapexValidationHelpModal'
});

/** TYPE IMPORTS */
import type {
  MapexValidationHelpProps,
  MapexValidationHelpEmits,
} from './interfaces';

/** VUE IMPORTS */
import { ref, computed } from 'vue';
import { copyToClipboard } from 'quasar';

/** COMPONENTS */
import { AppTooltip } from '@components/tooltips';
import { AppTabs } from '@components/tabs';

/** COMPOSABLES */
import { useCommonPlaceholders } from '@composables/i18n';

/** UTILS */
import { notifySuccess, notifyFail } from '@utils/alert/notify';

/** PROPS & EMITS */
const props = defineProps<MapexValidationHelpProps>();
const emit = defineEmits<MapexValidationHelpEmits>();

/** COMPOSABLES & STORES */
const { placeholders } = useCommonPlaceholders();

// State
const activeTab = ref('overview');
const searchQuery = ref('');

// Computed
const isOpen = computed({
  get: () => props.modelValue,
  set: (value: boolean) => emit('update:modelValue', value),
});

// Methods
const copyCode = async (code: string) => {
  try {
    await copyToClipboard(code);
    notifySuccess({
      message: 'Code copied to clipboard!',
      timeout: 2000,
    });
  } catch {
    notifyFail({
      message: 'Failed to copy code',
      timeout: 2000,
    });
  }
};

const closeModal = () => {
  isOpen.value = false;
};

/** LOCAL IMPORTS */
import type { TabContent } from './interfaces/mapexValidationHelp.interface';

const tabsContent: Record<string, TabContent> = {
  overview: {
    icon: 'info',
    color: 'primary',
    sections: [
      {
        title: 'Introduction',
        description:
          'MapexValidation ($mv) is a built-in validation library available in all validation scripts. It provides a fluent API for validating data structures.',
        code: `// Access $mv in any validation script
const schema = $mv.object({
  name: $mv.string().required(),
  age: $mv.number().min(18).required()
});

const { value, error } = schema.validate(payload);

if (error) {
  throw new Error('Validation failed: ' + error.message);
}`,
        keywords: ['basic', 'usage', 'intro', 'validation'],
      },
      {
        title: 'Available Types',
        description: 'All validation types available in MapexValidation:',
        code: `// String validation
$mv.string()

// Number validation
$mv.number()

// Boolean validation
$mv.boolean()

// Date validation
$mv.date()

// Array validation
$mv.array()

// Object validation
$mv.object()

// Any type (always passes)
$mv.any()`,
        keywords: ['types', 'string', 'number', 'boolean', 'date', 'array', 'object', 'any'],
      },
      {
        title: 'Common Modifiers',
        description: 'Modifiers that can be applied to most validation types:',
        code: `// Make field required
$mv.string().required()

// Make field optional (default behavior)
$mv.string().optional()

// Chain multiple modifiers
$mv.string().min(3).max(50).required()`,
        keywords: ['required', 'optional', 'modifiers'],
      },
      {
        title: 'Validation Result',
        description: 'The validate() method returns an object with value and error:',
        code: `const schema = $mv.string().required();
const { value, error } = schema.validate(payload);

if (error) {
  // error.message contains the validation error
  console.error(error.message);
} else {
  // value contains the validated data
  console.log(value);
}`,
        keywords: ['result', 'error', 'value', 'validate'],
      },
    ],
  },
  string: {
    icon: 'text_fields',
    color: 'blue',
    sections: [
      {
        title: 'Basic String Validation',
        description: 'Create a string validator:',
        code: `const schema = $mv.string();
const { value, error } = schema.validate('Hello World');`,
        keywords: ['string', 'basic'],
      },
      {
        title: 'Required String',
        description: 'Ensure the string is provided and not empty:',
        code: `const schema = $mv.string().required();

// Valid
schema.validate('Hello'); // { value: 'Hello', error: null }

// Invalid
schema.validate(''); // { value: null, error: Error }
schema.validate(null); // { value: null, error: Error }`,
        keywords: ['string', 'required'],
      },
      {
        title: 'Optional String',
        description: 'Allow the string to be optional (default behavior):',
        code: `const schema = $mv.string().optional();

// Both valid
schema.validate('Hello'); // { value: 'Hello', error: null }
schema.validate(null); // { value: null, error: null }`,
        keywords: ['string', 'optional'],
      },
      {
        title: 'Minimum Length',
        description: 'Require a minimum string length:',
        code: `const schema = $mv.string().min(3);

// Valid
schema.validate('Hello'); // { value: 'Hello', error: null }

// Invalid
schema.validate('Hi'); // { value: null, error: Error }`,
        keywords: ['string', 'min', 'length'],
      },
      {
        title: 'Maximum Length',
        description: 'Limit the maximum string length:',
        code: `const schema = $mv.string().max(10);

// Valid
schema.validate('Hello'); // { value: 'Hello', error: null }

// Invalid
schema.validate('Hello World!'); // { value: null, error: Error }`,
        keywords: ['string', 'max', 'length'],
      },
      {
        title: 'Min and Max Length',
        description: 'Combine minimum and maximum length:',
        code: `const schema = $mv.string().min(3).max(20);

// Valid
schema.validate('Hello'); // { value: 'Hello', error: null }

// Invalid
schema.validate('Hi'); // Too short
schema.validate('This is a very long string that exceeds limit'); // Too long`,
        keywords: ['string', 'min', 'max', 'length'],
      },
      {
        title: 'Email Validation',
        description: 'Validate email format:',
        code: `const schema = $mv.string().email();

// Valid
schema.validate('user@example.com'); // { value: 'user@example.com', error: null }
schema.validate('john.doe@company.co.uk'); // Valid

// Invalid
schema.validate('not-an-email'); // { value: null, error: Error }
schema.validate('missing@domain'); // { value: null, error: Error }`,
        keywords: ['string', 'email'],
      },
      {
        title: 'Pattern Validation (Regex)',
        description: 'Validate against a custom regular expression:',
        code: `// Only allow alphanumeric characters
const schema = $mv.string().pattern(/^[a-zA-Z0-9]+$/);

// Valid
schema.validate('Hello123'); // { value: 'Hello123', error: null }

// Invalid
schema.validate('Hello World!'); // Contains space and special char

// Phone number pattern
const phoneSchema = $mv.string().pattern(/^\\+?[1-9]\\d{1,14}$/);
phoneSchema.validate('+1234567890'); // Valid`,
        keywords: ['string', 'pattern', 'regex'],
      },
      {
        title: 'Complex String Validation',
        description: 'Combine multiple validators:',
        code: `// Username: 3-20 chars, alphanumeric with underscore
const usernameSchema = $mv.string()
  .min(3)
  .max(20)
  .pattern(/^[a-zA-Z0-9_]+$/)
  .required();

// Password: 8+ chars, at least one number
const passwordSchema = $mv.string()
  .min(8)
  .pattern(/^(?=.*\\d).+$/)
  .required();`,
        keywords: ['string', 'complex', 'username', 'password'],
      },
    ],
  },
  number: {
    icon: 'tag',
    color: 'green',
    sections: [
      {
        title: 'Basic Number Validation',
        description: 'Create a number validator:',
        code: `const schema = $mv.number();
const { value, error } = schema.validate(42);`,
        keywords: ['number', 'basic'],
      },
      {
        title: 'Required Number',
        description: 'Ensure the number is provided:',
        code: `const schema = $mv.number().required();

// Valid
schema.validate(42); // { value: 42, error: null }
schema.validate(0); // { value: 0, error: null }

// Invalid
schema.validate(null); // { value: null, error: Error }
schema.validate(undefined); // { value: null, error: Error }`,
        keywords: ['number', 'required'],
      },
      {
        title: 'Optional Number',
        description: 'Allow the number to be optional:',
        code: `const schema = $mv.number().optional();

// Both valid
schema.validate(42); // { value: 42, error: null }
schema.validate(null); // { value: null, error: null }`,
        keywords: ['number', 'optional'],
      },
      {
        title: 'Minimum Value',
        description: 'Require a minimum numeric value:',
        code: `const schema = $mv.number().min(0);

// Valid
schema.validate(10); // { value: 10, error: null }
schema.validate(0); // { value: 0, error: null }

// Invalid
schema.validate(-5); // { value: null, error: Error }`,
        keywords: ['number', 'min'],
      },
      {
        title: 'Maximum Value',
        description: 'Limit the maximum numeric value:',
        code: `const schema = $mv.number().max(100);

// Valid
schema.validate(50); // { value: 50, error: null }
schema.validate(100); // { value: 100, error: null }

// Invalid
schema.validate(150); // { value: null, error: Error }`,
        keywords: ['number', 'max'],
      },
      {
        title: 'Range Validation',
        description: 'Combine minimum and maximum values:',
        code: `const schema = $mv.number().min(1).max(100);

// Valid
schema.validate(50); // { value: 50, error: null }

// Invalid
schema.validate(0); // Too small
schema.validate(101); // Too large`,
        keywords: ['number', 'min', 'max', 'range'],
      },
      {
        title: 'Integer Validation',
        description: 'Ensure the number is an integer (no decimals):',
        code: `const schema = $mv.number().integer();

// Valid
schema.validate(42); // { value: 42, error: null }
schema.validate(-10); // { value: -10, error: null }

// Invalid
schema.validate(3.14); // { value: null, error: Error }
schema.validate(1.5); // { value: null, error: Error }`,
        keywords: ['number', 'integer'],
      },
      {
        title: 'Complex Number Validation',
        description: 'Combine multiple validators:',
        code: `// Age validation: integer between 0 and 120
const ageSchema = $mv.number()
  .integer()
  .min(0)
  .max(120)
  .required();

// Price validation: positive number, max 2 decimals
const priceSchema = $mv.number()
  .min(0)
  .required();

// Percentage: 0-100
const percentageSchema = $mv.number()
  .min(0)
  .max(100)
  .required();`,
        keywords: ['number', 'complex', 'age', 'price', 'percentage'],
      },
    ],
  },
  boolean: {
    icon: 'toggle_on',
    color: 'purple',
    sections: [
      {
        title: 'Basic Boolean Validation',
        description: 'Create a boolean validator:',
        code: `const schema = $mv.boolean();
const { value, error } = schema.validate(true);`,
        keywords: ['boolean', 'basic'],
      },
      {
        title: 'Required Boolean',
        description: 'Ensure the boolean is provided:',
        code: `const schema = $mv.boolean().required();

// Valid
schema.validate(true); // { value: true, error: null }
schema.validate(false); // { value: false, error: null }

// Invalid
schema.validate(null); // { value: null, error: Error }
schema.validate(undefined); // { value: null, error: Error }`,
        keywords: ['boolean', 'required'],
      },
      {
        title: 'Optional Boolean',
        description: 'Allow the boolean to be optional:',
        code: `const schema = $mv.boolean().optional();

// All valid
schema.validate(true); // { value: true, error: null }
schema.validate(false); // { value: false, error: null }
schema.validate(null); // { value: null, error: null }`,
        keywords: ['boolean', 'optional'],
      },
      {
        title: 'Boolean in Forms',
        description: 'Common use cases for boolean validation:',
        code: `// Terms and conditions acceptance
const termsSchema = $mv.boolean().required();

// Optional newsletter subscription
const newsletterSchema = $mv.boolean().optional();

// Feature flags
const schema = $mv.object({
  darkMode: $mv.boolean().optional(),
  notifications: $mv.boolean().optional(),
  analytics: $mv.boolean().required(),
});`,
        keywords: ['boolean', 'forms', 'checkbox'],
      },
    ],
  },
  date: {
    icon: 'calendar_today',
    color: 'orange',
    sections: [
      {
        title: 'Basic Date Validation',
        description: 'Create a date validator:',
        code: `const schema = $mv.date();
const { value, error } = schema.validate(new Date());`,
        keywords: ['date', 'basic'],
      },
      {
        title: 'Required Date',
        description: 'Ensure the date is provided:',
        code: `const schema = $mv.date().required();

// Valid
schema.validate(new Date()); // { value: Date, error: null }
schema.validate('2024-01-01'); // { value: Date, error: null }

// Invalid
schema.validate(null); // { value: null, error: Error }`,
        keywords: ['date', 'required'],
      },
      {
        title: 'Optional Date',
        description: 'Allow the date to be optional:',
        code: `const schema = $mv.date().optional();

// Both valid
schema.validate(new Date()); // { value: Date, error: null }
schema.validate(null); // { value: null, error: null }`,
        keywords: ['date', 'optional'],
      },
      {
        title: 'Minimum Date',
        description: 'Require a date to be after a specific date:',
        code: `const minDate = new Date('2024-01-01');
const schema = $mv.date().min(minDate);

// Valid
schema.validate(new Date('2024-06-01')); // { value: Date, error: null }

// Invalid
schema.validate(new Date('2023-12-31')); // { value: null, error: Error }`,
        keywords: ['date', 'min'],
      },
      {
        title: 'Maximum Date',
        description: 'Require a date to be before a specific date:',
        code: `const maxDate = new Date('2024-12-31');
const schema = $mv.date().max(maxDate);

// Valid
schema.validate(new Date('2024-06-01')); // { value: Date, error: null }

// Invalid
schema.validate(new Date('2025-01-01')); // { value: null, error: Error }`,
        keywords: ['date', 'max'],
      },
      {
        title: 'Date Range Validation',
        description: 'Validate a date is within a specific range:',
        code: `const minDate = new Date('2024-01-01');
const maxDate = new Date('2024-12-31');
const schema = $mv.date().min(minDate).max(maxDate);

// Valid
schema.validate(new Date('2024-06-15')); // Within range

// Invalid
schema.validate(new Date('2023-12-31')); // Before range
schema.validate(new Date('2025-01-01')); // After range`,
        keywords: ['date', 'min', 'max', 'range'],
      },
      {
        title: 'Future Date Validation',
        description: 'Ensure a date is in the future:',
        code: `const now = new Date();
const schema = $mv.date().min(now);

// Valid if date is in the future
const futureDate = new Date();
futureDate.setDate(futureDate.getDate() + 7);
schema.validate(futureDate); // { value: Date, error: null }

// Invalid if date is in the past
schema.validate(new Date('2020-01-01')); // { value: null, error: Error }`,
        keywords: ['date', 'future'],
      },
      {
        title: 'Past Date Validation',
        description: 'Ensure a date is in the past:',
        code: `const now = new Date();
const schema = $mv.date().max(now);

// Valid if date is in the past
schema.validate(new Date('2020-01-01')); // { value: Date, error: null }

// Invalid if date is in the future
const futureDate = new Date();
futureDate.setDate(futureDate.getDate() + 7);
schema.validate(futureDate); // { value: null, error: Error }`,
        keywords: ['date', 'past'],
      },
    ],
  },
  array: {
    icon: 'list',
    color: 'teal',
    sections: [
      {
        title: 'Basic Array Validation',
        description: 'Create an array validator:',
        code: `const schema = $mv.array();
const { value, error } = schema.validate([1, 2, 3]);`,
        keywords: ['array', 'basic'],
      },
      {
        title: 'Required Array',
        description: 'Ensure the array is provided:',
        code: `const schema = $mv.array().required();

// Valid
schema.validate([]); // { value: [], error: null }
schema.validate([1, 2, 3]); // { value: [1, 2, 3], error: null }

// Invalid
schema.validate(null); // { value: null, error: Error }`,
        keywords: ['array', 'required'],
      },
      {
        title: 'Optional Array',
        description: 'Allow the array to be optional:',
        code: `const schema = $mv.array().optional();

// Both valid
schema.validate([1, 2, 3]); // { value: [1, 2, 3], error: null }
schema.validate(null); // { value: null, error: null }`,
        keywords: ['array', 'optional'],
      },
      {
        title: 'Array Item Validation',
        description: 'Validate each item in the array:',
        code: `// Array of strings
const stringArraySchema = $mv.array().items($mv.string());

// Valid
stringArraySchema.validate(['a', 'b', 'c']); // Valid

// Invalid
stringArraySchema.validate([1, 2, 3]); // Items are not strings

// Array of numbers
const numberArraySchema = $mv.array().items($mv.number().min(0));
numberArraySchema.validate([1, 2, 3]); // Valid`,
        keywords: ['array', 'items'],
      },
      {
        title: 'Minimum Array Length',
        description: 'Require a minimum number of items:',
        code: `const schema = $mv.array().min(2);

// Valid
schema.validate([1, 2]); // { value: [1, 2], error: null }
schema.validate([1, 2, 3]); // { value: [1, 2, 3], error: null }

// Invalid
schema.validate([]); // { value: null, error: Error }
schema.validate([1]); // { value: null, error: Error }`,
        keywords: ['array', 'min', 'length'],
      },
      {
        title: 'Maximum Array Length',
        description: 'Limit the maximum number of items:',
        code: `const schema = $mv.array().max(5);

// Valid
schema.validate([1, 2, 3]); // { value: [1, 2, 3], error: null }

// Invalid
schema.validate([1, 2, 3, 4, 5, 6]); // { value: null, error: Error }`,
        keywords: ['array', 'max', 'length'],
      },
      {
        title: 'Array Length Range',
        description: 'Combine minimum and maximum length:',
        code: `const schema = $mv.array().min(1).max(10);

// Valid
schema.validate([1, 2, 3]); // Within range

// Invalid
schema.validate([]); // Too few items
schema.validate([1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11]); // Too many items`,
        keywords: ['array', 'min', 'max', 'length', 'range'],
      },
      {
        title: 'Array of Objects',
        description: 'Validate an array of objects with schemas:',
        code: `const itemSchema = $mv.object({
  id: $mv.number().required(),
  name: $mv.string().required(),
});

const schema = $mv.array().items(itemSchema);

// Valid
schema.validate([
  { id: 1, name: 'Item 1' },
  { id: 2, name: 'Item 2' },
]);

// Invalid
schema.validate([
  { id: 1, name: 'Item 1' },
  { id: 2 }, // Missing name
]);`,
        keywords: ['array', 'objects', 'items'],
      },
      {
        title: 'Complex Array Validation',
        description: 'Combine multiple validators:',
        code: `// Array of 1-10 email addresses
const emailListSchema = $mv.array()
  .items($mv.string().email())
  .min(1)
  .max(10)
  .required();

// Array of positive numbers
const positiveNumbersSchema = $mv.array()
  .items($mv.number().min(0))
  .optional();`,
        keywords: ['array', 'complex'],
      },
    ],
  },
  object: {
    icon: 'data_object',
    color: 'indigo',
    sections: [
      {
        title: 'Basic Object Validation',
        description: 'Create an object validator with a schema:',
        code: `const schema = $mv.object({
  name: $mv.string().required(),
  age: $mv.number().required(),
});

const { value, error } = schema.validate({
  name: 'John',
  age: 30,
});`,
        keywords: ['object', 'basic'],
      },
      {
        title: 'Required Object',
        description: 'Ensure the object is provided:',
        code: `const schema = $mv.object({
  name: $mv.string().required(),
}).required();

// Valid
schema.validate({ name: 'John' }); // { value: { name: 'John' }, error: null }

// Invalid
schema.validate(null); // { value: null, error: Error }`,
        keywords: ['object', 'required'],
      },
      {
        title: 'Optional Object',
        description: 'Allow the object to be optional:',
        code: `const schema = $mv.object({
  name: $mv.string().required(),
}).optional();

// Both valid
schema.validate({ name: 'John' }); // Valid
schema.validate(null); // Valid`,
        keywords: ['object', 'optional'],
      },
      {
        title: 'Nested Objects',
        description: 'Validate nested object structures:',
        code: `const schema = $mv.object({
  user: $mv.object({
    name: $mv.string().required(),
    email: $mv.string().email().required(),
  }).required(),
  settings: $mv.object({
    theme: $mv.string().optional(),
    notifications: $mv.boolean().optional(),
  }).optional(),
});

// Valid
schema.validate({
  user: {
    name: 'John',
    email: 'john@example.com',
  },
  settings: {
    theme: 'dark',
    notifications: true,
  },
});`,
        keywords: ['object', 'nested'],
      },
      {
        title: 'Mixed Field Types',
        description: 'Objects with various field types:',
        code: `const schema = $mv.object({
  name: $mv.string().required(),
  age: $mv.number().min(0).required(),
  email: $mv.string().email().required(),
  isActive: $mv.boolean().required(),
  joinDate: $mv.date().required(),
  tags: $mv.array().items($mv.string()).optional(),
});`,
        keywords: ['object', 'mixed', 'types'],
      },
      {
        title: 'Updating Object Schema',
        description: 'Use keys() to add or update fields:',
        code: `const baseSchema = $mv.object({
  name: $mv.string().required(),
});

// Extend with additional fields
const extendedSchema = baseSchema.keys({
  age: $mv.number().required(),
  email: $mv.string().email().required(),
});`,
        keywords: ['object', 'keys', 'extend'],
      },
      {
        title: 'Optional vs Required Fields',
        description: 'Mix optional and required fields:',
        code: `const schema = $mv.object({
  // Required fields
  id: $mv.number().required(),
  name: $mv.string().required(),

  // Optional fields
  nickname: $mv.string().optional(),
  age: $mv.number().optional(),
  bio: $mv.string().optional(),
});

// Valid - only required fields
schema.validate({
  id: 1,
  name: 'John',
});

// Valid - with optional fields
schema.validate({
  id: 1,
  name: 'John',
  nickname: 'Johnny',
  age: 30,
});`,
        keywords: ['object', 'optional', 'required'],
      },
      {
        title: 'Complex Object Validation',
        description: 'Real-world complex object structure:',
        code: `const userSchema = $mv.object({
  id: $mv.number().required(),
  username: $mv.string().min(3).max(20).required(),
  email: $mv.string().email().required(),
  profile: $mv.object({
    firstName: $mv.string().required(),
    lastName: $mv.string().required(),
    age: $mv.number().min(0).max(120).optional(),
    avatar: $mv.string().optional(),
  }).required(),
  roles: $mv.array().items($mv.string()).min(1).required(),
  settings: $mv.object({
    theme: $mv.string().optional(),
    language: $mv.string().optional(),
    notifications: $mv.boolean().optional(),
  }).optional(),
  createdAt: $mv.date().required(),
  isActive: $mv.boolean().required(),
});`,
        keywords: ['object', 'complex', 'nested'],
      },
    ],
  },
  examples: {
    icon: 'code',
    color: 'deep-orange',
    sections: [
      {
        title: 'Example 1: Sensor Data Validation',
        description: 'Validate IoT sensor data payload:',
        code: `// Sensor data validation schema
const sensorSchema = $mv.object({
  deviceId: $mv.string().required(),
  timestamp: $mv.date().required(),
  readings: $mv.object({
    temperature: $mv.number().min(-50).max(150).required(),
    humidity: $mv.number().min(0).max(100).required(),
    pressure: $mv.number().min(0).optional(),
  }).required(),
  location: $mv.object({
    latitude: $mv.number().min(-90).max(90).required(),
    longitude: $mv.number().min(-180).max(180).required(),
  }).optional(),
  batteryLevel: $mv.number().min(0).max(100).required(),
  status: $mv.string().required(),
});

// Usage in validation script
const { value, error } = sensorSchema.validate(payload);

if (error) {
  throw new Error(\`Sensor data validation failed: \${error.message}\`);
}

// Data is valid, continue processing
console.log('Valid sensor data:', value);`,
        keywords: ['example', 'sensor', 'iot'],
      },
      {
        title: 'Example 2: User Registration',
        description: 'Validate user registration form:',
        code: `// User registration validation
const registrationSchema = $mv.object({
  username: $mv.string()
    .min(3)
    .max(20)
    .pattern(/^[a-zA-Z0-9_]+$/)
    .required(),
  email: $mv.string().email().required(),
  password: $mv.string()
    .min(8)
    .pattern(/^(?=.*[A-Za-z])(?=.*\\d).+$/)
    .required(),
  confirmPassword: $mv.string().required(),
  age: $mv.number().min(18).integer().required(),
  acceptTerms: $mv.boolean().required(),
  newsletter: $mv.boolean().optional(),
});

// Usage
const { value, error } = registrationSchema.validate(payload);

if (error) {
  throw new Error(\`Registration validation failed: \${error.message}\`);
}

// Additional validation: password match
if (value.password !== value.confirmPassword) {
  throw new Error('Passwords do not match');
}

console.log('Registration data is valid:', value);`,
        keywords: ['example', 'user', 'registration', 'form'],
      },
      {
        title: 'Example 3: API Request Validation',
        description: 'Validate complex API request payload:',
        code: `// API request validation for creating an order
const orderSchema = $mv.object({
  customerId: $mv.number().required(),
  items: $mv.array()
    .items($mv.object({
      productId: $mv.number().required(),
      quantity: $mv.number().min(1).integer().required(),
      price: $mv.number().min(0).required(),
      discount: $mv.number().min(0).max(100).optional(),
    }))
    .min(1)
    .required(),
  shippingAddress: $mv.object({
    street: $mv.string().required(),
    city: $mv.string().required(),
    state: $mv.string().required(),
    zipCode: $mv.string().pattern(/^\\d{5}(-\\d{4})?$/).required(),
    country: $mv.string().required(),
  }).required(),
  billingAddress: $mv.object({
    street: $mv.string().required(),
    city: $mv.string().required(),
    state: $mv.string().required(),
    zipCode: $mv.string().pattern(/^\\d{5}(-\\d{4})?$/).required(),
    country: $mv.string().required(),
  }).optional(),
  notes: $mv.string().max(500).optional(),
  preferredDeliveryDate: $mv.date().min(new Date()).optional(),
});

// Usage
const { value, error } = orderSchema.validate(payload);

if (error) {
  throw new Error(\`Order validation failed: \${error.message}\`);
}

console.log('Order is valid:', value);`,
        keywords: ['example', 'api', 'order', 'complex'],
      },
      {
        title: 'Example 4: Configuration Validation',
        description: 'Validate application configuration:',
        code: `// Application configuration validation
const configSchema = $mv.object({
  app: $mv.object({
    name: $mv.string().required(),
    version: $mv.string().pattern(/^\\d+\\.\\d+\\.\\d+$/).required(),
    environment: $mv.string().required(),
    debug: $mv.boolean().required(),
  }).required(),
  database: $mv.object({
    host: $mv.string().required(),
    port: $mv.number().min(1).max(65535).integer().required(),
    username: $mv.string().required(),
    password: $mv.string().min(8).required(),
    database: $mv.string().required(),
    ssl: $mv.boolean().optional(),
  }).required(),
  api: $mv.object({
    baseUrl: $mv.string().required(),
    timeout: $mv.number().min(0).required(),
    retries: $mv.number().min(0).integer().required(),
    apiKey: $mv.string().optional(),
  }).required(),
  features: $mv.object({
    enableAnalytics: $mv.boolean().optional(),
    enableNotifications: $mv.boolean().optional(),
    maintenanceMode: $mv.boolean().optional(),
  }).optional(),
});

// Usage
const { value, error } = configSchema.validate(payload);

if (error) {
  throw new Error(\`Configuration invalid: \${error.message}\`);
}

console.log('Configuration loaded successfully');`,
        keywords: ['example', 'configuration', 'config'],
      },
      {
        title: 'Example 5: Data Import Validation',
        description: 'Validate bulk data import:',
        code: `// Bulk import validation
const importSchema = $mv.object({
  metadata: $mv.object({
    importId: $mv.string().required(),
    timestamp: $mv.date().required(),
    source: $mv.string().required(),
  }).required(),
  records: $mv.array()
    .items($mv.object({
      id: $mv.string().required(),
      name: $mv.string().min(1).max(200).required(),
      email: $mv.string().email().optional(),
      phone: $mv.string().pattern(/^\\+?[1-9]\\d{1,14}$/).optional(),
      category: $mv.string().required(),
      status: $mv.string().required(),
      tags: $mv.array().items($mv.string()).optional(),
      metadata: $mv.object({}).optional(),
    }))
    .min(1)
    .max(1000)
    .required(),
});

// Usage
const { value, error } = importSchema.validate(payload);

if (error) {
  throw new Error(\`Import validation failed: \${error.message}\`);
}

console.log(\`Importing \${value.records.length} records...\`);

// Process each record
value.records.forEach((record, index) => {
  console.log(\`Processing record \${index + 1}: \${record.name}\`);
});`,
        keywords: ['example', 'import', 'bulk', 'array'],
      },
      {
        title: 'Example 6: Webhook Payload Validation',
        description: 'Validate incoming webhook data:',
        code: `// Webhook validation schema
const webhookSchema = $mv.object({
  event: $mv.string().required(),
  timestamp: $mv.date().required(),
  signature: $mv.string().required(),
  data: $mv.object({
    id: $mv.string().required(),
    type: $mv.string().required(),
    attributes: $mv.object({}).required(),
    relationships: $mv.object({}).optional(),
  }).required(),
  metadata: $mv.object({
    webhookId: $mv.string().required(),
    attemptNumber: $mv.number().min(1).integer().required(),
  }).optional(),
});

// Usage in webhook handler
const { value, error } = webhookSchema.validate(payload);

if (error) {
  // Log error and return 400
  console.error('Invalid webhook payload:', error.message);
  throw new Error('Invalid webhook payload');
}

// Verify signature (pseudo-code)
// if (!verifySignature(value.signature, value)) {
//   throw new Error('Invalid webhook signature');
// }

console.log(\`Processing webhook event: \${value.event}\`);`,
        keywords: ['example', 'webhook', 'api'],
      },
    ],
  },
};

// Filtered content based on search
const filteredContent = computed(() => {
  const query = searchQuery.value.toLowerCase().trim();
  if (!query) {
    return tabsContent;
  }

  const filtered: Record<string, TabContent> = {};

  Object.entries(tabsContent).forEach(([key, content]) => {
    const matchingSections = content.sections.filter((section) => {
      return (
        section.title.toLowerCase().includes(query) ||
        section.description.toLowerCase().includes(query) ||
        section.code.toLowerCase().includes(query) ||
        section.keywords.some((keyword) => keyword.includes(query))
      );
    });

    if (matchingSections.length > 0) {
      filtered[key] = {
        ...content,
        sections: matchingSections,
      };
    }
  });

  return filtered;
});

// Check if current tab has content after filtering
const hasCurrentTabContent = computed(() => {
  return !!filteredContent.value[activeTab.value];
});

/**
 * Transform filtered content into AppTabs format
 */
const helpTabs = computed(() => {
  return Object.entries(filteredContent.value).map(([key, content]) => ({
    name: key,
    label: key.charAt(0).toUpperCase() + key.slice(1),
    icon: content.icon,
  }));
});

// Auto-switch to first available tab if current tab is empty
const switchToFirstAvailableTab = () => {
  if (!hasCurrentTabContent.value) {
    const firstAvailableTab = Object.keys(filteredContent.value)[0];
    if (firstAvailableTab) {
      activeTab.value = firstAvailableTab;
    }
  }
};

// Watch search query and switch tabs if needed
const performSearch = () => {
  switchToFirstAvailableTab();
};
</script>

<template>
  <q-dialog v-model="isOpen" maximized transition-show="slide-up" transition-hide="slide-down">
    <q-card class="mapex-validation-help-modal">
      <!-- Header -->
      <q-card-section class="row items-center q-pb-none bg-primary text-white">
        <div class="text-h5 q-mr-md">
          <q-icon name="help" size="32px" class="q-mr-sm" />
          MapexValidation ($mv) Guide
        </div>
        <q-space />
        <q-btn icon="close" flat round dense v-close-popup />
      </q-card-section>

      <!-- Search Bar -->
      <q-card-section class="q-pt-md q-pb-sm">
        <q-input
          v-model="searchQuery"
          outlined
          dense
          :placeholder="placeholders.search.value"
          clearable
          @update:model-value="performSearch"
        >
          <template #prepend>
            <q-icon name="search" />
          </template>
        </q-input>
      </q-card-section>

      <!-- Tabs -->
      <AppTabs v-model="activeTab" :tabs="helpTabs" />

      <!-- Tab Panels -->
      <q-tab-panels v-model="activeTab" animated style="height: calc(100vh - 220px)">
        <q-tab-panel
          v-for="(content, key) in filteredContent"
          :key="key"
          :name="key"
          class="q-pa-md"
        >
          <q-scroll-area style="height: calc(100vh - 240px)">
            <div class="q-gutter-md">
              <!-- Section Cards -->
              <q-card
                v-for="(section, index) in content.sections"
                :key="index"
                bordered
                flat
                class="section-card"
              >
                <q-card-section>
                  <div class="row items-center q-mb-md">
                    <div class="text-h6 text-weight-medium">{{ section.title }}</div>
                    <q-space />
                    <q-btn
                      icon="content_copy"
                      flat
                      round
                      dense
                      color="primary"
                      @click="copyCode(section.code)"
                    >
                      <AppTooltip content="Copy code to clipboard" />
                    </q-btn>
                  </div>
                  <div class="text-body2 text-grey-8 q-mb-md">
                    {{ section.description }}
                  </div>
                  <q-card bordered flat class="bg-grey-10 text-white code-block">
                    <q-card-section class="q-pa-md">
                      <pre class="q-ma-none"><code>{{ section.code }}</code></pre>
                    </q-card-section>
                  </q-card>
                </q-card-section>
              </q-card>
            </div>
          </q-scroll-area>
        </q-tab-panel>
      </q-tab-panels>

      <!-- Footer -->
      <q-separator />
      <q-card-section class="row items-center q-py-sm bg-grey-2">
        <q-icon name="info" size="20px" color="primary" class="q-mr-sm" />
        <div class="text-caption text-grey-7">
          Tip: Use the search bar above to quickly find specific methods or examples
        </div>
        <q-space />
        <q-btn color="primary" label="Close" @click="closeModal" />
      </q-card-section>
    </q-card>
  </q-dialog>
</template>

<style scoped lang="scss">
.mapex-validation-help-modal {
  .section-card {
    border-radius: var(--mapex-radius-md);
    transition: box-shadow 0.3s ease;

    &:hover {
      box-shadow: var(--mapex-shadow-md);
    }
  }

  .code-block {
    border-radius: var(--mapex-radius-sm);
    overflow: hidden;

    pre {
      font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
      font-size: 13px;
      line-height: 1.6;
      overflow-x: auto;
      white-space: pre-wrap;
      word-wrap: break-word;
    }

    code {
      font-family: inherit;
    }
  }

  :deep(.q-tab) {
    text-transform: capitalize;
  }

  :deep(.q-tab-panel) {
    padding: 0;
  }
}
</style>
