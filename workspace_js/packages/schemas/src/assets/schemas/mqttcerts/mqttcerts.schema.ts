import { z } from 'zod';

const Base64Bytes = z.string();

export const IssueCertRequestSchema = z.object({
	assetUUID: z.string().min(1),
	force: z.boolean().optional().default(false),
});

export const IssueCertResponseSchema = z.object({
	serial: z.string(),
	fingerprint: z.string(),
	subjectCN: z.string(),
	issuedAt: z.coerce.date(),
	expiresAt: z.coerce.date(),
	certPEM: Base64Bytes,
	keyPEM: Base64Bytes,
	caChainPEM: Base64Bytes,
});

export const RevokedCertResponseSchema = z.object({
	serial: z.string(),
	fingerprint: z.string(),
	assetUUID: z.string(),
	orgId: z.string(),
	subjectCN: z.string(),
	issuedAt: z.coerce.date(),
	revokedAt: z.coerce.date(),
	reason: z.string(),
});

export const ListRevokedQuerySchema = z.object({
	assetUUID: z.string().min(1),
});
