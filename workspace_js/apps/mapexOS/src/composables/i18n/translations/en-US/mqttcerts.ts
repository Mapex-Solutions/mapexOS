export const mqttcerts = {
	tab: { title: 'Certificates' },
	current: {
		title: 'Active Certificate',
		empty: 'No active certificate. Click Generate to issue a new one.',
		serial: 'Serial',
		fingerprint: 'Fingerprint',
		subjectCN: 'Subject CN',
		issuedAt: 'Issued',
		expiresAt: 'Expires',
	},
	actions: {
		generate: 'Generate certificate',
		revoke: 'Revoke',
	},
	dialog: {
		generate: {
			title: 'Generate MQTT certificate',
			warning:
				'The private key is shown only once. Download the zip and store it securely; the platform never persists it.',
			confirmButton: 'Generate and download',
		},
		replace: {
			title: 'Replace existing certificate?',
			body: 'This asset already has an active certificate. Generating a new one revokes the existing certificate immediately.',
			confirmButton: 'Replace',
		},
		revoke: {
			title: 'Revoke certificate?',
			body: 'The asset will lose MQTT access until a new certificate is issued.',
			confirmButton: 'Revoke',
		},
	},
	revoked: {
		title: 'Revoked certificates',
		retentionNotice:
			'Revoked certificates are retained for 30 days for audit. After that, audit data moves to the long-term archive (future).',
		empty: 'No revoked certificates for this asset.',
		columns: {
			serial: 'Serial',
			reason: 'Reason',
			revokedAt: 'Revoked at',
		},
	},
	errors: {
		caNotReady: 'PKI subsystem not ready. Try again shortly.',
		replaceRequired: 'Confirmation required to replace the existing certificate.',
		generic: 'Operation failed. See logs for details.',
	},
	success: {
		issued: 'Certificate issued. Download started.',
		revoked: 'Certificate revoked.',
	},
};
