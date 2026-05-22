import { describe, it, expect } from 'vitest';
import { buildDefaultConfig } from './buildDefaultConfig';
import type { PluginNodeType } from '@src/components/workflow/interfaces';

/**
 * Create a minimal PluginNodeType for testing
 */
function createNodeType(overrides: Partial<PluginNodeType> = {}): PluginNodeType {
  return {
    type: 'test/node',
    label: 'Test Node',
    icon: 'test',
    color: '#000',
    description: 'Test',
    inputs: [],
    outputs: [],
    configSchema: {},
    defaults: {},
    ...overrides,
  };
}

describe('buildDefaultConfig', () => {
  it('returns empty object when no properties and no defaults', () => {
    const nodeType = createNodeType();
    expect(buildDefaultConfig(nodeType)).toEqual({});
  });

  it('builds config from property defaults', () => {
    const nodeType = createNodeType({
      properties: [
        { name: 'operation', displayName: 'Operation', type: 'options', default: 'sendText' },
        { name: 'text', displayName: 'Text', type: 'string', default: '' },
        { name: 'silent', displayName: 'Silent', type: 'boolean', default: false },
      ],
    });

    const config = buildDefaultConfig(nodeType);

    expect(config).toEqual({
      operation: 'sendText',
      text: '',
      silent: false,
    });
  });

  it('deep clones object defaults (no reference sharing)', () => {
    const nodeType = createNodeType({
      properties: [
        { name: 'source', displayName: 'Source', type: 'fieldSource', default: { type: 'literal', value: '' } },
      ],
    });

    const config1 = buildDefaultConfig(nodeType);
    const config2 = buildDefaultConfig(nodeType);

    // Should be equal but NOT the same reference
    expect(config1.source).toEqual(config2.source);
    expect(config1.source).not.toBe(config2.source);
  });

  it('skips properties without default', () => {
    const nodeType = createNodeType({
      properties: [
        { name: 'withDefault', displayName: 'With', type: 'string', default: 'hello' },
        { name: 'noDefault', displayName: 'No', type: 'string', default: undefined },
      ],
    });

    const config = buildDefaultConfig(nodeType);

    expect(config).toEqual({ withDefault: 'hello' });
    expect('noDefault' in config).toBe(false);
  });

  it('falls back to legacy nodeType.defaults when no properties', () => {
    const nodeType = createNodeType({
      defaults: { operation: 'getInfo', chatId: '' },
    });

    const config = buildDefaultConfig(nodeType);

    expect(config).toEqual({ operation: 'getInfo', chatId: '' });
  });

  it('prefers properties over legacy defaults', () => {
    const nodeType = createNodeType({
      properties: [
        { name: 'operation', displayName: 'Op', type: 'string', default: 'fromProps' },
      ],
      defaults: { operation: 'fromDefaults' },
    });

    const config = buildDefaultConfig(nodeType);

    expect(config.operation).toBe('fromProps');
  });
});
