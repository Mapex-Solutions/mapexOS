<script setup lang="ts">
defineOptions({
	name: 'Step9Review'
});

/** TYPE IMPORTS */
import type { Step9ReviewProps, Step9ReviewEmits } from './interfaces/Step9Review.interface';
import type { ReviewSectionDef } from '@components/forms/review/interfaces';

/** VUE IMPORTS */
import { computed } from 'vue';

/** COMPONENTS */
import { FormReview } from '@components/forms';
import { AvailableFieldsList } from '@components/lists/availableFieldsList';
import { DetailChip } from '@components/chips';

/** COMPOSABLES */
import { useCommonActions } from '@composables/i18n';
import { useAddAssetTemplateTranslations } from '@src/composables/i18n/pages/assets/addAssetTemplate/useAddAssetTemplateTranslations';

/** LOCAL IMPORTS */
import { DYNAMIC_FIELD_TYPE_OPTIONS } from '../../constants';

/** PROPS & EMITS */
const props = defineProps<Step9ReviewProps>();
const emit = defineEmits<Step9ReviewEmits>();

/** COMPOSABLES & STORES */
const t = useAddAssetTemplateTranslations();
const { actions: commonActions } = useCommonActions();

/** COMPUTED */
const previewData = computed((): ReviewSectionDef[] => {
	const data = props.modelValue;

	return [
		// Basic Information Section
		{
			stepNumber: 1,
			label: t.steps.step9.reviewSection.basicInfo.value,
			icon: { name: 'info', color: 'primary' },
			fields: [
				{
					label: t.steps.step9.reviewSection.fields.name.value,
					value: data.name,
					type: 'text',
					colSize: 6
				},
				{
					label: t.steps.step9.reviewSection.fields.status.value,
					value: data.enabled
						? t.statusOptions.active.label.value
						: t.statusOptions.inactive.label.value,
					type: 'badge',
					badgeColors: {
						[t.statusOptions.active.label.value]: 'positive',
						[t.statusOptions.inactive.label.value]: 'negative'
					},
					colSize: 6,
				},
				{
					label: t.steps.step9.reviewSection.fields.description.value,
					value: data.description || '-',
					type: 'text',
					colSize: 12
				},
				{
					label: t.steps.step9.reviewSection.fields.manufacturer.value,
					value: data.manufacturerName || '-',
					type: 'text',
					colSize: 4
				},
				{
					label: t.steps.step9.reviewSection.fields.model.value,
					value: data.modelName || '-',
					type: 'text',
					colSize: 4
				},
				{
					label: t.steps.step9.reviewSection.fields.category.value,
					value: data.categoryName || '-',
					type: 'text',
					colSize: 4
				},
			],
		},
		// Asset ID Path Section
		{
			stepNumber: 2,
			label: t.steps.step2.title.value,
			icon: { name: 'route', color: 'secondary' },
			fields: [
				{
					label: t.steps.step9.reviewSection.fields.assetIdPath.value,
					value: data.assetIdPath,
					type: 'text',
					colSize: 12,
				},
			],
		},
		// Scripts Summary Section
		{
			stepNumber: 3,
			label: t.steps.step9.reviewSection.scriptSummary.value,
			icon: { name: 'code', color: 'primary' },
			fields: [
				{
					label: t.steps.step9.reviewSection.fields.preprocessorScript.value,
					value: data.scriptProcessor
						? t.steps.step9.reviewSection.fields.configured.value
						: t.steps.step9.reviewSection.fields.notConfigured.value,
					type: 'badge',
					badgeColors: {
						[t.steps.step9.reviewSection.fields.configured.value]: 'positive',
						[t.steps.step9.reviewSection.fields.notConfigured.value]: 'grey'
					},
					colSize: 6,
				},
				{
					label: t.steps.step9.reviewSection.fields.validationScript.value,
					value: data.scriptValidator
						? t.steps.step9.reviewSection.fields.configured.value
						: t.steps.step9.reviewSection.fields.notConfigured.value,
					type: 'badge',
					badgeColors: {
						[t.steps.step9.reviewSection.fields.configured.value]: 'positive',
						[t.steps.step9.reviewSection.fields.notConfigured.value]: 'negative'
					},
					colSize: 6,
				},
				{
					label: t.steps.step9.reviewSection.fields.conversionScript.value,
					value: data.scriptConversion
						? t.steps.step9.reviewSection.fields.configured.value
						: t.steps.step9.reviewSection.fields.notConfigured.value,
					type: 'badge',
					badgeColors: {
						[t.steps.step9.reviewSection.fields.configured.value]: 'positive',
						[t.steps.step9.reviewSection.fields.notConfigured.value]: 'negative'
					},
					colSize: 6,
				},
				{
					label: t.steps.step9.reviewSection.fields.testScript.value,
					value: data.scriptTest
						? t.steps.step9.reviewSection.fields.configured.value
						: t.steps.step9.reviewSection.fields.notConfigured.value,
					type: 'badge',
					badgeColors: {
						[t.steps.step9.reviewSection.fields.configured.value]: 'positive',
						[t.steps.step9.reviewSection.fields.notConfigured.value]: 'grey'
					},
					colSize: 6,
				},
			],
		},
	];
});

const dynamicFieldsCount = computed(() => {
	return (props.modelValue.dynamicFields || []).length;
});

const hasDynamicFields = computed(() => dynamicFieldsCount.value > 0);

/**
 * Get type label from type value
 *
 * @param {string} type - Field type value
 * @returns {string} Type label
 */
function getTypeLabel(type: string): string {
	const option = DYNAMIC_FIELD_TYPE_OPTIONS.find(opt => opt.value === type);
	return option?.label || type;
}

/**
 * Handle edit section click
 *
 * @param {number} stepNumber - Step number to navigate to
 */
function handleEditSection(stepNumber: number): void {
	emit('editSection', stepNumber);
}
</script>

<template>
	<div>
		<!-- Configuration Review Sections (Steps 1-6) -->
		<FormReview
			:sections="previewData"
			:description="t.steps.step9.subtitle.value"
			:show-success-banner="false"
			@edit-section="handleEditSection"
		/>

		<!-- Step 7: Available Fields -->
		<q-card flat bordered class="q-mb-md">
			<q-card-section>
				<div class="row items-center q-mb-md">
					<div class="col">
						<div class="row items-center">
							<q-icon name="mdi-format-list-bulleted" size="sm" color="primary" class="q-mr-sm" />
							<span class="text-h6 text-weight-medium text-dark">
								{{ t.steps.step9.availableFields.title.value }}
							</span>
							<DetailChip
								:label="String((modelValue.availableFields || []).length)"
								color="grey"
								size="sm"
								dense
								class="q-ml-sm"
							/>
						</div>
					</div>
					<div class="col-auto">
						<q-btn
							flat
							dense
							size="sm"
							icon="edit"
							color="primary"
							:label="commonActions.edit.value"
							@click="handleEditSection(7)"
						/>
					</div>
				</div>

				<!-- Available Fields Content -->
				<div v-if="modelValue.availableFields && modelValue.availableFields.length > 0">
					<AvailableFieldsList
						:fields="modelValue.availableFields"
						:max-height="200"
					/>
				</div>
				<div v-else class="text-center text-grey-6 q-py-md">
					<q-icon name="mdi-format-list-bulleted-type" size="32px" class="q-mb-sm" />
					<div class="text-body2">{{ t.steps.step9.noFieldsWarning.value }}</div>
				</div>
			</q-card-section>
		</q-card>

		<!-- Step 8: Dynamic Fields -->
		<q-card flat bordered class="q-mb-md">
			<q-card-section>
				<div class="row items-center q-mb-md">
					<div class="col">
						<div class="row items-center">
							<q-icon name="mdi-database-cog" size="sm" color="primary" class="q-mr-sm" />
							<span class="text-h6 text-weight-medium text-dark">
								{{ t.steps.step9.dynamicFields.title.value }}
							</span>
							<DetailChip
								:label="String(dynamicFieldsCount)"
								color="grey"
								size="sm"
								dense
								class="q-ml-sm"
							/>
						</div>
					</div>
					<div class="col-auto">
						<q-btn
							flat
							dense
							size="sm"
							icon="edit"
							color="primary"
							:label="commonActions.edit.value"
							@click="handleEditSection(8)"
						/>
					</div>
				</div>

				<!-- Dynamic Fields Content -->
				<div v-if="hasDynamicFields">
					<q-list dense separator>
						<q-item v-for="(field, index) in modelValue.dynamicFields" :key="`dyn-${index}`">
							<q-item-section avatar>
								<DetailChip
									:label="getTypeLabel(field.type)"
									:color="(field.type === 'string' ? 'blue' :
										field.type === 'number' ? 'green' :
										field.type === 'bool' ? 'orange' :
										field.type === 'date' ? 'purple' : 'red') as any"
									size="sm"
									dense
								/>
							</q-item-section>
							<q-item-section>
								<q-item-label class="text-weight-medium">{{ field.field }}</q-item-label>
								<q-item-label caption>
									<template v-if="field.type === 'geo'">
										lat: {{ field.latitudePath }} | lng: {{ field.longitudePath }}
									</template>
									<template v-else>
										{{ field.value }}
									</template>
								</q-item-label>
							</q-item-section>
						</q-item>
					</q-list>
				</div>
				<div v-else class="text-center text-grey-6 q-py-md">
					<q-icon name="mdi-database-off" size="32px" class="q-mb-sm" />
					<div class="text-body2">{{ t.steps.step9.dynamicFields.empty.value }}</div>
				</div>
			</q-card-section>
		</q-card>

		<!-- Success Banner -->
		<q-banner rounded class="bg-green-1 text-green-9 q-mt-lg">
			<template #avatar>
				<q-icon name="check_circle" color="green-7" />
			</template>
			<div class="text-body2">
				<strong>{{ t.steps.step9.completeBanner.value }}</strong>
			</div>
		</q-banner>
	</div>
</template>
