import { getEnv, streamName, subject, durable } from './naming';

describe('Naming Helpers', () => {
	let originalGoEnv: string | undefined;

	beforeEach(() => {
		originalGoEnv = process.env.GO_ENV;
	});

	afterEach(() => {
		if (originalGoEnv === undefined) {
			delete process.env.GO_ENV;
		} else {
			process.env.GO_ENV = originalGoEnv;
		}
	});

	describe('getEnv', () => {
		it('returns "dev" when GO_ENV is undefined', () => {
			delete process.env.GO_ENV;
			expect(getEnv()).toBe('dev');
		});

		it('returns "dev" when GO_ENV is empty string', () => {
			process.env.GO_ENV = '';
			expect(getEnv()).toBe('dev');
		});

		it('returns the explicit value when GO_ENV is set', () => {
			process.env.GO_ENV = 'prod';
			expect(getEnv()).toBe('prod');
		});

		it('returns "qa" when GO_ENV=qa', () => {
			process.env.GO_ENV = 'qa';
			expect(getEnv()).toBe('qa');
		});
	});

	describe('streamName', () => {
		it('builds canonical name with default env', () => {
			delete process.env.GO_ENV;
			expect(streamName('JSWORKFLOWEXECUTOR', 'CODE')).toBe('DEV-MAPEXOS-JSWORKFLOWEXECUTOR-CODE');
		});

		it('builds canonical name with explicit prod env', () => {
			process.env.GO_ENV = 'prod';
			expect(streamName('JSWORKFLOWEXECUTOR', 'CODE')).toBe('PROD-MAPEXOS-JSWORKFLOWEXECUTOR-CODE');
		});

		it('uppercases mixed-case service and context', () => {
			delete process.env.GO_ENV;
			expect(streamName('jsworkflowexecutor', 'Code')).toBe('DEV-MAPEXOS-JSWORKFLOWEXECUTOR-CODE');
		});

		it('omits trailing dash when context is empty', () => {
			delete process.env.GO_ENV;
			expect(streamName('DLQ', '')).toBe('DEV-MAPEXOS-DLQ');
		});

		it('preserves multi-token context segments', () => {
			delete process.env.GO_ENV;
			expect(streamName('ASSETS', 'HEALTH-MONITOR')).toBe('DEV-MAPEXOS-ASSETS-HEALTH-MONITOR');
		});
	});

	describe('subject', () => {
		it('builds lowercase subject with default env', () => {
			delete process.env.GO_ENV;
			expect(subject('workflow', 'code')).toBe('dev.mapexos.workflow.code');
		});

		it('builds lowercase subject with explicit prod env', () => {
			process.env.GO_ENV = 'prod';
			expect(subject('workflow', 'code')).toBe('prod.mapexos.workflow.code');
		});

		it('lowercases mixed-case service and action', () => {
			delete process.env.GO_ENV;
			expect(subject('Workflow', 'Code')).toBe('dev.mapexos.workflow.code');
		});

		it('preserves dotted action tokens', () => {
			delete process.env.GO_ENV;
			expect(subject('mapexos', 'fanout.workflow.definition.invalidate')).toBe(
				'dev.mapexos.mapexos.fanout.workflow.definition.invalidate'
			);
		});
	});

	describe('durable', () => {
		it('builds lowercase durable with default env', () => {
			delete process.env.GO_ENV;
			expect(durable('jsworkflowexecutor', 'workflow-code')).toBe(
				'dev-jsworkflowexecutor-workflow-code-consumer'
			);
		});

		it('builds lowercase durable with explicit prod env', () => {
			process.env.GO_ENV = 'prod';
			expect(durable('jsworkflowexecutor', 'workflow-code')).toBe(
				'prod-jsworkflowexecutor-workflow-code-consumer'
			);
		});

		it('lowercases mixed-case service and context', () => {
			delete process.env.GO_ENV;
			expect(durable('JsWorkflowExecutor', 'Code')).toBe(
				'dev-jsworkflowexecutor-code-consumer'
			);
		});

		it('preserves multi-token context', () => {
			delete process.env.GO_ENV;
			expect(durable('jsworkflowexecutor', 'definition-invalidate')).toBe(
				'dev-jsworkflowexecutor-definition-invalidate-consumer'
			);
		});
	});
});
