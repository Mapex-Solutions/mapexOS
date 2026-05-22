import type { ConfigDefinition } from './types';
import { findSensitiveDefaultsInUse, isDevEnv, paintRed, resolveNodeEnv } from './validation';
import { ConfigModule } from './config';

// --- findSensitiveDefaultsInUse (pure, env-agnostic) ---

describe('findSensitiveDefaultsInUse', () => {
	it('returns violations for sensitive defaults', () => {
		const defs: ConfigDefinition[] = [
			{ key: 'auth_secret', env: 'AUTH_SECRET', type: 'string', default: 'dev-secret', sensitive: true },
			{ key: 'nats_password', env: 'NATS_PASSWORD', type: 'string', default: 'service_secret', sensitive: true },
		];
		const current = {
			auth_secret: 'dev-secret',
			nats_password: 'service_secret',
		};

		const violations = findSensitiveDefaultsInUse(defs, current);
		expect(violations).toHaveLength(2);
		expect(violations.map((v) => v.env).sort()).toEqual(['AUTH_SECRET', 'NATS_PASSWORD']);
	});

	it('ignores non-sensitive keys even when equal to default', () => {
		const defs: ConfigDefinition[] = [
			{ key: 'http_port', env: 'HTTP_PORT', type: 'int', default: 5000 },
			{ key: 'service_name', env: 'SERVICE_NAME', type: 'string', default: 'js-executor' },
		];
		const current = { http_port: 5000, service_name: 'js-executor' };

		expect(findSensitiveDefaultsInUse(defs, current)).toEqual([]);
	});

	it('ignores overridden sensitive values', () => {
		const defs: ConfigDefinition[] = [
			{ key: 'auth_secret', env: 'AUTH_SECRET', type: 'string', default: 'dev-secret', sensitive: true },
		];
		const current = { auth_secret: 'prod-real-9f4a' };

		expect(findSensitiveDefaultsInUse(defs, current)).toEqual([]);
	});

	it('handles all supported value types via deep equality', () => {
		const defs: ConfigDefinition[] = [
			{ key: 'str_key', env: 'STR_KEY', type: 'string', default: 'default-str', sensitive: true },
			{ key: 'int_key', env: 'INT_KEY', type: 'int', default: 42, sensitive: true },
			{ key: 'bool_key', env: 'BOOL_KEY', type: 'bool', default: true, sensitive: true },
			{ key: 'arr_key', env: 'ARR_KEY', type: 'array', default: ['a', 'b'], sensitive: true },
			{ key: 'json_key', env: 'JSON_KEY', type: 'json', default: { k: 'v' }, sensitive: true },
		];
		const current = {
			str_key: 'default-str',
			int_key: 42,
			bool_key: true,
			arr_key: ['a', 'b'],
			json_key: { k: 'v' },
		};

		expect(findSensitiveDefaultsInUse(defs, current)).toHaveLength(5);
	});
});

// --- isDevEnv ---

describe('isDevEnv', () => {
	test.each([
		['', true],
		['dev', true],
		['development', true],
		['DEV', false], // case-sensitive on purpose
		['local', false], // not an alias — explicit decision
		['prod', false],
		['production', false],
		['staging', false],
		['qa', false],
		['test', false],
		['develpoment', false], // typo defaults to fatal
		['dev ', false], // trailing space is suspicious
	])('isDevEnv(%j) === %s', (input, expected) => {
		expect(isDevEnv(input as string)).toBe(expected);
	});
});

// --- paintRed ---

describe('paintRed', () => {
	const originalNoColor = process.env.NO_COLOR;
	afterEach(() => {
		if (originalNoColor === undefined) delete process.env.NO_COLOR;
		else process.env.NO_COLOR = originalNoColor;
	});

	it('wraps text in ANSI bold-red codes', () => {
		delete process.env.NO_COLOR;
		const out = paintRed('hi');
		expect(out).toContain('\x1b[1;31m');
		expect(out).toContain('\x1b[0m');
		expect(out).toContain('hi');
	});

	it('honors NO_COLOR by returning plain text', () => {
		process.env.NO_COLOR = '1';
		expect(paintRed('hi')).toBe('hi');
	});
});

// --- resolveNodeEnv ---

describe('resolveNodeEnv', () => {
	const originalNodeEnv = process.env.NODE_ENV;
	afterEach(() => {
		if (originalNodeEnv === undefined) delete process.env.NODE_ENV;
		else process.env.NODE_ENV = originalNodeEnv;
	});

	it('prefers config[node_env] when present and non-empty', () => {
		process.env.NODE_ENV = 'test';
		expect(resolveNodeEnv({ node_env: 'dev' })).toBe('dev');
	});

	it('falls back to process.env.NODE_ENV when config is missing', () => {
		process.env.NODE_ENV = 'staging';
		expect(resolveNodeEnv({})).toBe('staging');
	});

	it('falls back to empty string when both sources are absent', () => {
		delete process.env.NODE_ENV;
		expect(resolveNodeEnv({})).toBe('');
	});

	it('treats empty config value as falsy and falls back', () => {
		process.env.NODE_ENV = 'qa';
		expect(resolveNodeEnv({ node_env: '' })).toBe('qa');
	});
});

// --- ConfigModule integration ---

describe('ConfigModule security guard', () => {
	beforeEach(() => {
		ConfigModule.__resetForTest();
	});
	afterEach(() => {
		ConfigModule.__resetForTest();
	});

	it('warns (does not fatal) in dev when sensitive defaults are in use', () => {
		let warnMsg = '';
		let fatalCalled = false;

		ConfigModule.init(
			[
				{ key: 'node_env', env: 'FAKE_NODE_ENV_DEV_UNSET', type: 'string', default: 'dev' },
				{
					key: 'auth_secret',
					env: 'FAKE_AUTH_SECRET_DEV_UNSET',
					type: 'string',
					default: 'dev-secret',
					sensitive: true,
				},
			],
			{
				onWarn: (m) => {
					warnMsg = m;
				},
				onFatal: () => {
					fatalCalled = true;
				},
			},
		);

		expect(fatalCalled).toBe(false);
		expect(warnMsg).toContain('SECURITY WARNING');
		expect(warnMsg).toContain('FAKE_AUTH_SECRET_DEV_UNSET');
	});

	it('fatals when NODE_ENV is non-dev and sensitive default is in use', () => {
		process.env.FAKE_NODE_ENV_PROD_SET = 'prod';
		try {
			let fatalMsg = '';
			let warnCalled = false;

			ConfigModule.init(
				[
					{ key: 'node_env', env: 'FAKE_NODE_ENV_PROD_SET', type: 'string', default: 'dev' },
					{
						key: 'auth_secret',
						env: 'FAKE_AUTH_SECRET_PROD_UNSET',
						type: 'string',
						default: 'dev-secret',
						sensitive: true,
					},
				],
				{
					onFatal: (m) => {
						fatalMsg = m;
					},
					onWarn: () => {
						warnCalled = true;
					},
				},
			);

			expect(warnCalled).toBe(false);
			expect(fatalMsg).toContain('[SECURITY]');
			expect(fatalMsg).toContain('FAKE_AUTH_SECRET_PROD_UNSET');
			expect(fatalMsg).toContain('prod');
		} finally {
			delete process.env.FAKE_NODE_ENV_PROD_SET;
		}
	});

	it('fatals on typo NODE_ENV — fail-closed posture', () => {
		process.env.FAKE_NODE_ENV_TYPO = 'develpoment';
		try {
			let fatalCalled = false;

			ConfigModule.init(
				[
					{ key: 'node_env', env: 'FAKE_NODE_ENV_TYPO', type: 'string', default: 'dev' },
					{
						key: 'auth_secret',
						env: 'FAKE_AUTH_SECRET_TYPO',
						type: 'string',
						default: 'dev-secret',
						sensitive: true,
					},
				],
				{
					onFatal: () => {
						fatalCalled = true;
					},
					onWarn: () => {},
				},
			);

			expect(fatalCalled).toBe(true);
		} finally {
			delete process.env.FAKE_NODE_ENV_TYPO;
		}
	});

	it('passes silently when prod sets all sensitive env vars to non-default values', () => {
		process.env.FAKE_NODE_ENV_OK = 'prod';
		process.env.FAKE_AUTH_SECRET_OK = 'prod-real-9f4a';
		try {
			let fatalCalled = false;
			let warnCalled = false;

			const cfg = ConfigModule.init(
				[
					{ key: 'node_env', env: 'FAKE_NODE_ENV_OK', type: 'string', default: 'dev' },
					{
						key: 'auth_secret',
						env: 'FAKE_AUTH_SECRET_OK',
						type: 'string',
						default: 'dev-secret',
						sensitive: true,
					},
				],
				{
					onFatal: () => {
						fatalCalled = true;
					},
					onWarn: () => {
						warnCalled = true;
					},
				},
			);

			expect(fatalCalled).toBe(false);
			expect(warnCalled).toBe(false);
			expect(cfg.get('auth_secret')).toBe('prod-real-9f4a');
		} finally {
			delete process.env.FAKE_NODE_ENV_OK;
			delete process.env.FAKE_AUTH_SECRET_OK;
		}
	});
});
