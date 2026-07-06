import axios from 'axios';

/**
 * Computes the SHA-256 of a firmware blob via WebCrypto and returns both
 * encodings the OTA flow needs: hex for the init request (the platform
 * validates the stored object against it on finalize) and base64 for the
 * presigned PUT's x-amz-checksum-sha256 header (the object store validates
 * the bytes on write).
 *
 * Browser-only (WebCrypto `crypto.subtle`); firmware upload is a UI flow.
 *
 * @param data - The firmware bytes (File/Blob from the picker or a raw buffer)
 * @returns The digest as { hex, base64 }
 */
export async function sha256OfFirmware(data: Blob | ArrayBuffer): Promise<{ hex: string; base64: string }> {
	const buffer = data instanceof Blob ? await data.arrayBuffer() : data;
	const digest = await crypto.subtle.digest('SHA-256', buffer);
	const bytes = new Uint8Array(digest);

	const hex = Array.from(bytes)
		.map((b) => b.toString(16).padStart(2, '0'))
		.join('');
	const base64 = btoa(String.fromCharCode(...bytes));
	return { hex, base64 };
}

/**
 * PUTs the firmware .bin straight to object storage via the presigned URL
 * returned by `firmwareInit` — the bytes never traverse the platform API.
 *
 * The presigned URL was signed WITH the x-amz-checksum-sha256 header, so the
 * store rejects the write when the uploaded bytes do not match the declared
 * checksum (base64 of the raw SHA-256 digest — use `sha256OfFirmware`).
 *
 * Uses a bare axios call on purpose: the URL is absolute, pre-authorized, and
 * must NOT carry the platform JWT or interceptors.
 *
 * @param uploadUrl - The presigned PUT URL from firmwareInit
 * @param data - The firmware bytes
 * @param checksumBase64 - Base64-encoded SHA-256 of the bytes
 */
export async function uploadFirmwareBinary(uploadUrl: string, data: Blob | ArrayBuffer, checksumBase64: string): Promise<void> {
	await axios.put(uploadUrl, data, {
		headers: {
			'Content-Type': 'application/octet-stream',
			'x-amz-checksum-sha256': checksumBase64,
		},
	});
}
