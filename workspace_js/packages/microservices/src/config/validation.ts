/**
 * Production-default guard for the config singleton.
 *
 * Many services declare dev-friendly `default` values for credentials and
 * secrets so they boot out of the box in local development. Those defaults
 * must never reach production. This module provides the detection logic
 * used by ConfigModule to refuse startup when a sensitive key is still
 * using its hardcoded default in a non-dev environment, and to emit a
 * visible warning when the same condition is detected in local dev.
 */

import type { ConfigDefinition } from './types';

/** Violation identifies a sensitive key whose resolved value still equals its default. */
export type Violation = {
	key: string;
	env: string;
};

/**
 * Returns every sensitive=true definition whose resolved value in
 * `current` still equals its hardcoded `default`. The function is pure:
 * it reads no environment variables and touches no module state, so
 * callers decide what to do with the result based on the runtime
 * environment (see {@link isDevEnv}).
 */
export function findSensitiveDefaultsInUse(
	defs: ConfigDefinition[],
	current: Record<string, unknown>,
): Violation[] {
	const out: Violation[] = [];
	for (const d of defs) {
		if (!d.sensitive) continue;
		if (deepEqual(current[d.key], d.default)) {
			out.push({ key: d.key, env: d.env });
		}
	}
	return out;
}

/**
 * Reports whether `nodeEnv` represents a local/development environment
 * where hardcoded dev defaults are tolerated (with a visible warning).
 * Recognized values: "" (uninitialized / tests), "dev", "development".
 * Any other value — "staging", "qa", "prod", or typos like "develpoment"
 * — is treated as non-dev and triggers a fatal abort when sensitive
 * defaults are still in use.
 */
export function isDevEnv(nodeEnv: string): boolean {
	return nodeEnv === '' || nodeEnv === 'dev' || nodeEnv === 'development';
}

/**
 * Wraps text in ANSI bold-red so the prefix tag stands out in terminal
 * output. Honors the NO_COLOR convention (https://no-color.org) so log
 * capture in CI or piped output stays clean.
 */
export function paintRed(s: string): string {
	if (process.env.NO_COLOR) return s;
	return `[1;31m${s}[0m`;
}

/**
 * Resolves the effective NODE_ENV the security guard should branch on.
 * Precedence: config['node_env'] (the resolved value, which already
 * factors in env-var overrides) → process.env.NODE_ENV → "".
 */
export function resolveNodeEnv(current: Record<string, unknown>): string {
	const fromConfig = current['node_env'];
	if (typeof fromConfig === 'string' && fromConfig !== '') return fromConfig;
	return process.env.NODE_ENV ?? '';
}

function deepEqual(a: unknown, b: unknown): boolean {
	if (a === b) return true;
	if (a === null || b === null) return false;
	if (typeof a !== typeof b) return false;
	if (Array.isArray(a) && Array.isArray(b)) {
		if (a.length !== b.length) return false;
		for (let i = 0; i < a.length; i++) {
			if (!deepEqual(a[i], b[i])) return false;
		}
		return true;
	}
	if (typeof a === 'object' && typeof b === 'object') {
		const ak = Object.keys(a as object);
		const bk = Object.keys(b as object);
		if (ak.length !== bk.length) return false;
		for (const k of ak) {
			if (!deepEqual((a as Record<string, unknown>)[k], (b as Record<string, unknown>)[k])) {
				return false;
			}
		}
		return true;
	}
	return false;
}
