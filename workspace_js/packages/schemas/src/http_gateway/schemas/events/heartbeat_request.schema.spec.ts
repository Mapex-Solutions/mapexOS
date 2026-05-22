import { ZodHeartbeatRequestSchema } from './heartbeat_request.schema';

it('Should validate successfully when assetUUID is a non-empty string', () => {
  const validPayload = { assetUUID: '12345678' };
  expect(() => ZodHeartbeatRequestSchema.parse(validPayload)).not.toThrow();
});

it('Should fail validation when assetUUID is missing', () => {
  const invalidPayload = {};
  const result = ZodHeartbeatRequestSchema.safeParse(invalidPayload);
  expect(result.success).toBe(false);
});

it('Should fail validation when assetUUID is an empty string', () => {
  const invalidPayload = { assetUUID: '' };
  const result = ZodHeartbeatRequestSchema.safeParse(invalidPayload);
  expect(result.success).toBe(false);
});
