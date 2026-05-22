/**
 * Shutdowner is implemented by any component that needs graceful cleanup.
 * Infrastructure clients (NATS, Redis, Piscina, etc.) can implement this
 * to participate in the ordered shutdown sequence.
 */
export interface Shutdowner {
	shutdown(signal?: AbortSignal): Promise<void>;
}

/**
 * Holds a named cleanup function with an execution priority.
 */
export interface ShutdownHook {
	name: string;
	priority: number;
	fn: (signal?: AbortSignal) => Promise<void>;
}
