import { computed } from 'vue';
import { useTS } from '@utils/translation';

/**
 * OTA Plan detail page translations.
 *
 * Structure mirrors:
 * - File: src/pages/otaPlans/otaPlans/otaPlanDetailPage/OtaPlanDetailPage.vue
 * - JSON: src/i18n/{locale}/pages/otaPlans/otaPlanDetail.json
 */
export function useOtaPlanDetailTranslations() {
	const ts = useTS({ capitalize: true });
	const tsTitle = useTS({ titleCase: true });

	const base = 'pages.otaPlans.otaPlanDetail';

	return {
		/** Page-level strings */
		page: {
			back: computed(() => ts(`${base}.page.back`)),
			notFound: computed(() => ts(`${base}.page.notFound`)),
			loadError: computed(() => ts(`${base}.page.loadError`)),
			apiNotInitialized: computed(() => ts(`${base}.page.apiNotInitialized`)),
		},

		/** Left-panel sections */
		sections: {
			overview: computed(() => tsTitle(`${base}.sections.overview`)),
			overviewSubtitle: computed(() => ts(`${base}.sections.overviewSubtitle`)),
			status: computed(() => tsTitle(`${base}.sections.status`)),
			templates: computed(() => tsTitle(`${base}.sections.templates`)),
			schedule: computed(() => tsTitle(`${base}.sections.schedule`)),
		},

		/** Left-panel field labels */
		fields: {
			status: computed(() => ts(`${base}.fields.status`)),
			sourceTemplate: computed(() => ts(`${base}.fields.sourceTemplate`)),
			targetTemplate: computed(() => ts(`${base}.fields.targetTemplate`)),
			schedule: computed(() => ts(`${base}.fields.schedule`)),
			created: computed(() => ts(`${base}.fields.created`)),
			deadline: computed(() => ts(`${base}.fields.deadline`)),
		},

		/** Firmware sub-panel */
		firmware: {
			title: computed(() => tsTitle(`${base}.firmware.title`)),
			version: computed(() => ts(`${base}.firmware.version`)),
			filename: computed(() => ts(`${base}.firmware.filename`)),
			size: computed(() => ts(`${base}.firmware.size`)),
			checksum: computed(() => ts(`${base}.firmware.checksum`)),
			artifactStatus: computed(() => ts(`${base}.firmware.artifactStatus`)),
			empty: computed(() => ts(`${base}.firmware.empty`)),
			download: computed(() => ts(`${base}.firmware.download`)),
			notFound: computed(() => ts(`${base}.firmware.notFound`)),
			downloadError: computed(() => ts(`${base}.firmware.downloadError`)),
		},

		/** Progress footer */
		progress: {
			title: computed(() => ts(`${base}.progress.title`)),
			total: computed(() => ts(`${base}.progress.total`)),
			succeeded: computed(() => ts(`${base}.progress.succeeded`)),
			failed: computed(() => ts(`${base}.progress.failed`)),
			timedOut: computed(() => ts(`${base}.progress.timedOut`)),
		},

		/** Executions panel */
		executions: {
			title: computed(() => tsTitle(`${base}.executions.title`)),
			itemLabel: computed(() => ts(`${base}.executions.itemLabel`)),
			itemLabelPlural: computed(() => ts(`${base}.executions.itemLabelPlural`)),
			filterStatus: computed(() => ts(`${base}.executions.filterStatus`)),
			allStatus: computed(() => ts(`${base}.executions.allStatus`)),
			empty: computed(() => ts(`${base}.executions.empty`)),
			emptyTitle: computed(() => ts(`${base}.executions.emptyTitle`)),
			resetFilter: computed(() => ts(`${base}.executions.resetFilter`)),
			columns: {
				device: computed(() => ts(`${base}.executions.columns.device`)),
				state: computed(() => ts(`${base}.executions.columns.state`)),
				progress: computed(() => ts(`${base}.executions.columns.progress`)),
				attempts: computed(() => ts(`${base}.executions.columns.attempts`)),
				error: computed(() => ts(`${base}.executions.columns.error`)),
				updated: computed(() => ts(`${base}.executions.columns.updated`)),
			},
		},

		/** Header helpers */
		header: {
			immediate: computed(() => ts(`${base}.header.immediate`)),
		},

		/** Actions */
		actions: {
			cancel: computed(() => ts(`${base}.actions.cancel`)),
			refresh: computed(() => ts(`${base}.actions.refresh`)),
		},

		/** Cancel confirmation dialog */
		cancelDialog: {
			title: computed(() => ts(`${base}.cancelDialog.title`)),
			message: computed(() => ts(`${base}.cancelDialog.message`)),
			confirm: computed(() => ts(`${base}.cancelDialog.confirm`)),
			success: computed(() => ts(`${base}.cancelDialog.success`)),
			error: computed(() => ts(`${base}.cancelDialog.error`)),
		},

		/** Plan status display name, keyed by the wire value */
		planStatus: (status: string) => ts(`${base}.statuses.${status}`),

		/** Execution state display name, keyed by the wire value */
		executionState: (state: string) => ts(`${base}.states.${state}`),
	};
}
