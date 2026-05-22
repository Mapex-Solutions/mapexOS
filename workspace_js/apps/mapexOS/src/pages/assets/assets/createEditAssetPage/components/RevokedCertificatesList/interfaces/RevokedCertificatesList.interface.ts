export interface RevokedCertificatesListProps {
	assetUuid: string;
	rows: ReadonlyArray<{
		serial: string;
		fingerprint: string;
		revokedAt: string | Date;
		reason: string;
	}>;
	loading: boolean;
}
