import { computed } from 'vue';
import { useTS } from '@utils/translation';

/**
 * Group detail page translations
 *
 * Structure mirrors:
 * - File: src/pages/administrations/groups/groupDetailPage/GroupDetailPage.vue
 * - JSON: src/i18n/{locale}/pages/administrations/groups.json
 * - Composable: src/composables/i18n/pages/administrations/groups/useGroupDetailTranslations.ts
 */
export function useGroupDetailTranslations() {
	const ts = useTS({ capitalize: true });
	const tsTitle = useTS({ titleCase: true });

	return {
		/**
		 * Page header translations
		 */
		page: {
			title: computed(() => tsTitle('pages.administrations.groups.detail.title')),
			description: computed(() => ts('pages.administrations.groups.detail.description')),
			loading: computed(() => ts('pages.administrations.groups.detail.loading')),
			loadError: computed(() => ts('pages.administrations.groups.detail.loadError')),
			backButton: computed(() => ts('pages.administrations.groups.detail.backButton')),
			editButton: computed(() => ts('pages.administrations.groups.detail.editButton')),
		},

		/**
		 * Tabs translations
		 */
		tabs: {
			info: computed(() => ts('pages.administrations.groups.detail.tabs.info')),
			roles: computed(() => ts('pages.administrations.groups.detail.tabs.roles')),
			members: computed(() => ts('pages.administrations.groups.detail.tabs.members')),
		},

		/**
		 * Info tab translations
		 */
		info: {
			basicInfo: {
				title: computed(() => ts('pages.administrations.groups.detail.info.basicInfo.title')),
			},
			dates: {
				title: computed(() => ts('pages.administrations.groups.detail.info.dates.title')),
			},
			fields: {
				name: computed(() => ts('pages.administrations.groups.detail.info.fields.name')),
				description: computed(() => ts('pages.administrations.groups.detail.info.fields.description')),
				status: computed(() => ts('pages.administrations.groups.detail.info.fields.status')),
				membersCount: computed(() => ts('pages.administrations.groups.detail.info.fields.membersCount')),
				organization: computed(() => ts('pages.administrations.groups.detail.info.fields.organization')),
				created: computed(() => ts('pages.administrations.groups.detail.info.fields.created')),
				updated: computed(() => ts('pages.administrations.groups.detail.info.fields.updated')),
			},
			status: {
				enabled: computed(() => ts('pages.administrations.groups.detail.info.status.enabled')),
				disabled: computed(() => ts('pages.administrations.groups.detail.info.status.disabled')),
			},
			noData: computed(() => ts('pages.administrations.groups.detail.info.noData')),
		},

		/**
		 * Roles tab translations
		 */
		roles: {
			title: computed(() => ts('pages.administrations.groups.detail.roles.title')),
			loadError: computed(() => ts('pages.administrations.groups.detail.roles.loadError')),
			columns: {
				name: computed(() => ts('pages.administrations.groups.detail.roles.columns.name')),
				type: computed(() => ts('pages.administrations.groups.detail.roles.columns.type')),
				scope: computed(() => ts('pages.administrations.groups.detail.roles.columns.scope')),
			},
			empty: {
				title: computed(() => ts('pages.administrations.groups.detail.roles.empty.title')),
				description: computed(() => ts('pages.administrations.groups.detail.roles.empty.description')),
			},
		},

		/**
		 * Members tab translations
		 */
		members: {
			title: computed(() => ts('pages.administrations.groups.detail.members.title')),
			loadError: computed(() => ts('pages.administrations.groups.detail.members.loadError')),
			columns: {
				name: computed(() => ts('pages.administrations.groups.detail.members.columns.name')),
				email: computed(() => ts('pages.administrations.groups.detail.members.columns.email')),
				addedAt: computed(() => ts('pages.administrations.groups.detail.members.columns.addedAt')),
			},
			empty: {
				title: computed(() => ts('pages.administrations.groups.detail.members.empty.title')),
				description: computed(() => ts('pages.administrations.groups.detail.members.empty.description')),
			},
		},

		/**
		 * Error translations
		 */
		errors: {
			apiNotInitialized: computed(() => ts('pages.administrations.groups.errors.apiNotInitialized')),
			idMissing: computed(() => ts('pages.administrations.groups.errors.idMissing')),
		},
	};
}
