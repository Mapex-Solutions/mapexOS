<script setup lang="ts">
/** TYPE IMPORTS */
import type { RevokedCertificatesListProps } from './interfaces/RevokedCertificatesList.interface';

/** VUE IMPORTS */
import { computed } from 'vue';

/** COMPOSABLES */
import { useAddAssetTranslations } from '@src/composables/i18n/pages/assets/addAsset/useAddAssetTranslations';

/** PROPS */
const props = defineProps<RevokedCertificatesListProps>();

/** COMPOSABLES */
const t = useAddAssetTranslations();

/** COMPUTED */
const isEmpty = computed(() => !props.loading && props.rows.length === 0);
const columns = [
	{ name: 'serial', label: 'Serial', field: 'serial', align: 'left' as const },
	{ name: 'reason', label: 'Reason', field: 'reason', align: 'left' as const },
	{ name: 'revokedAt', label: 'Revoked at', field: 'revokedAt', align: 'left' as const },
];
</script>

<template>
	<q-card class="revoked-list">
		<q-card-section>
			<div class="text-h6">{{ t.steps.step4.certificate.revokedTitle.value }}</div>
			<div class="retention-notice">
				Revoked certificates are retained for 30 days for audit. After that, audit data moves to the long-term archive (future).
			</div>
		</q-card-section>
		<q-card-section v-if="isEmpty" class="empty">{{ t.steps.step4.certificate.revokedEmpty.value }}</q-card-section>
		<q-table
			v-else
			:rows="rows"
			:columns="columns"
			:loading="loading"
			flat dense row-key="serial"
		/>
	</q-card>
</template>

<style scoped lang="scss">
.revoked-list { background: var(--mapex-surface-primary); border-radius: var(--mapex-radius-md); }
.retention-notice {
	margin-top: var(--mapex-space-sm);
	padding: var(--mapex-space-sm);
	background: var(--mapex-surface-info-soft);
	border-radius: var(--mapex-radius-sm);
	color: var(--mapex-text-info);
	font-size: 0.9em;
}
.empty { color: var(--mapex-text-secondary); font-style: italic; }
</style>
