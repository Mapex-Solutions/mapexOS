/**
 * GenerateCertificateDialog is a pure presentation component — labels
 * are injected by the parent so the dialog can be reused across
 * contexts (post-save wizard, asset details drawer, list row-menu)
 * without coupling to a single i18n composable. Each label is a
 * pre-resolved string at render time.
 */
export interface GenerateCertificateDialogLabels {
	title: string;
	warning: string;
	replaceWarning: string;
	generateButton: string;
	skipButton: string;
}

export interface GenerateCertificateDialogProps {
	show: boolean;
	assetUuid: string;
	assetName: string;
	hasExistingCert: boolean;
	labels: GenerateCertificateDialogLabels;
}

export interface GenerateCertificateDialogEmits {
	(e: 'update:show', value: boolean): void;
	(e: 'issued'): void;
}
