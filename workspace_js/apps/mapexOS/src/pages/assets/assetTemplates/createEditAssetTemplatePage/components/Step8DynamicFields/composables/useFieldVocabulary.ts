/** TYPE IMPORTS */
import type { FieldVocabularyGroup } from '@mapexos/schemas';
import type { UseFieldVocabularyReturn } from './interfaces';

/** VUE IMPORTS */
import { ref } from 'vue';
import { useI18n } from 'vue-i18n';

/** SERVICES */
import { apis } from '@services/mapex';

/** COMPOSABLES */
import { useLogger } from '@composables/useLogger';

const logger = useLogger('useFieldVocabulary');

/** STATE (module-level singleton — fetched once, shared across all callers) */

/**
 * Cached suggestion groups, keyed by the resolved language so a locale switch
 * triggers a fresh fetch rather than serving stale labels and hints.
 */
const cache = new Map<string, FieldVocabularyGroup[]>();

/**
 * Resolve the backend language code from the active i18n locale.
 *
 * The backend keys its `hint` and group `label` translations by the full
 * locale (`en-US`, `pt-BR`) and falls back to `en-US` on a miss, so the locale
 * is passed through unchanged — dropping the region (e.g. `pt`) would miss the
 * `pt-BR` key and silently fall back to English.
 *
 * @param {string} locale - Full i18n locale (e.g. `en-US`, `pt-BR`).
 * @returns {string} The same locale, used verbatim as the backend `lang`.
 */
function toLang(locale: string): string {
	return locale;
}

/**
 * Loads the canonical field-name vocabulary used to suggest dynamic-field names.
 *
 * Fetches the grouped vocabulary once per language and caches it. On any
 * failure (network, unreachable endpoint, or schema mismatch) it resolves to
 * empty groups with `isAvailable` set to false and never throws, so the
 * consuming control silently falls back to free-text custom entry.
 *
 * @returns {UseFieldVocabularyReturn} Reactive groups, loading, and availability.
 */
export function useFieldVocabulary(): UseFieldVocabularyReturn {
	const { locale } = useI18n();

	const groups = ref<FieldVocabularyGroup[]>([]);
	const loading = ref(false);
	const isAvailable = ref(false);

	/**
	 * Fetch the vocabulary for the active locale, honoring the per-language cache.
	 *
	 * @returns {Promise<void>}
	 */
	async function load(): Promise<void> {
		const lang = toLang(String(locale.value));

		const cached = cache.get(lang);
		if (cached) {
			groups.value = cached;
			isAvailable.value = cached.length > 0;
			return;
		}

		loading.value = true;

		try {
			const data = await apis.assets.assetTemplate.fieldVocabulary({ lang });
			const fetched = data?.groups ?? [];

			cache.set(lang, fetched);
			groups.value = fetched;
			isAvailable.value = fetched.length > 0;
		} catch (error) {
			logger.warn('Field vocabulary unavailable, falling back to free-text entry', error);
			groups.value = [];
			isAvailable.value = false;
		} finally {
			loading.value = false;
		}
	}

	void load();

	return {
		groups,
		loading,
		isAvailable,
	};
}
