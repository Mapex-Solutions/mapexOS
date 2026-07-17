<script setup lang="ts">
defineOptions({
	name: 'Step8DynamicFields'
});

/** TYPE IMPORTS */
import type { Step8DynamicFieldsProps, Step8DynamicFieldsEmits, FieldNameOption } from './interfaces/Step8DynamicFields.interface';
import type { AssetTemplateData, DynamicFieldMapping, DynamicFieldType } from '../../interfaces';
import type { QSelect } from 'quasar';

/** VUE IMPORTS */
import { computed, ref } from 'vue';

/** COMPONENTS */
import { DetailChip } from '@components/chips';
import { AppTooltip } from '@components/tooltips';
import { InfoBanner } from '@components/banners';

/** COMPOSABLES */
import { useAddAssetTemplateTranslations } from '@src/composables/i18n/pages/assets/addAssetTemplate/useAddAssetTemplateTranslations';
import { useFieldVocabulary } from './composables';

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
const { groups: vocabularyGroups } = useFieldVocabulary();

/** STATE */
const newField = ref<DynamicFieldMapping>({
	field: '',
	value: '',
	type: 'string',
	latitudePath: '',
	longitudePath: '',
});

/**
 * Currently displayed field-name suggestions after debounced input filtering.
 * Seeded from the full flattened vocabulary and narrowed by `filterFieldName`.
 */
const fieldNameOptions = ref<FieldNameOption[]>([]);

/**
 * Caption surfaced under the field-name control after picking a suggestion:
 * the localized hint plus unit. Cleared when the user types a custom name.
 */
const selectedFieldHint = ref('');

/**
 * Latest term typed into the field-name control. Kept in sync by
 * `filterFieldName` so the empty-result slot can offer to commit it verbatim
 * as a custom field name.
 */
const fieldNameSearch = ref('');

/**
 * Template ref to the field-name select, used to programmatically close its
 * menu after committing a custom name from the empty-result slot.
 */
const fieldNameSelect = ref<QSelect | null>(null);

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

/**
 * Flattened suggestion list: each backend group contributes a non-selectable
 * header (its localized label) followed by its selectable canonical fields,
 * preserving the order received so categories render as visual separators.
 */
const fieldNameSuggestions = computed<FieldNameOption[]>(() => {
	const options: FieldNameOption[] = [];

	for (const group of vocabularyGroups.value) {
		options.push({
			value: group.label,
			hint: '',
			unit: '',
			isHeader: true,
		});

		for (const field of group.fields) {
			options.push({
				value: field.value,
				hint: field.hint,
				unit: field.unit,
				type: field.type,
				isHeader: false,
			});
		}
	}

	return options;
});

/**
 * Lookup from a canonical field name to its suggestion, used to pre-select the
 * type and surface the hint/unit caption when a suggestion is picked.
 */
const suggestionByValue = computed<Map<string, FieldNameOption>>(() => {
	const map = new Map<string, FieldNameOption>();
	for (const option of fieldNameSuggestions.value) {
		if (!option.isHeader) map.set(option.value, option);
	}
	return map;
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

/**
 * Reason the Add button is disabled while the user is still below the field
 * limit, surfaced as a tooltip so the missing input is obvious. Empty when the
 * button is enabled or when the block is the hard limit (that case keeps its own
 * limit tooltip).
 */
const addFieldDisabledReason = computed(() => {
	if (canAddField.value || isAtLimit.value) return '';
	return t.steps.step8.addField.disabledReason.value;
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
	selectedFieldHint.value = '';
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

/**
 * Filter the field-name suggestions by the typed term.
 *
 * Matching is case-insensitive against the canonical value and localized hint.
 * Category headers are kept only when at least one of their fields survives the
 * filter, so empty separators never appear.
 *
 * @param {string} val - Search term typed by the user.
 * @param {Function} update - Quasar callback that commits the narrowed options.
 */
function filterFieldName(val: string, update: (fn: () => void) => void): void {
	fieldNameSearch.value = val.trim();

	update(() => {
		const term = val.trim().toLowerCase();

		if (!term) {
			fieldNameOptions.value = fieldNameSuggestions.value;
			return;
		}

		const result: FieldNameOption[] = [];
		let pendingHeader: FieldNameOption | null = null;

		for (const option of fieldNameSuggestions.value) {
			if (option.isHeader) {
				pendingHeader = option;
				continue;
			}

			const matches =
				option.value.toLowerCase().includes(term) ||
				option.hint.toLowerCase().includes(term);

			if (!matches) continue;

			if (pendingHeader) {
				result.push(pendingHeader);
				pendingHeader = null;
			}
			result.push(option);
		}

		fieldNameOptions.value = result;
	});
}

/**
 * Build the hint caption for a field name: the localized hint plus unit when
 * the name matches a known canonical suggestion, or an empty string for a
 * custom name (no suggestion).
 *
 * @param {string} name - The current field name.
 * @returns {string} The caption to show under the control.
 */
function buildFieldHint(name: string): string {
	const suggestion = suggestionByValue.value.get(name);
	if (!suggestion) return '';

	return suggestion.unit
		? `${suggestion.hint} (${t.steps.step8.addField.fieldName.unitLabel.value}: ${suggestion.unit})`
		: suggestion.hint;
}

/**
 * Capture the text typed into the field-name control on every keystroke so a
 * custom name commits without needing Enter or a selection. The type is left
 * untouched here — typing implies a custom field, so the user keeps control of
 * the type; only an explicit suggestion pick pre-selects it.
 *
 * @param {string | null} value - The raw text currently in the input.
 */
function onFieldNameTyped(value: string | null): void {
	const name = value ?? '';
	newField.value.field = name;
	selectedFieldHint.value = buildFieldHint(name);
}

/**
 * Apply a field name picked from the suggestion list: set the name, pre-select
 * the suggestion's type (still editable afterwards), and surface its localized
 * hint/unit caption.
 *
 * @param {string | null} value - The canonical name picked by the user.
 */
function onFieldNameSelected(value: string | null): void {
	const name = value ?? '';
	newField.value.field = name;

	const suggestion = suggestionByValue.value.get(name);
	if (suggestion?.type) newField.value.type = suggestion.type;

	selectedFieldHint.value = buildFieldHint(name);
}

/**
 * Commit the current search term as a custom field name from the empty-result
 * slot, reusing the typed handler so the name and its (empty) hint stay in sync,
 * then close the select menu.
 */
function commitCustomFieldName(): void {
	if (!fieldNameSearch.value) return;
	onFieldNameTyped(fieldNameSearch.value);
	fieldNameSelect.value?.hidePopup();
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
		<InfoBanner variant="info" :title="t.steps.step8.banner.title.value" class="q-mb-md">
			{{ t.steps.step8.banner.description.value }}
		</InfoBanner>

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
						{{ t.steps.step8.progress.count(fieldCount, DYNAMIC_FIELDS_MAX) }}
					</span>
					<span v-if="!isAtLimit" class="text-grey-6">
						{{ t.steps.step8.progress.remaining(remainingFields) }}
					</span>
				</div>
			</div>
		</div>

		<!-- Warning: Limit Reached -->
		<InfoBanner v-if="isAtLimit" variant="danger" :title="t.steps.step8.limitReached.title.value" class="q-mb-md">
			{{ t.steps.step8.limitReached.description(DYNAMIC_FIELDS_MAX) }}
		</InfoBanner>

		<!-- Warning: Near Limit -->
		<InfoBanner v-else-if="isNearLimit" variant="warning" :title="t.steps.step8.nearLimit.title.value" class="q-mb-md">
			{{ t.steps.step8.nearLimit.description(fieldCount, DYNAMIC_FIELDS_MAX, remainingFields) }}
		</InfoBanner>

		<!-- Warning if no available fields -->
		<InfoBanner v-else-if="!hasAvailableFields" variant="info" :title="t.steps.step8.noFieldsWarning.title.value" class="q-mb-md">
			{{ t.steps.step8.noFieldsWarning.description.value }}
		</InfoBanner>

		<!-- Add New Field Form -->
		<q-card flat bordered class="q-mb-md">
			<q-card-section>
				<div class="text-subtitle2 text-weight-medium q-mb-xs">
					<q-icon name="add" color="primary" class="q-mr-xs" />
					{{ t.steps.step8.addField.title.value }}
				</div>
				<div class="add-field-helper q-mb-md">
					{{ t.steps.step8.addField.helper.value }}
				</div>

				<div class="row q-col-gutter-md">
					<!-- Field Name -->
					<div class="col-12 col-sm-6" :class="isGeoType ? 'col-md-3' : 'col-md-4'">
						<div class="step-input-label q-mb-xs">
							<span class="step-badge">1</span>
							<span>{{ t.steps.step8.addField.fieldName.label.value }}</span>
						</div>
						<q-select
							ref="fieldNameSelect"
							:model-value="newField.field"
							outlined
							dense
							use-input
							fill-input
							hide-selected
							new-value-mode="add-unique"
							input-debounce="200"
							class="rounded-borders"
							:placeholder="t.steps.step8.addField.fieldName.placeholder.value"
							:hint="selectedFieldHint"
							:options="fieldNameOptions"
							option-value="value"
							option-label="value"
							option-disable="isHeader"
							emit-value
							map-options
							@filter="filterFieldName"
							@input-value="onFieldNameTyped"
							@update:model-value="onFieldNameSelected"
						>
							<template v-slot:prepend>
								<q-icon name="mdi-form-textbox" color="primary" />
							</template>
							<template v-slot:option="scope">
								<q-item-label
									v-if="scope.opt.isHeader"
									header
									class="text-weight-medium text-primary"
								>
									{{ scope.opt.value }}
								</q-item-label>
								<q-item v-else v-bind="scope.itemProps">
									<q-item-section>
										<q-item-label>{{ scope.opt.value }}</q-item-label>
										<q-item-label v-if="scope.opt.hint" caption>
											{{ scope.opt.hint }}
											<span v-if="scope.opt.unit"> ({{ scope.opt.unit }})</span>
										</q-item-label>
									</q-item-section>
									<q-item-section v-if="scope.opt.type" side>
										<q-icon :name="getTypeIcon(scope.opt.type)" color="primary" size="xs" />
									</q-item-section>
								</q-item>
							</template>
							<template v-slot:no-option>
								<q-item
									v-if="fieldNameSearch"
									clickable
									@click="commitCustomFieldName"
								>
									<q-item-section avatar>
										<q-icon name="add" color="primary" />
									</q-item-section>
									<q-item-section>
										{{ t.steps.step8.addField.fieldName.createCustomPrefix.value }}
										«{{ fieldNameSearch }}»
									</q-item-section>
								</q-item>
								<q-item v-else>
									<q-item-section class="text-grey">
										{{ t.steps.step8.addField.fieldName.customHint.value }}
									</q-item-section>
								</q-item>
							</template>
						</q-select>
					</div>

					<!-- Field Type -->
					<div class="col-12 col-sm-6 col-md-2">
						<div class="step-input-label q-mb-xs">
							<span class="step-badge">2</span>
							<span>{{ t.steps.step8.addField.fieldType.label.value }}</span>
						</div>
						<q-select
							v-model="newField.type"
							outlined
							dense
							class="rounded-borders"
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
					<div v-if="!isGeoType" class="col-12 col-sm-6 col-md-5">
						<div class="step-input-label q-mb-xs">
							<span class="step-badge">3</span>
							<span>{{ t.steps.step8.addField.valuePath.label.value }}</span>
						</div>
						<q-select
							v-model="newField.value"
							outlined
							dense
							use-input
							fill-input
							hide-selected
							input-debounce="0"
							class="rounded-borders"
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
							<template v-slot:append>
								<q-icon name="info" color="grey-6" size="xs">
									<AppTooltip :content="t.steps.step8.addField.valuePath.helper.value" />
								</q-icon>
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
						<div class="step-input-label q-mb-xs">
							<span class="step-badge">3</span>
							<span>{{ t.steps.step8.addField.latitudePath.label.value }}</span>
						</div>
						<q-select
							v-model="newField.latitudePath"
							outlined
							dense
							use-input
							fill-input
							hide-selected
							input-debounce="0"
							class="rounded-borders"
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
						<div class="step-input-label q-mb-xs">
							<span class="step-badge step-badge--continuation" aria-hidden="true">3</span>
							<span>{{ t.steps.step8.addField.longitudePath.label.value }}</span>
						</div>
						<q-select
							v-model="newField.longitudePath"
							outlined
							dense
							use-input
							fill-input
							hide-selected
							input-debounce="0"
							class="rounded-borders"
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
					<div class="col-12 col-sm-auto column">
						<div class="step-input-label q-mb-xs" aria-hidden="true">&nbsp;</div>
						<q-btn
							unelevated
							round
							:color="isAtLimit ? 'grey' : 'primary'"
							icon="add"
							:disable="!canAddField"
							@click="addField"
						>
							<AppTooltip v-if="isAtLimit" :content="t.steps.step8.limitReached.tooltip.value" />
							<AppTooltip v-else-if="addFieldDisabledReason" :content="addFieldDisabledReason" />
							<AppTooltip v-else :content="t.steps.step8.addField.addButton.value" />
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

.add-field-helper {
	color: var(--mapex-text-secondary);
	font-size: var(--mapex-font-sm);
	line-height: var(--mapex-line-height-base);
}

.step-input-label {
	display: flex;
	align-items: center;
	gap: var(--mapex-spacing-xs);
	color: var(--mapex-text-primary);
	font-size: var(--mapex-font-sm);
	font-weight: var(--mapex-font-weight-medium);
}

.step-badge {
	display: inline-flex;
	align-items: center;
	justify-content: center;
	width: 20px;
	height: 20px;
	border-radius: var(--mapex-radius-full);
	background: var(--mapex-primary);
	color: var(--mapex-wf-text-on-accent);
	font-size: var(--mapex-font-2xs);
	font-weight: var(--mapex-font-weight-bold);
	line-height: 1;
}

.step-badge--continuation {
	background: var(--mapex-neutral-soft);
	color: var(--mapex-text-secondary);
}




</style>
