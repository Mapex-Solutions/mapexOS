import { describe, it, expect, vi } from 'vitest';
import { mountWithPlugins } from '@src/test/helpers';
import type { RawTriggerProps } from './interfaces';

// Mock @components/drawers barrel to avoid transitive monaco-editor resolution
vi.mock('@components/drawers', () => ({
  JsonDrawer: { name: 'JsonDrawer', template: '<div />' },
}));

import TriggerEventRow from './TriggerEventRow.vue';

const baseEvent: RawTriggerProps = {
  id: '1',
  triggerType: 'HTTP',
  triggerName: 'My Trigger',
  status: 'success',
  created: '2025-01-15T10:30:00Z',
};

function factory(eventOverrides: Partial<RawTriggerProps> = {}) {
  return mountWithPlugins(TriggerEventRow, {
    props: {
      event: { ...baseEvent, ...eventOverrides },
    },
  });
}

describe('TriggerEventRow', () => {
  it('renders without errors', () => {
    const wrapper = factory();
    expect(wrapper.exists()).toBe(true);
  });

  it('jsonDrawerOpen starts as false', () => {
    const wrapper = factory();
    expect(wrapper.vm.jsonDrawerOpen).toBe(false);
  });

  it('getCardClass returns class based on status', () => {
    const wrapper = factory({ status: 'success' });
    expect(wrapper.vm.getCardClass()).toBe('event-card--success');
  });

  it('getCardClass returns failed class', () => {
    const wrapper = factory({ status: 'failed' });
    expect(wrapper.vm.getCardClass()).toBe('event-card--failed');
  });

  it('getBorderClass returns class based on status', () => {
    const wrapper = factory({ status: 'success' });
    expect(wrapper.vm.getBorderClass()).toBe('event-card__border--success');
  });

  it('getTriggerTypeColor returns correct color for HTTP', () => {
    const wrapper = factory();
    expect(wrapper.vm.getTriggerTypeColor('HTTP')).toBe('blue-6');
  });

  it('getTriggerTypeColor returns correct color for MQTT', () => {
    const wrapper = factory();
    expect(wrapper.vm.getTriggerTypeColor('MQTT')).toBe('purple-6');
  });

  it('getTriggerTypeColor returns grey-6 for unknown type', () => {
    const wrapper = factory();
    expect(wrapper.vm.getTriggerTypeColor('Unknown')).toBe('grey-6');
  });

  it('getTriggerTypeIcon returns correct icon for HTTP', () => {
    const wrapper = factory();
    expect(wrapper.vm.getTriggerTypeIcon('HTTP')).toBe('language');
  });

  it('getTriggerTypeIcon returns correct icon for Notification', () => {
    const wrapper = factory();
    expect(wrapper.vm.getTriggerTypeIcon('Notification')).toBe('notifications');
  });

  it('getTriggerTypeIcon returns help for unknown type', () => {
    const wrapper = factory();
    expect(wrapper.vm.getTriggerTypeIcon('Unknown')).toBe('help');
  });

  it('formatDate returns formatted date string', () => {
    const wrapper = factory();
    const formatted = wrapper.vm.formatDate('2025-01-15T10:30:00Z');
    // Verify it contains expected parts (locale-dependent exact format)
    expect(formatted).toContain('2025');
    expect(formatted).toContain('Jan');
  });
});
