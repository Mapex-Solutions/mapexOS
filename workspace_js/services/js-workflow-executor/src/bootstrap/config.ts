import { container } from 'tsyringe';

import { ConfigModule, APP_TOKEN, CONFIG_TOKEN } from '@mapexos/microservices';
import { initLogger, LOGGER_TOKEN } from '@mapexos/infrastructure';

import { defaultConfiguration } from '@shared/configuration';

// InitConfig initializes the application configuration and logger.
// Registers CONFIG_TOKEN and LOGGER_TOKEN in the DI container.
/** Initializes ConfigModule and Logger, registers in DI container. */
export function initConfig() {
	// Initialize configuration (singleton)
	const configModule = ConfigModule.init(defaultConfiguration);

	// Initialize logger (LOG_LEVEL env overrides auto-derived level from NODE_ENV)
	const loggerOptions = configModule.getLoggerConfig();
	const logLevelOverride = configModule.get('log_level') as string;
	if (logLevelOverride) {
		loggerOptions.level = logLevelOverride;
	}
	const logger = initLogger(loggerOptions);

	// Register core dependencies in DI container
	container.register(CONFIG_TOKEN, { useValue: configModule });
	container.register(LOGGER_TOKEN, { useValue: logger });

	logger.info('[APP:BOOTSTRAP] JS Workflow Executor Service starting');
	logger.info('[APP:BOOTSTRAP] Logger initialized');

	return { configModule, logger };
}
