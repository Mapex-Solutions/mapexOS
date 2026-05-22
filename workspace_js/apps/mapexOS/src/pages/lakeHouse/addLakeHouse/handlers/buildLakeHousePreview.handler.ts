import type { LakeHouseConfigProps } from '@components/forms/lakeHouse';

// Basic type definitions
export type IconDef = { name: string; color: string };
export type FieldDef = {
  label: string;
  value: any;
  type: 'text' | 'badge' | 'datetime';
  badgeColors?: string | Record<string, string>;
  format?: string;
  colSize?: number;
};

export type SectionDef = { stepNumber: number, label: string; icon: IconDef; fields: FieldDef[] };

// Builds the "General Information" section
export function buildGeneralSection(data: LakeHouseConfigProps): SectionDef {
  const fields: FieldDef[] = [
    {
      label: 'Type',
      value: data.type,
      type: 'badge',
    },
    { label: 'Name', value: data.name, type: 'text', colSize: 6 },
    {
      label: 'Status',
      value: data.status ? 'Active' : 'Inactive',
      type: 'badge',
      badgeColors: { Active: 'positive', Inactive: 'grey' },
      colSize: 6,
    },
    { label: 'Description', value: data.description, type: 'text', colSize: 6 },
    { label: 'Tenant ID', value: data.tenantId ?? '—', type: 'text', colSize: 6 },
  ];

  return {
    stepNumber: 1,
    label: 'General Information',
    icon: { name: 'info', color: 'primary' },
    fields,
  };
}

// Builds the "Credentials" section
export function buildCredentialsSection(data: LakeHouseConfigProps): SectionDef {
  const creds = data.credentials;
  const fields: FieldDef[] = [];

  if (creds.accessKey) {
    fields.push({ label: 'Access Key', value: creds.accessKey, type: 'text', colSize: 6 });
  }
  if (creds.secretKey) {
    fields.push({ label: 'Secret Key', value: creds.secretKey, type: 'text', colSize: 6 });
  }
  if (creds.bucket) {
    fields.push({ label: 'Bucket', value: creds.bucket, type: 'text', colSize: 6 });
  }
  if (creds.containerName) {
    fields.push({ label: 'Container Name', value: creds.containerName, type: 'text', colSize: 6 });
  }
  if (creds.region) {
    fields.push({ label: 'Region', value: creds.region, type: 'text', colSize: 6 });
  }
  if (creds.endpoint) {
    fields.push({ label: 'Endpoint', value: creds.endpoint, type: 'text', colSize: 6 });
  }
  if (typeof creds.useSSL === 'boolean') {
    fields.push({
      label: 'Use SSL',
      value: creds.useSSL ? 'Yes' : 'No',
      type: 'badge',
      badgeColors: { Yes: 'positive', No: 'negative' },
      colSize: 6
    });
  }
  if (creds.accountName) {
    fields.push({ label: 'Account Name', value: creds.accountName, type: 'text', colSize: 6 });
  }
  if (creds.accountKey) {
    fields.push({ label: 'Account Key', value: creds.accountKey, type: 'text', colSize: 6 });
  }
  if (creds.projectId) {
    fields.push({ label: 'Project ID', value: creds.projectId, type: 'text', colSize: 6 });
  }
  if (creds.keyFile) {
    fields.push({ label: 'Key File', value: creds.keyFile, type: 'text', colSize: 6 });
  }
  if (creds.tags) {
    fields.push({
      label: 'Tags',
      value: Object.entries(creds.tags).map(([k, v]) => `${k}: ${v}`).join(', '),
      type: 'text',
      colSize: 6
    });
  }

  return {
    stepNumber: 3,
    label: 'Credentials',
    icon: { name: 'key', color: 'secondary' },
    fields,
  };
}

// Builds the "Schedule" section
export function buildFrequencySection(data: LakeHouseConfigProps): SectionDef {
  const freq = data.frequency;
  const fields: FieldDef[] = [{ label: 'Type', value: freq.type, type: 'text', colSize: 6 }];

  if ('cron' in freq && freq.cron) {
    fields.push({ label: 'Cron', value: freq.cron, type: 'text', colSize: 6 });
  }

  return {
    stepNumber: 5,
    label: 'Schedule',
    icon: { name: 'schedule', color: 'primary' },
    fields,
  };
}

// Builds the "Path Configuration" section
export function buildPathConfigSection(data: LakeHouseConfigProps): SectionDef {
  const path = data.pathConfig;
  const partitions = path.partitions.map((p) => `${p}`);
  const fields: FieldDef[] = [
    { label: 'Base Path', value: path.basePath, type: 'text', colSize: 6 },
    { label: 'Partitions', value: partitions.join('/'), type: 'text', colSize: 6 },
  ];

  if (path.compression) {
    fields.push({
      label: 'Compression',
      value: path.compression,
      type: 'badge',
      colSize: 6
    });
  }
  if (path.maxFileSize) {
    fields.push({ label: 'Max File Size', value: `${path.maxFileSize} MB`, type: 'text', colSize: 6 });
  }
  if (path.filePrefix) {
    fields.push({ label: 'File Prefix', value: path.filePrefix, type: 'text', colSize: 6 });
  }

  return {
    stepNumber: 4,
    label: 'Path Configuration',
    icon: { name: 'folder', color: 'primary' },
    fields,
  };
}

// Main builder that composes all sections
export function buildLakeHousePreview(
  data: LakeHouseConfigProps
): any[] {
  return [
    buildGeneralSection(data),
    buildCredentialsSection(data),
    buildPathConfigSection(data),
    buildFrequencySection(data),
  ];
}
