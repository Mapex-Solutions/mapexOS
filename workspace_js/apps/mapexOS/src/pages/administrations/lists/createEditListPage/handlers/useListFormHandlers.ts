import type { Ref } from 'vue';
import type { QForm } from 'quasar';
import type { ListFormData, ListType } from '../interfaces';

import { computed } from 'vue';
import { useRouter } from 'vue-router';

import { notifySuccess, notifyFail } from '@utils/alert/notify';
import { useLogger } from '@composables/useLogger';
import { apis } from '@services/mapex';
import { PARENT_TYPE_FOR } from '../constants';

const logger = useLogger('useListFormHandlers');

/**
 * Step component ref with validate method
 */
interface StepComponentRef {
	formRef?: QForm | null;
	validate: () => boolean | Promise<boolean>;
}

/**
 * Parameters for useListFormHandlers composable
 */
export interface UseListFormHandlersParams {
	listData: Ref<ListFormData>;
	currentStep: Ref<number>;
	isEditMode: Ref<boolean>;
	listId: Ref<string | undefined>;
	isSaving: Ref<boolean>;
	step1Ref: Ref<StepComponentRef | null>;
	t: any;
}

/**
 * Composable that owns step navigation, validation and form submission for
 * the list create/edit page. Mirrors useRoleFormHandlers / useCustomerFormHandlers.
 *
 * @param params - composable parameters
 * @returns form handler functions and reactive flags
 */
export function useListFormHandlers(params: UseListFormHandlersParams) {
	const {
		listData,
		currentStep,
		isEditMode,
		listId,
		isSaving,
		step1Ref,
		t,
	} = params;

	const router = useRouter();

	/**
	 * Whether the chosen type requires a parent selection
	 */
	const parentRequired = computed((): boolean => {
		const type = listData.value.type;
		if (!type) return false;
		return PARENT_TYPE_FOR[type] !== null;
	});

	/**
	 * Whether the Next button should be disabled on the current step
	 */
	const isNextButtonDisabled = computed((): boolean => {
		if (currentStep.value === 1) {
			const { type, parentId, name, value } = listData.value;
			if (!type) return true;
			if (parentRequired.value && !parentId) return true;
			if (!name?.trim()) return true;
			if (!value?.trim()) return true;
		}
		return false;
	});

	/**
	 * Validate the current step (delegated to the step's form ref) and move
	 * forward only if it passes; moving back is always allowed.
	 *
	 * @param step - target step number
	 */
	async function changeStep(step: number): Promise<void> {
		if (step > currentStep.value && currentStep.value === 1 && step1Ref.value) {
			const valid = await step1Ref.value.validate();
			if (!valid) return;
		}
		currentStep.value = step;
	}

	/**
	 * Sync wrapper for use as a Vue event handler
	 *
	 * @param step - target step number
	 */
	function handleStepChange(step: number): void {
		void changeStep(step);
	}

	/**
	 * Build the POST /api/v1/lists payload from form state.
	 * `isSystem` is hard-coded to false: the UI never lets users create system items.
	 *
	 * @returns request body for `apis.mapexOS.lists.create`
	 */
	function buildCreatePayload() {
		const { type, parentId, name, value, enabled, isTemplate } = listData.value;
		return {
			type: type as ListType,
			name: name.trim(),
			value: value.trim(),
			enabled,
			isSystem: false,
			isTemplate,
			...(parentId ? { parentId } : {}),
		};
	}

	/**
	 * Build the PATCH /api/v1/lists/:id payload. The backend only accepts
	 * name/value/enabled/parentId/metadata on update, so type is never sent.
	 *
	 * @returns request body for `apis.mapexOS.lists.update`
	 */
	function buildUpdatePayload() {
		const { parentId, name, value, enabled } = listData.value;
		return {
			name: name.trim(),
			value: value.trim(),
			enabled,
			...(parentId ? { parentId } : {}),
		};
	}

	/**
	 * Submit the form. Routes to create or update based on `isEditMode`,
	 * then navigates back to the lists list on success.
	 */
	async function submitForm(): Promise<void> {
		if (!apis.mapexOS?.lists) {
			notifyFail({ message: t.errors.apiNotInitialized.value });
			return;
		}

		isSaving.value = true;
		try {
			if (isEditMode.value && listId.value) {
				const payload = buildUpdatePayload();
				logger.debug('Updating List:', { listId: listId.value, payload });
				await apis.mapexOS.lists.update({ listId: listId.value }, payload);
				notifySuccess({ message: t.createEditNotifications.updated.value, timeout: 3000 });
			} else {
				const payload = buildCreatePayload();
				logger.debug('Creating List:', payload);
				await apis.mapexOS.lists.create(payload);
				notifySuccess({ message: t.createEditNotifications.created.value, timeout: 3000 });
			}
			await router.push('/admin/lists');
		} catch (error: any) {
			logger.error('List form submission error:', error);
			const status = error?.response?.status || error?.code;
			let message = isEditMode.value
				? t.createEditNotifications.updateFailed.value
				: t.createEditNotifications.createFailed.value;
			if (status === 409) message = t.createEditNotifications.alreadyExists.value;
			else if (status === 403) message = t.createEditNotifications.forbidden.value;
			notifyFail({ message, timeout: 5000 });
		} finally {
			isSaving.value = false;
		}
	}

	return {
		parentRequired,
		isNextButtonDisabled,
		changeStep,
		handleStepChange,
		submitForm,
	};
}
