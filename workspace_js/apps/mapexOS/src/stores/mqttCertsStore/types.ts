export interface RevokedCertRow {
	serial: string;
	fingerprint: string;
	assetUUID: string;
	orgId: string;
	subjectCN: string;
	issuedAt: string | Date;
	revokedAt: string | Date;
	reason: string;
}
