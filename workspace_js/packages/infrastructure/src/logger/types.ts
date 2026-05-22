/**
 * Logger types for infrastructure package.
 */
export { Logger } from 'pino';

export interface LoggerOptions {
	level: string;
	environment: string;
	serviceName: string;
	serviceVersion: string;
}

export const LOGGER_TOKEN = Symbol('logger');
