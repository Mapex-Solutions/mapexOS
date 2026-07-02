import type { Ref } from 'vue';
import type { FieldVocabularyGroup } from '@mapexos/schemas';

/**
 * Reactive surface returned by the field-vocabulary composable.
 *
 * The control consuming this degrades gracefully: when `isAvailable` is false
 * the suggestion list is empty and the field name falls back to free text.
 */
export interface UseFieldVocabularyReturn {
	/**
	 * Canonical-field suggestion groups, in the order received from the backend.
	 * Empty when the endpoint is unreachable or returned no groups.
	 */
	groups: Ref<FieldVocabularyGroup[]>;

	/**
	 * Whether a fetch is currently in flight.
	 */
	loading: Ref<boolean>;

	/**
	 * Whether the vocabulary was fetched successfully and has at least one group.
	 * False on any network, parse, or empty-result condition.
	 */
	isAvailable: Ref<boolean>;
}
