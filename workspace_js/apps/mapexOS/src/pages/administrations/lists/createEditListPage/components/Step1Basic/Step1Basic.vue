<template>
	<q-form ref="formRef" greedy>
		<div class="q-mb-md">
			<div class="text-subtitle1 text-weight-medium q-mb-xs">
				<q-icon name="badge" color="primary" class="q-mr-xs" />
				{{ t.sections.basicInfo.value }}
			</div>
			<div class="text-body2 text-grey-7">
				{{ t.formDescriptions.basicInfo.value }}
			</div>
		</div>

		<div class="row q-col-gutter-md">
			<!-- Type -->
			<div class="col-12 col-md-6">
				<q-select
					v-model="localData.type"
					outlined
					dense
					emit-value
					map-options
					class="rounded-borders"
					:label="`${t.fields.type.value} *`"
					:options="typeOptions"
					:disable="isEditMode || isSystemList"
					:rules="[(val) => !!val || t.validation.typeRequired.value]"
					@update:model-value="onTypeChange"
				>
					<template #prepend>
						<q-icon name="category" color="primary" />
					</template>
				</q-select>
				<div v-if="isEditMode" class="text-caption text-grey-6 q-mt-xs">
					<q-icon name="lock" size="xs" class="q-mr-xs" />
					{{ t.messages.typeLocked.value }}
				</div>
			</div>

			<!-- Parent (cascade) -->
			<div class="col-12 col-md-6">
				<q-select
					v-if="parentRequired"
					v-model="localData.parentId"
					outlined
					dense
					emit-value
					map-options
					class="rounded-borders"
					:label="`${t.fields.parent.value} *`"
					:options="parentSelectOptions"
					:loading="loadingParents"
					:disable="!localData.type || isSystemList"
					:rules="[(val) => !!val || t.validation.parentRequired.value]"
					@update:model-value="updateValue"
				>
					<template #prepend>
						<q-icon name="account_tree" color="primary" />
					</template>
				</q-select>
				<q-banner
					v-else-if="localData.type === 'asset_category'"
					rounded
					dense
					class="bg-grey-2 text-grey-8"
				>
					<template #avatar>
						<q-icon name="info" color="primary" />
					</template>
					{{ t.messages.noParentForCategory.value }}
				</q-banner>
			</div>

			<!-- Name -->
			<div class="col-12 col-md-6">
				<q-input
					v-model="localData.name"
					outlined
					dense
					class="rounded-borders"
					:label="`${t.fields.name.value} *`"
					:maxlength="NAME_MAX_LENGTH"
					:disable="isSystemList"
					:rules="[
						(val) => !!val?.trim() || t.validation.nameRequired.value,
						(val) => (val?.length ?? 0) <= NAME_MAX_LENGTH || t.validation.nameMaxLength.value,
					]"
					@update:model-value="updateValue"
				>
					<template #prepend>
						<q-icon name="label" color="primary" />
					</template>
				</q-input>
			</div>

			<!-- Value -->
			<div class="col-12 col-md-6">
				<q-input
					v-model="localData.value"
					outlined
					dense
					class="rounded-borders"
					:label="`${t.fields.value.value} *`"
					:maxlength="VALUE_MAX_LENGTH"
					:disable="isSystemList"
					:rules="[
						(val) => !!val?.trim() || t.validation.valueRequired.value,
						(val) => (val?.length ?? 0) <= VALUE_MAX_LENGTH || t.validation.valueMaxLength.value,
						(val) => !val || VALUE_PATTERN.test(val) || t.validation.valuePattern.value,
					]"
					@update:model-value="updateValue"
				>
					<template #prepend>
						<q-icon name="code" color="primary" />
					</template>
				</q-input>
			</div>

			<!-- Enabled & isTemplate toggles -->
			<div class="col-12 col-md-6">
				<div class="q-py-sm">
					<q-toggle
						v-model="localData.enabled"
						color="primary"
						:label="t.fields.enabled.value"
						:disable="isSystemList"
						@update:model-value="updateValue"
					/>
				</div>
			</div>

			<div class="col-12 col-md-6">
				<div class="q-py-sm">
					<q-toggle
						v-model="localData.isTemplate"
						color="primary"
						:label="t.fields.isTemplate.value"
						:disable="isSystemList"
						@update:model-value="updateValue"
					/>
					<div class="text-caption text-grey-7 q-pl-lg">
						{{ t.formDescriptions.isTemplate.value }}
					</div>
				</div>
			</div>
		</div>
	</q-form>
</template>

<script setup lang="ts">
defineOptions({
	name: 'Step1Basic',
});

/** TYPE IMPORTS */
import type { Step1BasicProps } from './interfaces';
import type { QForm } from 'quasar';
import type { ListFormData, ListType } from '../../interfaces';

/** VUE IMPORTS */
import { ref, reactive, computed, watch } from 'vue';

/** COMPOSABLES */
import { useListsTranslations } from '@composables/i18n/pages/administrations/lists/useListsTranslations';

/** LOCAL IMPORTS */
import {
	NAME_MAX_LENGTH,
	VALUE_MAX_LENGTH,
	VALUE_PATTERN,
	PARENT_TYPE_FOR,
} from '../../constants';

/** PROPS & EMITS */
const props = withDefaults(defineProps<Step1BasicProps>(), {
	isEditMode: false,
	isSystemList: false,
	loadingParents: false,
});

const emit = defineEmits<{
	(e: 'update:modelValue', value: Partial<ListFormData>): void;
}>();

/** COMPOSABLES & STORES */
const t = useListsTranslations();

/** STATE */
const formRef = ref<QForm | null>(null);
const localData = reactive<ListFormData>({
	type: props.modelValue.type,
	parentId: props.modelValue.parentId,
	name: props.modelValue.name || '',
	value: props.modelValue.value || '',
	enabled: props.modelValue.enabled ?? true,
	isTemplate: props.modelValue.isTemplate ?? false,
});

/** COMPUTED */

/**
 * Whether the current type requires a parent (everything except asset_category)
 */
const parentRequired = computed(() => {
	if (!localData.type) return false;
	return PARENT_TYPE_FOR[localData.type] !== null;
});

/**
 * Localised type options for the type select
 */
const typeOptions = computed(() => [
	{ label: t.typeOptions.asset_category.value, value: 'asset_category' },
	{ label: t.typeOptions.asset_manufacturer.value, value: 'asset_manufacturer' },
	{ label: t.typeOptions.asset_model.value, value: 'asset_model' },
]);

/**
 * Parent options mapped to the q-select shape (`{ label, value }`)
 */
const parentSelectOptions = computed(() =>
	(props.parentOptions ?? []).map((opt) => ({
		label: `${opt.name} (${opt.value})`,
		value: opt.id,
	})),
);

/** WATCHERS */
watch(
	() => props.modelValue,
	(newVal) => {
		localData.type = newVal.type;
		localData.parentId = newVal.parentId;
		localData.name = newVal.name || '';
		localData.value = newVal.value || '';
		localData.enabled = newVal.enabled ?? true;
		localData.isTemplate = newVal.isTemplate ?? false;
	},
	{ deep: true },
);

/** FUNCTIONS */

/**
 * Emit current local state to the parent (used by every field update)
 */
function updateValue(): void {
	emit('update:modelValue', { ...localData });
}

/**
 * When the type changes, reset the parent selection so we never carry a
 * parentId pointing at the wrong level of the tree.
 *
 * @param value - newly selected type
 */
function onTypeChange(value: ListType | null): void {
	localData.type = value;
	localData.parentId = null;
	updateValue();
}

/**
 * Validate the embedded q-form. Required by the parent stepper.
 *
 * @returns whether all fields pass their rules
 */
async function validate(): Promise<boolean> {
	const result = await formRef.value?.validate();
	return !!result;
}

/** EXPOSE */
defineExpose({
	formRef,
	validate,
});
</script>

<style scoped lang="scss">
.rounded-borders {
	border-radius: var(--mapex-radius-md);
}
</style>
