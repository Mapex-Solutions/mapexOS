import type { ConfigDefinition } from './types';
import type { LoggerOptions } from '@src/logger';
import type { NatsConnectionOptions } from '@mapexos/infrastructure';
import { findSensitiveDefaultsInUse, isDevEnv, paintRed, resolveNodeEnv } from './validation';

/**
 * Optional hooks for the security guard. Callers (typically test code)
 * may inject non-exiting stubs to assert the production-default
 * validator without killing the process. When omitted, ConfigModule
 * falls back to console.error+process.exit(1) for fatal and console.warn
 * for the dev banner.
 */
export type ConfigSecurityHandlers = {
	onFatal?: (message: string) => void;
	onWarn?: (message: string) => void;
};

const defaultOnFatal = (message: string): void => {
	console.error(message);
	process.exit(1);
};

const defaultOnWarn = (message: string): void => {
	console.warn(message);
};

/**
 * Constructs a new instance of the ConfigModule class.
 * This constructor is private to enforce the singleton pattern.
 *
 * @param definitions - An array of configuration definitions. Each definition specifies:
 *                      - `key`: The key for the configuration item.
 *                      - `env`: The environment variable name to fetch the value from.
 *                      - `type`: The type of the configuration item (e.g., 'string', 'int', 'bool', 'array', 'json').
 *                      - `default`: The default value to use if the environment variable is not set.
 *                      - `sensitive`: Optional flag that triggers the production-default guard.
 */
export class ConfigModule {
	private static instance: ConfigModule;
	private config: Record<string, any> = {};

	private constructor(
		definitions: ConfigDefinition[],
		handlers: Required<ConfigSecurityHandlers>,
	) {
		definitions.forEach((def) => {
			let value: any;

			switch (def.type) {
				case 'string':
					value = process.env[def.env] ?? def.default;
					break;
				case 'int':
					value = process.env[def.env]
						? parseInt(process.env[def.env]!, 10)
						: def.default;
					break;
				case 'bool':
					value = process.env[def.env]
						? process.env[def.env]!.toLowerCase() === 'true'
						: def.default;
					break;
				case 'array':
					value = process.env[def.env]
						? process.env[def.env]!.split(',')
						: def.default;
					break;
				case 'json':
					try {
						value = process.env[def.env]
							? JSON.parse(process.env[def.env]!)
							: def.default;
					} catch {
						value = def.default;
					}
					break;
				default:
					console.warn(`Type not supported: ${def.type} for key ${def.key}`);
					return;
			}

			this.config[def.key] = value;
		});

		const violations = findSensitiveDefaultsInUse(definitions, this.config);
		if (violations.length > 0) {
			const nodeEnv = resolveNodeEnv(this.config);
			const envs = violations.map((v) => v.env).join(', ');
			if (isDevEnv(nodeEnv)) {
				handlers.onWarn(
					`${paintRed('[SECURITY WARNING]')} NODE_ENV=${nodeEnv} — ${violations.length} sensitive env var(s) using DEV defaults: ${envs}. Safe for local development. NEVER deploy with these values.`,
				);
			} else {
				handlers.onFatal(
					`${paintRed('[SECURITY]')} refusing to start in NODE_ENV=${nodeEnv} — sensitive env vars using DEV defaults: ${envs}. Set them to production values before deploying.`,
				);
			}
		}
	}

	/**
	 * Initializes the ConfigModule singleton instance with the provided configuration definitions.
	 * If the instance already exists, it returns the existing instance.
	 *
	 * @param definitions - An array of configuration definitions that specify the key, environment variable name,
	 *                      type, and default value for each configuration item.
	 * @param hooks - Optional onFatal/onWarn callbacks for testing. Production callers omit this argument.
	 * @returns The singleton instance of the ConfigModule.
	 */
	public static init(
		definitions: ConfigDefinition[],
		hooks?: ConfigSecurityHandlers,
	): ConfigModule {
		if (!ConfigModule.instance) {
			const handlers: Required<ConfigSecurityHandlers> = {
				onFatal: hooks?.onFatal ?? defaultOnFatal,
				onWarn: hooks?.onWarn ?? defaultOnWarn,
			};
			ConfigModule.instance = new ConfigModule(definitions, handlers);
		}
		return ConfigModule.instance;
	}

	/**
	 * Test-only escape hatch — clears the singleton so a fresh instance
	 * can be built in another `init()` call. Never call from production
	 * code; the singleton invariant is intentional.
	 *
	 * @internal
	 */
	public static __resetForTest(): void {
		ConfigModule.instance = undefined as unknown as ConfigModule;
	}

	/**
	 * Retrieves the value associated with the specified configuration key.
	 *
	 * @param key - The key of the configuration item to retrieve.
	 * @returns The value of the configuration item associated with the given key.
	 */
	public get(key: string): any {
		return this.config[key];
	}

	/**
	 * Retrieves all configuration items as a key-value record.
	 *
	 * @returns A record containing all configuration items, where each key is a configuration key
	 *          and the value is the corresponding configuration value.
	 */
	public getAll(): Record<string, any> {
		return this.config;
	}

	/**
	 * Retrieves the logger configuration options.
	 *
	 * @returns {LoggerOptions} An object containing the logger configuration, including:
	 * - `level`: The logging level, set to 'trace' if the environment is 'dev', otherwise 'warning'.
	 * - `environment`: The current environment setting.
	 * - `serviceName`: The name of the service.
	 * - `serviceVersion`: The version of the service.
	 */
	public getLoggerConfig(): LoggerOptions {
		return {
			level: this.get('node_env') === 'dev' ? 'trace' : 'warning',
			environment: this.get('node_env'),
			serviceName: this.get('service_name'),
			serviceVersion: this.get('service_version'),
		};
	}

	/**
	 * Retrieves the NATS connection configuration options.
	 *
	 * @returns {NatsConnectionOptions} An object containing the NATS connection configuration, including:
	 * - `servers`: The URL of the NATS server.
	 * - `name`: The name of the NATS client.
	 * - `user`: The username for the NATS connection.
	 * - `pass`: The password for the NATS connection.
	 * - `maxReconnectAttempts`: The maximum number of reconnection attempts, set to -1 for unlimited.
	 * - `timeout`: The connection timeout in milliseconds.
	 */
	public getNatsConfig(): NatsConnectionOptions {
		return {
			servers: this.get('nats_url'),
			name: this.get('nats_client_name'),
			user: this.get('nats_username'),
			pass: this.get('nats_password'),
			maxReconnectAttempts: -1,
			timeout: 5000,
		} as NatsConnectionOptions;
	}
}
