import { describe, it, expect } from 'vitest';
import type { LakeHouseConfigProps } from '@components/forms/lakeHouse';
import {
  buildGeneralSection,
  buildCredentialsSection,
  buildFrequencySection,
  buildPathConfigSection,
  buildLakeHousePreview,
} from './buildLakeHousePreview.handler';

function makeData(overrides: Partial<LakeHouseConfigProps> = {}): LakeHouseConfigProps {
  return {
    name: 'Test Lake',
    description: 'A test lake house',
    type: 'aws-s3',
    status: true,
    credentials: {
      accessKey: 'AKIA...',
      secretKey: 'secret',
      bucket: 'my-bucket',
      region: 'us-east-1',
    },
    frequency: {
      type: 'hour',
    },
    pathConfig: {
      basePath: '/data/raw',
      partitions: ['year', 'month', 'day'],
    },
    ...overrides,
  };
}

describe('buildLakeHousePreview.handler', () => {
  // ── buildGeneralSection ──────────────────────────────────────────────
  describe('buildGeneralSection', () => {
    it('should return general info section with stepNumber 1', () => {
      const section = buildGeneralSection(makeData());

      expect(section.stepNumber).toBe(1);
      expect(section.label).toBe('General Information');
      expect(section.icon.name).toBe('info');
    });

    it('should include name, type, status, and description fields', () => {
      const section = buildGeneralSection(makeData({ name: 'My Lake', status: false }));
      const fieldLabels = section.fields.map((f) => f.label);

      expect(fieldLabels).toContain('Name');
      expect(fieldLabels).toContain('Type');
      expect(fieldLabels).toContain('Status');
      expect(fieldLabels).toContain('Description');

      const statusField = section.fields.find((f) => f.label === 'Status');
      expect(statusField?.value).toBe('Inactive');

      const nameField = section.fields.find((f) => f.label === 'Name');
      expect(nameField?.value).toBe('My Lake');
    });

    it('should show tenantId as dash when not provided', () => {
      const section = buildGeneralSection(makeData({}));
      const tenantField = section.fields.find((f) => f.label === 'Tenant ID');

      expect(tenantField?.value).toBe('\u2014');
    });

    it('should show tenantId when provided', () => {
      const section = buildGeneralSection(makeData({ tenantId: 'tenant-123' }));
      const tenantField = section.fields.find((f) => f.label === 'Tenant ID');

      expect(tenantField?.value).toBe('tenant-123');
    });
  });

  // ── buildCredentialsSection ──────────────────────────────────────────
  describe('buildCredentialsSection', () => {
    it('should return credentials section with stepNumber 3', () => {
      const section = buildCredentialsSection(makeData());

      expect(section.stepNumber).toBe(3);
      expect(section.label).toBe('Credentials');
    });

    it('should include AWS S3 credential fields', () => {
      const section = buildCredentialsSection(makeData());
      const fieldLabels = section.fields.map((f) => f.label);

      expect(fieldLabels).toContain('Access Key');
      expect(fieldLabels).toContain('Secret Key');
      expect(fieldLabels).toContain('Bucket');
      expect(fieldLabels).toContain('Region');
    });

    it('should include Azure fields when provided', () => {
      const section = buildCredentialsSection(makeData({
        credentials: {
          accountName: 'myaccount',
          accountKey: 'key123',
          containerName: 'container1',
        },
      }));
      const fieldLabels = section.fields.map((f) => f.label);

      expect(fieldLabels).toContain('Account Name');
      expect(fieldLabels).toContain('Account Key');
      expect(fieldLabels).toContain('Container Name');
    });

    it('should include GCP fields when provided', () => {
      const section = buildCredentialsSection(makeData({
        credentials: {
          projectId: 'my-project',
          keyFile: '/path/to/key.json',
        },
      }));
      const fieldLabels = section.fields.map((f) => f.label);

      expect(fieldLabels).toContain('Project ID');
      expect(fieldLabels).toContain('Key File');
    });

    it('should include useSSL badge when provided', () => {
      const section = buildCredentialsSection(makeData({
        credentials: {
          accessKey: 'key',
          useSSL: true,
        },
      }));

      const sslField = section.fields.find((f) => f.label === 'Use SSL');
      expect(sslField?.value).toBe('Yes');
      expect(sslField?.type).toBe('badge');
    });

    it('should include endpoint and tags when provided', () => {
      const section = buildCredentialsSection(makeData({
        credentials: {
          endpoint: 'https://minio.local:9000',
          tags: { env: 'prod', team: 'data' },
        },
      }));

      const endpointField = section.fields.find((f) => f.label === 'Endpoint');
      expect(endpointField?.value).toBe('https://minio.local:9000');

      const tagsField = section.fields.find((f) => f.label === 'Tags');
      expect(tagsField?.value).toContain('env: prod');
      expect(tagsField?.value).toContain('team: data');
    });

    it('should return empty fields for minimal credentials', () => {
      const section = buildCredentialsSection(makeData({
        credentials: {},
      }));

      expect(section.fields).toHaveLength(0);
    });
  });

  // ── buildFrequencySection ────────────────────────────────────────────
  describe('buildFrequencySection', () => {
    it('should return schedule section with stepNumber 5', () => {
      const section = buildFrequencySection(makeData());

      expect(section.stepNumber).toBe(5);
      expect(section.label).toBe('Schedule');
    });

    it('should include type field', () => {
      const section = buildFrequencySection(makeData({ frequency: { type: 'day' } }));
      const typeField = section.fields.find((f) => f.label === 'Type');

      expect(typeField?.value).toBe('day');
    });

    it('should include cron field when provided', () => {
      const section = buildFrequencySection(makeData({
        frequency: { type: 'minute', cron: '*/5 * * * *' } as any,
      }));

      const cronField = section.fields.find((f) => f.label === 'Cron');
      expect(cronField?.value).toBe('*/5 * * * *');
    });

    it('should not include cron field when not provided', () => {
      const section = buildFrequencySection(makeData({ frequency: { type: 'hour' } }));
      const cronField = section.fields.find((f) => f.label === 'Cron');

      expect(cronField).toBeUndefined();
    });
  });

  // ── buildPathConfigSection ───────────────────────────────────────────
  describe('buildPathConfigSection', () => {
    it('should return path config section with stepNumber 4', () => {
      const section = buildPathConfigSection(makeData());

      expect(section.stepNumber).toBe(4);
      expect(section.label).toBe('Path Configuration');
    });

    it('should include basePath and partitions', () => {
      const section = buildPathConfigSection(makeData({
        pathConfig: {
          basePath: '/data/raw',
          partitions: ['year', 'month'],
        },
      }));

      const basePathField = section.fields.find((f) => f.label === 'Base Path');
      expect(basePathField?.value).toBe('/data/raw');

      const partitionsField = section.fields.find((f) => f.label === 'Partitions');
      expect(partitionsField?.value).toBe('year/month');
    });

    it('should include compression when provided', () => {
      const section = buildPathConfigSection(makeData({
        pathConfig: {
          basePath: '/data',
          partitions: [],
          compression: 'gzip',
        },
      }));

      const compressionField = section.fields.find((f) => f.label === 'Compression');
      expect(compressionField?.value).toBe('gzip');
      expect(compressionField?.type).toBe('badge');
    });

    it('should include maxFileSize with MB suffix', () => {
      const section = buildPathConfigSection(makeData({
        pathConfig: {
          basePath: '/data',
          partitions: [],
          maxFileSize: 256,
        },
      }));

      const maxSizeField = section.fields.find((f) => f.label === 'Max File Size');
      expect(maxSizeField?.value).toBe('256 MB');
    });

    it('should include filePrefix when provided', () => {
      const section = buildPathConfigSection(makeData({
        pathConfig: {
          basePath: '/data',
          partitions: [],
          filePrefix: 'events_',
        },
      }));

      const prefixField = section.fields.find((f) => f.label === 'File Prefix');
      expect(prefixField?.value).toBe('events_');
    });
  });

  // ── buildLakeHousePreview (main composer) ────────────────────────────
  describe('buildLakeHousePreview', () => {
    it('should return 4 sections in correct order', () => {
      const sections = buildLakeHousePreview(makeData());

      expect(sections).toHaveLength(4);
      expect(sections[0].label).toBe('General Information');
      expect(sections[1].label).toBe('Credentials');
      expect(sections[2].label).toBe('Path Configuration');
      expect(sections[3].label).toBe('Schedule');
    });

    it('should return step numbers 1, 3, 4, 5', () => {
      const sections = buildLakeHousePreview(makeData());
      const stepNumbers = sections.map((s: any) => s.stepNumber);

      expect(stepNumbers).toEqual([1, 3, 4, 5]);
    });
  });
});
