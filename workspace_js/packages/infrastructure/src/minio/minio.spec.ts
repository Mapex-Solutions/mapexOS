import { describe, it, expect } from 'vitest';
import { validateConfig } from './minio';
import type { Config } from './types';

/**
 * Mirrors the Go kit's minio_test.go: the credential keys are required only for
 * static auth (AuthIsNeeded=true); with ambient IAM (false) they are optional.
 */
describe('validateConfig — static auth vs ambient IAM', () => {
	const base: Omit<Config, 'AuthIsNeeded' | 'AccessKeyID' | 'SecretAccessKey'> = {
		Endpoint: 'localhost:9000',
		BucketName: 'mapex-assets',
	};

	it('valid with static auth + keys', () => {
		expect(
			validateConfig({ ...base, AuthIsNeeded: true, AccessKeyID: 'admin', SecretAccessKey: 'secret' }),
		).toBeNull();
	});

	it('ambient IAM needs no keys', () => {
		expect(
			validateConfig({ ...base, AuthIsNeeded: false, AccessKeyID: '', SecretAccessKey: '' }),
		).toBeNull();
	});

	it('static auth missing access key -> error', () => {
		expect(
			validateConfig({ ...base, AuthIsNeeded: true, AccessKeyID: '', SecretAccessKey: 'secret' }),
		).toMatch(/access key/i);
	});

	it('static auth missing secret key -> error', () => {
		expect(
			validateConfig({ ...base, AuthIsNeeded: true, AccessKeyID: 'admin', SecretAccessKey: '' }),
		).toMatch(/secret/i);
	});

	it('missing endpoint -> error regardless of auth mode', () => {
		expect(
			validateConfig({ ...base, Endpoint: '', AuthIsNeeded: false, AccessKeyID: '', SecretAccessKey: '' }),
		).toMatch(/endpoint/i);
	});
});
