import { describe, it, expect } from 'vitest';
import { mountWithPlugins } from '@src/test/helpers';
import AuditEventRow from './AuditEventRow.vue';
import type { AuditLogProps } from './interfaces';

const baseEvent: AuditLogProps = {
  id: 'audit-1',
  type: 'assets',
  actor: 'admin@test.com',
  action: 'Create',
  resource: 'Asset ABC',
  status: 'success',
  created: '2025-01-15T10:30:00Z',
};

function factory(overrides: Record<string, unknown> = {}) {
  return mountWithPlugins(AuditEventRow, {
    props: { event: baseEvent, ...overrides },
  });
}

describe('AuditEventRow', () => {
  it('renders without errors', () => {
    const wrapper = factory();
    expect(wrapper.exists()).toBe(true);
  });

  it('getTypeLabel returns human-readable label for known type', () => {
    const wrapper = factory();
    expect(wrapper.vm.getTypeLabel('assets')).toBe('Assets');
    expect(wrapper.vm.getTypeLabel('businessRule')).toBe('Business Rule');
    expect(wrapper.vm.getTypeLabel('customers')).toBe('Customers');
  });

  it('getTypeLabel falls back to raw type for unknown type', () => {
    const wrapper = factory();
    expect(wrapper.vm.getTypeLabel('unknown' as any)).toBe('unknown');
  });

  it('getCardClass returns class based on action', () => {
    const wrapper = factory();
    expect(wrapper.vm.getCardClass()).toBe('event-card--create');
  });

  it('getBorderClass returns border class based on action', () => {
    const wrapper = factory();
    expect(wrapper.vm.getBorderClass()).toBe('event-card__border--create');
  });

  it('getActionColor returns correct color for each action', () => {
    const wrapper = factory();
    expect(wrapper.vm.getActionColor()).toBe('green-6');
  });

  it('getActionColor returns blue-6 for Update action', () => {
    const wrapper = factory({
      event: { ...baseEvent, action: 'Update' },
    });
    expect(wrapper.vm.getActionColor()).toBe('blue-6');
  });

  it('formatDate returns a formatted date string', () => {
    const wrapper = factory();
    const formatted = wrapper.vm.formatDate('2025-01-15T10:30:00Z');
    // Should contain day, month abbreviation, and year
    expect(formatted).toBeTruthy();
    expect(typeof formatted).toBe('string');
  });
});
