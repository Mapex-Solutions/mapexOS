import { computed } from 'vue';
import { useTS } from '@utils/translation';

/**
 * OTA Plans list page translations.
 *
 * Structure mirrors:
 * - File: src/pages/otaPlans/otaPlans/otaPlanListPage/OtaPlanListPage.vue
 * - JSON: src/i18n/{locale}/pages/otaPlans/otaPlanList.json
 *
 * Also carries the plan-detail drawer strings (the drawer opens from this
 * list, so its vocabulary lives in the same JSON).
 */
export function useOtaPlanListTranslations() {
	const ts = useTS({ capitalize: true });
	const tsTitle = useTS({ titleCase: true });
	const tsRaw = useTS({ capitalize: false });

	const base = 'pages.otaPlans.otaPlanList';

	return {
		/** Page header */
		page: {
			title: computed(() => tsTitle(`${base}.title`)),
			description: computed(() => ts(`${base}.description`)),
			addButton: computed(() => ts(`${base}.addButton`)),
			listTitle: computed(() => ts(`${base}.listTitle`)),
			itemLabel: computed(() => tsRaw(`${base}.itemLabel`)),
			itemLabelPlural: computed(() => tsRaw(`${base}.itemLabelPlural`)),
		},

		/** Filter labels */
		filters: {
			label: computed(() => ts(`${base}.filters.label`)),
			searchPlaceholder: computed(() => ts(`${base}.filters.searchPlaceholder`)),
			status: computed(() => ts(`${base}.filters.status`)),
			clearAll: computed(() => ts(`${base}.filters.clearAll`)),
		},

		/** Table column labels */
		columns: {
			name: computed(() => ts(`${base}.columns.name`)),
			status: computed(() => ts(`${base}.columns.status`)),
			progress: computed(() => ts(`${base}.columns.progress`)),
			startAt: computed(() => ts(`${base}.columns.startAt`)),
			maxTime: computed(() => ts(`${base}.columns.maxTime`)),
			updated: computed(() => ts(`${base}.columns.updated`)),
		},

		/** Plan status display names, keyed by the wire status value */
		planStatus: (status: string) => ts(`${base}.statuses.${status}`),

		/** Empty state */
		empty: {
			title: computed(() => ts(`${base}.empty.title`)),
			description: computed(() => ts(`${base}.empty.description`)),
			action: computed(() => ts(`${base}.empty.action`)),
		},

		/** Row actions */
		actions: {
			view: computed(() => ts(`${base}.actions.view`)),
			cancel: computed(() => ts(`${base}.actions.cancel`)),
		},

		/** Cancel confirmation dialog */
		cancelDialog: {
			title: computed(() => ts(`${base}.cancelDialog.title`)),
			message: computed(() => ts(`${base}.cancelDialog.message`)),
			confirm: computed(() => ts(`${base}.cancelDialog.confirm`)),
			keep: computed(() => ts(`${base}.cancelDialog.keep`)),
		},

		/** Toast messages */
		messages: {
			loadError: computed(() => ts(`${base}.messages.loadError`)),
			cancelSuccess: computed(() => ts(`${base}.messages.cancelSuccess`)),
			cancelError: computed(() => ts(`${base}.messages.cancelError`)),
		},

		/** Plan detail drawer */
		drawer: {
			title: computed(() => tsTitle(`${base}.drawer.title`)),
			info: {
				firmware: computed(() => ts(`${base}.drawer.info.firmware`)),
				window: computed(() => ts(`${base}.drawer.info.window`)),
				counters: computed(() => ts(`${base}.drawer.info.counters`)),
				source: computed(() => ts(`${base}.drawer.info.source`)),
				target: computed(() => ts(`${base}.drawer.info.target`)),
				rate: computed(() => ts(`${base}.drawer.info.rate`)),
				abort: computed(() => ts(`${base}.drawer.info.abort`)),
			},
			executions: {
				title: computed(() => ts(`${base}.drawer.executions.title`)),
				columns: {
					device: computed(() => ts(`${base}.drawer.executions.columns.device`)),
					state: computed(() => ts(`${base}.drawer.executions.columns.state`)),
					progress: computed(() => ts(`${base}.drawer.executions.columns.progress`)),
					attempts: computed(() => ts(`${base}.drawer.executions.columns.attempts`)),
					updated: computed(() => ts(`${base}.drawer.executions.columns.updated`)),
				},
				empty: computed(() => ts(`${base}.drawer.executions.empty`)),
				refresh: computed(() => ts(`${base}.drawer.executions.refresh`)),
			},
			messages: {
				loadError: computed(() => ts(`${base}.drawer.messages.loadError`)),
			},
		},

		/** Execution state display names, keyed by the wire state value */
		executionState: (state: string) => ts(`${base}.drawer.executions.states.${state}`),
	};
}
