<script setup lang="ts">
/** TYPE IMPORTS */
import type { CertificateInfoProps, CertificateInfoEmits } from './interfaces/CertificateInfo.interface';

/** VUE IMPORTS */
import { computed, ref } from 'vue';

/** COMPOSABLES */
import { useAddAssetTranslations } from '@src/composables/i18n/pages/assets/addAsset/useAddAssetTranslations';

/** PROPS & EMITS */
const props = defineProps<CertificateInfoProps>();
const emit = defineEmits<CertificateInfoEmits>();

/** COMPOSABLES */
const t = useAddAssetTranslations();

/** STATE */
const revoking = ref(false);

/** COMPUTED */
const hasCert = computed(() => props.asset.currentCert != null);

/** FUNCTIONS */

/**
 * Surface a Revoke confirmation dialog and call the API on confirm.
 * The actual revoke call is wired by the parent (Step4Connectivity)
 * via the @revoked emit + the parent's API plumbing.
 */
function onRevoke(): void {
	if (!props.asset.currentCert) return;
	if (window.confirm(t.steps.step4.certificate.revokeConfirm.value)) {
		revoking.value = true;
		emit('revoked');
		revoking.value = false;
	}
}
</script>

<template>
	<q-card class="certificate-info">
		<q-card-section>
			<div class="text-h6">{{ t.steps.step4.certificate.activeTitle.value }}</div>
		</q-card-section>

		<q-card-section v-if="!hasCert" class="empty">
			{{ t.steps.step4.certificate.noActiveCert.value }}
		</q-card-section>

		<q-card-section v-else class="meta">
			<div class="row"><span class="label">{{ t.steps.step4.certificate.fieldSerial.value }}</span><code>{{ asset.currentCert!.serial }}</code></div>
			<div class="row"><span class="label">{{ t.steps.step4.certificate.fieldFingerprint.value }}</span><code class="fp">{{ asset.currentCert!.fingerprint }}</code></div>
			<div class="row"><span class="label">{{ t.steps.step4.certificate.fieldSubjectCN.value }}</span><span>{{ asset.currentCert!.subjectCN }}</span></div>
			<div class="row"><span class="label">{{ t.steps.step4.certificate.fieldIssued.value }}</span><span>{{ asset.currentCert!.issuedAt }}</span></div>
			<div class="row"><span class="label">{{ t.steps.step4.certificate.fieldExpires.value }}</span><span>{{ asset.currentCert!.expiresAt }}</span></div>
		</q-card-section>

		<q-card-actions v-if="hasCert" align="right">
			<q-btn :label="t.steps.step4.certificate.revokeButton.value" color="negative" flat :loading="revoking" @click="onRevoke" />
		</q-card-actions>
	</q-card>
</template>

<style scoped lang="scss">
.certificate-info {
	background: var(--mapex-surface-primary);
	border-radius: var(--mapex-radius-md);
	box-shadow: var(--mapex-shadow-sm);
}
.empty {
	color: var(--mapex-text-secondary);
	font-style: italic;
}
.meta .row {
	display: flex;
	gap: var(--mapex-space-sm);
	padding: var(--mapex-space-xs) 0;
	align-items: baseline;
	.label { color: var(--mapex-text-secondary); min-width: 120px; }
	code { background: var(--mapex-surface-secondary); padding: 2px 6px; border-radius: var(--mapex-radius-sm); }
	.fp { font-size: 0.85em; word-break: break-all; }
}
</style>
