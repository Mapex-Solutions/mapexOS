/**
 * Minimal Logger interface compatible with the Mapex standard.
 * Mirrors the shape of @mapexos/microservices Logger for consistency.
 *
 * This mock service does not depend on @mapexos/microservices, so a local
 * console-backed implementation is provided below.
 */
export interface Logger {
	info(msg: string): void;
	warn(msg: string): void;
	error(msg: string, err?: Error): void;
	debug(msg: string): void;
}

/**
 * Creates a console-backed logger that emits structured output.
 * All log lines emitted by the service must carry a `[LAYER:Component]` prefix
 * as dictated by the /js-arch-back standard.
 */
export function createConsoleLogger(): Logger {
	return {
		info(msg: string): void {
			// eslint-disable-next-line no-console
			console.log(msg);
		},
		warn(msg: string): void {
			// eslint-disable-next-line no-console
			console.warn(msg);
		},
		error(msg: string, err?: Error): void {
			// eslint-disable-next-line no-console
			console.error(msg, err ?? '');
		},
		debug(msg: string): void {
			if (process.env.DEBUG) {
				// eslint-disable-next-line no-console
				console.log(msg);
			}
		},
	};
}
