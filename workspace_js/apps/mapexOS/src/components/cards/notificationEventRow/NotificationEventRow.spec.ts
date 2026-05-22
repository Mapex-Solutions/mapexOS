import { describe, it, expect } from 'vitest';
import { mountWithPlugins } from '@src/test/helpers';
import NotificationEventRow from './NotificationEventRow.vue';
import type { RawNotificationProps } from './interfaces';

const baseEvent: RawNotificationProps = {
  id: 'notif-1',
  notificationType: 'slack',
  notificationName: 'Alert Channel',
  status: 'success',
  tenantId: 'tenant-1',
  created: '2025-01-15T12:00:00Z',
};

function factory(overrides: Record<string, unknown> = {}) {
  return mountWithPlugins(NotificationEventRow, {
    props: { event: baseEvent, ...overrides },
  });
}

describe('NotificationEventRow', () => {
  it('renders without errors', () => {
    const wrapper = factory();
    expect(wrapper.exists()).toBe(true);
  });

  it('getCardClass returns class based on status', () => {
    const wrapper = factory();
    expect(wrapper.vm.getCardClass()).toBe('event-card--success');
  });

  it('getBorderClass returns border class based on status', () => {
    const wrapper = factory();
    expect(wrapper.vm.getBorderClass()).toBe('event-card__border--success');
  });

  it('getNotificationTypeColor returns correct color for slack', () => {
    const wrapper = factory();
    expect(wrapper.vm.getNotificationTypeColor('slack')).toBe('purple-6');
  });

  it('getNotificationTypeColor returns correct color for teams', () => {
    const wrapper = factory();
    expect(wrapper.vm.getNotificationTypeColor('teams')).toBe('blue-6');
  });

  it('getNotificationTypeColor returns grey-6 for unknown type', () => {
    const wrapper = factory();
    expect(wrapper.vm.getNotificationTypeColor('unknown')).toBe('grey-6');
  });

  it('getNotificationTypeIcon returns correct icon for email', () => {
    const wrapper = factory();
    expect(wrapper.vm.getNotificationTypeIcon('email')).toBe('mdi-email');
  });

  it('getNotificationTypeIcon returns fallback for unknown type', () => {
    const wrapper = factory();
    expect(wrapper.vm.getNotificationTypeIcon('unknown')).toBe('mdi-bell-outline');
  });

  it('formatDate returns a formatted date string', () => {
    const wrapper = factory();
    const formatted = wrapper.vm.formatDate('2025-01-15T12:00:00Z');
    expect(typeof formatted).toBe('string');
    expect(formatted.length).toBeGreaterThan(0);
  });
});
