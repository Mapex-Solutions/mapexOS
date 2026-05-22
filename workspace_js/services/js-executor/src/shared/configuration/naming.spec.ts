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
			expect(streamName('ASSETS', 'HEARTBEAT')).toBe('DEV-MAPEXOS-ASSETS-HEARTBEAT');
		});

		it('builds canonical name with explicit prod env', () => {
			process.env.GO_ENV = 'prod';
			expect(streamName('ASSETS', 'HEARTBEAT')).toBe('PROD-MAPEXOS-ASSETS-HEARTBEAT');
		});

		it('uppercases mixed-case service and context', () => {
			delete process.env.GO_ENV;
			expect(streamName('assets', 'Heartbeat')).toBe('DEV-MAPEXOS-ASSETS-HEARTBEAT');
		});

		it('omits trailing dash when context is empty', () => {
			delete process.env.GO_ENV;
			expect(streamName('DLQ', '')).toBe('DEV-MAPEXOS-DLQ');
		});

		it('preserves multi-token context segments', () => {
			delete process.env.GO_ENV;
			expect(streamName('ASSETS', 'HEALTH-MONITOR')).toBe('DEV-MAPEXOS-ASSETS-HEALTH-MONITOR');
		});

		it('returns "DEV-MAPEXOS-JSEXECUTOR-MQTTDATA" for the js-executor MQTT data stream', () => {
			delete process.env.GO_ENV;
			expect(streamName('JSEXECUTOR', 'MQTTDATA')).toBe('DEV-MAPEXOS-JSEXECUTOR-MQTTDATA');
		});
	});

	describe('subject', () => {
		it('builds lowercase subject with default env', () => {
			delete process.env.GO_ENV;
			expect(subject('events', 'save')).toBe('dev.mapexos.events.save');
		});

		it('builds lowercase subject with explicit prod env', () => {
			process.env.GO_ENV = 'prod';
			expect(subject('events', 'save')).toBe('prod.mapexos.events.save');
		});

		it('lowercases mixed-case service and action', () => {
			delete process.env.GO_ENV;
			expect(subject('Events', 'Save')).toBe('dev.mapexos.events.save');
		});

		it('preserves dotted action tokens', () => {
			delete process.env.GO_ENV;
			expect(subject('mapexos', 'fanout.asset.invalidate')).toBe(
				'dev.mapexos.mapexos.fanout.asset.invalidate'
			);
		});
	});

	describe('durable', () => {
		it('builds lowercase durable with default env', () => {
			delete process.env.GO_ENV;
			expect(durable('events', 'save')).toBe('dev-events-save-consumer');
		});

		it('builds lowercase durable with explicit prod env', () => {
			process.env.GO_ENV = 'prod';
			expect(durable('events', 'save')).toBe('prod-events-save-consumer');
		});

		it('lowercases mixed-case service and context', () => {
			delete process.env.GO_ENV;
			expect(durable('Events', 'Save')).toBe('dev-events-save-consumer');
		});

		it('preserves multi-token context', () => {
			delete process.env.GO_ENV;
			expect(durable('assets', 'mqtt-presence')).toBe('dev-assets-mqtt-presence-consumer');
		});
	});
});
