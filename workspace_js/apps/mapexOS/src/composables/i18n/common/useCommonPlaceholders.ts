import { computed } from 'vue';
import { useTS } from '@utils/translation';

/**
 * Common placeholder translations
 * Used for search inputs and generic empty-input prompts across the app.
 *
 * @example
 * ```ts
 * const { placeholders } = useCommonPlaceholders();
 * <q-input :placeholder="placeholders.searchFields" />
 * ```
 */
export function useCommonPlaceholders() {
	const ts = useTS({ capitalize: false });

	return {
		placeholders: {
			searchFields: computed(() => ts('common.placeholders.searchFields')),
			search: computed(() => ts('common.placeholders.search')),
			typeToSearch: computed(() => ts('common.placeholders.typeToSearch')),
			searchByName: computed(() => ts('common.placeholders.searchByName')),
		},
	};
}
