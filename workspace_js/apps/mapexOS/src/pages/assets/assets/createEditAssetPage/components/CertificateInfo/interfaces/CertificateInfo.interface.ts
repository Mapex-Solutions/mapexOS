// Props/emits for the CertificateInfo component.
// AssetResponse comes from @mapexos/schemas via the existing asset schema.
// The asset.currentCert subdoc (when present) carries the metadata
// displayed by this component; nil = empty-state.

export interface CertificateInfoProps {
	asset: {
		uuid: string;
		currentCert?: {
			serial: string;
			fingerprint: string;
			subjectCN: string;
			issuedAt: string | Date;
			expiresAt: string | Date;
		} | null;
	};
}

export interface CertificateInfoEmits {
	(e: 'revoked'): void;
}
