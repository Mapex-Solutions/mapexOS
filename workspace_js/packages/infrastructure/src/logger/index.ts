/**
 * Infrastructure Logger Module
 *
 * Single source of truth for logging across all packages.
 * When initLogger() is called, it automatically configures the global
 * infrastructure logger used by NATS, MinIO, Redis, TieredCache, etc.
 *
 * Usage:
 * ```typescript
 * import { initLogger, LOGGER_TOKEN } from '@mapexos/infrastructure';
 *
 * const logger = initLogger({
 *   level: 'info',
 *   environment: 'production',
 *   serviceName: 'my-service',
 *   serviceVersion: '1.0.0',
 * });
 *
 * // Register in DI container
 * container.register(LOGGER_TOKEN, { useValue: logger });
 *
 * // All infrastructure logs now use this logger automatically
 * ```
 */

import pino from 'pino';
import type { Logger, LoggerOptions } from './types';

export { Logger, LoggerOptions, LOGGER_TOKEN } from './types';

/** Silent logger - discards all messages */
const silentLogger = pino({ level: 'silent' });

/** Global infrastructure logger instance */
let _logger: Logger = silentLogger;

/**
 * Initializes and returns a logger instance configured with the specified options.
 * Follows workspace_go pattern with structured JSON logs.
 *
 * IMPORTANT: This function automatically sets the global infrastructure logger.
 * All infrastructure modules (NATS, MinIO, Redis, TieredCache) will use this logger.
 *
 * @param options - Configuration options for the logger.
 * @param options.level - The logging level (e.g., 'info', 'debug', 'error').
 * @param options.environment - The environment in which the service is running.
 * @param options.serviceName - The name of the service using the logger.
 * @param options.serviceVersion - The version of the service using the logger.
 * @returns A configured pino logger instance.
 */
export function initLogger(options: LoggerOptions): Logger {
	const logger = pino({
		level: options.level || 'info',
		base: {
			env: options.environment,
			service: options.serviceName,
			version: options.serviceVersion,
		},
		timestamp: () => `,"time":"${new Date().toISOString()}"`,
		formatters: {
			level: (label) => {
				return { level: label };
			},
		},
		messageKey: 'message',
	});

	// Automatically set as global infrastructure logger
	_logger = logger;

	return logger;
}

/**
 * Get the current infrastructure logger.
 * Returns silent logger if initLogger() was not called.
 */
export function getInfraLogger(): Logger {
	return _logger;
}

/**
 * Manually set the global infrastructure logger.
 * Use this if you need to set a custom logger instance.
 *
 * @param logger - pino Logger instance
 */
export function setInfraLogger(logger: Logger): void {
	_logger = logger;
}
