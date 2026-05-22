<template>
	<q-item>
		<q-item-section>
			<q-item-label>{{ displayName }}</q-item-label>
			<q-item-label caption>{{ props.policy.type }}</q-item-label>
		</q-item-section>
		<q-item-section side style="width: 200px">
			<q-input
				v-model.number="retentionDays"
				type="number"
				:min="minDays"
				:max="maxDays"
				:suffix="t.columns.retentionDays.value"
				dense
				outlined
				:error="hasError"
				:error-message="errorMessage"
			/>
		</q-item-section>
		<q-item-section side>
			<q-btn
				color="primary"
				dense
				flat
				:label="t.actions.save.value"
				:disable="hasError || retentionDays === props.policy.retentionDays"
				@click="onSave"
			/>
		</q-item-section>
	</q-item>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue';
import type { RetentionPolicyResponse } from '@mapexos/schemas';
import { useRetentionPoliciesTranslations } from '@composables/i18n/pages/administrations/retentionPolicies/useRetentionPoliciesTranslations';

interface Props {
	policy: RetentionPolicyResponse;
}

const props = defineProps<Props>();

const emit = defineEmits<{
	(
		e: 'save',
		payload: { type: string; name: string; retentionDays: number },
	): void;
}>();

const t = useRetentionPoliciesTranslations();

const retentionDays = ref<number>(props.policy.retentionDays ?? 0);

watch(
	() => props.policy.retentionDays,
	(v) => {
		retentionDays.value = v ?? 0;
	},
);

// Range enforcement. asset_status_history is the only type that the current
// UI actively edits with the 7–90 constraint; other types fall back to a
// permissive range and rely on the backend's ValidateRetentionPolicy.
const minDays = computed(() =>
	props.policy.type === 'asset_status_history' ? 1 : 1,
);
const maxDays = computed(() =>
	props.policy.type === 'asset_status_history' ? 90 : 3650,
);

const hasError = computed(() => {
	return retentionDays.value < minDays.value || retentionDays.value > maxDays.value;
});

const errorMessage = computed(() =>
	hasError.value ? t.validation.ttlRange(minDays.value, maxDays.value) : '',
);

const displayName = computed(() => {
	if (props.policy.type === 'asset_status_history') {
		return t.specialTypes.assetStatusHistory.value;
	}
	return props.policy.name ?? props.policy.type ?? '';
});

function onSave() {
	if (hasError.value) return;
	emit('save', {
		type: props.policy.type ?? '',
		name: props.policy.name ?? displayName.value,
		retentionDays: retentionDays.value,
	});
}
</script>
