<template>
	<q-page padding>
		<div class="q-mb-md">
			<div class="text-h5">{{ t.page.title.value }}</div>
			<div class="text-caption text-grey-7">{{ t.page.description.value }}</div>
		</div>

		<q-card flat bordered>
			<q-card-section>
				<div v-if="!loading && items.length === 0" class="text-center q-pa-xl">
					<div class="text-h6">{{ t.empty.title.value }}</div>
					<div class="text-caption text-grey-7 q-mt-sm">
						{{ t.empty.description.value }}
					</div>
				</div>
				<q-list v-else separator>
					<RetentionRowEditor
						v-for="policy in items"
						:key="policy.id ?? policy.type ?? ''"
						:policy="policy"
						@save="onRowSave"
					/>
				</q-list>
				<div v-if="loading" class="row justify-center q-pa-md">
					<q-spinner />
				</div>
			</q-card-section>
		</q-card>

		<q-banner v-if="error" class="bg-red-2 text-red-9 q-mt-md" rounded>
			{{ error }}
		</q-banner>
	</q-page>
</template>

<script setup lang="ts">
import { onMounted } from 'vue';
import { useQuasar } from 'quasar';
import RetentionRowEditor from './components/RetentionRowEditor.vue';
import { useRetentionPolicies } from './composables/useRetentionPolicies';
import { useRetentionPoliciesTranslations } from '@composables/i18n/pages/administrations/retentionPolicies/useRetentionPoliciesTranslations';

const $q = useQuasar();
const t = useRetentionPoliciesTranslations();

const { items, loading, error, fetch, update } = useRetentionPolicies();

async function onRowSave(payload: { type: string; name: string; retentionDays: number }) {
	try {
		await update(payload.type, payload.name, payload.retentionDays);
		$q.notify({ type: 'positive', message: t.messages.saveSuccess.value });
	} catch {
		$q.notify({ type: 'negative', message: t.messages.saveError.value });
	}
}

onMounted(async () => {
	await fetch();
});
</script>
