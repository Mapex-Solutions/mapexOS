import { describe, it, expect } from 'vitest';
import { cleanQueryParams } from './cleanQueryParams';

describe('cleanQueryParams', () => {
  it('removes undefined values', () => {
    expect(cleanQueryParams({ a: 1, b: undefined })).toEqual({ a: 1 });
  });

  it('removes null values', () => {
    expect(cleanQueryParams({ a: 'x', b: null })).toEqual({ a: 'x' });
  });

  it('removes empty strings', () => {
    expect(cleanQueryParams({ a: 1, b: '' })).toEqual({ a: 1 });
  });

  it('keeps false', () => {
    expect(cleanQueryParams({ enabled: false })).toEqual({ enabled: false });
  });

  it('keeps 0', () => {
    expect(cleanQueryParams({ page: 0 })).toEqual({ page: 0 });
  });

  it('keeps valid values', () => {
    const params = { page: 1, perPage: 15, name: 'test', enabled: true };
    expect(cleanQueryParams(params)).toEqual(params);
  });

  it('returns empty object when all values are empty', () => {
    expect(cleanQueryParams({ a: undefined, b: null, c: '' })).toEqual({});
  });

  it('handles mixed valid and invalid', () => {
    expect(cleanQueryParams({ page: 1, name: undefined, enabled: false, search: '' }))
      .toEqual({ page: 1, enabled: false });
  });
});
