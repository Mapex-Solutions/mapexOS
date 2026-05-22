import 'reflect-metadata';
import { findSensitiveDefaultsInUse } from '@mapexos/microservices';
import { defaultConfiguration } from './shared/configuration/application/configMap';

const sensitiveKeysForService = [
	'nats_password',
	'auth_secret',
	'internal_api_key',
	'minio_access_key',
	'minio_secret_key',
];

describe('js-workflow-executor config security', () => {
	it('marks every curated sensitive key with sensitive: true', () => {
		const marked = new Set<string>();
		for (const def of defaultConfiguration) {
			if (def.sensitive) marked.add(def.key);
		}
		for (const key of sensitiveKeysForService) {
			expect(marked.has(key)).toBe(true);
		}
	});

	it('flags every sensitive key when no env vars are overridden', () => {
		const resolved: Record<string, unknown> = {};
		for (const def of defaultConfiguration) {
			resolved[def.key] = def.default;
		}

		const violations = findSensitiveDefaultsInUse(defaultConfiguration, resolved);
		expect(violations).toHaveLength(sensitiveKeysForService.length);
	});

	it('passes validation when all sensitive env vars are overridden', () => {
		const resolved: Record<string, unknown> = {};
		for (const def of defaultConfiguration) {
			if (def.sensitive) {
				resolved[def.key] = `PROD_OVERRIDE_${def.key}`;
			} else {
				resolved[def.key] = def.default;
			}
		}

		expect(findSensitiveDefaultsInUse(defaultConfiguration, resolved)).toEqual([]);
	});
});
