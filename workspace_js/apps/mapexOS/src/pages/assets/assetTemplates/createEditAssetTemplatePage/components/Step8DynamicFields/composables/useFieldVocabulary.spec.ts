/** VUE IMPORTS */
import { defineComponent, createApp } from 'vue';

import { describe, it, expect, beforeEach, vi } from 'vitest';

/**
 * Mock the field-vocabulary endpoint so each test controls success vs failure.
 */
const fieldVocabulary = vi.fn();

vi.mock('@services/mapex', () => ({
	apis: {
		assets: {
			assetTemplate: {
				fieldVocabulary,
			},
		},
	},
}));

/** Lazy import — resolved after vi.mock and module reset (clears the cache). */
// eslint-disable-next-line @typescript-eslint/consistent-type-imports
let useFieldVocabulary: typeof import('./useFieldVocabulary').useFieldVocabulary;

/**
 * Invoke a composable inside a real Vue setup context so injection-based
 * dependencies (like the i18n locale) resolve.
 *
 * @param {() => T} composable - Composable factory.
 * @returns {{ result: T; app: ReturnType<typeof createApp> }}
 */
function withSetup<T>(composable: () => T): { result: T; app: ReturnType<typeof createApp> } {
	let result!: T;
	const app = createApp(
		defineComponent({
			setup() {
				result = composable();
				return () => null;
			},
		}),
	);
	app.mount(document.createElement('div'));
	return { result, app };
}

/**
 * Flush pending microtasks so the composable's fire-and-forget fetch settles.
 *
 * @returns {Promise<void>}
 */
async function flush(): Promise<void> {
	await Promise.resolve();
	await Promise.resolve();
}

beforeEach(async () => {
	fieldVocabulary.mockReset();
	// Reset module graph to clear the per-language cache between tests.
	vi.resetModules();
	const mod = await import('./useFieldVocabulary');
	useFieldVocabulary = mod.useFieldVocabulary;
});

describe('useFieldVocabulary', () => {
	it('falls back to empty + unavailable when the endpoint rejects', async () => {
		fieldVocabulary.mockRejectedValueOnce(new Error('network down'));

		const { result } = withSetup(() => useFieldVocabulary());
		await flush();

		expect(result.groups.value).toEqual([]);
		expect(result.isAvailable.value).toBe(false);
		expect(result.loading.value).toBe(false);
	});

	it('exposes the fetched groups and marks the vocabulary available', async () => {
		fieldVocabulary.mockResolvedValueOnce({
			groups: [
				{
					category: 'climate',
					label: 'Climate',
					fields: [{ value: 'temperature', hint: 'Ambient temperature', type: 'number', unit: '°C' }],
				},
			],
		});

		const { result } = withSetup(() => useFieldVocabulary());
		await flush();

		expect(result.groups.value).toHaveLength(1);
		expect(result.groups.value[0]?.fields[0]?.value).toBe('temperature');
		expect(result.isAvailable.value).toBe(true);
	});
});
