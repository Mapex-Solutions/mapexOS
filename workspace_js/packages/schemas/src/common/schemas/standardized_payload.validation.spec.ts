
import { ZodStandardizedPayloadSchema } from './standardized_payload.validation';

it('Should validate successfully when all fields are non-empty strings and objects', () => {
  const validPayload = {
    eventType: 'user.signup',
    eventId: '12345',
    data: { userId: 'abc123', plan: 'premium' },
    metadata: { source: 'web', campaign: 'spring_sale' },
    created: '2023-10-01T12:00:00Z'
  };

  expect(() => ZodStandardizedPayloadSchema.parse(validPayload)).not.toThrow();
});

it('Should fail validation when "eventType" is an empty string', () => {
  const invalidPayload = {
    eventType: '',
    eventId: 'validEventId',
    data: { key: 'value' },
    metadata: { key: 'value' },
    created: '2023-10-10T10:00:00Z'
  };

  const result = ZodStandardizedPayloadSchema.safeParse(invalidPayload);
  expect(result.success).toBe(false);
  expect(result.error?.issues[0].message).toContain('must not be empty');
});

it('Should fail validation when eventId is an empty string', () => {
  const invalidPayload = {
    eventType: 'someEventType',
    eventId: '',
    data: { key: 'value' },
    metadata: { metaKey: 'metaValue' },
    created: '2023-10-10T10:00:00Z'
  };

  const result = ZodStandardizedPayloadSchema.safeParse(invalidPayload);

  expect(result.success).toBe(false);
  if (!result.success) {
    expect(result.error.issues[0].path).toContain('eventId');
  }
});

it('should fail validation when "data" is an empty object', () => {
  const invalidPayload = {
    eventType: 'someEvent',
    eventId: '12345',
    data: {},
    metadata: { key: 'value' },
    created: '2023-10-10T10:00:00Z'
  };

  const result = ZodStandardizedPayloadSchema.safeParse(invalidPayload);
  expect(result.success).toBe(false);
  expect(result.error?.issues[0].path).toContain('data');
});

it('should fail validation when metadata is an empty object', () => {
  const invalidPayload = {
    eventType: 'eventTypeValue',
    eventId: 'eventIdValue',
    data: { key: 'value' },
    metadata: {}, // Empty object
    created: 'createdValue'
  };

  const result = ZodStandardizedPayloadSchema.safeParse(invalidPayload);
  expect(result.success).toBe(false);
  if (!result.success) {
    expect(result.error.issues[0].message).toContain('Object cannot be empty');
  }
});

it('Should fail validation when "eventType" is missing from the payload', () => {
  const payload = {
    eventId: '12345',
    data: { key: 'value' },
    metadata: { source: 'test' },
    created: '2023-10-01T12:00:00Z'
  };

  const result = ZodStandardizedPayloadSchema.safeParse(payload);
  expect(result.success).toBe(false);
  if (!result.success) {
    expect(result.error.issues[0].path).toContain('eventType');
  }
});

it('should fail validation when "metadata" is missing from the payload', () => {
	const invalidPayload = {
		eventType: 'test-event',
		eventId: 'test-id',
		data: { test: 'data' },
		created: '2024-01-01T00:00:00Z',
	};

	const result = ZodStandardizedPayloadSchema.safeParse(invalidPayload);
	expect(result.success).toBe(false);
});

test('Should fail validation when "created" is an empty string', () => {
  const invalidPayload = {
    eventType: 'eventTypeValue',
    eventId: 'eventIdValue',
    data: { key: 'value' },
    metadata: { key: 'value' },
    created: ''
  };

  const result = ZodStandardizedPayloadSchema.safeParse(invalidPayload);
  expect(result.success).toBe(false);
  if (!result.success) {
    expect(result.error.issues[0].message).toContain('must not be empty');
  }
});

it('should fail validation when "data" is missing from the payload', () => {
  const invalidPayload = {
    eventType: 'test-event',
    eventId: 'test-id',
    metadata: { meta: 'data' },
    created: '2024-01-01T00:00:00Z'
  };

  const result = ZodStandardizedPayloadSchema.safeParse(invalidPayload);
  expect(result.success).toBe(false);
});
