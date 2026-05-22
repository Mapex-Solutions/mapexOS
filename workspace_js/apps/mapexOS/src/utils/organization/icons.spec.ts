import { describe, it, expect } from 'vitest';
import { getOrganizationIcon, getOrganizationColor } from './icons';

describe('getOrganizationIcon', () => {
  it('returns store for vendor', () => {
    expect(getOrganizationIcon('vendor')).toBe('store');
  });

  it('returns business for customer', () => {
    expect(getOrganizationIcon('customer')).toBe('business');
  });

  it('returns place for site', () => {
    expect(getOrganizationIcon('site')).toBe('place');
  });

  it('returns apartment for building', () => {
    expect(getOrganizationIcon('building')).toBe('apartment');
  });

  it('returns layers for floor', () => {
    expect(getOrganizationIcon('floor')).toBe('layers');
  });

  it('returns room for zone', () => {
    expect(getOrganizationIcon('zone')).toBe('room');
  });

  it('returns help_outline for unknown type', () => {
    expect(getOrganizationIcon('unknown' as any)).toBe('help_outline');
  });
});

describe('getOrganizationColor', () => {
  it('returns purple for vendor', () => {
    expect(getOrganizationColor('vendor')).toBe('purple');
  });

  it('returns primary for customer', () => {
    expect(getOrganizationColor('customer')).toBe('primary');
  });

  it('returns grey for unknown type', () => {
    expect(getOrganizationColor('unknown' as any)).toBe('grey');
  });
});
