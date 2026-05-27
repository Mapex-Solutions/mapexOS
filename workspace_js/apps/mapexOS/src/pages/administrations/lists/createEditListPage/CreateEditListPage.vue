<script setup lang="ts">
defineOptions({
	name: 'CreateEditListPage',
});

/** TYPE IMPORTS */
import type { FormCardHeader } from '@components/cards';
import type { ListResponse } from '@mapexos/schemas';
import type { ListType, ParentOption } from './interfaces';

/** VUE IMPORTS */
import { ref, computed, watch, onMounted } from 'vue';
import { useRouter, useRoute } from 'vue-router';

/** COMPONENTS */
import { PageHeader } from '@components/headers';
import { FormCard } from '@components/cards';
import { StepperVertical } from '@components/steppers';
import { Step1Basic, Step2Review } from './components';

/** COMPOSABLES */
import { useStepperNavigation } from '@composables/shared/form';
import { useListsTranslations } from '@composables/i18n/pages/administrations/lists/useListsTranslations';
import { useLogger } from '@composables/useLogger';

/** UTILS */
import { notifyFail } from '@utils/alert';

/** SERVICES */
import { apis } from '@services/mapex';

/** LOCAL IMPORTS */
import {
	INITIAL_LIST_FORM_DATA,
	TOTAL_STEPS,
	STEP,
	PARENT_TYPE_FOR,
	TYPE_ICON,
} from './constants';
import { useListFormHandlers } from './handlers';

/** COMPOSABLES & STORES */
const t = useListsTranslations();
const router = useRouter();
const route = useRoute();
const logger = useLogger('CreateEditListPage');

/** EDIT MODE DETECTION */
const isEditMode = ref(!!route.params.id);
const listId = ref(route.params.id as string | undefined);

/** LOADING STATES */
const isLoading = ref(false);
const isSaving = ref(false);
const loadingParents = ref(false);

/** SYSTEM LIST CHECK */
const isSystemList = ref(false);

/** STATE */
const step1Ref = ref<InstanceType<typeof Step1Basic> | null>(null);
const currentStep = ref(1);
const listData = ref({ ...INITIAL_LIST_FORM_DATA });
const parentOptions = ref<ParentOption[]>([]);

/** COMPUTED */

/**
 * Page title — switches between create and edit labels
 */
const pageTitle = computed(() =>
	isEditMode.value ? t.page.titleEdit.value : t.page.titleCreate.value,
);

/**
 * Page icon — reflects the selected type once known, falls back to "list"
 */
const pageIcon = computed(() => {
	if (isEditMode.value) return 'edit';
	const type = listData.value.type;
	return type ? TYPE_ICON[type] : 'list';
});

/**
 * Steps array passed to the StepperVertical
 */
const translatedSteps = computed(() => [
	{
		title: t.steps.basicInfo.value,
		icon: 'badge',
		description: t.stepDescriptions.basicInfo.value,
	},
	{
		title: t.steps.review.value,
		icon: 'check_circle',
		description: t.stepDescriptions.review.value,
	},
]);

const handlers = useListFormHandlers({
	listData,
	currentStep,
	isEditMode,
	listId,
	isSaving,
	step1Ref: computed(() => step1Ref.value),
	t,
});

const formNavigation = computed(() => ({
	currentStep: currentStep.value,
	totalSteps: TOTAL_STEPS,
	showPreviousButton: true,
	showNextButton: true,
	showSaveButton: true,
	showCancelButton: true,
	disableNextButton: handlers.isNextButtonDisabled.value || isSaving.value,
	disableSaveButton: handlers.isNextButtonDisabled.value || isSaving.value,
	loadingSaveButton: isSaving.value,
}));

const buttonLabels = computed(() => ({
	previous: t.buttons.back.value,
	next: t.buttons.next.value,
	save: isEditMode.value ? t.buttons.updateList.value : t.buttons.createList.value,
}));

/**
 * Resolved parent display name for the review step
 */
const parentName = computed(() => {
	if (!listData.value.parentId) return null;
	const match = parentOptions.value.find((opt) => opt.id === listData.value.parentId);
	return match?.name ?? null;
});

/** FUNCTIONS */

/**
 * Update list form data with a partial patch from a step component.
 *
 * @param partial - partial form data to merge into state
 */
function updateListData(partial: Partial<typeof listData.value>): void {
	listData.value = { ...listData.value, ...partial };
}

/**
 * Load the parents available for the currently selected list type.
 * Returns an empty list for asset_category (top of the tree) and skips the API call.
 */
async function loadParentOptions(): Promise<void> {
	const currentType = listData.value.type;
	if (!currentType) {
		parentOptions.value = [];
		return;
	}
	const parentType = PARENT_TYPE_FOR[currentType];
	if (!parentType) {
		parentOptions.value = [];
		return;
	}

	if (!apis.mapexOS?.lists) {
		notifyFail({ message: t.errors.apiNotInitialized.value });
		return;
	}

	loadingParents.value = true;
	try {
		const response = await apis.mapexOS.lists.list({
			type: parentType,
			perPage: 100,
			projection: 'name,value,type',
		});
		parentOptions.value = (response?.items ?? []).map((item: ListResponse) => ({
			id: item.id ?? '',
			name: item.name ?? '',
			value: item.value ?? '',
		}));
	} catch (error: any) {
		logger.error('Failed to load parent options:', error);
		notifyFail({ message: t.messages.parentLoadFailed.value });
		parentOptions.value = [];
	} finally {
		loadingParents.value = false;
	}
}

/**
 * Load the list item being edited and populate form state.
 * On failure, navigates back to the list page so the user is not stuck on a broken form.
 */
async function loadListData(): Promise<void> {
	if (!isEditMode.value || !listId.value) return;

	if (!apis.mapexOS?.lists) {
		notifyFail({ message: t.errors.apiNotInitialized.value });
		return;
	}

	isLoading.value = true;
	try {
		const data: ListResponse = await apis.mapexOS.lists.getById({ listId: listId.value });

		isSystemList.value = data.isSystem || false;

		listData.value = {
			type: (data.type as ListType) ?? null,
			parentId: data.parentId ?? null,
			name: data.name ?? '',
			value: data.value ?? '',
			enabled: data.enabled ?? true,
			isTemplate: data.isTemplate ?? false,
		};

		// In edit mode start the user at the review step; everything is pre-filled
		currentStep.value = STEP.REVIEW;

		await loadParentOptions();
	} catch (error: any) {
		logger.error('Failed to load list:', error);
		notifyFail({ message: t.createEditNotifications.loadFailed.value });
		await router.push('/admin/lists');
	} finally {
		isLoading.value = false;
	}
}

/** WATCHERS */

/**
 * Reload the parent dropdown whenever the type changes (cascade behaviour)
 */
watch(
	() => listData.value.type,
	(newType, oldType) => {
		if (newType !== oldType) {
			void loadParentOptions();
		}
	},
);

/** COMPOSABLES USAGE */
useStepperNavigation({
	currentStep,
	totalSteps: TOTAL_STEPS,
	changeStep: handlers.handleStepChange,
});

/** LIFECYCLE HOOKS */
onMounted(() => {
	void loadListData();
});
</script>

<template>
	<q-page class="q-pa-lg">
		<!-- Loading state (edit mode initial fetch) -->
		<div v-if="isLoading" class="flex flex-center q-pa-xl">
			<q-spinner color="primary" size="50px" />
			<span class="q-ml-md text-grey-7">{{ t.messages.loadingList.value }}</span>
		</div>

		<!-- Form content -->
		<div v-else>
			<!-- Header -->
			<PageHeader
				:icon="pageIcon"
				icon-color="primary"
				:title="pageTitle"
				:description="t.page.descriptionCreate.value"
				:button="{
					label: t.page.backButton.value,
					icon: 'arrow_back',
					flat: true,
					to: '/admin/lists',
				}"
			/>

			<!-- System list warning -->
			<q-banner
				v-if="isEditMode && isSystemList"
				rounded
				class="bg-warning text-white q-mb-lg"
			>
				<template #avatar>
					<q-icon name="lock" color="white" />
				</template>
				{{ t.messages.systemListWarning.value }}
			</q-banner>

			<!-- Body -->
			<div class="row q-col-gutter-lg">
				<!-- Progress stepper -->
				<div class="col-12 col-md-4">
					<StepperVertical
						:title="t.sections.progressSteps.value"
						:subtitle="t.messages.completeAllSteps.value"
						:info-text="t.messages.allFieldsRequired.value"
						:current-step-label="t.messages.currentStep.value"
						:current-step="currentStep"
						:steps="translatedSteps"
						:allow-step-navigation="isEditMode"
						@step-click="handlers.handleStepChange"
					/>
				</div>

				<!-- Form card -->
				<div class="col-12 col-md-8">
					<FormCard
						:header="(translatedSteps[currentStep - 1] as unknown as FormCardHeader)"
						:navigation="formNavigation"
						:button-labels="buttonLabels"
						@previous="handlers.handleStepChange"
						@next="handlers.handleStepChange"
						@save="handlers.submitForm"
					>
						<template #form>
							<Step1Basic
								v-if="currentStep === STEP.BASIC_INFO"
								ref="step1Ref"
								:model-value="listData"
								:is-edit-mode="isEditMode"
								:is-system-list="isSystemList"
								:loading-parents="loadingParents"
								:parent-options="parentOptions"
								@update:model-value="updateListData"
							/>

							<Step2Review
								v-else-if="currentStep === STEP.REVIEW"
								:model-value="listData"
								:is-edit-mode="isEditMode"
								:parent-name="parentName"
								@edit-section="handlers.handleStepChange"
							/>
						</template>
					</FormCard>
				</div>
			</div>
		</div>
	</q-page>
</template>
