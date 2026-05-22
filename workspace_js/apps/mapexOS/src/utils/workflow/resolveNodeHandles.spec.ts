import { describe, it, expect } from 'vitest';
import { resolveNodeHandles } from './resolveNodeHandles';
import type { PluginNodeType } from '@src/components/workflow/interfaces';

function createNodeType(overrides: Partial<PluginNodeType> = {}): PluginNodeType {
  return {
    type: 'test/node',
    label: 'Test',
    icon: 'test',
    color: '#000',
    description: 'Test',
    inputs: [{ id: 'in', label: 'In', position: 'top' }],
    outputs: [{ id: 'out', label: 'Out', position: 'bottom' }],
    configSchema: {},
    defaults: {},
    ...overrides,
  };
}

describe('resolveNodeHandles', () => {
  it('returns static handles when no resolvers defined', () => {
    const nodeType = createNodeType();
    const result = resolveNodeHandles(nodeType, {});
    expect(result.inputs).toHaveLength(1);
    expect(result.outputs).toHaveLength(1);
    expect(result.inputs[0]!.id).toBe('in');
    expect(result.outputs[0]!.id).toBe('out');
  });

  it('calls resolveInputs when defined', () => {
    const nodeType = createNodeType({
      resolveInputs: () => [
        { id: 'dynamic_in', label: 'Dynamic', position: 'top' },
      ],
    });
    const result = resolveNodeHandles(nodeType, {});
    expect(result.inputs).toHaveLength(1);
    expect(result.inputs[0]!.id).toBe('dynamic_in');
  });

  it('calls resolveOutputs when defined', () => {
    const nodeType = createNodeType({
      resolveOutputs: () => [
        { id: 'success', label: 'OK', position: 'bottom' },
        { id: 'error', label: 'Err', position: 'bottom' },
      ],
    });
    const result = resolveNodeHandles(nodeType, {});
    expect(result.outputs).toHaveLength(2);
  });

  it('applies handle overrides from config', () => {
    const nodeType = createNodeType();
    const config = {
      __handleOverrides: {
        out: { label: 'Custom Label' },
      },
    };
    const result = resolveNodeHandles(nodeType, config);
    expect(result.outputs[0]!.label).toBe('Custom Label');
  });

  it('preserves handle id even with overrides', () => {
    const nodeType = createNodeType();
    const config = {
      __handleOverrides: {
        out: { label: 'Renamed' },
      },
    };
    const result = resolveNodeHandles(nodeType, config);
    expect(result.outputs[0]!.id).toBe('out');
  });
});
