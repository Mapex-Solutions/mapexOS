/**
 * Browser-side zip-download helper for the device-cert issue flow.
 *
 * Bundles cert + key + ca-chain PEM byte arrays into a single zip blob
 * and triggers download via an invisible <a> click. Run-once: the
 * server returns the PEMs only at issue time and discards immediately,
 * so the operator MUST save the zip on this single response.
 *
 * Uses jszip — added to apps/mapexOS/package.json by T30.
 */
import JSZip from 'jszip';

export interface CertZipParams {
	filename: string;
	certPEM: Uint8Array | string;
	keyPEM: Uint8Array | string;
	caChainPEM: Uint8Array | string;
}

export async function downloadCertZip(params: CertZipParams): Promise<void> {
	const zip = new JSZip();
	zip.file('cert.pem', params.certPEM);
	zip.file('key.pem', params.keyPEM);
	zip.file('ca-chain.pem', params.caChainPEM);

	const blob = await zip.generateAsync({ type: 'blob' });
	const url = URL.createObjectURL(blob);
	try {
		const a = document.createElement('a');
		a.href = url;
		a.download = params.filename;
		a.style.display = 'none';
		document.body.appendChild(a);
		a.click();
		document.body.removeChild(a);
	} finally {
		setTimeout(() => URL.revokeObjectURL(url), 60_000);
	}
}

/**
 * Decode Go-style base64 ([]byte JSON) into a Uint8Array.
 */
export function decodeBase64ToBytes(b64: string): Uint8Array {
	const bin = atob(b64);
	const out = new Uint8Array(bin.length);
	for (let i = 0; i < bin.length; i++) out[i] = bin.charCodeAt(i);
	return out;
}
