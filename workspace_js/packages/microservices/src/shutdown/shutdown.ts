import type { Logger } from '@mapexos/infrastructure';
import type { Shutdowner, ShutdownHook } from './types';

/**
 * Groups hooks by their priority value.
 * Input must be sorted by priority ascending.
 */
function groupByPriority(hooks: ShutdownHook[]): ShutdownHook[][] {
	if (hooks.length === 0) return [];

	const groups: ShutdownHook[][] = [[hooks[0]]];

	for (let i = 1; i < hooks.length; i++) {
		if (hooks[i].priority !== hooks[i - 1].priority) {
			groups.push([hooks[i]]);
		} else {
			groups[groups.length - 1].push(hooks[i]);
		}
	}

	return groups;
}

/**
 * Manages the ordered graceful shutdown sequence.
 * Hooks are executed by priority (ascending): P0 first, P5 last.
 * Hooks with the same priority run concurrently.
 *
 * Recommended priorities:
 *
 *   P0 — HTTP Server (stop accepting, drain in-flight requests)
 *   P1 — Message Consumers (stop fetching, finish current batch)
 *   P2 — Background workers (Piscina pools, sweep loops)
 *   P3 — Publishers/flush (ensure pending messages are sent)
 *   P4 — Caches (TieredCache, in-memory caches)
 *   P5 — Connections (Redis, NATS, MinIO)
 */
export class ShutdownManager {
	private hooks: ShutdownHook[] = [];
	private logger: Logger;

	constructor(logger: Logger) {
		this.logger = logger;
	}

	/**
	 * Adds a Shutdowner component to the shutdown sequence.
	 *
	 * @param name - Hook identifier for logging
	 * @param priority - Execution order (ascending: P0 first, P5 last)
	 * @param shutdowner - Component implementing the Shutdowner interface
	 */
	register(name: string, priority: number, shutdowner: Shutdowner): void {
		this.registerFunc(name, priority, (signal) => shutdowner.shutdown(signal));
	}

	/**
	 * Adds a cleanup function to the shutdown sequence.
	 *
	 * @param name - Hook identifier for logging
	 * @param priority - Execution order (ascending: P0 first, P5 last)
	 * @param fn - Cleanup function that receives an optional AbortSignal
	 */
	registerFunc(name: string, priority: number, fn: (signal?: AbortSignal) => Promise<void>): void {
		this.hooks.push({ name, priority, fn });
	}

	/**
	 * Listens for SIGTERM/SIGINT and executes all registered hooks in priority
	 * order with the given timeout. Hooks with the same priority run concurrently.
	 *
	 * @param timeoutMs - Max duration in milliseconds for the entire shutdown sequence
	 */
	waitForSignal(timeoutMs: number): void {
		const handler = (signal: string) => {
			// Remove listeners to prevent double-fire
			process.removeListener('SIGTERM', handler);
			process.removeListener('SIGINT', handler);

			this.logger.info(`[SHUTDOWN] Received ${signal}, starting graceful shutdown (timeout: ${timeoutMs}ms)...`);

			this.executeShutdown(timeoutMs).then(() => {
				process.exit(0);
			}).catch(() => {
				process.exit(1);
			});
		};

		process.on('SIGTERM', () => handler('SIGTERM'));
		process.on('SIGINT', () => handler('SIGINT'));
	}

	/**
	 * Executes all registered hooks in priority order.
	 * Exposed for testing — in production, use waitForSignal().
	 *
	 * @param timeoutMs - Max duration in milliseconds for the entire shutdown sequence
	 */
	async executeShutdown(timeoutMs: number): Promise<void> {
		const ac = new AbortController();
		const timer = setTimeout(() => ac.abort(), timeoutMs);

		const hooks = [...this.hooks].sort((a, b) => a.priority - b.priority);
		const groups = groupByPriority(hooks);

		const start = Date.now();

		for (const group of groups) {
			if (ac.signal.aborted) {
				this.logger.warn('[SHUTDOWN] Timeout reached, aborting remaining hooks');
				break;
			}

			await Promise.all(
				group.map(async (hook) => {
					const hookStart = Date.now();
					try {
						await hook.fn(ac.signal);
						this.logger.info(`[SHUTDOWN] ${hook.name} done (${Date.now() - hookStart}ms)`);
					} catch (err) {
						this.logger.warn(`[SHUTDOWN] ${hook.name} failed (${Date.now() - hookStart}ms): ${err}`);
					}
				}),
			);
		}

		clearTimeout(timer);
		this.logger.info(`[SHUTDOWN] Graceful shutdown complete (${Date.now() - start}ms)`);
	}
}
