<script setup lang="ts">
/** TYPE IMPORTS */
import type {
	GenerateCertificateDialogProps,
	GenerateCertificateDialogEmits,
} from './interfaces/GenerateCertificateDialog.interface';

/** VUE IMPORTS */
import { ref } from 'vue';

/** PROPS & EMITS */
const props = defineProps<GenerateCertificateDialogProps>();
const emit = defineEmits<GenerateCertificateDialogEmits>();

/** STATE */
const issuing = ref(false);

/** FUNCTIONS */

/**
 * Trigger the issue call. The parent owns the actual API call + zip
 * download via the @issued emit — this component focuses on the
 * dialog UX (confirm + warning copy) and emits up. The dialog is
 * closed by the parent through the v-model: when the parent has
 * successfully issued (or wants to keep the dialog open on error)
 * it controls visibility, so we only emit `issued` here.
 */
function onConfirm(): void {
	emit('issued');
}

function onClose(): void {
	emit('update:show', false);
}

// Surface props.* in the template through this name to silence the
// no-unused-vars warning when defineProps return value isn't bound.
const labels = props.labels;
</script>

<template>
	<q-dialog :model-value="show" persistent @update:model-value="onClose">
		<q-card class="generate-cert-dialog">
			<q-card-section>
				<div class="text-h6">{{ labels.title }}</div>
			</q-card-section>

			<q-card-section>
				<p class="warning">{{ labels.warning }}</p>
				<p v-if="hasExistingCert" class="warning">{{ labels.replaceWarning }}</p>
			</q-card-section>

			<q-card-actions align="right">
				<q-btn flat :label="labels.skipButton" :disable="issuing" @click="onClose" />
				<q-btn
					:label="labels.generateButton"
					color="primary"
					:loading="issuing"
					@click="onConfirm"
				/>
			</q-card-actions>
		</q-card>
	</q-dialog>
</template>

<style scoped lang="scss">
.generate-cert-dialog {
	min-width: 480px;
	background: var(--mapex-surface-primary);
}
.warning {
	color: var(--mapex-text-warning);
	padding: var(--mapex-space-md);
	background: var(--mapex-surface-warning-soft);
	border-radius: var(--mapex-radius-sm);
	line-height: 1.5;
	margin-bottom: var(--mapex-space-sm);
}
</style>
