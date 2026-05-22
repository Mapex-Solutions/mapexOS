<script setup lang="ts">
defineOptions({
	name: 'Step8DynamicFields'
});

/** TYPE IMPORTS */
import type { Step8DynamicFieldsProps, Step8DynamicFieldsEmits } from './interfaces/Step8DynamicFields.interface';
import type { AssetTemplateData, DynamicFieldMapping, DynamicFieldType } from '../../interfaces';

/** VUE IMPORTS */
import { computed, ref } from 'vue';

/** COMPONENTS */
import { DetailChip } from '@components/chips';
import { AppTooltip } from '@components/tooltips';

/** COMPOSABLES */
import { useAddAssetTemplateTranslations } from '@src/composables/i18n/pages/assets/addAssetTemplate/useAddAssetTemplateTranslations';

/** LOCAL IMPORTS */
import {
	DYNAMIC_FIELD_TYPE_OPTIONS,
	DYNAMIC_FIELDS_MAX,
	DYNAMIC_FIELDS_WARNING_THRESHOLD,
} from '../../constants';

const props = defineProps<Step8DynamicFieldsProps>();
const emit = defineEmits<Step8DynamicFieldsEmits>();

/** COMPOSABLES */
const t = useAddAssetTemplateTranslations();

/** STATE */
const newField = ref<DynamicFieldMapping>({
	field: '',
	value: '',
	type: 'string',
	latitudePath: '',
	longitudePath: '',
});

/** COMPUTED */
const data = computed({
	get: () => props.modelValue,
	set: (value: AssetTemplateData) => emit('update:modelValue', value)
});

const dynamicFields = computed({
	get: () => data.value.dynamicFields || [],
	set: (value: DynamicFieldMapping[]) => {
		data.value = { ...data.value, dynamicFields: value };
	}
});

const availableFieldsOptions = computed(() => {
	return (data.value.availableFields || []).map(field => ({
		label: field,
		value: field,
	}));
});

const hasAvailableFields = computed(() => {
	return (data.value.availableFields || []).length > 0;
});

const isGeoType = computed(() => newField.value.type === 'geo');

const fieldCount = computed(() => dynamicFields.value.length);
const isAtLimit = computed(() => fieldCount.value >= DYNAMIC_FIELDS_MAX);
const isNearLimit = computed(() => fieldCount.value >= DYNAMIC_FIELDS_WARNING_THRESHOLD && !isAtLimit.value);
const remainingFields = computed(() => DYNAMIC_FIELDS_MAX - fieldCount.value);

const canAddField = computed(() => {
	// Check limit first
	if (isAtLimit.value) return false;

	if (!newField.value.field.trim()) return false;

	if (newField.value.type === 'geo') {
		return !!newField.value.latitudePath && !!newField.value.longitudePath;
	}

	return !!newField.value.value;
});

/** FUNCTIONS */

/**
 * Add a new dynamic field mapping
 */
function addField(): void {
	if (!canAddField.value) return;

	const fieldToAdd: DynamicFieldMapping = {
		field: newField.value.field.trim(),
		value: newField.value.type === 'geo' ? '' : newField.value.value,
		type: newField.value.type,
	};

	if (newField.value.type === 'geo') {
		fieldToAdd.latitudePath = newField.value.latitudePath || '';
		fieldToAdd.longitudePath = newField.value.longitudePath || '';
	}

	dynamicFields.value = [...dynamicFields.value, fieldToAdd];

	// Reset form
	newField.value = {
		field: '',
		value: '',
		type: 'string',
		latitudePath: '',
		longitudePath: '',
	};
}

/**
 * Remove a dynamic field mapping by index
 *
 * @param {number} index - Index of field to remove
 */
function removeField(index: number): void {
	const updated = [...dynamicFields.value];
	updated.splice(index, 1);
	dynamicFields.value = updated;
}

/**
 * Get icon for field type
 *
 * @param {DynamicFieldType} type - Field type
 * @returns {string} Icon name
 */
function getTypeIcon(type: DynamicFieldType): string {
	const option = DYNAMIC_FIELD_TYPE_OPTIONS.find(opt => opt.value === type);
	return option?.icon || 'mdi-help';
}

/**
 * Get label for field type
 *
 * @param {DynamicFieldType} type - Field type
 * @returns {string} Type label
 */
function getTypeLabel(type: DynamicFieldType): string {
	const option = DYNAMIC_FIELD_TYPE_OPTIONS.find(opt => opt.value === type);
	return option?.label || type;
}

/**
 * Filter function for autocomplete
 *
 * @param {string} val - Search value
 * @param {Function} update - Update function
 */
function filterFn(val: string, update: (fn: () => void) => void): void {
	update(() => {
		// The filtering is handled by QSelect's filter prop
	});
}
</script>

<template>
	<div>
		<!-- Header -->
		<div class="q-mb-md">
			<div class="text-subtitle1 text-weight-medium q-mb-xs">
				<q-icon name="mdi-database-cog" color="primary" class="q-mr-xs" />
				{{ t.steps.step8.title.value }}
			</div>
			<div class="text-body2 text-grey-7">
				{{ t.steps.step8.subtitle.value }}
			</div>
		</div>

		<!-- Info Banner -->
		<q-banner rounded class="bg-blue-1 text-primary q-mb-md">
			<template v-slot:avatar>
				<q-icon name="info" color="primary" />
			</template>
			<div class="text-subtitle2 text-weight-medium q-mb-xs">
				{{ t.steps.step8.banner.title.value }}
			</div>
			<div class="text-body2">
				{{ t.steps.step8.banner.description.value }}
			</div>
		</q-banner>

		<!-- Field Count Progress -->
		<div class="row items-center q-mb-md">
			<div class="col">
				<q-linear-progress
					:value="fieldCount / DYNAMIC_FIELDS_MAX"
					:color="isAtLimit ? 'negative' : isNearLimit ? 'warning' : 'primary'"
					rounded
					size="8px"
					class="q-mb-xs"
				/>
				<div class="row justify-between text-caption">
					<span :class="isAtLimit ? 'text-negative' : isNearLimit ? 'text-warning' : 'text-grey-7'">
						{{ fieldCount }} / {{ DYNAMIC_FIELDS_MAX }} fields
					</span>
					<span v-if="!isAtLimit" class="text-grey-6">
						{{ remainingFields }} remaining
					</span>
				</div>
			</div>
		</div>

		<!-- Warning: Limit Reached -->
		<q-banner v-if="isAtLimit" rounded class="bg-red-1 text-negative q-mb-md">
			<template v-slot:avatar>
				<q-icon name="error" color="negative" />
			</template>
			<div class="text-subtitle2 text-weight-medium q-mb-xs">
				Field limit reached
			</div>
			<div class="text-body2">
				You have reached the maximum of {{ DYNAMIC_FIELDS_MAX }} dynamic fields.
				Remove existing fields to add new ones.
			</div>
		</q-banner>

		<!-- Warning: Near Limit -->
		<q-banner v-else-if="isNearLimit" rounded class="bg-orange-1 text-warning q-mb-md">
			<template v-slot:avatar>
				<q-icon name="warning" color="warning" />
			</template>
			<div class="text-subtitle2 text-weight-medium q-mb-xs">
				Approaching field limit
			</div>
			<div class="text-body2">
				You have {{ fieldCount }} fields. Maximum is {{ DYNAMIC_FIELDS_MAX }}.
				{{ remainingFields }} more fields can be added.
			</div>
		</q-banner>

		<!-- Warning if no available fields -->
		<q-banner v-else-if="!hasAvailableFields" rounded class="bg-orange-1 text-warning q-mb-md">
			<template v-slot:avatar>
				<q-icon name="warning" color="warning" />
			</template>
			<div class="text-subtitle2 text-weight-medium q-mb-xs">
				{{ t.steps.step8.noFieldsWarning.title.value }}
			</div>
			<div class="text-body2">
				{{ t.steps.step8.noFieldsWarning.description.value }}
			</div>
		</q-banner>

		<!-- Add New Field Form -->
		<q-card flat bordered class="q-mb-md">
			<q-card-section>
				<div class="text-subtitle2 text-weight-medium q-mb-md">
					<q-icon name="add" color="primary" class="q-mr-xs" />
					{{ t.steps.step8.addField.title.value }}
				</div>

				<div class="row q-col-gutter-md">
					<!-- Field Name -->
					<div class="col-12 col-sm-6 col-md-3">
						<q-input
							v-model="newField.field"
							outlined
							dense
							class="rounded-borders"
							:label="t.steps.step8.addField.fieldName.label.value"
							:placeholder="t.steps.step8.addField.fieldName.placeholder.value"
						>
							<template v-slot:prepend>
								<q-icon name="mdi-form-textbox" color="primary" />
							</template>
						</q-input>
					</div>

					<!-- Field Type -->
					<div class="col-12 col-sm-6 col-md-3">
						<q-select
							v-model="newField.type"
							outlined
							dense
							class="rounded-borders"
							:label="t.steps.step8.addField.fieldType.label.value"
							:options="DYNAMIC_FIELD_TYPE_OPTIONS"
							option-value="value"
							option-label="label"
							emit-value
							map-options
						>
							<template v-slot:prepend>
								<q-icon :name="getTypeIcon(newField.type)" color="primary" />
							</template>
							<template v-slot:option="scope">
								<q-item v-bind="scope.itemProps">
									<q-item-section avatar>
										<q-icon :name="scope.opt.icon" color="primary" />
									</q-item-section>
									<q-item-section>
										<q-item-label>{{ scope.opt.label }}</q-item-label>
										<q-item-label caption>{{ scope.opt.description }}</q-item-label>
									</q-item-section>
								</q-item>
							</template>
						</q-select>
					</div>

					<!-- Value Path (for non-geo types) -->
					<div v-if="!isGeoType" class="col-12 col-sm-6 col-md-4">
						<q-select
							v-model="newField.value"
							outlined
							dense
							use-input
							fill-input
							hide-selected
							input-debounce="0"
							class="rounded-borders"
							:label="t.steps.step8.addField.valuePath.label.value"
							:placeholder="t.steps.step8.addField.valuePath.placeholder.value"
							:options="availableFieldsOptions"
							option-value="value"
							option-label="label"
							emit-value
							map-options
							:disable="!hasAvailableFields"
							@filter="filterFn"
						>
							<template v-slot:prepend>
								<q-icon name="mdi-code-json" color="primary" />
							</template>
							<template v-slot:no-option>
								<q-item>
									<q-item-section class="text-grey">
										{{ t.steps.step8.addField.valuePath.noOptions.value }}
									</q-item-section>
								</q-item>
							</template>
						</q-select>
					</div>

					<!-- Geo: Latitude Path -->
					<div v-if="isGeoType" class="col-12 col-sm-6 col-md-3">
						<q-select
							v-model="newField.latitudePath"
							outlined
							dense
							use-input
							fill-input
							hide-selected
							input-debounce="0"
							class="rounded-borders"
							:label="t.steps.step8.addField.latitudePath.label.value"
							:placeholder="t.steps.step8.addField.latitudePath.placeholder.value"
							:options="availableFieldsOptions"
							option-value="value"
							option-label="label"
							emit-value
							map-options
							:disable="!hasAvailableFields"
							@filter="filterFn"
						>
							<template v-slot:prepend>
								<q-icon name="mdi-latitude" color="primary" />
							</template>
						</q-select>
					</div>

					<!-- Geo: Longitude Path -->
					<div v-if="isGeoType" class="col-12 col-sm-6 col-md-3">
						<q-select
							v-model="newField.longitudePath"
							outlined
							dense
							use-input
							fill-input
							hide-selected
							input-debounce="0"
							class="rounded-borders"
							:label="t.steps.step8.addField.longitudePath.label.value"
							:placeholder="t.steps.step8.addField.longitudePath.placeholder.value"
							:options="availableFieldsOptions"
							option-value="value"
							option-label="label"
							emit-value
							map-options
							:disable="!hasAvailableFields"
							@filter="filterFn"
						>
							<template v-slot:prepend>
								<q-icon name="mdi-longitude" color="primary" />
							</template>
						</q-select>
					</div>

					<!-- Add Button -->
					<div class="col-12 col-sm-6 col-md-2 flex items-center">
						<q-btn
							unelevated
							:color="isAtLimit ? 'grey' : 'primary'"
							icon="add"
							:label="t.steps.step8.addField.addButton.value"
							:disable="!canAddField"
							class="full-width"
							@click="addField"
						>
							<AppTooltip v-if="isAtLimit" content="Maximum field limit reached" />
						</q-btn>
					</div>
				</div>
			</q-card-section>
		</q-card>

		<!-- Mapped Fields List -->
		<div v-if="dynamicFields.length > 0" class="q-mb-md">
			<div class="text-subtitle2 text-weight-medium q-mb-sm row items-center">
				<q-icon name="mdi-format-list-bulleted" color="primary" class="q-mr-xs" />
				{{ t.steps.step8.mappedFields.title.value }}
				<q-badge
					:color="isAtLimit ? 'negative' : isNearLimit ? 'warning' : 'primary'"
					class="q-ml-sm"
				>
					{{ fieldCount }} / {{ DYNAMIC_FIELDS_MAX }}
				</q-badge>
			</div>

			<q-list bordered separator class="rounded-borders">
				<q-item v-for="(field, index) in dynamicFields" :key="`field-${index}`">
					<q-item-section avatar>
						<q-icon :name="getTypeIcon(field.type)" color="primary" />
					</q-item-section>

					<q-item-section>
						<q-item-label class="text-weight-medium">
							{{ field.field }}
						</q-item-label>
						<q-item-label caption>
							<template v-if="field.type === 'geo'">
								<span class="text-grey-7">lat:</span> {{ field.latitudePath }}
								<span class="q-mx-xs">|</span>
								<span class="text-grey-7">lng:</span> {{ field.longitudePath }}
							</template>
							<template v-else>
								<span class="text-grey-7">path:</span> {{ field.value }}
							</template>
						</q-item-label>
					</q-item-section>

					<q-item-section side>
						<DetailChip
							:label="getTypeLabel(field.type)"
							:color="(field.type === 'string' ? 'blue' :
								field.type === 'number' ? 'green' :
								field.type === 'bool' ? 'orange' :
								field.type === 'date' ? 'purple' : 'red') as any"
							dense
						/>
					</q-item-section>

					<q-item-section side>
						<q-btn
							flat
							round
							dense
							icon="delete"
							color="negative"
							@click="removeField(index)"
						>
							<AppTooltip :content="t.steps.step8.mappedFields.removeTooltip.value" />
						</q-btn>
					</q-item-section>
				</q-item>
			</q-list>
		</div>

		<!-- Empty State -->
		<div v-else class="text-center q-pa-lg text-grey-6">
			<q-icon name="mdi-database-off" size="48px" class="q-mb-md" />
			<div class="text-subtitle1">{{ t.steps.step8.emptyState.title.value }}</div>
			<div class="text-body2">{{ t.steps.step8.emptyState.description.value }}</div>
		</div>
	</div>
</template>

<style scoped>
.rounded-borders {
	border-radius: var(--mapex-radius-md);
}
</style>
