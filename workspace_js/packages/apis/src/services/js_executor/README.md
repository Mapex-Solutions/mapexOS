# JSExecutor Service API

This package provides a type-safe API wrapper for the JSExecutor service, following the MapexOS API integration pattern.

## Structure

```
js_executor/
├── scripts/
│   ├── scripts.api.ts    # Script test API implementation
│   └── index.ts          # Scripts module exports
├── api.ts                # Main service API factory
├── index.ts              # Service exports
└── README.md             # This file
```

## Usage

### Import the API

```typescript
import { apis } from '@services/mapex';

// Access jsExecutor API
const jsExecutor = apis.jsExecutor;
```

### Test Script Execution

The jsExecutor service provides a `test` endpoint for testing JavaScript code execution with decode, validation, and transform scripts.

```typescript
import type { ScriptTest } from '@mapexos/schemas';

// Prepare test data
const scriptTest: ScriptTest = {
  decode: 'function decode(input) { return input; }',
  validation: 'function validate(data) { return true; }',
  transform: 'function transform(data) { return data; }',
  event: {
    // Your event data here
    deviceId: 'device-123',
    payload: { temperature: 25 }
  }
};

// Call the test endpoint
try {
  const result = await apis.jsExecutor.scripts.test(scriptTest);
  console.log('Test result:', result);
} catch (error) {
  console.error('Test failed:', error);
}
```

## Schema

The API uses `ZodScriptTestSchema` for validation:

```typescript
{
  decode: string | optional,      // Decode script (optional)
  validation: string | optional,  // Validation script (optional)
  transform: string,              // Transform script (required)
  event: object                   // Event data (required)
}
```

## Configuration

The service is configured in `/workspace_js/apps/mapexOS/src/services/mapex/index.ts`:

```typescript
jsExecutor: { baseURL: jsExecutorBaseURL }
```

Base URL is set via environment variable `JSEXECUTOR_API_BASE_URL` (default: `http://localhost:5003`)

## API Methods

### scripts.test()

Tests JavaScript code execution with provided scripts and event data.

**Endpoint**: `POST /api/v1/scripts/test`

**Parameters**:
- `body` (ScriptTest): Script test configuration
  - `decode` (string, optional): Decode script
  - `validation` (string, optional): Validation script
  - `transform` (string, required): Transform script
  - `event` (object, required): Event data to process

**Returns**: Response data from backend (type to be defined)

**Example**:

```typescript
const result = await apis.jsExecutor.scripts.test({
  decode: '',
  validation: '',
  transform: 'function transform(data) { return { ...data, processed: true }; }',
  event: { deviceId: 'test-001', value: 42 }
});
```

## Error Handling

All API calls follow the standard MapexOS error handling pattern:

```typescript
async function testScript(scriptData: ScriptTest) {
  try {
    loading.value = true;

    const result = await apis.jsExecutor.scripts.test(scriptData);

    // Handle success
    notifySuccess({ message: 'Script test successful' });
    return result;

  } catch (err: any) {
    console.error('Script test failed:', err);

    // Handle validation errors
    if (err.name === 'SchemaError') {
      notifyFail({ message: 'Validation error: ' + err.message });
    } else {
      notifyFail({ message: err.message || 'Failed to test script' });
    }
  } finally {
    loading.value = false;
  }
}
```

## Authentication

The API automatically includes:
- **Authorization Header**: JWT token from auth store
- **X-Org-Context Header**: Selected organization ID

This is handled by the global API interceptor in `/workspace_js/apps/mapexOS/src/services/mapex/index.ts`.

## Type Safety

All methods are fully typed using schemas from `@mapexos/schemas`:

```typescript
import type { ScriptTest } from '@mapexos/schemas';
import { ZodScriptTestSchema } from '@mapexos/schemas';
```

The `createApiFactory` automatically validates:
- Body parameters against `ZodScriptTestSchema`
- Returns type-safe responses
- Provides TypeScript autocomplete

## Development Notes

1. **Schema Location**: `/workspace_js/packages/schemas/src/jsexecutor/schemas/scripts/scripts.schema.ts`
2. **Type Location**: `/workspace_js/packages/schemas/src/jsexecutor/types/scripts/scripts.type.ts`
3. **Backend Contract**: To be added in `/workspace_go/packages/contracts/services/jsexecutor/`

## TODO

- [ ] Define response type based on backend implementation (currently `any`)
- [ ] Add backend contract documentation
- [ ] Add integration tests
- [ ] Document actual backend endpoint behavior
