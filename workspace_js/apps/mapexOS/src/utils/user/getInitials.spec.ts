import { describe, it, expect } from 'vitest';
import { getInitials } from './getInitials';

describe('getInitials', () => {
  it('returns initials from first and last name', () => {
    expect(getInitials('Thiago', 'Anselmo')).toBe('TA');
  });

  it('returns first initial only when no last name', () => {
    expect(getInitials('Thiago')).toBe('T');
  });

  it('returns last initial only when no first name', () => {
    expect(getInitials(undefined, 'Anselmo')).toBe('A');
  });

  it('returns uppercase initials', () => {
    expect(getInitials('thiago', 'anselmo')).toBe('TA');
  });

  it('falls back to email first char when no names', () => {
    expect(getInitials(undefined, undefined, 'thiago@mapex.io')).toBe('T');
  });

  it('returns ? when nothing provided', () => {
    expect(getInitials()).toBe('?');
  });

  it('returns ? for empty strings', () => {
    expect(getInitials('', '', '')).toBe('?');
  });
});
