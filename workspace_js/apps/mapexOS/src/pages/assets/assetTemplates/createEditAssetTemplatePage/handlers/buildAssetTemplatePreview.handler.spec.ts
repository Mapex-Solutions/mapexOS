import { describe, it, expect } from 'vitest';
import { buildAssetTemplatePreview } from './buildAssetTemplatePreview.handler';
import type { AssetTemplateData } from '../interfaces';

/**
 * Create a minimal AssetTemplateData for testing
 */
function createMockData(overrides: Partial<AssetTemplateData> = {}): AssetTemplateData {
  return {
    name: 'Temperature Sensor',
    enabled: true,
    description: 'Reads temperature',
    categoryName: 'Sensors',
    manufacturerName: 'Acme',
    modelName: 'T-100',
    version: '1.0.0',
    isSystem: false,
    isTemplate: false,
    assetIdPath: 'data.deviceId',
    scriptProcessor: 'function decode(payload) {}',
    scriptValidator: 'function validate(payload) {}',
    scriptConversion: 'function transform(payload) {}',
    scriptTest: '{"deviceId": "123"}',
    ...overrides,
  };
}

describe('buildAssetTemplatePreview', () => {
  it('returns exactly 3 sections', () => {
    const result = buildAssetTemplatePreview(createMockData());
    expect(result).toHaveLength(3);
  });

  it('returns sections in correct order with step numbers', () => {
    const result = buildAssetTemplatePreview(createMockData());

    expect(result[0]!.stepNumber).toBe(1);
    expect(result[0]!.label).toBe('Basic Information');
    expect(result[1]!.stepNumber).toBe(2);
    expect(result[1]!.label).toBe('Asset ID Path');
    expect(result[2]!.stepNumber).toBe(3);
    expect(result[2]!.label).toBe('Scripts Summary');
  });

  describe('Basic Information section', () => {
    it('includes name field', () => {
      const result = buildAssetTemplatePreview(createMockData({ name: 'My Sensor' }));
      const section = result[0]!;
      const nameField = section.fields.find((f) => f.label === 'Name');

      expect(nameField).toBeDefined();
      expect(nameField!.value).toBe('My Sensor');
      expect(nameField!.type).toBe('text');
    });

    it('shows Active badge when enabled', () => {
      const result = buildAssetTemplatePreview(createMockData({ enabled: true }));
      const statusField = result[0]!.fields.find((f) => f.label === 'Status');

      expect(statusField!.value).toBe('Active');
      expect(statusField!.type).toBe('badge');
    });

    it('shows Inactive badge when disabled', () => {
      const result = buildAssetTemplatePreview(createMockData({ enabled: false }));
      const statusField = result[0]!.fields.find((f) => f.label === 'Status');

      expect(statusField!.value).toBe('Inactive');
    });

    it('shows fallback for empty description', () => {
      const result = buildAssetTemplatePreview(createMockData({ description: undefined }));
      const descField = result[0]!.fields.find((f) => f.label === 'Description');

      expect(descField!.value).toBe('No description provided');
    });

    it('shows fallback for missing category, manufacturer, model, version', () => {
      const result = buildAssetTemplatePreview(createMockData({
        categoryName: undefined,
        manufacturerName: undefined,
        modelName: undefined,
        version: undefined,
      }));
      const section = result[0]!;

      expect(section.fields.find((f) => f.label === 'Category')!.value).toBe('Not specified');
      expect(section.fields.find((f) => f.label === 'Manufacturer')!.value).toBe('Not specified');
      expect(section.fields.find((f) => f.label === 'Model')!.value).toBe('Not specified');
      expect(section.fields.find((f) => f.label === 'Version')!.value).toBe('Not specified');
    });

    it('has correct icon', () => {
      const result = buildAssetTemplatePreview(createMockData());
      expect(result[0]!.icon).toEqual({ name: 'info', color: 'primary' });
    });
  });

  describe('Asset ID Path section', () => {
    it('includes assetIdPath field', () => {
      const result = buildAssetTemplatePreview(createMockData({ assetIdPath: 'payload.id' }));
      const section = result[1]!;

      expect(section.fields).toHaveLength(1);
      expect(section.fields[0]!.value).toBe('payload.id');
      expect(section.fields[0]!.colSize).toBe(12);
    });

    it('has correct icon', () => {
      const result = buildAssetTemplatePreview(createMockData());
      expect(result[1]!.icon).toEqual({ name: 'route', color: 'secondary' });
    });
  });

  describe('Scripts Summary section', () => {
    it('shows Configured for scripts that have content', () => {
      const result = buildAssetTemplatePreview(createMockData({
        scriptProcessor: 'code',
        scriptValidator: 'code',
        scriptConversion: 'code',
        scriptTest: 'code',
      }));
      const section = result[2]!;

      expect(section.fields.find((f) => f.label === 'Preprocessor Script')!.value).toBe('Configured');
      expect(section.fields.find((f) => f.label === 'Validation Script')!.value).toBe('Configured');
      expect(section.fields.find((f) => f.label === 'Conversion Script')!.value).toBe('Configured');
      expect(section.fields.find((f) => f.label === 'Test Script')!.value).toBe('Configured');
    });

    it('shows Not configured for empty/missing scripts', () => {
      const result = buildAssetTemplatePreview(createMockData({
        scriptProcessor: '',
        scriptValidator: '',
        scriptConversion: '',
        scriptTest: '',
      }));
      const section = result[2]!;

      expect(section.fields.find((f) => f.label === 'Preprocessor Script')!.value).toBe('Not configured');
      expect(section.fields.find((f) => f.label === 'Validation Script')!.value).toBe('Not configured');
      expect(section.fields.find((f) => f.label === 'Conversion Script')!.value).toBe('Not configured');
      expect(section.fields.find((f) => f.label === 'Test Script')!.value).toBe('Not configured');
    });

    it('has badge type for all script fields', () => {
      const result = buildAssetTemplatePreview(createMockData());
      const section = result[2]!;

      section.fields.forEach((field) => {
        expect(field.type).toBe('badge');
      });
    });

    it('has correct badge colors for configured scripts', () => {
      const result = buildAssetTemplatePreview(createMockData());
      const section = result[2]!;

      section.fields.forEach((field) => {
        expect(field.badgeColors).toHaveProperty('Configured', 'positive');
      });
    });

    it('has correct icon', () => {
      const result = buildAssetTemplatePreview(createMockData());
      expect(result[2]!.icon).toEqual({ name: 'code', color: 'primary' });
    });
  });
});
