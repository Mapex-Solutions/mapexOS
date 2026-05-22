import { computed } from 'vue';
import { useTS } from '@utils/translation';

/**
 * i18n composable for the Retention Policies admin page.
 *
 * Structure mirrors:
 * - File: src/pages/administrations/retentionPoliciesPage/RetentionPoliciesPage.vue
 * - JSON: src/i18n/{locale}/pages/administrations/retentionPolicies.json
 */
export function useRetentionPoliciesTranslations() {
	const ts = useTS({ capitalize: true });
	const tsTitle = useTS({ titleCase: true });

	return {
		page: {
			title: computed(() => tsTitle('pages.administrations.retentionPolicies.page.title')),
			description: computed(() => ts('pages.administrations.retentionPolicies.page.description')),
		},
		columns: {
			type: computed(() => ts('pages.administrations.retentionPolicies.columns.type')),
			name: computed(() => ts('pages.administrations.retentionPolicies.columns.name')),
			retentionDays: computed(() =>
				ts('pages.administrations.retentionPolicies.columns.retentionDays'),
			),
			enabled: computed(() => ts('pages.administrations.retentionPolicies.columns.enabled')),
		},
		actions: {
			save: computed(() => ts('pages.administrations.retentionPolicies.actions.save')),
			cancel: computed(() => ts('pages.administrations.retentionPolicies.actions.cancel')),
		},
		messages: {
			saveSuccess: computed(() =>
				ts('pages.administrations.retentionPolicies.messages.saveSuccess'),
			),
			saveError: computed(() =>
				ts('pages.administrations.retentionPolicies.messages.saveError'),
			),
		},
		validation: {
			ttlRange: (min: number, max: number) =>
				ts('pages.administrations.retentionPolicies.validation.ttlRange', { min, max }),
		},
		specialTypes: {
			assetStatusHistory: computed(() =>
				ts('pages.administrations.retentionPolicies.specialTypes.asset_status_history'),
			),
		},
		empty: {
			title: computed(() => ts('pages.administrations.retentionPolicies.empty.title')),
			description: computed(() =>
				ts('pages.administrations.retentionPolicies.empty.description'),
			),
		},
	};
}
