import { describe, it, expect } from 'vitest';
import { isFieldSourceEmpty } from './fieldSourceValidation';

describe('isFieldSourceEmpty', () => {
  it('returns true for undefined source', () => {
    expect(isFieldSourceEmpty(undefined)).toBe(true);
  });

  it('returns true for source with no type', () => {
    expect(isFieldSourceEmpty({ type: '' as any, value: '' })).toBe(true);
  });

  it('returns true for source with empty value', () => {
    expect(isFieldSourceEmpty({ type: 'literal', value: '' })).toBe(true);
  });

  it('returns true for source with whitespace-only value', () => {
    expect(isFieldSourceEmpty({ type: 'literal', value: '   ' })).toBe(true);
  });

  it('returns false for literal with value', () => {
    expect(isFieldSourceEmpty({ type: 'literal', value: 'hello' })).toBe(false);
  });

  it('returns false for state with value', () => {
    expect(isFieldSourceEmpty({ type: 'state', value: 'counter' })).toBe(false);
  });

  it('returns false for event with value', () => {
    expect(isFieldSourceEmpty({ type: 'event', value: 'payload.temperature' })).toBe(false);
  });

  describe('nodeOutput type', () => {
    it('returns true when nodeId is missing', () => {
      expect(isFieldSourceEmpty({ type: 'nodeOutput', value: 'output.data' })).toBe(true);
    });

    it('returns true when nodeId is empty', () => {
      expect(isFieldSourceEmpty({ type: 'nodeOutput', value: 'output.data', nodeId: '' })).toBe(true);
    });

    it('returns true when nodeId is whitespace', () => {
      expect(isFieldSourceEmpty({ type: 'nodeOutput', value: 'output.data', nodeId: '  ' })).toBe(true);
    });

    it('returns false when nodeId and value are present', () => {
      expect(isFieldSourceEmpty({ type: 'nodeOutput', value: 'output.data', nodeId: 'node_1' })).toBe(false);
    });
  });

  describe('fetchOptions type', () => {
    it('returns true when value is empty', () => {
      expect(isFieldSourceEmpty({ type: 'fetchOptions', value: '' })).toBe(true);
    });

    it('returns false when value is selected', () => {
      expect(isFieldSourceEmpty({ type: 'fetchOptions', value: '-100123456' })).toBe(false);
    });
  });
});
