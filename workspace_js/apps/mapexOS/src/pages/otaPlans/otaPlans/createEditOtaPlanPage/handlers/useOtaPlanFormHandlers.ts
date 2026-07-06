import type { Ref } from 'vue';
import type { GroupedStep } from '@components/steppers';
import type { OtaPlanFormData } from '../interfaces/createEditOtaPlan.interface';

import { computed } from 'vue';

import { apis } from '@services/mapex';
import { notifySuccess } from '@utils/alert/notify';
import { handleApiError } from '@utils/error';
import { useCreateEditOtaPlanTranslations } from '@composables/i18n/pages/otaPlans/createEditOtaPlan/useCreateEditOtaPlanTranslations';
import { useLogger } from '@composables/useLogger';

import { STEP, TOTAL_STEPS } from '../constants/createEditOtaPlan.constant';

const logger = useLogger('useOtaPlanFormHandlers');

interface UseOtaPlanFormHandlersParams {
	formData: Ref<OtaPlanFormData>;
	currentStep: Ref<number>;
	isSaving: Ref<boolean>;
	scheduleValid: Ref<boolean>;
}

/**
 * Wizard orchestration: grouped step metadata for the vertical stepper,
 * per-step validation gating, guarded navigation, and the create-plan submit.
 */
export function useOtaPlanFormHandlers(params: UseOtaPlanFormHandlersParams) {
	const { formData, currentStep, isSaving, scheduleValid } = params;

	const t = useCreateEditOtaPlanTranslations();

	/**
	 * Flat, ordered leaf steps drive navigation and the form-card header;
	 * `group` collapses them into the stepper tree. Convention: Setup first,
	 * Finalization last. The source template + devices come before the target
	 * template + firmware.
	 */
	const translatedSteps = computed<GroupedStep[]>(() => {
		const setup = { id: 'setup', label: t.groups.setup.value, icon: 'mdi-cog-outline' };
		const rollout = { id: 'rollout', label: t.groups.rollout.value, icon: 'mdi-rocket-launch-outline' };
		const finalization = { id: 'finalization', label: t.groups.finalization.value, icon: 'mdi-flag-checkered' };
		return [
			{ title: t.steps.details.label.value, icon: 'mdi-information', description: t.steps.details.description.value, group: setup },
			{ title: t.steps.schedule.label.value, icon: 'mdi-calendar-clock', description: t.steps.schedule.description.value, group: setup },
			{ title: t.steps.fromTemplate.label.value, icon: 'mdi-database-arrow-right-outline', description: t.steps.fromTemplate.description.value, group: rollout },
			{ title: t.steps.devices.label.value, icon: 'mdi-devices', description: t.steps.devices.description.value, group: rollout },
			{ title: t.steps.toTemplate.label.value, icon: 'mdi-database-arrow-left-outline', description: t.steps.toTemplate.description.value, group: rollout },
			{ title: t.steps.firmware.label.value, icon: 'mdi-memory', description: t.steps.firmware.description.value, group: rollout },
			{ title: t.steps.review.label.value, icon: 'mdi-check-all', description: t.steps.review.description.value, group: finalization },
		];
	});

	/** The details step needs a plan name */
	const isDetailsStepValid = computed(() => !!formData.value.name.trim());

	/** The schedule step owns its validity (run-now vs future start + deadline) */
	const isScheduleStepValid = computed(() => scheduleValid.value && !!formData.value.maxTime);

	/** The source-template step needs a template chosen */
	const isFromTemplateStepValid = computed(() => !!formData.value.sourceTemplateId);

	/** The devices step needs at least one device selected */
	const isDevicesStepValid = computed(() => formData.value.selectedAssetIds.length > 0);

	/** The target-template step needs a template chosen (distinct from the source) */
	const isToTemplateStepValid = computed(
		() => !!formData.value.targetTemplateId && formData.value.targetTemplateId !== formData.value.sourceTemplateId,
	);

	/** The firmware step is complete only when the artifact is confirmed READY */
	const isFirmwareStepValid = computed(() => formData.value.firmwareReady && !!formData.value.firmwareId);

	/** Whether a given step's requirements are met */
	function isStepValid(step: number): boolean {
		if (step === STEP.DETAILS) return isDetailsStepValid.value;
		if (step === STEP.SCHEDULE) return isScheduleStepValid.value;
		if (step === STEP.FROM_TEMPLATE) return isFromTemplateStepValid.value;
		if (step === STEP.DEVICES) return isDevicesStepValid.value;
		if (step === STEP.TO_TEMPLATE) return isToTemplateStepValid.value;
		if (step === STEP.FIRMWARE) return isFirmwareStepValid.value;
		return true;
	}

	/** Next is disabled while the current step is incomplete */
	const isNextButtonDisabled = computed(() => !isStepValid(currentStep.value));

	/**
	 * Guarded navigation: moving forward requires every step before the target
	 * to be valid; moving backward is always allowed.
	 */
	function changeStep(step: number): void {
		if (step < 1 || step > TOTAL_STEPS) return;
		if (step > currentStep.value) {
			for (let s = currentStep.value; s < step; s++) {
				if (!isStepValid(s)) return;
			}
		}
		currentStep.value = step;
	}

	/**
	 * Create the plan from the form state. A null startAt is omitted — the
	 * platform runs the plan immediately. The target template travels on the
	 * firmware, so only the firmware id is sent.
	 */
	async function submitForm(): Promise<boolean> {
		const f = formData.value;
		if (!f.firmwareId || !f.sourceTemplateId || !f.maxTime) return false;

		try {
			isSaving.value = true;
			await apis.assets.otaPlans.createPlan({
				firmwareId: f.firmwareId,
				sourceTemplateId: f.sourceTemplateId,
				assetIds: f.selectedAssetIds,
				...(f.startAt ? { startAt: f.startAt } : {}),
				maxTime: f.maxTime,
				name: f.name,
				...(f.description ? { description: f.description } : {}),
				rolloutConfig: { ...f.rolloutConfig },
			});
			notifySuccess({ message: t.messages.createSuccess.value });
			return true;
		} catch (err: unknown) {
			logger.error('Error creating OTA plan:', err);
			handleApiError(err, { defaultMessage: t.messages.createError.value });
			return false;
		} finally {
			isSaving.value = false;
		}
	}

	return {
		translatedSteps,
		isDetailsStepValid,
		isScheduleStepValid,
		isFromTemplateStepValid,
		isDevicesStepValid,
		isToTemplateStepValid,
		isFirmwareStepValid,
		isNextButtonDisabled,
		isStepValid,
		changeStep,
		submitForm,
	};
}
