import type { z } from 'zod';
import type {
	IssueCertRequestSchema,
	IssueCertResponseSchema,
	RevokedCertResponseSchema,
	ListRevokedQuerySchema,
} from '@/assets/schemas/mqttcerts/mqttcerts.schema';

export type IssueCertRequest    = z.infer<typeof IssueCertRequestSchema>;
export type IssueCertResponse   = z.infer<typeof IssueCertResponseSchema>;
export type RevokedCertResponse = z.infer<typeof RevokedCertResponseSchema>;
export type ListRevokedQuery    = z.infer<typeof ListRevokedQuerySchema>;
