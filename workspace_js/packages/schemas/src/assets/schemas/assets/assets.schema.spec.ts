import { ZodAssetResponseSchema, ZodHealthMonitorConfigSchema } from './assets.schema';

/**
 * Refinement invariant: when HealthMonitor.enabled === true, the admin MUST
 * provide thresholdMinutes (≥ 10). Empty route-group arrays are valid —
 * monitor-only mode persists healthStatus to Mongo + ClickHouse without
 * publishing to mapexos.route.execute. Covered end-to-end at the backend
 * layer too (validateHealthMonitorConfig + 422) — this spec pins the
 * client-side immediate-feedback rule.
 */

describe('ZodHealthMonitorConfigSchema — enabled refinement', () => {
	it('passes when enabled=false even if both arrays are empty', () => {
		const result = ZodHealthMonitorConfigSchema.safeParse({
			enabled: false,
			thresholdMinutes: 10,
		});
		expect(result.success).toBe(true);
	});

	it('passes when enabled=true and both arrays are empty (monitor-only mode)', () => {
		const result = ZodHealthMonitorConfigSchema.safeParse({
			enabled: true,
			thresholdMinutes: 10,
			offlineRouteGroupIds: [],
			onlineRouteGroupIds: [],
		});
		expect(result.success).toBe(true);
	});

	it('passes when enabled=true with at least one offline route group', () => {
		const result = ZodHealthMonitorConfigSchema.safeParse({
			enabled: true,
			thresholdMinutes: 10,
			offlineRouteGroupIds: ['rg-offline-1'],
			onlineRouteGroupIds: [],
		});
		expect(result.success).toBe(true);
	});

	it('passes when enabled=true with at least one online route group', () => {
		const result = ZodHealthMonitorConfigSchema.safeParse({
			enabled: true,
			thresholdMinutes: 10,
			offlineRouteGroupIds: [],
			onlineRouteGroupIds: ['rg-online-1'],
		});
		expect(result.success).toBe(true);
	});

	it('passes when enabled=true with both offline and online arrays populated', () => {
		const result = ZodHealthMonitorConfigSchema.safeParse({
			enabled: true,
			thresholdMinutes: 10,
			offlineRouteGroupIds: ['rg-offline-1'],
			onlineRouteGroupIds: ['rg-online-1'],
		});
		expect(result.success).toBe(true);
	});

	it('passes when enabled=true and arrays omitted (defaults to [] via Zod default)', () => {
		// When the caller omits the arrays, Zod's .default([]) fills them in
		// to []. With the route-group-required rule dropped, this is a valid
		// monitor-only payload.
		const result = ZodHealthMonitorConfigSchema.safeParse({
			enabled: true,
			thresholdMinutes: 10,
		});
		expect(result.success).toBe(true);
	});

	/**
	 * Go-contract alignment (TKT-2026 zod-drift fix): all fields are
	 * field-level optional, mirroring HealthMonitorConfig pointers in
	 * packages/contracts/services/assets/assets/dto.go. The conditional
	 * "required when enabled=true" rule is enforced by the superRefine.
	 */
	it('passes when only { enabled: false } is provided (minimal toggle-off payload)', () => {
		const result = ZodHealthMonitorConfigSchema.safeParse({ enabled: false });
		expect(result.success).toBe(true);
	});

	it('passes when an entirely empty object {} is provided (no monitoring)', () => {
		const result = ZodHealthMonitorConfigSchema.safeParse({});
		expect(result.success).toBe(true);
	});

	it('fails ONLY on thresholdMinutes when enabled=true and nothing else is provided', () => {
		const result = ZodHealthMonitorConfigSchema.safeParse({ enabled: true });
		expect(result.success).toBe(false);
		if (!result.success) {
			const paths = result.error.issues.map((i) => i.path.join('.'));
			expect(paths).toContain('thresholdMinutes');
		}
	});

	it('passes when enabled=true with thresholdMinutes set and both arrays empty (monitor-only)', () => {
		const result = ZodHealthMonitorConfigSchema.safeParse({
			enabled: true,
			thresholdMinutes: 10,
			offlineRouteGroupIds: [],
			onlineRouteGroupIds: [],
		});
		expect(result.success).toBe(true);
	});

	it('passes when enabled=true with thresholdMinutes and one offline route group', () => {
		const result = ZodHealthMonitorConfigSchema.safeParse({
			enabled: true,
			thresholdMinutes: 10,
			offlineRouteGroupIds: ['x'],
		});
		expect(result.success).toBe(true);
	});

	it('passes for a full valid payload (no regression for existing callers)', () => {
		const result = ZodHealthMonitorConfigSchema.safeParse({
			enabled: true,
			thresholdMinutes: 10,
			requiredMisses: 3,
			offlineRouteGroupIds: ['x'],
		});
		expect(result.success).toBe(true);
	});
});

/**
 * ZodAssetResponseSchema — healthStatusChangedAt coverage.
 * Added alongside the Go-side `healthStatusChangedAt` persistence (TKT-2026-0030):
 * the API surfaces the last-flip timestamp so the UI can render "offline for X".
 */
describe('ZodAssetResponseSchema — healthStatusChangedAt', () => {
	it('accepts a valid ISO string', () => {
		const result = ZodAssetResponseSchema.safeParse({
			healthStatus: 'offline',
			healthStatusChangedAt: '2026-04-22T10:30:00Z',
		});
		expect(result.success).toBe(true);
	});

	it('accepts null (asset never transitioned)', () => {
		const result = ZodAssetResponseSchema.safeParse({
			healthStatus: 'unknown',
			healthStatusChangedAt: null,
		});
		expect(result.success).toBe(true);
	});

	it('accepts when healthStatusChangedAt is omitted', () => {
		const result = ZodAssetResponseSchema.safeParse({
			healthStatus: 'online',
		});
		expect(result.success).toBe(true);
	});
});

/**
 * heartbeatMode (TKT-2026-0034) controls who emits heartbeats:
 *   - 'implicit' (default): js-executor emits per data event
 *   - 'explicit': device sends heartbeat via MQTT or HTTP
 * Field is optional with .default('implicit') applied at parse time.
 */
describe('ZodHealthMonitorConfigSchema — heartbeatMode', () => {
	it('parses successfully when heartbeatMode is explicit', () => {
		const result = ZodHealthMonitorConfigSchema.safeParse({
			enabled: true,
			thresholdMinutes: 10,
			heartbeatMode: 'explicit',
		});
		expect(result.success).toBe(true);
		if (result.success) {
			expect(result.data.heartbeatMode).toBe('explicit');
		}
	});

	it('rejects an invalid heartbeatMode with invalid_enum_value', () => {
		const result = ZodHealthMonitorConfigSchema.safeParse({
			enabled: true,
			thresholdMinutes: 10,
			heartbeatMode: 'other',
		});
		expect(result.success).toBe(false);
		if (!result.success) {
			const heartbeatModeIssue = result.error.issues.find(
				(i) => i.path.join('.') === 'heartbeatMode',
			);
			expect(heartbeatModeIssue).toBeDefined();
			// Zod v4 emits 'invalid_value' for enum mismatches (was 'invalid_enum_value' in v3).
			expect(heartbeatModeIssue?.code).toBe('invalid_value');
		}
	});

	it("applies default 'implicit' when heartbeatMode is omitted", () => {
		const result = ZodHealthMonitorConfigSchema.safeParse({
			enabled: true,
			thresholdMinutes: 10,
		});
		expect(result.success).toBe(true);
		if (result.success) {
			expect(result.data.heartbeatMode).toBe('implicit');
		}
	});
});
