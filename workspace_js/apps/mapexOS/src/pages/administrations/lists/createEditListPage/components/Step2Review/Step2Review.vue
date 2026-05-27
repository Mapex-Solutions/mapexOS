<template>
	<FormReview
		:sections="reviewSections"
		:description="t.reviewStep.subtitle.value"
		:show-success-banner="true"
		:success-message="isEditMode ? t.reviewStep.successMessageEdit.value : t.reviewStep.successMessage.value"
		@edit-section="emit('editSection', $event)"
	/>
</template>

<script setup lang="ts">
defineOptions({
	name: 'Step2Review',
});

/** TYPE IMPORTS */
import type { Step2ReviewProps, Step2ReviewEmits } from './interfaces';
import type { ReviewSectionDef } from '@components/forms/review/interfaces';

/** VUE IMPORTS */
import { computed } from 'vue';

/** COMPONENTS */
import { FormReview } from '@components/forms';

/** COMPOSABLES */
import { useListsTranslations } from '@composables/i18n/pages/administrations/lists/useListsTranslations';

/** LOCAL IMPORTS */
import { STEP } from '../../constants';

/** PROPS & EMITS */
const props = withDefaults(defineProps<Step2ReviewProps>(), {
	isEditMode: false,
	parentName: null,
});

const emit = defineEmits<Step2ReviewEmits>();

/** COMPOSABLES & STORES */
const t = useListsTranslations();

/** COMPUTED */

/**
 * Pretty label for the chosen type ("category", "manufacturer", "model")
 */
const typeLabel = computed(() => {
	const type = props.modelValue.type;
	if (!type) return t.reviewStep.values.notSelected.value;
	const map: Record<string, string> = {
		asset_category: t.typeOptions.asset_category.value,
		asset_manufacturer: t.typeOptions.asset_manufacturer.value,
		asset_model: t.typeOptions.asset_model.value,
	};
	return map[type] || type;
});

/**
 * Parent display value — falls back to "no parent" for top-of-tree types
 */
const parentDisplay = computed(() => {
	if (props.modelValue.type === 'asset_category') {
		return t.reviewStep.values.noParent.value;
	}
	return props.parentName || t.reviewStep.values.notSelected.value;
});

/**
 * Build the FormReview sections from the form data
 */
const reviewSections = computed((): ReviewSectionDef[] => {
	const data = props.modelValue;
	return [
		{
			stepNumber: STEP.BASIC_INFO,
			label: t.reviewStep.sections.basicInfo.value,
			icon: { name: 'badge', color: 'primary' },
			fields: [
				{
					label: t.reviewStep.fields.type.value,
					value: typeLabel.value,
					type: 'chip',
					badgeColors: 'primary',
					icon: 'category',
					colSize: 6,
				},
				{
					label: t.reviewStep.fields.parent.value,
					value: parentDisplay.value,
					type: 'chip',
					badgeColors: data.parentId ? 'indigo' : 'grey',
					icon: 'account_tree',
					colSize: 6,
				},
				{
					label: t.reviewStep.fields.name.value,
					value: data.name || t.reviewStep.values.notProvided.value,
					type: 'text',
					colSize: 6,
				},
				{
					label: t.reviewStep.fields.value.value,
					value: data.value || t.reviewStep.values.notProvided.value,
					type: 'text',
					colSize: 6,
				},
				{
					label: t.reviewStep.fields.enabled.value,
					value: data.enabled,
					type: 'boolean',
					colSize: 6,
				},
				{
					label: t.reviewStep.fields.isTemplate.value,
					value: data.isTemplate,
					type: 'boolean',
					colSize: 6,
				},
			],
		},
	];
});
</script>
