/**
 * Logger Composable
 *
 * Provides centralized logging functionality with context prefix.
 *
 * Usage:
 * ```typescript
 * const logger = useLogger('MyComponent');
 * logger.info('User logged in', { userId: '123' });
 * logger.error('Failed to fetch data', error);
 * ```
 */

/** Log levels supported by the logger */
export type LogLevel = 'debug' | 'info' | 'warn' | 'error';

/**
 * Logger interface with methods for each log level
 */
export interface Logger {
  /** Log debug message */
  debug: (message: string, data?: unknown) => void;
  /** Log info message */
  info: (message: string, data?: unknown) => void;
  /** Log warning message */
  warn: (message: string, data?: unknown) => void;
  /** Log error message */
  error: (message: string, data?: unknown) => void;
}

/**
 * Create a log function for a specific level
 *
 * @param {string} context - Component/module context
 * @param {LogLevel} level - Log level
 * @returns {Function} Log function
 */
function createLogFn(context: string, level: LogLevel): (message: string, data?: unknown) => void {
  const prefix = `[${context}]`;

  return (message: string, data?: unknown): void => {
    if (data !== undefined) {
      console[level](prefix, message, data);
    } else {
      console[level](prefix, message);
    }
  };
}

/**
 * Logger composable
 *
 * Creates a logger instance with context prefix.
 *
 * @param {string} context - Context name (usually component/module name)
 * @returns {Logger} Logger instance with debug, info, warn, error methods
 *
 * @example
 * ```typescript
 * // In a component
 * const logger = useLogger('TriggerLogsPage');
 *
 * // Log with different levels
 * logger.debug('Fetching data...');
 * logger.info('Data loaded', { count: items.length });
 * logger.warn('Cache miss, fetching from API');
 * logger.error('Failed to load', error);
 * ```
 */
export function useLogger(context: string): Logger {
  return {
    debug: createLogFn(context, 'debug'),
    info: createLogFn(context, 'info'),
    warn: createLogFn(context, 'warn'),
    error: createLogFn(context, 'error'),
  };
}
