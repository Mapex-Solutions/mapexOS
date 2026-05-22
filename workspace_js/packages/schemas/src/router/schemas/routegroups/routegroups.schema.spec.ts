import { ZodRouteGroupQuerySchema } from './routegroups.schema';

/**
 * ZodRouteGroupQuerySchema — `kinds` filter coverage.
 * Mirrors Go: packages/contracts/services/router/routegroups/dto.go::RouteGroupQuery.Kinds.
 * The asset wizard's Health step (HealthMonitoringSection) relies on this
 * filter to surface only RouteGroups whose every router.kind is acceptable
 * to validateHealthMonitorConfig.
 */
describe('ZodRouteGroupQuerySchema — kinds field', () => {
	it('passes and preserves the array when kinds contains valid enum values', () => {
		const result = ZodRouteGroupQuerySchema.safeParse({
			kinds: ['trigger', 'workflow'],
		});
		expect(result.success).toBe(true);
		if (result.success) {
			expect(result.data.kinds).toEqual(['trigger', 'workflow']);
		}
	});

	it('fails when kinds contains an invalid enum value', () => {
		const result = ZodRouteGroupQuerySchema.safeParse({
			kinds: ['invalid_kind'],
		});
		expect(result.success).toBe(false);
		if (!result.success) {
			expect(result.error.issues[0].path).toContain('kinds');
		}
	});

	it('passes with kinds undefined (omitted) — field is optional', () => {
		const result = ZodRouteGroupQuerySchema.safeParse({});
		expect(result.success).toBe(true);
		if (result.success) {
			expect(result.data.kinds).toBeUndefined();
		}
	});
});
