import { computed } from 'vue';
import { useTS } from '@utils/translation';

/**
 * Create OTA Plan stepper translations.
 *
 * Structure mirrors:
 * - File: src/pages/otaPlans/otaPlans/createEditOtaPlanPage/CreateEditOtaPlanPage.vue
 * - JSON: src/i18n/{locale}/pages/otaPlans/createEditOtaPlan.json
 *
 * The wizard is GROUPED (Setup / Rollout / Finalization); step labels and
 * descriptions feed the stepper tree, while each step section carries its own
 * field vocabulary.
 */
export function useCreateEditOtaPlanTranslations() {
	const ts = useTS({ capitalize: true });
	const tsTitle = useTS({ titleCase: true });

	const base = 'pages.otaPlans.createEditOtaPlan';

	return {
		/** Page header */
		page: {
			title: computed(() => tsTitle(`${base}.title`)),
			description: computed(() => ts(`${base}.description`)),
			button: computed(() => ts(`${base}.button`)),
		},

		/** Stepper group labels (Setup / Rollout / Finalization) */
		groups: {
			setup: computed(() => ts(`${base}.groups.setup`)),
			rollout: computed(() => ts(`${base}.groups.rollout`)),
			finalization: computed(() => ts(`${base}.groups.finalization`)),
		},

		/** Stepper leaves: label + description per step */
		steps: {
			details: {
				label: computed(() => ts(`${base}.steps.details.label`)),
				description: computed(() => ts(`${base}.steps.details.description`)),
				title: computed(() => ts(`${base}.steps.details.title`)),
				subtitle: computed(() => ts(`${base}.steps.details.subtitle`)),
				fields: {
					name: {
						label: computed(() => ts(`${base}.steps.details.fields.name.label`)),
						placeholder: computed(() => ts(`${base}.steps.details.fields.name.placeholder`)),
						hint: computed(() => ts(`${base}.steps.details.fields.name.hint`)),
						required: computed(() => ts(`${base}.steps.details.fields.name.required`)),
					},
					description: {
						label: computed(() => ts(`${base}.steps.details.fields.description.label`)),
						placeholder: computed(() => ts(`${base}.steps.details.fields.description.placeholder`)),
						hint: computed(() => ts(`${base}.steps.details.fields.description.hint`)),
					},
				},
			},
			schedule: {
				label: computed(() => ts(`${base}.steps.schedule.label`)),
				description: computed(() => ts(`${base}.steps.schedule.description`)),
				title: computed(() => ts(`${base}.steps.schedule.title`)),
				subtitle: computed(() => ts(`${base}.steps.schedule.subtitle`)),
				modeLabel: computed(() => ts(`${base}.steps.schedule.modeLabel`)),
				options: {
					now: computed(() => ts(`${base}.steps.schedule.options.now`)),
					scheduled: computed(() => ts(`${base}.steps.schedule.options.scheduled`)),
				},
				field: {
					label: computed(() => ts(`${base}.steps.schedule.field.label`)),
					placeholder: computed(() => ts(`${base}.steps.schedule.field.placeholder`)),
					hint: computed(() => ts(`${base}.steps.schedule.field.hint`)),
					required: computed(() => ts(`${base}.steps.schedule.field.required`)),
					future: computed(() => ts(`${base}.steps.schedule.field.future`)),
					close: computed(() => ts(`${base}.steps.schedule.field.close`)),
				},
				maxTime: {
					label: computed(() => ts(`${base}.steps.schedule.maxTime.label`)),
					placeholder: computed(() => ts(`${base}.steps.schedule.maxTime.placeholder`)),
					hint: computed(() => ts(`${base}.steps.schedule.maxTime.hint`)),
					required: computed(() => ts(`${base}.steps.schedule.maxTime.required`)),
					afterStart: computed(() => ts(`${base}.steps.schedule.maxTime.afterStart`)),
				},
				pacing: {
					title: computed(() => ts(`${base}.steps.schedule.pacing.title`)),
					ratePerMinute: computed(() => ts(`${base}.steps.schedule.pacing.ratePerMinute`)),
					abortThresholdPct: computed(() => ts(`${base}.steps.schedule.pacing.abortThresholdPct`)),
					abortThresholdHint: computed(() => ts(`${base}.steps.schedule.pacing.abortThresholdHint`)),
					abortMinExecuted: computed(() => ts(`${base}.steps.schedule.pacing.abortMinExecuted`)),
				},
			},
			fromTemplate: {
				label: computed(() => ts(`${base}.steps.fromTemplate.label`)),
				description: computed(() => ts(`${base}.steps.fromTemplate.description`)),
				title: computed(() => ts(`${base}.steps.fromTemplate.title`)),
				subtitle: computed(() => ts(`${base}.steps.fromTemplate.subtitle`)),
			},
			devices: {
				label: computed(() => ts(`${base}.steps.devices.label`)),
				description: computed(() => ts(`${base}.steps.devices.description`)),
				title: computed(() => ts(`${base}.steps.devices.title`)),
				subtitle: computed(() => ts(`${base}.steps.devices.subtitle`)),
			},
			toTemplate: {
				label: computed(() => ts(`${base}.steps.toTemplate.label`)),
				description: computed(() => ts(`${base}.steps.toTemplate.description`)),
				title: computed(() => ts(`${base}.steps.toTemplate.title`)),
				subtitle: computed(() => ts(`${base}.steps.toTemplate.subtitle`)),
			},
			firmware: {
				label: computed(() => ts(`${base}.steps.firmware.label`)),
				description: computed(() => ts(`${base}.steps.firmware.description`)),
			},
			review: {
				label: computed(() => ts(`${base}.steps.review.label`)),
				description: computed(() => ts(`${base}.steps.review.description`)),
			},
		},

		/** Firmware step fields */
		firmware: {
			version: computed(() => ts(`${base}.firmware.version`)),
			versionHint: computed(() => ts(`${base}.firmware.versionHint`)),
			dropzone: {
				title: computed(() => ts(`${base}.firmware.dropzone.title`)),
				hint: computed(() => ts(`${base}.firmware.dropzone.hint`)),
			},
			fileFacts: {
				name: computed(() => ts(`${base}.firmware.fileFacts.name`)),
				size: computed(() => ts(`${base}.firmware.fileFacts.size`)),
				checksum: computed(() => ts(`${base}.firmware.fileFacts.checksum`)),
			},
			size: computed(() => ts(`${base}.firmware.size`)),
			checksum: computed(() => ts(`${base}.firmware.checksum`)),
			upload: computed(() => ts(`${base}.firmware.upload`)),
			computing: computed(() => ts(`${base}.firmware.computing`)),
			uploading: computed(() => ts(`${base}.firmware.uploading`)),
			confirming: computed(() => ts(`${base}.firmware.confirming`)),
			ready: computed(() => ts(`${base}.firmware.ready`)),
			uploadSuccess: computed(() => ts(`${base}.firmware.uploadSuccess`)),
			uploadError: computed(() => ts(`${base}.firmware.uploadError`)),
		},

		/** Devices step selection labels */
		devices: {
			selectAll: computed(() => ts(`${base}.devices.selectAll`)),
			selected: computed(() => ts(`${base}.devices.selected`)),
			noAssets: computed(() => ts(`${base}.devices.noAssets`)),
		},

		/** Review step */
		review: {
			title: computed(() => ts(`${base}.review.title`)),
			firmware: computed(() => ts(`${base}.review.firmware`)),
			devices: computed(() => ts(`${base}.review.devices`)),
			window: computed(() => ts(`${base}.review.window`)),
			runNow: computed(() => ts(`${base}.review.runNow`)),
			pacing: computed(() => ts(`${base}.review.pacing`)),
			deviceCount: computed(() => ts(`${base}.review.deviceCount`)),
		},

		/** Wizard buttons */
		buttons: {
			next: computed(() => ts(`${base}.buttons.next`)),
			back: computed(() => ts(`${base}.buttons.back`)),
			submit: computed(() => ts(`${base}.buttons.submit`)),
		},

		/** Per-step validation messages */
		validation: {
			required: computed(() => ts(`${base}.validation.required`)),
			fileRequired: computed(() => ts(`${base}.validation.fileRequired`)),
			assetsRequired: computed(() => ts(`${base}.validation.assetsRequired`)),
		},

		/** Toast messages */
		messages: {
			createSuccess: computed(() => ts(`${base}.messages.createSuccess`)),
			createError: computed(() => ts(`${base}.messages.createError`)),
		},
	};
}
