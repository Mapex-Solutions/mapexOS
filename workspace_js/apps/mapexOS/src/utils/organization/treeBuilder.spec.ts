import { describe, it, expect } from 'vitest';
import { getDepthFromPathKey, getParentPathKey, isAncestor } from './treeBuilder';

describe('getDepthFromPathKey', () => {
  it('returns 1 for root level', () => {
    expect(getDepthFromPathKey('000001')).toBe(1);
  });

  it('returns 2 for second level', () => {
    expect(getDepthFromPathKey('000001/000002')).toBe(2);
  });

  it('returns 3 for third level', () => {
    expect(getDepthFromPathKey('000001/000002/000003')).toBe(3);
  });
});

describe('getParentPathKey', () => {
  it('returns null for root level', () => {
    expect(getParentPathKey('000001')).toBeNull();
  });

  it('returns parent for second level', () => {
    expect(getParentPathKey('000001/000002')).toBe('000001');
  });

  it('returns parent for third level', () => {
    expect(getParentPathKey('000001/000002/000003')).toBe('000001/000002');
  });
});

describe('isAncestor', () => {
  it('returns true for direct parent', () => {
    expect(isAncestor('000001', '000001/000002')).toBe(true);
  });

  it('returns true for grandparent', () => {
    expect(isAncestor('000001', '000001/000002/000003')).toBe(true);
  });

  it('returns false for same path', () => {
    expect(isAncestor('000001', '000001')).toBe(false);
  });

  it('returns false for non-ancestor', () => {
    expect(isAncestor('000002', '000001/000003')).toBe(false);
  });

  it('returns false for similar prefix but not ancestor', () => {
    expect(isAncestor('000001', '0000012/000003')).toBe(false);
  });
});
